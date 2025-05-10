package server

import (
	"testing"

	"github.com/cnosuke/mcp-greeting/config"
	"github.com/cnosuke/mcp-greeting/greeter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// TestSetupServerComponents - Test server setup logic
func TestSetupServerComponents(t *testing.T) {
	// Set up test logger
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Test configuration
	cfg := &config.Config{}
	cfg.Greeting.DefaultMessage = "Test greeting"

	// Create and test greeter
	greeterInstance, err := greeter.NewGreeter(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, greeterInstance)

	// Test greeting generation functionality
	greeting, err := greeterInstance.GenerateGreeting("Test User")
	assert.NoError(t, err)
	assert.Equal(t, "Test greeting Test User!", greeting)
}
