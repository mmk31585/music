---
description: Run a read-only evidence-based Muse release readiness review
agent: muse-release-manager
---

Perform a read-only release-gate review for $ARGUMENTS.

Read docs/engineering-os/WORKFLOW.md and SAFETY_POLICY.md, inspect the full diff and recorded test output, and use independent code/security/QA reports when available.

Return READY FOR HUMAN RELEASE APPROVAL, NOT READY, or INSUFFICIENT EVIDENCE. List missing gates and concrete risks. Do not execute deployment, merge, migration or release commands.
