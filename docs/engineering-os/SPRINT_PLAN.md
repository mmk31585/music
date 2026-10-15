# Proposed Muse Sprint Plan

Status: proposed, pending Sprint 0 evidence. No backlog item is active until current checkout verifies it.

## Sprint 0 — Baseline and truth reconciliation
Goal: establish actual state without broad code changes.
- Verify OpenCode version and agent/command/skill discovery.
- Inspect git status, branch, recent commits, submodule SHA and pre-existing changes in both repositories.
- Read AGENTS.md, STATE.md, BACKLOG.md, CONCERNS.md, ROADMAP.md and recent learning notes.
- Run available backend tests, vet, build and lint; run frontend type-check, unit tests and build if dependencies/environment permit.
- Reconcile docs/codebase/CONCERNS.md against docs/agent/STATE.md and current behavior.
- Confirm frontend architecture target and update the frontend's missing/stale guidance.
- Label each candidate verified open / verified done / blocked / stale / needs decision.

Exit: exact commands/results, submodule state, current blockers and deduplicated backlog. No broad refactor before this gate.

## Sprint 1 — Bounded accessibility and test gaps
Only activate verified-open items:
- Admin icon-only accessible names and loading/disabled states.
- Focused RTL audit backed by screenshots/browser evidence and regression checks.
- Missing auth/catalog/player store tests.
- Missing upload/admin handler tests.
- Playwright smoke-test feasibility and deterministic test data.

Exit: intended behavior only, positive/negative/error tests, no unrelated cleanup.

## Sprint 2 — Authorization and architecture inventory
- Document current RequireRole behavior and intended role semantics.
- Do not let higher numeric roles inherit all permissions unless product/security explicitly intends it.
- Inventory duplicated middleware, ordering and failure responses.
- Map frontend components/composables/services/pages/stores and import graph.
- Propose incremental feature migration; no bulk moves before approval.
- Treat backend DTO/domain separation as a separate design and migration proposal.

Exit: accepted proposals and test plans before structural changes.

## Sprint 3 — Enrichment status via SSE
- Capture current polling flow and status lifecycle.
- Agree event names, payload, auth, terminal states and compatibility.
- Native EventSource does not support arbitrary Authorization headers. Do not put bearer tokens in query strings; evaluate same-site cookie/session auth or fetch-based streaming.
- Plan reconnect/backoff, last-event handling, cancellation, proxy buffering/timeouts and resource limits.
- Add backend integration and frontend state tests before removing polling.
- Keep polling fallback until SSE reliability is evidenced.

Exit: authorized end-to-end test, reconnect/disconnect and terminal-state coverage.

## Sprint rules
- One primary outcome, one owner, one reviewer per task.
- Parallelize only after contracts and file ownership are explicit.
- Record acceptance criteria, test commands/results, changed paths and risks.
- Agent completion messages are not evidence.
- Do not perform production changes without human approval.
