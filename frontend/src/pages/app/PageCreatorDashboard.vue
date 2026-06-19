<template>
  <div class="mx-auto w-full max-w-7xl px-4 pt-4 pb-32 md:px-6 lg:px-8">
    <!-- Loading skeleton -->
    <div v-if="loading" class="space-y-6">
      <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
        <SkeletonLoader v-for="i in 4" :key="i" variant="card" />
      </div>
      <SkeletonLoader variant="card" class="h-[300px]" />
    </div>

    <!-- Not a creator -->
    <div v-else-if="!isCreator" class="flex flex-col items-center gap-4 py-24 text-center">
      <i aria-hidden="true" class="pi pi-megaphone text-4xl text-slate-500" />
      <h2 class="text-xl font-bold text-white">Creator Studio</h2>
      <p class="max-w-md text-sm text-slate-400">
        Upload tracks and build your audience to unlock creator analytics.
      </p>
      <RouterLink to="/" class="text-sm font-medium text-[#1db954] underline underline-offset-2"
        >Browse music</RouterLink
      >
    </div>

    <!-- Creator dashboard -->
    <template v-else>
      <!-- Header -->
      <div class="relative mb-8 overflow-hidden rounded-[2rem] p-8 md:p-10">
        <div
          class="absolute inset-0 bg-gradient-to-br from-[#1db954]/8 via-transparent to-transparent"
        />
        <div
          class="relative z-10 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
        >
          <div>
            <p class="text-xs font-bold tracking-[0.25em] text-[#1db954] uppercase">Studio</p>
            <h1 class="mt-2 text-3xl font-black text-white md:text-4xl">Creator Studio</h1>
            <p class="mt-1 text-sm text-white/40">
              Analytics, earnings, audience, and content management
            </p>
          </div>
          <div class="flex items-center gap-3">
            <RouterLink
              to="/admin/media"
              class="inline-flex items-center gap-2 rounded-full bg-[#1db954] px-6 py-2.5 text-sm font-bold text-black transition hover:scale-105 hover:bg-[#1ed760]"
            >
              <i aria-hidden="true" class="pi pi-upload text-xs" /> Upload
            </RouterLink>
            <button
              type="button"
              :disabled="refreshing"
              @click="refreshStats"
              class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/[0.04] px-5 py-2.5 text-sm font-medium text-white backdrop-blur transition hover:bg-white/[0.10] disabled:opacity-50"
            >
              <i aria-hidden="true" :class="refreshing ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" />
              {{ refreshing ? 'Refreshing...' : 'Refresh' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="flex gap-1 rounded-xl bg-white/[0.04] p-1">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="flex-1 rounded-lg py-2.5 text-sm font-medium transition-all duration-200"
          :class="
            activeTab === tab.key
              ? 'bg-white/10 text-white shadow-lg'
              : 'text-white/30 hover:text-white/50'
          "
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- Tab: Overview -->
      <div v-show="activeTab === 'overview'" class="mt-6 space-y-6">
        <!-- Metric Cards -->
        <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
          <CreatorMetricCard label="Total Plays" :value="totalPlays" icon="pi pi-play" />
          <CreatorMetricCard label="Listeners" :value="uniqueListeners" icon="pi pi-users" />
          <CreatorMetricCard label="Followers" :value="totalFollowers" icon="pi pi-heart" />
          <CreatorMetricCard
            label="Revenue"
            :value="estimatedRevenue"
            icon="pi pi-dollar"
            prefix="$"
          />
        </div>
        <div class="grid grid-cols-2 gap-4 md:grid-cols-3">
          <CreatorMetricCard label="Tracks" :value="totalTracks" icon="pi pi-music" />
          <CreatorMetricCard label="Albums" :value="totalAlbums" icon="pi pi-compact-disc" />
          <CreatorMetricCard label="Playlists" :value="totalPlaylists" icon="pi pi-list" />
        </div>

        <!-- Daily bar chart -->
        <section>
          <h3 class="mb-3 text-lg font-bold text-white">Daily Plays</h3>
          <div class="rounded-2xl bg-white/[0.04] p-6">
            <div
              v-if="dailyStats.length === 0"
              class="flex flex-col items-center gap-3 py-16 text-center"
            >
              <i aria-hidden="true" class="pi pi-inbox text-4xl text-slate-500" />
              <p class="text-sm text-slate-400">No daily data yet</p>
            </div>
            <div v-else class="space-y-2">
              <div
                v-for="day in dailyStats.slice(0, 14)"
                :key="day.id"
                class="flex items-center gap-3"
              >
                <span class="w-24 shrink-0 text-xs text-slate-400">{{ formatDate(day.date) }}</span>
                <div class="h-6 flex-1 overflow-hidden rounded-full bg-white/[0.06]">
                  <div
                    class="h-full rounded-full bg-gradient-to-r from-[#1db954] to-[#1ed760] transition-all duration-500"
                    :style="{ width: barWidth(day.plays) + '%' }"
                  />
                </div>
                <span class="w-12 text-right text-xs font-medium text-white tabular-nums">{{
                  day.plays
                }}</span>
              </div>
            </div>
          </div>
        </section>

        <!-- Track list -->
        <section>
          <h3 class="mb-3 text-lg font-bold text-white">Top Tracks</h3>
          <div v-if="!trackStats.length" class="flex flex-col items-center gap-3 py-16 text-center">
            <i aria-hidden="true" class="pi pi-inbox text-4xl text-slate-500" />
            <p class="text-sm text-slate-400">No track stats yet</p>
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="track in trackStats.slice(0, 10)"
              :key="track.track_id"
              @click="playTrack(track)"
              class="group flex cursor-pointer items-center gap-4 rounded-xl bg-white/[0.03] px-4 py-3 transition hover:bg-white/[0.06]"
            >
              <div
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-white/10 text-xs text-white/30"
              >
                {{ track.total_plays }}
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                <p class="text-xs text-slate-500">
                  {{ formatNumber(track.total_plays) }} plays ·
                  {{ formatNumber(track.total_likes) }} likes
                </p>
              </div>
              <p class="text-sm font-bold text-white tabular-nums">
                {{ formatNumber(track.total_plays) }}
              </p>
            </div>
          </div>
        </section>
      </div>

      <!-- Tab: Earnings -->
      <div v-show="activeTab === 'earnings'" class="mt-6 space-y-6">
        <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
          <CreatorMetricCard
            label="Total Revenue"
            :value="earnings.total_revenue"
            icon="pi pi-dollar"
            prefix="$"
          />
          <CreatorMetricCard
            label="Stream Revenue"
            :value="earnings.stream_revenue"
            icon="pi pi-play"
            prefix="$"
          />
          <CreatorMetricCard
            label="Tips"
            :value="earnings.tip_revenue"
            icon="pi pi-heart"
            prefix="$"
          />
          <CreatorMetricCard
            label="Subscriptions"
            :value="earnings.subscription_revenue"
            icon="pi pi-users"
            prefix="$"
          />
        </div>
        <div class="grid grid-cols-2 gap-4 md:grid-cols-3">
          <CreatorMetricCard
            label="Pending Payout"
            :value="earnings.pending_payout"
            icon="pi pi-clock"
            prefix="$"
          />
          <CreatorMetricCard
            label="Last Payout"
            :value="earnings.last_payout"
            icon="pi pi-check-circle"
            prefix="$"
          />
        </div>

        <section>
          <h3 class="mb-3 text-lg font-bold text-white">Revenue Breakdown</h3>
          <div class="rounded-2xl bg-white/[0.04] p-6">
            <div class="space-y-3">
              <div
                v-for="item in revenueBreakdown"
                :key="item.label"
                class="flex items-center gap-3"
              >
                <span class="w-32 shrink-0 text-xs text-slate-400">{{ item.label }}</span>
                <div class="h-6 flex-1 overflow-hidden rounded-full bg-white/[0.06]">
                  <div
                    class="h-full rounded-full transition-all duration-500"
                    :class="item.color"
                    :style="{ width: revenuePercent(item.value) + '%' }"
                  />
                </div>
                <span class="w-20 text-right text-xs font-medium text-white tabular-nums"
                  >${{ formatNumber(item.value) }}</span
                >
              </div>
            </div>
          </div>
        </section>

        <section>
          <h3 class="mb-3 text-lg font-bold text-white">Payout History</h3>
          <div v-if="payouts.length === 0" class="flex flex-col items-center gap-3 py-16 text-center">
            <i aria-hidden="true" class="pi pi-inbox text-4xl text-slate-500" />
            <p class="text-sm text-slate-400">No payouts yet</p>
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="p in payouts"
              :key="p.id"
              class="flex items-center justify-between rounded-xl bg-white/[0.03] px-4 py-3"
            >
              <div>
                <p class="text-sm font-medium text-white">${{ formatNumber(p.amount) }}</p>
                <p class="text-xs text-slate-500">
                  {{ p.method }} · {{ formatDate(p.created_at) }}
                </p>
              </div>
              <span
                class="rounded-full px-2.5 py-0.5 text-[10px] font-medium tracking-wider uppercase"
                :class="
                  p.status === 'paid'
                    ? 'bg-green-500/10 text-green-400'
                    : 'bg-amber-500/10 text-amber-400'
                "
              >
                {{ p.status }}
              </span>
            </div>
          </div>
        </section>
      </div>

      <!-- Tab: Audience -->
      <div v-show="activeTab === 'audience'" class="mt-6 space-y-6">
        <div class="grid grid-cols-3 gap-4">
          <CreatorMetricCard
            label="Total Listeners"
            :value="audience.overview?.total_listeners || 0"
            icon="pi pi-users"
          />
          <CreatorMetricCard
            label="New (7d)"
            :value="audience.overview?.new_listeners_7d || 0"
            icon="pi pi-user-plus"
          />
          <div class="rounded-2xl bg-white/[0.04] p-5">
            <p class="text-2xl font-black text-white">{{ audienceRepeatRate }}%</p>
            <p class="mt-1 text-xs font-medium text-white/40">Repeat Rate</p>
          </div>
        </div>

        <section>
          <h3 class="mb-3 text-lg font-bold text-white">Top Listeners</h3>
          <div
            v-if="audience.top_listeners?.length === 0"
            class="flex flex-col items-center gap-3 py-16 text-center"
          >
            <i aria-hidden="true" class="pi pi-inbox text-4xl text-slate-500" />
            <p class="text-sm text-slate-400">No listener data yet</p>
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="(l, i) in audience.top_listeners?.slice(0, 10)"
              :key="l.user_id"
              class="flex items-center gap-3 rounded-xl bg-white/[0.03] px-4 py-3"
            >
              <span class="w-6 text-xs font-bold text-white/20">#{{ i + 1 }}</span>
              <div class="h-8 w-8 shrink-0 overflow-hidden rounded-full bg-white/10">
                <img
                  v-if="l.avatar_url"
                  :src="l.avatar_url"
                  :alt="l.username"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div
                  v-else
                  class="flex h-full w-full items-center justify-center text-xs font-bold text-white/30"
                >
                  {{ l.username.charAt(0).toUpperCase() }}
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-white">{{ l.username }}</p>
              </div>
              <p class="text-sm font-bold text-white/60 tabular-nums">{{ l.play_count }}</p>
            </div>
          </div>
        </section>

        <section>
          <h3 class="mb-3 text-lg font-bold text-white">Geographic</h3>
          <div
            v-if="audience.geographics?.length === 0"
            class="flex flex-col items-center gap-3 py-16 text-center"
          >
            <i aria-hidden="true" class="pi pi-inbox text-4xl text-slate-500" />
            <p class="text-sm text-slate-400">No geo data yet</p>
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="g in audience.geographics?.slice(0, 10)"
              :key="g.country + g.city"
              class="flex items-center gap-3 rounded-xl bg-white/[0.03] px-4 py-3"
            >
              <span class="text-lg">{{ countryFlag(g.country) }}</span>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-white">{{ g.city }}, {{ g.country }}</p>
                <p class="text-xs text-slate-500">{{ g.listeners }} listeners</p>
              </div>
              <p class="text-sm font-bold text-white/60 tabular-nums">{{ g.plays }}</p>
            </div>
          </div>
        </section>
      </div>

      <!-- Tab: Content -->
      <div v-show="activeTab === 'content'" class="mt-6 space-y-6">
        <section>
          <h3 class="mb-3 text-lg font-bold text-white">
            Tracks ({{ contentData.tracks?.length || 0 }})
          </h3>
          <div v-if="!contentData.tracks?.length" class="flex flex-col items-center gap-3 py-16 text-center">
            <i aria-hidden="true" class="pi pi-inbox text-4xl text-slate-500" />
            <p class="text-sm text-slate-400">No tracks uploaded yet</p>
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="track in contentData.tracks"
              :key="track.track_id"
              class="group flex items-center gap-3 rounded-xl bg-white/[0.03] px-4 py-3 transition hover:bg-white/[0.06]"
            >
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-white">{{ track.title }}</p>
                <p class="text-xs text-slate-500">{{ formatNumber(track.total_plays) }} plays</p>
              </div>
              <button
                class="rounded-lg bg-white/5 px-3 py-1.5 text-xs text-white/40 transition hover:bg-white/10 hover:text-white"
                @click="editTrack = track; showTrackModal = true"
              >
                Edit
              </button>
            </div>
          </div>
        </section>

        <section>
          <h3 class="mb-3 text-lg font-bold text-white">
            Albums ({{ contentData.albums?.length || 0 }})
          </h3>
          <div v-if="!contentData.albums?.length" class="flex flex-col items-center gap-3 py-16 text-center">
            <i aria-hidden="true" class="pi pi-inbox text-4xl text-slate-500" />
            <p class="text-sm text-slate-400">No albums yet</p>
          </div>
          <div v-else class="grid grid-cols-2 gap-4 md:grid-cols-4">
            <div
              v-for="album in contentData.albums"
              :key="album.id"
              class="rounded-xl bg-white/[0.03] p-4 transition hover:bg-white/[0.06]"
            >
              <div class="mb-2 h-24 w-full overflow-hidden rounded-lg bg-white/10">
                <img
                  v-if="album.cover_url"
                  :src="album.cover_url"
                  :alt="album.title"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center text-white/20">
                  <i aria-hidden="true" class="pi pi-compact-disc text-2xl" />
                </div>
              </div>
              <p class="truncate text-sm font-medium text-white">{{ album.title }}</p>
              <p class="text-xs text-slate-500">
                {{ album.track_count }} tracks · {{ formatNumber(album.total_plays) }} plays
              </p>
            </div>
          </div>
        </section>
      </div>
    </template>

    <!-- Track edit modal -->
    <Teleport to="body">
      <div
        v-if="showTrackModal"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
        @click.self="showTrackModal = false"
      >
        <div class="glass-strong mx-4 w-full max-w-lg rounded-2xl p-8">
          <h2 class="mb-6 text-xl font-bold text-white">Edit Track</h2>
          <div v-if="editTrack" class="space-y-4">
            <div>
              <label class="mb-1 block text-xs text-white/40">Title</label>
              <input
                v-model="editForm.title"
                aria-label="Track title"
                class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white outline-none focus:border-white/20"
              />
            </div>
            <div>
              <label class="mb-1 block text-xs text-white/40">Persian Title</label>
              <input
                v-model="editForm.persian_title"
                aria-label="Persian title"
                class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white outline-none focus:border-white/20"
              />
            </div>
            <div>
              <label class="mb-1 block text-xs text-white/40">Lyrics</label>
              <textarea
                v-model="editForm.lyrics"
                rows="4"
                aria-label="Lyrics"
                class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm text-white outline-none focus:border-white/20"
              />
            </div>
            <label class="flex items-center gap-3">
              <input
                v-model="editForm.explicit"
                type="checkbox"
                class="h-5 w-5 rounded border-white/10 bg-white/5 accent-[#1db954]"
              />
              <span class="text-sm text-white/60">Explicit</span>
            </label>
          </div>
          <div class="mt-6 flex gap-3">
            <button
              class="flex-1 rounded-xl bg-white/5 py-3 text-sm font-medium text-white/50 hover:bg-white/10"
              @click="showTrackModal = false"
            >
              Cancel
            </button>
            <button
              class="flex-1 rounded-xl bg-[#1db954] py-3 text-sm font-bold text-black hover:bg-[#1db954]/90"
              @click="saveTrack"
            >
              Save
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { SkeletonLoader } from '@/components/common'
import { useCreatorApi } from '@/services/api/creator'
import { usePlayerApi } from '@/services/api/player'
import { usePlayer } from '@/composables/player'
import { onImgError } from '@/utils/helpers'
import CreatorMetricCard from '@/components/creator/CreatorMetricCard.vue'
import type {
  TrackStats,
  EarningsBreakdown,
  AudienceData,
  CreatorContentData,
  Payout,
} from '@/services/api/creator/types'

const creatorApi = useCreatorApi()
const player = usePlayer()
const playerApi = usePlayerApi()

const loading = ref(true)
const refreshing = ref(false)
const isCreator = ref(false)
const activeTab = ref('overview')
const showTrackModal = ref(false)
const editTrack = ref<TrackStats | null>(null)

const tabs = [
  { key: 'overview', label: 'Overview' },
  { key: 'earnings', label: 'Earnings' },
  { key: 'audience', label: 'Audience' },
  { key: 'content', label: 'Content' },
]

const stats = ref<Record<string, unknown> | null>(null)
const dailyStats = ref<Array<Record<string, unknown>>>([])
const trackStats = ref<TrackStats[]>([])
const totalPlays = ref(0)
const uniqueListeners = ref(0)
const totalFollowers = ref(0)
const totalTracks = ref(0)
const totalAlbums = ref(0)
const estimatedRevenue = ref(0)
const totalPlaylists = ref(0)
let maxPlays = 1

const earnings = reactive<EarningsBreakdown>({
  total_revenue: 0,
  stream_revenue: 0,
  tip_revenue: 0,
  subscription_revenue: 0,
  pending_payout: 0,
  last_payout: 0,
})
const payouts = ref<Payout[]>([])
const audience = reactive<AudienceData>({
  overview: { total_listeners: 0, new_listeners_7d: 0, repeat_rate: 0 },
  top_listeners: [],
  geographics: [],
})
const contentData = reactive<CreatorContentData>({ tracks: [], albums: [], playlists: [] })

const editForm = reactive({ title: '', persian_title: '', lyrics: '', explicit: false })

const audienceRepeatRate = computed(() => {
  if (!audience.overview) return '0'
  return (audience.overview.repeat_rate * 100).toFixed(1)
})

const revenueBreakdown = computed(() => [
  { label: 'Streams', value: earnings.stream_revenue, color: 'bg-[#1db954]' },
  { label: 'Tips', value: earnings.tip_revenue, color: 'bg-[#e91e63]' },
  { label: 'Subscriptions', value: earnings.subscription_revenue, color: 'bg-[#60a5fa]' },
])

function revenuePercent(value: number): number {
  const total = earnings.total_revenue || 1
  return (value / total) * 100
}

async function fetchDashboard() {
  loading.value = true
  try {
    const [overview, daily, tracks, creatorCheck] = await Promise.all([
      creatorApi.getOverview(),
      creatorApi.getDailyStats({ limit: 30 }),
      creatorApi.getTrackStats(),
      creatorApi.isCreator(),
    ])
    if (overview.success && overview.data) {
      stats.value = overview.data
      totalPlays.value = overview.data.total_plays
      uniqueListeners.value = overview.data.unique_listeners
      totalFollowers.value = overview.data.total_followers
      totalTracks.value = overview.data.total_tracks
      totalAlbums.value = overview.data.total_albums
      estimatedRevenue.value = overview.data.estimated_revenue
      totalPlaylists.value = overview.data.total_playlists
    }
    if (daily.success && daily.data) {
      dailyStats.value = daily.data
      maxPlays = Math.max(1, ...daily.data.map((d: Record<string, unknown>) => (d as { plays: number }).plays))
    }
    if (tracks.success && tracks.data) trackStats.value = tracks.data
    isCreator.value = creatorCheck.is_creator
  } catch (err) {
    console.error('Failed to load dashboard:', err)
  } finally {
    loading.value = false
  }
}

async function fetchSecondaryData() {
  try {
    const [e, p, a, c] = await Promise.all([
      creatorApi.getEarnings(),
      creatorApi.getPayouts({ limit: 20 }),
      creatorApi.getAudience({ top_limit: 20 }),
      creatorApi.getContent(),
    ])
    if (e?.data) Object.assign(earnings, e.data)
    if (p?.data) payouts.value = p.data
    if (a?.data) Object.assign(audience, a.data)
    if (c?.data) Object.assign(contentData, c.data)
  } catch (err) {
    console.error('Failed to load secondary data:', err)
  }
}

async function refreshStats() {
  refreshing.value = true
  try {
    await creatorApi.refreshStats()
    await fetchDashboard()
    await fetchSecondaryData()
  } catch (err) {
    console.error('Failed to refresh:', err)
  } finally {
    refreshing.value = false
  }
}

function playTrack(track: Record<string, unknown>) {
  const pb = {
    id: String(track.track_id),
    title: track.title,
    artistName: track.artist_name || 'Unknown',
    albumTitle: null,
    coverUrl: track.cover_url || null,
    durationSeconds: track.duration ?? null,
    streamUrl: playerApi.getTrackStreamUrl(String(track.track_id)),
  }
  player.playTrack(pb)
}

async function saveTrack() {
  if (!editTrack.value) return
  try {
    await creatorApi.updateTrack(editTrack.value.track_id, {
      title: editForm.title,
      persian_title: editForm.persian_title,
      lyrics: editForm.lyrics,
      explicit: editForm.explicit,
    })
    showTrackModal.value = false
    await fetchSecondaryData()
  } catch (err) {
    console.error('Failed to save track:', err)
  }
}

function barWidth(plays: number) {
  return (plays / maxPlays) * 100
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function formatNumber(n: number) {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`
  if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`
  return String(n)
}

function countryFlag(country: string): string {
  const flags: Record<string, string> = {
    Iran: '🇮🇷',
    'United States': '🇺🇸',
    'United Kingdom': '🇬🇧',
    Canada: '🇨🇦',
    Germany: '🇩🇪',
    France: '🇫🇷',
    Turkey: '🇹🇷',
    Afghanistan: '🇦🇫',
    Tajikistan: '🇹🇯',
    Azerbaijan: '🇦🇿',
    Iraq: '🇮🇶',
    UAE: '🇦🇪',
    Sweden: '🇸🇪',
    Netherlands: '🇳🇱',
    Australia: '🇦🇺',
  }
  return flags[country] || '🌍'
}

onMounted(async () => {
  await fetchDashboard()
  if (isCreator.value) await fetchSecondaryData()
})
</script>
