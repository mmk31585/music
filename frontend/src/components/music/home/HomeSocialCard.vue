<template>
  <div class="w-48 shrink-0 rounded-2xl border border-white/[0.06] bg-white/[0.04] p-4 backdrop-blur">
    <div class="flex items-center gap-2">
      <div
        class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-[#1db954] to-[#60a5fa] text-sm font-bold text-white"
      >
        {{ initials }}
      </div>
      <div class="min-w-0 flex-1">
        <p class="truncate text-sm font-semibold text-white">
          {{ username }}
        </p>
        <p class="text-[11px] text-white/60">
          {{ timeAgo }}
        </p>
      </div>
    </div>
    <div class="mt-3 flex justify-center">
      <div class="h-20 w-20 overflow-hidden rounded-xl bg-white/10 ring-1 ring-white/10">
        <img
          v-if="coverUrl"
          :src="coverUrl"
          alt=""
          class="h-full w-full object-cover"
          loading="lazy"
        />
        <div v-else class="flex h-full items-center justify-center">
          <i aria-hidden="true" class="pi pi-music text-xl text-white/30" />
        </div>
      </div>
    </div>
    <div class="mt-2 text-center">
      <p class="truncate text-sm font-semibold text-white">
        {{ trackTitle }}
      </p>
      <p class="truncate text-xs text-white/60">
        {{ artistName }}
      </p>
    </div>
    <button
      type="button"
      class="mt-3 w-full rounded-full border border-white/10 bg-white/[0.04] px-3 py-1.5 text-xs font-medium text-white/80 transition hover:bg-white/10 focus-visible:ring-2 focus-visible:ring-[#1db954] focus-visible:outline-none"
      @click="$emit('listen-together', activity)"
    >
      پخش با هم
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface SocialActivity {
  username?: string
  timeAgo?: string
  cover_url?: string | null
  track_title?: string | null
  artist_name?: string | null
  user?: {
    username?: string
  }
  track?: {
    cover_url?: string | null
    title?: string | null
    artist_name?: string | null
  }
}

const props = defineProps<{
  activity: SocialActivity
}>()

defineEmits<{
  'listen-together': [activity: SocialActivity]
}>()

const initials = computed(() => {
  const name = props.activity.username || props.activity.user?.username || '?'
  return name.slice(0, 2).toUpperCase()
})

const username = computed(() => props.activity.username || props.activity.user?.username || 'کاربر')
const timeAgo = computed(() => props.activity.timeAgo || 'الان')
const coverUrl = computed(() => props.activity.cover_url || props.activity.track?.cover_url || null)
const trackTitle = computed(() => props.activity.track_title || props.activity.track?.title || 'آهنگ')
const artistName = computed(() => props.activity.artist_name || props.activity.track?.artist_name || '')
</script>
