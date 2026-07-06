# Journey 10: Playlist Lifecycle

> Full trace: create → name → add tracks → reorder → edit metadata → delete.

---

## Step-by-Step: Create a Playlist

```
User clicks "Create" button on Library page (or "+" in sidebar)
  → Dialog opens (PrimeVue Dialog, modal, 440px max-width)
  → Fields:
      [Name]            "My awesome playlist"  — required, min 1 char
      [Description]     "Optional description" — optional textarea
      [Public]          ☑ checkbox — default: checked
  → User clicks "Create"
    → API: POST /api/v1/playlists { name, description, is_public }
    → Button shows spinner (creating=true), disabled
    → On success:
      → Dialog closes
      → Form resets
      → If response has .id: router.push(/playlist/{id})
      → Else: refetch all library data
    → On error:
      → Toast: "Failed to create playlist"
```

### Create Entry Points

| Entry Point | Location | UX |
|-------------|----------|-----|
| Library playlists tab | "Create" button in header | Opens inline dialog |
| Library playlists empty state | "Create Playlist" CTA button | Same |
| Sidebar (MusicSidebar) | "+" icon | Opens same dialog or navigates to library |
| User page | "Create Playlist" | Same |

### Steps to First Playable Playlist

```
1. Click Create (1 click)
2. Type name (variable)
3. Click Create (1 click)
4. API call (~200-500ms)
5. Redirect to playlist detail page (~500ms load)
6. Now on empty playlist — must add tracks
7. Find tracks to add (navigate to search/catalog)
8. Select tracks — 1 click per track via "Add to Playlist"
```

**Minimum clicks to first track in a new playlist: 4 clicks + typing + 2 page loads**

---

## Adding Tracks to a Playlist

### Entry Points for "Add to Playlist"

| Surface | Trigger | UX |
|---------|---------|-----|
| TrackRow context menu | Right-click → "Add to Playlist" | Opens `AddToPlaylistDialog` |
| NowPlayingBar | Overflow → "Add to Playlist" | Same dialog |
| FullscreenPlayer | Overflow → "Add to Playlist" | Same dialog |
| Track detail page | "Add to Playlist" button | Same dialog |
| Playlist detail page | "Add tracks" search interface | Inline search + add |

### AddToPlaylistDialog

When triggered from a track:
```
Dialog appears with:
  - Title: "Add to Playlist"
  - Search bar: "Find a playlist..."
  - List of user's playlists (fetched via getMyPlaylists())
  - Each playlist shows: name + track count
  - Click playlist → API: POST /playlist/:id/tracks { track_id }
  - Checkmark appears briefly on clicked playlist
  - Selected count updates: "Added to N playlists"
  - "New Playlist" button at bottom to create on-the-fly
```

**No duplicate check**: The dialog doesn't warn if a track is already in the target playlist. It silently adds a duplicate entry.

### Playlist Detail "Add Tracks" Search

```
On playlist detail page:
  Button: "Add tracks"
    → Shows inline search bar
    → User types query → debounced search → results appear
    → Click "+" on a result → API: POST /playlist/:id/tracks { track_id }
    → Track appears in playlist list
    → Input stays open for more additions
```

### Adding from Search Results / Track Page

```
On any track (search result, album track list, etc.):
  Three-dot menu → "Add to Playlist"
    → Opens AddToPlaylistDialog
    → User selects playlist
    → Toast: "Added to playlist"
```

---

## Reordering Playlist Tracks

### Drag-and-Drop Reorder

On the playlist detail page, tracks are rendered in a `TrackRow` list with drag handles. Drag-and-drop uses `vuedraggable` (based on SortableJS):

```
User grabs drag handle → drags track up/down
  → Visual feedback: ghost element, drop zone highlight
  → On drop:
    → API: PUT /playlist/:id/reorder { track_ids: [new order] }
    → Track order updates in UI
    → No confirmation — immediate
```

**No "undo"** on reorder. Once dropped, the new order is saved.

### When Playing from Playlist

When a user plays a track from a playlist, the entire playlist is used as the queue (see J05). Shuffle/reorder in the player is independent of the playlist's stored order.

---

## Editing Playlist Metadata

| Field | Edit Location | UX |
|-------|--------------|-----|
| Name | Playlist detail page — pencil icon near title | Inline edit? Or dialog? Click → text input appears |
| Description | Same area | Inline textarea |
| Cover image | Hover overlay on cover → "Change cover" / "Add cover" | File picker → upload via `mediaApi` → API: `PUT /playlist/:id` |
| Collaborative toggle | Toggle button near header | API: `PUT /playlist/:id { collaborative: true/false }` |
| Public/Private | N/A — set at creation, not editable in current UI | Not visible on detail page |

### Cover Image Flow

```
User hovers playlist cover → sees overlay: camera icon + "Change cover" / "Add cover"
  → Click → hidden <input type="file" accept="image/*"> opens
  → Select file → uploadCover() called
    → Uploads via mediaApi.uploadImage()
    → Updates playlist via playlistsApi.updatePlaylist(id, { cover_url })
    → UI refreshes with new cover
  → No crop/resize step
  → Supported file types: image/* (any image MIME type — no validation)
```

---

## Deleting a Playlist

On the playlist detail page, if the user is the owner:

```
User clicks "Delete" (usually in overflow menu or header actions)
  → Confirmation dialog:
    ┌──────────────────────────────────────┐
    │  Delete playlist "[name]"?            │
    │                                      │
    │  This will remove it from your        │
    │  library. This can't be undone.       │
    │                                      │
    │  [  Cancel  ]    [  Delete  ]         │
    └──────────────────────────────────────┘
  → User clicks "Delete"
    → API: DELETE /api/v1/playlists/:id
    → Toast: "Playlist deleted"
    → Router pushes to /library
  → If user clicks "Cancel" → dialog closes, no change
```

**Destructive action has confirmation** ✅ — one of the few places with proper safety.

---

## Duplicate Track Handling

| Scenario | Behavior |
|----------|----------|
| Add same track twice from dialog | ✅ No warning — silently creates duplicate entry |
| Add same track via drag from search | ✅ Same — no check |
| Remove one instance of a duplicate | ✅ Removes only that instance (specific track entry ID) |
| Remove all instances | Must remove each one individually |

**No duplicate detection in the UI**. The backend may or may not enforce uniqueness — the frontend doesn't check.

---

## State Matrix Findings

| State | Create Flow | Add Track | Reorder | Edit | Delete |
|-------|------------|-----------|---------|------|--------|
| 🟢 Loading | ✅ Button spinner | ✅ Search loading | ✅ Ghost | ✅ Upload spinner | ✅ N/A |
| 🟢 Success | ✅ Redirect + toast | ✅ Toast "Added" | ✅ Instant | ✅ New cover visible | ✅ Toast + redirect |
| 🟢 Error | ✅ Toast + form stays | ❌ Silent catch | ❌ Order reverts? | ✅ Toast | ✅ Dialog closes |
| 🔴 Empty playlist | ✅ "No tracks — add some" | N/A | N/A | N/A | N/A |
| 🔴 Duplicate detection | N/A | ❌ No warning | N/A | N/A | N/A |
| 🔴 Undo any operation | ❌ No undo | ❌ No undo | ❌ No undo | ❌ No undo | ✅ Confirmed delete only |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1001 | ⚠️ MAJOR | `AddToPlaylistDialog.vue` | No warning when adding a duplicate track | Users may accidentally duplicate tracks | Check playlist contents before add |
| F-1002 | ⚠️ MAJOR | Playlist detail page | No "undo" snackbar after reorder or adding/removing tracks | Mistakes are permanent | Add undo toast with 5s timeout |
| F-1003 | 💡 IMPROVE | Create flow | Create → redirect → add tracks is 3+ steps + page loads | High friction for first playlist | Allow adding tracks directly in the create flow |
| F-1004 | 💡 IMPROVE | Playlist detail page | No filter/search within a playlist's tracks | Long playlists hard to navigate | Add inline search within playlist |
| F-1005 | 💡 IMPROVE | `PagePlaylistDetail.vue` | No batch select for remove/move operations | Cleaning up a playlist is one-by-one | Add multi-select mode |
| F-1006 | 💡 IMPROVE | `PagePlaylistDetail.vue` | Cover upload has no crop/position tools | Cover may display poorly | Add simple crop step |
| F-1007 | 💡 IMPROVE | Playlist detail page | No "Make Public/Private" toggle in detail view (only at creation) | Can't change visibility after creation | Add visibility toggle |
| F-1008 | 💡 IMPROVE | Playlist detail page | No playlist length/size indicator beyond track count | Users don't know total duration | Add "X hr Y min" total duration |
| F-1009 | 💡 IMPROVE | `PagePlaylistDetail.vue` | Delete confirmation is text-only — no playlist name shown in title? | User might delete wrong playlist | Show playlist name in confirmation |

## RTL / A11y / Mobile Notes

- ✅ Create dialog is modal with focus trap
- ✅ Delete has confirmation (rare in this app — note as good pattern)
- ✅ Cover upload has hover overlay (visible to keyboard users via focus styles)
- ❌ Drag reorder not keyboard accessible (vuedraggable mouse/touch only)
- ❌ No `aria-label` on track remove buttons in playlist
- ✅ Toast notifications for add/delete actions
- ❌ No "undo" on destructive actions (only confirmation before delete, but no undo after)
