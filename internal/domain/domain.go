package domain

import "time"

// Feed represents an RSS feed subscription.
type Feed struct {
	ID          string
	Title       string
	URL         string
	UnreadCount int
}

// Article represents a single item from an RSS feed.
type Article struct {
	ID        string
	FeedID    string
	Title     string
	URL       string
	Summary   string
	Published time.Time
	IsRead    bool
}
