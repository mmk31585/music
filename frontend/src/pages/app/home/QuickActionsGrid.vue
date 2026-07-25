<template>
  <section class="mt-12">
    <div class="mb-5">
      <p class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">Quick access</p>
      <h2 class="mt-1 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">Discover</h2>
    </div>
    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
      <RouterLink
        v-for="action in actions"
        :key="action.to"
        :to="action.to"
        class="group flex items-center gap-4 rounded-xl border border-white/[3%] bg-white/[1%] p-4 transition-all hover:-translate-y-0.5 hover:border-white/[10%] hover:bg-white/[4%]"
      >
        <div
          class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl"
          :class="action.bgClass"
        >
          <component :is="action.icon" class="h-5 w-5" :class="action.iconClass" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-bold text-white">{{ action.title }}</p>
          <p class="mt-0.5 text-xs text-white/40">{{ action.subtitle }}</p>
        </div>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="h-4 w-4 shrink-0 text-white/20 transition group-hover:text-white/50"><polyline points="9 18 15 12 9 6"/></svg>
      </RouterLink>
    </div>
  </section>
</template>

<script setup lang="ts">
import { h } from 'vue'
import { Library, Radio, Search, Sparkles } from 'lucide-vue-next'

interface QuickAction {
  to: string
  title: string
  subtitle: string
  icon: ReturnType<typeof h>
  bgClass: string
  iconClass: string
}

const actions: QuickAction[] = [
  {
    to: '/library',
    title: 'Library',
    subtitle: 'Your saved music',
    icon: h(Library),
    bgClass: 'bg-violet-500/20',
    iconClass: 'text-violet-300',
  },
  {
    to: '/search',
    title: 'Search',
    subtitle: 'Find any track or artist',
    icon: h(Search),
    bgClass: 'bg-spotify/20',
    iconClass: 'text-spotify',
  },
  {
    to: '/ai/playlist-generator',
    title: 'AI Playlist',
    subtitle: 'Generate with AI',
    icon: h(Sparkles),
    bgClass: 'bg-amber-500/20',
    iconClass: 'text-amber-300',
  },
  {
    to: '/radio',
    title: 'Radio',
    subtitle: 'Endless music',
    icon: h(Radio),
    bgClass: 'bg-rose-500/20',
    iconClass: 'text-rose-300',
  },
]
</script>
