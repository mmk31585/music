# Autonomy and Safety Policy

High autonomy is for routine, bounded work in development. It is not permission to make irreversible changes or operate production.

## Autonomous by default
- Read instructions, inspect Git state and source, compare diffs, and inspect redacted logs.
- Research primary sources and record URLs/dates.
- Plan tasks, write tests, edit scoped files on a feature branch, and run local tests/type checks/static analysis/builds.
- Fix in-scope failures introduced by the task.
- Prepare focused Conventional Commit messages; stage only explicit task paths after reviewing the diff.

## Ask first
- Install/update/remove dependencies or change lockfiles.
- Run migrations or change shared-environment data/schema.
- Start/stop services/containers that may be shared, change persistent volumes, or run cleanup scripts.
- Read outside the approved workspace.
- Change credentials, payment settings, auth strategy, CORS, firewall/network policy or security exceptions.
- Change public API contracts incompatibly.
- Push, merge, tag, deploy, or change hosting/provider settings.
- Any action with unclear blast radius.

## Never autonomously
- Deploy to production or claim production success.
- Run destructive Git cleanup, database drops, volume deletion or broad data cleanup.
- Print secrets, tokens, passwords, cookies, .env contents or raw user data.
- Disable/delete tests to make checks pass or misrepresent unexecuted checks.
- Weaken security checks to bypass failures.
- Change unrelated user edits or stage all files when the worktree has unrelated modifications.
- Merge competing changes without reviewing diffs and compatibility.

## Baseline protection
1. Record git status before edits; identify pre-existing changes.
2. Never use git add -A or git commit -a.
3. Review git diff --check, git diff --stat, the full diff and staged paths before commit.
4. If ownership is unclear, stop and report.

OpenCode permission syntax varies by major version. Verify the installed version before depending on configuration permissions. Prompt instructions do not replace OS/container/provider access controls.
