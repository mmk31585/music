<template>
  <div class="mx-auto max-w-4xl space-y-6 px-4 pt-20 pb-24 md:px-8">
    <button
      class="inline-flex items-center gap-1.5 text-sm text-white/40 transition hover:text-white/70"
      @click="goBack"
    >
      &larr; Back to Social
    </button>

    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="card" />
      <SkeletonLoader variant="card" class="h-48" />
    </div>

    <div v-else-if="error" class="rounded-2xl bg-white/[0.03] p-12 text-center">
      <p class="text-sm text-white/40">{{ error }}</p>
    </div>

    <template v-else-if="party">
      <!-- Party header -->
      <div class="glass-strong rounded-2xl p-6">
        <div class="flex items-start gap-4">
          <div class="flex h-16 w-16 shrink-0 items-center justify-center rounded-xl bg-[#1db954]/10 text-2xl">
            🎉
          </div>
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-3">
              <h1 class="truncate text-2xl font-black text-white">{{ party.title }}</h1>
              <span
                class="rounded-full px-3 py-1 text-xs font-semibold"
                :class="statusClass"
              >
                {{ party.status }}
              </span>
            </div>
            <p v-if="party.description" class="mt-1.5 line-clamp-2 text-sm text-white/40">
              {{ party.description }}
            </p>
            <div class="mt-3 flex items-center gap-4 text-xs text-white/30">
              <span>{{ party.participant_count }} participants</span>
              <span>Hosted by {{ userName(party.host_id) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Now Playing -->
      <div class="now-playing-card">
        <h2 class="mb-4 text-xs font-bold uppercase tracking-wider text-white/30">Now Playing</h2>
        <div v-if="currentTrack" class="flex flex-col sm:flex-row sm:items-center gap-4 sm:gap-6">
          <div class="relative mx-auto sm:mx-0">
            <div class="disc" :class="{ spinning: isPlayingTrack }">
              <div class="disc-inner">
                <img
                  v-if="currentTrack.coverUrl"
                  :src="currentTrack.coverUrl"
                  :alt="currentTrack.albumTitle || currentTrack.title"
                />
                <div v-else class="flex items-center justify-center h-full w-full">
                  <i class="pi pi-headphones text-2xl text-white/40" />
                </div>
              </div>
              <div class="disc-hole" />
            </div>
          </div>
          <div class="min-w-0 flex-1 text-center sm:text-left">
            <p class="truncate text-lg font-bold text-white">{{ currentTrack.title }}</p>
            <p class="truncate text-sm text-white/50">{{ currentTrack.artistName }}</p>
            <p v-if="currentTrack.albumTitle" class="truncate text-xs text-white/30 mt-0.5">{{ currentTrack.albumTitle }}</p>
            <div class="mt-3 flex items-center justify-center sm:justify-start gap-3">
              <button
                class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/20 hover:scale-105 active:scale-95"
                @click="togglePlay"
              >
                <i v-if="isPlayingTrack" class="pi pi-pause text-lg" />
                <i v-else class="pi pi-play ml-0.5 text-lg" />
              </button>
              <span class="text-xs text-white/30">
                <span v-if="isPlayingTrack" class="text-green-400">● Playing</span>
                <span v-else-if="party?.status === 'paused'" class="text-yellow-400">⏸ Paused</span>
                <span v-else-if="party?.status === 'ended'" class="text-white/40">Ended</span>
              </span>
            </div>
          </div>
        </div>
        <div v-else class="flex flex-col items-center gap-3 py-8 text-sm text-white/30">
          <i class="pi pi-headphones text-3xl" />
          <span>Waiting for host to start playing...</span>
        </div>
      </div>

      <!-- Controls -->
      <div v-if="party.host_id === auth.user?.id" class="glass-strong rounded-2xl p-6">
        <h2 class="mb-4 text-sm font-bold uppercase tracking-wider text-white/30">Host Controls</h2>
        <div class="flex flex-wrap gap-3">
          <button
            class="rounded-xl bg-green-500/10 px-4 py-2 text-sm font-semibold text-green-400 transition hover:bg-green-500/20"
            @click="updateStatus('active')"
          >Play</button>
          <button
            class="rounded-xl bg-yellow-500/10 px-4 py-2 text-sm font-semibold text-yellow-400 transition hover:bg-yellow-500/20"
            @click="updateStatus('paused')"
          >Pause</button>
          <button
            class="rounded-xl bg-red-500/10 px-4 py-2 text-sm font-semibold text-red-400 transition hover:bg-red-500/20"
            @click="updateStatus('ended')"
          >End</button>
          <button
            class="rounded-xl bg-blue-500/10 px-4 py-2 text-sm font-semibold text-blue-400 transition hover:bg-blue-500/20"
            @click="showTrackPicker = true"
          ><i class="pi pi-music mr-1" />Select Track</button>
        </div>
      </div>

      <!-- Join/Leave -->
      <button
        class="w-full rounded-xl py-3 text-sm font-bold transition"
        :class="isParticipant
          ? 'border border-red-500/20 text-red-400 hover:bg-red-500/10'
          : 'bg-[#1db954] text-black hover:bg-[#1db954]/90'"
        @click="isParticipant ? handleLeave() : handleJoin()"
      >
        {{ isParticipant ? 'Leave Party' : 'Join Party' }}
      </button>
      <!-- Track picker -->
      <TrackPickerDialog
        :visible="showTrackPicker"
        title="Choose a track for the party"
        @update:visible="showTrackPicker = false"
        @select="onTrackSelected"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { usePlayerApi } from '@/services/api/player'
import { usePlayerStore } from '@/stores/player'
import { getPublicUserProfile } from '@/services/api/users'
import { useUserAuthStore } from '@/stores'
import { TrackPickerDialog } from '@/components/social'
import type { ListeningParty } from '@/services/api/social'
import type { PlaybackTrack } from '@/services/api/player'
import type { SearchTrack } from '@/services/api/catalog/search'

const router = useRouter()
const route = useRoute()
const api = useSocialApi()
const playerApi = usePlayerApi()
const playerStore = usePlayerStore()
const auth = useUserAuthStore()

const partyId = route.params.id as string

const loading = ref(true)
const error = ref('')
const party = ref<ListeningParty | null>(null)
const currentTrack = ref<PlaybackTrack | null>(null)
const userNames = ref<Record<string, string>>({})

const isParticipant = ref(false)
const showTrackPicker = ref(false)

function onTrackSelected(track: SearchTrack) {
  showTrackPicker.value = false
  const trackId = String(track.id)
  updateStatus('active', trackId)
  party.value!.current_track_id = trackId
  currentTrack.value = {
    id: trackId,
    title: track.title,
    artistName: track.artist_name || 'Unknown',
    albumTitle: track.album_title || null,
    coverUrl: track.cover_url || null,
    durationSeconds: track.duration_seconds ?? null,
    streamUrl: playerApi.getTrackStreamUrl(trackId),
  }
}

async function fetchUserName(userId: string) {
  if (userNames.value[userId]) return
  try {
    const profile = await getPublicUserProfile(userId)
    userNames.value[userId] = profile.displayName || profile.username || userId.slice(0, 8)
  } catch {
    userNames.value[userId] = userId.slice(0, 8)
  }
}

function userName(userId: string): string {
  return userNames.value[userId] || userId?.slice(0, 8) || ''
}

const statusClass = computed(() => {
  if (!party.value) return ''
  switch (party.value.status) {
    case 'active': return 'bg-green-500/10 text-green-400'
    case 'paused': return 'bg-yellow-500/10 text-yellow-400'
    case 'ended': return 'bg-white/10 text-white/40'
    default: return 'bg-white/10 text-white/40'
  }
})

const isPlayingTrack = computed(() => {
  return playerStore.isPlaying && playerStore.currentTrack?.id === currentTrack.value?.id
})

async function loadParty() {
  try {
    const partyData = await api.getParty(partyId)
    party.value = partyData
    await fetchUserName(partyData.host_id)

    if (partyData.current_track_id) {
      try {
        currentTrack.value = await playerApi.getPlaybackTrack(partyData.current_track_id)
      } catch { /* ignore */ }
    }

    isParticipant.value = false
  } catch (e: any) {
    error.value = e?.message || 'Failed to load party'
  } finally {
    loading.value = false
  }
}

function togglePlay() {
  if (isPlayingTrack.value) {
    playerStore.pause()
  } else if (currentTrack.value && party.value?.current_track_id) {
    if (playerStore.currentTrack?.id === party.value.current_track_id) {
      playerStore.resume()
    } else {
      playerStore.playTrackById(party.value.current_track_id)
    }
  }
}

async function updateStatus(status: string, trackId?: string) {
  try {
    await api.updatePartyStatus(partyId, status, trackId)
    if (party.value) party.value.status = status as ListeningParty['status']
  } catch { /* ignore */ }
}

async function handleJoin() {
  try {
    await api.joinParty(partyId)
    isParticipant.value = true
  } catch { /* ignore */ }
}

async function handleLeave() {
  try {
    await api.leaveParty(partyId)
    isParticipant.value = false
  } catch { /* ignore */ }
}

function goBack() {
  router.push({ name: 'social' })
}

onMounted(loadParty)
</script>

<style scoped>
.now-playing-card {
  border-radius: 20px;
  padding: 24px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.07);
  transition: border-color 0.3s ease;
}

.now-playing-card:has(.spinning) {
  border-color: rgba(29, 185, 84, 0.2);
}

.disc {
  width: 120px;
  height: 120px;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.disc-inner {
  width: 100%;
  height: 100%;
  border-radius: 999px;
  overflow: hidden;
  background: #1a1a1a;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
  border: 2px solid rgba(255, 255, 255, 0.1);
}

.disc-inner img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.disc.spinning .disc-inner {
  animation: spinDisc 6s linear infinite;
}

.disc-hole {
  position: absolute;
  width: 28px;
  height: 28px;
  border-radius: 999px;
  background: rgba(0, 0, 0, 0.5);
  border: 2px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(4px);
}

@keyframes spinDisc {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
