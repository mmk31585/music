---
name: unit-test-vue-pinia
category: testing
description: 'Write and review unit tests for Vue 3 + TypeScript + Vitest + Pinia codebases. Covers createTestingPinia patterns, Vue Test Utils, and black-box assertions.'
---

# unit-test-vue-pinia

Write or review unit tests for Vue components, composables, and Pinia stores. Keep tests small, deterministic, and behavior-first.

## Core Rules
- Test one behavior per test
- Assert observable input/output first (rendered text, emitted events, store state)
- Avoid implementation-coupled assertions
- Prefer explicit `beforeEach()` setup; reset mocks every test
- Use `references/pinia-patterns.md` as local source of truth

## Pinia Testing Patterns
- Default: `createTestingPinia({ createSpy: vi.fn })` as global plugin
- Seed state: `createTestingPinia({ initialState: { storeName: { key: val } } })`
- Real actions only when needed: `stubActions: false`
- Pure store tests: `setActivePinia(createPinia())` for state transitions

## Vue Test Utils
- Mount shallow by default for focused unit tests
- Drive through props, user interactions, emitted events
- Use `nextTick` only for async updates
- Prefer `wrapper.emitted(...)` for event assertions

## References
- `references/pinia-patterns.md`
- Pinia testing cookbook: https://pinia.vuejs.org/cookbook/testing.html
- Vue Test Utils guide: https://test-utils.vuejs.org/guide/
