# Muse — Domain Glossary

Muse is a full-stack Persian music streaming platform. This document defines the key domain concepts used across the codebase.

## Player & Audio

### Player Engine
The module that composes audio playback, queue management, track preloading, and Media Session API integration behind a single seam. Exposes actions (`play`, `pause`, `seek`, `next`, `previous`) and emits events (`trackchange`, `playstate`, `timeupdate`, `buffering`, `error`). Accepts optional provider interfaces for network-dependent operations (track fetching, play history). Testable by substituting providers and the audio backend.

### Track
A single audio entity with id, title, artist, album, duration, streamUrl, and coverUrl.

### PlaybackTrack
The runtime representation of a Track used by the Player Engine. Adds normalized media URLs and runtime-only fields.

### Queue
An ordered list of PlaybackTracks awaiting playback. Managed by QueueManager within the Player Engine.

### Shuffle
Four modes: off (sequential), queue (shuffled order of current queue), catalog (random tracks from full catalog), similar (random tracks similar to current).

### Repeat
Three modes: off, one (repeat current track), all (loop queue).

### Radio Mode
Infinite playback mode generating a continuous stream of tracks based on seed tracks, genre, or mood. Uses a separate endpoint from the queue.

### PiP (Picture-in-Picture)
Desktop PiP mode using the Document Picture-in-Picture API to float a mini player above other windows.

## Catalog & Search

### Catalog
The central repository of tracks, albums, artists, and genres. Managed by the `catalog` module (Go backend) and consumed by `@catalog` (Vue frontend).

### Artist
A music artist with bio, image, genre affiliations, and a list of albums/tracks. Can be a creator on the platform.

### Album
A collection of tracks by one or more artists, with cover art, release date, and genre.

### Genre / Subgenre
Music categories (e.g., Pop, Rock, Classical) and their subcategories. Used for browsing, recommendations, and the mood explorer.

### Playlist
A user-curated or auto-generated ordered list of tracks. Supports collaborative editing, privacy levels (public, private, unlisted), and playlist art.

### Library
A user's personal collection of saved tracks, albums, playlists, and followed artists.

### History
The chronological record of a user's listening activity. Used for recommendations and recently played views.

### Search
Full-text search across tracks, albums, artists, and playlists. Uses OpenSearch for indexing.

## Auth & Security

### Authentication
JWT-based authentication with access/refresh token pairs. Tokens are issued on login/register, refreshed transparently, and validated on every API request.

### Authorization
Role-based access control (RBAC) with roles like `user`, `creator`, `moderator`, `admin`. Auth middleware checks required roles on protected routes.

### Session Restoration
On page reload, the frontend restores the session by validating the stored token and refreshing if needed. The `@auth` agent handles this flow.

## Social & Community

### Listening Party
Real-time synchronized listening sessions where participants hear the same track at the same time. Uses WebSocket for state synchronization.

### Live Room
Persistent audio chat rooms with track sharing, reactions, and real-time interaction via WebSocket.

### Music Club
Long-lived groups centered around music genres or themes, with shared playlists, discussion, and events.

### Activity Feed
A timeline of user actions (listened to X, liked Y, created playlist Z). Powers the social feed and notifications.

### Reactions
Emoji-based reactions to tracks, comments, and activities. Stored efficiently as bitmaps for performance.

## Creator Economy & Gamification

### XP (Experience Points)
Points earned by listening, engaging, and contributing. Drives level progression and badge unlocks.

### Badge
Achievement badges awarded for milestones (e.g., "Listened to 100 tracks", "First contribution"). Badges are displayed on user profiles.

### Challenge / Daily Challenge
Time-limited goals (daily, weekly) that reward XP upon completion. Challenges rotate and are tracked per user.

### Leaderboard
Ranking of users by XP, listening time, or other metrics. Supports weekly, monthly, and all-time views.

### Contribution
User-submitted content (lyrics, translations, corrections). Contributions go through moderation before being accepted. Contributors earn XP and tips.

### Tips & Subscriptions
Monetary support for creators via one-time tips (using the payment gateway) or recurring subscriptions.

## Admin & Moderation

### Moderation Queue
Pending user contributions that need moderator review. Moderators can approve, reject, or request changes.

### User Management
Admin CRUD for users: view profiles, manage roles, suspend/ban accounts.

### Content Management
Admin CRUD for tracks, albums, artists, genres. Includes media upload and metadata editing.

## Infrastructure

### WebSocket
Real-time bidirectional communication for player sync, social features, presence, and notifications. Powered by the `ws` Hub in Go with gorilla/websocket.

### Redis
Used for caching, rate limiting, rate-limit counters, and WebSocket pub/sub bridging across instances.

### OpenSearch
Full-text search index for tracks, albums, artists, and playlists. Separate from the primary PostgreSQL database.

### Event Bus
Internal pub/sub event system for cross-module communication (e.g., "track played" triggers history recording + XP award + recommendation update).

### Audio Engine (Frontend)
The browser-side audio playback engine using the Web Audio API. Handles streaming, buffering, and audio visualizer data.

### ML Service (Python)
Python microservice for Whisper speech-to-text (lyrics transcription), cover art optimization, and Celery-based async task processing.

## Architecture Terms

### Seam
A boundary where one module can be replaced without changing its consumers. In the frontend, the PlayerEngine interface is a seam between playback orchestration and UI state.

### Deep Module
A module whose interface is significantly simpler than its implementation. The PlayerEngine is deep: 12 action methods + 6 events replace 10+ direct singleton dependencies spread across the store.

### Locality
When related logic lives in one module rather than scattered. Playback orchestration (audio + queue + preload + media-session + shuffle + repeat) has locality inside the Player Engine.

### Enricher
Background pipeline that fetches metadata (cover art, lyrics, similar tracks) from external sources (MusicBrainz, Last.fm, Spotify, Deezer, LRCLib) after tracks are ingested.
