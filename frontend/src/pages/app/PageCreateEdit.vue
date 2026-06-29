<template>
  <div class="mx-auto w-full max-w-2xl px-4 pt-6 pb-36 md:px-6" :dir="isRTL ? 'rtl' : 'ltr'">
    <!-- Back -->
    <button
      type="button"
      class="mb-6 inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-white/50 transition hover:bg-white/6 hover:text-white focus-visible:outline-2 focus-visible:outline-spotify"
      @click="goBack"
    >
      <i aria-hidden="true" class="pi pi-arrow-right text-xs" />
      {{ t('back') }}
    </button>

    <!-- Step progress indicator -->
    <div class="mb-8 flex items-center gap-2" role="tablist" :aria-label="t('steps')">
      <div
        v-for="(s, i) in steps"
        :key="s.step"
        class="flex items-center gap-2"
        role="tab"
        :aria-selected="currentStep === s.step"
      >
        <div
          class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-bold transition-all duration-300"
          :class="stepClass(i)"
        >
          <i v-if="isStepComplete(i)" aria-hidden="true" class="pi pi-check text-xs" />
          <span v-else>{{ i + 1 }}</span>
        </div>
        <span
          class="hidden text-xs font-medium sm:block"
          :class="currentStep >= s.step ? 'text-white/70' : 'text-white/20'"
        >
          {{ s.label }}
        </span>
        <i
          v-if="i < steps.length - 1"
          aria-hidden="true"
          class="pi pi-chevron-left text-[10px] text-white/8"
        />
      </div>
    </div>

    <!-- ═══ Step 1: Pick track ═══ -->
    <div v-if="currentStep === 1" class="space-y-6" role="tabpanel">
      <h1 class="text-2xl font-black text-white">{{ t('pickTrack') }}</h1>

      <!-- Search input -->
      <div class="relative">
        <i
          aria-hidden="true"
          class="pi pi-search absolute right-3 top-1/2 -translate-y-1/2 text-sm text-slate-400"
        />
        <input
          ref="searchInputRef"
          v-model="searchQuery"
          type="text"
          :placeholder="t('searchPlaceholder')"
          class="h-11 w-full rounded-2xl border border-white/10 bg-white/4 px-10 text-sm text-white outline-hidden backdrop-blur-xs transition placeholder:text-slate-500 focus:border-spotify/40 focus:bg-white/6"
          aria-label="جستجوی آهنگ"
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
      <div class="space-y-1" role="listbox" aria-label="نتایج جستجو">
        <div
          v-if="!searchQuery"
          class="flex items-center justify-center py-16 text-sm text-white/30"
        >
          <i aria-hidden="true" class="pi pi-headphones ml-2" />
          {{ t('typeToSearch') }}
        </div>

        <div
          v-else-if="searching"
          class="flex items-center justify-center py-16 text-sm text-white/40"
        >
          <i aria-hidden="true" class="pi pi-spin pi-spinner ml-2" />
          {{ t('searching') }}
        </div>

        <div
          v-else-if="!searchResults.length"
          class="flex items-center justify-center py-16 text-sm text-white/30"
        >
          <i aria-hidden="true" class="pi pi-info-circle ml-2" />
          {{ t('noResults') }}
        </div>

        <TransitionGroup v-else name="track-list" tag="div" class="space-y-1">
          <button
            v-for="(track, i) in searchResults"
            :key="String(track.id)"
            type="button"
            role="option"
            :aria-selected="focusedIdx === i"
            class="flex w-full items-center gap-3 rounded-2xl px-4 py-3 text-right transition"
            :class="[
              focusedIdx === i
                ? 'bg-white/12 ring-1 ring-white/20'
                : 'hover:bg-white/6',
              selectedTrack?.id === track.id ? 'ring-1 ring-spotify/30 bg-spotify/5' : '',
            ]"
            @click="selectTrack(track)"
            @mouseenter="focusedIdx = i"
          >
            <div class="h-12 w-12 shrink-0 overflow-hidden rounded-xl bg-white/10">
              <img
                v-if="track.cover_url"
                :src="track.cover_url"
                :alt="track.title"
                class="h-full w-full object-cover"
                loading="lazy"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-music text-sm text-slate-500" />
              </div>
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-semibold text-white">{{ track.title }}</p>
              <p class="mt-0.5 truncate text-xs text-white/40">{{ track.artist_name || '' }}</p>
            </div>
            <span class="text-xs text-white/30 tabular-nums">{{ fmtDuration(track.duration_seconds || 0) }}</span>
          </button>
        </TransitionGroup>
      </div>
    </div>

    <!-- ═══ Step 2: Upload video ═══ -->
    <div v-if="currentStep === 2" class="space-y-6" role="tabpanel">
      <h1 class="text-2xl font-black text-white">{{ t('uploadVideo') }}</h1>

      <!-- Selected track info card -->
      <div class="flex items-center gap-3 rounded-2xl border border-white/8 bg-white/3 px-4 py-3 backdrop-blur-xs">
        <div class="h-12 w-12 shrink-0 overflow-hidden rounded-xl bg-white/10 shadow-lg">
          <img
            v-if="selectedTrack?.cover_url"
            :src="selectedTrack.cover_url"
            :alt="selectedTrack.title"
            class="h-full w-full object-cover"
            loading="lazy"
          />
          <div v-else class="flex h-full items-center justify-center">
            <i aria-hidden="true" class="pi pi-music text-sm text-slate-500" />
          </div>
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-bold text-white">{{ selectedTrack?.title }}</p>
          <p class="truncate text-xs text-white/40">{{ selectedTrack?.artist_name }}</p>
          <p v-if="selectedTrack?.duration_seconds" class="text-[10px] text-white/30">
            {{ fmtDuration(selectedTrack.duration_seconds) }}
          </p>
        </div>
        <button
          type="button"
          class="shrink-0 rounded-lg bg-white/8 px-3 py-1.5 text-xs font-medium text-spotify transition hover:bg-white/12 focus-visible:outline-2 focus-visible:outline-spotify"
          @click="currentStep = 1"
        >
          {{ t('change') }}
        </button>
      </div>

      <!-- File upload dropzone -->
      <div
        class="flex cursor-pointer flex-col items-center justify-center gap-3 rounded-2xl border-2 border-dashed px-6 py-12 transition"
        :class="videoFile
          ? 'border-spotify/40 bg-spotify/5'
          : 'border-white/10 bg-white/2 hover:border-spotify/30 hover:bg-white/4'"
        role="button"
        :tabindex="0"
        :aria-label="t('uploadFile')"
        @click="triggerFileInput"
        @keydown.enter="triggerFileInput"
        @keydown.space.prevent="triggerFileInput"
        @drop.prevent="onDrop"
      >
        <!-- No file state -->
        <template v-if="!videoFile">
          <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/5">
            <i aria-hidden="true" class="pi pi-video text-2xl text-slate-400" />
          </div>
          <p class="text-sm font-medium text-white/60">{{ t('selectVideo') }}</p>
          <p class="text-xs text-white/30">MP4, WebM, MOV — {{ t('maxSize') }} 200MB</p>
        </template>

        <!-- File selected state -->
        <template v-else>
          <!-- Video preview player -->
          <div class="relative w-full overflow-hidden rounded-2xl bg-black/60 shadow-2xl ring-1 ring-white/6">
            <video
              ref="videoPreviewRef"
              :src="videoPreviewUrl"
              class="w-full max-h-80 object-contain"
              :aria-label="t('videoPreview')"
              controls
              playsinline
              preload="metadata"
            />
          </div>

          <!-- File info bar -->
          <div class="flex w-full items-center gap-3 rounded-xl border border-white/6 bg-white/3 px-4 py-3 backdrop-blur-xs">
            <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-spotify/15">
              <i aria-hidden="true" class="pi pi-video text-sm text-spotify" />
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-white">{{ videoFile.name }}</p>
              <p class="text-xs text-white/40">{{ fmtFileSize(videoFile.size) }}</p>
            </div>
            <button
              type="button"
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-white/8 text-slate-400 transition hover:bg-red-500/20 hover:text-red-400 focus-visible:outline-2 focus-visible:outline-white"
              :aria-label="t('remove')"
              @click.stop="removeVideo"
            >
              <i aria-hidden="true" class="pi pi-trash text-xs" />
            </button>
          </div>

          <!-- Upload progress bar -->
          <div v-if="uploadProgress > 0 && uploadProgress < 100" class="h-1.5 w-full overflow-hidden rounded-full bg-white/8">
            <div
              class="h-full rounded-full bg-spotify transition-all duration-500 ease-out"
              :style="{ width: uploadProgress + '%' }"
            />
          </div>
        </template>

        <input
          ref="fileInputRef"
          type="file"
          accept="video/mp4,video/webm,video/quicktime"
          class="hidden"
          aria-hidden="true"
          @change="onFileChange"
        />
      </div>

      <!-- Track info card — just shows the selected track, no sync -->
      <div v-if="selectedTrack" class="space-y-3 rounded-2xl border border-white/6 bg-white/2 px-4 py-4 backdrop-blur-xs">
        <p class="text-xs font-medium text-white/50">{{ t('editWillUse') }}</p>
        <div class="flex items-center gap-3">
          <div class="h-10 w-10 shrink-0 overflow-hidden rounded-xl bg-white/10">
            <img
              v-if="selectedTrack.cover_url"
              :src="selectedTrack.cover_url"
              :alt="selectedTrack.title"
              class="h-full w-full object-cover"
            />
            <div v-else class="flex h-full items-center justify-center">
              <i aria-hidden="true" class="pi pi-music text-sm text-slate-500" />
            </div>
          </div>
          <div class="min-w-0 flex-1">
            <p class="truncate text-sm font-semibold text-white">{{ selectedTrack.title }}</p>
            <p class="truncate text-xs text-white/40">{{ selectedTrack.artist_name }}</p>
          </div>
          <span v-if="selectedTrack.duration_seconds" class="text-xs text-white/30 tabular-nums">
            {{ fmtDuration(selectedTrack.duration_seconds) }}
          </span>
        </div>
      </div>

      <!-- Continue button -->
      <button
        type="button"
        class="w-full rounded-2xl bg-spotify py-3.5 text-sm font-bold text-black transition hover:bg-spotify-hover disabled:opacity-40 focus-visible:outline-2 focus-visible:outline-white"
        :disabled="!videoFile"
        @click="currentStep = 3"
      >
        {{ t('continue') }}
      </button>
    </div>

    <!-- ═══ Step 3: Caption + publish ═══ -->
    <div v-if="currentStep === 3" class="space-y-6" role="tabpanel">
      <h1 class="text-2xl font-black text-white">{{ t('editDetails') }}</h1>

      <!-- Selected track summary -->
      <div class="flex items-center gap-3 rounded-2xl border border-white/8 bg-white/3 px-4 py-3 backdrop-blur-xs">
        <div class="h-11 w-11 shrink-0 overflow-hidden rounded-xl bg-white/10">
          <img
            v-if="selectedTrack?.cover_url"
            :src="selectedTrack.cover_url"
            :alt="selectedTrack.title"
            class="h-full w-full object-cover"
            loading="lazy"
          />
          <div v-else class="flex h-full items-center justify-center">
            <i aria-hidden="true" class="pi pi-music text-sm text-slate-500" />
          </div>
        </div>
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm font-semibold text-white">{{ selectedTrack?.title }}</p>
          <p class="truncate text-xs text-white/40">{{ selectedTrack?.artist_name }}</p>
        </div>
      </div>

      <!-- Title (required) -->
      <div class="space-y-1.5">
        <label for="edit-title" class="text-xs font-medium text-white/50">{{ t('title') }} *</label>
        <input
          id="edit-title"
          v-model="editTitle"
          type="text"
          maxlength="100"
          required
          :placeholder="t('titlePlaceholder')"
          class="w-full rounded-2xl border bg-white/4 px-4 py-3 text-sm text-white outline-hidden backdrop-blur-xs transition placeholder:text-slate-500 focus:bg-white/6"
          :class="titleError ? 'border-red-400/40 focus:border-red-400/60' : 'border-white/10 focus:border-spotify/40'"
          aria-describedby="edit-title-char-count edit-title-error"
          @blur="titleTouched = true"
          @input="titleTouched = true"
        />
        <div class="flex items-center justify-between">
          <p v-if="titleError" id="edit-title-error" class="text-[10px] text-red-400" role="alert">
            {{ titleError }}
          </p>
          <p id="edit-title-char-count" class="text-[10px] text-white/30 tabular-nums">{{ editTitle.length }}/100</p>
        </div>
      </div>

      <!-- Description (optional) -->
      <div class="space-y-1.5">
        <label for="edit-desc" class="text-xs font-medium text-white/50">{{ t('description') }}</label>
        <textarea
          id="edit-desc"
          v-model="editDescription"
          maxlength="300"
          rows="3"
          :placeholder="t('descPlaceholder')"
          class="w-full resize-none rounded-2xl border border-white/10 bg-white/4 px-4 py-3 text-sm text-white outline-hidden backdrop-blur-xs transition placeholder:text-slate-500 focus:border-spotify/40 focus:bg-white/6"
        />
        <p class="text-left text-[10px] text-white/30 tabular-nums">{{ editDescription.length }}/300</p>
      </div>

      <!-- Publish button -->
      <button
        type="button"
        class="mt-4 flex w-full items-center justify-center gap-2 rounded-2xl bg-spotify py-3.5 text-sm font-bold text-black transition hover:bg-spotify-hover disabled:opacity-40 focus-visible:outline-2 focus-visible:outline-white"
        :disabled="publishing || !isPublishValid"
        @click="publish"
      >
        <i v-if="publishing" aria-hidden="true" class="pi pi-spin pi-spinner" />
        {{ publishing ? t('publishing') : t('publish') }}
      </button>

      <!-- Back to step 2 -->
      <button
        type="button"
        class="w-full text-center text-xs text-white/30 transition hover:text-white/60"
        @click="currentStep = 2"
      >
        {{ t('backToUpload') }}
      </button>
    </div>

    <!-- ═══ Processing overlay ═══ -->
    <Transition name="fade">
      <div
        v-if="processingState"
        class="fixed inset-0 z-[9999] flex flex-col items-center justify-center gap-4 bg-black/80 backdrop-blur-xs"
        role="alertdialog"
        :aria-label="t('processing')"
      >
        <!-- Processing -->
        <template v-if="processingState === 'processing'">
          <div class="flex flex-col items-center gap-4">
            <div class="flex h-16 w-16 items-center justify-center rounded-full bg-spotify/15">
              <i aria-hidden="true" class="pi pi-spin pi-spinner text-3xl text-spotify" />
            </div>
            <p class="text-sm font-medium text-white/70">{{ t('processingMessage') }}</p>
            <div class="h-1 w-48 overflow-hidden rounded-full bg-white/10">
              <div class="h-full animate-pulse rounded-full bg-spotify/60" style="width: 60%" />
            </div>
          </div>
        </template>

        <!-- Success -->
        <template v-else-if="processingState === 'success'">
          <div class="flex flex-col items-center gap-4">
            <div class="flex h-16 w-16 items-center justify-center rounded-full bg-spotify/20">
              <i aria-hidden="true" class="pi pi-check-circle text-3xl text-spotify" />
            </div>
            <p class="text-sm font-bold text-white">{{ t('successMessage') }}</p>
            <p class="text-xs text-white/50">{{ t('redirecting') }}</p>
          </div>
        </template>

        <!-- Timeout -->
        <template v-else-if="processingState === 'timeout'">
          <div class="flex flex-col items-center gap-4">
            <div class="flex h-16 w-16 items-center justify-center rounded-full bg-amber-400/20">
              <i aria-hidden="true" class="pi pi-clock text-3xl text-amber-400" />
            </div>
            <p class="text-sm font-medium text-white/70">{{ t('timeoutMessage') }}</p>
            <button
              type="button"
              class="rounded-full bg-white/10 px-6 py-2 text-sm font-bold text-white transition hover:bg-white/20 focus-visible:outline-2 focus-visible:outline-white"
              @click="goToProfile"
            >
              {{ t('goToProfile') }}
            </button>
          </div>
        </template>

        <!-- Error -->
        <template v-else-if="processingState === 'error'">
          <div class="flex flex-col items-center gap-4">
            <div class="flex h-16 w-16 items-center justify-center rounded-full bg-red-500/20">
              <i aria-hidden="true" class="pi pi-exclamation-circle text-3xl text-red-400" />
            </div>
            <p class="text-sm font-medium text-white/70">{{ t('errorMessage') }}</p>
            <button
              type="button"
              class="rounded-full bg-spotify px-6 py-2 text-sm font-bold text-black transition hover:bg-spotify-hover focus-visible:outline-2 focus-visible:outline-white"
              @click="resetAfterError"
            >
              {{ t('tryAgain') }}
            </button>
          </div>
        </template>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
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

// ── RTL / i18n helpers ─────────────────────────────────────────────
const isRTL = computed(() => true) // Persian-only page

function t(key: string): string {
  const msgs: Record<string, string> = {
    back: 'بازگشت',
    steps: 'مراحل',
    pickTrack: 'ادیت برای کدوم آهنگ؟',
    searchPlaceholder: 'جستجوی آهنگ...',
    typeToSearch: 'برای جستجو تایپ کن',
    searching: 'در حال جستجو...',
    noResults: 'آهنگی یافت نشد',
    uploadVideo: 'ویدیوت رو آپلود کن',
    uploadFile: 'انتخاب فایل ویدیو',
    selectVideo: 'انتخاب از گالری',
    maxSize: 'حداکثر ۲۰۰ مگابایت',
    videoPreview: 'پیش‌نمایش ویدیو',
    remove: 'حذف',
    editWillUse: 'این ویدیو به‌عنوان ادیت برای این آهنگ منتشر می‌شود',
    continue: 'ادامه',
    change: 'تغییر',
    editDetails: 'جزئیات ادیت',
    title: 'عنوان (اختیاری)',
    titlePlaceholder: 'نام ادیتت رو بنویس...',
    description: 'توضیحات (اختیاری)',
    descPlaceholder: 'یه توضیح کوتاه...',
    publish: 'انتشار',
    publishing: 'در حال انتشار...',
    titleTooShort: 'عنوان باید حداقل ۲ حرف باشد',
    backToUpload: 'بازگشت به آپلود ویدیو',
    processing: 'در حال پردازش',
    processingMessage: 'ویدیوت داره پردازش می‌شه... ممکنه چند دقیقه طول بکشه',
    successMessage: 'ادیت با موفقیت منتشر شد!',
    redirecting: 'در حال انتقال...',
    timeoutMessage: 'پردازش طولانی شد. ادیتت تو پروفایلت نمایش داده می‌شه وقتی آماده بشه',
    goToProfile: 'رفتن به پروفایل',
    errorMessage: 'خطا در انتشار ادیت',
    tryAgain: 'تلاش مجدد',
  }
  return msgs[key] || key
}

// ── State ──────────────────────────────────────────────────────────

const currentStep = ref(1)

const steps = [
  { step: 1, label: 'انتخاب آهنگ' },
  { step: 2, label: 'آپلود ویدیو' },
  { step: 3, label: 'انتشار' },
]

function isStepComplete(i: number): boolean {
  return currentStep.value > steps[i].step
}

function stepClass(i: number): string {
  if (isStepComplete(i)) return 'bg-spotify text-black shadow-sm shadow-spotify/30'
  if (currentStep.value === steps[i].step) return 'bg-spotify/20 text-spotify ring-1 ring-spotify/30'
  return 'bg-white/10 text-white/40'
}

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
const videoPreviewRef = ref<HTMLVideoElement | null>(null)
const videoFile = ref<File | null>(null)
const videoPreviewUrl = ref('')
const uploadProgress = ref(0)

// Metadata
const editTitle = ref('')
const editDescription = ref('')

// Validation
const titleTouched = ref(false)

const isPublishValid = computed(() => {
  // Title must be at least 2 characters if provided
  if (editTitle.value && editTitle.value.trim().length < 2) return false
  // Title max 100 chars
  if (editTitle.value && editTitle.value.length > 100) return false
  // Description max 300 chars
  if (editDescription.value && editDescription.value.length > 300) return false
  return selectedTrack.value !== null && videoFile.value !== null
})

const titleError = computed(() => {
  if (!titleTouched.value || !editTitle.value) return ''
  if (editTitle.value.trim().length < 2) return t('titleTooShort')
  return ''
})

// Publishing
const publishing = ref(false)
const processingState = ref<'processing' | 'success' | 'timeout' | 'error' | null>(null)

// ── Lifecycle ──────────────────────────────────────────────────────

onMounted(() => {
  if (preSelectedTrackId) {
    fetchTrackById(preSelectedTrackId)
  }
})

onUnmounted(() => {
  abortController?.abort()
  if (videoPreviewUrl.value) URL.revokeObjectURL(videoPreviewUrl.value)
})

async function fetchTrackById(trackId: string) {
  try {
    const track = await tracksApi.getTrack(trackId)
    if (track) selectTrack(track as unknown as Track)
  } catch {
    // Track not found — user can search manually
  }
}

// ── Navigation ─────────────────────────────────────────────────────

function goBack() {
  if (processingState.value) return
  if (currentStep.value > 1) {
    currentStep.value--
  } else {
    router.back()
  }
}

// ── Step 1: Search ─────────────────────────────────────────────────

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
  editTitle.value = ''
  editDescription.value = ''
  titleTouched.value = false
  currentStep.value = 2
}

function fmtDuration(s: number) {
  const m = Math.floor(s / 60)
  const sec = s % 60
  return `${m}:${String(sec).padStart(2, '0')}`
}

// ── Step 2: Upload video ──────────────────────────────────────────

function triggerFileInput() {
  fileInputRef.value?.click()
}

function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) setVideoFile(file)
}

function onDrop(e: DragEvent) {
  dragOver.value = false
  const file = e.dataTransfer?.files?.[0]
  if (file) setVideoFile(file)
}

function setVideoFile(file: File) {
  // Validate video MIME type
  if (!file.type.startsWith('video/')) {
    toast.add({ severity: 'error', summary: 'خطا', detail: 'فایل انتخاب شده ویدیو نیست', life: 4000 })
    return
  }
  // Validate size (200MB max)
  if (file.size > 200 * 1024 * 1024) {
    toast.add({ severity: 'error', summary: 'خطا', detail: 'حجم فایل بیشتر از ۲۰۰ مگابایت است', life: 4000 })
    return
  }
  videoFile.value = file
  uploadProgress.value = 0

  // Create preview URL
  if (videoPreviewUrl.value) URL.revokeObjectURL(videoPreviewUrl.value)
  videoPreviewUrl.value = URL.createObjectURL(file)
}

function removeVideo() {
  videoFile.value = null
  if (videoPreviewUrl.value) {
    URL.revokeObjectURL(videoPreviewUrl.value)
    videoPreviewUrl.value = ''
  }
  if (fileInputRef.value) fileInputRef.value.value = ''
  uploadProgress.value = 0
}

function fmtFileSize(bytes: number) {
  if (bytes >= 1_000_000) return `${(bytes / 1_000_000).toFixed(1)} MB`
  if (bytes >= 1_000) return `${(bytes / 1_000).toFixed(1)} KB`
  return `${bytes} B`
}

// ── Step 3: Publish ───────────────────────────────────────────────

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
      },
    )

    // Transition to processing overlay
    processingState.value = 'processing'

    // Check for job_id (async processing) or video id (immediate)
    const jobId = (result as Record<string, unknown>)?.job_id as string | undefined
    const videoId = (result as Record<string, unknown>)?.id as string | undefined

    if (jobId) {
      await pollJob(jobId)
    } else {
      // No job ID — show success based on video ID
      processingState.value = 'success'
      await sleep(2000)
      if (videoId) {
        router.push(`/music-video/${videoId}`)
      } else {
        goToProfile()
      }
    }
  } catch (err: any) {
    const msg = err?.message || err?.data?.message || 'خطا در انتشار ادیت'
    toast.add({ severity: 'error', summary: 'خطا', detail: msg, life: 5000 })
    if (processingState.value === 'processing') {
      processingState.value = 'error'
    }
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
        processingState.value = 'error'
        return
      }
    } catch {
      // Network error — continue polling
    }
  }

  // Timeout exceeded
  processingState.value = 'timeout'
}

function resetAfterError() {
  processingState.value = null
}

function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

function goToProfile() {
  router.push('/profile')
}
</script>

<style scoped>
/* Track list transitions */
.track-list-enter-active {
  transition: all 0.25s ease-out;
}
.track-list-leave-active {
  transition: all 0.2s ease-in;
}
.track-list-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}
.track-list-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

/* Fade transition for overlay */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 300ms ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (prefers-reduced-motion: reduce) {
  .track-list-enter-active,
  .track-list-leave-active,
  .fade-enter-active,
  .fade-leave-active {
    transition: none;
  }
  .track-list-enter-from,
  .track-list-leave-to,
  .fade-enter-from,
  .fade-leave-to {
    opacity: 1;
    transform: none;
  }
}
</style>
