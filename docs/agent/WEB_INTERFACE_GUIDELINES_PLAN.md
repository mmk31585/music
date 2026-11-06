# Web Interface Guidelines — Frontend Fix Plan

Source: https://raw.githubusercontent.com/vercel-labs/web-interface-guidelines/main/command.md

Scope: all `src/**/*.vue` components + `src/App.vue`. TS-only files are out of scope (no UI).

## Epic 1: `transition: all` elimination (anti-pattern)

Guideline: "Never `transition: all`—list properties explicitly."

| # | File | Location |
|---|------|----------|
| 1.1 | `src/components/auth/LoginForm.vue` | line 122 (`transition-all`), line 186 (`transition: all 0.2s ease`) |
| 1.2 | `src/components/auth/RegisterForm.vue` | line 158 (`transition-all`), line 203 (`transition: all 0.2s ease`) |
| 1.3 | `src/components/music/home/HomeHero.vue` | lines 69, 78 (`transition-all`), line 284 (`transition: all 250ms`) |
| 1.4 | `src/components/music/home/HomeTrackCard.vue` | line 10 (`transition-all`) |
| 1.5 | `src/components/music/catalog/TrackRow.vue` | lines 3, 9 (`transition-all`) |
| 1.6 | `src/components/music/player/FullscreenPlayer.vue` | lines 18, 67, 73, 80, 136-139, 168, 181 |
| 1.7 | `src/components/music/player/ExpandedPlayer.vue` | lines 30, 44, 62, 70, 125, 127, 153, 167, 178, 193, 204, 219, 245, 255, 266, 275, 292 |
| 1.8 | `src/components/music/player/MobilePlayer.vue` | lines 13, 20 (`transition-all`) |
| 1.9 | `src/components/music/player/QueuePanel.vue` | line 289 (`transition: all 0.25s ease`) |
| 1.10 | `src/components/music/player/QueuePreview.vue` | lines 6, 69, 107, 128 (`transition-all`) |
| 1.11 | `src/components/music/player/PlaybackControls.vue` | lines 10, 62, 74, 95, 108, 168 (`transition-all`) |
| 1.12 | `src/components/music/player/VolumeControl.vue` | line 5 (`transition-all`) |
| 1.13 | `src/components/music/player/NowPlayingBarRoot.vue` | lines 12, 75, 93, 168, 179, 209, 230 (`transition-all`) |
| 1.14 | `src/components/music/player/ProgressBar.vue` | line 14 (`transition-all`) |
| 1.15 | `src/components/music/player/KaraokeLyrics.vue` | lines 35, 48, 60 (`transition-all`) |
| 1.16 | `src/components/music/layout/MusicRightPane.vue` | lines 5, 10, 94 (`transition-all`) |
| 1.17 | `src/components/music/layout/RightPaneNowPlaying.vue` | lines 3, 22, 84, 96, 120, 133, 142, 162 (`transition-all`) |
| 1.18 | `src/components/music/layout/RightPaneQueue.vue` | lines 16, 28, 58 (`transition-all`) |
| 1.19 | `src/components/music/layout/RightPaneSuggestions.vue` | lines 14, 37 (`transition-all`) |
| 1.20 | `src/components/admin/AdminStatCard.vue` | line 4 (`transition-all`) |

**Approach:** Replace `transition-all` / `transition: all Xs` with explicit property lists (e.g. `transition: color 150ms ease, background 150ms ease` or `transition: opacity 150ms ease, transform 150ms ease`). Each file = one story.

## Epic 2: Animations honor `prefers-reduced-motion`

Guideline: "Honor `prefers-reduced-motion` (provide reduced variant or disable)." + "Autoplay motion >5 seconds alongside other content needs pause, stop, or hide controls."

| # | File | Animation | Lines |
|---|------|-----------|-------|
| 2.1 | `src/components/auth/AuthHeroPanel.vue` | `eq-bounce` infinite equalizer bars | 57-97 |
| 2.2 | `src/components/music/home/HomeTrackCard.vue` | `eqBeat` infinite pulse | 130-139 |
| 2.3 | `src/components/music/catalog/TrackRow.vue` | `equalizer` infinite pulse | 271-290 |
| 2.4 | `src/components/music/player/QueuePreview.vue` | `animate-equalizer` / `eq-wave` | 226-252 |
| 2.5 | `src/pages/app/PageTrack.vue` | `now-playing`, `viz-pulse`, `similar-equalizer` | 858-904 |
| 2.6 | `src/components/common/AppLoader.vue` | `animate-spin` | 9-10 |
| 2.7 | `src/pages/app/PageHome.vue` | `animate-spin` pull-to-refresh | 20-21 |
| 2.8 | `src/components/music/home/HomeHero.vue` | `autoplay` loop video (>5s) | 14-23 |

**Approach:** Add `@media (prefers-reduced-motion: reduce)` rules that set `animation: none` / `transition: none` for each. For the hero video, add a pause control or stop autoplay under reduced motion. Each file = one story.

## Epic 3: Images need explicit dimensions (CLS)

Guideline: "`` needs explicit `width` and `height` (prevents CLS)." + "Below-fold images: `loading="lazy"`."

| # | File | Lines |
|---|------|-------|
| 3.1 | `src/components/common/CoverImage.vue` | 11-21 |
| 3.2 | `src/components/common/AppImage.vue` | 16-26 |

**Approach:** Add `width`/`height` props (default sensible values) and set them on the ``. Two stories.

## Epic 4: Skip link on every layout

Guideline: "Headings hierarchical `<h1>`–`<h6>`; include skip link for main content."

| # | File | Lines |
|---|------|-------|
| 4.1 | `src/layouts/LayoutAuth.vue` | 1-18 |

**Approach:** Add `<a href="#main-content" class="skip-link">Skip to main content</a>` matching the pattern in LayoutMusicApp.vue/LayoutAdmin.vue. One story.

## Epic 5: Error messages include remediation

Guideline: "Error messages include fix/next step, not just problem."

| # | File | Lines |
|---|------|-------|
| 5.1 | `src/components/auth/LoginForm.vue` | 104-113 (apiError alert) |
| 5.2 | `src/components/auth/RegisterForm.vue` | 139-148 (apiError alert) |

**Approach:** Append a next-step hint (e.g. "Try again" / "Check your credentials") to the displayed error message. Two stories.

## Epic 6: Placeholder typography

Guideline: "Placeholders end with `…` and show example pattern." + "`…` not `...`"

| # | File | Lines |
|---|------|-------|
| 6.1 | `src/components/common/AppLoader.vue` | 10 (`Loading...`) |
| 6.2 | `src/components/music/layout/MusicTopbar.vue` | 16 (`Search tracks, artists, albums...`) |

**Approach:** Replace `...` with `…` and make placeholders show an example pattern. Two stories.

## Epic 7: `outline: none` without focus replacement

Guideline: "Never `outline-none` / `outline: none` without focus replacement."

| # | File | Lines |
|---|------|-------|
| 7.1 | `src/components/common/ContextMenu.vue` | 547 (`outline: none` on `.ctx-menu` container) |

**Approach:** Replace `outline: none` with a `focus-visible` ring replacement. One story.

---

## Definition of Done per story

- `go test ./... -race -count=1` still passes (backend unaffected, run to confirm)
- Frontend: `npm run build` succeeds; `vue-tsc --build` has no **new** errors vs. baseline (10 pre-existing errors on `main` are known and out of scope)
- The specific guideline violation is resolved; manual spot-check in browser confirms
- Commit with Conventional Commit message: `fix(ui): <component> <guideline rule>`
- Commit size < ~200 changed lines

---

## Execution order (one story at a time)

Start with Epic 7 (single file, smallest), then Epic 3 (two small files), Epic 4, Epic 6, Epic 5, then Epic 2 (per-file), then Epic 1 (largest, per-file). After every 3-4 stories, run `npm run build` and `vue-tsc --build` to catch regressions early.
