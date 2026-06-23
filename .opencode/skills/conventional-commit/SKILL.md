---
name: conventional-commit
description: 'Generate conventional commit messages using structured format. Follows the Conventional Commits specification.'
---

# Conventional Commit

Generate standardized, descriptive commit messages following the Conventional Commits specification.

## Workflow
1. Run `git status` to review changed files
2. Run `git diff --cached` to inspect changes
3. Construct commit message using the structure below
4. Execute: `git commit -m "type(scope): description"`

## Structure
```
type(scope): description

[optional body]

[optional footer]
```

## Types
- `feat` — new feature
- `fix` — bug fix
- `docs` — documentation only
- `style` — formatting, whitespace
- `refactor` — code change with no feature/fix
- `perf` — performance improvement
- `test` — adding tests
- `build` — build system changes
- `ci` — CI configuration
- `chore` — maintenance tasks
- `revert` — revert previous commit

## Examples
```
feat(player): add crossfade between tracks
fix(catalog): correct album art URL for RTL layout
refactor(auth): extract JWT validation into middleware
feat(api)!: migrate to v2 endpoints (BREAKING CHANGE)
```

## Validation
- Description in imperative mood ("add" not "added")
- Scope is optional but recommended
- Breaking changes: append `!` after type/scope and add `BREAKING CHANGE:` in footer
