# Journey 06: Player Deep Dive

> Exhaustive trace of every player control, every state transition, every UI surface.
> Covers: NowPlayingBar, ExpandedPlayer, FullscreenPlayer, MobileBottomSheet, FloatingMiniPlayer, PiP.

---

## Player Surface Map

| Surface | Screen Type | Trigger | Key Controls |
|---------|-------------|---------|--------------|
| **NowPlayingBar** | Desktop + Tablet | Always visible (minimized) | Play/pause, prev/next, progress bar, volume, like, add-to-playlist, overflow menu, expand, queue toggle |
| **NowPlayingBar (expanded)** | Desktop | Click "now playing" area | Full controls: shuffle, repeat, queue preview, lyrics toggle |
| **NowPlayingBar (mobile)** | Mobile | Always visible (compact) | Cover, title, play/pause, next — tap to open MobileBottomSheet |
| **FullscreenPlayer** | All | Click expand / keyboard F | Tabs: Now Playing, Lyrics, Queue. Full controls + visualizer + PiP |
| **ExpandedPlayer** | All | Same (behind `redesignedPlayer` flag) | Same structure but redesigned layout |
| **MobileBottomSheet** | Mobile | Tap NowPlayingBar on mobile | Cover art, seek bar, controls, volume, like |
| **FloatingMiniPlayer** | Desktop | Overflow menu → "Open in Mini Player" | Draggable window, mini/expanded modes, PiP toggle |
| **PiPPlayerContent** | Desktop | PiP button / auto on tab switch | Compact inline-styled player, always-on-top |
| **KeyboardShortcuts** | All | Press `?` | Overlay listing all shortcuts |
| **RadioMode** | All | Radio playback | Full-screen radio with like/dislike, skip, seed info |
| **TheatreMode** | All | Theatre mode toggle | Large album art + karaoke lyrics side by side |
| **AmbientMode** | All | Ambient mode | Canvas particle system, minimal controls |

---

## Step-by-Step: NowPlayingBar (Desktop)

### Layout
```
┌──────────────────────────────────────────────────────────────┐
│ ████████████████████████████████████████████████████████████  │ ← Progress bar (clickable, draggable)
│ ┌────┐  ┌───────┐            ┌──┐ ┌──┐ ┌──┐ ┌──┐ ┌──┐ ┌──┐ │
│ │Art │  │Title  │            │♡ │ │⏮ │ │▶ │ │⏭ │ │🔊│ │⋯ │ │
│ │    │  │Artist │            │   │ │  │ │  │ │  │ │  │ │  │ │
│ └────┘  └───────┘            └──┘ └──┘ └──┘ └──┘ └──┘ └──┘ │
└──────────────────────────────────────────────────────────────┘
```

### Controls Detail

| Control | Interaction | Behavior | Status |
|---------|-------------|----------|--------|
| **Progress bar** | Click anywhere / drag thumb | `store.seek(seconds)` → `engine.seek()` → audio `currentTime` | ✅ Works |
| **Cover art** | Click | Toggles fullscreen player | ✅ |
| **Title/Artist** | Click | Toggles fullscreen player | ✅ |
| **Like button (♡)** | Click | `useTrackLike().toggleLike()` → API call | ✅ |
| **Add to playlist** | Click | Opens `AddToPlaylistDialog` | ✅ |
| **Previous (⏮)** | Click | `player.playPrevious()` — if >4s into track, seeks to 0. Else goes to previous track | ✅ |
| **Play/Pause (▶/⏸)** | Click | Toggle via `player.resume()` / `player.pause()` | ✅ |
| **Next (⏭)** | Click | `player.playNext()` → `engine.next()` | ✅ |
| **Volume (🔊)** | Click icon to mute/unmute. Drag slider to adjust. | `store.toggleMute()` / `store.setVolume(v)` | ✅ Volume slider hardcoded LTR (F-025) |
| **Overflow menu (⋯)** | Click | Opens `PlayerOverflowMenu`: Audio quality, Sleep timer, Crossfade slider, Open in Mini Player | ✅ |

### Shuffle & Repeat (Expanded State)

Clicking shuffle or repeat icons opens a **popup selector**:

**Shuffle modes** (click to cycle or open popup):
| Mode | Behavior | Effect on next() |
|------|----------|-----------------|
| `off` | Sequential through queue | `next()` → `nextSequential()` |
| `queue` | Fisher-Yates shuffled queue | `next()` → `nextShuffleQueue()` |
| `catalog` | Fetch random from API | `next()` → `nextCatalog()` — fetches 20 random tracks |
| `similar` | Fetch similar to current | `next()` → `nextSimilar()` — fetches 10 similar tracks |

**Repeat modes** (click to cycle: off → all → one → off):
| Mode | End-of-track | End-of-queue |
|------|-------------|--------------|
| `off` | Stop after track ends | Stop |
| `one` | Seek to 0 and replay | N/A |
| `all` | Continue to next | Wrap to queue[0] |

### Overflow Menu Detail

| Item | Action | UX |
|------|--------|----|
| Audio Quality | Cycles: auto → low → medium → high → lossless | Stored in `player.audioQuality`, persisted to localStorage |
| Sleep Timer | Presets: None / 5 / 15 / 30 / 45 / 60 min | `setSleepTimer(minutes)` → pauses after N min |
| Crossfade | Slider 0-12s | Stored but **not actually implemented** in audio engine (no crossfade logic) |
| Open in Mini Player | Toggles PiP / FloatingMiniPlayer | Calls `usePlayerPiPController().toggle()` |

---

## Step-by-Step: FullscreenPlayer / ExpandedPlayer

### Opening

| Trigger | Action |
|---------|--------|
| Click cover art in NowPlayingBar | Sets `fullscreenOpen = true` |
| Keyboard `F` | Toggles `fullscreenOpen` |
| Keyboard `Ctrl/Cmd+L` | Opens fullscreen to lyrics tab |
| Keyboard `Q` | Opens queue panel (not fullscreen) |

### Tabs

**Tab 1: Now Playing**
```
┌──────────────────────────────────────────┐
│         ┌──────────────┐                 │
│         │  Album Art   │ ← Large, centered│
│         │  (pulse ring)│ ← Glowing border │
│         └──────────────┘                 │
│         Track Title                       │
│         Artist Name                       │
│                                          │
│   ────●──────────────────────────────    │ ← Seek bar with hover preview
│   1:23                          3:45     │
│                                          │
│  ┌──┐ ┌──┐ ┌──┐ ┌────┐ ┌──┐ ┌──┐ ┌──┐ │
│  │🔀│ │⏮│ │▶⏸│ │⏭│ │🔁│ │🔊│ │⋯ │ │
│  └──┘ └──┘ └──┘ └────┘ └──┘ └──┘ └──┘ │
│                                          │
│  🎤 Lyrics tab    📋 Queue tab           │
└──────────────────────────────────────────┘
```

**Tab 2: Lyrics** (via SyncedLyrics / KaraokeLyrics / LyricsDisplay)
- Auto-scrolling timed LRC lines
- Click any line → seek to that timestamp
- Word-by-word karaoke highlighting (KaraokeLyrics)
- Particle canvas background for karaoke mode

**Tab 3: Queue** (via QueuePanel)
- Shows "Now Playing" track
- Draggable "Up Next" list (via `vuedraggable`)
- Click any track → play it next
- Remove button on each track

### ExpandedPlayer vs FullscreenPlayer

| Feature | FullscreenPlayer (legacy) | ExpandedPlayer (redesigned) |
|---------|--------------------------|----------------------------|
| Background | Blurred cover art + dynamic color | Aurora gradient |
| Visualizer | ❌ Not present | ✅ VisualizerSystem component |
| Crossfade toggle | ❌ Not present | ✅ `cycleCrossfade()` button |
| Playback speed | ❌ Not present | ✅ `cycleSpeed()` button (0.5→0.75→1→1.25→1.5→2) |
| Layout | Teleported to `<body>` | Teleported to `<body>` |
| PiP | Auto-open on `visibilitychange` when fullscreen+playing | PiP button in controls |
| Lines | 1016 lines | 960 lines |

### VisualizerSystem (`VisualizerSystem.vue`)

| Mode | Description |
|------|-------------|
| `spectrum` | Vertical frequency bars (default for low-tier devices) |
| `waveform` | Horizontal waveform line |
| `circular` | Circular spectrum with center art |
| `radialBars` | Bars radiating from center |
| `particle` | Particle system reacting to audio |
| `fluid` | Fluid/organic simulation |

- Gets `AnalyserNode` from `audioEngine.getAnalyserNode()`
- Canvas render via `requestAnimationFrame` loop
- Respects `prefers-reduced-motion`
- Auto-fallback to `spectrum` on low-performance devices

### Picture-in-Picture (PiP)

**Architecture**:
- Uses Document Picture-in-Picture API (`window.documentPictureInPicture`)
- Creates a **mini Vue app** inside the PiP window, sharing the main Pinia instance via `window.__PINIA__`
- `PiPPlayerContent.vue` — fully inline-styled (no class dependencies)

**Triggers**:
| Trigger | From | Behavior |
|---------|------|----------|
| PiP button | FullscreenPlayer | Opens PiP window |
| PiP button | ExpandedPlayer | Opens PiP window |
| "Open in Mini Player" | Overflow menu | Opens FloatingMiniPlayer (not PiP — separate feature) |
| Auto-PiP | FullscreenPlayer on `visibilitychange` | If fullscreen is open and playing, auto-opens PiP on tab switch |

**Limitations**:
- `window.documentPictureInPicture` is Chrome-only (no Firefox/Safari)
- No fallback to "audio only" PiP for unsupported browsers
- PiP window dimensions are hardcoded in `useDocumentPictureInPicture.ts`

---

## Media Session Integration

```typescript
// media-session.ts
updateMediaSession(track, {
  play: () => player.resume(),
  pause: () => player.pause(),
  nexttrack: () => player.playNext(),
  previoustrack: () => player.playPrevious(),
  seekto: (details) => player.seek(details.seekTime),
})
```

| Feature | Supported? | Notes |
|---------|-----------|-------|
| Lock screen controls | ✅ Play/pause, prev/next, seek |
| Lock screen artwork | ✅ 512x512 JPEG from `track.coverUrl` |
| Playback state | ✅ `playing` / `paused` / `none` |
| Seek on lock screen | ✅ Via `seekto` action handler |
| `skipad` handler | ❌ Not registered (no ads implemented) |
| `stop` handler | ❌ Not registered |

---

## Keyboard Shortcuts

| Key | Modifier | Action | Source |
|-----|----------|--------|--------|
| `Space` | — | Toggle play/pause | `useKeyboardShortcuts.ts` |
| `N` | — | Next track | Same |
| `P` | — | Previous track | Same |
| `ArrowRight` | — | Seek +5s | Same |
| `ArrowLeft` | — | Seek -5s | Same |
| `ArrowUp` | — | Volume +0.1 | Same |
| `ArrowDown` | — | Volume -0.1 | Same |
| `M` | — | Toggle mute | Same |
| `S` | — | Cycle shuffle modes | Same |
| `R` | — | Toggle repeat mode | Same |
| `L` | — | Like/unlike current track | Same |
| `F` | — | Toggle fullscreen player | Same |
| `Q` | — | Toggle queue panel | Same |
| `?` | — | Toggle shortcuts help | Same |
| `Ctrl/Cmd+L` | Ctrl/Cmd | Toggle lyrics | Same |
| `Escape` | — | Close bottom sheet | Component-specific |

**Guard**: All shortcuts are ignored when focus is in `INPUT`, `TEXTAREA`, `SELECT`, or `contentEditable`.

---

## State Transitions

### Play → Pause → Resume
```
play → isPlaying=true, isBuffering=false, engine.playstate='playing'
pause → isPlaying=false, engine.playstate='paused'
resume → isPlaying=true, engine.playstate='playing'
```

### Track End → Next
```
audio 'ended' event
  → if repeatMode='one': seek(0) + play (handled in engine)
  → else: engine.next()
    → nextSequential(): queue.next() → new track or null
    → nextShuffleQueue(): advance shuffle position → new track
    → nextCatalog(): fetch random track from API
    → nextSimilar(): fetch similar tracks from API
  → if no next track:
    → if repeatMode='all': wrap to queue[0]
    → else: stop playback
```

### Audio Error Mid-Play
```
audio 'error' event → engine emits 'error'
  → store.error = message
  → isBuffering = false, isPlaying = false
  → consecutiveFailures++
  → if consecutiveFailures < 3: engine.next() (skip)
  → if consecutiveFailures >= 3: playbackStopped = true, engine.stop()
```

### Volume Change
```
Volume slider drag → store.setVolume(value)
  → engine.setVolume(value) → audio.volume = value
  → audio 'volumechange' → engine re-emits → store syncs volume, muted
  → persisted to localStorage ('player-volume', 'player-muted')
```

---

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-601 | ⚠️ MAJOR | `FullscreenPlayer.vue` | FullscreenPlayer (1016 lines) and ExpandedPlayer (960 lines) are both extremely large (F-021 continuation) | Maintainability risk, slow re-renders | Decompose into sub-components |
| F-602 | ⚠️ MAJOR | `player-engine.ts` 'error' handler | Consecutive failure guard skips to next track automatically — no user choice | User may want to retry the current track, not skip | Show retry option before auto-skip |
| F-603 | ⚠️ MAJOR | `ExpandedPlayer.vue` | Crossfade toggle exists but crossfade is **not implemented** in audio engine | Setting has no effect — misleading UX | Implement crossfade or remove toggle |
| F-604 | ⚠️ MAJOR | `NowPlayingBar.vue` | Volume slider hardcoded `dir="ltr"` (F-025) | RTL users see reversed slider direction | Use logical CSS properties |
| F-605 | 💡 IMPROVE | `NowPlayingBar.vue` | Shuffle/repeat popups close on click — can't compare options easily | User must open popup, read, close, open other | Show both modes simultaneously or with tooltip |
| F-606 | 💡 IMPROVE | `PlayerOverflowMenu.vue` | Crossfade slider exists but feature not implemented (F-603) | User can set crossfade but nothing happens | Implement or hide |
| F-607 | 💡 IMPROVE | `FloatingMiniPlayerContent.vue` | Floating mini player has no close/minimize-to-bar button | Once opened, user must toggle off from overflow menu | Add minimize button |
| F-608 | 💡 IMPROVE | `media-session.ts` | Media Session `stop` action handler not registered | Lock screen may not clear player UI on stop | Register stop handler |
| F-609 | 💡 IMPROVE | `FullscreenPlayer.vue` | Auto-PiP on `visibilitychange` is intrusive — user may not expect it | Surprising behavior when switching tabs | Opt-in toggle or confirmation |
| F-610 | 💡 IMPROVE | `useKeyboardShortcuts.ts` | Shortcut help (`?`) shows all shortcuts but no way to customize or disable them | Power users may accidentally trigger shortcuts | Add shortcut disable in settings |
| F-611 | 💡 IMPROVE | `PiPPlayerContent.vue` | PiP uses Chrome-only API — no fallback for Firefox/Safari | Non-Chrome users can't use PiP | Add mini-player fallback |
| F-612 | ✨ DELIGHT | `VisualizerSystem.vue` | 6 visualizer modes but no mode-cycling button in the player UI | Users can't try different visualizers without knowing the modes exist | Add visualizer mode toggle to fullscreen player |

## RTL / A11y / Mobile Notes

- ✅ `aria-label` on all player controls (play, pause, next, prev, volume, shuffle, repeat)
- ✅ Touch targets ≥ 44px on all player buttons
- ✅ Media Session provides lock screen controls
- ✅ PiP window has minimal but functional controls
- ❌ Volume slider is hardcoded LTR (F-025)
- ❌ No `aria-valuetext` on seek bar — screen reader just hears "50%"
- ❌ Visualizer canvas has no `aria-label` describing it
- ❌ Progress bar seek: keyboard arrow keys work (from global shortcuts) but no visible focus indicator on the seek bar itself
- ✅ `prefers-reduced-motion` respected in VisualizerSystem and CSS animations
- ✅ Mobile bottom sheet has touch-drag close gesture
