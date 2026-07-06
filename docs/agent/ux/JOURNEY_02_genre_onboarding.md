# Journey 02: Genre Onboarding

## When Does It Trigger?

| Trigger | Condition | Result |
|---------|-----------|--------|
| **After registration** | `useAuth().register()` calls `router.push({ name: 'onboarding.genres' })` | ✅ Always redirects to `/onboarding/genres` after successful signup |
| **Manual navigation** | User types `/onboarding/genres` | ✅ Route exists, requires `requiresAuth: true` |
| **Login from existing user** | User logs in (not new) | ❌ Never redirects — always goes to home |
| **Skipped onboarding** | User clicked "بعداً انجام می‌دم" (Do it later) | ✅ Redirects to home — **no way to return** from settings or profile |
| **From settings** | User wants to change genre preferences | ❌ **No option exists** in PageUserSettings — no genre/preference tab for re-running onboarding |

## Route & Guard Details

```typescript
// router/routes/index.ts
{
  path: '/onboarding',
  meta: { layout: 'layout-empty' },  // No sidebar, no player
  children: [
    {
      path: 'genres',
      name: 'onboarding.genres',
      component: () => import('@/pages/onboarding/PageGenreOnboarding.vue'),
      meta: { title: 'Genre Onboarding', requiresAuth: true },
    },
  ],
},
```

- Uses `layout-empty` — no sidebar, no player bar, no app chrome
- Requires authentication (guard redirects to login if not logged in)
- After register, user is auto-logged in → guard passes → onboarding renders

## Page Design (`PageGenreOnboarding.vue`)

```
┌──────────────────────────────────────┐
│                                      │
│       🎵 (music icon)               │
│   "چه سبک موسیقی رو دوست داری؟"       │  ← Persian: "What music genre do you like?"
│                                      │
│   حداقل ۳ تا انتخاب کن تا بتونیم      │  ← "Select at least 3 so we can give
│   پیشنهادهای بهتری برات داشته باشیم    │     you better suggestions"
│                                      │
│  ┌──────┐ ┌──────┐ ┌──────┐ ...    │
│  │ Pop  │ │ Rock │ │ Jazz │         │  ← Genre chips (toggle on click)
│  └──────┘ └──────┘ └──────┘         │
│                                      │
│  [          ادامه          ]         │  ← "Continue" button (disabled if <3)
│                                      │
│         بعداً انجام می‌دم             │  ← "Do it later" link
│                                      │
└──────────────────────────────────────┘
```

**Fully Persian UI** — this is the only page in the user-facing app that is entirely in Persian.

## Step-by-Step Flow

### Step 1: Page Load
| Aspect | Detail |
|--------|--------|
| Loading state | ✅ Pulse skeleton chips (12 animated placeholders) |
| API call | `genresApi.getGenres()` — `GET /api/v1/catalog/genres` |
| Error state | ✅ Persian error: "بارگیری سبک‌ها با مشکل مواجه شد. بعداً تلاش کن." |
| Empty state | ✅ "در حال حاضر سبکی برای نمایش وجود ندارد" |
| Files | `PageGenreOnboarding.vue:103-112` |

### Step 2: Genre Selection
| Aspect | Detail |
|--------|--------|
| Interaction | Tap/click genre chip — toggles selection (green highlight + glow) |
| Min selection | "ادامه" button disabled until ≥3 genres selected |
| Visual feedback | Selected: green border + green bg + shadow glow. Unselected: subtle white border + light bg |
| Animation | Spring transition on chips (`transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)`) |

### Step 3: Save & Redirect
| Aspect | Detail |
|--------|--------|
| Save action | `POST /onboarding/genres` with `{ genre_ids: ["id1", "id2", "id3"] }` |
| Loading state | Button shows spinner, disabled |
| Error state | "ذخیره‌سازی با مشکل مواجه شد. دوباره تلاش کن." (Persian) |
| Success | `router.replace({ name: 'app.home' })` — replaces, not pushes |
| Skip | `router.replace({ name: 'app.home' })` — same destination |

### Step 4: What Happens to Selections?

The genres are sent to `POST /onboarding/genres`. This presumably stores them on the user profile to seed recommendations. The API response is not used on the frontend — the page simply navigates away.

**No confirmation toast** is shown on save or skip.

## State Matrix Findings

| State | Present? | Notes |
|-------|----------|-------|
| 🟢 Loading (skeleton) | ✅ | 12 animated skeleton chips |
| 🟢 Empty genres (API returns none) | ✅ | "در حال حاضر سبکی برای نمایش وجود ندارد" |
| 🟢 Error loading genres | ✅ | Persian error message |
| 🟢 Selection in progress | ✅ | Toggle chips with visual feedback |
| 🟢 Min-selection gate | ✅ | Button disabled if <3 selected |
| 🟢 Saving | ✅ | Button spinner, disabled |
| 🟢 Save error | ✅ | Persian error message |
| 🟢 Skip | ✅ | "بعداً انجام می‌دم" link |
| 🔴 Success feedback (save) | ❌ | No toast, no animation — abrupt redirect |
| 🔴 Success feedback (skip) | ❌ | Same — immediate redirect to home |
| 🔴 Re-run from settings | ❌ | **No way to restart onboarding** — no settings option, no profile link |
| 🔴 API contract (uses `client.post` directly) | ⚠️ | Uses raw `client.post('/onboarding/genres', ...)` instead of typed API service — no schema validation |

## Impact of Skipping

If the user skips genre onboarding:
- **`/users/me/preferences`**: No genre preference is stored for this user
- **Home feed recommendations**: `getPersonalized()` still fires but has no genre signal to work with — returns generic content or empty results
- **BecauseOf sections**: Sections like `because_of_[genre]` won't appear (no genres to base them on)
- **Album/artist discovery**: Still works via global catalog — just not personalized
- **No persistent reminder**: No banner, toast, or notification to nudge user back to onboarding

## Friction Points

| # | Severity | Location | Problem | User Impact |
|---|----------|----------|---------|-------------|
| F-201 | ⚠️ MAJOR | `PageGenreOnboarding.vue` | No way to re-run onboarding from Settings/Profile | Users who skipped or want to update preferences are stuck |
| F-202 | ⚠️ MAJOR | `PageGenreOnboarding.vue:129-130` | No success feedback after save — immediate redirect | User may wonder if it saved |
| F-203 | 💡 IMPROVE | `PageGenreOnboarding.vue:138-139` | Skip also has no feedback | Same as above |
| F-204 | 💡 IMPROVE | `PageGenreOnboarding.vue` | Uses raw `client.post()` — not a typed service | No API contract validation, inconsistent with other API calls |

## RTL / A11y / Mobile Notes

- ✅ **Fully Persian UI** — this is the only fully localized page in the app
- ✅ **Loading skeleton** present for initial data fetch
- ✅ **Min-selection enforced** — prevents premature submission
- ❌ **No skip-link** — LayoutEmpty doesn't have one
- ❌ **Focus management** — After saving/skipping, focus is lost (redirect)
- ✅ **Touch targets** — Genre chips have adequate tap areas (px-5 py-2.5)
- ℹ️ **Layout** — Uses `layout-empty` (no distracting chrome, clean focus on task)

## Delight Opportunities

- ✨ **Genre visual preview**: Show a sample track or artist for each genre when hovered
- ✨ **Smart defaults**: Pre-select popular Persian genres (Pop, Traditional, Classic) based on locale
- ✨ **Animation on complete**: Confetti or celebration animation when minimum 3 selected and saved
- ✨ **Add "Edit genres" to settings**: Profile/Preferences tab could include a "Music preferences" section to re-run genre selection

## Open Questions

1. Where exactly are the selected genres stored? `POST /onboarding/genres` — is this a user_preferences table or the user_genres join table?
2. Is there a backend endpoint to GET the user's current genre selections? (`GET /users/me/genres`?)
3. Should there be a persisted "nudge" for users who skipped onboarding? (e.g., show on home page: "Pick your genres for better recommendations")
4. Can the user select ALL genres or is there a max? Current code only enforces min (3) — no max.
