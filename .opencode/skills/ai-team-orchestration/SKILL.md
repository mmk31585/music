---
name: ai-team-orchestration
description: 'Bootstrap and run a multi-agent AI development team. Sprint planning, parallel dev/QA, brainstorm prompts, context recovery.'
---

# AI Team Orchestration

Bootstrap and run a multi-agent AI development team for software projects.

## Team Roles
| Role | Name | Focus |
|------|------|-------|
| Producer | Remy | Sprint planning, coordination, merging PRs |
| Product Designer | Kira | UX, mechanics, user experience |
| Visual Director | Milo | CSS, animations, visual identity |
| Frontend Engineer | Nova | UI framework, state, components |
| Backend Engineer | Sage | API, database, auth, security |
| DevOps Engineer | Dash | CI/CD, cloud deployment, pipelines |
| QA Engineer | Ivy | E2E tests, automation, playtesting |

## Chat Architecture
Human (CEO) acts as message bus between parallel chats:
- Producer → Plans/Merges (never writes code)
- Dev Team → Feature branches (Nova + Sage + Milo)
- QA Team → Test branches (Ivy)
- DevOps → On-demand (Dash)

## Project Bootstrap
1. Create `PROJECT_BRIEF.md` — single source of truth
2. Run brainstorm with distinct agent perspectives
3. Create sprint plans (`docs/sprint-N/plan.md`)
4. Execute sprints with progress tracking
5. QA sign-off per sprint

## Context Recovery
- Save progress before context overflow
- Cold start: read PROJECT_BRIEF.md + sprint progress
- Anti-patterns: no rebasing, no batch commits, file GitHub Issues

## Reference Files
- `references/project-brief-template.md`
- `references/brainstorm-format.md`
- `references/sprint-plan-template.md`
- `references/anti-patterns.md`
