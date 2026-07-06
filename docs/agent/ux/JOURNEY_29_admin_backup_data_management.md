# Journey 29: Admin — Backup & Data Management

> Full trace: backup infrastructure → data migration → storage management → cleanup operations.

---

## Current Backup Infrastructure

### What Exists

There is **no backup management UI** in the application. Backup infrastructure, if it exists, is handled outside the application (PostgreSQL dumps, cloud storage snapshots, etc.).

### Database Level

```sql
-- PostgreSQL database
-- Standard migrations in /migrations/ directory
-- No backup/restore endpoints in the API
```

### Storage Level

```go
// internal/modules/media/storage.go
// Files stored in object storage (S3-compatible)
// No backup management API
```

### Known Cleanup Worker

```go
// internal/app/worker.go
// cleanup worker — runs every 15 minutes
// Purpose: Expired sessions, stale drafts, old notifications
```

This is the **only automated data management** in the system.

---

## What's Missing

### Backup Operations

| Feature | Status | Importance |
|---------|--------|------------|
| On-demand database backup trigger | ❌ Missing | High — disaster recovery |
| Backup listing and download | ❌ Missing | High — restore capability |
| Scheduled backup configuration | ❌ Missing | Medium — automation |
| Backup restoration UI | ❌ Missing | Low — usually manual |
| Media file backup / sync | ❌ Missing | Medium — content protection |

### Storage Management

| Feature | Status | Importance |
|---------|--------|------------|
| Storage usage dashboard | ❌ Missing | Medium — cost monitoring |
| Orphaned file detection | ❌ Missing | Medium — storage waste |
| Media migration between providers | ❌ Missing | Low — infrastructure change |
| Storage quota management | ❌ Missing | Low — usage limits |
| CDN cache invalidation | ❌ Missing | Medium — content updates |

### Cleanup Operations

| Feature | Status | Importance |
|---------|--------|------------|
| Orphaned media cleanup | ❌ Missing | Medium — storage waste |
| Unused draft cleanup | ✅ Auto (15 min) | Already handled |
| Expired session cleanup | ✅ Auto (15 min) | Already handled |
| Old notification cleanup | ✅ Auto (15 min) | Already handled |
| Analytics event archive | ❌ Missing | Medium — data management |
| Temporary file cleanup | ❌ Missing | Low |
| Soft-deleted content purge | ❌ Missing | Medium — storage |

---

## Ideal Data Management Dashboard

```
Data Management Dashboard (missing)

├── Backup Status
│     ├── Last backup: [date/time]
│     ├── Backup size: [size]
│     ├── Backup frequency: [config]
│     └── [Trigger Backup Now] button
│
├── Storage Usage
│     ├── Total storage used
│     ├── By category: audio, images, videos
│     ├── Orphaned files count + size
│     └── [Clean Up Orphans] button
│
├── Database
│     ├── Table sizes
│     ├── Row counts per table
│     ├── Index size
│     └── [Analyze] / [Vacuum] buttons
│
└── Cleanup Operations
      ├── Orphaned media
      ├── Expired data
      ├── Temporary files
      └── [Run Cleanup] button with report
```

---

## Media Migration

Currently, media is stored in a single storage provider. There is **no migration tool** within the app:

- Storage provider is configured via environment variables
- Changing providers requires manual data transfer
- No dual-write or migration mode
- No CDN URL swapping mechanism

---

## State Matrix

| State | Backup UI | Storage Dashboard | Cleanup Tools | Migration Tools |
|-------|-----------|-------------------|---------------|-----------------|
| 🟢 Feature exists | ❌ None | ❌ None | ⚠️ Partial (auto-cleanup worker only) | ❌ None |
| 🟢 Manual trigger | ❌ | ❌ | ❌ | ❌ |
| 🟢 View status | ❌ | ❌ | ❌ | ❌ |
| 🟢 Scheduled | ❌ | ❌ | ✅ Auto every 15 min | ❌ |
| 🔴 Need on-demand backup | ❌ No trigger | N/A | N/A | N/A |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-2901 | ⚠️ MAJOR | Platform | **No on-demand backup trigger** — no API or UI to start a backup | Can't take a pre-deployment snapshot | Add `POST /admin/backup` endpoint with notification |
| F-2902 | ⚠️ MAJOR | Platform | **No storage usage dashboard** — can't see how much space is used | Storage costs can grow unbounded without visibility | Add storage usage page with per-category breakdown |
| F-2903 | ⚠️ MAJOR | `internal/modules/media/` | **No orphaned file detection** — media records without corresponding files (and vice versa) | Storage waste from orphaned uploads | Add orphan detection + cleanup tool |
| F-2904 | 💡 IMPROVE | Platform | **No backup listing** — can't see existing backups, their size, or date | Can't verify backup health | Add backup list with metadata (size, date, status) |
| F-2905 | 💡 IMPROVE | Platform | **No backup download** — can't retrieve a backup from the admin UI | Must access storage directly | Add one-time download URL generation |
| F-2906 | 💡 IMPROVE | Platform | **No scheduled backup configuration** — frequency, retention, time of day | Backups may be inconsistent | Add backup schedule UI |
| F-2907 | 💡 IMPROVE | Platform | **No CDN cache invalidation UI** — must use provider console | Delayed content updates after media replacement | Add "Purge CDN for this file" button |
| F-2908 | 💡 IMPROVE | Platform | **No database maintenance UI** — analyze, vacuum, reindex | DB performance degrades over time | Add "Run Maintenance" button with progress |
| F-2909 | 💡 IMPROVE | Platform | **No export/import for catalog data** — can't bulk export/import tracks as CSV/JSON | Manual data migration between instances | Add catalog export/import wizard |
| F-2910 | 💡 IMPROVE | Platform | **Soft-delete purge missing** — deleted tracks, artists, albums may still occupy space | Storage grows even after content deletion | Add soft-delete cleanup job with configurable TTL |

## RTL / A11y / Mobile Notes

- ❌ Backup page would need progress indicators with proper `aria-valuenow` for long-running operations
- ✅ Storage usage charts should use accessible data tables as fallback for screen readers
- ❌ "Trigger Backup" button should have confirmation with "This may take several minutes" warning
