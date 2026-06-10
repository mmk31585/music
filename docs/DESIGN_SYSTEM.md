# Muse — Design System & UX Specification

> **Vision:** The most beautiful, immersive, and social music streaming platform ever created.
> More premium than Spotify. More beautiful than Apple Music. More social than SoundCloud.

---

## Table of Contents

1. [Brand Identity](#1-brand-identity)
2. [Design Tokens](#2-design-tokens)
3. [Visual Language](#3-visual-language)
4. [Component System](#4-component-system)
5. [Page Architecture](#5-page-architecture)
6. [Player UX Specifications](#6-player-ux-specifications)
7. [Lyrics UX Specifications](#7-lyrics-ux-specifications)
8. [Discovery System](#8-discovery-system)
9. [Social Ecosystem](#9-social-ecosystem)
10. [Navigation System](#10-navigation-system)
11. [Motion System](#11-motion-system)
12. [Mobile UX Specifications](#12-mobile-ux-specifications)
13. [Creator Studio UX](#13-creator-studio-ux)
14. [Admin UX Specifications](#14-admin-ux-specifications)
15. [Figma Implementation Guide](#15-figma-implementation-guide)

---

## 1. Brand Identity

### Name & Logo

- **Product Name:** Muse
- **Tagline:** "Feel the music. Share the moment."
- **Logo:** Minimalist lettermark "M" in a custom wordmark, with a subtle soundwave integrated into the left stroke of the M. The logo pulses gently when music is playing.

### Color Palette

```
Primary Green:    #1DB954  (Spotify-inspired, our accent)
Primary Purple:   #B646FF  (creative energy)
Aurora Green:     #00FF87  (energy, freshness)
Aurora Blue:      #60A5FA  (depth, calm)
Aurora Pink:      #F472B6  (warmth, emotion)
Aurora Purple:    #A855F7  (mystery, creativity)

Surface Dark:     #050505  (deepest background)
Surface Base:     #0A0A0A  (main background)
Surface Raised:   #121212  (cards, containers)
Surface Overlay:  #1A1A1A  (hover states)
Surface Border:   rgba(255,255,255,0.06)
```

### Typography

```css
--font-display: 'Cabinet Grotesk', sans-serif;  /* Headlines, display text */
--font-ui:      'Inter', sans-serif;             /* UI elements, body */
--font-music:   'Satoshi', sans-serif;           /* Music metadata, lyrics */
```

Scale:

| Level | Size | Weight | Usage |
|-------|------|--------|-------|
| Display XL | 4.5rem (72px) | 800 | Hero titles |
| Display L | 3rem (48px) | 800 | Section headers |
| Display M | 2rem (32px) | 700 | Page titles |
| Heading L | 1.5rem (24px) | 700 | Card titles |
| Heading M | 1.25rem (20px) | 600 | Subsection headers |
| Body L | 1rem (16px) | 400 | Primary text |
| Body M | 0.875rem (14px) | 400 | Secondary text |
| Body S | 0.75rem (12px) | 500 | Metadata |
| Caption | 0.625rem (10px) | 600 | Labels, timestamps |

---

## 2. Design Tokens

### Spacing

4px base unit. Scale: 4, 8, 12, 16, 20, 24, 32, 40, 48, 64, 80, 96, 128.

### Border Radius

```css
--radius-sm:    8px;    /* Buttons, inputs */
--radius-md:    12px;   /* Cards, small containers */
--radius-lg:    16px;   /* Album art, modals */
--radius-xl:    24px;   /* Large cards, sections */
--radius-full:  9999px; /* Pills, avatars */
```

### Elevation

```css
--shadow-sm:    0 2px 8px rgba(0,0,0,0.3);
--shadow-md:    0 4px 16px rgba(0,0,0,0.4);
--shadow-lg:    0 8px 32px rgba(0,0,0,0.5);
--shadow-xl:    0 16px 64px rgba(0,0,0,0.6);
--shadow-glow:  0 0 30px rgba(29,185,84,0.3);
```

### Glassmorphism

```css
/* Base glass */
.glass {
  background: rgba(255,255,255,0.04);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255,255,255,0.06);
}

/* Strong glass (modals, overlays) */
.glass-strong {
  background: rgba(255,255,255,0.08);
  backdrop-filter: blur(32px);
  border: 1px solid rgba(255,255,255,0.1);
}

/* Dark glass (bottom bars, nav) */
.glass-darker {
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(24px);
  border: 1px solid rgba(255,255,255,0.05);
}
```

---

## 3. Visual Language

### Aurora System

Every major page features a subtle aurora background system:

- **Primary aurora:** Large radial gradient blobs (600px+) with extreme blur (120px)
- **Secondary aurora:** Smaller blobs offset with different timing
- **Colors blend** from the album/artist's dominant color in context
- **Global background:** Deep `#050505` with dark green/blue/purple radial gradients

```
body {
  background:
    radial-gradient(ellipse 80% 50% at 0% 0%, rgba(29,185,84,0.12), transparent 60%),
    radial-gradient(ellipse 60% 40% at 100% 0%, rgba(96,165,250,0.08), transparent 50%),
    radial-gradient(ellipse 50% 30% at 50% 100%, rgba(168,85,247,0.06), transparent 40%),
    linear-gradient(180deg, #101010 0%, #050505 45%, #000 100%);
}
```

### Dynamic Color System

Album/artist pages extract dominant colors from the cover art and use them for:

1. Page background gradients
2. Button accent colors
3. Text highlight colors
4. Progress bar colors
5. Shadow colors
6. Aurora blob colors

This creates a unique visual identity for every album and artist page, making each feel like a custom-branded experience.

### Glassmorphism Hierarchy

```
Floating elements (modals, tooltips)  → strongest blur (32px)
Side panels (sidebar, queue)          → medium blur (24px)
Content containers (cards, sections)  → subtle blur (20px)
Background aurora                     → extreme blur (120px)
```

### Typography Rhythm

- Headlines use tight tracking (-0.02em to -0.04em)
- Body text uses comfortable leading (1.5-1.6)
- Music metadata uses tabular-nums for time displays
- Lyrics use expanded letter-spacing for readability
- All caps used sparingly for section labels (tracking: 0.08em)

---

## 4. Component System

### 4.1 Audio Player (Core)

```
<AudioEngine>
├── Audio Source (HTML5 <audio> or Media Source Extensions)
├── Audio Analyzer (Web Audio API AnalyserNode)
│   ├── Frequency data (for spectrum visualizer)
│   └── Waveform data (for waveform visualization)
├── Queue Manager
│   ├── Current queue (ordered)
│   ├── History stack
│   ├── Shuffle (2 modes: smart + true random)
│   └── Repeat (off / all / one)
└── Preload Manager
    ├── Next track preload
    └── Buffer management
```

### 4.2 Media Session API Integration

```typescript
navigator.mediaSession.metadata = new MediaMetadata({
  title: track.title,
  artist: track.artistName,
  album: track.albumTitle,
  artwork: [
    { src: coverUrl, sizes: '96x96', type: 'image/png' },
    { src: coverUrl, sizes: '128x128', type: 'image/png' },
    { src: coverUrl, sizes: '256x256', type: 'image/png' },
    { src: coverUrl, sizes: '512x512', type: 'image/png' },
  ],
})

// Action handlers
navigator.mediaSession.setActionHandler('play', () => { /* ... */ })
navigator.mediaSession.setActionHandler('pause', () => { /* ... */ })
navigator.mediaSession.setActionHandler('previoustrack', () => { /* ... */ })
navigator.mediaSession.setActionHandler('nexttrack', () => { /* ... */ })
navigator.mediaSession.setActionHandler('seekto', () => { /* ... */ })
```

### 4.3 Component Tree

```
App.vue
├── LayoutAuth.vue
│   └── Auth pages (Login, Register)
├── LayoutMusicApp.vue  ← MAIN APP SHELL
│   ├── MusicSidebar.vue
│   │   ├── AppLogo
│   │   ├── NavLinks (Home, Search, Library, Discover)
│   │   ├── PlaylistSection
│   │   └── UserSection
│   ├── MusicTopbar.vue
│   │   ├── PageTitle
│   │   ├── SearchTrigger (Ctrl+K)
│   │   └── UserMenu
│   ├── RouterView (pages)
│   ├── NowPlayingBar.vue  ← Sticky bottom bar
│   ├── FloatingMiniPlayer.vue  ← Draggable floating player
│   ├── QueuePanel.vue  ← Slide-in queue
│   ├── SearchOverlay.vue  ← Fullscreen search
│   ├── FullscreenPlayer.vue  ← Immersive player
│   ├── MobileBottomSheet.vue  ← Mobile player
│   └── LyricsDisplay.vue
├── LayoutAdmin.vue
│   ├── AdminSidebar
│   ├── AdminTopbar
│   └── RouterView (admin pages)
└── LayoutEmpty.vue
```

### 4.4 Key Components

#### AlbumCard
```
┌─────────────────────────────────┐
│         ┌───────────┐           │
│         │           │           │
│         │  ALBUM    │           │
│         │  ART      │           │
│         │           │           │
│         │    ▶      │  ← hover  │
│         └───────────┘           │
│  Album Title                    │
│  Artist Name · Year             │
└─────────────────────────────────┘
- Aspect ratio: 1:1
- Play button overlay on hover (spring animation)
- Rounded corners: 16px
- Lift on hover: -4px translateY + shadow increase
```

#### ArtistCard
```
┌─────────────────────────────────┐
│         ┌───────────┐           │
│         │   ARTIST  │           │
│         │  AVATAR   │           │
│         │   (圆形)   │           │
│         └───────────┘           │
│  Artist Name                    │
│  X monthly listeners            │
└─────────────────────────────────┘
- Avatar: circular (100%)
- Subtle green ring on hover
- Follow button on hover
```

#### TrackRow
```
┌──────────────────────────────────────────────────────┐
│  №  │  ▷  │  [COVER]  │  Title           │  ♡  │  3:45 │
│      │     │  [32x32]  │  Artist          │     │       │
└──────────────────────────────────────────────────────┘
- Number or play icon in first column
- Album art (32x32 or 40x40)
- Track title + artist
- Like button
- Duration
- Explicit badge (E icon) if applicable
```

---

## 5. Page Architecture

### 5.1 Home (`/`)

```
┌─────────────────────────────────────────────────────────┐
│ Hero Section                                             │
│ ┌─────────────────────────────────────────────────────┐ │
│ │  [Aurora Flow Background]                           │ │
│ │                                                     │ │
│ │  Good evening, Alex                                 │ │
│ │  Your vibe today: Focus & Chill                     │ │
│ │                                                     │ │
│ │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐     │ │
│ │  │Mood 1│ │Mood 2│ │Mood 3│ │Mood 4│ │Mood 5│     │ │
│ │  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘     │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Made For You (Section)                                   │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ┌──┐ ┌──┐ ┌──┐ ┌──┐ ┌──┐ ┌──┐                    │ │
│ │ │D1│ │D2│ │D3│ │D4│ │D5│ │D6│                    │ │
│ │ └──┘ └──┘ └──┘ └──┘ └──┘ └──┘                    │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Continue Listening                                       │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ┌──┐ ┌──┐ ┌──┐ ┌──┐                               │ │
│ │ │C1│ │C2│ │C3│ │C4│                               │ │
│ │ └──┘ └──┘ └──┘ └──┘                               │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Trending Now + New Releases + Popular Artists            │
└─────────────────────────────────────────────────────────┘
```

### 5.2 Discover (`/discover`)

```
┌─────────────────────────────────────────────────────────┐
│ Hero Section: "Discover new music"                       │
│ ┌─────────────────────────────────────────────────────┐ │
│ │  Large aurora background with personalized message  │ │
│ │  Today's Top Pick: [Featured Track]                  │ │
│ │  Based on your listening this week                   │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Trending (Horizontal scrollable carousel)                │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ [1] [2] [3] [4] [5] [6] →                         │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ New Releases (Grid 4-6)                                  │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ ┌──┐ ┌──┐ ┌──┐ ┌──┐                               │ │
│ │ │N1│ │N2│ │N3│ │N4│                               │ │
│ │ └──┘ └──┘ └──┘ └──┘                               │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ For You (Personalized)                                   │
│ Moods │ Genres │ AI Mixes │ Weekly Discovery            │
│                                                         │
│ Similar to what you've been listening to                 │
│ Similar Artists │ Similar Tracks                         │
│                                                         │
│ Daily Mix [1-6]                                          │
└─────────────────────────────────────────────────────────┘
```

### 5.3 Search (`/search`)

```
┌─────────────────────────────────────────────────────────┐
│ Search Bar (large, centered)                             │
│ ┌─────────────────────────────────────────────────────┐ │
│ │  🔍 Search for songs, artists, albums, playlists... │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Filter Tabs: All | Tracks | Artists | Albums | Playlists │
│ ───────────────────────────────────────────────────────  │
│                                                         │
│ Results:                                                 │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Top Result                                          │ │
│ │ ┌──────────────────────────────────────────────┐   │ │
│ │ │ [Large] Artist / Track card                   │   │ │
│ │ └──────────────────────────────────────────────┘   │ │
│ │                                                     │ │
│ │ Tracks (list)                                      │ │
│ │ Artists (horizontal scroll)                        │ │
│ │ Albums (grid)                                      │ │
│ │ Playlists (grid)                                   │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### 5.4 Library (`/library`)

```
┌─────────────────────────────────────────────────────────┐
│ Header: "Your Library"                                  │
│                                                         │
│ Tabs: Tracks | Albums | Artists | Playlists              │
│ ───────────────────────────────────────────────────────  │
│                                                         │
│ Search bar (filter library)                              │
│ Sort: Recent | A-Z | Recently Added                     │
│                                                         │
│ Content Grid/List                                        │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ Filter chips: Liked | Downloaded | Offline          │ │
│ │                                                     │ │
│ │ [Content based on active tab]                       │ │
│ │                                                     │ │
│ │ Stats: X tracks · Y albums · Z artists              │ │
│ └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

### 5.5 Album (`/album/:id`)

```
┌─────────────────────────────────────────────────────────┐
│ Hero Section (dynamic album color background)            │
│ ┌─────────────────────────────────────────────────────┐ │
│ │  [Large Album Art]    Album Title                    │ │
│ │  300x300             Artist Name                     │ │
│ │                      Year · Genre · X songs          │ │
│ │                      ▶ Play All   ♡ Like  ↓ Download │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Track List                                               │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ №  Title                    ♡  Streams   Duration   │ │
│ │ 1   Track 1                  ♡  1.2M      3:45      │ │
│ │ 2   Track 2                  ♡  856K      4:02      │ │
│ │ 3   Track 3 (feat. Artist)   ♡  2.1M      3:12      │ │
│ │ ...                                                  │ │
│ │                                     Total: 42:30     │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ More by Artist (horizontal scroll)                       │
│ Fans also like                                           │
└─────────────────────────────────────────────────────────┘
```

### 5.6 Artist (`/artist/:id`)

```
┌─────────────────────────────────────────────────────────┐
│ Hero Section (giant, cinematic)                          │
│ ┌─────────────────────────────────────────────────────┐ │
│ │                                                     │ │
│ │  [Large parallax artist image]                      │ │
│ │                                                     │ │
│ │  Artist Name (big, bold)                            │ │
│ │  X monthly listeners                                │ │
│ │  [Follow] [Play All] [Share]                        │ │
│ │                                                     │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Popular Tracks (top 10)                                  │
│ ┌─────────────────────────────────────────────────────┐ │
│ │ 1. Hit Song 1           ▶  ♡  3:45  123M plays    │ │
│ │ 2. Hit Song 2           ▶  ♡  4:02  98M plays     │ │
│ │ 3. Hit Song 3           ▶  ♡  3:12  85M plays     │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Discography (grid of albums)                             │
│ Similar Artists (horizontal scroll)                      │
│ Appears On (guest appearances)                           │
│ Fans also like                                           │
└─────────────────────────────────────────────────────────┘
```

### 5.7 Playlist (`/playlist/:id`)

```
┌─────────────────────────────────────────────────────────┐
│ Hero Section (gradient based on playlist color)          │
│ ┌─────────────────────────────────────────────────────┐ │
│ │  [Cover]    Playlist Name                           │ │
│ │  200x200    Description                              │ │
│ │             Created by User · X songs · duration     │ │
│ │             ▶ Play  ♡ Like  ↓ Download  ⋮ More      │ │
│ │                                                     │ │
│ │             Collaborative: X users editing            │ │
│ └─────────────────────────────────────────────────────┘ │
│                                                         │
│ Track List (same as album)                               │
│                                                         │
│ Recommended tracks (at bottom)                           │
└─────────────────────────────────────────────────────────┘
```

---

## 6. Player UX Specifications

### 6.1 Core Player Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    PLAYER SYSTEM                          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  NowPlayingBar (persistent bottom bar)            │   │
│  │  Height: 72px desktop, 64px mobile               │   │
│  │  Glass-darker background                          │   │
│  │  3 columns: track | controls | extras             │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  FullscreenPlayer (overlay, z-index: 200)         │   │
│  │  2 columns: art+controls | lyrics/queue/spectrum  │   │
│  │  Dynamic background from album art                │   │
│  │  Transition: scale(0.96) + fade (250ms)           │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  FloatingMiniPlayer (draggable PiP)               │   │
│  │  Sizes: mini (64px) | compact (120px) | full     │   │
│  │  Draggable anywhere on screen                     │   │
│  │  Picture-in-Picture API fallback                  │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  MobileBottomSheet (mobile player)                │   │
│  │  Snap points: collapsed(64px) | half | full      │   │
│  │  Gesture: drag up to expand, down to collapse     │   │
│  └──────────────────────────────────────────────────┘   │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │  QueuePanel (slide-in from right)                 │   │
│  │  Width: 360px                                     │   │
│  │  Tabs: Queue | Recommendations                    │   │
│  │  Glass-strong background                           │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### 6.2 NowPlayingBar — Detailed

```
┌──────────────────────────────────────────────────────────────┐
│  [progress bar: 4px, green gradient]                         │
│ ┌──────────────────────────────────────────────────────────┐ │
│ │ ┌─────────────────┐ ┌──────────────────────┐ ┌────────┐ │ │
│ │ │ [48x48]          │ │    ⏮  ▶⏸  ⏭         │ │  🔲     │ │
│ │ │ Track Title  ⚡  │ │ ───●─────────────── │ │  📋  🔊 │ │
│ │ │ Artist Name      │ │  1:23         3:45   │ │  ───●─ │ │
│ │ └─────────────────┘ └──────────────────────┘ └────────┘ │ │
│ │ ← Track Info         ← Controls + Seek        ← Extras   │ │
│ └──────────────────────────────────────────────────────────┘ │
│  Left (minmax 0, 1fr)                                        │
│  Center (minmax 340px, 1fr)                                  │
│  Right (minmax 0, 1fr)                                       │
└──────────────────────────────────────────────────────────────┘
```

**States:**
- **Empty:** Show placeholder text "Select a track to start listening"
- **Loading:** Small spinner badge next to title
- **Buffering:** Pulsing progress bar
- **Error:** Red error text, retry button
- **Playing:** Spinning album art + equalizer animation
- **Paused:** Static album art, paused icon

**Keyboard Shortcuts:**
- `Space` — Play/Pause
- `→` — Seek forward 5s
- `←` — Seek backward 5s
- `↑` — Volume up 5%
- `↓` — Volume down 5%
- `Ctrl+→` — Next track
- `Ctrl+←` — Previous track
- `F` — Toggle fullscreen player
- `M` — Mute/unmute

### 6.3 FullscreenPlayer — Detailed

```
┌──────────────────────────────────────────────────────────────┐
│  [Dynamic background: extracted album colors + aurora blobs] │
│                                                              │
│  Top bar (transparent)                                        │
│  ┌────────────────────────────────────────────────────────┐  │
│  │  ↓ Close      [NOW PLAYING badge]    📋 🎤 📋  ⏏   │  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────┬───────────────────────────┐│
│  │     ALBUM ART + CONTROLS     │      LYRICS / QUEUE       ││
│  │                              │      / SPECTRUM           ││
│  │         ┌────────┐          │                           ││
│  │         │        │          │  ┌─────────────────────┐  ││
│  │         │  ART   │          │  │                     │  ││
│  │         │  CD    │          │  │  Lyrics or          │  ││
│  │         │  SPIN  │          │  │  Spectrum or        │  ││
│  │         │        │          │  │  Queue              │  ││
│  │         └────────┘          │  │                     │  ││
│  │                              │  └─────────────────────┘  ││
│  │  Track Title                │                           ││
│  │  Artist — Album             │                           ││
│  │                              │                           ││
│  │  ──●──────────────────────  │                           ││
│  │  1:23              3:45    │                           ││
│  │                              │                           ││
│  │  🔀 ⏮ ▶⏸ ⏭ 🔁   │                           ││
│  │                              │                           ││
│  │  🔊 ──●──  [1x]   🕐      │                           ││
│  └──────────────────────────────┴───────────────────────────┘│
│                                                              │
│  [Sleep timer overlay - bottom center]                        │
└──────────────────────────────────────────────────────────────┘
```

**Album Art States:**
- **Playing:** Spinning at 60s/rotation, subtle green glow shadow
- **Paused:** Paused spin, dimmed glow
- **Transition:** Art pop animation (scale 0.92 → 1, 250ms ease-out)

**Controls Detail:**

| Control | Icon | Behavior |
|---------|------|----------|
| Shuffle | 🔀 | Toggle on/off, green when active |
| Previous | ⏮ | Go to previous track in history |
| Play/Pause | ▶/⏸ | Large center button, 64px |
| Next | ⏭ | Skip to next track in queue |
| Repeat | 🔁 | Cycle: off → all → one → off |
| Volume | 🔊/🔇 | Slider with mute toggle |
| Speed | 1x | Cycle: 0.5x → 0.75x → 1x → 1.25x → 1.5x → 2x |
| Sleep Timer | 🕐 | Dropdown: 5/15/30/45/60 min, end of track |

**Spectrum Visualizer:**
- 64 bars, animated when playing
- Gradient from green (high) → blue (mid) → purple (low)
- Smooth transitions using requestAnimationFrame
- Bars fade to 8% height when paused
- Rounded bar caps (barW/2 radius)

### 6.4 FloatingMiniPlayer — Detailed

```
┌─────────────┐    ┌──────────────────┐    ┌──────────────────────┐
│  ████████    │    │  ┌────┐ ▶⏸ ⏭   │    │                      │
│  ████████    │    │  │    │ Track    │    │  ┌────┐              │
│  ████████    │    │  │art │ Artist   │    │  │    │              │
│  ████████    │    │  └────┘          │    │  │art │ ▶⏸  ⏭      │
│  ████████    │    │           🔊──●─ │    │  │    │              │
│  MINI        │    │  COMPACT          │    │  └────┘              │
│  64x64       │    │  200x120          │    │  Track Title         │
└─────────────┘    └──────────────────┘    │  Artist Name          │
             Draggable + Resizable          │  ──●─────────────     │
                                            │  1:23        3:45     │
                                            │                      │
                                            │  EXTENDED            │
                                            │  320x200             │
                                            └──────────────────────┘
```

**Modes:**
- **Mini:** Just album art (64x64), tap to show controls
- **Compact:** Small art + basic controls (play/pause, prev/next, volume)
- **Extended:** Full mini player with seekbar and track info
- **PiP:** Browser Picture-in-Picture API (video fallback with album art)

**Interaction:**
- Drag by the header area
- Double-click to toggle between mini and compact
- Right-click for context menu (close, go to fullscreen, dock)
- Fade to 50% opacity when not hovered for 3s (ambient mode)

### 6.5 QueuePanel — Detailed

```
┌──────────────────────────────────┐
│  Queue 📋                        │
│  ─────────────────────────       │
│  [Queue] [Recommendations]       │
│  ─────────────────────────       │
│                                  │
│  Now Playing:                    │
│  ┌────────────────────────────┐  │
│  │ [40x40] Track Name 🔈      │  │
│  │          Artist            │  │
│  └────────────────────────────┘  │
│                                  │
│  Next Up:                       │
│  ┌────────────────────────────┐  │
│  │ [40x40] Track 2           │  │
│  │          Artist           │  │
│  ├────────────────────────────┤  │
│  │ [40x40] Track 3           │  │
│  │          Artist           │  │
│  ├────────────────────────────┤  │
│  │ [40x40] Track 4           │  │
│  │          Artist           │  │
│  └────────────────────────────┘  │
│                                  │
│  Drag to reorder                 │
└──────────────────────────────────┘
```

---

## 7. Lyrics UX Specifications

### 7.1 Lyrics Architecture

```
Lyrics System
├── Parsers
│   ├── LRC Parser (timestamped [mm:ss.xx] format)
│   ├── Synced Lyrics Parser (JSON format with word-level timing)
│   └── Plain Text Parser (fallback)
├── Display Modes
│   ├── Scroll (standard, line-by-line synced)
│   ├── Karaoke (word-level highlighting)
│   └── Translation (side-by-side original + translated)
└── Visual Effects
    ├── Particle system (floating particles)
    ├── Background aurora
    └── Dynamic color from album art
```

### 7.2 Lyrics Display Modes

**Scroll Mode:**
```
┌──────────────────────────────────────────────┐
│                                              │
│   以前的人們                                  │
│   (previous lines in white/15)                │
│                                              │
│   從哪來的                                        │  ← ACTIVE LINE (bold, white, scale 1.05)
│   Where did it all come from                 │  ← Translation below
│                                              │
│   歌聲就從這個人的                            │
│   (next lines in white/25)                   │
│                                              │
│   嘴巴裡唱出來喔                              │
│                                              │
│   ────                                      │
│   🎤 Artist's note about this lyric           │
└──────────────────────────────────────────────┘
```

**Karaoke Mode:**
```
┌──────────────────────────────────────────────┐
│                                              │
│   以前的人們                                  │
│                                              │
│   從哪來的                                        │
│   ┌──────────────────────────────────────┐   │
│   │  [Where] [did] [it] [all] [come]      │   │
│   │              ↑ GREEN HIGHLIGHT         │   │
│   │  [from]                               │   │
│   └──────────────────────────────────────┘   │
│                                              │
│   歌聲就從這個人的                            │
│                                              │
│   Word by word highlighting                  │
│   Current word: green + glow                 │
│   Future words: white/40                     │
│   Past words: white/15                       │
└──────────────────────────────────────────────┘
```

**Translation Mode:**
```
┌──────────────────────────────┬───────────────┐
│  ORIGINAL                    │  TRANSLATION   │
│                              │               │
│  以前的人們                   │  People before │
│                              │               │
│  從哪來的                     │  Where did it  │
│                              │  all come      │
│  歌聲就從這個人的              │  from          │
│                              │               │
│  嘴巴裡唱出來喔               │  The voice     │
│                              │  came out of   │
│                              │  this person's │
│                              │  mouth         │
└──────────────────────────────┴───────────────┘
```

### 7.3 Particle System

Canvas-based particle system for lyrics background:

- **Particles:** 50-100 small circles (2-4px radius)
- **Colors:** Extracted from album art palette
- **Behavior:** Float upward slowly, gentle wind effect
- **Interaction:** Particles slightly avoid mouse cursor
- **Performance:** requestAnimationFrame, canvas 2D context
- **Audio reactivity:** Particle speed correlates with volume/energy

### 7.4 Lyrics States

| State | Visual | Interaction |
|-------|--------|-------------|
| Loading | Skeleton lines (6 bars, varied widths) | None |
| No lyrics | Centered icon + "No lyrics available" | None |
| Synced available | Full karaoke experience | Click line to seek |
| Plain text only | Full text, no timestamps | Manual scroll |
| Error | Error message + retry button | Retry action |

---

## 8. Discovery System

### 8.1 Discovery Page Sections

```
┌─────────────────────────────────────────────────────────────┐
│ DISCOVER PAGE                                                │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ 1. HERO: Personalized greeting + Today's Top Pick           │
│    [Large aurora background with dynamic gradient]            │
│                                                              │
│ 2. TRENDING: Top 50 tracks right now                         │
│    [Horizontal scroll, numbered list cards with cover art]    │
│                                                              │
│ 3. VIRAL: Fastest rising tracks                              │
│    [Same format, "VIRAL" badge with fire icon]                │
│                                                              │
│ 4. NEW RELEASES: Latest from followed + popular artists      │
│    [Grid 4 columns, album cards with "NEW" badge]             │
│                                                              │
│ 5. FOR YOU: AI-curated recommendations                       │
│    [Mixed format: tracks list + album cards]                  │
│                                                              │
│ 6. MOODS: Browse by mood                                     │
│    [Grid of mood cards with gradient backgrounds]             │
│    Moods: Energize | Relax | Focus | Workout | Chill |       │
│           Party | Romance | Sleep | Sad | Happy              │
│                                                              │
│ 7. GENRE: Browse by genre                                    │
│    [Horizontal scroll of genre pills/cards]                   │
│                                                              │
│ 8. AI MIXES: Generated playlists based on your taste         │
│    [Carousel "AI-curated for you"]                            │
│                                                              │
│ 9. WEEKLY DISCOVERY: Personalized new music (updated Fri)    │
│    [Special card with weekly badge + track list preview]      │
│                                                              │
│ 10. DAILY MIX [1-6]: Curated mixes by style                  │
│     [Grid of mix cards, each with unique gradient]            │
│                                                              │
│ 11. SIMILAR ARTISTS: Based on recent listens                 │
│     [Artist cards horizontal scroll]                          │
│                                                              │
│ 12. SIMILAR TRACKS: Based on current mood                    │
│     [Track list]                                              │
│                                                              │
│ 13. RECENTLY PLAYED (for signed-in users)                    │
│     [History-based horizontal scroll]                         │
└─────────────────────────────────────────────────────────────┘
```

### 8.2 Mood Cards

10 moods, each with unique gradient + icon:

| Mood | Gradient | Icon |
|------|----------|------|
| Energize | orange → red | ⚡ |
| Relax | teal → blue | 🌊 |
| Focus | blue → purple | 🎯 |
| Workout | red → orange | 💪 |
| Chill | green → blue | 🧘 |
| Party | purple → pink | 🎉 |
| Romance | pink → rose | 💕 |
| Sleep | indigo → navy | 🌙 |
| Sad | blue → gray | 🌧️ |
| Happy | yellow → green | ☀️ |

### 8.3 AI Playlist Generator

```
┌─────────────────────────────────────────────────────────────┐
│ AI Playlist Generator                                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │                                                         │ │
│ │  "Describe your perfect playlist..."                     │ │
│ │                                                         │ │
│ │  [Multiline input field with glass background]           │ │
│ │                                                         │ │
│ │  e.g., "Rainy day jazz with lo-fi beats and coffee"     │ │
│ │                                                         │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                              │
│ Mood Selectors (icon-based)                                  │
│ Mood: [Energize] [Relax] [Focus] [Workout] [Chill]          │
│ Genre: [select dropdown]                                     │
│ Duration: [15m] [30m] [1h] [2h]                              │
│                                                              │
│ [Generate Playlist] button (prominent, green gradient)       │
│                                                              │
│ ── Results (appear after generation) ──                      │
│                                                              │
│ Generated: "Rainy Day Jazz" (12 tracks, 48 min)              │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ 1. Track 1 - Artist 1           3:45                   │ │
│ │ 2. Track 2 - Artist 2           4:02                   │ │
│ │ ...                                                    │ │
│ │ [Save to Library] [Play All] [Regenerate]               │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

## 9. Social Ecosystem

### 9.1 Social Features Matrix

| Feature | Visibility | Interaction | Location |
|---------|-----------|-------------|----------|
| Follow User | Public | Button on profile | Profile, Search |
| Follow Artist | Public | Button on artist page | Artist Page |
| Like Track | Public | Heart icon | Every track, player |
| Love Track | Public | Double-tap heart | Player |
| Comment | Public | Text input | Track page |
| Share | External | Share sheet | Player, track |
| Activity Feed | Friends | Timeline | Home, Profile |
| Stories | 24h | Story rings | Home top bar |
| Listening Together | Real-time | Join session | Player |
| Collaborative Playlist | Invite | Multi-editor | Playlist page |
| Profile Badges | Public | Badge display | Profile |
| Artist Verification | Verified only | Blue checkmark | Artist page |
| Reactions | Public | Emoji picker | Track, comment |

### 9.2 Activity Feed

```
┌─────────────────────────────────────────────────────────────┐
│ Following                                                    │
│ ───────────────────────────────                              │
│                                                              │
│ [User Avatar] Alex listened to "Bohemian Rhapsody"    2m ago │
│              ┌──────────────────────────────────────┐       │
│              │ ▶ Listen now                          │       │
│              └──────────────────────────────────────┘       │
│                                                              │
│ [Artist Badge] Taylor Swift released new album         15m ago│
│              [Cover] "The Tortured Poets Department"         │
│              [Listen] [Save] [Share]                         │
│                                                              │
│ [User Avatar] Sarah created playlist "Summer Vibes"   1h ago │
│              "Chill summer tracks for road trips"            │
│              [View Playlist]                                 │
│                                                              │
│ [User Avatar] Mike followed 3 new artists             2h ago │
│              Artist1 · Artist2 · Artist3                     │
│                                                              │
│ [User Avatar] You reached 1000 listening hours!       3h ago │
│              🎉 Milestone achieved!                          │
└─────────────────────────────────────────────────────────────┘
```

### 9.3 Music Sharing

```
Share Sheet (overlay, glass-strong):
┌──────────────────────────────┐
│  Share "Track Name"          │
│  ─────────────────────────── │
│                              │
│  ┌────────────────────────┐  │
│  │ [40x40] Track          │  │
│  │         Artist         │  │
│  └────────────────────────┘  │
│                              │
│  Share to:                   │
│  ┌──┐ ┌──┐ ┌──┐ ┌──┐      │
│  │IG│ │TG│ │WA│ │X │      │
│  └──┘ └──┘ └──┘ └──┘      │
│                              │
│  Copy Link    Embed          │
│  Q R Code                    │
│                              │
│  Include:                    │
│  [☑] Album art              │
│  [☑] My reaction            │
│  [☐] Timestamp (1:23)       │
└──────────────────────────────┘
```

### 9.4 Collaborative Playlists

```
┌─────────────────────────────────────────────────────────────┐
│ Collaborative Playlist                                       │
│                                                              │
│ [Collaborative Badge] ✨ 5 friends editing                    │
│                                                              │
│ Recent activity:                                             │
│   Sarah added "Track X" · 2m ago                             │
│   Mike removed "Track Y" · 5m ago                            │
│   You added "Track Z" · 10m ago                              │
│                                                              │
│ Collaborators: [Alex] [Sarah] [Mike] [Anna] [+ Invite]       │
│                                                              │
│ Track List (with who added each)                              │
│ 1. Track A — Artist A             Added by Alex       3:45   │
│ 2. Track B — Artist B             Added by Sarah      4:02   │
│ 3. Track C — Artist C             Added by You        3:12   │
│ ...                                                          │
└─────────────────────────────────────────────────────────────┘
```

---

## 10. Navigation System

### 10.1 Sidebar Navigation (Desktop)

```
┌──────────────────┐
│  [MUSE LOGO]      │  ← 32x32, pulses when playing
│                   │
│  🏠 Home          │  ← Active state: green text + left border
│  🔍 Search        │
│  📚 Library       │
│  ✨ Discover      │
│                   │
│  ── LIBRARY ──    │  ← Section label
│  ▶ Liked Songs    │
│  📋 Playlists     │
│  ⏱ Recently      │
│                   │
│  ── SOCIAL ──     │
│  👥 Following     │
│  🔔 Notifications │
│                   │
│  [AI Playlist Gen]│  ← Special CTA button
│                   │
│  [User Avatar]    │
│  Alex          ↓  │  ← User menu dropdown
└──────────────────┘
Width: 280px (collapsible to 72px icon-only mode)
```

### 10.2 Top Bar

```
┌─────────────────────────────────────────────────────────────┐
│  ← Page Title            [Search (Ctrl+K)]  🔔  [👤 Alex]  │
│  Subtitle/breadcrumb                                        │
└─────────────────────────────────────────────────────────────┘
Height: 64px
Sticky at top
Glass-darker background with blur
```

### 10.3 Bottom Navigation (Mobile)

```
┌─────────────────────────────────────────────────────────────┐
│  🏠     🔍     📚     ✨     👤                              │
│ Home  Search  Library  Discover Profile                     │
└─────────────────────────────────────────────────────────────┘
Height: 56px + safe area inset
Glass-darker with top border
Active: green text + subtle indicator dot
```

### 10.4 Gesture Navigation

| Gesture | Action |
|---------|--------|
| Swipe left (track) | Add to queue |
| Swipe right (track) | Like track |
| Swipe down (player) | Close mini player |
| Swipe up (player) | Open fullscreen player |
| Double tap (art) | Love track |
| Long press (track) | Context menu |
| Pinch (album art) | Enter/exit fullscreen |
| Two-finger swipe | Volume adjust |

---

## 11. Motion System

### 11.1 Duration Principles

| Context | Duration | Easing |
|---------|----------|--------|
| Page transitions | 350ms | cubic-bezier(0.4, 0, 0.2, 1) |
| Card hover | 200ms | cubic-bezier(0.4, 0, 0.2, 1) |
| Button press | 100ms | ease-out |
| Modal enter | 250ms | cubic-bezier(0.16, 1, 0.3, 1) |
| Modal exit | 150ms | ease-in |
| Progress bar | 100ms | linear |
| Player transition | 300ms | cubic-bezier(0.16, 1, 0.3, 1) |
| Like animation | 400ms | spring(0.4, 1) |

### 11.2 Defined Animations

```css
/* Page entry: content slides up with fade */
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Scale entry for modals/overlays */
@keyframes scale-in {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}

/* Slide-in for queue panel */
@keyframes slide-in-right {
  from { transform: translateX(100%); }
  to { transform: translateX(0); }
}

/* Vinyl spin for album art */
@keyframes spin-slow {
  to { transform: rotate(360deg); }
}
/* Duration: 60s/rotation (playing), paused (not playing) */

/* Equalizer bars animation */
@keyframes equalizer {
  0%, 100% { transform: scaleY(0.4); }
  50% { transform: scaleY(1); }
}

/* Album art glow pulse */
@keyframes vinyl-glow {
  0% { box-shadow: 0 0 30px rgba(29,185,84,0.1); }
  50% { box-shadow: 0 0 60px rgba(29,185,84,0.25); }
  100% { box-shadow: 0 0 30px rgba(29,185,84,0.1); }
}

/* Aurora flow (background gradient shift) */
@keyframes aurora-flow {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

/* Heart like animation */
@keyframes heart-pop {
  0% { transform: scale(1); }
  30% { transform: scale(1.3); }
  60% { transform: scale(0.9); }
  100% { transform: scale(1); }
}
```

### 11.3 Spring Animations

Used for interactive elements (buttons, toggles, drag interactions):

```css
.spring {
  transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}
```

---

## 12. Mobile UX Specifications

### 12.1 Mobile Layout Architecture

```
┌──────────────────────────────┐
│  Status Bar (transparent)     │
├──────────────────────────────┤
│  Top Bar (contextual)         │
├──────────────────────────────┤
│                              │
│  CONTENT                     │
│  (scrollable)                │
│                              │
│                              │
├──────────────────────────────┤
│  Mini Player (64px)          │
├──────────────────────────────┤
│  Bottom Nav (56px)           │
└──────────────────────────────┘
```

### 12.2 Mobile Bottom Sheet Player

Snap points:
1. **Collapsed (64px):** Mini player bar with track info + play/pause
2. **Half (~40%):** Album art (small), controls, progress
3. **Full (100%):** Fullscreen mobile player

```
SNAP: COLLAPSED (64px)
┌────────────────────────────────┐
│ [40x40] Track Title      ▶⏸ ⏭│
│         Artist Name            │
└────────────────────────────────┘

SNAP: HALF
┌────────────────────────────────┐
│  [← Drag handle]               │
│                                │
│     ┌────────────────┐        │
│     │                │        │
│     │  ALBUM ART     │        │
│     │  200x200       │        │
│     │                │        │
│     └────────────────┘        │
│                                │
│  Track Title                   │
│  Artist Name                   │
│                                │
│  ──●────────────────         │
│  1:23             3:45        │
│                                │
│  ⏮  ▶⏸  ⏭                    │
│                                │
│  🔊 ──●──  📋  ♡              │
└────────────────────────────────┘

SNAP: FULL
┌────────────────────────────────┐
│  [Dynamic background]          │
│                                │
│     ┌────────────────┐        │
│     │                │        │
│     │  ALBUM ART     │        │
│     │  FULL WIDTH    │        │
│     │                │        │
│     └────────────────┘        │
│                                │
│  Track Title                   │
│  Artist — Album                │
│                                │
│  ──●────────────────         │
│  1:23             3:45        │
│                                │
│  🔀  ⏮  ▶⏸  ⏭  🔁          │
│                                │
│  🔊 ──●──  🕐  📋  ♡  🔽    │
│                                │
│  [Lyrics toggle]               │
└────────────────────────────────┘
```

### 12.3 Gesture Interactions

| Gesture | Element | Action |
|---------|---------|--------|
| Swipe up | Mini player | Expand to half/full |
| Swipe down | Full player | Collapse to mini |
| Swipe left | Track row | Add to queue |
| Swipe right | Track row | Like/unlike |
| Long press | Track row | Context menu |
| Tap twice | Album art | Toggle like |
| Pinch | Album art | Toggle fullscreen |
| Horizontal swipe | Now playing | Switch to lyrics view |

### 12.4 Mobile-Specific Components

**SwipeableTrackRow:**
```
┌────────────────────────────────────┐
│ ← Add to Queue [Track] Like →      │
│    (gray)          (green)          │
│    ┌────┬──────────────────┬──┐    │
│    │    │ Title            │Heart│ │
│    │Cover│ Artist           │Menu │ │
│    │    │                  │    │  │
│    └────┴──────────────────┴──┘    │
└────────────────────────────────────┘
Swipe feedback: spring animation
Action buttons revealed on swipe
```

---

## 13. Creator Studio UX

### 13.1 Layout

```
┌─────────────────────────────────────────────────────────────┐
│  [Logo] Creator Studio                                       │
│  ─────────────────────────────────────────────────────────  │
│                                                              │
│  ┌───────┐ ┌─────────────────────────────────────────────┐  │
│  │       │ │ Dashboard                                    │  │
│  │ Dash  │ │                                              │  │
│  │ Music │ │ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐        │  │
│  │ Albums│ │ │ 1.2M │ │ 45K  │ │ 12K  │ │$8.4K │        │  │
│  │ Stats │ │ │Streams│ │Follower│ │Engage.│ │Revenue│        │  │
│  │ Fans  │ │ └──────┘ └──────┘ └──────┘ └──────┘        │  │
│  │       │ │                                              │  │
│  │Profile│ │ Streams Chart (7d/30d/12m)                   │  │
│  │       │ │ ┌────────────────────────────────────────┐   │  │
│  │Verific.│ │ │ [Chart: line chart with gradient fill] │   │  │
│  │       │ │ └────────────────────────────────────────┘   │  │
│  │Settings│ │                                              │  │
│  │       │ │ Top Tracks (this month)                       │  │
│  │       │ │ 1. Track A — 450K streams                     │  │
│  │       │ │ 2. Track B — 320K streams                     │  │
│  │       │ │ 3. Track C — 210K streams                     │  │
│  │       │ │                                              │  │
│  │       │ │ Recent Releases                               │  │
│  │       │ │ [Album Card] [Album Card] [Album Card]       │  │
│  └───────┘ └─────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 13.2 Upload Flow

```
Step 1: Upload Track
┌────────────────────────────────────────────┐
│  Drop audio file here or click to browse    │
│  Supported: MP3, FLAC, WAV, AAC, OGG       │
│  Max: 200MB per track                       │
│                                            │
│  ┌────────────────────────────────────┐    │
│  │  [Drag & drop zone with dashed     │    │
│  │   border and glass background]     │    │
│  │                                    │    │
│  │  📂 Browse Files    or    Paste URL │    │
│  └────────────────────────────────────┘    │
└────────────────────────────────────────────┘

Step 2: Track Details
┌────────────────────────────────────────────┐
│  Track Title: [________________]            │
│  Artist Name: [________________]            │
│  Featuring:  [________________]             │
│  Album:      [Select Album ▼] [+ New]      │
│  Genre:      [Select Genre ▼]               │
│  Mood:       [Select Mood ▼]                │
│  Language:   [Select ▼]                     │
│                                            │
│  Lyrics:     [Paste lyrics...]              │
│                                            │
│  Cover Art:  [Drop image]                   │
│              Recommended: 3000x3000         │
│                                            │
│  Explicit:   [☐] This track contains E     │
└────────────────────────────────────────────┘

Step 3: Publishing
┌────────────────────────────────────────────┐
│  Release Date: [Now] [Schedule ▼]          │
│                                            │
│  Distribution:                             │
│  [☑] Muse platform                         │
│  [☐] Distribute to partners (coming soon)  │
│                                            │
│  Visibility: [Public] [Unlisted] [Private] │
│                                            │
│  [← Back]           [Upload & Publish]     │
└────────────────────────────────────────────┘
```

---

## 14. Admin UX Specifications

### 14.1 Admin Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│  Admin Dashboard                                             │
│  ─────────────────────────────────────────────────────────  │
│                                                              │
│  KPI Cards:                                                  │
│  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐             │
│  │ 2.4M │ │ 45K  │ │ 120K │ │ 8.2K │ │ 892  │             │
│  │Users │ │Artist │ │Tracks │ │Albums │ │Reports│             │
│  └──────┘ └──────┘ └──────┘ └──────┘ └──────┘             │
│                                                              │
│  Charts Row:                                                 │
│  ┌────────────────────────┐ ┌────────────────────────┐      │
│  │ New Users (7d)          │ │ Streams (7d)            │      │
│  │ [Line chart]            │ │ [Area chart]            │      │
│  └────────────────────────┘ └────────────────────────┘      │
│                                                              │
│  Recent Activity:                     Pending Moderation:    │
│  ┌────────────────────────┐          ┌────────────────┐     │
│  │ New user: Alex (2m ago)│          │ 12 reports      │     │
│  │ New track: X (5m ago)  │          │ 3 pending       │     │
│  │ Report: spam (10m ago) │          │ [Review Now]    │     │
│  └────────────────────────┘          └────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

### 14.2 Moderation Queue

```
┌─────────────────────────────────────────────────────────────┐
│  Moderation Queue                             [Filters ▼]   │
│  ─────────────────────────────────────────────────────────  │
│                                                              │
│  Tabs: Pending (12) | Reviewed (892) | Appeals (3)          │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ Report #1,234  │ Spam         │ 5m ago │ [Resolve] [Dismiss]│
│  │ Track: "Song" by Artist       │ User: @spammer          │ │
│  ├────────────────────────────────────────────────────────┤ │
│  │ Report #1,233  │ Copyright    │ 15m ago│ [Resolve] [Dismiss]│
│  │ Album: "Album" by Artist      │ User: @copyright_holder │ │
│  ├────────────────────────────────────────────────────────┤ │
│  │ Report #1,232  │ Explicit     │ 1h ago │ [Resolve] [Dismiss]│
│  │ Track: "Track" by Artist      │ Reported by 3 users     │ │
│  └────────────────────────────────────────────────────────┘ │
│                                                              │
│  Bulk Actions: [Delete Selected] [Dismiss Selected]          │
└─────────────────────────────────────────────────────────────┘
```

### 14.3 User Management

```
┌─────────────────────────────────────────────────────────────┐
│  Users                              Search [_________]       │
│  ─────────────────────────────────────────────────────────  │
│                                                              │
│  Filters: [All] [Active] [Suspended] [Premium] [Admin]       │
│                                                              │
│  ┌────┬──────────┬────────┬─────────┬──────────┬─────────┐ │
│  │ ID │ Username │ Email  │ Status  │ Plan     │ Joined  │ │
│  ├────┼──────────┼────────┼─────────┼──────────┼─────────┤ │
│  │ 1  │ alex     │ ...    │ Active  │ Premium  │ Jan 24  │ │
│  │ 2  │ sarah    │ ...    │ Active  │ Free     │ Mar 24  │ │
│  │ 3  │ mike     │ ...    │ Suspend │ —        │ Feb 24  │ │
│  └────┴──────────┴────────┴─────────┴──────────┴─────────┘ │
│                                                              │
│  Pagination: [<] 1 2 3 ... 24 [>]    Per page: 25           │
└─────────────────────────────────────────────────────────────┘
```

---

## 15. Figma Implementation Guide

### 15.1 File Structure

```
Muse Design System.fig
├── 🎨 Cover / Thumbnail
├── 🏷️ 01 - Brand Identity
│   ├── Logo (variants: light, dark, favicon)
│   ├── Color Palette
│   ├── Typography
│   └── Iconography
├── 🧩 02 - Design Tokens
│   ├── Colors (with dark mode variants)
│   ├── Spacing (4px grid)
│   ├── Border Radius
│   ├── Shadows / Elevation
│   └── Glassmorphism styles
├── 🧱 03 - Components
│   ├── Buttons (primary, secondary, ghost, icon)
│   ├── Inputs (text, search, select, range)
│   ├── Cards (album, artist, playlist, track)
│   ├── Navigation (sidebar, topbar, bottom nav)
│   ├── Player (now playing, fullscreen, mini, PiP)
│   ├── Modals (share, context menu, dialogs)
│   ├── Progress (bars, spinners, skeleton)
│   └── Social (avatar, badge, reaction, comment)
├── 📄 04 - Pages
│   ├── Home
│   ├── Discover
│   ├── Search
│   ├── Library
│   ├── Album Detail
│   ├── Artist Detail
│   ├── Playlist Detail
│   ├── Track Detail
│   ├── Player (fullscreen + mini)
│   ├── Lyrics (scroll + karaoke)
│   ├── Notifications
│   ├── Profile
│   ├── Creator Studio
│   ├── Admin Dashboard
│   └── Settings
├── 📱 05 - Mobile
│   ├── Bottom Sheet Player (3 snap states)
│   ├── Mobile Now Playing
│   ├── Mobile Home
│   ├── Mobile Search
│   ├── Mobile Library
│   ├── Swipeable Track Row
│   └── Gesture interactions
├── 🎬 06 - Motion
│   ├── Page transitions
│   ├── Card hover states
│   ├── Player transitions
│   ├── Loading states
│   └── Micro-interactions
└── 🔗 07 - Export / Dev Handoff
    ├── Component specs
    ├── CSS variables
    ├── Responsive breakpoints
    └── Asset exports
```

### 15.2 Responsive Breakpoints

| Breakpoint | Width | Layout |
|-----------|-------|--------|
| Mobile S | 320px | Single column, bottom nav |
| Mobile L | 414px | Single column, bottom nav |
| Tablet | 768px | Two columns, collapsible sidebar |
| Laptop | 1024px | Full sidebar, 3 columns |
| Desktop | 1440px | Maximum content width |
| Ultra | 1920+ | Expanded with larger art |

### 15.3 Key Figma Techniques

1. **Auto Layout** for all components (padding, gap, alignment)
2. **Variants** for component states (default, hover, active, disabled)
3. **Component Properties** for text overrides, icon swaps
4. **Local Variables** for colors, spacing, radius (maps to CSS tokens)
5. **Interactive Components** for hover/click/drag prototypes
6. **Smart Animate** for page transitions and player states
7. **Overlays** for modals, queue panel, context menus
8. **Scrollable containers** for long lists with overflow

### 15.4 Color Variable Mapping

```
Figma Variable              → CSS Custom Property
----------------------------→-----------------------
Color/Primary/Main          → var(--color-primary-main)
Color/Neutral/Gray-13       → var(--color-neutral-gray-13)
Color/Spotify               → var(--color-spotify)
Color/Aurora/Green          → var(--color-aurora-green)
Color/Aurora/Blue           → var(--color-aurora-blue)
Color/Aurora/Purple         → var(--color-aurora-purple)
Color/Aurora/Pink           → var(--color-aurora-pink)
Spacing/4                   → 0.25rem
Spacing/16                  → 1rem
Spacing/24                  → 1.5rem
Radius/Sm                   → var(--radius-sm)  [8px]
Radius/Md                   → var(--radius-md)  [12px]
Radius/Lg                   → var(--radius-lg)  [16px]
Radius/Xl                   → var(--radius-xl)  [24px]
```

---

## Appendix A: Player Controls Quick Reference

| Control | Desktop | Mobile | Gesture |
|---------|---------|--------|---------|
| Play/Pause | Space / Click | Tap center | Double tap art |
| Next | Ctrl+→ / Click | Tap next btn | Swipe right |
| Previous | Ctrl+← / Click | Tap prev btn | Swipe left |
| Volume ↑ | ↑ | — | Two-finger up |
| Volume ↓ | ↓ | — | Two-finger down |
| Seek → | → | — | Drag seekbar |
| Seek ← | ← | — | Drag seekbar |
| Fullscreen | F / Click | Swipe up | Pinch open |
| Like | Click heart | Tap heart | Swipe right (track) |
| Queue | Click queue btn | — | — |
| Mute | M / Click | — | — |
| Speed | Click speed btn | — | — |
| Sleep Timer | Click timer btn | — | — |
| Close Player | Escape / ↓ | Swipe down | Pinch close |

## Appendix B: Page Load Performance Targets

| Metric | Target |
|--------|--------|
| First Contentful Paint | < 1.5s |
| Largest Contentful Paint | < 2.5s |
| First Input Delay | < 100ms |
| Time to Interactive | < 3.5s |
| Page transition animation | 350ms |
| Player toggle (fullscreen) | 250ms |
| Search results appear | < 200ms |
| Track playback start | < 500ms |
| Lyrics load | < 1s |

## Appendix C: Accessibility Requirements

- All interactive elements must be keyboard navigable
- Focus indicators visible (green ring, 2px)
- Color contrast: AA minimum (4.5:1 for text)
- ARIA labels on all player controls
- Skip-to-content link on every page
- Reduce motion media query respects user preferences
- Screen reader announcements for track changes
- Volume controls accessible via keyboard
- Lyrics font size adjustable
- Closed captions for any video content
