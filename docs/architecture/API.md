# API Design — Persian Music Ecosystem

> **Document**: API Design, Endpoints, Patterns, Conventions
> **Status**: v1.0 — Final
> **Target**: 50M users, 10M tracks, 1B streams/month

---

## 1. Design Conventions

### 1.1 URL Structure

```
/api/v1/{service}/{resource}[/{id}][/{subresource}]
```

### 1.2 HTTP Methods

| Method | Operation | Idempotent | Safe |
|--------|-----------|------------|------|
| GET | Read / List | Yes | Yes |
| POST | Create / Action | No | No |
| PUT | Full update / Replace | Yes | No |
| PATCH | Partial update | Yes* | No |
| DELETE | Delete | Yes | No |

### 1.3 Response Envelope

**Success:**
```json
{
  "data": { ... },
  "meta": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

**Error:**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Title is required",
    "details": [
      { "field": "title", "message": "must not be empty" }
    ]
  },
  "request_id": "req_abc123"
}
```

### 1.4 Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| VALIDATION_ERROR | 422 | Request body validation failed |
| UNAUTHORIZED | 401 | Missing or invalid auth token |
| FORBIDDEN | 403 | Authenticated but not permitted |
| NOT_FOUND | 404 | Resource does not exist |
| RATE_LIMITED | 429 | Too many requests |
| CONFLICT | 409 | Duplicate resource |
| INTERNAL_ERROR | 500 | Server error |
| SERVICE_UNAVAILABLE | 503 | Downstream service down |

### 1.5 Pagination

- **Cursor-based** for real-time feeds (activity, comments, notifications)
- **Offset-based** for static lists (library, search results, artist discography)
- Default page_size: 20, max: 100
- Response includes `meta.page`, `meta.page_size`, `meta.total`, `meta.total_pages`

### 1.6 Request ID

- Every request tagged with `X-Request-ID` (UUID v7)
- Propagated to downstream services via headers
- Included in all responses and error payloads
- Logged by all services for tracing

### 1.7 API Versioning

- Version in URL path: `/api/v1/...`
- Minor/backward-compatible changes: additive only (new fields, new endpoints)
- Breaking changes: new version (`/api/v2/...`), old version deprecated for 6 months
- Deprecation communicated via `Sunset` header

### 1.8 OpenAPI / Swagger

- All endpoints documented in OpenAPI 3.0 specs per service
- `api/openapi/` directory with per-service spec files
- Combined spec at `api/openapi/openapi.yaml`
- Swagger UI served at `/api/docs` in development

---

## 2. Authentication Endpoints

### 2.1 Auth Service

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/auth/register | Register new user |
| POST | /api/v1/auth/login | Login with email/phone + password |
| POST | /api/v1/auth/refresh | Refresh access token |
| POST | /api/v1/auth/logout | Logout (revoke refresh token) |
| POST | /api/v1/auth/forgot-password | Send password reset email |
| POST | /api/v1/auth/reset-password | Reset password with token |
| POST | /api/v1/auth/verify-email | Verify email with code |
| POST | /api/v1/auth/verify-phone | Verify phone with SMS code |
| GET | /api/v1/auth/oauth/{provider} | Initiate OAuth flow |
| GET | /api/v1/auth/oauth/{provider}/callback | OAuth callback |

**POST /auth/register:**
```json
{
  "username": "user123",
  "email": "user@example.com",
  "password": "SecurePass123!",
  "display_name": "User",
  "locale": "fa"
}
→ 201
{
  "data": {
    "user": { "id": "uuid", "username": "user123", "display_name": "User" },
    "tokens": {
      "access_token": "jwt...",
      "expires_in": 900,
      "refresh_token": "jwt...",
      "refresh_expires_in": 2592000
    }
  }
}
```

---

## 3. User Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/users/me | Get current user profile |
| PATCH | /api/v1/users/me | Update profile |
| DELETE | /api/v1/users/me | Delete account (soft) |
| GET | /api/v1/users/me/library | Get user's library |
| POST | /api/v1/users/me/library | Add to library |
| DELETE | /api/v1/users/me/library/{type}/{id} | Remove from library |
| GET | /api/v1/users/{id} | Get user profile (public) |
| GET | /api/v1/users/{id}/tracks | Get user's liked tracks |
| GET | /api/v1/users/{id}/playlists | Get user's playlists |
| GET | /api/v1/users/{id}/followers | Get user's followers |
| GET | /api/v1/users/{id}/following | Who user follows |
| POST | /api/v1/users/{id}/follow | Follow user |
| DELETE | /api/v1/users/{id}/follow | Unfollow user |
| GET | /api/v1/users/{id}/activity | Get user's public activity |
| GET | /api/v1/users/me/feed | Get personalized activity feed |
| GET | /api/v1/users/me/stats | Get personal listening stats |
| GET | /api/v1/users/me/contributions | Get user's contributions |

---

## 4. Music Endpoints

### 4.1 Tracks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/tracks | List tracks (filtered/browsed) |
| GET | /api/v1/tracks/{id} | Get track details |
| GET | /api/v1/tracks/{id}/stream | Get streaming URL |
| POST | /api/v1/tracks/{id}/like | Like track |
| DELETE | /api/v1/tracks/{id}/like | Unlike track |
| GET | /api/v1/tracks/{id}/related | Get related tracks |
| GET | /api/v1/tracks/{id}/lyrics | Get lyrics (with translation) |
| GET | /api/v1/tracks/{id}/comments | Get track comments |
| POST | /api/v1/tracks/{id}/comments | Add comment |
| GET | /api/v1/tracks/{id}/ratings | Get rating distribution |
| POST | /api/v1/tracks/{id}/ratings | Rate track (1-10) |
| GET | /api/v1/tracks/{id}/waveform | Get waveform data |

### 4.2 Albums

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/albums | List albums |
| GET | /api/v1/albums/{id} | Get album details |
| GET | /api/v1/albums/{id}/tracks | Get album tracks |
| POST | /api/v1/albums/{id}/save | Save to library |
| DELETE | /api/v1/albums/{id}/save | Remove from library |

### 4.3 Artists

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/artists | List artists |
| GET | /api/v1/artists/{id} | Get artist profile |
| GET | /api/v1/artists/{id}/tracks | Get artist's top tracks |
| GET | /api/v1/artists/{id}/albums | Get artist's discography |
| GET | /api/v1/artists/{id}/related | Get related artists |
| GET | /api/v1/artists/{id}/stats | Get artist stats (monthly listeners, etc.) |
| POST | /api/v1/artists/{id}/follow | Follow artist |
| DELETE | /api/v1/artists/{id}/follow | Unfollow artist |

### 4.4 Playlists

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/playlists | List playlists (browse/trending) |
| GET | /api/v1/playlists/{id} | Get playlist details |
| POST | /api/v1/playlists | Create playlist |
| PATCH | /api/v1/playlists/{id} | Update playlist metadata |
| DELETE | /api/v1/playlists/{id} | Delete playlist |
| GET | /api/v1/playlists/{id}/tracks | Get playlist tracks |
| POST | /api/v1/playlists/{id}/tracks | Add track to playlist |
| DELETE | /api/v1/playlists/{id}/tracks/{position} | Remove track |
| PATCH | /api/v1/playlists/{id}/tracks/reorder | Reorder tracks |
| POST | /api/v1/playlists/{id}/follow | Follow playlist |
| DELETE | /api/v1/playlists/{id}/follow | Unfollow playlist |

**POST /playlists:**
```json
{
  "title": "Persian Chill",
  "description": "Best Persian lo-fi vibes",
  "is_public": true,
  "is_collaborative": false
}
→ 201
```

### 4.5 Genres & Moods

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/genres | List all genres (tree) |
| GET | /api/v1/genres/{slug} | Get genre details |
| GET | /api/v1/genres/{slug}/tracks | Get genre tracks |
| GET | /api/v1/moods | List mood categories |
| GET | /api/v1/moods/{mood}/tracks | Get mood tracks |
| GET | /api/v1/dastgahs | List Persian dastgahs |
| GET | /api/v1/dastgahs/{dastgah}/tracks | Get dastgah tracks |

---

## 5. Search Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/search | Full-text search |
| GET | /api/v1/search/suggestions | Autocomplete suggestions |
| GET | /api/v1/search/recent | Recent searches (auth) |
| DELETE | /api/v1/search/recent | Clear recent searches |

**GET /search?q=گل&type=track,artist&genre=pop&limit=20&offset=0**
```json
{
  "data": {
    "tracks": {
      "items": [...],
      "total": 150
    },
    "artists": {
      "items": [...],
      "total": 12
    },
    "albums": {
      "items": [...],
      "total": 8
    },
    "playlists": {
      "items": [...],
      "total": 5
    }
  },
  "meta": {
    "query": "گل",
    "corrected_query": null,
    "suggestions": ["گل و نارنج", "گل یخ", "گل پونه"],
    "took_ms": 45
  }
}
```

---

## 6. Social Endpoints

### 6.1 Activity Feed

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/feed | Get user's activity feed |
| GET | /api/v1/feed/global | Get global activity feed |
| GET | /api/v1/feed/trending | Get trending activity |

### 6.2 Listening Parties

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/social/parties | Create party |
| GET | /api/v1/social/parties/{id} | Get party details |
| POST | /api/v1/social/parties/{id}/join | Join party |
| POST | /api/v1/social/parties/{id}/leave | Leave party |
| POST | /api/v1/social/parties/{id}/sync | Broadcast sync position |
| GET | /api/v1/social/parties/active | List active public parties |
| PATCH | /api/v1/social/parties/{id}/queue | Update party queue |

### 6.3 Live Rooms

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/social/rooms | Create live room |
| GET | /api/v1/social/rooms/{id} | Get room details |
| POST | /api/v1/social/rooms/{id}/join | Join room |
| POST | /api/v1/social/rooms/{id}/leave | Leave room |
| GET | /api/v1/social/rooms/live | List active live rooms |

### 6.4 Music Clubs

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/social/clubs | Create club |
| GET | /api/v1/social/clubs/{id} | Get club details |
| PATCH | /api/v1/social/clubs/{id} | Update club |
| POST | /api/v1/social/clubs/{id}/join | Join club |
| POST | /api/v1/social/clubs/{id}/leave | Leave club |
| GET | /api/v1/social/clubs | Browse clubs |
| GET | /api/v1/social/clubs/{id}/members | List members |

### 6.5 Comments & Ratings

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/comments/{id} | Get comment thread |
| POST | /api/v1/comments/{id}/reply | Reply to comment |
| DELETE | /api/v1/comments/{id} | Delete own comment |
| POST | /api/v1/comments/{id}/like | Like comment |
| POST | /api/v1/comments/{id}/report | Report comment |

---

## 7. Recommendation Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/recommendations/home | Home page recommendations |
| GET | /api/v1/recommendations/track/{id} | "More like this" for track |
| GET | /api/v1/recommendations/artist/{id} | Artist radio |
| GET | /api/v1/recommendations/album/{id} | Album radio |
| GET | /api/v1/recommendations/discover | Discovery page recommendations |
| GET | /api/v1/recommendations/mix/{type} | Get a mix (daily, chill, focus, party) |
| GET | /api/v1/trending | Trending tracks |
| GET | /api/v1/trending/{country} | Trending by country |
| GET | /api/v1/trending/genre/{slug} | Trending by genre |
| GET | /api/v1/new-releases | New releases |

---

## 8. Contribution Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/contributions | Submit contribution |
| GET | /api/v1/contributions | List user's contributions |
| GET | /api/v1/contributions/{id} | Get contribution details |
| PATCH | /api/v1/contributions/{id} | Update pending contribution |
| DELETE | /api/v1/contributions/{id} | Withdraw pending contribution |
| GET | /api/v1/contributions/{id}/history | Get version history |
| GET | /api/v1/contributions/leaderboard | Top contributors |

**POST /contributions:**
```json
{
  "contribution_type": "lyrics",
  "target_type": "track",
  "target_id": "uuid",
  "locale": "fa",
  "data": {
    "lyrics": "متن کامل شعر...",
    "timestamps": [
      { "time": 0, "text": "..." },
      { "time": 5.2, "text": "..." }
    ]
  },
  "is_minor": false
}
→ 201
{
  "data": {
    "id": "uuid",
    "status": "pending",
    "ai_verdict": "pass",
    "ai_confidence": 0.97
  }
}
```

---

## 9. Moderation Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/moderation/queue | Get moderation queue |
| POST | /api/v1/moderation/queue/{id}/review | Submit moderation decision |
| GET | /api/v1/moderation/stats | Moderator stats |
| POST | /api/v1/moderation/reports | Submit report |
| GET | /api/v1/moderation/reports | List user's reports |
| GET | /api/v1/moderation/appeals | List user's appeals |
| POST | /api/v1/moderation/appeals | Submit appeal |

---

## 10. Gamification Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/gamification/profile | Get current XP, level, rank |
| GET | /api/v1/gamification/badges | List all badges |
| GET | /api/v1/gamification/badges/earned | Get earned badges |
| GET | /api/v1/gamification/challenges/daily | Get today's challenges |
| GET | /api/v1/gamification/challenges/weekly | Get weekly challenges |
| GET | /api/v1/gamification/challenges/{id}/progress | Get challenge progress |
| GET | /api/v1/gamification/leaderboard/{type} | Get leaderboard (xp, streams, etc.) |
| GET | /api/v1/gamification/levels | Get level definitions |

---

## 11. Analytics Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/analytics/me/overview | User's personal analytics |
| GET | /api/v1/analytics/me/top-tracks | User's most played tracks |
| GET | /api/v1/analytics/me/top-artists | User's most played artists |
| GET | /api/v1/analytics/me/listening-time | Listening time chart data |
| GET | /api/v1/analytics/artist/{id}/overview | Artist dashboard overview |
| GET | /api/v1/analytics/artist/{id}/streams | Artist stream chart data |
| GET | /api/v1/analytics/artist/{id}/listeners | Artist listener demographics |
| GET | /api/v1/analytics/artist/{id}/tracks | Artist per-track analytics |

---

## 12. WebSocket Endpoints

| Endpoint | Protocol | Purpose |
|----------|----------|---------|
| /ws/notifications | WSS | Real-time notifications |
| /ws/party/{id} | WSS | Listening party sync |
| /ws/room/{id} | WSS + WebRTC | Live room audio |
| /ws/feed | WSS | Live activity feed updates |

### WebSocket Message Format

```json
{
  "type": "event_type",
  "payload": { ... },
  "id": "msg_uuid",
  "timestamp": "2026-06-06T12:00:00Z"
}
```

---

## 13. Admin Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/admin/dashboard | Admin dashboard KPIs |
| GET | /api/v1/admin/users | List/manage users |
| POST | /api/v1/admin/users/{id}/ban | Ban user |
| POST | /api/v1/admin/users/{id}/role | Change user role |
| GET | /api/v1/admin/moderation/queue | Full moderation queue |
| GET | /api/v1/admin/moderation/stats | Moderation system stats |
| GET | /api/v1/admin/reports | Report management |
| GET | /api/v1/admin/gamification | Game config management |
| PATCH | /api/v1/admin/gamification/xp/{userId} | Adjust user XP |
| GET | /api/v1/admin/content/tracks | Track management |
| DELETE | /api/v1/admin/content/tracks/{id} | Remove track |
| GET | /api/v1/admin/system/health | System health check |
| GET | /api/v1/admin/system/metrics | System metrics |

---

## 14. API Rate Limits

| Endpoint Group | Auth Required | Limit | Burst |
|----------------|---------------|-------|-------|
| /auth/* | No | 10/min | 15 |
| /auth/register | No | 3/min | 5 |
| /users/* | Yes | 100/min | 150 |
| /tracks/* | Yes | 200/min | 300 |
| /search* | Yes | 60/min | 100 |
| /search/suggestions | Yes | 120/min | 200 |
| /social/* | Yes | 60/min | 100 |
| /recommendations/* | Yes | 30/min | 50 |
| /contributions | Yes | 20/min | 30 |
| /moderation/* | Yes | 60/min | 100 |
| /gamification/* | Yes | 30/min | 50 |
| /admin/* | Admin | 200/min | 300 |
| /analytics/* | Creator+ | 30/min | 50 |
| Static assets | No | 1000/min | 2000 |

---

## 15. Webhook API (For Creators)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/webhooks | Register webhook |
| GET | /api/v1/webhooks | List webhooks |
| DELETE | /api/v1/webhooks/{id} | Delete webhook |
| POST | /api/v1/webhooks/{id}/test | Send test event |

**Webhook Events:**
- `track.published` — New track release
- `stream.milestone` — Creator reached stream milestone
- `follower.milestone` — Creator reached follower milestone
- `contribution.approved` — Community contribution approved on creator's content

Event payload signed with HMAC-SHA256 using webhook secret.

---

## 16. API Health & Status

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/health | Basic health check |
| GET | /api/v1/health/ready | Readiness probe |
| GET | /api/v1/health/live | Liveness probe |
| GET | /api/v1/health/dependencies | Downstream dependency status |

---

## 17. Feature Flags

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/features | Return all module feature flags |

**GET /api/v1/features** — Returns the current state of all feature flags. Each flag is controlled by the corresponding `FEATURE_*_ENABLED` environment variable (default: `true`). When a feature is disabled, its backend routes are not registered and its frontend pages redirect to home.

```json
{
  "success": true,
  "data": {
    "analytics": true,
    "recommendation": true,
    "search": true,
    "social": true,
    "reactions": true,
    "creator": true,
    "moderation": true,
    "ai": true,
    "contribution": true,
    "gamification": true,
    "tips": true,
    "subscription": true,
    "notification": true
  }
}
```

### Available Flags

| Flag | Env Variable | Default | Description |
|------|-------------|---------|-------------|
| analytics | `FEATURE_ANALYTICS_ENABLED` | `true` | Listening stats, user analytics |
| recommendation | `FEATURE_RECOMMENDATION_ENABLED` | `true` | Track recommendations, mixes |
| search | `FEATURE_SEARCH_ENABLED` | `true` | Full-text & OpenSearch search |
| social | `FEATURE_SOCIAL_ENABLED` | `true` | Parties, rooms, clubs, comments |
| reactions | `FEATURE_REACTIONS_ENABLED` | `true` | Likes, saves, reposts |
| creator | `FEATURE_CREATOR_ENABLED` | `true` | Creator dashboard & analytics |
| moderation | `FEATURE_MODERATION_ENABLED` | `true` | Content moderation, flagging |
| ai | `FEATURE_AI_ENABLED` | `true` | AI playlist gen, mood explorer |
| contribution | `FEATURE_CONTRIBUTION_ENABLED` | `true` | Community metadata contributions |
| gamification | `FEATURE_GAMIFICATION_ENABLED` | `true` | XP, badges, leaderboards |
| tips | `FEATURE_TIPS_ENABLED` | `true` | Monetary tipping between users |
| subscription | `FEATURE_SUBSCRIPTION_ENABLED` | `true` | Paid subscription plans |
| notification | `FEATURE_NOTIFICATION_ENABLED` | `true` | Real-time WebSocket notifications |
