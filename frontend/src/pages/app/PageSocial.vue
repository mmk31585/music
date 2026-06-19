<template>
  <div class="mx-auto max-w-5xl space-y-8 px-4 pt-20 pb-24 md:px-8">
    <!-- Hero -->
    <div class="relative overflow-hidden rounded-[2rem] border border-white/[0.06] bg-[#0C0C14] p-10 text-center">
      <div class="absolute -top-20 -right-20 h-60 w-60 rounded-full bg-[#1db954]/10 blur-3xl" />
      <div class="relative">
        <p class="text-[10px] font-bold tracking-[0.35em] text-white/30 uppercase">Community</p>
        <h1 class="mt-2 text-4xl font-black text-white md:text-5xl">Social</h1>
        <p class="mt-3 mx-auto max-w-md text-sm text-white/50">
          Listening parties, live rooms, music clubs, and discussions.
        </p>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex gap-1 rounded-xl bg-white/[0.04] p-1">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="flex-1 rounded-lg py-2.5 text-sm font-medium transition-all duration-200"
        :class="activeTab === tab.key
          ? 'bg-white/10 text-white shadow-lg'
          : 'text-white/30 hover:text-white/50'"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Create + Refresh row -->
    <div v-if="activeTab !== 'discussions'" class="flex items-center justify-end gap-3">
      <button
        class="inline-flex items-center gap-1.5 rounded-lg border border-white/10 bg-white/[0.03] px-3 py-2 text-xs font-medium text-white/40 transition hover:border-white/20 hover:bg-white/[0.06] hover:text-white/60"
        @click="loadData"
        :disabled="loading"
        aria-label="Refresh"
      >
        <span :class="loading ? 'animate-spin' : ''">⟳</span>
        Refresh
      </button>
      <button
        class="inline-flex items-center gap-1.5 rounded-xl bg-[#1db954] px-5 py-2.5 text-sm font-bold text-black transition hover:bg-[#1ed760]"
        @click="showCreateModal = true"
        :aria-label="'Create ' + (activeTab === 'parties' ? 'Party' : activeTab === 'rooms' ? 'Room' : 'Club')"
      >
        <i aria-hidden="true" class="pi pi-plus text-xs" />
        {{ activeTab === 'parties' ? 'Party' : activeTab === 'rooms' ? 'Room' : 'Club' }}
      </button>
    </div>

    <!-- Loading spinner -->
    <div v-if="loading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3" role="status">
      <SkeletonLoader v-for="i in 6" :key="i" variant="card" />
    </div>

    <template v-else>
      <!-- Tab: Listening Parties -->
      <div v-show="activeTab === 'parties'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <ListeningPartyCard
          v-for="party in parties"
          :key="party.id"
          :party="party"
          :host-name="userNames[party.host_id]"
          :current-track-name="trackNames[party.current_track_id || '']"
          @join="handleJoinParty"
        />
        <div
          v-if="!parties.length"
          class="col-span-full flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.06] py-16 text-center"
          role="status"
        >
          <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/[0.04]">
            <i aria-hidden="true" class="pi pi-users text-xl text-slate-500" />
          </div>
          <p class="text-sm font-medium text-white/60">No active listening parties. Create one!</p>
        </div>
      </div>

      <!-- Tab: Live Rooms -->
      <div v-show="activeTab === 'rooms'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <LiveRoomCard
          v-for="room in rooms"
          :key="room.id"
          :room="room"
          :host-name="userNames[room.host_id]"
          @join="handleJoinRoom"
        />
        <div
          v-if="!rooms.length"
          class="col-span-full flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.06] py-16 text-center"
          role="status"
        >
          <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/[0.04]">
            <i aria-hidden="true" class="pi pi-megaphone text-xl text-slate-500" />
          </div>
          <p class="text-sm font-medium text-white/60">No live rooms right now. Start one!</p>
        </div>
      </div>

      <!-- Tab: Music Clubs -->
      <div v-show="activeTab === 'clubs'" class="space-y-6">
        <div class="flex items-center justify-end">
          <button
            class="inline-flex items-center gap-1.5 rounded-lg border border-white/10 bg-white/[0.03] px-3 py-2 text-xs font-medium text-white/40 transition hover:border-white/20 hover:bg-white/[0.06] hover:text-white/60"
            @click="router.push({ name: 'social.clubs.browse' })"
          >
            مرور همه کلاب‌ها
          </button>
        </div>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <MusicClubCard
            v-for="club in clubs"
            :key="club.id"
            :club="club"
            :creator-name="userNames[club.created_by]"
            @join="handleJoinClub"
          />
          <div
            v-if="!clubs.length"
            class="col-span-full flex flex-col items-center gap-3 rounded-2xl border border-dashed border-white/[0.06] py-16 text-center"
            role="status"
          >
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-white/[0.04]">
              <i aria-hidden="true" class="pi pi-building text-xl text-slate-500" />
            </div>
            <p class="text-sm font-medium text-white/60">No music clubs yet. Create one!</p>
          </div>
        </div>
      </div>

      <!-- Tab: Discussions -->
      <div v-show="activeTab === 'discussions'" class="glass-strong rounded-2xl p-6 md:p-8" aria-live="polite">
        <div class="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center">
          <input
            v-model="discussionFilter.target_type"
            type="text"
            placeholder="Type (track, album, playlist, artist)"
            aria-label="Discussion type filter"
            class="flex-1 rounded-lg border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white placeholder-white/20 outline-none focus:border-white/20"
          />
          <input
            v-model="discussionFilter.target_id"
            type="text"
            placeholder="ID"
            aria-label="Discussion ID filter"
            class="flex-1 rounded-lg border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white placeholder-white/20 outline-none focus:border-white/20"
          />
          <button
            class="inline-flex items-center gap-1.5 rounded-lg bg-[#1db954]/10 px-4 py-2.5 text-sm font-medium text-[#1db954] transition hover:bg-[#1db954]/20 disabled:opacity-40"
            :disabled="discussionsLoading"
            @click="loadDiscussions"
          >
            <span v-if="discussionsLoading" class="inline-block h-3.5 w-3.5 animate-spin rounded-full border-2 border-[#1db954] border-t-transparent" />
            <template v-else>Browse</template>
          </button>
        </div>
        <div v-if="discussionsError" class="mb-4 rounded-xl bg-red-500/10 px-4 py-3 text-sm text-red-400">
          {{ discussionsError }}
        </div>
        <DiscussionThread
          :discussions="discussions"
          :user-names="userNames"
          :is-posting="discussionPosting"
          @create="handleCreateDiscussion"
        />
      </div>
    </template>

    <!-- Create Modal -->
    <Teleport to="body">
      <div
        v-if="showCreateModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
        @click.self="showCreateModal = false"
      >
        <div class="glass-strong mx-4 w-full max-w-md rounded-2xl p-8">
          <h2 class="mb-6 text-xl font-bold text-white">
            Create {{
              activeTab === 'parties' ? 'Listening Party'
                : activeTab === 'rooms' ? 'Live Room'
                : 'Music Club'
            }}
          </h2>
          <div class="space-y-4">
            <input
              v-model="createForm.title"
              type="text"
              :placeholder="activeTab === 'clubs' ? 'Club name' : 'Title'"
              aria-label="Title"
              class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 transition outline-none focus:border-white/20"
            />
            <textarea
              v-model="createForm.description"
              placeholder="Description (optional)"
              rows="3"
              aria-label="Description"
              class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 transition outline-none focus:border-white/20"
            />
            <label class="flex items-center gap-3">
              <input
                v-model="createForm.is_public"
                type="checkbox"
                class="h-5 w-5 rounded border-white/10 bg-white/5 accent-[#1db954]"
              />
              <span class="text-sm text-white/60">Public</span>
            </label>
          </div>
          <div class="mt-6 flex gap-3">
            <button
              class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50 transition hover:bg-white/10"
              @click="showCreateModal = false"
            >Cancel</button>
            <button
              class="flex-1 rounded-xl bg-[#1db954] py-3 text-sm font-bold text-black transition hover:bg-[#1db954]/90 disabled:opacity-40"
              :disabled="!createForm.title.trim() || isCreating"
              @click="handleCreate"
            >
              <span v-if="isCreating" class="inline-flex items-center gap-2">
                <span class="inline-block h-4 w-4 animate-spin rounded-full border-2 border-black border-t-transparent" />
                Creating...
              </span>
              <span v-else>Create</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAppToast } from '@/composables/useAppToast'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { usePlayerApi } from '@/services/api/player'
import { useUserApi } from '@/services/api/users'
import type { ListeningParty, LiveRoom, MusicClub, Discussion } from '@/services/api/social'
import ListeningPartyCard from '@/components/social/ListeningPartyCard.vue'
import LiveRoomCard from '@/components/social/LiveRoomCard.vue'
import MusicClubCard from '@/components/social/MusicClubCard.vue'
import DiscussionThread from '@/components/social/DiscussionThread.vue'

const router = useRouter()
const toast = useAppToast()
const api = useSocialApi()
const playerApi = usePlayerApi()

const loading = ref(true)
const activeTab = ref('parties')
const showCreateModal = ref(false)
const isCreating = ref(false)
const discussionsLoading = ref(false)
const discussionsError = ref('')
const discussionPosting = ref(false)

const tabs = [
  { key: 'parties', label: 'Listening Parties' },
  { key: 'rooms', label: 'Live Rooms' },
  { key: 'clubs', label: 'Music Clubs' },
  { key: 'discussions', label: 'Discussions' },
]

const parties = ref<ListeningParty[]>([])
const rooms = ref<LiveRoom[]>([])
const clubs = ref<MusicClub[]>([])
const discussions = ref<Discussion[]>([])
const userNames = ref<Record<string, string>>({})
const trackNames = ref<Record<string, string>>({})

const discussionFilter = reactive({ target_type: '', target_id: '' })

const createForm = reactive({
  title: '',
  description: '',
  is_public: true,
})

async function fetchUserName(userId: string) {
  if (userNames.value[userId]) return
  try {
    const profile = await useUserApi().getPublicUserProfile(userId)
    const p = profile as Record<string, unknown>
    userNames.value[userId] = String(p.full_name || p.username || userId.slice(0, 8))
  } catch (err) {
    console.error('Failed to fetch user name:', err)
    userNames.value[userId] = userId.slice(0, 8)
  }
}

async function fetchTrackName(trackId: string) {
  if (!trackId || trackNames.value[trackId]) return
  try {
    const track = await playerApi.getPlaybackTrack(trackId)
    trackNames.value[trackId] = `${track.artistName} - ${track.title}`
  } catch (err) {
    console.error('Failed to fetch track name:', err)
    trackNames.value[trackId] = trackId.slice(0, 12)
  }
}

async function loadData() {
  loading.value = true
  const [partiesRes, roomsRes, clubsRes] = await Promise.allSettled([
    api.listParties({ limit: 50 }),
    api.listRooms({ limit: 50 }),
    api.listClubs({ limit: 50 }),
  ])
  const rawParties = partiesRes.status === 'fulfilled' ? partiesRes.value : null
  const rawRooms = roomsRes.status === 'fulfilled' ? roomsRes.value : null
  const rawClubs = clubsRes.status === 'fulfilled' ? clubsRes.value : null
  parties.value = Array.isArray(rawParties) ? rawParties.filter(Boolean) : []
  rooms.value = Array.isArray(rawRooms) ? rawRooms.filter(Boolean) : []
  clubs.value = Array.isArray(rawClubs) ? rawClubs.filter(Boolean) : []

  const userIds = new Set<string>()
  if (Array.isArray(rawParties)) rawParties.forEach((p: ListeningParty) => userIds.add(p.host_id))
  if (Array.isArray(rawRooms)) rawRooms.forEach((r: LiveRoom) => userIds.add(r.host_id))
  if (Array.isArray(rawClubs)) rawClubs.forEach((c: MusicClub) => userIds.add(c.created_by))
  await Promise.all(Array.from(userIds).map(fetchUserName))

  const trackIds = new Set<string>()
  if (Array.isArray(rawParties)) {
    rawParties.forEach((p: ListeningParty) => {
      if (p.current_track_id) trackIds.add(p.current_track_id)
    })
  }
  await Promise.all(Array.from(trackIds).map(fetchTrackName))

  loading.value = false
}

async function loadDiscussions() {
  if (!discussionFilter.target_type || !discussionFilter.target_id) {
    toast.warn('Please enter both type and ID to browse discussions')
    return
  }
  discussionsLoading.value = true
  discussionsError.value = ''
  try {
    const res = await api.getDiscussions({
      target_type: discussionFilter.target_type,
      target_id: discussionFilter.target_id,
      limit: 50,
    })
    discussions.value = (res || []).filter(Boolean)
    if (!discussions.value.length) {
      toast.info('No discussions found for this filter')
    }
    const userIds = new Set<string>()
    discussions.value.forEach((d: Discussion) => userIds.add(d.user_id))
    await Promise.all(Array.from(userIds).map(fetchUserName))
  } catch (err: any) {
    discussionsError.value = err?.response?.data?.message || err.message || 'Failed to load discussions'
    discussions.value = []
    toast.error(discussionsError.value)
  } finally {
    discussionsLoading.value = false
  }
}

async function handleCreate() {
  if (!createForm.title.trim() || isCreating.value) return
  isCreating.value = true
  try {
    if (activeTab.value === 'parties') {
      await api.createParty({
        title: createForm.title,
        description: createForm.description,
        is_public: createForm.is_public,
      })
      toast.success('Listening party created!')
      loadData()
    } else if (activeTab.value === 'rooms') {
      await api.createRoom({
        title: createForm.title,
        description: createForm.description,
        is_public: createForm.is_public,
      })
      toast.success('Live room created!')
      loadData()
    } else if (activeTab.value === 'clubs') {
      await api.createClub({
        name: createForm.title,
        description: createForm.description,
        is_public: createForm.is_public,
      })
      toast.success('Music club created!')
      loadData()
    }
    showCreateModal.value = false
    createForm.title = ''
    createForm.description = ''
    createForm.is_public = true
  } catch (err: any) {
    toast.apiError(err, 'Failed to create')
  } finally {
    isCreating.value = false
  }
}

async function handleJoinParty(id: string) {
  try {
    await api.joinParty(id)
    toast.success('Joined party!')
    router.push({ name: 'social.party', params: { id } })
  } catch (err: any) {
    toast.apiError(err, 'Failed to join party')
  }
}

async function handleJoinRoom(id: string) {
  try {
    await api.joinRoom(id)
    toast.success('Joined room!')
    router.push({ name: 'social.room', params: { id } })
  } catch (err: any) {
    toast.apiError(err, 'Failed to join room')
  }
}

async function handleJoinClub(id: string) {
  try {
    await api.joinClub(id)
    toast.success('Joined club!')
    router.push({ name: 'social.club', params: { id } })
  } catch (err: any) {
    toast.apiError(err, 'Failed to join club')
  }
}

async function handleCreateDiscussion(content: string, parentId?: string) {
  discussionPosting.value = true
  try {
    await api.createDiscussion({
      target_type: discussionFilter.target_type,
      target_id: discussionFilter.target_id,
      content,
      parent_id: parentId,
    })
    toast.success(parentId ? 'Reply posted!' : 'Comment posted!')
    loadDiscussions()
  } catch (err: any) {
    toast.apiError(err, 'Failed to post')
  } finally {
    discussionPosting.value = false
  }
}

onMounted(loadData)
</script>
