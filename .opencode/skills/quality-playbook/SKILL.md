---
name: quality-playbook
description: 'Run a complete quality engineering audit on any codebase. Produces requirements, tests, code review, spec audit, TDD-verified patches. By Andrew Stellman.'
license: MIT
metadata:
  version: 1.5.6
  author: Andrew Stellman
  github: https://github.com/andrewstellman/quality-playbook
---

# Quality Playbook

Run a complete quality engineering audit on any codebase — produces requirements, tests, code review, spec audit, and TDD-verified patches.

## Phases
1. **Explore** — Codebase exploration, domain analysis, risk identification
2. **Generate** — Requirements, constitution, functional tests, protocols
3. **Code Review** — Three-pass review with regression tests for confirmed bugs
4. **Spec Audit** — Three independent AI auditors (Council of Three)
5. **Reconciliation** — TDD red-green verification for all bugs
6. **Verify** — Self-check benchmarks against all artifacts

## Output Artifacts
| File | Purpose |
|------|---------|
| `quality/EXPLORATION.md` | Phase 1 findings |
| `quality/REQUIREMENTS.md` | Testable requirements with use cases |
| `quality/QUALITY.md` | Quality constitution |
| `quality/CONTRACTS.md` | Behavioral contracts |
| `quality/test_functional.*` | Automated functional tests |
| `quality/RUN_CODE_REVIEW.md` | Three-pass review protocol |
| `quality/BUGS.md` | Consolidated bug report with patches |
| `quality/COMPLETENESS_REPORT.md` | Final gate verdict |

## Usage
Run one phase at a time for best results. Default: Phase 1 only.
- "Run quality playbook phase N"
- "Continue to next phase"
- "Run all phases" — runs 1-6 sequentially

## Runs On
Any language. Reads `reference_docs/` for enriched analysis.
