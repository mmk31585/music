<template>
  <div class="relative overflow-hidden rounded-2xl">
    <!-- Aurora gradient — driven by genreColor or green fallback -->
    <div class="pointer-events-none absolute inset-0" aria-hidden="true">
      <div
        class="absolute inset-0 bg-linear-to-br"
        :style="heroGradient"
      />
      <div
        class="absolute -top-20 -right-20 h-80 w-80 rounded-full blur-3xl opacity-20"
        :style="{ background: orbColor }"
      />
      <div
        class="absolute -bottom-16 -left-16 h-60 w-60 rounded-full blur-3xl opacity-15"
        :style="{ background: orbColor2 }"
      />
      <div class="absolute inset-0 bg-linear-to-t from-surface-base via-surface-base/20 to-transparent" />
    </div>

    <div class="relative z-10 flex flex-col gap-6 px-6 pt-16 pb-8 md:flex-row md:items-end md:gap-10 md:pt-12 md:pb-10">
      <!-- Avatar with floating ring -->
      <div class="relative shrink-0 self-center md:self-end">
        <div
          v-if="musicStatus?.playing"
          class="absolute -inset-2 motion-safe:animate-ping rounded-full border-2 border-spotify opacity-30"
        />
        <div
          class="h-36 w-36 overflow-hidden rounded-full border-4 shadow-2xl md:h-44 md:w-44"
          :class="musicStatus?.playing ? 'border-spotify' : 'border-white/10'"
        >
          <img
            v-if="avatarUrl"
            :src="avatarUrl"
            :alt="displayName"
            loading="lazy"
            class="h-full w-full object-cover"
            @error="onImgError"
          />
          <div
            v-else
            class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30 text-4xl text-white/40"
          >
            <i aria-hidden="true" class="pi pi-user" />
          </div>
        </div>
        <!-- Creator badge -->
        <div
          v-if="isCreator"
          class="absolute -bottom-1 -right-1 flex h-7 w-7 items-center justify-center rounded-full bg-spotify ring-2 ring-surface-base shadow-lg"
          aria-label="Creator"
        >
          <i aria-hidden="true" class="pi pi-check text-[11px] text-black" />
        </div>
      </div>

      <div class="flex flex-col items-center text-center md:items-start md:text-start">
        <!-- Name + handle -->
        <h1 class="text-3xl font-black text-white md:text-5xl">{{ displayName }}</h1>
        <p v-if="handle" class="mt-1 text-sm text-white/40">@{{ handle }}</p>
        <p v-if="bio" class="mt-2 max-w-md text-sm text-white/50 leading-relaxed">{{ bio }}</p>

        <!-- Stats row -->
        <div class="mt-4 flex flex-wrap items-center justify-center gap-5 text-sm md:justify-start">
          <button
            class="transition hover:text-white focus-visible:outline-2 focus-visible:outline-[#1db954]"
            aria-label="View followers"
            @click="$emit('showFollowers')"
          >
            <span class="font-bold text-white tabular-nums">{{ formatNumber(followerCount) }}</span>
            <span class="text-white/40"> followers</span>
          </button>
          <button
            class="transition hover:text-white focus-visible:outline-2 focus-visible:outline-[#1db954]"
            aria-label="View following"
            @click="$emit('showFollowing')"
          >
            <span class="font-bold text-white tabular-nums">{{ formatNumber(followingCount) }}</span>
            <span class="text-white/40"> following</span>
          </button>
          <span class="text-white/30" aria-hidden="true">·</span>
          <span class="text-white/40">
            <span class="font-medium text-white/60">Joined</span>
            {{ formatJoinDate }}
          </span>
        </div>

        <!-- NOW PLAYING — prominent, animated -->
        <div
          v-if="musicStatus?.playing && currentTrackInfo"
          class="mt-4 flex w-full max-w-md items-center gap-3 rounded-2xl bg-white/6 p-3 ring-1 ring-white/10 motion-safe:animate-fade-in-up"
        >
          <div class="h-12 w-12 shrink-0 overflow-hidden rounded-xl shadow-lg">
            <img
              v-if="currentTrackInfo.coverUrl"
              :src="currentTrackInfo.coverUrl"
              alt=""
              class="h-full w-full object-cover"
            />
            <div v-else class="flex h-full w-full items-center justify-center bg-white/5">
              <i aria-hidden="true" class="pi pi-music text-sm text-white/30" />
            </div>
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <span class="flex gap-0.5" aria-hidden="true">
                <span
                  v-for="i in 4" :key="i"
                  class="h-3 w-0.5 rounded-full bg-spotify motion-safe:animate-equalizer"
                  :style="{ animationDelay: `${i * 100}ms` }"
                />
              </span>
              <p class="truncate text-xs font-semibold text-white">{{ currentTrackInfo.title }}</p>
            </div>
            <p class="truncate text-[10px] text-white/40 ms-4">
              {{ currentTrackInfo.artist }}
            </p>
          </div>
          <button
            v-if="!isOwnProfile"
            aria-label="Listen along"
            class="shrink-0 rounded-full bg-spotify px-4 py-1.5 text-[10px] font-bold text-black transition hover:bg-spotify-hover hover:scale-105 active:scale-95 focus-visible:outline-2 focus-visible:outline-white"
            @click="$emit('listenAlong', currentTrackInfo.id)"
          >
            Listen
          </button>
        </div>

        <!-- Not playing -->
        <div
          v-else-if="musicStatus && !musicStatus.playing"
          class="mt-4 text-xs text-white/30 flex items-center gap-2"
        >
          <i aria-hidden="true" class="pi pi-pause-circle text-sm" />
          Not listening right now
        </div>

        <!-- Action buttons -->
        <div class="mt-5 flex flex-wrap items-center gap-3">
          <button
            v-if="!isOwnProfile"
            :class="[
              'inline-flex items-center gap-2 rounded-full px-6 py-3 text-sm font-bold transition',
              'focus-visible:outline-2 focus-visible:outline-white active:scale-[0.97]',
              isFollowing
                ? 'border border-spotify/50 bg-spotify/10 text-spotify hover:bg-spotify/20'
                : 'bg-spotify text-black hover:bg-spotify-hover hover:scale-105',
            ]"
            @click="$emit('toggleFollow')"
          >
            <i aria-hidden="true" :class="isFollowing ? 'pi pi-check' : 'pi pi-plus'" class="text-xs" />
            {{ isFollowing ? 'Following' : 'Follow' }}
          </button>
          <button
            aria-label="Share profile"
            class="inline-flex items-center justify-center rounded-full border border-white/15 bg-white/4 p-3 text-white/60 backdrop-blur-xs transition hover:bg-white/10 hover:text-white active:scale-95 focus-visible:outline-2 focus-visible:outline-[#1db954]"
            @click="$emit('share')"
          >
            <i aria-hidden="true" class="pi pi-share-alt text-sm" />
          </button>
          <button
            v-if="isOwnProfile"
            class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/4 px-5 py-3 text-sm font-medium text-white/60 backdrop-blur-xs transition hover:bg-white/10 hover:text-white focus-visible:outline-2 focus-visible:outline-[#1db954]"
            @click="$emit('editProfile')"
          >
            <i aria-hidden="true" class="pi pi-pencil text-xs" />
            Edit Profile
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { onImgError } from '@/utils/helpers'

const props = defineProps<{
  displayName: string
  handle?: string
  bio?: string
  avatarUrl?: string | null
  followerCount: number
  followingCount: number
  isOwnProfile: boolean
  isFollowing: boolean
  isCreator?: boolean
  joinDate?: string
  genreColor?: string
  musicStatus?: {
    playing: boolean
    currentTrack?: {
      id: string
      title: string
      artist: string
      coverUrl: string
    } | null
  } | null
}>()

defineEmits<{
  toggleFollow: []
  showFollowers: []
  showFollowing: []
  share: []
  editProfile: []
  listenAlong: [trackId: string]
}>()

const currentTrackInfo = computed(() => {
  return props.musicStatus?.currentTrack || null
})

const heroGradient = computed(() => {
  const c = props.genreColor || '#1db954'
  return `from-${c}20 via-transparent to-transparent`
})

const orbColor = computed(() => {
  return props.genreColor || '#1db954'
})

const orbColor2 = computed(() => {
  return props.genreColor ? `${props.genreColor}88` : '#a855f788'
})

const formatJoinDate = computed(() => {
  if (props.joinDate!) return ''
  const d = new Date(props.joinDate)
  try {
    return d.toLocaleDateString('fa-IR', { year: 'numeric', month: 'long' })
  } catch {
    return d.toLocaleDateString()
  }
})

function formatNumber(n: number) {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`
  if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`
  return String(n)
}
</script>

<style scoped>
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in-up {
  animation: fade-in-up 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes equalizer {
  0%, 100% { transform: scaleY(1); }
  50% { transform: scaleY(2.5); }
}
.animate-equalizer {
  animation: equalizer 0.6s ease-in-out infinite;
  transform-origin: bottom;
}
@keyframes ping {
  75%, 100% { transform: scale(2); opacity: 0; }
}
.animate-ping {
  animation: ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite;
}

@media (prefers-reduced-motion: reduce) {
  .animate-fade-in-up { animation: none; }
  .animate-equalizer { animation: none; }
  .animate-ping { animation: none; }
}
</style>
