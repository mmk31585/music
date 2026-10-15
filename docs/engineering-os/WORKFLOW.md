# Engineering Workflow

## Task states
BACKLOG → INTAKE → RESEARCH → READY → IN_PROGRESS → VERIFYING → REVIEW → INTEGRATION → DONE

Exceptional states: BLOCKED, REJECTED, CANCELLED. Failed verification returns to implementation with a concrete report. Never skip gates to make status look green.

## Intake
- Read applicable AGENTS.md, STATE.md, BACKLOG.md, CONCERNS.md and recent LEARNINGS.md.
- Inspect repository boundaries, branch, git status, recent commits and submodule state.
- Separate pre-existing failures from task-introduced failures.
- State problem, outcome, scope, non-goals and unknowns.
- Do not infer that an old backlog item is open or fixed from its title alone.

## Research and product framing
- Product Manager defines user outcome, story, acceptance criteria, priority and success signal.
- UX Researcher uses actual screens, code, user feedback, browser evidence and current primary sources.
- Separate observations, repository evidence, external facts and hypotheses.
- Record source URLs, alternatives, trade-offs and research date.
- Never claim browser/usability observations without evidence.

## Architecture and planning
- Architect checks boundaries, API contract, data flow, security, compatibility and rollback.
- Cross-repo tasks begin with an agreed contract and owned file lists.
- Breaking API/config/data changes require a decision record and explicit approval.
- Every task has one owner, independent reviewer, dependencies, allowed paths, forbidden changes, criteria, tests and rollback notes.
- Parallelize only disjoint file ownership and explicit interface boundaries. One writer per file.

## Implementation
- Work on a feature branch from a confirmed base; preserve pre-existing changes.
- Add tests with behavior changes; follow actual repository scripts.
- Never hand-edit generated API types.
- Ask before dependency/lockfile changes, auth/security policy changes, public API changes or shared-data operations.

## Verification
- Record exact commands and exit statuses. Run focused tests, then required gates.
- Backend commands are in Makefile/AGENTS.md; commonly make test, go vet ./..., make build/build-all.
- Inspect frontend package.json before choosing commands. Current package defines type-check, test:unit:run and build.
- WARNING: frontend npm run lint includes eslint . --fix --cache and mutates files. Use non-fixing ESLint diagnostics.
- E2E needs an isolated account/environment and deterministic cleanup.
- A passing build does not prove runtime, auth, audio, SSE or deployment behavior.

## Independent review and completion
Code Review checks correctness, regressions, architecture and tests. Security reviews relevant threat surfaces. QA checks acceptance criteria and observable behavior, including loading/error/retry, RTL, keyboard and accessibility where relevant.

DONE requires criteria satisfied, checks executed and recorded, blocking findings resolved, only intended files changed, docs updated and approvals met. Report outcome, paths, exact checks/results, review verdict, branch/commit, blockers, risks and next step.
