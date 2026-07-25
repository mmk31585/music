---
name: acquire-codebase-knowledge
description: 'Map, document, and onboard into an existing codebase. Produces 7 structured docs in docs/codebase/ covering stack, structure, architecture, conventions, integrations, testing, and concerns.'
license: MIT
compatibility: 'Cross-platform. Requires Python 3.8+ and git.'
metadata:
  version: "1.3"
---

# Acquire Codebase Knowledge

Produces seven populated documents in `docs/codebase/` covering everything needed to work effectively on the project.

## Workflow
1. Run `python3 scripts/scan.py --output docs/codebase/.codebase-scan.txt` from project root
2. Search for PRD, README, ROADMAP, SPEC, DESIGN files
3. Investigate using inquiry checkpoints
4. Populate all seven templates: STACK.md, STRUCTURE.md, ARCHITECTURE.md, CONVENTIONS.md, INTEGRATIONS.md, TESTING.md, CONCERNS.md
5. Validate all docs against checkpoints

## Phase 1: Scan
Run scan script, read intent documents, summarize project intent.

## Phase 2: Investigate
Use inquiry checkpoints to probe each of the seven documentation areas.

## Phase 3: Populate Templates
Copy templates and fill in order: STACK → STRUCTURE → ARCHITECTURE → CONVENTIONS → INTEGRATIONS → TESTING → CONCERNS

## Phase 4: Validate
Cross-reference every claim, fix gaps, present summary with [ASK USER] items.

## Bundled Assets
| Asset | Use |
|-------|-----|
| `scripts/scan.py` | Run first in Phase 1 |
| `references/inquiry-checkpoints.md` | Phase 2 investigation questions |
| `references/stack-detection.md` | Only if stack is ambiguous |
| `assets/templates/` | Seven templates for Phase 3 |
