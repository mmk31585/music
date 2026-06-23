# Muse — Frontend

> Premium Persian Music Ecosystem — Vue 3 SPA

## Stack

Vue 3 (Composition API + `<script setup>`) · TypeScript (strict) · Vite 7 · Pinia 3 · PrimeVue 4 · Tailwind CSS v4 · Vue Router 4 · Axios · Zod 4

## Setup

```bash
cp .env.example .env   # Configure VITE_API_BASE_URL
npm install
npm run dev            # Vite dev server → :5173
```

## Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Vite dev server with HMR |
| `npm run build` | Type-check + production build |
| `npm run lint` | ESLint check |
| `npm run type-check` | vue-tsc type check |
| `npm run test:unit` | Vitest unit tests |
| `npm run audit:ci` | Production-only npm audit |

## Directory Structure

```
src/
├── assets/          # CSS, fonts, images
├── components/      # UI components by domain (admin/, auth/, common/, music/, social/)
├── composables/     # 27 composables (useAuth, usePlayer, useCatalogSearch, etc.)
├── layouts/         # LayoutMusicApp, LayoutAuth, LayoutAdmin, LayoutEmpty
├── pages/           # Page components (app/, admin/, auth/, errors/)
├── plugins/         # Axios client, request factory
├── router/          # Vue Router config + domain route files
├── services/        # API modules, audio engine, WebSocket, storage
├── stores/          # Pinia stores (player, user-auth, feature-flags)
├── types/           # TypeScript type definitions
└── utils/           # PrimeVue preset, utility functions
```

## CI/CD

GitHub Actions pipeline in `.github/workflows/ci.yml`. Stages: go-lint → go-test → frontend-lint → frontend-test → frontend-build → security-scan → docker-build.
Also includes deploy workflow in `.github/workflows/deploy.yml` triggered on version tags.

## Docs

- [`AGENTS.md`](AGENTS.md) — Agent team architecture (file ownership, contracts)
- [`PERFORMANCE.md`](PERFORMANCE.md) — Performance-sensitive code guidelines
- [`../docs/README.md`](../docs/README.md) — Project documentation index
