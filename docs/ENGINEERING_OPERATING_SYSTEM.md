# Muse — Engineering Operating System

> **Territory:** SRE/ops practices · culture/org defaults · strategy/market layer · post-mortems.
> Applies to every repo in the polyrepo (muse-contracts/web/api/streamer/ai/infra).
> This is the *how-we-operate* layer.
> - Map (why/what/architecture): [MASTER_PLAN.md](MASTER_PLAN.md)
> - Execution (phases/tests/CI): [ROADMAP.md](ROADMAP.md)
> - Architecture (repos/contracts/migration): [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md)
>
> **Spine:** *you own what you ship, end to end — write, deploy, monitor, on-call, learn.*
> Everything below is that sentence made concrete.

---

## 0. The three layers — and when each is real

| Layer | Applies | What it is |
|---|---|---|
| **1. Ops / SRE** | **NOW (solo)** | Observability, SLOs, incident response, launch-vs-land, tech-debt pipeline. Buildable this week. |
| **2. Culture / Org** | When a team forms | Ladder, goal-based structure, PM role, code review culture, incentives. Documented now, enforced later. |
| **3. Strategy / Market** | When traction + funding | Survival vs optimization, moat, capital, connecting tech/product/business. |

A senior knows that applying Layer 3 to a pre-traction product is self-deception, and applying Layer 2 to a solo dev is ceremony. Layer 1 is the whole job right now.

---

# LAYER 1 — OPS / SRE (operational now)

## 1. End-to-end ownership

Every module you write, you own **until it is deleted**. Owning means:

1. **You wrote it** — it follows your conventions, has tests, has a README line.
2. **You deployed it** — the deploy runbook runs from your machine and is a script, not a memory.
3. **You monitor it** — it emits logs, metrics, and traces. If it can't be observed, it doesn't ship.
4. **You're on call for it** — a P0 at 3am is yours. Solo or not, define severity now:
   - **P0** — playback/search down for all users → wake up, fix now.
   - **P1** — degraded for some (slow search, lyrics missing) → fix same day.
   - **P2** — cosmetic/edge → next release.
5. **You learn from it** — every incident produces a post-mortem (below), and the action items are *tracked to done*, not filed.

Rule: **a module with no on-call owner is an orphan; orphans get deleted.**

## 2. Observability — the triad

Every service emits all three. Nothing ships with fewer.

| Signal | Tool | Must answer | Muse specifics |
|---|---|---|---|
| **Logs** | Structured JSON, one library (zap in Go, winston/pino in web, logging in Python) | *What happened?* | Every request logs a request-id; every event logs the correlation header. |
| **Metrics** | Prometheus + Grafana (self-hosted on the VM — no vendor dependency in the Iranian DC) | *How many / how fast / how healthy?* | Golden signals per service: latency, traffic, errors, saturation. Plus business: plays, searches, uploads. |
| **Traces** | One **correlation header** (`X-Request-ID` / `traceparent`) flows web → API → outbox → worker → AI | *Where did the time go?* | One incident across four repos = one header search. |

**ELK alternative:** Elasticsearch+Logstash+Kibana is heavy. Start with **Loki + Grafana** (same Grafana, log aggregation without running a cluster). Upgrade to ELK only when log volume makes Loki insufficient — same trigger discipline as everything else.

**The golden signals (define per service before writing it):**

```
muse-web:      request rate, p50/p95/p99 latency, error rate, SSR render time
muse-api:      request rate, p95 latency, error rate, DB query time, cache hit ratio
muse-streamer: throughput MB/s, range-request count, cache hit ratio, origin latency
muse-ai:       queue depth, job duration, job failure rate, model inference time
muse-infra:    CPU/mem/disk on Postgres/Redis/MinIO, connection pool saturation
```

**One business dashboard** (Grafana): MAU, daily active listeners, plays/day, search success rate, uploads/day, top tracks. If the dashboard is empty, you don't have a product — you have a project.

## 3. SLOs / SLAs — with real numbers

**SLA** = customer-facing commitment (you'd pay a penalty). **SLO** = your internal target. Most things are SLOs. Define them per service, and define the **error budget** (how much "bad" you're allowed before you stop shipping risky changes).

**Phase-1 SLOs (concrete, achievable on one VM):**

| Signal | SLO target | Error budget (per 30d) |
|---|---|---|
| API availability | 99.9% | ~43 min down allowed |
| Search p95 latency | < 200 ms | — |
| Click-to-play (to audio first byte) | < 1 s p95 | — |
| SSR page render (artist/album/track) | < 500 ms p50 | — |
| API error rate | < 0.5% of requests | — |
| Playback success | > 99% of play attempts produce audio | — |

**The error-budget policy:** if you exhaust the budget, the rule is **no risky releases** until you earn it back — you don't ship that week, you stabilize. This is the *entropy management* part of SRE: the budget decides, not vibes, when the codebase is too unstable to touch.

**The outbox is a service too:** give the relay a metric (queue depth, processed/sec, failed count) and an alert. A silent outbox is a silent data loss.

## 4. Incident response + blameless post-mortem

**Process (Solo or team — same five steps):**
1. **Detect** — alert fires (or a user tells you). Log the timestamp.
2. **Contain** — stop the bleeding first: rollback, disable feature flag, fail-over. Never investigate before containing.
3. **Mitigate** — restore service; announce status.
4. **Investigate** — root cause analysis *after* the fire is out.
5. **Post-mortem** — within 48h, blameless, action items tracked.

**Blameless post-mortem template** (`docs/runbooks/POSTMORTEM.md` — create it): a post-mortem is not about *who*, it's about *what in the system allowed this*. If the answer is "I made a mistake," that's a wrong post-mortem — the right question is *"what made the mistake possible and undetected?"* (missing alert, no test, no canary, no runbook).

```
## Post-mortem — <title>
Date / Duration / Severity (P0-P2) / Reporter
### Summary (2 sentences, what users experienced)
### Timeline (UTC)
   detect → contain → mitigate → investigate
### Root cause (system-level, not person-level)
### What allowed it to happen undetected?   ← the real question
### Impact (users, availability, error budget consumed)
### Corrective actions (each: owner, deadline, tracked to done)
   [ ] Action A — CI gate added
   [ ] Action B — alert threshold fixed
   [ ] Action C — runbook written
### Lessons (one line: what behavior change prevents this class)
```

**The promise:** *no incident happens twice.* Every post-mortem ends in action items that are tracked to done, not a nice document that gets filed. If the same incident recurs, that's a failure of the process, not of the person.

## 5. Launch vs. Land

**Launch** = ship a feature. **Land** = the feature is stable, monitored, documented, and someone's responsible. Your trap (as a builder) is falling in love with Launch and never Landing.

**The "Land" checklist — a feature isn't done until ALL of these:**

- [ ] Golden-path e2e passes (Playwright)
- [ ] Logged, metered (Prometheus), traced (correlation header)
- [ ] Alert exists for its failure mode
- [ ] Error budget assigned / updated
- [ ] Runbook exists (how to debug/rollback it)
- [ ] Backed up (if it writes data)
- [ ] Feature-flagged if risky
- [ ] Written down (README line or ADR if a hard decision)

**Phase 1's own Launch/Land example:** the outbox table *launches* in Phase 1 (created, wired to the relay). It *lands* when: the relay has queue-depth metrics, a "silent outbox" alert, and a runbook. Land it before adding Kafka — a Kafka pipeline on top of a silent outbox is a distributed way to lose data.

## 6. Tech debt — the financial tool

**Tech debt is a loan you take against the future to move fast now.** Taking it is correct. The discipline is: (a) know what you borrowed, (b) pay the interest, (c) pay down principal in a managed pipeline.

**Your current repo's debt ledger (real, from the audit):**

| Debt | Type | Interest you're paying | Paydown trigger |
|---|---|---|---|
| 86 MB `api` + 65 MB `worker` binaries committed to git | hygiene | repo is unbounded; clones slow; history bloated | NOW — delete + gitignore in the polyrepo split |
| DTOs double as domain models (M-002) | design | Go↔contracts drift later makes the Rust/TS split painful | Phase 1, as you write new modules |
| `recommender_v2` duplicates `recommendation` worker | duplication | two workers to maintain, ambiguous ownership | NOW — pick one, delete the other |
| OpenSearch cluster before search load | premature infra | a whole cluster to operate for a catalog Postgres FTS handles | defer (trigger: Postgres FTS measurably slow) |
| `internal/common/middleware` + `internal/middleware` coexist (M-003) | hygiene | two homes for the same concept | first refactor that touches middleware |
| Frontend organized by type not domain (M-001) | design | hard to navigate at scale | Phase 2+ when player grows |

**The pipeline:**
1. **Log every debt you take** — a one-line note in the PR or a `docs/tech-debt.md` ledger. Unlogged debt is invisible debt.
2. **Pay interest continuously** — a small refactor *in the same PR* where you touch the code. Never leave a spot worse than you found it.
3. **Pay principal on triggers** — rewrites happen on a schedule or a trigger (the canonical deferred-infra trigger table lives in [ARCHITECTURE_REPOS.md](ARCHITECTURE_REPOS.md#7-deferred-infra-triggers-canonical)). A managed rewrite is a business decision; a reactive one is a fire.
4. **Never refactor hot code in a crisis** — stabilize first (error budget), then refactor. This is the survival-vs-optimization discipline below, applied to code.

---

# LAYER 2 — CULTURE / ORG (documented now, enforced when the team forms)

These are written here as the *defaults* so that when a second person arrives, the norms pre-exist the team. None of this is buildable solo — do not burn energy implementing it before people exist.

## 7. Engineering ladder (no dead-ends)

Two parallel tracks, both senior, both paid:

```
Individual Contributor (IC)        Management
Staff / Principal                  Head of Eng
Senior IC (IC3+)                   Engineering Manager
Engineer (IC1-2)                   Tech Lead (hybrid)
```

Rules:
- **Moving tracks is a lateral move, never a demotion.**
- Promotion criteria are *behavioral and documented*, not tenure: "owns a system end-to-end," "raises the bar on code review," "teaches others," "unblocks a launch."
- A manager who can't code is a bureaucrat; an IC with no path up is a resignation risk. The ladder is the answer to both.

## 8. Goal-based structure, not ritual Scrum

- Replace ceremony (standup-grooming-retro as religion) with **one goal per week** and **one north-star metric**.
- No task is assigned *how* — engineers own their approach and report *outcome*, not activity.
- If a standup takes 5 minutes, it's useful. If it takes 30, it's a status theater — cancel it and use async updates (the project's task list, PRs, a shared notes doc).
- **Direct communication:** the shortest message that carries the full intent. No "just FYI" emails; no status theater; no passive-aggression. Praise in public, correct in private, feedback in the moment.

## 9. Product ↔ Engineering friction

- **PM's job is value-to-customer, not throughput.** The PM owns *what* and *why*; engineers own *how* and *when*.
- A PM who becomes a queue (specs in, features out, no feedback loop) is a cost, not a role.
- Default: engineers talk to **real users** directly. No proxy for the product. The PM's value is *focus* — protecting engineers from the 20 features that won't matter — not *toll-boothing* the ones that will.
- **Incentive alignment:** individual goals, the team's metric, and the company's metric are the same number. If someone's promotion depends on shipping features and the company's survival depends on stability, they will ship unstable things. Fix the incentive, not the person.

## 10. Code review culture

- **Purpose:** catch bugs early + share ownership + teach. Not gatekeeping, not approval theater.
- Small PRs, one concern each, CI green, reviewer is a colleague (or future-you).
- Review rules: correctness first, then simplicity, then style. "Can I understand this in 5 minutes?" is the bar.
- **Never ship unreviewed code to production** — the moment you have two engineers, review is non-negotiable, because it is the cheapest post-mortem prevention there is.

---

# LAYER 3 — STRATEGY / MARKET (when there is traction)

## 11. Survival vs. Optimization

Every decision gets the test: **does this improve our odds of still existing in 12 months, or does it optimize for a world we're already in?**

- **Survival moves** (aggressive, financial, sometimes ugly): get paying users, cut unprofitable scope, raise/borrow at the right moment, price against the actual market.
- **Optimization moves** (polish, scale, refactor): only affordable *from a position of survival*.

The trap: solo builders optimize (architecture, refactors, "the right stack") while the company is still in survival (no users, no revenue). **The roadmap's phase order *is* the survival strategy**: Phase 1 ships to market; Phase 4 (Kafka/ClickHouse/Rust) is optimization that only fires on measured triggers. If a phase isn't user-visible or metric-justified, it's optimization before survival — postpone it.

## 12. Moat & market dynamics (regional reality)

A **moat** is what stops a competitor from taking your users. For a Persian music platform, be honest about where moats actually live in this market:

- **Not** the tech stack (Go vs Rust, ClickHouse vs Postgres — a competitor copies it in a quarter).
- **Content/catalog access** — rights and relationships with Persian artists/labels. Hard to copy, legally layered.
- **Community & social graph** — rooms, clubs, creators, taste data. Hard to copy, compound.
- **Distribution/trust** — the brand users believe streams reliably in Iran.

So the strategy implication: **spend engineering on catalog ingestion, artist relationships, and community — not on OLAP elegance.** The big-data stack (ClickHouse, Pinot, Spark, Flink, MapReduce — all the concepts in your list) is what you deploy when you have the data and the competition; it is not what creates the moat.

**Capital/competition layer:** if you're competing regionally and/or raising money, the market dynamics (anti-competitive behavior by incumbents, legal layers, fundraising timing) are a *business* function. That is a different conversation from engineering — and it should start only when traction says you're in the game.

## 13. The executive connection

The senior in the room connects three lines at once, and only acts where they meet:

```
technology ─── what's buildable
product     ─── what's useful
business    ─── what's viable
```

Every decision is the intersection. A feature that's buildable and useful but not viable is a hobby. A feature that's viable and buildable but not useful is waste. **The roadmap's phase plan is exactly this intersection, applied ahead of time.**

---

## 14. The "when" matrix (keep this on your wall)

| Practice | Solo (now) | Small team | Scaling |
|---|---|---|---|
| Observability triad | ✅ build | ✅ enforce | ✅ scale |
| SLOs + error budget | ✅ define | ✅ enforce | ✅ automate |
| Blameless post-mortems | ✅ (on yourself) | ✅ enforce | ✅ process |
| Launch vs Land checklist | ✅ enforce | ✅ enforce | ✅ enforce |
| Tech-debt ledger + pipeline | ✅ start | ✅ enforce | ✅ tool it |
| E2E ownership | ✅ natural | ✅ formalize on-call | ✅ rotations |
| Code review | solo = self-review | ✅ mandatory | ✅ automated gates |
| Engineering ladder | — | ✅ design | ✅ evolve |
| Goal-based structure | ✅ (one goal/week) | ✅ replace Scrum | ✅ scale |
| PM role / incentives | — | ✅ when PM exists | ✅ critical |
| Moat / market / capital | think, don't spend | think | ✅ full focus |

---

## 15. Where this lands in execution

This document defines *how we operate*; [ROADMAP.md](ROADMAP.md) is *what we build and when*. Two artifacts from Layer 1 are created as part of Phase 0, before the first line of product code:

1. **`docs/runbooks/POSTMORTEM.md`** — the template above, ready for the first incident (there will be one).
2. **`docs/tech-debt.md`** — the debt ledger; every PR that takes debt logs one line.

Observability scaffolding (structured logs + `X-Request-ID`, `/metrics`, Prometheus + Grafana + Loki, the *business* dashboard with panels defined even if empty) is wired in Phase 0/1 per the ROADMAP. The moment the first real incident happens, run the first blameless post-mortem on yourself — the habit is the point.

> The difference between a project and a product is that a product *has an operator*. These practices — observability, SLOs, post-mortems, land-over-launch, managed debt — are what turn you from a builder into an operator. Everything in this document is the sentence *"you own what you ship"* written large.
