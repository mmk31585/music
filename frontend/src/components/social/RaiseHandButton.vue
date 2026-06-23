<template>
  <button
    v-if="myRole === 'listener'"
    class="raise-hand-btn fixed bottom-24 left-1/2 z-50 flex -translate-x-1/2 items-center gap-2.5 rounded-full px-5 py-3 text-sm font-semibold shadow-lg backdrop-blur-xl transition-all duration-250"
    :class="handRaised
      ? 'border border-amber-500/50 bg-amber-500/10 text-amber-500 hover:bg-amber-500/20'
      : 'border border-white/10 bg-white/5 text-white hover:bg-white/10'"
    @click="handRaised ? $emit('lower') : $emit('raise')"
  >
    <!-- Spinner when hand is raised -->
    <svg
      v-if="handRaised"
      class="h-4 w-4 animate-spin"
      viewBox="0 0 24 24"
      fill="none"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
    </svg>
    <i aria-hidden="true" v-else class="pi pi-hand text-base" />
    <span>{{ handRaised ? 'در انتظار تأیید میزبان...' : 'دست بلند کن' }}</span>
  </button>
</template>

<script setup lang="ts">
defineProps<{
  handRaised: boolean
  myRole: 'host' | 'speaker' | 'listener'
}>()

defineEmits<{
  raise: []
  lower: []
}>()
</script>

<style scoped>
.duration-250 {
  transition-duration: 250ms;
}
</style>
