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
          <div class="rounded-2xl bg-white/[0.04] p-6 ring-1 ring-white/[0.07]">
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
                  <span>Host: {{ room.host_id ? userName(room.host_id) : 'Unknown' }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Now Playing Hero (from Phase 3) -->
          <RoomNowPlayingHero
            :now-playing="queueSocket.queueState.value?.now_playing || null"
            :is-loading="loading"
            :is-playing="isPlayingTrack"
            @toggle-play="togglePlay"
          />

          <!-- Player controls (seek + volume, below hero) -->
          <div class="rounded-2xl bg-white/[0.04] p-4 ring-1 ring-white/[0.07]">
            <div class="flex items-center gap-3">
              <button
                class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 transition hover:text-white/70"
                @click="playerStore.toggleMute()"
              >
                {{ playerStore.muted || playerStore.volume === 0 ? '🔇' : playerStore.volume < 0.5 ? '🔉' : '🔊' }}
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
              <span class="text-xs text-white/30">Volume</span>
            </div>
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

          <!-- Chat -->
          <div class="rounded-2xl bg-white/[0.04] p-6 ring-1 ring-white/[0.07]">
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

          <!-- Stage (speaker grid) -->
          <LiveRoomStage
            v-if="stageState"
            :stage-state="{ ...stageState }"
            :speaking-user-ids="speakingUserIds"
            :show-host-actions="isHost"
            @remove-speaker="handleRemoveSpeaker"
            @toggle-mute="handleToggleMute"
          />

          <!-- Host-only: hand raise queue -->
          <HandRaiseQueue
            v-if="isHost"
            :pending-raises="[...(stageState?.pending_hand_raises || [])]"
            @approve="handleApproveHand"
            @deny="handleDenyHand"
          />

          <!-- Queue (Phase 3) -->
          <RoomQueueList
            :candidates="[...(queueSocket.queueState.value?.candidates || [])]"
            @vote="handleVote"
            @unvote="handleUnvote"
            @suggest-clicked="showTrackPicker = true"
          />

          <!-- Leave button -->
          <button
            class="w-full rounded-xl border border-red-500/20 py-3 text-sm font-bold text-red-400 transition hover:bg-red-500/10"
            @click="handleLeave"
          >
            Leave Room
          </button>
        </div>
      </div>

      <!-- Floating raise-hand button (visible to listeners) -->
      <RaiseHandButton
        :hand-raised="stageState?.my_hand_raised || false"
        :my-role="stageState?.my_role || 'listener'"
        @raise="handleRaiseHand"
        @lower="handleLowerHand"
      />
    </template>

    <!-- Track picker -->
    <TrackPickerDialog
      :visible="showTrackPicker"
      title="پیشنهاد آهنگ"
      @update:visible="showTrackPicker = false"
      @select="onTrackSelected"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { usePlayerStore } from '@/stores/player'
import { useUserAuthStore } from '@/stores'
import { wsClient } from '@/services/socket'
import { TrackPickerDialog, RoomNowPlayingHero, RoomQueueList, LiveRoomStage, RaiseHandButton, HandRaiseQueue } from '@/components/social'
import { useRoomQueueSocket } from '@/composables/social/useRoomQueueSocket'
import { useStageSocket } from '@/composables/social/useStageSocket'
import type { LiveRoom } from '@/services/api/social'

const router = useRouter()
const route = useRoute()
const api = useSocialApi()
const playerStore = usePlayerStore()
const auth = useUserAuthStore()

const roomId = String(route.params.id)
const currentUserId = computed(() => String(auth.user?.id ?? ''))

const loading = ref(true)
const error = ref('')
const room = ref<LiveRoom | null>(null)
const showTrackPicker = ref(false)
const messages = ref<{ user_id: string; userName: string; content: string }[]>([])
const chatInput = ref('')
const chatContainer = ref<HTMLElement | null>(null)

const userNames = ref<Record<string, string>>({})

// Stage composable
const stage = useStageSocket(roomId, currentUserId.value)
const stageState = stage.stageState

// Queue composable (Phase 3)
const queueSocket = useRoomQueueSocket(roomId)
const queueState = queueSocket.queueState

const currentTrack = computed(() => queueState.value?.now_playing?.track || null)

const isHost = computed(() => stageState.value?.my_role === 'host')
const isPlayingTrack = computed(() => {
  return playerStore.isPlaying && playerStore.currentTrack?.id === currentTrack.value?.id
})

// Speaking user IDs (stub — wire to real audio levels later)
const speakingUserIds = computed(() => new Set<string>())

async function fetchUserName(userId: string) {
  if (userNames.value[userId]) return
  const { useUserApi } = await import('@/services/api/users')
  try {
    const profile = await useUserApi().getPublicUserProfile(userId)
    const p = profile as any
    userNames.value[userId] = p.full_name || p.username || userId.slice(0, 8)
  } catch {
    userNames.value[userId] = userId.slice(0, 8)
  }
}

function userName(userId: string): string {
  return userNames.value[userId] || userId?.slice(0, 8) || ''
}

async function loadRoom() {
  try {
    const roomData = await api.getRoom(roomId)
    room.value = roomData
    if (roomData.host_id) {
      await fetchUserName(roomData.host_id)
    }
  } catch (e: any) {
    error.value = e?.message || 'Failed to load room'
  } finally {
    loading.value = false
  }
}

function togglePlay() {
  const np = queueState.value?.now_playing
  if (!np?.track?.id) return
  const trackId = np.track.id
  if (isPlayingTrack.value) {
    playerStore.pause()
  } else {
    if (playerStore.currentTrack?.id === trackId) {
      playerStore.resume()
    } else {
      playerStore.playTrackById(trackId)
    }
  }
}

function setupWebSocket() {
  wsClient.connect()
  wsClient.subscribe(`room:${roomId}`)

  wsClient.on('room.participant_joined', (msg: { payload?: { user_id: string } }) => {
    if (msg.payload?.user_id) fetchUserName(msg.payload.user_id)
  })

  wsClient.on('room.message', (msg: { payload?: { user_id: string; content: string } }) => {
    const uid = msg.payload?.user_id || ''
    if (uid && !userNames.value[uid]) fetchUserName(uid)
    messages.value.push({
      user_id: uid,
      userName: userName(uid),
      content: msg.payload?.content || '',
    })
    scrollChat()
  })

  wsClient.on('room.ended', () => {
    loadRoom()
  })
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
    await queueSocket.suggest(String(track.id))
  } catch { /* ignore */ }
}

function handleVote(candidateId: string) {
  queueSocket.vote(candidateId)
}

function handleUnvote(candidateId: string) {
  queueSocket.unvote(candidateId)
}

function handleRaiseHand() {
  stage.raiseHand()
}

function handleLowerHand() {
  stage.lowerHand()
}

function handleApproveHand(userId: string) {
  stage.approveHand(userId)
}

function handleDenyHand(userId: string) {
  stage.denyHand(userId)
}

function handleRemoveSpeaker(userId: string) {
  stage.removeFromStage(userId)
}

function handleToggleMute(payload: { userId: string; muted: boolean }) {
  stage.toggleMute(payload.userId, payload.muted)
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
  wsClient.unsubscribe(`room:${roomId}`)
})
</script>

<style scoped>
@media (prefers-reduced-motion: reduce) {
  .animate-pulse {
    animation: none !important;
  }
}
</style>
