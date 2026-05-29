<template>
  <div class="rounded-3xl border border-white/10 bg-white/5 p-6 shadow-lg">
    <div class="mb-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-xl font-bold text-white">Albums</h2>
        <p class="mt-2 text-sm text-slate-400">Manage album titles, covers, and artists.</p>
      </div>

      <Button
        label="Add album"
        icon="pi pi-plus"
        class="border-0 bg-[#1db954] text-black"
        @click="openCreate"
      />
    </div>

    <DataTable
      :value="albums"
      :loading="loading"
      data-key="id"
      responsive-layout="scroll"
      class="overflow-hidden rounded-2xl"
    >
      <Column header="Cover">
        <template #body="{ data }">
          <img
            v-if="data.cover_url"
            :src="data.cover_url"
            :alt="data.title"
            class="h-12 w-12 rounded-xl object-cover"
          />

          <div
            v-else
            class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/10 text-xs text-slate-400"
          >
            N/A
          </div>
        </template>
      </Column>

      <Column field="title" header="Title" sortable />

      <Column header="Artist">
        <template #body="{ data }">
          {{ getArtistName(data.artist_id) }}
        </template>
      </Column>

      <Column field="cover_url" header="Cover URL">
        <template #body="{ data }">
          <span class="line-clamp-1 max-w-md text-sm text-slate-300">
            {{ data.cover_url || '-' }}
          </span>
        </template>
      </Column>

      <Column header="Actions">
        <template #body="{ data }">
          <div class="flex gap-2">
            <Button icon="pi pi-pencil" severity="secondary" text rounded @click="openEdit(data)" />

            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              rounded
              :loading="deleting"
              @click="handleDelete(data.id)"
            />
          </div>
        </template>
      </Column>

      <template #empty>
        <div class="py-8 text-center text-sm text-slate-400">No albums found.</div>
      </template>
    </DataTable>

    <Dialog
      v-model:visible="dialogVisible"
      modal
      :header="editingAlbum ? 'Edit album' : 'Add album'"
      class="w-full max-w-lg"
    >
      <div class="flex flex-col gap-4">
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300">Title</label>

          <InputText
            v-model="form.title"
            class="w-full"
            placeholder="Album title"
            autofocus
            :disabled="saving || uploadingImage"
          />
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300">Artist</label>

          <Dropdown
            v-model="form.artist_id"
            :options="artistOptions"
            option-label="label"
            option-value="value"
            placeholder="Select artist"
            class="w-full"
            :loading="artistsLoading"
            :disabled="saving || uploadingImage"
            show-clear
          />
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300">Album cover</label>

          <input
            type="file"
            accept="image/*"
            class="block w-full rounded-xl border border-white/10 bg-black/20 px-3 py-2 text-sm text-slate-300 file:mr-4 file:rounded-lg file:border-0 file:bg-[#1db954] file:px-4 file:py-2 file:text-sm file:font-semibold file:text-black"
            :disabled="saving || uploadingImage"
            @change="handleCoverImageChange"
          />

          <p v-if="uploadingImage" class="mt-2 text-xs text-slate-400">Uploading cover...</p>

          <p v-if="uploadImageError" class="mt-2 text-xs text-red-300">
            {{ uploadImageError }}
          </p>

          <p v-if="form.cover_url" class="mt-2 line-clamp-1 text-xs text-slate-500">
            {{ form.cover_url }}
          </p>
        </div>

        <div v-if="form.cover_url" class="rounded-2xl bg-black/20 p-3">
          <img
            :src="form.cover_url"
            alt="Album preview"
            class="h-32 w-32 rounded-2xl object-cover"
          />

          <Button
            type="button"
            label="Remove cover"
            severity="danger"
            text
            size="small"
            class="mt-3 px-0"
            :disabled="saving || uploadingImage"
            @click="removeCover"
          />
        </div>

        <div
          v-if="errorMessage"
          class="rounded-2xl border border-red-500/20 bg-red-500/10 px-4 py-3 text-sm text-red-300"
        >
          {{ errorMessage }}
        </div>
      </div>

      <template #footer>
        <Button
          label="Cancel"
          severity="secondary"
          outlined
          :disabled="saving || uploadingImage"
          @click="dialogVisible = false"
        />

        <Button
          :label="editingAlbum ? 'Save changes' : 'Create album'"
          icon="pi pi-check"
          :loading="saving"
          :disabled="!canSubmit"
          class="border-0 bg-[#1db954] text-black"
          @click="submit"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import { useToast } from 'primevue/usetoast'

import { useAdminAlbums, useAdminArtists } from '@/composables/admin'
import { useCatalogImageUpload } from '@/composables/media/useCatalogImageUpload'
import type { Album } from '@/services/api/catalog'

const toast = useToast()

const { albums, loading, saving, deleting, fetchAlbums, createAlbum, updateAlbum, deleteAlbum } =
  useAdminAlbums()

const { artists, loading: artistsLoading, fetchArtists } = useAdminArtists()

const { uploadingImage, uploadImageError, uploadCatalogImage } = useCatalogImageUpload()

const dialogVisible = ref(false)
const editingAlbum = ref<Album | null>(null)
const errorMessage = ref('')

const form = reactive({
  title: '',
  artist_id: null as string | number | null,
  cover_url: '',
})

const canSubmit = computed(() => {
  return form.title.trim().length > 0 && !saving.value && !uploadingImage.value
})

const artistOptions = computed(() => {
  return artists.value.map((artist) => ({
    label: artist.name,
    value: artist.id,
  }))
})

onMounted(async () => {
  await Promise.all([fetchAlbums(), fetchArtists()])
})

function getArtistName(id: string | number | null | undefined) {
  if (id == null) return '-'

  const artist = artists.value.find((item) => String(item.id) === String(id))
  return artist?.name || '-'
}

function resetForm() {
  form.title = ''
  form.artist_id = null
  form.cover_url = ''
  errorMessage.value = ''
  editingAlbum.value = null
}

function openCreate() {
  resetForm()
  dialogVisible.value = true
}

function openEdit(album: Album) {
  editingAlbum.value = album
  form.title = album.title
  form.artist_id = album.artist_id ?? null
  form.cover_url = album.cover_url || ''
  errorMessage.value = ''
  dialogVisible.value = true
}

async function handleCoverImageChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) return

  errorMessage.value = ''

  try {
    form.cover_url = await uploadCatalogImage(file, 'album-cover')

    toast.add({
      severity: 'success',
      summary: 'Cover uploaded',
      life: 2000,
    })
  } catch {
    errorMessage.value = 'Could not upload album cover. Please try again.'
  } finally {
    input.value = ''
  }
}

function removeCover() {
  form.cover_url = ''
}

function getPayload() {
  return {
    title: form.title.trim(),
    artist_id: form.artist_id,
    cover_url: form.cover_url.trim() || null,
  }
}

async function submit() {
  if (!canSubmit.value) return

  errorMessage.value = ''

  try {
    if (editingAlbum.value) {
      const updated = await updateAlbum(editingAlbum.value.id, getPayload())

      toast.add({
        severity: 'success',
        summary: 'Album updated',
        detail: updated.title,
        life: 2500,
      })
    } else {
      const created = await createAlbum(getPayload())

      toast.add({
        severity: 'success',
        summary: 'Album created',
        detail: created.title,
        life: 2500,
      })
    }

    dialogVisible.value = false
    resetForm()
  } catch {
    errorMessage.value = 'Could not save album. Please try again.'
  }
}

async function handleDelete(id: string | number) {
  const confirmed = window.confirm('Delete this album?')
  if (!confirmed) return

  try {
    await deleteAlbum(id)

    toast.add({
      severity: 'success',
      summary: 'Album deleted',
      life: 2500,
    })
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Delete failed',
      detail: 'Could not delete album.',
      life: 3000,
    })
  }
}
</script>
