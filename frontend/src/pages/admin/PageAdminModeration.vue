<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      title="Moderation"
      eyebrow="Content reports"
      description="Manage reported content, flags, and audit activity"
    >
      <template #actions>
        <div class="flex items-center gap-2">
          <button
            type="button"
            class="inline-flex items-center gap-1.5 rounded-lg border border-white/10 bg-white/[0.04] px-3 py-1.5 text-xs font-medium text-white/60 backdrop-blur transition hover:bg-white/[0.08] disabled:opacity-40"
            :disabled="refreshing"
            @click="refreshAll"
          >
            <i aria-hidden="true" :class="refreshing ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs" />
            {{ refreshing ? 'Loading...' : 'Refresh' }}
          </button>
        </div>
      </template>
    </AdminSectionHeader>

    <!-- Tabs -->
    <div class="mt-6 flex gap-1 rounded-xl bg-white/[0.04] p-1">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="flex items-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium transition-all duration-200"
        :class="
          activeTab === tab.key
            ? 'bg-white/10 text-white shadow-lg'
            : 'text-white/30 hover:text-white/50'
        "
        @click="activeTab = tab.key"
      >
        <i aria-hidden="true" :class="tab.icon" class="text-xs" />
        {{ tab.label }}
        <span
          v-if="tab.badge"
          class="rounded-full bg-red-500/15 px-1.5 py-0.5 text-[10px] font-bold text-red-400"
          >{{ tab.badge }}</span
        >
      </button>
    </div>

    <!-- === TAB: QUEUE === -->
    <div v-show="activeTab === 'queue'" class="mt-6 space-y-4">
      <!-- Bar: filter + bulk -->
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="flex items-center gap-2">
          <button
            v-for="f in typeFilters"
            :key="f"
            class="rounded-lg px-3 py-1.5 text-xs font-medium transition"
            :class="
              typeFilter === f
                ? 'bg-white/10 text-white'
                : 'bg-white/[0.04] text-white/40 hover:text-white/60'
            "
            @click="typeFilter = f; fetchQueue()"
          >
            {{ f === 'all' ? 'All Types' : f }}
          </button>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-xs text-white/30">{{ selectedIds.length }} selected</span>
          <button
            v-if="selectedIds.length > 0"
            class="rounded-lg bg-green-500/10 px-3 py-1.5 text-xs font-medium text-green-400 transition hover:bg-green-500/20"
            @click="bulkResolve('resolved')"
          >
            Resolve
          </button>
          <button
            v-if="selectedIds.length > 0"
            class="rounded-lg bg-amber-500/10 px-3 py-1.5 text-xs font-medium text-amber-400 transition hover:bg-amber-500/20"
            @click="bulkResolve('dismissed')"
          >
            Dismiss
          </button>
        </div>
      </div>

      <!-- Loading state -->
      <div v-if="queueLoading" class="space-y-2">
        <div v-for="i in 4" :key="i" class="h-24 animate-pulse rounded-xl bg-white/[0.06]" />
      </div>

      <!-- Empty state -->
      <AdminEmptyState
        v-else-if="queueItems.length === 0"
        title="Queue is clear"
        description="No pending reports. All content is clean."
        icon="pi pi-check-circle"
      />

      <!-- Queue items -->
      <div v-else class="space-y-2">
        <div
          v-for="report in filteredQueue"
          :key="report.id"
          class="group rounded-xl border border-white/[0.05] bg-white/[0.02] p-4 transition hover:bg-white/[0.04]"
          :class="{ 'border-l-2 border-l-red-500/40': isHighPriority(report) }"
        >
          <div class="flex items-start gap-3">
            <!-- Checkbox -->
            <label class="mt-1 flex shrink-0 cursor-pointer items-center">
              <input
                type="checkbox"
                :checked="selectedIds.includes(report.id)"
                class="h-4 w-4 rounded border-white/20 bg-white/5 accent-[#1db954]"
                @change="toggleSelect(report.id)"
              />
            </label>
            <!-- Content -->
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center gap-2">
                <span
                  class="rounded-full bg-red-500/10 px-2 py-0.5 text-[10px] font-bold text-red-400 uppercase"
                  >{{ report.reason }}</span
                >
                <span
                  class="rounded-full bg-white/[0.06] px-2 py-0.5 text-[10px] font-medium text-white/40"
                  >{{ report.target_type }}</span
                >
                <span class="font-mono text-[10px] text-white/20"
                  >#{{ shortId(report.target_id) }}</span
                >
                <span
                  v-if="isHighPriority(report)"
                  class="rounded-full bg-red-500/15 px-2 py-0.5 text-[10px] font-bold text-red-400"
                  >HIGH</span
                >
              </div>
              <p v-if="report.description" class="mt-1.5 line-clamp-2 text-sm text-white/70">
                {{ report.description }}
              </p>
              <p class="mt-1.5 text-xs text-white/30">
                <span class="font-medium text-white/40">reported</span>
                {{ timeAgo(report.created_at) }}
                <span class="mx-1">·</span>
                <span class="font-medium text-white/40">by</span> {{ shortId(report.reporter_id) }}
              </p>
              <!-- Inline resolution note input (shown when resolving) -->
              <div v-if="resolvingId === report.id" class="mt-3 flex items-center gap-2">
                <input
                  v-model="resolveNote"
                  placeholder="Resolution note (optional)"
                  class="flex-1 rounded-lg border border-white/10 bg-white/5 px-3 py-1.5 text-xs text-white outline-none focus:border-white/20"
                  @keyup.esc="resolvingId = null"
                />
                <button
                  class="rounded-lg bg-green-500/10 px-3 py-1.5 text-xs text-green-400 hover:bg-green-500/20"
                  @click="confirmResolve(report.id, 'resolved')"
                >
                  Resolve
                </button>
                <button
                  class="rounded-lg bg-amber-500/10 px-3 py-1.5 text-xs text-amber-400 hover:bg-amber-500/20"
                  @click="confirmResolve(report.id, 'dismissed')"
                >
                  Dismiss
                </button>
                <button
                  class="text-xs text-white/30 hover:text-white/50"
                  @click="resolvingId = null"
                >
                  Cancel
                </button>
              </div>
            </div>
            <!-- Actions -->
            <div class="flex shrink-0 gap-1.5">
              <button
                v-if="resolvingId !== report.id"
                title="Resolve"
                class="rounded-lg bg-white/5 p-2 text-xs text-white/30 transition hover:bg-green-500/10 hover:text-green-400"
                @click="resolvingId = report.id; resolveNote = ''"
              >
                <i aria-hidden="true" class="pi pi-check" />
              </button>
              <button
                title="Flag content"
                class="rounded-lg bg-white/5 p-2 text-xs text-white/30 transition hover:bg-red-500/10 hover:text-red-400"
                @click="openFlagDialog(report)"
              >
                <i aria-hidden="true" class="pi pi-flag" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- === TAB: FLAGGED === -->
    <div v-show="activeTab === 'flagged'" class="mt-6 space-y-4">
      <div class="flex items-center gap-2">
        <label class="flex items-center gap-2 text-xs text-white/40">
          <input
            v-model="includeExpiredFlags"
            type="checkbox"
            class="h-4 w-4 rounded border-white/20 bg-white/5 accent-[#1db954]"
            @change="fetchFlags"
          />
          Include expired
        </label>
      </div>

      <div v-if="flagsLoading" class="space-y-2">
        <div v-for="i in 3" :key="i" class="h-16 animate-pulse rounded-xl bg-white/[0.06]" />
      </div>

      <AdminEmptyState
        v-else-if="flagItems.length === 0"
        title="No flagged content"
        description="No content flags active right now."
        icon="pi pi-flag"
      />

      <div v-else class="space-y-2">
        <div
          v-for="flag in flagItems"
          :key="flag.id"
          class="flex items-center justify-between rounded-xl border border-white/[0.05] bg-white/[0.02] px-4 py-3 transition hover:bg-white/[0.04]"
          :class="{ 'opacity-40': isExpired(flag) }"
        >
          <div class="flex min-w-0 items-center gap-3">
            <i aria-hidden="true" class="pi pi-flag text-sm text-red-400/60" />
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <span
                  class="rounded-full bg-white/[0.06] px-2 py-0.5 text-[10px] font-medium text-white/40"
                  >{{ flag.target_type }}</span
                >
                <span class="text-xs font-medium text-white/60">{{ flag.flag_type }}</span>
                <span class="font-mono text-[10px] text-white/20"
                  >#{{ shortId(flag.target_id) }}</span
                >
              </div>
              <p class="mt-0.5 text-[10px] text-white/30">
                Flagged {{ timeAgo(flag.flagged_at) }}
                <span v-if="flag.expires_at"> · expires {{ timeAgo(flag.expires_at) }}</span>
                <span v-else> · permanent</span>
              </p>
            </div>
          </div>
          <span v-if="isExpired(flag)" class="text-[10px] font-medium text-white/20">Expired</span>
        </div>
      </div>
    </div>

    <!-- === TAB: HISTORY === -->
    <div v-show="activeTab === 'history'" class="mt-6 space-y-4">
      <div class="flex items-center gap-2">
        <button
          v-for="s in historyStatusFilters"
          :key="s"
          class="rounded-lg px-3 py-1.5 text-xs font-medium transition"
          :class="
            historyFilter === s
              ? 'bg-white/10 text-white'
              : 'bg-white/[0.04] text-white/40 hover:text-white/60'
          "
          @click="historyFilter = s; fetchHistory()"
        >
          {{ s === 'all' ? 'All' : s.charAt(0).toUpperCase() + s.slice(1) }}
        </button>
      </div>

      <div v-if="historyLoading" class="space-y-2">
        <div v-for="i in 3" :key="i" class="h-20 animate-pulse rounded-xl bg-white/[0.06]" />
      </div>

      <AdminEmptyState
        v-else-if="historyItems.length === 0"
        title="No history"
        description="No resolved or dismissed reports yet."
        icon="pi pi-history"
      />

      <div v-else class="space-y-2">
        <div
          v-for="report in historyItems"
          :key="report.id"
          class="rounded-xl border border-white/[0.05] bg-white/[0.02] px-4 py-3 transition hover:bg-white/[0.04]"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center gap-2">
                <span
                  class="rounded-full px-2 py-0.5 text-[10px] font-bold uppercase"
                  :class="
                    report.status === 'resolved'
                      ? 'bg-green-500/10 text-green-400'
                      : 'bg-slate-500/10 text-slate-400'
                  "
                >
                  {{ report.status }}
                </span>
                <span
                  class="rounded-full bg-red-500/10 px-2 py-0.5 text-[10px] font-bold text-red-400"
                  >{{ report.reason }}</span
                >
                <span
                  class="rounded-full bg-white/[0.06] px-2 py-0.5 text-[10px] font-medium text-white/40"
                  >{{ report.target_type }}</span
                >
              </div>
              <p v-if="report.description" class="mt-1 line-clamp-1 text-xs text-white/50">
                {{ report.description }}
              </p>
              <p class="mt-1 text-[10px] text-white/20">
                Resolved {{ timeAgo(report.resolved_at || report.created_at) }}
                <span v-if="report.resolution_note"> · Note: {{ report.resolution_note }}</span>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- === TAB: STATS === -->
    <div v-show="activeTab === 'stats'" class="mt-6 space-y-6">
      <div v-if="statsLoading" class="grid grid-cols-2 gap-4 md:grid-cols-4">
        <div v-for="i in 5" :key="i" class="h-28 animate-pulse rounded-2xl bg-white/[0.06]" />
      </div>

      <template v-else>
        <!-- Metric cards -->
        <div class="grid grid-cols-2 gap-4 md:grid-cols-4">
          <div
            class="rounded-2xl bg-white/[0.04] p-5 backdrop-blur transition hover:bg-white/[0.06]"
          >
            <p class="text-[10px] font-medium tracking-wider text-white/30 uppercase">
              Total Reports
            </p>
            <p class="mt-2 text-3xl font-black text-white tabular-nums">
              {{ n(stats.total_reports) }}
            </p>
          </div>
          <div
            class="rounded-2xl bg-white/[0.04] p-5 backdrop-blur transition hover:bg-white/[0.06]"
          >
            <p class="text-[10px] font-medium tracking-wider text-white/30 uppercase">Pending</p>
            <p class="mt-2 text-3xl font-black text-red-400 tabular-nums">
              {{ n(stats.pending_reports) }}
            </p>
          </div>
          <div
            class="rounded-2xl bg-white/[0.04] p-5 backdrop-blur transition hover:bg-white/[0.06]"
          >
            <p class="text-[10px] font-medium tracking-wider text-white/30 uppercase">
              Resolved Today
            </p>
            <p class="mt-2 text-3xl font-black text-green-400 tabular-nums">
              {{ n(stats.resolved_today) }}
            </p>
          </div>
          <div
            class="rounded-2xl bg-white/[0.04] p-5 backdrop-blur transition hover:bg-white/[0.06]"
          >
            <p class="text-[10px] font-medium tracking-wider text-white/30 uppercase">
              Flagged Content
            </p>
            <p class="mt-2 text-3xl font-black text-amber-400 tabular-nums">
              {{ n(stats.flagged_content) }}
            </p>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4 md:grid-cols-3">
          <div
            class="rounded-2xl bg-white/[0.04] p-5 backdrop-blur transition hover:bg-white/[0.06]"
          >
            <p class="text-[10px] font-medium tracking-wider text-white/30 uppercase">
              Unique Reporters
            </p>
            <p class="mt-2 text-2xl font-black text-white tabular-nums">
              {{ n(stats.unique_reporters) }}
            </p>
          </div>
          <div
            class="rounded-2xl bg-white/[0.04] p-5 backdrop-blur transition hover:bg-white/[0.06]"
          >
            <p class="text-[10px] font-medium tracking-wider text-white/30 uppercase">
              Avg Resolution
            </p>
            <p class="mt-2 text-2xl font-black text-white tabular-nums">
              {{ formatHours(stats.avg_resolution_hours) }}
            </p>
          </div>
        </div>

        <!-- By Reason -->
        <section v-if="reasonEntries.length > 0">
          <h4 class="mb-3 text-sm font-bold text-white/60">Reports by Reason</h4>
          <div class="rounded-2xl bg-white/[0.04] p-6">
            <div class="space-y-3">
              <div
                v-for="[reason, count] in reasonEntries"
                :key="reason"
                class="flex items-center gap-3"
              >
                <span class="w-32 shrink-0 truncate text-xs text-slate-400">{{ reason }}</span>
                <div class="h-5 flex-1 overflow-hidden rounded-full bg-white/[0.06]">
                  <div
                    class="h-full rounded-full bg-gradient-to-r from-red-400 to-amber-400 transition-all duration-500"
                    :style="{ width: reasonPercent(count) + '%' }"
                  />
                </div>
                <span class="w-12 text-right text-xs font-medium text-white/80 tabular-nums">{{
                  count
                }}</span>
              </div>
            </div>
          </div>
        </section>

        <!-- By Target Type -->
        <section v-if="typeEntries.length > 0">
          <h4 class="mb-3 text-sm font-bold text-white/60">Reports by Target Type</h4>
          <div class="rounded-2xl bg-white/[0.04] p-6">
            <div class="space-y-3">
              <div v-for="[tt, count] in typeEntries" :key="tt" class="flex items-center gap-3">
                <span class="w-24 shrink-0 text-xs text-slate-400">{{ tt }}</span>
                <div class="h-5 flex-1 overflow-hidden rounded-full bg-white/[0.06]">
                  <div
                    class="h-full rounded-full bg-gradient-to-r from-blue-400 to-purple-400 transition-all duration-500"
                    :style="{ width: reasonPercent(count) + '%' }"
                  />
                </div>
                <span class="w-12 text-right text-xs font-medium text-white/80 tabular-nums">{{
                  count
                }}</span>
              </div>
            </div>
          </div>
        </section>
      </template>
    </div>

    <!-- Flag dialog -->
    <Dialog
      v-model:visible="showFlagDialog"
      modal
      :draggable="false"
      :style="{ width: '440px' }"
      :pt="{
        root: { class: '!border-white/[0.06] !bg-[#141414] !rounded-2xl !shadow-2xl' },
        header: { class: '!bg-transparent !border-0 !pb-2' },
        content: { class: '!bg-transparent !px-6 !pt-0 !pb-2' },
        footer: { class: '!bg-transparent !border-0' },
        mask: { class: '!backdrop-blur-sm' },
      }"
    >
      <template #header>
        <div class="flex items-center gap-3">
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-red-500/10">
            <i aria-hidden="true" class="pi pi-flag text-red-400" />
          </div>
          <div>
            <h3 class="text-base font-semibold text-white">Flag Content</h3>
            <p class="text-xs text-slate-500">
              Add a flag to {{ flagTarget?.target_type || '' }} #{{ shortId(flagTarget?.target_id || '') }}
            </p>
          </div>
        </div>
      </template>

      <div class="mt-2 space-y-4">
        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-400">Flag Type</label>
          <select
            v-model="flagForm.flag_type"
            class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white outline-none focus:border-white/20"
          >
            <option value="inappropriate">Inappropriate</option>
            <option value="copyright">Copyright Violation</option>
            <option value="spam">Spam</option>
            <option value="misinformation">Misinformation</option>
            <option value="hate_speech">Hate Speech</option>
            <option value="explicit">Explicit Content</option>
          </select>
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-400">Duration</label>
          <select
            v-model.number="flagForm.expires_in_hours"
            class="w-full rounded-xl border border-white/10 bg-white/5 px-4 py-2.5 text-sm text-white outline-none focus:border-white/20"
          >
            <option :value="0">Permanent</option>
            <option :value="24">24 hours</option>
            <option :value="72">3 days</option>
            <option :value="168">7 days</option>
            <option :value="720">30 days</option>
          </select>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            label="Cancel"
            text
            class="!text-slate-400 hover:!text-white"
            @click="showFlagDialog = false"
          />
          <Button
            label="Flag Content"
            icon="pi pi-flag"
            class="!rounded-xl !bg-red-500/20 !text-red-400 !ring-1 !ring-red-500/20 hover:!bg-red-500/30"
            @click="confirmFlag"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useModerationApi } from '@/services/api/moderation'
import { useToast } from 'primevue/usetoast'
import type { ContentReport, ContentFlag, ModerationStats } from '@/services/api/moderation/types'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import { AdminSectionHeader, AdminEmptyState } from '@/components/admin'

const moderationApi = useModerationApi()
const toast = useToast()

const activeTab = ref('queue')
const refreshing = ref(false)

// Queue state
const queueItems = ref<ContentReport[]>([])
const queueLoading = ref(true)
const typeFilter = ref('all')
const typeFilters = ['all', 'track', 'album', 'artist', 'playlist', 'comment', 'user']
const selectedIds = ref<string[]>([])
const resolvingId = ref<string | null>(null)
const resolveNote = ref('')

// Flag state
const showFlagDialog = ref(false)
const flagTarget = ref<ContentReport | null>(null)
const flagForm = reactive({ flag_type: 'inappropriate', expires_in_hours: 0 })

// Flags state
const flagItems = ref<ContentFlag[]>([])
const flagsLoading = ref(true)
const includeExpiredFlags = ref(false)

// History state
const historyItems = ref<ContentReport[]>([])
const historyLoading = ref(true)
const historyFilter = ref('resolved')
const historyStatusFilters = ['all', 'resolved', 'dismissed']

// Stats state
const statsLoading = ref(true)
const stats = reactive<ModerationStats>({
  total_reports: 0,
  pending_reports: 0,
  resolved_today: 0,
  flagged_content: 0,
  unique_reporters: 0,
  avg_resolution_hours: 0,
  by_reason: {},
  by_target_type: {},
})

const tabs = computed(() => [
  {
    key: 'queue',
    label: 'Queue',
    icon: 'pi pi-inbox',
    badge: queueItems.value.length || undefined,
  },
  {
    key: 'flagged',
    label: 'Flagged',
    icon: 'pi pi-flag',
    badge: flagItems.value.length || undefined,
  },
  { key: 'history', label: 'History', icon: 'pi pi-history' },
  { key: 'stats', label: 'Stats', icon: 'pi pi-chart-bar' },
])

const filteredQueue = computed(() => {
  if (typeFilter.value === 'all') return queueItems.value
  return queueItems.value.filter((r) => r.target_type === typeFilter.value)
})

const reasonEntries = computed(() =>
  Object.entries(stats.by_reason ?? {}).sort((a, b) => (b[1] as number) - (a[1] as number)),
)
const typeEntries = computed(() =>
  Object.entries(stats.by_target_type ?? {}).sort((a, b) => (b[1] as number) - (a[1] as number)),
)

function isHighPriority(r: ContentReport): boolean {
  return ['copyright', 'hate_speech', 'explicit'].includes(r.reason)
}

function shortId(id: string): string {
  return id.length > 8 ? id.slice(0, 8) + '…' : id
}

function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  return `${days}d ago`
}

function isExpired(flag: ContentFlag): boolean {
  if (!flag.expires_at) return false
  return new Date(flag.expires_at) < new Date()
}

function n(v: number): string {
  if (v >= 1_000_000) return `${(v / 1_000_000).toFixed(1)}M`
  if (v >= 1_000) return `${(v / 1_000).toFixed(1)}K`
  return String(v)
}

function formatHours(h: number): string {
  if (h < 1) return `${Math.round(h * 60)}m`
  return `${h.toFixed(1)}h`
}

function reasonPercent(count: number): number {
  const total = stats.total_reports || 1
  return (count / total) * 100
}

function toggleSelect(id: string) {
  const idx = selectedIds.value.indexOf(id)
  if (idx >= 0) selectedIds.value.splice(idx, 1)
  else selectedIds.value.push(id)
}

async function fetchQueue() {
  queueLoading.value = true
  try {
    const data = await moderationApi.getPendingReports({ limit: 50 })
    queueItems.value = data?.items ?? []
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to load queue', life: 3000 })
  } finally {
    queueLoading.value = false
  }
}

async function fetchFlags() {
  flagsLoading.value = true
  try {
    const data = await moderationApi.getFlags({ limit: 50, expired: includeExpiredFlags.value })
    flagItems.value = data?.items ?? []
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to load flags', life: 3000 })
  } finally {
    flagsLoading.value = false
  }
}

async function fetchHistory() {
  historyLoading.value = true
  try {
    if (historyFilter.value === 'all') {
      const [resolved, dismissed] = await Promise.all([
        moderationApi.getReportsByStatus('resolved', { limit: 25 }),
        moderationApi.getReportsByStatus('dismissed', { limit: 25 }),
      ])
      historyItems.value = [...(resolved?.items || []), ...(dismissed?.items || [])].sort(
        (a, b) =>
          new Date(b.resolved_at || b.created_at).getTime() -
          new Date(a.resolved_at || a.created_at).getTime(),
      )
    } else {
      const data = await moderationApi.getReportsByStatus(historyFilter.value, { limit: 50 })
      historyItems.value = data?.items ?? []
    }
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to load history', life: 3000 })
  } finally {
    historyLoading.value = false
  }
}

async function fetchStats() {
  statsLoading.value = true
  try {
    const data = await moderationApi.getStats()
    if (data?.stats) Object.assign(stats, data.stats)
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to load stats', life: 3000 })
  } finally {
    statsLoading.value = false
  }
}

async function confirmResolve(id: string, status: string) {
  try {
    await moderationApi.resolveReport(id, status, resolveNote.value)
    queueItems.value = queueItems.value.filter((r) => r.id !== id)
    selectedIds.value = selectedIds.value.filter((s) => s !== id)
    resolvingId.value = null
    resolveNote.value = ''
    toast.add({ severity: 'success', summary: `Report ${status}`, life: 2000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to resolve', life: 3000 })
  }
}

async function bulkResolve(action: string) {
  if (selectedIds.value.length === 0) return
  try {
    await moderationApi.bulkAction(selectedIds.value, action)
    queueItems.value = queueItems.value.filter((r) => !selectedIds.value.includes(r.id))
    toast.add({
      severity: 'success',
      summary: `${selectedIds.value.length} reports ${action}`,
      life: 2000,
    })
    selectedIds.value = []
  } catch {
    toast.add({ severity: 'error', summary: 'Bulk action failed', life: 3000 })
  }
}

function openFlagDialog(report: ContentReport) {
  flagTarget.value = report
  flagForm.flag_type = 'inappropriate'
  flagForm.expires_in_hours = 0
  showFlagDialog.value = true
}

async function confirmFlag() {
  if (!flagTarget.value) return
  try {
    await moderationApi.flagContent({
      target_id: flagTarget.value.target_id,
      target_type: flagTarget.value.target_type,
      flag_type: flagForm.flag_type,
      expires_in_hours: flagForm.expires_in_hours > 0 ? flagForm.expires_in_hours : undefined,
    })
    showFlagDialog.value = false
    toast.add({ severity: 'success', summary: 'Content flagged', life: 2000 })
  } catch {
    toast.add({ severity: 'error', summary: 'Failed to flag content', life: 3000 })
  }
}

async function refreshAll() {
  refreshing.value = true
  await Promise.all([fetchQueue(), fetchFlags(), fetchHistory(), fetchStats()])
  refreshing.value = false
}

onMounted(async () => {
  await Promise.all([fetchQueue(), fetchFlags(), fetchHistory(), fetchStats()])
})
</script>
