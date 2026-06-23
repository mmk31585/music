---
name: breakdown-plan
description: 'Generate comprehensive project plans with Epic > Feature > Story/Enabler > Test hierarchy, dependencies, priorities, and automated issue tracking.'
---

# Project Planning & Issue Breakdown

Generate comprehensive GitHub project plans with Agile work item hierarchy and automated tracking.

## Work Item Hierarchy
- **Epic**: Large business capability (milestone level)
- **Feature**: Deliverable functionality within an epic
- **Story**: User-focused requirement delivering independent value
- **Enabler**: Technical infrastructure supporting stories
- **Task**: Implementation work breakdown

## Process
1. Collect feature artifacts (PRD, UX design, technical breakdown, test plan)
2. Create work item hierarchy with dependencies
3. Prioritize using value/effort matrix
4. Set up GitHub project board (Backlog → Ready → In Progress → Review → Done)
5. Generate GitHub issues with templates

## Issue Templates
- Epic: Business value, success metrics, feature list, DoD
- Feature: User stories, enablers, dependencies, acceptance criteria
- Story: INVEST criteria, acceptance criteria, technical tasks
- Enabler: Technical requirements, enabled stories, DoD

## Estimation
- Stories: Fibonacci (1, 2, 3, 5, 8, 13)
- Epics/Features: T-shirt sizes (XS-XXL)
- Sprint capacity with 20% buffer

## Output
- `docs/plan/{epic}/{feature}/project-plan.md`
- `docs/plan/{epic}/{feature}/issues-checklist.md`
