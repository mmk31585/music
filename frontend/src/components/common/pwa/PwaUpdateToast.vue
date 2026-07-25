<template>
  <div
    v-if="needRefresh && !dismissed"
    role="alert"
    class="fixed bottom-20 inset-x-0 z-50 flex items-center justify-center pointer-events-none"
  >
    <div class="pointer-events-auto flex items-center gap-3 rounded-2xl bg-surface-raised/90 backdrop-blur-xl px-5 py-3 shadow-2xl ring-1 ring-white/10 motion-safe:animate-fade-in-up">
      <div class="flex h-8 w-8 items-center justify-center rounded-full bg-aurora-blue/20">
        <RefreshCw aria-hidden="true" class="text-xs text-aurora-blue"  />
      </div>
      <p class="text-sm text-white">New version available</p>
      <button
        class="rounded-full bg-spotify px-4 py-1.5 text-xs font-bold text-black transition-all duration-200 hover:bg-spotify-hover active:scale-95 focus-visible:outline-2 focus-visible:outline-white"
        @click="refresh"
      >
        RefreshCw
      </button>
      <button
        aria-label="Dismiss"
        class="flex h-8 w-8 items-center justify-center rounded-full text-white/30 transition hover:bg-white/10 hover:text-white focus-visible:outline-2 focus-visible:outline-accent"
        @click="dismissed = true"
      >
        <X aria-hidden="true" class="text-xs"  />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RefreshCw, X } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRegisterSW } from 'virtual:pwa-register/vue'

const { needRefresh, updateServiceWorker } = useRegisterSW({
  onRegisteredSW() {},
  onRegisterError() {},
})

const dismissed = ref(false)

function refresh() {
  dismissed.value = true
  updateServiceWorker()
}
</script>

<style scoped>
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}
.motion-safe\:animate-fade-in-up {
  animation: fade-in-up 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
@media (prefers-reduced-motion: reduce) {
  .motion-safe\:animate-fade-in-up { animation: none; }
}
</style>
