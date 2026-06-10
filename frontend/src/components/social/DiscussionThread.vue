<template>
  <div class="space-y-4">
    <div v-if="!discussions.length" class="py-12 text-center text-sm text-white/30">
      No comments yet. Be the first to share your thoughts!
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="discussion in discussions"
        :key="discussion.id"
        class="rounded-xl border border-white/5 bg-white/[0.02] p-4 transition hover:border-white/10"
      >
        <div class="flex items-start gap-3">
          <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-[#1db954]/20 text-xs font-bold text-[#1db954]">
            {{ displayName(discussion.user_id)?.charAt(0).toUpperCase() || '?' }}
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2 text-xs">
              <span class="font-semibold text-white/70">{{ displayName(discussion.user_id) || 'Unknown' }}</span>
              <span class="text-white/20">{{ formatTime(discussion.created_at) }}</span>
            </div>
            <p class="mt-1 text-sm text-white/80">{{ discussion.content }}</p>
            <div class="mt-2 flex items-center gap-3">
              <button
                v-if="!showReplyInputs[discussion.id]"
                class="text-[11px] text-blue-400/50 transition hover:text-blue-400"
                @click="toggleReply(discussion.id)"
              >
                Reply
              </button>
            </div>

            <!-- Reply input -->
            <div v-if="showReplyInputs[discussion.id]" class="mt-3 flex gap-2">
              <input
                v-model="replyTexts[discussion.id]"
                type="text"
                placeholder="Write a reply..."
                class="min-w-0 flex-1 rounded-lg border border-white/10 bg-white/5 px-3 py-2 text-xs text-white placeholder-white/20 outline-none focus:border-white/20"
                @keydown.enter="submitReply(discussion.id)"
              />
              <button
                class="rounded-lg bg-blue-500/10 px-3 py-2 text-xs font-medium text-blue-400 transition hover:bg-blue-500/20"
                @click="submitReply(discussion.id)"
              >
                Send
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- New comment input -->
    <div class="flex gap-3">
      <input
        v-model="newComment"
        type="text"
        placeholder="Write a comment..."
        class="min-w-0 flex-1 rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 outline-none transition focus:border-white/20"
        @keydown.enter="submitComment"
      />
      <button
        class="rounded-xl bg-[#1db954]/10 px-5 py-3 text-sm font-semibold text-[#1db954] transition hover:bg-[#1db954]/20 disabled:opacity-30"
        :disabled="!newComment.trim()"
        @click="submitComment"
      >
        Post
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { Discussion } from '@/services/api/social'

const props = defineProps<{
  discussions: Discussion[]
  userNames?: Record<string, string>
}>()

const emit = defineEmits<{
  create: [content: string, parentId?: string]
}>()

const newComment = ref('')
const showReplyInputs = reactive<Record<string, boolean>>({})
const replyTexts = reactive<Record<string, string>>({})

function displayName(userId: string): string {
  return props.userNames?.[userId] || userId?.slice(0, 8) || 'Unknown'
}

function submitComment() {
  if (!newComment.value.trim()) return
  emit('create', newComment.value.trim())
  newComment.value = ''
}

function toggleReply(id: string) {
  showReplyInputs[id] = !showReplyInputs[id]
  if (!showReplyInputs[id]) {
    replyTexts[id] = ''
  }
}

function submitReply(parentId: string) {
  const text = replyTexts[parentId]?.trim()
  if (!text) return
  emit('create', text, parentId)
  replyTexts[parentId] = ''
  showReplyInputs[parentId] = false
}

function formatTime(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}d ago`
  return d.toLocaleDateString()
}
</script>
