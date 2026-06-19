<template>
  <Transition name="menu-scale">
    <div
      v-show="true"
      ref="menuRef"
      role="menu"
      aria-label="Player options"
      class="glass-strong absolute bottom-full right-0 mb-2 min-w-[260px] origin-bottom-right rounded-2xl p-2 shadow-2xl"
      @keydown="onKeydown"
      tabindex="-1"
    >
      <!-- Audio Quality -->
      <button
        type="button"
        role="menuitem"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-1 focus-visible:ring-[#1db954]"
        @click="cycleQuality"
      >
        <i aria-hidden="true" class="pi pi-waveform text-base text-white/40" />
        <span class="flex-1 text-left">Audio Quality</span>
        <span class="text-xs font-medium text-white/40 tabular-nums">{{ qualityLabel }}</span>
      </button>

      <!-- Sleep Timer -->
      <button
        type="button"
        role="menuitem"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-1 focus-visible:ring-[#1db954]"
        @click="showSleepPicker = !showSleepPicker"
      >
        <i aria-hidden="true" class="pi pi-clock text-base text-white/40" />
        <span class="flex-1 text-left">Sleep Timer</span>
        <span class="text-xs font-medium" :class="sleepTimerMinutes > 0 ? 'text-[#1db954]' : 'text-white/40'">
          {{ sleepTimerMinutes > 0 ? `${sleepTimerMinutes}m` : 'Off' }}
        </span>
      </button>

      <!-- Sleep sub-picker -->
      <div
        v-if="showSleepPicker"
        class="mb-1 ml-9 flex flex-wrap gap-1.5"
      >
        <button
          v-for="opt in sleepOptions"
          :key="opt.value"
          type="button"
          class="rounded-lg px-2.5 py-1 text-xs font-medium transition"
          :class="sleepTimerMinutes === opt.value ? 'bg-[#1db954] text-black' : 'bg-white/10 text-white/50 hover:bg-white/20 hover:text-white/80'"
          @click.stop="setTimer(opt.value)"
        >
          {{ opt.label }}
        </button>
      </div>

      <!-- Crossfade -->
      <div class="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70">
        <i aria-hidden="true" class="pi pi-arrows-alt text-base text-white/40" />
        <span class="flex-1 text-left">Crossfade</span>
        <span class="text-xs font-medium text-white/40 tabular-nums">{{ crossfadeDuration }}s</span>
        <input
          type="range"
          min="0"
          max="12"
          step="1"
          aria-label="Crossfade duration"
          class="w-16 h-1 cursor-pointer appearance-none rounded-full bg-white/10 accent-[#1db954]"
          :value="crossfadeDuration"
          @input="onCrossfadeChange"
        />
      </div>

      <!-- Open in Mini Player -->
      <button
        type="button"
        role="menuitem"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-1 focus-visible:ring-[#1db954]"
        @click="onTogglePiP"
      >
        <i aria-hidden="true" class="pi pi-window-maximize text-base text-white/40" />
        <span>Open in Mini Player</span>
      </button>

      <div class="my-1 border-t border-white/5" />

      <!-- Save to Library -->
      <button
        type="button"
        role="menuitem"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-1 focus-visible:ring-[#1db954]"
        @click="toggleLike"
      >
        <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill text-[#1db954]' : 'pi pi-heart text-white/40'" class="text-base" />
        <span>{{ liked ? 'Saved to Library' : 'Save to Library' }}</span>
      </button>

      <!-- Add to Playlist -->
      <button
        type="button"
        role="menuitem"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-1 focus-visible:ring-[#1db954]"
        @click="onAddToPlaylist"
      >
        <i aria-hidden="true" class="pi pi-plus-circle text-base text-white/40" />
        <span>Add to Playlist</span>
      </button>

      <!-- Share -->
      <button
        type="button"
        role="menuitem"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-1 focus-visible:ring-[#1db954]"
        @click="onShare"
      >
        <i aria-hidden="true" class="pi pi-share-alt text-base text-white/40" />
        <span>Share Track</span>
      </button>

      <div class="my-1 border-t border-white/5" />

      <!-- Track Info -->
      <button
        type="button"
        role="menuitem"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-white/70 transition-all hover:bg-white/10 hover:text-white focus-visible:ring-1 focus-visible:ring-[#1db954]"
        @click="onTrackInfo"
      >
        <i aria-hidden="true" class="pi pi-info-circle text-base text-white/40" />
        <span>Track Info</span>
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayerControls, usePlayer } from '@/composables/player'
import type { AudioQuality } from '@/services/player/audio-engine'

const emit = defineEmits<{
  close: []
  'toggle-pip': []
}>()

const router = useRouter()
const menuRef = ref<HTMLElement | null>(null)

const {
  currentTrack,
  sleepTimerMinutes,
  crossfadeDuration,
  audioQuality,
  setSleepTimer,
  clearSleepTimer,
} = usePlayerControls()

const showSleepPicker = ref(false)
const liked = ref(false)

const sleepOptions = [
  { label: 'Off', value: 0 },
  { label: '5m', value: 5 },
  { label: '15m', value: 15 },
  { label: '30m', value: 30 },
  { label: '60m', value: 60 },
] as const

function setTimer(minutes: number) {
  if (minutes === 0) clearSleepTimer()
  else setSleepTimer(minutes)
  showSleepPicker.value = false
}

function onCrossfadeChange(e: Event) {
  const val = parseInt((e.target as HTMLInputElement).value) || 0
  usePlayer().crossfadeDuration = val
}

const qualityOptions: AudioQuality[] = ['low', 'medium', 'high']
const qualityLabel = computed(() => {
  const q = audioQuality
  if (q === 'low') return 'Low'
  if (q === 'medium') return 'Medium'
  return 'High'
})

function cycleQuality() {
  const p = usePlayer()
  const idx = qualityOptions.indexOf(p.audioQuality)
  const nextIdx = (idx + 1) % qualityOptions.length
  p.audioQuality = qualityOptions[nextIdx]!
}

function toggleLike() {
  liked.value = !liked.value
}

function onTogglePiP() {
  emit('toggle-pip')
  emit('close')
}

function onAddToPlaylist() {
  emit('close')
}

function onShare() {
  if (currentTrack.value) {
    const url = `${window.location.origin}/track/${currentTrack.value.id}`
    navigator.clipboard?.writeText(url)
  }
  emit('close')
}

function onTrackInfo() {
  if (currentTrack.value?.id) {
    router.push(`/track/${currentTrack.value.id}`)
  }
  emit('close')
}

function closeMenu() {
  showSleepPicker.value = false
  emit('close')
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    closeMenu()
    return
  }
  const items = menuRef.value?.querySelectorAll('[role="menuitem"]')
  if (!items?.length) return
  const currentIndex = Array.from(items).indexOf(document.activeElement as HTMLElement)
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    const next = (currentIndex + 1) % items.length
    ;(items[next] as HTMLElement).focus()
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    const prev = (currentIndex - 1 + items.length) % items.length
    ;(items[prev] as HTMLElement).focus()
  }
}

function onClickOutside(e: MouseEvent) {
  if (menuRef.value && !menuRef.value.contains(e.target as Node)) {
    closeMenu()
  }
}

function onScroll() {
  closeMenu()
}

onMounted(async () => {
  await nextTick()
  document.addEventListener('click', onClickOutside)
  document.addEventListener('scroll', onScroll, true)
  menuRef.value?.focus()
  // Focus first menuitem
  const firstItem = menuRef.value?.querySelector('[role="menuitem"]') as HTMLElement
  firstItem?.focus()
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onClickOutside)
  document.removeEventListener('scroll', onScroll, true)
})
</script>

<style scoped>
.menu-scale-enter-active {
  transition: opacity 200ms cubic-bezier(0.34, 1.56, 0.64, 1), transform 200ms cubic-bezier(0.34, 1.56, 0.64, 1);
}
.menu-scale-leave-active {
  transition: opacity 100ms ease, transform 100ms ease;
}
.menu-scale-enter-from {
  opacity: 0;
  transform: scale(0.95);
}
.menu-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

@media (prefers-reduced-motion: reduce) {
  .menu-scale-enter-active,
  .menu-scale-leave-active {
    transition: none;
  }
  .menu-scale-enter-from,
  .menu-scale-leave-to {
    opacity: 0;
    transform: none;
  }
}
</style>
