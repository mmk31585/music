<template>
  <div class="rounded-2xl bg-white/6 p-4 ring-1 ring-white/10 transition hover:bg-white/8">
    <div class="flex items-center justify-between">
      <i aria-hidden="true" :class="`pi pi-${icon}`" :style="iconStyle" class="text-lg" />
      <div v-if="trend !== undefined" class="flex items-center gap-0.5 text-[10px]" :class="trend >= 0 ? 'text-spotify' : 'text-red-400'">
        <i aria-hidden="true" :class="trend >= 0 ? 'pi pi-arrow-up' : 'pi pi-arrow-down'" class="text-[9px]" />
        <span>{{ Math.abs(trend) }}%</span>
      </div>
    </div>
    <p class="mt-3 text-2xl font-black text-white tabular-nums">{{ value }}</p>
    <p class="mt-0.5 text-xs text-white/40">{{ label }}</p>
    <p v-if="sub" class="text-[10px] text-white/25">{{ sub }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  icon: string
  value: string
  label: string
  sub?: string
  trend?: number
  color?: string
}>()

const iconStyle = computed(() => ({
  color: props.color || 'rgba(255,255,255,0.4)',
}))
</script>
