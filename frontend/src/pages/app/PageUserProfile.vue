<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-4 pb-32 md:px-6 lg:px-8">
    <!-- Loading -->
    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="hero" />
      <div class="grid grid-cols-2 gap-3 md:grid-cols-4">
        <SkeletonLoader v-for="i in 4" :key="i" variant="card" />
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="flex flex-col items-center gap-4 py-24 text-center">
      <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/[0.04]">
        <i aria-hidden="true" class="pi pi-exclamation-circle text-4xl text-slate-500" />
      </div>
      <h2 class="text-xl font-bold text-white">Failed to load profile</h2>
      <p class="text-sm text-white/40">Something went wrong. Try again?</p>
      <button
        class="rounded-full bg-[#1db954] px-6 py-2.5 text-sm font-bold text-black transition hover:bg-[#1ed760]"
        @click="fetchProfile"
      >
        Try again
      </button>
    </div>

    <template v-else>
      <!-- ── Cinematic Hero ── -->
      <UserHero
        :display-name="displayName"
        :handle="profileHandle"
        :bio="profileBio"
        :avatar-url="profileAvatar"
        :follower-count="followerCount"
        :following-count="followingCount"
        :is-own-profile="isOwnProfile"
        :is-following="isFollowing"
        :is-creator="isCreator"
        :join-date="profileJoinDate"
        :music-status="musicStatus"
        @toggle-follow="toggleFollow"
        @show-followers="activeTab = 'followers'"
        @show-following="activeTab = 'following'"
        @share="handleShare"
        @edit-profile="handleEditProfile"
        @listen-along="handleListenAlong"
      />

      <!-- ── Stats Gallery ── -->
      <div class="mt-6 grid grid-cols-2 gap-3 md:grid-cols-4">
        <StatCard
          icon="clock"
          :value="formattedListeningHours"
          label="Listening"
          :trend="statsTrends.listening"
        />
        <StatCard
          icon="music"
          :value="topGenreName"
          :sub="topGenrePercent"
          label="Top Genre"
          :color="topGenreColor"
        />
        <StatCard
          icon="award"
          :value="String(badgeCount)"
          label="Badges earned"
        />
        <StatCard
          icon="heart"
          :value="String(likedTracks.length)"
          label="Liked tracks"
        />
      </div>

      <!-- ── Tabs ── -->
      <div class="mt-8">
        <ProfileTabs
          :tabs="tabs"
          :active-tab="activeTab"
          :aria-label="`${displayName}'s profile`"
          @update:active-tab="activeTab = $event"
        />
      </div>

      <div class="mt-6 space-y-3">
        <!-- TRACKS TAB -->
        <div v-show="activeTab === 'tracks'" role="tabpanel" id="tabpanel-tracks">
          <!-- View toggle + count -->
          <div class="mb-4 flex items-center justify-between">
            <p class="text-xs text-white/30 tabular-nums">{{ likedTracks.length }} tracks</p>
            <div class="flex gap-1 rounded-lg bg-white/[0.04] p-0.5">
              <button
                aria-label="Mosaic view"
                class="rounded-md p-1.5 transition"
                :class="trackViewMode === 'mosaic' ? 'bg-white/10 text-white' : 'text-white/30 hover:text-white/50'"
                @click="trackViewMode = 'mosaic'"
              >
                <i aria-hidden="true" class="pi pi-th-large text-xs" />
              </button>
              <button
                aria-label="List view"
                class="rounded-md p-1.5 transition"
                :class="trackViewMode === 'list' ? 'bg-white/10 text-white' : 'text-white/30 hover:text-white/50'"
                @click="trackViewMode = 'list'"
              >
                <i aria-hidden="true" class="pi pi-list text-xs" />
              </button>
            </div>
          </div>

          <!-- MOSAIC VIEW -->
          <div v-if="likedTracks.length > 0 && trackViewMode === 'mosaic'" class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <div
              v-for="track in likedTracks" :key="track.id || track.track_id"
              class="group relative overflow-hidden rounded-2xl bg-white/[0.04] ring-1 ring-white/[0.06]
                     transition-all duration-300 hover:ring-[#1db954]/30 hover:bg-white/[0.06]
                     focus-within:ring-[#1db954]"
            >
              <!-- Cover art -->
              <div class="aspect-square overflow-hidden">
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.track_title"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
                />
                <div v-else class="flex h-full items-center justify-center bg-gradient-to-br from-white/[0.04] to-white/[0.02]">
                  <i aria-hidden="true" class="pi pi-music text-2xl text-white/20" />
                </div>
              </div>

              <!-- Play overlay -->
              <div class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 transition-opacity duration-200">
                <button
                  aria-label="Play"
                  class="flex h-10 w-10 items-center justify-center rounded-full bg-[#1db954] text-black shadow-lg transition hover:scale-110 active:scale-95 focus-visible:outline-2 focus-visible:outline-white"
                  @click="playLikedTrack(track)"
                >
                  <i aria-hidden="true" class="pi pi-play-fill text-sm ms-0.5" />
                </button>
              </div>

              <!-- Info -->
              <div class="p-3">
                <p class="truncate text-sm font-medium text-white">{{ track.track_title || track.title }}</p>
                <p class="truncate text-xs text-white/40">{{ track.artist_name || track.artistName || '' }}</p>
              </div>

              <!-- Visibility badge (own profile) -->
              <div v-if="isOwnProfile" class="absolute top-2 start-2">
                <span
                  class="flex h-6 w-6 items-center justify-center rounded-full backdrop-blur-sm text-[10px]"
                  :class="isTrackPublic(track) ? 'bg-black/30 text-white/60' : 'bg-amber-500/30 text-amber-300'"
                  :title="isTrackPublic(track) ? 'Public' : 'Private'"
                >
                  <i aria-hidden="true" :class="isTrackPublic(track) ? 'pi pi-globe' : 'pi pi-lock'" class="text-[9px]" />
                </span>
              </div>
            </div>
          </div>

          <!-- LIST VIEW -->
          <div v-else-if="likedTracks.length > 0 && trackViewMode === 'list'" class="space-y-1">
            <div
              v-for="(track, i) in likedTracks" :key="track.id || track.track_id || i"
              class="group flex items-center gap-3 rounded-xl px-3 py-2.5 transition hover:bg-white/[0.04] focus-within:bg-white/[0.04]"
            >
              <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  alt=""
                  class="h-full w-full object-cover"
                  loading="lazy"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-music text-xs text-white/30" />
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-white">{{ track.track_title }}</p>
                <p class="truncate text-xs text-white/40">{{ track.artist_name || '' }}</p>
              </div>
              <button
                v-if="isOwnProfile"
                class="flex h-8 w-8 items-center justify-center rounded-lg text-sm opacity-0 group-hover:opacity-100 group-focus-within:opacity-100 transition hover:bg-white/10 focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-[#1db954]"
                :class="isTrackPublic(track) ? 'text-white/40' : 'text-amber-400/70'"
                :aria-label="isTrackPublic(track) ? 'Set private' : 'Set public'"
                @click="toggleTrackVisibility(track, i)"
              >
                <i aria-hidden="true" :class="isTrackPublic(track) ? 'pi pi-globe' : 'pi pi-lock'" class="text-xs" />
              </button>
              <button
                aria-label="Play track"
                class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 transition hover:bg-white/10 hover:text-white focus-visible:outline-2 focus-visible:outline-[#1db954]"
                @click="playLikedTrack(track)"
              >
                <i aria-hidden="true" class="pi pi-play-fill text-xs ms-0.5" />
              </button>
            </div>
          </div>

          <!-- Empty state -->
          <div v-else class="flex flex-col items-center gap-4 py-20 text-center">
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-heart text-3xl text-white/10" />
            </div>
            <div>
              <p class="text-base font-semibold text-white/40">
                {{ isOwnProfile ? 'No liked tracks yet' : 'Tracks are private' }}
              </p>
              <p class="mt-1 text-sm text-white/30">
                {{ isOwnProfile ? 'Heart tracks to save them here' : 'This user keeps their likes private' }}
              </p>
            </div>
            <RouterLink v-if="isOwnProfile" to="/"
              class="rounded-full bg-[#1db954] px-6 py-2.5 text-sm font-bold text-black transition hover:bg-[#1ed760]"
            >
              <i aria-hidden="true" class="pi pi-compass text-xs me-1.5" />
              Discover music
            </RouterLink>
          </div>
        </div>

        <!-- ALBUMS TAB -->
        <div v-show="activeTab === 'albums'" role="tabpanel" id="tabpanel-albums">
          <div v-if="likedAlbums.length > 0" class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <AlbumCard v-for="album in likedAlbums" :key="album.id || album.album_id" :album="album" />
          </div>
          <div v-else class="flex flex-col items-center gap-4 py-20 text-center">
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-compact-disc text-3xl text-white/10" />
            </div>
            <p class="text-base font-semibold text-white/40">No liked albums yet</p>
          </div>
        </div>

        <!-- EDITS TAB -->
        <div v-show="activeTab === 'edits'" role="tabpanel" id="tabpanel-edits">
          <div v-if="userVideos.length > 0" class="grid grid-cols-2 gap-3 sm:gap-4">
            <VideoCard
              v-for="(video, index) in userVideos"
              :key="video.id"
              :video="video"
              @open="openVideoPlayer(index)"
            />
          </div>
          <div v-else class="flex flex-col items-center gap-4 py-20 text-center">
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-video text-3xl text-white/10" />
            </div>
            <div>
              <p class="text-base font-semibold text-white/40">
                {{ isOwnProfile ? 'No edits yet' : 'No edits from this user' }}
              </p>
            </div>
            <RouterLink v-if="isOwnProfile" to="/create-edit"
              class="rounded-full bg-[#1db954] px-5 py-2 text-xs font-bold text-black transition hover:bg-[#1ed760]"
            >
              <i aria-hidden="true" class="pi pi-plus text-xs me-1" />
              Create edit
            </RouterLink>
          </div>
          <button
            v-if="userVideosHasMore"
            class="mt-4 w-full rounded-xl py-3 text-sm font-medium text-slate-400 transition hover:bg-white/5 hover:text-white"
            @click="loadMoreVideos"
          >
            بارگذاری بیشتر
          </button>
        </div>

        <!-- FOLLOWERS TAB — Avatar Wall -->
        <div v-show="activeTab === 'followers'" role="tabpanel" id="tabpanel-followers">
          <div v-if="followers.length > 0" class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <RouterLink
              v-for="f in followers"
              :key="f.follower_id"
              :to="`/profile/${f.follower_id}`"
              class="flex flex-col items-center gap-3 rounded-2xl bg-white/[0.04] p-5 ring-1 ring-white/[0.06] transition-all duration-200 hover:bg-white/[0.08] hover:ring-white/[0.12] focus-visible:outline-2 focus-visible:outline-[#1db954]"
            >
              <div class="h-16 w-16 overflow-hidden rounded-full ring-2 ring-white/10">
                <img
                  v-if="f.avatar_url"
                  :src="f.avatar_url"
                  :alt="f.follower_name || f.follower_id"
                  class="h-full w-full object-cover"
                  loading="lazy"
                />
                <div class="flex h-full w-full items-center justify-center bg-gradient-to-br from-white/[0.08] to-white/[0.02] text-lg font-bold text-white/50">
                  {{ (f.follower_name || f.follower_id).charAt(0).toUpperCase() }}
                </div>
              </div>
              <div class="text-center">
                <p class="truncate text-sm font-semibold text-white max-w-[100px]">{{ f.follower_name || f.follower_id }}</p>
                <p class="text-[10px] text-white/30">Follows you</p>
              </div>
            </RouterLink>
          </div>
          <div v-else class="flex flex-col items-center gap-4 py-20 text-center">
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-users text-3xl text-white/10" />
            </div>
            <p class="text-base font-semibold text-white/40">No followers yet</p>
            <p class="text-sm text-white/30">Share your profile to grow your community</p>
          </div>
        </div>

        <!-- FOLLOWING TAB — Avatar Wall -->
        <div v-show="activeTab === 'following'" role="tabpanel" id="tabpanel-following">
          <div v-if="following.length > 0" class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <RouterLink
              v-for="f in following"
              :key="f.followed_id"
              :to="`/profile/${f.followed_id}`"
              class="flex flex-col items-center gap-3 rounded-2xl bg-white/[0.04] p-5 ring-1 ring-white/[0.06] transition-all duration-200 hover:bg-white/[0.08] hover:ring-white/[0.12] focus-visible:outline-2 focus-visible:outline-[#1db954]"
            >
              <div class="h-16 w-16 overflow-hidden rounded-full ring-2 ring-white/10">
                <img
                  v-if="f.avatar_url"
                  :src="f.avatar_url"
                  :alt="f.followed_name || f.followed_id"
                  class="h-full w-full object-cover"
                  loading="lazy"
                />
                <div class="flex h-full w-full items-center justify-center bg-gradient-to-br from-white/[0.08] to-white/[0.02] text-lg font-bold text-white/50">
                  {{ (f.followed_name || f.followed_id).charAt(0).toUpperCase() }}
                </div>
              </div>
              <div class="text-center">
                <p class="truncate text-sm font-semibold text-white max-w-[100px]">{{ f.followed_name || f.followed_id }}</p>
                <p class="text-[10px] text-white/30">Following</p>
              </div>
            </RouterLink>
          </div>
          <div v-else class="flex flex-col items-center gap-4 py-20 text-center">
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-user-plus text-3xl text-white/10" />
            </div>
            <p class="text-base font-semibold text-white/40">Not following anyone yet</p>
          </div>
        </div>
      </div>
    </template>

    <!-- Video player overlay -->
    <VerticalVideoPlayer
      v-if="playerOpen"
      :videos="allUserVideos"
      :start-index="playerStartIndex"
      @close="playerOpen = false"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { useReactionsApi } from '@/services/api/reactions'
import { useVideoApi } from '@/services/api/video'
import { useUserAuthStore } from '@/stores'
import { usePlayerStore } from '@/stores/player'
import { AlbumCard } from '@/components/music'
import UserHero from '@/components/music/profile/UserHero.vue'
import ProfileTabs from '@/components/music/profile/ProfileTabs.vue'
import StatCard from '@/components/music/profile/StatCard.vue'

import VideoCard from '@/components/video/VideoCard.vue'
import VerticalVideoPlayer from '@/components/video/VerticalVideoPlayer.vue'
import type { VideoItem, LikedTrackItem } from '@/services/api/video/types'

const route = useRoute()
const auth = useUserAuthStore()
const playerStore = usePlayerStore()
const socialApi = useSocialApi()
const reactionsApi = useReactionsApi()
const videoApi = useVideoApi()

const userId = route.params.id as string | undefined
const targetUserId = userId || String(auth.user?.id || '')
const isOwnProfile = !userId || userId === String(auth.user?.id)

const loading = ref(false)
const error = ref<any>(null)
const activeTab = ref('tracks')
const trackViewMode = ref<'mosaic' | 'list'>('mosaic')

const followers = ref<any[]>([])
const following = ref<any[]>([])
const followerCount = ref(0)
const followingCount = ref(0)
const isFollowing = ref(false)
const likedTracks = ref<(LikedTrackItem & { is_public?: boolean })[]>([])
const likedAlbums = ref<any[]>([])
const musicStatus = ref<{
  playing: boolean
  currentTrack?: {
    id: string
    title: string
    artist: string
    coverUrl: string
  } | null
  current_track_id?: string | null
  updated_at?: string | null
} | null>(null)

// Edits tab
const userVideos = ref<VideoItem[]>([])
const allUserVideos = ref<VideoItem[]>([])
const userVideosOffset = ref(0)
const userVideosHasMore = ref(false)
const VIDEOS_LIMIT = 20

// Video player
const playerOpen = ref(false)
const playerStartIndex = ref(0)

// ── Computed ──

const displayName = computed(() => {
  if (isOwnProfile) return auth.user?.displayName || auth.user?.username || auth.user?.name || 'User'
  return (profileData.value as any)?.full_name || (profileData.value as any)?.username || targetUserId || 'User'
})

const profileData = ref<any>(null)

const profileHandle = computed(() => {
  return isOwnProfile ? (auth.user as any)?.handle || null : (profileData.value as any)?.handle || null
})

const profileBio = computed(() => {
  return isOwnProfile ? (auth.user as any)?.bio || null : (profileData.value as any)?.bio || null
})

const profileAvatar = computed(() => {
  if (isOwnProfile) return (auth.user as any)?.avatarUrl
  return (profileData.value as any)?.avatar_url || (profileData.value as any)?.avatarUrl
})

const profileJoinDate = computed(() => {
  return (profileData.value as any)?.created_at || (profileData.value as any)?.createdAt || null
})

const isCreator = computed(() => {
  return (profileData.value as any)?.is_creator || (profileData.value as any)?.isCreator || false
})

const formattedListeningHours = computed(() => {
  // Placeholder — real data from API
  return '—'
})

const statsTrends = computed(() => ({
  listening: undefined as number | undefined,
}))

const topGenreName = computed(() => '—')
const topGenrePercent = computed(() => '')
const topGenreColor = computed(() => '#1db954')
const badgeCount = computed(() => 0)

const tabs = [
  { key: 'tracks', label: 'آهنگ‌ها', icon: 'pi pi-music' },
  { key: 'albums', label: 'آلبوم‌ها', icon: 'pi pi-compact-disc' },
  { key: 'edits', label: 'ادیت‌ها', icon: 'pi pi-video' },
  { key: 'followers', label: 'دنبال‌کننده‌ها', icon: 'pi pi-users' },
  { key: 'following', label: 'دنبال‌شونده‌ها', icon: 'pi pi-user-plus' },
]

// ── Data fetching ──

async function fetchProfile() {
  loading.value = true
  error.value = null
  try {
    const [
      followersData,
      followingData,
      likedTracksData,
      likedAlbumsData,
      userVideosData,
      musicStatusData,
      profile,
    ] = await Promise.all([
      socialApi.getFollowers(targetUserId).catch(() => null),
      socialApi.getFollowing(targetUserId).catch(() => null),
      reactionsApi.getLikedTracks({ user_id: isOwnProfile ? undefined : targetUserId, limit: 50 }).catch(() => null),
      reactionsApi.getLikedAlbums({ user_id: isOwnProfile ? undefined : targetUserId, limit: 10 }).catch(() => null),
      videoApi.getUserVideos(targetUserId, { limit: VIDEOS_LIMIT, offset: 0 }).catch(() => null),
      isOwnProfile ? Promise.resolve(null) : videoApi.getMusicStatus(targetUserId).catch(() => null),
      socialApi.getPublicProfile?.(targetUserId).catch(() => null) || Promise.resolve(null),
    ])

    if (followersData) {
      followers.value = followersData.items || []
      followerCount.value = followersData.total_count || 0
    }
    if (followingData) {
      following.value = followingData.items || []
      followingCount.value = followingData.total_count || 0
    }
    if (likedTracksData) {
      likedTracks.value = (likedTracksData.items ?? []).map((t: any) => ({
        ...t,
        is_public: t.is_public !== false,
      }))
    }
    if (likedAlbumsData) likedAlbums.value = likedAlbumsData.items ?? []

    if (userVideosData) {
      const items: VideoItem[] = userVideosData.items ?? []
      userVideos.value = items
      allUserVideos.value = items
      userVideosOffset.value = items.length
      userVideosHasMore.value = items.length >= VIDEOS_LIMIT
    }

    if (musicStatusData) {
      const ms = musicStatusData as any
      if (ms?.playing && ms?.current_track) {
        musicStatus.value = {
          playing: true,
          currentTrack: {
            id: ms.current_track.id || '',
            title: ms.current_track.title || '',
            artist: ms.current_track.artist || ms.current_track.artist_name || '',
            coverUrl: ms.current_track.cover_url || ms.current_track.coverUrl || '',
          },
          current_track_id: ms.current_track.id,
        }
      } else if (ms?.playing && ms?.current_track_id) {
        musicStatus.value = {
          playing: true,
          currentTrack: null,
          current_track_id: ms.current_track_id,
        }
      } else {
        musicStatus.value = { playing: false, currentTrack: null }
      }
    } else if (isOwnProfile) {
      musicStatus.value = null
    }

    if (profile) profileData.value = profile

    if (!isOwnProfile) {
      const f = await socialApi.isFollowing(targetUserId).catch(() => null)
      if (f) isFollowing.value = (f as any).is_following
    }
  } catch (err) {
    error.value = err
  } finally {
    loading.value = false
  }
}

// ── Actions ──

async function toggleFollow() {
  if (!userId) return
  try {
    if (isFollowing.value) {
      await socialApi.unfollow(userId)
      isFollowing.value = false
      followerCount.value = Math.max(0, followerCount.value - 1)
    } else {
      await socialApi.follow(userId)
      isFollowing.value = true
      followerCount.value += 1
    }
  } catch { /* silent */ }
}

function isTrackPublic(track: any): boolean {
  return track.is_public !== false
}

async function toggleTrackVisibility(track: any, index: number) {
  const newVisibility = isTrackPublic(track) ? 'private' : 'public'
  const trackId = String(track.track_id || track.id)
  const prev = track.is_public
  likedTracks.value[index] = { ...track, is_public: newVisibility === 'public' }
  try {
    await videoApi.setLikeVisibility(trackId, { visibility: newVisibility as 'public' | 'private' })
  } catch {
    likedTracks.value[index] = { ...track, is_public: prev }
  }
}

function playLikedTrack(track: any) {
  const trackId = String(track.track_id || track.id)
  playerStore.playTrackById(trackId)
}

function handleShare() {
  navigator.clipboard?.writeText(window.location.href)
}

function handleEditProfile() {
  // Navigate to settings
}

function handleListenAlong(trackId: string) {
  playerStore.playTrackById(trackId)
}

async function loadMoreVideos() {
  if (!userVideosHasMore.value) return
  try {
    const data = await videoApi.getUserVideos(targetUserId, {
      limit: VIDEOS_LIMIT,
      offset: userVideosOffset.value,
    })
    if (data?.items?.length) {
      const newItems = data.items as VideoItem[]
      userVideos.value = [...userVideos.value, ...newItems]
      allUserVideos.value = [...allUserVideos.value, ...newItems]
      userVideosOffset.value += newItems.length
      userVideosHasMore.value = newItems.length >= VIDEOS_LIMIT
    } else {
      userVideosHasMore.value = false
    }
  } catch {
    userVideosHasMore.value = false
  }
}

function openVideoPlayer(index: number) {
  playerStartIndex.value = index
  playerOpen.value = true
}

onMounted(fetchProfile)
</script>
