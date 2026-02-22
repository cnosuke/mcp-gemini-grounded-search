package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cnosuke/mcp-gemini-grounded-search/config"
	ierrors "github.com/cnosuke/mcp-gemini-grounded-search/internal/errors"
	"github.com/cnosuke/mcp-gemini-grounded-search/logger"
	"github.com/cnosuke/mcp-gemini-grounded-search/server"
	"github.com/urfave/cli/v3"
)

var (
	// Version and Revision are replaced when building.
	// To set specific version, edit Makefile.
	Version  = "0.0.1"
	Revision = "xxx"

	Name  = "mcp-gemini-grounded-search"
	Usage = "MCP server for Gemini grounded search"
)

func main() {
	app := &cli.Command{
		Name:    Name,
		Usage:   Usage,
		Version: fmt.Sprintf("%s (%s)", Version, Revision),
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "Start the MCP server for Gemini grounded search (stdio)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "config.yml",
						Usage:   "path to the configuration file",
					},
					&cli.StringFlag{
						Name:    "log",
						Aliases: []string{"l"},
						Usage:   "path to log file (overrides config file)",
						Sources: cli.EnvVars("LOG_PATH"),
					},
					&cli.BoolFlag{
						Name:    "debug",
						Aliases: []string{"d"},
						Usage:   "enable debug logging (overrides config file)",
						Sources: cli.EnvVars("DEBUG"),
					},
					&cli.StringFlag{
						Name:    "api-key",
						Aliases: []string{"k"},
						Usage:   "Gemini API key (overrides config file)",
						Sources: cli.EnvVars("GEMINI_API_KEY"),
					},
					&cli.StringFlag{
						Name:    "model",
						Aliases: []string{"m"},
						Usage:   "Gemini model name (overrides config file)",
						Sources: cli.EnvVars("GEMINI_MODEL_NAME"),
					},
					&cli.StringFlag{
						Name:    "thinking-level",
						Usage:   "Gemini thinking level: MINIMAL, LOW, MEDIUM, HIGH (overrides config file)",
						Sources: cli.EnvVars("GEMINI_THINKING_LEVEL"),
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					cfg, err := config.LoadConfig(cmd.String("config"))
					if err != nil {
						return ierrors.Wrap(err, "failed to load configuration file")
					}
					if cmd.IsSet("log") {
						cfg.Log = cmd.String("log")
					}
					if cmd.IsSet("debug") {
						cfg.Debug = cmd.Bool("debug")
					}
					if cmd.IsSet("api-key") {
						cfg.Gemini.APIKey = cmd.String("api-key")
					}
					if cmd.IsSet("model") {
						cfg.Gemini.ModelName = cmd.String("model")
					}
					if cmd.IsSet("thinking-level") {
						cfg.Gemini.ThinkingLevel = cmd.String("thinking-level")
					}
					if cfg.Gemini.APIKey == "" {
						return fmt.Errorf("Gemini API key is required. Set it in config.yml or use --api-key flag or GEMINI_API_KEY environment variable")
					}
					if err := logger.InitLogger(cfg.Debug, cfg.Log); err != nil {
						return ierrors.Wrap(err, "failed to initialize logger")
					}
					defer logger.Sync()
					return server.RunStdio(cfg, Name, Version, Revision)
				},
			},
			{
				Name:  "httpserver",
				Usage: "Start the MCP server for Gemini grounded search (Streamable HTTP)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "config.yml",
						Usage:   "path to the configuration file",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					cfg, err := config.LoadConfig(cmd.String("config"))
					if err != nil {
						return ierrors.Wrap(err, "failed to load configuration file")
					}
					if cfg.Gemini.APIKey == "" {
						return fmt.Errorf("Gemini API key is required. Set it in config.yml or GEMINI_API_KEY environment variable")
					}
					if err := logger.InitLogger(cfg.Debug, cfg.Log); err != nil {
						return ierrors.Wrap(err, "failed to initialize logger")
					}
					defer logger.Sync()
					return server.RunHTTP(cfg, Name, Version, Revision)
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
