<template>
  <button
    v-if="item.future"
    type="button"
    disabled
    class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-tertiary opacity-50 cursor-not-allowed"
    :title="`${item.label} (coming soon)`"
  >
    <component :is="item.icon" :size="18" class="shrink-0" />
    <span class="truncate">{{ item.label }}</span>
    <span class="ml-auto text-[9px] font-bold text-muted uppercase tracking-wider">Soon</span>
  </button>

  <RouterLink
    v-else-if="item.to"
    :to="item.to"
    class="group relative flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all duration-150"
    :class="activeClass"
    :aria-current="isActive ? 'page' : undefined"
  >
    <div
      v-if="isActive"
      class="absolute start-0 top-1/2 h-5 w-0.5 -translate-y-1/2 rounded-full bg-accent"
    />
    <component :is="item.icon" :size="18" class="shrink-0" :class="iconClass" />
    <span class="truncate" :class="labelClass">{{ item.label }}</span>
    <SidebarBadge v-if="item.badge" :value="item.badge" :variant="item.badgeVariant" class="ml-auto" />
  </RouterLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import type { NavItem } from './types'
import SidebarBadge from './SidebarBadge.vue'

const props = defineProps<{
  item: NavItem
}>()

const route = useRoute()

const isActive = computed(() => {
  if (props.item.to === '/') return route.path === '/'
  if (!props.item.to) return false
  return route.path.startsWith(props.item.to)
})

const activeClass = computed(() => {
  if (!isActive.value) return 'text-secondary hover:bg-surface-hover hover:text-primary'
  return 'bg-surface-active text-primary shadow-xs'
})

const iconClass = computed(() => {
  return isActive.value ? 'text-accent' : 'text-tertiary group-hover:text-accent'
})

const labelClass = computed(() => {
  return isActive.value ? 'font-bold' : ''
})
</script>
