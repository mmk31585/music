<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="visible"
        role="dialog"
        aria-modal="true"
        aria-label="Keyboard Shortcuts"
        class="fixed inset-0 z-[300] flex items-center justify-center bg-bg-overlay/60 backdrop-blur-xs"
        @click.self="visible = false"
        @keydown.escape="visible = false"
      >
        <div class="mx-4 w-full max-w-md rounded-2xl border border-border-default bg-surface-raised p-6 shadow-2xl">
          <div class="mb-6 flex items-center justify-between">
            <h2 class="text-lg font-bold text-primary">Keyboard Shortcuts</h2>
            <button
              ref="closeBtnRef"
              type="button"
              aria-label="Close shortcuts"
              class="flex h-8 w-8 items-center justify-center rounded-full text-secondary transition hover:bg-surface-active hover:text-primary"
              @click="visible = false"
            >
              <X aria-hidden="true" class="text-sm"  />
            </button>
          </div>

          <div class="space-y-1">
            <div
              v-for="shortcut in shortcuts"
              :key="shortcut.label"
              class="flex items-center justify-between rounded-lg px-3 py-2.5 transition hover:bg-surface-overlay/60"
            >
              <span class="text-sm text-secondary">{{ shortcut.label }}</span>
              <kbd class="flex items-center gap-1">
                <span
                  v-for="(key, i) in shortcut.keys"
                  :key="i"
                  class="inline-flex items-center rounded-md border border-border-default bg-surface-overlay px-2 py-0.5 text-xs font-medium text-secondary"
                >
                  {{ key }}
                </span>
              </kbd>
            </div>
          </div>

          <p class="mt-4 text-center text-xs text-muted">Press <kbd class="rounded-sm bg-surface-overlay px-1.5 py-0.5 text-xs text-secondary">?</kbd> to toggle this panel</p>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { X } from 'lucide-vue-next'
import { watch } from 'vue'
import { useTemplateRef } from 'vue'

const visible = defineModel<boolean>('visible', { default: false })
const closeBtnRef = useTemplateRef<HTMLButtonElement>('closeBtnRef')

// Focus trap: focus close button when dialog opens
watch(visible, (val) => {
  if (val) {
    requestAnimationFrame(() => closeBtnRef.value?.focus())
  }
})

const shortcuts = [
  { label: 'Play / Pause', keys: ['Space'] },
  { label: 'Next track', keys: ['N'] },
  { label: 'Previous track', keys: ['P'] },
  { label: 'Seek forward', keys: ['→'] },
  { label: 'Seek backward', keys: ['←'] },
  { label: 'Volume up', keys: ['↑'] },
  { label: 'Volume down', keys: ['↓'] },
  { label: 'Toggle mute', keys: ['M'] },
  { label: 'Cycle shuffle modes', keys: ['S'] },
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
