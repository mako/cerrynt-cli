package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/mako/cerrynt-cli/internal/config"
)

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected os.ErrNotExist to be wrapped in error, got: %v", err)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	content := `
api_base_url: https://api.cerrynt.example.com
auth_token: tok_test123
feeds:
  - id: "1"
    title: Hacker News
    url: https://news.ycombinator.com/rss
  - id: "2"
    title: Go Blog
    url: https://go.dev/blog/feed.atom
`
	path := writeTempConfig(t, content)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.APIBaseURL != "https://api.cerrynt.example.com" {
		t.Errorf("APIBaseURL = %q, want %q", cfg.APIBaseURL, "https://api.cerrynt.example.com")
	}
	if cfg.AuthToken != "tok_test123" {
		t.Errorf("AuthToken = %q, want %q", cfg.AuthToken, "tok_test123")
	}
	if len(cfg.Feeds) != 2 {
		t.Fatalf("len(Feeds) = %d, want 2", len(cfg.Feeds))
	}
	if cfg.Feeds[0].ID != "1" || cfg.Feeds[0].Title != "Hacker News" {
		t.Errorf("Feeds[0] = %+v, unexpected", cfg.Feeds[0])
	}
}

func TestLoad_UnknownField(t *testing.T) {
	content := `
feeds:
  - id: "1"
    title: Test
    url: https://example.com
    unknown_field: oops
`
	path := writeTempConfig(t, content)

	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for unknown field, got nil")
	}
}

func TestLoad_EmptyFeeds(t *testing.T) {
	content := `
api_base_url: https://api.cerrynt.example.com
feeds: []
`
	path := writeTempConfig(t, content)

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Feeds) != 0 {
		t.Errorf("len(Feeds) = %d, want 0", len(cfg.Feeds))
	}
}

// writeTempConfig writes content to a temporary YAML file and returns its path.
func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}
