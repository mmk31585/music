
<template>
  <div class="rounded-3xl border border-white/10 bg-white/5 p-6 shadow-lg">
    <div class="mb-6">
      <h2 class="text-xl font-bold text-white">Upload track</h2>
      <p class="mt-2 text-sm text-slate-400">
        Upload an audio file, optional cover image, and save it as a track in the catalog.
      </p>
    </div>

    <div class="rounded-2xl border border-dashed border-white/15 bg-black/20 p-6">
      <div class="flex flex-col gap-5">
        <!-- Track title -->
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300">
            Track title
          </label>

          <InputText
            v-model="form.title"
            placeholder="Enter track title"
            class="w-full"
            :disabled="submitting"
          />

          <p v-if="showTitleError" class="mt-2 text-xs text-red-300">
            Track title is required.
          </p>
        </div>

        <!-- Artist -->
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300">
            Artist
          </label>

          <Dropdown
            v-model="form.artist_id"
            :options="artistOptions"
            option-label="label"
            option-value="value"
            placeholder="Select artist"
            class="w-full"
            :loading="artistsLoading"
            :disabled="submitting"
            show-clear
          />

          <p v-if="showArtistError" class="mt-2 text-xs text-red-300">
            Artist is required.
          </p>
        </div>

        <!-- Album -->
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300">
            Album
          </label>

          <Dropdown
            v-model="form.album_id"
            :options="albumOptions"
            option-label="label"
            option-value="value"
            placeholder="Select album"
            class="w-full"
            :loading="albumsLoading"
            :disabled="submitting"
            show-clear
          />
        </div>

        <!-- Genre -->
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300">
            Genre
          </label>

          <Dropdown
            v-model="form.genre_id"
            :options="genreOptions"
            option-label="label"
            option-value="value"
            placeholder="Select genre"
            class="w-full"
            :loading="genresLoading"
            :disabled="submitting"
            show-clear
          />
        </div>

        <!-- Duration -->
        <div>
          <div class="mb-2 flex items-center justify-between gap-3">
            <label class="block text-sm font-medium text-slate-300">
              Duration seconds
            </label>

            <span v-if="durationStatus" class="text-xs text-slate-500">
              {{ durationStatus }}
            </span>
          </div>

          <InputNumber
            v-model="form.duration_seconds"
            placeholder="Auto detected from selected audio"
            class="w-full"
            input-class="w-full"
            :min="0"
            :disabled="submitting || detectingDuration"
          />

          <p v-if="formattedDuration" class="mt-2 text-xs text-slate-400">
            Duration: {{ formattedDuration }}
          </p>
        </div>

        <!-- Audio file -->
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300">
            Choose audio file
          </label>

          <input
            ref="audioInputRef"
            type="file"
            accept="audio/*"
            class="block w-full rounded-xl border border-white/10 bg-white/5 px-3 py-3 text-sm text-slate-300 file:mr-4 file:rounded-lg file:border-0 file:bg-[#1db954] file:px-4 file:py-2 file:text-sm file:font-semibold file:text-black disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="submitting"
            @change="onAudioFileChange"
          />

          <p v-if="showFileError" class="mt-2 text-xs text-red-300">
            Audio file is required.
          </p>
        </div>

        <div
          v-if="audioFileName"
          class="rounded-xl bg-white/5 px-4 py-3 text-sm text-slate-300"
        >
          Audio:
          <span class="font-medium text-white">{{ audioFileName }}</span>

          <span v-if="audioFileSizeLabel" class="ml-2 text-slate-500">
            {{ audioFileSizeLabel }}
          </span>
        </div>

        <!-- Cover image -->
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300">
            Cover image
            <span class="text-xs text-slate-500">(optional)</span>
          </label>

          <div class="grid gap-4 md:grid-cols-[140px_1fr] md:items-start">
            <div
              class="flex aspect-square w-full max-w-[140px] items-center justify-center overflow-hidden rounded-2xl border border-white/10 bg-white/5"
            >
              <img
                v-if="coverPreviewUrl"
                :src="coverPreviewUrl"
                alt="Cover preview"
                class="h-full w-full object-cover"
              />

              <div v-else class="px-4 text-center text-xs text-slate-500">
                No cover selected
              </div>
            </div>

            <div class="flex flex-col gap-3">
              <input
                ref="coverInputRef"
                type="file"
                accept="image/*"
                class="block w-full rounded-xl border border-white/10 bg-white/5 px-3 py-3 text-sm text-slate-300 file:mr-4 file:rounded-lg file:border-0 file:bg-white/10 file:px-4 file:py-2 file:text-sm file:font-semibold file:text-white disabled:cursor-not-allowed disabled:opacity-60"
                :disabled="submitting"
                @change="onCoverFileChange"
              />

              <div
                v-if="coverFileName"
                class="rounded-xl bg-white/5 px-4 py-3 text-sm text-slate-300"
              >
                Cover:
                <span class="font-medium text-white">{{ coverFileName }}</span>

                <span v-if="coverFileSizeLabel" class="ml-2 text-slate-500">
                  {{ coverFileSizeLabel }}
                </span>
              </div>

              <Button
                v-if="selectedCoverFile"
                label="Remove cover"
                icon="pi pi-times"
                severity="secondary"
                outlined
                size="small"
                :disabled="submitting"
                @click="removeCover"
              />
            </div>
          </div>
        </div>

        <!-- Upload progress -->
        <div v-if="uploading || totalProgress > 0" class="flex flex-col gap-2">
          <div class="flex items-center justify-between text-xs text-slate-400">
            <span>{{ uploadStatusLabel }}</span>
            <span>{{ totalProgress }}%</span>
          </div>

          <ProgressBar :value="totalProgress" />
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <Button
            label="Upload and create track"
            icon="pi pi-upload"
            :loading="submitting"
            :disabled="!canSubmit"
            class="border-0 bg-[#1db954] text-black"
            @click="submit"
          />

          <Button
            label="Clear"
            severity="secondary"
            outlined
            :disabled="submitting"
            @click="clear"
          />
        </div>

        <div
          v-if="createdTrack"
          class="rounded-2xl border border-emerald-500/20 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-300"
        >
          Track created:
          <span class="font-semibold text-emerald-200">
            {{ createdTrack.title }}
          </span>
        </div>

        <div
          v-if="resultMessage"
          class="rounded-2xl border border-emerald-500/20 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-300"
        >
          {{ resultMessage }}
        </div>

        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-300"
        >
          {{ errorMessage }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import Button from 'primevue/button'
import ProgressBar from 'primevue/progressbar'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Dropdown from 'primevue/dropdown'
import { useToast } from 'primevue/usetoast'

import { useMediaApi } from '@/services/api/media'
import type { Track, TrackFormPayload } from '@/services/api/catalog'

import {
  useAdminTracks,
  useAdminArtists,
  useAdminAlbums,
  useAdminGenres,
} from '@/composables/admin'

const AUDIO_MEDIA_TYPE = 'trackAudio'
const COVER_MEDIA_TYPE = 'trackCover'

const toast = useToast()
const mediaApi = useMediaApi()

const { createTrack, saving: trackSaving } = useAdminTracks()
const { artists, loading: artistsLoading, fetchArtists } = useAdminArtists()
const { albums, loading: albumsLoading, fetchAlbums } = useAdminAlbums()
const { genres, loading: genresLoading, fetchGenres } = useAdminGenres()

const audioInputRef = ref<HTMLInputElement | null>(null)
const coverInputRef = ref<HTMLInputElement | null>(null)

const selectedAudioFile = ref<File | null>(null)
const audioFileName = ref('')
const audioFileSizeLabel = ref('')

const selectedCoverFile = ref<File | null>(null)
const coverFileName = ref('')
const coverFileSizeLabel = ref('')
const coverPreviewUrl = ref('')

const uploading = ref(false)
const uploadStep = ref<'idle' | 'audio' | 'cover' | 'track'>('idle')
const audioUploadProgress = ref(0)
const coverUploadProgress = ref(0)

const detectingDuration = ref(false)
const submitted = ref(false)

const resultMessage = ref('')
const errorMessage = ref('')
const createdTrack = ref<Track | null>(null)

let audioObjectUrl: string | null = null
let coverObjectUrl: string | null = null

const form = reactive<TrackFormPayload>({
  title: '',
  artist_id: null,
  album_id: null,
  genre_id: null,
  duration_seconds: null,
  audio_url: null,
  cover_url: null,
})

const submitting = computed(() => uploading.value || trackSaving.value)

const hasCoverUpload = computed(() => Boolean(selectedCoverFile.value))

const totalProgress = computed(() => {
  if (!uploading.value) {
    return Math.max(audioUploadProgress.value, coverUploadProgress.value)
  }

  if (!hasCoverUpload.value) {
    return audioUploadProgress.value
  }

  // Audio is 50%, cover is 50%.
  return Math.round((audioUploadProgress.value + coverUploadProgress.value) / 2)
})

const uploadStatusLabel = computed(() => {
  if (uploadStep.value === 'audio') return 'Uploading audio...'
  if (uploadStep.value === 'cover') return 'Uploading cover...'
  if (uploadStep.value === 'track') return 'Creating track...'
  if (totalProgress.value >= 100) return 'Upload complete'

  return 'Preparing upload...'
})

const canSubmit = computed(() => {
  return Boolean(
    selectedAudioFile.value &&
    form.title &&
    String(form.title).trim().length > 0 &&
    form.artist_id &&
    !submitting.value &&
    !detectingDuration.value,
  )
})

const showTitleError = computed(() => {
  return submitted.value && !String(form.title || '').trim()
})

const showArtistError = computed(() => {
  return submitted.value && !form.artist_id
})

const showFileError = computed(() => {
  return submitted.value && !selectedAudioFile.value
})

const durationStatus = computed(() => {
  if (detectingDuration.value) return 'Detecting duration...'

  if (form.duration_seconds !== null && form.duration_seconds !== undefined) {
    return 'Auto detected'
  }

  return ''
})

const formattedDuration = computed(() => {
  if (form.duration_seconds === null || form.duration_seconds === undefined) return ''

  const total = Number(form.duration_seconds)
  if (!Number.isFinite(total)) return ''

  const minutes = Math.floor(total / 60)
  const seconds = total % 60

  return `${minutes}:${String(seconds).padStart(2, '0')}`
})

const artistOptions = computed(() => {
  return artists.value.map((artist) => ({
    label: artist.name,
    value: artist.id,
  }))
})

const albumOptions = computed(() => {
  return albums.value.map((album) => ({
    label: album.title,
    value: album.id,
  }))
})

const genreOptions = computed(() => {
  return genres.value.map((genre) => ({
    label: genre.name,
    value: genre.id,
  }))
})

onMounted(async () => {
  await Promise.all([fetchArtists(), fetchAlbums(), fetchGenres()])
})

onBeforeUnmount(() => {
  revokeAudioObjectUrl()
  revokeCoverObjectUrl()
})

async function onAudioFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0] || null

  selectedAudioFile.value = file
  audioFileName.value = file?.name || ''
  audioFileSizeLabel.value = file ? formatFileSize(file.size) : ''

  resultMessage.value = ''
  errorMessage.value = ''
  createdTrack.value = null
  audioUploadProgress.value = 0
  coverUploadProgress.value = 0

  if (!file) {
    form.duration_seconds = null
    return
  }

  if (!file.type.startsWith('audio/')) {
    errorMessage.value = 'Please select a valid audio file.'

    selectedAudioFile.value = null
    audioFileName.value = ''
    audioFileSizeLabel.value = ''
    form.duration_seconds = null
    target.value = ''

    return
  }

  if (!String(form.title || '').trim()) {
    form.title = file.name.replace(/\.[^/.]+$/, '')
  }

  await detectAudioDuration(file)
}

function onCoverFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0] || null

  resultMessage.value = ''
  errorMessage.value = ''
  createdTrack.value = null
  coverUploadProgress.value = 0

  if (!file) {
    removeCover()
    return
  }

  if (!file.type.startsWith('image/')) {
    errorMessage.value = 'Please select a valid image file.'

    selectedCoverFile.value = null
    coverFileName.value = ''
    coverFileSizeLabel.value = ''
    revokeCoverObjectUrl()
    target.value = ''

    return
  }

  selectedCoverFile.value = file
  coverFileName.value = file.name
  coverFileSizeLabel.value = formatFileSize(file.size)

  revokeCoverObjectUrl()
  coverObjectUrl = URL.createObjectURL(file)
  coverPreviewUrl.value = coverObjectUrl
}

function removeCover() {
  selectedCoverFile.value = null
  coverFileName.value = ''
  coverFileSizeLabel.value = ''
  coverUploadProgress.value = 0
  form.cover_url = null

  revokeCoverObjectUrl()

  if (coverInputRef.value) {
    coverInputRef.value.value = ''
  }
}

function detectAudioDuration(file: File) {
  return new Promise<void>((resolve) => {
    detectingDuration.value = true
    form.duration_seconds = null

    revokeAudioObjectUrl()

    audioObjectUrl = URL.createObjectURL(file)

    const audio = document.createElement('audio')
    audio.preload = 'metadata'
    audio.src = audioObjectUrl

    audio.onloadedmetadata = () => {
      if (Number.isFinite(audio.duration)) {
        form.duration_seconds = Math.max(0, Math.round(audio.duration))
      }

      detectingDuration.value = false
      resolve()
    }

    audio.onerror = () => {
      detectingDuration.value = false
      errorMessage.value =
        'Could not detect audio duration automatically. You can enter it manually.'

      resolve()
    }
  })
}

function revokeAudioObjectUrl() {
  if (audioObjectUrl) {
    URL.revokeObjectURL(audioObjectUrl)
    audioObjectUrl = null
  }
}

function revokeCoverObjectUrl() {
  if (coverObjectUrl) {
    URL.revokeObjectURL(coverObjectUrl)
    coverObjectUrl = null
  }

  coverPreviewUrl.value = ''
}

function formatFileSize(bytes: number) {
  if (bytes < 1024) return `(${bytes} B)`
  if (bytes < 1024 * 1024) return `(${(bytes / 1024).toFixed(1)} KB)`

  return `(${(bytes / 1024 / 1024).toFixed(1)} MB)`
}

function resetForm() {
  form.title = ''
  form.artist_id = null
  form.album_id = null
  form.genre_id = null
  form.duration_seconds = null
  form.audio_url = null
  form.cover_url = null
  submitted.value = false
}

function clear() {
  selectedAudioFile.value = null
  audioFileName.value = ''
  audioFileSizeLabel.value = ''

  removeCover()

  resultMessage.value = ''
  errorMessage.value = ''
  createdTrack.value = null

  audioUploadProgress.value = 0
  coverUploadProgress.value = 0
  uploadStep.value = 'idle'

  revokeAudioObjectUrl()
  resetForm()

  if (audioInputRef.value) {
    audioInputRef.value.value = ''
  }
}

function getUploadedDurationSeconds(uploadedMedia: {
  durationSeconds?: number | null
  duration_seconds?: number | null
}) {
  return uploadedMedia.durationSeconds ?? uploadedMedia.duration_seconds ?? null
}

function getUploadedMediaUrl(uploadedMedia: {
  url?: string
  file_url?: string
  fileUrl?: string
  path?: string
}) {
  return (
    uploadedMedia.url ||
    uploadedMedia.fileUrl ||
    uploadedMedia.file_url ||
    uploadedMedia.path ||
    null
  )
}

function getApiErrorMessage(error: any) {
  return (
    error?.response?.data?.details?.artistId ||
    error?.response?.data?.details?.artistID ||
    error?.response?.data?.details?.artist_id ||
    error?.response?.data?.message ||
    error?.message ||
    'Upload or track creation failed. Please try again.'
  )
}

async function uploadAudioFile(file: File) {
  uploadStep.value = 'audio'
  audioUploadProgress.value = 0

  const uploadedMedia = await mediaApi.uploadAdminMedia(
    AUDIO_MEDIA_TYPE,
    file,
    undefined,
    {
      onUploadProgress(event) {
        if (!event.total) return
        audioUploadProgress.value = Math.round((event.loaded / event.total) * 100)
      },
    },
  )

  audioUploadProgress.value = 100

  return uploadedMedia
}

async function uploadCoverFile(file: File) {
  uploadStep.value = 'cover'
  coverUploadProgress.value = 0

  const uploadedMedia = await mediaApi.uploadAdminMedia(
    COVER_MEDIA_TYPE,
    file,
    undefined,
    {
      onUploadProgress(event) {
        if (!event.total) return
        coverUploadProgress.value = Math.round((event.loaded / event.total) * 100)
      },
    },
  )

  coverUploadProgress.value = 100

  return uploadedMedia
}

async function submit() {
  submitted.value = true
  errorMessage.value = ''
  resultMessage.value = ''
  createdTrack.value = null

  if (!selectedAudioFile.value) {
    errorMessage.value = 'Please choose an audio file.'
    return
  }

  if (!form.artist_id) {
    errorMessage.value = 'Please select an artist.'
    return
  }

  if (!String(form.title || '').trim()) {
    errorMessage.value = 'Please enter a track title.'
    return
  }

  if (detectingDuration.value) {
    errorMessage.value = 'Please wait until duration detection finishes.'
    return
  }

  uploading.value = true
  uploadStep.value = 'audio'
  audioUploadProgress.value = 0
  coverUploadProgress.value = 0

  try {
    const uploadedAudio = await uploadAudioFile(selectedAudioFile.value)

    const uploadedAudioUrl = getUploadedMediaUrl(uploadedAudio)

    if (!uploadedAudioUrl) {
      throw new Error('Audio upload succeeded but no media URL was returned.')
    }

    let uploadedCoverUrl: string | null = null

    if (selectedCoverFile.value) {
      const uploadedCover = await uploadCoverFile(selectedCoverFile.value)
      uploadedCoverUrl = getUploadedMediaUrl(uploadedCover)

      if (!uploadedCoverUrl) {
        throw new Error('Cover upload succeeded but no media URL was returned.')
      }
    }

    const uploadedDurationSeconds = getUploadedDurationSeconds(uploadedAudio)

    uploadStep.value = 'track'

    const payload: TrackFormPayload = {
      title: String(form.title).trim(),
      artist_id: form.artist_id ?? null,
      album_id: form.album_id ?? null,
      genre_id: form.genre_id ?? null,
      duration_seconds: uploadedDurationSeconds ?? form.duration_seconds ?? null,
      audio_url: uploadedAudioUrl,
      cover_url: uploadedCoverUrl ?? form.cover_url ?? null,
    }

    const track = await createTrack(payload)

    createdTrack.value = track

    const audioWasDuplicate = Boolean(uploadedAudio?.duplicate)

    resultMessage.value = audioWasDuplicate
      ? 'This audio file was already uploaded. Existing media URL reused and track was created.'
      : selectedCoverFile.value
        ? 'Audio and cover uploaded successfully. Track was created.'
        : 'Audio uploaded successfully. Track was created.'

    toast.add({
      severity: 'success',
      summary: 'Track created',
      detail: track.title,
      life: 2500,
    })

    clear()

    createdTrack.value = track
    resultMessage.value = audioWasDuplicate
      ? 'This audio file was already uploaded. Existing media URL reused and track was created.'
      : selectedCoverFile.value
        ? 'Audio and cover uploaded successfully. Track was created.'
        : 'Audio uploaded successfully. Track was created.'
  } catch (error: any) {
    errorMessage.value = getApiErrorMessage(error)

    toast.add({
      severity: 'error',
      summary: 'Failed',
      detail: errorMessage.value,
      life: 3000,
    })
  } finally {
    uploading.value = false
    uploadStep.value = 'idle'

    if (audioUploadProgress.value > 0) {
      audioUploadProgress.value = 100
    }

    if (selectedCoverFile.value && coverUploadProgress.value > 0) {
      coverUploadProgress.value = 100
    }
  }
}
</script>
