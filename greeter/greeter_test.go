package greeter

import (
	"testing"

	"github.com/cnosuke/mcp-greeting/config"
	"github.com/stretchr/testify/assert"
)

func TestGreeter_GenerateGreeting(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		defaultMessage string
		inputName      string
		expected       string
	}{
		{
			name:           "Default message only",
			defaultMessage: "Hello!",
			inputName:      "",
			expected:       "Hello!",
		},
		{
			name:           "Greeting with name",
			defaultMessage: "Hello!",
			inputName:      "Tanaka",
			expected:       "Hello! Tanaka!",
		},
		{
			name:           "Different default message",
			defaultMessage: "Hi!",
			inputName:      "Smith",
			expected:       "Hi! Smith!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test configuration
			cfg := &config.Config{}
			cfg.Greeting.DefaultMessage = tc.defaultMessage

			// Initialize Greeter
			greeter, err := NewGreeter(cfg)
			assert.NoError(t, err)

			// Generate greeting
			greeting, err := greeter.GenerateGreeting(tc.inputName)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, greeting)
		})
	}
}
