# CLAUDE.md

## Project Overview

Go-based MCP server that provides grounded search via Google Gemini API.
Supports two transport modes: **stdio** (for Claude Desktop / Claude Code) and **Streamable HTTP**.

## Build & Run

```bash
# Build binary
make bin/mcp-gemini-grounded-search

# stdio mode
./bin/mcp-gemini-grounded-search server --config config.yml

# HTTP mode
./bin/mcp-gemini-grounded-search httpserver --config config.yml

# Tests
make test

# Lint
make inspect
```

## Key Files

| Path | Role |
|------|------|
| `main.go` | CLI entrypoint (`server` / `httpserver` subcommands) |
| `config/config.go` | Config loading: defaults → YAML → env vars |
| `config.yml` | Default config (do NOT commit secrets) |
| `server/server.go` | `createMCPServer()` + `RunStdio()` |
| `server/http.go` | `RunHTTP()` — Streamable HTTP transport |
| `server/middleware.go` | Auth token + CORS origin validation middleware |
| `server/tools.go` | MCP tool definitions (`search`) |
| `searcher/searcher.go` | Gemini API client |
| `internal/errors/wrap.go` | Thin error wrapper (replaces cockroachdb/errors) |

## Config & Environment Variables

Configuration priority: **defaults → config.yml → env vars**

| Env Var | Config Key | Notes |
|---------|-----------|-------|
| `GEMINI_API_KEY` | `gemini.api_key` | Required |
| `GEMINI_MODEL_NAME` | `gemini.model_name` | Default: `gemini-3.1-pro-preview` |
| `GEMINI_MAX_TOKENS` | `gemini.max_tokens` | Default: 5000 |
| `GEMINI_THINKING_LEVEL` | `gemini.thinking_level` | MINIMAL/LOW/MEDIUM/HIGH (Gemini 3.x series) |
| `GEMINI_THINKING_BUDGET` | `gemini.thinking_budget` | Token count (Gemini 2.5 series) — fails fast on invalid |
| `GEMINI_QUERY_TEMPLATE` | `gemini.query_template` | Must contain `%s` placeholder |
| `HTTP_PORT` | `http.port` | Default: 8080 |
| `HTTP_AUTH_TOKEN` | `http.auth_token` | Bearer token — set via env, not config.yml |
| `HTTP_ENDPOINT_PATH` | `http.endpoint_path` | Default: `/mcp` |
| `HTTP_ALLOWED_ORIGINS` | `http.allowed_origins` | Comma-separated, e.g. `https://a.com,https://b.com` |
| `HTTP_HEARTBEAT_SECONDS` | `http.heartbeat_seconds` | Default: 30 |
| `LOG_PATH` | `log` | Log file path |
| `DEBUG` | `debug` | `true` or `1` |

## HTTP Middleware Order

```
Client → withOriginValidation (CORS) → withAuthMiddleware (Bearer) → httpServer (/mcp)
                                                                ↑
                                        /health bypasses both middlewares
```

## Coding Conventions

- KISS / YAGNI / Fail-Fast — no speculative abstractions
- No try-catch-style recovery; surface errors immediately
- Comments only where logic is non-obvious
- All source code and comments in English
