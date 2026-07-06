# Journey 11: Collaborative Playlists

> Full trace: enable collaboration → invite → real-time sync → permissions → conflict handling.

---

## Enabling Collaboration

On a playlist you own, the detail page shows a "Collaborative" toggle:

```
Playlist info header:
  ┌─────────────────────────────────────────┐
  │  PLAYLIST            ⏺ Collaborative    │
  │  My Road Trip Mix                       │
  │  42 tracks • Updated 2h ago             │
  │                                         │
  │  [☐ Collaborative]  [Share]             │
  └─────────────────────────────────────────┘
```

| Action | API | Effect |
|--------|-----|--------|
| Toggle ON | `PUT /playlists/:id { collaborative: true }` | Playlist becomes editable by anyone with the link |
| Toggle OFF | `PUT /playlists/:id { collaborative: false }` | Only owner can edit |

**No confirmation** when toggling collaborative mode on/off. One click, immediate effect.

---

## Inviting Collaborators

### Share Link Mechanism

The primary sharing mechanism is via `useSocialShare().copyLink()`:

```
User clicks "Share" → copies playlist URL to clipboard
  → Toast: "Link copied to clipboard"
  → User sends link to friends (via WhatsApp, Telegram, etc.)
  → Recipient opens link → sees playlist page
  → If collaborative is ON: "Add track" button is visible to all viewers
  → If collaborative is OFF: viewers see read-only view
```

### Explicit Collaborator Management

The playlist detail page shows collaborators:
```
Collaborators section:
  ┌──────────────────────────────────────┐
  │  Collaborators                        │
  │                                      │
  │  [User avatar]  Alice (Owner)        │
  │  [User avatar]  Bob                  │ [X]
  │  [User avatar]  Charlie              │ [X]
  │                                      │
  │  [ + Add Collaborator ]              │
  └──────────────────────────────────────┘
```

| Action | API | UX |
|--------|-----|-----|
| Add by username/ID | `POST /playlists/:id/collaborators { user_id }` | Text input + search → select user → added |
| Remove collaborator | `DELETE /playlists/:id/collaborators/:userId` | Click X → removed immediately |
| Leave (as collaborator) | Same API (self) | Button: "Leave playlist" |

**No confirmation** on remove collaborator — immediate removal with no undo. The removed user is not notified.

---

## Real-Time Sync Architecture

### WebSocket Integration

```typescript
// From useCollaborativePlaylist.ts
wsClient.connect()
wsClient.subscribe(`playlist:${playlistId}`)

// Events listened to:
wsClient.on('playlist.track_added', (msg) => handleTrackAdded(msg.payload))
wsClient.on('playlist.track_removed', (msg) => handleTrackRemoved(msg.payload))
wsClient.on('playlist.track_reordered', (msg) => handleTrackReordered(msg.payload))
wsClient.on('playlist.updated', (msg) => handlePlaylistUpdated(msg.payload))
```

### What Users See in Real Time

| Event | Sender Sees | Other Collaborators See |
|-------|-------------|------------------------|
| Track added | Track appears in list + toast | **Playlist auto-updates** — track appears in list without refresh. No visible notification of WHO added it |
| Track removed | Track disappears from list | **Track disappears** — no notification or undo banner |
| Track reordered | Tracks in new order | **Order updates silently** — no indication of change |
| Playlist metadata updated | New metadata shown | **Metadata updates silently** — no toast |

### What's Missing in the Sync UX

| Missing Feature | Impact |
|----------------|--------|
| **Presence indicators**: No "Alice is editing" or "Bob is currently viewing" | No awareness of others in the playlist |
| **Who added/removed**: No user attribution on events | Can't tell who changed the playlist |
| **Notification to removed user**: None | Can be removed without knowing |
| **Conflict handling**: No guard if two people reorder simultaneously | Last write wins — silent overwrite |
| **Cursor/selection sync**: No awareness of what others are selecting | Can accidentally disrupt someone's reorder |

### Disconnect Behavior

```
WebSocket disconnects → user experiences:
  • No real-time updates (playlist appears static)
  • Own changes are still applied via REST API
  • On reconnect → missing events are NOT replayed
  • User must refresh page to see full state
```

**No event replay on WebSocket reconnect**. If a collaborator adds 5 tracks while your connection is down, you won't see them until you refresh.

---

## Permissions Model

| Action | Owner | Collaborator | Viewer |
|--------|-------|-------------|--------|
| View playlist | ✅ | ✅ | ✅ |
| Add tracks | ✅ | ✅ | ❌ |
| Remove tracks | ✅ | ✅ | ❌ |
| Reorder tracks | ✅ | ✅ | ❌ |
| Edit name/description | ✅ | ❌ | ❌ |
| Change cover | ✅ | ❌ | ❌ |
| Toggle collaborative | ✅ | ❌ | ❌ |
| Add/remove collaborators | ✅ | ❌ | ❌ |
| Delete playlist | ✅ | ❌ | ❌ |

**No granular permissions**: All collaborators can add/remove/reorder any track. No way to restrict "add only" or "no remove".

---

## State Matrix Findings

| State | Present? | Notes |
|-------|----------|-------|
| 🟢 Enable collaborative | ✅ One-click toggle | |
| 🟢 Add collaborator | ✅ Search + add | |
| 🟢 Remove collaborator | ✅ X button | |
| 🔴 Presence indicators | ❌ Not implemented | |
| 🔴 Who-added attribution | ❌ Not shown in sync | |
| 🔴 Event replay on reconnect | ❌ Missing events lost | |
| 🔴 Conflict resolution | ❌ Last write wins | |
| 🔴 Undo on collaborator removal | ❌ Immediate, no undo | |
| 🔴 Notify removed user | ❌ No notification | |
| 🔴 Granular permissions | ❌ All-or-nothing | |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1101 | 🚨 BLOCKER | `useCollaborativePlaylist.ts` | **No event replay on WebSocket reconnect** — missing events are lost permanently | Users returning from offline/disconnect miss all changes made by others | Replay events on reconnect or re-fetch playlist |
| F-1102 | ⚠️ MAJOR | `useCollaborativePlaylist.ts` | No presence indicators — users don't know if others are editing | Multiple people may edit simultaneously causing conflicts | Add "N people editing" badge |
| F-1103 | ⚠️ MAJOR | WebSocket events | `track_added` and `track_removed` events include `user_id` but UI doesn't display it | Can't tell who added/removed tracks | Show attribution: "Added by Alice" |
| F-1104 | ⚠️ MAJOR | Collaborative | No conflict handling for simultaneous reorder | Last write silently overwrites the other person's arrangement | Implement operational transform or locking |
| F-1105 | 💡 IMPROVE | Collaborate toggle | No confirmation when enabling/disabling collaborative mode | Accidental toggle makes playlist public/private | Add confirmation dialog |
| F-1106 | 💡 IMPROVE | Collaborator removal | No undo or confirmation on removing a collaborator | Might accidentally remove someone | Add confirmation + undo |
| F-1107 | 💡 IMPROVE | Invite flow | Only share-by-link — no in-app notification/invite | Friends don't know they've been added to a playlist | Send in-app notification + email |
| F-1108 | 💡 IMPROVE | Collaborate toggle | No granular permissions (add-only, view-only) | All collaborators can delete any track | Add "can edit" vs "can add" roles |
| F-1109 | 💡 IMPROVE | UI | No visual indicator on playlist cards that a playlist is collaborative | Users don't know which playlists are shared | Add "Collaborative" badge to playlist grid items |

## RTL / A11y / Mobile Notes

- ✅ Collaborator list with avatar + name is readable
- ❌ No keyboard accessibility for drag-reorder (vuedraggable limitation)
- ✅ WebSocket events auto-update UI without page refresh
- ❌ Screen readers are not notified when remote tracks are added (no `aria-live` on playlist list)
