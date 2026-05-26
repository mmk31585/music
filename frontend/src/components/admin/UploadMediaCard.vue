<template>
  <div class="rounded-3xl border border-white/10 bg-white/5 p-6 shadow-lg">
    <div class="mb-6">
      <h2 class="text-xl font-bold text-white">Upload media</h2>
      <p class="mt-2 text-sm text-slate-400">
        Upload an audio file or asset through the admin media endpoint.
      </p>
    </div>

    <div class="rounded-2xl border border-dashed border-white/15 bg-black/20 p-6">
      <div class="flex flex-col gap-4">
        <div>
          <label class="mb-2 block text-sm font-medium text-slate-300"> Choose file </label>

          <input
            type="file"
            class="block w-full rounded-xl border border-white/10 bg-white/5 px-3 py-3 text-sm text-slate-300 file:mr-4 file:rounded-lg file:border-0 file:bg-[#1db954] file:px-4 file:py-2 file:text-sm file:font-semibold file:text-black"
            @change="onFileChange"
          />
        </div>

        <div v-if="fileName" class="rounded-xl bg-white/5 px-4 py-3 text-sm text-slate-300">
          Selected: <span class="font-medium text-white">{{ fileName }}</span>
        </div>

        <ProgressBar v-if="loading || progress > 0" :value="progress" />

        <div class="flex flex-wrap items-center gap-3">
          <Button
            label="Upload"
            icon="pi pi-upload"
            :loading="loading"
            :disabled="!selectedFile"
            class="border-0 bg-[#1db954] text-black"
            @click="submit"
          />

          <Button
            label="Clear"
            severity="secondary"
            outlined
            :disabled="loading || !selectedFile"
            @click="clear"
          />
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
import { ref } from 'vue'
import Button from 'primevue/button'
import ProgressBar from 'primevue/progressbar'
import { useToast } from 'primevue/usetoast'
import { useMediaApi } from '@/services/api/media'

const toast = useToast()
const mediaApi = useMediaApi()

const selectedFile = ref<File | null>(null)
const fileName = ref('')
const loading = ref(false)
const resultMessage = ref('')
const errorMessage = ref('')
const progress = ref(0)

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0] || null

  selectedFile.value = file
  fileName.value = file?.name || ''
  resultMessage.value = ''
  errorMessage.value = ''
  progress.value = 0
}

function clear() {
  selectedFile.value = null
  fileName.value = ''
  resultMessage.value = ''
  errorMessage.value = ''
  progress.value = 0
}

async function submit() {
  if (!selectedFile.value) return

  loading.value = true
  resultMessage.value = ''
  errorMessage.value = ''
  progress.value = 0

  try {
    const result = await mediaApi.uploadAdminMedia(
      'trackAudio',
      selectedFile.value,
      undefined,
      {
        onUploadProgress(event) {
          if (!event.total) return
          progress.value = Math.round((event.loaded / event.total) * 100)
        },
      },
    )

    resultMessage.value = result.duplicate
      ? 'This file was already uploaded. Existing media URL reused.'
      : 'File uploaded successfully.'
    toast.add({
      severity: 'success',
      summary: 'Upload complete',
      detail: selectedFile.value.name,
      life: 2500,
    })
    clear()
  } catch (error) {
    errorMessage.value = 'Upload failed. Please try again.'
    toast.add({
      severity: 'error',
      summary: 'Upload failed',
      detail: 'Could not upload the selected file',
      life: 3000,
    })
  } finally {
    loading.value = false
    progress.value = progress.value === 0 ? 0 : 100
  }
}
</script>
