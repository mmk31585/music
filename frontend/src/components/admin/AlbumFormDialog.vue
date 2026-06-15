<template>
  <Dialog
    v-model:visible="visible"
    modal
    :closable="!saving"
    :draggable="false"
    :style="{ width: '520px' }"
    :pt="{
      root: { class: '!border-white/[0.06] !bg-[#141414] !rounded-2xl !shadow-2xl' },
      header: { class: '!bg-transparent !border-0 !pb-2' },
      content: { class: '!bg-transparent !px-6 !pt-0 !pb-2' },
      footer: { class: '!bg-transparent !border-0' },
      mask: { class: '!backdrop-blur-sm' },
    }"
  >
    <template #header>
      <div class="flex items-center gap-3">
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-500/10">
          <i class="pi pi-book text-blue-400" />
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
          class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white placeholder:!text-slate-600 focus:!border-emerald-500/40"
          :invalid="!!errors.title"
          autofocus
        />
        <small v-if="errors.title" class="mt-1 block text-xs text-red-400">{{ errors.title }}</small>
      </div>

      <!-- Artist ID -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">Artist ID</label>
        <InputText
          v-model="form.artist_id"
          placeholder="Artist ID (optional)"
          class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white placeholder:!text-slate-600 focus:!border-emerald-500/40"
        />
      </div>

      <!-- Cover URL -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">Cover URL</label>
        <div class="flex gap-2">
          <InputText
            v-model="form.cover_url"
            placeholder="https://example.com/cover.jpg"
            class="flex-1 !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white placeholder:!text-slate-600 focus:!border-emerald-500/40"
          />
          <Button
            icon="pi pi-upload"
            severity="secondary"
            outlined
            :loading="uploading"
            class="!rounded-xl !border-white/[0.08] !bg-white/[0.03] hover:!bg-white/[0.08]"
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
        class="overflow-hidden rounded-xl border border-white/[0.06]"
      >
        <img
          :src="form.cover_url"
          :alt="form.title"
          class="h-40 w-full object-cover"
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
          class="!text-slate-400 hover:!text-white"
          @click="visible = false"
        />
        <Button
          :label="isEditing ? 'Save changes' : 'Create album'"
          :icon="isEditing ? 'pi pi-check' : 'pi pi-plus'"
          :loading="saving"
          class="!rounded-xl !bg-emerald-500 !text-black hover:!bg-emerald-400"
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
import { useToast } from 'primevue/usetoast'
import { useMediaApi } from '@/services/api/media/routes'
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

const isEditing = computed(() => !!props.album)
const toast = useToast()
const mediaApi = useMediaApi()
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

const form = reactive({
  title: '',
  cover_url: null as string | null,
  artist_id: null as string | null,
})

const errors = reactive<Record<string, string>>({})

watch(
  () => props.album,
  (album) => {
    if (album) {
      form.title = album.title
      form.cover_url = album.cover_url
      form.artist_id = album.artist_id != null ? String(album.artist_id) : null
    } else {
      form.title = ''
      form.cover_url = null
      form.artist_id = null
    }
    Object.keys(errors).forEach((key) => delete errors[key])
  },
  { immediate: true },
)

function handleSubmit() {
  Object.keys(errors).forEach((key) => delete errors[key])
  if (!form.title?.trim()) {
    errors.title = 'Title is required'
    return
  }
  emit('submit', { ...form })
}

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
</script>
