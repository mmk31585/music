<template>
  <div
    class="flex gap-1 bg-white/4 rounded-2xl p-1 overflow-x-auto scrollbar-none"
    role="tablist"
    :aria-label="ariaLabel"
  >
    <button
      v-for="tab in tabs"
      :key="tab.key"
      role="tab"
      :aria-selected="activeTab === tab.key"
      :aria-controls="`tabpanel-${tab.key}`"
      :id="`tab-${tab.key}`"
      class="rounded-xl px-4 py-2 text-sm transition-all flex items-center gap-1.5 whitespace-nowrap focus-visible:outline-2 focus-visible:outline-primary"
      :class="activeTab === tab.key
        ? 'font-semibold text-white bg-white/10 shadow shadow-black/20'
        : 'text-white/40 hover:text-white/70'"
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
