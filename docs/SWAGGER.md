# Swagger / OpenAPI — Handler Annotation Guide

Muse uses [swaggo](https://github.com/swaggo/swag) to auto-generate OpenAPI 3.0 specs from Go handler annotations. The generated spec is served at `/api/v1/swagger/*any`.

## Quick Start

```bash
# Install
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs from handler annotations
swag init -g main.go -d cmd/api,internal/modules,internal/common

# Output (gitignored):
#   docs/docs.go        — compiled Go swagger spec
#   docs/swagger.json   — OpenAPI JSON
#   docs/swagger.yaml   — OpenAPI YAML
```

## How It Works

1. **Annotate handler methods** with swaggo comment directives
2. **Run `swag init`** to generate `docs/docs.go`, `swagger.json`, `swagger.yaml`
3. **Server imports `docs` package** — `docs.SwaggerInfo` is set at runtime in `cmd/api/main.go`
4. **Swagger UI** is served at `GET /api/v1/swagger/*any` (registered in `internal/app/routes.go:42`)

## Entrypoint Metadata

Set in `cmd/api/main.go`:

```go
// @title           Music API
// @version         1.0
// @description     This is the API server for the Music Streaming Application.
// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer " followed by your token
```

Runtime overrides in `main()`:

```go
docs.SwaggerInfo.Title = "Music API"
docs.SwaggerInfo.Description = "Music streaming API documentation"
docs.SwaggerInfo.Version = "1.0"
docs.SwaggerInfo.Host = "localhost:8080"
docs.SwaggerInfo.BasePath = "/api/v1"
```

## Handler Annotation Format

Each handler method gets a doc comment block before the function:

```go
// ListPublic godoc
// @Summary      List public tracks
// @Description  Returns a paginated list of publicly visible tracks.
// @Tags         tracks
// @Produce      json
// @Param        limit   query  int  false  "Maximum number of items to return"
// @Param        offset  query  int  false  "Number of items to skip"
// @Success      200  {array}   Track
// @Failure      500  {object}  map[string]interface{}
// @Router       /tracks [get]
func (h *Handler) ListPublic(c *gin.Context) { ... }
```

### Directive Reference

| Directive | Format | Required | Purpose |
|-----------|--------|----------|---------|
| `@Summary` | `text` | Yes | Short endpoint name |
| `@Description` | `text` | No | Extended description |
| `@Tags` | `tag` | Yes | Group endpoints by module |
| `@Accept` | `json\|xml\|mpfd` | No | Request content type |
| `@Produce` | `json\|xml` | No | Response content type |
| `@Param` | `name {type} location required description` | Varies | URL/query/body/header params |
| `@Success` | `code {type} model` | Yes | Successful response |
| `@Failure` | `code {type} model` | No | Error response |
| `@Security` | `Bearer` | If auth | Requires Bearer token |
| `@Router` | `path [method]` | Yes | URL path + HTTP method |

### @Param Location Values

| Location | Description | Example |
|----------|-------------|---------|
| `path` | URL path segment | `@Param trackID path string true "Track ID"` |
| `query` | URL query param | `@Param limit query int false "Page size"` |
| `body` | Request body | `@Param request body CreateRequest true "Payload"` |
| `header` | Header | `@Param Authorization header string true "Bearer token"` |
| `formData` | Form field | `@Param file formData file true "Audio file"` |

### Response Object Types

```go
// Simple model
// @Success 200 {object} Track

// Array of models
// @Success 200 {array} Track

// Generic map (for errors)
// @Failure 400 {object} map[string]interface{}
```

Models are resolved by swaggo from the type definitions. The model must be in one of the directories passed to `swag init -d`.

## Annotated Modules

All 28 domain modules are annotated. See specific handlers for reference:

| Module | Handler File | Endpoints |
|--------|-------------|-----------|
| auth | `internal/modules/auth/handler.go` | login, register, refresh, logout, me |
| catalog/track | `internal/modules/catalog/track/handler.go` | CRUD tracks |
| catalog/album | `internal/modules/catalog/album/handler.go` | CRUD albums |
| catalog/artist | `internal/modules/catalog/artist/handler.go` | CRUD artists |
| catalog/genre | `internal/modules/catalog/genre/handler.go` | CRUD genres |
| library | `internal/modules/library/handler.go` | Like/unlike, history |
| playlist | `internal/modules/playlist/handler.go` | CRUD playlists, add/remove tracks |
| queue | `internal/modules/queue/handler.go` | Queue management |
| player | `internal/modules/player/handler.go` | Playback info, streaming |
| search | `internal/modules/search/handler.go` | Search catalog |
| recommendation | `internal/modules/recommendation/handler.go` | Popular, best, similar, for-you |
| media | `internal/modules/media/handler.go` | Upload, list, delete |

### Missing Annotations

The following modules have no swagger annotations yet. Add them when adding/modifying routes:

- social (parties, rooms, clubs, discussions)
- follow
- reactions
- lyrics
- analytics
- notification
- moderation
- subscription
- ai
- creator
- gamification
- contribution
- ingestion
- dashboard
- features
- health
- importcmd

## Best Practices

1. **All public handler methods must have annotations** — swaggo ignores un-annotated handlers
2. **Use `{object} map[string]interface{}` for error responses** unless a typed error DTO exists
3. **Model types must be in the `-d` directory scope** — swaggo resolves types by scanning those paths
4. **Regenerate after changing models or handler signatures**:
   ```bash
   swag init -g main.go -d cmd/api,internal/modules,internal/common
   ```
5. **Keep annotations in sync with actual handler logic** — stale annotations are misleading
6. **Use `@Security Bearer` for authenticated-only endpoints** — omitting it marks the endpoint as public
7. **Use `@Tags` consistent with the module package name** for logical grouping in Swagger UI
