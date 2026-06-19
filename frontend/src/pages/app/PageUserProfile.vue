<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-4 pb-32 md:px-6 lg:px-8">
    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="hero" />
      <div class="space-y-3">
        <SkeletonLoader v-for="i in 3" :key="i" variant="track" />
      </div>
    </div>

    <div v-else-if="error" class="flex flex-col items-center gap-4 py-24 text-center">
      <i aria-hidden="true" class="pi pi-exclamation-circle text-4xl text-slate-500" />
      <h2 class="text-xl font-bold text-white">Failed to load profile</h2>
      <button
        class="text-sm font-medium text-[#1db954] underline underline-offset-2"
        @click="fetchProfile"
      >
        Try again
      </button>
    </div>

    <template v-else>
      <UserHero
        :display-name="displayName"
        :avatar-url="(auth.user as any)?.avatarUrl"
        :follower-count="followerCount"
        :following-count="followingCount"
        :is-own-profile="isOwnProfile"
        :is-following="isFollowing"
        @toggle-follow="toggleFollow"
        @show-followers="activeTab = 'followers'"
        @show-following="activeTab = 'following'"
      />

      <div class="mt-8">
        <ProfileTabs :tabs="tabs" :active-tab="activeTab" @update:active-tab="activeTab = $event" />
      </div>

      <div class="mt-6 space-y-3">
        <!-- Feed Tab -->
        <div v-show="activeTab === 'feed'">
          <div v-if="feed.length === 0" class="flex flex-col items-center gap-3 py-16 text-center">
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-clock text-xl text-slate-500" />
            </div>
            <p class="text-sm font-medium text-white/60">No activity yet</p>
          </div>
          <ActivityItem v-for="item in feed" :key="item.id" :item="item" />
          <button
            v-if="feedHasMore"
            class="mt-4 w-full rounded-xl py-3 text-sm font-medium text-slate-400 transition hover:bg-white/5 hover:text-white"
          >
            Load more
          </button>
        </div>

        <!-- Liked Tracks Tab -->
        <div v-show="activeTab === 'liked-tracks'">
          <div
            v-if="likedTracks.length === 0"
            class="flex flex-col items-center gap-3 py-16 text-center"
          >
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-heart text-xl text-slate-500" />
            </div>
            <p class="text-sm font-medium text-white/60">No liked tracks yet</p>
          </div>
          <TrackList v-else :tracks="likedTracks as any" />
        </div>

        <!-- Liked Albums Tab -->
        <div v-show="activeTab === 'liked-albums'">
          <div
            v-if="likedAlbums.length === 0"
            class="flex flex-col items-center gap-3 py-16 text-center"
          >
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-compact-disc text-xl text-slate-500" />
            </div>
            <p class="text-sm font-medium text-white/60">No liked albums yet</p>
          </div>
          <div v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <AlbumCard
              v-for="album in likedAlbums"
              :key="album.id || album.album_id"
              :album="album"
            />
          </div>
        </div>

        <!-- Followers Tab -->
        <div v-show="activeTab === 'followers'">
          <div
            v-if="followers.length === 0"
            class="flex flex-col items-center gap-3 py-16 text-center"
          >
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-users text-xl text-slate-500" />
            </div>
            <p class="text-sm font-medium text-white/60">No followers yet</p>
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="f in followers"
              :key="f.follower_id"
              class="group flex items-center gap-4 rounded-2xl border border-white/[0.04] bg-white/[0.02] px-4 py-3 transition hover:border-white/[0.08] hover:bg-white/[0.04]"
            >
              <div
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-white/[0.08] to-white/[0.02] text-sm font-bold text-white/70 ring-1 ring-white/[0.04]"
              >
                {{ (f.follower_name || f.follower_id).charAt(0).toUpperCase() }}
              </div>
              <div class="min-w-0">
                <p class="truncate text-sm font-semibold text-white">{{ f.follower_name || f.follower_id }}</p>
                <p class="mt-0.5 text-xs text-white/40">
                  Following since {{ new Date(f.created_at).toLocaleDateString() }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Following Tab -->
        <div v-show="activeTab === 'following'">
          <div
            v-if="following.length === 0"
            class="flex flex-col items-center gap-3 py-16 text-center"
          >
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-users text-xl text-slate-500" />
            </div>
            <p class="text-sm font-medium text-white/60">Not following anyone yet</p>
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="f in following"
              :key="f.followed_id"
              class="group flex items-center gap-4 rounded-2xl border border-white/[0.04] bg-white/[0.02] px-4 py-3 transition hover:border-white/[0.08] hover:bg-white/[0.04]"
            >
              <div
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-white/[0.08] to-white/[0.02] text-sm font-bold text-white/70 ring-1 ring-white/[0.04]"
              >
                {{ (f.followed_name || f.followed_id).charAt(0).toUpperCase() }}
              </div>
              <div class="min-w-0">
                <p class="truncate text-sm font-semibold text-white">{{ f.followed_name || f.followed_id }}</p>
                <p class="mt-0.5 text-xs text-white/40">
                  Following since {{ new Date(f.created_at).toLocaleDateString() }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { useReactionsApi } from '@/services/api/reactions'
import { useUserAuthStore } from '@/stores'
import { TrackList, AlbumCard } from '@/components/music'

const route = useRoute()
const auth = useUserAuthStore()
const socialApi = useSocialApi()
const reactionsApi = useReactionsApi()

const userId = route.params.id as string | undefined
const targetUserId = userId || String(auth.user?.id || '')
const isOwnProfile = !userId || userId === String(auth.user?.id)

const loading = ref(false)
const error = ref<any>(null)
const activeTab = ref('feed')

const followers = ref<any[]>([])
const following = ref<any[]>([])
const followerCount = ref(0)
const followingCount = ref(0)
const isFollowing = ref(false)
const feed = ref<any[]>([])
const likedTracks = ref<any[]>([])
const likedAlbums = ref<any[]>([])
const feedHasMore = ref(false)

const displayName = isOwnProfile
  ? auth.user?.displayName || auth.user?.username || auth.user?.name || 'User'
  : targetUserId || 'User'
const tabs = [
  { key: 'feed', label: 'Activity' },
  { key: 'liked-tracks', label: 'Tracks' },
  { key: 'liked-albums', label: 'Albums' },
  { key: 'followers', label: 'Followers' },
  { key: 'following', label: 'Following' },
]

async function fetchProfile() {
  loading.value = true
  error.value = null
  try {
    const [followersData, followingData, feedData, likedTracksData, likedAlbumsData] =
      await Promise.all([
        socialApi.getFollowers(targetUserId).catch(() => null),
        socialApi.getFollowing(targetUserId).catch(() => null),
        socialApi
          .getFeed({ user_id: isOwnProfile ? undefined : targetUserId, limit: 10 })
          .catch(() => null),
        reactionsApi
          .getLikedTracks({ user_id: isOwnProfile ? undefined : targetUserId, limit: 10 })
          .catch(() => null),
        reactionsApi
          .getLikedAlbums({ user_id: isOwnProfile ? undefined : targetUserId, limit: 10 })
          .catch(() => null),
      ])

    if (followersData) {
      followers.value = followersData.items
      followerCount.value = followersData.total_count
    }
    if (followingData) {
      following.value = followingData.items
      followingCount.value = followingData.total_count
    }
    if (feedData) {
      feed.value = feedData.items
      feedHasMore.value = feedData.pagination.has_more
    }
    if (likedTracksData) likedTracks.value = likedTracksData.items ?? []
    if (likedAlbumsData) likedAlbums.value = likedAlbumsData.items ?? []

    if (!isOwnProfile) {
      const f = await socialApi.isFollowing(targetUserId).catch(() => null)
      if (f) isFollowing.value = f.is_following
    }
  } catch (err) {
    error.value = err
  } finally {
    loading.value = false
  }
}

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
  } catch {
    /* silent */
  }
}

onMounted(fetchProfile)
</script>
