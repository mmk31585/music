<template>
  <div class="rounded-2xl border border-border-subtle bg-surface-hover/30 px-2 py-2">
    <p class="px-2 pb-1 pt-1 text-[11px] font-bold tracking-[0.15em] text-muted uppercase select-none">
      {{ section.label }}
    </p>
    <div class="space-y-0.5">
      <SidebarItem
        v-for="item in section.items"
        :key="item.id"
        :item="item"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { NavSection } from './types'
import SidebarItem from './SidebarItem.vue'

defineProps<{
  section: NavSection
}>()
</script>
