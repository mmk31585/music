<template>
  <div
    class="flex items-start gap-3 rounded-xl p-3 transition hover:bg-white/4 focus-within:ring-2 focus-within:ring-spotify"
    role="article"
    :aria-label="`Activity: ${item.userName} ${item.action}`"
  >
    <!-- Avatar with link to profile -->
    <RouterLink
      :to="`/profile/${item.userId}`"
      class="shrink-0 h-8 w-8 overflow-hidden rounded-full ring-1 ring-white/10 focus-visible:outline-2 focus-visible:outline-[#1db954]"
    >
      <img
        v-if="item.avatarUrl"
        :src="item.avatarUrl"
        :alt="item.userName"
        class="h-full w-full object-cover"
        loading="lazy"
      />
      <div
        v-else
        class="flex h-full w-full items-center justify-center bg-white/10 text-[10px] font-bold text-white"
      >
        {{ (item.userName || '?')[0] }}
      </div>
    </RouterLink>

    <!-- Content -->
    <div class="min-w-0 flex-1">
      <p class="text-xs text-white/70 leading-relaxed">
        <RouterLink
          :to="`/profile/${item.userId}`"
          class="font-semibold text-white hover:underline focus-visible:outline-2 focus-visible:outline-[#1db954]"
        >
          {{ item.userName }}
        </RouterLink>
        {{ ' ' + item.action + ' ' }}
        <template v-if="item.targetName">
          <RouterLink
            v-if="item.targetUrl"
            :to="item.targetUrl"
            class="font-medium text-spotify hover:underline focus-visible:outline-2 focus-visible:outline-[#1db954]"
          >
            {{ item.targetName }}
          </RouterLink>
          <span v-else class="font-medium text-white/80">{{ item.targetName }}</span>
        </template>
      </p>
      <p class="mt-0.5 text-[10px] text-white/30">{{ displayTimeAgo }}</p>
    </div>

    <!-- Contextual action button -->
    <button
      v-if="item.actionType === 'party'"
      aria-label="Join party"
      class="shrink-0 rounded-lg bg-spotify/10 px-3 py-1.5 text-[10px] font-semibold text-spotify transition hover:bg-spotify/20 focus-visible:outline-2 focus-visible:outline-[#1db954]"
      @click="$emit('action', item)"
    >
      Join
    </button>
    <button
      v-else-if="item.actionType === 'room'"
      aria-label="Listen live"
      class="shrink-0 rounded-lg bg-red-500/10 px-3 py-1.5 text-[10px] font-semibold text-red-400 transition hover:bg-red-500/20 focus-visible:outline-2 focus-visible:outline-red-400"
      @click="$emit('action', item)"
    >
      Listen
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface ActivityItemData {
  id: string
  userId: string
  userName: string
  avatarUrl?: string | null
  action: string
  targetName?: string | null
  targetUrl?: string | null
  actionType?: string | null
  timeAgo?: string | null
  createdAt?: string | null
}

const props = defineProps<{
  item: ActivityItemData
}>()

defineEmits<{
  action: [item: ActivityItemData]
}>()

const displayTimeAgo = computed(() => {
  if (props.item.timeAgo) return props.item.timeAgo
  if (!props.item.createdAt) return ''
  const d = new Date(props.item.createdAt)
  const now = Date.now()
  const diff = now - d.getTime()
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
