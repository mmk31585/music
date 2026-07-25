<template>
  <Dialog
    v-model:visible="visible"
    modal
    :closable="saving!"
    :draggable="false"
    :style="{ width: '560px' }"
    :pt="{
      root: { class: 'border-white/6! bg-surface-raised! rounded-2xl! shadow-2xl!' },
      header: { class: 'bg-transparent! border-0! pb-2!' },
      content: { class: 'bg-transparent! px-6! pt-0! pb-2!' },
      footer: { class: 'bg-transparent! border-0!' },
      mask: { class: 'backdrop-blur-xs!' },
    }"
  >
    <template #header>
      <div class="flex items-center gap-3">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-500/10">
          <Book aria-hidden="true" class="text-blue-400"  />
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
          <Loader2 aria-hidden="true" v-if="artistSearchLoading"
            class="text-[10px] text-emerald-400 animate-spin" />
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
          panel-class="bg-surface-overlay! border-white/8!"
          @complete="searchArtists"
        >
          <template #option="{ option }">
            <div class="flex items-center gap-2 text-sm text-white">
              <User aria-hidden="true" class="text-xs text-slate-500"  />
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
          panel-class="bg-surface-overlay! border-white/8!"
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
          v-if="isEditing"
          icon="pi pi-refresh"
          label="Enrich data"
          severity="info"
          text
          :loading="enriching"
          class="mr-auto text-slate-500! hover:text-emerald-400!"
          @click="handleEnrich"
        />
        <Button
          label="Cancel"
          text
          :disabled="saving || enriching"
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
import { Book, Loader2, User } from 'lucide-vue-next'
import { reactive, computed, watch, ref } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import AutoComplete from 'primevue/autocomplete'
import Select from 'primevue/select'
import { useToast } from 'primevue/usetoast'
import { useMediaApi } from '@/services/api/media/routes'
import { useArtistsApi } from '@/services/api/catalog/artists'
import { useAlbumsApi } from '@/services/api/catalog/albums'
import type { Album } from '@/services/api/catalog/albums'
import type { AlbumFormPayload } from '@/composables/admin/useAdminAlbums'

const visible = defineModel<boolean>({ default: false })

const props = defineProps<{
  album?: Album | null
  saving?: boolean
  prefillTitle?: string
}>()

const emit = defineEmits<{
  submit: [payload: AlbumFormPayload]
}>()

const isEditing = computed(() => props.album!)
const toast = useToast()
const mediaApi = useMediaApi()
const artistsApi = useArtistsApi()
const albumsApi = useAlbumsApi()
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)
const enriching = ref(false)

// ── Artist search ──
const filteredArtists = ref<Array<{ id: string | number; name: string }>>([])
const artistModels = ref<Array<{ id: string | number; name: string }>>([])
const artistSearchLoading = ref(false)
let artistSearchTimer: ReturnType<typeof setTimeout> | null = null
let lastSearchTerm = ''

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

// ── Refetch album when dialog opens to always show latest persisted state ──
watch(visible, async (val) => {
  if (val && props.album?.id) {
    try {
      const refreshed = await albumsApi.getAlbum(props.album.id)
      if (refreshed) {
        applyAlbum(refreshed)
      }
    } catch {
      // Silently fall back to the prop data
    }
  }
  if (val) {
    clearErrors()
  }
})

function applyAlbum(album: Album) {
  form.title = album.title
  form.cover_url = album.cover_url
  form.release_date = album.release_date ?? null
  form.album_type = album.album_type ?? null

  // Hydrate artist models from album artist name
  if (album.artist_name) {
    searchArtistByName(album.artist_name)
  } else {
    artistModels.value = []
  }
}

// ── Watch album prop & prefillTitle for edit / create hydration ──
watch(
  [() => props.album, () => props.prefillTitle],
  ([album, prefillTitle]) => {
    clearErrors()
    if (album) {
      applyAlbum(album)
    } else {
      form.title = prefillTitle ?? ''
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

// ── Search artist by name and populate artistModels ──
async function searchArtistByName(name: string) {
  const trimmed = name.trim()
  if (!trimmed) {
    artistModels.value = []
    return
  }
  try {
    const results = await artistsApi.searchArtists(trimmed)
    const list = Array.isArray(results) ? results : []
    const match = list.find((a: any) => {
      const an = (a.name ?? a.artist_name ?? '').trim().toLowerCase()
      return an === trimmed.toLowerCase()
    })
    if (match) {
      artistModels.value = [{
        id: match.id ?? (match as any).artist_id,
        name: match.name ?? (match as any).artist_name,
      }]
    } else if (list.length > 0) {
      // Fall back to first result
      artistModels.value = [{
        id: list[0]?.id ?? (list[0] as any)?.artist_id,
        name: list[0]?.name ?? (list[0] as any)?.artist_name,
      }]
    }
  } catch {
    // Silently skip
  }
}

// ── Submit ──
function handleSubmit() {
  clearErrors()
  if (!form.title?.trim()) {
    errors.title = 'Title is required'
    return
  }

  const payload: AlbumFormPayload = {
    title: form.title.trim(),
    artists: artistModels.value.length
      ? artistModels.value.map((artist, index) => ({
          artist_id: artist.id,
          role: index === 0 ? 'primary' : 'featured',
          position: index,
        }))
      : [],
    cover_url: form.cover_url?.trim() || null,
    release_date: form.release_date?.trim() || null,
    album_type: form.album_type || null,
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
  if (!file) return

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

  if (!q) {
    filteredArtists.value = []
    return
  }

  // Debounced API search with stale-response guard
  if (artistSearchTimer) clearTimeout(artistSearchTimer)
  lastSearchTerm = q
  artistSearchTimer = setTimeout(async () => {
    artistSearchLoading.value = true
    try {
      const results = await artistsApi.searchArtists(q)
      // Discard stale responses — only apply if this is still the latest search
      if (lastSearchTerm !== q) return
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
    } catch (err) {
      console.error('[AlbumFormDialog] artist search failed', err)
    } finally {
      if (lastSearchTerm === q) artistSearchLoading.value = false
    }
  }, 300)
}

// ── Enrich from external source ──
async function handleEnrich() {
  if (!props.album?.id) return
  enriching.value = true
  try {
    await albumsApi.adminEnrichAlbum(props.album.id)
    // Refetch the album from the API to sync the form with persisted state
    const refreshed = await albumsApi.getAlbum(props.album.id)
    if (refreshed) {
      form.title = refreshed.title
      form.cover_url = refreshed.cover_url
      form.release_date = refreshed.release_date ?? null
      form.album_type = refreshed.album_type ?? null
      if (refreshed.artist_name) {
        searchArtistByName(refreshed.artist_name)
      }
    }
    toast.add({ severity: 'success', summary: 'Enriched', detail: 'Album data fetched from external sources', life: 3000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Enrich failed', detail: 'Could not fetch album data. Check external service connectivity.', life: 4000 })
  } finally {
    enriching.value = false
  }
}
</script>
