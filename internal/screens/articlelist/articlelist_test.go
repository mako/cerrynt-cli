package articlelist

import (
	"testing"
	"time"

	"github.com/mako/cerrynt-cli/internal/domain"
)

// ---------------------------------------------------------------------------
// formatAge
// ---------------------------------------------------------------------------

func TestFormatAge(t *testing.T) {
	t.Parallel()

	// Each case specifies how far in the past the timestamp is.
	// Values are chosen well away from bucket boundaries so the tiny
	// elapsed time between time.Now() and the formatAge call cannot
	// push the result into a different bucket.
	tests := []struct {
		name string
		age  time.Duration // how long ago the article was published
		want string
	}{
		{"just now – 0 seconds", 0, "just now"},
		{"just now – 30 seconds", 30 * time.Second, "just now"},
		{"just now – 59 seconds", 59 * time.Second, "just now"},
		{"minutes – 1 minute", 90 * time.Second, "1m ago"},
		{"minutes – 5 minutes", 5 * time.Minute, "5m ago"},
		{"minutes – 59 minutes", 59 * time.Minute, "59m ago"},
		{"hours – 1 hour", 90 * time.Minute, "1h ago"},
		{"hours – 3 hours", 3 * time.Hour, "3h ago"},
		{"hours – 23 hours", 23 * time.Hour, "23h ago"},
		{"days – 1 day", 36 * time.Hour, "1d ago"},
		{"days – 2 days", 48 * time.Hour, "2d ago"},
		{"days – 7 days", 7 * 24 * time.Hour, "7d ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			published := time.Now().Add(-tt.age)
			got := formatAge(published)
			if got != tt.want {
				t.Errorf("formatAge(%v ago) = %q, want %q", tt.age, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// MarkRead
// ---------------------------------------------------------------------------

func TestMarkRead(t *testing.T) {
	t.Parallel()

	articles := []domain.Article{
		{ID: "1", Title: "First", IsRead: false},
		{ID: "2", Title: "Second", IsRead: false},
		{ID: "3", Title: "Third", IsRead: true},
	}
	original := New(domain.Feed{ID: "f1"}).SetArticles(articles)

	t.Run("marks the correct article as read", func(t *testing.T) {
		t.Parallel()
		updated := original.MarkRead("2")

		if !updated.articles[1].IsRead {
			t.Error("article 2 should be marked read after MarkRead")
		}
		if updated.articles[0].IsRead {
			t.Error("article 1 should not be affected")
		}
		if !updated.articles[2].IsRead {
			t.Error("article 3 should still be read (was already read)")
		}
	})

	t.Run("returns new model – original is not mutated", func(t *testing.T) {
		t.Parallel()
		_ = original.MarkRead("1")

		if original.articles[0].IsRead {
			t.Error("MarkRead must not mutate the receiver; original article 1 should still be unread")
		}
	})

	t.Run("non-existent id is a no-op", func(t *testing.T) {
		t.Parallel()
		updated := original.MarkRead("does-not-exist")

		for _, a := range updated.articles {
			if a.ID != "3" && a.IsRead {
				t.Errorf("article %q should not have been marked read", a.ID)
			}
		}
	})

	t.Run("already-read article stays read", func(t *testing.T) {
		t.Parallel()
		updated := original.MarkRead("3")

		if !updated.articles[2].IsRead {
			t.Error("article 3 should still be read after MarkRead on an already-read article")
		}
	})

	t.Run("empty article list is safe", func(t *testing.T) {
		t.Parallel()
		empty := New(domain.Feed{ID: "f2"}).SetArticles(nil)
		updated := empty.MarkRead("anything")

		if len(updated.articles) != 0 {
			t.Errorf("expected 0 articles, got %d", len(updated.articles))
		}
	})
}
