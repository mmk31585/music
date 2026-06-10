# Agent Team Architecture

Parallel agent system for the Soundify frontend. Each agent owns a set of files and can work independently. A coordinator agent manages cross-cutting changes.

## How to Use

Invoke agents in parallel by describing the task and which agent(s) it belongs to:

```
Task: "Add album art color extraction to AlbumCard" → @catalog
Task: "Fix queue skip on error" → @player
Task: "Add rate limit display to login form" → @auth @ui
Task: "End-to-end: add new badge type" → @admin @creator @ui
```

## Roles

### @auth — Auth & Security
- **Focus**: Login/register flows, JWT, token refresh, auth guards, session restoration
- **Files**: `stores/user-auth.ts`, `composables/auth/*`, `components/auth/*`, `pages/auth/*`, `services/api/auth/*`, `plugins/client/*`, `router/index.ts` (guard section)
- **Contracts**: Exports `useAuth()`, `useUserAuthStore`, auth guard `requireAuth`
- **Deps on**: `@infra` (client, router), `@ui` (layout, toast)

### @player — Player & Audio
- **Focus**: Audio engine, queue, radio mode, visualizer, theatre, mini player, PiP, lyrics
- **Files**: `stores/player.ts`, `composables/player/*`, `composables/lyrics/*`, `services/player/*`, `services/api/player/*`, `services/api/lyrics/*`, `components/music/*Player*.vue`, `components/music/*Queue*.vue`, `components/music/*Radio*.vue`, `components/music/*Visualizer*.vue`, `components/music/*Theatre*.vue`, `components/music/*Ambient*.vue`, `components/music/*Mini*.vue`, `components/music/MobileBottomSheet.vue`, `components/music/NowPlayingBar.vue`, `components/music/KaraokeLyrics.vue`, `components/music/LyricsDisplay.vue`, `components/music/MusicSidebar.vue`, `components/music/MusicTopbar.vue`, `types/document-pip.ts`, `composables/usePlayerPiPController.ts`, `composables/useDocumentPictureInPicture.ts`
- **Contracts**: Exports `usePlayer()`, `usePlayerControls()`, `useQueueManager()`, `audioEngine`
- **Deps on**: `@infra` (player store, API), `@ui` (layout shell), `@catalog` (track type)

### @catalog — Catalog & Search
- **Focus**: Tracks, albums, artists, genres, search, playlists, library, history, recommendations, AI features
- **Files**: `services/api/catalog/*`, `services/api/recommendation/*`, `services/api/playlist/*`, `services/api/library/*`, `services/api/history/*`, `services/api/ai/*`, `composables/catalog/*`, `composables/useDiscover.ts`, `composables/useLibrary.ts`, `composables/useAlbumColors.ts`, `composables/useCollaborativePlaylist.ts`, `pages/app/PageSearch.vue`, `pages/app/PageTrack.vue`, `pages/app/PageAlbum.vue`, `pages/app/PageArtist.vue`, `pages/app/PageHome.vue`, `pages/app/PageDiscover.vue`, `pages/app/PagePlaylists.vue`, `pages/app/PagePlaylistDetail.vue`, `pages/app/PageLibrary.vue`, `pages/app/PageRecentlyPlayed.vue`, `pages/app/PageAIPlaylistGenerator.vue`, `pages/app/PageAIMoodExplorer.vue`, `pages/app/PageRecommendations.vue`, `pages/app/recommendations/*`, `components/music/SearchOverlay.vue`, `components/music/TrackList.vue`, `components/music/TrackRow.vue`, `components/music/AlbumCard.vue`, `components/music/ArtistCard.vue`, `components/music/ArtistHero.vue`, `components/music/HomeCarousel.vue`, `components/music/HomeSection.vue`, `components/music/UserHero.vue`, `components/music/ProfileTabs.vue`
- **Contracts**: Exports `useCatalogSearch()`, `useTrack()`, `useAlbum()`, `useArtist()`, search composable
- **Deps on**: `@infra` (API barrel, router), `@player` (play action)

### @social — Social & Community
- **Focus**: Social features, rooms, clubs, parties, notifications, reactions
- **Files**: `services/api/social/*`, `services/api/notification/*`, `services/api/reactions/*`, `composables/social/*`, `pages/app/PageRoomLive.vue`, `pages/app/PageClubDetail.vue`, `pages/app/PagePartyDetail.vue`, `pages/app/PageNotifications.vue`, `pages/app/PageSocial.vue`, `components/social/*`, `components/music/ActivityItem.vue`, `components/music/ReactionButton.vue`, `components/music/NotificationItem.vue`
- **Contracts**: Exports social API modules, notification composable
- **Deps on**: `@infra` (API, socket), `@ui` (common components)

### @creator — Creator & Gamification
- **Focus**: Creator dashboard, badges, challenges, leaderboard, XP, tips, contributions, subscriptions
- **Files**: `pages/app/PageCreatorDashboard.vue`, `pages/app/PageGamification.vue`, `pages/app/PageContributions.vue`, `pages/app/PageSubscriptions.vue`, `components/creator/*`, `components/gamification/*`, `components/contribution/*`, `components/tips/*`, `services/api/creator/*`, `services/api/gamification/*`, `services/api/contribution/*`, `services/api/tips/*`, `services/api/subscription/*`, `composables/useCreatorDashboard.ts`
- **Contracts**: Exports creator/gamification API modules
- **Deps on**: `@infra` (API), `@ui` (common components, layout)

### @admin — Admin & Moderation
- **Focus**: Admin panel, user management, track/album/artist CRUD, media upload, moderation
- **Files**: `pages/admin/*`, `components/admin/*`, `composables/admin/*`, `composables/media/*`, `services/api/moderation/*`, `services/api/media/*`, `services/api/users/*`, `router/routes/admin.ts`, `layouts/LayoutAdmin.vue`, `components/layouts/admin/*`
- **Contracts**: Exports admin composables, admin API modules
- **Deps on**: `@infra` (API, router), `@ui` (common components)

### @ui — UI/UX & Layout
- **Focus**: App shell, common components, error pages, styling, utils, global stores
- **Files**: `App.vue`, `main.ts`, `layouts/*` (except `LayoutAdmin.vue`), `pages/errors/*`, `components/common/*`, `components/layouts/*` (except admin), `assets/*`, `utils/*`, `composables/useSwipe.ts`, `composables/useUserProfile.ts`, `stores/feature-flags.ts`, `stores/maintenance.ts`, `stores/page-loader.ts`, `plugins/index.ts`, `plugins/query-builder/*`
- **Contracts**: Exports `AppLoader`, `AppPageContainer`, `AppLogo`, layout components
- **Deps on**: `@infra` (router, stores)

### @infra — Infrastructure
- **Focus**: API barrel, router config, stores (except auth/player), services, plugins, composables root
- **Files**: `services/api/index.ts`, `services/api/common/*`, `services/api/feature-flags/*`, `services/storage/*`, `services/socket/*`, `router/index.ts` (route config), `router/types.ts`, `router/routes/index.ts`, `router/routes/app.ts`, `router/routes/auth.ts`, `stores/index.ts`, `plugins/index.ts`, `plugins/client/*`, `plugins/query-builder/*`, `composables/index.ts`, `composables/useRequest.ts`, `composables/useLoading.ts`, `composables/useFeatureFlags.ts`, `composables/useMaintenance.ts`, `tests/setup.ts`
- **Contracts**: Exports configured API modules, router instance, common stores
- **Deps on**: None (base layer)

## Shared Contract

Each agent must follow these interface contracts when depending on another agent:

| Export | Provider | Consumer |
|--------|----------|----------|
| `PlaybackTrack` type | `@catalog` | `@player` |
| `Track`/`Album`/`Artist` types | `@catalog` | `@admin`, `@social`, `@creator` |
| `useUserAuthStore` + `requireAuth` | `@auth` | All agents |
| `usePlayerStore`, `usePlayer` | `@player` | `@catalog`, `@ui` |
| `audioEngine.play(track)` | `@player` | `@catalog` (play buttons) |
| `useToast()` | `@ui` | All agents |
| `AppPageContainer`, layout slots | `@ui` | All agents |
| API modules (`useXxxApi()`) | `@infra` | All agents |
| Router instance | `@infra` | All agents |

## Workflow

1. **Coordinator** receives a task and decomposes it into agent-specific subtasks
2. Agents run in parallel on their file sets
3. Each agent returns a diff summary and any new/changed contracts
4. Coordinator resolves cross-agent conflicts and merges

### Rules

- No agent modifies files outside its `Files:` list without coordinator approval
- Contract changes (types, exports, composable signatures) must be flagged for coordinator review
- Each agent task should target 1-3 files max; larger tasks are decomposed
- Performance-sensitive code (canvas, animation, heavy computed chains) must follow `PERFORMANCE.md`
- All component exports go through their module's `index.ts` barrel
