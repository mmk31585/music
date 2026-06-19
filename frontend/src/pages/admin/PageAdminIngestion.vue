<!-- TODO MEDIUM: This page uses old PrimeVue surface-* / primary-* theme classes instead of the app's admin dark theme (slate-*, white/*, emerald-*). Inconsistent styling. -->
<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Content"
      title="Music Ingestion"
      description="Upload audio files to extract metadata and create ingestion drafts."
    />

    <!-- Upload Zone -->
    <Card class="mb-8 overflow-hidden">
      <template #content>
        <div
          class="relative flex cursor-pointer flex-col items-center justify-center rounded-xl border-2 border-dashed p-10 transition-all"
          :class="isDragOver
            ? 'border-primary bg-primary/10 scale-[1.02]'
            : 'border-surface-300 dark:border-surface-600 hover:border-primary/50 hover:bg-surface-50/50 dark:hover:bg-surface-800/50'"
          @dragover.prevent="isDragOver = true"
          @dragleave.prevent="isDragOver = false"
          @drop.prevent="handleDrop"
          @click="triggerFileInput"
        >
          <div v-if="!uploading" class="flex flex-col items-center gap-3">
            <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-surface-100 dark:bg-surface-800">
              <i aria-hidden="true" class="pi pi-cloud-upload text-3xl text-surface-400"></i>
            </div>
            <div class="text-center">
              <p class="text-sm font-medium text-surface-700 dark:text-surface-300">
                Drag and drop an audio file here
              </p>
              <p class="mt-1 text-xs text-surface-500">
                or click to browse &mdash; MP3, FLAC, OGG, M4A, AAC (max 200MB)
              </p>
            </div>
            <Button label="Select File" icon="pi pi-folder-open" severity="secondary" outlined size="small" />
          </div>

          <div v-else class="flex flex-col items-center gap-3">
            <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-primary/10">
              <i aria-hidden="true" class="pi pi-spin pi-spinner text-3xl text-primary"></i>
            </div>
            <div class="text-center">
              <p class="text-sm font-medium text-surface-700 dark:text-surface-300">
                Uploading &amp; extracting metadata&hellip;
              </p>
              <p class="mt-1 text-xs text-surface-500">{{ uploadProgress }}</p>
            </div>
            <Button
              label="Cancel"
              icon="pi pi-times"
              severity="danger"
              size="small"
              text
              @click="cancelUpload"
            />
          </div>

          <input
            ref="fileInputRef"
            type="file"
            accept=".mp3,.flac,.ogg,.m4a,.aac,audio/mpeg,audio/flac,audio/ogg,audio/mp4,audio/aac"
            class="hidden"
            @change="handleFileSelect"
          />
        </div>

        <Transition name="fade">
          <div v-if="uploadError" class="mt-4">
            <Message severity="error" :closable="true" @close="uploadError = null">
              {{ uploadError }}
            </Message>
          </div>
        </Transition>
      </template>
    </Card>

    <!-- Enrichment Progress -->
    <Transition name="fade">
      <Card v-if="enrichingPoll" class="mb-8">
        <template #content>
          <div class="flex items-center justify-between gap-4">
            <div class="flex items-center gap-3">
              <i aria-hidden="true" class="pi pi-spin pi-spinner text-info"></i>
              <div>
                <p class="text-sm font-medium text-surface-700 dark:text-surface-300">
                  Enriching metadata
                </p>
                <p class="text-xs text-surface-400 mt-0.5">
                  Querying MusicBrainz, Last.fm, and Spotify&hellip;
                </p>
              </div>
            </div>
            <div class="flex items-center gap-3 shrink-0">
              <ProgressBar
                :value="enrichPercent"
                class="h-1.5 w-32 sm:w-48"
                :show-value="false"
              />
              <span class="text-xs text-surface-400 w-14 text-right tabular-nums">
                {{ enrichAttempt }}/{{ enrichMaxAttempts }}
              </span>
            </div>
          </div>
        </template>
      </Card>
    </Transition>

    <!-- Extraction Results -->
    <Transition name="fade">
      <div v-if="uploadResult" class="mb-8 space-y-4">
        <Card>
          <template #title>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <h3 class="text-lg font-semibold text-surface-900 dark:text-surface-100">
                  Extraction Results
                </h3>
                <Badge :value="uploadResult.status" :severity="statusSeverity(uploadResult.status)" />
              </div>
              <Button
                v-if="uploadResult.status === 'pending'"
                label="Enrich Now"
                icon="pi pi-magic"
                size="small"
                severity="info"
                @click="triggerEnrich(uploadResult.draftId)"
              />
            </div>
          </template>

          <template #content>
            <!-- Enrichment Suggestions -->
            <div v-if="showSuggestions" class="mb-6 rounded-lg border border-primary/20 bg-primary/5 p-5">
              <h4 class="mb-4 flex items-center gap-2 text-sm font-semibold text-primary-700 dark:text-primary-300">
                <i aria-hidden="true" class="pi pi-magic"></i>Enrichment Suggestions
              </h4>

              <div v-if="uploadResult.coverArtUrl || spotifyAlbumCover" class="mb-4 flex flex-wrap gap-4">
                <div v-if="uploadResult.coverArtUrl">
                  <p class="mb-1.5 text-xs font-medium text-surface-500">Embedded</p>
                  <img
                    :src="uploadResult.coverArtUrl"
                    alt="Embedded Cover"
                    class="h-24 w-24 rounded-lg object-cover shadow-sm ring-1 ring-surface-200"
                    loading="lazy"
                  />
                </div>
                <div v-if="spotifyAlbumCover">
                  <p class="mb-1.5 text-xs font-medium text-surface-500">Spotify Art</p>
                  <img
                    :src="spotifyAlbumCover"
                    alt="Spotify Album Art"
                    class="h-24 w-24 rounded-lg object-cover shadow-sm ring-2 ring-green-500/50"
                    loading="lazy"
                  />
                </div>
                <div v-if="spotifyArtistImage">
                  <p class="mb-1.5 text-xs font-medium text-surface-500">Artist Image</p>
                  <img
                    :src="spotifyArtistImage"
                    alt="Artist Image"
                    class="h-24 w-24 rounded-lg object-cover shadow-sm ring-2 ring-purple-500/50"
                    loading="lazy"
                  />
                </div>
              </div>

              <div class="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
                <div
                  v-for="suggestion in suggestions"
                  :key="suggestion.field"
                  class="flex items-center gap-2 rounded-lg bg-surface-50/80 px-3 py-2 dark:bg-surface-800/50"
                >
                  <span class="w-24 shrink-0 text-xs font-medium text-surface-500">{{ suggestion.field }}</span>
                  <span class="truncate text-xs font-semibold text-surface-900 dark:text-surface-100">
                    {{ suggestion.value }}
                  </span>
                  <div class="ml-auto flex shrink-0 gap-1">
                    <Badge :value="sourceLabel(suggestion.source)" :severity="sourceSeverity(suggestion.source)" size="small" />
                    <Badge :value="confidenceLabel(suggestion.confidence)" :severity="confidenceSeverity(suggestion.confidence)" size="small" />
                  </div>
                </div>
              </div>

              <div v-if="musicBrainzInfo || lastFmInfo" class="mt-4 flex flex-wrap gap-x-6 gap-y-1 border-t border-primary/10 pt-4 text-xs text-surface-500">
                <span v-if="musicBrainzInfo"><i aria-hidden="true" class="pi pi-book mr-1"></i>{{ musicBrainzInfo }}</span>
                <span v-if="lastFmInfo"><i aria-hidden="true" class="pi pi-star mr-1"></i>{{ lastFmInfo }}</span>
              </div>
            </div>

            <!-- Extracted Metadata -->
            <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              <div v-for="(group, gIdx) in metadataGroups" :key="gIdx" class="space-y-2.5">
                <div v-for="field in group" :key="field.label" class="flex items-center gap-2">
                  <span class="w-20 shrink-0 text-xs font-medium text-surface-500">{{ field.label }}</span>
                  <span
                    v-if="field.value"
                    class="flex items-center gap-1 text-xs text-emerald-600 dark:text-emerald-400"
                  >
                    <i aria-hidden="true" class="pi pi-check-circle text-[10px]"></i>{{ field.value }}
                  </span>
                  <span v-else class="flex items-center gap-1 text-xs text-orange-500">
                    <i aria-hidden="true" class="pi pi-exclamation-circle text-[10px]"></i>Not found
                  </span>
                </div>
              </div>
            </div>

            <div v-if="lyricsSnippet" class="mt-5 border-t border-surface-200 pt-5 dark:border-surface-700">
              <h4 class="mb-2 text-sm font-medium text-surface-700 dark:text-surface-300">Lyrics Snippet</h4>
              <pre class="max-h-32 overflow-y-auto whitespace-pre-wrap rounded-lg bg-surface-100 p-3 text-xs leading-relaxed text-surface-600 break-words dark:bg-surface-700 dark:text-surface-400">{{ lyricsSnippet }}</pre>
            </div>
          </template>
        </Card>
      </div>
    </Transition>

    <!-- Drafts List -->
    <div class="mt-8">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
        <h3 class="text-lg font-semibold text-surface-900 dark:text-surface-100">
          Drafts
          <span v-if="totalItems > 0" class="ml-1.5 text-sm font-normal text-surface-400">({{ totalItems }})</span>
        </h3>
        <div class="flex items-center gap-3">
          <Button
            v-if="bulkReviewHint"
            label="Review All Pending"
            icon="pi pi-list"
            size="small"
            severity="info"
            outlined
            @click="goBulkReview"
          />
          <SelectButton
            v-model="statusFilter"
            :options="statusOptions"
            option-label="label"
            option-value="value"
            allow-empty
            class="text-sm"
            @change="onFilterChange"
          />
        </div>
      </div>

      <Transition name="fade">
        <Message
          v-if="bulkReviewHint"
          severity="info"
          :closable="true"
          class="mb-4"
        >
          <div class="flex items-center gap-2">
            <i aria-hidden="true" class="pi pi-info-circle"></i>
            <span>
              You have <strong>{{ draftsInReview }}</strong> draft{{ draftsInReview > 1 ? 's' : '' }} pending review.
              <a class="cursor-pointer underline" @click="goBulkReview">Review all pending</a>
              to process them sequentially.
            </span>
          </div>
        </Message>
      </Transition>

      <Transition name="fade">
        <Message
          v-if="draftsError"
          severity="error"
          :closable="true"
          class="mb-4"
          @close="draftsError = null"
        >
          <div class="flex items-center gap-3">
            <span>Failed to load drafts.</span>
            <Button label="Retry" icon="pi pi-refresh" size="small" severity="danger" text @click="loadDrafts" />
          </div>
        </Message>
      </Transition>

      <Card>
        <template #content>
          <DataTable
            :value="drafts"
            :loading="loadingDrafts"
            striped-rows
            size="small"
            class="w-full"
            :empty-message="loadingDrafts ? ' ' : 'No drafts found. Upload an audio file to get started.'"
          >
            <Column field="originalFilename" header="Filename" sortable>
              <template #body="{ data }">
                <div class="flex items-center gap-2">
                  <i
                    :class="data.hasCoverArt ? 'pi pi-check-circle text-emerald-500' : 'pi pi-circle-thin text-orange-400'"
                    class="text-sm"
                  ></i>
                  <span class="truncate max-w-48">{{ data.originalFilename }}</span>
                </div>
              </template>
            </Column>
            <Column field="title" header="Title" sortable>
              <template #body="{ data }">
                <span v-if="data.title" class="text-emerald-600 dark:text-emerald-400">
                  <i aria-hidden="true" class="pi pi-check-circle mr-1 text-xs"></i>{{ data.title }}
                </span>
                <span v-else class="text-orange-500">
                  <i aria-hidden="true" class="pi pi-exclamation-circle mr-1 text-xs"></i>Not found
                </span>
              </template>
            </Column>
            <Column field="artist" header="Artist">
              <template #body="{ data }">
                <span v-if="data.artist" class="text-emerald-600 dark:text-emerald-400">
                  <i aria-hidden="true" class="pi pi-check-circle mr-1 text-xs"></i>{{ data.artist }}
                </span>
                <span v-else class="text-orange-500">
                  <i aria-hidden="true" class="pi pi-exclamation-circle mr-1 text-xs"></i>Not found
                </span>
              </template>
            </Column>
            <Column field="format" header="Format" class="w-20" />
            <Column field="durationSeconds" header="Duration">
              <template #body="{ data }">
                {{ formatDuration(data.durationSeconds) }}
              </template>
            </Column>
            <Column field="status" header="Status" class="w-28">
              <template #body="{ data }">
                <Badge :value="data.status" :severity="statusSeverity(data.status)" />
              </template>
            </Column>
            <Column field="createdAt" header="Uploaded" class="w-36">
              <template #body="{ data }">
                {{ formatDate(data.createdAt) }}
              </template>
            </Column>
            <Column header="Actions" class="w-24">
              <template #body="{ data }">
                <Button
                  v-if="data.status === 'review'"
                  label="Review"
                  icon="pi pi-eye"
                  size="small"
                  severity="info"
                  @click="openReview(data.id)"
                />
              </template>
            </Column>
          </DataTable>

          <Paginator
            v-if="totalItems > 0"
            :rows="pageSize"
            :total-records="totalItems"
            :first="(currentPage - 1) * pageSize"
            class="mt-4"
            @page="handlePageChange"
          />
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import { useIngestionApi } from '@/services/api/ingestion/routes'
import type { UploadResponse, DraftListItem, EnrichmentResult } from '@/services/api/ingestion/types'

const router = useRouter()
const ingestionApi = useIngestionApi()

const fileInputRef = ref<HTMLInputElement | null>(null)
const isDragOver = ref(false)
const uploading = ref(false)
let uploadAbort: AbortController | null = null
const uploadProgress = ref('')
const uploadError = ref<string | null>(null)
const uploadResult = ref<UploadResponse | null>(null)

const enrichingPoll = ref(false)
const enrichAttempt = ref(0)
const enrichMaxAttempts = 30
let enrichPollTimer: ReturnType<typeof setInterval> | null = null

const enrichmentResult = ref<EnrichmentResult | null>(null)

const statusOptions = [
  { label: 'All', value: '' },
  { label: 'Pending', value: 'pending' },
  { label: 'Enriching', value: 'enriching' },
  { label: 'Review', value: 'review' },
  { label: 'Accepted', value: 'accepted' },
  { label: 'Rejected', value: 'rejected' },
  { label: 'Published', value: 'published' },
  { label: 'Enrich. Failed', value: 'enrichment_failed' },
]
const statusFilter = ref('')
const drafts = ref<DraftListItem[]>([])
const loadingDrafts = ref(false)
const draftsError = ref<string | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const totalItems = ref(0)

const enrichPercent = computed(() => (enrichAttempt.value / enrichMaxAttempts) * 100)

const bulkReviewHint = computed(() => draftsInReview.value >= 10)
const draftsInReview = computed(() => drafts.value.filter(d => d.status === 'review').length)

const lyricsSnippet = computed(() => {
  const lyrics = uploadResult.value?.extractedMetadata?.lyrics
  if (!lyrics) return ''
  return lyrics.split('\n').slice(0, 5).join('\n')
})

const metadataGroups = computed(() => [
  [
    { label: 'Title', value: uploadResult.value?.extractedMetadata?.title },
    { label: 'Artist', value: uploadResult.value?.extractedMetadata?.artist },
    { label: 'Album', value: uploadResult.value?.extractedMetadata?.album },
    { label: 'Album Artist', value: uploadResult.value?.extractedMetadata?.albumArtist },
  ],
  [
    { label: 'Year', value: uploadResult.value?.extractedMetadata?.year ? String(uploadResult.value.extractedMetadata.year) : '' },
    { label: 'Genre', value: uploadResult.value?.extractedMetadata?.genre },
    {
      label: 'Track',
      value: formatTrack(uploadResult.value?.extractedMetadata?.trackNumber, uploadResult.value?.extractedMetadata?.trackTotal),
    },
    {
      label: 'Disc',
      value: formatTrack(uploadResult.value?.extractedMetadata?.discNumber, uploadResult.value?.extractedMetadata?.discTotal),
    },
  ],
  [
    { label: 'Duration', value: formatDuration(uploadResult.value?.durationSeconds) },
    { label: 'Bitrate', value: uploadResult.value?.bitrate ? `${uploadResult.value.bitrate} kbps` : '' },
    { label: 'Format', value: uploadResult.value?.format },
    { label: 'Composer', value: uploadResult.value?.extractedMetadata?.composer },
  ],
])

const suggestions = computed(() => enrichmentResult.value?.suggestions || [])
const showSuggestions = computed(() => suggestions.value.length > 0)

const spotifyAlbumCover = computed(() => enrichmentResult.value?.spotify?.albumCoverUrl || null)
const spotifyArtistImage = computed(() => enrichmentResult.value?.spotify?.artistImageUrl || null)

const musicBrainzInfo = computed(() => {
  const mb = enrichmentResult.value?.musicbrainz
  if (!mb) return null
  if (!mb.artistName && !mb.albumName && !mb.releaseYear) return null
  const parts: string[] = []
  if (mb.artistName) parts.push(`Artist: ${mb.artistName}`)
  if (mb.albumName) parts.push(`Album: ${mb.albumName}`)
  if (mb.releaseYear) parts.push(`Year: ${mb.releaseYear}`)
  if (mb.genres?.length) parts.push(`Genres: ${mb.genres.slice(0, 4).join(', ')}`)
  return parts.join(' \u00B7 ')
})

const lastFmInfo = computed(() => {
  const lf = enrichmentResult.value?.lastfm
  if (!lf) return null
  if (!lf.playCount && !lf.listenerCount && !lf.tags?.length) return null
  const parts: string[] = []
  if (lf.playCount) parts.push(`${lf.playCount.toLocaleString()} plays`)
  if (lf.listenerCount) parts.push(`${lf.listenerCount.toLocaleString()} listeners`)
  if (lf.tags?.length) parts.push(`Tags: ${lf.tags.slice(0, 4).join(', ')}`)
  return parts.join(' \u00B7 ')
})

function sourceLabel(source: string): string {
  switch (source) {
    case 'file': return 'Embedded'
    case 'musicbrainz': return 'MusicBrainz'
    case 'lastfm': return 'Last.fm'
    case 'spotify': return 'Spotify'
    default: return source
  }
}

function sourceSeverity(source: string) {
  switch (source) {
    case 'file': return 'info'
    case 'musicbrainz': return 'warn'
    case 'lastfm': return 'help'
    case 'spotify': return 'success'
    default: return undefined
  }
}

function confidenceLabel(confidence: string): string {
  switch (confidence) {
    case 'exact_match': return 'Exact'
    case 'fuzzy': return 'Fuzzy'
    case 'fallback': return 'Fallback'
    default: return confidence
  }
}

function confidenceSeverity(confidence: string) {
  switch (confidence) {
    case 'exact_match': return 'success'
    case 'fuzzy': return 'warn'
    case 'fallback': return 'danger'
    default: return undefined
  }
}

function openReview(id: string) {
  router.push({ name: 'admin.ingestion.review', params: { id } })
}

function goBulkReview() {
  const reviewId = drafts.value.find(d => d.status === 'review')?.id
  if (reviewId) {
    router.push({ name: 'admin.ingestion.review', params: { id: reviewId } })
  }
}

function triggerFileInput() {
  if (!uploading.value) {
    fileInputRef.value?.click()
  }
}

function handleDrop(e: DragEvent) {
  isDragOver.value = false
  if (uploading.value) return
  const files = e.dataTransfer?.files
  if (files && files.length > 0 && files[0]) {
    uploadFile(files[0])
  }
}

function handleFileSelect(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    uploadFile(file)
  }
  if (fileInputRef.value) {
    fileInputRef.value.value = ''
  }
}

async function uploadFile(file: File) {
  uploading.value = true
  enrichPollingStop()
  uploadError.value = null
  uploadResult.value = null
  enrichmentResult.value = null

  uploadAbort = new AbortController()

  try {
    const result = await ingestionApi.uploadAudio(file, undefined, {
      signal: uploadAbort.signal,
      onUploadProgress(progressEvent) {
        if (progressEvent.total) {
          const pct = Math.round((progressEvent.loaded / progressEvent.total) * 100)
          const mbLoaded = (progressEvent.loaded / 1024 / 1024).toFixed(1)
          const mbTotal = (progressEvent.total / 1024 / 1024).toFixed(1)
          uploadProgress.value = `${pct}% (${mbLoaded} MB / ${mbTotal} MB)`
        }
      },
    })
    uploadResult.value = result
    if (result.status === 'enriching') {
      enrichPollingStart(result.draftId)
    }
    await loadDrafts()
  } catch (err: any) {
    const errObj = err as { code?: string; message?: string } | null
    if (errObj?.code !== 'ERR_CANCELED' && errObj?.message !== 'canceled') {
      uploadError.value = errObj?.message || 'Upload failed. Please try again.'
    }
  } finally {
    uploading.value = false
    uploadAbort = null
  }
}

function cancelUpload() {
  uploadAbort?.abort()
  uploadAbort = null
  uploading.value = false
  enrichPollingStop()
}

async function triggerEnrich(draftId: string) {
  try {
    await ingestionApi.enrichDraft(draftId)
    enrichPollingStart(draftId)
  } catch {
    uploadError.value = 'Failed to start enrichment.'
  }
}

function enrichPollingStart(draftId: string) {
  enrichingPoll.value = true
  enrichAttempt.value = 0

  enrichPollTimer =   // TODO MEDIUM: Heavy polling (2s interval x 30 = 60s). Use WebSocket or SSE instead.
  setInterval(async () => {
    enrichAttempt.value++

    try {
      const detail = await ingestionApi.getDraftDetail(draftId)
      if (detail && detail.status !== 'enriching') {
        enrichPollingStop()
        if (detail.enrichedMetadata?.enrichment_attempted) {
          enrichmentResult.value = detail.enrichedMetadata
        }
        uploadResult.value = { ...uploadResult.value!, status: detail.status }
        await loadDrafts()
        return
      }
    } catch {
      // ignore polling errors
    }

    if (enrichAttempt.value >= enrichMaxAttempts) {
      enrichPollingStop()
    }
  }, 2000)
}

function enrichPollingStop() {
  if (enrichPollTimer) {
    clearInterval(enrichPollTimer)
    enrichPollTimer = null
  }
  enrichingPoll.value = false
  enrichAttempt.value = 0
}

async function loadDrafts() {
  loadingDrafts.value = true
  draftsError.value = null
  try {
    const response = await ingestionApi.listDrafts({
      status: statusFilter.value || undefined,
      page: currentPage.value,
      limit: pageSize.value,
    })
    drafts.value = response.items || []
    totalItems.value = response.total || 0
  } catch {
    drafts.value = []
    totalItems.value = 0
    draftsError.value = 'Failed to load drafts.'
  } finally {
    loadingDrafts.value = false
  }
}

function onFilterChange() {
  currentPage.value = 1
  loadDrafts()
}

function handlePageChange(event: any) {
  const ev = event as { first: number; rows: number }
  currentPage.value = Math.floor(ev.first / ev.rows) + 1
  pageSize.value = ev.rows
  loadDrafts()
}

function formatDuration(seconds?: number): string {
  if (!seconds) return ''
  const m = Math.floor(seconds / 60)
  const s = Math.round(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

function formatTrack(num?: number, total?: number): string {
  if (!num && !total) return ''
  if (num && total) return `${num} / ${total}`
  if (num) return String(num)
  return ''
}

function statusSeverity(status: string) {
  switch (status) {
    case 'pending': return 'warn'
    case 'enriching': return 'info'
    case 'review': return 'info'
    case 'accepted': return 'success'
    case 'rejected': return 'danger'
    case 'published': return 'success'
    case 'enrichment_failed': return 'danger'
    default: return undefined
  }
}

function formatDate(dateStr: string): string {
  try {
    const d = new Date(dateStr)
    if (isNaN(d.getTime())) return dateStr
    return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
  } catch {
    return dateStr
  }
}

onMounted(() => {
  loadDrafts()
})

onUnmounted(() => {
  enrichPollingStop()
  uploadAbort?.abort()
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
