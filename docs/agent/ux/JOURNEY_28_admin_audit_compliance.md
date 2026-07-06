# Journey 28: Admin — Audit & Compliance

> Full trace: existing audit trails → compliance gaps → GDPR/privacy → data retention.

---

## Current Audit Infrastructure

### What Exists

The only audit trail in the system is **moderation actions**:

```go
// internal/modules/moderation/model.go
type ModerationAction struct {
    ID             uuid.UUID       `json:"id"`
    ReportID       uuid.UUID       `json:"report_id"`
    ModeratorID    uuid.UUID       `json:"moderator_id"`
    Action         string          `json:"action"`          // resolve, dismiss, flag, unflag
    TargetID       uuid.UUID       `json:"target_id"`
    TargetType     string          `json:"target_type"`
    PreviousStatus string          `json:"previous_status"`
    NewStatus      string          `json:"new_status"`
    Note           string          `json:"note"`
    Metadata       json.RawMessage `json:"metadata"`
    CreatedAt      time.Time       `json:"created_at"`
}

// GET /moderation/actions?report_id=:id → fetches audit trail for a specific report
```

### What's Missing

| Audit Category | Exists? | Where? |
|---------------|---------|--------|
| Moderation actions | ✅ Yes | `moderation_actions` table |
| Catalog changes (track/album/artist edits) | ❌ No | No history of who changed what |
| User management actions (suspend/ban/role change) | ❌ No | No trail of admin actions on users |
| Admin login/logout | ❌ No | No admin session audit |
| Content deletion | ❌ No | No record of who deleted what and when |
| Flag changes | ❌ No | No history of flag toggles |
| Configuration changes | ❌ No | No audit of setting changes |
| Data exports | ❌ No | No record of who exported data |
| Impersonation events | ❌ No | (feature doesn't exist yet) |

---

## Compliance Requirements

### GDPR / Privacy Features

| Requirement | Status | Notes |
|-------------|--------|-------|
| Data export (user downloads their data) | ❌ Missing | F-1910 noted in settings |
| Account deletion | ❌ Missing | F-1911 — labeled "Coming soon" |
| Data retention policy enforcement | ❌ Missing | No auto-purge of old data |
| Consent management | ❌ Missing | No consent records |
| Right to be forgotten | ❌ Missing | No cascade delete for user data |
| Cookie consent banner | ❌ Not found | No cookie consent UI found |
| Privacy policy page | ❌ Not found | No `/privacy` route |
| Terms of service page | ❌ Not found | No `/terms` route |

### CCPA / Additional Compliance

| Requirement | Status |
|-------------|--------|
| Opt-out of data sale | ❌ Missing |
| California-specific disclosures | ❌ Missing |
| Data inventory / map | ❌ Missing |

---

## Admin Audit Log Needs

### What a Comprehensive Audit Log Should Track

```
AdminAction:
  ├── id (UUID)
  ├── admin_id (who did it)
  ├── action_type (catalog_update, user_suspend, settings_change, data_export, etc.)
  ├── target_type (track, album, artist, user, setting, etc.)
  ├── target_id
  ├── previous_value (JSON)
  ├── new_value (JSON)
  ├── ip_address
  ├── user_agent
  ├── created_at
  └── metadata (JSON — additional context)
```

### Audit Log View

An admin audit log page should provide:

| Feature | Purpose |
|---------|---------|
| **Timeline view** | Chronological list of all admin actions |
| **Filters** | By admin, action type, target type, date range |
| **Search** | Free-text search across actions |
| **Detail expand** | Show previous/new values as diff |
| **Export** | Download audit log as CSV |
| **Retention config** | How long to keep audit logs |
| **Real-time feed** | Live stream of admin actions via WebSocket |

---

## Data Retention

### Current State

There is **no data retention policy enforcement** in the codebase:
- Listening history grows indefinitely
- Analytics events accumulate without purge
- Media files are never cleaned up
- User accounts persist after deletion request (feature doesn't exist)
- Moderation actions persist forever (this is probably desired)

### Recommended Retention Policies

| Data Type | Retention | Justification |
|-----------|-----------|---------------|
| Listening history | 12 months | User preference, privacy |
| Raw analytics events | 90 days | Aggregated rollups retained permanently |
| Search queries | 30 days | Aggregated search terms retained |
| User accounts (after deletion request) | 30 days grace, then purge | GDPR right to be forgotten |
| Media files (unused/orphaned) | 30 days | Storage cleanup |
| Audit logs (admin) | Permanent | Compliance |
| Moderation records | Permanent | Compliance |
| Session tokens | Until expiry | Security |
| Email notifications | 90 days | Operational |
| Chat messages | 12 months | Social features |

---

## State Matrix

| State | Audit Log | GDPR Features | Data Retention |
|-------|-----------|---------------|----------------|
| 🟢 Page exists | ❌ No audit page | ❌ No GDPR page | ❌ No retention page |
| 🟢 Moderation audit | ✅ Limited (per-report only) | N/A | N/A |
| 🟢 Catalog change tracking | ❌ Nonexistent | N/A | N/A |
| 🔴 Deletion request | N/A | ❌ Not implemented (F-1911) | N/A |
| 🔴 Data export | N/A | ❌ Not implemented (F-1910) | N/A |
| 🔴 Auto-purge | N/A | N/A | ❌ No purge jobs |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-2801 | ⚠️ MAJOR | All admin | **No general-purpose audit log** — only moderation actions tracked | No accountability for catalog changes, user management, or settings | Add AdminAction model + store for all admin operations |
| F-2802 | ⚠️ MAJOR | Platform | **No GDPR compliance features** — no data export, no account deletion, no consent | Legal risk, user trust erosion | Implement data export + account deletion + consent management |
| F-2803 | ⚠️ MAJOR | Platform | **No data retention policy** — no auto-purge of old data | Storage grows unbounded, privacy risk | Add configurable retention policies with automated cleanup |
| F-2804 | 💡 IMPROVE | Audit | **No diff view** for changes — can't see before/after | Hard to understand what changed | Add JSON diff viewer in audit detail |
| F-2805 | 💡 IMPROVE | Audit | **No admin session audit** — no record of who logged in and when | Can't detect unauthorized admin access | Add admin login/logout event tracking |
| F-2806 | 💡 IMPROVE | Compliance | **No privacy policy page** — no `/privacy` route | Legal requirement unmet | Add privacy policy page |
| F-2807 | 💡 IMPROVE | Compliance | **No terms of service page** — no `/terms` route | Legal requirement unmet | Add terms of service page |
| F-2808 | 💡 IMPROVE | Compliance | **No cookie consent banner** | GDPR requirement for EU users | Add cookie consent banner with preference management |
| F-2809 | 💡 IMPROVE | Audit | **No pagination/search on moderation audit** — per-report only, no global view | Can't search across all moderation actions | Add global audit log view with search/filter |
| F-2810 | 💡 IMPROVE | Compliance | **No "right to be forgotten" cascade** — deleting user doesn't clean up all related data | User data persists after account deletion | Add cascade delete for all user-associated data |
| F-2811 | 💡 IMPROVE | Audit | **No export of audit log** — can't download for compliance review | Must manually screenshot | Add CSV export with date range filter |

## RTL / A11y / Mobile Notes

- ❌ Privacy policy and terms pages would need full RTL support
- ❌ Audit log timeline needs proper `aria-live` for real-time updates
- ✅ JSON diff viewer can use standard `<del>` / `<ins>` with accessible color contrast
