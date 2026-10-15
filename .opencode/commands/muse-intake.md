---
description: Inspect current Muse state and produce a source-backed intake report
agent: muse
---

Run the Muse engineering intake.

1. Read AGENTS.md and repo-specific instructions. Read docs/engineering-os/README.md, WORKFLOW.md, SAFETY_POLICY.md, BACKLOG_RECONCILIATION.md, docs/agent/STATE.md, docs/agent/BACKLOG.md, docs/codebase/CONCERNS.md and recent LEARNINGS.md when present.
2. Inspect current branch, git status, recent commits, submodule state and OpenCode config. Never expose secret contents.
3. Identify contradictions between current code/commits and status documents.
4. Run safe, available, read-only baseline checks only. Do not start/stop services, install dependencies, migrate databases, deploy or modify code.
5. Delegate read-only research where it adds distinct evidence, then consolidate written handoffs.
6. Report verified facts, pre-existing blockers, priorities, unknowns and next smallest task.

Do not mark stale backlog items done/open without evidence. Intake is read-only.
