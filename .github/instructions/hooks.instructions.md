---
description: 'Portable guidance for authoring safe, fast hooks in Muse. Covers session-logger, secrets-scanner, and guardrails'
applyTo: '.github/hooks/**'
---

# Hook Authoring Guidelines for Muse

Hooks are **small, deterministic commands or scripts** that run at specific lifecycle events in GitHub Copilot.

## Folder Structure

```
.github/
└── hooks/
    ├── session-logger.json       ← hook config
    ├── secrets-scanner.json      ← hook config
    ├── tool-guardian.json        ← hook config
    └── scripts/
        ├── session-logger.sh
        ├── secrets-scanner.sh
        └── tool-guardian.sh
```

## Config Format

```json
{
  "version": 1,
  "hooks": {
    "preToolUse": [
      {
        "type": "command",
        "bash": "./.github/hooks/scripts/tool-guardian.sh",
        "cwd": ".",
        "timeoutSec": 5
      }
    ]
  }
}
```

## Event Types

| Event | Use Case |
|-------|----------|
| `sessionStart` | Setup, validation, context injection |
| `sessionEnd` | Cleanup, summaries |
| `userPromptSubmitted` | Auditing, prompt blocking |
| `preToolUse` | Guardrails, deny/block, argument modification |
| `postToolUse` | Logging, formatting |
| `errorOccurred` | Diagnostics, alerts |

## Script Contract

- Read JSON from stdin
- Respond through exit code and stdout
- Exit 0 = allow; non-zero = block
- For `preToolUse`: stdout can carry `{"permissionDecision":"deny","permissionDecisionReason":"..."}`

## Muse Hooks Plan

### 1. Session Logger
Logs session start/end and user prompts for audit.

### 2. Secrets Scanner
Scans for accidentally committed secrets in tool output.

### 3. Tool Guardian
Blocks dangerous commands (rm -rf /, git push --force, etc.).

## Universal Design Rules

- One hook, one responsibility
- Keep hooks synchronous, bounded, and non-interactive
- Make hooks deterministic and idempotent
- Treat prompts, tool arguments, and tool output as untrusted
- Redact secrets from logs
- Use strict modes: Bash `set -euo pipefail`
