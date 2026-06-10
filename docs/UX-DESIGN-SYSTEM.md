# UX Design System — Next-Generation Music Platform

> **Role**: Principal Product Designer (ex-Spotify, Apple Music, SoundCloud, Tidal)
> **Platform**: Web-first progressive SPA with mobile-responsive layers
> **Philosophy**: "The platform should feel like a living, breathing music organism — not a database with a player attached."

---

## 1. Information Architecture

```
Root
├── Home                          # / (discovery landing)
├── Search                        # /search
├── Library                       # /library
│   ├── Playlists                 # /library/playlists
│   ├── Albums                    # /library/albums
│   ├── Artists                   # /library/artists
│   ├── Liked Songs               # /library/tracks
│   └── History                   # /library/history
├── Discover                      # /discover
│   ├── Trending                  # /discover/trending
│   ├── Viral                     # /discover/viral
│   ├── Moods                     # /discover/moods
│   ├── Genres                    # /discover/genres
│   └── AI Mixes                  # /discover/ai-mixes
├── Playlist                      # /playlist/:id
├── Album                         # /album/:id
├── Artist                        # /artist/:id
├── Track                         # /track/:id
├── Notifications                 # /notifications
├── Profile                       # /profile/:id
├── Creator Studio                # /creator
│   ├── Dashboard                 # /creator
│   ├── Upload                    # /creator/upload
│   ├── Releases                  # /creator/releases
│   ├── Analytics                 # /creator/analytics
│   └── Settings                  # /creator/settings
├── Admin                         # /admin
│   ├── Dashboard                 # /admin
│   ├── Users                     # /admin/users
│   ├── Artists                   # /admin/artists
│   ├── Albums                    # /admin/albums
│   ├── Tracks                    # /admin/tracks
│   ├── Genres                    # /admin/genres
│   ├── Moderation                # /admin/moderation
│   ├── Reports                   # /admin/reports
│   └── Media                     # /admin/media
├── Settings                      # /settings
└── Analytics                     # /analytics
```

### Navigation Hierarchy

```
Primary Nav (Sidebar)
├── Home
├── Discover
├── Search
├── Library
│   ├── Playlists
│   ├── Liked Songs
│   └── History
└── Notifications [count]

Secondary Nav (Top bar)
├── Creator Studio [if artist]
├── Admin [if admin]
├── Settings
└── Profile

Persistent UI
├── NowPlayingBar (bottom)
├── FloatingMiniPlayer (above bar)
├── QueuePanel (right drawer)
└── FullscreenPlayer (overlay)
```

---

## 2. User Flow Diagrams

### Core Listening Flow
```
Land (Home) → Browse Section → Click Track
  → Track plays in NowPlayingBar
  → FloatingMiniPlayer appears
  → Continue browsing / Navigate elsewhere
  → Track continues (persistent playback)
  → [Optional] Click floating player → Fullscreen Player
```

### Exploration Flow
```
Discover Page → Browse section (Trending/Moods/Genres)
  → Hover card → Preview plays (5s)
  → Click card → Album/Playlist page
  → Click track → Adds to queue / Plays immediately
  → Follow artist → Adds to library
```

### Social Flow
```
User Profile → Activity Feed
  → See "User liked Track X"
  → Click → Listen to Track X
  → React / Comment
  → Share to feed
  → Follow user
```

### Creator Flow
```
Creator Studio Dashboard → View stats
  → Upload Track → Fill metadata → Publish
  → Track goes to moderation
  → Published → Appears on artist page
  → View real-time analytics
```

---

## 3. Page Hierarchy

### Home Page
```
┌─────────────────────────────────────────────────┐
│  TopBar: Search | Notifications | Profile       │
├─────────────────────────────────────────────────┤
│  Hero Carousel: Curated featured content        │
│  (Full-bleed, auto-playing video/gradient bg)   │
├─────────────────────────────────────────────────┤
│  "Good evening" greeting section                │
│  Recently Played (horizontal scroll row)        │
├─────────────────────────────────────────────────┤
│  "Made for You"                                 │
│  AI Mixes | Daily Mix | Weekly Discovery        │
├─────────────────────────────────────────────────┤
│  "Trending Now"                                 │
│  Track list with live listener count            │
├─────────────────────────────────────────────────┤
│  "New Releases"                                 │
│  Album cards in grid                            │
├─────────────────────────────────────────────────┤
│  "Recommended Playlists"                        │
│  Playlist cards in grid                         │
├─────────────────────────────────────────────────┤
│  "Popular Artists"                              │
│  Artist cards in horizontal scroll              │
├─────────────────────────────────────────────────┤
│  Footer: Links | Language | Legal              │
├─────────────────────────────────────────────────┤
│  NowPlayingBar (fixed bottom)                   │
│  FloatingMiniPlayer (above bar, draggable)     │
└─────────────────────────────────────────────────┘
```

### Fullscreen Player
```
┌─────────────────────────────────────────────────┐
│  ╳ Close      ░░░░ Progress ░░░░     Queue ☰   │
├─────────────────────────────────────────────────┤
│                                                 │
│           ┌─────────────────┐                    │
│           │                 │                    │
│           │   Album Art     │                    │
│           │   (animated /   │                    │
│           │   spectrum)     │                    │
│           │                 │                    │
│           └─────────────────┘                    │
│                                                 │
│   Track Title (large, dynamic color)            │
│   Artist Name → link                            │
│   Album Name → link                             │
│                                                 │
│   ♥ Like   🔄 Share   ➕ Add to Library         │
│                                                 │
│   ─────────────────●────────────────────        │
│   1:23                           3:45            │
│                                                 │
│   ⏮   ⏸   ⏭                              │
│   ↺ Shuffle   ↻ Repeat                         │
│                                                 │
│   🔊 ──────●─── Volume                          │
│                                                 │
│   Lyrics panel (scroll sync)                    │
│   │ Text line 1 (highlighted)                  │
│   │ Text line 2                                │
│   │ Text line 3 (dim, upcoming)                │
│                                                 │
├─────────────────────────────────────────────────┤
│  Bottom: Related tracks | Queue peek            │
└─────────────────────────────────────────────────┘
```

---

## 4. Wireframes (Text)

### FloatingMiniPlayer — 3 Modes

**Mini (64×64)**: Small album art square, play/pause overlay on tap. Draggable.

**Compact (288×80)**: Small album art + 1-line title + 1-line artist + minimal progress bar + prev/play/next + close. Draggable.

**Expanded (320×auto)**: 96×96 album art + title + artist + progress bar with time + full controls + pop-out + fullscreen buttons. Draggable.

### PiP Window (380×220)
Title: "Now Playing"
Full viewport fill of Compact/Expanded mode, no border-radius, expanded controls. The PiP composable teleports the component tree into the pip document, preserving Pinia reactivity.

---

## 5. Component System

### Atomic Design Levels

**Atoms**
- Button (primary, secondary, ghost, icon)
- Input (text, search, range, select)
- Typography (heading, body, caption, label, number)
- Icon (lucide-vue-next wrappers)
- Avatar (user, artist, placeholder)
- Badge (verified, new, exclusive, explicit)
- Progress bar (linear, circular, stepped)
- Slider (volume, seek, crossfade)

**Molecules**
- TrackRow (cover + title + artist + album + duration + actions)
- AlbumCard (cover + title + artist + year)
- ArtistCard (avatar + name + followers + genre tags)
- PlaylistCard (cover + title + track count + owner)
- SearchResultGroup (section header + results list)
- ControlButtonGroup (play/pause + next + prev)
- VolumeControl (icon + slider)
- LikeButton (heart toggle with animation)
- ShareButton (share sheet trigger)
- ProgressDisplay (seekable bar + time labels)

**Organisms**
- NowPlayingBar (full-width bottom bar with track info + controls + volume)
- FloatingMiniPlayerContent (draggable player in 3 modes)
- FloatingMiniPlayerHost (PiP orchestrator)
- FullscreenPlayer (immersive overlay)
- QueuePanel (right-side drawer with current queue + recommendations)
- LyricsDisplay (scroll-synced lyrics with highlight)
- SearchOverlay (full-screen modal search)
- MusicSidebar (primary navigation)
- HomeHero (featured content carousel)
- SectionRow (section header + "Show all" + horizontal scroll)
- TrackList (sortable, filterable track table)
- ArtistHero (image + name + stats + actions)

**Templates**
- LayoutMusicApp (sidebar + topbar + content + player stack)
- LayoutAdmin (compact sidebar + topbar + content)
- LayoutAuth (centered card on gradient background)

**Pages**
- See Section 3 for full page hierarchy.

---

## 6. Design System

### Color Architecture

```
Background (base):        #0a0a0a
Surface (card/panel):     #121212
Surface Elevated:         #1a1a1a / #1e1e1e
Surface Overlay:          rgba(255,255,255,0.05-0.15)

Primary (Spotify green):  #1db954 → #1ed760 (hover)
Primary Subtle:           rgba(29,185,84,0.15)

Accent (Apple-style):     Gradient(#1db954 → #4ade80)
Aurora (ambient):         Gradient(#1db954 → #6366f1 → #a855f7)

Text Primary:             #ffffff
Text Secondary:           rgba(255,255,255,0.6-0.8)
Text Tertiary:            rgba(255,255,255,0.3-0.4)

Border:                   rgba(255,255,255,0.06-0.10)
Border Hover:             rgba(255,255,255,0.15-0.20)

Error:                    #ef4444
Warning:                  #f59e0b
Success:                  #22c55e

Dynamic Colors (per-album):
  Extracted from album art dominant palette
  Applied to: backgrounds, gradients, glow effects
```

### Glassmorphism Tokens

```
.glass-light {
  background: rgba(255,255,255,0.05);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255,255,255,0.08);
}

.glass-heavy {
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(24px);
  border: 1px solid rgba(255,255,255,0.06);
}

.glass-player {
  background: rgba(18,18,18,0.85);
  backdrop-filter: blur(32px) saturate(1.2);
  border-top: 1px solid rgba(255,255,255,0.06);
  box-shadow: 0 -20px 60px rgba(0,0,0,0.65);
}
```

### Typography Stack

```
Family: IranYekan (primary), system-ui (fallback)

Display/Headline: 700 (Black), 2rem–4rem
Heading 1:        800 (ExtraBold), 1.5rem–2rem
Heading 2:        700 (Bold), 1.25rem–1.5rem
Heading 3:        700 (Bold), 1rem–1.125rem
Body:             400 (Regular), 0.875rem
Body Small:       400 (Regular), 0.75rem–0.8125rem
Caption:          400 (Regular), 0.6875rem
Label:            700 (Bold), 0.6875rem–0.75rem (uppercase, tracked)
Number:           500 (Medium), 0.75rem (tabular-nums, monospace)
```

### Spacing Scale

```
Base unit: 4px
Scale: 0, 4, 8, 12, 16, 20, 24, 32, 40, 48, 56, 64, 80, 96, 128
```

### Border Radius

```
none:    0px
sm:      6px
md:      10px
lg:      16px
xl:      20px
2xl:     24px
full:    9999px
```

### Shadows

```
sm:   0 2px 8px rgba(0,0,0,0.3)
md:   0 4px 16px rgba(0,0,0,0.4)
lg:   0 8px 32px rgba(0,0,0,0.5)
xl:   0 16px 48px rgba(0,0,0,0.6)
glow: 0 0 30px rgba(29,185,84,0.3)     (Spotify green glow)
glow-aurora: 0 0 60px rgba(99,102,241,0.2) (purple ambient glow)
```

### Z-Index Layers

```
sidebar:      40
topbar:       30
overlay:      50
modal:        60
toast:        70
player-bar:   100
floating:     150
pip-window:   2147483647  (browser maximum)
```

---

## 7. Motion System

### Timing

```
micro (hover/down):    100ms ease
fast (toggle/close):   160ms ease
normal (enter/leave):  250ms ease-out
slow (page transition):400ms cubic-bezier(0.34, 1.56, 0.64, 1)
hero (carousel):       600ms cubic-bezier(0.16, 1, 0.3, 1)
```

### Motion Patterns

```
.fade-enter-active, .fade-leave-active {
  transition: opacity 160ms ease;
}

.slide-up-enter-active {
  transition: all 250ms cubic-bezier(0.34, 1.56, 0.64, 1);
}
.slide-up-leave-active {
  transition: all 200ms ease-in;
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(16px) scale(0.95);
}

.scale-in-enter-active {
  transition: all 200ms cubic-bezier(0.34, 1.56, 0.64, 1);
}
.scale-in-leave-active {
  transition: all 150ms ease-in;
}
.scale-in-enter-from,
.scale-in-leave-to {
  opacity: 0;
  transform: scale(0.92);
}
```

### Micro-interactions

- **Hover card**: Slight lift (translateY(-2px)) + shadow increase + 250ms
- **Play button**: Scale pulse to 1.1 on press, back to 1 on release
- **Like toggle**: Heart icon scales 1→1.3→1 with color transition (300ms spring)
- **Progress thumb**: Opacity 0→1 on hover, scale 1→1.2 on active
- **Page transition**: Content fades + slides up 12px (200ms)
- **NowPlayingBar art spin**: 8s linear infinite when playing (like a vinyl)
- **Equalizer bars**: 600ms ease-in-out infinite alternate wave animation
- **Album art pop**: 250ms ease-out scale(0.92→1) on track change
- **PiP enter**: Browser-native PiP transition (smooth zoom-out effect)
- **Toast**: Slides in from top-right, auto-dismisses after 3-4s with fade

---

## 8. Player UX Specifications

### NowPlayingBar (Bottom Bar)
```
Height: 64px (mobile), 72px (desktop)
Layout: 3-column grid
  Left:   Cover (48-56px) + Title + Artist + loading state
  Center: Controls row (shuffle, prev, play, next, repeat)
          + Seekbar with time labels (max-width 480px)
  Right:  Fullscreen btn + Queue btn + Volume slider (120-160px)
States: Default, Active (track playing), Empty (no track)
Behavior:
  - Always visible when a track is loaded
  - Progress bar at top edge (4px, green gradient)
  - Controls centered, responsive
  - Volume persists via localStorage
```

### FloatingMiniPlayer (Spotify-style)
```
Position: Fixed, bottom: 80px (above NowPlayingBar), right: 16px
Z-index:  150
Modes:
  - Mini (64×64): Album art square with play/pause overlay
  - Compact (288×80): Art + 1-line info + minimal controls + progress
  - Expanded (320×auto): Full art + full info + progress + expanded controls
Drag:      Mouse + touch, clamped to viewport bounds
PiP:       Document Picture-in-Picture API, teleports component tree
           Sizing: 380-400×220-240px
           Fallback: Toast notification when unsupported
Transitions: Float-up spring animation (enter/leave)
```

### Fullscreen Player
```
Trigger:    Click floating player / NowPlayingBar expand btn
Overlay:    Full viewport, z-index 200
Background: Dynamic gradient from album art dominant color
            + blurred album art as texture layer (opacity 0.15)
Layout:     2-column (player left, lyrics/queue right)
  Left pane (55%):
    - Album art: 280-320px, animated rotation when playing
    - Optional: Audio spectrum visualizer overlay
    - Track title: Large dynamic-color heading
    - Artist / Album: Clickable links
    - Action row: Like, Share, Add to Library
    - Seekbar: Full-width with time + duration
    - Controls: Shuffle, Prev, Play/Pause, Next, Repeat
      - Play button: 56px circle, white bg, hover -> green
    - Volume: Slider with mute toggle
    - Extended: Sleep timer, Quality, Speed
    
  Right pane (45%):
    - Tabbed: Lyrics | Queue | Related
    - Lyrics tab: Synced scrolling, current line highlighted
    - Queue tab: Up next + recommendations to add
    - Related tab: "Fans also like" tracks + artist discography
    
Controls actions:
  - Play/Pause: Space keybinding
  - Next/Prev: Arrow Right/Left
  - Seek: Click on bar or Arrow Up/Down (skip 5s)
  - Volume: Scroll wheel over player
  - Like: Heart toggle (sends to backend)
  - Share: Copy link / Share sheet
```

### Media Session Integration
```
navigator.mediaSession.metadata = {
  title, artist, album, artwork
}
navigator.mediaSession.setActionHandler('play', resume)
navigator.mediaSession.setActionHandler('pause', pause)
navigator.mediaSession.setActionHandler('next', playNext)
navigator.mediaSession.setActionHandler('previous', playPrevious)
navigator.mediaSession.setActionHandler('seekto', seek)
```

---

## 9. Lyrics UX Specifications

### Default View (Album Page / Track Page)
```
Container: Max-width 680px, centered
Current line: Large (1.5rem), white, bold
Past lines: Dimmed, smaller (0.875rem), above
Future lines: Dim, smaller (0.875rem), below
Scroll: Auto-syncs with playback position
        100ms offset for anticipation
Transition: Current line smoothly fades in at 1.2x scale
```

### Fullscreen Karaoke Mode
```
Full viewport, dark background
Current line: Very large (2.5-3rem), bold, centered
Word-level highlight: Each word lights up in sync
  - Uses word-level timestamps from LRC or API
  - Highlight color: Primary green gradient
Past words: Dimmed white
Future words: Very dim
Background: Abstract particles / waveform reacting to audio
Translation: Floating panel (bottom-right), synced current line
Annotations: Click a line → opens annotation panel
  - Artist commentary on this line
  - Community annotations (reddit-style upvoted)
  - "What does this mean?" explainer
```

### Non-Synced Lyrics
```
Display as scrolling text block
Current position: Auto-scrolls to keep roughly centered
Can still scroll manually (override auto-scroll)
Show timestamps on each line when available
```

### Lyrics States
```
Loading: Skeleton shimmer (3 lines)
Empty: "No lyrics available" + contribute CTA
Offline: Cached lyrics from last view
Error: "Lyrics unavailable" in dim text
End of track: "Thank you for listening" with repeat button
```

---

## 10. Mobile UX Specifications

### Bottom Sheet Player
```
Trigger: Swipe up on NowPlayingBar
State: Sheet at 30% height (controls) / 60% (lyrics peek) / 90% (full)
Content:
  - Handle bar at top (drag indicator)
  - Album art: 120px circle (30%), 180px (60%)
  - Title + Artist
  - Seekbar + time
  - Controls (prev, play, pause, next)
  - Lyrics preview (60%+)
  - Queue button
  - Volume (hidden behind expand)
Gestures:
  - Swipe up: Expand sheet
  - Swipe down: Collapse to bar
  - Swipe left on track: Queue next / Add to playlist sheet
  - Swipe right on track: Like
  - Tap art: Go to album
```

### Mobile Fullscreen Player
```
Trigger: Rotate to landscape / Tap expand in sheet
Full viewport, landscape orientation
Album art: Large (60% width), left
Controls: Right side, vertically stacked
Volume: Hardware buttons preferred, on-screen fallback
Lyrics: Overlay on right half
```

### Gesture Navigation (Touch)
```
Swipe left: Queue this track
Swipe right: Like track (with haptic feedback)
Swipe up on track row: Show context menu
Swipe down on NowPlayingBar: Collapse to mini player
Long press album art: Share
Double tap album art: Toggle fullscreen
Pinch on album art: Enter immersive mode
```

### Responsive Breakpoints
```
Mobile:    < 640px
Tablet:    640px - 1024px
Desktop:   > 1024px
Wide:      > 1440px

Layout shifts:
  <640:  Single column, bottom sheet player, floating mini hidden
  640+:  Two-column, NowPlayingBar visible, floating mini on right
  1024+: Full sidebar + content + player bar
  1440+: Max-width content container (1400px)
```

---

## 11. Admin UX Specifications

### Layout
```
Sidebar: Compact (64px icons) or full (240px with labels)
  Sections: Dashboard | Catalog | Users | Reports | Moderation | Media
Content: Full remaining width
TopBar: Breadcrumb + Search + "View Site" + Quick actions
```

### Dashboard
```
4 KPI cards row:
  Total Users | Active Artists | Total Streams | Revenue
Stacked line chart: Streams over time (7d, 30d, 90d, 1y)
Bar chart: Top 5 genres by streams
Table: Recent uploads needing review
List: Pending moderation items (with quick-action buttons)
```

### Catalog Management
```
Data table with:
  - Advanced multi-select filters (status, genre, date range, etc.)
  - Sortable columns
  - Bulk actions (publish, archive, delete, assign genre)
  - Inline editing for status/title
  - Pagination (50 per page) + "Load more"
  - Export CSV
Row actions: Edit, Preview, Delete, Feature
Create flow: Side panel form (not separate page)
```

### Moderation Center
```
Queue: Oldest first, priority flagged items at top
Moderation card:
  - Track/Album metadata preview
  - Audio player (quick preview)
  - Cover art
  - Flags/reports count
  - Action buttons: Approve, Reject (with reason), Flag for review
Statistics: Average review time, queue size, actions taken today
```

### Design Constraints
```
- Always maintain dark theme (matching main platform)
- Different accent: Indigo/blue instead of green
  - Primary: #6366f1 (Indigo)
  - Charts: Multi-color palette (not green-only)
- Denser layouts (more info per screen)
- Table-heavy where appropriate
- Batch actions clearly separated
- Toast confirmations for all mutations
```

---

## 12. Creator Studio UX Specifications

### Dashboard
```
Header: "Welcome back, [Artist Name]" + date range selector
First row (4 stat cards, auto-animated on mount):
  Streams (total)
  Listeners (unique)
  Followers
  Revenue (if monetized)
Second row:
  Line chart: Streams over time (daily granularity)
  Bar chart: Top tracks by streams
Third row:
  "Latest Release" card + "Upload New" CTA
  "Recent Activity" feed (new followers, playlist adds, comments)
Empty state: "Upload your first track to start seeing analytics"
```

### Upload Flow
```
Multi-step wizard:

Step 1 - Audio:
  - Drop zone (drag/drop audio file)
  - File validation (format, size, duration)
  - Upload progress bar
  - Auto-tagging (metadata extraction)

Step 2 - Details:
  - Title (required)
  - Artist name (auto-filled)
  - Featuring artists (multi-select, searchable)
  - Album (new/existing dropdown)
  - Genre(s) (multi-select)
  - Mood tags
  - Explicit content toggle
  - Lyrics (textarea or .lrc upload)
  - Credits (write, produce, mix, master)

Step 3 - Artwork:
  - Image drop zone (min 1400×1400, max 3000×3000)
  - Auto-crop tool
  - Preview (track listing mockup + fullscreen player mockup)

Step 4 - Review:
  - Summary card
  - "Publish now" / "Schedule release" (date picker)
  - "Save as draft"

Steps: Progress indicator at top (4 circles with labels)
State: Save draft at any step (auto-save every 30s)
```

### Artist Verification Flow
```
Settings → Verification tab
Upload: Government ID, Social links, Proof of ownership
Status: Pending / Verified / Rejected
Badge: Blue checkmark (displayed on artist page, search, comments)
```

### Analytics
```
Tabbed sub-nav: Overview | Tracks | Listeners | Revenue

Tracks tab:
  Sortable table: Track, Streams, Listeners, Saves, Skips, Completion rate
  Click row → Detailed track analytics page

Listeners tab:
  Geo map (heatmap of listener locations)
  Top cities table
  Age/gender demographics (pie chart)
  Listening time distribution (hourly heatmap)
```

---

## 13. Social UX Specifications

### Activity Feed
```
Timeline layout (reverse chronological)
Each activity is a card:
  - Avatar + Name + Action + Timestamp
  - Context preview (album art, track card, playlist card)
  - Action buttons (Listen, Like, Reply, Share)
  - Like count + Comments preview
Activity types:
  - "Listened to [Track]" 
  - "Liked [Track]"
  - "Added [Track] to playlist [Name]"
  - "Followed [User]"
  - "Created playlist [Name]"
  - "Commented on [Track]"
  - "Shared [Track]"
```

### Follow System
```
Follow/Unfollow: Toggle button (pill shape)
  States: "Follow" (outline) → "Following" (filled with checkmark)
Notifications: 
  - "X started following you"
  - "X liked your track"
  - Email/push digest (optional)
Mutual follows: Special badge on profile
```

### Reactions & Comments
```
Reactions: Like (heart), Fire, Clap, Mind-blown, Sad
  - Mouse hover: Shows expanded reaction picker (5 options)
  - Tap/hold (mobile): Opens reaction sheet
  - Count shown as horizontal icon row

Comments:
  - Threaded (top-level + replies)
  - Rich text (mentions @username, #tags)
  - Timestamps
  - Like on comments
  - Delete own comment
  - Report inappropriate (3-dot menu)
```

### Music Sharing
```
Share sheet (triggered from any track/album/playlist):
  - Copy link (auto-format as rich preview)
  - Share to platform feed
  - Share to social (Twitter, WhatsApp, Telegram)
  - Generate share card (beautiful image with album art + QR)
  - Embed code (for websites/blogs)
```

### Listening Together (Future)
```
Real-time sync:
  - Create room → Share invite link
  - All participants hear same position
  - Chat panel alongside player
  - Emoji reactions fly across screen
  - Host controls play/pause/queue
Avatar stack: Show participant avatars in circle
```

---

## 14. Discovery UX Specifications

### Sections (Home & Discover)

```
Hero: Full-width auto-playing video/gradient carousel
  - Featured artist/album of the moment
  - Click → Album/Artist page

"Made for You" (AI powered):
  - Daily Mix 1, 2, 3 (algorithmic)
  - Weekly Discovery (fresh, never-heard-before)
  - AI Mix Generator (describe mood → get playlist)
  - Each card: Cover + title + short description (why this?)

"Trending Now":
  - Tracks sorted by velocity (recent spike in plays)
  - Each row: Position + Cover + Title + Artist + daily gain %
  - Live listener count badge

"New Releases":
  - Grid of album cards (most recent first)
  - Filter: This week, This month, All
  - Badge: "New" (fades after 7 days)

"Moods & Moments":
  - Mood grid (2×3 or 3×2 cards with gradient/ambient background)
  - Moods: Focus, Chill, Energy, Sad, Happy, Workout, Late Night, Road Trip
  - Each mood → AI-curated infinite mix

"Genre Worlds":
  - Genre cards in grid (each with distinct color/gradient)
  - Click → Genre page: Top tracks, top artists, new releases, description

"Community Favorites":
  - Most liked this week
  - Most shared this week
  - Top commented tracks

"Because you listened to [Track]":
  - Similar tracks (audio features + user behavior)
  - Similar artists
  - "Fans also like" (collaborative filtering)

"Viral":
  - Tracks with highest share/listen ratio
  - Social proof: "X people shared this today"
  - Platform-exclusive viral chart
```

### Search UX
```
Full-screen overlay on trigger (Ctrl+K or click)
Focused input at top, auto-focused
Instant results (debounced 150ms)

Sections (live as you type):
  1. Top result (best match across all)
  2. Artists (max 3)
  3. Tracks (max 5)
  4. Albums (max 3)
  5. Playlists (max 3)
  6. Users (max 3)

Empty state: Past searches, trending searches, browse genres
Keyboard navigation: Arrow keys + Enter, Escape to close
Each result: Cover/avatar + title + subtitle + type badge
```

---

## 15. Visual Descriptions (Figma-ready)

### Home Page — Visual Design
```
A dark, infinite canvas. The page is layered with translucent
glass panels floating over a subtle animated gradient background.

The hero is a full-bleed carousel with a left-aligned text overlay.
Behind each slide, a cinematic image or looped video plays at 30%
opacity. The active slide has a soft glow ring around the cover art.

Section rows are separated by 8px of negative space. Each row has a
small uppercase section label (12px, tracked 0.15em, #1db954) with
a "Show all" link on the right. Content scrolls horizontally with
overflow hidden, revealing a soft gradient fade on the right edge
to indicate scrollability.

Cards have a 10px radius, dark surface background (#121212) with
a subtle 1px white/5% border. On hover, they lift 2px with an
increased shadow. The cover art inside is flush to the card edges.

The NowPlayingBar at the bottom is a thick glass panel with a
4px progress line across its top edge. The album art spins slowly
for the current track. Equalizer bars pulse on the right edge.

The FloatingMiniPlayer sits in the bottom-right corner, a compact
glass orb with the current album art, pulsing gently.

It's midnight. The room is dark. The screen glows.
```

### Fullscreen Player — Visual Design
```
A cathedral of light and sound. The background is a massive,
gently shifting gradient extracted from the current album's
dominant colors, with a subtle layer of the album art blurred
to 80px and reduced to 15% opacity.

The album art sits center-left, 300px square, with rounded
corners (16px). It rotates at 4s per revolution, a god ray
sweeping across its surface. When paused, it stops instantly.

Track title floats below, set in ExtraBold 28px, colored with
a gradient from white to 70% white. Artist name beneath,
18px, linking to their page.

Controls are minimal, floating, spaced with 40px gaps. The
play button is a 56px white circle with a black play icon.
Hovering over it turns the circle into Spotify green.

The seekbar is 2px tall with a 14px circular thumb. On hover,
the thumb grows to 18px and glows green.

On the right, lyrics scroll. Each line is 16px, dim white.
The current line is 22px, bold, pure white, with a subtle
green glow on the left edge.

When the audio spectrum visualizer is enabled, soft bars
pulse behind the album art, reacting to frequency bands.
Low frequencies are wider, higher narrower.

The experience is visceral. You feel the music in the pixels.
```

### Discover Page — Visual Design
```
A living magazine of sound. The discover page is designed as
a scrolling feed of musical discovery, each section a new
visual treat.

The hero is a panoramic gradient with the featured artist's
colors, a large typographic overlay reading the artist name
in 64px bold. A subtle parallax effect as you scroll.

"Moods" are presented as full-bleed horizontal cards, each
240px wide, 160px tall. Each has a distinct gradient:
  Focus: Cool blues (#0f172a → #1e293b)
  Energy: Hot reds (#7f1d1d → #991b1b)  
  Chill: Sunset oranges (#78350f → #92400e)
  Sad: Muted purples (#3b0764 → #581c87)
  Happy: Bright yellows (#713f12 → #854d0e)

Genre cards are larger (280×200) with the genre name
overlaid in bold white text, a subtle texture pattern
behind, and the representative accent color as a
gradient overlay.

"Trending" shows an ordered list with rank numbers in
tabular-nums, a bright green upward arrow for trending
up, and a small sparkline chart for each track showing
velocity over the last 24h.

Loading states show shimmering skeletons with the same
card dimensions. Empty states show a music note icon
in 30% white with "Nothing here yet" in secondary text.
```

### Lyric View — Visual Design
```
Pure typographic magic. The lyrics view transforms text into
a visual experience.

For synced lyrics, each line is its own moment. The current
line floats in the center of the screen, 2.5rem, bold.
Behind it, the album art blooms as an enormous blurred
background (100px blur, 20% opacity), shifting hues with
each new track.

Words highlight in sequence: a green gradient sweeps across
each syllable in perfect time with the vocal. The highlight
isn't a simple color change — it's a smooth gradient that
shifts from dim to bright, creating a miniature sunrise
on each word.

Past lines fade upward, shrinking to 0.875rem and 30% opacity.
Future lines fade in from below, 0.875rem, 20% opacity.
The overall effect is a vertical river of text, with the
current line being a bright stepping stone.

In karaoke mode, the background becomes a particle field.
Thousands of tiny luminescent dots pulse and drift in sync
with the audio's energy — brighter during chorus, calmer
during verses. The particles react to bass hits with
explosive radial bursts.

When translations are enabled, a subtle panel slides in
from the bottom-right showing the translated current line
in a smaller, italicized secondary typeface.

It's not just reading lyrics. It's experiencing the song
through its own words, visualized.
```

### Mobile — Visual Design
```
Thumb-first. Everything is within reach of a comfortable grip.

The bottom sheet player is a smooth glass panel that peeks
up 60px above the bottom edge. The album art is a 120px
circle at the top of the sheet. Swiping up reveals controls,
then lyrics.

Gestures are fluid. Swiping left on a track in the queue
shows a "Add to playlist" action with a satisfying rubber-band
bounce. Swiping right triggers a heart animation, the icon
scales up, flashes, and settles into a filled state.

The fullscreen player in landscape is a split view: album
art on the left (60% width), controls on the right (40%).
Volume overlay appears as a vertical slider on the right
edge when swiped.

Navigation uses a bottom tab bar (4 tabs) with icons only,
the active tab highlighted in green. The bar is glass,
blurring the content behind it.

Everything is designed for one-handed use. Key actions are
in the lower half of the screen. The top half is content.
This is mobile-first, thumb-optimized.

Music in your pocket. The world at your thumb.
```

---

## Design Principles

1. **Music First, UI Second** — Every pixel either serves the music or gets out of the way
2. **Dark is the Canvas** — Dark backgrounds make album art and lyrics shine
3. **Glass, Not Solid** — Translucency creates depth without weight
4. **Motion with Meaning** — Every animation has a reason (feedback, spatial orientation, delight)
5. **Typography is Texture** — Words don't just convey meaning, they create visual rhythm
6. **Context is King** — Show the right information at the right moment, hide everything else
7. **Playback is Sacred** — Never interrupt, never reload, always remember position
8. **Social Without Noise** — Social features enhance discovery, not distract from listening
9. **Snap, Don't Lag** — Every interaction under 100ms, optimistic updates everywhere
10. **Premium Feeling** — The UI should feel expensive, like a high-end audio device brought to software

---

> "Design is not just what it looks like and feels like. Design is how it works."
> — This platform is designed for people who live and breathe music.
> Every decision serves the moment when a song hits you and nothing else matters.
