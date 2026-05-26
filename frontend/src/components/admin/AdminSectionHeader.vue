<template>
  <div class="mb-6 flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
    <div>
      <p v-if="eyebrow" class="text-xs tracking-[0.3em] text-slate-500 uppercase">
        {{ eyebrow }}
      </p>
      <h1 class="mt-1 text-3xl font-extrabold text-white">
        {{ title }}
      </h1>
      <p v-if="description" class="mt-2 max-w-2xl text-sm text-slate-400">
        {{ description }}
      </p>
    </div>

    <div v-if="$slots.actions" class="flex items-center gap-3">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  eyebrow?: string
  description?: string
}>()
</script>
