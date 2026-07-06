# Journey 26: Admin — Reporting & Analytics

> Full trace: analytics event pipeline → dashboard stats → metrics → data gaps.

---

## Analytics Architecture

### Event Pipeline

```
User action (play, pause, skip, like, etc.)
  → POST /api/v1/analytics/event { event_type, properties, timestamp }
  → Analytics module stores event
  → Worker runs daily rollup (aggregates into summary tables)
  → Queries available via analytics API
```

### Events Tracked

```go
// internal/modules/analytics/
// Events: play, pause, skip, seek, like, unlike, share, follow, unfollow,
//         playlist_create, playlist_add, playlist_remove, search, login, register
```

### Aggregation

```go
// Worker: analytics worker runs every 1 minute
// DailyRollup: aggregates events into daily summary tables
// Types: plays_by_track, plays_by_artist, plays_by_hour, user_activity
```

### Current Admin Dashboard Stats

```go
// GET /admin/dashboard/stats → { total_tracks, total_artists, total_albums, total_genres }
// From internal/modules/dashboard/handler.go

// NOTE: Only 4 basic counts. No trend data, no growth rates, no charts.
// The analytics module has richer data but no admin frontend to query it.
```

**Key Gap**: The analytics module ingests a wealth of event data and runs daily rollups, but there is **no admin-facing analytics UI** to query or visualize this data. The dashboard shows only 4 static counts.

---

## What an Admin Analytics Dashboard Should Include

Based on the data available in the analytics module, these queries are possible but have no UI:

| Metric | Data Source | Currently Visible? |
|--------|------------|-------------------|
| Total Users (current) | User table | ❌ No |
| Total Tracks (current) | Track table | ✅ Yes (dashboard) |
| Total Artists (current) | Artist table | ✅ Yes (dashboard) |
| Total Albums (current) | Album table | ✅ Yes (dashboard) |
| Total Plays (all time) | Analytics events | ❌ No |
| Plays Today / This Week | Analytics events | ❌ No |
| Active Users (daily/weekly/monthly) | Analytics events | ❌ No |
| Top Tracks (by plays, period) | Analytics events + rollup | ❌ No |
| Top Artists (by plays, period) | Analytics events + rollup | ❌ No |
| Growth Rate (users, tracks, plays) | Historical comparison | ❌ No |
| Retention / Engagement | Analytics events | ❌ No |
| Search Queries (popular, failed) | Search analytics events | ❌ No |
| Geographic Distribution | User locations | ❌ No |
| Device / Platform Breakdown | User agent data | ❌ No |
| Subscription Conversion Rate | Subscription data | ❌ No |

---

## Metrics Endpoint

```go
// GET /api/v1/metrics — Prometheus metrics
// Exposed but no admin UI for it
// Counters: HTTP requests, worker jobs, DB operations, storage operations
// Histograms: request duration, worker job duration
```

The metrics endpoint exists for Prometheus scraping but has **no admin UI** to display it. An admin would need to access Prometheus/Grafana separately.

---

## Existing Data Points That Could Feed Analytics

| Admin Page | Visible Metrics | Missing Context |
|-----------|----------------|-----------------|
| Dashboard | Total tracks, artists, albums, genres | No trends, no growth, no user counts |
| Moderation Stats | Total reports, pending, resolved today, flagged, unique reporters, avg resolution time | No trend over time |
| Creator Dashboard (per-creator) | Plays, listeners, followers, revenue, daily plays chart | No aggregate across all creators |
| Subscription Admin | Plans, subscriber counts | No MRR, churn, conversion rates |

---

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-2601 | ⚠️ MAJOR | `internal/modules/analytics/` | **Analytics module has rich event data but ZERO admin UI to query it** | Admins can't see user growth, top content, engagement trends | Build admin analytics dashboard with time-series charts |
| F-2602 | ⚠️ MAJOR | `internal/modules/dashboard/` | **Dashboard stats endpoint returns only 4 counts** — no users, no plays, no trends | Admin dashboard is nearly useless for monitoring platform health | Expand dashboard endpoint with richer stats (users, plays, DAU/WAU/MAU) |
| F-2603 | ⚠️ MAJOR | All admin pages | **No CSV/JSON export** for any data (users, tracks, reports) | Can't do offline analysis | Add "Export" button to every list page |
| F-2604 | 💡 IMPROVE | All admin pages | **No date range filter** on any admin report | Can't compare periods | Add date range picker to all data views |
| F-2605 | 💡 IMPROVE | Admin dashboard | **No top content lists** (top tracks, top artists, top albums) | Can't see what's popular | Add "Top 10" sections to dashboard |
| F-2606 | 💡 IMPROVE | Admin analytics | **No growth rate indicators** — no % change week-over-week or month-over-month | Can't spot trends | Add Δ% badges to all metric cards |
| F-2607 | 💡 IMPROVE | `internal/modules/analytics/` | **No custom query endpoint** — can't filter analytics by date range or event type | Can't answer ad-hoc questions | Add query endpoint with type + date filters |
| F-2608 | 💡 IMPROVE | Admin analytics | **No search query analytics** — can't see what users search for | Missed opportunities for content gaps | Add search term analytics page |
| F-2609 | 💡 IMPROVE | Admin analytics | **No retention/cohort analysis** | Can't measure user retention | Add cohort tables (day 1, 7, 30 retention) |
| F-2610 | 💡 IMPROVE | Admin analytics | **No geographic analytics** — can't see where users are | Can't target regional content | Add geo-map or country list |
| F-2611 | 💡 IMPROVE | Admin analytics | **No exportable report scheduling** — can't get weekly PDF/CSV emailed | Manual reporting only | Add scheduled report delivery |

## RTL / A11y / Mobile Notes

- ❌ Charts and graphs don't exist yet — accessibility requirements TBD
- ✅ Numeric stat cards in existing admin use proper `aria-label`
- ❌ No admin analytics page has skip link to content
