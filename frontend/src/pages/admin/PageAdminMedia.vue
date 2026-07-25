<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
    title="Media"
    eyebrow="Upload & manage"
    description="Upload audio files and cover images for your catalog"
  />

  <div class="mt-6 grid gap-5 md:grid-cols-2 xl:grid-cols-4 reveal-stagger">
    <div
      v-for="item in uploadTypes"
      :key="item.kind"
      class="group cursor-pointer rounded-xl border border-white/6 bg-white/3 p-5 transition hover:border-white/12 hover:bg-white/5"
      @click="openUpload(item.kind)"
    >
      <div class="flex items-center gap-4">
        <div
          class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl"
          :class="item.iconBg"
        >
          <i aria-hidden="true" :class="[item.icon, item.iconColor]" />
        </div>
        <div class="min-w-0 flex-1">
          <h3 class="text-sm font-bold text-white">{{ item.title }}</h3>
          <p class="mt-0.5 text-xs text-slate-500">{{ item.subtitle }}</p>
        </div>
        <Upload aria-hidden="true" class="text-xs text-slate-600 opacity-0 transition group-hover:opacity-100"  />
      </div>
    </div>
  </div>

  <div class="mt-6 rounded-xl border border-white/6 bg-white/3 p-5 reveal-fade">
    <div class="mb-3 flex items-center justify-between">
      <h3 class="text-sm font-bold text-white">All Media ({{ mediaList.length }})</h3>
      <Button
        icon="pi pi-refresh"
        text
        size="small"
        class="text-slate-500! hover:text-white!"
        @click="fetchMediaList"
      />
    </div>
    <div v-if="loading" class="space-y-2">
      <div v-for="i in 3" :key="i" class="h-14 animate-pulse rounded-lg bg-white/4" />
    </div>
    <AdminEmptyState
      v-else-if="loadError"
      icon="pi pi-exclamation-triangle"
      title="Failed to load media"
      :description="loadError"
    >
      <template #action>
        <Button label="Retry" icon="pi pi-refresh" @click="fetchMediaList" />
      </template>
    </AdminEmptyState>
    <div v-else-if="paginatedItems.length" class="space-y-2">
      <div
        v-for="item in paginatedItems"
        :key="item.id"
        class="group flex items-center gap-3 rounded-lg bg-white/4 px-4 py-2.5"
      >
        <Image aria-hidden="true" :class="item.mediaType === 'audio' ? 'pi text-blue-400' : 'pi text-green-400'"
          class="text-sm" />
        <div class="min-w-0 flex-1">
          <p class="truncate text-sm text-white">{{ item.originalFilename || item.objectKey.split('/').pop() }}</p>
          <p class="text-xs text-slate-500">
            {{ formatFileSize(item.fileSize) }} &bull; {{ item.mimeType || '—' }} &bull; {{ formatDate(item.createdAt) }}
          </p>
        </div>
        <Button
          icon="pi pi-external-link"
          text
          rounded
          size="small"
          class="h-7! w-7! text-slate-600! opacity-0 transition hover:text-white! group-hover:opacity-100"
          @click="openUrl(item.publicUrl)"
          v-tooltip.top="'Open'"
        />
        <Button
          icon="pi pi-trash"
          text
          rounded
          size="small"
          class="h-7! w-7! text-slate-600! opacity-0 transition hover:text-red-400! group-hover:opacity-100"
          @click="confirmDelete(item)"
          v-tooltip.top="'Delete'"
        />
      </div>
    </div>
    <AdminEmptyState
      v-else
      icon="pi pi-inbox"
      title="No media yet"
      description="Upload audio files or cover images to get started."
    />
    <div v-if="mediaList.length > pageSize" class="mt-4">
      <Paginator
        v-model:first="firstRecord"
        :rows="pageSize"
        :total-records="mediaList.length"
        :pt="{
          root: 'bg-transparent! border-white/6!',
          page: 'text-slate-400!',
          pageSelected: 'bg-spotify! text-black!',
        }"
      />
    </div>
  </div>

  <UploadMediaDialog
    v-model="showUploadDialog"
    :kind="selectedKind"
    @uploaded="onUploaded"
  />

  <AdminDeleteConfirm
    v-model="showDeleteConfirm"
    title="Delete Media"
    :item-name="deleteTarget?.originalFilename || 'this file'"
    :deleting="deleting"
    @confirm="handleDelete"
  />
  </div>
</template>

<script setup lang="ts">
import { Image, Upload } from 'lucide-vue-next'
// TODO MEDIUM: Media from ingestion flow (uploaded audio/covers) don't appear here — they're in ingestion_drafts, not media table.
import { ref, computed, onMounted } from 'vue'
import AdminDeleteConfirm from '@/components/admin/AdminDeleteConfirm.vue'
import AdminEmptyState from '@/components/admin/AdminEmptyState.vue'
import { useToast } from 'primevue/usetoast'
import { AdminSectionHeader } from '@/components/admin'
import UploadMediaDialog from '@/components/admin/UploadMediaDialog.vue'
import Paginator from 'primevue/paginator'
import { useMediaApi } from '@/services/api/media'
import type { Media } from '@/services/api/media'

type UploadKind = 'track-audio' | 'track-cover' | 'album-cover' | 'artist-image'

interface UploadTypeConfig {
  kind: UploadKind
  title: string
  subtitle: string
  icon: string
  iconBg: string
  iconColor: string
}

const uploadTypes: UploadTypeConfig[] = [
  {
    kind: 'track-audio',
    title: 'Track Audio',
    subtitle: 'MP3, FLAC, OGG, WAV up to 30MB',
    icon: 'pi pi-music',
    iconBg: 'bg-blue-500/10',
    iconColor: 'text-blue-400',
  },
  {
    kind: 'track-cover',
    title: 'Track Cover',
    subtitle: 'JPEG, PNG, WebP up to 5MB',
    icon: 'pi pi-image',
    iconBg: 'bg-emerald-500/10',
    iconColor: 'text-emerald-400',
  },
  {
    kind: 'album-cover',
    title: 'Album Cover',
    subtitle: 'JPEG, PNG, WebP up to 5MB',
    icon: 'pi pi-book',
    iconBg: 'bg-purple-500/10',
    iconColor: 'text-purple-400',
  },
  {
    kind: 'artist-image',
    title: 'Artist Image',
    subtitle: 'JPEG, PNG, WebP up to 5MB',
    icon: 'pi pi-user',
    iconBg: 'bg-amber-500/10',
    iconColor: 'text-amber-400',
  },
]

const mediaApi = useMediaApi()
const toast = useToast()

const showUploadDialog = ref(false)
const selectedKind = ref<UploadKind>('track-audio')

const mediaList = ref<Media[]>([])
const loading = ref(false)
const loadError = ref('')
const showDeleteConfirm = ref(false)
const deleteTarget = ref<Media | null>(null)
const deleting = ref(false)

// ── Client-side pagination ──
const pageSize = 10
const firstRecord = ref(0)

const paginatedItems = computed(() => {
  const end = firstRecord.value + pageSize
  return mediaList.value.slice(firstRecord.value, end)
})

async function fetchMediaList() {
  loading.value = true
  loadError.value = ''
  try {
    mediaList.value = await mediaApi.listAdminMedia()
    firstRecord.value = 0
  } catch (err) {
    loadError.value = err instanceof Error ? err.message : 'Failed to load media'
  } finally {
    loading.value = false
  }
}

onMounted(fetchMediaList)

function openUpload(kind: UploadKind) {
  selectedKind.value = kind
  showUploadDialog.value = true
}

function onUploaded() {
  fetchMediaList()
}

function confirmDelete(item: Media) {
  deleteTarget.value = item
  showDeleteConfirm.value = true
}

async function handleDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await mediaApi.deleteAdminMedia(String(deleteTarget.value.id))
    toast.add({ severity: 'success', summary: 'Media deleted', life: 2500 })
    showDeleteConfirm.value = false
    deleteTarget.value = null
    fetchMediaList()
  } catch {
    toast.add({ severity: 'error', summary: 'Delete failed', life: 3000 })
  } finally {
    deleting.value = false
  }
}

function openUrl(url?: string | null) {
  if (url) window.open(url, '_blank')
}

function formatFileSize(bytes?: number | null): string {
  if (bytes == null) return '—'
  if (bytes >= 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  if (bytes >= 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${bytes} B`
}

function formatDate(dateStr?: string | null): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString()
}
</script>
