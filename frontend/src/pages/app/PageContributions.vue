<template>
  <div class="mx-auto max-w-6xl space-y-8 px-4 py-6 md:px-6 md:py-8">
    <!-- Header -->
    <div class="flex items-start justify-between gap-4">
      <div>
        <h1 class="text-3xl font-bold text-white md:text-4xl">Contributions</h1>
        <p class="mt-1 text-sm text-white/40">
          Help improve the community by contributing lyrics, translations, credits, and more
        </p>
      </div>
      <router-link
        to="/upload"
        class="spring shrink-0 rounded-xl bg-spotify/10 px-4 py-2.5 text-sm font-medium text-spotify transition-all hover:bg-spotify/20"
      >
        <Upload aria-hidden="true" class="mr-1.5 inline text-xs" />
        Upload Music
      </router-link>
    </div>

    <!-- Reputation Strip -->
    <div
      v-if="reputation"
      class="flex flex-wrap items-center gap-4 rounded-2xl border border-white/6 bg-white/3 px-5 py-3"
    >
      <div class="flex items-center gap-3">
        <div
          class="flex h-9 w-9 items-center justify-center rounded-full"
          :class="tierBgClass(reputation.tier)"
        >
          <Shield aria-hidden="true" class="text-sm" :class="tierIconClass(reputation.tier)" />
        </div>
        <div>
          <p class="text-xs font-medium text-white/70">
            {{ reputation.tierLabel || reputation.tier || 'Newcomer' }}
          </p>
          <p class="text-[10px] text-white/30">Trust Tier</p>
        </div>
      </div>
      <div class="h-8 w-px bg-white/6" />
      <div class="flex items-center gap-2 text-sm tabular-nums">
        <span class="font-semibold text-white">{{ reputation.trustScore ?? reputation.score ?? 0 }}</span>
        <span class="text-xs text-white/30">Score</span>
      </div>
      <div class="h-8 w-px bg-white/6" />
      <div class="flex items-center gap-2 text-sm tabular-nums">
        <span class="font-semibold text-white">{{ reputation.acceptedContributions ?? 0 }}</span>
        <span class="text-xs text-white/30">Accepted</span>
      </div>
      <div class="h-8 w-px bg-white/6" />
      <div class="flex items-center gap-2 text-sm tabular-nums">
        <span class="font-semibold text-white">{{ reputation.uploadSlots ?? 0 }}</span>
        <span class="text-xs text-white/30">Upload Slots</span>
      </div>
      <div class="h-8 w-px bg-white/6" />
      <div
        class="rounded-full px-2.5 py-0.5 text-[10px] font-medium"
        :class="reputation.autoPublish ? 'bg-emerald-500/10 text-emerald-400' : 'bg-amber-500/10 text-amber-400'"
      >
        {{ reputation.autoPublish ? 'Auto-publish' : 'Needs review' }}
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex items-center gap-1 rounded-xl bg-white/4 p-1">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        class="spring flex-1 rounded-lg px-4 py-2 text-sm font-medium transition-all"
        :class="
          activeTab === tab.key
            ? 'bg-white/10 text-white shadow-xs'
            : 'text-white/40 hover:text-white/60'
        "
        @click="activeTab = tab.key; loadTab(tab.key)"
      >
        <i aria-hidden="true" :class="tab.icon" class="mr-1.5" />
        {{ tab.label }}
      </button>
    </div>

    <!-- Tab Content -->
    <div v-if="activeTab === 'submit'">
      <ContributionSubmitForm @submitted="onSubmitted" />
    </div>

    <div v-if="activeTab === 'my'">
      <div v-if="loadingMy" class="space-y-3">
        <SkeletonLoader v-for="i in 3" :key="i" variant="lines" :lines="3" />
      </div>

      <div
        v-else-if="myContributions.length === 0"
        class="glass-strong flex flex-col items-center gap-3 rounded-2xl py-12 text-center"
      >
        <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/5">
          <Inbox aria-hidden="true" class="text-2xl text-white/15"  />
        </div>
        <p class="text-sm text-white/25">No contributions yet</p>
        <button
          type="button"
          class="text-xs text-spotify transition-colors hover:text-spotify-hover"
          @click="activeTab = 'submit'"
        >
          Make your first contribution
        </button>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="c in myContributions"
          :key="c.id"
          class="contribution-card glass-strong spring rounded-2xl p-5 transition-all hover:bg-white/6"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span
                  class="rounded-lg px-2 py-0.5 text-[10px] font-semibold tracking-wider uppercase"
                  :class="typeBadgeClass(c.contribution_type)"
                  >{{ c.contribution_type }}</span
                >
                <span
                  class="rounded-lg px-2 py-0.5 text-[10px] font-semibold"
                  :class="statusBadgeClass(c.status)"
                  >{{ c.status }}</span
                >
                <span v-if="c.is_minor" class="text-[10px] text-white/20">Minor</span>
              </div>
              <p class="mt-2 text-sm font-medium text-white/70">
                <router-link
                  :to="{ name: c.target_type, params: { id: c.target_id } }"
                  class="text-blue-400/60 hover:text-blue-400"
                >
                  {{ targetLabel(c.target_type) }} · {{ c.target_id }}
                </router-link>
              </p>
              <p v-if="c.summary" class="mt-0.5 text-xs text-white/30">{{ c.summary }}</p>
              <div class="mt-2 flex items-center gap-3 text-[10px] text-white/20">
                <span>v{{ c.version }}</span>
                <span>{{ formatDate(c.created_at) }}</span>
                <span v-if="c.xp_awarded > 0" class="text-spotify">+{{ c.xp_awarded }} XP</span>
              </div>
            </div>

            <div class="flex shrink-0 items-center gap-2">
              <button
                type="button"
                class="spring flex h-8 w-8 items-center justify-center rounded-full text-white/30 transition-all hover:bg-white/10 hover:text-white/60"
                @click="viewHistory(c)"
                aria-label="View history"
              >
                <History aria-hidden="true" class="text-xs"  />
              </button>
            </div>
          </div>

          <!-- Expanded history -->
          <div v-if="expandedContribution === c.id" class="mt-4 border-t border-white/6 pt-4">
            <ContributionHistory :items="historyItems" :loading="loadingHistory" />
          </div>
        </div>

        <!-- Pagination -->
        <div
          v-if="myMeta && myMeta.total > myMeta.page_size"
          class="flex items-center justify-center gap-2 pt-4"
        >
          <button
            type="button"
            :disabled="myMeta.page <= 1"
            class="spring rounded-lg bg-white/5 px-4 py-2 text-xs font-medium text-white/50 transition-all hover:bg-white/10 disabled:opacity-30"
            @click="loadMyContributions(myMeta.page - 1)"
          >
            Previous
          </button>
          <span class="text-xs text-white/30"
            >Page {{ myMeta.page }} of {{ Math.ceil(myMeta.total / myMeta.page_size) }}</span
          >
          <button
            type="button"
            :disabled="myMeta.page * myMeta.page_size >= myMeta.total"
            class="spring rounded-lg bg-white/5 px-4 py-2 text-xs font-medium text-white/50 transition-all hover:bg-white/10 disabled:opacity-30"
            @click="loadMyContributions(myMeta.page + 1)"
          >
            Next
          </button>
        </div>
      </div>
    </div>

    <div v-if="activeTab === 'leaderboard'">
      <ContributionLeaderboard :contributors="leaderboardData" :loading="loadingLeaderboard" />
    </div>

    <div v-if="activeTab === 'moderate'">
      <div v-if="loadingPending" class="space-y-3">
        <SkeletonLoader v-for="i in 3" :key="i" variant="lines" :lines="4" />
      </div>

      <div
        v-else-if="pendingItems.length === 0"
        class="glass-strong flex flex-col items-center gap-3 rounded-2xl py-12 text-center"
      >
        <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-white/5">
          <CheckCircle class="text-2xl text-white/15"<i aria-hidden="true"  /> />
        </div>
        <p class="text-sm text-white/25">All caught up! No pending contributions to review</p>
      </div>

      <div v-else class="space-y-3">
        <div v-for="c in pendingItems" :key="c.id" class="glass-strong rounded-2xl p-5">
          <div class="mb-3 flex items-center gap-2">
            <span
              class="rounded-lg px-2 py-0.5 text-[10px] font-semibold tracking-wider uppercase"
              :class="typeBadgeClass(c.contribution_type)"
              >{{ c.contribution_type }}</span
            >
            <router-link
              :to="{ name: c.target_type, params: { id: c.target_id } }"
              class="text-xs text-blue-400/60 hover:text-blue-400"
            >
              {{ targetLabel(c.target_type) }} · {{ c.target_id }}
            </router-link>
            <span
              v-if="c.ai_verdict"
              class="rounded-sm bg-white/5 px-1.5 py-0.5 text-[10px] text-white/30"
              >AI: {{ c.ai_verdict }} ({{ ((c.ai_confidence ?? 0) * 100).toFixed(0) }}%)</span
            >
          </div>

          <pre
            v-if="showData === c.id"
            class="mb-3 overflow-x-auto rounded-lg bg-black/40 p-3 font-mono text-[11px] text-white/40"
            >{{ formatContributionData(c.data) }}</pre
          >
          <button
            type="button"
            class="mb-3 text-xs text-white/30 transition-colors hover:text-white/50"
            @click="showData = showData === c.id ? null : c.id"
            aria-label="Toggle contribution data"
          >
            <ChevronRight aria-hidden="true" class="mr-1 text-[10px]"<i
              
              :class="{ 'rotate-90': showData === c.id }"
            /> />
            {{ showData === c.id ? 'Hide' : 'View' }} data
          </button>

          <div class="flex items-center gap-3">
            <button
              type="button"
              :disabled="reviewingId === c.id"
              class="spring rounded-lg bg-spotify/10 px-4 py-2 text-xs font-medium text-spotify transition-all hover:bg-spotify/20 disabled:opacity-40"
              @click="reviewContribution(c.id, 'approve')"
            >
              <Loader2 class="mr-1"<i aria-hidden="true" v-if="reviewingId === c.id"  /> />
              Approve
            </button>
            <button
              type="button"
              :disabled="reviewingId === c.id"
              class="spring rounded-lg bg-red-500/10 px-4 py-2 text-xs font-medium text-red-400 transition-all hover:bg-red-500/20 disabled:opacity-40"
              @click="reviewContribution(c.id, 'reject')"
            >
              <Loader2 class="mr-1"<i aria-hidden="true" v-if="reviewingId === c.id"  /> />
              Reject
            </button>
            <span class="text-[10px] text-white/20">{{ formatDate(c.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  CheckCircle,
  ChevronRight,
  History,
  Inbox,
  Loader2,
  Shield,
  Upload,
} from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { useUserAuthStore } from '@/stores'
import { useAppToast } from '@/composables/useAppToast'
import { SkeletonLoader } from '@/components/common'
import { useContributionApi } from '@/services/api/contribution'
import { useReputationApi } from '@/services/api/reputation'
import type {
  Contribution,
  ContributionHistoryItem,
  ContributorStats,
} from '@/services/api/contribution'
import ContributionSubmitForm from '@/components/contribution/ContributionSubmitForm.vue'
import ContributionHistory from '@/components/contribution/ContributionHistory.vue'
import ContributionLeaderboard from '@/components/contribution/ContributionLeaderboard.vue'

const api = useContributionApi()
const reputationApi = useReputationApi()
const auth = useUserAuthStore()
const toast = useAppToast()

const reputation = ref<Record<string, any> | null>(null)

const tabs = [
  { key: 'submit', label: 'Submit', icon: 'pi pi-plus' },
  { key: 'my', label: 'My Contributions', icon: 'pi pi-inbox' },
  { key: 'leaderboard', label: 'Leaderboard', icon: 'pi pi-trophy' },
  { key: 'moderate', label: 'Moderate', icon: 'pi pi-shield' },
]

const activeTab = ref('submit')

const myContributions = ref<Contribution[]>([])
const loadingMy = ref(false)
const myMeta = ref<{ page: number; page_size: number; total: number } | null>(null)

const pendingItems = ref<Contribution[]>([])
const loadingPending = ref(false)

const leaderboardData = ref<ContributorStats[]>([])
const loadingLeaderboard = ref(false)

const expandedContribution = ref<string | null>(null)
const historyItems = ref<ContributionHistoryItem[]>([])
const loadingHistory = ref(false)

const showData = ref<string | null>(null)
const reviewingId = ref<string | null>(null)

function loadTab(key: string) {
  if (key === 'my') loadMyContributions()
  if (key === 'moderate') loadPending()
  if (key === 'leaderboard') loadLeaderboard()
}

async function loadMyContributions(page = 1) {
  loadingMy.value = true
  try {
    const res = await api.listMy({ page, page_size: 20 })
    myContributions.value = res?.data ?? []
    myMeta.value = res?.meta ?? null
  } catch (err: unknown) {
    toast.apiError(err, 'Failed to load contributions')
    myContributions.value = []
  } finally {
    loadingMy.value = false
  }
}

async function loadPending() {
  loadingPending.value = true
  try {
    const res = await api.listPending({ page: 1, page_size: 50 })
    pendingItems.value = res?.data ?? []
  } catch (err: unknown) {
    toast.apiError(err, 'Failed to load pending')
    pendingItems.value = []
  } finally {
    loadingPending.value = false
  }
}

async function loadLeaderboard() {
  loadingLeaderboard.value = true
  try {
    const res = await api.getLeaderboard({ limit: 20 })
    leaderboardData.value = res ?? []
  } catch (err: unknown) {
    toast.apiError(err, 'Failed to load leaderboard')
    leaderboardData.value = []
  } finally {
    loadingLeaderboard.value = false
  }
}

async function viewHistory(c: Contribution) {
  if (expandedContribution.value === c.id) {
    expandedContribution.value = null
    return
  }
  expandedContribution.value = c.id
  loadingHistory.value = true
  try {
    const res = await api.getHistory(c.id)
    historyItems.value = res ?? []
  } catch (err: unknown) {
    toast.apiError(err, 'Failed to load history')
    historyItems.value = []
  } finally {
    loadingHistory.value = false
  }
}

async function reviewContribution(id: string, action: 'approve' | 'reject') {
  reviewingId.value = id
  try {
    await api.review(id, { action })
    pendingItems.value = pendingItems.value.filter((c) => c.id !== id)
    toast.success(`Contribution ${action}d`)
  } catch (err: unknown) {
    toast.apiError(err, `Failed to ${action} contribution`)
  } finally {
    reviewingId.value = null
  }
}

function onSubmitted() {
  loadMyContributions()
}

function typeBadgeClass(type: string) {
  const classes: Record<string, string> = {
    lyrics: 'bg-spotify/10 text-spotify',
    translation: 'bg-aurora-blue/10 text-aurora-blue',
    credits: 'bg-aurora-purple/10 text-aurora-purple',
    metadata: 'bg-amber-500/10 text-amber-500',
    album_art: 'bg-aurora-pink/10 text-aurora-pink',
    bio: 'bg-emerald-400/10 text-emerald-400',
  }
  return classes[type] || 'bg-white/5 text-white/40'
}

function statusBadgeClass(status: string) {
  switch (status) {
    case 'approved':
      return 'bg-spotify/10 text-spotify'
    case 'rejected':
      return 'bg-red-500/10 text-red-400'
    case 'pending':
      return 'bg-amber-500/10 text-amber-500'
    case 'needs_review':
      return 'bg-aurora-blue/10 text-aurora-blue'
    default:
      return 'bg-white/5 text-white/40'
  }
}

function targetLabel(type: string) {
  return type.charAt(0).toUpperCase() + type.slice(1)
}

function formatDate(dateStr: string) {
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatContributionData(data: unknown): string {
  if (typeof data === 'string') {
    try {
      return JSON.stringify(JSON.parse(data), null, 2)
    } catch {
      return data
    }
  }
  return JSON.stringify(data, null, 2)
}

function tierBgClass(tier?: string) {
  const map: Record<string, string> = {
    newcomer: 'bg-white/5',
    contributor: 'bg-blue-500/10',
    trusted: 'bg-emerald-500/10',
    verified: 'bg-purple-500/10',
    elite: 'bg-yellow-500/10',
    legend: 'bg-gradient-to-br from-yellow-500/20 to-amber-500/20',
  }
  return map[tier?.toLowerCase() ?? ''] || 'bg-white/5'
}

function tierIconClass(tier?: string) {
  const map: Record<string, string> = {
    newcomer: 'text-white/40',
    contributor: 'text-blue-400',
    trusted: 'text-emerald-400',
    verified: 'text-purple-400',
    elite: 'text-yellow-400',
    legend: 'text-amber-400',
  }
  return map[tier?.toLowerCase() ?? ''] || 'text-white/40'
}

onMounted(() => {
  loadTab('submit')
  if (auth.isAuthenticated) {
    reputationApi.getUserReputation(String(auth.user?.id)).then((r) => {
      reputation.value = r as Record<string, any>
    }).catch(() => {})
  }
})
</script>
