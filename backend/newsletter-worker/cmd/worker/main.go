package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/robfig/cron/v3"
)

type config struct {
	UsersDBURL      string
	NewsletterDBURL string
	Neo4jURI        string
	Neo4jUser       string
	Neo4jPassword   string
	CronSchedule    string
	RunMode         string
}

func loadConfig() config {
	return config{
		UsersDBURL:      getEnv("USERS_DB_URL", "postgres://postgres:postgres@localhost:5432/users_db?sslmode=disable"),
		NewsletterDBURL: getEnv("NEWSLETTER_DB_URL", "postgres://postgres:postgres@localhost:5432/newsletter_db?sslmode=disable"),
		Neo4jURI:        getEnv("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:       getEnv("NEO4J_USERNAME", "neo4j"),
		Neo4jPassword:   getEnv("NEO4J_PASSWORD", "neo4jpassword"),
		CronSchedule:    getEnv("CRON_SCHEDULE", "*/2 * * * *"),
		RunMode:         getEnv("RUN_MODE", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	cfg := loadConfig()

	usersDB, err := sqlx.Connect("postgres", cfg.UsersDBURL)
	if err != nil {
		log.Fatalf("newsletter-worker: failed to connect to users-db: %v", err)
	}
	defer usersDB.Close()
	log.Println("newsletter-worker: connected to users-db")

	newsletterDB, err := sqlx.Connect("postgres", cfg.NewsletterDBURL)
	if err != nil {
		log.Fatalf("newsletter-worker: failed to connect to newsletter-db: %v", err)
	}
	defer newsletterDB.Close()
	log.Println("newsletter-worker: connected to newsletter-db")

	if err := migrate(newsletterDB); err != nil {
		log.Fatalf("newsletter-worker: migration failed: %v", err)
	}

	neo4jDriver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jURI,
		neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""),
	)
	if err != nil {
		log.Fatalf("newsletter-worker: failed to create neo4j driver: %v", err)
	}
	defer neo4jDriver.Close(context.Background())

	if err := neo4jDriver.VerifyConnectivity(context.Background()); err != nil {
		log.Fatalf("newsletter-worker: neo4j not reachable: %v", err)
	}
	log.Println("newsletter-worker: connected to neo4j")

	// Kubernetes: einmal ausführen und beenden
	if cfg.RunMode == "oneshot" {
		log.Println("newsletter-worker: oneshot mode, running once...")
		runGeneration(usersDB, newsletterDB, neo4jDriver)
		log.Println("newsletter-worker: done")
		return
	}

	// Docker Compose: Cron Modus
	log.Println("newsletter-worker: running initial generation...")
	runGeneration(usersDB, newsletterDB, neo4jDriver)

	c := cron.New()
	c.AddFunc(cfg.CronSchedule, func() {
		log.Println("newsletter-worker: cron triggered, generating newsletters...")
		runGeneration(usersDB, newsletterDB, neo4jDriver)
	})
	c.Start()
	log.Printf("newsletter-worker: cron started with schedule '%s'", cfg.CronSchedule)

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
	<-sigch
	log.Println("newsletter-worker: shutting down...")
	c.Stop()
}

func migrate(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS newsletters (
			firebase_uid TEXT PRIMARY KEY,
			content      JSONB NOT NULL,
			generated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func runGeneration(usersDB *sqlx.DB, newsletterDB *sqlx.DB, driver neo4j.DriverWithContext) {
	ctx := context.Background()

	var uids []string
	if err := usersDB.SelectContext(ctx, &uids, `SELECT firebase_uid FROM users`); err != nil {
		log.Printf("newsletter-worker: failed to fetch users: %v", err)
		return
	}
	log.Printf("newsletter-worker: generating for %d users", len(uids))

	for _, uid := range uids {
		sections, err := generateForUser(ctx, driver, uid)
		if err != nil {
			log.Printf("newsletter-worker: error for user %s: %v", uid, err)
			continue
		}

		content, err := json.Marshal(sections)
		if err != nil {
			log.Printf("newsletter-worker: marshal error for user %s: %v", uid, err)
			continue
		}

		_, err = newsletterDB.ExecContext(ctx, `
			INSERT INTO newsletters (firebase_uid, content, generated_at)
			VALUES ($1, $2, NOW())
			ON CONFLICT (firebase_uid) DO UPDATE
			SET content = EXCLUDED.content,
			    generated_at = EXCLUDED.generated_at
		`, uid, content)
		if err != nil {
			log.Printf("newsletter-worker: db write error for user %s: %v", uid, err)
			continue
		}
	}
	log.Println("newsletter-worker: generation complete")
}
