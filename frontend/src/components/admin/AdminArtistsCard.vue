<template>
  <div class="rounded-3xl border border-white/10 bg-white/5 p-6 shadow-lg">
    <div class="mb-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
      <div>
        <h2 class="text-xl font-bold text-white">Artists</h2>
        <p class="mt-2 text-sm text-slate-400">Manage artist profiles and images.</p>
      </div>

      <Button
        label="Add artist"
        icon="pi pi-plus"
        class="border-0 bg-[#1db954] text-black"
        @click="openCreate"
      />
    </div>

    <DataTable
      :value="artists"
      :loading="loading"
      data-key="id"
      responsive-layout="scroll"
      class="overflow-hidden rounded-2xl"
    >
      <Column header="Image">
        <template #body="{ data }">
          <img
            v-if="data.image_url"
            :src="data.image_url"
            :alt="data.name"
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

      <Column field="name" header="Name" sortable />

      <Column field="bio" header="Bio">
        <template #body="{ data }">
          <span class="line-clamp-2 max-w-md text-sm text-slate-300">
            {{ data.bio || '-' }}
          </span>
        </template>
      </Column>

      <Column field="image_url" header="Image URL">
        <template #body="{ data }">
          <span class="line-clamp-1 max-w-md text-sm text-slate-300">
            {{ data.image_url || '-' }}
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
        <div class="py-8 text-center text-sm text-slate-400">No artists found.</div>
      </template>
    </DataTable>

    <Dialog
      v-model:visible="dialogVisible"
      modal
      :header="editingArtist ? 'Edit artist' : 'Add artist'"
      class="w-full max-w-lg"
    >
      <div class="flex flex-col gap-4">
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300"> Name </label>

          <InputText
            v-model="form.name"
            class="w-full"
            placeholder="Artist name"
            autofocus
            :disabled="saving || uploadingImage"
          />
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300"> Bio </label>

          <Textarea
            v-model="form.bio"
            class="w-full"
            rows="4"
            auto-resize
            placeholder="Short artist bio"
            :disabled="saving || uploadingImage"
          />

          <p class="mt-1 text-xs text-slate-500">{{ form.bio.length }}/5000</p>
        </div>

        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300"> Artist image </label>

          <input
            type="file"
            accept="image/*"
            class="block w-full rounded-xl border border-white/10 bg-black/20 px-3 py-2 text-sm text-slate-300 file:mr-4 file:rounded-lg file:border-0 file:bg-[#1db954] file:px-4 file:py-2 file:text-sm file:font-semibold file:text-black"
            :disabled="saving || uploadingImage"
            @change="handleArtistImageChange"
          />

          <p v-if="uploadingImage" class="mt-2 text-xs text-slate-400">Uploading image...</p>

          <p v-if="uploadImageError" class="mt-2 text-xs text-red-300">
            {{ uploadImageError }}
          </p>

          <p v-if="form.image_url" class="mt-2 line-clamp-1 text-xs text-slate-500">
            {{ form.image_url }}
          </p>
        </div>

        <div v-if="form.image_url" class="rounded-2xl bg-black/20 p-3">
          <img
            :src="form.image_url"
            alt="Artist preview"
            class="h-32 w-32 rounded-2xl object-cover"
          />

          <Button
            type="button"
            label="Remove image"
            severity="danger"
            text
            size="small"
            class="mt-3 px-0"
            :disabled="saving || uploadingImage"
            @click="removeImage"
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
          :label="editingArtist ? 'Save changes' : 'Create artist'"
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
import Textarea from 'primevue/textarea'
import { useToast } from 'primevue/usetoast'

import { useAdminArtists } from '@/composables/admin'
import { useCatalogImageUpload } from '@/composables/media/useCatalogImageUpload'
import type { Artist } from '@/services/api/catalog'

const toast = useToast()

const {
  artists,
  loading,
  saving,
  deleting,
  fetchArtists,
  createArtist,
  updateArtist,
  deleteArtist,
} = useAdminArtists()

const { uploadingImage, uploadImageError, uploadCatalogImage } = useCatalogImageUpload()

const dialogVisible = ref(false)
const editingArtist = ref<Artist | null>(null)
const errorMessage = ref('')

const form = reactive({
  name: '',
  bio: '',
  image_url: '',
})

const canSubmit = computed(() => {
  return form.name.trim().length > 0 && !saving.value && !uploadingImage.value
})

onMounted(() => {
  fetchArtists()
})

function resetForm() {
  form.name = ''
  form.bio = ''
  form.image_url = ''
  errorMessage.value = ''
  editingArtist.value = null
}

function openCreate() {
  resetForm()
  dialogVisible.value = true
}

function openEdit(artist: Artist) {
  editingArtist.value = artist
  form.name = artist.name
  form.bio = artist.bio || ''
  form.image_url = artist.image_url || ''
  errorMessage.value = ''
  dialogVisible.value = true
}

async function handleArtistImageChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  if (!file) return

  errorMessage.value = ''

  try {
    form.image_url = await uploadCatalogImage(file, 'artist-image')
      console.log(form.image_url)
    toast.add({
      severity: 'success',
      summary: 'Image uploaded',
      life: 2000,
    })
  } catch {
    errorMessage.value = 'Could not upload artist image. Please try again.'
  } finally {
    input.value = ''
  }
}

function removeImage() {
  form.image_url = ''
}

async function submit() {
  if (!canSubmit.value) return

  errorMessage.value = ''

  try {
    if (editingArtist.value) {
      const updated = await updateArtist(editingArtist.value.id, {
        name: form.name.trim(),
        bio: form.bio.trim() || null,
        image_url: form.image_url.trim() || null,
      })

      toast.add({
        severity: 'success',
        summary: 'Artist updated',
        detail: updated.name,
        life: 2500,
      })
    } else {
      const created = await createArtist({
        name: form.name.trim(),
        bio: form.bio.trim() || null,
        image_url: form.image_url.trim() || null,
      })

      toast.add({
        severity: 'success',
        summary: 'Artist created',
        detail: created.name,
        life: 2500,
      })
    }

    dialogVisible.value = false
    resetForm()
  } catch {
    errorMessage.value = 'Could not save artist. Please try again.'
  }
}

async function handleDelete(id: string | number) {
  const confirmed = window.confirm('Delete this artist?')
  if (!confirmed) return

  try {
    await deleteArtist(id)

    toast.add({
      severity: 'success',
      summary: 'Artist deleted',
      life: 2500,
    })
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Delete failed',
      detail: 'Could not delete artist.',
      life: 3000,
    })
  }
}
</script>
