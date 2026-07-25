<template>
  <div class="mx-auto max-w-7xl space-y-8 px-4 pt-4 pb-32 md:px-6 lg:px-8">
    <!-- ── Dynamic Hero ── -->
    <div class="relative overflow-hidden rounded-2xl p-8 md:p-10">
      <div class="pointer-events-none absolute inset-0" aria-hidden="true">
        <div class="absolute -top-40 -right-40 h-125 w-125 rounded-full bg-spotify/8 blur-3xl" />
        <div class="absolute -bottom-20 -left-20 h-75 w-75 rounded-full bg-aurora-purple/6 blur-3xl" />
        <div class="absolute inset-0 bg-linear-to-br from-black/40 via-transparent to-black/60" />
      </div>

      <div class="relative z-10 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <p class="text-[10px] font-bold tracking-[0.35em] text-white/30 uppercase">Social Hub</p>
          <h1 class="mt-1 text-3xl font-black text-white md:text-4xl">
            {{ greeting }}
          </h1>
          <p class="mt-1 text-sm text-white/50">با دوستات گوش بده — Listen with friends</p>
        </div>

        <!-- Inline "Start a Party from current track" -->
        <div
          v-if="currentTrack"
          class="flex items-center gap-3 rounded-2xl bg-white/6 p-3 ring-1 ring-white/10"
        >
          <div class="h-10 w-10 shrink-0 overflow-hidden rounded-xl">
            <img
              v-if="currentTrack.coverUrl"
              :src="currentTrack.coverUrl"
              alt=""
              class="h-full w-full object-cover"
            />
            <div v-else class="flex h-full w-full items-center justify-center bg-white/5">
              <Music aria-hidden="true" class="text-sm text-white/30"  />
            </div>
          </div>
          <div class="min-w-0 max-w-45">
            <p class="truncate text-xs font-medium text-white">{{ currentTrack.title }}</p>
            <p class="truncate text-[10px] text-white/40">Now Playing</p>
          </div>
          <button
            aria-label="Start a listening party with this track"
            class="rounded-full bg-spotify px-4 py-2 text-[11px] font-bold text-black transition hover:bg-spotify-hover hover:scale-105 active:scale-95 focus-visible:outline-2 focus-visible:outline-white"
            @click="startPartyFromTrack"
          >
            <span class="flex items-center gap-1.5">
              <Users aria-hidden="true" class="text-[10px]"  />
              Party
            </span>
          </button>
        </div>
      </div>
    </div>

    <!-- loading state -->
    <template v-if="loading">
      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <div v-for="i in 3" :key="i" class="rounded-2xl bg-white/4 p-5 ring-1 ring-white/6">
          <div class="h-4 w-20 rounded-full bg-white/5 animate-pulse mb-4" />
          <div class="h-24 rounded-xl bg-white/5 animate-pulse mb-3" />
          <div class="h-8 w-full rounded-lg bg-white/5 animate-pulse" />
        </div>
      </div>
      <div class="space-y-2">
        <div v-for="i in 4" :key="i" class="h-14 rounded-xl bg-white/5 animate-pulse" />
      </div>
    </template>

    <!-- loaded content -->
    <template v-else>
      <!-- ── Bento Grid: Live Now / Trending Parties / Your Clubs ── -->
      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <!-- Live Now -->
        <div class="rounded-2xl bg-white/6 p-5 ring-1 ring-white/10 transition hover:bg-white/8 group">
          <div class="flex items-center gap-2 mb-3">
            <span class="flex h-2 w-2 rounded-full bg-red-500 motion-safe:animate-pulse" />
            <span class="text-[10px] font-bold tracking-wider text-red-400 uppercase">Live Now</span>
          </div>

          <!-- Waveform stub -->
          <div class="flex items-end gap-0.5 h-12 mb-3" aria-hidden="true">
            <div
              v-for="i in 16" :key="i"
              class="w-1.5 rounded-full bg-spotify/60"
              :style="{ height: `${12 + Math.random() * 36}px`, animationDelay: `${i * 100}ms` }"
            />
          </div>

          <template v-if="featuredRoom">
            <h3 class="text-lg font-bold text-white">{{ featuredRoom.title }}</h3>
            <p class="mt-1 text-xs text-white/40 line-clamp-2">{{ featuredRoom.description }}</p>
            <div class="mt-3 flex items-center gap-2">
              <div class="flex -space-x-2" aria-label="Listeners">
                <div
                  v-for="i in Math.min(5, featuredRoom.listener_count || 0)" :key="i"
                  class="h-7 w-7 overflow-hidden rounded-full border-2 border-surface-base bg-white/10"
                />
                <div
                  v-if="(featuredRoom.listener_count || 0) > 5"
                  class="flex h-7 w-7 items-center justify-center rounded-full border-2 border-surface-base bg-white/10 text-[9px] font-bold text-white/60"
                >
                  +{{ (featuredRoom.listener_count || 0) - 5 }}
                </div>
              </div>
              <span class="text-[11px] text-white/40 tabular-nums">{{ featuredRoom.listener_count || 0 }} listening</span>
            </div>
            <button
              aria-label="Join live room"
              class="mt-4 w-full rounded-xl bg-red-500/10 py-2.5 text-xs font-bold text-red-400 transition hover:bg-red-500/20 active:scale-[0.98] focus-visible:outline-2 focus-visible:outline-red-400"
              @click="handleJoinRoom(featuredRoom.id)"
            >
              Listen Live
            </button>
          </template>
          <div v-else class="flex flex-col items-center gap-3 py-6 text-center">
            <Megaphone aria-hidden="true" class="text-xl text-white/20"  />
            <p class="text-xs text-white/30">No live rooms right now</p>
            <button
              class="rounded-full bg-red-500/10 px-4 py-1.5 text-[10px] font-semibold text-red-400 transition hover:bg-red-500/20"
              @click="openCreateWizard('room')"
            >
              Go Live
            </button>
          </div>
        </div>

        <!-- Trending Parties -->
        <div class="rounded-2xl bg-white/6 p-5 ring-1 ring-white/10">
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-sm font-bold text-white">Trending Parties</h3>
            <button
              v-if="parties.length > 3"
              class="text-[10px] text-white/30 hover:text-white/60 transition"
              @click="activeTab = 'parties'"
            >
              See All
            </button>
          </div>
          <div v-if="parties.length" class="space-y-2">
            <div
              v-for="party in parties.slice(0, 3)" :key="party.id"
              class="flex items-center gap-3 rounded-xl p-2 transition hover:bg-white/4 cursor-pointer"
              @click="handleJoinParty(party.id)"
            >
              <div class="h-8 w-8 shrink-0 overflow-hidden rounded-lg bg-white/10">
                <div v-if="(party as any).currentTrackCover" class="h-full w-full bg-cover bg-center" :style="{ backgroundImage: `url(${(party as any).currentTrackCover})` }" />
                <div v-else class="flex h-full items-center justify-center">
                  <Music aria-hidden="true" class="text-xs text-white/30"  />
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-xs font-medium text-white">{{ party.title }}</p>
                <p class="text-[10px] text-white/30">{{ party.participant_count || 0 }} listening</p>
              </div>
              <span class="flex h-2 w-2 rounded-full bg-spotify motion-safe:animate-pulse shrink-0" />
            </div>
          </div>
          <div v-else class="flex flex-col items-center gap-3 py-8 text-center">
            <Users aria-hidden="true" class="text-xl text-white/20"  />
            <p class="text-xs text-white/30">No active parties</p>
            <button
              class="rounded-full bg-spotify/10 px-4 py-1.5 text-[10px] font-semibold text-spotify transition hover:bg-spotify/20"
              @click="openCreateWizard('party')"
            >
              Start One
            </button>
          </div>
        </div>

        <!-- Your Clubs -->
        <div class="rounded-2xl bg-white/6 p-5 ring-1 ring-white/10">
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-sm font-bold text-white">Your Clubs</h3>
            <button
              v-if="clubs.length > 3"
              class="text-[10px] text-white/30 hover:text-white/60 transition"
              @click="activeTab = 'clubs'"
            >
              See All
            </button>
          </div>
          <div v-if="myClubs.length" class="space-y-2">
            <div
              v-for="club in myClubs.slice(0, 3)" :key="club.id"
              class="flex items-center gap-3 rounded-xl p-2 transition hover:bg-white/4 cursor-pointer"
              @click="router.push(`/social/clubs/${club.id}`)"
            >
              <div
                class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-purple-500/10"
              >
                <Building2 aria-hidden="true" class="text-xs text-purple-400"  />
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-xs font-medium text-white">{{ club.name }}</p>
                <p class="flex items-center gap-1 text-[10px] text-white/30">
                  <Users aria-hidden="true" class="text-[8px]"  />
                  {{ club.member_count || 0 }} / {{ club.max_members || '∞' }}
                </p>
              </div>
              <div class="h-1 w-12 overflow-hidden rounded-full bg-white/5 shrink-0">
                <div
                  class="h-full rounded-full bg-purple-500/40"
                  :style="{ width: `${Math.min(100, ((club.member_count || 0) / (club.max_members || 1)) * 100)}%` }"
                />
              </div>
            </div>
          </div>
          <div v-else class="flex flex-col items-center gap-3 py-8 text-center">
            <Building2 aria-hidden="true" class="text-xl text-white/20"  />
            <p class="text-xs text-white/30">No clubs yet</p>
            <button
              class="rounded-full bg-purple-500/10 px-4 py-1.5 text-[10px] font-semibold text-purple-400 transition hover:bg-purple-500/20"
              @click="openCreateWizard('club')"
            >
              Create One
            </button>
          </div>
        </div>
      </div>

      <!-- ── Activity River ── -->
      <div class="rounded-2xl bg-white/6 p-5 ring-1 ring-white/10" role="feed" aria-label="Live activity feed">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-sm font-bold text-white">Activity</h3>
          <span class="flex items-center gap-1.5 text-[10px] text-white/30">
            <span class="h-1.5 w-1.5 rounded-full bg-spotify motion-safe:animate-pulse" aria-hidden="true" />
            Live
          </span>
        </div>
        <div v-if="activities.length" class="space-y-1">
          <ActivityItem
            v-for="item in visibleActivities"
            :key="item.id"
            :item="item"
            @action="handleActivityAction"
          />
        </div>
        <div v-else class="flex flex-col items-center gap-3 py-10 text-center">
          <Clock aria-hidden="true" class="text-2xl text-white/10"  />
          <p class="text-xs text-white/30">No activity yet — be the first!</p>
        </div>
        <button
          v-if="activities.length > 5"
          class="mt-3 w-full rounded-lg py-2 text-xs text-white/40 transition hover:bg-white/4 hover:text-white/60"
          @click="showAllActivities"
        >
          Show {{ activities.length - 5 }} more
        </button>
      </div>

      <!-- ── Discussions ── -->
      <div class="rounded-2xl bg-white/6 p-5 ring-1 ring-white/10">
        <h3 class="text-sm font-bold text-white mb-4">Discussions</h3>
        <DiscussionThread
          :discussions="discussions"
          :user-names="userNames"
          :user-avatars="userAvatars"
          :is-posting="discussionPosting"
          @create="handleCreateDiscussion"
          @reply="handleReplyDiscussion"
        />
      </div>
    </template>

    <!-- Create Wizard -->
    <CreateWizard
      :visible="showCreateWizard"
      :entity-type="createEntityType"
      @close="closeCreateWizard"
      @created="handleCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { Building2, Clock, Megaphone, Music, Users } from 'lucide-vue-next'
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppToast } from '@/composables/useAppToast'
import { SkeletonLoader } from '@/components/common'
import { useSocialApi } from '@/services/api/social'
import { usePlayerApi } from '@/services/api/player'
import { useUserApi } from '@/services/api/users'
import { useUserAuthStore } from '@/stores'
import { usePlayerStore } from '@/stores/player'
import type { ListeningParty, LiveRoom, MusicClub, Discussion } from '@/services/api/social'
import ListeningPartyCard from '@/components/social/ListeningPartyCard.vue'
import LiveRoomCard from '@/components/social/LiveRoomCard.vue'
import MusicClubCard from '@/components/social/MusicClubCard.vue'
import DiscussionThread from '@/components/social/DiscussionThread.vue'
import CreateWizard from '@/components/social/CreateWizard.vue'
import ActivityItem from '@/components/music/ActivityItem.vue'

const router = useRouter()
const toast = useAppToast()
const api = useSocialApi()
const playerApi = usePlayerApi()
const authStore = useUserAuthStore()
const playerStore = usePlayerStore()

const loading = ref(true)
const activeTab = ref('parties')
const discussionPosting = ref(false)

const parties = ref<ListeningParty[]>([])
const rooms = ref<LiveRoom[]>([])
const clubs = ref<MusicClub[]>([])
const myClubs = ref<MusicClub[]>([])
const discussions = ref<Discussion[]>([])
const activities = ref<any[]>([])
const userNames = ref<Record<string, string>>({})
const userAvatars = ref<Record<string, string | null>>({})
const trackNames = ref<Record<string, string>>({})

// Create wizard state
const showCreateWizard = ref(false)
const createEntityType = ref<'party' | 'room' | 'club'>('party')

// Computed
const currentTrack = computed(() => playerStore.currentTrack)

const greeting = computed(() => {
  if (!authStore.isAuthenticated) return 'Community'
  const name = authStore.user?.displayName || authStore.user?.username || ''
  const hour = new Date().getHours()
  let timeGreeting = 'Hello'
  if (hour < 12) timeGreeting = 'صبح بخیر'
  else if (hour < 17) timeGreeting = 'ظهر بخیر'
  else timeGreeting = 'عصر بخیر'
  return name ? `${timeGreeting}، ${name}` : 'Community'
})

const featuredRoom = computed(() => {
  return rooms.value.length > 0 ? rooms.value[0] : null
})

const visibleActivities = computed(() => {
  return activities.value.slice(0, 5)
})

const tabs = [
  { key: 'parties', label: 'Listening Parties' },
  { key: 'rooms', label: 'Live Rooms' },
  { key: 'clubs', label: 'Music Clubs' },
  { key: 'discussions', label: 'Discussions' },
]

// ── Data fetching ──

async function fetchUserName(userId: string) {
  if (userNames.value[userId]) return
  try {
    const profile = await useUserApi().getPublicUserProfile(userId)
    const p = profile as Record<string, unknown>
    userNames.value[userId] = String(p.full_name || p.username || userId.slice(0, 8))
    userAvatars.value[userId] = (p as any).avatar_url || (p as any).avatarUrl || null
  } catch {
    userNames.value[userId] = userId.slice(0, 8)
  }
}

async function fetchTrackName(trackId: string) {
  if (trackId! || trackNames.value[trackId]) return
  try {
    const track = await playerApi.getPlaybackTrack(trackId)
    trackNames.value[trackId] = `${track.artistName} - ${track.title}`
  } catch {
    trackNames.value[trackId] = trackId.slice(0, 12)
  }
}

async function loadData() {
  loading.value = true
  const [partiesRes, roomsRes, clubsRes, activitiesRes] = await Promise.allSettled([
    api.listParties({ limit: 50 }),
    api.listRooms({ limit: 50 }),
    api.listClubs({ limit: 50 }),
    api.getFeed({ limit: 20 }).catch(() => null),
  ])

  const rawParties = partiesRes.status === 'fulfilled' ? partiesRes.value : null
  const rawRooms = roomsRes.status === 'fulfilled' ? roomsRes.value : null
  const rawClubs = clubsRes.status === 'fulfilled' ? clubsRes.value : null
  const rawActivities = activitiesRes.status === 'fulfilled' ? (activitiesRes.value as any)?.items || [] : []

  parties.value = Array.isArray(rawParties) ? rawParties.filter(Boolean) : []
  rooms.value = Array.isArray(rawRooms) ? rawRooms.filter(Boolean) : []
  clubs.value = Array.isArray(rawClubs) ? rawClubs.filter(Boolean) : []
  myClubs.value = Array.isArray(rawClubs) ? rawClubs.filter(Boolean).slice(0, 10) : []
  activities.value = Array.isArray(rawActivities) ? rawActivities.filter(Boolean) : []

  // Fetch user names
  const userIds = new Set<string>()
  if (Array.isArray(rawParties)) rawParties.forEach((p: any) => userIds.add(p.host_id))
  if (Array.isArray(rawRooms)) rawRooms.forEach((r: any) => userIds.add(r.host_id))
  if (Array.isArray(rawClubs)) rawClubs.forEach((c: any) => userIds.add(c.created_by))
  rawActivities.forEach((a: any) => {
    if (a.userId || a.user_id) userIds.add(a.userId || a.user_id)
  })
  await Promise.all(Array.from(userIds).map(fetchUserName))

  // Fetch track names
  const trackIds = new Set<string>()
  if (Array.isArray(rawParties)) {
    rawParties.forEach((p: any) => {
      if (p.current_track_id) trackIds.add(p.current_track_id)
    })
  }
  await Promise.all(Array.from(trackIds).map(fetchTrackName))

  loading.value = false
}

// ── Actions ──

function openCreateWizard(type: 'party' | 'room' | 'club') {
  createEntityType.value = type
  showCreateWizard.value = true
}

function closeCreateWizard() {
  showCreateWizard.value = false
}

function handleCreated(id: string, type: string) {
  toast.success(`${type === 'party' ? 'Party' : type === 'room' ? 'Room' : 'Club'} created!`)
  loadData()
  if (type === 'party') router.push(`/social/parties/${id}`)
  else if (type === 'room') router.push(`/social/rooms/${id}`)
  else router.push(`/social/clubs/${id}`)
}

function startPartyFromTrack() {
  if (currentTrack.value) {
    createEntityType.value = 'party'
    showCreateWizard.value = true
  }
}

function showAllActivities() {
  activities.value = [...activities.value]
}

async function handleJoinParty(id: string) {
  try {
    await api.joinParty(id)
    toast.success('Joined party!')
    router.push(`/social/parties/${id}`)
  } catch (err: any) {
    toast.apiError(err, 'Failed to join party')
  }
}

async function handleJoinRoom(id: string) {
  try {
    await api.joinRoom(id)
    toast.success('Joined room!')
    router.push(`/social/rooms/${id}`)
  } catch (err: any) {
    toast.apiError(err, 'Failed to join room')
  }
}

async function handleJoinClub(id: string) {
  try {
    await api.joinClub(id)
    toast.success('Joined club!')
    router.push(`/social/clubs/${id}`)
  } catch (err: any) {
    toast.apiError(err, 'Failed to join club')
  }
}

async function handleCreateDiscussion(content: string, parentId?: string) {
  discussionPosting.value = true
  try {
    await api.createDiscussion({
      target_type: 'track',
      target_id: '0',
      content,
      parent_id: parentId,
    })
    toast.success(parentId ? 'Reply posted!' : 'Comment posted!')
  } catch (err: any) {
    toast.apiError(err, 'Failed to post')
  } finally {
    discussionPosting.value = false
  }
}

function handleReplyDiscussion(discussionId: string) {
  // opens reply input — DiscussionThread handles it internally
}

function handleActivityAction(item: any) {
  if (item.actionType === 'party' && item.targetUrl) {
    router.push(item.targetUrl)
  } else if (item.actionType === 'room' && item.targetUrl) {
    router.push(item.targetUrl)
  }
}

onMounted(loadData)
</script>
