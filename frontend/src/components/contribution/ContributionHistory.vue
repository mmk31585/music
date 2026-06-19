<template>
  <div class="glass-strong rounded-2xl p-6 md:p-8">
    <h3 class="mb-6 text-lg font-bold text-white">Version History</h3>

    <div v-if="loading" class="space-y-3">
      <div v-for="i in 4" :key="i" class="flex gap-4">
        <div class="h-8 w-8 shrink-0 animate-pulse rounded-full bg-white/5" />
        <div class="flex-1 space-y-2">
          <div class="h-4 w-3/4 animate-pulse rounded bg-white/5" />
          <div class="h-3 w-1/2 animate-pulse rounded bg-white/5" />
        </div>
      </div>
    </div>

    <div v-else-if="items.length === 0" class="flex flex-col items-center gap-3 py-8 text-center">
      <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-white/5">
        <i aria-hidden="true" class="pi pi-history text-xl text-white/15" />
      </div>
      <p class="text-sm text-white/25">No version history yet</p>
    </div>

    <div v-else class="relative space-y-0">
      <div class="absolute top-3 bottom-3 left-4 w-px bg-white/[0.06]" />

      <div v-for="(item, idx) in items" :key="item.id" class="relative flex gap-4 pb-6 last:pb-0">
        <div
          class="relative z-10 mt-1.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full"
          :class="badgeClass(item.change_type)"
        >
          <i aria-hidden="true" :class="iconClass(item.change_type)" class="text-xs" />
        </div>

        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-white/70">{{
              changeLabel(item.change_type)
            }}</span>
            <span class="rounded bg-white/5 px-1.5 py-0.5 text-[10px] font-medium text-white/30">{{
              item.changed_by
            }}</span>
          </div>
          <p class="mt-1 text-xs text-white/30">{{ formatDate(item.created_at) }}</p>

          <div v-if="item.data" class="mt-2">
            <button
              type="button"
              class="text-xs text-white/30 transition-colors hover:text-white/50"
              @click="toggleExpand(idx)"
            >
              <i
                class="pi pi-chevron-right mr-1 text-[10px]"
                :class="{ 'rotate-90': expanded === idx }"
              />
              View data
            </button>
            <pre
              v-if="expanded === idx"
              class="mt-2 overflow-x-auto rounded-lg bg-black/40 p-3 font-mono text-[11px] text-white/40"
              >{{ formatData(item.data) }}</pre
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { ContributionHistoryItem } from '@/services/api/contribution'

defineProps<{
  items: ContributionHistoryItem[]
  loading?: boolean
}>()

const expanded = ref<number | null>(null)

function toggleExpand(idx: number) {
  expanded.value = expanded.value === idx ? null : idx
}

function badgeClass(changeType: string) {
  switch (changeType) {
    case 'approve':
      return 'bg-[#1db954]/15 text-[#1db954]'
    case 'reject':
      return 'bg-red-500/15 text-red-400'
    case 'create':
      return 'bg-[#60a5fa]/15 text-[#60a5fa]'
    case 'revert':
      return 'bg-[#f59e0b]/15 text-[#f59e0b]'
    default:
      return 'bg-white/5 text-white/40'
  }
}

function iconClass(changeType: string) {
  switch (changeType) {
    case 'approve':
      return 'pi pi-check'
    case 'reject':
      return 'pi pi-times'
    case 'create':
      return 'pi pi-plus'
    case 'revert':
      return 'pi pi-undo'
    default:
      return 'pi pi-circle'
  }
}

function changeLabel(changeType: string) {
  switch (changeType) {
    case 'approve':
      return 'Approved'
    case 'reject':
      return 'Rejected'
    case 'create':
      return 'Created'
    case 'revert':
      return 'Reverted'
    default:
      return changeType
  }
}

function formatDate(dateStr: string) {
  const d = new Date(dateStr)
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatData(data: unknown): string {
  if (typeof data === 'string') {
    try {
      return JSON.stringify(JSON.parse(data), null, 2)
    } catch {
      return data
    }
  }
  return JSON.stringify(data, null, 2)
}
</script>
