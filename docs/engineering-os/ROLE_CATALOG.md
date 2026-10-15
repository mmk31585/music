# Agent Role Catalog

Reuse existing OpenCode specialists. Add new agents only when they provide a distinct capability.

| Role | Existing agent(s) | Responsibility |
|---|---|---|
| CTO / orchestrator | muse | Intake, scope, delegation, conflict resolution, gates, final report |
| Engineering manager | muse-team | Dependency-aware task list, ownership and integration |
| Product Manager | muse-product-manager (new) | User outcomes, priority, acceptance criteria, scope |
| UX Researcher | muse-ux-researcher (new) | Evidence-backed friction research and journey analysis |
| UX / UI | se-ux-ui-designer.agent, muse-ui | Interaction, responsive UI, RTL and accessibility |
| Architect | api-architect.agent, arch.agent | Boundaries, contracts, data flow, ADRs and migration strategy |
| Backend / Go | muse-auth, muse-catalog, muse-admin, go-mcp-expert.agent | APIs, backend logic, persistence and tests |
| Frontend / Vue | vuejs-expert.agent, muse-ui, muse-player, muse-infra | Vue, stores, player, routing and API integration |
| QA | test-engineer, playwright-tester.agent, accessibility-runtime-tester.agent | Test plans, browser tests, accessibility and reproducible defects |
| Security | security-auditor, se-security-reviewer.agent | Threat model, auth, secrets, abuse cases |
| Independent reviewer | code-reviewer, doublecheck.agent | Review-only correctness/regression verdict |
| DevOps / SRE | muse-sre, muse-infra, devops-expert.agent, github-actions-expert.agent | CI, operational risk and observability |
| Release Manager | muse-release-manager (new) | Read-only release gate and evidence review |
| Knowledge curator | documentation-writer, documentation-and-adrs skills | Status, lessons, ADRs and handoff |

## Coordination rules
- The parent orchestrator owns durable task state and conflict resolution.
- Each specialist returns findings, paths changed, exact checks/results, blockers, risks and recommendation.
- Do not assume direct agent messaging or a shared live task list unless a concrete tool implements it.
- PM/UX recommend; owner/PM approves scope; Architect approves technical approach; QA and Review verify independently.
- Release and review-only agents must not edit the patch they approve.
