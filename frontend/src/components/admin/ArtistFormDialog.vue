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
        <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-500/10">
          <i class="pi pi-user text-emerald-400" />
        </div>
        <div>
          <h3 class="text-base font-semibold text-white">
            {{ isEditing ? 'Edit artist' : 'New artist' }}
          </h3>
          <p class="text-xs text-slate-500">
            {{ isEditing ? 'Update artist information' : 'Add a new artist to the catalog' }}
          </p>
        </div>
      </div>
    </template>

    <form @submit.prevent="handleSubmit" class="mt-4 space-y-5">
      <!-- Name -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">
          Name <span class="text-red-400">*</span>
        </label>
        <InputText
          v-model="form.name"
          placeholder="Artist name"
          class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white placeholder:!text-slate-600 focus:!border-emerald-500/40 focus:!ring-1 focus:!ring-emerald-500/20"
          :invalid="!!errors.name"
          autofocus
        />
        <small v-if="errors.name" class="mt-1 block text-xs text-red-400">{{ errors.name }}</small>
      </div>

      <!-- Bio -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">Bio</label>
        <Textarea
          v-model="form.bio"
          placeholder="Short biography..."
          rows="3"
          auto-resize
          class="w-full !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white placeholder:!text-slate-600 focus:!border-emerald-500/40 focus:!ring-1 focus:!ring-emerald-500/20"
        />
      </div>

      <!-- Image URL -->
      <div>
        <label class="mb-1.5 block text-xs font-medium text-slate-400">Image URL</label>
        <div class="flex gap-2">
          <InputText
            v-model="form.image_url"
            placeholder="https://example.com/artist.jpg"
            class="flex-1 !rounded-xl !border-white/[0.08] !bg-white/[0.03] !text-white placeholder:!text-slate-600 focus:!border-emerald-500/40 focus:!ring-1 focus:!ring-emerald-500/20"
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

      <!-- Image Preview -->
      <div
        v-if="form.image_url"
        class="overflow-hidden rounded-xl border border-white/[0.06] bg-white/[0.02]"
      >
        <img
          :src="form.image_url"
          :alt="form.name"
          class="h-32 w-full object-cover"
          @error="onImageError"
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
          :label="isEditing ? 'Save changes' : 'Create artist'"
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
import Textarea from 'primevue/textarea'
import { useToast } from 'primevue/usetoast'
import { useMediaApi } from '@/services/api/media/routes'
import type { Artist } from '@/services/api/catalog/artists'
import type { ArtistFormPayload } from '@/composables/admin/useAdminArtists'

const visible = defineModel<boolean>({ default: false })

const props = defineProps<{
  artist?: Artist | null
  saving?: boolean
}>()

const emit = defineEmits<{
  submit: [payload: ArtistFormPayload]
}>()

const isEditing = computed(() => !!props.artist)
const toast = useToast()
const mediaApi = useMediaApi()
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

const form = reactive<ArtistFormPayload>({
  name: '',
  bio: null,
  image_url: null,
})

const errors = reactive<Record<string, string>>({})

watch(
  () => props.artist,
  (artist) => {
    if (artist) {
      form.name = artist.name
      form.bio = artist.bio
      form.image_url = artist.image_url
    } else {
      form.name = ''
      form.bio = null
      form.image_url = null
    }
    clearErrors()
  },
  { immediate: true },
)

watch(visible, (val) => {
  if (!val) {
    clearErrors()
  }
})

function clearErrors() {
  Object.keys(errors).forEach((key) => delete errors[key])
}

function validate(): boolean {
  clearErrors()
  if (!form.name?.trim()) {
    errors.name = 'Name is required'
    return false
  }
  return true
}

function handleSubmit() {
  if (!validate()) return
  emit('submit', { ...form })
}

function onImageError(e: Event) {
  const img = e.target as HTMLImageElement
  img.style.display = 'none'
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
    const res = await mediaApi.uploadAdminMedia('artistImage', file)
    if (res?.url) {
      form.image_url = res.url
      toast.add({ severity: 'success', summary: 'Uploaded', detail: 'Image uploaded successfully', life: 3000 })
    }
  } catch {
    toast.add({ severity: 'error', summary: 'Upload failed', detail: 'Could not upload image', life: 4000 })
  } finally {
    uploading.value = false
    input.value = ''
  }
}
</script>
