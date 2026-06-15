<template>
  <div class="mx-auto w-full max-w-6xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Content"
      title="Music Ingestion"
      description="Upload audio files to extract metadata and create ingestion drafts."
    />

    <div class="mb-8">
      <div
        class="flex cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed p-12 transition-colors"
        :class="isDragOver
          ? 'border-primary bg-primary/5'
          : 'border-surface-300 hover:border-surface-400 dark:border-surface-600 dark:hover:border-surface-500'"
        @dragover.prevent="isDragOver = true"
        @dragleave.prevent="isDragOver = false"
        @drop.prevent="handleDrop"
        @click="triggerFileInput"
      >
        <i class="pi pi-cloud-upload mb-4 text-4xl text-surface-400"></i>
        <p class="mb-2 text-sm font-medium text-surface-700 dark:text-surface-300">
          Drag and drop an audio file here
        </p>
        <p class="mb-4 text-xs text-surface-500">
          or click to browse (MP3, FLAC, OGG, M4A, AAC — max 200MB)
        </p>
        <Button label="Select File" icon="pi pi-folder-open" severity="secondary" outlined />
        <input
          ref="fileInputRef"
          type="file"
          accept=".mp3,.flac,.ogg,.m4a,.aac,audio/mpeg,audio/flac,audio/ogg,audio/mp4,audio/aac"
          class="hidden"
          @change="handleFileSelect"
        />
      </div>
    </div>

    <div v-if="uploading" class="mb-8">
      <ProgressBar mode="indeterminate" class="h-2" />
      <p class="mt-2 text-center text-sm text-surface-500">Extracting metadata from file...</p>
    </div>

    <div v-if="enrichingPoll" class="mb-8">
      <ProgressBar mode="indeterminate" class="h-2" />
      <p class="mt-2 text-center text-sm text-surface-400">
        <i class="pi pi-spin pi-spinner mr-2"></i>Enriching metadata from external sources...
      </p>
    </div>

    <div v-if="uploadError" class="mb-8">
      <Message severity="error" :closable="false">
        {{ uploadError }}
      </Message>
    </div>

    <div v-if="uploadResult" class="mb-8">
      <div class="mb-4 flex items-center justify-between">
        <h3 class="text-lg font-semibold text-surface-900 dark:text-surface-100">
          Extraction Results
        </h3>
        <Badge :value="uploadResult.status" :severity="statusSeverity(uploadResult.status)" />
      </div>

      <div v-if="showSuggestions" class="mb-4 rounded-lg border border-primary/30 bg-primary/5 p-4">
        <h4 class="mb-3 text-sm font-semibold text-primary-700 dark:text-primary-300">
          <i class="pi pi-magic mr-1"></i>Enrichment Suggestions
        </h4>

        <div class="mb-4 flex flex-col gap-6 md:flex-row">
          <div class="flex-shrink-0">
            <p class="mb-2 text-xs font-medium text-surface-500">Embedded</p>
            <img
              v-if="uploadResult.coverArtUrl"
              :src="uploadResult.coverArtUrl"
              alt="Embedded Cover"
              class="h-40 w-40 rounded-lg object-cover shadow-md"
            />
            <div
              v-else
              class="flex h-40 w-40 items-center justify-center rounded-lg bg-surface-200 dark:bg-surface-700"
            >
              <i class="pi pi-image text-3xl text-surface-400"></i>
            </div>
          </div>

          <div v-if="spotifyAlbumCover" class="flex-shrink-0">
            <p class="mb-2 text-xs font-medium text-surface-500">Spotify</p>
            <img
              :src="spotifyAlbumCover"
              alt="Spotify Album Art"
              class="h-40 w-40 rounded-lg object-cover shadow-md ring-2 ring-green-500"
            />
          </div>

          <div v-if="spotifyArtistImage" class="flex-shrink-0">
            <p class="mb-2 text-xs font-medium text-surface-500">Artist Image</p>
            <img
              :src="spotifyArtistImage"
              alt="Artist Image"
              class="h-40 w-40 rounded-lg object-cover shadow-md ring-2 ring-purple-500"
            />
          </div>
        </div>

        <div class="space-y-2">
          <div v-for="suggestion in suggestions" :key="suggestion.field" class="flex items-center gap-2 rounded bg-surface-50 p-2 dark:bg-surface-800">
            <span class="w-24 text-xs font-medium text-surface-500">{{ suggestion.field }}</span>
            <span class="text-xs font-semibold text-surface-900 dark:text-surface-100">{{ suggestion.value }}</span>
            <Badge
              :value="sourceLabel(suggestion.source)"
              :severity="sourceSeverity(suggestion.source)"
              size="small"
            />
            <Badge
              :value="suggestion.confidence"
              :severity="confidenceSeverity(suggestion.confidence)"
              size="small"
            />
          </div>
        </div>

        <div v-if="musicBrainzInfo" class="mt-3 border-t border-primary/20 pt-3">
          <p class="mb-1 text-xs text-surface-500">
            <i class="pi pi-book mr-1"></i>MusicBrainz
          </p>
          <p class="text-xs text-surface-700 dark:text-surface-300">
            {{ musicBrainzInfo }}
          </p>
        </div>

        <div v-if="lastFmInfo" class="mt-2">
          <p class="mb-1 text-xs text-surface-500">
            <i class="pi pi-star mr-1"></i>Last.fm
          </p>
          <p class="text-xs text-surface-700 dark:text-surface-300">
            {{ lastFmInfo }}
          </p>
        </div>
      </div>

      <div class="rounded-lg border border-surface-200 bg-surface-50 p-6 dark:border-surface-700 dark:bg-surface-800">
        <div class="mb-6 flex flex-col gap-6 md:flex-row">
          <div v-if="uploadResult.coverArtUrl" class="flex-shrink-0">
            <img
              :src="uploadResult.coverArtUrl"
              alt="Cover Art"
              class="h-40 w-40 rounded-lg object-cover shadow-md"
            />
          </div>
          <div v-else class="flex h-40 w-40 flex-shrink-0 items-center justify-center rounded-lg bg-surface-200 dark:bg-surface-700">
            <i class="pi pi-image text-3xl text-surface-400"></i>
            <span class="ml-2 text-sm text-surface-500">No Cover</span>
          </div>

          <div class="flex-1 space-y-3">
            <div v-for="field in col1" :key="field.label" class="flex items-center gap-2">
              <span class="w-20 text-xs font-medium text-surface-500">{{ field.label }}</span>
              <span v-if="field.value" class="flex items-center gap-1 text-xs text-green-600 dark:text-green-400">
                <i class="pi pi-check-circle text-xs"></i>{{ field.value }}
              </span>
              <span v-else class="flex items-center gap-1 text-xs text-orange-500">
                <i class="pi pi-exclamation-circle text-xs"></i>Not found
              </span>
            </div>
          </div>

          <div class="flex-1 space-y-3">
            <div v-for="field in col2" :key="field.label" class="flex items-center gap-2">
              <span class="w-20 text-xs font-medium text-surface-500">{{ field.label }}</span>
              <span v-if="field.value" class="flex items-center gap-1 text-xs text-green-600 dark:text-green-400">
                <i class="pi pi-check-circle text-xs"></i>{{ field.value }}
              </span>
              <span v-else class="flex items-center gap-1 text-xs text-orange-500">
                <i class="pi pi-exclamation-circle text-xs"></i>Not found
              </span>
            </div>
          </div>

          <div class="flex-1 space-y-3">
            <div v-for="field in col3" :key="field.label" class="flex items-center gap-2">
              <span class="w-20 text-xs font-medium text-surface-500">{{ field.label }}</span>
              <span v-if="field.value" class="flex items-center gap-1 text-xs text-green-600 dark:text-green-400">
                <i class="pi pi-check-circle text-xs"></i>{{ field.value }}
              </span>
              <span v-else class="flex items-center gap-1 text-xs text-orange-500">
                <i class="pi pi-exclamation-circle text-xs"></i>Not found
              </span>
            </div>
          </div>
        </div>

        <div v-if="uploadResult.extractedMetadata?.lyrics" class="mt-4 border-t border-surface-200 pt-4 dark:border-surface-700">
          <h4 class="mb-2 text-sm font-medium text-surface-700 dark:text-surface-300">Lyrics Snippet</h4>
          <pre class="max-h-32 overflow-y-auto rounded bg-surface-100 p-3 text-xs text-surface-600 dark:bg-surface-700 dark:text-surface-400">{{ lyricsSnippet }}</pre>
        </div>
      </div>
    </div>

    <div class="mt-8">
      <div class="mb-4 flex items-center justify-between">
        <h3 class="text-lg font-semibold text-surface-900 dark:text-surface-100">
          Recent Drafts
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
            @change="loadDrafts"
          />
        </div>
      </div>

      <Message
        v-if="bulkReviewHint"
        severity="info"
        :closable="true"
        class="mb-4"
      >
        <div class="flex items-center gap-2">
          <i class="pi pi-info-circle"></i>
          <span>You have <strong>{{ draftsInReview }}</strong> drafts pending review.
            <a class="cursor-pointer underline" @click="goBulkReview">Review all pending</a>
            to process them sequentially.</span>
        </div>
      </Message>

      <DataTable
        :value="drafts"
        :loading="loadingDrafts"
        striped-rows
        size="small"
        class="w-full"
      >
        <Column field="originalFilename" header="Filename" sortable>
          <template #body="{ data }">
            <div class="flex items-center gap-2">
              <i
                :class="data.hasCoverArt ? 'pi pi-check-circle text-green-500' : 'pi pi-circle text-orange-400'"
                class="text-sm"
              ></i>
              <span class="truncate">{{ data.originalFilename }}</span>
            </div>
          </template>
        </Column>
        <Column field="title" header="Title" sortable>
          <template #body="{ data }">
            <span v-if="data.title" class="text-green-600 dark:text-green-400">
              <i class="pi pi-check-circle mr-1 text-xs"></i>{{ data.title }}
            </span>
            <span v-else class="text-orange-500">
              <i class="pi pi-exclamation-circle mr-1 text-xs"></i>Not found
            </span>
          </template>
        </Column>
        <Column field="artist" header="Artist">
          <template #body="{ data }">
            <span v-if="data.artist" class="text-green-600 dark:text-green-400">
              <i class="pi pi-check-circle mr-1 text-xs"></i>{{ data.artist }}
            </span>
            <span v-else class="text-orange-500">
              <i class="pi pi-exclamation-circle mr-1 text-xs"></i>Not found
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import { useIngestionApi } from '@/services/api/ingestion/routes'
import type { UploadResponse, DraftListItem, EnrichmentResult } from '@/services/api/ingestion/types'

const router = useRouter()

const ingestionApi = useIngestionApi()

const fileInputRef = ref<HTMLInputElement | null>(null)
const isDragOver = ref(false)
const uploading = ref(false)
const uploadError = ref<string | null>(null)
const uploadResult = ref<UploadResponse | null>(null)

const enrichingPoll = ref(false)
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
const currentPage = ref(1)
const pageSize = ref(20)
const totalItems = ref(0)

const bulkReviewHint = computed(() => draftsInReview.value >= 10)
const draftsInReview = computed(() => drafts.value.filter(d => d.status === 'review').length)

const lyricsSnippet = computed(() => {
  const lyrics = uploadResult.value?.extractedMetadata?.lyrics
  if (!lyrics) return ''
  const lines = lyrics.split('\n').slice(0, 5)
  return lines.join('\n')
})

const col1 = computed(() => [
  { label: 'Title', value: uploadResult.value?.extractedMetadata?.title },
  { label: 'Artist', value: uploadResult.value?.extractedMetadata?.artist },
  { label: 'Album', value: uploadResult.value?.extractedMetadata?.album },
  { label: 'Album Artist', value: uploadResult.value?.extractedMetadata?.albumArtist },
])

const col2 = computed(() => [
  { label: 'Year', value: uploadResult.value?.extractedMetadata?.year ? String(uploadResult.value.extractedMetadata.year) : '' },
  { label: 'Genre', value: uploadResult.value?.extractedMetadata?.genre },
  { label: 'Track', value: formatTrack(uploadResult.value?.extractedMetadata?.trackNumber, uploadResult.value?.extractedMetadata?.trackTotal) },
  { label: 'Disc', value: formatTrack(uploadResult.value?.extractedMetadata?.discNumber, uploadResult.value?.extractedMetadata?.discTotal) },
])

const col3 = computed(() => [
  { label: 'Duration', value: formatDuration(uploadResult.value?.durationSeconds) },
  { label: 'Bitrate', value: uploadResult.value?.bitrate ? `${uploadResult.value.bitrate} kbps` : '' },
  { label: 'Format', value: uploadResult.value?.format },
  { label: 'Composer', value: uploadResult.value?.extractedMetadata?.composer },
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
  fileInputRef.value?.click()
}

function handleDrop(e: DragEvent) {
  isDragOver.value = false
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
}

async function uploadFile(file: File) {
  uploading.value = true
  enrichPollingStop()
  uploadError.value = null
  uploadResult.value = null
  enrichmentResult.value = null

  try {
    const result = await ingestionApi.uploadAudio(file)
    uploadResult.value = result
    if (fileInputRef.value) {
      fileInputRef.value.value = ''
    }
    if (result.status === 'enriching') {
      enrichPollingStart(result.draftId)
    }
    await loadDrafts()
  } catch (err: any) {
    uploadError.value = err?.message || 'Upload failed. Please try again.'
  } finally {
    uploading.value = false
  }
}

function enrichPollingStart(draftId: string) {
  enrichingPoll.value = true
  let attempt = 0
  const maxAttempts = 30

  enrichPollTimer = setInterval(async () => {
    attempt++
    if (attempt > maxAttempts) {
      enrichPollingStop()
      return
    }

    try {
      const detail = await ingestionApi.getDraftDetail(draftId)
      if (detail && detail.status !== 'enriching') {
        enrichingPoll.value = false
        if (enrichPollTimer) {
          clearInterval(enrichPollTimer)
          enrichPollTimer = null
        }
        if (detail.enrichedMetadata && detail.enrichedMetadata.enrichment_attempted) {
          enrichmentResult.value = detail.enrichedMetadata
        }
        uploadResult.value = { ...uploadResult.value!, status: detail.status }
        await loadDrafts()
      }
    } catch {
      // ignore polling errors
    }
  }, 2000)
}

function enrichPollingStop() {
  if (enrichPollTimer) {
    clearInterval(enrichPollTimer)
    enrichPollTimer = null
  }
  enrichingPoll.value = false
}

async function loadDrafts() {
  loadingDrafts.value = true
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
  } finally {
    loadingDrafts.value = false
  }
}

function handlePageChange(event: any) {
  currentPage.value = Math.floor(event.first / event.rows) + 1
  pageSize.value = event.rows
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
})
</script>

<style scoped>
pre {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
