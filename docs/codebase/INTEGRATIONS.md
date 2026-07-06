# Muse — External Integrations

## Music Data Sources (Import)

| Service | Integration | Direction | Auth | Purpose |
|---------|-------------|-----------|------|---------|
| **MusicBrainz** | REST API | Outbound | None (rate-limited) | Artist/track metadata, MBID lookup, cover art |
| **LastFM** | REST API | Outbound | API Key | Track stats, similar artists, album art |
| **Spotify** | REST API | Outbound | Client Credentials | Track metadata, search, import |
| **Deezer** | REST API | Outbound | None | Track metadata, search, import |
| **iTunes/Apple Music** | REST API | Outbound | None | Track metadata, search, import |
| **CoverArt Archive** | REST API | Outbound | None | Album cover images |
| **LRCLib** | REST API | Outbound | None | Synchronized lyrics (LRC format) |

## AI/ML Services

| Service | Integration | Direction | Auth | Purpose |
|---------|-------------|-----------|------|---------|
| **OpenAI** | REST API (Go client) | Outbound | API Key | Text embeddings, mood analysis, song matching |
| **OpenRouter** | REST API | Outbound | API Key | AI-powered lyric generation/analysis |
| **Whisper (self-hosted)** | Python FastAPI (moja-ml-service) | Internal | HMAC-signed | Speech-to-text for lyrics transcription |
| **Custom ML Service** | REST API (moja-ml-service) | Internal | HMAC-signed | Cover image optimization, AI moderation |

## Infrastructure

| Service | Integration | Direction | Auth | Purpose |
|---------|-------------|-----------|------|---------|
| **PostgreSQL** | pgx/sqlx | Internal | User/Pass | Primary database |
| **Redis** | go-redis | Internal | Password | Cache, sessions, pub/sub |
| **OpenSearch** | OpenSearch Go client | Internal | None (dev) | Full-text search index |
| **MinIO / S3** | AWS SDK v2 | Internal | Access Key/Secret | Object storage (media files) |

## Payment Gateways

| Gateway | Type | Status | Purpose |
|---------|------|--------|---------|
| **IDPay** | Iranian payment gateway | Implemented | Subscription payments |
| **ZarinPal** | Iranian payment gateway | Implemented | Subscription payments |

## Social/External

| Service | Integration | Purpose |
|---------|-------------|---------|
| **YouTube-dl / yt-dlp** | Subprocess | Audio download for import |
| **FFmpeg** | Subprocess | Audio transcoding, thumbnail extraction |
| **Audio Tag Library** | dhowden/tag (Go) | Audio metadata parsing (ID3, FLAC, etc.) |

## Frontend Integrations

| Service | Method | Purpose |
|---------|--------|---------|
| **Media Session API** | Browser API | OS-level playback controls (lock screen, notification) |
| **Document Picture-in-Picture** | Browser API | Floating video/mini-player window |
| **Web Audio API** | Browser API | Audio playback, FFT visualization, equalizer |
| **WebSocket** | gorilla/websocket (backend) + native WebSocket (frontend) | Real-time features (social rooms, collaborative playlists) |

## API Contracts (Internal)

- **Go ↔ Go**: Direct function calls via event bus or service interfaces
- **Go ↔ Python (ML)**: REST API with HMAC request signing (shared secret)
- **Go ↔ Frontend**: REST JSON API (`/api/v1/...`) + WebSocket
- **Frontend ↔ ML Service**: Via Go API reverse proxy (no direct frontend→ML calls)

## Environment Variables

Key environment variables for integrations:
```
POSTGRES_URL=postgres://user:pass@host:5432/db?sslmode=disable
REDIS_ADDR=host:6379
REDIS_PASSWORD=
JWT_ACCESS_SECRET=...
JWT_REFRESH_SECRET=...
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
IMPORT_PROXY=socks5://127.0.0.1:1081
LASTFM_API_KEY=...
SPOTIFY_CLIENT_ID=...
SPOTIFY_CLIENT_SECRET=...
```

## WebSocket Events

| Event | Direction | Purpose |
|-------|-----------|---------|
| `party:state` | Server→Client | Party room state updates |
| `party:track:queue` | Bidirectional | Party queue management |
| `club:stage:join` | Server→Client | Stage participant join |
| `room:chat` | Bidirectional | Room chat messages |
| `player:sync` | Server→Client | Collaborative playlist sync |
| `notification:*` | Server→Client | Real-time notifications |
