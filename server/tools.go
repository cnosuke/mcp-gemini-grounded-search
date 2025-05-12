package server

import (
	"context"

	"github.com/cnosuke/mcp-gemini-grounded-search/searcher"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// RegisterAllTools - Register all tools with the server
func RegisterAllTools(mcpServer *server.MCPServer, searcher *searcher.Searcher) error {
	// Register search tool
	if err := registerSearchTool(mcpServer, searcher); err != nil {
		return err
	}

	return nil
}

// registerSearchTool - Register the search tool
func registerSearchTool(mcpServer *server.MCPServer, searcher *searcher.Searcher) error {
	zap.S().Debugw("registering search tool")

	// Define the tool
	tool := mcp.NewTool("search",
		mcp.WithDescription("Search the web with Gemini grounded search"),
		mcp.WithString("query",
			mcp.Description("The search query"),
			mcp.Required(), // 修正：引数なしで使用
		),
		mcp.WithNumber("max_token",
			mcp.Description("Maximum number of tokens for the response"),
		),
	)

	// Add the tool handler
	mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		response, err := searcher.Search(ctx, query, maxToken)
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
