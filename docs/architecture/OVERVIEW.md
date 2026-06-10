# Architecture Overview — Persian Music Ecosystem

> **Document**: Architecture Overview
> **Status**: v1.0 — Final
> **Target**: 50M users, 10M tracks, 1B streams/month across Iran and MENA

---

## 1. Vision

A world-class Persian music ecosystem combining the best of Spotify, SoundCloud, Genius, Discord, TikTok Discovery, Last.fm, and RateYourMusic — designed for Persian-speaking audiences with native RTL support, Persian calendar, Dastgah modal system, and community-owned content moderation.

### Core Pillars

- **Streaming**: Premium audio at up to 320kbps AAC / FLAC, gapless playback, crossfade, spatial audio
- **Discovery**: AI-driven recommendations, mood/genre exploration, viral charts, trending tracks
- **Social**: Activity feeds, listening parties, live rooms, music clubs, discussions
- **Community**: Wikipedia-style contributions (metadata, lyrics, translations), multi-tier moderation
- **Creator Economy**: Artist profiles, release scheduling, analytics dashboard, fan tipping, ad revenue share
- **Gamification**: XP system, badges, leaderboards, daily challenges, level-based privileges
- **Persian Identity**: Dastgah modal classification, Persian poetry integration, Jalali calendar, Vazirmatn typography

---

## 2. UX Architecture

### 2.1 Navigation Model

```
Root
├── App Shell (persistent)
│   ├── Top Navigation Bar (search, notifications, profile)
│   ├── Sidebar (home, discover, library, playlist, social, settings)
│   ├── Now Playing Bar (persistent bottom bar)
│   ├── Floating Mini Player (overlays when scrolling)
│   └── Main Content Area (router-view)
├── Fullscreen Views
│   ├── FullscreenPlayer (cinematic mode, triggered via space/N)
│   ├── SearchOverlay (modal overlay, triggered via / or Ctrl+K)
│   ├── KaraokeLyrics (word-level sync, fullscreen)
│   └── TheatreMode (album art + lyrics side by side)
├── Bottom Sheet (mobile)
│   └── MobileBottomSheet (snap-drag player, gesture controlled)
└── Modals / Sheets
    ├── Queue Panel (slide-in from right)
    ├── Share Sheet
    ├── Create Playlist
    └── Report / Moderate
```

### 2.2 Player Architecture (8 Modes)

| Mode | Trigger | Description |
|------|---------|-------------|
| Compact | Default | NowPlayingBar — minimal controls, progress, like |
| Fullscreen | Space / N / tap art | Cinematic with album colors, spectrum, queue |
| Floating Mini | Swipe down / minimize | Draggable PiP, 3 sizes (mini/compact/expanded) |
| Theatre | Ctrl+T / button | Album art left, lyrics right, dimmed background |
| Karaoke | K / button | Word-level highlight, particle background |
| Lyrics Focus | Ctrl+L / button | Fullscreen lyrics with scroll, album art thumbnail |
| Ambient | Ctrl+A / button | Fullscreen album art + particle effects, minimal UI |
| Visualizer | V / button | WebGL spectrum/circular/particle/fluid visualizations |

### 2.3 Page Architecture

```
/ (Home)
├── Personalized greeting
├── "Made For You" hero section
├── Continue Listening carousel
├── Mood cards (8 moods with gradients)
├── Trending Now (ranked list)
├── Genre pills strip
└── Recommended Albums carousel

/discover
├── Mood pills (8 interactive)
├── Trending Recently
├── Genre Worlds (5 with gradient banners)
├── Viral Hits (ranked with badges)
├── For You Mixes
└── From the Community feed

/library
├── Playlists (created, saved, followed)
├── Liked Tracks (with smart sort)
├── Albums (saved)
├── Artists (followed)
├── Podcasts (subscribed)
└── Local Files

/search (overlay)
├── Top Result hero card
├── Tracks section
├── Artists section
├── Albums section
├── Playlists section
└── Keyboard navigation (↑↓↩ Esc)

/artist/:id
├── Hero banner (dynamic colors from art)
├── Top tracks (play all)
├── Discography (by year)
├── Related artists
├── About (bio, social links)
├── Fan insights (monthly listeners, charts)
└── Community contributions (lyrics, translations)

/album/:id
├── Cover art (dynamic color extraction)
├── Track list (with BPM indicators, explicit)
├── Credits (artists, producers, writers)
├── Release info (label, date, format)
├── Related albums
└── User reviews / ratings

/playlist/:id
├── Collaborative indicator
├── Track list with drag reorder
├── Playlist stats (duration, followers)
├── Made by badge (with avatar)
└── Smart suggestions (bottom)

/social
├── Activity Feed (friends' listening activity)
├── Listening Parties (active, scheduled)
├── Live Rooms (voice chat + synchronized playback)
├── Music Clubs (genre/artist groups)
├── Discussions (per track/album)
└── Friend recommendations

/creator (artist dashboard)
├── Analytics (streams, listeners, revenue)
├── Releases (manage, schedule)
├── Posts (composer, schedule)
├── Fans (demographics, growth)
├── Royalties (earnings, payouts)
└── Settings (profile, verification)

/contribute
├── My contributions
├── Submit (lyrics, translation, credits, metadata)
├── Moderation queue
├── History (approved, rejected, pending)
└── Leaderboard (top contributors)

/admin
├── Dashboard (KPI cards, charts)
├── Moderation queue (AI + human decisions)
├── Reports (user, content, copyright)
├── Users (search, ban, role management)
├── Content (tracks, albums, playlists management)
├── Appeals (review decisions)
├── Gamification (adjust XP, badges)
└── Settings (platform configuration)
```

---

## 3. Technology Stack

### 3.1 Backend

| Component | Technology | Purpose |
|-----------|------------|---------|
| API Gateway | Go + Gin | Request routing, auth, rate limiting |
| Services | Go (14 microservices) | Business logic per domain |
| Database | PostgreSQL 16 | Primary data store |
| Cache | Redis 7 | Session, rate limits, hot data |
| Search | OpenSearch 2.x | Full-text, autocomplete, faceted search |
| Object Storage | MinIO / S3 | Audio files, images, assets |
| Event Bus | Apache Kafka | Async communication between services |
| Monitoring | Prometheus + Grafana | Metrics, alerting |
| Logging | Zap (structured) | Application logs |
| Tracing | OpenTelemetry | Distributed tracing |
| CI/CD | GitHub Actions + ArgoCD | Automated deployment |
| Container | Docker + Kubernetes | Orchestration |
| Service Mesh | Istio | Traffic management, security |

### 3.2 Frontend

| Component | Technology | Purpose |
|-----------|------------|---------|
| Framework | Vue 3 + Composition API | SPA |
| Language | TypeScript (strict) | Type safety |
| State | Pinia | Reactive state management |
| UI Library | PrimeVue 4 | Accessible components |
| Styling | Tailwind CSS v4 | Utility-first CSS |
| Design System | Custom (Persian Gem) | Glassmorphism, aurora, animations |
| Build | Vite 6 | Fast HMR, optimized builds |
| PWA | vite-plugin-pwa | Offline support, install prompt |
| Audio | Howler.js / Web Audio API | Audio playback, visualization |
| Canvas | PixiJS 8 / Three.js | Visualizers, particles |
| Router | Vue Router 4 | SPA routing |

### 3.3 Database Layer

| Database | Purpose | Sharding |
|----------|---------|----------|
| PostgreSQL (Primary) | Users, tracks, artists, playlists, social graph | Vertical per service |
| PostgreSQL (Analytics) | Streams, listening history, aggregated metrics | Time-partitioned monthly |
| Redis (Cache) | Session, hot tracks, trending, top charts | Cluster mode |
| Redis (Rate Limit) | API rate limiting, anti-abuse | Standalone |
| OpenSearch | Tracks, artists, playlists, lyrics search | Multi-node cluster |
| Kafka (Stream) | Event log, real-time analytics | Partitioned by entity ID |

---

## 4. High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         CDN (CloudFront)                      │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                    API Gateway (Gin)                         │
│         Auth  ·  Rate Limit  ·  Request Validation          │
│         CORS  ·  Compression  ·  Request Logging            │
└──┬──────────┬──────────┬──────────┬──────────┬────────────┘
   │          │          │          │          │
   ▼          ▼          ▼          ▼          ▼
┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
│User  │ │Music │ │Social│ │Search│ │Contrl│  ... 14 services
│Svc   │ │Svc   │ │Svc   │ │Svc   │ │Svc   │
└──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘
   │        │        │        │        │
   └────────┴────────┴────────┴────────┘
                │
        ┌───────▼───────┐
        │    Kafka       │
        │  (Event Bus)   │
        └───┬───┬───┬───┘
            │   │   │
   ┌────────┘   │   └────────┐
   ▼            ▼            ▼
┌──────┐  ┌──────────┐  ┌──────────┐
│Post- │  │  Redis    │  │OpenSearch│
│greSQL│  │  (Cache)  │  │  (Search)│
└──────┘  └──────────┘  └──────────┘
   │
   ▼
┌──────────┐
│   S3 /   │
│  MinIO   │
└──────────┘
```

### 4.1 Request Flow

```
Client → CDN → API Gateway → Auth Middleware → Rate Limiter → Service Handler
  → (if read) → Redis Cache check → PostgreSQL query → Response
  → (if write) → PostgreSQL write → Kafka event → Cache invalidation → Response
```

### 4.2 CQRS Pattern

- **Command Side**: PostgreSQL writes through service handlers
- **Query Side**: Redis cache for hot data, OpenSearch for search, PostgreSQL for complex queries
- **Eventual Consistency**: Kafka events trigger cache invalidation and search index updates
- **Read Models**: Materialized views in PostgreSQL for dashboard/complex queries

### 4.3 Data Flow for Key Operations

**Streaming a Track:**
```
1. Client requests stream URL → Gateway → Music Service
2. Music Service validates subscription/license → Generates pre-signed S3 URL (TTL: 1h)
3. Music Service records stream start → Kafka event → Streams Service
4. Client streams directly from CDN/S3
5. On track end → Client sends complete event → Streams Service finalizes
6. Async: Update artist stats, user history, recommendations, trending
```

**Search:**
```
1. Client sends query → Gateway → Search Service
2. Search Service applies: transliteration, Persian normalization, synonym expansion
3. Queries OpenSearch with edge_ngram autocomplete + boosted fields
4. Returns grouped results: tracks, artists, albums, playlists, lyrics
5. Caches frequent queries in Redis (TTL: 5min)
```

**Contribution:**
```
1. User submits (lyrics, translation, metadata) → Contribution Service
2. AI auto-moderation: spam, profanity, plagiarism check (pass/fail/review)
3. If fail → rejected with reason; if pass → auto-approved; if review → moderation queue
4. Community moderators review → approve/reject with reason
5. If approved → content published, contributor gains XP, event broadcast
6. Version history recorded for audit trail
```

---

## 5. Security Architecture

### 5.1 Authentication

- JWT-based with access token (15min) + refresh token (30 days)
- Refresh token rotation — old token invalidated on use
- OAuth 2.0 / OpenID Connect for social login (Google, Spotify, GitHub)
- Device fingerprinting for suspicious login detection
- Rate limiting: 5 attempts/min per IP, lockout after 10 failures

### 5.2 Authorization

- RBAC with 6 roles: Guest, Listener, Creator, Moderator, Admin, Superadmin
- Permission inheritance: Creator inherits Listener, Admin inherits Moderator
- Fine-grained permissions (e.g., `track:edit`, `user:ban`, `content:moderate`)
- API-level enforcement via middleware by service

### 5.3 Data Protection

- All data encrypted at rest (AES-256) and in transit (TLS 1.3)
- PII encrypted with per-user keys stored in vault
- Audio files: DRM via encrypted HLS streams for premium content
- GDPR-compliant data export and account deletion

### 5.4 Anti-Abuse

- Rate limiting at API Gateway (per IP, per user, per endpoint)
- Behavioral analysis: abnormal patterns trigger review
- Persian NLP: auto-detect spam in Farsi/Arabic text
- CAPTCHA on signup, password reset, and contribution submission
- Device fingerprinting across requests to detect bots
- Content fingerprinting (acoustic + perceptual hashing) for piracy detection

---

## 6. Performance Targets

| Metric | Target |
|--------|--------|
| P99 API response | < 200ms |
| P50 API response | < 50ms |
| Search response | < 300ms (P99) |
| Stream start time | < 1.5s (including auth) |
| Page load (LCP) | < 1.5s |
| Time to interactive | < 3s |
| Concurrent streams | 1M+ |
| Uptime SLA | 99.9% |
| Cache hit ratio | > 90% (hot data) |
| Database query P99 | < 100ms |

---

## 7. Directory Structure

```
/
├── cmd/                    # Entry points
│   ├── api/                # API server
│   ├── worker/             # Background workers
│   └── migrate/            # Database migrations
├── internal/               # Private application code
│   ├── domain/             # Domain models, interfaces
│   ├── service/            # Business logic
│   ├── repository/         # Data access
│   ├── handler/            # HTTP handlers
│   ├── middleware/         # HTTP middleware
│   ├── event/              # Event handlers
│   └── pkg/                # Shared utilities
├── pkg/                    # Shared libraries
│   ├── auth/               # JWT, OAuth
│   ├── storage/            # S3/MinIO abstraction
│   ├── search/             # OpenSearch client
│   ├── cache/              # Redis client
│   └── bus/                # Kafka client
├── api/                    # API definitions
│   ├── openapi/            # OpenAPI specs
│   └── proto/              # Protobuf definitions
├── migrations/             # SQL migrations
├── deploy/                 # Deployment configs
│   ├── k8s/                # Kubernetes manifests
│   └── docker/             # Dockerfiles
├── frontend/               # Vue 3 SPA
│   ├── src/
│   │   ├── components/     # Reusable components
│   │   ├── pages/          # Page components
│   │   ├── stores/         # Pinia stores
│   │   ├── composables/    # Composition API hooks
│   │   ├── services/       # API clients
│   │   └── assets/         # CSS, images, fonts
│   └── public/             # Static assets
└── docs/                   # Documentation
    ├── architecture/       # Architecture documents
    └── DESIGN_SYSTEM.md    # Design system spec
```
