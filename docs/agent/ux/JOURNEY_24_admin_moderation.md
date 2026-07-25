# Journey 24: Admin — Moderation Queue

> Full trace: content reporting → moderation queue → flagging → resolution → stats.

---

## Moderation Architecture

The moderation system handles user-generated content reports. Users flag content (tracks, albums, artists, playlists, comments, other users), and admins review, resolve, or dismiss reports.

### Report Lifecycle

```
1. User reports content → POST /moderation/report { target_id, target_type, reason, description }
2. Report enters queue → status = 'pending'
3. Admin reviews in moderation page
4. Admin resolves → POST /moderation/resolve/:id { status: 'resolved'|'dismissed', note }
5. Audit trail created → ModerationAction record
```

---

## Moderation Page (`PageAdminModeration.vue` — 768 lines)

### 4-Tab Interface

```
├── Queue tab (pending reports)
│     ├── Type filter: all / track / album / artist / playlist / comment / user
│     ├── Bulk select: checkbox per row + "Select all"
│     ├── Bulk actions: Resolve, Dismiss
│     ├── Report rows:
│     │     ├── Target type icon + name
│     │     ├── Reporter info (username)
│     │     ├── Reason (with priority highlighting)
│     │     ├── Description (truncated, expandable)
│     │     ├── Date/time
│     │     └── Inline action: View → detail dialog, Resolve, Dismiss
│     └── Pagination (offset/limit)
│
├── Flagged tab (active content flags)
│     ├── "Include expired" toggle
│     ├── Flag rows:
│     │     ├── Target type + name
│     │     ├── Flag type
│     │     ├── Expiry indicator (time remaining or "Expired")
│     │     └── Remove flag action
│     └── Pagination
│
├── History tab (resolved/dismissed reports)
│     ├── Status filter: all / resolved / dismissed
│     ├── Report rows + resolution info (moderator, note, date)
│     └── Pagination
│
└── Stats tab
      ├── Metric cards:
      │     ├── Total Reports
      │     ├── Pending
      │     ├── Resolved Today
      │     ├── Flagged Content
      │     ├── Unique Reporters
      │     └── Avg Resolution Time
      ├── Chart: Reports by reason (pie/donut)
      └── Chart: Reports by target type (bar)
```

### Priority Reasons

Certain reasons get highlighted with higher visual priority:

| Reason | Priority | Color |
|--------|----------|-------|
| `copyright` | 🔴 High | Red badge |
| `hate_speech` | 🔴 High | Red badge |
| `explicit` | 🟡 Medium | Yellow badge |
| `spam` | 🟡 Medium | Yellow badge |
| `inappropriate` | 🟢 Normal | Gray badge |
| `misinformation` | 🟢 Normal | Gray badge |

### Report Detail Dialog

```
PrimeVue Dialog with:
  ├── Target: type + name + link to view
  ├── Reported by: username + date
  ├── Reason: badge + description text
  ├── Status: current status
  ├── Resolution history: timeline of actions
  ├── Action buttons:
  │     ├── View Content → opens target in new tab
  │     ├── Resolve → inline form with resolution note
  │     └── Dismiss → inline form with dismissal note
  └── Flag section:
        ├── Flag type: dropdown (inappropriate, copyright, spam, misinformation, hate_speech, explicit)
        ├── Duration: permanent / 24h / 3d / 7d / 30d
        └── Apply Flag button
```

### Flag Dialog

```typescript
// Flag form within report detail:
{
  flag_type: 'inappropriate' | 'copyright' | 'spam' | 'misinformation' | 'hate_speech' | 'explicit',
  expires_in_hours?: number,  // null = permanent
  // duration presets:
  //   null → permanent
  //   24  → 24h
  //   72  → 3d
  //   168 → 7d
  //   720 → 30d
}
```

---

## Reporting Flow (User-Facing)

### Report Button

Users can report content from various surfaces:

| Location | Action | How |
|----------|--------|-----|
| Track context menu | "Report" → reason picker | Right-click or overflow menu |
| Album/Artist page | "Report" button | In overflow menu or actions bar |
| Comment | "Report" icon | Inline on each comment |
| User profile | "Report" → "Report user" | In overflow menu |

### Report Dialog (User-Facing)

```
Dialog with:
  ├── Header: "Report [target type]"
  ├── Reason: dropdown (inappropriate, copyright, spam, misinformation, hate_speech, explicit)
  ├── Description: textarea (optional, max 500 chars)
  ├── CTA: "Submit Report"
  └── Success: "Thanks for reporting — our team will review this content."
```

### Rate Limiting

```typescript
// POST /moderation/report
// Rate limit: 5 reports per hour per user
// On limit hit: "You've submitted too many reports. Please try again later."
```

---

## Backend Moderation API

### Routes

```go
// routes.go
r := rg.Group("/moderation")
{
    r.POST("/report", handler.Report)                    // public — any user can report
    r.GET("/pending", authMW, handler.ListPending)        // admin
    r.GET("/status/:status", authMW, handler.ListByStatus) // admin
    r.GET("/reports/:id", authMW, handler.GetReport)      // admin
    r.POST("/resolve/:id", authMW, handler.Resolve)       // admin
    r.POST("/bulk", authMW, handler.BulkAction)           // admin
    r.POST("/flag", authMW, handler.FlagContent)          // admin
    r.GET("/flags", authMW, handler.ListFlags)            // admin
    r.GET("/stats", authMW, handler.GetStats)             // admin
    r.GET("/actions", authMW, handler.GetActions)         // admin
}
```

### Report Model

```go
// model.go
type ContentReport struct {
    ID             uuid.UUID `json:"id"`
    ReporterID     uuid.UUID `json:"reporter_id"`
    TargetID       uuid.UUID `json:"target_id"`
    TargetType     string    `json:"target_type"` // track, album, artist, playlist, comment, user
    Reason         string    `json:"reason"`
    Description    string    `json:"description"`
    Status         string    `json:"status"` // pending, resolved, dismissed
    ModeratorID    *uuid.UUID `json:"moderator_id"`
    ResolvedAt     *time.Time `json:"resolved_at"`
    ResolutionNote string    `json:"resolution_note"`
    CreatedAt      time.Time `json:"created_at"`
}

type ContentFlag struct {
    ID         uuid.UUID  `json:"id"`
    TargetID   uuid.UUID  `json:"target_id"`
    TargetType string     `json:"target_type"`
    FlagType   string     `json:"flag_type"`
    FlaggedAt  time.Time  `json:"flagged_at"`
    ExpiresAt  *time.Time `json:"expires_at"`
}

type ModerationAction struct {
    ID             uuid.UUID  `json:"id"`
    ReportID       uuid.UUID  `json:"report_id"`
    ModeratorID    uuid.UUID  `json:"moderator_id"`
    Action         string     `json:"action"` // resolve, dismiss, flag, unflag
    TargetID       uuid.UUID  `json:"target_id"`
    TargetType     string     `json:"target_type"`
    PreviousStatus string     `json:"previous_status"`
    NewStatus      string     `json:"new_status"`
    Note           string     `json:"note"`
    Metadata       json.RawMessage `json:"metadata"`
    CreatedAt      time.Time  `json:"created_at"`
}
```

### Moderation Stats

```go
// Aggregated query:
type ModerationStats struct {
    TotalReports        int            `json:"total_reports"`
    PendingReports      int            `json:"pending_reports"`
    ResolvedToday       int            `json:"resolved_today"`
    FlaggedContent      int            `json:"flagged_content"`
    UniqueReporters     int            `json:"unique_reporters"`
    AvgResolutionHours  float64        `json:"avg_resolution_hours"`
    ByReason            map[string]int `json:"by_reason"`
    ByTargetType        map[string]int `json:"by_target_type"`
}
```

### Audit Trail

Every moderation action creates an audit record:

```
Resolve → creates ModerationAction { action: "resolve", previous_status: "pending", new_status: "resolved" }
Dismiss → creates ModerationAction { action: "dismiss", previous_status: "pending", new_status: "dismissed" }
Flag    → creates ModerationAction { action: "flag", metadata: { flag_type, expires_at } }
Bulk    → creates ModerationAction per report { action: "bulk_resolve" or "bulk_dismiss" }
```

---

## State Matrix

| State | Queue Tab | Flagged Tab | History Tab | Stats Tab |
|-------|-----------|-------------|-------------|-----------|
| 🟢 Loading | ✅ Skeleton rows | ✅ Skeleton rows | ✅ Skeleton rows | ✅ Skeleton cards + chart placeholder |
| 🟢 Loaded (has data) | ✅ Report rows with actions | ✅ Flag rows with expiry | ✅ History with filters | ✅ Metrics + charts |
| 🟢 Loaded (no data) | ✅ "No pending reports — all clear!" | ✅ "No active flags" | ✅ "No resolved reports yet" | ✅ All zeros with "No data yet" |
| 🟢 Error loading | ❌ Silent — empty table | ❌ Silent | ❌ Silent | ❌ Silent — zeros shown |
| 🟢 Resolve in progress | ✅ Spinner on action button | N/A | N/A | N/A |
| 🟢 Bulk action | ✅ "Resolving X of Y" progress text | N/A | N/A | N/A |
| 🟢 Flag applied | ✅ Toast + report moves to flagged | ✅ New flag appears | N/A | N/A |
| 🟢 404 on target (deleted content) | ⚠️ Shows "Content deleted" placeholder | ⚠️ Same | ⚠️ Same | N/A |
| 🔴 Rate-limited reporter | ❌ No indicator on report row | N/A | N/A | N/A |

## Friction Points

| # | Severity | Journey | Location | Problem | User Impact | Fix |
|---|----------|---------|----------|---------|-------------|-----|
| F-2401 | ⚠️ MAJOR | J24 | `PageAdminModeration.vue` | **No real-time updates** — page must be refreshed to see new reports | Admins miss urgent copyright reports | Add WebSocket push for new reports or auto-poll |
| F-2402 | ⚠️ MAJOR | J24 | `PageAdminModeration.vue` | **No notification** for admins when new report arrives | High-priority reports sit unnoticed | Add admin notification bell badge |
| F-2403 | ⚠️ MAJOR | J24 | `PageAdminModeration.vue` | **Bulk action has no progress bar** — "Resolving X of Y" text only | No visual feedback for long batch operations | Add progress bar with percentage |
| F-2404 | 💡 IMPROVE | J24 | Report flow (user) | **No "already reported" check** — user can report same content multiple times | Duplicate reports for same user | Check existing pending report before creating |
| F-2405 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No sort options** in queue — always newest first | Can't prioritize by severity | Add sort by reason priority, date, reporter |
| F-2406 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No moderator assignment** — any admin can act on any report | No accountability for resolution | Add "Assign to me" button |
| F-2407 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No searching within reports** — can't find reports about specific content | Hard to investigate repeat offenders | Add search by target name, reporter name |
| F-2408 | 💡 IMPROVE | J24 | Report flow (user) | **No edit/cancel** after submitting report | User submits wrong reason, can't correct | Add "Edit report" or "Cancel" within 5 min |
| F-2409 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No appeal mechanism** for resolved reports (especially dismissed) | User can't contest a dismissal | Add "Appeal" button with reason |
| F-2410 | 💡 IMPROVE | J24 | Flag system | **No auto-unflag** when expired flag's content is still problematic | Content remains unflagged after expiry | Add auto-flag renewal for repeat offenders |
| F-2411 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **No keyboard shortcuts** — all actions require mouse | Slow workflow for high-volume moderation | Add keyboard shortcuts (r=resolve, d=dismiss, j/k=navigate) |
| F-2412 | 💡 IMPROVE | J24 | `PageAdminModeration.vue` | **Stats tab not auto-refreshing** — stale until page reload | Can't monitor moderation effectiveness | Add periodic stats refresh (every 60s) |

## RTL / A11y / Mobile Notes

- ❌ Report row action buttons (Resolve/Dismiss) are icon-only with no `aria-label`
- ✅ `AdminDeleteConfirm` has proper focus management
- ❌ Stats charts have no text alternative for screen readers
- ✅ Pagination has proper `aria-label` — "Page X of Y"
- ❌ Bulk select checkbox doesn't announce selection count to screen readers
- ✅ Report description truncation has expand/collapse with `aria-expanded`
