# ./

This template should help get you started developing with Vue 3 in Vite.

## Recommended IDE Setup

[VS Code](https://code.visualstudio.com/) + [Vue (Official)](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Recommended Browser Setup

- Chromium-based browsers (Chrome, Edge, Brave, etc.):
  - [Vue.js devtools](https://chromewebstore.google.com/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)
  - [Turn on Custom Object Formatter in Chrome DevTools](http://bit.ly/object-formatters)
- Firefox:
  - [Vue.js devtools](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)
  - [Turn on Custom Object Formatter in Firefox DevTools](https://fxdx.dev/firefox-devtools-custom-object-formatters/)

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
npm run test:unit
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```

## GitLab CI Pipeline

This repository includes a production-oriented pipeline in `.gitlab-ci.yml` for npm + Vue 3 + Vite + TypeScript.

### Pipeline triggers

The pipeline runs for:

- Merge request pipelines
- Default branch pipelines
- `main`/`master` branch pipelines
- Tag pipelines

### Stage order

1. `prepare`
2. `quality`
3. `test`
4. `build`
5. `docs`
6. `security`

### Runtime and install strategy

- Default image: `node:22.12.0-bookworm-slim` (override with `NODE_IMAGE`)
- Deterministic installs: `npm ci`
- npm cache path: `$CI_PROJECT_DIR/.npm` via `NPM_CONFIG_CACHE`
- Jobs are `interruptible` to reduce wasted CI time on superseded pipelines

### Jobs

- `install`:
  - Validates Node/npm versions
  - Runs `npm ci --prefer-offline --no-audit`
- `typecheck`:
  - Runs `npm run -s type-check`
- `lint_oxlint`:
  - Runs check-only Oxlint via `ci/lint-oxlint-check.sh`
  - No `--fix` in CI
- `lint_eslint`:
  - Runs check-only ESLint via `ci/lint-eslint-check.sh`
  - Keeps `--cache`, no `--fix` in CI
- `prettier_check`:
  - Runs `prettier --check` via `ci/prettier-check.sh`
  - No `--write` in CI
- `test_unit`:
  - Runs Vitest in CI mode
  - Publishes `reports/junit.xml` and `reports/vitest.log`
  - Artifacts retained for 2 weeks
- `build_app`:
  - Runs `npm run -s build`
  - Publishes `dist/` artifacts for 1 week
- `docs_build`:
  - Runs `npm run -s docs:gen` then `npm run -s docs:build`
  - Runs on default branch/tags, and on MRs only when docs-related files change
  - Publishes `docs/.vitepress/dist/` for 1 week
- `security_audit`:
  - Runs `npm run -s audit:ci` (`npm audit --omit=dev`)
  - Uses `NPM_CONFIG_AUDIT_LEVEL=high`

## Security Audit Policy

Use production-only audit results for fail/pass gates:

```sh
npm run audit:ci
```

Useful audit commands:

```sh
npm run audit:prod
npm run audit:full
```

Notes:

- `audit:prod` checks runtime dependencies only and is the required security gate.
- `audit:full` includes dev tooling and is informational.
- Current dev-only advisories are from the ESLint 9 toolchain (`ajv`/`minimatch`) and are tracked for a later ESLint 10 migration.
