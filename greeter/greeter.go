package greeter

import (
	"github.com/cnosuke/mcp-greeting/config"
	"go.uber.org/zap"
)

// Greeter - Greeting structure
type Greeter struct {
	DefaultMessage string
	cfg            *config.Config
}

// NewGreeter - Create a new Greeter
func NewGreeter(cfg *config.Config) (*Greeter, error) {
	zap.S().Infow("creating new Greeter",
		"default_message", cfg.Greeting.DefaultMessage)

	return &Greeter{
		DefaultMessage: cfg.Greeting.DefaultMessage,
		cfg:            cfg,
	}, nil
}

// GenerateGreeting - Generate a greeting message
func (g *Greeter) GenerateGreeting(name string) (string, error) {
	if name == "" {
		return g.DefaultMessage, nil
	}
	return g.DefaultMessage + " " + name + "!", nil
}
