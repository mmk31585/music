# Music Streaming Platform - Architecture Overview

## Project Structure

```
/workspace
├── backend/ (Go)
│   ├── cmd/api/                 # Application entry point
│   ├── internal/
│   │   ├── app/                 # App initialization, routes
│   │   ├── common/              # Shared utilities
│   │   │   ├── errors/          # Error handling
│   │   │   ├── middleware/      # Global middleware
│   │   │   ├── pagination/      # Pagination helpers
│   │   │   ├── request/         # Request parsing
│   │   │   ├── response/        # Response formatting
│   │   │   └── validator/       # Validation
│   │   ├── config/              # Configuration management
│   │   ├── domain/              # Domain entities
│   │   ├── handler/             # HTTP handlers
│   │   ├── middleware/          # Custom middleware
│   │   ├── modules/             # Feature modules
│   │   │   ├── auth/            # Authentication & Authorization
│   │   │   ├── catalog/         # Music catalog (tracks, artists, albums, genres)
│   │   │   ├── playlist/        # Playlist management
│   │   │   ├── library/         # User library (liked, saved, followed)
│   │   │   ├── history/         # Listening history
│   │   │   ├── queue/           # Playback queue
│   │   │   ├── search/          # Search & discovery
│   │   │   ├── recommendation/  # AI recommendations
│   │   │   ├── social/          # Social features (follow, friends, activity)
│   │   │   ├── podcast/         # Podcasts & audiobooks
│   │   │   ├── download/        # Offline downloads
│   │   │   ├── audio/           # Audio processing (EQ, crossfade)
│   │   │   ├── premium/         # Subscription & billing
│   │   │   ├── artist/          # Artist tools & analytics
│   │   │   └── admin/           # Admin dashboard
│   │   ├── repository/          # Data access layer
│   │   ├── service/             # Business logic
│   │   └── platform/
│   │       ├── cache/           # Redis cache
│   │       ├── database/        # PostgreSQL connection & migrations
│   │       └── logger/          # Logging
│   └── migrations/              # SQL migrations
│
├── frontend/ (Vue 3 + TypeScript)
│   ├── src/
│   │   ├── assets/              # Static assets
│   │   ├── components/
│   │   │   ├── admin/           # Admin components
│   │   │   ├── auth/            # Auth components
│   │   │   ├── common/          # Reusable UI components
│   │   │   ├── forms/           # Form components
│   │   │   ├── layouts/         # Layout components
│   │   │   ├── music/           # Music-specific components
│   │   │   │   ├── TrackRow.vue
│   │   │   │   ├── TrackCard.vue
│   │   │   │   ├── AlbumCard.vue
│   │   │   │   ├── ArtistCard.vue
│   │   │   │   ├── PlaylistCard.vue
│   │   │   │   ├── NowPlayingBar.vue
│   │   │   │   ├── MusicSidebar.vue
│   │   │   │   ├── MusicTopbar.vue
│   │   │   │   ├── PlayerControls.vue
│   │   │   │   ├── QueuePanel.vue
│   │   │   │   ├── Equalizer.vue
│   │   │   │   └── Visualizer.vue
│   │   │   ├── playlist/        # Playlist components
│   │   │   ├── library/         # Library components
│   │   │   ├── search/          # Search components
│   │   │   ├── social/          # Social components
│   │   │   ├── podcast/         # Podcast components
│   │   │   └── premium/         # Premium components
│   │   ├── composables/         # Composable functions
│   │   │   ├── auth/            # Auth composables
│   │   │   ├── admin/           # Admin composables
│   │   │   ├── catalog/         # Catalog composables
│   │   │   ├── playlist/        # Playlist composables
│   │   │   ├── library/         # Library composables
│   │   │   ├── search/          # Search composables
│   │   │   ├── social/          # Social composables
│   │   │   └── usePlayer.ts     # Player composable
│   │   ├── layouts/             # Page layouts
│   │   ├── pages/               # Page components
│   │   │   ├── app/             # Main app pages
│   │   │   │   ├── PageHome.vue
│   │   │   │   ├── PageSearch.vue
│   │   │   │   ├── PageLibrary.vue
│   │   │   │   ├── PagePlaylist.vue
│   │   │   │   ├── PageArtist.vue
│   │   │   │   ├── PageAlbum.vue
│   │   │   │   ├── PageTrack.vue
│   │   │   │   ├── PageGenre.vue
│   │   │   │   ├── PagePodcasts.vue
│   │   │   │   └── PagePremium.vue
│   │   │   ├── auth/            # Auth pages
│   │   │   ├── admin/           # Admin pages
│   │   │   └── errors/          # Error pages
│   │   ├── plugins/             # Vue plugins
│   │   │   ├── client/          # API client
│   │   │   └── query-builder/   # Query builder
│   │   ├── router/              # Vue Router config
│   │   │   ├── middleware/      # Route guards
│   │   │   └── routes/          # Route definitions
│   │   ├── services/
│   │   │   ├── api/             # API service modules
│   │   │   │   ├── auth/
│   │   │   │   ├── catalog/
│   │   │   │   ├── playlist/
│   │   │   │   ├── library/
│   │   │   │   ├── search/
│   │   │   │   ├── social/
│   │   │   │   ├── podcast/
│   │   │   │   ├── premium/
│   │   │   │   └── common/
│   │   │   ├── socket/          # WebSocket service
│   │   │   ├── storage/         # Local storage service
│   │   │   └── analytics/       # Analytics service
│   │   ├── stores/              # Pinia stores
│   │   │   ├── index.ts
│   │   │   ├── user-auth.ts
│   │   │   ├── player.ts
│   │   │   ├── queue.ts
│   │   │   ├── playlist.ts
│   │   │   ├── library.ts
│   │   │   ├── social.ts
│   │   │   ├── download.ts
│   │   │   └── page-loader.ts
│   │   ├── utils/               # Utility functions
│   │   ├── App.vue
│   │   └── main.ts
│   └── package.json
│
└── migrations/                  # Database migrations
```

## API Routes Structure

### Public Routes (`/api/v1`)
```
/catalog
  /artists
  /albums
  /tracks
  /genres
  /search

/playlist
  /public/:id

/podcast
  /shows
  /episodes
```

### Authenticated Routes (`/api/v1`)
```
/me
  /profile
  /library
    /tracks
    /albums
    /artists
  /playlists
  /history
  /queue
  /downloads
  /following
  /followers

/playlist
  / (create)
  /:id
    / (get, update, delete)
    /tracks
    /collaborators
    /folder

/library
  /tracks/:id (like/unlike)
  /albums/:id (save/unsave)
  /artists/:id (follow/unfollow)

/history
  / (get history)
  /clear

/queue
  / (get/set queue)
  /next
  /clear

/search
  / (search all)
  /tracks
  /artists
  /albums
  /playlists
  /podcasts

/recommendation
  /daily-mix
  /discover-weekly
  /release-radar
  /based-on/:trackId
  /mood/:mood

/social
  /friends
  /activity
  /blend
  /session

/podcast
  /subscriptions
  /history

/download
  / (list downloads)
  /track/:id
  /playlist/:id
  /sync

/premium
  /subscribe
  /status
  /billing

/artist
  /analytics
  /profile
  /tour
  /merch

/admin
  /users
  /catalog
  /playlists
  /reports
  /analytics
```

## Database Schema

### Core Tables
- users
- artists
- albums
- tracks
- genres
- track_genres (junction)

### Playlist Tables
- playlists
- playlist_tracks (junction)
- playlist_collaborators (junction)
- playlist_folders

### Library Tables
- liked_tracks (user_track junction)
- saved_albums (user_album junction)
- followed_artists (user_artist junction)

### History & Queue
- listening_history
- playback_queue

### Social Tables
- follows (user-user)
- friend_activity
- blend_playlists
- group_sessions

### Podcast Tables
- podcast_shows
- podcast_episodes
- podcast_subscriptions
- podcast_history

### Download Tables
- downloads
- download_queue

### Premium Tables
- subscriptions
- payments
- subscription_features

### Analytics Tables
- track_plays
- user_listening_stats
- artist_analytics
- daily_charts
- trending_tracks

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL 15+
- **Cache**: Redis
- **Search**: PostgreSQL full-text search (can extend to Elasticsearch)
- **Queue**: Redis Streams / NATS
- **Storage**: Local/S3 for media files
- **Auth**: JWT with refresh tokens

### Frontend
- **Framework**: Vue 3.4+
- **Language**: TypeScript 5+
- **State**: Pinia
- **Router**: Vue Router 4
- **HTTP**: Axios
- **Validation**: Zod
- **Build**: Vite
- **UI**: PrimeVue + TailwindCSS
- **Audio**: Web Audio API

## Key Features Implementation Strategy

### Phase 1: Core Enhancement (Week 1-2)
- Enhanced database schema
- Playlist CRUD
- Library management
- Listening history

### Phase 2: Search & Discovery (Week 3)
- Full-text search
- Filters & sorting
- Charts & trending

### Phase 3: Recommendations (Week 4)
- Collaborative filtering
- Content-based recommendations
- Auto-playlists generation

### Phase 4: Social Features (Week 5)
- Follow system
- Activity feed
- Blend playlists
- Group sessions

### Phase 5: Offline & Sync (Week 6)
- Download manager
- Offline playback
- Cross-device sync

### Phase 6: Audio Enhancement (Week 7)
- Equalizer
- Crossfade
- Volume normalization
- Visualizations

### Phase 7: Premium & Monetization (Week 8)
- Subscription tiers
- Payment integration
- Ad system
- Artist tools
