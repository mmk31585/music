# Database Architecture — Persian Music Ecosystem

> **Document**: Database Schema & Data Architecture
> **Status**: v1.0 — Final
> **Target**: 50M users, 10M tracks, 1B streams/month

---

## 1. Strategy

### 1.1 Principles

- **Domain-driven per-service databases** — each microservice owns its schema, no cross-service foreign keys
- **CQRS** — write to PostgreSQL, read from Redis (hot) + OpenSearch (search) + materialized views (analytics)
- **Partitioning** — large tables partitioned by time (listening_history monthly, analytics_events monthly, tracks yearly)
- **Indexing strategy** — B-tree for PKs/FKs/exact lookups, GiST for text search, BRIN for time-series
- **Connection pooling** — PgBouncer in transaction mode, max 100 pooled connections per service
- **Migration strategy** — golang-migrate with forward-only, backward-compatible changes, zero-downtime deploys

### 1.2 Data Distribution

| Database | Owner Service | Data Volume | Growth Rate |
|----------|---------------|-------------|-------------|
| user_db | User Service | 50M users | 500k/month |
| music_db | Music Service | 10M tracks, 2M artists | 200k tracks/month |
| social_db | Social Service | 500M edges | 10M/month |
| stream_db | Streams Service | 50B rows (partitioned) | 2B/month |
| contribution_db | Contribution Service | 50M submissions | 2M/month |
| moderation_db | Moderation Service | 10M actions | 500k/month |
| gamification_db | Gamification Service | 50M XP records | 10M/month |
| analytics_db | Analytics Service | Aggregated metrics | 100GB/month |

---

## 2. Schema Definitions

### 2.1 User Service (`user_db`)

```sql
-- Users core table
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username        VARCHAR(32) NOT NULL UNIQUE,
    email           VARCHAR(255) NOT NULL UNIQUE,
    phone           VARCHAR(20),
    password_hash   VARCHAR(255) NOT NULL,
    display_name    VARCHAR(100),
    avatar_url      TEXT,
    cover_url       TEXT,
    bio             TEXT,
    locale          VARCHAR(10) DEFAULT 'fa',
    is_verified     BOOLEAN DEFAULT false,
    is_creator      BOOLEAN DEFAULT false,
    role            VARCHAR(20) DEFAULT 'listener' -- guest, listener, creator, moderator, admin, superadmin
    trust_score     SMALLINT DEFAULT 0, -- -100 to 100
    level           INT DEFAULT 1,
    xp              BIGINT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_trust_score ON users(trust_score DESC);

-- User profiles (extended)
CREATE TABLE user_profiles (
    user_id         UUID PRIMARY KEY REFERENCES users(id),
    gender          VARCHAR(10),
    birth_date      DATE,
    country         VARCHAR(2), -- ISO 3166-1 alpha-2
    city            VARCHAR(100),
    preferred_lang  VARCHAR(10) DEFAULT 'fa',
    theme           VARCHAR(20) DEFAULT 'dark',
    explicit_filter BOOLEAN DEFAULT false,
    autoplay        BOOLEAN DEFAULT true,
    audio_quality   VARCHAR(10) DEFAULT 'high', -- low, normal, high, flac
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Social auth providers
CREATE TABLE user_auth_providers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    provider        VARCHAR(20) NOT NULL, -- google, spotify, github, apple
    provider_id     VARCHAR(255) NOT NULL,
    access_token    TEXT,
    refresh_token   TEXT,
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(provider, provider_id)
);

-- User sessions
CREATE TABLE user_sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    refresh_token   VARCHAR(512) NOT NULL,
    device_info     JSONB,
    ip_address      INET,
    user_agent      TEXT,
    last_used_at    TIMESTAMPTZ DEFAULT NOW(),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    expires_at      TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sessions_user ON user_sessions(user_id);
CREATE INDEX idx_sessions_expires ON user_sessions(expires_at) WHERE expires_at < NOW();

-- Follows (social graph)
CREATE TABLE follows (
    follower_id     UUID NOT NULL REFERENCES users(id),
    followee_id     UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (follower_id, followee_id)
);

CREATE INDEX idx_follows_followee ON follows(followee_id, created_at DESC);

-- Blocks
CREATE TABLE user_blocks (
    blocker_id      UUID NOT NULL REFERENCES users(id),
    blocked_id      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (blocker_id, blocked_id)
);
```

### 2.2 Music Service (`music_db`)

```sql
-- Artists
CREATE TABLE artists (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    persian_name    VARCHAR(255),
    bio             TEXT,
    persian_bio     TEXT,
    avatar_url      TEXT,
    cover_url       TEXT,
    genre_ids       UUID[],
    monthly_listeners BIGINT DEFAULT 0,
    follower_count  BIGINT DEFAULT 0,
    verified        BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_artists_name ON artists USING GIN (to_tsvector('simple', name));
CREATE INDEX idx_artists_monthly ON artists(monthly_listeners DESC);

-- Artist social links
CREATE TABLE artist_links (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artist_id       UUID NOT NULL REFERENCES artists(id),
    platform        VARCHAR(50) NOT NULL, -- instagram, twitter, youtube, spotify, website
    url             TEXT NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Tracks (partitioned by year)
CREATE TABLE tracks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(255) NOT NULL,
    persian_title   VARCHAR(255),
    artist_id       UUID NOT NULL REFERENCES artists(id),
    album_id        UUID REFERENCES albums(id),
    duration        INT NOT NULL, -- seconds
    track_number    SMALLINT,
    disc_number     SMALLINT DEFAULT 1,
    genre_id        UUID REFERENCES genres(id),
    subgenres       UUID[],
    moods           VARCHAR(50)[],
    bpm             SMALLINT,
    key             VARCHAR(10), -- musical key
    mode            VARCHAR(20), -- major, minor, dastgah value
    dastgah         VARCHAR(50), -- Persian modal system
    language        VARCHAR(10) DEFAULT 'fa',
    explicit        BOOLEAN DEFAULT false,
    lyrics          TEXT,
    lyrics_format   VARCHAR(20) DEFAULT 'lrc', -- lrc, plain, synced
    has_instrumental BOOLEAN DEFAULT false,
    audio_url       TEXT NOT NULL,
    audio_320_url   TEXT,
    audio_flac_url  TEXT,
    duration_ms     INT NOT NULL,
    play_count      BIGINT DEFAULT 0,
    like_count      BIGINT DEFAULT 0,
    comment_count   BIGINT DEFAULT 0,
    repost_count    BIGINT DEFAULT 0,
    is_published    BOOLEAN DEFAULT false,
    release_date    DATE,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
) PARTITION BY RANGE (release_date);

-- Create yearly partitions
CREATE TABLE tracks_2024 PARTITION OF tracks
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
CREATE TABLE tracks_2025 PARTITION OF tracks
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
CREATE TABLE tracks_2026 PARTITION OF tracks
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');
-- Default partition for older tracks
CREATE TABLE tracks_legacy PARTITION OF tracks DEFAULT;

CREATE INDEX idx_tracks_artist ON tracks(artist_id, release_date DESC);
CREATE INDEX idx_tracks_album ON tracks(album_id, track_number);
CREATE INDEX idx_tracks_title ON tracks USING GIN (to_tsvector('simple', title || ' ' || COALESCE(persian_title, '')));
CREATE INDEX idx_tracks_plays ON tracks(play_count DESC);
CREATE INDEX idx_tracks_release ON tracks(release_date DESC);
CREATE INDEX idx_tracks_mood ON tracks USING GIN(moods);
CREATE INDEX idx_tracks_dastgah ON tracks(dastgah) WHERE dastgah IS NOT NULL;

-- Albums
CREATE TABLE albums (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(255) NOT NULL,
    persian_title   VARCHAR(255),
    artist_id       UUID NOT NULL REFERENCES artists(id),
    cover_url       TEXT NOT NULL,
    cover_color     VARCHAR(7), -- extracted dominant color #RRGGBB
    album_type      VARCHAR(20) DEFAULT 'album', -- album, single, ep, compilation, live
    genre_id        UUID REFERENCES genres(id),
    label           VARCHAR(255),
    upc             VARCHAR(20),
    track_count     SMALLINT DEFAULT 0,
    duration_ms     INT DEFAULT 0,
    release_date    DATE,
    is_published    BOOLEAN DEFAULT false,
    play_count      BIGINT DEFAULT 0,
    like_count      BIGINT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_albums_artist ON albums(artist_id, release_date DESC);
CREATE INDEX idx_albums_release ON albums(release_date DESC);

-- Album artists (for compilations)
CREATE TABLE album_artists (
    album_id        UUID NOT NULL REFERENCES albums(id),
    artist_id       UUID NOT NULL REFERENCES artists(id),
    role            VARCHAR(50) DEFAULT 'main', -- main, featured, producer, writer
    PRIMARY KEY (album_id, artist_id)
);

-- Genres
CREATE TABLE genres (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    persian_name    VARCHAR(100),
    slug            VARCHAR(100) NOT NULL UNIQUE,
    parent_id       UUID REFERENCES genres(id),
    description     TEXT,
    image_url       TEXT,
    color           VARCHAR(7), -- accent color
    track_count     BIGINT DEFAULT 0,
    listener_count  BIGINT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_genres_parent ON genres(parent_id);

-- Playlists
CREATE TABLE playlists (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    cover_url       TEXT,
    cover_color     VARCHAR(7),
    owner_id        UUID NOT NULL REFERENCES users(id),
    is_public       BOOLEAN DEFAULT true,
    is_collaborative BOOLEAN DEFAULT false,
    is_curated      BOOLEAN DEFAULT false, -- official editorial playlists
    track_count     SMALLINT DEFAULT 0,
    duration_ms     BIGINT DEFAULT 0,
    follower_count  BIGINT DEFAULT 0,
    like_count      BIGINT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_playlists_owner ON playlists(owner_id);
CREATE INDEX idx_playlists_followers ON playlists(follower_count DESC);

-- Playlist tracks (ordered)
CREATE TABLE playlist_tracks (
    playlist_id     UUID NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    track_id        UUID NOT NULL REFERENCES tracks(id),
    position        INT NOT NULL,
    added_by        UUID REFERENCES users(id),
    added_at        TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (playlist_id, position)
);

CREATE INDEX idx_playlist_tracks_playlist ON playlist_tracks(playlist_id, position);

-- User likes
CREATE TABLE user_likes (
    user_id         UUID NOT NULL REFERENCES users(id),
    track_id        UUID NOT NULL REFERENCES tracks(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, track_id)
);

CREATE INDEX idx_likes_track ON user_likes(track_id, created_at DESC);

-- User library (saved albums, playlists)
CREATE TABLE user_library (
    user_id         UUID NOT NULL REFERENCES users(id),
    item_type       VARCHAR(20) NOT NULL, -- album, playlist, artist, podcast
    item_id         UUID NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, item_type, item_id)
);
```

### 2.3 Streams Service (`stream_db`)

```sql
-- Listening history (partitioned monthly)
CREATE TABLE listening_history (
    id              BIGSERIAL,
    user_id         UUID NOT NULL,
    track_id        UUID NOT NULL,
    album_id        UUID,
    artist_id       UUID NOT NULL,
    played_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    duration_played INT NOT NULL, -- seconds actually played
    duration_total  INT NOT NULL, -- total track seconds
    completion_rate FLOAT GENERATED ALWAYS AS (duration_played::FLOAT / NULLIF(duration_total, 0)) STORED,
    source          VARCHAR(50) NOT NULL, -- playlist, album, search, radio, recommendation
    device_type     VARCHAR(50), -- mobile, desktop, tablet, tv
    os              VARCHAR(50),
    country         VARCHAR(2),
    city            VARCHAR(100),
    ip_address      INET,
    user_agent      TEXT,
    session_id      UUID,
    is_offline      BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (played_at);

-- Monthly partitions (created by cron job)
-- CREATE TABLE listening_history_2025_01 PARTITION OF listening_history
--     FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
-- CREATE TABLE listening_history_2025_02 PARTITION OF listening_history
--     FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
-- ... automated via pg_partman

CREATE INDEX idx_history_user ON listening_history(user_id, played_at DESC);
CREATE INDEX idx_history_track ON listening_history(track_id, played_at DESC);
CREATE INDEX idx_history_artist ON listening_history(artist_id, played_at DESC);
CREATE INDEX idx_history_played ON listening_history(played_at DESC);

-- Stream analytics (aggregated hourly)
CREATE TABLE stream_analytics_hourly (
    bucket          TIMESTAMPTZ NOT NULL, -- truncated to hour
    track_id        UUID NOT NULL,
    artist_id       UUID NOT NULL,
    album_id        UUID,
    stream_count    BIGINT DEFAULT 0,
    unique_listeners BIGINT DEFAULT 0,
    total_duration  BIGINT DEFAULT 0, -- total seconds streamed
    avg_completion  FLOAT DEFAULT 0,
    source          VARCHAR(50),
    country         VARCHAR(2),
    PRIMARY KEY (bucket, track_id, source, country)
);

CREATE INDEX idx_stream_analytics_track ON stream_analytics_hourly(track_id, bucket DESC);
CREATE INDEX idx_stream_analytics_artist ON stream_analytics_hourly(artist_id, bucket DESC);

-- Daily artist analytics
CREATE TABLE artist_analytics_daily (
    date            DATE NOT NULL,
    artist_id       UUID NOT NULL,
    stream_count    BIGINT DEFAULT 0,
    unique_listeners BIGINT DEFAULT 0,
    follower_count  INT DEFAULT 0,
    follower_growth INT DEFAULT 0, -- net change
    avg_completion  FLOAT DEFAULT 0,
    top_track_id    UUID,
    top_country     VARCHAR(2),
    PRIMARY KEY (date, artist_id)
);
```

### 2.4 Social Service (`social_db`)

```sql
-- Listening parties
CREATE TABLE listening_parties (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id         UUID NOT NULL REFERENCES users(id),
    title           VARCHAR(255),
    track_id        UUID REFERENCES tracks(id),
    playlist_id     UUID REFERENCES playlists(id),
    is_live         BOOLEAN DEFAULT true,
    started_at      TIMESTAMPTZ DEFAULT NOW(),
    ended_at        TIMESTAMPTZ,
    participant_count INT DEFAULT 0,
    max_participants INT DEFAULT 50,
    is_public       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_parties_live ON listening_parties(is_live, started_at DESC) WHERE is_live = true;

-- Party participants
CREATE TABLE party_participants (
    party_id        UUID NOT NULL REFERENCES listening_parties(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id),
    joined_at       TIMESTAMPTZ DEFAULT NOW(),
    is_host         BOOLEAN DEFAULT false,
    sync_offset_ms  INT DEFAULT 0,
    PRIMARY KEY (party_id, user_id)
);

-- Live rooms
CREATE TABLE live_rooms (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id         UUID NOT NULL REFERENCES users(id),
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    cover_url       TEXT,
    is_live         BOOLEAN DEFAULT true,
    listener_count  INT DEFAULT 0,
    max_listeners   INT DEFAULT 500,
    tags            VARCHAR(50)[],
    started_at      TIMESTAMPTZ DEFAULT NOW(),
    ended_at        TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_liverooms_live ON live_rooms(is_live, listener_count DESC) WHERE is_live = true;

-- Music clubs
CREATE TABLE music_clubs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    cover_url       TEXT,
    owner_id        UUID NOT NULL REFERENCES users(id),
    genre_id        UUID REFERENCES genres(id),
    is_public       BOOLEAN DEFAULT true,
    member_count    INT DEFAULT 0,
    post_count      INT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_clubs_members ON music_clubs(member_count DESC);

-- Club members
CREATE TABLE club_members (
    club_id         UUID NOT NULL REFERENCES music_clubs(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id),
    role            VARCHAR(20) DEFAULT 'member', -- member, moderator, admin, owner
    joined_at       TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (club_id, user_id)
);

-- Activity feed
CREATE TABLE activity_events (
    id              BIGSERIAL,
    user_id         UUID NOT NULL,
    event_type      VARCHAR(50) NOT NULL, -- listen, like, follow, playlist_create, share, comment, review
    target_type     VARCHAR(20), -- track, album, artist, playlist, user
    target_id       UUID,
    metadata        JSONB, -- flexible payload
    created_at      TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Monthly partitions for activity (high volume)

CREATE INDEX idx_activity_user ON activity_events(user_id, created_at DESC);
CREATE INDEX idx_activity_type ON activity_events(event_type, created_at DESC);
CREATE INDEX idx_activity_feed ON activity_events(created_at DESC);

-- Track discussions / comments
CREATE TABLE track_comments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    track_id        UUID NOT NULL REFERENCES tracks(id),
    user_id         UUID NOT NULL REFERENCES users(id),
    parent_id       UUID REFERENCES track_comments(id), -- for replies
    body            TEXT NOT NULL,
    like_count      BIGINT DEFAULT 0,
    reply_count     BIGINT DEFAULT 0,
    is_edited       BOOLEAN DEFAULT false,
    is_deleted      BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_comments_track ON track_comments(track_id, created_at DESC);
CREATE INDEX idx_comments_parent ON track_comments(parent_id);

-- Track reviews / ratings
CREATE TABLE track_ratings (
    user_id         UUID NOT NULL REFERENCES users(id),
    track_id        UUID NOT NULL REFERENCES tracks(id),
    rating          SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 10),
    review          TEXT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, track_id)
);

CREATE INDEX idx_ratings_track ON track_ratings(track_id, rating DESC);
```

### 2.5 Contribution Service (`contribution_db`)

```sql
-- Contributions
CREATE TABLE contributions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contributor_id  UUID NOT NULL REFERENCES users(id),
    contribution_type VARCHAR(50) NOT NULL, -- lyrics, translation, credits, metadata, album_art, bio
    target_type     VARCHAR(20) NOT NULL, -- track, album, artist
    target_id       UUID NOT NULL,
    locale          VARCHAR(10), -- for translations
    data            JSONB NOT NULL, -- the contribution payload
    status          VARCHAR(20) DEFAULT 'pending', -- pending, approved, rejected, needs_review
    ai_verdict      VARCHAR(20), -- pass, fail, review
    ai_confidence   FLOAT, -- 0-1
    ai_reason       TEXT,
    moderator_id    UUID REFERENCES users(id),
    moderator_note  TEXT,
    version         INT DEFAULT 1,
    is_minor        BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    decided_at      TIMESTAMPTZ,
    applied_at      TIMESTAMPTZ
);

CREATE INDEX idx_contributions_status ON contributions(status, created_at DESC);
CREATE INDEX idx_contributions_type ON contributions(contribution_type, status);
CREATE INDEX idx_contributions_target ON contributions(target_type, target_id);
CREATE INDEX idx_contributions_contributor ON contributions(contributor_id, created_at DESC);

-- Contribution history (version tracking)
CREATE TABLE contribution_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contribution_id UUID NOT NULL REFERENCES contributions(id),
    data            JSONB NOT NULL, -- snapshot of what changed
    previous_data   JSONB, -- previous state (for rollback)
    changed_by      UUID NOT NULL REFERENCES users(id),
    change_type     VARCHAR(20) NOT NULL, -- create, approve, reject, revert
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Content versions (published state snapshots)
CREATE TABLE content_versions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_type     VARCHAR(20) NOT NULL,
    target_id       UUID NOT NULL,
    version         INT NOT NULL,
    data            JSONB NOT NULL, -- full snapshot
    applied_by      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(target_type, target_id, version)
);
```

### 2.6 Moderation Service (`moderation_db`)

```sql
-- Moderation queue
CREATE TABLE moderation_queue (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type    VARCHAR(20) NOT NULL, -- contribution, comment, playlist, track, user, report
    content_id      UUID NOT NULL, -- polymorphic reference
    reporter_id     UUID REFERENCES users(id),
    reason          TEXT,
    category        VARCHAR(50), -- spam, profanity, copyright, harassment, misinformation
    priority        VARCHAR(10) DEFAULT 'normal', -- low, normal, high, critical
    status          VARCHAR(20) DEFAULT 'pending', -- pending, reviewing, resolved, escalated
    ai_verdict      JSONB, -- AI moderation result
    assigned_to     UUID REFERENCES users(id),
    resolved_by     UUID REFERENCES users(id),
    resolution      VARCHAR(20), -- approved, rejected, removed, warned
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    resolved_at     TIMESTAMPTZ
);

CREATE INDEX idx_modqueue_status ON moderation_queue(status, priority, created_at DESC);
CREATE INDEX idx_modqueue_content ON moderation_queue(content_type, content_id);

-- Reports
CREATE TABLE reports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_id     UUID NOT NULL REFERENCES users(id),
    content_type    VARCHAR(20) NOT NULL,
    content_id      UUID NOT NULL,
    reason          TEXT NOT NULL,
    category        VARCHAR(50) NOT NULL,
    status          VARCHAR(20) DEFAULT 'open', -- open, investigating, resolved, dismissed
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- User warnings
CREATE TABLE user_warnings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    issued_by       UUID NOT NULL REFERENCES users(id),
    reason          TEXT NOT NULL,
    points          SMALLINT DEFAULT 1, -- accumulates toward ban
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_warnings_user ON user_warnings(user_id, created_at DESC);
```

### 2.7 Gamification Service (`gamification_db`)

```sql
-- XP transactions (ledger)
CREATE TABLE xp_transactions (
    id              BIGSERIAL,
    user_id         UUID NOT NULL REFERENCES users(id),
    amount          INT NOT NULL,
    balance_after   INT NOT NULL,
    source          VARCHAR(50) NOT NULL, -- listen, like, share, contribute, moderate, review, daily_login
    metadata        JSONB,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_xp_user ON xp_transactions(user_id, created_at DESC);
CREATE INDEX idx_xp_source ON xp_transactions(source, created_at DESC);

-- XP sources and rewards
-- listen: +1/track, like: +5, share: +10, contribute(approved): +50,
-- translate: +30, moderate: +20, review: +10, daily_login: +15,
-- first_stream_day: +25, playlist_create: +20, invite_friend: +100

-- Level definitions
CREATE TABLE level_definitions (
    level           INT PRIMARY KEY,
    xp_required     BIGINT NOT NULL,
    title           VARCHAR(100), -- "Novice Listener", "Melody Master", ...
    title_persian   VARCHAR(100),
    privileges      JSONB, -- { "create_playlist": true, "upload_track": false, ... }
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Badges
CREATE TABLE badges (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    icon_url        TEXT,
    category        VARCHAR(20) NOT NULL, -- listener, contributor, social, special
    rarity          VARCHAR(20) DEFAULT 'common', -- common, rare, epic, legendary, hidden
    criteria        JSONB NOT NULL, -- { "type": "stream_count", "target": 10000 }
    xp_reward       INT DEFAULT 0,
    is_hidden       BOOLEAN DEFAULT false,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- User badges
CREATE TABLE user_badges (
    user_id         UUID NOT NULL REFERENCES users(id),
    badge_id        UUID NOT NULL REFERENCES badges(id),
    earned_at        TIMESTAMPTZ DEFAULT NOW(),
    is_displayed    BOOLEAN DEFAULT true,
    PRIMARY KEY (user_id, badge_id)
);

-- Daily challenges
CREATE TABLE daily_challenges (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    challenge_type  VARCHAR(50) NOT NULL, -- stream_count, explore_genre, contribute, social_share
    target_count    INT NOT NULL,
    xp_reward       INT NOT NULL,
    badge_reward_id UUID REFERENCES badges(id),
    is_active       BOOLEAN DEFAULT true,
    valid_from      DATE NOT NULL,
    valid_until     DATE NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- User challenge progress
CREATE TABLE user_challenges (
    user_id         UUID NOT NULL REFERENCES users(id),
    challenge_id    UUID NOT NULL REFERENCES daily_challenges(id),
    progress        INT DEFAULT 0,
    is_completed    BOOLEAN DEFAULT false,
    completed_at    TIMESTAMPTZ,
    PRIMARY KEY (user_id, challenge_id)
);

-- Leaderboard snapshots (weekly)
CREATE TABLE leaderboard_snapshots (
    id              BIGSERIAL,
    leaderboard_type VARCHAR(30) NOT NULL, -- xp_weekly, xp_monthly, streams, contributions
    rank            INT NOT NULL,
    user_id         UUID NOT NULL REFERENCES users(id),
    score           BIGINT NOT NULL,
    snapshot_date   DATE NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(leaderboard_type, snapshot_date, rank)
);

CREATE INDEX idx_leaderboard_snapshot ON leaderboard_snapshots(leaderboard_type, snapshot_date, rank);
```

### 2.8 Analytics Service (`analytics_db`)

```sql
-- Analytics events (raw, partitioned monthly)
CREATE TABLE analytics_events (
    id              BIGSERIAL,
    event_type      VARCHAR(50) NOT NULL, -- page_view, search, click, stream_start, stream_end, share, login
    user_id         UUID,
    session_id      UUID,
    page            VARCHAR(255),
    referrer        VARCHAR(255),
    metadata        JSONB,
    device_type     VARCHAR(50),
    os              VARCHAR(50),
    browser         VARCHAR(50),
    ip_address      INET,
    country         VARCHAR(2),
    city            VARCHAR(100),
    created_at      TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (created_at);

CREATE INDEX idx_analytics_event ON analytics_events(event_type, created_at DESC);
CREATE INDEX idx_analytics_page ON analytics_events(page, created_at DESC);

-- Pre-aggregated dashboards (materialized views refreshed hourly)
-- Active users, top tracks, top artists, top countries, retention, conversion
```

---

## 3. Partitioning Strategy

| Table | Partition Key | Interval | Partitions | Retention |
|-------|--------------|----------|------------|-----------|
| tracks | release_date | Yearly | 10+ | Permanent |
| listening_history | played_at | Monthly | 36+ | 3 years → aggregate |
| activity_events | created_at | Monthly | 36+ | 1 year → delete |
| analytics_events | created_at | Monthly | 12+ | 6 months → delete |

Partitions are managed by `pg_partman` with automatic creation 3 months ahead and automatic detach/delete after retention period.

---

## 4. Indexing Strategy

### 4.1 B-tree Indexes
- Primary keys, unique constraints
- Foreign keys (though no cross-service FKs enforced)
- Columns used in `WHERE`, `ORDER BY`, `JOIN` with high cardinality
- Composite indexes for common query patterns: `(user_id, created_at DESC)`

### 4.2 GIN / GiST Indexes
- Full-text search on `tracks.title`, `artists.name`, `albums.title` (GIN with `to_tsvector`)
- Array columns: `tracks.moods`, `tracks.subgenres` (GIN)
- JSONB columns on `contributions.data`, `activity_events.metadata` (GIN)

### 4.3 BRIN Indexes
- Time-series columns on partitioned tables: `played_at`, `created_at`
- Much smaller than B-tree for naturally-ordered append-only data

### 4.4 Partial Indexes
- `WHERE deleted_at IS NULL` on soft-delete tables
- `WHERE is_live = true` on listening_parties and live_rooms
- `WHERE status = 'pending'` on moderation_queue

---

## 5. Caching Strategy (Redis)

### 5.1 Cache Layers

| Layer | TTL | Data | Pattern |
|-------|-----|------|---------|
| L1 (hot) | 30s-2min | Trending tracks, top charts, current user session | Write-through |
| L2 (warm) | 5-15min | Artist pages, album pages, playlist metadata | Cache-aside |
| L3 (cold) | 1-6h | Genre pages, search suggestions, editorial content | Cache-aside |
| Session | 15min | Access tokens, refresh tokens | Write-through |

### 5.2 Key Naming Convention

```
{service}:{entity}:{id}[:field]
```

Examples:
```
music:track:a1b2c3d4
music:artist:e5f6g7h8:top-tracks
music:album:i9j0k1l2:tracks
user:session:token_abc123
social:feed:user_x1y2z3:page:1
analytics:trending:tracks:ir
```

### 5.3 Cache Invalidation

- **Write-through**: On successful write to PostgreSQL, update cache immediately
- **Event-driven**: On Kafka event, invalidate related cache keys
- **TTL expiry**: Automatic expiration based on data staleness tolerance
- **Manual**: Admin action to flush specific cache keys

### 5.4 Redis Data Structures

- **Strings**: Session tokens, API responses (JSON), rate limits
- **Sorted Sets**: Trending tracks (score = weighted stream count), leaderboards
- **Sets**: User likes (for fast membership check), user sessions
- **Lists**: Recent listens, activity feed (capped at 100)
- **Hashes**: User profiles, track metadata (for partial updates)

---

## 6. Full-Text Search (OpenSearch)

### 6.1 Index Mapping Strategy

**tracks index:**
```json
{
  "settings": {
    "analysis": {
      "analyzer": {
        "persian_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "arabic_normalization", "persian_normalization", "stop_persian"]
        },
        "autocomplete": {
          "tokenizer": "edge_ngram",
          "filter": ["lowercase"]
        }
      },
      "normalization": {
        "persian": {
          "type": "persian_normalization"
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "title": { "type": "text", "analyzer": "persian_analyzer", "fields": { "autocomplete": { "type": "text", "analyzer": "autocomplete" } } },
      "persian_title": { "type": "text", "analyzer": "persian_analyzer", "boost": 2 },
      "artist_name": { "type": "text", "analyzer": "persian_analyzer", "boost": 1.5 },
      "album_name": { "type": "text", "analyzer": "persian_analyzer" },
      "lyrics": { "type": "text", "analyzer": "persian_analyzer" },
      "genre": { "type": "keyword" },
      "tags": { "type": "keyword" },
      "release_date": { "type": "date" },
      "play_count": { "type": "long" },
      "language": { "type": "keyword" },
      "dastgah": { "type": "keyword" }
    }
  }
}
```

### 6.2 Search Features

- **Autocomplete**: edge_ngram on title + persian_title (2-20 chars)
- **Transliteration**: English → Farsi phonetic matching (e.g., "gol" → "گل")
- **Synonym expansion**: Persian synonym filter (e.g., "عاشقانه" ↔ "romantic")
- **Fuzzy matching**: Levenshtein distance 1-2 for typos
- **Function score**: Boost by play_count, recency, user preference
- **Faceted search**: By genre, mood, dastgah, language, year
- **Multi-index search**: Across tracks, artists, albums, playlists, lyrics

### 6.3 Query Types

- `multi_match` with `best_fields` for primary search
- `term` + `terms` for faceted filtering
- `match_phrase` for exact title matching
- `function_score` for popularity-weighted results
- `more_like_this` for recommendations ("more like this track")

---

## 7. Migration Strategy

### 7.1 Tooling

- Use `golang-migrate` for SQL migrations
- Migrations are forward-only (no down migrations in production)
- Rollback via new migration that reverses changes

### 7.2 Zero-Downtime Principles

1. **Additive changes only**: Never remove columns in a single deployment
2. **Expand-contract pattern**: Add new column → dual-write → backfill → remove old column
3. **Index creation**: `CONCURRENTLY` for production indexes
4. **Column changes**: Add nullable column first, backfill, then add NOT NULL
5. **Table changes**: Create new table, dual-write, backfill, swap via trigger

### 7.3 Migration Naming

```
{YYYYMMDDHHMMSS}_{description}.sql
e.g., 20250601000001_create_users.up.sql
```

---

## 8. Backup & Disaster Recovery

| Component | Backup | RPO | RTO |
|-----------|--------|-----|-----|
| PostgreSQL | WAL streaming + pg_dump daily | 1 sec (WAL) | 1 hour |
| Redis | RDB snapshots every 5min | 5 min | 15 min |
| OpenSearch | Snapshot to S3 daily | 1 hour | 2 hours |
| S3/MinIO | Cross-region replication | 15 min | 30 min |
| Kafka | Log retention 7 days | 7 days | 1 hour |

- DR site in different region with warm standby
- Regular DR drills (quarterly)
- Automated failover testing (monthly)
