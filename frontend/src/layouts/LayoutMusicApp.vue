<template>
  <div class="min-h-screen bg-black text-white">
    <div class="grid min-h-screen grid-cols-1 lg:grid-cols-[260px_1fr]">
      <MusicSidebar class="hidden lg:block" />

      <div class="flex min-h-screen flex-col">
        <MusicTopbar />
        <main class="flex-1 overflow-y-auto bg-gradient-to-b from-[#1a1a1a] to-black pb-28">
          <router-view />
        </main>
      </div>
    </div>

    <NowPlayingBar />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { MusicSidebar, MusicTopbar, NowPlayingBar } from '@/components/music'
import { usePlayerStore } from '@/stores'

const player = usePlayerStore()

function isTypingTarget(target: EventTarget | null) {
  if (!(target instanceof HTMLElement)) return false

  const tag = target.tagName.toLowerCase()
  return tag === 'input' || tag === 'textarea' || tag === 'select' || target.isContentEditable
}

function onKeydown(event: KeyboardEvent) {
  if (isTypingTarget(event.target)) return

  if (event.code === 'Space') {
    event.preventDefault()
    void player.togglePlay()
  }

  if (event.code === 'ArrowRight' && event.shiftKey) {
    event.preventDefault()
    void player.playNext()
  }

  if (event.code === 'ArrowLeft' && event.shiftKey) {
    event.preventDefault()
    void player.playPrevious()
  }
}

onMounted(() => {
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>
