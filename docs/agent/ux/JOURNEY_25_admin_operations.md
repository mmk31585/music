# Journey 25: Admin — Dashboard & Operations

> Full trace: admin dashboard → user management → subscriptions → contributions → video management → reporting gaps.

---

## Admin Dashboard (`PageAdminDashboard.vue`)

### Route

```typescript
// Route: /admin → admin.dashboard
// Meta: { requiresAuth: true, requiresRole: 'admin' }
```

### Layout

```
AdminDashboard.vue
  ├── AdminSectionHeader ("Dashboard" / "داشبورد")
  ├── Stats row:
  │     ├── Total Users
  │     ├── Total Tracks
  │     ├── Total Artists
  │     ├── Total Albums
  │     └── Total Plays
  ├── Catalog Quick Actions (CatalogQuickActions.vue):
  │     ├── Tracks → /admin/tracks
  │     ├── Artists → /admin/artists
  │     ├── Albums → /admin/albums
  │     └── Genres → /admin/genres
  └── Recent activity (placeholder or limited)
```

### Data Source

```typescript
// GET /admin/dashboard/stats
// Returns: { total_users, total_tracks, total_artists, total_albums, total_plays }
```

### Concerns (From Inline Comment)

```
// TODO: Fetching ALL catalog items just for preview counts is wasteful.
// Use a dedicated stats endpoint instead.
```

The dashboard previously fetched entire track/artist/album lists to compute counts. A dedicated stats endpoint exists now but the comment suggests the issue may not be fully resolved.

---

## User Management (`PageAdminUsers.vue` — 724 lines)

### User List

```
Searchable user table with:
  ├── Avatar + username
  ├── Email
  ├── Role: admin / user / moderator
  ├── Status: active / suspended / banned
  ├── Verified badge
  ├── Created date
  ├── Last login
  └── Actions: Edit, Suspend, Delete

Filters:
  ├── Role filter (all / admin / user / moderator)
  ├── Status filter (all / active / suspended / banned)
  └── Sortable columns

Pagination: server-side, offset/limit
```

### User Detail Dialog

```
Dialog opens on "View" or inline:
  ├── Full user info
  ├── Edit role dropdown
  ├── Toggle active/suspended/verified
  ├── Delete user (with confirmation)
  └── Activity log (limited)
```

### User Actions

| Action | UX | Confirmation |
|--------|----|-------------|
| Edit role | Dropdown → save | None |
| Suspend | Toggle → "Suspend user?" dialog | Yes — "They won't be able to log in" |
| Ban | Toggle → "Ban user?" dialog | Yes — "This is permanent and irreversible" |
| Verify | Toggle → instant | None |
| Delete | Button → "Delete user?" dialog | Yes — "This will permanently delete all user data" |

---

## Subscription Management (`PageAdminSubscriptions.vue` — 150 lines)

### Route

```typescript
// Route: /admin/subscriptions
```

### Layout

```
PageAdminSubscriptions.vue
  ├── AdminSectionHeader ("Subscriptions" / "اشتراک‌ها")
  ├── Plan cards (current subscription plans):
  │     ├── Plan name
  │     ├── Price
  │     ├── Active subscribers count
  │     └── Status (active / inactive)
  └── User subscriptions list:
        ├── User
        ├── Plan
        ├── Status (active / trial / canceled / expired)
        ├── Next billing date
        └── Actions: Cancel, Refund
```

### State

| State | Behavior |
|-------|----------|
| 🟢 Loading | ✅ Skeleton cards + list |
| 🟢 Loaded (has data) | ✅ Plan cards + subscription list |
| 🟢 Loaded (no subscriptions) | ✅ "No subscriptions yet" |
| 🟢 Cancel subscription | ✅ Confirmation dialog |
| 🟢 Error | ❌ Silent — empty list |
| 🔴 Refund flow | ❌ No refund UI — manual process only |

---

## Contributions Management (`PageAdminContributions.vue` — 615 lines)

### Route

```typescript
// Route: /admin/contributions
```

### Layout

```
PageAdminContributions.vue
  ├── AdminSectionHeader ("Contributions" / "مشارکت‌ها")
  ├── Stats row: Total, Pending review, Approved, Rejected, Total XP awarded
  ├── Tabs:
  │     ├── Pending (awaiting review)
  │     │     ├── Contribution cards: user, type, submitted date, status
  │     │     ├── Inline actions: Approve (with XP award), Reject (with reason)
  │     │     └── Review detail: expandable view of the submitted content
  │     ├── Approved
  │     │     └── List of approved contributions with XP awarded
  │     └── Rejected
  │           └── List of rejected contributions with rejection reason
  └── Filters: contribution type, date range, user
```

### Contribution Review Flow

```
Admin reviews pending contribution:
  → Views submitted content (lyrics, translation, metadata)
  → Compares with original
  → Approve: optional XP award amount → contribution accepted, XP credited
  → Reject: required rejection reason → contribution rejected, user notified
```

---

## Video Management (`PageAdminVideos.vue` — 714 lines)

### Route

```typescript
// Route: /admin/videos
```

### Layout

```
PageAdminVideos.vue
  ├── AdminSectionHeader ("Videos" / "ویدیوها")
  ├── Search: by title, artist
  ├── Status filter: all / pending / approved / rejected
  ├── Video cards:
  │     ├── Thumbnail
  │     ├── Title
  │     ├── Artist
  │     ├── Duration
  │     ├── Status badge
  │     └── Actions: Approve, Reject, Delete, Edit
  └── Pagination
```

---

## Missing Admin Features

| Feature | Status | Note |
|---------|--------|------|
| **Playlist curation** | ❌ Missing | No admin UI for featured/curated playlists |
| **Recommendation management** | ❌ Missing | No way to boost/pin tracks or influence recommendations |
| **Analytics dashboard** | ⚠️ Limited | Only basic counts — no charts, trends, or user growth |
| **Activity log** | ⚠️ Limited | Moderation has audit trail, but no global admin activity log |
| **Bulk operations** | ⚠️ Partial | Bulk moderation exists, but no bulk catalog operations (bulk tag, bulk genre) |
| **System health** | ❌ Missing | No server health, queue depth, or error rate display |
| **Content migration** | ❌ Missing | No tool to move content between storage providers |
| **Audit trails (catalog)** | ❌ Missing | No "who changed what" for tracks/albums/artists |

---

## Admin Sidebar Navigation

The sidebar has 14 navigation items in 4 sections:

```
Overview
  ├── Dashboard (/admin)

Catalog
  ├── Catalog (/admin/catalog)
  ├── Tracks (/admin/tracks)
  ├── Artists (/admin/artists)
  ├── Albums (/admin/albums)
  └── Genres (/admin/genres)

Management
  ├── Users (/admin/users)
  ├── Media (/admin/media)
  ├── Videos (/admin/videos)
  ├── Video Upload (/admin/video-upload)
  ├── Import from Internet (/admin/import)
  ├── Import by Artist (/admin/import/artist)
  ├── Music Ingestion (/admin/ingestion)
  └── Moderation (/admin/moderation)

Finance
  ├── Subscriptions (/admin/subscriptions)
  └── Contributions (/admin/contributions)
```

---

## State Matrix

| State | Dashboard | Users | Subscriptions | Contributions | Videos |
|-------|-----------|-------|---------------|---------------|--------|
| 🟢 Loading | ✅ Skeleton stats | ✅ Table skeleton | ✅ Skeleton | ✅ Skeleton cards | ✅ Card skeleton |
| 🟢 Loaded (has data) | ✅ Stats + quick actions | ✅ Paginated table | ✅ Plans + list | ✅ Contribution tabs | ✅ Video cards |
| 🟢 Loaded (no data) | ✅ Zero stats | ✅ "No users found" | ✅ "No subscriptions" | ✅ "No contributions" | ✅ "No videos" |
| 🟢 Error loading | ❌ Silent zeros | ❌ Silent empty table | ❌ Silent | ❌ Silent | ❌ Silent |
| 🟢 Search/filter | ✅ N/A | ✅ Working | ❌ No search | ⚠️ Basic filters | ✅ Working |
| 🟢 Pagination | N/A | ✅ Server-side | ❌ No pagination | ✅ Paginated | ✅ Paginated |
| 🔴 Missing playlist curation | ❌ No admin UI | N/A | N/A | N/A | N/A |
| 🔴 No analytics trends | ❌ Only current counts | N/A | N/A | N/A | N/A |

## Friction Points

| # | Severity | Journey | Location | Problem | User Impact | Fix |
|---|----------|---------|----------|---------|-------------|-----|
| F-2501 | ⚠️ MAJOR | J25 | All admin pages | **No admin notifications** — no bell, no real-time alerts for moderation/contributions | Admins must manually check each section | Add admin notification system with role-based alerts |
| F-2502 | ⚠️ MAJOR | J25 | Admin dashboard | **Dashboard shows only current counts** — no trends, charts, or growth | Can't monitor platform health over time | Add time-series charts (users/tracks/plays over 7/30/90 days) |
| F-2503 | ⚠️ MAJOR | J25 | Admin sidebar | **No playlist curation route** — playlists managed only on public side | No way to feature/promote playlists | Add `/admin/playlists` for playlist management and curation |
| F-2504 | 💡 IMPROVE | J25 | `PageAdminUsers.vue` | **No user activity log** — can't see user's recent actions | Hard to investigate problem users | Add "Activity" tab per user showing recent listens, reports, contributions |
| F-2505 | 💡 IMPROVE | J25 | `PageAdminSubscriptions.vue` | **No subscription analytics** — no MRR, churn, conversion rates | Can't measure business health | Add subscription metrics dashboard |
| F-2506 | 💡 IMPROVE | J25 | `PageAdminUsers.vue` | **No impersonation mode** — can't see what user sees | Support can't debug user issues | Add "Log in as user" with audit trail |
| F-2507 | 💡 IMPROVE | J25 | `PageAdminSubscriptions.vue` | **No refund flow** — must process externally | Support friction for payment issues | Add "Issue refund" button with confirmation |
| F-2508 | 💡 IMPROVE | J25 | Admin sidebar | **No search across admin** — must navigate to each section | Slow to find specific content | Add global admin search (Cmd+K for admin palette) |
| F-2509 | 💡 IMPROVE | J25 | All admin pages | **No bulk operations on users** — suspend/ban/verify one at a time | Slow for spam waves | Add multi-select + bulk actions on user list |
| F-2510 | 💡 IMPROVE | J25 | All admin pages | **No "last edited by"** on any catalog item | No accountability for changes | Add audit metadata to all catalog items |
| F-2511 | 💡 IMPROVE | J25 | Admin dashboard | **No system health section** — no server status, queue depth, error rates | Can't monitor infrastructure | Add health panel (API uptime, WS status, queue depth, error rate) |
| F-2512 | 💡 IMPROVE | J25 | `PageAdminContributions.vue` | **No bulk contribution approval** — must approve one at a time | Slow for high-quality bulk submissions | Add "Approve all pending" with XP presets |
| F-2513 | 💡 IMPROVE | J25 | All admin pages | **No export/download** for any list (users, tracks, subscriptions) | Can't do offline analysis | Add CSV export per admin page |
| F-2514 | 💡 IMPROVE | J25 | `PageAdminVideos.vue` | **No bulk video operations** — approve/reject one at a time | Slow for video moderation | Add multi-select + bulk approve/reject |

## RTL / A11y / Mobile Notes

- ✅ Admin sidebar collapses on mobile to hamburger drawer
- ✅ Admin topbar shows route title for screen readers
- ❌ No skip link on admin layout — keyboard users tab through full sidebar
- ✅ Admin stat cards have `aria-label` for numeric values
- ❌ User management table rows have no row-level `aria-label` — screen reader just reads cells
- ✅ All dialogs use `<Teleport>` for proper stacking
- ❌ No keyboard shortcut for "Back to app" — must click sidebar link
