---
description: Execute one bounded Muse task with tests and independent review
agent: muse
---

Execute exactly one task: $ARGUMENTS

1. Confirm/create a task contract from docs/engineering-os/TASK_CONTRACT.md.
2. Confirm scope, criteria, owned paths, dependencies, base branch and baseline git status.
3. Stop if scope or baseline is unsafe/unclear.
4. Research only what is needed; record sources and assumptions.
5. Delegate distinct specialist tasks; never assign multiple writers to the same files or assume peer messaging.
6. Implement the smallest correct change on a feature branch.
7. Run focused tests and required gates; record exact commands/results.
8. Request independent code/security/QA review according to risk.
9. Resolve blockers and rerun affected checks.
10. Update status/learning/ADR artifacts and report paths, tests, branch/commit, risks and next action.

Never deploy, push, merge, run shared-environment migrations, rotate secrets or perform destructive cleanup without explicit human approval. Do not claim completion if required checks did not run.
