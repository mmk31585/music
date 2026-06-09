<template>
  <button
    type="button"
    class="flex items-center gap-2 rounded-full bg-white/10 px-4 py-2 text-sm text-white transition-all hover:bg-white/15 active:scale-95 disabled:cursor-not-allowed disabled:opacity-50"
    :disabled="disabled"
    @click="onClick"
  >
    <i class="pi pi-external-link text-xs" />
    <span>{{ label }}</span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    disabled?: boolean
    supported?: boolean
    active?: boolean
  }>(),
  { disabled: false, supported: true, active: false },
)

const emit = defineEmits<{ click: [] }>()

const label = computed(() =>
  props.active ? 'Return Player' : props.supported ? 'Pop Out Player' : 'Pop Out Unavailable',
)

function onClick() {
  if (!props.disabled) emit('click')
}
</script>
