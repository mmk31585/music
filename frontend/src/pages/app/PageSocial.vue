<template>
  <div class="mx-auto max-w-5xl space-y-8 px-4 pt-20 pb-24 md:px-8">
    <!-- Hero -->
    <div class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-[#1db954]/20 via-transparent to-purple-500/10 p-8 text-center">
      <div class="relative">
        <h1 class="text-4xl font-black text-white md:text-5xl">Social</h1>
        <p class="mt-2 text-sm text-white/40 max-w-md mx-auto">
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
        class="flex items-center gap-1.5 rounded-lg bg-white/5 px-3 py-2 text-xs text-white/40 transition hover:bg-white/10 hover:text-white/60"
        @click="loadData"
        :disabled="loading"
      >
        <span :class="loading ? 'animate-spin' : ''">⟳</span>
        Refresh
      </button>
      <button
        class="rounded-xl bg-[#1db954]/10 px-5 py-2.5 text-sm font-semibold text-[#1db954] transition hover:bg-[#1db954]/20"
        @click="showCreateModal = true"
      >
        + Create {{ activeTab === 'parties' ? 'Party' : activeTab === 'rooms' ? 'Room' : 'Club' }}
      </button>
    </div>

    <!-- Loading spinner -->
    <div v-if="loading" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
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
          class="col-span-full rounded-2xl border border-dashed border-white/5 py-16 text-center"
        >
          <p class="text-3xl mb-2">🎉</p>
          <p class="text-sm text-white/30">No active listening parties. Create one!</p>
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
          class="col-span-full rounded-2xl border border-dashed border-white/5 py-16 text-center"
        >
          <p class="text-3xl mb-2">🎤</p>
          <p class="text-sm text-white/30">No live rooms right now. Start one!</p>
        </div>
      </div>

      <!-- Tab: Music Clubs -->
      <div v-show="activeTab === 'clubs'" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <MusicClubCard
          v-for="club in clubs"
          :key="club.id"
          :club="club"
          :creator-name="userNames[club.created_by]"
          @join="handleJoinClub"
        />
        <div
          v-if="!clubs.length"
          class="col-span-full rounded-2xl border border-dashed border-white/5 py-16 text-center"
        >
          <p class="text-3xl mb-2">🏛</p>
          <p class="text-sm text-white/30">No music clubs yet. Create one!</p>
        </div>
      </div>

      <!-- Tab: Discussions -->
      <div v-show="activeTab === 'discussions'" class="glass-strong rounded-2xl p-6 md:p-8">
        <div class="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center">
          <input
            v-model="discussionFilter.target_type"
            type="text"
            placeholder="Type (track, album, playlist, artist)"
            class="flex-1 rounded-lg border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white placeholder-white/20 outline-none focus:border-white/20"
          />
          <input
            v-model="discussionFilter.target_id"
            type="text"
            placeholder="ID"
            class="flex-1 rounded-lg border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white placeholder-white/20 outline-none focus:border-white/20"
          />
          <button
            class="rounded-lg bg-[#1db954]/10 px-4 py-2.5 text-sm font-medium text-[#1db954] transition hover:bg-[#1db954]/20"
            @click="loadDiscussions"
          >
            Browse
          </button>
        </div>
        <DiscussionThread :discussions="discussions" :user-names="userNames" @create="handleCreateDiscussion" />
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
              class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white placeholder-white/20 transition outline-none focus:border-white/20"
            />
            <textarea
              v-model="createForm.description"
              placeholder="Description (optional)"
              rows="3"
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
              :disabled="!createForm.title.trim()"
              @click="handleCreate"
            >Create</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
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
const api = useSocialApi()
const playerApi = usePlayerApi()

const loading = ref(true)
const activeTab = ref('parties')
const showCreateModal = ref(false)

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
    const p = profile as any
    userNames.value[userId] = p.full_name || p.username || userId.slice(0, 8)
  } catch {
    userNames.value[userId] = userId.slice(0, 8)
  }
}

async function fetchTrackName(trackId: string) {
  if (!trackId || trackNames.value[trackId]) return
  try {
    const track = await playerApi.getPlaybackTrack(trackId)
    trackNames.value[trackId] = `${track.artistName} - ${track.title}`
  } catch {
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
  if (!discussionFilter.target_type || !discussionFilter.target_id) return
  loading.value = true
  try {
    const res = await api.getDiscussions({
      target_type: discussionFilter.target_type,
      target_id: discussionFilter.target_id,
      limit: 50,
    })
    discussions.value = (res || []).filter(Boolean)
    const userIds = new Set<string>()
    discussions.value.forEach((d: Discussion) => userIds.add(d.user_id))
    await Promise.all(Array.from(userIds).map(fetchUserName))
  } catch {
    discussions.value = []
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  try {
    if (activeTab.value === 'parties') {
      await api.createParty({
        title: createForm.title,
        description: createForm.description,
        is_public: createForm.is_public,
      })
      loadData()
    } else if (activeTab.value === 'rooms') {
      await api.createRoom({
        title: createForm.title,
        description: createForm.description,
        is_public: createForm.is_public,
      })
      loadData()
    } else if (activeTab.value === 'clubs') {
      await api.createClub({
        name: createForm.title,
        description: createForm.description,
        is_public: createForm.is_public,
      })
      loadData()
    }
    showCreateModal.value = false
    createForm.title = ''
    createForm.description = ''
    createForm.is_public = true
  } catch { /* silent */ }
}

async function handleJoinParty(id: string) {
  try {
    await api.joinParty(id)
    router.push({ name: 'social.party', params: { id } })
  } catch { /* silent */ }
}

async function handleJoinRoom(id: string) {
  try {
    await api.joinRoom(id)
    router.push({ name: 'social.room', params: { id } })
  } catch { /* silent */ }
}

async function handleJoinClub(id: string) {
  try {
    await api.joinClub(id)
    router.push({ name: 'social.club', params: { id } })
  } catch { /* silent */ }
}

async function handleCreateDiscussion(content: string, parentId?: string) {
  try {
    await api.createDiscussion({
      target_type: discussionFilter.target_type,
      target_id: discussionFilter.target_id,
      content,
      parent_id: parentId,
    })
    loadDiscussions()
  } catch { /* silent */ }
}

onMounted(loadData)
</script>
