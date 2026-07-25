---
name: secret-scanning
description: 'Guide for configuring GitHub secret scanning, push protection, custom patterns, and alert remediation.'
---

# Secret Scanning

Guide for configuring GitHub secret scanning — detecting leaked credentials, preventing pushes, defining custom patterns, and managing alerts.

## Enable Secret Scanning
1. Settings → Advanced Security → Enable Secret Protection
2. Enable Push Protection
3. Configure exclusions in `.github/secret_scanning.yml`
4. Enable non-provider patterns, AI detection, validity checks

## Resolve Blocked Pushes
- **Remove secret**: amend commit or interactive rebase
- **Bypass**: visit URL from error, select reason, allow push
- **Request bypass**: if delegated bypass is enabled

## Custom Patterns
Organization-specific regex patterns with dry-run validation before publishing.

## Alert Management
- Rotate credential immediately (critical)
- Check validity status (active/inactive/unknown)
- Dismiss with reason (false positive, revoked, used in tests)

## Pre-Commit Scanning via Agent
Install Advanced Security plugin for in-agent secret scanning:
- Copilot CLI: `/plugin install advanced-security@copilot-plugins`
- VS Code: Chat: Plugins → install advanced-security

## Reference Files
- `references/push-protection.md` — Bypass mechanics
- `references/custom-patterns.md` — Regex creation
- `references/alerts-and-remediation.md` — Alert management
