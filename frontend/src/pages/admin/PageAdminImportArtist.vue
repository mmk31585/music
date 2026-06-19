<template>
  <div class="mx-auto w-full max-w-5xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Import by Artist"
      title="Artist Discography Import"
      description="Enter an artist name to discover all albums and tracks, then select which to import."
    />

    <!-- Search bar -->
    <div class="mt-6">
      <div class="flex gap-3">
        <IconField class="flex-1">
          <InputIcon><i aria-hidden="true" class="pi pi-search"></i></InputIcon>
          <InputText
            v-model="artistName"
            placeholder="Enter artist name — e.g. 'd4vd', 'Arctic Monkeys'"
            class="w-full"
            @keydown.enter="doSearch"
          />
        </IconField>
        <Button
          label="Search Artist"
          icon="pi pi-search"
          :loading="searching"
          :disabled="!artistName.trim()"
          @click="doSearch"
        />
      </div>
    </div>

    <!-- Loading -->
    <div v-if="searching" class="mt-12 text-center">
      <i aria-hidden="true" class="pi pi-spin pi-spinner text-3xl text-slate-400"></i>
      <p class="mt-3 text-sm text-slate-500">Searching Deezer and MusicBrainz for discography...</p>
    </div>

    <!-- Error -->
    <div v-else-if="searchError" class="mt-6">
      <Message severity="error" :closable="false">{{ searchError }}</Message>
    </div>

    <!-- Batch in progress -->
    <div v-else-if="batchId" class="mt-6">
      <div class="rounded-xl border border-white/[0.06] bg-white/[0.03] p-6">
        <div class="text-center">
          <i aria-hidden="true" class="pi pi-spin pi-spinner text-3xl text-emerald-400"></i>
          <p class="mt-3 text-sm font-medium text-white">Importing {{ batchTotal }} tracks...</p>
          <p class="mt-1 text-xs text-slate-400">
            {{ batchCompleted + batchFailed }} / {{ batchTotal }} processed
            <span v-if="batchFailed > 0">({{ batchFailed }} failed)</span>
          </p>
          <div class="mx-auto mt-4 h-2 w-full max-w-md overflow-hidden rounded-full bg-white/[0.06]">
            <div
              class="h-full rounded-full bg-emerald-500 transition-all duration-500"
              :style="{ width: batchProgressPct + '%' }"
            />
          </div>
          <p class="mt-2 text-xs text-slate-500">{{ batchProgressPct }}%</p>
        </div>

        <!-- Per-track status list -->
        <div class="mt-6 space-y-2 border-t border-white/[0.06] pt-4">
          <p class="text-xs font-medium text-slate-400 uppercase tracking-wider">Tracks</p>
          <div
            v-for="job in batchJobs"
            :key="job.title + job.artist"
            class="flex items-center justify-between rounded-lg bg-white/[0.03] px-3 py-2"
          >
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm text-white">{{ job.title }}</p>
              <p class="truncate text-xs text-slate-500">{{ job.artist }}</p>
            </div>
            <span
              class="ml-3 shrink-0 text-xs font-medium"
              :class="batchJobStatusClass(job)"
            >
              {{ batchJobStatusLabel(job) }}
            </span>
          </div>
        </div>

        <!-- Completed actions -->
        <div v-if="batchProgressPct >= 100" class="mt-6 text-center">
          <div class="flex items-center justify-center gap-3">
            <Button
              label="Review in Ingestion"
              icon="pi pi-eye"
              @click="goToIngestion"
            />
            <Button
              label="Search Another Artist"
              icon="pi pi-refresh"
              severity="secondary"
              @click="resetAll"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Results -->
    <template v-else-if="discography">
      <!-- Artist header -->
      <div class="mt-6 flex items-center gap-4">
        <img
          v-if="discography.artist_info.image"
          :src="discography.artist_info.image"
          :alt="discography.artist_info.name"
          class="h-16 w-16 shrink-0 rounded-full object-cover ring-2 ring-white/[0.08]"
        />
        <div v-else class="flex h-16 w-16 shrink-0 items-center justify-center rounded-full bg-white/[0.06]">
          <i aria-hidden="true" class="pi pi-user text-2xl text-slate-500"></i>
        </div>
        <div>
          <h2 class="text-xl font-bold text-white">{{ discography.artist_info.name }}</h2>
          <p class="text-sm text-slate-400">
            {{ totalTracks }} track{{ totalTracks !== 1 ? 's' : '' }} across {{ discography.albums.length }} album{{ discography.albums.length !== 1 ? 's' : '' }}
          </p>
        </div>
      </div>

      <!-- Batch actions bar -->
      <div class="mt-4 flex items-center gap-3 rounded-xl border border-white/[0.06] bg-white/[0.03] px-4 py-3">
        <label class="flex items-center gap-2 text-sm text-slate-300">
          <Checkbox
            :binary="true"
            :model-value="allTracksSelected"
            :indeterminate="someTracksSelected && !allTracksSelected"
            @update:model-value="toggleSelectAll"
          />
          Select all {{ totalTracks }} tracks
        </label>
        <span class="text-xs text-slate-500">{{ selectedCount }} selected</span>
        <div class="ml-auto flex gap-2">
          <Button
            label="Import Selected"
            icon="pi pi-download"
            :disabled="selectedCount === 0"
            :loading="importing"
            size="small"
            @click="doBatchImport"
          />
        </div>
      </div>

      <!-- Albums -->
      <div class="mt-4 space-y-4">
        <div
          v-for="(album, ai) in discography.albums"
          :key="album.title + ai"
          class="overflow-hidden rounded-xl border border-white/[0.06] bg-white/[0.02]"
        >
          <!-- Album header (clickable collapse) -->
          <button
            class="flex w-full items-center gap-4 px-4 py-3 text-left transition hover:bg-white/[0.03]"
            @click="toggleAlbum(ai)"
          >
            <img
              v-if="album.cover"
              :src="album.cover"
              :alt="album.title"
              class="h-12 w-12 shrink-0 rounded-lg object-cover"
            />
            <div v-else class="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg bg-white/[0.06]">
              <i aria-hidden="true" class="pi pi-compact-disc text-lg text-slate-500"></i>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-sm font-medium text-white truncate">{{ album.title || 'Unknown Album' }}</p>
              <p class="text-xs text-slate-500">{{ album.tracks.length }} track{{ album.tracks.length !== 1 ? 's' : '' }}</p>
            </div>
            <div class="flex items-center gap-2">
              <label class="flex items-center gap-1.5 text-xs text-slate-400" @click.stop>
                <Checkbox
                  :binary="true"
                  :model-value="albumTracksSelected(ai)"
                  :indeterminate="albumPartiallySelected(ai)"
                  @update:model-value="(v: boolean) => toggleAlbumTracks(ai, v)"
                />
                Album
              </label>
              <i
                aria-hidden="true"
                class="pi text-sm text-slate-500 transition-transform"
                :class="openAlbums[ai] ? 'pi-chevron-up' : 'pi-chevron-down'"
              ></i>
            </div>
          </button>

          <!-- Tracks (collapsible) -->
          <div v-if="openAlbums[ai]" class="border-t border-white/[0.06]">
            <div
              v-for="(track, ti) in album.tracks"
              :key="track.title + ti"
              class="flex items-center gap-3 px-4 py-2.5 transition hover:bg-white/[0.03]"
            >
              <Checkbox
                :binary="true"
                :model-value="isTrackSelected(ai, ti)"
                @update:model-value="(v: boolean) => toggleTrack(ai, ti, v)"
              />
              <div class="min-w-0 flex-1">
                <p class="text-sm text-white truncate">{{ track.title }}</p>
                <p class="text-xs text-slate-500">
                  {{ formatDuration(track.duration) }}
                  <span v-if="track.source" class="ml-2">via {{ track.source }}</span>
                </p>
              </div>
              <span
                class="shrink-0 rounded-full px-2 py-0.5 text-[10px] font-medium uppercase tracking-wider"
                :class="sourceBadge(track.source)"
              >
                {{ track.source }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div v-if="discography.albums.length === 0" class="mt-12 text-center">
        <i aria-hidden="true" class="pi pi-compact-disc text-3xl text-slate-500"></i>
        <p class="mt-3 text-sm text-slate-500">No albums found for this artist.</p>
      </div>
    </template>

    <!-- Empty state -->
    <div v-else-if="searched" class="mt-12 text-center">
      <i aria-hidden="true" class="pi pi-search text-3xl text-slate-500"></i>
      <p class="mt-3 text-sm text-slate-500">No results found. Try a different artist name.</p>
    </div>

    <div v-else class="mt-12 text-center">
      <i aria-hidden="true" class="pi pi-cloud-download text-3xl text-slate-500"></i>
      <p class="mt-3 text-sm text-slate-500">Enter an artist name to discover their discography.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import Checkbox from 'primevue/checkbox'
import Message from 'primevue/message'
import { AdminSectionHeader } from '@/components/admin'
import { useImportApi } from '@/services/api/importcmd'
import type { AlbumGroup, TrackResult, BatchJobResult, ImportResponse } from '@/services/api/importcmd'

const router = useRouter()
const toast = useToast()
const importApi = useImportApi()

// Search state
const artistName = ref('')
const searching = ref(false)
const searchError = ref('')
const searched = ref(false)
const discography = ref<{
  artist_info: { name: string; image: string }
  albums: AlbumGroup[]
} | null>(null)

// UI state
const openAlbums = ref<Record<number, boolean>>({})
const selectedTracks = ref<Set<string>>(new Set())

// Batch import state
const importing = ref(false)
const batchId = ref('')
const batchTotal = ref(0)
const batchCompleted = ref(0)
const batchFailed = ref(0)
const batchProgressPct = ref(0)
const batchJobs = ref<BatchJobResult[]>([])
let progressTimer: ReturnType<typeof setInterval> | null = null

const totalTracks = computed(() => {
  if (!discography.value) return 0
  return discography.value.albums.reduce((sum, a) => sum + a.tracks.length, 0)
})

const selectedCount = computed(() => selectedTracks.value.size)

const allTracksSelected = computed(() => {
  if (!discography.value || totalTracks.value === 0) return false
  return selectedCount.value === totalTracks.value
})

const someTracksSelected = computed(() => {
  return selectedCount.value > 0 && !allTracksSelected.value
})

function trackKey(albumIdx: number, trackIdx: number): string {
  return `${albumIdx}:${trackIdx}`
}

function isTrackSelected(albumIdx: number, trackIdx: number): boolean {
  return selectedTracks.value.has(trackKey(albumIdx, trackIdx))
}

function toggleTrack(albumIdx: number, trackIdx: number, val: boolean) {
  const key = trackKey(albumIdx, trackIdx)
  if (val) {
    selectedTracks.value.add(key)
  } else {
    selectedTracks.value.delete(key)
  }
  // Trigger reactivity
  selectedTracks.value = new Set(selectedTracks.value)
}

function albumTracksSelected(albumIdx: number): boolean {
  if (!discography.value) return false
  const album = discography.value.albums[albumIdx]
  if (!album) return false
  for (let ti = 0; ti < album.tracks.length; ti++) {
    if (!isTrackSelected(albumIdx, ti)) return false
  }
  return album.tracks.length > 0
}

function albumPartiallySelected(albumIdx: number): boolean {
  if (!discography.value) return false
  const album = discography.value.albums[albumIdx]
  if (!album) return false
  let count = 0
  for (let ti = 0; ti < album.tracks.length; ti++) {
    if (isTrackSelected(albumIdx, ti)) count++
  }
  return count > 0 && count < album.tracks.length
}

function toggleAlbumTracks(albumIdx: number, val: boolean) {
  if (!discography.value) return
  const album = discography.value.albums[albumIdx]
  if (!album) return
  for (let ti = 0; ti < album.tracks.length; ti++) {
    const key = trackKey(albumIdx, ti)
    if (val) {
      selectedTracks.value.add(key)
    } else {
      selectedTracks.value.delete(key)
    }
  }
  selectedTracks.value = new Set(selectedTracks.value)
}

function toggleSelectAll(val: boolean) {
  if (!discography.value) return
  selectedTracks.value = new Set()
  if (val) {
    for (let ai = 0; ai < discography.value.albums.length; ai++) {
      for (let ti = 0; ti < discography.value.albums[ai].tracks.length; ti++) {
        selectedTracks.value.add(trackKey(ai, ti))
      }
    }
  }
}

function toggleAlbum(ai: number) {
  openAlbums.value[ai] = !openAlbums.value[ai]
  openAlbums.value = { ...openAlbums.value }
}

function sourceBadge(source: string): string {
  const map: Record<string, string> = {
    deezer: '!bg-purple-500/20 !text-purple-400',
    musicbrainz: '!bg-blue-500/20 !text-blue-400',
    spotify: '!bg-emerald-500/20 !text-emerald-400',
  }
  return map[source?.toLowerCase()] || '!bg-slate-500/20 !text-slate-400'
}

async function doSearch() {
  const name = artistName.value.trim()
  if (!name) return

  searching.value = true
  searchError.value = ''
  searched.value = false
  discography.value = null
  selectedTracks.value = new Set()
  openAlbums.value = {}

  try {
    const result = await importApi.searchArtist(name)
    discography.value = result
    searched.value = true

    // Open all albums by default
    for (let i = 0; i < result.albums.length; i++) {
      openAlbums.value[i] = true
    }
  } catch (err: unknown) {
    searchError.value = err instanceof Error ? err.message : 'Search failed.'
    searched.value = true
  } finally {
    searching.value = false
  }
}

async function doBatchImport() {
  if (!discography.value || selectedCount.value === 0) return

  const tracks: Array<{
    title: string
    artist: string
    album: string | undefined
    duration: number
    source: string
    external_ids: Record<string, string> | undefined
  }> = []

  for (const key of selectedTracks.value) {
    const [ai, ti] = key.split(':').map(Number)
    const album = discography.value.albums[ai]
    const track = album?.tracks[ti]
    if (album && track) {
      tracks.push({
        title: track.title,
        artist: discography.value.artist_info.name,
        album: album.title,
        duration: track.duration,
        source: track.source || 'deezer',
        external_ids: track.external_ids,
      })
    }
  }

  if (tracks.length === 0) return

  importing.value = true

  try {
    const res = await importApi.batchImport(tracks)
    batchId.value = res.batchId
    batchTotal.value = tracks.length
    batchCompleted.value = 0
    batchFailed.value = 0
    batchProgressPct.value = 0
    batchJobs.value = res.jobs

    toast.add({
      severity: 'info',
      summary: 'Batch import started',
      detail: `Importing ${tracks.length} track${tracks.length !== 1 ? 's' : ''}.`,
      life: 3000,
    })

    // Start progress polling
    startBatchProgressPolling(res.batchId)
  } catch (err: unknown) {
    importing.value = false
    toast.add({
      severity: 'error',
      summary: 'Batch import failed',
      detail: err instanceof Error ? err.message : 'Could not start batch import.',
      life: 5000,
    })
  }
}

function startBatchProgressPolling(batchIdVal: string) {
  if (progressTimer) clearInterval(progressTimer)

  progressTimer = setInterval(async () => {
    try {
      const progress = await importApi.getBatchProgress(batchIdVal)
      if (!progress) return

      batchCompleted.value = progress.completed
      batchFailed.value = progress.failed
      batchProgressPct.value = progress.progressPct

      // Check per-job progress to update status labels
      if (batchJobs.value.length > 0) {
        const updatedJobs = [...batchJobs.value]
        for (let i = 0; i < updatedJobs.length; i++) {
          const job = updatedJobs[i]
          if (job.jobId) {
            try {
              const jobProgress = await importApi.getProgress(job.jobId)
              if (jobProgress) {
                updatedJobs[i] = { ...job, status: jobProgress.status }
              }
            } catch { /* ignore polling errors */ }
          }
        }
        batchJobs.value = updatedJobs
      }

      if (progress.progressPct >= 100) {
        stopProgressPolling()
        importing.value = false

        const failed = progress.failed
        if (failed === 0) {
          toast.add({
            severity: 'success',
            summary: 'All tracks imported!',
            detail: `${progress.completed} track${progress.completed !== 1 ? 's' : ''} added to ingestion.`,
            life: 6000,
          })
        } else {
          toast.add({
            severity: 'warn',
            summary: 'Import completed with errors',
            detail: `${progress.completed} succeeded, ${failed} failed.`,
            life: 6000,
          })
        }
      }
    } catch {
      // Ignore polling errors — keep trying
    }
  }, 2000)
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

function batchJobStatusClass(job: BatchJobResult & { status?: string }): string {
  if (job.error) return '!text-red-400'
  if (!job.jobId) return '!text-slate-500'
  if (job.status === 'complete') return '!text-emerald-400'
  if (job.status === 'failed') return '!text-red-400'
  if (job.status === 'downloading' || job.status === 'uploading') return '!text-amber-400'
  if (job.status === 'queued' || job.status === 'resolving') return '!text-blue-400'
  return '!text-slate-400'
}

function batchJobStatusLabel(job: BatchJobResult & { status?: string }): string {
  if (job.error) return 'Failed'
  if (!job.jobId) return 'Skipped'
  if (job.status === 'complete') return 'Complete ✓'
  if (job.status === 'failed') return 'Failed ✗'
  if (job.status) return job.status.charAt(0).toUpperCase() + job.status.slice(1)
  return 'Queued'
}

function goToIngestion() {
  router.push({ name: 'admin.ingestion' })
}

function resetAll() {
  batchId.value = ''
  batchTotal.value = 0
  batchProgressPct.value = 0
  batchJobs.value = []
  importing.value = false
  discography.value = null
  searched.value = false
  artistName.value = ''
  selectedTracks.value = new Set()
}

function formatDuration(seconds: number): string {
  if (!seconds || seconds <= 0) return '—'
  const m = Math.floor(seconds / 60)
  const s = seconds % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}
</script>
