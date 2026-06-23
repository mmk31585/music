<template>
  <Dialog v-model:visible="visibleInternal" modal :header="dialogTitle" :style="{ width: '520px' }">
    <div class="space-y-4 py-2">
      <p class="text-sm text-slate-400">
        {{ description }}
      </p>

      <FileUpload
        :name="fieldName"
        :choose-label="chooseLabel"
        :auto="false"
        :custom-upload="true"
        :disabled="uploading"
        :accept="accept"
        :max-file-size="maxFileSize"
        drag-drop
        @uploader="onCustomUpload"
        class="w-full"
      />

      <ProgressBar v-if="uploading || progress > 0" :value="progress" />

      <div v-if="uploadedUrl" class="rounded-xl bg-white/5 px-3 py-2 text-xs text-slate-300">
        Uploaded URL:
        <span class="font-mono break-all text-emerald-300">{{ uploadedUrl }}</span>
      </div>

      <div class="flex justify-end gap-2 pt-4">
        <Button label="Close" severity="secondary" text :disabled="uploading" @click="close" />
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import FileUpload, { type FileUploadUploaderEvent } from 'primevue/fileupload'
import ProgressBar from 'primevue/progressbar'
import { useToast } from 'primevue/usetoast'
import { useMediaApi, type UploadResponse } from '@/services/api/media'
import type { UploadFieldName } from '@/services/api/media/routes'

type UploadKind = 'track-audio' | 'track-cover' | 'album-cover' | 'artist-image'

const props = defineProps<{
  modelValue: boolean
  kind: UploadKind
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'uploaded', url: string): void
}>()

const toast = useToast()
const mediaApi = useMediaApi()
const uploading = ref(false)
const uploadedUrl = ref('')
const progress = ref(0)

const visibleInternal = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
})

const fieldName = computed(() => {
  switch (props.kind) {
    case 'track-audio':
      return 'trackAudio'
    case 'track-cover':
      return 'trackCover'
    case 'album-cover':
      return 'albumCover'
    case 'artist-image':
      return 'artistImage'
    default:
      return ''
  }
})

const accept = computed(() => {
  switch (props.kind) {
    case 'track-audio':
      return 'audio/*'
    default:
      return 'image/*'
  }
})

const maxFileSize = computed(() => (props.kind === 'track-audio' ? 30 : 5) * 1024 * 1024)

const chooseLabel = computed(() => {
  switch (props.kind) {
    case 'track-audio':
      return 'Choose audio'
    default:
      return 'Choose image'
  }
})

const dialogTitle = computed(() => {
  switch (props.kind) {
    case 'track-audio':
      return 'Upload track audio'
    case 'track-cover':
      return 'Upload track cover'
    case 'album-cover':
      return 'Upload album cover'
    case 'artist-image':
      return 'Upload artist image'
    default:
      return 'Upload media'
  }
})

const description = computed(() => {
  switch (props.kind) {
    case 'track-audio':
      return 'Upload an audio file for the track.'
    case 'track-cover':
      return 'Upload a cover image for the track.'
    case 'album-cover':
      return 'Upload a cover image for the album.'
    case 'artist-image':
      return 'Upload an image for the artist.'
    default:
      return ''
  }
})

watch(visibleInternal, (val) => {
  if (!val) {
    uploading.value = false
    uploadedUrl.value = ''
    progress.value = 0
  }
})

async function onCustomUpload(event: FileUploadUploaderEvent) {
  const files = Array.isArray(event.files) ? event.files : event.files ? [event.files] : []
  const file = files[0]
  if (!file) return

  uploading.value = true
  uploadedUrl.value = ''
  progress.value = 0

  try {
    const response: UploadResponse = await mediaApi.uploadAdminMedia(
      fieldName.value as UploadFieldName,
      file,
      {
        silent: false,
      },
      {
        onUploadProgress(progressEvent) {
          if (!progressEvent.total) return
          progress.value = Math.round((progressEvent.loaded / progressEvent.total) * 100)
        },
      },
    )

    const url = response?.url
    if (!url) {
      throw new Error('No URL returned from upload')
    }

    uploadedUrl.value = url
    emit('uploaded', url)

    toast.add({
      severity: 'success',
      summary: response.duplicate ? 'Already uploaded' : 'Upload complete',
      detail: file.name,
      life: 2500,
    })
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Upload failed',
      detail: 'Could not upload file',
      life: 3000,
    })
  } finally {
    uploading.value = false
    progress.value = progress.value === 0 ? 0 : 100
  }
}

function close() {
  visibleInternal.value = false
}
</script>
