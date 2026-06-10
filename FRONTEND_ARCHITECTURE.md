# Frontend Architecture & Product Plan

## Vision
A premium, immersive music streaming platform combining the best of Spotify, Apple Music, and SoundCloud with its own identity.

---

## 1. Information Architecture

### Top Navigation
```
Auth (Login/Register/Settings) → Main App
├── Home           → Discovery, Continue Listening, Made For You
├── Search         → Instant overlay, keyboard nav, typed results
├── Library        → Tracks, Albums, Artists, Playlists, History
├── Playlists      → All playlists CRUD, collaborative
├── Recommendations→ For You, Trending, Popular, Recent, Daily Mix
├── Notifications  → Activity feed, system messages
└── Profile        → User profile, followers, following, stats
```

### Page Hierarchy
```
/app                          # Authenticated app shell
  /                           # Home (discovery)
  /search                     # Full search page (fallback)
  /library                    # User library (tabs: tracks|albums|artists|playlists)
  /playlists                  # All playlists
  /playlists/:id              # Single playlist detail
  /artist/:id                 # Artist profile
  /album/:id                  # Album detail
  /track/:id                  # Track detail with lyrics
  /recently-played            # Listening history
  /discover                   # Discovery hub
  /recommendations            # Recs hub (parent)
    /popular                  # Trending
    /best                     # Top rated
    /recent                   # New releases
    /for-you                  # Personalized
  /notifications              # Notification center
  /profile                    # My profile
  /profile/:id                # Other user profile
  /creator                    # Creator dashboard (if artist)
  /ai/playlist-generator      # AI-generated playlists
  /ai/mood-explorer           # Mood-based discovery
```

---

## 2. Design System (Already Built)

### Colors
- `#1db954` - Spotify green (primary accent)
- `#121212` - Base dark surface
- `#0a0a0a` - Deeper dark (player bar)
- `#1a1a2e` - Elevated surfaces
- Slate grays for text hierarchy (50-500)

### Typography
- Primary: IRANYekanWeb (Persian, also works with Latin)
- Font weights: Light(300), Regular(400), Bold(700), Black(900)
- Base: 14px on body

### Spacing
- 4px base unit, scales: 1, 2, 3, 4, 5, 6, 8, 10, 12, 16, 20, 24
- Tailwind utility classes throughout

### Components
- PrimeVue component library (auto-imported)
- Custom: Player, Queue, Lyrics, Carousels, TrackRows, Cards
- Dark theme only (`.app-dark`), light theme structure ready

---

## 3. Route Structure (Current + Planned)

```
# Current routes (all working):
/auth/login          → LayoutAuth
/auth/register       → LayoutAuth
/                    → LayoutMusicApp (home)
/search              → LayoutMusicApp
/library             → LayoutMusicApp
/playlists           → LayoutMusicApp
/playlists/:id       → LayoutMusicApp
/recently-played     → LayoutMusicApp
/track/:id           → LayoutMusicApp
/album/:id           → LayoutMusicApp
/artist/:id          → LayoutMusicApp
/discover            → LayoutMusicApp
/recommendations     → LayoutMusicApp
/notifications       → LayoutMusicApp
/profile             → LayoutMusicApp
/profile/:id         → LayoutMusicApp
/creator             → LayoutMusicApp
/ai/playlist-generator  → LayoutMusicApp
/ai/mood-explorer    → LayoutMusicApp
/admin/**            → LayoutAdmin
/404                 → LayoutEmpty

# Planned routes (composables/services exist, need pages):
/recommendations/for-you → LayoutMusicApp (exists as tab, needs dedicated)
/library/tracks     → LayoutMusicApp (tabs exist, needs routing)
/library/albums     → LayoutMusicApp
/library/artists    → LayoutMusicApp
```

---

## 4. Component Architecture (Current + Gaps)

### Built & Working
```
LayoutMusicApp         # App shell (sidebar + topbar + player)
├── MusicSidebar       # Desktop nav
├── MusicTopbar        # Sticky header with search/notifications
├── SearchOverlay      # NEW: Instant search overlay (Ctrl+K)
├── NowPlayingBar      # Bottom player bar
├── QueuePanel         # Queue drawer
├── FullscreenPlayer   # NOT YET BUILT
└── LyricsDisplay      # Synced lyrics viewer

Music Components:
├── HomeSection        # Section header with eyebrow
├── HomeCarousel       # Horizontal scrollable carousel
├── TrackRow           # Single track row (playable)
├── TrackList          # Track list container
├── ArtistHero         # Artist page header
├── ArtistCard         # Artist card for carousels
├── AlbumCard          # Album card for carousels
├── ActivityItem       # Social activity feed item
├── NotificationItem   # Notification card
├── ReactionButton     # Like/love/dislike
├── UserHero           # User profile hero
└── ProfileTabs        # Profile section tabs
```

### Gaps to Fill (Priority Order)

#### P0 - Must Have
- **FullscreenPlayer** - `/components/music/FullscreenPlayer.vue`
  - Large cover art with vinyl animation
  - Full controls (prev, play/pause, next, shuffle, repeat)
  - Seek bar with time
  - Volume control
  - Queue toggle
  - Lyrics mode toggle
  - Like/share buttons
  - Background blur from album art

- **LyricsDisplay Enhancement** - Real-time synced lyrics
  - Current line highlight with green accent
  - Auto-scroll
  - Karaoke mode (word-by-word highlight)
  - Fullscreen mode
  - Smooth scroll animations

#### P1 - Engagement Drivers
- **LibraryPage Enhancement** - `/pages/app/PageLibrary.vue`
  - Tabbed: Liked Songs | Albums | Artists | Playlists
  - Grid/list toggle
  - Sort by recent, title, artist
  - Search within library
  - Pinned items

- **AlbumPage Enhancement** - `/pages/app/PageAlbum.vue`
  - Credits section (writers, producers)
  - Featured artists
  - Play count
  - Like/unlike album
  - Add to playlist
  - Share
  - Related albums carousel

- **TrackPage Enhancement** - `/pages/app/PageTrack.vue`
  - Full lyrics display
  - Credits section
  - Related tracks
  - Play count
  - Add to playlist/share

#### P2 - Social & Community
- **Social Features**
  - Share track/album/playlist (copy link)
  - Activity feed on profile
  - Following/followers lists
  - Listening activity toggle

- **Collaborative Playlists**
  - Add/remove collaborators
  - Real-time updates via WebSocket
  - Drag-and-drop track reordering

#### P3 - Advanced
- **Recommendations Hub Overhaul**
  - Daily Mix cards
  - Discovery Weekly
  - Genre-based mixes
  - Mood-based playlists

- **Years Wrapped**
  - Top tracks, artists, genres
  - Listening stats
  - Shareable cards

---

## 5. State Architecture

### Pinia Stores (Current)
```
user-auth.ts    # Auth tokens, user profile, admin check
  ├── login/logout/me/refresh
  ├── device ID persistence
  └── admin role detection

player.ts       # Audio playback engine bridge
  ├── currentTrack, queue, isPlaying
  ├── shuffle/repeat modes
  ├── currentTime, duration, volume
  └── AudioEngine, QueueManager, PreloadManager

maintenance.ts  # Maintenance mode flag
page-loader.ts  # Global loading state
```

### Planned Stores (when needed)
```
library.ts      # Cached library items (liked tracks, albums, artists)
search.ts       # Search history, recent searches state
ui.ts           # Sidebar open, fullscreen player, queue drawer
```

### Service Layer (22 API modules - all built)
```
auth/       → login, register, logout, refresh, me
catalog/    → tracks, artists, albums, genres, search
playlist/   → CRUD, collaborative, add/remove tracks
library/    → liked tracks, albums, artists, history
history/    → listening history
recommendation/ → popular, best, recent, for-you, similar
social/     → follow, unfollow, feed, followers, following
reactions/  → like/love/dislike tracks, albums
creator/    → dashboard analytics, daily stats, track stats
ai/         → playlist generation, mood analysis
lyrics/     → fetch lyrics
player/     → backend player sync
subscription/ → plans, checkout, cancel
moderation/ → reports, pending, resolve
notifications/ → list, mark read
media/      → file upload
users/      → user profiles
```

---

## 6. Implementation Progress

### ✅ Phase 1 - Foundation (COMPLETE)
- Home page: 7 sections with skeleton loading
- Search overlay: keyboard nav, recent searches, type-ahead results
- Admin tracks + catalog: play buttons on rows/cards
- Discover page: data fetching fix

### ✅ Phase 2 - Core Pages (COMPLETE)
- Artist page: hero, top tracks, albums, related
- Recently played: track details + play
- Playlist page: create/delete with toasts, navigate to detail

### 🔄 Phase 3 - Immersive Experience (IN PROGRESS)
- Fullscreen Player (not yet built)
- Real-time lyrics sync (not yet built)
- Enhanced Album page (not yet built)
- Enhanced Track page (not yet built)
- Collaborative playlists (not yet built)

### ⏳ Phase 4 - Social & Discovery
- Social features (share, activity feed)
- Notification center
- Library tabs
- Recommendations hub

### ⏳ Phase 5 - Creator & Analytics
- Creator dashboard polish
- Years Wrapped
- AI features

---

## 7. Key Metrics & Priorities

| Metric | Target | Key Feature |
|--------|--------|-------------|
| Time-to-first-play | <2s from landing | Home page autoload popular tracks |
| Search-to-result | <300ms | Search overlay debounced at 250ms |
| Page transitions | <100ms | Skeleton loading + route prefetch |
| Track play initiation | <1 click from anywhere | Play buttons on all cards/rows |
| Daily active users | ↑ via push notifications | Notification center |
| Playlist creation rate | ↑ via improved UX | Toast feedback + auto-navigate |
| Social engagement | ↑ via sharing/shoutouts | Share buttons on all entities |
