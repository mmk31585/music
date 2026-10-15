# Muse AI Engineering Operating System

This governance layer extends the existing OpenCode setup; it does not replace it. Reuse existing agents, skills and references under .opencode rather than duplicating personas.

## Entry points
- Workflow: WORKFLOW.md
- Autonomy and safety: SAFETY_POLICY.md
- Agent mapping: ROLE_CATALOG.md
- Backlog reconciliation: BACKLOG_RECONCILIATION.md
- Proposed sprint: SPRINT_PLAN.md
- Task contract: TASK_CONTRACT.md
- Product strategy: ../MASTER_PLAN.md
- Delivery roadmap: ../ROADMAP.md
- Operating culture: ../ENGINEERING_OPERATING_SYSTEM.md
- Current state: ../agent/STATE.md and ../agent/BACKLOG.md
- Known concerns: ../codebase/CONCERNS.md

## Commands
- /muse-intake — read-only source-backed intake report.
- /muse-plan-sprint — reconcile backlog and propose a dependency-aware sprint.
- /muse-task — execute one scoped task.
- /muse-release-check — read-only release readiness verdict.

Use @muse as primary coordinator and @muse-team for cross-domain coordination. Product Manager, UX Researcher and Release Manager are distinct new roles.

## Principles
1. Current code and executed checks outrank stale status docs.
2. High autonomy applies to reversible, in-scope work on a feature branch.
3. Human approval is mandatory for production deployment, destructive data operations, shared-environment migrations, credential rotation, security exceptions, and externally visible merges/releases.
4. Do not claim direct agent messaging or shared live state unless a real integration is verified.
5. Never claim a test passed without recorded evidence.
6. Preserve uncommitted work; never stage unrelated files.

## OpenCode compatibility
The checked-in config uses a V1-style configuration shape. Do not migrate its schema as part of backlog work. Verify the installed version before relying on permission syntax.

Official references:
- https://opencode.ai/docs/agents/
- https://opencode.ai/docs/config/
- https://opencode.ai/v2/docs/migrate-v1/
