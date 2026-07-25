<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUploader } from '@/composables/useUploader'
import { useReputation } from '@/composables/useReputation'
import { usePermissions } from '@/composables/usePermissions'
import { useUserAuthStore } from '@/stores'
import { Upload, FileAudio, Shield, CheckCircle, Clock, AlertTriangle, X, Loader2, Info } from 'lucide-vue-next'

const router = useRouter()
const auth = useUserAuthStore()
const { canUpload, canReview } = usePermissions()
const { drafts, slots, loading, fetchMyDrafts, fetchSlots, submitUpload } = useUploader()
const rep = useReputation(auth.user?.id)

const activeTab = ref<'upload' | 'drafts'>('upload')
const uploadFile = ref<File | null>(null)
const uploading = ref(false)
const uploadError = ref<string | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

const hasSlots = computed(() => slots.value?.available ?? false)
const maxSlots = computed(() => slots.value?.max ?? 0)
const usedSlots = computed(() => slots.value?.used ?? 0)
const slotPercent = computed(() => maxSlots.value > 0 ? Math.round((usedSlots.value / maxSlots.value) * 100) : 0)

const acceptedDrafts = computed(() => drafts.value.filter(d => d.status === 'accepted').length)
const pendingDrafts = computed(() => drafts.value.filter(d => d.status === 'pending' || d.status === 'needs_review').length)
const rejectedDrafts = computed(() => drafts.value.filter(d => d.status === 'rejected').length)

async function handleUpload() {
  if (!uploadFile.value) return
  uploading.value = true
  uploadError.value = null
  try {
    await submitUpload(uploadFile.value.name, uploadFile.value)
    uploadFile.value = null
  } catch (err: unknown) {
    uploadError.value = err instanceof Error ? err.message : 'Upload failed. Try again.'
  } finally {
    uploading.value = false
  }
}

function handleFileSelected(e: Event) {
  const target = e.target as HTMLInputElement
  uploadError.value = null
  if (target.files?.length) uploadFile.value = target.files[0] as File
}

function selectFile() {
  fileInput.value?.click()
}

onMounted(async () => {
  await Promise.all([fetchMyDrafts(), fetchSlots(), rep.fetchReputation()])
})
</script>

<template>
  <div class="mx-auto w-full max-w-4xl px-4 pb-32 pt-8 md:px-6">
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-3xl font-black text-white">Upload Music</h1>
      <p class="mt-1 text-sm text-white/40">
        Share your music with the community. Tracks go through review based on your trust tier.
      </p>
    </div>

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="h-24 animate-pulse rounded-2xl bg-white/5" />
    </div>

    <template v-else>
      <!-- Reputation bar -->
      <div class="mb-6 flex flex-wrap items-center gap-4 rounded-2xl border border-white/6 bg-white/3 px-5 py-3">
        <div class="flex items-center gap-3">
          <div
            class="flex h-9 w-9 items-center justify-center rounded-full"
            :class="rep.canAutoPublish ? 'bg-emerald-500/10' : 'bg-blue-500/10'"
          >
            <Shield aria-hidden="true" class="text-sm" :class="rep.canAutoPublish ? 'text-emerald-400' : 'text-blue-400'" />
          </div>
          <div>
            <p class="text-xs font-medium text-white/70">{{ rep.tierLabel }}</p>
            <p class="text-[10px] text-white/30">Trust Tier</p>
          </div>
        </div>
        <div class="h-8 w-px bg-white/6" />
        <div class="text-sm tabular-nums">
          <span class="font-semibold text-white">{{ rep.trustScore }}</span>
          <span class="ml-1 text-xs text-white/30">Score</span>
        </div>
        <div class="h-8 w-px bg-white/6" />
        <div class="text-sm tabular-nums">
          <span class="font-semibold text-white">{{ rep.acceptedContributions }}</span>
          <span class="ml-1 text-xs text-white/30">Accepted</span>
        </div>
        <div class="ml-auto">
          <span
            v-if="rep.canAutoPublish"
            class="rounded-full bg-emerald-500/10 px-3 py-1 text-[10px] font-medium text-emerald-400"
          >
            Auto-publish
          </span>
          <span v-else class="rounded-full bg-amber-500/10 px-3 py-1 text-[10px] font-medium text-amber-400">
            Needs review
          </span>
        </div>
      </div>

      <!-- Slot gauge -->
      <div class="mb-6 rounded-2xl border border-white/6 bg-white/3 p-5">
        <div class="mb-3 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <Upload aria-hidden="true" class="text-sm text-white/40" />
            <span class="text-sm font-medium text-white/70">Upload Slots</span>
          </div>
          <span class="text-sm tabular-nums text-white/70">
            <strong class="text-white">{{ usedSlots }}</strong> / {{ maxSlots }} used
          </span>
        </div>
        <div class="h-2 overflow-hidden rounded-full bg-white/8">
          <div
            class="h-full rounded-full transition-all duration-500"
            :class="slotPercent >= 80 ? 'bg-amber-500' : slotPercent >= 50 ? 'bg-spotify' : 'bg-blue-500'"
            :style="{ width: slotPercent + '%' }"
          />
        </div>
        <p v-if="!hasSlots" class="mt-2 flex items-center gap-1.5 text-xs text-amber-400">
          <AlertTriangle aria-hidden="true" class="text-[10px]" />
          All slots full. Wait for pending reviews to complete.
        </p>
        <p v-else class="mt-2 text-xs text-white/30">
          {{ maxSlots - usedSlots }} slot{{ maxSlots - usedSlots !== 1 ? 's' : '' }} available
        </p>
      </div>

      <!-- Tabs -->
      <div class="mb-6 flex gap-1 rounded-xl bg-white/4 p-1">
        <button
          class="flex flex-1 items-center justify-center gap-2 rounded-lg py-2.5 text-sm font-medium transition"
          :class="activeTab === 'upload' ? 'bg-white/10 text-white shadow-xs' : 'text-white/40 hover:text-white/60'"
          @click="activeTab = 'upload'"
        >
          <Upload aria-hidden="true" class="text-xs" />
          Upload
        </button>
        <button
          class="flex flex-1 items-center justify-center gap-2 rounded-lg py-2.5 text-sm font-medium transition"
          :class="activeTab === 'drafts' ? 'bg-white/10 text-white shadow-xs' : 'text-white/40 hover:text-white/60'"
          @click="activeTab = 'drafts'"
        >
          <FileAudio aria-hidden="true" class="text-xs" />
          My Drafts
          <span v-if="drafts.length" class="rounded-full bg-white/10 px-1.5 py-0.5 text-[10px]">{{ drafts.length }}</span>
        </button>
      </div>

      <!-- Upload tab -->
      <div v-show="activeTab === 'upload'" class="space-y-6">
        <div
          class="flex flex-col items-center justify-center gap-5 rounded-2xl border-2 border-dashed border-white/10 px-8 py-16 text-center transition hover:border-white/20"
          @drag.prevent
          @drop.prevent
        >
          <input
            type="file"
            accept="audio/*,.mp3,.flac,.wav,.ogg,.aac,.m4a"
            class="hidden"
            @change="handleFileSelected"
            ref="fileInput"
          />

          <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/5">
            <Upload aria-hidden="true" class="text-2xl text-white/20" />
          </div>

          <div v-if="!uploadFile">
            <p class="text-base font-semibold text-white/70">Drop your audio file here</p>
            <p class="mt-1 text-sm text-white/30">or click to browse — MP3, FLAC, WAV, OGG, AAC</p>
            <button
              class="mt-5 rounded-full bg-spotify px-8 py-2.5 text-sm font-bold text-black transition hover:scale-105 hover:bg-spotify-hover disabled:opacity-50"
              :disabled="!hasSlots || uploading"
              @click="selectFile"
            >
              Select Audio File
            </button>
          </div>

          <div v-else class="w-full max-w-sm">
            <div class="flex items-center gap-3 rounded-xl bg-white/5 p-4">
              <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-spotify/10">
                <FileAudio aria-hidden="true" class="text-spotify" />
              </div>
              <div class="min-w-0 flex-1 text-left">
                <p class="truncate text-sm font-medium text-white">{{ uploadFile.name }}</p>
                <p class="text-xs text-white/40">{{ (uploadFile.size / 1024 / 1024).toFixed(1) }} MB</p>
              </div>
              <button
                class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-white/30 transition hover:bg-white/10 hover:text-white/60"
                @click="uploadFile = null"
                aria-label="Remove file"
              >
                <X aria-hidden="true" class="text-xs" />
              </button>
            </div>

            <button
              class="mt-4 w-full rounded-full bg-spotify px-8 py-2.5 text-sm font-bold text-black transition hover:bg-spotify-hover disabled:opacity-50"
              :disabled="uploading"
              @click="handleUpload"
            >
              <Loader2 v-if="uploading" aria-hidden="true" class="mr-2 inline animate-spin text-xs" />
              {{ uploading ? 'Uploading...' : 'Start Upload' }}
            </button>

            <p v-if="uploadError" class="mt-2 flex items-center justify-center gap-1.5 text-xs text-red-400">
              <AlertTriangle aria-hidden="true" class="text-[10px]" />
              {{ uploadError }}
            </p>
          </div>

          <p v-if="!hasSlots && !uploadFile" class="flex items-center gap-1.5 text-xs text-amber-400">
            <AlertTriangle aria-hidden="true" class="text-[10px]" />
            All upload slots are full. Wait for pending reviews to complete.
          </p>
        </div>

        <!-- Trust info -->
        <div class="flex items-start gap-3 rounded-xl bg-white/3 p-4">
          <Info aria-hidden="true" class="mt-0.5 shrink-0 text-xs text-white/30" />
          <div class="text-xs text-white/30 leading-relaxed">
            <strong class="text-white/50">How uploads work:</strong>
            Based on your trust tier, uploads are auto-published or queued for review.
            Higher reputation earns more slots and faster publishing.
            <router-link to="/gamification" class="text-spotify hover:underline">Learn more</router-link>
          </div>
        </div>
      </div>

      <!-- Drafts tab -->
      <div v-show="activeTab === 'drafts'" class="space-y-4">
        <!-- Summary -->
        <div class="flex gap-3">
          <div class="flex-1 rounded-xl bg-white/4 px-4 py-3 text-center">
            <p class="text-lg font-bold text-white">{{ acceptedDrafts }}</p>
            <p class="text-[10px] text-white/30">Accepted</p>
          </div>
          <div class="flex-1 rounded-xl bg-white/4 px-4 py-3 text-center">
            <p class="text-lg font-bold text-amber-400">{{ pendingDrafts }}</p>
            <p class="text-[10px] text-white/30">Pending</p>
          </div>
          <div class="flex-1 rounded-xl bg-white/4 px-4 py-3 text-center">
            <p class="text-lg font-bold text-red-400">{{ rejectedDrafts }}</p>
            <p class="text-[10px] text-white/30">Rejected</p>
          </div>
        </div>

        <div v-if="drafts.length" class="space-y-2">
          <div
            v-for="d in drafts"
            :key="d.id"
            class="flex items-center gap-4 rounded-xl border border-white/6 bg-white/3 p-4"
          >
            <div
              class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg"
              :class="
                d.status === 'accepted'
                  ? 'bg-emerald-500/10'
                  : d.status === 'rejected'
                    ? 'bg-red-500/10'
                    : 'bg-amber-500/10'
              "
            >
              <CheckCircle
                v-if="d.status === 'accepted'"
                aria-hidden="true"
                class="text-emerald-400"
              />
              <X
                v-else-if="d.status === 'rejected'"
                aria-hidden="true"
                class="text-red-400"
              />
              <Clock v-else aria-hidden="true" class="text-amber-400" />
            </div>
            <div class="min-w-0 flex-1">
              <p class="truncate text-sm font-medium text-white">{{ d.originalFilename }}</p>
              <p class="text-xs text-white/40">
                {{ new Date(d.createdAt).toLocaleDateString() }}
                <span v-if="d.fileSize" class="ml-2">· {{ (d.fileSize / 1024 / 1024).toFixed(1) }} MB</span>
              </p>
            </div>
            <span
              class="shrink-0 rounded-full px-3 py-1 text-[10px] font-bold"
              :class="
                d.status === 'accepted'
                  ? 'bg-emerald-500/15 text-emerald-400'
                  : d.status === 'rejected'
                    ? 'bg-red-500/15 text-red-400'
                    : 'bg-amber-500/15 text-amber-400'
              "
            >
              {{ d.status === 'needs_review' ? 'Needs Review' : d.status }}
            </span>
          </div>
        </div>

        <div v-else class="flex flex-col items-center gap-3 rounded-2xl border border-white/6 bg-white/3 py-14 text-center">
          <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/5">
            <FileAudio aria-hidden="true" class="text-2xl text-white/15" />
          </div>
          <p class="text-sm text-white/30">No drafts yet</p>
          <button
            class="text-xs text-spotify transition-colors hover:text-spotify-hover"
            @click="activeTab = 'upload'"
          >
            Upload your first track
          </button>
        </div>
      </div>
    </template>
  </div>
</template>
