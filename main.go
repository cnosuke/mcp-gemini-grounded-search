package main

import (
	"fmt"
	"os"

	"github.com/cnosuke/mcp-gemini-grounded-search/config"
	"github.com/cnosuke/mcp-gemini-grounded-search/logger"
	"github.com/cnosuke/mcp-gemini-grounded-search/server"
	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"
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
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s (%s)", Version, Revision)
	app.Name = Name
	app.Usage = Usage

	app.Commands = []*cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Start the MCP server for Gemini grounded search",
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
					EnvVars: []string{"LOG_PATH"},
				},
				&cli.BoolFlag{
					Name:    "debug",
					Aliases: []string{"d"},
					Usage:   "enable debug logging (overrides config file)",
					EnvVars: []string{"DEBUG"},
				},
				&cli.StringFlag{
					Name:    "api-key",
					Aliases: []string{"k"},
					Usage:   "Gemini API key (overrides config file)",
					EnvVars: []string{"GEMINI_API_KEY"},
				},
				&cli.StringFlag{
					Name:    "model",
					Aliases: []string{"m"},
					Usage:   "Gemini model name (overrides config file)",
					EnvVars: []string{"GEMINI_MODEL_NAME"},
				},
			},
			Action: func(c *cli.Context) error {
				configPath := c.String("config")

				// Read the configuration file
				cfg, err := config.LoadConfig(configPath)
				if err != nil {
					return errors.Wrap(err, "failed to load configuration file")
				}

				// Override config with command line flags
				if c.IsSet("log") {
					cfg.Log = c.String("log")
				}
				if c.IsSet("debug") {
					cfg.Debug = c.Bool("debug")
				}
				if c.IsSet("api-key") {
					cfg.Gemini.APIKey = c.String("api-key")
				}
				if c.IsSet("model") {
					cfg.Gemini.ModelName = c.String("model")
				}

				// Verify required configurations
				if cfg.Gemini.APIKey == "" {
					return errors.New("Gemini API key is required. Set it in config.yml or use --api-key flag or GEMINI_API_KEY environment variable")
				}

				// Initialize logger
				if err := logger.InitLogger(cfg.Debug, cfg.Log); err != nil {
					return errors.Wrap(err, "failed to initialize logger")
				}
				defer logger.Sync()

				return server.Run(cfg, Name, Version, Revision)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
