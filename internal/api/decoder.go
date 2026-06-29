package api

import (
	"encoding/json"
	"fmt"
	"io"
)

// DecodeFeeds reads JSON from r and decodes it into a slice of FeedDTO.
func DecodeFeeds(r io.Reader) ([]FeedDTO, error) {
	var feeds []FeedDTO
	if err := json.NewDecoder(r).Decode(&feeds); err != nil {
		return nil, fmt.Errorf("api: decode feeds: %w", err)
	}
	return feeds, nil
}

// DecodeArticles reads JSON from r and decodes it into a slice of ArticleDTO.
func DecodeArticles(r io.Reader) ([]ArticleDTO, error) {
	var articles []ArticleDTO
	if err := json.NewDecoder(r).Decode(&articles); err != nil {
		return nil, fmt.Errorf("api: decode articles: %w", err)
	}
	return articles, nil
}
