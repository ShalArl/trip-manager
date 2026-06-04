package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tenantdb"

	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/robfig/cron/v3"
)

type config struct {
	UsersDBURL      string `envconfig:"USERS_DB_URL"`
	NewsletterDBURL string `envconfig:"NEWSLETTER_DB_URL"`
	Neo4jURI        string `envconfig:"NEO4J_URI"`
	Neo4jUser       string `envconfig:"NEO4J_USERNAME"`
	Neo4jPassword   string `envconfig:"NEO4J_PASSWORD"`
	CronSchedule    string `envconfig:"CRON_SCHEDULE"`
	RunMode         string `envconfig:"RUN_MODE"`
	LogLevel        string `envconfig:"LOG_LEVEL"`
}

func load() (*config, error) {
	var config config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	cfg, err := load()
	if err != nil {
		log.Printf("Failed to load config")

	}

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
	_, err = c.AddFunc(cfg.CronSchedule, func() {
		log.Println("newsletter-worker: cron triggered, generating newsletters...")
		runGeneration(usersDB, newsletterDB, neo4jDriver)
	})
	if err != nil {
		log.Printf("newsletter-worker: failed to add cron job: %v", err)
		return
	}
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
			firebase_uid TEXT         NOT NULL,
			tenant_id    VARCHAR(255) NOT NULL DEFAULT 'default',
			content      JSONB        NOT NULL,
			generated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
			PRIMARY KEY (firebase_uid, tenant_id)
		);
		CREATE INDEX IF NOT EXISTS idx_newsletters_tenant_id ON newsletters(tenant_id);
		ALTER TABLE newsletters ENABLE ROW LEVEL SECURITY;
		DO $$ BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_policies 
				WHERE tablename = 'newsletters' AND policyname = 'tenant_isolation_newsletters'
			) THEN
				CREATE POLICY tenant_isolation_newsletters ON newsletters
					USING (tenant_id = current_setting('app.tenant_id', true));
			END IF;
		END $$;
	`)
	return err
}

func runGeneration(usersDB *sqlx.DB, newsletterDB *sqlx.DB, driver neo4j.DriverWithContext) {
	ctx := context.Background()

	var users []struct {
		FirebaseUID string `db:"firebase_uid"`
		TenantID    string `db:"tenant_id"`
	}
	if err := usersDB.SelectContext(ctx, &users, `SELECT firebase_uid, tenant_id FROM users`); err != nil {
		log.Printf("newsletter-worker: failed to fetch users: %v", err)
		return
	}
	log.Printf("newsletter-worker: generating for %d users", len(users))

	for _, u := range users {
		tenantCtx := tenantdb.WithTenantID(ctx, u.TenantID)

		sections, err := generateForUser(tenantCtx, driver, u.FirebaseUID)
		if err != nil {
			log.Printf("newsletter-worker: error for user %s: %v", u.FirebaseUID, err)
			continue
		}

		content, err := json.Marshal(sections)
		if err != nil {
			log.Printf("newsletter-worker: marshal error for user %s: %v", u.FirebaseUID, err)
			continue
		}

		err = tenantdb.WithTenant(tenantCtx, newsletterDB, func(tx *sqlx.Tx) error {
			_, err := tx.ExecContext(tenantCtx, `
				INSERT INTO newsletters (firebase_uid, tenant_id, content, generated_at)
				VALUES ($1, $2, $3, NOW())
				ON CONFLICT (firebase_uid, tenant_id) DO UPDATE
				SET content = EXCLUDED.content,
				    generated_at = EXCLUDED.generated_at
			`, u.FirebaseUID, u.TenantID, content)
			return err
		})
		if err != nil {
			log.Printf("newsletter-worker: db write error for user %s: %v", u.FirebaseUID, err)
			continue
		}
	}
	log.Println("newsletter-worker: generation complete")
}
