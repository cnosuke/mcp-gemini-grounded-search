package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cnosuke/mcp-gemini-grounded-search/config"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// RunHTTP starts the MCP server with Streamable HTTP transport.
func RunHTTP(cfg *config.Config, name string, version string, revision string) error {
	zap.S().Infow("starting MCP Gemini Grounded Search Server (HTTP)")

	s, _, err := createMCPServer(cfg, name, version, revision)
	if err != nil {
		return err
	}

	opts := []mcpserver.StreamableHTTPOption{
		mcpserver.WithEndpointPath(cfg.HTTP.EndpointPath),
	}
	if cfg.HTTP.HeartbeatSeconds > 0 {
		opts = append(opts, mcpserver.WithHeartbeatInterval(time.Duration(cfg.HTTP.HeartbeatSeconds)*time.Second))
	}

	httpServer := mcpserver.NewStreamableHTTPServer(s, opts...)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.Handle(cfg.HTTP.EndpointPath, httpServer)

	var handler http.Handler = mux
	handler = withOriginValidation(handler, cfg.HTTP.AllowedOrigins)
	handler = withAuthMiddleware(handler, cfg.HTTP.AuthToken)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: handler,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		zap.S().Infow("HTTP server listening", "addr", srv.Addr, "endpoint", cfg.HTTP.EndpointPath)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case sig := <-quit:
		zap.S().Infow("received signal, shutting down", "signal", sig)
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		zap.S().Errorw("MCP server shutdown error", "error", err)
	}
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Errorw("HTTP server shutdown error", "error", err)
		return err
	}

	zap.S().Infow("server shut down gracefully")
	return nil
}
