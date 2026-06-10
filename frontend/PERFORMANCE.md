# Performance Guide

## Rendering Optimizations

### 1. TrackRow — Avoid Cascading Computed Chains

Don't wrap track data into a single `PlaybackTrack` computed that 3+ other computeds depend on.
Instead, compute primitives directly (`trackId`) and only build the full object at call time.

**Before:** 4 computed re-evaluations per prop change
**After:** 1 computed (`trackId`) + inline object creation in `handlePlay`

### 2. Karaoke/Lyrics — No Array Allocation on Time Tick

Never `.map()` on every `props.currentTime` update (4x/sec during playback).
Instead, compute only `activeLineIdx` as a single number, and use `:class="idx === activeLineIdx"` in the template.

**Before:** New 100+ item array every 250ms → DOM re-render of all lyrics
**After:** Single number comparison per line

### 3. v-for Keys

Every `v-for` loop must have `:key`:
- Dynamic lists → use unique `item.id`
- Static arrays → use the item itself or a stable index
- Skeleton loaders → use `:key="i"` (iteration variable)

All 100+ `v-for` loops in this project are properly keyed.

## Reactivity

### 4. Avoid `deep: true` on Primitive Watchers

`deep` only matters for nested objects/arrays. Watching arrays of primitives (string, number, boolean) with `deep: true` adds unnecessary traversal overhead.

```ts
// Bad — deep adds no value for primitives
watch([a, b, c], fn, { deep: true })

// Good
watch([a, b, c], fn)
```

### 5. Combine Related Watchers

When two refs trigger the same handler, use a single watcher over an array:

```ts
// Bad — two watcher registrations
watch(query, () => search())
watch(type, () => search())

// Good
watch([query, type], () => search())
```

### 6. Avoid Spread of Entire Store

Spreading a store (`{ ...usePlayerStore() }`) into a component watches ALL its refs.
Destructure only what you need.

```ts
// Bad — binds to 30+ reactive refs
const { currentTime, isPlaying, ...rest } = usePlayerStore()

// Good — only what you use
const currentTime = usePlayerStore().currentTime
const isPlaying = usePlayerStore().isPlaying
```

## Canvas & Animation

### 7. Throttle Resize Listeners

Canvas/particle resize handlers should use `requestAnimationFrame`:

```ts
let raf: number
function onResize() {
  cancelAnimationFrame(raf)
  raf = requestAnimationFrame(actualResize)
}
window.addEventListener('resize', onResize)
```

Applied in: `VisualizerSystem.vue`, `useLyricsParticles.ts`

### 8. Self-Cleaning Composables

Every composable that registers DOM listeners or animation frames should call `onUnmounted()` internally:

```ts
export function useFeature() {
  onUnmounted(() => {
    removeEventListener(...)
    cancelAnimationFrame(...)
  })
}
```

This prevents memory leaks if a consumer forgets to call `stop()`.

## Bundle Size

### 9. Prefer Named Imports Over Barrels for Heavy Modules

Barrel files (`index.ts` with `export *`) prevent tree-shaking in many bundlers.
Import from the specific module path when you only need one export:

```ts
// Instead of:
import { useAuthApi } from '@/services/api'

// Use:
import { useAuthApi } from '@/services/api/auth'
```

The top-level `@/services/api` barrel is useful for discoverability, but
critical paths (auth, player) should use direct imports.

### 10. Dynamic Imports on All Routes

All 40+ route components use `() => import(...)` for automatic code splitting.
This keeps the initial bundle under 200KB (gzipped).

## Monitoring

| Metric | Target | How |
|--------|--------|-----|
| Time-to-interactive | <2s | Route splitting + skeleton loading |
| Search response | <300ms | 250ms debounce + keyboard nav |
| Track play | 1 click | Play buttons on all cards/rows |
| Animation FPS | 60 | requestAnimationFrame-throttled canvas |
| Bundle size | <200KB gzip | Dynamic imports on all routes |
