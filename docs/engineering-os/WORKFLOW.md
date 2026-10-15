# Engineering Workflow

## Task states
BACKLOG → INTAKE → RESEARCH → READY → IN_PROGRESS → VERIFYING → REVIEW → INTEGRATION → DONE

Exceptional states: BLOCKED, REJECTED, CANCELLED. Failed verification returns to implementation with a concrete report. Never skip gates to make status look green.

## Intake
- Read applicable AGENTS.md, STATE.md, BACKLOG.md, CONCERNS.md and recent LEARNINGS.md.
- Inspect repository boundaries, current branch, git status, recent commits and submodule state.
- Separate pre-existing failures from task-introduced failures.
- State problem, outcome, scope, non-goals and unknowns.
- Do not infer that an old backlog item is open or fixed from its title alone.

## Research and product framing
- Product Manager defines user outcome, user story, acceptance criteria, priority and success signal.
- UX Researcher uses actual screens, code, user feedback, browser evidence and current primary sources where useful.
- Separate observations, repository evidence, external facts and hypotheses.
- Record source URLs, alternatives, trade-offs and research date.
- Never claim browser/usability observations without the relevant evidence.

## Architecture and planning
- Architect checks boundaries, API contract, data flow, security, compatibility and rollback.
- Cross-repo tasks begin with an agreed contract and owned file lists.
- Breaking API/config/data changes require a decision record and explicit approval.
- Every task has one accountable owner, independent reviewer, dependencies, allowed paths, forbidden changes, acceptance criteria, tests and rollback notes.
- Parallelize only tasks with disjoint file ownership and explicit interface boundaries. One writer per file.
- Prefer focused changes and small reviewable commits.

## Implementation
- Work on a feature branch from a confirmed base. Preserve pre-existing changes.
- Add tests with behavior changes; follow actual repo scripts.
- Do not hand-edit generated API types.
- Ask before changing dependencies, lockfiles, auth/security policy, public API contracts or shared data.

## Verification
- Record exact commands and exit statuses. Run focused tests first, then required gates.
- Backend canonical commands are documented in Makefile/AGENTS.md; commonly make test, go vet ./..., make build/build-all.
- Inspect frontend package.json before choosing commands. The current package defines type-check, test:unit:run and build.
- WARNING: frontend npm run lint includes eslint . --fix --cache and mutates files. Use a non-fixing ESLint invocation for diagnostics.
- E2E needs an isolated test account/environment and deterministic cleanup.
- A passing build does not prove runtime, auth, audio, SSE or deployment behavior.

## Independent review and completion
Code Review checks correctness, regressions, architecture and tests. Security reviews relevant threat surfaces. QA checks acceptance criteria and observable behavior, including error/empty/loading/retry, RTL, keyboard and accessibility where relevant.

DONE requires acceptance criteria met, required checks executed and recorded, blocking findings resolved, only intended files changed, docs updated and approvals satisfied. Report outcome, paths, exact checks/results, review verdict, branch/commit, blockers, risks and next step.
