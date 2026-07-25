<template>
  <header class="mb-8">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div class="min-w-0">
        <p
          v-if="eyebrow"
          class="mb-1.5 text-xs font-bold uppercase tracking-[0.15em] text-emerald-400"
        >
          {{ eyebrow }}
        </p>

        <h1 class="text-2xl font-bold tracking-tight text-white md:text-3xl">
          {{ title }}
        </h1>

        <p v-if="description" class="mt-2 max-w-2xl text-sm leading-relaxed text-slate-400">
          {{ description }}
        </p>
      </div>

      <div v-if="$slots.actions" class="flex shrink-0 items-center gap-3">
        <slot name="actions" />
      </div>
    </div>

    <div class="mt-6 h-px bg-linear-to-r from-white/10 via-white/5 to-transparent" />
  </header>
</template>

<script setup lang="ts">
defineProps<{
  eyebrow?: string
  title: string
  description?: string
}>()
</script>
