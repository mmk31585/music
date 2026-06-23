<template>
  <div class="mx-auto w-full max-w-2xl px-4 pt-6 pb-36 md:px-6" dir="rtl">
    <!-- Back -->
    <button
      type="button"
      class="mb-6 inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-white/50 transition hover:bg-white/[0.06] hover:text-white"
      @click="goBack"
    >
      <i aria-hidden="true" class="pi pi-arrow-right text-xs" />
      بازگشت
    </button>

    <!-- Step indicator -->
    <div class="mb-8 flex items-center gap-2">
      <div
        v-for="s in steps"
        :key="s.step"
        class="flex items-center gap-2"
      >
        <div
          class="flex h-7 w-7 items-center justify-center rounded-full text-xs font-bold transition"
          :class="currentStep === s.step
            ? 'bg-[#1db954] text-black'
            : currentStep > s.step
              ? 'bg-[#1db954]/30 text-[#1db954]'
              : 'bg-white/10 text-white/40'"
        >
          <i v-if="currentStep > s.step" aria-hidden="true" class="pi pi-check text-[10px]" />
          <span v-else>{{ s.step }}</span>
        </div>
        <span class="text-xs font-medium" :class="currentStep === s.step ? 'text-white' : 'text-white/40'">
          {{ s.label }}
        </span>
        <i v-if="s.step < 3" aria-hidden="true" class="pi pi-chevron-left text-[10px] text-white/20" />
      </div>
    </div>

    <!-- Step 1: Pick track -->
    <div v-if="currentStep === 1" class="space-y-6">
      <h1 class="text-2xl font-black text-white">ادیت برای کدوم آهنگ؟</h1>

      <!-- Search input -->
      <div class="relative">
        <i aria-hidden="true" class="pi pi-search absolute right-3 top-1/2 -translate-y-1/2 text-sm text-slate-400" />
        <input
          ref="searchInputRef"
          v-model="searchQuery"
          type="text"
          placeholder="جستجوی آهنگ..."
          class="w-full rounded-2xl border border-white/10 bg-white/[0.04] px-10 py-3.5 text-sm text-white outline-none backdrop-blur-sm transition placeholder:text-slate-500 focus:border-[#1db954]/40 focus:bg-white/[0.06]"
          @input="onSearchInput"
          @keydown="onSearchKeydown"
        />
        <i
          v-if="searching"
          aria-hidden="true"
          class="pi pi-spin pi-spinner absolute left-3 top-1/2 -translate-y-1/2 text-xs text-slate-400"
        />
      </div>

      <!-- Search results -->
      <div class="space-y-1">
        <div v-if="!searchQuery" class="flex items-center justify-center py-16 text-sm text-white/30">
          <i aria-hidden="true" class="pi pi-headphones ml-2" /> برای جستجو تایپ کن
        </div>

        <div v-else-if="searching" class="flex items-center justify-center py-16 text-sm text-white/40">
          <i aria-hidden="true" class="pi pi-spin pi-spinner ml-2" /> در حال جستجو...
        </div>

        <div v-else-if="!searchResults.length" class="flex items-center justify-center py-16 text-sm text-white/30">
          <i aria-hidden="true" class="pi pi-info-circle ml-2" /> آهنگی یافت نشد
        </div>

        <div v-else class="space-y-1">
          <button
            v-for="(track, i) in searchResults"
            :key="String(track.id)"
            type="button"
            class="flex w-full items-center gap-3 rounded-2xl px-4 py-3 text-right transition"
            :class="focusedIdx === i ? 'bg-white/[0.12] ring-1 ring-white/20' : 'hover:bg-white/[0.06]'"
            @click="selectTrack(track)"
            @mouseenter="focusedIdx = i"
          >
            <div class="h-12 w-12 shrink-0 overflow-hidden rounded-xl bg-white/10">
              <img
                v-if="track.cover_url"
                :src="track.cover_url"
                :alt="track.title"
                class="h-full w-full object-cover"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-music text-sm text-slate-500" />
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold text-white">{{ track.title }}</p>
              <p class="mt-0.5 truncate text-xs text-white/40">{{ track.artist_name || '' }}</p>
            </div>
            <span class="text-xs text-white/30">{{ fmtDuration(track.duration_seconds || 0) }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Step 2: Upload video + set start time -->
    <div v-if="currentStep === 2" class="space-y-6">
      <h1 class="text-2xl font-black text-white">ویدیوت رو آپلود کن</h1>

      <!-- Selected track info -->
      <div class="flex items-center gap-3 rounded-2xl border border-white/[0.06] bg-white/[0.02] px-4 py-3">
        <div class="h-11 w-11 shrink-0 overflow-hidden rounded-xl bg-white/10">
          <img
            v-if="selectedTrack?.cover_url"
            :src="selectedTrack.cover_url"
            :alt="selectedTrack.title"
            class="h-full w-full object-cover"
          />
          <div v-else class="flex h-full items-center justify-center">
            <i aria-hidden="true" class="pi pi-music text-sm text-slate-500" />
          </div>
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-semibold text-white">{{ selectedTrack?.title }}</p>
          <p class="truncate text-xs text-white/40">{{ selectedTrack?.artist_name }}</p>
        </div>
        <button
          type="button"
          class="text-xs font-medium text-[#1db954] transition hover:text-[#1ed760]"
          @click="currentStep = 1"
        >
          تغییر
        </button>
      </div>

      <!-- File upload area -->
      <div
        class="flex cursor-pointer flex-col items-center justify-center gap-3 rounded-2xl border-2 border-dashed px-6 py-12 transition"
        :class="videoFile
          ? 'border-[#1db954]/40 bg-[#1db954]/5'
          : 'border-white/10 bg-white/[0.02] hover:border-white/20 hover:bg-white/[0.04]'"
        @click="triggerFileInput"
        @dragover.prevent="dragOver = true"
        @dragleave.prevent="dragOver = false"
        @drop.prevent="onDrop"
      >
        <div v-if="!videoFile" class="flex flex-col items-center gap-2">
          <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/5">
            <i aria-hidden="true" class="pi pi-video text-2xl text-slate-400" />
          </div>
          <p class="text-sm font-medium text-white/60">+ انتخاب از گالری</p>
          <p class="text-xs text-white/30">MP4, WebM, MOV</p>
        </div>
        <div v-else class="flex flex-col items-center gap-2">
          <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-[#1db954]/10">
            <i aria-hidden="true" class="pi pi-check-circle text-2xl text-[#1db954]" />
          </div>
          <p class="text-sm font-medium text-white">{{ videoFile.name }}</p>
          <p class="text-xs text-white/40">{{ fmtFileSize(videoFile.size) }}</p>
          <button
            type="button"
            class="text-xs font-medium text-slate-400 transition hover:text-white"
            @click.stop="removeVideo"
          >
            حذف
          </button>
        </div>
        <input
          ref="fileInputRef"
          type="file"
          accept="video/*"
          class="hidden"
          @change="onFileChange"
        />
      </div>

      <!-- Seek bar for start time -->
      <div v-if="selectedTrack" class="space-y-3">
        <p class="text-sm font-medium text-white/70">از کجای آهنگ شروع بشه؟</p>
        <div class="flex items-center gap-3">
          <button
            type="button"
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/20"
            @click="toggleAudioPreview"
          >
            <i aria-hidden="true" :class="isPreviewPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-xs" />
          </button>
          <div class="flex-1">
            <input
              type="range"
              min="0"
              :max="maxSeekMs"
              step="100"
              class="w-full accent-[#1db954]"
              :value="trackStartMs"
              @input="onSeekChange"
            />
          </div>
          <span class="min-w-[3rem] text-right text-xs font-medium text-white/50 tabular-nums">
            {{ fmtMs(trackStartMs) }}
          </span>
        </div>
      </div>

      <!-- Audio element for preview (hidden) -->
      <audio
        ref="audioPreviewRef"
        class="hidden"
        :controls="false"
        @timeupdate="onAudioTimeUpdate"
        @ended="isPreviewPlaying = false"
      />
    </div>

    <!-- Step 3: Caption + publish -->
    <div v-if="currentStep === 3" class="space-y-6">
      <h1 class="text-2xl font-black text-white">جزئیات ادیت</h1>

      <!-- Selected track info -->
      <div class="flex items-center gap-3 rounded-2xl border border-white/[0.06] bg-white/[0.02] px-4 py-3">
        <div class="h-11 w-11 shrink-0 overflow-hidden rounded-xl bg-white/10">
          <img
            v-if="selectedTrack?.cover_url"
            :src="selectedTrack.cover_url"
            :alt="selectedTrack.title"
            class="h-full w-full object-cover"
          />
          <div v-else class="flex h-full items-center justify-center">
            <i aria-hidden="true" class="pi pi-music text-sm text-slate-500" />
          </div>
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-semibold text-white">{{ selectedTrack?.title }}</p>
          <p class="truncate text-xs text-white/40">{{ selectedTrack?.artist_name }}</p>
        </div>
        <div class="text-xs text-white/40">
          {{ fmtMs(trackStartMs) }}
        </div>
      </div>

      <!-- Title -->
      <div class="space-y-1.5">
        <label class="text-xs font-medium text-white/50">عنوان (اختیاری)</label>
        <input
          v-model="editTitle"
          type="text"
          maxlength="100"
          placeholder="نام ادیتت رو بنویس..."
          class="w-full rounded-2xl border border-white/10 bg-white/[0.04] px-4 py-3 text-sm text-white outline-none backdrop-blur-sm transition placeholder:text-slate-500 focus:border-[#1db954]/40 focus:bg-white/[0.06]"
        />
        <p class="text-left text-[10px] text-white/30">{{ editTitle.length }}/100</p>
      </div>

      <!-- Description -->
      <div class="space-y-1.5">
        <label class="text-xs font-medium text-white/50">توضیحات (اختیاری)</label>
        <textarea
          v-model="editDescription"
          maxlength="300"
          rows="3"
          placeholder="یه توضیح کوتاه..."
          class="w-full resize-none rounded-2xl border border-white/10 bg-white/[0.04] px-4 py-3 text-sm text-white outline-none backdrop-blur-sm transition placeholder:text-slate-500 focus:border-[#1db954]/40 focus:bg-white/[0.06]"
        />
        <p class="text-left text-[10px] text-white/30">{{ editDescription.length }}/300</p>
      </div>

      <!-- Publish button -->
      <button
        type="button"
        class="mt-4 flex w-full items-center justify-center gap-2 rounded-2xl bg-[#1db954] py-3.5 text-sm font-bold text-black transition hover:bg-[#1ed760] disabled:opacity-40"
        :disabled="publishing"
        @click="publish"
      >
        <i v-if="publishing" aria-hidden="true" class="pi pi-spin pi-spinner" />
        {{ publishing ? 'در حال انتشار...' : 'انتشار' }}
      </button>
    </div>

    <!-- Processing state overlay -->
    <Transition name="fade">
      <div
        v-if="processingState"
        class="fixed inset-0 z-[9999] flex flex-col items-center justify-center gap-4 bg-black/80 backdrop-blur-sm"
      >
        <div v-if="processingState === 'processing'" class="flex flex-col items-center gap-4">
          <i aria-hidden="true" class="pi pi-spin pi-spinner text-3xl text-[#1db954]" />
          <p class="text-sm font-medium text-white/70">ویدیوت داره پردازش می‌شه... ممکنه چند دقیقه طول بکشه</p>
        </div>
        <div v-else-if="processingState === 'success'" class="flex flex-col items-center gap-4">
          <div class="flex h-16 w-16 items-center justify-center rounded-full bg-[#1db954]/20">
            <i aria-hidden="true" class="pi pi-check-circle text-3xl text-[#1db954]" />
          </div>
          <p class="text-sm font-bold text-white">ادیت با موفقیت منتشر شد!</p>
          <p class="text-xs text-white/50"> redirecting...</p>
        </div>
        <div v-else-if="processingState === 'timeout'" class="flex flex-col items-center gap-4">
          <div class="flex h-16 w-16 items-center justify-center rounded-full bg-amber-400/20">
            <i aria-hidden="true" class="pi pi-clock text-3xl text-amber-400" />
          </div>
          <p class="text-sm font-medium text-white/70">پردازش طولانی شد. ادیتت تو پروفایلت نمایش داده می‌شه وقتی آماده بشه</p>
          <button
            type="button"
            class="rounded-full bg-white/10 px-6 py-2 text-sm font-bold text-white transition hover:bg-white/20"
            @click="goToProfile"
          >
            رفتن به پروفایل
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSearchApi } from '@/services/api/catalog/search'
import { useVideoApi } from '@/services/api/video'
import { useTracksApi } from '@/services/api/catalog/tracks'
import { useToast } from 'primevue/usetoast'
import type { Track } from '@/services/api/catalog/tracks/types'

const route = useRoute()
const router = useRouter()
const searchApi = useSearchApi()
const tracksApi = useTracksApi()
const videoApi = useVideoApi()
const toast = useToast()

// ── State ─────────────────────────────────────────────────────────────

const currentStep = ref(1)

const steps = [
  { step: 1, label: 'انتخاب آهنگ' },
  { step: 2, label: 'آپلود ویدیو' },
  { step: 3, label: 'انتشار' },
]

// Pre-selected track from query param
const preSelectedTrackId = route.query.trackId as string | undefined

// Search
const searchQuery = ref('')
const searchResults = ref<Track[]>([])
const searching = ref(false)
let debounceTimer: ReturnType<typeof setTimeout> | null = null
let abortController: AbortController | null = null
const searchInputRef = ref<HTMLInputElement | null>(null)
const focusedIdx = ref(-1)

// Selected track
const selectedTrack = ref<Track | null>(null)

// Video file
const fileInputRef = ref<HTMLInputElement | null>(null)
const videoFile = ref<File | null>(null)
const dragOver = ref(false)

// Audio preview
const audioPreviewRef = ref<HTMLAudioElement | null>(null)
const isPreviewPlaying = ref(false)
const trackStartMs = ref(0)
const maxSeekMs = ref(0)

// Metadata
const editTitle = ref('')
const editDescription = ref('')

// Publishing
const publishing = ref(false)
const processingState = ref<'processing' | 'success' | 'timeout' | null>(null)

// ── Lifecycle ─────────────────────────────────────────────────────────

onMounted(() => {
    if (preSelectedTrackId) {
    fetchTrackById(preSelectedTrackId)
  }
})

onUnmounted(() => {
  abortController?.abort()
  audioPreviewRef.value?.pause()
})

async function fetchTrackById(trackId: string) {
  try {
    const track = await tracksApi.getTrack(trackId)
    if (track) selectTrack(track as unknown as Track)
  } catch {
    // Track not found — user can search manually
  }
}

// ── Navigation ───────────────────────────────────────────────────────

function goBack() {
  if (processingState.value) return
  if (currentStep.value > 1) {
    currentStep.value--
  } else {
    router.back()
  }
}

// ── Step 1: Search ───────────────────────────────────────────────────

function onSearchInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(doSearch, 300)
  focusedIdx.value = -1
}

async function doSearch() {
  const term = searchQuery.value.trim()
  if (!term) {
    searchResults.value = []
    searching.value = false
    return
  }

  abortController?.abort()
  abortController = new AbortController()

  searching.value = true
  try {
    const res = await searchApi.searchCatalog(
      { query: term, type: 'tracks', limit: 10 },
      { signal: abortController.signal } as Record<string, unknown>,
    )
    searchResults.value = res.tracks ?? []
  } catch (err) {
    if ((err as Record<string, unknown>)?.name === 'AbortError') return
    searchResults.value = []
  } finally {
    searching.value = false
  }
}

function onSearchKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (focusedIdx.value < searchResults.value.length - 1) focusedIdx.value++
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (focusedIdx.value > 0) focusedIdx.value--
  } else if (e.key === 'Enter' && focusedIdx.value >= 0) {
    e.preventDefault()
    const t = searchResults.value[focusedIdx.value]
    if (t) selectTrack(t)
  }
}

function selectTrack(track: Track) {
  selectedTrack.value = track
  searchResults.value = []
  searchQuery.value = ''
  maxSeekMs.value = (track.duration_seconds || 0) * 1000
  trackStartMs.value = 0
  currentStep.value = 2
}

function fmtDuration(s: number) {
  const m = Math.floor(s / 60)
  const sec = s % 60
  return `${m}:${String(sec).padStart(2, '0')}`
}

// ── Step 2: Upload video ─────────────────────────────────────────────

function triggerFileInput() {
  fileInputRef.value?.click()
}

function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    videoFile.value = file
  }
}

function onDrop(e: DragEvent) {
  dragOver.value = false
  const file = e.dataTransfer?.files?.[0]
  if (file) {
    videoFile.value = file
  }
}

function removeVideo() {
  videoFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
}

function fmtFileSize(bytes: number) {
  if (bytes >= 1_000_000) return `${(bytes / 1_000_000).toFixed(1)} MB`
  if (bytes >= 1_000) return `${(bytes / 1_000).toFixed(1)} KB`
  return `${bytes} B`
}

function onSeekChange(e: Event) {
  const val = parseInt((e.target as HTMLInputElement).value, 10)
  trackStartMs.value = val

  // Preview audio from this position
  const audio = audioPreviewRef.value
  if (audio && selectedTrack.value) {
    const wasPlaying = !audio.paused
    audio.currentTime = val / 1000
    if (wasPlaying) audio.play()
  }
}

function toggleAudioPreview() {
  const audio = audioPreviewRef.value
  if (!audio || !selectedTrack.value) return

  if (isPreviewPlaying.value) {
    audio.pause()
    isPreviewPlaying.value = false
  } else {
    if (!audio.src) {
      audio.src = selectedTrack.value.audio_url || ''
    }
    audio.currentTime = trackStartMs.value / 1000
    audio.play().then(() => {
      isPreviewPlaying.value = true
    }).catch(() => {
      // Autoplay may be blocked
    })
  }
}

function onAudioTimeUpdate() {
  // Keep duration sync
}

// ── Step 3: Publish ──────────────────────────────────────────────────

async function publish() {
  if (!selectedTrack.value || !videoFile.value) return
  publishing.value = true

  try {
    const result = await videoApi.uploadVideoFile(
      videoFile.value,
      {
        track_id: selectedTrack.value.id,
        title: editTitle.value || undefined,
        description: editDescription.value || undefined,
        track_start_ms: trackStartMs.value || undefined,
      },
    )

    // Transition to processing overlay (no flash)
    processingState.value = 'processing'
    await nextTick()

    // Poll for job completion
    const jobId = result?.job_id
    if (jobId) {
      await pollJob(jobId)
    } else {
      // No job ID — assume immediate processing, show success
      processingState.value = 'success'
      await sleep(1500)
      goToProfile()
    }
  } catch (err: any) {
    const msg = err?.message || 'خطا در انتشار ادیت'
    toast.add({ severity: 'error', summary: 'خطا', detail: msg, life: 5000 })
  } finally {
    publishing.value = false
  }
}

async function pollJob(jobId: string) {
  const MAX_ATTEMPTS = 10
  const INTERVAL_MS = 5000

  for (let i = 0; i < MAX_ATTEMPTS; i++) {
    await sleep(INTERVAL_MS)
    try {
      const status = await videoApi.getVideoJobStatus(jobId)
      const jobStatus = status?.status as string
      if (jobStatus === 'completed' || jobStatus === 'success') {
        processingState.value = 'success'
        await sleep(1500)
        goToProfile()
        return
      }
      if (status?.status === 'failed') {
        processingState.value = 'timeout'
        return
      }
      // Continue polling for 'pending' or 'processing'
    } catch {
      // Network error — continue polling
    }
  }

  // Timeout exceeded — show timeout state
  processingState.value = 'timeout'
}

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function fmtMs(ms: number): string {
  const totalSec = Math.floor(ms / 1000)
  const m = Math.floor(totalSec / 60)
  const s = totalSec % 60
  return `${m}:${String(s).padStart(2, '0')}`
}

function goToProfile() {
  router.push('/profile')
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 300ms ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
