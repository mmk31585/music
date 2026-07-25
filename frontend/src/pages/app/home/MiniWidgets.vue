<template>
  <section class="mt-12">
    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <div
        v-for="widget in widgets"
        :key="widget.label"
        class="flex items-center gap-3 rounded-xl border border-white/[4%] bg-white/[1%] px-4 py-3"
      >
        <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg" :class="widget.bgClass">
          <component :is="widget.icon" class="h-4 w-4" :class="widget.iconClass" />
        </div>
        <div>
          <p class="text-lg font-bold text-white">{{ widget.value }}</p>
          <p class="text-[11px] text-white/40">{{ widget.label }}</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { Award, BarChart3, Clock, Headphones } from 'lucide-vue-next'

interface Widget {
  label: string
  value: string
  icon: ReturnType<typeof h>
  bgClass: string
  iconClass: string
}

defineProps<{
  minutesListened?: number
  tracksPlayed?: number
  topGenre?: string
  streak?: number
}>()

const widgets: Widget[] = [
  { label: 'Minutes Listened', value: '—', icon: h(Clock), bgClass: 'bg-spotify/15', iconClass: 'text-spotify' },
  { label: 'Tracks Played', value: '—', icon: h(Headphones), bgClass: 'bg-aurora-purple/15', iconClass: 'text-aurora-purple' },
  { label: 'Top Genre', value: '—', icon: h(Award), bgClass: 'bg-amber-500/15', iconClass: 'text-amber-400' },
  { label: 'Day Streak', value: '—', icon: h(BarChart3), bgClass: 'bg-rose-500/15', iconClass: 'text-rose-300' },
]
</script>
