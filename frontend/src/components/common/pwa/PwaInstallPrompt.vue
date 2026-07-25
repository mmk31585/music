<template>
  <!-- Android/Chrome install banner -->
  <div
    v-if="showInstall"
    role="alert"
    class="fixed bottom-0 inset-x-0 z-50 flex items-center gap-3 bg-surface-raised/90 backdrop-blur-xl border-t border-white/10 px-4 py-3 motion-safe:animate-slide-up"
    style="padding-bottom: max(0.75rem, env(safe-area-inset-bottom, 0.75rem))"
  >
    <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-spotify/20">
      <Download aria-hidden="true" class="text-sm text-spotify"  />
    </div>
    <div class="min-w-0 flex-1">
      <p class="text-sm font-semibold text-white">Install Muse</p>
      <p class="text-xs text-white/40">Get the best experience with our app</p>
    </div>
    <button
      class="rounded-full bg-spotify px-4 py-2 text-xs font-bold text-black transition-all duration-200 hover:bg-spotify-hover active:scale-95 focus-visible:outline-2 focus-visible:outline-white"
      @click="install"
    >
      Install
    </button>
    <button
      aria-label="Dismiss"
      class="flex h-8 w-8 items-center justify-center rounded-full text-white/30 transition hover:bg-white/10 hover:text-white focus-visible:outline-2 focus-visible:outline-accent"
      @click="dismiss"
    >
      <X aria-hidden="true" class="text-xs"  />
    </button>
  </div>

  <!-- iOS instructions -->
  <div
    v-if="showIOS"
    role="alert"
    class="fixed bottom-0 inset-x-0 z-50 flex items-start gap-3 bg-surface-raised/90 backdrop-blur-xl border-t border-white/10 px-4 py-4 motion-safe:animate-slide-up"
    style="padding-bottom: max(0.75rem, env(safe-area-inset-bottom, 0.75rem))"
  >
    <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-aurora-blue/20">
      <Apple aria-hidden="true" class="text-sm text-aurora-blue"  />
    </div>
    <div class="min-w-0 flex-1">
      <p class="text-sm font-semibold text-white">Add Muse to Home Screen</p>
      <p class="text-xs text-white/40 leading-relaxed">
        Tap <span class="inline-flex items-center gap-0.5 rounded bg-white/10 px-1.5 py-0.5 text-[10px] font-medium text-white/70"><Share2 aria-hidden="true" class="text-[8px]"  /> Share</span>
        then <span class="inline-flex items-center gap-0.5 rounded bg-white/10 px-1.5 py-0.5 text-[10px] font-medium text-white/70">Add to Home Screen</span>
      </p>
    </div>
    <button
      aria-label="Dismiss"
      class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-white/30 transition hover:bg-white/10 hover:text-white focus-visible:outline-2 focus-visible:outline-accent"
      @click="dismissIOS"
    >
      <X aria-hidden="true" class="text-xs"  />
    </button>
  </div>
</template>

<script setup lang="ts">
import { Apple, Download, Share2, X } from 'lucide-vue-next'
import { ref, computed, onMounted } from 'vue'

const isInstallable = ref(false)
const showIOSInstructions = ref(false)
const installDismissed = ref(false)
const iosDismissed = ref(false)

interface BeforeInstallPromptEvent extends Event {
  prompt: () => Promise<void>
  userChoice: Promise<{ outcome: 'accepted' | 'dismissed' }>
}

let deferredPrompt: BeforeInstallPromptEvent | null = null

// Capture the install prompt immediately at module load (before mount)
// so we don't miss the one-time event if it fires during page load.
function handleBeforeInstallPrompt(e: Event) {
  e.preventDefault()
  deferredPrompt = e as BeforeInstallPromptEvent
  isInstallable.value = true
}

if (typeof window !== 'undefined') {
  window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
  window.addEventListener('appinstalled', () => {
    isInstallable.value = false
    deferredPrompt = null
  })
}

onMounted(() => {
  const dismissedTime = localStorage.getItem('pwa-install-dismissed')
  if (dismissedTime) {
    const elapsed = Date.now() - Number(dismissedTime)
    if (elapsed < 7 * 24 * 60 * 60 * 1000) {
      installDismissed.value = true
    }
  }
  const iosTime = localStorage.getItem('pwa-ios-dismissed')
  if (iosTime) {
    const elapsed = Date.now() - Number(iosTime)
    if (elapsed < 7 * 24 * 60 * 60 * 1000) {
      iosDismissed.value = true
    }
  }

  const isIOS = /iPad|iPhone|iPod/.test(navigator.userAgent) && !(window as any).MSStream
  const isStandalone = window.matchMedia('(display-mode: standalone)').matches || (navigator as any).standalone
  if (isIOS && !isStandalone) {
    showIOSInstructions.value = true
  }
})

const showInstall = computed(() => isInstallable.value && !installDismissed.value)
const showIOS = computed(() => showIOSInstructions.value && !iosDismissed.value)

async function install() {
  if (!deferredPrompt) return
  deferredPrompt.prompt()
  const { outcome } = await deferredPrompt.userChoice
  if (outcome === 'accepted') {
    isInstallable.value = false
  }
  deferredPrompt = null
}

function dismiss() {
  installDismissed.value = true
  localStorage.setItem('pwa-install-dismissed', String(Date.now()))
}

function dismissIOS() {
  iosDismissed.value = true
  localStorage.setItem('pwa-ios-dismissed', String(Date.now()))
}
</script>

<style scoped>
@keyframes slide-up {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}
.motion-safe\:animate-slide-up {
  animation: slide-up 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}
@media (prefers-reduced-motion: reduce) {
  .motion-safe\:animate-slide-up { animation: none; }
}
</style>
