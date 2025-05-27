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
		mcp.WithString("question",
			mcp.Description("The question to be examined. Formulate the question as a complete sentence in natural language. Questions should not be a list of space-separated keywords. Example: [What are the most contributive biological factors to human civilizational evolution, according to the latest research?]"),
			mcp.Required(),
		),
		mcp.WithNumber("max_token",
			mcp.Description(fmt.Sprintf("Maximum number of tokens for the response (default: %d)", s.DefaultMaxTokens)),
		),
	)

	// Add the tool handler
	m.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments
		args, argsOK := request.Params.Arguments.(map[string]any)
		if !argsOK {
			return mcp.NewToolResultError("Invalid arguments format"), nil
		}

		// Extract question parameter
		questionVal, questionValOK := args["question"]
		if !questionValOK {
			return mcp.NewToolResultError("Missing question parameter"), nil
		}
		question, questionStrOK := questionVal.(string)
		if !questionStrOK || question == "" {
			return mcp.NewToolResultError("Missing or empty question parameter"), nil
		}

		// Extract max_token parameter (optional)
		var maxToken int
		if maxTokenVal, maxTokenValOK := args["max_token"]; maxTokenValOK {
			if mt, maxTokenFloatOK := maxTokenVal.(float64); maxTokenFloatOK {
				maxToken = int(mt)
			}
		}

		zap.S().Debugw("executing search",
			"question", question,
			"max_token", maxToken)

		// Perform search
		response, err := s.Search(ctx, question, maxToken)
		if err != nil {
			zap.S().Errorw("failed to search",
				"question", question,
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
