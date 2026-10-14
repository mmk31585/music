# Muse — Documentation

> Premium Persian Music Ecosystem
> "Feel the music. Share the moment."

---

## The Plan (start here)

One spine, four documents. `MASTER_PLAN.md` is the map; the three below are the territory, each owning exactly one concern. Read the map first.

```
MASTER_PLAN.md                    ← the map — why/what (CEO), how (CTO), when (Dev), governance, success/failure
├── ROADMAP.md                    ← territory: phases · tests · CI · exit criteria · week-1 sprint
├── ARCHITECTURE_REPOS.md         ← territory: polyrepo · contracts & the seam · topology · triggers · migration
└── ENGINEERING_OPERATING_SYSTEM.md ← territory: SRE · culture · strategy · post-mortem
```

| Doc | Read it to answer | Skip it if |
|---|---|---|
| [MASTER_PLAN.md](MASTER_PLAN.md) | "What are we building and why?" "How do the three lenses stay aligned?" | you're mid-execution on a specific task |
| [ROADMAP.md](ROADMAP.md) | "What ships this week/phase?" "What tests and CI gates?" "Am I done?" | you need to know how repos are structured |
| [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md) | "Why polyrepo?" "How do contracts keep 4 languages aligned?" "When does Kafka/ClickHouse/Rust arrive?" | you're deciding product scope |
| [ENGINEERING_OPERATING_SYSTEM.md](ENGINEERING_OPERATING_SYSTEM.md) | "How do we own what we ship?" "SLOs, post-mortems, tech debt, culture, strategy" | you're picking phase deliverables |

---

## Documentation Index (project reference)

### Architecture

| Document | Description |
|----------|-------------|
| [`architecture/OVERVIEW.md`](architecture/OVERVIEW.md) | High-level architecture, navigation model, tech stack, request flow |
| [`architecture/API.md`](architecture/API.md) | API design, endpoints, conventions, rate limits, feature flags |
| [`architecture/SERVICES.md`](architecture/SERVICES.md) | Microservices, event bus, recommendation, search, moderation, gamification |
| [`architecture/DATABASE.md`](architecture/DATABASE.md) | Database schema, partitioning, indexing, caching, search engine mapping |
| [`architecture/OPS.md`](architecture/OPS.md) | Deployment, scaling, monitoring, CI/CD, security, incident response |

### Design & UX

| Document | Description |
|----------|-------------|
| [`DESIGN_SYSTEM.md`](DESIGN_SYSTEM.md) | Visual design tokens, component specs, player UX, navigation, motion system |

### API Documentation

| Document | Description |
|----------|-------------|
| [`SWAGGER.md`](SWAGGER.md) | Swagger/OpenAPI handler annotation guide |

### Frontend

| Document | Description |
|----------|-------------|
| [`../frontend/AGENTS.md`](../frontend/AGENTS.md) | Agent team architecture & file ownership |
| [`../frontend/PERFORMANCE.md`](../frontend/PERFORMANCE.md) | Performance-sensitive code guidelines |
| [`../frontend/README.md`](../frontend/README.md) | Frontend setup, development, CI/CD pipeline |
| [`../FRONTEND_ARCHITECTURE.md`](../FRONTEND_ARCHITECTURE.md) | Frontend architecture & product plan |
| [`../TODO.md`](../TODO.md) | Active project TODO items & known issues |
