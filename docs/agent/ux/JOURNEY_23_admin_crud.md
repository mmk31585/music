# Journey 23: Admin — Track / Album / Artist CRUD

> Full trace: admin catalog list views → detail pages → create/edit flows → delete flows.

---

## Admin Catalog Architecture

All catalog CRUD follows a consistent pattern across tracks, artists, albums, and genres:

```
List Page (table or grid)
  ├── Search bar
  ├── Sort / Filter controls
  ├── "Add New" button
  └── Item list with inline actions (Edit / Delete)

Detail Page (optional)
  ├── Full item detail
  ├── Edit button → opens form dialog
  └── Delete button → confirmation dialog

Form Dialog (modal)
  ├── Fields for the entity
  ├── Save / Cancel
  └── On save → refreshes list
```

Each entity has: a page component, a card/list component, a form dialog, and a composable CRUD service.

---

## Tracks Management (`PageAdminTracks.vue` — 969 lines)

### Track List View

```
Search bar (by title, artist, album, genres)
Two view modes:
  ├── Table mode: columns (play, cover+title, artists, album, genres, duration, actions)
  └── Card mode: grid of cards

Table columns:
  ├── Play button → plays track preview
  ├── Title + cover art thumbnail
  ├── Artists (comma-separated names)
  ├── Album link
  ├── Genres (tag chips)
  ├── Duration (mm:ss)
  └── Actions: Edit → TrackFormDialog, Delete → AdminDeleteConfirm

Action buttons:
  ├── "Add Track" → TrackFormDialog (create mode)
  ├── "Enrich All" → POST /admin/tracks/enrich-all (batch enrichment)
  └── Per-track "Enrich" → POST /admin/tracks/:id/enrich
```

### Data Fetching

```typescript
// useAdminTracks.ts composable
const tracks = ref<Track[]>([])
const loading = ref(false)
const saving = ref(false)
const error = ref<string | null>(null)

async function fetchTracks() {
  loading.value = true
  try {
    const { data } = await adminTracksApi.list()
    tracks.value = data
  } catch (e) {
    error.value = 'Failed to load tracks'
  } finally {
    loading.value = false
  }
}
```

### Track Detail (`PageAdminTrackDetail.vue` — 524 lines)

```
Layout:
  ├── Cover art (large)
  ├── Title (primary + Persian)
  ├── Artist link → admin.artist.detail
  ├── Album link → admin.album.detail
  ├── Genre chips
  ├── Duration, play count
  ├── "Also appears on" section (other albums containing this track)
  └── Action buttons:
        ├── Edit Track → TrackFormDialog
        ├── Fetch LRC → lyrics pipeline
        ├── Fetch All → enrich all metadata
        └── Delete → AdminDeleteConfirm
```

### Track Form Dialog (`TrackFormDialog.vue` — 1728+ lines)

**Largest form in the entire app.** Sections:

```
Audio file:
  ├── Upload new audio (or shows current file)
  ├── Auto-metadata extraction via music-metadata-browser
  ├── Audio preview with seek bar
  └── Duration auto-filled from file

Cover image:
  ├── Upload new image
  ├── Shows current cover (if editing)
  └── Accepted: jpg, png, max 5MB

Basic info:
  ├── Title (required)
  ├── Persian title (optional)
  ├── Artist: AutoComplete with API search (can create new inline)
  ├── Featured artists: multi-select AutoComplete
  ├── Album: AutoComplete with search (can create new inline)
  └── Duration (auto-filled, editable)

Classification:
  ├── Genre(s): multi-select dropdown
  ├── Track number
  ├── Disc number
  ├── Language: dropdown
  └── Explicit: toggle

Credits:
  ├── Composer
  ├── Lyricist
  ├── Producer
  └── Label

Identifiers:
  ├── ISRC code
  └── Language

Lyrics:
  ├── Fetch from LRCLIB (with iTunes cover art fallback)
  ├── AI generation pipeline (with polling)
  ├── AI sync (plain text → LRC timed format)
  ├── AI review
  ├── Create/update/delete per language
  └── Multi-language support
```

---

## Artists Management

### List View (`AdminArtistsCard.vue` — 281 lines)

```
Searchable artist grid:
  ├── Artist image (or fallback)
  ├── Name
  ├── Verified badge
  └── Actions: Edit → ArtistFormDialog, Delete → AdminDeleteConfirm

"Add Artist" button → ArtistFormDialog
```

### Detail (`PageAdminArtistDetail.vue` — 363 lines)

```
Layout:
  ├── Artist image (large)
  ├── Name + verified badge
  ├── Bio (expandable)
  ├── Enrichment support button
  └── Track list (tracks by this artist)
```

### Artist Form Dialog (`ArtistFormDialog.vue` — 271 lines)

```
Fields:
  ├── Name (required)
  ├── Bio (textarea)
  ├── Image: upload via media API or URL
  └── "Enrich data" button → fetches external artist data
```

---

## Albums Management

### List View (`AdminAlbumsCard.vue` — 252 lines)

```
Searchable album grid:
  ├── Cover art
  ├── Title
  ├── Artist name
  ├── Track count
  └── Actions: Edit → AlbumFormDialog, Delete → AdminDeleteConfirm

"Add Album" button → AlbumFormDialog
```

### Detail (`PageAdminAlbumDetail.vue` — 327 lines)

```
Layout:
  ├── Cover art (large)
  ├── Title + artist
  ├── Year, track count, total duration
  ├── Track list (ordered)
  └── Actions: Edit, Delete
```

### Album Form Dialog (`AlbumFormDialog.vue` — 442 lines)

```
Fields:
  ├── Title (required)
  ├── Persian title (optional)
  ├── Artist (AutoComplete with search, can create inline)
  ├── Cover: upload via media API
  ├── Release date: date picker
  ├── Genre(s): multi-select
  ├── Label (optional)
  └── Type: album / single / EP / compilation
```

---

## Genres Management

### List View (`AdminGenresCard.vue` — 226 lines)

```
Tag-style genre display:
  ├── Genre name
  ├── Color indicator (swatch)
  ├── Track count
  └── Actions: Edit → dialog, Delete → confirmation

"Add Genre" button → dialog
```

### Genre Form Dialog

```
Fields:
  ├── Name (required, unique)
  ├── Persian name (optional)
  ├── Color: color picker (swatch for visual identification)
  ├── Description (optional)
  └── Parent genre: optional hierarchy
```

---

## CRUD Composables

Each composable follows the same pattern:

```typescript
// useAdminTracks.ts — 5 exports:
export function useAdminTracks() {
  const tracks = ref<Track[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const deleting = ref(false)
  const error = ref<string | null>(null)

  async function fetchTracks() { /* ... */ }
  async function createTrack(payload) { /* ... */ }
  async function updateTrack(id, payload) { /* ... */ }
  async function deleteTrack(id) { /* ... */ }

  return { tracks, loading, saving, deleting, error, fetchTracks, createTrack, updateCard, deleteTrack }
}

// Aggregated into useAdminCatalog.ts:
export function useAdminCatalog() {
  const tracks = useAdminTracks()
  const artists = useAdminArtists()
  const albums = useAdminAlbums()
  const genres = useAdminGenres()

  async function fetchAll() {
    await Promise.all([
      tracks.fetchTracks(),
      artists.fetchArtists(),
      albums.fetchAlbums(),
      genres.fetchGenres(),
    ])
  }

  return { tracks, artists, albums, genres, fetchAll }
}
```

### API Contracts

| Entity | List | Create | Update | Delete | Enrich |
|--------|------|--------|--------|--------|--------|
| Track | GET `/admin/tracks` | POST `/admin/tracks` | PATCH `/admin/tracks/:id` | DELETE `/admin/tracks/:id` | POST `/admin/tracks/:id/enrich` |
| Artist | GET `/admin/artists` | POST `/admin/artists` | PATCH `/admin/artists/:id` | DELETE `/admin/artists/:id` | POST `/admin/artists/:id/enrich` |
| Album | GET `/admin/albums` | POST `/admin/albums` | PATCH `/admin/albums/:id` | DELETE `/admin/albums/:id` | — |
| Genre | GET `/admin/genres` | POST `/admin/genres` | PATCH `/admin/genres/:id` | DELETE `/admin/genres/:id` | — |

### Payload Builders

```typescript
// Tracks have complex payload builders:
function toCreatePayload(form: TrackFormData): CreateTrackPayload {
  return {
    title: form.title,
    persian_title: form.persianTitle,
    primary_artist_id: form.primaryArtist.id,
    featured_artist_ids: form.featuredArtists.map(a => a.id),
    album_id: form.album?.id,
    duration_seconds: form.duration,
    track_number: form.trackNumber,
    disc_number: form.discNumber,
    genre_ids: form.genres.map(g => g.id),
    isrc: form.isrc,
    language: form.language,
    explicit: form.explicit,
    credits: form.credits,
  }
}

function toUpdatePayload(form, currentTrack) {
  // Only includes changed fields
  // Uses clear_fields for arrays when needed
}
```

---

## Type Safety Concerns

Inline comments in `PageAdminTracks.vue` note:

```
// TODO: AnyTrack defeats TypeScript safety — 10+ fragile getter functions
// for normalizing API responses
```

The page has 10+ getter functions like:
- `getArtistNames(track)` — joins artist names
- `getGenres(track)` — genre name list
- `getCredits(track)` — formatted credit strings
- `getSearchText(track)` — combined searchable text

These operate on `any` typed track data, bypassing TypeScript's safety.

---

## State Matrix

| State | Track List | Track Detail | Form Dialog | Delete Flow |
|-------|-----------|-------------|-------------|-------------|
| 🟢 Loading | ✅ Table skeleton | ✅ Detail skeleton | ✅ Spinner on save | N/A |
| 🟢 Loaded (has data) | ✅ Table/cards | ✅ Full detail | ✅ Form populated | N/A |
| 🟢 Loaded (no data) | ✅ "No tracks yet" + "Add Track" CTA | N/A | N/A | N/A |
| 🟢 Creating | N/A | N/A | ✅ Save spinner, fields disabled | N/A |
| 🟢 Editing | N/A | N/A | ✅ Pre-filled form + save spinner | N/A |
| 🟢 Deleting | ✅ Item dims + spinner | N/A | N/A | ✅ Dialog confirms |
| 🟢 Delete success | ✅ Item removed + toast | Redirect to list | N/A | N/A |
| 🟢 Error on create | N/A | N/A | ✅ Inline error message | N/A |
| 🟢 Error on edit | N/A | N/A | ✅ Inline error message | N/A |
| 🟢 Error on delete | ❌ Silent — item stays | ❌ Silent | N/A | ❌ Dialog stays open |
| 🔴 Enrich all (batch) | ⚠️ No progress indicator | N/A | N/A | N/A |
| 🔴 Large track list | ❌ No pagination | N/A | N/A | N/A |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-2301 | ⚠️ MAJOR | `PageAdminTracks.vue` | **No pagination on track list** — fetches all tracks at once | Extremely slow with 1000+ tracks | Add server-side pagination |
| F-2302 | ⚠️ MAJOR | `PageAdminTracks.vue` | **AnyTrack type defeats TypeScript safety** — 10+ fragile getter functions | Runtime errors from unexpected API shapes | Create proper Track type, eliminate `any` |
| F-2303 | ⚠️ MAJOR | `TrackFormDialog.vue` (1728 lines) | **Extremely large form component** — audio, metadata, lyrics, enrichment all in one | Maintainability risk, slow to render | Decompose into sub-forms (basic, credits, lyrics) |
| F-2304 | 💡 IMPROVE | All CRUD pages | **No "Enrich All" progress indicator** — fires batch POST with no feedback | Admin doesn't know if it's working | Add progress bar with item count |
| F-2305 | 💡 IMPROVE | `TrackFormDialog.vue` | **Artist/Album inline creation dispatches events** — fragile event-based coupling | Race conditions, hard to debug | Use callback props instead of events |
| F-2306 | 💡 IMPROVE | All form dialogs | **No unsaved changes warning** when closing dialog | Accidental close loses form data | Add `beforeClose` check with confirmation |
| F-2307 | 💡 IMPROVE | `PageAdminTracks.vue` | **No column customization** in table view — fixed columns | Admin may want different data visible | Add column toggle/visibility |
| F-2308 | 💡 IMPROVE | `AdminArtistsCard.vue` | **No album/track count** shown in artist grid | Can't gauge artist popularity at a glance | Add track/album count badges |
| F-2309 | 💡 IMPROVE | `AdminAlbumsCard.vue` | **No release year filter** in album grid | Hard to find albums from a specific year | Add year filter or facet |
| F-2310 | 💡 IMPROVE | `PageAdminTrackDetail.vue` | **Delete from detail page returns to list** — scroll position lost | Must find track again in list | Preserve scroll position or return to same spot |
| F-2311 | 💡 IMPROVE | All CRUD pages | **No "duplicate" action** for tracks/albums/artists | Manually re-entering similar metadata | Add "Duplicate" button to pre-fill form |
| F-2312 | 💡 IMPROVE | All CRUD pages | **No audit log** shown for edits | Can't see what changed or who changed it | Add "Last edited by X on date" metadata |

## RTL / A11y / Mobile Notes

- ✅ All form dialogs have proper `<label>` associations
- ❌ `TrackFormDialog` has no `aria-live` region for save status — screen reader doesn't announce completion
- ✅ Delete confirmation dialog has focus trap
- ❌ Table rows are not keyboard navigable — can't arrow-key through tracks
- ✅ Artist/album search Autocomplete has proper `aria-activedescendant`
- ❌ Genre color picker has no text label for screen readers
