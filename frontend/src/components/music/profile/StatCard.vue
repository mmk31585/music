<template>
  <div class="relative rounded-2xl bg-surface-overlay/60 ring-1 ring-border-subtle p-4 flex flex-col gap-1 transition hover:bg-surface-hover">
    <div class="flex items-center justify-between">
      <div
        class="w-8 h-8 rounded-lg bg-surface-active/80 flex items-center justify-center text-sm"
        :style="iconStyle"
      >
        <i aria-hidden="true" :class="icon" />
      </div>
      <div
        v-if="trend !== undefined"
        class="absolute top-3 end-3 flex items-center gap-0.5 text-[10px]"
        :class="trend >= 0 ? 'text-success' : 'text-danger'"
      >
        <span>{{ trend >= 0 ? '▲' : '▼' }}</span>
        <span>{{ Math.abs(trend) }}%</span>
      </div>
    </div>
    <p
      class="text-2xl font-black tabular-nums"
      :class="value === '\u2014' ? 'text-muted' : 'text-primary'"
    >
      {{ animatedValue }}
    </p>
    <p v-if="sub" class="text-xs text-tertiary mt-0.5 truncate">{{ sub }}</p>
    <p class="text-[10px] uppercase tracking-widest text-muted mt-auto pt-2">{{ label }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

const props = defineProps<{
  icon: string
  value: string
  label: string
  sub?: string
  trend?: number
  color?: string
  animated?: boolean
}>()

const displayValue = ref(props.value)

function animateCountUp(target: string, duration = 800) {
  const num = parseInt(target.replace(/[^0-9]/g, ''), 10)
  if (isNaN(num)) {
    displayValue.value = target
    return
  }
  const suffix = target.replace(/[0-9.]/g, '')
  const start = 0
  const startTime = performance.now()
  function tick(now: number) {
    const elapsed = now - startTime
    const progress = Math.min(elapsed / duration, 1)
    const eased = 1 - Math.pow(1 - progress, 3)
    const current = Math.round(start + (num - start) * eased)
    displayValue.value = current.toLocaleString() + suffix
    if (progress < 1) requestAnimationFrame(tick)
  }
  requestAnimationFrame(tick)
}

onMounted(() => {
  if (props.animated) animateCountUp(props.value)
  else displayValue.value = props.value
})

const animatedValue = computed(() => displayValue.value)

const iconStyle = computed(() => ({
  color: props.color || 'rgba(255,255,255,0.4)',
}))
</script>
