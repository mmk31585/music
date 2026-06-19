<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="visible"
        class="fixed inset-0 z-[300] flex items-center justify-center bg-black/60 backdrop-blur-sm"
        @click.self="visible = false"
      >
        <div class="mx-4 w-full max-w-md rounded-2xl border border-white/[0.08] bg-[#141414] p-6 shadow-2xl">
          <div class="mb-6 flex items-center justify-between">
            <h2 class="text-lg font-bold text-white">Keyboard Shortcuts</h2>
            <button
              type="button"
              aria-label="Close shortcuts"
              class="flex h-8 w-8 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
              @click="visible = false"
            >
              <i aria-hidden="true" class="pi pi-times text-sm" />
            </button>
          </div>

          <div class="space-y-1">
            <div
              v-for="shortcut in shortcuts"
              :key="shortcut.label"
              class="flex items-center justify-between rounded-lg px-3 py-2.5 transition hover:bg-white/[0.04]"
            >
              <span class="text-sm text-slate-300">{{ shortcut.label }}</span>
              <kbd class="flex items-center gap-1">
                <span
                  v-for="(key, i) in shortcut.keys"
                  :key="i"
                  class="inline-flex items-center rounded-md border border-white/[0.10] bg-white/[0.06] px-2 py-0.5 text-xs font-medium text-slate-300"
                >
                  {{ key }}
                </span>
              </kbd>
            </div>
          </div>

          <p class="mt-4 text-center text-xs text-slate-600">Press <kbd class="rounded bg-white/[0.06] px-1.5 py-0.5 text-xs text-slate-400">?</kbd> to toggle this panel</p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
const visible = defineModel<boolean>('visible', { default: false })

const shortcuts = [
  { label: 'Play / Pause', keys: ['Space'] },
  { label: 'Next track', keys: ['N'] },
  { label: 'Previous track', keys: ['P'] },
  { label: 'Seek forward', keys: ['→'] },
  { label: 'Seek backward', keys: ['←'] },
  { label: 'Volume up', keys: ['↑'] },
  { label: 'Volume down', keys: ['↓'] },
  { label: 'Toggle mute', keys: ['M'] },
  { label: 'Toggle shuffle', keys: ['S'] },
  { label: 'Toggle repeat', keys: ['R'] },
  { label: 'Toggle sidebar', keys: ['Ctrl', 'B'] },
  { label: 'Search', keys: ['Ctrl', 'K'] },
  { label: 'Fullscreen player', keys: ['F'] },
  { label: 'Close / Escape', keys: ['Esc'] },
  { label: 'Like track', keys: ['L'] },
  { label: 'Queue', keys: ['Q'] },
  { label: 'Lyrics', keys: ['Ctrl', 'L'] },
]
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
