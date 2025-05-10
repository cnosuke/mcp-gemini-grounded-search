# MCP Greeting Server

MCP Greeting Server is a Go-based MCP server implementation that provides basic greeting functionality, allowing MCP clients (e.g., Claude Desktop) to generate greeting messages.

## Features

* MCP Compliance: Provides a JSON‐RPC based interface for tool execution according to the MCP specification.
* Greeting Operations: Supports generating greeting messages, with options for personalization.

## Requirements

- Docker (recommended)

For local development:

- Go 1.24 or later

## Using with Docker (Recommended)

```bash
docker pull cnosuke/mcp-greeting:latest

docker run -i --rm cnosuke/mcp-greeting:latest
```

### Using with Claude Desktop (Docker)

To integrate with Claude Desktop using Docker, add an entry to your `claude_desktop_config.json` file:

```json
{
  "mcpServers": {
    "greeting": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "cnosuke/mcp-greeting:latest"]
    }
  }
}
```

## Building and Running (Go Binary)

Alternatively, you can build and run the Go binary directly:

```bash
# Build the server
make bin/mcp-greeting

# Run the server
./bin/mcp-greeting server --config=config.yml
```

### Using with Claude Desktop (Go Binary)

To integrate with Claude Desktop using the Go binary, add an entry to your `claude_desktop_config.json` file:

```json
{
  "mcpServers": {
    "greeting": {
      "command": "./bin/mcp-greeting",
      "args": ["server"],
      "env": {
        "LOG_PATH": "mcp-greeting.log",
        "DEBUG": "false",
        "GREETING_DEFAULT_MESSAGE": "こんにちは"
      }
    }
  }
}
```

## Configuration

The server is configured via a YAML file (default: config.yml). For example:

```yaml
log: 'path/to/mcp-greeting.log' # Log file path, if empty no log will be produced
debug: false # Enable debug mode for verbose logging

greeting:
  default_message: "こんにちは！"
```

Note: The default greeting message can also be injected via an environment variable `GREETING_DEFAULT_MESSAGE`. If this environment variable is set, it will override the value in the configuration file.

You can override configurations using environment variables:
- `LOG_PATH`: Path to log file
- `DEBUG`: Enable debug mode (true/false)
- `GREETING_DEFAULT_MESSAGE`: Default greeting message

## Logging

Logging behavior is controlled through configuration:

- If `log` is set in the config file, logs will be written to the specified file
- If `log` is empty, no logs will be produced
- Set `debug: true` for more verbose logging

## MCP Server Usage

MCP clients interact with the server by sending JSON‐RPC requests to execute various tools. The following MCP tools are supported:

* `greeting/hello`: Generates a greeting message, with an optional name parameter for personalization.

## Command-Line Parameters

When starting the server, you can specify various settings:

```bash
./bin/mcp-greeting server [options]
```

Options:

- `--config`, `-c`: Path to the configuration file (default: "config.yml").

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests for improvements or bug fixes. For major changes, open an issue first to discuss your ideas.

## License

This project is licensed under the MIT License.

Author: cnosuke ( x.com/cnosuke )