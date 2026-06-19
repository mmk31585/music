<template>
  <div class="mx-auto w-full max-w-5xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Import"
      title="Import from Internet"
      description="Discover tracks via Spotify, Deezer, MusicBrainz, Last.fm, or your local catalog, then import into the ingestion pipeline."
    />

    <div class="mt-6">
      <div class="flex gap-3">
        <IconField class="flex-1">
          <InputIcon><i aria-hidden="true" class="pi pi-search"></i></InputIcon>
          <InputText
            v-model="query"
            placeholder="Search for a track — e.g. 'feel it from d4vd'"
            class="w-full"
            @keydown.enter="doSearch"
          />
        </IconField>
        <Button
          label="Search"
          icon="pi pi-search"
          :loading="searching"
          :disabled="!query.trim()"
          @click="doSearch"
        />
      </div>
    </div>

    <div v-if="searching" class="mt-12 text-center">
      <i aria-hidden="true" class="pi pi-spin pi-spinner text-3xl text-slate-400"></i>
      <p class="mt-3 text-sm text-slate-500">Searching across Spotify, Deezer, MusicBrainz, Last.fm, and local catalog...</p>
    </div>

    <div v-else-if="searchError" class="mt-6">
      <Message severity="error" :closable="false">{{ searchError }}</Message>
    </div>

    <div v-else-if="importJobId" class="mt-6">
      <div class="rounded-xl border border-white/[0.06] bg-white/[0.03] p-6 text-center">
        <i aria-hidden="true" class="pi pi-spin pi-spinner text-3xl text-emerald-400"></i>
        <p class="mt-3 text-sm font-medium text-white">Import in progress...</p>
        <p class="mt-1 text-xs text-slate-500">{{ importStage }}</p>
        <div class="mx-auto mt-4 h-2 w-full max-w-md overflow-hidden rounded-full bg-white/[0.06]">
          <div
            class="h-full rounded-full bg-emerald-500 transition-all duration-500"
            :style="{ width: importProgress + '%' }"
          />
        </div>
        <p class="mt-2 text-xs text-slate-500">{{ importProgress }}%</p>
      </div>
    </div>

    <div v-else-if="results.length > 0" class="mt-6 space-y-3">
      <p class="text-sm text-slate-500">{{ results.length }} result{{ results.length !== 1 ? 's' : '' }}</p>

      <div
        v-for="(r, i) in results"
        :key="r.url"
        class="group flex items-center gap-4 rounded-xl border border-white/[0.06] bg-white/[0.03] p-4 transition hover:border-white/[0.12]"
      >
        <img
          v-if="r.thumbnail"
          :src="r.thumbnail"
          alt=""
          class="h-16 w-16 shrink-0 rounded-lg object-cover"
        />
        <div v-else class="flex h-16 w-16 shrink-0 items-center justify-center rounded-lg bg-white/[0.06]">
          <i aria-hidden="true" class="pi pi-music text-xl text-slate-500"></i>
        </div>

        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <p class="truncate text-sm font-medium text-white">{{ r.title }}</p>
            <span
              class="shrink-0 rounded-full px-2 py-0.5 text-[10px] font-medium uppercase tracking-wider"
              :class="sourceBadge(r.source)"
            >{{ r.source }}</span>
            <span
              v-if="r.score && r.score > 0"
              class="shrink-0 text-[10px] text-slate-500"
            >{{ (r.score * 100).toFixed(0) }}%</span>
          </div>
          <p class="mt-0.5 text-xs text-slate-400">{{ r.artist }}</p>
          <p class="mt-0.5 text-xs text-slate-500">
            {{ formatDuration(r.duration) }}
            <template v-if="r.isrc"> &bull; ISRC: {{ r.isrc }}</template>
          </p>
        </div>

        <Button
          :label="importingUrl === r.url ? 'Importing...' : 'Import'"
          :icon="importingUrl === r.url ? 'pi pi-spin pi-spinner' : 'pi pi-download'"
          :loading="importingUrl === r.url"
          :disabled="!!importingUrl"
          size="small"
          @click="doImport(r)"
        />
      </div>
    </div>

    <div v-else-if="searched" class="mt-12 text-center">
      <i aria-hidden="true" class="pi pi-search text-3xl text-slate-500"></i>
      <p class="mt-3 text-sm text-slate-500">No results found. Try a different search term.</p>
    </div>

    <div v-else class="mt-12 text-center">
      <i aria-hidden="true" class="pi pi-cloud-download text-3xl text-slate-500"></i>
      <p class="mt-3 text-sm text-slate-500">Search for a track to get started.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import Message from 'primevue/message'
import { AdminSectionHeader } from '@/components/admin'
import { useImportApi } from '@/services/api/importcmd'
import type { SearchResult } from '@/services/api/importcmd'

const router = useRouter()
const toast = useToast()
const importApi = useImportApi()

const query = ref('')
const results = ref<SearchResult[]>([])
const searching = ref(false)
const searchError = ref('')
const searched = ref(false)
const importingUrl = ref('')
const importJobId = ref('')
const importProgress = ref(0)
const importStage = ref('')
let progressTimer: ReturnType<typeof setInterval> | null = null

function sourceBadge(source: string): string {
  const map: Record<string, string> = {
    spotify: '!bg-emerald-500/20 !text-emerald-400',
    deezer: '!bg-purple-500/20 !text-purple-400',
    musicbrainz: '!bg-blue-500/20 !text-blue-400',
    lastfm: '!bg-red-500/20 !text-red-400',
    local: '!bg-slate-500/20 !text-slate-400',
    bandcamp: '!bg-cyan-500/20 !text-cyan-400',
    soundcloud: '!bg-orange-500/20 !text-orange-400',
    youtube: '!bg-rose-500/20 !text-rose-400',
    archiveorg: '!bg-amber-500/20 !text-amber-400',
  }
  return map[source.toLowerCase()] || '!bg-slate-500/20 !text-slate-400'
}

async function doSearch() {
  const q = query.value.trim()
  if (!q) return

  searching.value = true
  searchError.value = ''
  results.value = []
  searched.value = false

  try {
    results.value = await importApi.search(q)
    searched.value = true
  } catch (err: unknown) {
    searchError.value = err instanceof Error ? err.message : 'Search failed.'
  } finally {
    searching.value = false
  }
}

async function doImport(r: SearchResult) {
  importingUrl.value = r.url || r.title
  try {
    const res = await importApi.importTrack(r)
    if (res.jobId) {
      importJobId.value = res.jobId
      importProgress.value = 0
      importStage.value = 'Queued...'
      startProgressPolling(res.jobId)
      toast.add({
        severity: 'info',
        summary: 'Import queued',
        detail: `Track added to import queue.`,
        life: 3000,
      })
    } else if (res.draftId) {
      toast.add({
        severity: 'success',
        summary: 'Import complete',
        detail: `${res.title} added to ingestion.`,
        life: 4000,
      })
      router.push({ name: 'admin.ingestion.review', params: { id: res.draftId } })
    }
  } catch (err: unknown) {
    toast.add({
      severity: 'error',
      summary: 'Import failed',
      detail: err instanceof Error ? err.message : 'Could not import track.',
      life: 5000,
    })
  } finally {
    importingUrl.value = ''
  }
}

function startProgressPolling(jobId: string) {
  progressTimer = setInterval(async () => {
    try {
      const progress = await importApi.getProgress(jobId)
      importProgress.value = progress.progress
      importStage.value = progress.stage || progress.status

      if (progress.status === 'complete') {
        stopProgressPolling()
        toast.add({
          severity: 'success',
          summary: 'Import complete',
          detail: 'Track added to ingestion.',
          life: 4000,
        })
        if (progress.draftId) {
          router.push({ name: 'admin.ingestion.review', params: { id: progress.draftId } })
        }
      } else if (progress.status === 'failed') {
        stopProgressPolling()
        toast.add({
          severity: 'error',
          summary: 'Import failed',
          detail: progress.error || 'Unknown error',
          life: 5000,
        })
        importJobId.value = ''
      }
    } catch {
      stopProgressPolling()
      importJobId.value = ''
    }
  }, 1000)
}

function stopProgressPolling() {
  if (progressTimer) {
    clearInterval(progressTimer)
    progressTimer = null
  }
}

onUnmounted(() => {
  stopProgressPolling()
})

function formatDuration(seconds: number): string {
  if (!seconds || seconds <= 0) return '—'
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}
</script>
