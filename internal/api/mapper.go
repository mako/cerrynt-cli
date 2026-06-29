package api

import (
	"time"

	"github.com/mako/cerrynt-cli/internal/domain"
)

// MapFeeds transforms a slice of FeedDTOs into domain.Feeds.
// It performs a strict 1:1 field mapping with no business logic.
func MapFeeds(input []FeedDTO) []domain.Feed {
	if input == nil {
		return nil
	}

	result := make([]domain.Feed, len(input))
	for i, dto := range input {
		result[i] = domain.Feed{
			ID:    dto.ID,
			Title: dto.Title,
			URL:   dto.URL,
		}
	}

	return result
}

// MapArticles transforms a slice of ArticleDTOs into domain.Articles.
// It maps matching fields 1:1. The Published string is parsed as RFC3339
// because the domain model requires a time.Time; if parsing fails, it
// falls back to the zero value.
func MapArticles(input []ArticleDTO) []domain.Article {
	if input == nil {
		return nil
	}

	result := make([]domain.Article, len(input))
	for i, dto := range input {
		published, err := time.Parse(time.RFC3339, dto.Published)
		if err != nil {
			published = time.Now().UTC()
		}

		result[i] = domain.Article{
			ID:        dto.ID,
			FeedID:    dto.FeedID,
			Title:     dto.Title,
			URL:       dto.URL,
			Published: published,
			// Summary and IsRead are left as zero values since they
			// are not present in the DTO layer.
		}
	}

	return result
}
