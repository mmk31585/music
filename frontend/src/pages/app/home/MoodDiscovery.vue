<template>
  <section class="mt-12">
    <div class="mb-5">
      <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Discover by mood</p>
      <h2 class="mt-1 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">What's your mood?</h2>
    </div>
    <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
      <button
        v-for="mood in moods"
        :key="mood.id"
        type="button"
        class="group relative flex aspect-square cursor-pointer flex-col items-center justify-center gap-2 overflow-hidden rounded-2xl p-4 text-center transition-all duration-300 hover:scale-[1.03] hover:brightness-110 focus-visible:ring-2 focus-visible:ring-spotify focus-visible:ring-offset-2 focus-visible:outline-hidden"
        :style="{ background: mood.gradient }"
        @click="$emit('select', mood.id)"
      >
        <div class="absolute inset-0 rounded-2xl border-2 border-white/0 transition-all duration-300 group-hover:border-white/15" />
        <component :is="mood.icon" class="relative z-10 h-8 w-8 text-white drop-shadow-lg" />
        <span class="relative z-10 text-sm font-bold text-white drop-shadow-xl">{{ mood.label }}</span>
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { CloudSun, Coffee, Heart, Moon, Music, Sunset, Waves, Zap } from 'lucide-vue-next'

interface MoodItem {
  id: string
  label: string
  icon: ReturnType<typeof h>
  gradient: string
}

defineEmits<{
  select: [moodId: string]
}>()

const moods: MoodItem[] = [
  { id: 'happy', label: 'Happy', icon: h(Heart), gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { id: 'chill', label: 'Chill', icon: h(CloudSun), gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' },
  { id: 'energetic', label: 'Energetic', icon: h(Zap), gradient: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)' },
  { id: 'focus', label: 'Focus', icon: h(Music), gradient: 'linear-gradient(135deg, #a18cd1 0%, #fbc2eb 100%)' },
  { id: 'night', label: 'Night', icon: h(Moon), gradient: 'linear-gradient(135deg, #0c3483 0%, #a2b6df 100%)' },
  { id: 'morning', label: 'Morning', icon: h(Coffee), gradient: 'linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%)' },
  { id: 'sunset', label: 'Sunset', icon: h(Sunset), gradient: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)' },
  { id: 'ambient', label: 'Ambient', icon: h(Waves), gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' },
  { id: 'romantic', label: 'Romantic', icon: h(Heart), gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' },
  { id: 'workout', label: 'Workout', icon: h(Zap), gradient: 'linear-gradient(135deg, #f12711 0%, #f5af19 100%)' },
]
</script>
