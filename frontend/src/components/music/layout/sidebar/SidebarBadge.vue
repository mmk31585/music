<template>
  <span
    v-if="value !== undefined && value !== null"
    class="inline-flex items-center justify-center rounded-full px-1.5 py-0.5 text-[10px] font-bold leading-none"
    :class="badgeClass"
  >
    {{ displayValue }}
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  value?: string | number | null
  variant?: 'default' | 'new' | 'count'
}>()

const badgeClass = computed(() => {
  switch (props.variant) {
    case 'new':
      return 'bg-accent-subtle text-accent'
    case 'count':
      return 'bg-surface-active text-secondary'
    default:
      return 'bg-surface-active text-secondary'
  }
})

const displayValue = computed(() => {
  if (props.variant === 'new') return 'NEW'
  return String(props.value ?? '')
})
</script>
