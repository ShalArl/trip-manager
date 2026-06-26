package tenant

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GitHubProvisioner struct {
	Token  string
	Repo   string // "ShalArl/trip-manager"
	Branch string // "gitops-prod"
}

func NewGitHubProvisioner(token, repo, branch string) *GitHubProvisioner {
	return &GitHubProvisioner{Token: token, Repo: repo, Branch: branch}
}

func (g *GitHubProvisioner) ProvisionEnterpriseTenant(ctx context.Context, tenantSlug, tenantID, dbPassword string) error {
	// Helm values für den Enterprise-Tenant generieren
	valuesContent := fmt.Sprintf(`tenantSlug: %s
tenantId: %s
postgres:
  db: trips
  user: trips_enterprise
  password: %s
persistence:
  enabled: true
  size: 5Gi
  storageClass: standard-rwo
`, tenantSlug, tenantID, dbPassword)

	// Datei in GitHub erstellen
	path := fmt.Sprintf("gitops/enterprise/%s/values.yaml", tenantSlug)
	return g.createOrUpdateFile(ctx, path, valuesContent,
		fmt.Sprintf("feat: provision enterprise tenant %s", tenantSlug))
}

func (g *GitHubProvisioner) DeprovisionEnterpriseTenant(ctx context.Context, tenantSlug string) error {
	// Alle Dateien im enterprise/<slug>/ Verzeichnis löschen
	files, err := g.listFiles(ctx, fmt.Sprintf("gitops/enterprise/%s", tenantSlug))
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := g.deleteFile(ctx, file.Path, file.SHA,
			fmt.Sprintf("feat: deprovision enterprise tenant %s", tenantSlug)); err != nil {
			return err
		}
	}
	return nil
}

type githubFile struct {
	Path string
	SHA  string
}

func (g *GitHubProvisioner) listFiles(ctx context.Context, path string) ([]githubFile, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s", g.Repo, path, g.Branch)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+g.Token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}

	var items []struct {
		Path string `json:"path"`
		SHA  string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	var files []githubFile
	for _, item := range items {
		files = append(files, githubFile{Path: item.Path, SHA: item.SHA})
	}
	return files, nil
}

func (g *GitHubProvisioner) createOrUpdateFile(ctx context.Context, path, content, message string) error {
	// Prüfen ob Datei bereits existiert
	existingSHA := ""
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s", g.Repo, path, g.Branch)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+g.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	if resp, err := http.DefaultClient.Do(req); err == nil && resp.StatusCode == 200 {
		var existing struct {
			SHA string `json:"sha"`
		}
		json.NewDecoder(resp.Body).Decode(&existing)
		existingSHA = existing.SHA
		resp.Body.Close()
	}

	body := map[string]interface{}{
		"message": message,
		"content": base64.StdEncoding.EncodeToString([]byte(content)),
		"branch":  g.Branch,
	}
	if existingSHA != "" {
		body["sha"] = existingSHA
	}

	bodyJSON, _ := json.Marshal(body)
	putURL := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", g.Repo, path)
	putReq, _ := http.NewRequestWithContext(ctx, "PUT", putURL, bytes.NewReader(bodyJSON))
	putReq.Header.Set("Authorization", "Bearer "+g.Token)
	putReq.Header.Set("Accept", "application/vnd.github+json")
	putReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(putReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("github API error: %d", resp.StatusCode)
	}
	return nil
}

func (g *GitHubProvisioner) deleteFile(ctx context.Context, path, sha, message string) error {
	body := map[string]interface{}{
		"message": message,
		"sha":     sha,
		"branch":  g.Branch,
	}
	bodyJSON, _ := json.Marshal(body)
	url := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", g.Repo, path)
	req, _ := http.NewRequestWithContext(ctx, "DELETE", url, bytes.NewReader(bodyJSON))
	req.Header.Set("Authorization", "Bearer "+g.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// generateDBPassword generiert ein zufälliges Passwort
func generateDBPassword() string {
	return fmt.Sprintf("ent-%d", time.Now().UnixNano())
}
