<template>
  <div class="rounded-lg border border-white/6 bg-white/3 p-3">
    <div class="flex items-center gap-3">
      <Button
        :icon="playing ? 'pi pi-pause' : 'pi pi-play'"
        :severity="playing ? 'secondary' : 'primary'"
        rounded
        size="small"
        @click="toggle"
        :aria-label="playing ? 'Pause preview' : 'Play preview'"
      />
      <div class="flex-1">
        <div
          ref="seekBar"
          class="relative h-1.5 cursor-pointer rounded-full bg-white/8"
          role="slider"
          :aria-label="'Audio preview seek'"
          :aria-valuemin="0"
          :aria-valuemax="duration"
          :aria-valuenow="currentTime"
          tabindex="0"
          @click="seek"
          @keydown.left.prevent="skip(-5)"
          @keydown.right.prevent="skip(5)"
        >
          <div
            class="absolute left-0 top-0 h-full rounded-full bg-emerald-500 transition-all duration-150"
            :style="{ width: `${duration ? (currentTime / duration) * 100 : 0}%` }"
          />
        </div>
      </div>
      <span class="w-20 text-right text-[11px] text-slate-400 tabular-nums">
        {{ formatTime(currentTime) }} / {{ formatTime(duration) }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  src: string | null
}>()

const emit = defineEmits<{
  loaded: [duration: number]
  timeupdate: [currentTime: number]
  playstate: [playing: boolean]
}>()

/** Exposed for external seek control (e.g. synced lyrics) */
function setCurrentTime(seconds: number) {
  audio.currentTime = seconds
  currentTime.value = seconds
}

defineExpose({ setCurrentTime })

const audio = new Audio()
audio.preload = 'metadata'

const playing = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const seekBar = ref<HTMLElement | null>(null)
let animFrameId = 0

watch(
  () => props.src,
  (url) => {
    if (url) {
      audio.src = url
      audio.load()
    }
  },
  { immediate: true },
)

function onAudioLoaded() {
  duration.value = audio.duration
  emit('loaded', audio.duration)
}

function onAudioPlay() {
  playing.value = true
  emit('playstate', true)
  animFrameId = requestAnimationFrame(onAudioTime)
}

function onAudioPause() {
  playing.value = false
  emit('playstate', false)
  cancelAnimationFrame(animFrameId)
}

function onAudioEnded() {
  playing.value = false
  currentTime.value = 0
  emit('playstate', false)
  emit('timeupdate', 0)
}

function onAudioTime() {
  currentTime.value = audio.currentTime
  emit('timeupdate', audio.currentTime)
  animFrameId = requestAnimationFrame(onAudioTime)
}

function toggle() {
  if (playing.value) {
    audio.pause()
  } else {
    audio.play()
  }
}

function seek(e: MouseEvent) {
  const bar = seekBar.value
  if (!bar || !duration.value) return
  const rect = bar.getBoundingClientRect()
  const ratio = (e.clientX - rect.left) / rect.width
  audio.currentTime = ratio * duration.value
  currentTime.value = audio.currentTime
}

function skip(seconds: number) {
  audio.currentTime = Math.max(0, Math.min(audio.currentTime + seconds, duration.value))
  currentTime.value = audio.currentTime
}

function formatTime(seconds: number): string {
  if (!seconds || !Number.isFinite(seconds)) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

onMounted(() => {
  audio.addEventListener('loadedmetadata', onAudioLoaded)
  audio.addEventListener('play', onAudioPlay)
  audio.addEventListener('pause', onAudioPause)
  audio.addEventListener('ended', onAudioEnded)
})

onUnmounted(() => {
  audio.pause()
  audio.src = ''
  audio.load()
  cancelAnimationFrame(animFrameId)
  audio.removeEventListener('loadedmetadata', onAudioLoaded)
  audio.removeEventListener('play', onAudioPlay)
  audio.removeEventListener('pause', onAudioPause)
  audio.removeEventListener('ended', onAudioEnded)
})
</script>
