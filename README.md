# MCP Gemini Grounded Search

MCP Gemini Grounded Search is a Go-based MCP server that provides grounded search functionality using Google's Gemini API. MCP clients such as Claude Desktop and Claude Code can perform real-time web searches and retrieve up-to-date information with source attribution.

## Features

- **MCP Compliance**: JSON-RPC based interface for tool execution per the MCP specification
- **Grounded Search**: Gemini API generates answers with source attributions
- **Two Transport Modes**: stdio (for Claude Desktop / Claude Code) and Streamable HTTP
- **Flexible Configuration**: config file, environment variables, or command-line flags

## Requirements

- Docker (recommended)

For local development:

- Go 1.24 or later
- Gemini API key

## Using with Docker (Recommended)

```bash
docker pull cnosuke/mcp-gemini-grounded-search:latest

docker run -i --rm -e GEMINI_API_KEY="your-api-key" cnosuke/mcp-gemini-grounded-search:latest server
```

### Using with Claude Desktop (Docker)

Add an entry to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "gemini-search": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "-e", "GEMINI_API_KEY=your-api-key", "cnosuke/mcp-gemini-grounded-search:latest", "server"]
    }
  }
}
```

### Using with Claude Code (Docker)

```sh
claude mcp add-json mcp-gemini-grounded-search '{
  "command": "docker",
  "args": [
    "run", "-i", "--rm",
    "-e", "GEMINI_API_KEY",
    "-e", "GEMINI_MODEL_NAME",
    "-e", "GEMINI_THINKING_LEVEL",
    "cnosuke/mcp-gemini-grounded-search:latest",
    "server"
  ],
  "env": {
    "GEMINI_MODEL_NAME": "gemini-3.1-pro-preview",
    "GEMINI_THINKING_LEVEL": "LOW",
    "GEMINI_API_KEY": "<your-gemini-api-key>"
  }
}'
```

## Building and Running (Go Binary)

```bash
# Build
make bin/mcp-gemini-grounded-search

# stdio mode (for Claude Desktop / Claude Code)
./bin/mcp-gemini-grounded-search server --config config.yml

# Streamable HTTP mode
./bin/mcp-gemini-grounded-search httpserver --config config.yml
```

### Using with Claude Desktop (Go Binary)

```json
{
  "mcpServers": {
    "gemini-search": {
      "command": "/path/to/mcp-gemini-grounded-search",
      "args": ["server", "--config", "/path/to/config.yml"],
      "env": {
        "GEMINI_API_KEY": "your-api-key"
      }
    }
  }
}
```

## Streamable HTTP Mode

The `httpserver` subcommand starts an HTTP server compatible with the MCP Streamable HTTP transport.

```bash
HTTP_AUTH_TOKEN=secret GEMINI_API_KEY=your-key \
  ./bin/mcp-gemini-grounded-search httpserver --config config.yml

# Health check (no auth required)
curl http://localhost:8080/health

# MCP endpoint (auth required)
curl -H "Authorization: Bearer secret" http://localhost:8080/mcp
```

HTTP-specific settings can be configured entirely via environment variables — no need to put secrets in config.yml.

## Configuration

### config.yml

```yaml
log: 'path/to/mcp-gemini-grounded-search.log'  # empty = no log output
debug: false

gemini:
  api_key: ''                      # Set via GEMINI_API_KEY env var
  model_name: 'gemini-3.1-pro-preview'
  max_tokens: 5000
  thinking_level: 'LOW'            # Gemini 3.x series: MINIMAL, LOW, MEDIUM, HIGH
  # thinking_budget: 0             # Gemini 2.5 series: token count (0 = disable thinking)

http:
  port: 8080
  endpoint_path: /mcp
  auth_token: ''                   # Set via HTTP_AUTH_TOKEN env var
  allowed_origins: []              # e.g. ['https://example.com'] — empty = allow all
  heartbeat_seconds: 30
```

### Environment Variables

Configuration priority: **defaults → config.yml → environment variables**

| Variable | Description |
|----------|-------------|
| `GEMINI_API_KEY` | Gemini API key (required) |
| `GEMINI_MODEL_NAME` | Model name (default: `gemini-3.1-pro-preview`) |
| `GEMINI_MAX_TOKENS` | Max response tokens (default: 5000) |
| `GEMINI_THINKING_LEVEL` | `MINIMAL` / `LOW` / `MEDIUM` / `HIGH` (Gemini 3.x) |
| `GEMINI_THINKING_BUDGET` | Token budget for thinking (Gemini 2.5; integer required) |
| `GEMINI_QUERY_TEMPLATE` | Custom query template (must contain `%s`) |
| `HTTP_PORT` | HTTP server port (default: 8080) |
| `HTTP_AUTH_TOKEN` | Bearer token for MCP endpoint authentication |
| `HTTP_ENDPOINT_PATH` | MCP endpoint path (default: `/mcp`) |
| `HTTP_ALLOWED_ORIGINS` | Comma-separated allowed CORS origins |
| `HTTP_HEARTBEAT_SECONDS` | SSE heartbeat interval in seconds (default: 30) |
| `LOG_PATH` | Log file path |
| `DEBUG` | Enable debug logging (`true` or `1`) |

## Command-Line Options

### `server` subcommand (stdio)

```bash
./bin/mcp-gemini-grounded-search server [options]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Path to config file (default: `config.yml`) |
| `--log` | `-l` | Log file path |
| `--debug` | `-d` | Enable debug logging |
| `--api-key` | `-k` | Gemini API key |
| `--model` | `-m` | Gemini model name |
| `--thinking-level` | | `MINIMAL` / `LOW` / `MEDIUM` / `HIGH` |

### `httpserver` subcommand (Streamable HTTP)

```bash
./bin/mcp-gemini-grounded-search httpserver [options]
```

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | `-c` | Path to config file (default: `config.yml`) |

All HTTP settings (`port`, `auth_token`, etc.) are configured via environment variables or config.yml.

## MCP Tools

### `search`

Performs a web search using the Gemini API and returns a grounded answer with sources.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `question` | string | Yes | Natural language question to search |
| `max_token` | number | No | Max tokens for the response |
| `thinking_level` | string | No | Override thinking level for this call |

**Response:**

```json
{
  "text": "Generated answer text",
  "groundings": [
    {
      "title": "Source title",
      "domain": "example.com",
      "url": "https://example.com/article"
    }
  ]
}
```

## Logging

- Set `log` in config.yml or `LOG_PATH` env var to write logs to a file
- If `log` is empty, no log file is produced
- Set `debug: true` or `DEBUG=true` for verbose logging

## Contributing

Contributions are welcome. Please fork the repository and submit pull requests for improvements or bug fixes. For major changes, open an issue first to discuss your ideas.

## License

This project is licensed under the MIT License.

Author: cnosuke ( x.com/cnosuke )
