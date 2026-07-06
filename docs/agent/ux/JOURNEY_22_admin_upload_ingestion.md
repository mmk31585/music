# Journey 22: Admin — Upload & Ingestion Pipeline

> Full trace: media upload → ingestion drafts → enrichment → finalization.

---

## Upload Architecture

### Two Upload Systems

The platform has **two parallel upload systems** that serve different purposes:

| System | Route | Purpose | Storage |
|--------|-------|---------|---------|
| **Media Upload** | `/admin/media` | Simple file upload (tracks, covers, images) | `media` table in DB |
| **Ingestion Pipeline** | `/admin/ingestion` | Full-featured: upload → enrich → review → finalize | `ingestion_drafts` table |

---

## Media Upload (`PageAdminMedia.vue`)

### Page Layout

```
PageAdminMedia.vue
  ├── AdminSectionHeader ("Media" / "مدیریت رسانه")
  ├── 4 Upload Type Cards:
  │     ├── Track Audio (30MB max, .mp3/.flac/.wav)
  │     ├── Track Cover (5MB max, .jpg/.png)
  │     ├── Album Cover (5MB max, .jpg/.png)
  │     └── Artist Image (5MB max, .jpg/.png)
  └── Media List
        ├── Each item: filename, size, type, date
        └── Delete button per item
```

### Upload Dialog (`UploadMediaDialog.vue`)

```
PrimeVue Dialog wrapping FileUpload component:
  ├── Drag-and-drop zone
  ├── File type validation (check on select)
  ├── File size validation (check on select)
  ├── ProgressBar during upload
  ├── Success: shows uploaded URL + "Copy" button
  ├── Duplicate detection: shows warning if already exists
  └── Error: inline error message
```

**Upload kinds mapped to backend field names:**

| Upload Kind | Field Name | Backend Category |
|-------------|-----------|-----------------|
| Track Audio | `trackAudio` | `track-audio` |
| Track Cover | `trackCover` | `track-covers` |
| Album Cover | `albumCover` | `album-covers` |
| Artist Image | `artistImage` | `artist-images` |

### Upload Flow (Backend)

```
POST /admin/media/upload (multipart/form-data)
  → Parse multipart form (60MB hard limit)
  → Detect field name → map to category
  → Stream file with inline SHA256 checksum
  → Validate MIME type per category
  → Validate file size per category
  → Check duplicates via SHA256 checksum
  → Store file (object storage + DB record)
  → Return UploadResponse { mediaId, url, path, fileName, duplicate, durationSeconds }
```

### Media List

```
GET /admin/media → Media[]
  → Rows: filename, mime_type, file_size, created_at, delete button
  → No pagination (all items fetched at once)
  → No search/filter
  → Delete: confirmation dialog → DELETE /admin/media/:id
```

---

## Ingestion Pipeline (`PageAdminIngestion.vue` — 817 lines)

### Pipeline Stages

```
1. UPLOAD → 2. EXTRACT → 3. ENRICH → 4. REVIEW → 5. FINALIZE
```

### Stage 1: Audio Upload

```
Drag-and-drop file upload zone:
  → Accepts: .mp3, .flac, .wav, .m4a, .ogg
  → Max: 200MB per file
  → Multiple files supported (queued upload)
  → Progress bar per file
  → Status: uploading / processing / done / error
  → On upload complete → auto-extracts metadata
```

### Stage 2: Automatic Metadata Extraction

```typescript
// Uses music-metadata-browser library on frontend
// Extracts: title, artist, album, track number, duration, genre, year, cover art
// Shown immediately in the draft card
// Status: "Extracted automatically"
```

### Stage 3: Enrichment

```
POST /admin/ingestion/drafts/:id/enrich
  → Queries: MusicBrainz, LastFM, Spotify, LRCLib (lyrics), ML service (genre/style)
  → Returns: { musicbrainz, lastfm, spotify, lrclib, ml }
  → Each source: suggestions for title, artist, album, genre, cover art
  → UI: Shows enrichment results per source
  → Admin selects which values to accept
```

### Stage 4: Review & Edit

```
Draft Detail page (PageAdminIngestionReview.vue):
  ├── Cover art preview
  ├── Track metadata fields:
  │     ├── Title (editable)
  │     ├── Persian title (editable)
  │     ├── Artist (searchable from catalog)
  │     ├── Album (searchable from catalog or "New Album")
  │     ├── Genre (dropdown)
  │     ├── Track number, disc number
  │     ├── Duration (auto-detected, editable)
  │     ├── ISRC (editable)
  │     ├── Language (dropdown)
  │     ├── Label (editable)
  │     ├── Composer, Lyricist, Producer (credits)
  │     └── Explicit flag
  ├── Lyrics editor:
  │     ├── Fetch from LRCLIB
  │     ├── AI generation
  │     ├── AI sync (plain text → LRC)
  │     └── AI review
  ├── Image management:
  │     └── Upload/select cover art per entity (track, album, artist)
  └── Action buttons:
        ├── Save Draft
        ├── Finalize (creates actual Track/Album/Artist records)
        └── Reject (deletes draft)
```

### Stage 5: Finalization

```
POST /admin/ingestion/drafts/:id/finalize
  → Creates Track record
  → Creates/updates Artist record(s)
  → Creates/updates Album record (if specified)
  → Links track to album
  → Assigns genre(s)
  → Saves lyrics (per language)
  → Moves media to permanent storage
  → Returns FinalizeResult { track_id, album_id }
```

### Draft List View

```
Ingestion page main view:
  ├── Stats bar: Total drafts, Pending review, Finalized today, Failed
  ├── Draft list: cards with cover, title, artist, status badge, duration
  ├── Filters: Status (all/pending/completed/failed), Date range
  ├── Sort: newest first
  └── Bulk actions: Delete selected
```

| Status | Badge Color | Meaning |
|--------|------------|---------|
| `pending` | Yellow | Uploaded, awaiting review |
| `enriched` | Blue | External data fetched, needs admin review |
| `editing` | Purple | Admin is editing metadata |
| `finalized` | Green | Successfully published to catalog |
| `failed` | Red | Finalization error, needs re-review |

---

## Catalog Import (`PageAdminImport.vue` — 336 lines)

### Import from External Sources

```
Search and import from:
  ├── Spotify — search tracks/albums/artists
  ├── Deezer — search tracks/albums/artists
  ├── MusicBrainz — search by MBID
  └── Last.fm — artist/album info

Search → results list → "Import" button
  → Fetches metadata
  → Creates ingestion draft automatically
  → Takes admin to review page
```

### Import by Artist (`PageAdminImportArtist.vue`)

```
Artist search → artist detail → album/track list
  → Select tracks to import
  → Batch import → creates drafts for all selected
  → Review individually
```

---

## State Matrix

| State | Media Upload | Ingestion Pipeline | Import |
|-------|-------------|-------------------|--------|
| 🟢 Loading | ✅ Upload cards + media list skeleton | ✅ Skeleton draft list + stats | ✅ Skeleton |
| 🟢 Empty (no media/drafts) | ✅ "No media uploaded yet" | ✅ "No drafts yet — upload audio files" | ✅ Search prompt |
| 🟢 Upload in progress | ✅ Progress bar | ✅ Per-file progress | ✅ N/A |
| 🟢 Upload success | ✅ URL + "Copy" button | ✅ Draft created + enriched | ✅ Draft created |
| 🟢 Upload error (wrong type) | ✅ Inline error | ✅ Inline error | ✅ N/A |
| 🟢 Upload error (too large) | ✅ Inline error | ✅ Inline error | ✅ N/A |
| 🟢 Duplicate detected | ✅ Warning + existing info | ✅ N/A (draft system) | ✅ N/A |
| 🟢 Enrichment in progress | N/A | ✅ Spinner per source | ✅ Spinner |
| 🟢 Finalize success | N/A | ✅ Toast + redirect to track | ✅ Toast |
| 🟢 Finalize failure | N/A | ❌ Silent — draft stays "pending" | ❌ Silent |
| 🔴 Large media list (500+) | ❌ No pagination — all items | ✅ Paginated | ✅ Paginated |
| 🔴 Network error during upload | ❌ Generic error, no retry | ❌ File stuck at "uploading" | ❌ Generic error |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-2201 | ⚠️ MAJOR | `PageAdminMedia.vue` | **Media list has no pagination** — fetches all items at once | Slow page load with 500+ media items | Add pagination or virtual scroll |
| F-2202 | ⚠️ MAJOR | `UploadMediaDialog.vue` | Upload error has **no retry button** — user must re-select file | Frustrating for large files on flaky connections | Add retry button |
| F-2203 | ⚠️ MAJOR | `PageAdminIngestion.vue` | Ingestion stats bar uses **"old PrimeVue surface-* theme classes"** instead of admin dark theme | Visual inconsistency with rest of admin | Migrate to admin dark theme tokens |
| F-2204 | 💡 IMPROVE | `UploadMediaDialog.vue` | No **batch upload progress summary** — only per-file progress | Can't see overall progress for multi-file upload | Add batch progress bar |
| F-2205 | 💡 IMPROVE | `PageAdminMedia.vue` | No **search/filter** on media list | Can't find specific file among many | Add search by filename and type filter |
| F-2206 | 💡 IMPROVE | `PageAdminIngestion.vue` | No **notification when enrichment completes** — admin must refresh | Admin doesn't know when external data is ready | Add WebSocket push or polling badge |
| F-2207 | 💡 IMPROVE | `PageAdminIngestionReview.vue` | No **version comparison** for enriched fields | Can't see what changed between original and enriched | Show diff view (original → enriched) |
| F-2208 | 💡 IMPROVE | Ingestion pipeline | **Media and ingestion are disconnected** — media uploads don't appear in drafts | Admin confused about which tool to use | Unify media and ingestion into single workflow |
| F-2209 | 💡 IMPROVE | `PageAdminImport.vue` | Import search has **no rate-limit indicator** | Admin thinks search is broken when rate-limited | Show "API rate limited — try again in X seconds" |
| F-2210 | 💡 IMPROVE | `PageAdminIngestion.vue` | **No undo** on draft delete — permanent immediately | Accidental delete loses work | Add soft-delete or undo toast with 30s timeout |
| F-2211 | 💡 IMPROVE | `PageAdminIngestion.vue` | **formatDuration treats 0 as falsy** — 0-second tracks show "--" | Confusing for very short audio files | Fix formatDuration to handle 0 properly |
| F-2212 | 💡 IMPROVE | `PageAdminImport.vue` | Import creates drafts but **doesn't navigate to review** — stays on import page | Admin must find the draft manually | Auto-navigate to review page after import |
| F-2213 | 💡 IMPROVE | `UploadMediaDialog.vue` | No **bulk delete** for media items | Cleaning up is one-by-one | Add multi-select + bulk delete |
| F-2214 | 💡 IMPROVE | `PageAdminIngestionReview.vue` | **No "finalize and add another"** — finalize returns to draft list | Admin re-opening drafts one at a time | Add "Finalize & Next" button |

## RTL / A11y / Mobile Notes

- ✅ Upload drag-and-drop zone has proper `aria-label`
- ❌ Upload progress bar lacks `aria-valuenow` — screen reader can't announce percentage
- ✅ Media delete confirmation uses `AdminDeleteConfirm` with proper focus trap
- ❌ Ingestion draft cards lack keyboard-selectable actions — must click
- ✅ All upload forms have associated labels
