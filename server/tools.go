package server

import (
	"context"
	"fmt"

	"github.com/cnosuke/mcp-gemini-grounded-search/searcher"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// RegisterAllTools - Register all tools with the server
func RegisterAllTools(m *server.MCPServer, s *searcher.Searcher) error {
	// Register search tool
	if err := registerSearchTool(m, s); err != nil {
		return err
	}

	return nil
}

// registerSearchTool - Register the search tool
func registerSearchTool(m *server.MCPServer, s *searcher.Searcher) error {
	zap.S().Debugw("registering search tool")

	// Define the tool
	tool := mcp.NewTool("search",
		mcp.WithDescription("Searches the web using Gemini Grounded Search. Expect more accurate results by searching in a natural language question format rather than by keywords."),
		mcp.WithString("query",
			mcp.Description("The search query. Please describe it as if asking a question in natural language, rather than specifying keywords. Example: [What are the most contributive biological factors to human civilizational evolution, according to the latest research?]"),
			mcp.Required(),
		),
		mcp.WithNumber("max_token",
			mcp.Description(fmt.Sprintf("Maximum number of tokens for the response (default: %d)", s.DefaultMaxTokens)),
		),
	)

	// Add the tool handler
	m.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract query parameter
		query, ok := request.Params.Arguments["query"].(string)
		if !ok || query == "" {
			return mcp.NewToolResultError("Missing or empty query parameter"), nil
		}

		// Extract max_token parameter (optional)
		var maxToken int
		if maxTokenVal, ok := request.Params.Arguments["max_token"].(float64); ok {
			maxToken = int(maxTokenVal)
		}

		zap.S().Debugw("executing search",
			"query", query,
			"max_token", maxToken)

		// Perform search
		response, err := s.Search(ctx, query, maxToken)
		if err != nil {
			zap.S().Errorw("failed to search",
				"query", query,
				"error", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert response to JSON
		jsonResponse, err := response.ToJSON()
		if err != nil {
			zap.S().Errorw("failed to convert response to JSON",
				"error", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(jsonResponse), nil
	})

	return nil
}
