# Journey 30: Admin — Advanced Administration

> Full trace: bulk operations → worker management → import pipeline → content enrichment → admin UX patterns.

---

## Bulk Operations

### Current State

Bulk operations are **severely limited** in the admin panel:

| Operation | Status | Detail |
|-----------|--------|--------|
| **Bulk moderation (resolve/dismiss)** | ✅ Exists | Multi-select reports + bulk action |
| **Bulk track enrichment** | ⚠️ Partial | "Enrich All" button but no progress indicator |
| **Bulk user actions (suspend/ban)** | ❌ Missing | One at a time only |
| **Bulk track operations (genre, delete)** | ❌ Missing | One at a time only |
| **Bulk album operations** | ❌ Missing | One at a time only |
| **Bulk contribution approval** | ❌ Missing | F-2512 noted |
| **Bulk video operations** | ❌ Missing | F-2514 noted |
| **Bulk media deletion** | ❌ Missing | F-2213 noted |

### Ideal Bulk Operations Pattern

```
Select items via checkbox → action toolbar appears:
  → [Selected: 12 tracks]
  → Actions: [Add Genre] [Remove Genre] [Set Explicit] [Delete] [Export]
  → Confirmation dialog with summary:
      "You are about to delete 12 tracks. This cannot be undone."
  → Progress bar: "Processing 12 of 12..."
  → Result summary: "12 tracks deleted. 0 errors."
```

---

## Worker Management

### Worker Infrastructure

```go
// internal/app/worker.go
// 9 workers defined, 6 wired into the system:

type WorkerConfig struct {
    Name     string
    Interval time.Duration
    Handler  func(ctx context.Context) error
}

// Wired workers:
//   analytics    — 1 min  — daily rollup
//   recommendation — 5 min — content similarity
//   notification  — 30 sec — notification outbox
//   email        — ?      — email queue
//   cleanup      — 15 min — expired sessions, stale drafts
//   import       — event  — import job consumer (Redis stream)
//
// Defined but NOT wired:
//   ai, eventbus, indexer, recommender_v2, transcoder
```

### Worker Monitoring Gaps

| Need | Current State |
|------|---------------|
| See which workers are running | ❌ No admin visibility |
| View last run time + duration | ❌ Not tracked in accessible store |
| View error count | ❌ No per-worker error metrics |
| Trigger worker manually | ❌ No "Run Now" button |
| Pause/resume worker | ❌ No lifecycle control |
| View worker logs | ❌ No worker log viewer |
| Configure worker interval | ❌ Via code/config only |

### Ideal Worker Dashboard

```
Worker Management (missing)

├── Worker Status:
│     ├── ✅ analytics — Last run: 30s ago — Duration: 2.1s — 0 errors
│     ├── ✅ recommendation — Last run: 4m ago — Duration: 8.3s — 0 errors
│     ├── ⏳ notification — Last run: 12s ago — Duration: 0.3s — 2 errors
│     ├── ❌ ai — DISABLED — Not wired
│     └── ...
│
├── Per-worker actions:
│     ├── [Run Now] — triggers immediate execution
│     ├── [Pause] / [Resume] — lifecycle control
│     ├── [Configure] — interval, batch size, etc.
│     └── [Logs] — recent worker output
│
└── Global controls:
      ├── [Restart All Workers]
      └── Worker health: 🟢 6/6 running
```

---

## Import Pipeline

### Import from External Sources

```go
// internal/modules/importcmd/
// Multi-source import: Spotify, Deezer, MusicBrainz, Last.fm
// Frontend: PageAdminImport.vue + PageAdminImportArtist.vue
```

### Import Flow

```
Admin searches external source → results displayed → "Import" clicked
  → Creates ingestion draft
  → Enriches draft with external metadata
  → Admin reviews + finalizes
```

### Import Gaps

| Gap | Impact |
|-----|--------|
| **No import queue status** — can't see pending/in-progress/completed imports | Admin doesn't know if import is still working |
| **No batch import progress** — importing an artist's full discography shows no progress | Long imports appear stalled |
| **No import error details** — failed imports show generic error | Can't diagnose failures |
| **No import history** — can't see what was imported and when | Duplicate imports possible |
| **No scheduled import** — can't set "import new releases from followed artists weekly" | Manual process only |

---

## Content Enrichment Pipeline

### Current Enrichment Flow

```go
// POST /admin/tracks/:id/enrich
// POST /admin/artists/:id/enrich
// POST /admin/tracks/enrich-all
// POST /admin/ingestion/drafts/:id/enrich
```

Sources: MusicBrainz, LastFM, Spotify, LRCLib, ML service.

### Enrichment Gaps

| Gap | Impact |
|-----|--------|
| **No enrichment status dashboard** | Can't see how many items are enriched vs not |
| **No "enrich by criteria"** — can't say "enrich all tracks missing ISRC" | Manual per-artist enrichment |
| **No enrichment queue** — bulk enrich fires all at once without queue | Server strain on large catalogs |
| **No enrichment comparison** — can't preview changes before applying | Admin may accept wrong data |
| **No enrichment rollback** — accidental enrichment can't be undone | Data quality risk |

---

## Admin UX Pattern Library

### Shared Admin Patterns

The admin uses several shared patterns that should be cataloged:

| Pattern | Component | Usage |
|---------|-----------|-------|
| **Section Header** | `AdminSectionHeader.vue` | Eyebrow + title + description + actions slot — used on every page |
| **Stat Card** | `AdminStatCard.vue` | Metric with icon, value, hint, 4 color themes, loading state |
| **Empty State** | `AdminEmptyState.vue` | Icon + title + description + action CTA |
| **Delete Confirm** | `AdminDeleteConfirm.vue` | Reusable confirmation dialog |
| **Form Dialog** | TrackFormDialog, ArtistFormDialog, AlbumFormDialog | Modal CRUD forms |
| **Upload Dialog** | `UploadMediaDialog.vue` | Drag-and-drop with progress |
| **CRUD Composable** | `useAdminTracks`, `useAdminArtists`, etc. | Standard fetch/create/update/delete pattern |

### Admin UX Anti-Patterns

| Anti-Pattern | Location | Problem |
|-------------|----------|---------|
| **No pagination on lists** | Media list, track list (historically) | Slow with large datasets |
| **Any type used** | `PageAdminTracks.vue` | TypeScript safety defeated |
| **Event-based coupling** | TrackFormDialog inline creation | Race conditions |
| **Silent error catching** | Multiple pages | Admin doesn't know something failed |
| **Aggressive Promise.all** | Catalog preview fetching all items | Wasteful |
| **Old theme classes** | Ingestion page | Visual inconsistency |
| **No unsaved changes warning** | All form dialogs | Accidental data loss |

---

## State Matrix

| State | Bulk Operations | Worker Management | Import Pipeline | Enrichment |
|-------|----------------|-------------------|-----------------|------------|
| 🟢 Page exists | ⚠️ Partial (moderation only) | ❌ None | ✅ Import pages exist | ❌ No status page |
| 🟢 Bulk action UI | ✅ Moderation only | ❌ | ❌ | ❌ |
| 🟢 Progress indication | ❌ No progress bar | ❌ | ❌ | ❌ |
| 🟢 Error reporting | ❌ Silent | ❌ | ❌ Catch but no detail | ❌ Silent |
| 🟢 History/log | ❌ | ❌ | ❌ No import history | ❌ |
| 🟢 Manual trigger | ✅ Moderation bulk | ❌ | ✅ Search → Import | ✅ Per-item + enrich-all |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-3001 | ⚠️ MAJOR | All admin lists | **No bulk operations on tracks, users, albums, artists** — only moderation has bulk | Slow admin workflows for large catalogs | Add multi-select + bulk action toolbar to all list pages |
| F-3002 | ⚠️ MAJOR | `internal/app/worker.go` | **No worker monitoring UI** — can't see which workers are running, their status, or errors | Silent worker failures go undetected | Add worker status page with per-worker metrics |
| F-3003 | ⚠️ MAJOR | `PageAdminImport.vue` | **Import has no queue or progress** — no status for in-progress imports | Admin unsure if import is working | Add import queue with progress tracking |
| F-3004 | 💡 IMPROVE | Enrichment pipeline | **No enrichment status dashboard** — can't see what's been enriched vs not | Incomplete metadata goes unnoticed | Add enrichment coverage report |
| F-3005 | 💡 IMPROVE | All admin pages | **No bulk operation progress bar** — text-only status for batch actions | No visual feedback during long operations | Add progress bar with item count + ETA |
| F-3006 | 💡 IMPROVE | Admin patterns | **No keyboard shortcuts** in any admin page (except possibly moderation) | Slow for power users | Add admin-wide keyboard shortcuts (n=next, p=prev, /=search, ?=help) |
| F-3007 | 💡 IMPROVE | Import pipeline | **No import history** — can't see what was imported when | Accidental duplicate imports | Add import log with source, count, date, user |
| F-3008 | 💡 IMPROVE | Workers | **No "Run Now" for workers** — can't trigger a worker manually | Must wait for next interval to test changes | Add admin button to trigger worker immediately |
| F-3009 | 💡 IMPROVE | Workers | **No worker interval configuration** — hardcoded in code | Can't adjust frequency without redeploy | Add worker interval config in admin settings |
| F-3010 | 💡 IMPROVE | All admin | **No "recent actions" widget** on dashboard | Can't see latest admin activity | Add recent actions feed to dashboard |
| F-3011 | 💡 IMPROVE | All admin pages | **No bulk edit dialog** — can't edit multiple tracks' genres at once | One-by-one for large catalogs | Add bulk edit modal (change genre, explicit, etc.) |
| F-3012 | 💡 IMPROVE | Admin | **No admin tour / onboarding** — new admins shown no guidance | Learning curve for admin tools | Add admin onboarding overlay on first visit |
| F-3013 | 💡 IMPROVE | Enrichment | **No enrichment rollback** — can't undo an enrichment action | Data quality risk from bad enrichment | Add "Undo enrichment" with previous snapshot |
| F-3014 | 💡 IMPROVE | All admin | **No admin dark/light mode** — follows app theme only | Admin may prefer different theme | Add independent admin theme setting |

## RTL / A11y / Mobile Notes

- ✅ Existing admin bulk select pattern (moderation) has proper `aria-label` on checkboxes
- ❌ Worker management page would need accessible status indicators (color + icon + text)
- ✅ Import flow uses existing form patterns with proper labels
- ❌ Admin keyboard shortcuts should be documented in a help overlay triggered by `?`
- ✅ Bulk operation progress bar needs `aria-valuenow` for screen reader announcements
