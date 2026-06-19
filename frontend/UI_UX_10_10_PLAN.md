# 10/10 UI/UX Plan — Soundify

Target: All categories → 10/10. Priority: Player experience + Mobile responsiveness.

## Category Targets

| Category | Current | Target | Key Levers |
|----------|---------|--------|------------|
| Visual Quality | 9 | 10 | Consistent glassmorphism, micro-interactions, dynamic theming |
| Interaction | 9 | 10 | Swipe gestures, haptic feedback patterns, press states |
| Accessibility | 8 | 10 | Focus trapping, aria-live, skip links on all layouts |
| Layout | 8 | 10 | Perfect responsive, landscape, tablet optimization |
| **Overall** | **8.5** | **10** | |

## Phase Plan — Execution Order

### Phase 1: Fix ExpandedPlayer (highest impact)
- Remove hardcoded `dir="rtl"` — bug for English UI
- Connect shuffle/repeat/speed/sleep/crossfade to actual player store

### Phase 2: Fix FullscreenPlayer
- Add dynamic album art background color extraction
- Improve transition and visual polish

### Phase 3: Fix MobileBottomSheet
- Responsive height instead of fixed 560px
- Safe area awareness + bottom nav overlap fix
- Enhanced drag-to-dismiss

### Phase 4: Fix NowPlayingBar mobile
- Add volume slider on mobile
- Better compact layout for small screens
- Visual press feedback

### Phase 5: Fix LayoutMusicApp
- Bottom nav + player bar overlap resolution
- Consistent padding on all breakpoints

### Phase 6: Add swipe gestures to player
- Swipe left/right on album art → next/prev track
- Swipe down to dismiss fullscreen player

### Phase 7: Final review + score verification

## Technical Notes

- All changes must pass `vue-tsc --noEmit`
- Use `useAlbumColors` composable for dynamic colors
- Use `useSwipe` composable for gesture support
- Account for `env(safe-area-inset-bottom)` everywhere
