package api

// FeedDTO represents a feed as returned by the Cerrynt API.
// It is a direct mapping of the JSON response shape and contains no business logic.
type FeedDTO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// ArticleDTO represents an article as returned by the Cerrynt API.
// Published is a raw string; parsing it into a time.Time is the responsibility
// of the mapping layer (Step 3), not this DTO.
type ArticleDTO struct {
	ID        string `json:"id"`
	FeedID    string `json:"feed_id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Published string `json:"published"`
}
