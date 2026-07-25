---
description: 'Guidelines for Persian (Farsi) localization and RTL support in the Muse music platform'
applyTo: '**/*.vue,**/*.ts,**/*.go,**/*.md'
---

# Persian Localization & RTL Guidelines

Muse is a Persian music ecosystem. All user-facing content must support Persian (Farsi) language and Right-to-Left (RTL) layout.

## RTL Layout Principles

### CSS & Styling
- Use **logical CSS properties** instead of directional ones:
  - `margin-inline-start` / `margin-inline-end` instead of `margin-left` / `margin-right`
  - `padding-inline-start` / `padding-inline-end` instead of `padding-left` / `padding-right`
  - `border-inline-start` / `border-inline-end` instead of `border-left` / `border-right`
  - `inset-inline-start` / `inset-inline-end` instead of `left` / `right`
  - `text-align: start` / `text-align: end` instead of `text-align: left` / `text-align: right`
- Set `direction: rtl` on the root element when Persian is active
- Use `html[dir="rtl"]` or `[dir="rtl"]` selectors for RTL-specific overrides
- Test all UI at 320px width in both LTR and RTL modes

### Flexbox & Grid
- Use `gap` instead of margin on flex/grid children
- Use `justify-content: flex-start` / `flex-end` (these are RTL-aware)
- Avoid using `order` to rearrange elements for RTL — use logical properties

## Persian Text Handling

### Typography
- Use Persian-compatible fonts (e.g., Vazir, IRANSans, or system Persian fonts)
- Ensure font weights render correctly in Persian script
- Persian text uses a different baseline; account for this in vertical alignment
- Line-height for Persian text may need adjustment (typically 1.8x instead of 1.5x)

### Numbers & Dates
- Use `Intl.NumberFormat` with `fa-IR` locale for Persian numerals
- Use `Intl.DateTimeFormat` with `fa-IR` locale for Jalali calendar dates
- Support both Persian digits (۱۲۳) and Arabic digits (123) based on user preference
- Format currencies in Iranian Rial or Toman as appropriate

### String Handling
- Persian text may contain Arabic characters; normalize using `String.prototype.normalize()`
- Use `toLocaleLowerCase('fa')` for case conversion
- Be aware that Persian punctuation differs from English (e.g., comma, question mark)

## Go Backend Localization

### String Resources
- Store all user-facing strings in locale files (JSON or Go embedded)
- Structure: `locales/fa.json`, `locales/en.json`
- Use key-based lookups: `i18n.T("player.now_playing")`
- Support parameter interpolation: `i18n.T("tracks.count", count)`

### API Response Headers
- Set `Content-Language` header based on user's locale preference
- Accept `Accept-Language` header for content negotiation
- Store user locale preference in the database (JWT or user profile)

## Vue Frontend Localization

### i18n Setup
- Use `vue-i18n` for internationalization
- Configure with `locale` from user preference store
- Lazy-load locale messages for performance
- Use `$t()` or `t()` in templates and `setup()`

### RTL Detection
- Derive `dir` attribute from current locale: `locale === 'fa' ? 'rtl' : 'ltr'`
- Apply `dir` to `<html>` element reactively
- Store RTL preference in Pinia store

### Component Patterns
- Use a `useRTL()` composable that exposes `isRTL`, `dir`, and logical CSS class helpers
- RTL-aware components should accept a `dir` prop or derive from store
- Lottie animations may need to be flipped for RTL

## Testing RTL
- Test every component in both LTR and RTL modes
- Use Playwright to snapshot RTL layouts
- Verify text alignment, icon placement, and overflow in RTL
- Test keyboard navigation (arrow keys should respect RTL direction)
