<template>
  <div class="mx-auto max-w-5xl space-y-8 px-4 pt-20 pb-24 md:px-8">
    <!-- Header -->
    <div class="text-center">
      <h1 class="text-4xl font-black text-white md:text-5xl">Gamification</h1>
      <p class="mt-2 text-sm text-white/40">
        Earn XP, unlock badges, complete challenges, and climb the leaderboard.
      </p>
    </div>

    <!-- Profile & Stats -->
    <div class="glass-strong rounded-2xl p-6 md:p-8">
      <div class="flex flex-col gap-6 md:flex-row md:items-center md:justify-between">
        <!-- User Profile -->
        <div class="flex flex-1 items-center gap-4">
          <div class="relative">
            <div class="h-16 w-16 rounded-full bg-gradient-to-br from-spotify to-spotify/60 p-0.5">
              <div class="h-full w-full rounded-full bg-surface-overlay flex items-center justify-center">
                <span class="text-xl font-bold text-white">{{ profile?.username?.charAt(0) || 'U' }}</span>
              </div>
            </div>
            <div class="absolute -right-1 -bottom-1 flex h-5 w-5 items-center justify-center rounded-full bg-green-500 ring-2 ring-[var(--surface-2)]">
              <span class="text-xs">✓</span>
            </div>
          </div>
          <div class="min-w-0 flex-1">
            <h2 class="truncate text-lg font-bold text-white">{{ profile?.username || 'Guest' }}</h2>
            <p class="text-sm text-white/40">Level {{ profile?.level || 1 }} • Rank #{{ profile?.rank || '—' }}</p>
          </div>
        </div>

        <!-- Stats Grid -->
        <div class="grid flex-1 grid-cols-3 gap-4 text-center md:grid-cols-3">
          <!-- XP -->
          <div class="rounded-lg bg-white/5 p-3 transition-all hover:bg-white/10">
            <div class="text-xs uppercase tracking-wider text-white/40">XP</div>
            <div class="mt-1 text-lg font-bold text-white">{{ formatXP(profile?.total_xp ?? 0) }}</div>
          </div>
          <!-- Badges -->
          <div class="rounded-lg bg-white/5 p-3 transition-all hover:bg-white/10">
            <div class="text-xs uppercase tracking-wider text-white/40">Badges</div>
            <div class="mt-1 text-lg font-bold text-white">{{ profile?.badges?.length || 0 }}</div>
          </div>
          <!-- Challenges -->
          <div class="rounded-lg bg-white/5 p-3 transition-all hover:bg-white/10">
            <div class="text-xs uppercase tracking-wider text-white/40">Challenges</div>
            <div class="mt-1 text-lg font-bold text-white">{{ completedChallenges }}/{{ profile?.challenges?.length || 0 }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- XP Bar -->
    <GamificationXPBar :profile="profile" />

    <!-- Tabs -->
    <div class="flex gap-1 rounded-xl bg-white/4 p-1">
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

    <KeepAlive>
      <component :is="activeComponent" v-bind="(activeProps as any)" :key="activeTab" />
    </KeepAlive>

    <!-- Loading state -->
    <div v-if="loading" class="space-y-4">
      <SkeletonLoader variant="card" class="mx-auto max-w-md" />
      <SkeletonLoader variant="card" />
      <SkeletonLoader variant="card" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { SkeletonLoader } from '@/components/common'
import { useGamificationApi } from '@/services/api/gamification'
import GamificationXPBar from '@/components/gamification/GamificationXPBar.vue'
import GamificationBadges from '@/components/gamification/GamificationBadges.vue'
import GamificationChallenges from '@/components/gamification/GamificationChallenges.vue'
import GamificationLeaderboard from '@/components/gamification/GamificationLeaderboard.vue'
import type {
  GamificationProfile,
  BadgesResponse,
  ChallengesResponse,
  LeaderboardEntry,
} from '@/services/api/gamification'

const api = useGamificationApi()

const loading = ref(true)
const profile = ref<GamificationProfile | null>(null)
const badgesData = ref<BadgesResponse>({ all: [], mine: [] })
const challengesData = ref<ChallengesResponse>({ challenges: [], progress: [], earned_xp: 0 })
const leaderboardData = ref<LeaderboardEntry[]>([])

const completedChallenges = computed(() =>
  challengesData.value.progress.filter((p) => p.is_completed).length,
)

function formatXP(xp: number): string {
  if (xp >= 1000) return `${(xp / 1000).toFixed(1)}k`
  return String(xp)
}

const tabs = [
  { key: 'badges', label: 'Badges' },
  { key: 'challenges', label: 'Challenges' },
  { key: 'leaderboard', label: 'Leaderboard' },
]

const activeTab = ref('badges')

const activeComponent = computed(() => {
  switch (activeTab.value) {
    case 'badges': return GamificationBadges
    case 'challenges': return GamificationChallenges
    case 'leaderboard': return GamificationLeaderboard
    default: return GamificationBadges
  }
})

const leaderboardType = ref('all')

/** Map frontend leaderboard tab keys to backend API types. */
const LEADERBOARD_TYPE_MAP: Record<string, string> = {
  all: 'all',
  monthly: 'xp_monthly',
  weekly: 'xp_weekly',
  streams: 'streams',
  friends: 'friends',
}

const activeProps = computed(() => {
  switch (activeTab.value) {
    case 'badges':
      return { allBadges: badgesData.value.all ?? [], userBadges: badgesData.value.mine ?? [] }
    case 'challenges':
      return { challenges: challengesData.value.challenges, progress: challengesData.value.progress, earnedXP: challengesData.value.earned_xp }
    case 'leaderboard':
      return { entries: leaderboardData.value, type: leaderboardType.value, userRank: profile.value?.rank, userId: profile.value?.user_id, onTypeChange: onLeaderboardTypeChange }
    default:
      return {}
  }
})

function onLeaderboardTypeChange(type: string) {
  leaderboardType.value = type
  const apiType = LEADERBOARD_TYPE_MAP[type] ?? type
  api.getLeaderboard({ type: apiType, limit: 50 }).then((res) => {
    if (res) leaderboardData.value = res.entries
  })
}

// ─── Live polling for challenges ─────────────────────────────
// Challenges progress is updated in the backend as the user
// performs actions (playing tracks, liking songs, etc.). Poll
// every 12 seconds so the user sees their progress update live.
const CHALLENGE_POLL_MS = 12_000
let pollTimer: ReturnType<typeof setInterval> | undefined

async function refreshChallenges() {
  const res = await api.getChallenges()
  if (res) challengesData.value = res
}

// Poll only when the challenges tab is visible
watch(activeTab, (tab) => {
  clearInterval(pollTimer)
  if (tab === 'challenges') {
    refreshChallenges() // immediate refresh on tab switch
    pollTimer = setInterval(refreshChallenges, CHALLENGE_POLL_MS)
  }
})

onUnmounted(() => clearInterval(pollTimer))

// ─── Profile auto-refresh (XP/level updates after plays) ────
const PROFILE_POLL_MS = 30_000
let profileTimer: ReturnType<typeof setInterval> | undefined

async function refreshProfile() {
  const res = await api.getProfile()
  if (res) profile.value = res
}

onMounted(() => {
  profileTimer = setInterval(refreshProfile, PROFILE_POLL_MS)
})

onUnmounted(() => clearInterval(profileTimer))

// ─── Initial load ────────────────────────────────────────────
onMounted(async () => {
  try {
    const [profileRes, badgesRes, challengesRes, leaderboardRes] = await Promise.all([
      api.getProfile(),
      api.getBadges(),
      api.getChallenges(),
      api.getLeaderboard({ type: 'all', limit: 50 }),
    ])
    if (profileRes) profile.value = profileRes
    if (badgesRes) badgesData.value = badgesRes
    if (challengesRes) challengesData.value = challengesRes
    if (leaderboardRes) leaderboardData.value = leaderboardRes.entries
  } catch (err) {
    console.error('Failed to load gamification data:', err)
  } finally {
    loading.value = false
  }
})
</script>
