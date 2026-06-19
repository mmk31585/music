<template>
  <div class="mx-auto max-w-6xl space-y-8 px-4 py-6 md:px-6 md:py-8">
    <!-- Header -->
    <div>
      <h1 class="text-3xl font-bold text-white md:text-4xl">Contributions</h1>
      <p class="mt-1 text-sm text-white/40">
        Help improve the community by contributing lyrics, translations, credits, and more
      </p>
    </div>

    <!-- Tabs -->
    <div class="flex items-center gap-1 rounded-xl bg-white/[0.04] p-1">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        class="spring flex-1 rounded-lg px-4 py-2 text-sm font-medium transition-all"
        :class="
          activeTab === tab.key
            ? 'bg-white/10 text-white shadow-sm'
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
          <i aria-hidden="true" class="pi pi-inbox text-2xl text-white/15" />
        </div>
        <p class="text-sm text-white/25">No contributions yet</p>
        <button
          type="button"
          class="text-xs text-[#1db954] transition-colors hover:text-[#1ed760]"
          @click="activeTab = 'submit'"
        >
          Make your first contribution
        </button>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="c in myContributions"
          :key="c.id"
          class="contribution-card glass-strong spring rounded-2xl p-5 transition-all hover:bg-white/[0.06]"
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
                <span v-if="c.xp_awarded > 0" class="text-[#1db954]">+{{ c.xp_awarded }} XP</span>
              </div>
            </div>

            <div class="flex shrink-0 items-center gap-2">
              <button
                type="button"
                class="spring flex h-8 w-8 items-center justify-center rounded-full text-white/30 transition-all hover:bg-white/10 hover:text-white/60"
                @click="viewHistory(c)"
                aria-label="View history"
              >
                <i aria-hidden="true" class="pi pi-history text-xs" />
              </button>
            </div>
          </div>

          <!-- Expanded history -->
          <div v-if="expandedContribution === c.id" class="mt-4 border-t border-white/[0.06] pt-4">
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
          <i aria-hidden="true" class="pi pi-check-circle text-2xl text-white/15" />
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
              class="rounded bg-white/5 px-1.5 py-0.5 text-[10px] text-white/30"
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
            <i
              class="pi pi-chevron-right mr-1 text-[10px]"
              :class="{ 'rotate-90': showData === c.id }"
            />
            {{ showData === c.id ? 'Hide' : 'View' }} data
          </button>

          <div class="flex items-center gap-3">
            <button
              type="button"
              :disabled="reviewingId === c.id"
              class="spring rounded-lg bg-[#1db954]/10 px-4 py-2 text-xs font-medium text-[#1db954] transition-all hover:bg-[#1db954]/20 disabled:opacity-40"
              @click="reviewContribution(c.id, 'approve')"
            >
              <i aria-hidden="true" v-if="reviewingId === c.id" class="pi pi-spin pi-spinner mr-1" />
              Approve
            </button>
            <button
              type="button"
              :disabled="reviewingId === c.id"
              class="spring rounded-lg bg-red-500/10 px-4 py-2 text-xs font-medium text-red-400 transition-all hover:bg-red-500/20 disabled:opacity-40"
              @click="reviewContribution(c.id, 'reject')"
            >
              <i aria-hidden="true" v-if="reviewingId === c.id" class="pi pi-spin pi-spinner mr-1" />
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
import { onMounted, ref } from 'vue'
import { useToast } from 'primevue/usetoast'
import { SkeletonLoader } from '@/components/common'
import { useContributionApi } from '@/services/api/contribution'
import type {
  Contribution,
  ContributionHistoryItem,
  ContributorStats,
} from '@/services/api/contribution'
import ContributionSubmitForm from '@/components/contribution/ContributionSubmitForm.vue'
import ContributionHistory from '@/components/contribution/ContributionHistory.vue'
import ContributionLeaderboard from '@/components/contribution/ContributionLeaderboard.vue'

const api = useContributionApi()
const toast = useToast()

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
    myContributions.value = (res as Record<string, unknown>)?.data as Contribution[] ?? []
    myMeta.value = (res as Record<string, unknown>)?.meta as { page: number; page_size: number; total: number } | null ?? null
  } catch (err) {
    console.error('Failed to load contributions:', err)
    myContributions.value = []
  } finally {
    loadingMy.value = false
  }
}

async function loadPending() {
  loadingPending.value = true
  try {
    const res = await api.listPending({ page: 1, page_size: 50 })
    pendingItems.value = (res as Record<string, unknown>)?.data as Contribution[] ?? []
  } catch (err) {
    console.error('Failed to load pending:', err)
    pendingItems.value = []
  } finally {
    loadingPending.value = false
  }
}

async function loadLeaderboard() {
  loadingLeaderboard.value = true
  try {
    const res = await api.getLeaderboard({ limit: 20 })
    leaderboardData.value = (res as ContributorStats[]) ?? []
  } catch (err) {
    console.error('Failed to load leaderboard:', err)
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
    historyItems.value = (res as ContributionHistoryItem[]) ?? []
  } catch (err) {
    console.error('Failed to load history:', err)
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
    toast.add({ severity: 'success', summary: `Contribution ${action}d`, life: 2000 })
  } catch (err) {
    console.error(`Failed to ${action} contribution:`, err)
    toast.add({ severity: 'error', summary: `Failed to ${action} contribution`, life: 3000 })
  } finally {
    reviewingId.value = null
  }
}

function onSubmitted() {
  loadMyContributions()
}

function typeBadgeClass(type: string) {
  const classes: Record<string, string> = {
    lyrics: 'bg-[#1db954]/10 text-[#1db954]',
    translation: 'bg-[#60a5fa]/10 text-[#60a5fa]',
    credits: 'bg-[#a855f7]/10 text-[#a855f7]',
    metadata: 'bg-[#f59e0b]/10 text-[#f59e0b]',
    album_art: 'bg-[#f472b6]/10 text-[#f472b6]',
    bio: 'bg-[#34d399]/10 text-[#34d399]',
  }
  return classes[type] || 'bg-white/5 text-white/40'
}

function statusBadgeClass(status: string) {
  switch (status) {
    case 'approved':
      return 'bg-[#1db954]/10 text-[#1db954]'
    case 'rejected':
      return 'bg-red-500/10 text-red-400'
    case 'pending':
      return 'bg-[#f59e0b]/10 text-[#f59e0b]'
    case 'needs_review':
      return 'bg-[#60a5fa]/10 text-[#60a5fa]'
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

onMounted(() => {
  loadTab('submit')
})
</script>
