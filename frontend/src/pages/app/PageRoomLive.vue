<template>
  <div class="mx-auto max-w-5xl space-y-6 px-4 pt-20 pb-24 md:px-8">
    <button
      class="inline-flex items-center gap-1.5 text-sm text-white/40 transition hover:text-white/70"
      @click="goBack"
    >
      &larr; Back to Social
    </button>

    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="card" />
      <SkeletonLoader variant="card" class="h-48" />
      <SkeletonLoader variant="card" class="h-64" />
    </div>

    <div v-else-if="error" class="rounded-2xl bg-white/[0.03] p-12 text-center">
      <p class="text-sm text-white/40">{{ error }}</p>
    </div>

    <template v-else-if="room">
      <div class="grid gap-6 lg:grid-cols-3">

        <!-- Main column -->
        <div class="space-y-6 lg:col-span-2">

          <!-- Room header -->
          <div class="glass-strong rounded-2xl p-6">
            <div class="flex items-start gap-4">
              <div class="flex h-16 w-16 shrink-0 items-center justify-center rounded-xl bg-[#1db954]/10 text-2xl">
                🎤
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-3">
                  <h1 class="truncate text-2xl font-black text-white">{{ room.title }}</h1>
                  <span class="flex items-center gap-1.5 rounded-full bg-red-500/10 px-3 py-1 text-xs font-semibold text-red-400">
                    <span class="h-1.5 w-1.5 animate-pulse rounded-full bg-red-500" />
                    LIVE
                  </span>
                </div>
                <p v-if="room.description" class="mt-1.5 line-clamp-2 text-sm text-white/40">
                  {{ room.description }}
                </p>
                <div class="mt-3 flex items-center gap-4 text-xs text-white/30">
                  <span>{{ room.listener_count }} listening</span>
                  <span>Host: {{ userName(room.host_id) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Now Playing with full controls -->
          <div class="glass-strong rounded-2xl overflow-hidden">
            <div class="flex flex-col sm:flex-row">
              <!-- Album art -->
              <div class="flex items-center justify-center bg-white/[0.02] p-8 sm:w-56">
                <div
                  class="relative h-40 w-40 shrink-0"
                  :class="playerStore.isPlaying ? 'animate-spin-slow' : ''"
                >
                  <div class="absolute inset-0 rounded-full bg-gradient-to-br from-gray-800 to-gray-900 ring-2 ring-white/10" />
                  <img
                    v-if="currentTrack?.coverUrl"
                    :src="currentTrack.coverUrl"
                    class="h-full w-full rounded-full object-cover p-2"
                  />
                  <div v-else class="flex h-full w-full items-center justify-center rounded-full p-2">
                    <div class="flex h-full w-full items-center justify-center rounded-full bg-[#1db954]/10 text-4xl">
                      🎵
                    </div>
                  </div>
                  <div class="absolute inset-0 flex items-center justify-center">
                    <div class="h-8 w-8 rounded-full bg-black/60 ring-2 ring-white/20" />
                  </div>
                </div>
              </div>

              <!-- Track info + controls -->
              <div class="flex flex-1 flex-col justify-center p-6">
                <h2 class="mb-1 text-xs font-bold uppercase tracking-wider text-white/30">Now Playing</h2>
                <div v-if="currentTrack">
                  <p class="text-lg font-bold text-white">{{ currentTrack.title }}</p>
                  <p class="text-sm text-white/50">{{ currentTrack.artistName }}</p>
                  <p v-if="currentTrack.albumTitle" class="text-xs text-white/30">{{ currentTrack.albumTitle }}</p>
                </div>
                <p v-else class="text-sm text-white/30">Waiting for track...</p>

                <!-- Controls -->
                <div class="mt-4 flex items-center gap-3">
                  <button
                    class="flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white transition hover:bg-white/20"
                    @click="togglePlay"
                  >
                    <span v-if="isPlayingTrack" class="text-lg">⏸</span>
                    <span v-else class="ml-0.5 text-lg">▶</span>
                  </button>
                  <button
                    class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 transition hover:text-white/70"
                    @click="playerStore.toggleMute()"
                  >
                    <span>{{ playerStore.muted || playerStore.volume === 0 ? '🔇' : playerStore.volume < 0.5 ? '🔉' : '🔊' }}</span>
                  </button>
                  <input
                    type="range"
                    min="0"
                    max="1"
                    step="0.01"
                    :value="playerStore.volume"
                    @input="playerStore.setVolume(Number(($event.target as HTMLInputElement).value))"
                    class="h-1 w-20 cursor-pointer appearance-none rounded-full bg-white/10 accent-[#1db954] [&::-webkit-slider-thumb]:h-3 [&::-webkit-slider-thumb]:w-3 [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white"
                  />
                </div>

                <!-- Seek bar -->
                <div v-if="currentTrack" class="mt-3">
                  <input
                    type="range"
                    min="0"
                    :max="playerStore.duration || 0"
                    step="1"
                    :value="playerStore.currentTime"
                    @input="playerStore.seek(Number(($event.target as HTMLInputElement).value))"
                    class="h-1 w-full cursor-pointer appearance-none rounded-full bg-white/10 accent-[#1db954] [&::-webkit-slider-thumb]:h-3 [&::-webkit-slider-thumb]:w-3 [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-white"
                  />
                  <div class="mt-1 flex justify-between text-[10px] text-white/30">
                    <span>{{ formatTime(playerStore.currentTime) }}</span>
                    <span>{{ formatTime(playerStore.duration) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Chat -->
          <div class="glass-strong rounded-2xl p-6">
            <h2 class="mb-4 text-sm font-bold uppercase tracking-wider text-white/30">Chat</h2>
            <div ref="chatContainer" class="mb-4 max-h-64 space-y-2 overflow-y-auto">
              <div
                v-for="(msg, i) in messages"
                :key="i"
                class="rounded-lg bg-white/[0.03] px-3 py-2"
              >
                <span class="text-xs font-semibold text-[#1db954]">{{ msg.userName || msg.user_id?.slice(0, 8) || 'System' }}</span>
                <p class="mt-0.5 text-sm text-white/70">{{ msg.content }}</p>
              </div>
              <p v-if="!messages.length" class="text-center text-xs text-white/20">No messages yet</p>
            </div>
            <form class="flex gap-2" @submit.prevent="sendMessage">
              <input
                v-model="chatInput"
                type="text"
                placeholder="Type a message..."
                class="min-w-0 flex-1 rounded-xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white placeholder-white/20 outline-none transition focus:border-white/20"
              />
              <button
                type="submit"
                class="rounded-xl bg-[#1db954]/10 px-4 py-2.5 text-sm font-semibold text-[#1db954] transition hover:bg-[#1db954]/20 disabled:opacity-40"
                :disabled="!chatInput.trim()"
              >
                Send
              </button>
            </form>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">

          <!-- Participants -->
          <div class="glass-strong rounded-2xl p-6">
            <h2 class="mb-4 text-sm font-bold uppercase tracking-wider text-white/30">
              Participants ({{ participants.length }})
            </h2>
            <div class="space-y-2">
              <div
                v-for="p in participants"
                :key="p.id"
                class="flex items-center gap-3 rounded-lg bg-white/[0.03] px-3 py-2"
              >
                <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-[#1db954]/20 text-xs font-bold text-[#1db954]">
                  {{ userName(p.user_id)?.charAt(0).toUpperCase() || '?' }}
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-xs font-medium text-white">{{ userName(p.user_id) || 'Unknown' }}</p>
                  <p class="text-[10px] text-white/20">{{ p.role }}</p>
                </div>
                <span v-if="p.role === 'host'" class="rounded bg-yellow-500/10 px-2 py-0.5 text-[10px] font-semibold text-yellow-400">Host</span>
              </div>
              <p v-if="!participants.length" class="text-center text-xs text-white/20">No participants yet</p>
            </div>
          </div>

          <!-- Queue -->
          <div class="glass-strong rounded-2xl p-6">
            <h2 class="mb-4 text-sm font-bold uppercase tracking-wider text-white/30">Queue ({{ queue.length }})</h2>

            <!-- Now playing indicator in queue -->
            <div v-if="currentTrack" class="mb-3 flex items-center gap-3 rounded-lg bg-[#1db954]/10 px-3 py-2.5 border border-[#1db954]/20">
              <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-[#1db954]/20">
                <i class="pi pi-play text-xs text-[#1db954]" />
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-xs font-semibold text-green-400">Now Playing</p>
                <p class="truncate text-[11px] text-white/50">{{ currentTrack.title }} — {{ currentTrack.artistName }}</p>
              </div>
            </div>

            <div class="space-y-1">
              <div
                v-for="(item, i) in queue"
                :key="item.id"
                class="flex items-center gap-3 rounded-lg bg-white/[0.03] px-3 py-2 transition hover:bg-white/[0.06]"
              >
                <span class="w-5 shrink-0 text-center text-[11px] text-white/20">{{ i + 1 }}</span>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-xs font-medium text-white">{{ queueTrackName(item.track_id) }}</p>
                  <p class="truncate text-[10px] text-white/30">Added by {{ userName(item.added_by) || 'someone' }}</p>
                </div>
                <!-- Remove button can be added when API supports it -->
              </div>
              <p v-if="!queue.length" class="py-6 text-center text-xs text-white/20">Queue is empty</p>
              <button
                class="mt-2 flex w-full items-center justify-center gap-2 rounded-xl bg-[#1db954]/10 py-2.5 text-sm font-semibold text-[#1db954] transition hover:bg-[#1db954]/20"
                @click="showTrackPicker = true"
              >
                <i class="pi pi-plus-circle" /> Add to Queue
              </button>
            </div>
          </div>

          <!-- Leave button -->
          <button
            class="w-full rounded-xl border border-red-500/20 py-3 text-sm font-bold text-red-400 transition hover:bg-red-500/10"
            @click="handleLeave"
          >
            Leave Room
          </button>
        </div>
      </div>
    </template>
    <!-- Track picker -->
    <TrackPickerDialog
      :visible="showTrackPicker"
      title="Add track to queue"
      @update:visible="showTrackPicker = false"
      @select="onTrackSelected"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { usePlayerApi } from '@/services/api/player'
import { getPublicUserProfile } from '@/services/api/users'
import { usePlayerStore } from '@/stores/player'
import { useUserAuthStore } from '@/stores'
import { wsClient } from '@/services/socket'
import { TrackPickerDialog } from '@/components/social'
import type { LiveRoom, LiveRoomParticipant, LiveRoomQueueItem } from '@/services/api/social'
import type { PlaybackTrack } from '@/services/api/player'

const router = useRouter()
const route = useRoute()
const api = useSocialApi()
const playerApi = usePlayerApi()
const playerStore = usePlayerStore()
const auth = useUserAuthStore()

const roomId = route.params.id as string

const loading = ref(true)
const error = ref('')
const room = ref<LiveRoom | null>(null)
const participants = ref<LiveRoomParticipant[]>([])
const queue = ref<LiveRoomQueueItem[]>([])
const currentTrack = ref<PlaybackTrack | null>(null)
const showTrackPicker = ref(false)
const messages = ref<{ user_id: string; userName: string; content: string }[]>([])
const chatInput = ref('')
const chatContainer = ref<HTMLElement | null>(null)

const userNames = ref<Record<string, string>>({})
const trackNames = ref<Record<string, string>>({})

let unsubscribeWs: (() => void) | null = null

const isPlayingTrack = computed(() => {
  return playerStore.isPlaying && playerStore.currentTrack?.id === currentTrack.value?.id
})

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

async function fetchTrackName(trackId: string) {
  if (trackNames.value[trackId]) return
  try {
    const track = await playerApi.getPlaybackTrack(trackId)
    trackNames.value[trackId] = `${track.artistName} - ${track.title}`
  } catch {
    trackNames.value[trackId] = trackId.slice(0, 12)
  }
}

function queueTrackName(trackId: string): string {
  return trackNames.value[trackId] || trackId?.slice(0, 12) || 'Unknown'
}

async function loadRoom() {
  try {
    const [roomData, participantsData, queueData] = await Promise.all([
      api.getRoom(roomId),
      api.getRoomParticipants(roomId),
      api.getRoomQueue(roomId),
    ])
    room.value = roomData
    participants.value = Array.isArray(participantsData) ? participantsData : []
    queue.value = Array.isArray(queueData) ? queueData : []

    const userIds = new Set<string>()
    userIds.add(roomData.host_id)
    if (Array.isArray(participantsData)) {
      participantsData.forEach((p: LiveRoomParticipant) => userIds.add(p.user_id))
    }
    if (Array.isArray(queueData)) {
      queueData.forEach((q: LiveRoomQueueItem) => userIds.add(q.added_by))
    }
    await Promise.all(Array.from(userIds).map(fetchUserName))

    const trackIds = new Set<string>()
    if (roomData.current_track_id) trackIds.add(roomData.current_track_id)
    if (Array.isArray(queueData)) {
      queueData.forEach((q: LiveRoomQueueItem) => trackIds.add(q.track_id))
    }
    await Promise.all(Array.from(trackIds).map(fetchTrackName))

    if (roomData.current_track_id) {
      loadCurrentTrack(roomData.current_track_id)
    }
  } catch (e: any) {
    error.value = e?.message || 'Failed to load room'
  } finally {
    loading.value = false
  }
}

async function loadCurrentTrack(trackId: string) {
  try {
    currentTrack.value = await playerApi.getPlaybackTrack(trackId)
  } catch {
    // Track details not available
  }
}

function togglePlay() {
  if (isPlayingTrack.value) {
    playerStore.pause()
  } else if (currentTrack.value) {
    if (playerStore.currentTrack?.id === currentTrack.value.id) {
      playerStore.resume()
    } else {
      playerStore.playTrackById(room.value!.current_track_id!)
    }
  }
}

function setupWebSocket() {
  wsClient.connect()
  wsClient.subscribe(`room:${roomId}`)
  const unsub = wsClient.onAny((msg: any) => {
    if (msg.type === 'room.participant_joined') {
      loadParticipants()
      if (msg.payload?.user_id) fetchUserName(msg.payload.user_id)
    } else if (msg.type === 'room.participant_left') {
      loadParticipants()
    } else if (msg.type === 'room.track_changed') {
      loadRoom()
      if (msg.payload?.track_id) {
        fetchTrackName(msg.payload.track_id)
        loadCurrentTrack(msg.payload.track_id)
        playerStore.playTrackById(msg.payload.track_id)
      }
    } else if (msg.type === 'room.queue_updated') {
      loadQueue()
    } else if (msg.type === 'room.message') {
      const uid = msg.payload?.user_id || ''
      if (uid && !userNames.value[uid]) fetchUserName(uid)
      messages.value.push({
        user_id: uid,
        userName: userName(uid),
        content: msg.payload?.content || '',
      })
      scrollChat()
    } else if (msg.type === 'room.ended') {
      loadRoom()
    }
  })
  unsubscribeWs = unsub
}

async function loadParticipants() {
  try {
    const data = await api.getRoomParticipants(roomId)
    participants.value = Array.isArray(data) ? data : []
    if (Array.isArray(data)) {
      await Promise.all(data.map((p: LiveRoomParticipant) => fetchUserName(p.user_id)))
    }
  } catch { /* ignore */ }
}

async function loadQueue() {
  try {
    const data = await api.getRoomQueue(roomId)
    queue.value = Array.isArray(data) ? data : []
    if (Array.isArray(data)) {
      await Promise.all(data.map((q: LiveRoomQueueItem) => fetchTrackName(q.track_id)))
    }
  } catch { /* ignore */ }
}

function scrollChat() {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight
    }
  })
}

function sendMessage() {
  if (!chatInput.value.trim()) return
  wsClient.send('room.message', {
    room_id: roomId,
    content: chatInput.value.trim(),
  })
  chatInput.value = ''
}

async function handleLeave() {
  try {
    await api.leaveRoom(roomId)
  } catch { /* ignore */ }
  router.push({ name: 'social' })
}

function goBack() {
  router.push({ name: 'social' })
}

async function onTrackSelected(track: any) {
  showTrackPicker.value = false
  try {
    await api.addToRoomQueue(roomId, String(track.id))
    await loadQueue()
  } catch { /* ignore */ }
}

function formatTime(s: number): string {
  if (!s || !isFinite(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

onMounted(async () => {
  await loadRoom()
  setupWebSocket()
})

onUnmounted(() => {
  if (unsubscribeWs) unsubscribeWs()
  wsClient.unsubscribe(`room:${roomId}`)
})
</script>

<style scoped>
@keyframes spin-slow {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.animate-spin-slow {
  animation: spin-slow 8s linear infinite;
}
</style>
