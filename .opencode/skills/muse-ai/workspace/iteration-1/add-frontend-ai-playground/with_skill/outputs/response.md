# Walkthrough: AI Playground Page (`/ai/experiment`)

This guide covers the full stack — backend Go endpoint that proxies arbitrary prompts to the AI and returns the raw JSON, then the Vue frontend page, route registration, and navigation link.

---

## 1. Backend: New "experiment" endpoint

The page needs an endpoint that accepts an arbitrary prompt, feeds it to the AI, and returns **both** the parsed track selection **and** the raw JSON the LLM responded with. We'll add a new method to the existing AI module.

### 1a. Add DTOs — `internal/modules/ai/dto.go`

Append these types after the existing `TrackItem` struct:

```go
type AIExperimentRequest struct {
    Prompt string   `json:"prompt" binding:"required"`
    Tracks []string `json:"tracks,omitempty"` // optional seed track IDs
    Limit  int      `json:"limit,omitempty"`
}

type AIExperimentResponse struct {
    ID          string      `json:"id"`
    Prompt      string      `json:"prompt"`
    Tracks      []TrackItem `json:"tracks"`
    RawResponse string      `json:"raw_response"` // raw JSON from the LLM
    GeneratedAt string      `json:"generated_at"`
}
```

`RawResponse` holds the verbatim JSON the OpenAI chat completion returned so the frontend can display it in a collapsible panel.

### 1b. Add `Experiment` to the `AIClient` interface — `internal/modules/ai/client.go`

Add the method to the interface:

```go
type AIClient interface {
    GenerateEmbedding(ctx context.Context, input string) ([]float64, error)
    AnalyzeMood(ctx context.Context, title, artist, genre string) (*TrackMood, error)
    GeneratePlaylist(ctx context.Context, prompt string, tracks []TrackMeta) ([]string, error)
    Experiment(ctx context.Context, prompt string, tracks []TrackMeta) ([]string, string, error) // +
}
```

Implement on `openAIClient` — similar to `GeneratePlaylist` but also returns the raw content string:

```go
func (c *openAIClient) Experiment(ctx context.Context, prompt string, tracks []TrackMeta) ([]string, string, error) {
    trackList, _ := json.Marshal(tracks)
    systemMsg := `You are a music curator. Select track IDs matching the user's request. Return JSON with "track_ids" (array of strings).`
    userMsg := fmt.Sprintf(`User request: "%s"\nAvailable tracks: %s\nSelect up to 20 track IDs.`, prompt, string(trackList))

    messages := []openAIMessage{
        {Role: "system", Content: systemMsg},
        {Role: "user", Content: userMsg},
    }

    body := openAIChatReq{
        Model:       c.model,
        Messages:    messages,
        Temperature: 0.7,
        MaxTokens:   1000,
        ResponseFormat: &responseFormat{Type: "json_object"},
    }

    data, _ := json.Marshal(body)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/chat/completions", bytes.NewReader(data))
    req.Header.Set("Content-Type", "application/json")
    if c.apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+c.apiKey)
    }

    resp, _ := c.client.Do(req)
    defer resp.Body.Close()
    respBody, _ := io.ReadAll(resp.Body)

    var result openAIChatResp
    json.Unmarshal(respBody, &result)

    rawContent := result.Choices[0].Message.Content // raw JSON string

    var selection struct{ TrackIDs []string `json:"track_ids"` }
    json.Unmarshal([]byte(rawContent), &selection)

    return selection.TrackIDs, rawContent, nil
}
```

Implement on `fallbackClient` for when no API key is configured:

```go
func (f *fallbackClient) Experiment(ctx context.Context, prompt string, tracks []TrackMeta) ([]string, string, error) {
    if len(tracks) == 0 {
        return nil, `{"track_ids":[]}`, fmt.Errorf("no tracks available")
    }
    limit := 20
    if limit > len(tracks) { limit = len(tracks) }
    ids := make([]string, limit)
    for i := 0; i < limit; i++ { ids[i] = tracks[i].ID }
    raw, _ := json.Marshal(map[string]any{"track_ids": ids})
    return ids, string(raw), nil
}
```

### 1c. Add service method — `internal/modules/ai/service.go`

```go
func (s *Service) Experiment(ctx context.Context, req AIExperimentRequest) (*AIExperimentResponse, error) {
    if !s.enabled {
        return nil, fmt.Errorf("AI features are disabled")
    }

    candidates, err := s.repo.GetTracks(ctx, 50)
    if err != nil {
        return nil, fmt.Errorf("get candidate tracks: %w", err)
    }

    selectedIDs, rawJSON, err := s.ai.Experiment(ctx, req.Prompt, candidates)
    if err != nil {
        s.logger.Warn("AI experiment failed, using fallback", zap.Error(err))
        selectedIDs = s.fallbackSelect(candidates, req.limitOrDefault())
        rawJSON = fmt.Sprintf(`{"fallback":true,"track_ids":%s}`, toJSON(selectedIDs))
    }

    tracks, _ := s.repo.GetTracksByIDs(ctx, selectedIDs)
    items := makeTrackItems(tracks)

    return &AIExperimentResponse{
        ID:          fmt.Sprintf("exp_%d", time.Now().Unix()),
        Prompt:      req.Prompt,
        Tracks:      items,
        RawResponse: rawJSON,
        GeneratedAt: time.Now().UTC().Format(time.RFC3339),
    }, nil
}

func (r AIExperimentRequest) limitOrDefault() int {
    if r.Limit <= 0 || r.Limit > 100 { return 20 }
    return r.Limit
}
```

Helper functions to extract (or inline):

```go
func makeTrackItems(tracks []TrackMeta) []TrackItem {
    items := make([]TrackItem, len(tracks))
    for i, t := range tracks {
        items[i] = TrackItem{
            ID: t.ID, Title: t.Title, Artist: t.Artist,
            Duration: t.Duration, CoverURL: t.CoverURL,
        }
    }
    return items
}

func toJSON(v any) string { b, _ := json.Marshal(v); return string(b) }
```

### 1d. Add handler — `internal/modules/ai/handler.go`

```go
func (h *Handler) HandleExperiment(c *gin.Context) {
    var req AIExperimentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
        return
    }

    result, err := h.service.Experiment(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
```

### 1e. Register route — `internal/modules/ai/routes.go`

Add a single line inside the `aiGroup` block:

```go
aiGroup.POST("/experiment", h.HandleExperiment)
```

---

## 2. Frontend: New API route and composable method

### 2a. Add enum — `frontend/src/services/api/ai/enums.ts`

```ts
export enum AIApiRoutes {
  // ... existing routes ...
  EXPERIMENT = '/ai/experiment',
}
```

### 2b. Add types — `frontend/src/services/api/ai/types.ts`

```ts
export interface AIExperimentPayload {
  prompt: string
  tracks?: string[]
  limit?: number
}

export const AIExperimentResponseSchema = z.object({
  id: z.string(),
  prompt: z.string(),
  tracks: z.array(AITrackItemSchema),
  raw_response: z.string(),
  generated_at: z.string(),
})

export type AIExperimentResponse = z.infer<typeof AIExperimentResponseSchema>
```

### 2c. Add composable method — `frontend/src/services/api/ai/routes.ts`

Import the new types, then add to the `useAIApi` return object:

```ts
import {
  AIExperimentResponseSchema,
  type AIExperimentResponse,
  type AIExperimentPayload,
} from './types'

// Inside return { ... }:
const experiment = async (payload: AIExperimentPayload, config?: UseRequestConfig<AIExperimentResponse>) => {
  return useRequest<AIExperimentResponse>(
    AIApiRoutes.EXPERIMENT,
    { method: 'POST', data: payload },
    { schema: AIExperimentResponseSchema, silent: false, ...config },
  )
}

return {
  generatePlaylist,
  analyzeMood,
  getMood,
  generateEmbedding,
  similarByMood,
  similarByEmbedding,
  experiment,
}
```

---

## 3. Frontend: The `PageAIExperiment.vue` component

Create `frontend/src/pages/app/PageAIExperiment.vue` following the same structural patterns as `PageAIPlaylistGenerator.vue` and `PageAIMoodExplorer.vue`:

### Template structure

```vue
<template>
  <div class="mx-auto w-full max-w-5xl px-4 pt-6 pb-32 md:px-6 lg:px-8">
    <!-- Hero header (same pattern as PlaylistGenerator) -->
    <div class="mb-8">
      <p class="text-xs font-bold tracking-[0.25em] text-[#1db954] uppercase">AI-Powered</p>
      <h1 class="mt-2 text-3xl font-black text-white md:text-4xl">AI Playground</h1>
      <p class="mt-2 text-sm text-slate-400">
        Send a custom prompt to the AI and inspect the raw response.
      </p>
    </div>

    <div class="grid gap-8 lg:grid-cols-5">
      <!-- Left column: prompt input -->
      <div class="space-y-6 lg:col-span-2">
        <div class="rounded-2xl border border-white/[0.06] bg-white/[0.03] p-6">
          <h2 class="mb-5 text-sm font-bold text-white">Experiment</h2>
          <div class="space-y-4">
            <div>
              <label class="mb-2 block text-xs font-medium text-slate-400">Prompt</label>
              <Textarea
                v-model="prompt"
                placeholder="e.g. Recommend happy upbeat tracks for a party"
                :auto-resize="true"
                rows="4"
                class="w-full"
              />
            </div>
            <Button
              label="Send to AI"
              icon="pi pi-send"
              severity="success"
              class="mt-2 w-full"
              :loading="loading"
              :disabled="loading || !prompt.trim()"
              @click="handleExperiment"
            />
          </div>
        </div>
      </div>

      <!-- Right column: results + raw JSON -->
      <div class="lg:col-span-3 space-y-4">
        <!-- Loading state -->
        <div v-if="loading" class="flex flex-col items-center justify-center rounded-2xl border border-white/[0.06] bg-white/[0.03] px-6 py-24 text-center">
          <i aria-hidden="true" class="pi pi-spin pi-spinner text-3xl text-[#1db954]" />
          <p class="mt-4 text-sm font-medium text-white">Consulting the AI...</p>
        </div>

        <!-- Empty state (same dashed-card pattern as MoodExplorer) -->
        <div v-else-if="!result" class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-white/[0.08] bg-white/[0.02] px-6 py-24 text-center">
          <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-gradient-to-br from-[#1db954]/20 to-cyan-500/20">
            <i aria-hidden="true" class="pi pi-code text-2xl text-[#1db954]" />
          </div>
          <h3 class="mt-6 text-lg font-bold text-white">Enter a prompt</h3>
          <p class="mt-2 max-w-xs text-sm text-slate-400">
            Type a request above and the AI will pick tracks you might like.
          </p>
        </div>

        <!-- Results -->
        <template v-else>
          <!-- Track results (same track-list pattern as PlaylistGenerator) -->
          <div class="rounded-2xl border border-white/[0.06] bg-white/[0.02]">
            <div class="flex items-center justify-between border-b border-white/[0.06] px-4 py-3">
              <h2 class="text-sm font-bold text-white">Selected Tracks ({{ result.tracks.length }})</h2>
              <Button label="Play All" icon="pi pi-play" severity="success" size="small"
                :disabled="!result.tracks.length" @click="playAll" />
            </div>
            <!-- track rows -- same exact markup as PageAIPlaylistGenerator lines 171-214 -->
            <div v-for="(track, index) in result.tracks" :key="track.id"
              class="group flex items-center gap-3 px-4 py-2 transition hover:bg-white/[0.06]">
              <span class="w-6 text-right text-xs text-slate-500">{{ index + 1 }}</span>
              <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
                <img v-if="track.cover_url" :src="track.cover_url" :alt="track.title"
                  loading="lazy" class="h-full w-full object-cover" @error="onImgError" />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-music text-xs text-slate-500" />
                </div>
                <button type="button" aria-label="Play track"
                  class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 transition group-hover:opacity-100"
                  @click="playTrack(index)">
                  <i aria-hidden="true" class="pi pi-play-fill text-xs text-white" />
                </button>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                <p v-if="track.artist" class="truncate text-xs text-slate-400">{{ track.artist }}</p>
              </div>
              <span class="text-xs text-slate-500">{{ formatTime(track.duration) }}</span>
            </div>
          </div>

          <!-- Raw JSON accordion -->
          <div class="rounded-2xl border border-white/[0.06] bg-white/[0.03] overflow-hidden">
            <button type="button"
              class="flex w-full items-center justify-between px-4 py-3 text-sm font-bold text-white transition hover:bg-white/[0.04]"
              @click="showRaw = !showRaw">
              <span>Raw Response</span>
              <i aria-hidden="true" class="pi" :class="showRaw ? 'pi-chevron-up' : 'pi-chevron-down'" />
            </button>
            <Transition name="collapse">
              <pre v-if="showRaw"
                class="m-0 overflow-x-auto border-t border-white/[0.06] bg-black/30 p-4 text-xs text-slate-300 leading-relaxed">{{
                  formatJSON(result.raw_response) }}</pre>
            </Transition>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>
```

### Script setup

```ts
<script setup lang="ts">
import { ref } from 'vue'
import { useAIApi } from '@/services/api/ai'
import type { AIExperimentResponse } from '@/services/api/ai/types'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { onImgError } from '@/utils/helpers'
import { useToast } from 'primevue/usetoast'

const aiApi = useAIApi()
const player = usePlayer()
const playerApi = usePlayerApi()
const toast = useToast()

const prompt = ref('')
const loading = ref(false)
const result = ref<AIExperimentResponse | null>(null)
const showRaw = ref(false)

async function handleExperiment() {
  loading.value = true
  showRaw.value = false
  try {
    const res = await aiApi.experiment({ prompt: prompt.value, limit: 20 })
    if (res) result.value = res
  } catch {
    toast.add({ severity: 'error', summary: 'Experiment failed', detail: 'Could not get AI response.', life: 3000 })
  } finally {
    loading.value = false
  }
}

function playTrack(index: number) {
  if (!result.value) return
  const queue = result.value.tracks.map((t) => ({
    id: t.id, title: t.title, artistName: t.artist || 'Unknown',
    coverUrl: t.cover_url || null, durationSeconds: t.duration ?? null,
    streamUrl: playerApi.getTrackStreamUrl(t.id),
  }))
  player.setQueueAndPlay(queue, index)
}

function playAll() {
  if (result.value?.tracks.length) playTrack(0)
}

function formatTime(seconds?: number) {
  if (!seconds) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

function formatJSON(raw: string) {
  try { return JSON.stringify(JSON.parse(raw), null, 2) }
  catch { return raw }
}
</script>
```

Add the collapse transition in a `<style scoped>` block (same pattern used in `LayoutMusicApp.vue` for the mobile nav):

```vue
<style scoped>
.collapse-enter-active,
.collapse-leave-active {
  transition: all 0.2s ease;
  overflow: hidden;
}
.collapse-enter-from,
.collapse-leave-to {
  max-height: 0;
  opacity: 0;
}
</style>
```

---

## 4. Route registration — `frontend/src/router/routes/app.ts`

Add a new child route inside the `LayoutMusicApp` children array, following the pattern of the existing AI routes (lines 176-186):

```ts
{
  path: 'ai/experiment',
  name: 'ai.experiment',
  component: () => import('@/pages/app/PageAIExperiment.vue'),
  meta: { title: 'AI Playground' },
},
```

---

## 5. Navigation link — `frontend/src/layouts/LayoutMusicApp.vue`

Two places to update:

### 5a. `moreItems` array (line 236-245)

Add a new entry, keeping alphabetical order relative to the existing AI links:

```ts
const moreItems = [
  // ...
  { label: 'AI Mood Explorer', icon: 'pi pi-magic', to: '/ai/mood-explorer' },
  { label: 'AI Playground', icon: 'pi pi-code', to: '/ai/experiment' },
  { label: 'AI Playlist Generator', icon: 'pi pi-sync', to: '/ai/playlist-generator' },
  // ...
]
```

### 5b. `pageTitleMap` (line 247-265)

Add the title for the new route:

```ts
const pageTitleMap: Record<string, string> = {
  // ...
  '/ai/mood-explorer': 'AI Mood Explorer',
  '/ai/playlist-generator': 'AI Playlist Generator',
  '/ai/experiment': 'AI Playground',
  // ...
}
```

---

## Summary of file changes

| File | Change |
|---|---|
| `internal/modules/ai/dto.go` | Add `AIExperimentRequest`, `AIExperimentResponse` |
| `internal/modules/ai/client.go` | Add `Experiment()` to interface + 2 implementations |
| `internal/modules/ai/service.go` | Add `Experiment()` service method + helpers |
| `internal/modules/ai/handler.go` | Add `HandleExperiment` HTTP handler |
| `internal/modules/ai/routes.go` | Add `POST /experiment` route |
| `frontend/src/services/api/ai/enums.ts` | Add `EXPERIMENT` enum |
| `frontend/src/services/api/ai/types.ts` | Add `AIExperimentPayload`, `AIExperimentResponse` types + Zod schema |
| `frontend/src/services/api/ai/routes.ts` | Add `experiment()` method to `useAIApi()` |
| `frontend/src/pages/app/PageAIExperiment.vue` | **New** — playground page |
| `frontend/src/router/routes/app.ts` | Add route `ai/experiment` → `PageAIExperiment` |
| `frontend/src/layouts/LayoutMusicApp.vue` | Add nav link + page title |
