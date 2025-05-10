package main

import (
	"fmt"
	"os"

	"github.com/cnosuke/mcp-greeting/config"
	"github.com/cnosuke/mcp-greeting/logger"
	"github.com/cnosuke/mcp-greeting/server"
	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"
)

var (
	// Version and Revision are replaced when building.
	// To set specific version, edit Makefile.
	Version  = "0.0.1"
	Revision = "xxx"

	Name  = "mcp-greeting"
	Usage = "A simple MCP server implementation for greetings"
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
			Usage:   "A simple MCP server implementation for greetings",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "config",
					Aliases: []string{"c"},
					Value:   "config.yml",
					Usage:   "path to the configuration file",
				},
			},
			Action: func(c *cli.Context) error {
				configPath := c.String("config")

				// Read the configuration file
				cfg, err := config.LoadConfig(configPath)
				if err != nil {
					return errors.Wrap(err, "failed to load configuration file")
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
	}
}
