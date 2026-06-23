<script setup lang="ts">
import { computed, reactive, ref, watch, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import AutoComplete from 'primevue/autocomplete'
import Checkbox from 'primevue/checkbox'
import Select from 'primevue/select'
import { useToast } from 'primevue/usetoast'
import { parseBlob } from 'music-metadata-browser'
import { useLyricsApi } from '@/services/api/lyrics'
import { useArtistsApi } from '@/services/api/catalog/artists'
import { usePlayerApi } from '@/services/api/player'
import { useTracksApi } from '@/services/api/catalog/tracks'

type CatalogId = string | number

type CatalogOption = {
  id: CatalogId
  name: string
  slug?: string
  image_url?: string | null
  avatar_url?: string | null
  cover_url?: string | null
}

type AutoCompleteModel = CatalogOption | string | null

type CreditRole =
  | 'primary'
  | 'featured'
  | 'producer'
  | 'composer'
  | 'lyricist'
  | 'remixer'
  | 'arranger'
  | 'writer'
  | 'publisher'
  | string

type CreditPayload = {
  artist_id: CatalogId
  role: CreditRole
  position?: number
}

type CreditFormRow = {
  artist: CatalogOption | null
  role: CreditRole
}

export type TrackFormPayload = {
  title: string
  album_id: CatalogId | null
  album_artist_id: CatalogId | null

  duration_seconds: number | null
  track_number: number | null
  disc_number: number | null
  year: number | null

  genre_ids: CatalogId[]
  artist_ids: CatalogId[]
  featured_artist_ids: CatalogId[]
  credits: CreditPayload[]

  composer: string | null
  lyrics: string | null
  lyrics_language: string | null
  lyrics_type: 'plain' | 'synced'

  explicit: boolean
  isrc: string | null
  language: string | null
  release_date: string | null
  label: string | null

  cover_url: string | null
  audioFile: File | null
  coverFile: File | null
}

type TrackMetadataResult = {
  title?: string
  artists?: string[]
  album?: string
  albumArtist?: string
  genres?: string[]
  year?: number
  composer?: string
  lyrics?: string
  trackNumber?: number
  discNumber?: number
  durationSeconds?: number
  coverFile?: File
  coverPreviewUrl?: string
}

type ExistingTrack = Partial<{
  id: CatalogId
  title: string
  album_id: CatalogId | null
  album_artist_id: CatalogId | null

  duration_seconds: number | null
  track_number: number | null
  disc_number: number | null
  year: number | null

  composer: string | null
  lyrics: string | null
  lyrics_language: string | null
  lyrics_type: 'plain' | 'synced'

  explicit: boolean
  isrc: string | null
  language: string | null
  release_date: string | null
  label: string | null
  cover_url: string | null

  genre_ids: CatalogId[]
  artist_id: CatalogId | null
  artist_ids: CatalogId[]
  featured_artist_ids: CatalogId[]

  artists: Array<{
    artist_id: CatalogId
    role?: string
    is_primary?: boolean
    position?: number
  }>

  credits: Array<{
    artist_id: CatalogId
    role: string
    position?: number
  }>
}>

const props = defineProps<{
  visible: boolean
  loading?: boolean
  mode?: 'create' | 'edit'
  track?: ExistingTrack | null

  artistsOptions: CatalogOption[]
  albumsOptions: CatalogOption[]
  genresOptions: CatalogOption[]
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  submit: [payload: TrackFormPayload]
  cancel: []
}>()

const internalVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value),
})

const isEditMode = computed(() => props.mode === 'edit')

const form = reactive({
  title: '',
  duration_seconds: null as number | null,
  cover_url: null as string | null,

  track_number: null as number | null,
  disc_number: null as number | null,
  year: null as number | null,
  composer: '' as string | null,

  lyrics: '' as string | null,
  lyrics_language: 'en' as string | null,
  lyrics_type: 'plain' as 'plain' | 'synced',

  explicit: false,
  isrc: '' as string | null,
  language: '' as string | null,
  release_date: '' as string | null,
  label: '' as string | null,

  audioFile: null as File | null,
  coverFile: null as File | null,
})

const primaryArtistModels = ref<CatalogOption[]>([])
const featuredArtistModels = ref<CatalogOption[]>([])
const albumModel = ref<AutoCompleteModel>(null)
const albumArtistModel = ref<AutoCompleteModel>(null)
const genreModels = ref<CatalogOption[]>([])

const creditRows = ref<CreditFormRow[]>([])

const filteredArtists = ref<CatalogOption[]>([])
const filteredAlbums = ref<CatalogOption[]>([])
const filteredGenres = ref<CatalogOption[]>([])

const audioInputRef = ref<HTMLInputElement | null>(null)
const coverInputRef = ref<HTMLInputElement | null>(null)

const audioFileName = ref('')
const coverPreviewUrl = ref<string | null>(null)
const metadataLoading = ref(false)
const metadataError = ref('')
const aiGenerating = ref(false)
let aiPollTimer: ReturnType<typeof setInterval> | null = null

const toast = useToast()
const router = useRouter()
const lyricsApi = useLyricsApi()
const artistsApi = useArtistsApi()
const playerApi = usePlayerApi()
const tracksApi = useTracksApi()

const detectedArtistNames = ref<string[]>([])
const detectedAlbumName = ref('')
const detectedGenreNames = ref<string[]>([])
const detectedAlbumArtistName = ref('')

// ── API-backed artist search ──
const artistSearchLoading = ref(false)
let artistSearchTimer: ReturnType<typeof setTimeout> | null = null

// Track IDs whose artist data we've already fetched (hydration dedup)
const hydratedArtistIds = ref<Set<string>>(new Set())

// ── Audio preview (edit mode) ──
const audioPreviewUrl = ref<string | null>(null)
const audioPreviewPlaying = ref(false)
const audioPreviewRef = ref<HTMLAudioElement | null>(null)
const audioPreviewDuration = ref(0)
const audioPreviewCurrent = ref(0)

const creditRoleOptions = [
  { label: 'Producer', value: 'producer' },
  { label: 'Composer', value: 'composer' },
  { label: 'Lyricist', value: 'lyricist' },
  { label: 'Writer', value: 'writer' },
  { label: 'Arranger', value: 'arranger' },
  { label: 'Remixer', value: 'remixer' },
  { label: 'Publisher', value: 'publisher' },
]

const selectedAlbum = computed(() => optionFromModel(albumModel.value))
const selectedAlbumArtist = computed(() => optionFromModel(albumArtistModel.value))

const canSubmit = computed(() => {
  return form.title.trim().length > 0 && primaryArtistModels.value.length > 0
})

const unresolvedAlbumName = computed(() => {
  const typed = getModelText(albumModel.value) || detectedAlbumName.value
  if (!typed) return ''
  if (selectedAlbum.value) return ''
  if (findOptionByName(props.albumsOptions, typed)) return ''
  return typed
})

const unresolvedAlbumArtistName = computed(() => {
  const typed = getModelText(albumArtistModel.value) || detectedAlbumArtistName.value
  if (!typed) return ''
  if (selectedAlbumArtist.value) return ''
  if (findOptionByName(props.artistsOptions, typed)) return ''
  return typed
})

watch(
  () => props.visible,
  (visible) => {
    if (!visible) return

    resetForm()

    if (props.track) {
      hydrateFromTrack(props.track)
    }
  },
  { immediate: true },
)

watch(selectedAlbum, (album) => {
  // Auto-populate cover from the selected album if no cover is set yet
  if (album && !coverPreviewUrl.value && !form.coverFile) {
    const albumCover = album.cover_url || album.image_url || null
    if (albumCover) {
      form.cover_url = albumCover
      coverPreviewUrl.value = albumCover
    }
  }
})

watch(
  () => props.artistsOptions,
  () => {
    // Start with local options; API search will supplement
    filteredArtists.value = props.artistsOptions.slice(0, 20)
  },
  { immediate: true },
)

watch(
  () => props.albumsOptions,
  () => {
    filteredAlbums.value = props.albumsOptions.slice(0, 20)
  },
  { immediate: true },
)

watch(
  () => props.genresOptions,
  () => {
    filteredGenres.value = props.genresOptions.slice(0, 20)
  },
  { immediate: true },
)

function resetForm() {
  form.title = ''
  form.duration_seconds = null
  form.cover_url = null

  form.track_number = null
  form.disc_number = null
  form.year = null
  form.composer = ''

  form.lyrics = ''
  form.lyrics_language = 'en'
  form.lyrics_type = 'plain'

  form.explicit = false
  form.isrc = ''
  form.language = ''
  form.release_date = ''
  form.label = ''

  form.audioFile = null
  form.coverFile = null

  primaryArtistModels.value = []
  featuredArtistModels.value = []
  albumModel.value = null
  albumArtistModel.value = null
  genreModels.value = []
  creditRows.value = []

  audioFileName.value = ''
  coverPreviewUrl.value = null
  metadataLoading.value = false
  metadataError.value = ''

  detectedArtistNames.value = []
  detectedAlbumName.value = ''
  detectedGenreNames.value = []
  detectedAlbumArtistName.value = ''

  // Reset audio preview
  audioPreviewUrl.value = null
  audioPreviewPlaying.value = false
  audioPreviewCurrent.value = 0
  audioPreviewDuration.value = 0

  // Reset hydration tracker
  hydratedArtistIds.value = new Set()
}

function hydrateFromTrack(track: ExistingTrack) {
  form.title = track.title ?? ''
  form.duration_seconds = track.duration_seconds ?? null
  form.cover_url = track.cover_url ?? null

  form.track_number = track.track_number ?? null
  form.disc_number = track.disc_number ?? null
  form.year = track.year ?? null
  form.composer = track.composer ?? ''

  form.lyrics = track.lyrics ?? ''
  form.lyrics_language = track.lyrics_language ?? 'en'
  form.lyrics_type = track.lyrics_type ?? 'plain'

  form.explicit = Boolean(track.explicit)
  form.isrc = track.isrc ?? ''
  form.language = track.language ?? ''
  form.release_date = track.release_date ?? ''
  form.label = track.label ?? ''

  albumModel.value = findOptionById(props.albumsOptions, track.album_id ?? null)
  albumArtistModel.value = findOptionById(props.artistsOptions, track.album_artist_id ?? null)

  if (track.cover_url) {
    coverPreviewUrl.value = track.cover_url
  }

  const credits = track.credits?.length
    ? track.credits
    : normalizeLegacyArtistsToCredits(track)

  // Try to find primary/featured artists in local options
  primaryArtistModels.value = credits
    .filter((item) => item.role === 'primary')
    .sort((a, b) => Number(a.position ?? 0) - Number(b.position ?? 0))
    .map((item) => findOptionById(props.artistsOptions, item.artist_id))
    .filter((item): item is CatalogOption => Boolean(item))

  featuredArtistModels.value = credits
    .filter((item) => item.role === 'featured')
    .sort((a, b) => Number(a.position ?? 0) - Number(b.position ?? 0))
    .map((item) => findOptionById(props.artistsOptions, item.artist_id))
    .filter((item): item is CatalogOption => Boolean(item))

  creditRows.value = credits
    .filter((item) => item.role !== 'primary' && item.role !== 'featured')
    .map((item) => ({
      artist: findOptionById(props.artistsOptions, item.artist_id),
      role: item.role,
    }))
    .filter((item) => item.artist)

  genreModels.value = (track.genre_ids ?? [])
    .map((id) => findOptionById(props.genresOptions, id))
    .filter((item): item is CatalogOption => Boolean(item))

  // ── Fetch any artists not found in local options ──
  const allCreditIds = credits.map((c) => c.artist_id)
  hydrateMissingArtists(allCreditIds).then(() => {
    // Re-hydrate from credits now that we have more artist data
    const updatedCredits = credits

    if (primaryArtistModels.value.length === 0) {
      primaryArtistModels.value = updatedCredits
        .filter((item) => item.role === 'primary')
        .sort((a, b) => Number(a.position ?? 0) - Number(b.position ?? 0))
        .map((item) => findOptionById([...props.artistsOptions, ...filteredArtists.value], item.artist_id))
        .filter((item): item is CatalogOption => Boolean(item))
    }

    if (featuredArtistModels.value.length === 0) {
      featuredArtistModels.value = updatedCredits
        .filter((item) => item.role === 'featured')
        .sort((a, b) => Number(a.position ?? 0) - Number(b.position ?? 0))
        .map((item) => findOptionById([...props.artistsOptions, ...filteredArtists.value], item.artist_id))
        .filter((item): item is CatalogOption => Boolean(item))
    }

    // Re-hydrate credit rows that are still null
    creditRows.value = updatedCredits
      .filter((item) => item.role !== 'primary' && item.role !== 'featured')
      .map((item) => {
        const existing = creditRows.value.find(
          (r) => r.artist && String(r.artist.id) === String(item.artist_id),
        )
        return existing || {
          artist: findOptionById([...props.artistsOptions, ...filteredArtists.value], item.artist_id),
          role: item.role,
        }
      })
      .filter((item) => item.artist)
  })

  // ── Set up audio preview URL ──
  if (track.id) {
    const streamUrl = playerApi.getTrackStreamUrl(String(track.id))
    audioPreviewUrl.value = streamUrl
    audioPreviewCurrent.value = 0
    audioPreviewDuration.value = 0
    audioPreviewPlaying.value = false
  }
}

function normalizeLegacyArtistsToCredits(track: ExistingTrack) {
  const result: CreditPayload[] = []

  // Handle full artist objects from the API — supports both
  // camelCase (artistId) from the raw API and snake_case (artist_id)
  // from the ExistingTrack type.
  if (track.artists?.length) {
    for (const item of track.artists) {
      const artistId = (item as any).artist_id ?? (item as any).artistId
      if (artistId == null) continue
      result.push({
        artist_id: artistId,
        role: item.role || ((item as any).is_primary ? 'primary' : 'featured'),
        position: item.position,
      })
    }

    if (result.length > 0) return result
  }

  // Handle snake_case plural arrays (ExistingTrack shape)
  if (track.artist_ids?.length) {
    track.artist_ids.forEach((artistId, index) => {
      result.push({
        artist_id: artistId,
        role: 'primary',
        position: index,
      })
    })
  }

  if (track.featured_artist_ids?.length) {
    track.featured_artist_ids.forEach((artistId, index) => {
      result.push({
        artist_id: artistId,
        role: 'featured',
        position: index,
      })
    })
  }

  // Fallback: Track type uses singular artist_id + artist_name
  // (Zod drops the artists array, so neither track.artists nor
  // track.artist_ids are available).
  if (result.length === 0 && track.artist_id != null) {
    result.push({
      artist_id: track.artist_id,
      role: 'primary',
      position: 0,
    })
  }

  return result
}

function normalizeForSearch(value: string) {
  return value
    .toLowerCase()
    .trim()
    .replace(/\s+/g, ' ')
}

function normalizeText(value?: string | null) {
  if (!value) return undefined

  const normalized = value.trim().replace(/\s+/g, ' ')
  return normalized.length ? normalized : undefined
}

function getModelText(model: AutoCompleteModel) {
  if (!model) return ''
  if (typeof model === 'string') return model.trim()
  return model.name.trim()
}

function optionFromModel(model: AutoCompleteModel) {
  if (!model) return null
  if (typeof model === 'string') return null
  return model
}

function findOptionById(options: CatalogOption[], id?: CatalogId | null) {
  if (id === null || id === undefined) return null
  return options.find((item) => String(item.id) === String(id)) ?? null
}

function findOptionByName(options: CatalogOption[], name?: string | null) {
  if (!name) return null

  const normalized = normalizeForSearch(name)
  return options.find((item) => normalizeForSearch(item.name) === normalized) ?? null
}

function searchOptions(options: CatalogOption[], query?: string) {
  const normalized = normalizeForSearch(query ?? '')

  if (!normalized) {
    return options.slice(0, 20)
  }

  return options
    .filter((item) => normalizeForSearch(item.name).includes(normalized))
    .slice(0, 30)
}

function searchAlbums(event: { query: string }) {
  filteredAlbums.value = searchOptions(props.albumsOptions, event.query)
}

function searchGenres(event: { query: string }) {
  filteredGenres.value = searchOptions(props.genresOptions, event.query)
}

// ── Artist search with API fallback ──
async function searchArtists(event: { query: string }) {
  const q = (event.query || '').trim()
  const normalized = normalizeForSearch(q)

  // First, filter local options
  let local = q
    ? props.artistsOptions.filter((a) => normalizeForSearch(a.name).includes(normalized))
    : props.artistsOptions

  filteredArtists.value = local.slice(0, 30)

  // If local search is empty or we want more results, query the API
  if (q.length >= 2 || (!local.length && q.length > 0)) {
    // Debounce API calls
    if (artistSearchTimer) clearTimeout(artistSearchTimer)
    artistSearchTimer = setTimeout(async () => {
      artistSearchLoading.value = true
      try {
        const results = await artistsApi.searchArtists(q)
        if (Array.isArray(results)) {
          // Merge API results with local, deduplicate by ID
          const merged = [...filteredArtists.value]
          const seen = new Set(merged.map((a) => String(a.id)))
          for (const artist of results) {
            const opt = toCatalogOption(artist)
            if (opt && !seen.has(String(opt.id))) {
              merged.push(opt)
              seen.add(String(opt.id))
            }
          }
          filteredArtists.value = merged.slice(0, 30)
        }
      } catch {
        // API search failed — local results are already shown
      } finally {
        artistSearchLoading.value = false
      }
    }, 300)
  }
}

/** Convert raw Artist from API to CatalogOption */
function toCatalogOption(artist: Record<string, any>): CatalogOption | null {
  const id = artist.id ?? artist.artist_id
  const name = artist.name ?? artist.artist_name
  if (id == null || !name) return null
  return {
    id,
    name,
    slug: artist.slug,
    image_url: artist.image_url ?? artist.avatar_url ?? artist.cover_url ?? null,
    avatar_url: artist.avatar_url ?? artist.image_url ?? null,
    cover_url: artist.cover_url ?? artist.image_url ?? null,
  }
}

/** Fetch individual artists by ID to supplement the options during hydration */
async function hydrateMissingArtists(creditArtistIds: CatalogId[]) {
  const missingIds = creditArtistIds.filter(
    (id) => !findOptionById(props.artistsOptions, id) && !hydratedArtistIds.value.has(String(id)),
  )

  if (!missingIds.length) return

  const fetched: CatalogOption[] = []
  for (const id of missingIds) {
    try {
      const artist = await artistsApi.getArtist(String(id), { silent: true })
      if (artist) {
        const opt = toCatalogOption(artist as unknown as Record<string, any>)
        if (opt) {
          fetched.push(opt)
          hydratedArtistIds.value.add(String(id))
        }
      }
    } catch {
      // Artist may be deleted or inaccessible — skip silently
    }
  }

  if (!fetched.length) return

  // Add fetched artists to filtered artists for AutoComplete
  const existing = new Set(filteredArtists.value.map((a) => String(a.id)))
  for (const opt of fetched) {
    if (!existing.has(String(opt.id))) {
      filteredArtists.value.push(opt)
      existing.add(String(opt.id))
    }
  }
}

function splitArtists(value?: string | null) {
  if (!value) return []

  return value
    .replace(/\s+\((feat\.?|ft\.?|featuring)\s+/gi, ' feat. ')
    .replace(/\)$/g, '')
    .split(/\s+(?:feat\.?|ft\.?|featuring)\s+|,|&|;|\/|\+/gi)
    .map((item) => item.trim())
    .filter(Boolean)
}

function splitGenres(values?: string[] | string | null): string[] {
  if (!values) return []

  if (Array.isArray(values)) {
    return values
      .flatMap((item) => splitGenres(item))
      .map((item) => item.trim())
      .filter(Boolean)
  }

  return values
    .split(/,|;|\/|\|/g)
    .map((item) => item.trim())
    .filter(Boolean)
}

function uniqById(items: CatalogOption[]) {
  const seen = new Set<string>()

  return items.filter((item) => {
    const key = String(item.id)

    if (seen.has(key)) return false

    seen.add(key)
    return true
  })
}

function uniqCredits(items: CreditPayload[]) {
  const seen = new Set<string>()

  return items.filter((item) => {
    const key = `${String(item.artist_id)}:${item.role}`

    if (seen.has(key)) return false

    seen.add(key)
    return true
  })
}

function toOptionsByNames(options: CatalogOption[], names: string[]) {
  return uniqById(
    names
      .map((name) => findOptionByName(options, name))
      .filter((item): item is CatalogOption => Boolean(item)),
  )
}

function resolveAlbumInput() {
  const typed = getModelText(albumModel.value)
  if (!typed) return

  const found = findOptionByName(props.albumsOptions, typed)
  if (found) {
    albumModel.value = found
  }
}

function resolveAlbumArtistInput() {
  const typed = getModelText(albumArtistModel.value)
  if (!typed) return

  const found = findOptionByName(props.artistsOptions, typed)
  if (found) {
    albumArtistModel.value = found
  }
}

function triggerAudioInput() {
  audioInputRef.value?.click()
}

function triggerCoverInput() {
  coverInputRef.value?.click()
}

async function onAudioFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) return

  form.audioFile = file
  audioFileName.value = file.name

  await autoFillFromAudioFile(file)

  input.value = ''
}

function onCoverFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) return

  form.coverFile = file

  if (coverPreviewUrl.value?.startsWith('blob:')) {
    URL.revokeObjectURL(coverPreviewUrl.value)
  }

  coverPreviewUrl.value = URL.createObjectURL(file)

  input.value = ''
}

function removeCover() {
  form.coverFile = null
  form.cover_url = null

  if (coverPreviewUrl.value?.startsWith('blob:')) {
    URL.revokeObjectURL(coverPreviewUrl.value)
  }

  coverPreviewUrl.value = null
}

function imageMimeToExtension(mime: string) {
  switch (mime) {
    case 'image/jpeg':
      return 'jpg'
    case 'image/png':
      return 'png'
    case 'image/webp':
      return 'webp'
    default:
      return 'jpg'
  }
}

async function readAudioMetadata(file: File): Promise<TrackMetadataResult> {
  const metadata = await parseBlob(file)

  const common = metadata.common
  const format = metadata.format

  const firstPicture = common.picture?.[0]
  let coverFile: File | undefined
  let coverPreviewUrl: string | undefined

  if (firstPicture?.data?.length) {
    const extension = imageMimeToExtension(firstPicture.format)
    const coverBlob = new Blob([new Uint8Array(firstPicture.data)], {
      type: firstPicture.format,
    })

    coverFile = new File([coverBlob], `${file.name}-cover.${extension}`, {
      type: firstPicture.format,
    })

    coverPreviewUrl = URL.createObjectURL(coverFile)
  }

  const firstGenre = common.genre?.[0]

  return {
    title: normalizeText(common.title),
    artists: splitArtists(normalizeText(common.artist)),
    album: normalizeText(common.album),
    albumArtist: normalizeText(common.albumartist),
    genres: splitGenres(common.genre?.length ? common.genre : firstGenre ? [firstGenre] : []),
    year: common.year,
    composer: common.composer?.join(', '),
    lyrics: common.lyrics?.join('\n'),
    trackNumber: common.track?.no ?? undefined,
    discNumber: common.disk?.no ?? undefined,
    durationSeconds: format.duration ? Math.round(format.duration) : undefined,
    coverFile,
    coverPreviewUrl,
  }
}

async function autoFillFromAudioFile(file: File) {
  metadataLoading.value = true
  metadataError.value = ''

  try {
    const metadata = await readAudioMetadata(file)

    if (metadata.title && !form.title) {
      form.title = metadata.title
    }

    if (metadata.durationSeconds && !form.duration_seconds) {
      form.duration_seconds = metadata.durationSeconds
    }

    if (metadata.trackNumber && !form.track_number) {
      form.track_number = metadata.trackNumber
    }

    if (metadata.discNumber && !form.disc_number) {
      form.disc_number = metadata.discNumber
    }

    if (metadata.year && !form.year) {
      form.year = metadata.year
    }

    if (metadata.composer && !form.composer) {
      form.composer = metadata.composer
    }

    if (metadata.lyrics && !form.lyrics) {
      form.lyrics = metadata.lyrics
    }

    if (metadata.albumArtist) {
      detectedAlbumArtistName.value = metadata.albumArtist
      albumArtistModel.value =
        findOptionByName(props.artistsOptions, metadata.albumArtist) ?? metadata.albumArtist
    }

    if (metadata.artists?.length) {
      detectedArtistNames.value = metadata.artists

      const matched = toOptionsByNames(props.artistsOptions, metadata.artists)

      if (matched.length) {
        primaryArtistModels.value = [matched[0]!]
        featuredArtistModels.value = matched.slice(1)
      }
    }

    if (metadata.album) {
      detectedAlbumName.value = metadata.album
      albumModel.value = findOptionByName(props.albumsOptions, metadata.album) ?? metadata.album
    }

    if (metadata.genres?.length) {
      detectedGenreNames.value = metadata.genres
      genreModels.value = toOptionsByNames(props.genresOptions, metadata.genres)
    }

    if (metadata.coverFile && metadata.coverPreviewUrl) {
      form.coverFile = metadata.coverFile

      if (coverPreviewUrl.value?.startsWith('blob:')) {
        URL.revokeObjectURL(coverPreviewUrl.value)
      }

      coverPreviewUrl.value = metadata.coverPreviewUrl
    }
  } catch (error) {
    console.error(error)
    metadataError.value = 'Could not read metadata from this audio file.'
  } finally {
    metadataLoading.value = false
  }
}

function addCreditRow() {
  creditRows.value.push({
    artist: null,
    role: 'producer',
  })
}

function removeCreditRow(index: number) {
  creditRows.value.splice(index, 1)
}

function buildCreditsPayload() {
  const primaryCredits: CreditPayload[] = primaryArtistModels.value.map((artist, index) => ({
    artist_id: artist.id,
    role: 'primary',
    position: index,
  }))

  const featuredCredits: CreditPayload[] = featuredArtistModels.value.map((artist, index) => ({
    artist_id: artist.id,
    role: 'featured',
    position: index,
  }))

  const extraCredits: CreditPayload[] = creditRows.value
    .filter((row) => row.artist && String(row.role).trim())
    .map((row, index) => ({
      artist_id: row.artist!.id,
      role: String(row.role).trim(),
      position: index,
    }))

  return uniqCredits([...primaryCredits, ...featuredCredits, ...extraCredits])
}

function submitForm() {
  const title = form.title.trim()

  if (!title) return
  if (!primaryArtistModels.value.length) return

  const credits = buildCreditsPayload()

  const payload: TrackFormPayload = {
    title,

    album_id: selectedAlbum.value?.id ?? null,
    album_artist_id: selectedAlbumArtist.value?.id ?? null,

    duration_seconds: form.duration_seconds,
    track_number: form.track_number,
    disc_number: form.disc_number,
    year: form.year,

    genre_ids: genreModels.value.map((item) => item.id),
    artist_ids: primaryArtistModels.value.map((item) => item.id),
    featured_artist_ids: featuredArtistModels.value.map((item) => item.id),
    credits,

    composer: form.composer?.trim() || null,
    lyrics: form.lyrics?.trim() || null,
    lyrics_language: form.lyrics_language?.trim() || null,
    lyrics_type: form.lyrics_type ?? 'plain',

    explicit: form.explicit,
    isrc: form.isrc?.trim() || null,
    language: form.language?.trim() || null,
    release_date: form.release_date?.trim() || null,
    label: form.label?.trim() || null,

    cover_url: form.cover_url,
    audioFile: form.audioFile,
    coverFile: form.coverFile,
  }

  emit('submit', payload)
}

function stopAiPolling() {
  if (aiPollTimer) {
    clearInterval(aiPollTimer)
    aiPollTimer = null
  }
}

async function fillLyricsFromResult(content: string, type: 'synced' | 'plain', trackName: string) {
  form.lyrics = content
  form.lyrics_type = type
  toast.add({
    severity: type === 'synced' ? 'success' : 'info',
    summary: type === 'synced'
      ? `Synced LRC fetched for "${trackName}"`
      : `Plain lyrics fetched for "${trackName}"`,
    life: 3000,
  })
}

const LYRICS_TIMEOUT_MS = 8000

async function fetchLRCLyrics() {
  const trackName = form.title?.trim()
  const artistName = primaryArtistModels.value[0]?.name || detectedArtistNames.value[0] || ''
  if (!trackName) {
    toast.add({ severity: 'warn', summary: 'Enter a track title first', life: 2500 })
    return
  }

  metadataLoading.value = true
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), LYRICS_TIMEOUT_MS)

  try {
    const params = new URLSearchParams({ track_name: trackName })
    if (artistName) params.set('artist_name', artistName)
    if (form.duration_seconds) params.set('duration', String(Math.round(form.duration_seconds)))

    // Try exact /get first
    let data: any = null
    try {
      const resp = await fetch(`https://lrclib.net/api/get?${params}`, {
        headers: { Accept: 'application/json' },
        signal: controller.signal,
      })
      if (resp.ok) {
        data = await resp.json()
      }
    } catch {
      // /get failed (timeout / CORS / network) — clear data so we try /search
      data = null
    }

    // Fallback to /search if exact fails
    if (!data || (!data.syncedLyrics && !data.plainLyrics)) {
      const q = artistName ? `${artistName} ${trackName}` : trackName
      try {
        const searchController = new AbortController()
        const searchTimeout = setTimeout(() => searchController.abort(), LYRICS_TIMEOUT_MS)
        const resp = await fetch(`https://lrclib.net/api/search?q=${encodeURIComponent(q)}`, {
          headers: { Accept: 'application/json' },
          signal: searchController.signal,
        })
        clearTimeout(searchTimeout)
        if (resp.ok) {
          const results: any[] = await resp.json()
          // Prefer synced, then any with lyrics
          data = results.find((r: any) => r.syncedLyrics) || results.find((r: any) => r.plainLyrics) || null
        } else {
          data = null
        }
      } catch {
        data = null
      }
    }

    if (data && (data.syncedLyrics || data.plainLyrics)) {
      const synced = data.syncedLyrics as string | undefined
      if (synced) {
        await fillLyricsFromResult(synced, 'synced', trackName)
      } else if (data.plainLyrics) {
        await fillLyricsFromResult(data.plainLyrics, 'plain', trackName)
      }

      // Auto-select album from LRCLIB response
      if (data.albumName) {
        const matchedAlbum = findOptionByName(props.albumsOptions, data.albumName)
        if (matchedAlbum) {
          albumModel.value = matchedAlbum

          // Auto-populate cover from the matched album if none set
          if (!coverPreviewUrl.value && !form.coverFile) {
            const albumCover = matchedAlbum.cover_url || matchedAlbum.image_url || null
            if (albumCover) {
              form.cover_url = albumCover
              coverPreviewUrl.value = albumCover
            }
          }
        } else {
          detectedAlbumName.value = data.albumName
        }
      }

      // Auto-select primary artist from LRCLIB response (only if none selected)
      if (data.artistName && !primaryArtistModels.value.length) {
        const matchedArtist = findOptionByName(props.artistsOptions, data.artistName)
        if (matchedArtist) {
          primaryArtistModels.value = [matchedArtist]
        } else if (!detectedArtistNames.value.length) {
          detectedArtistNames.value = [data.artistName]
        }
      }

      // Try to fetch cover from iTunes API if we still don't have one
      if (!coverPreviewUrl.value && !form.coverFile) {
        fetchCoverFromiTunes().catch(() => {})
      }
      return
    }

    // ── LRCLIB returned nothing or timed out ── try server-side pipeline (LRCLIB → AI fallback)
    await fetchLyricsViaPipeline(trackName)
  } catch (err) {
    // Only show error for non-abort (timeout) errors — timeouts go to pipeline
    if (err instanceof DOMException && err.name === 'AbortError') {
      await fetchLyricsViaPipeline(trackName)
    } else {
      toast.add({ severity: 'error', summary: 'Failed to fetch lyrics', detail: String(err), life: 3000 })
    }
  } finally {
    clearTimeout(timeout)
    metadataLoading.value = false
  }
}

/**
 * Calls the server-side fetch-or-generate endpoint.
 * If the server enqueues an AI job, starts polling for the result.
 */
async function fetchLyricsViaPipeline(trackName: string) {
  // If we're in edit mode and have a track ID, use the server pipeline
  if (!props.track?.id) {
    toast.add({ severity: 'warn', summary: 'Save the track first, then try Fetch LRC again for AI-powered lyrics', life: 5000 })
    return
  }

  try {
    const result = await lyricsApi.fetchOrGenerateLyrics(props.track.id)

    if (result?.source === 'lrclib' && result?.data?.content) {
      await fillLyricsFromResult(
        result.data.content,
        result.data.type === 'lrc' ? 'synced' : 'plain',
        trackName,
      )
      return
    }

    if (result?.source === 'ai' && result?.job_id) {
      // AI job enqueued — start polling the AI status endpoint for progress
      aiGenerating.value = true
      toast.add({
        severity: 'info',
        summary: 'AI lyrics generation started',
        detail: 'This may take 30-60 seconds. We\'ll fill the form automatically when ready.',
        life: 8000,
      })

      const trackId = result.track_id || String(props.track.id)
      let attempts = 0
      const maxAttempts = 30 // ~2.5 minutes at 5s intervals

      stopAiPolling()
      aiPollTimer = setInterval(async () => {
        attempts++
        try {
          const statusResult = await lyricsApi.aiStatus(trackId)

          if (statusResult?.status === 'completed') {
            // AI finished — lyrics should be in DB now
            stopAiPolling()
            aiGenerating.value = false
            // Fetch the actual lyrics
            const lyricsResult = await lyricsApi.getTrackLyrics(trackId, undefined, { silent: true })
            if (lyricsResult?.content) {
              await fillLyricsFromResult(
                lyricsResult.content,
                lyricsResult.type === 'lrc' ? 'synced' : 'plain',
                trackName,
              )
            }
            toast.add({
              severity: 'success',
              summary: 'AI lyrics generated!',
              life: 3000,
            })
            return
          }

          if (statusResult?.status === 'failed') {
            stopAiPolling()
            aiGenerating.value = false
            toast.add({
              severity: 'error',
              summary: 'AI lyrics generation failed',
              detail: statusResult?.message || 'Unknown error',
              life: 5000,
            })
            return
          }

          // Still in progress — update status message
          if (statusResult?.message) {
            // The button already shows "AI Generating..." - update the toast
            toast.add({
              severity: 'info',
              summary: statusResult.message,
              life: 2000,
            })
          }
        } catch (err) {
          // Status endpoint unavailable — keep polling
        }

        if (attempts >= maxAttempts) {
          stopAiPolling()
          aiGenerating.value = false
          toast.add({
            severity: 'warn',
            summary: 'AI lyrics generation taking longer than expected',
            detail: 'The lyrics will be saved automatically. You can close this dialog and check back later.',
            life: 5000,
          })
        }
      }, 5000)
      return
    }

    // source === "none"
    toast.add({ severity: 'error', summary: result?.message || 'No lyrics found on LRCLIB and AI generation not available', life: 3000 })
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Failed to fetch or generate lyrics', detail: String(err?.message || err), life: 3000 })
  }
}

// Clean up polling on unmount
onUnmounted(() => {
  stopAiPolling()
})

async function fetchCoverFromiTunes() {
  const trackName = form.title?.trim()
  const artistName = primaryArtistModels.value[0]?.name || detectedArtistNames.value[0] || ''
  const albumName = getModelText(albumModel.value) || detectedAlbumName.value || ''
  if (!trackName && !artistName && !albumName) return

  try {
    const terms = [artistName, albumName || trackName].filter(Boolean).join(' ')
    const url = `https://itunes.apple.com/search?term=${encodeURIComponent(terms)}&entity=song&limit=5`
    const resp = await fetch(url, { headers: { Accept: 'application/json' } })
    if (!resp.ok) return

    const data: any = await resp.json()
    const result = data.results?.[0]
    if (!result?.artworkUrl100) return

    // Get larger artwork by replacing 100x100 with larger size
    const coverUrl = result.artworkUrl100.replace('/100x100bb.', '/600x600bb.')
    form.cover_url = coverUrl
    if (coverPreviewUrl.value?.startsWith('blob:')) {
      URL.revokeObjectURL(coverPreviewUrl.value)
    }
    coverPreviewUrl.value = coverUrl
  } catch {
    // Silent fail — cover fetch is best-effort
  }
}

function closeDialog() {
  // Stop audio preview if playing
  if (audioPreviewPlaying.value) {
    audioPreviewPlaying.value = false
  }
  audioPreviewUrl.value = null
  emit('cancel')
  internalVisible.value = false
}

// ── Audio preview ──
function toggleAudioPreview() {
  const audio = audioPreviewRef.value
  if (!audio || !audioPreviewUrl.value) return

  if (audio.paused) {
    audio.play().catch(() => {
      audioPreviewPlaying.value = false
    })
    audioPreviewPlaying.value = true
  } else {
    audio.pause()
    audioPreviewPlaying.value = false
  }
}

function onAudioPreviewTimeUpdate(event: Event) {
  const audio = event.target as HTMLAudioElement
  audioPreviewCurrent.value = audio.currentTime
}

function onAudioPreviewMetadata(event: Event) {
  const audio = event.target as HTMLAudioElement
  audioPreviewDuration.value = audio.duration || 0
}

function onAudioPreviewEnded() {
  audioPreviewPlaying.value = false
  audioPreviewCurrent.value = 0
}

function seekAudioPreview(event: Event) {
  const input = event.target as HTMLInputElement
  const value = Number(input.value)
  const audio = audioPreviewRef.value
  if (audio) {
    audio.currentTime = value
    audioPreviewCurrent.value = value
  }
}

function formatDuration(seconds: number): string {
  if (!seconds || !isFinite(seconds)) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}
</script>

<template>
  <Dialog
    v-model:visible="internalVisible"
    modal
    :header="isEditMode ? 'Edit track' : 'Create track'"
    class="w-[95vw] max-w-5xl"
    content-class="!bg-[#121212] !text-white"
    header-class="!bg-[#121212] !text-white"
  >
    <form class="space-y-6" @submit.prevent="submitForm">
      <!-- Uploads + Audio Preview -->
      <section class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <!-- Audio file upload -->
        <div class="rounded-2xl border border-white/[0.08] bg-white/[0.03] p-4">
          <label class="mb-2 block text-xs font-medium text-slate-400">
            Audio file
          </label>

          <input
            ref="audioInputRef"
            type="file"
            accept="audio/*"
            hidden
            @change="onAudioFileChange"
          />

          <Button
            type="button"
            icon="pi pi-upload"
            :label="audioFileName || 'Choose audio file'"
            class="!rounded-xl"
            outlined
            @click="triggerAudioInput"
          />

          <p v-if="metadataLoading" class="mt-2 text-xs text-slate-400">
            Reading metadata...
          </p>

          <p v-if="metadataError" class="mt-2 text-xs text-red-400">
            {{ metadataError }}
          </p>

          <!-- Audio preview (edit mode: existing track) -->
          <div v-if="audioPreviewUrl" class="mt-4 border-t border-white/[0.06] pt-3">
            <div class="mb-2 flex items-center gap-2">
              <Button
                type="button"
                :icon="audioPreviewPlaying ? 'pi pi-pause-circle' : 'pi pi-play-circle'"
                severity="secondary"
                text
                rounded
                aria-label="Preview track audio"
                class="!text-emerald-400 !text-xl"
                @click="toggleAudioPreview"
              />
              <span class="text-xs text-slate-500">
                Preview
              </span>
              <span class="ml-auto text-xs tabular-nums text-slate-500">
                {{ formatDuration(audioPreviewCurrent) }} / {{ formatDuration(audioPreviewDuration) }}
              </span>
            </div>

            <!-- Simple seek bar -->
            <input
              type="range"
              :min="0"
              :max="audioPreviewDuration || 0"
              :value="audioPreviewCurrent"
              class="h-1 w-full cursor-pointer appearance-none rounded-full bg-white/[0.08] accent-emerald-500"
              @input="seekAudioPreview"
              aria-label="Seek audio preview"
            />

            <audio
              ref="audioPreviewRef"
              :src="audioPreviewUrl"
              preload="metadata"
              hidden
              @timeupdate="onAudioPreviewTimeUpdate"
              @loadedmetadata="onAudioPreviewMetadata"
              @ended="onAudioPreviewEnded"
            />
          </div>
        </div>

        <!-- Cover image -->
        <div class="rounded-2xl border border-white/[0.08] bg-white/[0.03] p-4">
          <label class="mb-2 block text-xs font-medium text-slate-400">
            Cover image
          </label>

          <input
            ref="coverInputRef"
            type="file"
            accept="image/*"
            hidden
            @change="onCoverFileChange"
          />

          <div class="flex items-center gap-4">
            <div
              class="flex h-24 w-24 items-center justify-center overflow-hidden rounded-xl bg-white/[0.06]"
            >
              <img
                v-if="coverPreviewUrl"
                :src="coverPreviewUrl"
                alt="Cover"
                class="h-full w-full object-cover"
              />
              <i aria-hidden="true" v-else class="pi pi-image text-2xl text-slate-500" />
            </div>

            <div class="flex flex-col gap-2">
              <Button
                type="button"
                icon="pi pi-image"
                label="Choose cover"
                class="!rounded-xl"
                outlined
                @click="triggerCoverInput"
              />

              <Button
                v-if="coverPreviewUrl"
                type="button"
                icon="pi pi-trash"
                label="Remove"
                text
                class="!text-red-400"
                @click="removeCover"
              />
            </div>
          </div>
        </div>
      </section>

      <!-- Main metadata -->
      <section class="space-y-4">
        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-400">
            Title
          </label>

          <InputText
            v-model="form.title"
            placeholder="Track title"
            class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
          />
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1.5 flex items-center gap-2 text-xs font-medium text-slate-400">
              Primary artists
              <i
                v-if="artistSearchLoading"
                class="pi pi-spin pi-spinner text-[10px] text-emerald-400"
              />
            </label>

            <AutoComplete
              v-model="primaryArtistModels"
              :suggestions="filteredArtists"
              optionLabel="name"
              multiple
              dropdown
              completeOnFocus
              placeholder="Search primary artists"
              class="w-full"
              input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
              panel-class="!bg-[#181818] !border-white/[0.08]"
              @complete="searchArtists"
            >
              <template #option="{ option }">
                <div class="flex items-center gap-2 text-sm text-white">
                  <i aria-hidden="true" class="pi pi-user text-xs text-slate-500" />
                  <span>{{ option.name }}</span>
                </div>
              </template>
            </AutoComplete>

            <p
              v-if="detectedArtistNames.length && !primaryArtistModels.length"
              class="mt-1 text-xs text-yellow-400"
            >
              Detected:
              {{ detectedArtistNames.join(', ') }}
            </p>
          </div>

          <div>
            <label class="mb-1.5 flex items-center gap-2 text-xs font-medium text-slate-400">
              Featured artists
              <i
                v-if="artistSearchLoading"
                class="pi pi-spin pi-spinner text-[10px] text-emerald-400"
              />
            </label>

            <AutoComplete
              v-model="featuredArtistModels"
              :suggestions="filteredArtists"
              optionLabel="name"
              multiple
              dropdown
              completeOnFocus
              placeholder="Search featured artists"
              class="w-full"
              input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
              panel-class="!bg-[#181818] !border-white/[0.08]"
              @complete="searchArtists"
            >
              <template #option="{ option }">
                <div class="flex items-center gap-2 text-sm text-white">
                  <i aria-hidden="true" class="pi pi-user-plus text-xs text-slate-500" />
                  <span>{{ option.name }}</span>
                </div>
              </template>
            </AutoComplete>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">
              Album
            </label>

            <AutoComplete
              v-model="albumModel"
              :suggestions="filteredAlbums"
              optionLabel="name"
              placeholder="Search album"
              dropdown
              completeOnFocus
              class="w-full"
              input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
              panel-class="!bg-[#181818] !border-white/[0.08]"
              @complete="searchAlbums"
              @blur="resolveAlbumInput"
            />

            <p v-if="unresolvedAlbumName" class="mt-1 text-xs text-yellow-400">
              Album not found:
              {{ unresolvedAlbumName }}
            </p>
          </div>

          <div>
            <label class="mb-1.5 flex items-center gap-2 text-xs font-medium text-slate-400">
              Album artist
              <i
                v-if="artistSearchLoading"
                class="pi pi-spin pi-spinner text-[10px] text-emerald-400"
              />
            </label>

            <AutoComplete
              v-model="albumArtistModel"
              :suggestions="filteredArtists"
              optionLabel="name"
              placeholder="Search album artist"
              dropdown
              completeOnFocus
              class="w-full"
              input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
              panel-class="!bg-[#181818] !border-white/[0.08]"
              @complete="searchArtists"
              @blur="resolveAlbumArtistInput"
            />

            <p
              v-if="unresolvedAlbumArtistName"
              class="mt-1 text-xs text-yellow-400"
            >
              Album artist not found:
              {{ unresolvedAlbumArtistName }}
            </p>
          </div>
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-400">
            Genres
          </label>

          <AutoComplete
            v-model="genreModels"
            :suggestions="filteredGenres"
            optionLabel="name"
            multiple
            dropdown
            completeOnFocus
            placeholder="Search genres"
            class="w-full"
            input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
            panel-class="!bg-[#181818] !border-white/[0.08]"
            @complete="searchGenres"
          >
            <template #option="{ option }">
              <div class="flex items-center gap-2 text-sm text-white">
                <i aria-hidden="true" class="pi pi-tag text-xs text-slate-500" />
                <span>{{ option.name }}</span>
              </div>
            </template>
          </AutoComplete>

          <p
            v-if="detectedGenreNames.length && !genreModels.length"
            class="mt-1 text-xs text-yellow-400"
          >
            Detected:
            {{ detectedGenreNames.join(', ') }}
          </p>
        </div>
      </section>

      <!-- Numbers -->
      <section class="grid grid-cols-1 gap-4 md:grid-cols-4">
        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-400">
            Duration seconds
          </label>

          <InputNumber
            v-model="form.duration_seconds"
            class="w-full"
            input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
          />
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-400">
            Track number
          </label>

          <InputNumber
            v-model="form.track_number"
            class="w-full"
            input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
          />
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-400">
            Disc number
          </label>

          <InputNumber
            v-model="form.disc_number"
            class="w-full"
            input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
          />
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-400">
            Year
          </label>

          <InputNumber
            v-model="form.year"
            class="w-full"
            input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
          />
        </div>
      </section>

      <!-- Industry metadata -->
      <section class="space-y-4">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">
              ISRC
            </label>

            <InputText
              v-model="form.isrc"
              placeholder="e.g. USRC17607839"
              class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">
              Language
            </label>

            <InputText
              v-model="form.language"
              placeholder="en, fa, fr..."
              class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">
              Release date
            </label>

            <InputText
              v-model="form.release_date"
              placeholder="YYYY-MM-DD"
              class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">
              Label
            </label>

            <InputText
              v-model="form.label"
              placeholder="Record label"
              class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
            />
          </div>
        </div>

        <div class="flex items-center gap-2">
          <Checkbox v-model="form.explicit" binary inputId="explicit" />

          <label for="explicit" class="text-sm text-slate-300">
            Explicit content
          </label>
        </div>
      </section>

      <!-- Composer -->
      <section>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">
          Composer
        </label>

        <InputText
          v-model="form.composer"
          placeholder="Composer names"
          class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
        />
      </section>

      <!-- Credits -->
      <section class="space-y-3 rounded-2xl border border-white/[0.08] bg-white/[0.03] p-4">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-sm font-semibold text-white">Credits</h3>
            <p class="text-xs text-slate-400">
              Add producer, composer, lyricist, remixer, and other contributors.
            </p>
          </div>

          <Button
            type="button"
            icon="pi pi-plus"
            label="Add credit"
            size="small"
            outlined
            class="!rounded-lg"
            @click="addCreditRow"
          />
        </div>

        <div
          v-for="(row, index) in creditRows"
          :key="index"
          class="grid grid-cols-1 gap-3 rounded-xl border border-white/[0.06] bg-black/20 p-3 md:grid-cols-[1fr_220px_44px]"
        >
          <AutoComplete
            v-model="row.artist"
            :suggestions="filteredArtists"
            optionLabel="name"
            placeholder="Artist"
            dropdown
            completeOnFocus
            class="w-full"
            input-class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
            panel-class="!bg-[#181818] !border-white/[0.08]"
            @complete="searchArtists"
          />

          <Select
            v-model="row.role"
            :options="creditRoleOptions"
            optionLabel="label"
            optionValue="value"
            editable
            placeholder="Role"
            class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
            panel-class="!bg-[#181818] !border-white/[0.08]"
          />

          <Button
            type="button"
            icon="pi pi-trash"
            text
            class="!text-slate-500 hover:!text-red-400"
            @click="removeCreditRow(index)"
          />
        </div>

        <p v-if="!creditRows.length" class="text-xs text-slate-500">
          No extra credits added.
        </p>
      </section>

      <!-- Lyrics -->
      <section class="space-y-4">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">
              Lyrics language
            </label>

            <InputText
              v-model="form.lyrics_language"
              placeholder="en"
              class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
            />
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">
              Lyrics type
            </label>

            <Select
              v-model="form.lyrics_type"
              :options="[
                { label: 'Plain', value: 'plain' },
                { label: 'Synced', value: 'synced' },
              ]"
              optionLabel="label"
              optionValue="value"
              class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
              panel-class="!bg-[#181818] !border-white/[0.08]"
            />
          </div>
        </div>

        <div>
          <div class="mb-1.5 flex items-center justify-between">
            <label class="block text-xs font-medium text-slate-400">
              Lyrics
            </label>

            <Button
              type="button"
              icon="pi pi-cloud-download"
              :label="aiGenerating ? 'AI Generating...' : 'Fetch LRC'"
              size="small"
              text
              :loading="metadataLoading || aiGenerating"
              class="!text-teal-400 hover:!text-teal-300"
              @click="fetchLRCLyrics"
            />
          </div>

          <Textarea
            v-model="form.lyrics"
            rows="6"
            autoResize
            placeholder="Track lyrics..."
            class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white"
          />
        </div>
      </section>

      <!-- Actions -->
      <footer class="flex items-center justify-end gap-3 border-t border-white/[0.08] pt-4">
        <Button
          type="button"
          label="Cancel"
          text
          class="!text-slate-300"
          @click="closeDialog"
        />

        <Button
          type="submit"
          icon="pi pi-check"
          :label="isEditMode ? 'Save changes' : 'Create track'"
          :loading="props.loading"
          :disabled="!canSubmit || props.loading"
          class="!rounded-xl"
        />
      </footer>
    </form>
  </Dialog>
</template>
