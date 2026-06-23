<template>
  <div class="space-y-4">
    <!-- Category pill tabs -->
    <div
      class="flex gap-2 overflow-x-auto pb-2 scrollbar-none"
      role="tablist"
      aria-label="Discussion categories"
    >
      <button
        v-for="cat in discussionCategories"
        :key="cat.key"
        role="tab"
        :aria-selected="activeTab === cat.key"
        class="shrink-0 rounded-full px-4 py-2 text-xs font-medium transition focus-visible:outline-2 focus-visible:outline-[#1db954]"
        :class="activeTab === cat.key
          ? 'bg-white/15 text-white shadow-lg'
          : 'bg-white/4 text-white/40 hover:bg-white/8 hover:text-white/60'"
        @click="activeTab = cat.key"
      >
        <span class="flex items-center gap-2 whitespace-nowrap">
          <i aria-hidden="true" :class="cat.icon" class="text-[10px]" />
          {{ cat.label }}
        </span>
      </button>
    </div>

    <!-- Discussion feed -->
    <div v-if="!showStartInput" class="flex justify-end">
      <button
        aria-label="Start a new discussion"
        class="inline-flex items-center gap-1.5 rounded-full bg-spotify/10 px-4 py-2 text-xs font-semibold text-spotify transition hover:bg-spotify/20 focus-visible:outline-2 focus-visible:outline-[#1db954]"
        @click="showStartInput = true"
      >
        <i aria-hidden="true" class="pi pi-plus text-[10px]" />
        Start Discussion
      </button>
    </div>

    <!-- Start discussion input -->
    <div v-if="showStartInput" class="flex gap-2">
      <input
        v-model="newDiscussionContent"
        type="text"
        placeholder="Share your thoughts..."
        aria-label="New discussion content"
        class="min-w-0 flex-1 rounded-lg border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white placeholder-white/20 outline-hidden transition focus:border-white/20"
        @keydown.enter="submitDiscussion"
      />
      <button
        class="rounded-lg bg-spotify/10 px-4 py-2.5 text-sm font-semibold text-spotify transition hover:bg-spotify/20 disabled:opacity-30 inline-flex items-center gap-1.5"
        :disabled="newDiscussionContent.trim!() || isPosting"
        @click="submitDiscussion"
      >
        <span v-if="isPosting" class="inline-block h-3.5 w-3.5 animate-spin rounded-full border-2 border-spotify border-t-transparent" />
        {{ isPosting ? 'Posting...' : 'Post' }}
      </button>
      <button
        class="rounded-lg bg-white/5 px-3 py-2.5 text-xs text-white/40 transition hover:bg-white/10"
        @click="showStartInput = false; newDiscussionContent = ''"
        aria-label="Cancel"
      >
        Cancel
      </button>
    </div>

    <!-- Discussion cards -->
    <div v-if="filteredDiscussions.length" class="space-y-2" role="feed" aria-label="Discussions">
      <DiscussionCard
        v-for="discussion in filteredDiscussions"
        :key="discussion.id"
        :discussion="discussion"
        @reply="handleReply"
      />
    </div>

    <!-- Empty state -->
    <div
      v-else
      class="flex flex-col items-center gap-3 py-12 text-center"
    >
      <i aria-hidden="true" class="pi pi-comments text-2xl text-white/10" />
      <p class="text-sm text-white/30">No discussions yet in this category</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import DiscussionCard from './DiscussionCard.vue'
import type { Discussion } from '@/services/api/social'
import type { DiscussionCardData } from './DiscussionCard.vue'

export type { DiscussionCardData }

const props = defineProps<{
  discussions: Discussion[]
  userNames?: Record<string, string>
  userAvatars?: Record<string, string | null>
  isPosting?: boolean
}>()

const emit = defineEmits<{
  create: [content: string, parentId?: string]
  reply: [discussionId: string]
}>()

const activeTab = ref('all')
const showStartInput = ref(false)
const newDiscussionContent = ref('')

const discussionCategories = [
  { key: 'all', label: 'All', icon: 'pi pi-globe' },
  { key: 'track', label: 'Tracks', icon: 'pi pi-music' },
  { key: 'album', label: 'Albums', icon: 'pi pi-book' },
  { key: 'artist', label: 'Artists', icon: 'pi pi-user' },
  { key: 'playlist', label: 'Playlists', icon: 'pi pi-list' },
]

const filteredDiscussions = computed<DiscussionCardData[]>(() => {
  const items = props.discussions || []
  const filtered = activeTab.value === 'all'
    ? items
    : items.filter((d) => d.target_type === activeTab.value)

  return filtered.map((d) => ({
    id: d.id,
    user_id: d.user_id,
    userName: props.userNames?.[d.user_id],
    avatarUrl: props.userAvatars?.[d.user_id] || null,
    content: d.content,
    reply_count: d.parent_id ? undefined : (items.filter((r) => r.parent_id === d.id).length),
    created_at: d.created_at,
  }))
})

function handleReply(discussionId: string) {
  emit('reply', discussionId)
}

function submitDiscussion() {
  if (newDiscussionContent.value.trim!()) return
  emit('create', newDiscussionContent.value.trim())
  newDiscussionContent.value = ''
  showStartInput.value = false
}
</script>
