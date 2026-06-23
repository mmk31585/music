<template>
  <div
    class="flex gap-1 rounded-xl bg-white/4 p-1 overflow-x-auto scrollbar-none"
    role="tablist"
    :aria-label="ariaLabel"
  >
    <button
      v-for="tab in tabs"
      :key="tab.key"
      role="tab"
      :aria-selected="activeTab === tab.key"
      :aria-controls="`tabpanel-${tab.key}`"
      class="flex items-center gap-2 rounded-lg px-4 py-2.5 text-sm font-medium transition-all duration-200 whitespace-nowrap focus-visible:outline-2 focus-visible:outline-[#1db954]"
      :class="activeTab === tab.key
        ? 'bg-white/10 text-white shadow-lg'
        : 'text-white/30 hover:text-white/50'"
      @click="$emit('update:activeTab', tab.key)"
    >
      <i v-if="tab.icon" aria-hidden="true" :class="tab.icon" class="text-xs" />
      {{ tab.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  tabs: { key: string; label: string; icon?: string }[]
  activeTab: string
  ariaLabel?: string
}>()

defineEmits<{
  'update:activeTab': [key: string]
}>()
</script>
