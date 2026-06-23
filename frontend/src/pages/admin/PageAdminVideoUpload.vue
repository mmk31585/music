<template>
  <div class="mx-auto w-full max-w-3xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Media"
      title="Video Upload"
      description="Upload an official music video (MV) for a track."
    />

    <!-- Step progress indicator -->
    <div class="mb-8 flex items-center gap-2">
      <div
        v-for="(step, i) in steps"
        :key="i"
        class="flex items-center gap-2"
      >
        <div
          class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-bold transition-all duration-300"
          :class="stepClass(i)"
        >
          <i v-if="stepComplete(i)" class="pi pi-check text-xs" />
          <span v-else>{{ i + 1 }}</span>
        </div>
        <span
          class="hidden text-xs font-medium sm:block"
          :class="currentStep >= i ? 'text-white/70' : 'text-white/20'"
        >
          {{ step.label }}
        </span>
        <i
          v-if="i < steps.length - 1"
          class="pi pi-chevron-left text-[10px] text-white/8"
        />
      </div>
    </div>

    <!-- Step 1: Select Track -->
    <section class="mb-8">
      <h3 class="mb-3 text-sm font-semibold text-white/80">Select Track</h3>
      <div class="relative">
        <i
          class="pi pi-search absolute right-3 top-1/2 -translate-y-1/2 text-sm text-slate-500"
        />
        <InputText
          v-model="searchQuery"
          placeholder="Search tracks by title or artist..."
          class="h-11! w-full! rounded-xl! border-white/8! bg-white/3! pr-10! text-sm! text-white! placeholder:text-slate-600!"
          @input="onSearchInput"
        />
        <i
          v-if="searching"
          class="pi pi-spin pi-spinner absolute left-3 top-1/2 -translate-y-1/2 text-xs text-slate-500"
        />
      </div>

      <TransitionGroup
        name="track-list"
        tag="div"
        class="mt-3 space-y-1"
      >
        <button
          v-for="track in searchResults"
          :key="String(track.id)"
          type="button"
          class="flex w-full items-center gap-3 rounded-xl px-4 py-3 text-right transition"
          :class="selectedTrack?.id === track.id
            ? 'bg-spotify/10 ring-1 ring-spotify/30'
            : 'hover:bg-white/6'"
          @click="selectTrack(track)"
        >
          <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
            <img
              v-if="track.cover_url"
              :src="track.cover_url"
              :alt="track.title"
              class="h-full w-full object-cover"
            />
            <div v-else class="flex h-full items-center justify-center">
              <i class="pi pi-music text-xs text-slate-500" />
            </div>
          </div>
          <div class="min-w-0 flex-1 text-right">
            <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
            <p class="mt-0.5 truncate text-xs text-white/40">{{ track.artist_name || 'Unknown artist' }}</p>
          </div>
          <i
            v-if="selectedTrack?.id === track.id"
            class="pi pi-check-circle text-spotify text-lg"
          />
        </button>

        <div
          v-if="searchQuery && !searching && searchResults.length === 0"
          class="py-8 text-center text-sm text-slate-500"
        >
          No tracks found for "{{ searchQuery }}"
        </div>
      </TransitionGroup>

      <!-- Track Preview Player -->
      <Transition name="fade-slide">
        <div
          v-if="selectedTrack && selectedTrack.audio_url"
          class="mt-4 rounded-2xl border border-white/8 bg-white/3 p-4"
        >
          <div class="mb-2 flex items-center gap-2 text-xs text-white/60">
            <i class="pi pi-headphones text-spotify" />
            <span>Audio preview — {{ selectedTrack.title }}</span>
          </div>
          <audio
            :src="selectedTrack.audio_url"
            controls
            preload="metadata"
            class="w-full rounded-lg"
            style="height: 40px;"
          >
            Your browser does not support the audio element.
          </audio>
        </div>
      </Transition>
    </section>

    <!-- Step 2: Upload Video File -->
    <section class="mb-8">
      <h3 class="mb-3 text-sm font-semibold text-white/80">Upload Video</h3>

      <div
        class="flex cursor-pointer flex-col items-center justify-center rounded-2xl border-2 border-dashed px-6 py-12 transition-all duration-200"
        :class="dropZoneClass"
        @click="triggerFileInput"
        @dragover.prevent="dragOver = true"
        @dragleave.prevent="dragOver = false"
        @drop.prevent="onFileDrop"
      >
        <div
          class="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl transition-all duration-300"
          :class="uploadedFile
            ? 'bg-spotify/20'
            : dragOver
              ? 'bg-spotify/10 scale-110'
              : 'bg-white/5'"
        >
          <i
            class="text-2xl transition-all duration-300"
            :class="uploadedFile
              ? 'pi pi-check-circle text-spotify'
              : dragOver
                ? 'pi pi-arrow-down text-spotify'
                : 'pi pi-cloud-upload text-slate-400'"
          />
        </div>
        <template v-if="!uploadedFile">
          <p class="text-sm font-medium text-white/70">
            {{ dragOver ? 'Drop file here' : 'Click or drag a video file' }}
          </p>
          <p class="mt-1 text-xs text-slate-500">MP4, WebM, MOV — up to 200 MB</p>
        </template>
        <template v-else>
          <p class="text-sm font-medium text-white/80">{{ uploadedFile.name }}</p>
          <p class="mt-1 text-xs text-slate-500">
            {{ (uploadedFile.size / (1024 * 1024)).toFixed(1) }} MB
          </p>
        </template>
        <input
          ref="fileInputRef"
          type="file"
          accept="video/mp4,video/webm,video/quicktime"
          class="hidden"
          @change="onFileChange"
        />
      </div>

      <!-- Upload progress -->
      <Transition name="fade-slide">
        <div v-if="uploading" class="mt-4">
          <div class="mb-1.5 flex items-center justify-between text-xs">
            <span class="text-white/60">Uploading to server...</span>
            <span class="tabular-nums text-white/80">{{ uploadProgress }}%</span>
          </div>
          <div class="h-2 overflow-hidden rounded-full bg-white/6">
            <div
              class="h-full rounded-full bg-spotify transition-all duration-300"
              :style="{ width: uploadProgress + '%' }"
            />
          </div>
        </div>
      </Transition>

      <!-- Upload success -->
      <Transition name="fade-slide">
        <div
          v-if="uploadedFileUrl && !uploading"
          class="mt-3 flex items-center gap-2 rounded-xl bg-emerald-500/10 px-4 py-3 text-sm text-emerald-400"
        >
          <i class="pi pi-check-circle" />
          <span>Video uploaded successfully — ready to submit.</span>
        </div>
      </Transition>
    </section>

    <!-- Step 3: Sync & Submit -->
    <section class="mb-8">
      <h3 class="mb-3 text-sm font-semibold text-white/80">Sync &amp; Submit</h3>

      <!-- Sync controls -->
      <Transition name="fade-slide">
        <div v-if="selectedTrack" class="mb-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">
              Audio start offset (ms)
            </label>
            <InputNumber
              v-model="trackStartMs"
              :min="0"
              :max="300000"
              :step="100"
              class="w-full"
              input-class="rounded-xl! border-white/8! bg-white/3! text-white! w-full!"
              placeholder="0"
            />
            <p class="mt-1 text-[10px] text-slate-500">
              Delay track audio start (use if audio/video drift)
            </p>
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-400">
              Audio end offset (ms)
            </label>
            <InputNumber
              v-model="trackEndMs"
              :min="0"
              :max="300000"
              :step="100"
              class="w-full"
              input-class="rounded-xl! border-white/8! bg-white/3! text-white! w-full!"
              placeholder="0"
            />
            <p class="mt-1 text-[10px] text-slate-500">
              Trim track audio early from the end
            </p>
          </div>
        </div>
      </Transition>

      <div class="mb-4 flex items-start gap-2 rounded-xl bg-blue-500/5 px-4 py-3 text-sm text-white/50">
        <i class="pi pi-info-circle mt-0.5 shrink-0 text-xs text-blue-400" />
        <span>The video will be processed (audio replaced with the selected track) and published as an Official MV.</span>
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <Button
          label="Upload Official MV"
          icon="pi pi-cloud-upload"
          class="rounded-xl! bg-spotify! px-6! text-black! hover:bg-spotify/90!"
          :disabled="selectedTrack! || uploadedFileUrl! || submitting"
          :loading="submitting"
          @click="submitMV"
        />
        <Button
          v-if="submittedVideoId"
          label="View Video"
          icon="pi pi-external-link"
          severity="secondary"
          class="rounded-xl!"
          @click="viewVideo"
        />
        <Button
          v-if="submittedVideoId"
          label="Back to Tracks"
          icon="pi pi-arrow-left"
          severity="info"
          class="rounded-xl!"
          @click="goToTracks"
        />
      </div>
    </section>

    <!-- Status card -->
    <Transition name="fade-slide">
      <div
        v-if="statusMessage"
        class="rounded-2xl border p-4 text-sm"
        :class="statusError
          ? 'border-red-500/20 bg-red-500/10 text-red-400'
          : 'border-emerald-500/20 bg-emerald-500/10 text-emerald-400'"
      >
        <div class="flex items-center gap-2">
          <i
            class="text-base"
            :class="statusError ? 'pi pi-exclamation-circle' : 'pi pi-check-circle'"
          />
          <span>{{ statusMessage }}</span>
        </div>
      </div>
    </Transition>

    <!-- Processing info after successful submit -->
    <Transition name="fade-slide">
      <div
        v-if="submittedVideoId && !statusError"
        class="mt-4 overflow-hidden rounded-2xl border border-white/8 bg-white/3"
      >
        <div class="flex items-center gap-3 bg-spotify/5 p-4">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-spotify/10">
            <i class="pi pi-spin pi-spinner text-spotify" />
          </div>
          <div>
            <p class="text-sm font-medium text-white/90">Processing your video...</p>
            <p class="mt-0.5 text-xs text-slate-500">This usually takes 1–5 minutes</p>
          </div>
        </div>
        <div class="border-t border-white/6 px-4 py-3">
          <p class="text-xs text-slate-500">
            The audio from your selected track will be mixed into the video. Processing happens in the background — you can navigate away and come back later.
          </p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useRequest } from '@/composables/useRequest'
import { useTracksApi } from '@/services/api/catalog/tracks/routes'
import { AdminSectionHeader } from '@/components/admin'

interface TrackResult {
  id: string | number
  title: string
  artist_name?: string
  cover_url?: string | null
  duration_seconds?: number | null
  audio_url?: string | null
}

const tracksApi = useTracksApi()

const router = useRouter()
const toast = useToast()

// ── Steps ──
const steps = [
  { label: 'Select Track', key: 'track' },
  { label: 'Upload Video', key: 'upload' },
  { label: 'Submit', key: 'submit' },
]

const currentStep = computed(() => {
  if (submittedVideoId.value) return 3
  if (uploadedFileUrl.value && selectedTrack.value) return 2
  if (selectedTrack.value) return 1
  return 0
})

function stepClass(i: number) {
  if (currentStep.value > i) return 'bg-spotify text-black'
  if (currentStep.value === i) return 'bg-spotify/20 text-spotify ring-1 ring-spotify/30'
  return 'bg-white/10 text-white/40'
}

function stepComplete(i: number) {
  return currentStep.value > i
}

// ── Track search ──
const searchQuery = ref('')
const searching = ref(false)
const searchResults = ref<TrackResult[]>([])
const selectedTrack = ref<TrackResult | null>(null)
let searchTimer: ReturnType<typeof setTimeout> | null = null

function onSearchInput() {
  if (searchTimer) clearTimeout(searchTimer)
  const q = searchQuery.value.trim()
  if (q!) {
    searchResults.value = []
    return
  }
  searching.value = true
  searchTimer = setTimeout(async () => {
    try {
      const res = await tracksApi.getTracks({ q, limit: 10 })
      searchResults.value = (res as unknown as TrackResult[]) || []
    } catch {
      searchResults.value = []
    } finally {
      searching.value = false
    }
  }, 300)
}

function selectTrack(track: TrackResult) {
  selectedTrack.value = track
  if (track.id) {
    useRequest<TrackResult>(`/catalog/tracks/${track.id}`, { method: 'GET' })
      .then((fullTrack) => {
        if (fullTrack && selectedTrack.value?.id === track.id) {
          selectedTrack.value = { ...selectedTrack.value, ...fullTrack }
        }
      })
      .catch(() => {})
  }
}

// ── Sync offsets ──
const trackStartMs = ref(0)
const trackEndMs = ref(0)

// ── File upload ──
const fileInputRef = ref<HTMLInputElement | null>(null)
const uploadedFile = ref<File | null>(null)
const uploadedFileUrl = ref('')
const uploading = ref(false)
const uploadProgress = ref(0)
const dragOver = ref(false)

const dropZoneClass = computed(() => {
  if (uploadedFile) return 'border-spotify/30 bg-spotify/2'
  if (dragOver.value) return 'border-spotify/50 bg-spotify/4 scale-[1.01]'
  return 'border-white/8 hover:border-spotify/30 hover:bg-white/2'
})

function triggerFileInput() {
  fileInputRef.value?.click()
}

function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) handleFile(file)
}

function onFileDrop(e: DragEvent) {
  dragOver.value = false
  const file = e.dataTransfer?.files?.[0]
  if (file) handleFile(file)
}

function handleFile(file: File) {
  uploadedFile.value = file
  uploadedFileUrl.value = ''
  uploadToMedia(file)
}

async function uploadToMedia(file: File) {
  uploading.value = true
  uploadProgress.value = 0
  try {
    const formData = new FormData()
    formData.append('video', file)

    const result = await useRequest<{ url: string; mediaId: string }>(
      '/admin/media/upload',
      {
        method: 'POST',
        data: formData,
        onUploadProgress: (e: ProgressEvent) => {
          if (e.total) {
            uploadProgress.value = Math.round((e.loaded / e.total) * 100)
          }
        },
      } as any,
    )
    uploadedFileUrl.value = result.url
    toast.add({ severity: 'success', summary: 'Video uploaded', life: 3000 })
  } catch (err: any) {
    const msg = err?.message || 'Upload failed'
    toast.add({ severity: 'error', summary: msg, life: 5000 })
    uploadedFile.value = null
  } finally {
    uploading.value = false
  }
}

// ── Submit MV ──
const submitting = ref(false)
const submittedVideoId = ref('')
const statusMessage = ref('')
const statusError = ref(false)

async function submitMV() {
  if (selectedTrack.value! || uploadedFileUrl.value!) return
  submitting.value = true
  statusMessage.value = ''
  statusError.value = false
  try {
    const result = await useRequest<{ id: string }>(
      `/admin/tracks/${selectedTrack.value.id}/mv`,
      {
        method: 'POST',
        data: {
          raw_video_url: uploadedFileUrl.value,
          title: `${selectedTrack.value.title} - Official Music Video`,
          track_start_ms: trackStartMs.value,
          track_end_ms: trackEndMs.value,
        },
      },
    )
    submittedVideoId.value = result.id
    statusMessage.value = 'MV submitted for processing!'
    toast.add({ severity: 'success', summary: 'MV uploaded successfully', life: 4000 })
  } catch (err: any) {
    statusError.value = true
    statusMessage.value = err?.message || 'Failed to submit MV'
    toast.add({ severity: 'error', summary: statusMessage.value, life: 5000 })
  } finally {
    submitting.value = false
  }
}

function viewVideo() {
  if (submittedVideoId.value) {
    router.push(`/music-video/${submittedVideoId.value}`)
  }
}

function goToTracks() {
  router.push('/admin/tracks')
}
</script>

<style scoped>
/* Track list transitions */
.track-list-enter-active {
  transition: all 0.25s ease-out;
}
.track-list-leave-active {
  transition: all 0.2s ease-in;
}
.track-list-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}
.track-list-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

/* Fade + slide transitions */
.fade-slide-enter-active {
  transition: all 0.3s ease-out;
}
.fade-slide-leave-active {
  transition: all 0.2s ease-in;
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

@media (prefers-reduced-motion: reduce) {
  .track-list-enter-active,
  .track-list-leave-active,
  .fade-slide-enter-active,
  .fade-slide-leave-active {
    transition: none;
  }
  .track-list-enter-from,
  .track-list-leave-to,
  .fade-slide-enter-from,
  .fade-slide-leave-to {
    opacity: 1;
    transform: none;
  }
}
</style>
