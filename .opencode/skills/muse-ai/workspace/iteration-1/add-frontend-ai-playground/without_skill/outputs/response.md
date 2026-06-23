# Adding an AI Playground page at `/ai/experiment`

Based on existing patterns (PageAIPlaylistGenerator.vue, PageAIMoodExplorer.vue, app.ts routes).

---

## 1. Create the page file

`frontend/src/pages/app/PageAIExperiment.vue`

Follow the existing template pattern:
- Layout: `<div class="mx-auto w-full max-w-5xl px-4 pt-6 pb-32 md:px-6 lg:px-8">`
- Section heading with `AI-Powered` label, `<h1>`, and description `<p>`
- Import `useAIApi` from `@/services/api/ai` and call `generatePlaylist({ prompt: prompt.value })`
- Capture the raw response JSON (store the full response object, then `JSON.stringify` it for display)
- PrimeVue components (`Textarea`, `Button`) are **not** auto-imported — import explicitly:

```ts
import Textarea from 'primevue/textarea'
import Button from 'primevue/button'
```

Key states to cover:
- **Empty/initial**: prompt form with a textarea and "Generate" button
- **Loading**: show a spinner / disabled button with `loading` prop
- **Success**: render track list (reuse the track row pattern from PageAIPlaylistGenerator) + a collapsible `<pre>` block showing `JSON.stringify(response, null, 2)`
- **Error**: catch + toast (PrimeVue `useToast`)

The `generatePlaylist` payload type is:

```ts
interface GeneratePlaylistPayload {
  prompt: string
  mood?: string
  activity?: string
  seed_track_id?: string
  genre?: string
  limit?: number
  exclude_ids?: string[]
}
```

The response type is `AIPlaylistResponse` (import from `@/services/api/ai/types`):

```ts
interface AIPlaylistResponse {
  id: string
  name: string
  description: string
  tracks: Array<{ id: string; title: string; artist?: string; cover_url?: string; duration?: number }>
  generated_at: string
}
```

---

## 2. Register the route

Add a new child route entry inside the `LayoutMusicApp` children array in `frontend/src/router/routes/app.ts`:

```ts
{
  path: 'ai/experiment',
  name: 'ai.experiment',
  component: () => import('@/pages/app/PageAIExperiment.vue'),
  meta: { title: 'AI Playground' },
},
```

The name `ai.experiment` follows the existing dotted convention (`ai.mood-explorer`, `ai.playlist-generator`).

---

## 3. Wire up navigation

If you want a nav link, add it in the appropriate sidebar/menu component (likely under @ui or @catalog in the agent architecture). Look for existing AI-related nav items to follow the same pattern — search for `ai.mood-explorer` or `mood-explorer` in the sidebar/nav components.

---

## 4. Reference files to look at

| Purpose | File |
|---|---|
| Existing AI page with form → API → results | `frontend/src/pages/app/PageAIPlaylistGenerator.vue` |
| Route registration pattern | `frontend/src/router/routes/app.ts` |
| AI API composable + types | `frontend/src/services/api/ai/routes.ts` |
| Response type | `frontend/src/services/api/ai/types.ts` |
| Player integration (play from results) | `frontend/src/pages/app/PageAIPlaylistGenerator.vue:309-319` |

---

## 5. Notes

- The `useRequest` composable (used internally by `useAIApi`) may throw on error — always wrap calls in try/catch
- Use `useToast()` from `primevue/usetoast` for error/success feedback
- For the raw JSON display, a `<pre>` block with `text-xs` and a monospace utility class keeps it readable without polluting the main UI — consider putting it in a collapsible `<Accordion>` or a toggle section
- Agent boundary: `PageAIExperiment.vue` falls under **@catalog** (covers `pages/app/`, `services/api/ai/`). The route addition falls under **@infra** (covers `router/routes/app.ts`). Nav link would be **@ui**.
