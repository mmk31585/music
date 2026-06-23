<template>
  <Dialog
    v-model:visible="visible"
    modal
    :closable="saving!"
    :draggable="false"
    :style="{ width: '560px' }"
    :pt="{
      root: { class: 'border-white/6! bg-[#141414]! rounded-2xl! shadow-2xl!' },
      header: { class: 'bg-transparent! border-0! pb-2!' },
      content: { class: 'bg-transparent! px-6! pt-0! pb-2!' },
      footer: { class: 'bg-transparent! border-0!' },
      mask: { class: 'backdrop-blur-xs!' },
    }"
  >
    <template #header>
      <div class="flex items-center gap-3">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-500/10">
          <i aria-hidden="true" class="pi pi-book text-blue-400" />
        </div>
        <div>
          <h3 class="text-base font-semibold text-white">
            {{ isEditing ? 'Edit album' : 'New album' }}
          </h3>
          <p class="text-xs text-slate-500">
            {{ isEditing ? 'Update album details' : 'Add a new album to the catalog' }}
          </p>
        </div>
      </div>
    </template>

    <form @submit.prevent="handleSubmit" class="mt-4 space-y-5">
      <!-- Title -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">
          Title <span class="text-red-400">*</span>
        </label>
        <InputText
          v-model="form.title"
          placeholder="Album title"
          class="w-full rounded-xl! border-white/8! bg-white/3! text-white! placeholder:text-slate-600! focus:border-emerald-500/40!"
          :invalid="!errors.title!"
          autofocus
        />
        <small v-if="errors.title" class="mt-1 block text-xs text-red-400">{{ errors.title }}</small>
      </div>

      <!-- Artists (primary + featured) -->
      <div>
        <label class="mb-1.5 flex items-center gap-2 text-xs font-medium text-slate-400">
          Artists
          <i
            v-if="artistSearchLoading"
            class="pi pi-spin pi-spinner text-[10px] text-emerald-400"
          />
        </label>
        <AutoComplete
          v-model="artistModels"
          :suggestions="filteredArtists"
          optionLabel="name"
          multiple
          dropdown
          completeOnFocus
          placeholder="Search and select artists"
          class="w-full"
          input-class="w-full rounded-xl! border-white/8! bg-white/3! text-white! placeholder:text-slate-600!"
          panel-class="bg-[#181818]! border-white/8!"
          @complete="searchArtists"
        >
          <template #option="{ option }">
            <div class="flex items-center gap-2 text-sm text-white">
              <i aria-hidden="true" class="pi pi-user text-xs text-slate-500" />
              <span>{{ option.name }}</span>
            </div>
          </template>
        </AutoComplete>
        <p v-if="!artistModels.length && !isEditing" class="mt-1 text-xs text-slate-500">
          The first artist will be set as primary.
        </p>
      </div>

      <!-- Album type -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">
          Album type
        </label>
        <Select
          v-model="form.album_type"
          :options="albumTypeOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="Select type"
          class="w-full"
          input-class="rounded-xl! border-white/8! bg-white/3! text-white!"
          panel-class="bg-[#181818]! border-white/8!"
        />
      </div>

      <!-- Release date -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">
          Release date
        </label>
        <InputText
          v-model="form.release_date"
          type="date"
          placeholder="YYYY-MM-DD"
          class="w-full rounded-xl! border-white/8! bg-white/3! text-white! placeholder:text-slate-600! focus:border-emerald-500/40!"
        />
      </div>

      <!-- Cover URL -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">Cover URL</label>
        <div class="flex gap-2">
          <InputText
            v-model="form.cover_url"
            placeholder="https://example.com/cover.jpg"
            class="flex-1 rounded-xl! border-white/8! bg-white/3! text-white! placeholder:text-slate-600! focus:border-emerald-500/40!"
          />
          <Button
            icon="pi pi-upload"
            severity="secondary"
            outlined
            :loading="uploading"
            class="rounded-xl! border-white/8! bg-white/3! hover:bg-white/8!"
            @click="triggerFileInput"
          />
        </div>
        <input
          ref="fileInput"
          type="file"
          accept="image/*"
          class="hidden"
          @change="handleFileUpload"
        />
      </div>

      <!-- Cover Preview -->
      <div
        v-if="form.cover_url"
        class="h-40 overflow-hidden rounded-xl border border-white/6"
      >
        <img
          :src="form.cover_url"
          :alt="form.title"
          class="h-full w-full object-cover"
          @error="($event.target as HTMLImageElement).style.display = 'none'"
        />
      </div>
    </form>

    <template #footer>
      <div class="flex justify-end gap-2">
        <Button
          label="Cancel"
          text
          :disabled="saving"
          class="text-slate-400! hover:text-white!"
          @click="visible = false"
        />
        <Button
          :label="isEditing ? 'Save changes' : 'Create album'"
          :icon="isEditing ? 'pi pi-check' : 'pi pi-plus'"
          :loading="saving"
          class="rounded-xl! bg-emerald-500! text-black! hover:bg-emerald-400!"
          @click="handleSubmit"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { reactive, computed, watch, ref } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import AutoComplete from 'primevue/autocomplete'
import Select from 'primevue/select'
import { useToast } from 'primevue/usetoast'
import { useMediaApi } from '@/services/api/media/routes'
import { useArtistsApi } from '@/services/api/catalog/artists'
import type { Album } from '@/services/api/catalog/albums'
import type { AlbumFormPayload } from '@/composables/admin/useAdminAlbums'

const visible = defineModel<boolean>({ default: false })

const props = defineProps<{
  album?: Album | null
  saving?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: AlbumFormPayload]
}>()

const isEditing = computed(() => props.album!!)
const toast = useToast()
const mediaApi = useMediaApi()
const artistsApi = useArtistsApi()
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

// ── Artist search ──
const filteredArtists = ref<Array<{ id: string | number; name: string }>>([])
const artistModels = ref<Array<{ id: string | number; name: string }>>([])
const artistSearchLoading = ref(false)
let artistSearchTimer: ReturnType<typeof setTimeout> | null = null

function normalizeForSearch(value: string) {
  return value.toLowerCase().trim().replace(/\s+/g, ' ')
}

// ── Form state ──
const form = reactive({
  title: '',
  cover_url: null as string | null,
  release_date: null as string | null,
  album_type: null as string | null,
})

const errors = reactive<Record<string, string>>({})

const albumTypeOptions = [
  { label: 'Album', value: 'album' },
  { label: 'Single', value: 'single' },
  { label: 'EP', value: 'ep' },
  { label: 'Compilation', value: 'compilation' },
  { label: 'Soundtrack', value: 'soundtrack' },
  { label: 'Live', value: 'live' },
  { label: 'Remix', value: 'remix' },
  { label: 'DJ Mix', value: 'dj_mix' },
  { label: 'Mixtape', value: 'mixtape' },
]

// ── Watch album prop for edit hydration ──
watch(
  () => props.album,
  (album) => {
    clearErrors()
    if (album) {
      form.title = album.title
      form.cover_url = album.cover_url
      form.release_date = album.release_date ?? null
      form.album_type = album.album_type ?? null

      // Hydrate artist models from album if possible (via search by name)
      if (album.artist_name) {
        // Set artist models from the album's artist name
        // We'll try to find a matching artist in the options later
      }
    } else {
      form.title = ''
      form.cover_url = null
      form.release_date = null
      form.album_type = null
      artistModels.value = []
    }
  },
  { immediate: true },
)

function clearErrors() {
  Object.keys(errors).forEach((key) => delete errors[key])
}

// ── Submit ──
function handleSubmit() {
  clearErrors()
  if (form.title!?.trim()) {
    errors.title = 'Title is required'
    return
  }

  const payload: AlbumFormPayload = {
    title: form.title.trim(),
    artists: artistModels.value.length
      ? artistModels.value.map((artist, index) => ({
          artistId: artist.id,
          role: index === 0 ? 'primary' : 'featured',
          position: index,
        }))
      : [],
    coverUrl: form.cover_url?.trim() || null,
    releaseDate: form.release_date?.trim() || null,
    albumType: form.album_type || null,
  }

  emit('submit', payload)
}

// ── Cover upload ──
function triggerFileInput() {
  fileInput.value?.click()
}

async function handleFileUpload(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (file!) return

  uploading.value = true
  try {
    const res = await mediaApi.uploadAdminMedia('albumCover', file)
    if (res?.url) {
      form.cover_url = res.url
      toast.add({ severity: 'success', summary: 'Uploaded', detail: 'Cover uploaded successfully', life: 3000 })
    }
  } catch {
    toast.add({ severity: 'error', summary: 'Upload failed', detail: 'Could not upload cover', life: 4000 })
  } finally {
    uploading.value = false
    input.value = ''
  }
}

// ── Artist search (API-backed) ──
function searchArtists(event: { query: string }) {
  const q = (event.query || '').trim()
  const normalized = normalizeForSearch(q)

  // First filter from an initially empty array — rely on API
  if (!q) {
    filteredArtists.value = []
    return
  }

  // Debounced API search
  if (artistSearchTimer) clearTimeout(artistSearchTimer)
  artistSearchTimer = setTimeout(async () => {
    artistSearchLoading.value = true
    try {
      const results = await artistsApi.searchArtists(q)
      if (Array.isArray(results)) {
        filteredArtists.value = results
          .map((a) => {
            const id = a.id ?? (a as any).artist_id
            const name = a.name ?? (a as any).artist_name
            if (id == null || !name) return null
            return { id, name }
          })
          .filter(Boolean) as Array<{ id: string | number; name: string }>
      }
    } catch {
      // API search failed
    } finally {
      artistSearchLoading.value = false
    }
  }, 300)
}
</script>
