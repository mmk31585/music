---
description: 'Comprehensive Vue 3 development standards: Composition API, `<script setup>`, Pinia, Vue Router, TypeScript, testing, and performance for Muse Persian music platform'
applyTo: '**/*.vue, **/*.ts, **/*.js'
---

# Vue 3 Development Instructions

Authoritative guidance for building production-grade Vue 3 applications for the Muse platform. Default to the **Composition API** with `<script setup lang="ts">`, the modern reactivity system, and the official ecosystem (Pinia, Vue Router, Vite, Vitest).

## Project Context
- Vue 3.4+ with Composition API
- `<script setup lang="ts">` single-file components (SFCs) as the default authoring style
- TypeScript everywhere: components, composables, stores, and router
- Pinia for state management; Vue Router for routing; Vite for build/dev
- Vitest + Vue Test Utils for tests
- Glassmorphism design system with Persian (RTL) support
- RTL-aware layouts and components

## Authoring Style & Component Design
- Use `<script setup>` — it is more concise, faster, and has better type inference than `setup()` or the Options API
- One responsibility per component; split large components into smaller focused ones plus composables
- Order an SFC as `<script setup>`, then `<template>`, then `<style scoped>`
- Name components in PascalCase; use multi-word names (e.g. `UserCard`, not `Card`)
- Co-locate component-specific types, and lift shared types into a `types/` module

## Compiler Macros (no imports needed)
- `defineProps<T>()` — declare typed props from a TypeScript interface/type for full inference
- `withDefaults(defineProps<T>(), { ... })` — provide prop defaults
- `defineEmits<{ change: [id: number]; update: [value: string] }>()` — declare typed events
- `defineModel<T>()` (3.4+) — the canonical way to implement `v-model` on a component
- `defineExpose({ ... })` — explicitly expose a public imperative API
- `defineSlots<{ default(props: { item: T }): any }>()` — type named/scoped slots
- `defineOptions({ name, inheritAttrs })` — set component options inside `<script setup>`
- Never mutate props directly — emit an event, use `defineModel`, or derive local state with `computed`/`ref`

## Reactivity System
### Core primitives
- `ref()` for primitives and single replaceable references; access via `.value` in script (auto-unwrapped in templates)
- `reactive()` for deep-reactive objects/collections; never destructure it directly — use `toRefs()`/`toRef()`
- `computed()` for derived values; keep getters pure and side-effect free
- Prefer `computed` over `watch` whenever you are *deriving* a value rather than performing a side effect

### Watchers
- `watch(source, cb, options)` for explicit dependencies; `watchEffect(cb)` for auto-tracked dependencies
- Use watch options deliberately: `{ immediate: true }`, `{ deep: true }`, `{ once: true }` (3.4+), and `flush: 'post'`
- Register cleanup with the `onCleanup`/`onWatcherCleanup` callback to cancel stale async work

### Advanced reactivity (use intentionally)
- `shallowRef` / `shallowReactive` for large or externally-managed data
- `readonly()` to hand out immutable views of shared state
- `toRef` / `toRefs` to keep reactivity when destructuring
- `effectScope()` to group and dispose related effects together
- `customRef` for debounced/throttled or storage-backed refs

## Composables (reusable logic)
- Extract stateful, reusable logic into `useXxx()` functions under `composables/`
- Accept refs/getters as inputs and return refs/computed; use `toValue()`/`MaybeRefOrGetter` to normalize ref-or-plain inputs
- Set up and tear down inside the composable (`onMounted`/`onUnmounted` or `tryOnScopeDispose`)
- Keep composables synchronous in their setup phase; expose async actions as returned functions

## Lifecycle & Effects
- Use `onMounted`, `onBeforeMount`, `onUpdated`, `onBeforeUnmount`, `onUnmounted`, `onActivated`/`onDeactivated` (with `<KeepAlive>`), and `onErrorCaptured`
- Always clean up timers, listeners, observers, and subscriptions in `onUnmounted`
- Guard browser-only APIs (`window`, `document`) for SSR; run them in `onMounted`

## Template Best Practices
- Always set a stable, unique `:key` on `v-for`; never use the array index when items can reorder or mutate
- Never put `v-if` and `v-for` on the same element — filter via a `computed` instead
- `v-show` for frequently toggled elements; `v-if` for conditional mounting
- Use `v-memo` to skip re-rendering of expensive static subtrees, and `v-once` for content that renders a single time
- Use the `:` (v-bind) and `@` (v-on) shorthands consistently
- Avoid heavy expressions in templates — move them to `computed` or methods
- For RTL support: use logical CSS properties (`margin-inline-start`, `padding-inline-end`, `inset-inline-start`)

## Slots
- Use named slots for layout extension and scoped slots (`<slot :item="item" />` + `#default="{ item }"`)
- Provide sensible fallback slot content

## Built-in Components
- `<Teleport to="body">` for modals, toasts, and tooltips
- `<Suspense>` with `#default`/`#fallback` for async setup and lazy components
- `<Transition>` / `<TransitionGroup>` for enter/leave and list animations
- `<KeepAlive>` (with `include`/`exclude`/`max`) to cache component state
- `<component :is="...">` for dynamic components; `defineAsyncComponent(() => import('...'))` for code-split/lazy loading

## Provide / Inject (dependency injection)
- Type injections with an `InjectionKey<T>` (`Symbol`) for safety: `provide(key, value)` / `inject(key)`
- Provide a default or assert presence to avoid `undefined`
- Prefer `readonly()` when providing state that children should not mutate

## State Management with Pinia
- Use Pinia for shared/cross-component state; keep component-only state local with `ref`/`reactive`
- Prefer **setup stores**: `defineStore('user', () => { ... })`
- One store per domain; keep actions for async/side effects and getters pure & synchronous
- Destructure with `storeToRefs()` to preserve reactivity
- Use `$patch` for batched mutations, `$reset` to restore state

## Routing with Vue Router
- Define routes with lazy `component: () => import('...')` for automatic code splitting
- Use navigation guards (`beforeEach`, `beforeEnter`, `beforeRouteLeave`) for auth and unsaved-changes checks
- Use `route.meta` (typed) for per-route config like `requiresAuth`
- Read params/query via `useRoute()` and navigate via `useRouter()`

## TypeScript Integration
- Type props/emits/slots through generic compiler macros, not runtime object syntax
- Type refs explicitly when inference is too narrow: `ref<User | null>(null)`
- Type template refs with `useTemplateRef<HTMLInputElement>('input')` (3.5+) or `ref<HTMLInputElement | null>(null)`
- Build generic components with `<script setup lang="ts" generic="T">`

## Styling
- Default to `<style scoped>`; use `:deep()`, `:slotted()`, and `:global()` selectors deliberately
- Use `v-bind()` in `<style>` to drive CSS from reactive state
- Prefer CSS custom properties for theming (see design system at docs/DESIGN_SYSTEM.md)
- Use logical CSS properties for RTL support (e.g., `border-inline-start` instead of `border-left`)

## Accessibility (RTL-specific)
- Ensure full keyboard operability and visible focus states
- Set `dir="auto"` or `dir="rtl"` on RTL-aware components
- Manage focus on route changes and when opening/closing dialogs
- Give icon-only controls accessible names (`aria-label`)
- Associate labels with inputs; use semantic HTML

## Performance
- Code-split routes and heavy components (`defineAsyncComponent`, dynamic `import()`)
- Use `computed` for caching, `v-memo`/`v-once` for static subtrees
- Virtualize long lists; paginate or window large data
- Avoid unnecessary deep reactivity and creating new object/array literals inline in templates

## Anti-Patterns to Avoid
- Mixing Options API and Composition API arbitrarily in the same codebase
- Mutating props directly, or destructuring `reactive()`/Pinia stores without `toRefs`/`storeToRefs`
- Using `watch` for values that should be `computed`
- `v-if` together with `v-for` on one element; using array index as `:key`
- Heavy logic inside templates; unbounded deep reactivity on large data
- Leaking timers/listeners by skipping `onUnmounted` cleanup
- Rendering untrusted HTML via `v-html`
- Hardcoding LTR-specific CSS in RTL contexts
