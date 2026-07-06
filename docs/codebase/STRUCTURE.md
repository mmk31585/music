# Muse — Project Structure

```
music/
├── cmd/                          # Go entrypoints
│   ├── api/                      #   → main.go (API server)
│   └── worker/                   #   → main.go (background worker)
│
├── internal/                     # Go internal packages
│   ├── app/                      # App bootstrap, DI container, routes
│   │   ├── app.go                #   App struct (config, DB, Redis, logger, lifecycle)
│   │   ├── bootstrap.go          #   Bootstrap() initialization
│   │   ├── container.go          #   DI container (706 lines, manual wiring)
│   │   ├── http.go               #   HTTP server creation
│   │   ├── routes.go             #   All route registration (public/auth/admin/internal)
│   │   ├── media_routes.go       #   Media serving routes
│   │   ├── shutdown.go           #   Graceful shutdown
│   │   ├── worker_container.go   #   Background worker container
│   │   └── app_test.go           #   Tests
│   │
│   ├── common/                   # Shared utilities
│   │   ├── errors/               #   Domain error types
│   │   ├── middleware/            #   Auth, CORS, security, logging, HMAC
│   │   ├── pagination/           #   Paginated response helpers
│   │   ├── request/              #   Request parsing
│   │   ├── response/             #   Standardized JSON response format
│   │   └── validator/            #   Input validation wrapper
│   │
│   ├── config/                   # App configuration loading
│   │
│   ├── modules/                  # 30 domain modules
│   │   ├── ai/                   #   OpenAI embeddings, mood analysis, song matching
│   │   ├── analytics/            #   Event-driven analytics (plays, signups)
│   │   ├── auth/                 #   JWT auth, login/register, middleware
│   │   ├── catalog/              #   Core music catalog (track/album/artist/genre)
│   │   ├── contribution/         #   User-contributed metadata (AI-verified)
│   │   ├── covers/               #   Album cover image handling
│   │   ├── creator/              #   Creator dashboard & tools
│   │   ├── dashboard/            #   Admin dashboard stats
│   │   ├── features/             #   Feature flag toggles
│   │   ├── follow/               #   Follow/unfollow artists/users
│   │   ├── gamification/         #   Badges, XP, leaderboards
│   │   ├── health/               #   Health check endpoint
│   │   ├── history/              #   Listening history (signal-classified)
│   │   ├── importcmd/            #   Batch import (iTunes, Deezer, Spotify, etc.)
│   │   ├── ingestion/            #   Media ingestion pipeline
│   │   ├── library/              #   User music library & favorites
│   │   ├── lyrics/               #   Lyrics storage, ML client, OpenRouter AI
│   │   ├── media/                #   File upload management
│   │   ├── moderation/           #   Content moderation
│   │   ├── notification/         #   Push/in-app notifications
│   │   ├── player/               #   Audio streaming, playback state
│   │   ├── playlist/             #   Playlist CRUD, collaborative, WebSocket
│   │   ├── queue/                #   Play queue management
│   │   ├── reactions/            #   Like/reaction system
│   │   ├── recommendation/       #   Taste profiles, home feed, discover weekly, radio
│   │   ├── search/               #   OpenSearch + Persian-aware ranking
│   │   ├── social/               #   Parties, clubs, stages, discussions
│   │   ├── subscription/         #   Premium subscriptions
│   │   ├── tips/                 #   Creator tipping system
│   │   └── video/                #   Video content (CRUD, comments, ML)
│   │
│   ├── pkg/
│   │   └── audioinfo/            # Audio file analysis
│   │
│   ├── platform/                 # Infrastructure abstractions
│   │   ├── audio/                #   Audio analyzer
│   │   ├── cache/                #   Redis cache layer
│   │   ├── circuitbreaker/       #   Circuit breaker for external calls
│   │   ├── database/             #   PostgreSQL pool + migrations
│   │   ├── events/               #   In-process bus + transactional outbox
│   │   ├── fuzzy/                #   Fuzzy string matching
│   │   ├── logger/               #   Zap logger init
│   │   ├── metrics/              #   Prometheus metrics
│   │   ├── opensearch/           #   OpenSearch client + index mgmt
│   │   ├── payment/              #   Payment gateways (IDPay, ZarinPal)
│   │   ├── storage/              #   Local / S3 file storage abstraction
│   │   ├── web/                  #   User context from JWT
│   │   └── ws/                   #   WebSocket hub + presence
│   │
│   └── workers/                  # Background workers
│       ├── ai/                   #   AI/ML background tasks
│       ├── analytics/            #   Analytics aggregation
│       ├── cleanup/              #   Data cleanup jobs
│       ├── eventbus/             #   Event bus consumer
│       ├── indexer/              #   Search indexer
│       ├── notification/         #   Notification dispatch
│       ├── recommendation/       #   Recommendation computation
│       ├── recommender_v2/       #   V2 recommendation engine
│       └── transcoder/           #   Audio transcoding
│
├── frontend/                     # Vue 3 SPA
│   └── src/
│       ├── App.vue               # Root component (layout switching)
│       ├── main.ts               # App entry (Pinia, Router, PrimeVue, fonts)
│       ├── assets/               # Static assets (fonts, images, styles)
│       ├── components/           # UI components
│       │   ├── admin/            #   Admin panel components
│       │   ├── auth/             #   Auth form components
│       │   ├── common/           #   Shared/common components
│       │   ├── contribution/     #   Contribution UI
│       │   ├── creator/          #   Creator dashboard components
│       │   ├── forms/            #   Form components
│       │   ├── gamification/     #   Badges, XP, challenges UI
│       │   ├── layouts/          #   Layout shells (Empty, Sidebar, Admin)
│       │   ├── music/            #   Core music components (player, cards, rows)
│       │   ├── profile/          #   Profile components
│       │   ├── social/           #   Social feature components
│       │   ├── tips/             #   Tipping UI
│       │   ├── video/            #   Video components
│       │   └── widgets/          #   Reusable widgets
│       │
│       ├── composables/          # Vue composables
│       │   ├── auth/             #   Auth composables
│       │   ├── admin/            #   Admin composables
│       │   ├── catalog/          #   Catalog composables
│       │   ├── lyrics/           #   Lyrics composables
│       │   ├── media/            #   Media upload composables
│       │   ├── player/           #   Player composables
│       │   ├── recommendation/   #   Recommendation composables
│       │   ├── search/           #   Search composables
│       │   └── social/           #   Social composables
│       │   (plus root-level: useRequest, useRTL, useAuth, etc.)
│       │
│       ├── layouts/              # Layout components
│       │   ├── LayoutAuth.vue    #   Auth layout
│       │   ├── LayoutMusicApp.vue#   Main app shell
│       │   └── LayoutAdmin.vue   #   Admin shell
│       │
│       ├── pages/                # Page components
│       │   ├── admin/            #   ~20 admin pages
│       │   ├── app/              #   ~30 app pages (home, search, library, etc.)
│       │   ├── auth/             #   Login + Register
│       │   ├── errors/           #   404
│       │   └── onboarding/       #   Genre onboarding
│       │
│       ├── plugins/              # Plugin configurations
│       │   ├── client/           #   Axios client + request wrapper
│       │   └── query-builder/    #   URL query param builder
│       │
│       ├── router/               # Vue Router config
│       │   ├── index.ts          #   Router creation + guard chain
│       │   ├── middleware/       #   Guards (auth, maintenance, login)
│       │   └── routes/           #   Route definitions (auth, app, admin)
│       │
│       ├── services/             # Service layer
│       │   ├── api/              #   25 domain API modules
│       │   ├── player/           #   Audio engine, player engine, queue manager
│       │   ├── socket/           #   WebSocket client
│       │   ├── storage/          #   SSR-safe localStorage
│       │   └── analytics/        #   Analytics service
│       │
│       ├── stores/               # Pinia stores
│       │   ├── user-auth.ts      #   Auth tokens, user profile, admin check
│       │   ├── player.ts         #   Audio playback state
│       │   ├── feature-flags.ts  #   Feature flag cache
│       │   ├── maintenance.ts    #   Maintenance mode
│       │   └── page-loader.ts    #   Global loading state
│       │
│       ├── types/                # TypeScript declarations
│       │   └── document-pip.ts   #   PiP API types
│       │
│       ├── utils/                # Utility functions
│       └── tests/                # Test setup + helpers
│
├── moja-ml-service/              # Python ML microservice
│   ├── app/                      # FastAPI app
│   ├── alembic/                  # DB migrations
│   ├── tests/                    # Python tests
│   ├── Dockerfile
│   └── pyproject.toml
│
├── migrations/                   # 57 SQL migration files
├── deployments/                  # Docker Compose + deployment configs
│   └── docker-compose.yml
├── docs/                         # Documentation
│   ├── architecture/             #   Architecture docs (5 files)
│   ├── DESIGN_SYSTEM.md          #   Complete design system spec
│   └── codebase/                 #   Generated codebase docs
│
├── scripts/                      # Shell scripts
│   ├── dev.sh                    #   Unified dev server
│   ├── migrate.sh                #   DB migration runner
│   ├── watchdog.sh               #   Service health monitor
│   ├── autopilot.sh
│   └── start-proxy.sh
│
├── .github/                      # GitHub config
│   └── copilot-instructions.md   # AI coding standards
│
├── Dockerfile                    # Multi-stage Go build
├── Dockerfile.worker             # Worker-specific build
├── Makefile                      # 18 targets
├── go.mod / go.sum
├── AGENTS.md                     # AI agent system
├── CONTEXT.md                    # Domain glossary
├── FRONTEND_ARCHITECTURE.md      # Frontend architecture doc
├── TODO.md                       # Living TODO/critical issues
└── README.md                     # Project overview
```
