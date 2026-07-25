# Service Architecture — Persian Music Ecosystem

> **Document**: Microservices, Event Bus, Recommendation, Search, Moderation, Gamification, Trust
> **Status**: v1.1 — Current (updated 2026-06-22)
> **🚧 Aspirational**: This document describes a microservices architecture that is the long-term target. The current implementation is a **Go monolith** with all modules in a single binary (`cmd/api`).
> **Target**: 50M users, 10M tracks, 1B streams/month

## Current Status (as of June 2026)

The platform is currently a **Go monolith** hosted in a single `cmd/api` binary. All 30 modules live under `internal/modules/` and share the same database pool. Key architectural notes:

- **Not microservices**: There is no service mesh, no gRPC, no Kafka event bus in production. All inter-module communication is direct Go function calls.
- **Database**: A single PostgreSQL database with all tables. Vertical partitioning planned for future.
- **Search**: OpenSearch is configured in docker-compose but only basic search is implemented.
- **Cache**: Redis is configured but used primarily for sessions and rate limiting.
- **ML Service**: The `moja-ml-service/` directory is a separate Python microservice (FastAPI + Celery) for audio processing (Whisper, embeddings, mood analysis). It's the only genuinely separate service.
- **Frontend**: Single Vue 3 SPA communicating with the Go API via REST + WebSocket.

---

## 1. Microservices Overview

> 🚧 **Aspirational**: These microservices do not yet exist as separate processes. All functionality is implemented as Go packages within a single `internal/modules/` directory tree. When the monolith is decomposed, each module will become a standalone service.

| Service | Responsibility | Database | Replicas | Language |
|---------|---------------|----------|----------|----------|
| API Gateway | Auth, routing, rate limiting, compression | Redis | 5+ | Go |
| User Service | Registration, profiles, sessions, follows | user_db + Redis | 3+ | Go |
| Music Service | Tracks, albums, artists, genres, playlists | music_db + Redis | 5+ | Go |
| Streams Service | Playback sessions, stream recording, analytics | stream_db | 5+ | Go |
| Search Service | Full-text search, autocomplete, faceted search | OpenSearch | 3+ | Go |
| Social Service | Feed, parties, rooms, clubs, comments | social_db + Redis | 3+ | Go |
| Recommendation Service | ML-based recommendations, trending | Redis + OpenSearch | 3+ | Python |
| Contribution Service | User submissions, version history | contribution_db | 3+ | Go |
| Moderation Service | AI + human moderation pipeline | moderation_db | 3+ | Go |
| Gamification Service | XP, badges, challenges, leaderboards | gamification_db + Redis | 2+ | Go |
| Analytics Service | Event collection, aggregation, dashboards | analytics_db + ClickHouse | 3+ | Go |
| Notification Service | Push, email, in-app notifications | Redis | 2+ | Go |
| Payment Service | Subscriptions, tips, payouts | payment_db | 3+ | Go |
| Upload Service | Audio transcoding, image processing, CDN | S3/MinIO | 3+ | Go |

### Service Communication

- **Synchronous**: HTTP/gRPC between services for real-time queries (e.g., API Gateway → User Service)
- **Asynchronous**: Kafka events for eventual consistency (e.g., stream recorded → update play counts → refresh recommendations)
- **Service mesh**: Istio for traffic management, mTLS, circuit breaking, retries

---

## 2. API Gateway

### 2.1 Responsibilities

- SSL termination (TLS 1.3)
- JWT validation and extraction
- Rate limiting per IP (100/s) and per user (500/s)
- Request validation and sanitization
- Request logging (structured JSON via Zap)
- Response compression (gzip, brotli)
- CORS enforcement
- Route to appropriate service

### 2.2 Rate Limiting Algorithm

- **Token bucket** per user+endpoint combination
- Redis-backed with atomic increment + TTL
- Burst allowance: 2x sustained rate for 5 seconds
- Exceeded → 429 with `Retry-After` header
- Critical endpoints (login, signup): stricter limits

### 2.3 Middleware Stack

```go
router := gin.New()
router.Use(
    middleware.RequestID(),
    middleware.Logger(zapLogger),
    middleware.Recovery(),
    middleware.CORS(config.CORS),
    middleware.RateLimiter(redisClient),
    middleware.Auth(jwtService, publicPaths),
    middleware.RequestValidator(),
    middleware.Compression(ginzstd.Middleware),
)
```

---

## 3. User Service

### 3.1 Core Features

- User registration (email, phone, OAuth)
- Authentication (password + JWT + refresh token rotation)
- Profile CRUD
- Follow/unfollow with bidirectional notification
- Block/unblock
- Session management (list active, revoke)
- Account deletion (soft delete with 30-day grace period)

### 3.2 Authentication Flow

```
1. POST /auth/login → validate credentials → generate access (15min) + refresh (30d)
2. Client stores access in memory, refresh in httpOnly cookie
3. POST /auth/refresh → validate refresh → rotate (old invalidated) → new pair
4. POST /auth/logout → revoke refresh → clear server session
5. All protected routes: Authorization: Bearer <access_token>
```

### 3.3 Refresh Token Rotation

```
Request: POST /auth/refresh { refresh_token: "old_abc" }
Validate: Check signature, expiry, family_id
If valid: 
  - Mark old token as used
  - Issue new access + refresh (same family_id)
  - Return new pair
If reused (already used token detected):
  - Revoke entire family (all tokens with that family_id)
  - Require re-login (potential token theft)
```

---

## 4. Music Service

### 4.1 Core Features

- Track CRUD with S3 upload/pre-signed URLs
- Album management
- Artist profiles
- Playlist CRUD with reorder
- Genre tree management
- Like/unlike tracks
- Library management (save albums, playlists, artists)
- Browse by genre, mood, dastgah, year

### 4.2 Streaming Flow

```
1. GET /tracks/:id/stream → validate subscription/license → 
2. Generate pre-signed S3 URL (TTL: 1 hour) → 
3. Return URL + track metadata + DRM key (if premium)
4. Client streams directly from CDN/S3 (no proxying through service)
5. Client sends POST /streams/start → Streams Service records begin
6. Client sends POST /streams/end → Streams Service finalizes
```

### 4.3 Audio Transcoding Pipeline

```
Upload → S3 raw storage → Kafka event → Transcoding Worker:
  - FLAC source → 320kbps MP3 → 192kbps MP3 → 128kbps AAC → 48kbps OPUS
  - Audio fingerprinting (acoustic hashing)
  - Loudness normalization (EBU R128)
  - ReplayGain calculation
  - Waveform data generation (for visualizer)
```

---

## 5. Streams Service

### 5.1 Core Features

- Record stream start/end events
- Calculate completion rate
- Deduplicate (ignore streams < 5 seconds)
- Update play counts on tracks, albums, artists
- Feed data to Analytics Service (Kafka)
- Anti-cheat: detect automated streaming, cap per-user daily streams

### 5.2 Stream Deduplication

```go
// Ignore streams that are:
// - Shorter than 5 seconds (accidental plays)
// - Same user + same track within 30 seconds (double-taps)
// - Exceeding 24 hours of streaming per day (bot detection)
// - More than 500 streams/day (unusual activity)
```

### 5.3 Event Schema (Kafka)

```json
{
  "event_type": "stream.completed",
  "user_id": "uuid",
  "track_id": "uuid",
  "artist_id": "uuid",
  "album_id": "uuid",
  "duration_played": 180,
  "duration_total": 240,
  "completion_rate": 0.75,
  "source": "playlist",
  "device_type": "mobile",
  "country": "IR",
  "session_id": "uuid",
  "timestamp": "2026-06-06T12:00:00Z"
}
```

---

## 6. Social Service

### 6.1 Core Features

- Activity feed generation (user-centric and global)
- Listening parties (synchronized playback)
- Live rooms (voice chat + shared queue)
- Music clubs (interest-based groups)
- Track/album discussions (threaded comments)
- Track ratings (1-10 scale with optional review)
- Share functionality (internal + external links)

### 6.2 Activity Feed Generation

**Fan-out on write** (for active users with < 500 followers):
- When user listens/likes/follows, write event to each follower's feed list in Redis
- Feed stored as Redis List (capped at 200 events per user)
- Pagination via LRANGE with cursor

**Fan-out on read** (for users with 500+ followers):
- Store event once, read + merge at query time
- For power users: hybrid approach (fan-out to top 100 closest followers, read for rest)

### 6.3 Listening Party Sync

```
1. Host creates party → adds tracks to queue
2. Participants join → receive current position + timestamp
3. All clients sync to host's playback position via WebSocket
4. Every 5 seconds: host broadcasts { position_ms, paused, track_id, timestamp }
5. Participants adjust: their_position = host_position + (now - host_timestamp)
6. On host seek/pause/skip: immediate broadcast
7. Chat overlay for party participants
```

### 6.4 WebSocket Connections

- Each user maintains persistent WebSocket to Social Service
- Used for: feed updates, party sync, live room audio, notifications
- Connection pool managed per pod, state in Redis
- Sticky sessions via Istio (session cookie)

---

## 7. Search Service

### 7.1 Core Features

- Unified search across tracks, artists, albums, playlists, lyrics
- Edge_ngram autocomplete (2-20 characters)
- Persian transliteration (English → Farsi)
- Synonym expansion (Persian + English)
- Fuzzy matching (Levenshtein distance 1-2)
- Faceted filtering (genre, mood, dastgah, year, language, explicit)
- Search suggestions ("Did you mean...")
- Recent searches per user (Redis, capped at 10)

### 7.2 Persian NLP Pipeline

```
Input: "gol barg"
→ Normalize: "gol barg"
→ Transliterate: "گل برگ" (phonetic Farsi)
→ Synonym expansion: ["گل", "برگ", "flower", "leaf", "gol", "barg"]
→ Tokenize: ["گل", "برگ"]
→ Remove stop words: [] (none removed)
→ Query with boost: title^3, persian_title^4, artist^2, lyrics^1
```

### 7.3 Query Construction

```go
func BuildSearchQuery(query string, filters SearchFilters) map[string]interface{} {
    return map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{
                    {
                        "multi_match": map[string]interface{}{
                            "query":  query,
                            "fields": []string{
                                "persian_title^4",
                                "title^3",
                                "artist_name^2",
                                "album_name",
                                "lyrics",
                            },
                            "type":            "best_fields",
                            "fuzziness":       "AUTO",
                            "transpositions":  true,
                            "minimum_should_match": "70%",
                        },
                    },
                },
                "filter": buildFilters(filters),
                "should": []map[string]interface{}{
                    {"term": {"language": "fa"}}, // boost Persian results
                },
            },
        },
        "sort": []map[string]interface{}{
            {"_score": "desc"},
            {"play_count": "desc"},
        },
        "size": 20,
        "from": filters.Offset,
    }
}
```

---

## 8. Recommendation Service

### 8.1 Architecture

Hybrid ensemble combining multiple recommendation strategies:

```
┌─────────────────────────────────────────────────────────────┐
│                   Recommendation API                        │
├─────────────────────────────────────────────────────────────┤
│                    Ensemble Engine                          │
│         (weighted blending of candidates)                   │
├──────────┬──────────┬──────────┬──────────┬────────────────┤
│Collabora-│ Content- │ Context- │ Trending │  Freshness     │
│tive      │ Based    │ ual      │ (Popular │  (New          │
│Filtering │ (Audio   │ Bandits  │ + Recency│  Releases)     │
│(ALS)     │ Features)│ (Explore │ Weighted)│                │
│          │          │ /Exploit)│          │                │
└──────────┴──────────┴──────────┴──────────┴────────────────┘
```

### 8.2 Collaborative Filtering (ALS)

- **Algorithm**: Alternating Least Squares with implicit feedback
- **Training**: Weekly batch job on Spark, user-item matrix from listening history
- **Features**: 100 latent factors per user/track
- **Cold start**: New users get genre-based onboarding, new tracks get content-based
- **Output**: Top-N recommendations per user (cached in Redis, TTL: 6h)

### 8.3 Content-Based Filtering

- **Audio analysis**: MFCC, chromagram, spectral features extracted via librosa
- **Metadata features**: Genre, mood, BPM, key, dastgah, language, year
- **Text features**: TF-IDF on lyrics, title, tags
- **Similarity**: Cosine similarity on combined feature vector
- **Use case**: "More like this", "Fans also like", radio mode

### 8.4 Contextual Bandits

- **Algorithm**: LinUCB (contextual linear upper confidence bound)
- **Context**: Time of day, day of week, device, location, recent listens
- **Exploration rate**: ε = 0.1 (10% explore, 90% exploit)
- **Reward signal**: Completion rate > 50% = positive, like = strong positive
- **Use case**: Home page "Made For You", discovery feed, radio mode

### 8.5 Trending Algorithm

```
trending_score = 
    recent_streams * 0.4 + 
    stream_growth_rate * 0.3 +   // (streams_last_24h / streams_prev_24h - 1) * 10
    completion_rate * 0.15 + 
    like_rate * 0.1 +            // likes per stream
    freshness_boost * 0.05       // exp(-days_since_release / 30)

// Normalized per genre to ensure genre diversity
```

### 8.6 Recommendation Cache Strategy

| Endpoint | TTL | Refresh Trigger |
|----------|-----|----------------|
| Home feed | 30 min | On new stream > 3 |
| Artist radio | 1 hour | On artist release |
| Track radio | 1 hour | On like/unlike |
| Search suggestions | 5 min | N/A |
| Trending | 15 min | Every stream event |
| "Made For You" | 1 hour | Daily (or on 5+ new streams) |

---

## 9. Contribution Service

### 9.1 Core Features

- Submit contributions (lyrics, translations, credits, metadata, bio, album art)
- AI auto-moderation with confidence scoring
- Version history with full snapshot per version
- Rollback capability (restore previous version)
- Attribution tracking (primary + secondary contributors)
- XP reward on approval

### 9.2 Contribution Workflow

```
Submit → AI Auto-Moderation
         ├── Pass (confidence > 0.9) → Auto-approve → Publish → Grant XP
         ├── Fail (confidence > 0.95) → Reject → Notify user
         └── Review (else) → Community Moderation Queue
                              ├── Approve → Publish → Grant XP
                              └── Reject → Notify user + reason
```

### 9.3 AI Moderation Checks

| Check | Method | Threshold |
|-------|--------|-----------|
| Spam detection | Regex + ML classifier | > 0.9 confidence |
| Profanity filter | Persian + English word list | Any match |
| Plagiarism check | Levenshtein vs existing content | > 80% similarity |
| Gibberish detection | Entropy + language model | < 0.3 coherence |
| Metadata validation | Format + range checks | Strict |
| Image NSFW check | Vision model (for album art) | > 0.9 NSFW → reject |

---

## 10. Moderation Service

### 10.1 4-Layer Moderation Pipeline

```
Layer 1 — AI Auto-Moderation (real-time, < 100ms)
  - Spam, profanity, gibberish, NSFW, plagiarism
  - Handles ~80% of submissions automatically

Layer 2 — Community Moderation (~1-5min)
  - Trusted users (trust score > 50) can review queue
  - Weighted voting: higher trust = more weight
  - 3+ matching decisions → action taken

Layer 3 — Expert Moderation (~1-24h)
  - Assigned experts by content type (lyrics, music theory, translation)
  - Handle escalated items, edge cases, appeals

Layer 4 — Audit & Appeal (~24-72h)
  - Admin review of disputed decisions
  - User appeals with justification
  - Pattern analysis for moderator performance
```

### 10.2 Trust Scoring

| Action | Score Change |
|--------|-------------|
| Approved contribution | +5 |
| Accurate moderation vote | +2 |
| Inaccurate moderation vote | -3 |
| Report confirmed valid | +3 |
| Report confirmed false | -5 |
| Content removal | -10 |
| Multiple violations | -50 (temporary ban) |

### 10.3 Community Moderation Weight

```go
func ModerationWeight(user *User) float64 {
    baseWeight := 1.0
    trustBonus := float64(user.TrustScore) / 100.0 * 2.0 // up to +2x
    levelBonus := float64(user.Level) / 20.0 * 0.5       // up to +0.5x
    accuracyBonus := user.ModerationAccuracy * 0.5        // up to +0.5x
    return baseWeight + trustBonus + levelBonus + accuracyBonus
}

// Final weight range: 1.0x (new) to 4.0x (expert)
```

---

## 11. Gamification Service

### 11.1 XP System

**XP Sources:**
| Action | XP | Daily Cap |
|--------|----|-----------|
| Stream a track | +1 | 100 |
| Like a track | +5 | 30 |
| Share a track | +10 | 20 |
| Create a playlist | +20 | 5 |
| Complete daily challenge | +50-200 | N/A |
| Approved contribution | +50 | 10 |
| Moderation vote matches | +20 | 20 |
| Review/rate a track | +10 | 20 |
| Daily login streak (day 1-7) | +15, +30, +50, +75, +100, +150, +200 | 1 |
| Invite friend (verified) | +100 | 5 |
| First stream of the day | +25 | 1 |

**Level Progression:**
```go
func XpForLevel(level int) int64 {
    return int64(math.Pow(float64(level), 2.5) * 100)
}
// Level 1: 100 XP, Level 5: 5,590 XP, Level 10: 31,623 XP, Level 20: 178,885 XP
```

### 11.2 Level Titles (Persian-English)

| Level | XP Required | Title (EN) | Title (FA) |
|-------|-------------|------------|------------|
| 1 | 0 | Novice Listener | شنونده تازه‌کار |
| 2 | 100 | Curious Ear | گوش کنجکاو |
| 3 | 520 | Melody Seeker | جوینده ملودی |
| 4 | 1,600 | Rhythm Catcher | ریتم‌گیر |
| 5 | 3,900 | Harmony Lover | عاشق هارمونی |
| 6 | 8,000 | Tone Explorer | کاوشگر نغمه |
| 7 | 14,000 | Music Wanderer | سرگردان موسیقی |
| 8 | 22,000 | Note Collector | جمع‌آورنده نت |
| 9 | 33,000 | Groove Master | استاد گروو |
| 10 | 48,000 | Beat Connoisseur | خبره ضرب‌آهنگ |
| 11 | 66,000 | Vocal Virtuoso | استاد آواز |
| 12 | 88,000 | Lyric Sage | فرزانه شعر |
| 13 | 115,000 | Dastgah Navigator | راهبر دستگاه |
| 14 | 147,000 | Genre Weaver | بافنده سبک |
| 15 | 185,000 | Orchestra Commander | فرمانده ارکستر |
| 16 | 230,000 | Melody Architect | معمار ملودی |
| 17 | 283,000 | Music Scholar | دانشور موسیقی |
| 18 | 345,000 | Sound Alchemist | کیمیاگر صدا |
| 19 | 417,000 | Legendary Listener | شنونده افسانه‌ای |
| 20 | 500,000 | Persian Gem | گوهر پارسی |

### 11.3 Badge Categories

| Category | Example Badges | Rarity |
|----------|---------------|--------|
| Listener | Night Owl (100 streams 12am-5am), Explorer (20 genres), Marathon (1000 streams in a week) | Common→Epic |
| Contributor | Wordsmith (10 lyrics), Translator (5 translations), Curator (10 playlists) | Common→Legendary |
| Social | Networker (100 followers), Party Host (10 parties), Club President (create club with 100 members) | Rare→Legendary |
| Special | Early Bird (first 1000 users), Beta Tester, Community Hero (top moderator) | Hidden/Legendary |

### 11.4 Daily Challenges

| Challenge | Target | XP | Rotation |
|-----------|--------|----|----------|
| Stream 10 tracks | 10 streams | 50 | Daily |
| Explore a genre | 3 tracks in new genre | 75 | Daily |
| Like 5 tracks | 5 likes | 50 | Daily |
| Share a track | 1 share | 100 | Daily |
| Listen for 1 hour | 60 min total | 100 | Daily |
| Complete a playlist | Finish 1 full playlist | 75 | Daily |
| Weekly: 200 streams | 200 streams | 500 | Weekly |
| Weekly: 5 contributions | 5 approved | 750 | Weekly |

---

## 12. Trust System

### 12.1 Score Range & Tiers

| Tier | Score | Label | Permissions |
|------|-------|-------|-------------|
| 0 | -100 to -51 | Banned | No access |
| 1 | -50 to -1 | Restricted | Read-only, no contributions |
| 2 | 0 to 19 | New | Full read, limited write (3 contributions/day) |
| 3 | 20 to 49 | Regular | Full access, 10 contributions/day |
| 4 | 50 to 79 | Trusted | Moderation voting, 50 contributions/day |
| 5 | 80 to 100 | Expert | Skip AI review, direct publish, moderation veto power |

### 12.2 Score Modifiers

| Action | Score Delta | Notes |
|--------|-------------|-------|
| Account creation | +0 | Start at 0 |
| Email verified | +5 | One-time |
| Phone verified | +5 | One-time |
| Contribution approved | +2 | Per contribution |
| Daily active (7 day streak) | +1 | Per day |
| Followed by trusted user | +1 | Per follow (max 10/day) |
| Report confirmed valid | +3 | Per report |
| Moderation matches consensus | +2 | Per decision |
| Contribution rejected | -3 | Per rejection |
| Report confirmed false | -5 | Per report |
| Content removed (violation) | -15 | Per removal |
| Multiple violations (3+/30d) | -50 | Automatic restriction |
| Ban evasion detected | -100 | Permanent ban |

### 13.3 Anti-Spam & Abuse

| Layer | Mechanism | Threshold |
|-------|-----------|-----------|
| Rate limiting | Per-endpoint token bucket | 5x normal user rate |
| Behavioral analysis | Abnormal click speed, pattern detection | 3x standard deviation |
| Persian NLP | Spam classifier on text content | > 0.85 confidence |
| Device fingerprinting | Canvas, WebGL, audio context fingerprint | Same fingerprint, different accounts |
| CAPTCHA | ReCAPTCHA v3 + custom puzzles | Suspicious behavior trigger |
| IP reputation | Known proxy/VPN/Tor detection | Block or CAPTCHA |
| Account age gating | New accounts (< 7 days) | Reduced limits, mandatory CAPTCHA |
| Shadow banning | Content visible only to self | Pattern: likely spam but uncertain |

---

## 13. Kafka Event Bus

### 13.1 Topics & Partitions

| Topic | Partitions | Retention | Key | Consumers |
|-------|------------|-----------|-----|-----------|
| `stream.completed` | 12 | 7 days | track_id | Analytics, Recommendation, Music |
| `track.liked` | 6 | 7 days | user_id | Social, Recommendation, Gamification |
| `track.shared` | 6 | 7 days | user_id | Social, Gamification |
| `user.followed` | 6 | 7 days | followee_id | Social, Notification |
| `contribution.submitted` | 6 | 3 days | target_id | Moderation |
| `contribution.approved` | 6 | 30 days | target_id | Music, Gamification |
| `moderation.action` | 3 | 7 days | content_id | Moderation, Notification |
| `notification.send` | 6 | 1 day | user_id | Notification |
| `analytics.event` | 12 | 30 days | event_type | Analytics |
| `cache.invalidate` | 3 | 1 hour | cache_key | All services |

### 13.2 Event Schema Standard

```protobuf
message Event {
  string event_id = 1;
  string event_type = 2;
  string source_service = 3;
  string aggregate_type = 4; // user, track, artist, etc.
  string aggregate_id = 5;
  google.protobuf.Struct payload = 6;
  map<string, string> metadata = 7;
  google.protobuf.Timestamp timestamp = 8;
  int64 version = 9;
}
```

### 13.3 Exactly-Once Semantics

- Producer: `acks=all`, `enable.idempotence=true`
- Consumer: `isolation.level=read_committed`, manual commit after processing
- Idempotent handlers via event_id dedup in Redis (TTL: 7 days)
- Dead letter queue for failed events (3 retries, then DLQ)

---

## 14. Notification Service

### 14.1 Channel Support

| Channel | Latency | Volume | Content |
|---------|---------|--------|---------|
| In-app (WebSocket) | < 1s | Unlimited | All notifications |
| Push (FCM/APNs) | < 30s | 100/user/day | Likes, follows, milestones |
| Email (SMTP) | < 5min | 10/user/day | Weekly digest, account, security |
| SMS | < 30s | 2/user/day | Verification codes, critical alerts |

### 14.2 Notification Types

- **Social**: follow, like, comment, share, party invite, club invite
- **Contribution**: approved, rejected, needs revision
- **Milestone**: level up, badge earned, challenge completed
- **Creator**: new follower milestone, stream milestone, earnings report
- **System**: account security, subscription, feature announcements

### 14.3 Preference Management

Each user has notification preferences stored in Redis JSON:
```json
{
  "email": { "weekly_digest": true, "security": true, "marketing": false },
  "push": { "like": true, "follow": true, "milestone": true, "party": false },
  "in_app": { "all": true },
  "quiet_hours": { "enabled": true, "start": "23:00", "end": "08:00", "timezone": "Asia/Tehran" }
}
```

---

## 15. Payment Service (Future)

- **Subscription tiers**: Free (ads), Premium (no ads, 320kbps), Family (6 accounts), Student (50% off)
- **Acceptance**: Zarinpal, IDPay, Stripe (for MENA ex-Iran), cryptocurrency
- **Creator payouts**: Monthly, min 500,000 IRR / $10, based on pro-rata stream share
- **Tips**: One-time or monthly subscription to creator, platform takes 15%
- **Revenue split**: 70% artist / 30% platform for streams, 85% / 15% for tips
