package server

import (
	"testing"

	"github.com/cnosuke/mcp-greeting/greeter"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Test for GreetingHelloArgs
func TestGreetingHelloArgs(t *testing.T) {
	// When name is empty
	argsEmpty := GreetingHelloArgs{
		Name: "",
	}
	assert.Equal(t, "", argsEmpty.Name)

	// When name is set
	argsWithName := GreetingHelloArgs{
		Name: "Test User",
	}
	assert.Equal(t, "Test User", argsWithName.Name)
}

// TestRegisterAllTools - Test tool registration
func TestRegisterAllTools(t *testing.T) {
	// Set up test logger
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Create MCP server
	mcpServer := server.NewMCPServer("test-server", "0.0.1")

	// Create mock greeter instance
	cfg := &greeter.Greeter{
		DefaultMessage: "Hello!",
	}

	// Register tools
	err := RegisterAllTools(mcpServer, cfg)
	assert.NoError(t, err)

	// Verify that tools are registered
	// This is a basic test - we're just making sure registration doesn't fail
	// More comprehensive tests would verify the actual tool behavior
}
