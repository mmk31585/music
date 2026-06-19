<template>
  <div v-if="hasError" class="flex flex-col items-center gap-4 rounded-2xl border border-white/10 bg-white/[0.03] px-6 py-12 text-center backdrop-blur">
    <div class="flex h-14 w-14 items-center justify-center rounded-full bg-red-500/10">
      <i aria-hidden="true" class="pi pi-exclamation-circle text-2xl text-red-400" />
    </div>
    <div>
      <h3 class="text-lg font-bold text-white">Something went wrong</h3>
      <p class="mt-1 text-sm text-slate-400" role="alert">{{ errorMessage }}</p>
    </div>
    <button
      type="button"
      class="rounded-full bg-white/10 px-5 py-2 text-sm font-semibold text-white transition hover:bg-white/15"
      @click="retry"
    >
      Try again
    </button>
  </div>
  <slot v-else />
</template>

<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue'

interface Props {
  message?: string
}

const props = withDefaults(defineProps<Props>(), {
  message: 'An unexpected error occurred. Please try again.',
})

const emit = defineEmits<{
  error: [err: Error]
  retry: []
}>()

const hasError = ref(false)
const errorMessage = ref(props.message)

onErrorCaptured((err: Error) => {
  hasError.value = true
  errorMessage.value = err.message || props.message
  emit('error', err)
  return false
})

function retry() {
  hasError.value = false
  errorMessage.value = props.message
  emit('retry')
}
</script>
