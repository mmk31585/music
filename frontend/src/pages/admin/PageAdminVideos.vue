<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Management"
      title="Videos"
      description="Manage all music videos — approve, edit, or delete."
    >
      <template #actions>
        <Button
          label="Upload Video"
          icon="pi pi-plus"
          size="small"
          class="rounded-xl! bg-emerald-500! px-4! text-black! hover:bg-emerald-400!"
          @click="openUpload"
        />
      </template>
    </AdminSectionHeader>

    <!-- Search + stats bar -->
    <div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex flex-wrap items-center gap-3">
        <span class="text-sm text-slate-500 tabular-nums">
          {{ filteredVideos.length }} video{{ filteredVideos.length !== 1 ? 's' : '' }}
        </span>

        <div class="hidden h-4 w-px bg-white/10 sm:block" />

        <InputText
          v-model="searchQuery"
          placeholder="Search videos by title..."
          aria-label="Search videos"
          class="h-9! w-full! rounded-lg! border-white/8! bg-white/3! text-sm! text-white! placeholder:text-slate-600! sm:w-72!"
        />
      </div>

      <div class="flex items-center gap-2">
        <!-- RefreshCw -->
        <Button
          icon="pi pi-refresh"
          text
          rounded
          size="small"
          class="text-slate-400!"
          v-tooltip.top="'RefreshCw'"
          :loading="loading"
          @click="handleRefresh"
        />
      </div>
    </div>

    <!-- Loading state -->
    <div
      v-if="loading && !videos.length"
      class="overflow-hidden rounded-2xl border border-white/6 bg-white/2"
    >
      <div class="divide-y divide-white/4">
        <div v-for="i in 6" :key="i" class="flex items-center gap-4 px-5 py-4">
          <div class="h-12 w-8 animate-pulse rounded-lg bg-white/6" />
          <div class="flex-1 space-y-2">
            <div class="h-4 w-40 animate-pulse rounded bg-white/6" />
            <div class="h-3 w-24 animate-pulse rounded bg-white/4" />
          </div>
          <div class="h-4 w-12 animate-pulse rounded bg-white/4" />
          <div class="h-4 w-16 animate-pulse rounded bg-white/4" />
          <div class="h-4 w-16 animate-pulse rounded bg-white/4" />
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <AdminEmptyState
      v-else-if="filteredVideos.length === 0 && !searchQuery"
      icon="pi pi-play-circle"
      title="No videos yet"
      description="Upload your first music video to get started."
    >
      <template #action>
        <Button
          label="Upload Video"
          icon="pi pi-plus"
          size="small"
          class="rounded-xl! bg-emerald-500! px-4! text-black! hover:bg-emerald-400!"
          @click="openUpload"
        />
      </template>
    </AdminEmptyState>

    <!-- No results state -->
    <AdminEmptyState
      v-else-if="filteredVideos.length === 0 && searchQuery"
      icon="pi pi-search"
      title="No videos match your search"
      description="Try a different search term."
    />

    <!-- Table -->
    <div
      v-else
      class="overflow-hidden rounded-2xl border border-white/6 bg-white/2"
    >
      <!-- Table header (desktop) -->
      <div class="hidden border-b border-white/6 px-5 py-2.5 text-[11px] font-semibold tracking-wider text-slate-500 uppercase md:grid md:grid-cols-[48px_1fr_80px_100px_100px_80px_80px_100px]">
        <span />
        <span>Title</span>
        <span>Type</span>
        <span>Status</span>
        <span>Approved</span>
        <span>Views</span>
        <span>Likes</span>
        <span class="text-right">Actions</span>
      </div>

      <!-- Table body -->
      <div class="divide-y divide-white/4">
        <div
          v-for="v in filteredVideos"
          :key="String(v.id)"
          class="grid grid-cols-1 gap-2 px-4 py-4 text-sm transition hover:bg-white/3 md:grid-cols-[48px_1fr_80px_100px_100px_80px_80px_100px] md:items-center md:px-5 md:py-3"
        >
          <!-- Thumbnail (clickable) -->
          <div class="hidden md:block">
            <button
              type="button"
              class="flex h-12 w-8 items-center justify-center overflow-hidden rounded-md bg-white/6 transition hover:ring-2 hover:ring-sky-400/50"
              :title="'Watch ' + v.title"
              @click="watchVideo(v)"
            >
              <img
                v-if="v.thumbnail_url || v.thumbnail_path"
                :src="(v.thumbnail_url || v.thumbnail_path) ?? undefined"
                :alt="v.title"
                class="h-full w-full object-cover"
                loading="lazy"
                @error="($event.target as HTMLImageElement).style.display='none'"
              />
              <PlayCircle v-else aria-hidden="true" class="text-xs text-slate-500"  />
            </button>
          </div>

          <!-- Title + metadata (mobile) -->
          <div class="md:hidden">
            <div class="flex items-start gap-3">
              <button
                type="button"
                class="flex h-16 w-11 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-white/6 transition hover:ring-2 hover:ring-sky-400/50"
                :title="'Watch ' + v.title"
                @click="watchVideo(v)"
              >
                <img
                  v-if="v.thumbnail_url || v.thumbnail_path"
                  :src="(v.thumbnail_url || v.thumbnail_path) ?? undefined"
                  :alt="v.title"
                  class="h-full w-full object-cover"
                  loading="lazy"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
<PlayCircle v-else aria-hidden="true" class="text-base text-slate-500"  />
              </button>
              <div class="min-w-0 flex-1">
                <p class="truncate font-semibold text-white">{{ v.title }}</p>
                <div class="mt-1 flex flex-wrap items-center gap-2">
                  <span class="rounded-full px-2 py-0.5 text-[10px] font-bold" :class="v.type === 'official_mv' ? 'bg-spotify/15 text-spotify' : 'bg-blue-500/15 text-blue-400'">
                    {{ v.type === 'official_mv' ? 'MV' : 'Edit' }}
                  </span>
                  <span class="rounded-full px-2 py-0.5 text-[10px] font-bold" :class="statusClass(v.status)">
                    {{ v.status || '—' }}
                  </span>
                  <span class="text-[11px] text-slate-500">{{ formatCount(v.view_count) }} views</span>
                </div>
                <div class="mt-1 flex items-center gap-3 text-[11px] text-slate-500">
                  <span :class="v.is_approved ? 'text-emerald-400' : 'text-slate-600'">
                    {{ v.is_approved ? 'Approved' : 'Pending' }}
                  </span>
                  <span :class="v.is_public ? 'text-slate-400' : 'text-slate-600'">
                    {{ v.is_public ? 'Public' : 'Private' }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- Desktop columns -->
          <div class="hidden min-w-0 md:block">
            <p class="truncate font-medium text-white">{{ v.title }}</p>
          </div>
          <div class="hidden md:block">
            <span class="rounded-full px-2 py-0.5 text-[10px] font-bold" :class="v.type === 'official_mv' ? 'bg-spotify/15 text-spotify' : 'bg-blue-500/15 text-blue-400'">
              {{ v.type === 'official_mv' ? 'MV' : 'Edit' }}
            </span>
          </div>
          <div class="hidden md:block">
            <span class="rounded-full px-2 py-0.5 text-[10px] font-bold" :class="statusClass(v.status)">
              {{ v.status || '—' }}
            </span>
          </div>
          <div class="hidden md:block">
            <!-- Approve toggle -->
            <button
              type="button"
              class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[11px] font-bold transition"
              :class="v.is_approved
                ? 'bg-emerald-500/15 text-emerald-400 hover:bg-emerald-500/25'
                : 'bg-slate-500/10 text-slate-500 hover:bg-slate-500/20'"
              :disabled="saving"
              @click="handleToggleApprove(v)"
            >
              <component :is="v.is_approved ? CheckCircle : Clock" aria-hidden="true"<i  class="text-[10px]" /> />
              {{ v.is_approved ? 'Approved' : 'Approve' }}
            </button>
          </div>
          <div class="hidden text-slate-400 tabular-nums md:block">{{ formatCount(v.view_count) }}</div>
          <div class="hidden text-slate-400 tabular-nums md:block">{{ formatCount(v.like_count) }}</div>

          <!-- Actions -->
          <div class="flex items-center justify-end gap-2">
            <Button
              icon="pi pi-eye"
              text
              rounded
              size="small"
              class="text-sky-400/60! hover:text-sky-400!"
              v-tooltip.top="'Watch'"
              @click="watchVideo(v)"
            />
            <Button
              icon="pi pi-comments"
              text
              rounded
              size="small"
              class="text-sky-400/60! hover:text-sky-400!"
              v-tooltip.top="'Comments'"
              @click="openComments(v)"
            />
            <Button
              icon="pi pi-pencil"
              text
              rounded
              size="small"
              class="text-slate-400! hover:text-white!"
              v-tooltip.top="'Edit'"
              @click="openEdit(v)"
            />
            <Button
              icon="pi pi-trash"
              text
              rounded
              size="small"
              class="text-red-400/60! hover:text-red-400!"
              v-tooltip.top="'Delete'"
              @click="confirmDelete(v)"
            />
            <!-- Mobile approve action -->
            <Button
              v-if="!v.is_approved"
              icon="pi pi-check"
              text
              rounded
              size="small"
              class="text-emerald-400/60! hover:text-emerald-400! md:hidden!"
              v-tooltip.top="'Approve'"
              :loading="saving"
              @click="handleToggleApprove(v)"
            />
          </div>
        </div>
      </div>

      <!-- Load more -->
      <div v-if="hasMore" class="flex justify-center border-t border-white/6 px-5 py-4">
        <Button
          label="Load More"
          icon="pi pi-chevron-down"
          text
          size="small"
          class="text-slate-400! hover:text-white!"
          :loading="loading"
          @click="loadMore"
        />
      </div>
    </div>

    <!-- ── Edit Dialog ── -->
    <Dialog
      v-model:visible="editDialogVisible"
      :header="'Edit Video'"
      :modal="true"
      :dismissable-mask="true"
      class="w-full! max-w-lg!"
    >
      <div v-if="editingVideo" class="space-y-4">
        <div class="flex items-center gap-3">
          <div class="flex h-16 w-11 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-white/6">
            <img
              v-if="editingVideo.thumbnail_url || editingVideo.thumbnail_path"
              :src="(editingVideo.thumbnail_url || editingVideo.thumbnail_path) ?? undefined"
              :alt="editingVideo.title"
              class="h-full w-full object-cover"
              loading="lazy"
            />
            <PlayCircle v-else aria-hidden="true" class="text-base text-slate-500"  />
          </div>
          <div class="min-w-0">
            <p class="truncate font-semibold text-white">{{ editingVideo.title }}</p>
            <p class="text-xs text-slate-500">{{ editingVideo.type === 'official_mv' ? 'Official MV' : 'Fan Edit' }}</p>
          </div>
        </div>

        <div class="space-y-3">
          <div>
            <label class="mb-1 block text-xs font-semibold text-slate-400">Title</label>
            <InputText
              v-model="editForm.title"
              class="w-full! rounded-lg! border-white/8! bg-white/3! text-sm! text-white!"
              placeholder="Video title"
            />
          </div>

          <div>
            <label class="mb-1 block text-xs font-semibold text-slate-400">Description</label>
            <Textarea
              v-model="editForm.description"
              class="w-full! rounded-lg! border-white/8! bg-white/3! text-sm! text-white!"
              rows="3"
              placeholder="Video description"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="mb-1 block text-xs font-semibold text-slate-400">Type</label>
              <Select
                v-model="editForm.type"
                :options="typeOptions"
                option-label="label"
                option-value="value"
                class="w-full!"
              />
            </div>

            <div>
              <label class="mb-1 block text-xs font-semibold text-slate-400">Status</label>
              <Select
                v-model="editForm.status"
                :options="statusOptions"
                option-label="label"
                option-value="value"
                class="w-full!"
              />
            </div>
          </div>

          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2">
              <ToggleSwitch v-model="editForm.is_approved" input-id="edit-approved" />
              <label for="edit-approved" class="text-sm text-slate-300">Approved</label>
            </div>
            <div class="flex items-center gap-2">
              <ToggleSwitch v-model="editForm.is_public" input-id="edit-public" />
              <label for="edit-public" class="text-sm text-slate-300">Public</label>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            label="Cancel"
            text
            size="small"
            class="text-slate-400!"
            @click="editDialogVisible = false"
          />
          <Button
            label="Save"
            icon="pi pi-check"
            size="small"
            class="rounded-xl! bg-emerald-500! px-4! text-black! hover:bg-emerald-400!"
            :loading="saving"
            @click="handleSaveEdit"
          />
        </div>
      </template>
    </Dialog>

    <!-- ── Delete Confirm Dialog ── -->
    <AdminDeleteConfirm
      v-model:visible="deleteDialogVisible"
      title="Delete Video"
      :description="deleteTarget ? `Are you sure you want to delete «${deleteTarget.title}»? This cannot be undone.` : ''"
      :loading="deleting"
      @confirm="handleDelete"
      @cancel="deleteDialogVisible = false"
    />

    <!-- ── Comments Dialog ── -->
    <Dialog
      v-model:visible="commentsDialogVisible"
      :header="commentsVideo ? `Comments — ${commentsVideo.title}` : 'Comments'"
      :modal="true"
      :dismissable-mask="true"
      class="w-full! max-w-xl!"
    >
      <div v-if="loadingComments" class="space-y-4 py-4">
        <div v-for="i in 4" :key="i" class="flex animate-pulse gap-3">
          <div class="h-8 w-8 shrink-0 rounded-full bg-white/6" />
          <div class="flex-1 space-y-2">
            <div class="h-3 w-24 rounded bg-white/6" />
            <div class="h-4 w-3/4 rounded bg-white/4" />
          </div>
        </div>
      </div>
      <div v-else-if="commentItems.length === 0" class="flex flex-col items-center py-12 text-center">
        <MessageCircle aria-hidden="true" class="mb-3 block text-3xl text-slate-600" />
        <p class="text-sm text-slate-500">No comments on this video.</p>
      </div>
      <div v-else class="max-h-[60vh] space-y-3 overflow-y-auto">
        <div
          v-for="cm in commentItems"
          :key="cm.id"
          class="rounded-xl border border-white/6 bg-white/2 p-4 transition hover:bg-white/4"
        >
          <!-- Header -->
          <div class="flex items-start justify-between gap-2">
            <div class="flex items-center gap-2">
              <div class="flex h-7 w-7 items-center justify-center rounded-full bg-white/6 text-[10px] font-bold text-slate-500">
                {{ cm.author_name.charAt(0).toUpperCase() }}
              </div>
              <div>
                <span class="text-sm font-semibold text-white/80">{{ cm.author_name }}</span>
                <span class="ml-2 text-[10px] text-slate-500">{{ timeAgo(cm.created_at) }}</span>
              </div>
            </div>
            <div class="flex shrink-0 gap-1">
              <Button
                icon="pi pi-pencil"
                text
                rounded
                size="small"
                class="text-slate-500! hover:text-white! w-7! h-7!"
                v-tooltip.top="'Edit'"
                @click="startEditComment(cm)"
              />
              <Button
                icon="pi pi-trash"
                text
                rounded
                size="small"
                class="text-red-400/50! hover:text-red-400! w-7! h-7!"
                v-tooltip.top="'Delete'"
                :loading="deletingCommentId === cm.id"
                @click="handleDeleteComment(cm.id)"
              />
            </div>
          </div>
          <!-- Content / edit inline -->
          <div v-if="editingCommentId === cm.id" class="mt-3">
            <Textarea
              v-model="editingCommentContent"
              class="w-full! rounded-lg! border-white/8! bg-white/3! text-sm! text-white!"
              rows="2"
              maxlength="1000"
            />
            <div class="mt-2 flex justify-end gap-2">
              <Button
                label="Cancel"
                text
                size="small"
                class="text-slate-400! text-xs!"
                @click="editingCommentId = null"
              />
              <Button
                label="Save"
                icon="pi pi-check"
                size="small"
                class="rounded-lg! bg-emerald-500! px-3! text-black! hover:bg-emerald-400! text-xs!"
                :loading="savingComment"
                :disabled="!editingCommentContent.trim()"
                @click="handleSaveComment(cm.id)"
              />
            </div>
          </div>
          <p v-else class="mt-1.5 text-sm leading-relaxed text-white/60 break-words">{{ cm.content }}</p>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { CheckCircle, Clock, MessageCircle, PlayCircle } from 'lucide-vue-next'
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useAdminVideos } from '@/composables/admin'
import { useVideoApi } from '@/services/api/video'
import type { VideoItem } from '@/services/api/video/types'
import { AdminSectionHeader, AdminEmptyState, AdminDeleteConfirm } from '@/components/admin'
import { formatCount } from '@/utils/number'

const router = useRouter()
const toast = useToast()
const videoApi = useVideoApi()
const {
  videos,
  loading,
  saving,
  deleting,
  hasMore,
  fetchVideos,
  loadMore,
  updateVideo,
  deleteVideo,
  approveVideo,
} = useAdminVideos()

const searchQuery = ref('')

const filteredVideos = computed(() => {
  if (!searchQuery.value) return videos.value
  const q = searchQuery.value.toLowerCase()
  return videos.value.filter((v) => v.title.toLowerCase().includes(q))
})

// ── Edit dialog state ──
const editDialogVisible = ref(false)
const editingVideo = ref<VideoItem | null>(null)
const editForm = reactive({
  title: '',
  description: '',
  type: 'official_mv' as string,
  status: 'ready' as string,
  is_approved: false,
  is_public: false,
})

const typeOptions = [
  { label: 'Official MV', value: 'official_mv' },
  { label: 'Fan Edit', value: 'user_edit' },
]

const statusOptions = [
  { label: 'Ready', value: 'ready' },
  { label: 'Processing', value: 'processing' },
  { label: 'Failed', value: 'failed' },
]

function openEdit(v: VideoItem) {
  editingVideo.value = v
  editForm.title = v.title
  editForm.description = v.description || ''
  editForm.type = v.type
  editForm.status = v.status || 'ready'
  editForm.is_approved = v.is_approved ?? false
  editForm.is_public = v.is_public ?? false
  editDialogVisible.value = true
}

async function handleSaveEdit() {
  if (!editingVideo.value) return
  try {
    await updateVideo(String(editingVideo.value.id), {
      title: editForm.title,
      description: editForm.description,
      type: editForm.type,
      status: editForm.status,
      is_approved: editForm.is_approved,
      is_public: editForm.is_public,
    })
    editDialogVisible.value = false
    editingVideo.value = null
  } catch {
    // Toast handled by API layer
  }
}

// ── Delete state ──
const deleteDialogVisible = ref(false)
const deleteTarget = ref<VideoItem | null>(null)

function confirmDelete(v: VideoItem) {
  deleteTarget.value = v
  deleteDialogVisible.value = true
}

async function handleDelete() {
  if (!deleteTarget.value) return
  try {
    await deleteVideo(String(deleteTarget.value.id))
    deleteDialogVisible.value = false
    deleteTarget.value = null
  } catch {
    // Toast handled by API layer
  }
}

// ── Approve toggle ──
async function handleToggleApprove(v: VideoItem) {
  try {
    if (v.is_approved) {
      // Unapprove by updating
      await updateVideo(String(v.id), { is_approved: false })
    } else {
      await approveVideo(String(v.id))
    }
  } catch {
    // Toast handled by API layer
  }
}

// ── Actions ──
function watchVideo(v: VideoItem) {
  window.open(`/music-video/${v.id}`, '_blank')
}

function openUpload() {
  router.push('/admin/video-upload')
}

function handleRefresh() {
  fetchVideos(true)
}

// ── Comments dialog ──
interface AdminCommentItem {
  id: string
  video_id: string
  user_id: string
  author_name: string
  content: string
  created_at: string
  updated_at: string
}

const commentsDialogVisible = ref(false)
const commentsVideo = ref<VideoItem | null>(null)
const commentItems = ref<AdminCommentItem[]>([])
const loadingComments = ref(false)
const editingCommentId = ref<string | null>(null)
const editingCommentContent = ref('')
const savingComment = ref(false)
const deletingCommentId = ref<string | null>(null)

async function openComments(v: VideoItem) {
  commentsVideo.value = v
  commentsDialogVisible.value = true
  editingCommentId.value = null
  loadingComments.value = true
  try {
    const res = await videoApi.adminGetVideoComments(String(v.id))
    commentItems.value = res?.items || []
  } catch {
    commentItems.value = []
  } finally {
    loadingComments.value = false
  }
}

function startEditComment(cm: AdminCommentItem) {
  editingCommentId.value = cm.id
  editingCommentContent.value = cm.content
}

async function handleSaveComment(commentId: string) {
  if (editingCommentContent.value.trim!()) return
  savingComment.value = true
  try {
    await videoApi.adminUpdateComment(commentId, { content: editingCommentContent.value.trim() })
    const idx = commentItems.value.findIndex((c) => c.id === commentId)
    if (idx !== -1) commentItems.value[idx]!.content = editingCommentContent.value.trim()
    editingCommentId.value = null
    toast.add({ severity: 'success', summary: 'Comment updated', life: 3000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to update comment', life: 3000 })
  } finally {
    savingComment.value = false
  }
}

async function handleDeleteComment(commentId: string) {
  deletingCommentId.value = commentId
  try {
    await videoApi.adminDeleteComment(commentId)
    commentItems.value = commentItems.value.filter((c) => c.id !== commentId)
    toast.add({ severity: 'success', summary: 'Comment deleted', life: 3000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to delete comment', life: 3000 })
  } finally {
    deletingCommentId.value = null
  }
}

function timeAgo(dateStr: string): string {
  const now = Date.now()
  const then = new Date(dateStr).getTime()
  const diff = Math.floor((now - then) / 1000)
  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return `${Math.floor(diff / 86400)}d ago`
}

// ── Helpers ──
function statusClass(status?: string): string {
  switch (status) {
    case 'ready': return 'bg-emerald-500/15 text-emerald-400'
    case 'processing': return 'bg-yellow-500/15 text-yellow-400'
    case 'failed': return 'bg-red-500/15 text-red-400'
    default: return 'bg-slate-500/10 text-slate-500'
  }
}

// Initial fetch
fetchVideos(true)
</script>
