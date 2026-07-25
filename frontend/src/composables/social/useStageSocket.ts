import { ref, readonly, onMounted, onUnmounted } from 'vue'
import { wsClient } from '@/services/socket'
import { useUserApi } from '@/services/api/users'
import { useAppToast } from '@/composables/useAppToast'
import {
  getStageState,
  raiseHand,
  lowerHand,
  approveHand,
  denyHand,
  removeFromStage,
  leaveStage,
  toggleMute,
} from '@/services/api/social/stage'
import type {
  StageState,
  StageMember,
  HandRaiseData,
  StageStateResponse,
} from '@/services/api/social/stage'

export function useStageSocket(roomId: string, currentUserId: string) {
  const stageState = ref<StageState | null>(null)
  const toast = useAppToast()

  let unsubHandRaised: (() => void) | null = null
  let unsubStageUpdated: (() => void) | null = null
  let unsubHandApproved: (() => void) | null = null
  let unsubHandDenied: (() => void) | null = null

  // Cache resolved usernames
  const userNames = ref<Record<string, { username: string; avatar_url?: string }>>({})

  async function resolveUser(userId: string): Promise<{ username: string; avatar_url?: string }> {
    if (userNames.value[userId]) return userNames.value[userId]
    try {
      const profile = await useUserApi().getPublicUserProfile(userId) as Record<string, any>
      const username = String(profile.full_name || profile.username || userId.slice(0, 8))
      const avatar_url = profile.avatar_url as string | undefined
      userNames.value[userId] = { username, avatar_url }
      return { username, avatar_url }
    } catch {
      const fallback = { username: userId.slice(0, 8) }
      userNames.value[userId] = fallback
      return fallback
    }
  }

  async function resolveUsers(ids: string[]) {
    const unique = [...new Set(ids)]
    await Promise.all(unique.map(resolveUser))
  }

  function enrichState(raw: StageStateResponse, myRole: 'host' | 'speaker' | 'listener', handRaised: boolean): StageState {
    const host: StageMember & { username?: string; avatar_url?: string } = {
      ...raw.host,
      username: userNames.value[raw.host.user_id]?.username || raw.host.user_id.slice(0, 8),
      avatar_url: userNames.value[raw.host.user_id]?.avatar_url,
    }
    const speakers = raw.speakers.map((s) => ({
      ...s,
      username: userNames.value[s.user_id]?.username || s.user_id.slice(0, 8),
      avatar_url: userNames.value[s.user_id]?.avatar_url,
    }))
    const pending_hand_raises = raw.pending_requests.map((hr: HandRaiseData) => ({
      ...hr,
      username: userNames.value[hr.user_id]?.username || hr.user_id.slice(0, 8),
      avatar_url: userNames.value[hr.user_id]?.avatar_url,
    }))

    return {
      host,
      speakers,
      pending_hand_raises,
      my_role: myRole,
      my_hand_raised: handRaised,
    }
  }

  async function fetchAndSet() {
    try {
      const raw = await getStageState(roomId)
      // Resolve all user IDs in the response
      const userIds = [raw.host.user_id]
      raw.speakers.forEach((s) => userIds.push(s.user_id))
      raw.pending_requests.forEach((hr) => userIds.push(hr.user_id))
      await resolveUsers(userIds)

      // Determine my role
      let myRole: 'host' | 'speaker' | 'listener' = 'listener'
      if (raw.host.user_id === currentUserId) {
        myRole = 'host'
      } else if (raw.speakers.some((s) => s.user_id === currentUserId)) {
        myRole = 'speaker'
      }

      const myHandRaised = raw.pending_requests.some((hr) => hr.user_id === currentUserId)

      stageState.value = enrichState(raw, myRole, myHandRaised)
    } catch {
      // keep previous state
    }
  }

  function setupSocket() {
    wsClient.connect()
    wsClient.subscribe(`room:${roomId}`)

    unsubHandRaised = wsClient.on('room.hand_raised', async (msg: { payload?: { user_id: string } }) => {
      const uid = msg.payload?.user_id
      if (uid) {
        const info = await resolveUser(uid)
        stageState.value?.pending_hand_raises.push({
          id: '',
          room_id: roomId,
          user_id: uid,
          status: 'pending',
          created_at: new Date().toISOString(),
          username: info.username,
          avatar_url: info.avatar_url,
        })
        // Only show toast to host
        if (currentUserId && stageState.value?.my_role === 'host') {
          toast.info(`${info.username} دستش رو بالا برد`)
        }
      }
    })

    unsubStageUpdated = wsClient.on('room.stage_updated', () => {
      fetchAndSet()
    })

    unsubHandApproved = wsClient.on('room.hand_approved', (msg: { payload?: { user_id: string } }) => {
      if (msg.payload?.user_id === currentUserId) {
        toast.success('میزبان قبول کرد! می‌تونی صحبت کنی')
      }
      fetchAndSet()
    })

    unsubHandDenied = wsClient.on('room.hand_denied', (msg: { payload?: { user_id: string } }) => {
      if (msg.payload?.user_id === currentUserId) {
        toast.info('میزبان درخواستت رو رد کرد')
      }
      fetchAndSet()
    })
  }

  function teardownSocket() {
    if (unsubHandRaised) { unsubHandRaised(); unsubHandRaised = null }
    if (unsubStageUpdated) { unsubStageUpdated(); unsubStageUpdated = null }
    if (unsubHandApproved) { unsubHandApproved(); unsubHandApproved = null }
    if (unsubHandDenied) { unsubHandDenied(); unsubHandDenied = null }
    wsClient.unsubscribe(`room:${roomId}`)
  }

  onMounted(() => {
    fetchAndSet()
    setupSocket()
  })

  onUnmounted(() => {
    teardownSocket()
  })

  return {
    stageState: readonly(stageState),
    refresh: fetchAndSet,
    raiseHand: () => raiseHand(roomId),
    lowerHand: () => lowerHand(roomId),
    approveHand: (userId: string) => approveHand(roomId, userId),
    denyHand: (userId: string) => denyHand(roomId, userId),
    removeFromStage: (userId: string) => removeFromStage(roomId, userId),
    leaveStage: () => leaveStage(roomId),
    toggleMute: (userId: string, muted: boolean) => toggleMute(roomId, userId, muted),
  }
}
