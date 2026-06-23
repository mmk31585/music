---
name: go-mcp-server-generator
description: 'Generate a complete Go MCP server project with proper structure, dependencies, and implementation using the official go-sdk.'
---

# Go MCP Server Generator

Generate a complete, production-ready Model Context Protocol (MCP) server project in Go.

## Structure
```
myserver/
├── go.mod
├── main.go
├── tools/
│   ├── registry.go
│   └── tool1.go
├── resources/
│   └── resource1.go
├── config/
│   └── config.go
├── README.md
└── main_test.go
```

## Key Components
- **main.go**: Server setup, transport, graceful shutdown
- **tools/**: Type-safe tool handlers with JSON schema tags
- **config/**: Environment variable configuration
- **resources/**: MCP resource definitions

## Dependencies
```go
require github.com/modelcontextprotocol/go-sdk v1.0.0
```

## SDK Patterns
- `mcp.NewServer()` with Implementation and Options
- `mcp.AddTool()` with typed input/output structs
- JSON schema tags: `json:"param" jsonschema:"required,description=..."`
- Stdio transport (default), HTTP transport available

## Best Practices
- Validate inputs before processing
- Check `ctx.Err()` for cancellation
- Use table-driven tests for tool handlers
- Graceful shutdown via signal handling
- Structured logging with log/slog
