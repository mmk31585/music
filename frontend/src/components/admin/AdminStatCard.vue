<template>
  <div
    class="group relative overflow-hidden rounded-2xl border border-white/6 bg-white/3 p-5 transition-all duration-300 hover:border-white/10 hover:bg-white/5"
  >
    <!-- Glow effect on hover -->
    <div
      class="pointer-events-none absolute -right-6 -top-6 h-24 w-24 rounded-full opacity-0 blur-2xl transition-opacity duration-500 group-hover:opacity-100"
      :class="glowColor"
    />

    <div class="relative flex items-start justify-between">
      <div class="min-w-0 flex-1">
        <p class="text-xs font-medium uppercase tracking-wider text-slate-500">{{ label }}</p>

        <p class="mt-2 text-3xl font-bold tabular-nums text-white">
          <span v-if="loading" class="inline-block h-8 w-16 animate-pulse rounded-lg bg-white/10" />
          <template v-else>{{ formattedValue }}</template>
        </p>

        <p v-if="hint" class="mt-2 text-xs text-slate-500">{{ hint }}</p>
      </div>

      <div
        class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl transition-colors"
        :class="iconBgClass"
      >
        <i aria-hidden="true" :class="[icon, 'text-sm', iconColorClass]" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    label: string
    value: number
    hint?: string
    icon?: string
    color?: 'emerald' | 'blue' | 'purple' | 'amber'
    loading?: boolean
  }>(),
  {
    icon: 'pi pi-chart-bar',
    color: 'emerald',
    loading: false,
  },
)

const formattedValue = computed(() =>
  props.value >= 1000 ? `${(props.value / 1000).toFixed(1)}k` : String(props.value),
)

const colorMap = {
  emerald: {
    glow: 'bg-emerald-500/20',
    iconBg: 'bg-emerald-500/10 group-hover:bg-emerald-500/15',
    iconColor: 'text-emerald-400',
  },
  blue: {
    glow: 'bg-blue-500/20',
    iconBg: 'bg-blue-500/10 group-hover:bg-blue-500/15',
    iconColor: 'text-blue-400',
  },
  purple: {
    glow: 'bg-purple-500/20',
    iconBg: 'bg-purple-500/10 group-hover:bg-purple-500/15',
    iconColor: 'text-purple-400',
  },
  amber: {
    glow: 'bg-amber-500/20',
    iconBg: 'bg-amber-500/10 group-hover:bg-amber-500/15',
    iconColor: 'text-amber-400',
  },
}

const glowColor = computed(() => colorMap[props.color].glow)
const iconBgClass = computed(() => colorMap[props.color].iconBg)
const iconColorClass = computed(() => colorMap[props.color].iconColor)
</script>
