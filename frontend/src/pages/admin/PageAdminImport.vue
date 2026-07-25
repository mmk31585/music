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
          <InputIcon><Search aria-hidden="true" class=""></Search></InputIcon>
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
      <Loader2 aria-hidden="true" class="text-3xl text-slate-400 animate-spin"></Loader2>
      <p class="mt-3 text-sm text-slate-500">Searching across Spotify, Deezer, MusicBrainz, Last.fm, and local catalog...</p>
    </div>

    <div v-else-if="searchError" class="mt-6">
      <Message severity="error" :closable="false">{{ searchError }}</Message>
    </div>

    <!-- No-source error (shows alongside results, not replacing them) -->
    <div v-if="importError && !searching && !importJobId" class="mt-6">
      <div class="rounded-xl border border-amber-500/20 bg-amber-500/5 p-5">
        <div class="flex items-start gap-3">
          <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-amber-500/15">
            <Info aria-hidden="true" class="text-lg text-amber-400"></Info>
          </div>
          <div class="flex-1">
            <h4 class="text-sm font-semibold text-amber-300">No downloadable source found</h4>
            <p class="mt-1 text-xs text-slate-400">
              We couldn't find a downloadable version of this track. Try these alternatives:
            </p>
            <ul class="mt-2 space-y-1 text-xs text-slate-400">
              <li class="flex items-center gap-1.5">
                <Youtube aria-hidden="true" class="text-[10px] text-red-400"></Youtube>
                Paste a YouTube URL directly in the search box
              </li>
              <li class="flex items-center gap-1.5">
                <Cloud aria-hidden="true" class="text-[10px] text-orange-400"></Cloud>
                Paste a SoundCloud URL directly in the search box
              </li>
              <li class="flex items-center gap-1.5">
                <Search aria-hidden="true" class="text-[10px] text-slate-400"></Search>
                Try a different search term or check for spelling errors
              </li>
            </ul>
            <Button
              label="Dismiss"
              size="small"
              severity="secondary"
              text
              class="mt-3"
              @click="importError = ''"
            />
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="importJobId" class="mt-6">
      <div class="rounded-xl border border-white/6 bg-white/3 p-6 text-center">
        <Loader2 aria-hidden="true" class="text-3xl text-emerald-400 animate-spin"></Loader2>
        <p class="mt-3 text-sm font-medium text-white">Import in progress...</p>
        <p class="mt-1 text-xs text-slate-500">{{ importStage }}</p>
        <div class="mx-auto mt-4 h-2 w-full max-w-md overflow-hidden rounded-full bg-white/6">
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
        class="group flex items-center gap-4 rounded-xl border border-white/6 bg-white/3 p-4 transition hover:border-white/12"
      >
        <img
          v-if="r.thumbnail"
          :src="r.thumbnail"
          alt=""
          class="h-16 w-16 shrink-0 rounded-lg object-cover"
        />
        <div v-else class="flex h-16 w-16 shrink-0 items-center justify-center rounded-lg bg-white/6">
          <Music aria-hidden="true" class="text-xl text-slate-500"></Music>
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
          :disabled="!!importingUrl && importingUrl !== r.url"
          size="small"
          @click="doImport(r)"
        />
      </div>
    </div>

    <div v-else-if="searched" class="mt-12 text-center">
      <Search aria-hidden="true" class="text-3xl text-slate-500"></Search>
      <p class="mt-3 text-sm text-slate-500">No results found. Try a different search term.</p>
    </div>

    <div v-else class="mt-12 text-center">
      <CloudDownload aria-hidden="true" class="text-3xl text-slate-500"></CloudDownload>
      <p class="mt-3 text-sm text-slate-500">Search for a track to get started.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Cloud, CloudDownload, Info, Loader2, Music, Search, Youtube } from 'lucide-vue-next'
import { ref, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { AdminSectionHeader } from '@/components/admin'
import { useImportApi } from '@/services/api/importcmd'
import type { SearchResult } from '@/services/api/importcmd'
import { formatDuration } from '@/utils/format'

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
const importError = ref('')
let progressTimer: ReturnType<typeof setInterval> | null = null

function sourceBadge(source: string): string {
  const map: Record<string, string> = {
    spotify: 'bg-emerald-500/20! text-emerald-400!',
    deezer: 'bg-purple-500/20! text-purple-400!',
    musicbrainz: 'bg-blue-500/20! text-blue-400!',
    lastfm: 'bg-red-500/20! text-red-400!',
    local: 'bg-slate-500/20! text-slate-400!',
    bandcamp: 'bg-cyan-500/20! text-cyan-400!',
    soundcloud: 'bg-orange-500/20! text-orange-400!',
    youtube: 'bg-rose-500/20! text-rose-400!',
    archiveorg: 'bg-amber-500/20! text-amber-400!',
  }
  return map[source.toLowerCase()] || 'bg-slate-500/20! text-slate-400!'
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
  importError.value = ''
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
    const msg = err instanceof Error ? err.message : 'Could not import track.'
    const isNoSourceError = /could not find.*downloadable|no downloadable source|no suitable source|unable to find.*source/i.test(msg)
    if (isNoSourceError) {
      importError.value = msg
    } else {
      toast.add({
        severity: 'error',
        summary: 'Import failed',
        detail: msg,
        life: 5000,
      })
    }
  } finally {
    importingUrl.value = ''
  }
}

function stageLabel(stage: string): string {
  const labels: Record<string, string> = {
    queued: 'Waiting in queue...',
    resolving: 'Searching for a downloadable source...',
    downloading: 'Downloading track...',
    extracting: 'Extracting metadata...',
    uploading: 'Uploading...',
    complete: 'Import complete!',
    failed: 'Import failed',
  }
  return labels[stage] || stage || 'Processing...'
}

function stageIcon(stage: string): string {
  const icons: Record<string, string> = {
    queued: 'pi pi-clock',
    resolving: 'pi pi-search',
    downloading: 'pi pi-cloud-download',
    extracting: 'pi pi-file',
    uploading: 'pi pi-cloud-upload',
    complete: 'pi pi-check-circle',
    failed: 'pi pi-times-circle',
  }
  return icons[stage] || 'pi pi-spin pi-spinner'
}

function startProgressPolling(jobId: string) {
  progressTimer = setInterval(async () => {
    try {
      const progress = await importApi.getProgress(jobId)
      importProgress.value = progress.progress
      importStage.value = stageLabel(progress.stage || progress.status)

      if (progress.status === 'complete') {
        stopProgressPolling()
        toast.add({
          severity: 'success',
          summary: 'Import complete',
          detail: 'Track added to ingestion. Redirecting to review...',
          life: 3000,
        })
        if (progress.draftId) {
          setTimeout(() => {
            router.push({ name: 'admin.ingestion.review', params: { id: progress.draftId } })
          }, 1000)
        }
      } else if (progress.status === 'failed') {
        stopProgressPolling()
        const errorMsg = progress.error || 'Unknown error'
        const isNoSource = /could not find.*downloadable|no downloadable source|no suitable source|unable to find.*source/i.test(errorMsg)
        if (isNoSource) {
          importError.value = errorMsg
          importJobId.value = ''
        } else {
          toast.add({
            severity: 'error',
            summary: 'Import failed',
            detail: errorMsg,
            life: 5000,
          })
          importJobId.value = ''
        }
      }
    } catch {
      stopProgressPolling()
      importJobId.value = ''
    }
  }, 3000)
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

</script>
