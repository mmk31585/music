# Muse Backlog Reconciliation

Inspected remotely on 2026-10-16. This is a repository discovery report, not proof that tests currently pass.

## Verified
- Backend repo: mmk31585/music; main head at inspection: 14fd1029e085fc8e1be0e878c0f8bd0095c9acdc (2026-10-14).
- Backend .gitmodules points frontend/ to https://github.com/mmk31585/music-frontend.git.
- Backend instructions describe Go/Gin, PostgreSQL, Redis, OpenSearch, Python ML, and frontend/ as a separately owned submodule.
- docs/agent/STATE.md and docs/agent/BACKLOG.md claim all phases complete, but STATE's build table is dated 2026-07-09.
- docs/codebase/CONCERNS.md still describes several critical/high issues inconsistent with the later STATE/BACKLOG claims.
- Frontend README links to AGENTS.md, but that file was not present in the inspected music-frontend default branch.
- Frontend package.json currently uses Vue 3, Vite 7, TypeScript, Pinia, PrimeVue 4, Tailwind CSS 4, Axios, Zod, Vitest and vite-plugin-pwa.
- Frontend npm run lint includes eslint . --fix --cache, so it mutates source.
- Frontend README lists audit:ci, but the inspected package.json does not define that script.
- The frontend instructions supplied in the work request describe a materially different architecture (plain CSS tokens, generated OpenAPI client, Vercel Blob upload and other scripts). Do not silently impose them on the current checked-in repository.

## Consequence
Do not begin with a blind feature refactor. Sprint 0 must inspect the actual local checkout, confirm intended frontend architecture, run available checks and reconcile docs. Remote inspection is not a substitute for local tests, browser testing or runtime checks.

## Candidate backlog from supplied task list
| ID | Candidate | Initial risk | Disposition |
|---|---|---|---|
| MUSE-UX-01 | RTL consistency audit | Low-medium | Inspect actual UI/screens; fix verified issues with regression checks |
| MUSE-A11Y-01 | Admin icon-only accessible names | Low | Verify current labels and component structure first |
| MUSE-TEST-FE-01 | Pinia store unit tests | Low-medium | Inventory current tests and package scripts |
| MUSE-TEST-BE-01 | Upload/admin handler tests | Medium | Inspect handlers, DB seams and existing mocks |
| MUSE-E2E-01 | Playwright critical journey smoke tests | Medium | Confirm tooling, isolated test account and cleanup |
| MUSE-ARCH-FE-01 | Feature/domain frontend layout | Medium-high | Map actual imports and ownership; migrate one domain at a time |
| MUSE-ARCH-BE-01 | Backend middleware consolidation | Medium-high | Inventory order/behavior/tests before centralizing |
| MUSE-AUTHZ-01 | RequireRole hierarchy semantics | High | Decide hierarchy versus independent capabilities first |
| MUSE-STREAM-01 | Enrichment SSE | Medium-high | Contract-first; cover auth, reconnect, terminal state and proxy behavior |
| MUSE-ARCH-DOMAIN-01 | Backend DTO/domain separation | High | Separate architecture proposal; not first-sprint cleanup |

## Sprint 0 questions
1. Does the supplied frontend architecture describe this exact repository/target branch?
2. Which CONCERNS entries remain real versus stale?
3. Which candidate items are already implemented?
4. What OpenCode version is installed locally, and which agents/commands/skills are discovered?
5. What scripts and test gates actually work in the local checkout?
