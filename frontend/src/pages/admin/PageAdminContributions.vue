<template>
  <div class="mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      title="Contributions"
      eyebrow="User submitted content"
      description="Review, approve, reject, and apply contributions from users"
    >
      <template #actions>
        <div class="flex items-center gap-2">
          <span class="rounded-full bg-blue-500/10 px-3 py-1 text-xs font-medium text-blue-400"
            >{{ pendingCount }} pending</span
          >
          <button
            type="button"
            class="inline-flex items-center gap-1.5 rounded-lg border border-white/10 bg-white/[0.04] px-3 py-1.5 text-xs font-medium text-white/60 backdrop-blur transition hover:bg-white/[0.08] disabled:opacity-40"
            :disabled="loading"
            @click="refreshAll"
          >
            <i :class="loading ? 'pi pi-spin pi-spinner' : 'pi pi-refresh'" class="text-xs" />
            {{ loading ? 'Loading...' : 'Refresh' }}
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
        <i :class="tab.icon" class="text-xs" />
        {{ tab.label }}
        <span
          v-if="tab.badge"
          class="rounded-full bg-blue-500/15 px-1.5 py-0.5 text-[10px] font-bold text-blue-400"
          >{{ tab.badge }}</span
        >
      </button>
    </div>

    <!-- Search & Filter Bar -->
    <div class="mt-4 flex items-center gap-3">
      <div class="relative flex-1">
        <i class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-xs text-white/30" />
        <input
          v-model="searchQuery"
          placeholder="Search by summary or contributor name..."
          class="w-full rounded-lg border border-white/10 bg-white/[0.04] py-2 pl-9 pr-3 text-xs text-white outline-none placeholder:text-white/20 focus:border-white/20"
          @input="onSearchInput"
        />
      </div>
      <select
        v-model="typeFilter"
        class="rounded-lg border border-white/10 bg-white/[0.04] px-3 py-2 text-xs text-white/60 outline-none focus:border-white/20"
        @change="fetchCurrentTab"
      >
        <option value="">All types</option>
        <option value="lyrics">Lyrics</option>
        <option value="translation">Translation</option>
        <option value="credits">Credits</option>
        <option value="metadata">Metadata</option>
        <option value="album_art">Album Art</option>
        <option value="bio">Bio</option>
      </select>
    </div>

    <!-- === TAB: PENDING === -->
    <div v-show="activeTab === 'pending'" class="mt-6 space-y-4">
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 5" :key="i" class="h-24 animate-pulse rounded-xl bg-white/[0.06]" />
      </div>

      <AdminEmptyState
        v-else-if="pendingItems.length === 0"
        title="No pending contributions"
        description="All contributions have been reviewed."
        icon="pi pi-check-circle"
      />

      <div v-else class="space-y-2">
        <div
          v-for="c in pendingItems"
          :key="c.id"
          class="rounded-xl border border-white/[0.05] bg-white/[0.02] p-4 transition hover:bg-white/[0.04]"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center gap-2">
                <span
                  class="rounded-full px-2 py-0.5 text-[10px] font-bold uppercase"
                  :class="contributionTypeClass(c.contribution_type)"
                >
                  {{ c.contribution_type }}
                </span>
                <span
                  class="rounded-full bg-white/[0.06] px-2 py-0.5 text-[10px] font-medium text-white/40"
                  >{{ c.target_type }}</span
                >
                <router-link
                  :to="{ name: c.target_type, params: { id: c.target_id } }"
                  class="font-mono text-[10px] text-blue-400/60 hover:text-blue-400"
                  >#{{ shortId(c.target_id) }}</router-link
                >
              </div>
              <p class="mt-1.5 line-clamp-2 text-xs text-white/50">
                {{ c.summary || c.contribution_type + ' contribution' }}
              </p>
              <p class="mt-1.5 text-xs text-white/30">
                by <span class="font-medium text-white/40">{{ c.contributor_username || shortId(c.contributor_id) }}</span> ·
                {{ timeAgo(c.created_at) }}
                <span v-if="c.ai_verdict">
                  · AI: <span :class="aiVerdictClass(c.ai_verdict)">{{ c.ai_verdict }}</span>
                  <span v-if="c.ai_confidence">({{ (c.ai_confidence * 100).toFixed(0) }}%)</span>
                </span>
              </p>
            </div>
            <div class="flex shrink-0 gap-1.5">
              <button
                title="View details"
                class="rounded-lg bg-white/5 p-2 text-xs text-white/30 transition hover:bg-white/10 hover:text-white"
                @click="selected = c; showDetail = true"
              >
                <i class="pi pi-eye" />
              </button>
              <button
                v-if="reviewingId !== c.id"
                title="Review"
                class="rounded-lg bg-white/5 p-2 text-xs text-white/30 transition hover:bg-amber-500/10 hover:text-amber-400"
                @click="reviewingId = c.id; reviewNote = ''"
              >
                <i class="pi pi-check-circle" />
              </button>
            </div>
          </div>

          <!-- Inline review form -->
          <div v-if="reviewingId === c.id" class="mt-3 flex items-center gap-2">
            <input
              v-model="reviewNote"
              placeholder="Review note (optional)"
              class="flex-1 rounded-lg border border-white/10 bg-white/5 px-3 py-1.5 text-xs text-white outline-none focus:border-white/20"
              @keyup.esc="reviewingId = null"
            />
            <button
              class="rounded-lg bg-green-500/10 px-3 py-1.5 text-xs text-green-400 hover:bg-green-500/20"
              @click="doReview(c.id, 'approve')"
            >
              Approve
            </button>
            <button
              class="rounded-lg bg-red-500/10 px-3 py-1.5 text-xs text-red-400 hover:bg-red-500/20"
              @click="doReview(c.id, 'reject')"
            >
              Reject
            </button>
            <button class="text-xs text-white/30 hover:text-white/50" @click="reviewingId = null">
              Cancel
            </button>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="pendingTotal > pendingItems.length" class="flex justify-center pt-2">
          <button
            class="rounded-lg bg-white/5 px-4 py-2 text-xs text-white/50 hover:bg-white/10"
            @click="pendingPage++"
          >
            Load more
          </button>
        </div>
      </div>
    </div>

    <!-- === TAB: APPROVED === -->
    <div v-show="activeTab === 'approved'" class="mt-6 space-y-4">
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 3" :key="i" class="h-20 animate-pulse rounded-xl bg-white/[0.06]" />
      </div>

      <AdminEmptyState
        v-else-if="approvedItems.length === 0"
        title="No approved contributions"
        description="Approved contributions will appear here."
        icon="pi pi-check"
      />

      <div v-else class="space-y-2">
        <div
          v-for="c in approvedItems"
          :key="c.id"
          class="rounded-xl border border-white/[0.05] bg-white/[0.02] p-4 transition hover:bg-white/[0.04]"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center gap-2">
                <span
                  class="rounded-full bg-green-500/10 px-2 py-0.5 text-[10px] font-bold text-green-400 uppercase"
                  >Approved</span
                >
                <span
                  class="rounded-full px-2 py-0.5 text-[10px] font-bold uppercase"
                  :class="contributionTypeClass(c.contribution_type)"
                >
                  {{ c.contribution_type }}
                </span>
                <span
                  class="rounded-full bg-white/[0.06] px-2 py-0.5 text-[10px] font-medium text-white/40"
                  >{{ c.target_type }}</span
                >
              </div>
              <p class="mt-1 text-xs text-white/30">
                by {{ c.contributor_username || shortId(c.contributor_id) }} · {{ timeAgo(c.created_at) }}
              </p>
            </div>
            <div class="flex shrink-0 gap-1.5">
              <span
                v-if="c.applied_at"
                class="flex items-center gap-1 rounded-full bg-blue-500/10 px-2 py-1 text-[10px] font-medium text-blue-400"
              >
                <i class="pi pi-check" /> Applied
              </span>
              <button
                v-else
                class="rounded-lg bg-blue-500/10 px-3 py-1.5 text-xs font-medium text-blue-400 transition hover:bg-blue-500/20"
                :disabled="applyingId === c.id"
                @click="doApply(c.id)"
              >
                <i v-if="applyingId === c.id" class="pi pi-spin pi-spinner mr-1" />
                Apply
              </button>
              <button
                title="View"
                class="rounded-lg bg-white/5 p-2 text-xs text-white/30 transition hover:bg-white/10 hover:text-white"
                @click="selected = c; showDetail = true"
              >
                <i class="pi pi-eye" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- === TAB: REJECTED === -->
    <div v-show="activeTab === 'rejected'" class="mt-6 space-y-4">
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 3" :key="i" class="h-16 animate-pulse rounded-xl bg-white/[0.06]" />
      </div>

      <AdminEmptyState
        v-else-if="rejectedItems.length === 0"
        title="No rejected contributions"
        description="Rejected contributions will appear here."
        icon="pi pi-times"
      />

      <div v-else class="space-y-2">
        <div
          v-for="c in rejectedItems"
          :key="c.id"
          class="rounded-xl border border-white/[0.05] bg-white/[0.02] p-4 transition hover:bg-white/[0.04]"
        >
          <div class="flex items-start gap-4">
            <div class="min-w-0 flex-1">
              <div class="flex flex-wrap items-center gap-2">
                <span
                  class="rounded-full bg-red-500/10 px-2 py-0.5 text-[10px] font-bold text-red-400 uppercase"
                  >Rejected</span
                >
                <span
                  class="rounded-full px-2 py-0.5 text-[10px] font-bold uppercase"
                  :class="contributionTypeClass(c.contribution_type)"
                >
                  {{ c.contribution_type }}
                </span>
                <span
                  class="rounded-full bg-white/[0.06] px-2 py-0.5 text-[10px] font-medium text-white/40"
                  >{{ c.target_type }}</span
                >
              </div>
              <p class="mt-1 text-xs text-white/30">
                by {{ c.contributor_username || shortId(c.contributor_id) }} · {{ timeAgo(c.created_at) }}
                <span v-if="c.moderator_note"> · Note: {{ c.moderator_note }}</span>
              </p>
            </div>
            <button
              title="View"
              class="rounded-lg bg-white/5 p-2 text-xs text-white/30 transition hover:bg-white/10 hover:text-white"
              @click="selected = c; showDetail = true"
            >
              <i class="pi pi-eye" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Detail dialog -->
    <Dialog
      v-model:visible="showDetail"
      modal
      :draggable="false"
      :style="{ width: '520px' }"
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
          <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-amber-500/10">
            <i class="pi pi-pen-to-square text-amber-400" />
          </div>
          <div>
            <h3 class="text-base font-semibold text-white">Contribution Details</h3>
          </div>
        </div>
      </template>

      <div class="mt-2 space-y-3 text-sm">
        <div class="flex gap-2">
          <span class="w-28 shrink-0 text-white/40">Type</span>
          <span
            class="rounded-full px-2 py-0.5 text-xs font-bold uppercase"
            :class="contributionTypeClass(selected.contribution_type)"
            >{{ selected.contribution_type }}</span
          >
        </div>
        <div class="flex gap-2">
          <span class="w-28 shrink-0 text-white/40">Target</span>
          <router-link
            :to="{ name: selected.target_type, params: { id: selected.target_id } }"
            class="text-blue-400/60 hover:text-blue-400"
            >{{ selected.target_type }} #{{ shortId(selected.target_id) }}</router-link
          >
        </div>
        <div class="flex gap-2">
          <span class="w-28 shrink-0 text-white/40">Contributor</span>
          <span class="font-mono text-xs text-white/80">{{ selected.contributor_id }}</span>
        </div>
        <div class="flex gap-2">
          <span class="w-28 shrink-0 text-white/40">Status</span>
          <span
            class="rounded-full px-2 py-0.5 text-xs font-bold uppercase"
            :class="
              selected.status === 'approved'
                ? 'bg-green-500/10 text-green-400'
                : selected.status === 'rejected'
                  ? 'bg-red-500/10 text-red-400'
                  : 'bg-yellow-500/10 text-yellow-400'
            "
          >
            {{ selected.status }}
          </span>
        </div>
        <div class="flex gap-2">
          <span class="w-28 shrink-0 text-white/40">Created</span>
          <span class="text-white/60">{{ formatDate(selected.created_at) }}</span>
        </div>
        <div v-if="selected.applied_at" class="flex gap-2">
          <span class="w-28 shrink-0 text-white/40">Applied</span>
          <span class="text-green-400/60">{{ formatDate(selected.applied_at) }}</span>
        </div>
        <div v-if="selected.ai_verdict" class="flex gap-2">
          <span class="w-28 shrink-0 text-white/40">AI Verdict</span>
          <span :class="aiVerdictClass(selected.ai_verdict)" class="font-medium">
            {{ selected.ai_verdict }}
            <span v-if="selected.ai_confidence"
              >({{ (selected.ai_confidence * 100).toFixed(0) }}%)</span
            >
          </span>
        </div>
        <div v-if="selected.ai_confidence && selected.ai_confidence > 0" class="flex gap-2">
          <span class="w-28 shrink-0 text-white/40">AI Confidence</span>
          <span class="text-xs text-white/50"
            >{{ (selected.ai_confidence * 100).toFixed(0) }}%</span
          >
        </div>
        <div>
          <span class="mb-1 block text-white/40">Data</span>
          <pre
            class="max-h-40 overflow-x-auto rounded-xl bg-black/40 p-3 font-mono text-xs text-white/60"
            >{{ formatJSON(selected.data) }}</pre
          >
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            label="Close"
            text
            class="!text-slate-400 hover:!text-white"
            @click="showDetail = false"
          />
          <Button
            v-if="selected?.status === 'approved' && !selected?.applied_at"
            label="Apply Now"
            icon="pi pi-check"
            :loading="applyingId === selected?.id"
            class="!rounded-xl !bg-blue-500/20 !text-blue-400 !ring-1 !ring-blue-500/20 hover:!bg-blue-500/30"
            @click="doApply(selected.id)"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useContributionApi } from '@/services/api/contribution'
import { useToast } from 'primevue/usetoast'
import type { Contribution } from '@/services/api/contribution'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import { AdminSectionHeader, AdminEmptyState } from '@/components/admin'

const api = useContributionApi()
const toast = useToast()

const activeTab = ref('pending')
const loading = ref(true)
const reviewingId = ref<string | null>(null)
const reviewNote = ref('')
const applyingId = ref<string | null>(null)
const selected = ref<Contribution | null>(null)
const showDetail = ref(false)

// Filters
const searchQuery = ref('')
const typeFilter = ref('')
let searchTimer: ReturnType<typeof setTimeout> | null = null

function onSearchInput() {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => fetchCurrentTab(), 300)
}

// Pending
const pendingItems = ref<Contribution[]>([])
const pendingPage = ref(1)
const pendingTotal = ref(0)

// Approved (fetched from API with status filter - using listMy as placeholder)
const approvedItems = ref<Contribution[]>([])

// Rejected
const rejectedItems = ref<Contribution[]>([])

const tabs = computed(() => [
  { key: 'pending', label: 'Pending', icon: 'pi pi-inbox', badge: pendingCount.value || undefined },
  { key: 'approved', label: 'Approved', icon: 'pi pi-check' },
  { key: 'rejected', label: 'Rejected', icon: 'pi pi-times' },
])

const pendingCount = computed(() => pendingTotal.value)

function shortId(id: string): string {
  return id.length > 8 ? id.slice(0, 8) + '\u2026' : id
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

function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return dateStr
  }
}

function formatJSON(data: any): string {
  try {
    const obj = typeof data === 'string' ? JSON.parse(data) : data
    return JSON.stringify(obj, null, 2)
  } catch {
    return String(data)
  }
}

function contributionTypeClass(type: string): string {
  const map: Record<string, string> = {
    lyrics: 'bg-purple-500/10 text-purple-400',
    translation: 'bg-cyan-500/10 text-cyan-400',
    credits: 'bg-amber-500/10 text-amber-400',
    metadata: 'bg-pink-500/10 text-pink-400',
    album_art: 'bg-orange-500/10 text-orange-400',
    bio: 'bg-green-500/10 text-green-400',
  }
  return map[type] || 'bg-white/10 text-white/50'
}

function aiVerdictClass(verdict: string): string {
  const v = verdict?.toLowerCase()
  if (v === 'approved' || v === 'approve') return 'text-green-400'
  if (v === 'rejected' || v === 'reject') return 'text-red-400'
  return 'text-yellow-400'
}

async function fetchCurrentTab() {
  if (activeTab.value === 'pending') {
    pendingPage.value = 1
    await fetchPending()
  } else if (activeTab.value === 'approved') {
    await fetchApproved()
  } else {
    await fetchRejected()
  }
}

async function fetchPending() {
  try {
    const res = await api.listPending({ page: pendingPage.value, page_size: 20 })
    if (pendingPage.value === 1) {
      pendingItems.value = res?.data || []
    } else {
      pendingItems.value = [...pendingItems.value, ...(res?.data || [])]
    }
    pendingTotal.value = res?.meta.total || 0
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Failed to load pending contributions', life: 3000 })
  }
}

async function fetchApproved() {
  try {
    const params: { page: number; page_size: number; status: string; type?: string; q?: string } = { page: 1, page_size: 50, status: 'approved' }
    if (typeFilter.value) params.type = typeFilter.value
    if (searchQuery.value) params.q = searchQuery.value
    const res = await api.adminListContributions(params)
    approvedItems.value = res?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Failed to load approved contributions', life: 3000 })
  }
}

async function fetchRejected() {
  try {
    const params: { page: number; page_size: number; status: string; type?: string; q?: string } = { page: 1, page_size: 50, status: 'rejected' }
    if (typeFilter.value) params.type = typeFilter.value
    if (searchQuery.value) params.q = searchQuery.value
    const res = await api.adminListContributions(params)
    rejectedItems.value = res?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Failed to load rejected contributions', life: 3000 })
  }
}

async function doReview(id: string, action: 'approve' | 'reject') {
  try {
    await api.review(id, { action, note: reviewNote.value || undefined })
    pendingItems.value = pendingItems.value.filter((c) => c.id !== id)
    pendingTotal.value--
    reviewingId.value = null
    reviewNote.value = ''
    toast.add({ severity: 'success', summary: `Contribution ${action}d`, life: 2000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Failed to review contribution', life: 3000 })
  }
}

async function doApply(id: string) {
  applyingId.value = id
  try {
    await api.apply(id)
    toast.add({ severity: 'success', summary: 'Contribution applied to target', life: 2000 })
    // Refresh both tabs
    await Promise.all([fetchPending(), fetchApproved()])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Failed to apply contribution', life: 3000 })
  } finally {
    applyingId.value = null
  }
}

async function refreshAll() {
  loading.value = true
  await Promise.all([fetchPending(), fetchApproved(), fetchRejected()])
  loading.value = false
}

watch(activeTab, () => fetchCurrentTab())
watch(pendingPage, () => fetchPending())

onMounted(async () => {
  loading.value = true
  await Promise.all([fetchPending(), fetchApproved(), fetchRejected()])
  loading.value = false
})
</script>
