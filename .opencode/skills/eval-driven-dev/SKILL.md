---
name: eval-driven-dev
description: 'Build an automated evaluation pipeline for Python LLM applications. Define eval criteria, instrument the app, build golden datasets, run evaluations, analyze results. Uses pixie-qa framework.'
license: MIT
compatibility: Python 3.10+
metadata:
  version: 0.8.4
---

# Eval-Driven Development for Python LLM Applications

Build an automated evaluation pipeline that tests a Python-based AI application end-to-end, scoring outputs using evaluators.

## The Workflow (Steps 1-6)

### Step 1: Understand App & Define Eval Criteria
1a. Project analysis → `pixie_qa/00-project-analysis.md`
1b. Entry point analysis → `pixie_qa/01-entry-point.md`
1c. Eval criteria → `pixie_qa/02-eval-criteria.md`

### Step 2: Instrument & Capture Traces
2a. Add `wrap()` at data boundaries
2b. Implement Runnable class
2c. Capture reference trace → `pixie_qa/reference-trace.jsonl`

### Step 3: Define Evaluators
Map each eval criterion to built-in or custom evaluator functions.

### Step 4: Build Dataset
Create test scenarios tying together runnable, evaluators, and use cases.

### Step 5: Run `pixie test`
Fix mechanical issues until tests produce scores.

### Step 6: Analyze Outcomes
Score pending evaluations, produce action plan.

## Key Principles
- App's LLM calls go to **real LLM** — never mock
- Test the app's code path, not mocked internals
- Each eval criterion needs a corresponding evaluator
- Output: working `pixie test` run with real scores

## Reference Files
- `references/1-a-project-analysis.md` through `references/6-analyze-outcomes.md`
- `resources/setup.sh` — Initialize the evaluation environment
