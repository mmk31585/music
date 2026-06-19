<template>
  <div class="relative overflow-hidden rounded-[2rem]">
    <div
      class="absolute inset-0 bg-gradient-to-b from-[#1db954]/20 via-[#121212]/60 to-[#121212]"
    />
    <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />

    <div
      class="relative z-10 flex flex-col items-center gap-6 px-6 pt-16 pb-8 text-center md:flex-row md:items-end md:gap-8 md:pt-8 md:text-left"
    >
      <div
        class="h-36 w-36 shrink-0 overflow-hidden rounded-full border-4 border-white/10 shadow-2xl md:h-44 md:w-44"
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
          class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#121212] text-4xl text-white/40"
        >
          <i aria-hidden="true" class="pi pi-user" />
        </div>
      </div>

      <div class="flex-1">
        <h1 class="text-3xl leading-tight font-black text-white md:text-5xl">
          {{ displayName }}
        </h1>

        <div
          class="mt-3 flex flex-wrap items-center justify-center gap-4 text-sm text-slate-400 md:justify-start"
        >
          <button class="transition hover:text-white" @click="$emit('showFollowers')">
            <span class="font-semibold text-white">{{ formatNumber(followerCount) }}</span>
            followers
          </button>
          <button class="transition hover:text-white" @click="$emit('showFollowing')">
            <span class="font-semibold text-white">{{ formatNumber(followingCount) }}</span>
            following
          </button>
          <span
            v-if="isCreator"
            class="rounded-full bg-[#1db954]/10 px-3 py-0.5 text-xs font-semibold text-[#1db954]"
          >
            Creator
          </span>
        </div>

        <div
          v-if="!isOwnProfile"
          class="mt-6 flex flex-wrap items-center justify-center gap-3 md:justify-start"
        >
          <button
            type="button"
            :class="[
              'inline-flex items-center gap-2 rounded-full px-6 py-3 text-sm font-bold transition',
              isFollowing
                ? 'border border-[#1db954]/50 bg-[#1db954]/10 text-[#1db954] hover:bg-[#1db954]/20'
                : 'border border-white/15 bg-white/10 text-white hover:bg-white/15',
            ]"
            @click="$emit('toggleFollow')"
          >
            <i aria-hidden="true" :class="isFollowing ? 'pi pi-check' : 'pi pi-plus'" class="text-xs" />
            {{ isFollowing ? 'Following' : 'Follow' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onImgError } from '@/utils/helpers'

defineProps<{
  displayName: string
  avatarUrl?: string | null
  followerCount: number
  followingCount: number
  isOwnProfile: boolean
  isFollowing: boolean
  isCreator?: boolean
}>()

defineEmits<{
  toggleFollow: []
  showFollowers: []
  showFollowing: []
}>()

function formatNumber(n: number) {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`
  if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`
  return String(n)
}
</script>
