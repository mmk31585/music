# Muse — Master Plan

> The definitive plan for building and operating a Persian music streaming platform.
> Three lenses, one spine. **CEO** (why/what) → **CTO** (architecture/how) → **Senior Dev** (execution/when).
> Detail lives in: [ROADMAP.md](ROADMAP.md) · [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md) · [ENGINEERING_OPERATING_SYSTEM.md](ENGINEERING_OPERATING_SYSTEM.md).
> This document is the map; those are the territory.

---

## 0. Executive summary (read this in 2 minutes)

**What we're building:** the place where Persian music fans find, play, and share music — fast, reliable, and in Iran.

**How we win:** not with the tech stack (a competitor copies Go-vs-Rust in a quarter). We win with **catalog access, community, and trust** — the things that are hard to copy. Engineering's job is to make the product fast and reliable so those advantages compound.

**The discipline:** start small, ship to a real server early, and let *measured evidence* decide when to add Kafka, ClickHouse, Rust, or AI. Every heavy technology has a written trigger. Nothing arrives by momentum.

**The stack (phased):** Nuxt 3 SSR · Go modular monolith · Postgres + Redis · MinIO + CDN · (later) Rust streamer, Python AI, Kafka + ClickHouse. Six repos, one contracts repo as the seam, Iranian-DC deployable.

**The spine:** *you own what you ship, end to end.* Observability, SLOs, blameless post-mortems, launch-vs-land, managed tech debt. Product has an operator, not just a builder.

**Timeline:** Phase 0 (1 wk) → Phase 1 MVP (2–4 wk) → Phase 2 player/library (2–3 wk) → Phase 3 AI/recs (2–3 wk) → Phase 4 scale-on-evidence → Phase 5 ecosystem. **Shippable product in ~8 weeks.**

---

# LAYER A — CEO (why & what)

## A1. The one sentence and the north-star metric

**One sentence:** *A place where someone can find and play Persian music, and it loads fast enough that they stay.*

**North-star metric:** **weekly active listeners** — people who come back and play, not signups. Every decision is judged: does it move weekly active listeners, or is it decoration?

Supporting metrics that explain the north star (the "why is it moving" dashboard):

| Metric | Answers | Dashboard |
|---|---|---|
| Click-to-play < 1s p95 | is it fast enough to stay? | latency |
| Search success rate | can they find what they want? | % of searches that lead to a play |
| Play completion rate | do they stay to the end? | % of plays reaching 80%+ |
| Daily active listeners | is it habit-forming? | DAU/WAU |

## A2. Market reality (be honest, it shapes everything)

- **Iranian data center:** hosting, DNS, CDN (ArvanCloud), and payment gateways are *local* decisions. Global SaaS billing is not available — payments need a local-gateway abstraction (already designed in the old repo's `internal/platform/payment`).
- **Storage:** no blind trust in provider "S3." Ship **MinIO + CDN** first; provider object storage only when measured. This is the storage strategy in [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md).
- **The moat** (what a competitor can't copy quickly): catalog rights + artist relationships, community/social graph, and the brand users trust to stream reliably. The tech stack is *not* the moat.

## A3. Revenue model — when it's real

Phase 1–3 = **no monetization.** Focus is product-market fit. Revenue conversations start only when there are engaged listeners (the trigger: a meaningful weekly active listener count and a "we could charge for this" signal). Candidate models, in order of fit for this market:
1. **Subscriptions** (ad-free, offline, higher quality) — needs the local gateway.
2. **Creator tips / subscriptions** — monetizes the artist side of the moat.
3. **Ads** — later, only with scale.

> The CEO posture here: **survival before optimization.** Optimize only from a position of having users. The roadmap's phase order *is* the survival strategy.

## A4. What we are NOT building (written down so we don't build it)

- Social rooms, clubs, gamification, subscriptions → **Phase 5**, only with an audience.
- Kafka, ClickHouse, Rust streamer, OpenSearch, embeddings → **Phase 4**, only on measured triggers.
- ML-based recommendations → only when there's enough listening data to beat a popularity sort.

Unwritten ambitions become accidentally-built features. This list is the guardrail.

---

# LAYER B — CTO (how)

## B1. The polyrepo map

Six repos, split by **technology/lifecycle, never by domain.**

```
muse-contracts   OpenAPI + generated clients (TS/Go/Py/Rust)   — the seam, released on every change
muse-web         Nuxt 3 SSR                                    — SEO + app shell + player + admin
muse-api         Go modular monolith                           — system of record (auth/catalog/playlist/library/queue/social/admin/outbox)
muse-streamer    Rust media edge                               — range/HLS/DASH over S3/MinIO, behind CDN
muse-ai          Python FastAPI + Celery                       — Whisper lyrics, embeddings, similar tracks
muse-infra       Terraform + compose + runbooks                — Postgres/Redis/MinIO/(Kafka/ClickHouse later)
```

**Critical decision:** polyrepo ≠ microservices. The Go backend stays a **modular monolith**. We split by technology for operational independence; we do *not* split the API by domain (that's deferred indefinitely, trigger: team size + measured scaling pain).

## B2. The runtime topology (Iranian DC)

```
User ── CDN (ArvanCloud) ──┬── muse-web (Nuxt SSR, nginx)
                           └── muse-streamer (Rust media edge) ── MinIO / object storage
        muse-api (Go, internal VMs) ── Postgres, Redis
        muse-ai (Python, internal, HMAC-gated)
        observability: Prometheus + Grafana + Loki (self-hosted on the VM — no vendor dependency)
```

**Phase 1 runs without the dotted pieces:** nginx serves audio (range requests), events flow through the outbox → `events_log` table, no Kafka, no ClickHouse, no Rust, no Python service beyond an optional stub.

## B3. The data pipeline evolution — the core architectural bet

The single most important pattern, from day one, is the **event outbox** (copied from the old repo's `000044`):

```
Phase 1:  API writes domain + outbox atomically → relay logs "would publish" → events_log
Phase 2:  relay becomes real → publishes to events_log (still no Kafka)
Phase 4:  relay sink swaps → publishes to Kafka topics → ClickHouse ingest, workers, Python trainers
```

**Why this wins:** at every phase, the *boundary* is real (events exist, idempotently, atomically) while the *technology* is deferred. Kafka becomes a sink swap, not a rewrite. This is "architecture is about boundaries, not about which services you're running" made concrete.

## B4. What's deferred, and the trigger for each

This table is the CTO's answer to "why aren't you building X?" — decided in advance so the decision is never made by momentum or anxiety.

| Deferred item | Trigger (measured, user-visible or metric-justified) |
|---|---|
| **Kafka** | outbox `events_log` grows + a second consumer needs events (analytics, notifications) |
| **ClickHouse** | listen events make Postgres analytics slow + a user-facing chart/leaderboard needs it |
| **Rust streamer** | nginx range-serving is the bottleneck (huge catalog, HLS/DASH, high concurrency) |
| **OpenSearch / pgvector** | Postgres FTS measurably underperforms at real search load |
| **ML recs v2** | enough listening history that collaborative filtering beats popularity |
| **Microservices (domain split)** | team size + measured scaling pain — and even then, reconsider |
| **Payments / revenue** | engaged listener base + a chargeable signal |

**If a trigger doesn't fire, we don't build it — that's a feature, not a failure.** This is survival-before-optimization applied to architecture.

## B5. Reliability targets (SLOs, Phase 1)

| Signal | SLO | Error budget (30d) |
|---|---|---|
| API availability | 99.9% | ~43 min |
| Search p95 | < 200 ms | — |
| Click-to-play p95 | < 1 s | — |
| SSR render p50 | < 500 ms | — |
| API error rate | < 0.5% | — |
| Playback success | > 99% | — |

Error-budget rule: **exhausted budget ⇒ no risky releases that week.** The budget decides, not vibes, when the codebase is too unstable to touch.

---

# LAYER C — SENIOR DEV (execution)

## C1. The phased roadmap (detail: [ROADMAP.md](ROADMAP.md))

| Phase | Duration | Ships (user-visible) | New infra | Tests | CI |
|---|---|---|---|---|---|
| **0 Foundations** | 1 wk | tooling + CI + dev env | Postgres, Redis, MinIO | pipeline itself is the test | api/web/contracts jobs |
| **1 Vertical slice MVP** | 2–4 wk | search → play, end-to-end, on a VM | nginx, outbox (relay off) | unit (auth/search/storage), integration (real Postgres), Playwright golden path | + integration, e2e-smoke, contract-check |
| **2 Player & Library** | 2–3 wk | player, queue, playlists, liked, history | partitioned history, outbox relay real | unit (player engine hard), integration (playlist/history/outbox), e2e | + vitest coverage gate on player/ |
| **3 AI & Recs** | 2–3 wk | lyrics, similar tracks, daily mix | FastAPI + Celery | pytest (normalizer/LRC), contract Go↔Py, integration (upload→lyrics) | + ai job, contract job |
| **4 Scale on evidence** | trigger-driven | charts, HLS, better recs | Kafka, ClickHouse, Rust | consumer idempotency, ClickHouse rollups, Rust range | per-service jobs |
| **5 Ecosystem** | with users | social, gamification, subs | — (feature flags) | per-feature unit+integration+e2e | as added |

**Definition of Done per phase** (all must pass): golden-path e2e · unit+integration green in CI · type checks · `make deploy` runbook executed · observability wired (log+metric+trace for the headline action) · exit criteria *measured* · ADR if a hard decision was made.

## C2. What actually happens Week 1 (the concrete start)

```
Day 1  monorepo; .gitignore; commitlint; Makefile; infra/compose (postgres, redis, minio); make infra-up works.
Day 2  muse-api Go module; health module; /healthz + /metrics; CI api job green.
Day 3  migrations + outbox table (000044 pattern) + genre seed; migrate-up wired.
Day 4  muse-contracts OpenAPI skeleton (auth + tracks + search); drift check in CI.
Day 5  muse-web Nuxt app; dev proxy; one SSR page (/artist/:id) from the API; CI web job green.
Day 6  ADR-001 (stack+phases+deferred triggers), ADR-002 (monolith); README.
Day 7  Demo: fresh clone → make infra-up → make run → SSR page renders. CI green. VM notes written.
```

## C3. Engineering discipline (operational now)

- **Testing pyramid:** unit-test domain logic hard (player engine, search ranking, Persian normalization, queue); integration-test repos against *real* Postgres (never mock the DB); one Playwright e2e per golden path. Coverage as a floor, not a target.
- **Observability:** structured logs + `X-Request-ID` correlation header across web→API→outbox→worker→AI; `/metrics` Prometheus; Grafana + Loki self-hosted; **one business dashboard** (MAU, plays/day, search success, top tracks). Empty dashboard ≠ no product decision to define it.
- **Security baseline:** short-lived access tokens + rotating refresh in httpOnly cookie; secrets never in git; upload allowlists + size limits; rate limits on auth/search; parameterized SQL only.
- **Feature flags** for anything risky or half-done. Shipped-off beats not-shipped.
- **Code review:** small PRs, one concern, CI green, reviewer is future-you. When a second engineer arrives, review becomes mandatory.
- **Tech-debt ledger** (`docs/tech-debt.md`): every PR that takes debt logs one line; pay interest in the same PR; pay principal on triggers. Unlogged debt is invisible debt.

## C4. The operations playbook (when it breaks)

1. **Detect → contain → mitigate → investigate → post-mortem.** Containment before investigation, always.
2. **Blameless post-mortem within 48h** ([template](ENGINEERING_OPERATING_SYSTEM.md#4)): ask *"what made this mistake possible and undetected?"* not *"who made it?"* Action items tracked to done. **No incident happens twice** — if it recurs, the process failed, not the person.
3. **Severity:** P0 (playback/search down for all → wake up now) · P1 (degraded for some → same day) · P2 (cosmetic → next release).

---

# GOVERNANCE — who decides what

Solo today, defined now so the norms pre-exist the team:

| Decision type | Owner | Example |
|---|---|---|
| What we build (scope) | **CEO lens** | ship search before social rooms |
| How we build it (architecture) | **CTO lens** | monolith now, Kafka on trigger |
| When it's done (quality) | **Dev lens** | Definition of Done, SLOs, post-mortems |
| Deferred-infra triggers | **The trigger table** | no one gets to add ClickHouse by vibes |

Rules of the game:
- **Any of the three lenses can veto** — if a decision breaks the north star, the architecture spine, or the quality floor, it stops.
- **Direct communication:** the shortest message that carries the full intent. No status theater.
- **One goal per week**, one north-star metric. No ritual Scrum.
- **The engineering ladder** (IC + management tracks, lateral moves not demotions) is written in [ENGINEERING_OPERATING_SYSTEM.md](ENGINEERING_OPERATING_SYSTEM.md) and activates when the team forms.

---

# SUCCESS & FAILURE CRITERIA

**We are succeeding when (in order):**
1. Someone finds and plays music on our server, end to end (Phase 1 done).
2. Weekly active listeners grow, and the "why" (search success, completion rate) is visible on the dashboard (Phase 2–3).
3. A real user would be sad if we shut down (product-market fit signal).

**We are failing if:**
- We spend 3+ months without anything deployed (architecture theater).
- We add Kafka/ClickHouse/Rust *before* a trigger fires (optimization before survival).
- The business dashboard is empty and we keep building anyway (no operator).
- We ship features nobody plays (north star ignored).

**Kill criteria (the hard one):** if, after Phase 3 with real deployment and real usage attempts, there is no engagement signal at all — we stop, learn, and pivot the product rather than the stack. The stack was never the risk; the product is. That's the honest CEO read.

---

## Index of the supporting documents

- [ROADMAP.md](ROADMAP.md) — phased execution, tests, CI, exit criteria, week-1 sprint.
- [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md) — polyrepo map, contracts flow, migration from the old monolith.
- [ENGINEERING_OPERATING_SYSTEM.md](ENGINEERING_OPERATING_SYSTEM.md) — SRE practices, culture/org defaults, strategy layer, post-mortem template.
- (soon) `docs/tech-debt.md` — the debt ledger.
- (soon) `docs/runbooks/POSTMORTEM.md` — the incident template.
