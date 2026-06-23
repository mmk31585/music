<template>
  <div
    class="group glass hover:glass-hover rounded-2xl p-5 transition-all duration-300 hover:-translate-y-0.5"
  >
    <div class="flex items-center justify-between">
      <div>
        <p class="text-2xl font-black text-white">{{ prefix }}{{ formattedValue }}</p>
        <p class="mt-1 text-xs font-medium text-white/40">{{ label }}</p>
      </div>
      <div
        class="flex h-10 w-10 items-center justify-center rounded-xl bg-white/5 text-white/40 transition group-hover:bg-spotify/10 group-hover:text-spotify"
      >
        <i aria-hidden="true" :class="icon" />
      </div>
    </div>
    <div
      v-if="change !== undefined"
      class="mt-3 flex items-center gap-1 text-xs"
      :class="change >= 0 ? 'text-green-400' : 'text-red-400'"
    >
      <i aria-hidden="true" :class="change >= 0 ? 'pi pi-arrow-up' : 'pi pi-arrow-down'" class="text-[10px]" />
      <span>{{ Math.abs(change) }}% vs last period</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  label: string
  value: number
  icon: string
  prefix?: string
  change?: number
}>()

const formattedValue = computed(() => {
  if (props.value >= 1_000_000) return `${(props.value / 1_000_000).toFixed(1)}M`
  if (props.value >= 1_000) return `${(props.value / 1_000).toFixed(1)}K`
  return String(props.value)
})
</script>
