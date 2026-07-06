# Journey 18: User Profile — Public Profile, Stats, Activity

> Full trace: profile page → tabs → follow/unfollow → gamification display → other user profiles.

---

## Profile Architecture

### Route

```typescript
// Route: /profile → redirects to user.profile
// Route: /user/:id? → PageUserProfile.vue (name: 'user.profile')
// Optional :id — no id shows own profile
// Meta: { title: 'Profile' } — NO requiresAuth (public profiles)
```

The profile route is **public** — any guest can view any user's profile. This enables artist discovery and social browsing without login.

### Data Fetching — `fetchProfile()`

```typescript
async function fetchProfile() {
  isLoading.value = true
  try {
    const [profileRes, reactionsRes, gamificationRes, socialRes] = await Promise.all([
      profileId.value
        ? useUserApi().getPublicUserProfile(profileId.value)
        : Promise.resolve(auth.user),
      Promise.all([
        profileId.value
          ? reactionsApi.getLikedTracks({ user_id: profileId.value })
          : reactionsApi.getLikedTracks(),
        profileId.value
          ? reactionsApi.getLikedAlbums({ user_id: profileId.value })
          : reactionsApi.getLikedAlbums(),
      ]).catch(() => [{ tracks: [], total: 0 }, { albums: [], total: 0 }]),
      // Gamification data only for own profile
      !profileId.value
        ? gamificationApi.getProfile().catch(() => null)
        : Promise.resolve(null),
      profileId.value
        ? socialApi.isFollowing(profileId.value)
        : Promise.resolve(null),
    ])
    // ... assign to reactive refs
  } catch (e) {
    error.value = e
  } finally {
    isLoading.value = false
  }
}
```

Key detail: **Gamification data (XP, badges, level) only loads for the current user's own profile**. Other users' profiles show a "Now Playing" section instead.

---

## Profile Layout

```
PageUserProfile.vue (670 lines)
  ├── UserHero.vue (272 lines)
  │     ├── Aurora gradient background (genreColor-driven)
  │     ├── Avatar (with "now playing" animated ping ring for other users)
  │     ├── Name, handle (@username), bio
  │     ├── Follower/following counts + join date
  │     ├── "Now Playing" section (other users): eq bars, cover, title, artist, "Listen" button
  │     ├── Action buttons: Follow/Following toggle, Share, Edit Profile (own)
  │     └── Full keyboard focus + ARIA labels
  │
  ├── ProfileTabs.vue (35 lines)
  │     └── Tab bar: Tracks, Albums, Edits, Followers, Following
  │
  └── Tab content (conditional on active tab)
        ├── Tracks tab: Mosaic grid OR list view (toggleable)
        │     └── Public/private toggle for own profile tracks
        ├── Albums tab: Grid of liked albums
        ├── Edits tab: User's video edits (from videoApi)
        ├── Followers tab: Avatar grid + "View all" link
        └── Following tab: Same grid pattern
```

### UserHero Component (`UserHero.vue`)

| Section | Source | Behavior |
|---------|--------|----------|
| Aurora background | `genreColor` from genre preferences | Reacts to dominant genre, gradient animation |
| Avatar | `user.avatar_url` | Circular, fallback to initials |
| Display name | `user.full_name` | Public field |
| Handle | `user.username` | `@username`, clickable |
| Bio | `user.bio` | Expandable if >200 chars |
| Stats row | follower count, following count, join date | Counts link to tabs |
| Now Playing (other users) | `videoApi.getMusicStatus(userId)` | Animated equalizer bars, cover art, title |
| Follow button | `socialApi.isFollowing()` | Toggle Follow/Following, optimistic update |
| Share button | Web Share API (or clipboard fallback) | Share profile URL |
| Edit Profile button | Own profile only | Navigate to `/settings` or inline edit |

---

## Tab System

### Tab State Machine

```typescript
const activeTab = ref<'tracks' | 'albums' | 'edits' | 'followers' | 'following'>('tracks')

// Tabs have lazy loading — content only fetches when tab is active
watch(activeTab, (tab) => {
  if (tab === 'followers' && !followers.value.length) fetchFollowers()
  if (tab === 'following' && !following.value.length) fetchFollowing()
  if (tab === 'edits' && !edits.value.length) fetchEdits()
})
```

### Tab Loading States

| Tab | Loading | Empty | Error |
|-----|---------|-------|-------|
| Tracks | ✅ Skeleton grid | ✅ Heart icon + "No liked tracks yet" (own) / "Tracks are private" (other) | ❌ Silent |
| Albums | ✅ Skeleton grid | ✅ Disc icon + "No liked albums yet" | ❌ Silent |
| Edits | ✅ Skeleton grid | ✅ Video icon + "No edits yet" (own) / "No edits from this user" (other) | ❌ Silent |
| Followers | ✅ Skeleton grid | ✅ Users icon + "No followers yet" + share prompt | ❌ Silent |
| Following | ✅ Skeleton grid | ✅ User-plus icon + "Not following anyone yet" | ❌ Silent |

---

## Follow/Unfollow Flow

```typescript
async function toggleFollow() {
  if (isFollowing.value) {
    await socialApi.unfollow(profileId.value)
    isFollowing.value = false
    followerCount.value--
  } else {
    await socialApi.follow(profileId.value)
    isFollowing.value = true
    followerCount.value++
  }
}
```

- Optimistic update — UI changes immediately
- If API fails, state reverts (rollback)
- No confirmation dialog on unfollow (instant action)
- No toast feedback on follow/unfollow

### Follow State Visual

| State | Button Text | Button Style |
|-------|-------------|--------------|
| Not following | "Follow" | Primary gradient |
| Following | "Following" | Secondary/outline |
| Own profile | "Edit Profile" | Secondary/outline |
| Pending (API) | Spinner overlay | Disabled |
| Error | Reverts to previous state | Previous style |

---

## Other User Profile (Now Playing)

When viewing another user's profile who is currently listening to music:

```
GET /users/:userId/music-status → { is_listening, track, album, artist, album_art_url, timestamp }

UserHero displays:
  → Animated equalizer bars (if listening)
  → Album art thumbnail
  → Track title + artist name
  → "Listen" button → plays the track in user's own player
  → If not listening: "Not listening right now" with muted indicator
```

---

## Gamification Display (Own Profile Only)

```
Only shown on own profile — NOT on other user's profiles.

Display:
  → Level badge below avatar (e.g., "Level 12")
  → XP bar (GamificationXPBar component)
  → Mini stats: Tracks liked, Followers, Following
  → (Full gamification at /gamification route)
```

---

## Share Profile

```typescript
async function shareProfile() {
  const url = `${window.location.origin}/user/${user.value.id}`
  if (navigator.share) {
    await navigator.share({ title: user.value.full_name, url })
  } else {
    await navigator.clipboard.writeText(url)
    // Toast: "Profile link copied!"
  }
}
```

Uses Web Share API on supported browsers, clipboard fallback on others.

---

## State Matrix & Friction

| State | Own Profile | Other User Profile |
|-------|------------|-------------------|
| 🟢 Loading | ✅ Full skeleton hero + tab skeleton | ✅ Same |
| 🟢 Loaded (has data) | ✅ Full profile with all tabs | ✅ Profile minus gamification |
| 🟢 Loaded (no data) | ✅ Empty states per tab | ✅ Same |
| 🟢 Error loading | ✅ "Failed to load profile" + "Try again" | ✅ Same |
| 🟢 Follow error | ❌ Silent rollback, no toast | ❌ Same |
| 🔴 User not found | ❌ API 404 not differentiated from other errors | ❌ Shows generic error |
| 🔴 Avatar load fail | ✅ Fallback to initials | ✅ Same |
| 🔴 Very long bio | ✅ Expandable (>200 chars) | ✅ Same |
| 🔴 Guest viewing profile | ✅ Works (public profiles) | ✅ N/A |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1801 | ⚠️ MAJOR | `PageUserProfile.vue` | Profile fetch uses `Promise.all` — if one API fails, the entire profile fails | Partial profile data is better than none | Use `Promise.allSettled` per section |
| F-1802 | ⚠️ MAJOR | `PageUserProfile.vue` | **No edit profile modal** — "Edit Profile" navigates to `/settings` | User leaves the profile context to make changes | Inline edit or modal for basic fields |
| F-1803 | ⚠️ MAJOR | `PageUserProfile.vue` | Follow/unfollow has **no toast/feedback** on success or failure | User may not know if action registered | Add success toast, error toast on failure |
| F-1804 | 💡 IMPROVE | `PageUserProfile.vue` | **No profile header cover image** — only avatar + gradient | Profile feels sparse compared to competitors | Add optional cover image upload |
| F-1805 | 💡 IMPROVE | `PageUserProfile.vue` | **No "tracks in common"** section for other users | No shared-taste discovery | Add mutual liked tracks section |
| F-1806 | 💡 IMPROVE | `UserHero.vue` | Bio character limit **not shown** to user when editing | User may exceed limit and get error silently | Show "X/200" counter |
| F-1807 | 💡 IMPROVE | `PageUserProfile.vue` | **No pagination** on followers/following tabs | Large follow lists (500+) unfetchable | Add API pagination or virtual scroll |
| F-1808 | 💡 IMPROVE | `PageUserProfile.vue` | **No "sort by"** on tracks tab — always by date added | Can't find tracks alphabetically | Add sort options |
| F-1809 | 💡 IMPROVE | `PageUserProfile.vue` | **No "View as guest"** toggle for own profile | Can't preview what others see | Add preview mode toggle |
| F-1810 | 💡 IMPROVE | `UserHero.vue` | Now Playing section for other users **not updated in real-time** | Stale "now playing" after user changes track | Poll every 30s or use WebSocket push |
| F-1811 | 💡 IMPROVE | `PageUserProfile.vue` | **No activity feed** tab — when user joined, what they listened to (public) | No way to see user's recent activity | Add "Activity" tab with recent listens |
| F-1812 | 💡 IMPROVE | `PageUserProfile.vue` | **No "Report user"** action | Users can't report inappropriate profiles | Add report button with reason picker |
| F-1813 | 💡 IMPROVE | `PageUserProfile.vue` | **No "Block user"** action | Users can't block unwanted interactions | Add block button with confirmation |

## RTL / A11y / Mobile Notes

- ✅ Follow button has proper `aria-label` — "Follow [username]" / "Unfollow [username]"
- ✅ Tab system uses `aria-selected`, `aria-controls`, `tabindex`, `focus-visible`
- ❌ Tracks tab uses a mosaic grid that doesn't preserve focus order well in RTL
- ✅ Avatar has `alt` text with user's name
- ❌ Edit profile button on mobile may be hidden behind overflow if screen is too small
- ✅ "Now Playing" section uses `aria-live="polite"` for track changes
