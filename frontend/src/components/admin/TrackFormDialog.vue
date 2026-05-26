<template>
  <Dialog
    v-model:visible="visibleInternal"
    modal
    :header="isEdit ? 'Edit track' : 'Create track'"
    :style="{ width: '520px' }"
    class="track-form-dialog"
  >
    <div class="space-y-4 py-2">
      <div class="grid gap-4">
        <span class="p-float-label">
          <InputText id="title" v-model="form.title" class="w-full" />
          <label for="title">Title</label>
        </span>

        <span class="p-float-label">
          <Dropdown
            id="artist"
            v-model="form.artist_id"
            :options="artists"
            option-label="name"
            option-value="id"
            class="w-full"
            :loading="loadingLookups"
            show-clear
          />
          <label for="artist">Artist</label>
        </span>

        <span class="p-float-label">
          <Dropdown
            id="album"
            v-model="form.album_id"
            :options="albums"
            option-label="title"
            option-value="id"
            class="w-full"
            :loading="loadingLookups"
            show-clear
          />
          <label for="album">Album</label>
        </span>

        <span class="p-float-label">
          <Dropdown
            id="genre"
            v-model="form.genre_id"
            :options="genres"
            option-label="name"
            option-value="id"
            class="w-full"
            :loading="loadingLookups"
            show-clear
          />
          <label for="genre">Genre</label>
        </span>

        <span class="p-float-label">
          <InputNumber
            id="duration"
            v-model="form.duration_seconds"
            class="w-full"
            :min="0"
            :use-grouping="false"
          />
          <label for="duration">Duration (seconds)</label>
        </span>

        <div class="space-y-2">
          <label class="text-xs font-medium text-slate-300">Audio URL</label>
          <div class="flex flex-col gap-2 sm:flex-row">
            <InputText v-model="form.audio_url" class="flex-1" />
            <input
              v-if="!isEdit"
              type="file"
              accept="audio/*"
              class="block w-full rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-sm text-slate-300 file:mr-3 file:rounded-md file:border-0 file:bg-[#1db954] file:px-3 file:py-1.5 file:text-sm file:font-semibold file:text-black sm:max-w-56"
              @change="onAudioFileChange"
            />
          </div>
          <p v-if="selectedAudioFileName" class="truncate text-xs text-emerald-300">
            Selected audio: {{ selectedAudioFileName }}
          </p>
        </div>

        <div class="space-y-2">
          <label class="text-xs font-medium text-slate-300">Cover URL</label>
          <div class="flex gap-2">
            <InputText v-model="form.cover_url" class="flex-1" />
            <Button
              icon="pi pi-image"
              text
              rounded
              @click="showCoverUploadDialog = true"
              v-tooltip="'Upload cover'"
            />
          </div>
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-4">
        <Button label="Cancel" severity="secondary" text :disabled="saving" @click="close" />
        <Button
          :label="isEdit ? 'Save changes' : 'Create track'"
          :loading="saving"
          class="border-0 bg-[#1db954] text-black"
          @click="handleSubmit"
        />
      </div>
    </div>

    <UploadMediaDialog
      v-model="showCoverUploadDialog"
      kind="track-cover"
      @uploaded="(url) => (form.cover_url = url)"
    />
  </Dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import InputNumber from 'primevue/inputnumber'
import { useToast } from 'primevue/usetoast'
import { useCatalogApi } from '@/services/api/catalog'
import type { Track, Artist, Album, Genre, TrackFormPayload } from '@/services/api/catalog'
import UploadMediaDialog from './UploadMediaDialog.vue'

const props = defineProps<{
  modelValue: boolean
  track?: Track | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'submit', payload: TrackFormPayload): void
}>()
const showCoverUploadDialog = ref(false)

const toast = useToast()
const api = useCatalogApi()

const visibleInternal = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const isEdit = computed(() => !!props.track?.id)

const form = reactive<Partial<Track>>({
  id: undefined,
  title: '',
  duration_seconds: undefined,
  audio_url: '',
  cover_url: '',
  artist_id: undefined,
  album_id: undefined,
  genre_id: undefined,
})

const artists = ref<Artist[]>([])
const albums = ref<Album[]>([])
const genres = ref<Genre[]>([])
const loadingLookups = ref(false)
const saving = ref(false)
const selectedAudioFile = ref<File | null>(null)
const selectedAudioFileName = ref('')

watch(
  () => props.track,
  (track) => {
    if (track) {
      selectedAudioFile.value = null
      selectedAudioFileName.value = ''
      Object.assign(form, {
        id: track.id,
        title: track.title,
        duration_seconds: track.duration_seconds ?? undefined,
        audio_url: track.audio_url ?? '',
        cover_url: track.cover_url ?? '',
        artist_id: track.artist_id ?? undefined,
        album_id: track.album_id ?? undefined,
        genre_id: track.genre_id ?? undefined,
      })
    } else {
      resetForm()
    }
  },
  { immediate: true },
)

watch(
  () => props.modelValue,
  async (open) => {
    if (open) {
      await loadLookups()
    }
  },
)

function resetForm() {
  Object.assign(form, {
    id: undefined,
    title: '',
    duration_seconds: undefined,
    audio_url: '',
    cover_url: '',
    artist_id: undefined,
    album_id: undefined,
    genre_id: undefined,
  })
  selectedAudioFile.value = null
  selectedAudioFileName.value = ''
}

async function loadLookups() {
  if (loadingLookups.value) return
  loadingLookups.value = true
  try {
    const [artistList, albumList, genreList] = await Promise.all([
      api.getArtists(),
      api.getAlbums(),
      api.getGenres(),
    ])
    artists.value = artistList
    albums.value = albumList
    genres.value = genreList
  } finally {
    loadingLookups.value = false
  }
}

async function handleSubmit() {
  if (!form.title?.trim()) {
    toast.add({
      severity: 'warn',
      summary: 'Missing title',
      detail: 'Please enter a track title.',
      life: 2000,
    })
    return
  }

  if (!form.artist_id) {
    toast.add({
      severity: 'warn',
      summary: 'Missing artist',
      detail: 'Please choose an artist before creating the track.',
      life: 2500,
    })
    return
  }

  if (!isEdit.value && !selectedAudioFile.value && !form.audio_url) {
    toast.add({
      severity: 'warn',
      summary: 'Missing audio',
      detail: 'Please choose an audio file or paste an audio URL.',
      life: 2500,
    })
    return
  }

  saving.value = true
  try {
    const payload: TrackFormPayload = {
      id: form.id,
      title: form.title?.trim(),
      duration_seconds: form.duration_seconds,
      audio_url: form.audio_url || null,
      cover_url: form.cover_url || null,
      artist_id: form.artist_id || null,
      album_id: form.album_id || null,
      genre_id: form.genre_id || null,
      audioFile: selectedAudioFile.value,
    }
    emit('submit', payload)
    visibleInternal.value = false
    resetForm()
  } finally {
    saving.value = false
  }
}

function close() {
  visibleInternal.value = false
}

function onAudioFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0] ?? null

  selectedAudioFile.value = file
  selectedAudioFileName.value = file?.name ?? ''
}
</script>

<style scoped>
.track-form-dialog :deep(.p-dialog-header) {
  border-bottom: 1px solid rgb(255 255 255 / 0.1);
  background: #000;
  color: #fff;
}
.track-form-dialog :deep(.p-dialog-content) {
  background: #000;
  color: #fff;
}
</style>
