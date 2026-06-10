<template>
  <div class="mx-auto max-w-5xl space-y-8 px-4 pt-20 pb-24 md:px-8">
    <!-- Header -->
    <div class="text-center">
      <h1 class="text-4xl font-black text-white md:text-5xl">Gamification</h1>
      <p class="mt-2 text-sm text-white/40">
        Earn XP, unlock badges, complete challenges, and climb the leaderboard.
      </p>
    </div>

    <!-- XP Bar -->
    <GamificationXPBar :profile="profile" />

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

    <!-- Badges -->
    <GamificationBadges
      v-show="activeTab === 'badges'"
      :all-badges="badgesData.all"
      :user-badges="badgesData.mine"
    />

    <!-- Challenges -->
    <GamificationChallenges
      v-show="activeTab === 'challenges'"
      :challenges="challengesData.challenges"
      :progress="challengesData.progress"
      :earnedXP="challengesData.earned_xp"
    />

    <!-- Leaderboard -->
    <GamificationLeaderboard v-show="activeTab === 'leaderboard'" :entries="leaderboardData" />

    <!-- Loading state -->
    <div v-if="loading" class="space-y-4">
      <SkeletonLoader variant="card" class="mx-auto max-w-md" />
      <SkeletonLoader variant="card" />
      <SkeletonLoader variant="card" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
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

const tabs = [
  { key: 'badges', label: 'Badges' },
  { key: 'challenges', label: 'Challenges' },
  { key: 'leaderboard', label: 'Leaderboard' },
]

const activeTab = ref('badges')

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
  } catch {
    // Silently handle
  } finally {
    loading.value = false
  }
})
</script>
