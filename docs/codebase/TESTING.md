# Muse — Testing Strategy

## Current Testing State

### Backend (Go)
- **Framework**: testify + go-sqlmock + table-driven tests
- **Coverage threshold**: 30% (enforced in CI via `cover-ci`)
- **Test command**: `go test ./... -race -count=1`
- **Coverage command**: `go test ./... -race -coverprofile=coverage.out -covermode=atomic`
- Test files: `*_test.go` alongside source

### Modules with Tests
Modules that have `_test.go` files:
- `auth/` — auth flows, middleware
- `search/` — Persian-aware ranking
- `history/` — signal classification
- `playlist/` — collaborative features
- `features/` — feature flags
- `app/` — app bootstrap

### Frontend (Vue 3)
- **Framework**: Vitest + Vue Test Utils + jsdom
- **Test command**: `cd frontend && npm run test` (not specified, likely `vitest run`)
- **Setup**: `tests/setup.ts` in frontend/src
- Test files: co-located `__tests__/` or `.spec.ts`

### Python ML Service
- **Framework**: pytest
- **Test directory**: `moja-ml-service/tests/`

## Test Gaps (from TODO.md)

| Area | Gap | Priority |
|------|-----|----------|
| **API handlers** | Most handlers lack tests | CRITICAL |
| **Repository layer** | Only auth/features have repo tests | HIGH |
| **Frontend components** | No component tests found | HIGH |
| **Frontend stores** | No store unit tests | HIGH |
| **Integration tests** | No end-to-end API tests | HIGH |
| **WebSocket** | No socket tests | MEDIUM |
| **ML service** | Minimal test coverage | MEDIUM |

## Target Architecture (from docs/architecture/OPS.md)

### Testing in CI Pipeline
```
Commit → Lint → Unit Test → Build → Integration → Security Scan → Deploy
```

### Test Types (Aspirational)

| Type | Cadence | Scope |
|------|---------|-------|
| Unit | Every commit | Individual functions, models, repositories |
| Integration | Every merge to main | API endpoints, DB queries, service boundaries |
| Component | Every merge to main | Frontend components (Vitest + Vue Test Utils) |
| E2E | Weekly (CI) | Full user flows (Playwright) |
| Performance | Weekly (CI) | Load testing (k6) |
| Security | Weekly (CI) | SAST + dependency scan |
| Chaos | Monthly | Resilience testing |
| Accessibility | Per PR | WCAG 2.2 AA audit |

## Testing Patterns

### Go Table-Driven Tests
```go
func TestAuthenticate(t *testing.T) {
    tests := []struct {
        name     string
        email    string
        password string
        wantErr  bool
    }{
        {"valid credentials", "user@test.com", "password123", false},
        {"wrong password", "user@test.com", "wrong", true},
        {"empty email", "", "password123", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### Mocking Strategy
- **Database**: `go-sqlmock` for repository tests
- **External APIs**: Interface-based, mock implementations
- **Event bus**: In-process bus subscribable in tests
- **Storage**: `Local` storage implementation for test isolation

## Running Tests

```bash
# All Go tests with race detection
make test

# Coverage report
make cover-html

# CI coverage check (30% threshold)
make cover-ci

# Frontend tests
cd frontend && npx vitest run

# ML service tests
cd moja-ml-service && uv run pytest
```
