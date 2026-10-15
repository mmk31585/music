# Work Plan for Remaining Tasks

Based on analysis of BACKLOG.md and codebase inspection.

## Tasks Needing Work

### 1. Frontend: Reorganize by domain/features layout (architectural)
- Current structure: components/ grouped by type (admin, auth, forms, layouts, music, profile, social, etc.)
- Target structure: features/ grouping by domain (authentication, catalog, playback, etc.)
- Files to move: components/, composables/, services/, stores/ etc.

### 2. Backend: Middleware consolidation
- Audit middleware usage across modules
- Identify opportunities to move duplicated/shared middleware to common/
- Ensure middleware follows consistent patterns

### 3. Frontend: Ingestion SSE endpoint (implement full SSE instead of polling)
- Replace polling composable (`useEnrichmentPolling`) with SSE implementation
- Create backend SSE endpoint for draft enrichment status updates
- Update frontend to use EventSource instead of polling

### 4. Frontend: RTL consistency check and fixes
- Audit admin components for RTL issues
- Fix hardcoded left/right values that should be direction-aware
- Ensure icon directions flip appropriately in RTL
- Check text alignment in RTL contexts

### 5. Frontend: ARIA labels for admin table action buttons
- Add aria-label to icon-only buttons in table actions
- Affected components: AdminArtistsCard, AdminAlbumsCard, AdminGenresCard, etc.
- Buttons needing labels: View (pi pi-eye), Enrich (pi pi-refresh), Edit (pi pi-pencil), Delete (pi pi-trash), etc.

### 6. Backend: RequireRole hierarchy support
- Enhance AuthMiddleware to store role hierarchy level in context
- OR modify RequireRole to check hierarchy levels instead of exact string matching
- Role hierarchy: listener(10) < creator(30) < curator(50) < moderator(80) < admin(100)

### 7. Testing: Playwright E2E smoke test
- Create basic end-to-end test covering critical user flows
- Test: user registration, login, track upload, playback

### 8. Testing: Vitest frontend store tests
- Create unit tests for Pinia stores
- Focus on auth store, catalog store, player store

### 9. Testing: Backend handler tests (upload, admin CRUD)
- Create test files for handler modules
- Test authentication, validation, error handling, success cases

## Tasks Likely Already Done

- Backend DTO/domain separation ✅
- Empty states for all admin tables ✅ (AdminEmptyState component)
- VersionedEvent unused code ✅ (not found)
- RepositoryInterface unused in auth ❌ (actually used correctly as pattern)
- Album upsert embedded cover fallback ✅ (appears implemented in ingestion finalization)

## Execution Plan

Use subagents to work on independent tasks in parallel where possible.