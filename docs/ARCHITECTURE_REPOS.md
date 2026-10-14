# Muse — Polyrepo Architecture

> **Territory:** polyrepo map · contracts & the seam · runtime topology · data-pipeline evolution · deferred-infra triggers · migration.
> This is the *how-it-fits-together* layer.
> - Map (why/what/execution): [MASTER_PLAN.md](MASTER_PLAN.md)
> - Execution (phases/tests/CI): [ROADMAP.md](ROADMAP.md)
> - Culture/SRE/strategy: [ENGINEERING_OPERATING_SYSTEM.md](ENGINEERING_OPERATING_SYSTEM.md)

---

## 1. The decision — and what it does NOT mean

**Decision:** split the platform into separate repos by **technology/lifecycle**, with a versioned **contracts repo** as the single shared seam.

**Yes:** one repo per technology stack + one contracts repo + one infra repo. Independent pipelines, independent versioning, independent deploy cadence, independent blast radius.

**No — critically:** polyrepo does **not** mean microservices. The Go backend stays a **modular monolith in one repo**. We split by technology and lifecycle, never by domain (auth vs catalog vs playlist). Domain-splitting the Go service is deferred indefinitely — trigger: team size + measured scaling pain, and even then, reconsider.

> Why: a monolith-with-modules is far easier to reason about, deploy, and test than a distributed system. Separate repos give *operational* independence without paying the *distributed-systems* tax early.

---

## 2. The repo map

| Repo | Stack | Deploy cadence | Blast radius | Owns |
|---|---|---|---|---|
| **muse-contracts** | OpenAPI spec + generated clients (TS/Go/Py/Rust) | every change is a release (tagged) | the seam — everything depends on it | API surface, media URL contract, event schemas |
| **muse-web** | Nuxt 3 SSR | daily | user-facing, high traffic | SEO pages, app shell, player, admin |
| **muse-api** | Go modular monolith | weekly | system of record | auth, catalog, playlists, library, queue, social, admin, outbox |
| **muse-streamer** | Rust media edge | rare | media delivery | range/HLS/DASH serving over S3/MinIO, behind CDN |
| **muse-ai** | Python FastAPI + Celery | when pipeline/model changes | AI features degrade gracefully | lyrics/Whisper, embeddings, similar tracks, moderation |
| **muse-infra** | Terraform + compose + runbooks | manual, gated | shared infra | Postgres, Redis, MinIO, Kafka/ClickHouse (later), deploy docs |

**Layout:**
```
muse-contracts/
├── openapi/            # the source of truth (v1.0.0 → vN)
├── proto/              # (phase 4+ for internal RPC / Kafka schemas)
├── gen/
│   ├── typescript/     # @muse/contracts  (npm)
│   ├── go/             # go.muse.dev/contracts
│   ├── python/         # muse-contracts   (PyPI)
│   └── rust/           # muse-contracts   (crates.io)
└── scripts/release.sh  # tags + publishes all 4 packages
```

---

## 3. The contracts repo — what makes polyrepo sustainable

The most important file in the platform is `muse-contracts/openapi/v1.yaml` — the source of truth for what the API returns, what the frontend renders, what the streamer streams, and what the AI consumes.

**How it flows:**
```
muse-contracts v1.4.0 (released, tagged)
   ├── openapi/v1.yaml
   ├── @muse/contracts@1.4.0        → used by muse-web
   ├── go.muse.dev/contracts@v1.4.0 → used by muse-api
   ├── muse-contracts==1.4.0        → used by muse-ai
   └── muse-contracts@1.4.0 (crate) → used by muse-streamer
```

**Versioning rules:**
- **Additive** change (new optional field, new endpoint) → minor bump. No coordination; each repo upgrades on its own cadence.
- **Breaking** change (renamed/removed field, changed shape) → major bump. Coordinated window: consumers upgrade first, the contract goes live behind a feature flag, the old version is served during migration.

**Enforcement — non-negotiable:**
- Every app repo's CI has a **contract-check** job: it validates the app against its *pinned* contract version.
- Tooling: `openapi-generator` / `protoc` / `buf` generate clients; `renovate`/`dependabot` keep pinned versions current.
- A contract change is a *release*, not an edit. Tag it. No untagged contract changes.

> Without this, separate repos degrade into "it worked on my machine." With it, four languages can't drift — the spec is the referee.

### Change flow — add `is_explicit` to a track

1. **muse-contracts:** add `isExplicit?: boolean` → tag `v1.5.0` (additive) → publish 4 packages.
2. **muse-api:** add column, expose field, bump pin to `v1.5.0`. CI contract-check green.
3. **muse-web:** bump `@muse/contracts` to `1.5.0`, use the new TS type, deploy.
4. **muse-streamer / muse-ai:** no change needed; they upgrade later or never.

A **breaking** change: bump to `v2.0.0` → all consumers upgrade in a planned window → old spec stays served until all are on v2.

---

## 4. Runtime topology (Iranian DC)

```
User ── CDN (ArvanCloud) ──┬── muse-web (Nuxt SSR, nginx)
                           └── muse-streamer (Rust media edge) ── MinIO / object storage
        muse-api (Go, internal VMs) ── Postgres, Redis
        muse-ai (Python, internal, HMAC-gated)
        observability: Prometheus + Grafana + Loki (self-hosted on the VM — no vendor dependency)
```

**Phase 1 runs without the dotted pieces:** nginx serves audio (range requests), events flow through the outbox → `events_log`, no Kafka, no ClickHouse, no Rust, no Python beyond an optional stub.

**Storage strategy:** `local`/`s3` drivers behind the seam (reuse the current repo's storage abstraction). Iranian DC: **MinIO + CDN first**; provider object storage (ArvanCloud) only when measured. It's a config swap, not a rewrite.

---

## 5. Data-pipeline evolution — the core architectural bet

The single most important pattern, from day one, is the **event outbox** (copied from the current repo's `000044`):

```
Phase 1:  API writes domain + outbox atomically → relay logs "would publish" → events_log
Phase 2:  relay becomes real → publishes to events_log (still no Kafka)
Phase 4:  relay sink swaps → publishes to Kafka topics → ClickHouse ingest, workers, Python trainers
```

**Why this wins:** at every phase the *boundary* is real (events exist, atomically, idempotently) while the *technology* is deferred. Kafka becomes a sink swap, not a rewrite. Architecture is about boundaries, not about which services you're running.

**The outbox is a service too** — it needs queue-depth metrics and a "silent outbox" alert. Land it (see [ENGINEERING_OPERATING_SYSTEM.md](ENGINEERING_OPERATING_SYSTEM.md)) before adding Kafka; a Kafka pipeline on a silent outbox is a distributed way to lose data.

---

## 6. Reliability via the split

| Concern | How the polyrepo handles it |
|---|---|
| **Independent deploy** | each repo has its own pipeline + rollback; a frontend hotfix deploys without touching the API |
| **Isolated blast radius** | AI outage degrades lyrics, not playback; streamer outage is absorbed by CDN cache; API outage doesn't block serving cached media |
| **Independent scaling** | web scales with traffic, API with writes, streamer with media concurrency, AI with queue depth |
| **Ownership** | each repo owns its README, CODEOWNERS, issue templates, runbook |
| **Observability across the split** | one **correlation header** (`X-Request-ID`/`traceparent`) flows web → API → outbox → worker → AI; one incident across four repos is one header search |

> SLO targets, error budgets, and the observability triad are defined in [ENGINEERING_OPERATING_SYSTEM.md](ENGINEERING_OPERATING_SYSTEM.md#3-slos--slas).

---

## 7. Deferred-infra triggers (canonical)

The CTO's answer to "why aren't we building X?" — decided in advance so the decision is never made by momentum or anxiety. **This is the canonical table.** If a trigger doesn't fire, we don't build it — that's a feature, not a failure.

| Deferred item | Trigger (measured, user-visible or metric-justified) |
|---|---|
| **Kafka** | outbox `events_log` grows + a second consumer needs events (analytics, notifications) |
| **ClickHouse** | listen events make Postgres analytics slow + a user-facing chart/leaderboard needs it |
| **Rust streamer** | nginx range-serving is the bottleneck (huge catalog, HLS/DASH, high concurrency) |
| **OpenSearch / pgvector** | Postgres FTS measurably underperforms at real search load |
| **ML recs v2** | enough listening history that collaborative filtering beats popularity |
| **Microservices (domain split)** | team size + measured scaling pain — and even then, reconsider |
| **Payments / revenue** | engaged listener base + a chargeable signal |

---

## 8. Migration from the current monolith repo

The current repo is already mostly split internally — `frontend/`, `moja-ml-service/`, and the Go code are separate trees. The work is extraction, not re-architecture.

1. **muse-contracts** — extract the v1 OpenAPI from the current routes (`internal/app/routes.go` shows the full surface) + seed generated packages. The foundation; do it first.
2. **muse-web** — move `frontend/` as-is (it is self-contained). Own CI (lint/typecheck/build), `.env.example`, README.
3. **muse-api** — move the Go code (`cmd/`, `internal/`, `migrations/`). **Drop the committed binaries** (`api`, `worker` — ~150 MB total) and fix `.gitignore`. Keep the modular-monolith structure — it's already right.
4. **muse-ai** — move `moja-ml-service/`.
5. **muse-streamer** — create a fresh Rust skeleton implementing the media URL contract. Until it exists, nginx serves `uploads/` (the seam is the contract, not the service).
6. **muse-infra** — move `deployments/` + compose + runbooks.
7. Each repo must build green **standalone** before any integration wiring.

**Order matters:** contracts → api (system of record) → web → ai → streamer → infra.

---

## 9. Rules that keep it sustainable

1. **Never share code across app repos** except through `muse-contracts`.
2. **Never share `.env`** — each repo has its own `.env.example` and secrets story.
3. **Never commit binaries or artifacts** — the current repo has `api` (86 MB) and `worker` (65 MB) committed; that never happens again.
4. **A contract change is a release** — tagged, versioned, CI-gated.
5. **Each repo README has exactly three sections:** Run · Test · Deploy.
6. **The Go monolith stays a monolith** — domain-splitting is a deferred decision with its own trigger (above).
7. **Every repo observes the same correlation header** — no incident is ever "in one repo" because it can't be found.

---

## 10. What ships in Phase 1 of the polyrepo

- `muse-contracts` v1.0.0 (auth + tracks + search + upload surface) with the 4 packages publishing.
- `muse-api` Go monolith, green standalone, CI gate on vet/test/build/contract-check.
- `muse-web` Nuxt SSR, green standalone, CI gate on lint/typecheck/build/contract-check.
- `muse-infra` compose (Postgres, Redis, MinIO) + the one-VM deploy runbook.
- `muse-streamer` and `muse-ai` exist as **stubs** (a `TODO` contract implementation) — the seams are real, the services wait for their phase (see [ROADMAP.md](ROADMAP.md#7-phase-4--scale-on-evidence)).
