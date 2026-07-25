---
name: security-review
description: 'AI-powered codebase security scanner that traces data flows and catches vulnerabilities across JS, TS, Python, Java, PHP, Go, Ruby, Rust.'
---

# Security Review

An AI-powered security scanner that reasons about codebases to find vulnerabilities that pattern-matching tools miss.

## Workflow
1. **Scope Resolution** — determine what to scan, identify languages/frameworks
2. **Dependency Audit** — check for known CVEs in dependencies
3. **Secrets & Exposure Scan** — hardcoded keys, tokens, passwords
4. **Vulnerability Deep Scan** — SQLi, XSS, command injection, IDOR, crypto, SSRF
5. **Cross-File Data Flow** — trace user input from entry to sink
6. **Self-Verification** — re-read findings, confirm exploitability
7. **Report** — findings summary with file:line and fix suggestions

## Severity
- 🔴 CRITICAL: SQLi, RCE, auth bypass
- 🟠 HIGH: XSS, IDOR, hardcoded secrets
- 🟡 MEDIUM: CSRF, open redirect, weak crypto
- 🔵 LOW: Best practice violations

## Muse-Specific
- JWT handling (generation, validation, refresh)
- Media upload validation
- SQL injection in catalog search
- XSS in user-generated content
- RLS for data isolation
- WebSocket security for social features

## Reference Files
- `references/vuln-categories.md` — Detection signals for all categories
- `references/secret-patterns.md` — Regex patterns by provider
- `references/language-patterns.md` — Framework-specific patterns
- `references/vulnerable-packages.md` — Curated CVE watchlist
- `references/report-format.md` — Output template
