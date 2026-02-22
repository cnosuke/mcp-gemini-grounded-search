package server

import (
	"context"

	"github.com/cnosuke/mcp-gemini-grounded-search/config"
	ierrors "github.com/cnosuke/mcp-gemini-grounded-search/internal/errors"
	"github.com/cnosuke/mcp-gemini-grounded-search/searcher"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// createMCPServer creates and configures an MCP server instance with all tools registered.
func createMCPServer(cfg *config.Config, name, version, revision string) (*mcpserver.MCPServer, *searcher.Searcher, error) {
	versionString := version
	if revision != "" && revision != "xxx" {
		versionString = versionString + " (" + revision + ")"
	}

	zap.S().Debugw("creating Searcher")
	ctx := context.Background()
	searcherInstance, err := searcher.NewSearcher(ctx, cfg)
	if err != nil {
		zap.S().Errorw("failed to create Searcher", "error", err)
		return nil, nil, err
	}

	hooks := &mcpserver.Hooks{}
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		zap.S().Errorw("MCP error occurred",
			"id", id,
			"method", method,
			"error", err,
		)
	})

	zap.S().Debugw("creating MCP server", "name", name, "version", versionString)
	s := mcpserver.NewMCPServer(name, versionString, mcpserver.WithHooks(hooks))

	zap.S().Debugw("registering tools")
	if err := RegisterAllTools(s, searcherInstance); err != nil {
		zap.S().Errorw("failed to register tools", "error", err)
		return nil, nil, err
	}

	return s, searcherInstance, nil
}

// RunStdio starts the MCP server with stdio transport.
func RunStdio(cfg *config.Config, name string, version string, revision string) error {
	zap.S().Infow("starting MCP Gemini Grounded Search Server (stdio)")

	s, _, err := createMCPServer(cfg, name, version, revision)
	if err != nil {
		return err
	}

	zap.S().Infow("starting MCP server")
	if err := mcpserver.ServeStdio(s); err != nil {
		zap.S().Errorw("failed to start server", "error", err)
		return ierrors.Wrap(err, "failed to start server")
	}

	zap.S().Infow("server shutting down")
	return nil
}
