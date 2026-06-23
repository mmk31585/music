<template>
  <span>{{ formattedValue }}</span>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'

const props = withDefaults(
  defineProps<{
    to: number
    duration?: number
  }>(),
  { duration: 1000 },
)

const current = ref(0)
let animationId: number | null = null

const formattedValue = computed(() =>
  current.value.toLocaleString('fa-IR'),
)

function prefersReducedMotion(): boolean {
  if (typeof window === 'undefined') return false
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function animate() {
  if (prefersReducedMotion()) {
    current.value = props.to
    return
  }

  const startTime = performance.now()
  const startValue = 0
  const delta = props.to - startValue

  function step(now: number) {
    const elapsed = now - startTime
    const progress = Math.min(elapsed / props.duration, 1)
    const eased = 1 - (1 - progress) * (1 - progress)
    current.value = Math.round(startValue + delta * eased)

    if (progress < 1) {
      animationId = requestAnimationFrame(step)
    }
  }

  animationId = requestAnimationFrame(step)
}

watch(
  () => props.to,
  () => {
    if (animationId) cancelAnimationFrame(animationId)
    animate()
  },
)

onMounted(() => animate())
</script>
