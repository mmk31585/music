<template>
  <div class="mx-auto w-full max-w-5xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Content"
      title="Review Ingestion Draft"
      :description="`${draftDetail?.originalFilename || ''} — ${formatDuration(draftDetail?.durationSeconds)}`"
    >
      <template #actions>
        <Button
          v-if="!enriching && (draftDetail?.status === 'pending' || draftDetail?.status === 'enrichment_failed')"
          label="Enrich Now"
          icon="pi pi-magic"
          severity="warn"
          size="small"
          @click="triggerEnrich"
        />
        <Button
          v-if="enriching"
          label="Enriching..."
          icon="pi pi-spin pi-spinner"
          severity="warn"
          size="small"
          disabled
        />
        <Button
          label="Back to List"
          icon="pi pi-arrow-left"
          severity="secondary"
          text
          @click="goBack"
        />
      </template>
    </AdminSectionHeader>

    <div v-if="loading" class="py-20 text-center">
      <Loader2 aria-hidden="true" class="text-3xl text-slate-400 animate-spin"></Loader2>
      <p class="mt-4 text-slate-500">Loading draft...</p>
    </div>

    <div v-else-if="error" class="mb-8">
      <Message severity="error" :closable="false">{{ error }}</Message>
    </div>

    <IngestionReviewSuccess
      v-else-if="published && publishResult"
      :filename="draftDetail?.originalFilename || ''"
      :publish-result="publishResult"
      :final-metadata="finalMetadata"
      @back="goBack"
    />

    <template v-else>
      <!-- Extracted metadata summary -->
      <div v-if="extractedTitle || extractedArtist || extractedAlbum" class="mb-6 flex flex-wrap items-center gap-x-5 gap-y-1 rounded-lg border border-white/6 bg-white/3 px-4 py-2.5">
        <div v-if="extractedTitle" class="flex items-center gap-1.5 text-xs">
          <Music aria-hidden="true" class="text-slate-500"></Music>
          <span class="text-slate-500">Title:</span>
          <span class="font-medium text-white">{{ extractedTitle }}</span>
        </div>
        <div v-if="extractedArtist" class="flex items-center gap-1.5 text-xs">
          <User aria-hidden="true" class="text-slate-500"></User>
          <span class="text-slate-500">Artist:</span>
          <span class="font-medium text-white">{{ extractedArtist }}</span>
        </div>
        <div v-if="extractedAlbum" class="flex items-center gap-1.5 text-xs">
          <Book aria-hidden="true" class="text-slate-500"></Book>
          <span class="text-slate-500">Album:</span>
          <span class="font-medium text-white">{{ extractedAlbum }}</span>
        </div>
      </div>

      <div class="mb-8">
        <div class="flex items-center gap-2">
          <div
            v-for="(s, i) in steps"
            :key="i"
            class="flex items-center"
          >
            <div
              class="flex cursor-pointer items-center gap-2 rounded-full px-3 py-1.5 text-xs font-medium transition-colors"
              :class="stepperClass(i)"
              @click="step = i"
            >
              <span
                class="flex h-5 w-5 items-center justify-center rounded-full text-[10px] font-bold"
                :class="stepIconClass(i)"
              >
                <Check aria-hidden="true" v-if="i < step" />
                <span v-else>{{ i + 1 }}</span>
              </span>
              {{ s.label }}
            </div>
            <ChevronRight aria-hidden="true" v-if="i < steps.length - 1" class="mx-2 text-xs text-slate-500"></ChevronRight>
          </div>
        </div>
      </div>

      <div v-if="enriching" class="mb-8">
        <Message severity="info" :closable="false">
          <Loader2 aria-hidden="true" class="mr-2 animate-spin"></Loader2>Enrichment in progress — covers, bio, and suggestions will appear once complete.
        </Message>
      </div>

      <div v-if="enrichmentComplete" class="mb-8">
        <Message severity="success" :closable="true" @close="enrichmentComplete = false">
          <CheckCircle aria-hidden="true" class="mr-2"></CheckCircle>Enrichment complete. Found data from
          <template v-if="enrichment?.spotify"> Spotify,</template>
          <template v-if="enrichment?.lastfm"> Last.fm,</template>
          <template v-if="enrichment?.musicbrainz"> MusicBrainz,</template>
          <template v-if="enrichment?.ml"> ML service,</template>
          <template v-if="!enrichment?.spotify && !enrichment?.lastfm && !enrichment?.musicbrainz && !enrichment?.ml"> no external sources</template>
          and applied to empty fields. Click any image or "Apply" button to use specific results.
        </Message>
      </div>

      <div v-if="!enriching && (draftDetail?.status === 'pending' || draftDetail?.status === 'enrichment_failed')" class="mb-8">
        <Message severity="warn" :closable="false">
          <AlertTriangle aria-hidden="true" class="mr-2"></AlertTriangle>No enrichment data yet. Click "Enrich Now" above or fill in the fields manually.
        </Message>
      </div>

      <Transition name="fade" mode="out-in">
        <div :key="step">
          <IngestionReviewArtistStep v-if="step === 0" /> />

          <IngestionReviewAlbumStep v-if="step === 1" />

          <IngestionReviewTrackStep v-if="step === 2" />

          <IngestionReviewConfirmStep
            v-if="step === 3"
            :final-metadata="finalMetadata"
            :publishing-error="publishingError"
          />
        </div>
      </Transition>

      <div class="mt-8 flex items-center justify-between border-t border-white/6 pt-6">
        <div>
          <Button
            v-if="step > 0"
            label="Back"
            icon="pi pi-chevron-left"
            severity="secondary"
            text
            @click="step--"
          />
        </div>
        <div class="flex items-center gap-3">
          <Button
            v-if="step < 3"
            label="Next"
            icon="pi pi-chevron-right"
            icon-pos="right"
            :disabled="!canProceed"
            @click="step++"
          />
          <Button
            v-if="step === 3"
            label="Publish to Catalog"
            icon="pi pi-check"
            :loading="publishing"
            :disabled="!isValid"
            @click="publish"
          />
          <Button
            label="Reject"
            icon="pi pi-times"
            severity="danger"
            text
            :loading="rejecting"
            @click="confirmReject"
          />
        </div>
      </div>
    </template>

    <IngestionReviewRejectDialog
      v-model:visible="rejectDialogVisible"
      :rejecting="rejecting"
      @confirm="rejectDraft"
    />
  </div>
</template>

<script setup lang="ts">
import { AlertTriangle, Book, Check, CheckCircle, ChevronRight, Loader2, Music, User } from 'lucide-vue-next'
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router'
import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import IngestionReviewSuccess from '@/components/admin/IngestionReviewSuccess.vue'
import IngestionReviewConfirmStep from '@/components/admin/IngestionReviewConfirmStep.vue'
import IngestionReviewRejectDialog from '@/components/admin/IngestionReviewRejectDialog.vue'
import IngestionReviewArtistStep from '@/components/admin/IngestionReviewArtistStep.vue'
import IngestionReviewAlbumStep from '@/components/admin/IngestionReviewAlbumStep.vue'
import IngestionReviewTrackStep from '@/components/admin/IngestionReviewTrackStep.vue'
import { useIngestionApi } from '@/services/api/ingestion/routes'
import type {
  DraftDetailResponse,
  EnrichmentResult,
  EnrichedSuggestion,
  SaveFinalMetadataRequest,
  ArtistSearchResult,
  AlbumSearchResult,
  FinalizeResult,
} from '@/services/api/ingestion/types'
import { useToast } from 'primevue/usetoast'
import { formatDuration } from '@/utils/format'
import { provideIngestionReviewContext } from '@/composables/useIngestionReviewContext'
import { useEnrichmentPolling } from '@/composables/admin'

// TODO LOW: Album upsert in finalization doesn't use embedded cover as fallback when user doesn't provide one.
defineOptions({ name: 'PageAdminIngestionReview' })

const route = useRoute()
const router = useRouter()
const toast = useToast()
const ingestionApi = useIngestionApi()
const { start: startPolling } = useEnrichmentPolling()

const draftId = route.params.id as string

const loading = ref(true)
const error = ref<string | null>(null)
const publishing = ref(false)
const publishingError = ref<string | null>(null)
const rejecting = ref(false)
const published = ref(false)
const publishResult = ref<FinalizeResult | null>(null)
const rejectDialogVisible = ref(false)
const lyricsExpanded = ref(false)

const audioPlaying = ref(false)
const audioCurrentTime = ref(0)
const audioDuration = ref(0)
const showSyncedLyrics = ref(false)

const enriching = ref(false)
const enrichmentComplete = ref(false)
const enrichingArtist = ref(false)
const enrichingAlbum = ref(false)
const enrichingTrack = ref(false)

const draftDetail = ref<DraftDetailResponse | null>(null)
const enrichment = ref<EnrichmentResult | null>(null)
const step = ref(0)

const finalMetadata = reactive<SaveFinalMetadataRequest>({
  artists: [
    { action: 'create', name: '', bio: '', imageUrl: '', country: '', musicbrainzMbid: '' },
  ],
  album: { action: 'create', title: '', releaseYear: undefined, genre: '', coverUrl: '', musicbrainzReleaseId: '' },
  track: { title: '', trackNumber: undefined, durationSeconds: 0, genre: '', lyrics: '', explicit: false, spotifyPreviewUrl: '', coverUrl: '',
  },
})

const trackAudioUrl = computed(() => {
  return draftDetail.value?.assets?.find(a => a.assetType === 'audio')?.url || null
})

const lyricsType = computed(() => {
  return findSuggestion('lyrics_type')?.value || 'plain'
})

const steps = [
  { label: 'Artist', key: 'artist' },
  { label: 'Album', key: 'album' },
  { label: 'Track', key: 'track' },
  { label: 'Confirm', key: 'confirm' },
]

const hasUnsavedChanges = ref(false)

const embeddedCover = computed(() => {
  return draftDetail.value?.assets?.find(a => a.assetType === 'cover' || a.assetType === 'track_cover')?.url || null
})

const artistImageAsset = computed(() => {
  return draftDetail.value?.assets?.find(a => a.assetType === 'artist_image')?.url || null
})

const albumCoverAsset = computed(() => {
  return draftDetail.value?.assets?.find(a => a.assetType === 'album_cover')?.url || null
})

const suggestedAlbumCover = computed(() => {
  return enrichment.value?.spotify?.albumCoverUrl || enrichment.value?.lastfm?.albumCoverUrl || enrichment.value?.ml?.albumCoverUrl || null
})

const suggestedLastfmAlbumCover = computed(() => {
  return enrichment.value?.lastfm?.albumCoverUrl || null
})

const selectedArtistId = ref<string | undefined>()
const selectedAlbumId = ref<string | undefined>()
const artistSearchQuery = ref('')
const albumSearchQuery = ref('')
const artistSearchResults = ref<ArtistSearchResult[]>([])
const albumSearchResults = ref<AlbumSearchResult[]>([])
const artistSearching = ref(false)
const albumSearching = ref(false)
let artistSearchTimer: ReturnType<typeof setTimeout> | null = null
let albumSearchTimer: ReturnType<typeof setTimeout> | null = null

const canProceed = computed(() => {
  if (step.value === 0) return (finalMetadata.artists[0]?.name?.trim().length ?? 0) > 0
  if (step.value === 1) return finalMetadata.album.title.trim().length > 0
  if (step.value === 2) return finalMetadata.track.title.trim().length > 0
  return true
})

const extractedTitle = computed(() => draftDetail.value?.extractedMetadata?.title || '')
const extractedArtist = computed(() => draftDetail.value?.extractedMetadata?.artist || '')
const extractedAlbum = computed(() => draftDetail.value?.extractedMetadata?.album || '')

const isValid = computed(() => {
  return (
    (finalMetadata.artists[0]?.name?.trim().length ?? 0) > 0 &&
    finalMetadata.album.title.trim().length > 0 &&
    finalMetadata.track.title.trim().length > 0
  )
})

function findSuggestion(field: string): EnrichedSuggestion | undefined {
  return enrichment.value?.suggestions?.find(s => s.field === field)
}

function sourceLabel(source: string): string {
  const map: Record<string, string> = { file: 'Embedded', musicbrainz: 'MB', lastfm: 'Last.fm', spotify: 'Spotify' }
  return map[source] || source
}

function sourceSeverity(source: string) {
  const map: Record<string, string> = { file: 'info', musicbrainz: 'warn', lastfm: 'help', spotify: 'success' }
  return map[source] || undefined
}

function confidenceSeverity(c: string) {
  const map: Record<string, string> = { exact_match: 'success', fuzzy: 'warn', fallback: 'danger' }
  return map[c] || undefined
}

function seekAudio(seconds: number) {
  audioCurrentTime.value = seconds
}

function stepperClass(i: number): string {
  if (i === step.value) return 'bg-emerald-500/20 text-emerald-400'
  if (i < step.value) return 'bg-green-500/10 text-green-400'
  return 'text-slate-500 hover:text-white/70'
}

function stepIconClass(i: number): string {
  if (i === step.value) return 'bg-emerald-500 text-white'
  if (i < step.value) return 'bg-green-500 text-white'
  return 'bg-white/6 text-slate-400'
}

function truncateBio(bio: string): string {
  if (!bio) return ''
  return bio.length > 300 ? bio.slice(0, 300) + '...' : bio
}

function selectExistingArtist(a: ArtistSearchResult, index: number = 0) {
  selectedArtistId.value = a.id
  const artist = finalMetadata.artists[index]
  if (!artist) return
  artist.existingId = a.id
  artist.name = a.name
  artist.country = a.country || ''
  artist.imageUrl = a.imageUrl || ''
  artist.bio = a.bio || ''
  hasUnsavedChanges.value = true
}

// ── Featured artists management ──
const featuredArtistDefaults = () => ({
  action: 'create' as const,
  name: '',
  bio: '',
  imageUrl: '',
  country: '',
  musicbrainzMbid: '',
})

function addFeaturedArtist() {
  finalMetadata.artists.push(featuredArtistDefaults())
}

function removeFeaturedArtist(index: number) {
  // Don't remove the primary artist (index 0)
  if (index === 0) return
  finalMetadata.artists.splice(index, 1)
  hasUnsavedChanges.value = true
}

function moveFeaturedArtist(fromIndex: number, direction: -1 | 1) {
  const toIndex = fromIndex + direction
  if (toIndex < 1 || toIndex >= finalMetadata.artists.length) return
  const temp = finalMetadata.artists[fromIndex]!
  finalMetadata.artists[fromIndex] = finalMetadata.artists[toIndex]!
  finalMetadata.artists[toIndex] = temp
  hasUnsavedChanges.value = true
}

/** Parse "feat." / "ft." patterns from a raw artist string */
function parseFeatArtists(rawArtist: string): string[] {
  if (!rawArtist) return []
  // Split on common feat/ft patterns
  const parts = rawArtist.split(/\s+(?:feat\.|ft\.|featuring|Feat\.|Ft\.|Featuring)\s+/i)
  const mainArtist = parts[0]?.trim()
  if (!mainArtist) return []
  const result = [mainArtist]
  if (parts.length > 1) {
    // Split the feat part further by ",", "&", "and"
    const featPart = parts.slice(1).join(', ')
    const featArtists = featPart.split(/\s*[,&]\s*|\s+and\s+/i).map(s => s.trim()).filter(Boolean)
    result.push(...featArtists)
  }
  return result
}

function selectExistingAlbum(a: AlbumSearchResult) {
  selectedAlbumId.value = a.id
  finalMetadata.album.existingId = a.id
  finalMetadata.album.title = a.title
  finalMetadata.album.coverUrl = a.coverUrl || ''
  if (a.releaseYear) finalMetadata.album.releaseYear = a.releaseYear
  hasUnsavedChanges.value = true
}

function debouncedArtistSearch() {
  if (artistSearchTimer) clearTimeout(artistSearchTimer)
  const q = artistSearchQuery.value.trim()
  if (q.length < 2) {
    artistSearchResults.value = []
    return
  }
  artistSearchTimer = setTimeout(async () => {
    artistSearching.value = true
    try {
      const res = await ingestionApi.searchArtists(q)
      artistSearchResults.value = Array.isArray(res) ? res : []
    } catch (err) {
      console.error('Artist search failed:', err)
      artistSearchResults.value = []
    } finally {
      artistSearching.value = false
    }
  }, 300)
}

function debouncedAlbumSearch() {
  if (albumSearchTimer) clearTimeout(albumSearchTimer)
  const q = albumSearchQuery.value.trim()
  if (q.length < 2) {
    albumSearchResults.value = []
    return
  }
  albumSearchTimer = setTimeout(async () => {
    albumSearching.value = true
    try {
      const res = await ingestionApi.searchAlbums(q)
      albumSearchResults.value = Array.isArray(res) ? res : []
    } catch (err) {
      console.error('Album search failed:', err)
      albumSearchResults.value = []
    } finally {
      albumSearching.value = false
    }
  }, 300)
}

async function loadDraft() {
  loading.value = true
  error.value = null
  enrichmentComplete.value = false
  try {
    const detail = await ingestionApi.getDraftDetail(draftId)
    if (!detail) {
      error.value = 'Draft not found.'
      return
    }
    draftDetail.value = detail
    enrichment.value = detail.enrichedMetadata || null
    if (detail.status === 'enriching' || detail.status === 'refetching') {
      enriching.value = true
      pollEnrichment()
    }
    prefillForm(detail)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load draft.'
  } finally {
    loading.value = false
  }
}

function onEnrichmentComplete(detail: DraftDetailResponse) {
  draftDetail.value = detail
  enrichmentComplete.value = true
  enrichingArtist.value = false
  enrichingAlbum.value = false
  enrichingTrack.value = false
  enrichment.value = detail.enrichedMetadata || null
  applyEnrichmentToEmptyFields()
  if (!hasUnsavedChanges.value) {
    prefillForm(detail)
  }
  const parts: string[] = []
  const e = enrichment.value
  if (e?.spotify?.artistImageUrl) parts.push('Spotify artist image')
  if (e?.lastfm?.artistImageUrl) parts.push('Last.fm artist image')
  if (e?.lastfm?.artistBio) parts.push('artist bio')
  if (e?.ml?.artistImageUrl) parts.push('ML artist image')
  if (e?.ml?.albumCoverUrl) parts.push('ML cover art')
  if (e?.ml?.lyrics) parts.push('lyrics')
  if (e?.spotify?.albumCoverUrl) parts.push('Spotify cover art')
  if (e?.lastfm?.albumCoverUrl) parts.push('Last.fm cover art')
  if (e?.musicbrainz?.artistMbid || e?.musicbrainz?.albumMbid) parts.push('MusicBrainz IDs')
  const summary = parts.length > 0 ? parts.join(', ') + '.' : 'No additional data found.'
  toast.add({ severity: 'success', summary: 'Enrichment complete', detail: summary, life: 4000 })
}

function pollEnrichment() {
  startPolling(draftId, {
    onComplete: (detail) => {
      enriching.value = false
      onEnrichmentComplete(detail)
    },
    onTimeout: () => {
      enriching.value = false
      enrichingArtist.value = false
      enrichingAlbum.value = false
      enrichingTrack.value = false
      toast.add({ severity: 'warn', summary: 'Timed out', detail: 'Enrichment is taking longer than expected. The data may still appear.', life: 5000 })
    },
  })
}

async function triggerEnrich() {
  enriching.value = true
  enrichmentComplete.value = false
  try {
    await ingestionApi.enrichDraft(draftId)
    pollEnrichment()
  } catch (err) {
    console.error('Enrichment failed:', err)
    enriching.value = false
    toast.add({ severity: 'error', summary: 'Failed to start enrichment', life: 3000 })
  }
}

function applyEnrichmentToEmptyFields() {
  const mb = enrichment.value?.musicbrainz
  const lfm = enrichment.value?.lastfm
  const spot = enrichment.value?.spotify
  const ml = enrichment.value?.ml
  const sug = enrichment.value?.suggestions || []

  const sugMap = new Map<string, string>()
  for (const s of sug) {
    sugMap.set(s.field, s.value)
  }

  const artist = finalMetadata.artists[0]
  if (artist) {
    if (!artist.bio) artist.bio = lfm?.artistBio || sugMap.get('artist_bio') || ''
    if (!artist.imageUrl) artist.imageUrl = spot?.artistImageUrl || lfm?.artistImageUrl || ml?.artistImageUrl || sugMap.get('artist_image_url') || ''
    if (!artist.musicbrainzMbid && mb?.artistMbid) artist.musicbrainzMbid = mb.artistMbid
  }
  if (!finalMetadata.album.musicbrainzReleaseId && mb?.albumMbid) {
    finalMetadata.album.musicbrainzReleaseId = mb.albumMbid
  }
  if (!finalMetadata.track.lyrics && ml?.lyrics) {
    finalMetadata.track.lyrics = ml.lyrics
  }
}

function fetchSection(section: 'artist' | 'album' | 'track') {
  const tags = draftDetail.value?.extractedMetadata
  if (!tags) return

  if (section === 'artist' && tags.artist) {
    const parsedArtists = parseFeatArtists(tags.artist)
    if (parsedArtists.length > 0) {
      finalMetadata.artists = parsedArtists.map((name, i) => ({
        action: 'create' as const,
        name,
        bio: finalMetadata.artists[i]?.bio || '',
        imageUrl: finalMetadata.artists[i]?.imageUrl || '',
        country: finalMetadata.artists[i]?.country || '',
        musicbrainzMbid: finalMetadata.artists[i]?.musicbrainzMbid || '',
      }))
    }
  }
  if (section === 'album') {
    if (tags.album) finalMetadata.album.title = tags.album
    if (tags.year) finalMetadata.album.releaseYear = Number(tags.year) || undefined
  }
  if (section === 'track') {
    if (tags.title) finalMetadata.track.title = tags.title
    if (tags.lyrics) finalMetadata.track.lyrics = tags.lyrics
    if (tags.trackNumber) finalMetadata.track.trackNumber = Number(tags.trackNumber) || undefined
  }
  hasUnsavedChanges.value = true
}

function prefillForm(detail: DraftDetailResponse) {
  const tags = detail.extractedMetadata
  const sug = enrichment.value?.suggestions || []
  const mb = enrichment.value?.musicbrainz
  const lfm = enrichment.value?.lastfm
  const spot = enrichment.value?.spotify
  const ml = enrichment.value?.ml

  const sugMap = new Map<string, string>()
  for (const s of sug) {
    sugMap.set(s.field, s.value)
  }

  // ── Artists (multi-artist with feat. support) ──
  const rawArtist = tags?.artist || sugMap.get('artist') || ''
  const parsedArtists = parseFeatArtists(rawArtist)

  finalMetadata.artists = parsedArtists.map((name, i) => ({
    action: 'create' as const,
    name,
    bio: '',
    imageUrl: '',
    country: '',
    musicbrainzMbid: '',
  }))

  // Apply enrichment to primary artist (index 0)
  if (finalMetadata.artists[0]) {
    finalMetadata.artists[0].bio = lfm?.artistBio || sugMap.get('artist_bio') || ''
    finalMetadata.artists[0].imageUrl = artistImageAsset.value || spot?.artistImageUrl || lfm?.artistImageUrl || ml?.artistImageUrl || sugMap.get('artist_image_url') || ''
    if (mb?.artistMbid) finalMetadata.artists[0].musicbrainzMbid = mb.artistMbid
  }

  // Apply imageUrl to featured artists from enrichment if available
  const featImage = sugMap.get('artist_image_url') || spot?.artistImageUrl || lfm?.artistImageUrl || ''
  for (let i = 1; i < finalMetadata.artists.length; i++) {
    if (!finalMetadata.artists[i]!.imageUrl) {
      finalMetadata.artists[i]!.imageUrl = featImage
    }
  }

  // ── Album ──
  finalMetadata.album.title = tags?.album || sugMap.get('album') || ''
  finalMetadata.album.releaseYear = (sugMap.get('year') || tags?.year || undefined) as number | undefined
  finalMetadata.album.genre = sugMap.get('genre') || tags?.genre || ''
  finalMetadata.album.coverUrl = albumCoverAsset.value || spot?.albumCoverUrl || lfm?.albumCoverUrl || ml?.albumCoverUrl || sugMap.get('album_cover_url') || sugMap.get('album_cover_url_lastfm') || ''
  if (mb?.albumMbid) finalMetadata.album.musicbrainzReleaseId = mb.albumMbid

  // ── Track ──
  finalMetadata.track.title = tags?.title || sugMap.get('title') || ''
  finalMetadata.track.trackNumber = tags?.trackNumber || undefined
  finalMetadata.track.durationSeconds = Math.round(detail.durationSeconds || tags?.duration || 0)
  finalMetadata.track.genre = sugMap.get('genre') || tags?.genre || ''
  finalMetadata.track.lyrics = ml?.lyrics || sugMap.get('lyrics') || tags?.lyrics || ''
  finalMetadata.track.explicit = false
  finalMetadata.track.spotifyPreviewUrl = spot?.previewUrl || ''
  finalMetadata.track.coverUrl = embeddedCover.value || spot?.albumCoverUrl || lfm?.albumCoverUrl || ml?.albumCoverUrl || ''
}

function goBack() {
  if (hasUnsavedChanges.value) {
    if (window.confirm!('You have unsaved changes. Are you sure you want to leave?')) return
  }
  router.push({ name: 'admin.ingestion' })
}

async function uploadArtistImage(e: Event, index: number = 0) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  try {
    const res = await ingestionApi.uploadDraftImage(draftId, 'artist', file)
    if (res?.url) {
      const artist = finalMetadata.artists[index]
      if (artist) {
        artist.imageUrl = res.url
      }
      hasUnsavedChanges.value = true
      toast.add({ severity: 'success', summary: 'Image uploaded', detail: 'Artist image uploaded successfully.', life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Upload failed', detail: err instanceof Error ? err.message : 'Failed to upload image.', life: 5000 })
  }
  target.value = ''
}

async function uploadAlbumCover(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  try {
    const res = await ingestionApi.uploadDraftImage(draftId, 'album', file)
    if (res?.url) {
      finalMetadata.album.coverUrl = res.url
      hasUnsavedChanges.value = true
      toast.add({ severity: 'success', summary: 'Cover uploaded', detail: 'Album cover uploaded successfully.', life: 3000 })
    }
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Upload failed', detail: err instanceof Error ? err.message : 'Failed to upload image.', life: 5000 })
  }
  target.value = ''
}

async function publish() {
  publishing.value = true
  publishingError.value = null
  try {
    await ingestionApi.saveFinalMetadata(draftId, {
      artists: finalMetadata.artists,
      album: finalMetadata.album,
      track: finalMetadata.track,
    })
    const result = await ingestionApi.finalizeDraft(draftId)
    publishResult.value = result
    hasUnsavedChanges.value = false
    published.value = true
    toast.add({ severity: 'success', summary: 'Draft published', detail: 'The draft has been published to the catalog.', life: 5000 })
  } catch (err) {
    publishingError.value = err instanceof Error ? err.message : 'Failed to publish draft.'
  } finally {
    publishing.value = false
  }
}

function confirmReject() {
  rejectDialogVisible.value = true
}

async function rejectDraft(reason: string) {
  rejecting.value = true
  try {
    await ingestionApi.rejectDraft(draftId, reason)
    hasUnsavedChanges.value = false
    rejectDialogVisible.value = false
    toast.add({ severity: 'info', summary: 'Draft rejected', detail: 'The draft has been rejected.', life: 4000 })
    setTimeout(() => router.push({ name: 'admin.ingestion' }), 1500)
  } catch (err) {
    console.error('Reject draft failed:', err)
  } finally {
    rejecting.value = false
  }
}

provideIngestionReviewContext({
  draftId,
  finalMetadata: finalMetadata as SaveFinalMetadataRequest,
  draftDetail,
  enrichment,
  step,
  canProceed,
  isValid,
  hasUnsavedChanges,
  loading,
  error,
  publishing,
  publishingError,
  rejecting,
  published,
  publishResult,
  enrichmentComplete,
  enriching,
  enrichingArtist,
  enrichingAlbum,
  enrichingTrack,
  embeddedCover,
  artistImageAsset,
  albumCoverAsset,
  suggestedAlbumCover,
  suggestedLastfmAlbumCover,
  extractedTitle,
  extractedArtist,
  extractedAlbum,
  trackAudioUrl,
  lyricsType,
  audioCurrentTime,
  audioDuration,
  audioPlaying,
  showSyncedLyrics,
  lyricsExpanded,
  artistSearchQuery,
  albumSearchQuery,
  artistSearchResults,
  albumSearchResults,
  artistSearching,
  albumSearching,
  selectedArtistId,
  selectedAlbumId,
  findSuggestion,
  sourceLabel,
  sourceSeverity,
  confidenceSeverity,
  truncateBio,
  stepperClass,
  stepIconClass,
  seekAudio,
  addFeaturedArtist,
  removeFeaturedArtist,
  moveFeaturedArtist,
  selectExistingArtist,
  selectExistingAlbum,
  uploadArtistImage,
  uploadAlbumCover,
  fetchSection,
  debouncedArtistSearch,
  debouncedAlbumSearch,
  triggerEnrich,
  publish,
  confirmReject,
  goBack,
})

function handleKeydown(e: KeyboardEvent) {
  if (rejectDialogVisible.value || loading.value || published.value) return
  const tag = (e.target as HTMLElement)?.tagName
  const isInput = tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT'
  if (e.key === 'Escape') {
    if (!isInput && step.value > 0) {
      e.preventDefault()
      step.value--
    }
  } else if (e.key === 'Enter') {
    if (!isInput) {
      e.preventDefault()
      if (step.value < 3 && canProceed.value) {
        step.value++
      } else if (step.value === 3 && isValid.value) {
        publish()
      }
    }
  }
}

watch(
  () => ({
    artists: finalMetadata.artists.map(a => ({ ...a })),
    album: { ...finalMetadata.album },
    track: { ...finalMetadata.track },
  }),
  () => { hasUnsavedChanges.value = true },
  { deep: true }
)

onMounted(() => {
  loadDraft()
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  if (artistSearchTimer) clearTimeout(artistSearchTimer)
  if (albumSearchTimer) clearTimeout(albumSearchTimer)
  window.removeEventListener('keydown', handleKeydown)
})

onBeforeRouteLeave((_to, _from, next) => {
  if (hasUnsavedChanges.value) {
    const answer = window.confirm('You have unsaved changes. Leave anyway?')
    if (!answer) {
      next(false)
      return
    }
  }
  next()
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
