<template>
  <div class="rounded-3xl border border-white/10 bg-white/5 p-5 shadow-lg backdrop-blur">
    <div class="flex items-start justify-between gap-4">
      <div>
        <p class="text-sm text-slate-400">{{ label }}</p>
        <h3 class="mt-2 text-3xl font-bold text-white">{{ value }}</h3>
        <p v-if="hint" class="mt-2 text-xs text-slate-500">{{ hint }}</p>
      </div>

      <div
        class="flex h-12 w-12 items-center justify-center rounded-2xl bg-[#1db954]/15 text-[#1db954]"
      >
        <i :class="icon" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  value: string | number
  hint?: string
  icon?: string
}>()
</script>
