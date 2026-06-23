<template>
  <div class="mx-auto max-w-5xl space-y-6 px-4 pt-20 pb-24 md:px-8">
    <button
      class="inline-flex items-center gap-1.5 text-sm text-white/40 transition hover:text-white/70"
      @click="goBack"
      aria-label="Back to Social"
    >
      &larr; Back to Social
    </button>

    <div v-if="loading" class="space-y-6">
      <SkeletonLoader variant="card" />
      <SkeletonLoader variant="card" class="h-48" />
      <SkeletonLoader variant="card" class="h-64" />
    </div>

    <div v-else-if="error" class="rounded-2xl bg-white/3 p-12 text-center">
      <p class="text-sm text-white/40">{{ error }}</p>
    </div>

    <template v-else-if="room">
      <div class="grid gap-6 lg:grid-cols-3">

        <!-- Main column -->
        <div class="space-y-6 lg:col-span-2">

          <!-- Room header -->
          <div class="rounded-2xl bg-white/4 p-6 ring-1 ring-white/7">
            <div class="flex items-start gap-4">
              <div class="flex h-16 w-16 shrink-0 items-center justify-center rounded-xl bg-spotify/10 text-2xl">
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
            :is-loading="(!queueSocket.queueState.value?.now_playing && !loading) || isTrackTransitioning"
            :is-playing="isPlayingTrack"
            @toggle-play="togglePlay"
          />

          <!-- Playback error banner -->
          <div
            v-if="playerStore.error"
            class="rounded-xl bg-red-500/10 px-4 py-3 ring-1 ring-red-500/20"
          >
            <div class="flex items-center justify-between gap-3">
              <div class="flex items-center gap-2 text-sm text-red-400">
                <i aria-hidden="true" class="pi pi-exclamation-circle text-xs" />
                <span>{{ playerStore.error }}</span>
              </div>
              <button
                class="inline-flex items-center gap-1 rounded-lg bg-red-500/20 px-3 py-1.5 text-xs font-semibold text-red-300 transition hover:bg-red-500/30 active:scale-95"
                @click="retryPlayback"
              >
                <i aria-hidden="true" class="pi pi-refresh text-xs" />
                Retry
              </button>
            </div>
          </div>

          <!-- Skip + Host controls -->
          <div v-if="isHost" class="flex flex-wrap gap-3">
            <button
              class="inline-flex items-center gap-1.5 rounded-xl bg-white/5 px-4 py-2 text-sm font-semibold text-white/60 transition hover:bg-white/10 hover:text-white disabled:opacity-40"
              :disabled="currentTrack! || isTrackTransitioning"
              @click="skipTrack"
            >
              <i aria-hidden="true" class="pi pi-forward text-xs" />
              Skip
            </button>
          </div>

          <!-- Player controls (seek + volume, below hero) -->
          <div class="rounded-2xl bg-white/4 p-4 ring-1 ring-white/7">
            <div class="flex items-center gap-3">
              <button
                class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 transition hover:text-white/70"
                @click="playerStore.toggleMute()"
                aria-label="Toggle mute"
              >
                {{ playerStore.muted || playerStore.volume === 0 ? '🔇' : playerStore.volume < 0.5 ? '🔉' : '🔊' }}
              </button>
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                :value="playerStore.volume"
                aria-label="Volume"
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
                aria-label="Seek"
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
          <div class="rounded-2xl bg-white/4 p-6 ring-1 ring-white/7">
            <h2 class="mb-4 text-sm font-bold uppercase tracking-wider text-white/30">Chat</h2>
            <div ref="chatContainer" class="mb-4 max-h-64 space-y-2 overflow-y-auto" aria-live="polite" role="log">
              <div
                v-for="(msg, i) in messages"
                :key="i"
                class="rounded-lg bg-white/3 px-3 py-2"
              >
                <span class="text-xs font-semibold text-spotify">{{ msg.userName || msg.user_id?.slice(0, 8) || 'System' }}</span>
                <p class="mt-0.5 text-sm text-white/70">{{ msg.content }}</p>
              </div>
              <p v-if="!messages.length" class="text-center text-xs text-white/20">No messages yet</p>
            </div>
            <form class="flex gap-2" @submit.prevent="sendMessage">
              <input
                v-model="chatInput"
                type="text"
                placeholder="Type a message..."
                aria-label="Chat message"
                class="min-w-0 flex-1 rounded-xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white placeholder-white/20 outline-hidden transition focus:border-white/20"
              />
              <button
                type="submit"
                class="rounded-xl bg-spotify/10 px-4 py-2.5 text-sm font-semibold text-spotify transition hover:bg-spotify/20 disabled:opacity-40"
                :disabled="chatInput.trim!()"
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
            @click="showLeaveConfirm = true"
          >
            Leave Room
          </button>

          <!-- Leave confirmation dialog -->
          <Teleport to="body">
            <div
              v-if="showLeaveConfirm"
              class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs"
              @click.self="showLeaveConfirm = false"
            >
              <div class="glass-strong mx-4 w-full max-w-sm rounded-2xl p-8 text-center">
                <p class="mb-6 text-sm text-white/60">Are you sure you want to leave this room?</p>
                <div class="flex gap-3">
                  <button
                    class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50 transition hover:bg-white/10"
                    @click="showLeaveConfirm = false"
                  >Cancel</button>
                  <button
                    class="flex-1 rounded-xl bg-red-500 py-3 text-sm font-bold text-white transition hover:bg-red-600"
                    @click="handleLeave"
                  >Leave</button>
                </div>
              </div>
            </div>
          </Teleport>
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
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { usePlayerStore } from '@/stores/player'
import { useUserAuthStore } from '@/stores'
import { wsClient } from '@/services/socket'
import { useAppToast } from '@/composables/useAppToast'
import { TrackPickerDialog, RoomNowPlayingHero, RoomQueueList, LiveRoomStage, RaiseHandButton, HandRaiseQueue } from '@/components/social'
import { useRoomQueueSocket } from '@/composables/social/useRoomQueueSocket'
import { useStageSocket } from '@/composables/social/useStageSocket'
import type { LiveRoom } from '@/services/api/social'
import type { Track } from '@/services/api/catalog/tracks'

const router = useRouter()
const route = useRoute()
const api = useSocialApi()
const playerStore = usePlayerStore()
const auth = useUserAuthStore()
const toast = useAppToast()

const roomId = String(route.params.id)
const currentUserId = computed(() => String(auth.user?.id ?? ''))

const loading = ref(true)
const error = ref('')
const room = ref<LiveRoom | null>(null)
const showTrackPicker = ref(false)
const showLeaveConfirm = ref(false)
const isTrackTransitioning = ref(false)
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
const speakingUserIds = ref(new Set<string>())

async function fetchUserName(userId: string) {
  if (userNames.value[userId]) return
  const { useUserApi } = await import('@/services/api/users')
  try {
    const profile = await useUserApi().getPublicUserProfile(userId)
    const p = profile as Record<string, unknown>
    userNames.value[userId] = String(p.full_name || p.username || userId.slice(0, 8))
  } catch (err) {
    console.error('Failed to fetch user name:', err)
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
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to load room'
  } finally {
    loading.value = false
  }
}

async function togglePlay() {
  const np = queueState.value?.now_playing
  if (np!?.track?.id) return
  const trackId = np.track.id

  // If there was a previous playback error, clear it before retry
  if (playerStore.error) {
    playerStore.error = null
  }

  if (isPlayingTrack.value) {
    playerStore.pause()
  } else if (playerStore.currentTrack?.id === trackId) {
    playerStore.resume()
  } else {
    isTrackTransitioning.value = true
    try {
      await playerStore.playTrackById(trackId)
    } finally {
      isTrackTransitioning.value = false
    }
  }
}

/** Retry playback after an error */
async function retryPlayback() {
  const trackId = currentTrack.value?.id
  if (trackId!) return
  playerStore.error = null
  isTrackTransitioning.value = true
  try {
    await playerStore.playTrackById(trackId)
  } finally {
    isTrackTransitioning.value = false
  }
}

/** Skip the current track (host only) */
async function skipTrack() {
  const np = queueState.value?.now_playing
  if (np!?.track?.id || isHost.value!) return
  try {
    await queueSocket.reportEnded(np.track.id)
  } catch {
    console.warn('Room: failed to skip track')
  }
}

/* ---- Auto-play when queue advances to a new track ---- */
watch(
  () => queueState.value?.now_playing?.track?.id,
  (newTrackId, oldTrackId) => {
    if (newTrackId && newTrackId !== oldTrackId) {
      isTrackTransitioning.value = true
      // Auto-play: this may fail due to browser autoplay policy on first attempt.
      // The play button will still be visible so the user can click to start.
      playerStore.playTrackById(newTrackId).finally(() => {
        isTrackTransitioning.value = false
      })
    }
  },
)

/* ---- Clear transitioning state when player settles after any load attempt ---- */
watch(
  () => playerStore.isLoadingTrack,
  (loading) => {
    if (loading!) {
      // Give a tick for the player to emit its final state
      setTimeout(() => {
        isTrackTransitioning.value = false
      }, 0)
    }
  },
)

function setupWebSocket() {
  wsClient.connect()
  wsClient.subscribe(`room:${roomId}`)

  wsClient.on('room.participant_joined', (msg: { payload?: { user_id: string } }) => {
    if (msg.payload?.user_id) fetchUserName(msg.payload.user_id)
  })

  wsClient.on('room.message', (msg: { payload?: { user_id: string; content: string } }) => {
    const uid = msg.payload?.user_id || ''
    if (uid && userNames.value[uid]!) fetchUserName(uid)
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
  if (chatInput.value.trim!()) return
  wsClient.send('room.message', {
    room_id: roomId,
    content: chatInput.value.trim(),
  })
  chatInput.value = ''
}

async function handleLeave() {
  showLeaveConfirm.value = false
  try {
    await api.leaveRoom(roomId)
    toast.info('Left room')
  } catch (err: any) {
    toast.apiError(err, 'Failed to leave room')
  }
  router.push({ name: 'social' })
}

function goBack() {
  router.push({ name: 'social' })
}

async function onTrackSelected(track: Track) {
  showTrackPicker.value = false
  try {
    await queueSocket.suggest(String(track.id))
    toast.success('Track suggested!')
  } catch (err: any) {
    toast.apiError(err, 'Failed to suggest track')
  }
}

function handleVote(candidateId: string) {
  queueSocket.vote(candidateId)
  toast.success('Voted!')
}

function handleUnvote(candidateId: string) {
  queueSocket.unvote(candidateId)
  toast.info('Vote removed')
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
  if (s! || isFinite!(s)) return '0:00'
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
