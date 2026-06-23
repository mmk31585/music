# Sprint 01 — Social Hub + Profile Redesign

> Implementation sprint for the two major UX redesigns
> Design specs: `design-system/pages/social-hub.md`, `design-system/pages/user-profile.md`

---

## Team

| Agent | Role | Owns |
|-------|------|------|
| **@muse-social** | Social Hub | `PageSocial.vue`, `components/social/*`, `components/music/ActivityItem.vue` |
| **@muse-catalog** | User Profile | `PageUserProfile.vue`, `UserHero.vue`, `ProfileTabs.vue`, `StatCard.vue` |
| **@muse-infra** | Infrastructure | API barrel, stores, composables |

## Shared Contracts

| Contract | Provider | Consumer |
|----------|----------|----------|
| `ActivityItem` component | `@muse-social` | `@muse-catalog` (profile feed) |
| `IconName` convention | Both | Both — use Lucide-style names |
| API types (Party, Room, Club) | `@muse-infra` → `@muse-social` | Both |
| User profile types | `@muse-infra` → `@muse-catalog` | Both |

---

## Task List

### @muse-social — Social Hub (social-hub.md)

| # | File | Action | Depends On |
|---|------|--------|-----------|
| 1 | `PageSocial.vue` | Full rewrite: bento grid, activity river, inline party CTA, category discussion pills, 3-step wizard | — |
| 2 | `ListeningPartyCard.vue` | Rewrite: avatar grid, cover thumbnail, share button, SVG icons | — |
| 3 | `LiveRoomCard.vue` | Rewrite: listener avatars, live badge, SVG icons | — |
| 4 | `MusicClubCard.vue` | Rewrite: member bar, avatar row, SVG icons (🏛→Building icon) | — |
| 5 | `DiscussionThread.vue` | Rewrite as `DiscussionCard.vue` with category pill filtering | — |
| 6 | `CreateClubDialog.vue` → `CreateWizard.vue` | 3-step wizard (Name → Details → Review) | — |
| 7 | `ActivityItem.vue` | Enhance: contextual action buttons, timeAgo, avatar | — |
| 8 | `components/social/index.ts` | Update barrel exports | 1-7 |

### @muse-catalog — User Profile (user-profile.md)

| # | File | Action | Depends On |
|---|------|--------|-----------|
| 1 | `PageUserProfile.vue` | Full rewrite: stats gallery, mosaic/list toggle, avatar wall, rich empty states | — |
| 2 | `UserHero.vue` | Rewrite: cinematic aurora gradient, prominent "Now Playing", equalizer animation, creator badge | — |
| 3 | `ProfileTabs.vue` | Rewrite: pill-style with icons, `role="tablist"`, `aria-selected` | — |
| 4 | `StatCard.vue` (new) | Generic stat card component: icon, value, label, optional trend | — |
| 5 | `MusicStatusWidget.vue` | Simplify or merge into UserHero (now prominent) | 2 |
| 6 | `components/music/profile/index.ts` | Update barrel exports | 1-5 |

### @muse-infra — Infrastructure

| # | File | Action | Depends On |
|---|------|--------|-----------|
| 1 | `services/api/social/index.ts` | Ensure Party/Room/Club API methods match spec needs | — |
| 2 | `services/api/reactions/index.ts` | Ensure liked tracks/albums API matches profile data flow | — |
| 3 | Router | Verify no new routes needed (social sub-routes already exist) | — |

---

## Progress

```
[ ] @muse-social: PageSocial.vue rewrite
[ ] @muse-social: All card component rewrites
[ ] @muse-social: ActivityItem.vue enhance
[ ] @muse-social: CreateWizard.vue
[ ] @muse-catalog: PageUserProfile.vue rewrite
[ ] @muse-catalog: UserHero.vue rewrite
[ ] @muse-catalog: StatCard.vue + tabs rewrite
[ ] @muse-infra: API/barrel audit
[ ] VERIFY: integration test
```
