<template>
  <div
    class="flex items-start gap-3 rounded-xl bg-white/3 p-4 ring-1 ring-white/6 transition hover:bg-white/5 focus-within:ring-spotify focus-within:ring-2"
    role="article"
    :aria-label="`Discussion by ${discussion.userName || discussion.user_id}`"
  >
    <!-- Avatar -->
    <RouterLink
      :to="`/profile/${discussion.user_id}`"
      class="shrink-0 h-8 w-8 overflow-hidden rounded-full ring-1 ring-white/10 focus-visible:outline-2 focus-visible:outline-[#1db954]"
    >
      <img
        v-if="discussion.avatarUrl"
        :src="discussion.avatarUrl"
        :alt="discussion.userName || 'User'"
        class="h-full w-full object-cover"
        loading="lazy"
      />
      <div
        v-else
        class="flex h-full w-full items-center justify-center bg-white/10 text-[10px] font-bold text-white"
      >
        {{ String(discussion.userName || discussion.user_id || '?')[0].toUpperCase() }}
      </div>
    </RouterLink>

    <!-- Content -->
    <div class="min-w-0 flex-1">
      <div class="flex items-center gap-2 text-xs">
        <span class="font-semibold text-white/70">{{ discussion.userName || discussion.user_id?.slice(0, 8) || 'Unknown' }}</span>
        <span class="text-white/20">{{ formattedTime }}</span>
      </div>
      <p class="mt-1 text-sm text-white/80 leading-relaxed">{{ discussion.content }}</p>

      <!-- Footer: reply count + reply button -->
      <div class="mt-2 flex items-center gap-3">
        <button
          aria-label="Reply to discussion"
          class="inline-flex items-center gap-1 text-[11px] text-blue-400/50 transition hover:text-blue-400 focus-visible:outline-2 focus-visible:outline-[#1db954]"
          @click="$emit('reply', discussion.id)"
        >
          <i aria-hidden="true" class="pi pi-reply text-[10px]" />
          Reply
        </button>
        <span
          v-if="discussion.reply_count !== undefined"
          class="inline-flex items-center gap-1 text-[11px] text-white/20"
        >
          <i aria-hidden="true" class="pi pi-comments text-[10px]" />
          {{ discussion.reply_count }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface DiscussionCardData {
  id: string
  user_id: string
  userName?: string
  avatarUrl?: string | null
  content: string
  reply_count?: number
  created_at: string
}

const props = defineProps<{
  discussion: DiscussionCardData
}>()

defineEmits<{
  reply: [id: string]
}>()

const formattedTime = computed(() => {
  if (props.discussion.created_at!) return ''
  const d = new Date(props.discussion.created_at)
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
})
</script>
