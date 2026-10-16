---
description: Reconcile verified backlog items and propose a dependency-aware sprint
agent: muse-team
---

Prepare a sprint proposal for: $ARGUMENTS

Perform intake using docs/engineering-os/WORKFLOW.md. Use BACKLOG_RECONCILIATION.md and SPRINT_PLAN.md as hypotheses, then verify against the current checkout.

Request product framing, UX evidence, architecture dependencies and QA/security risk analysis where relevant. Do not assume agents share context or can message each other; collect written handoffs through the orchestrator.

Output sprint goal/success signal; verified-open and excluded stale/done items; task contracts with IDs/owners/file scopes/acceptance criteria; dependency graph and safe parallel groups; evidence-backed estimates or UNKNOWN; test gates/rollback; risks/approvals/exit criteria.

Planning does not authorize all tasks for implementation. Only explicitly selected tasks may enter IN_PROGRESS.
