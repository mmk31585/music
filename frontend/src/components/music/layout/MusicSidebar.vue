<template>
  <nav
    aria-label="Main navigation"
    class="hidden h-screen w-[260px] shrink-0 border-r border-border-subtle bg-surface-raised/60 backdrop-blur-2xl lg:flex lg:flex-col"
  >
    <div class="flex shrink-0 flex-col px-4 pb-4 pt-4">
      <SidebarLogo />
    </div>

    <div
      class="flex-1 overflow-y-auto px-4 scroll-bar"
      :class="bottomPadding"
    >
      <div class="flex flex-col gap-3 pt-1">
        <SidebarSection
          v-for="section in visibleSections"
          :key="section.id"
          :section="section"
        />
      </div>
      <div class="h-4" />
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { usePlayerStore } from '@/stores/player'
import SidebarLogo from './sidebar/SidebarLogo.vue'
import SidebarSection from './sidebar/SidebarSection.vue'
import { useNavigation } from './sidebar/useNavigation'

const playerStore = usePlayerStore()
const { visibleSections } = useNavigation()

const barCollapsed = ref(localStorage.getItem('player-bar-collapsed') === 'true')

onMounted(() => {
  const handler = (e: Event) => {
    barCollapsed.value = (e as CustomEvent).detail
  }
  window.addEventListener('playerbar-collapse', handler)
})

const bottomPadding = computed(() => {
  if (!playerStore.currentTrack) return 'pb-4'
  return barCollapsed.value ? 'pb-20' : 'pb-32'
})
</script>
