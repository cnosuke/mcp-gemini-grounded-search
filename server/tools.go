package server

import (
	"context"

	"github.com/cnosuke/mcp-greeting/greeter"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// GreetingHelloArgs - Arguments for greeting_hello tool (kept for testing compatibility)
type GreetingHelloArgs struct {
	Name string `json:"name" jsonschema:"description=Optional name for personalized greeting"`
}

// RegisterAllTools - Register all tools with the server
func RegisterAllTools(mcpServer *server.MCPServer, greeter *greeter.Greeter) error {
	// Register greeting_hello tool
	if err := registerGreetingHelloTool(mcpServer, greeter); err != nil {
		return err
	}

	return nil
}

// registerGreetingHelloTool - Register the greeting_hello tool
func registerGreetingHelloTool(mcpServer *server.MCPServer, greeter *greeter.Greeter) error {
	zap.S().Debugw("registering greeting_hello tool")

	// Define the tool
	tool := mcp.NewTool("greeting_hello",
		mcp.WithDescription("Generate a greeting message"),
		mcp.WithString("name",
			mcp.Description("Optional name for personalized greeting"),
		),
	)

	// Add the tool handler
	mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract name parameter
		var name string
		if nameVal, ok := request.Params.Arguments["name"].(string); ok {
			name = nameVal
		}

		zap.S().Debugw("executing greeting_hello",
			"name", name)

		// Generate greeting
		greeting, err := greeter.GenerateGreeting(name)
		if err != nil {
			zap.S().Errorw("failed to generate greeting",
				"name", name,
				"error", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(greeting), nil
	})

	return nil
}
