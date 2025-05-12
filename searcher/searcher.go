package searcher

import (
	"context"
	"encoding/json"

	search "github.com/cnosuke/go-gemini-grounded-search"
	"github.com/cnosuke/mcp-gemini-grounded-search/config"
	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

// Searcher - Search interface
type Searcher struct {
	client *search.Client
	cfg    *config.Config
}

// SearchResponse - Response for search results
type SearchResponse struct {
	Text       string       `json:"text"`
	Groundings []*Grounding `json:"groundings"`
}

// Grounding - Information about the source of the search content
type Grounding struct {
	Title  string `json:"title"`
	Domain string `json:"domain"`
	URL    string `json:"url"`
}

// NewSearcher - Create a new Searcher
func NewSearcher(ctx context.Context, cfg *config.Config) (*Searcher, error) {
	zap.S().Infow("creating new Searcher",
		"model_name", cfg.Gemini.ModelName)

	// Initialize the client
	client, err := search.NewClient(ctx, cfg.Gemini.APIKey,
		search.WithModelName(cfg.Gemini.ModelName),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Gemini client")
	}

	return &Searcher{
		client: client,
		cfg:    cfg,
	}, nil
}

// Search - Perform a search with the given query and max token limit
func (s *Searcher) Search(ctx context.Context, query string, maxTokens int) (*SearchResponse, error) {
	zap.S().Debugw("executing search",
		"query", query,
		"max_tokens", maxTokens)

	// Set parameters for the search
	params := &search.GenerationParams{
		Prompt: query,
	}

	// Apply max tokens if specified
	if maxTokens > 0 {
		maxTokensInt32 := int32(maxTokens)
		params.MaxOutputTokens = &maxTokensInt32
	}

	// Execute the search
	result, err := s.client.GenerateGroundedContentWithParams(ctx, params)
	if err != nil {
		if apiErr, ok := search.GetAPIError(err); ok {
			zap.S().Errorw("API error in search",
				"status_code", apiErr.StatusCode,
				"message", apiErr.Message)
		} else if search.IsContentBlockedError(err) {
			zap.S().Errorw("content blocked error in search",
				"error", err)
		}
		return nil, errors.Wrap(err, "failed to generate grounded content")
	}

	// Create response
	response := &SearchResponse{
		Text:       result.GeneratedText,
		Groundings: make([]*Grounding, 0, len(result.GroundingAttributions)),
	}

	// Add groundings
	for _, attr := range result.GroundingAttributions {
		response.Groundings = append(response.Groundings, &Grounding{
			Title:  attr.Title,
			Domain: attr.Domain,
			URL:    attr.URL,
		})
	}

	return response, nil
}

// ToJSON - Convert search response to JSON string
func (r *SearchResponse) ToJSON() (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal response to JSON")
	}
	return string(bytes), nil
}
