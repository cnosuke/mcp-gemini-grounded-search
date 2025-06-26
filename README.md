# MCP Gemini Grounded Search

MCP Gemini Grounded Search is a Go-based MCP server implementation that provides grounded search functionality using Google's Gemini API, allowing MCP clients (e.g., Claude Desktop) to perform web searches and retrieve up-to-date information with sources.

## Features

* MCP Compliance: Provides a JSON‐RPC based interface for tool execution according to the MCP specification.
* Grounded Search: Uses Gemini API to generate search results with source information (attributions).
* Customizable: Configure through config file, environment variables, or command line options.

## Requirements

- Docker (recommended)

For local development:

- Go 1.24 or later
- Gemini API key

## Using with Docker (Recommended)

```bash
docker pull cnosuke/mcp-gemini-grounded-search:latest

docker run -i --rm -e GEMINI_API_KEY="your-api-key" cnosuke/mcp-gemini-grounded-search:latest
```

### Using with Claude Desktop (Docker)

To integrate with Claude Desktop using Docker, add an entry to your `claude_desktop_config.json` file:

```json
{
  "mcpServers": {
    "gemini-search": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "-e", "GEMINI_API_KEY=your-api-key", "cnosuke/mcp-gemini-grounded-search:latest"]
    }
  }
}
```

### Usage with Claude Code (Docker)

To integrate with Claude Code using Docker, type the following command in the terminal:

```sh
claude mcp add-json mcp-gemini-grounded-search '{
  "command": "docker",
  "args": [
    "run",
    "-i",
    "--rm",
    "-e",
    "GEMINI_API_KEY",
    "-e",
    "GEMINI_MODEL_NAME",
    "cnosuke/mcp-gemini-grounded-search:latest"
  ],
  "env": {
    "GEMINI_MODEL_NAME": "gemini-2.5-flash",
    "GEMINI_API_KEY": "<your-gemini-api-key>"
  }
}'
```

## Building and Running (Go Binary)

Alternatively, you can build and run the Go binary directly:

```bash
# Build the server
make bin/mcp-gemini-grounded-search

# Run the server
./bin/mcp-gemini-grounded-search server --api-key="your-api-key"
```

### Using with Claude Desktop (Go Binary)

To integrate with Claude Desktop using the Go binary, add an entry to your `claude_desktop_config.json` file:

```json
{
  "mcpServers": {
    "gemini-search": {
      "command": "./bin/mcp-gemini-grounded-search",
      "args": ["server"],
      "env": {
        "LOG_PATH": "mcp-gemini-grounded-search.log",
        "DEBUG": "false",
        "GEMINI_API_KEY": "your-api-key",
        "GEMINI_MODEL_NAME": "gemini-2.5-flash-preview-04-17"
      }
    }
  }
}
```

## Configuration

The server is configured via a YAML file (default: config.yml). For example:

```yaml
log: 'path/to/mcp-gemini-grounded-search.log' # Log file path, if empty no log will be produced
debug: false # Enable debug mode for verbose logging

gemini:
  api_key: "your-api-key" # Gemini API key
  model_name: "gemini-2.5-flash-preview-04-17" # Gemini model to use
```

You can override configurations using environment variables:
- `LOG_PATH`: Path to log file
- `DEBUG`: Enable debug mode (true/false)
- `GEMINI_API_KEY`: Gemini API key
- `GEMINI_MODEL_NAME`: Gemini model name

## Logging

Logging behavior is controlled through configuration:

- If `log` is set in the config file, logs will be written to the specified file
- If `log` is empty, no logs will be produced
- Set `debug: true` for more verbose logging

## MCP Server Usage

MCP clients interact with the server by sending JSON‐RPC requests to execute various tools. The following MCP tools are supported:

* `search`: Performs a web search using the Gemini API and returns results with source information.
  * Parameters:
    * `query` (string, required): The search query
    * `max_token` (number, optional): Maximum number of tokens for the generated response

  * Response format:
    ```json
    {
      "text": "Generated text content",
      "groundings": [
        {
          "title": "Source title",
          "domain": "Source domain",
          "url": "Source URL"
        },
        ...
      ]
    }
    ```

## Command-Line Parameters

When starting the server, you can specify various settings:

```bash
./bin/mcp-gemini-grounded-search server [options]
```

Options:

- `--config`, `-c`: Path to the configuration file (default: "config.yml")
- `--log`, `-l`: Path to log file (overrides config file)
- `--debug`, `-d`: Enable debug mode (overrides config file)
- `--api-key`, `-k`: Gemini API key (overrides config file)
- `--model`, `-m`: Gemini model name (overrides config file)

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests for improvements or bug fixes. For major changes, open an issue first to discuss your ideas.

## License

This project is licensed under the MIT License.

Author: cnosuke ( x.com/cnosuke )
