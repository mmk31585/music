<template>
  <div class="space-y-4">
    <div v-if="!discussions.length" class="py-8 text-center text-sm text-white/30" role="status">
      هنوز بحثی ثبت نشده. اولین نفر باش!
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="d in discussions"
        :key="d.id"
        class="rounded-xl border border-white/5 bg-white/2 p-4 transition hover:border-white/10"
      >
        <div class="flex items-start gap-3">
          <div
            class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-spotify/20 text-xs font-bold text-spotify"
          >
            {{ displayName(d.author_id)?.charAt(0).toUpperCase() || '?' }}
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2 text-xs">
              <span class="font-semibold text-white/70">{{ displayName(d.author_id) || 'Unknown' }}</span>
              <span class="text-white/20">{{ formatTime(d.created_at) }}</span>
            </div>
            <h4 class="mt-1 text-sm font-bold text-white">{{ d.title }}</h4>
            <p class="mt-1 text-sm text-white/70 leading-relaxed">{{ d.body }}</p>
            <div class="mt-3 flex items-center gap-4 text-xs">
              <button
                class="inline-flex items-center gap-1 text-white/40 transition hover:text-red-400"
                @click="toggleLike(d.id)"
                aria-label="Like"
              >
                <Heart aria-hidden="true" class="text-sm" :class="likedDiscussions[d.id] ? 'text-red-400' : ''"  />
                {{ likeCounts[d.id] || 0 }}
              </button>
              <button
                class="inline-flex items-center gap-1 text-white/40 transition hover:text-blue-400"
                @click="toggleReplies(d.id)"
              >
                <MessageCircle aria-hidden="true" class="text-sm"  />
                {{ d.reply_count }} پاسخ
              </button>
              <button
                v-if="canDelete(d)"
                class="mr-auto text-white/20 transition hover:text-red-400"
                @click="emit('delete', d.id)"
                aria-label="Delete"
              >
                <Trash2 aria-hidden="true" class="text-xs"  />
              </button>
            </div>

            <!-- Replies -->
            <div v-if="openReplies[d.id]" class="mt-4 space-y-3 border-t border-white/5 pt-4">
              <div v-if="repliesLoading[d.id]" class="flex items-center justify-center py-3">
                <span class="inline-block h-4 w-4 animate-spin rounded-full border-2 border-white/20 border-t-accent" />
              </div>
              <div
                v-for="reply in loadedReplies[d.id] || []"
                v-else
                :key="reply.id"
                class="flex items-start gap-2 pr-6"
              >
                <div
                  class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-white/10 text-[10px] font-bold text-white/50"
                >
                  {{ displayName(reply.author_id)?.charAt(0).toUpperCase() || '?' }}
                </div>
                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-1.5 text-[11px]">
                    <span class="font-semibold text-white/50">{{ displayName(reply.author_id) || 'Unknown' }}</span>
                    <span class="text-white/20">{{ formatTime(reply.created_at) }}</span>
                  </div>
                  <p class="mt-0.5 text-xs text-white/60">{{ reply.body }}</p>
                </div>
              </div>

              <div v-if="replyError[d.id]" class="rounded-lg bg-red-500/10 px-3 py-2 text-[11px] text-red-400">
                {{ replyError[d.id] }}
              </div>

              <div class="flex gap-2 pr-6">
                <input
                  v-model="replyInputs[d.id]"
                  type="text"
                  placeholder="پاسخ خودت رو بنویس..."
                  aria-label="پاسخ"
                  class="min-w-0 flex-1 rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-xs text-white placeholder-white/20 outline-hidden focus:border-white/20"
                  dir="rtl"
                  @keydown.enter="submitReply(d.id)"
                />
                <button
                  class="inline-flex items-center gap-1 rounded-lg bg-spotify/10 px-3 py-2 text-xs font-medium text-spotify transition hover:bg-spotify/20 disabled:opacity-40"
                  :disabled="!replyInputs[d.id]?.trim() || repliesLoading[d.id]"
                  @click="submitReply(d.id)"
                >
                  <span v-if="repliesLoading[d.id]" class="inline-block h-3 w-3 animate-spin rounded-full border-2 border-spotify border-t-transparent" />
                  ارسال
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Heart, MessageCircle, Trash2 } from 'lucide-vue-next'
import { reactive, watch } from 'vue'
import { useSocialApi } from '@/services/api/social'
import { useReactionsApi } from '@/services/api/reactions'
import type { ClubDiscussion, ClubDiscussionReply } from '@/services/api/social'

const props = defineProps<{
  discussions: ClubDiscussion[]
  currentUserId?: string
  currentUserRole?: string
  userNames?: Record<string, string>
}>()

const emit = defineEmits<{
  delete: [id: string]
}>()

const socialApi = useSocialApi()
const reactionsApi = useReactionsApi()

const loadedReplies = reactive<Record<string, ClubDiscussionReply[]>>({})
const repliesLoading = reactive<Record<string, boolean>>({})
const replyError = reactive<Record<string, string>>({})
const openReplies = reactive<Record<string, boolean>>({})
const replyInputs = reactive<Record<string, string>>({})
const likedDiscussions = reactive<Record<string, boolean>>({})
const likeCounts = reactive<Record<string, number>>({})

function displayName(userId: string): string {
  return props.userNames?.[userId] || userId?.slice(0, 8) || 'Unknown'
}

function canDelete(d: ClubDiscussion): boolean {
  if (props.currentUserId === undefined) return false
  if (d.author_id === props.currentUserId) return true
  if (props.currentUserRole === 'admin' || props.currentUserRole === 'moderator') return true
  return false
}

async function toggleReplies(discussionId: string) {
  if (openReplies[discussionId]) {
    openReplies[discussionId] = false
    return
  }
  openReplies[discussionId] = true
  if (!loadedReplies[discussionId]) {
    await loadReplies(discussionId)
  }
}

async function loadReplies(discussionId: string) {
  repliesLoading[discussionId] = true
  replyError[discussionId] = ''
  try {
    const data = await socialApi.getClubDiscussionReplies(discussionId)
    loadedReplies[discussionId] = Array.isArray(data) ? data : []
  } catch (err: any) {
    const msg = err?.response?.data?.message || err.message || 'Failed to load replies'
    replyError[discussionId] = msg
    loadedReplies[discussionId] = []
  } finally {
    repliesLoading[discussionId] = false
  }
}

async function submitReply(discussionId: string) {
  const text = replyInputs[discussionId]?.trim()
  if (!text || repliesLoading[discussionId]) return
  repliesLoading[discussionId] = true
  replyError[discussionId] = ''
  try {
    await socialApi.createClubDiscussionReply(discussionId, { body: text })
    replyInputs[discussionId] = ''
    await loadReplies(discussionId)
  } catch (err: any) {
    const msg = err?.response?.data?.message || err.message || 'Failed to submit reply'
    replyError[discussionId] = msg
  } finally {
    repliesLoading[discussionId] = false
  }
}

async function toggleLike(discussionId: string) {
  if (likedDiscussions[discussionId]) {
    try {
      await reactionsApi.removeReaction('club_discussion', discussionId)
      likedDiscussions[discussionId] = false
      likeCounts[discussionId] = Math.max(0, (likeCounts[discussionId] || 0) - 1)
    } catch (err) {
      console.error('Failed to remove reaction:', err)
    }
  } else {
    try {
      await reactionsApi.react({
        target_type: 'club_discussion',
        target_id: discussionId,
        type: 'like',
      })
      likedDiscussions[discussionId] = true
      likeCounts[discussionId] = (likeCounts[discussionId] || 0) + 1
    } catch (err) {
      console.error('Failed to add reaction:', err)
    }
  }
}

async function loadLikeState(discussionId: string) {
  try {
    const [counts, userReaction] = await Promise.allSettled([
      reactionsApi.getCounts('club_discussion', discussionId),
      reactionsApi.getUserReaction('club_discussion', discussionId),
    ])
    if (counts.status === 'fulfilled' && counts.value) {
      likeCounts[discussionId] = counts.value.like
    }
    if (userReaction.status === 'fulfilled' && userReaction.value?.reaction) {
      likedDiscussions[discussionId] = true
    }
  } catch (err) {
    console.error('Failed to load like state:', err)
  }
}

watch(() => props.discussions, (list) => {
  list.forEach(d => { if (likedDiscussions[d.id] === undefined) loadLikeState(d.id) })
}, { immediate: true })

function formatTime(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'همین حالا'
  if (mins < 60) return `${mins} دقیقه پیش`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours} ساعت پیش`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days} روز پیش`
  return d.toLocaleDateString('fa-IR')
}
</script>
