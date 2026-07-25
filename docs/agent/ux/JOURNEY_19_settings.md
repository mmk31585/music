# Journey 19: Settings — Audio, Appearance, Privacy, Notifications

> Full trace: settings page → 3 tabs → profile editing → account management → preferences → persistence.

---

## Settings Architecture

### Route

```typescript
// Route: /settings → PageUserSettings.vue (name: 'settings')
// Meta: { title: 'Settings', requiresAuth: true }
```

Settings is **auth-guarded** — guests cannot access it (they have no profile to configure).

### Layout

```
PageUserSettings.vue (486 lines)
  ├── SkeletonLoader (loading state)
  └── Tab panel (on loaded)
        ├── Profile tab
        │     ├── Avatar display + edit (file input overlay)
        │     ├── Display Name input
        │     ├── Username input
        │     ├── Bio textarea
        │     ├── Location input
        │     ├── Website input
        │     └── Save changes button
        ├── Account tab
        │     ├── Email (read-only display)
        │     ├── Change Password section
        │     │     ├── Current password field
        │     │     ├── New password field (min 8 chars validation)
        │     │     └── Save password button
        │     └── Delete Account section
        │           └── "Coming soon" (disabled danger zone)
        └── Preferences tab
              ├── RTL layout toggle
              ├── Lyrics autoscroll toggle
              ├── Explicit content toggle
              └── Notification preferences
                    ├── New Releases toggle
                    └── Social Activity toggle
```

---

## Profile Settings

### Avatar Upload

```typescript
const avatarFile = ref<File | null>(null)

function handleAvatarUpload(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    // Validate: max 5MB, image type
    // Preview: createObjectURL
    avatarFile.value = file
  }
}

async function saveProfile() {
  const formData = new FormData()
  if (avatarFile.value) formData.append('avatar', avatarFile.value)
  formData.append('full_name', displayName.value)
  formData.append('username', username.value)
  formData.append('bio', bio.value)
  formData.append('location', location.value)
  formData.append('website', website.value)

  await api.put('/users/me/profile', formData)
  // Success: toast + update local state
}
```

| State | Avatar Display | Behavior |
|-------|---------------|----------|
| Default | Circular fallback (initials) | No image |
| Has avatar | Rounded image | Displayed |
| Hover | Opacity overlay + camera icon | File input click |
| Uploading | Loading spinner overlay | Form submitting |
| Upload error | Previous state + error toast | Rollback |
| File too large | Error toast "Max 5MB" | Reject before upload |
| Wrong format | Error toast "Images only" | Reject before upload |

### Profile Form Fields

| Field | Validation | Max Length | Notes |
|-------|-----------|------------|-------|
| Display Name | Required, non-empty | 100 | Used in header |
| Username | Required, unique, alphanumeric | 30 | `@handle`, used in URL |
| Bio | Optional | 500 | Expandable on profile |
| Location | Optional | 100 | Shown below name |
| Website | Optional, URL format | 200 | Clickable link |

### Save Behavior

```
Save button → PUT /users/me/profile
  → 200: Green success banner + "Profile updated" (aria-live="polite")
  → 400: Red error banner with specific message
  → 409 (username taken): Red error on username field specifically
  → 500: Red error banner "Something went wrong"
```

---

## Account Settings

### Change Password

```typescript
const currentPassword = ref('')
const newPassword = ref('')
const passwordError = ref('')

async function changePassword() {
  if (newPassword.value.length < 8) {
    passwordError.value = 'Password must be at least 8 characters'
    return
  }
  try {
    await api.put('/users/me/password', {
      current_password: currentPassword.value,
      new_password: newPassword.value,
    })
    // Success toast + clear fields
    currentPassword.value = ''
    newPassword.value = ''
  } catch (e) {
    passwordError.value = 'Current password is incorrect'
  }
}
```

| State | Behavior |
|-------|----------|
| Empty | Both fields empty, button disabled |
| Current filled, new empty | Button disabled |
| Current filled, new < 8 chars | Inline error "Min 8 characters" |
| Both filled, valid | Button enabled |
| Save success | Toast + clear fields |
| Save failure (wrong current) | Inline error on current password field |
| Save failure (server error) | Error toast |

### Delete Account

```
Red danger zone card:
  ⚠️ "Delete Account"
  "Permanently delete your account and all data. This action cannot be undone."
  [Delete Account] button (disabled)
  Badge: "Coming soon"
```

The delete account feature is **not implemented**. The button is visibly present but disabled with a "Coming soon" label. This is a deliberate placeholder.

---

## Preferences

### RTL Layout Toggle

```typescript
// composables/useRTL.ts
const RTL_STORAGE_KEY = 'muse-rtl'

const isRTL = ref(localStorage.getItem(RTL_STORAGE_KEY) === 'true')

// Watcher syncs to <html> dir attribute
watch(isRTL, (val) => {
  document.documentElement.dir = val ? 'rtl' : 'ltr'
  localStorage.setItem(RTL_STORAGE_KEY, String(val))
  // Also save to server profile preferences
  savePreferencesToServer()
})
```

**Persistence**: localStorage + server-side (embedded in profile `preferences` JSON).

**Instant effect**: No page reload needed — `dir` attribute changes reactively.

### Lyrics Autoscroll

```typescript
// localStorage key: muse-lyrics-autoscroll
// Default: true
// Controls: FullscreenPlayer lyrics tab auto-scroll behavior
```

### Explicit Content

```typescript
// localStorage key: muse-explicit-content
// Default: true (show explicit)
// Controls: Filter explicit tracks from search, browse, recommendations
// Not synced to server — localStorage only
```

### Notification Preferences

```typescript
// Serverside: stored in user profile preferences JSON
// Two toggles currently:
//   - New Releases: notify when followed artists release
//   - Social Activity: notify on likes, follows, comments

// Future: expand to per-category (see F-1607)
```

---

## Settings Persistence Matrix

| Setting | Storage Location | Synced to Server? | Instant Apply? |
|---------|-----------------|-------------------|----------------|
| Display Name | DB (profile) | ✅ Yes | ✅ Yes |
| Username | DB (profile) | ✅ Yes | ✅ Yes |
| Bio | DB (profile) | ✅ Yes | ✅ Yes |
| Location | DB (profile) | ✅ Yes | ✅ Yes |
| Website | DB (profile) | ✅ Yes | ✅ Yes |
| Avatar | Object storage + DB | ✅ Yes | ✅ Yes |
| Password | DB (hashed) | ✅ Yes | ✅ Yes |
| RTL | localStorage + DB | ✅ Yes | ✅ Yes (reactive) |
| Lyrics autoscroll | localStorage | ❌ No | ✅ Yes |
| Explicit content | localStorage | ❌ No | ✅ Yes |
| Notification prefs | DB (profile.preferences) | ✅ Yes | ✅ Yes |
| Delete account | — | ❌ Not implemented | N/A |

---

## Missing Settings Categories

The settings page is notably missing:

| Missing Setting | Why Important | User Impact |
|----------------|---------------|-------------|
| **Audio quality** | Premium users can't select bitrate | Always 128kbps or always 320kbps |
| **Theme (dark/light)** | No light mode toggle | Locked into dark-only design |
| **Keyboard shortcuts** | Can't disable or customize | F-610 continuation |
| **Crossfade duration** | Toggle exists in player but no setting | F-603 continuation |
| **Language** | No explicit language picker | Inherited from browser/system |
| **Playback** (gapless, normalization) | No audio behavior settings | Can't customize listening experience |
| **Privacy** (activity status, history pause) | No private session toggle | F-1209 continuation |
| **Notifications** (per-category) | Only 2 broad toggles | F-1607 continuation |
| **Download location** / storage | No offline management UI | F-1700+ (no offline yet) |

---

## State Matrix

| State | Settings Page |
|-------|--------------|
| 🟢 Loading | ✅ SkeletonLoader |
| 🟢 Loaded, has data | ✅ 3 tabs with current values |
| 🟢 Save success | ✅ Green success banner (aria-live) |
| 🟢 Save failure | ✅ Red error banner |
| 🟢 Validation error | ✅ Inline field error |
| 🟢 RTL toggle | ✅ Instant reactive change |
| 🟢 Network error on save | ❌ No explicit offline detection |
| 🔴 Concurrent save race | ❌ No debounce — double-click can fire twice |
| 🔴 Avatar upload progress | ❌ No upload progress bar |

## Friction Points

| # | Severity | Location | Problem | User Impact | Fix |
|---|----------|----------|---------|-------------|-----|
| F-1901 | ⚠️ MAJOR | `PageUserSettings.vue` | **No audio quality setting** — not present in UI at all | Premium users can't select 320kbps | Add audio quality selector (128/320/flac) in preferences |
| F-1902 | ⚠️ MAJOR | `PageUserSettings.vue` | **Theme (dark/light) toggle** completely absent | Users locked into dark-only | Add light mode with proper contrast ratios |
| F-1903 | ⚠️ MAJOR | `PageUserSettings.vue` | **No language/locale picker** — locale inherited from browser | Persian users may see English UI | Add language selector (fa/en) |
| F-1904 | 💡 IMPROVE | `PageUserSettings.vue` | Save button has **no debounce** — double-click fires two save requests | Duplicate saves, possibly overwriting data | Add debounce + loading state |
| F-1905 | 💡 IMPROVE | `PageUserSettings.vue` | Avatar upload has **no progress indicator** | User doesn't know if upload is happening | Add progress bar or spinner |
| F-1906 | 💡 IMPROVE | `PageUserSettings.vue` | **No "Reset to defaults"** button for preferences | Accidental changes are permanent | Add "Reset to defaults" per section |
| F-1907 | 💡 IMPROVE | `PageUserSettings.vue` | **No preview** when changing RTL — must save, navigate, and see | Can't A/B test layout preference | Add live preview area |
| F-1908 | 💡 IMPROVE | `PageUserSettings.vue` | Password change has **no strength meter** — only min 8 char check | Users may set weak passwords | Add password strength indicator |
| F-1909 | 💡 IMPROVE | `PageUserSettings.vue` | Notification preferences **only 2 toggles** — no per-category | Users can't granularly control notifications | Add per-category notification settings (F-1607 continuation) |
| F-1910 | 💡 IMPROVE | `PageUserSettings.vue` | **No "Export my data"** GDPR/Privacy option | Users can't download their data | Add data export request button |
| F-1911 | 💡 IMPROVE | `PageUserSettings.vue` | Delete account labeled "Coming soon" with **no timeline or feedback** | Users who want to leave can't | Add request deletion flow (support ticket) |
| F-1912 | 💡 IMPROVE | `PageUserSettings.vue` | **Preferences tab saves independently** — each toggle fires its own API call | Multiple rapid toggles spam the server | Debounce + batch preference saves |

## RTL / A11y / Mobile Notes

- ✅ RTL toggle uses existing `useRTL()` composable with reactive `dir` attribute
- ✅ All form inputs have associated `<label>` elements
- ✅ Success/error messages use `aria-live="polite"`
- ❌ Avatar upload input has poor keyboard accessibility — hidden file input needs explicit label
- ✅ Password fields have proper `type="password"` with toggle visibility
- ❌ No skip link to settings content — keyboard users tab through full nav first
- ✅ Touch targets for all toggle switches are adequate (>44px)
