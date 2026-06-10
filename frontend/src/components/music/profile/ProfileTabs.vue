<template>
  <div class="border-b border-white/5">
    <div class="flex gap-6 px-1">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        :class="[
          'relative pb-3 text-sm font-medium transition',
          activeTab === tab.key ? 'text-white' : 'text-slate-500 hover:text-slate-300',
        ]"
        @click="$emit('update:activeTab', tab.key)"
      >
        {{ tab.label }}
        <div
          v-if="activeTab === tab.key"
          class="absolute right-0 bottom-0 left-0 h-0.5 rounded-full bg-[#1db954]"
        />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  tabs: { key: string; label: string }[]
  activeTab: string
}>()

defineEmits<{
  'update:activeTab': [key: string]
}>()
</script>
