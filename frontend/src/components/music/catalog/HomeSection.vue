<template>
  <section class="space-y-4">
    <div class="flex items-end justify-between gap-4">
      <div>
        <p v-if="eyebrow" class="text-xs font-semibold tracking-[0.2em] text-[#1db954] uppercase">
          {{ eyebrow }}
        </p>
        <h2 class="mt-1 text-2xl font-black text-white">
          {{ title }}
        </h2>
      </div>
      <RouterLink
        v-if="actionTo"
        :to="actionTo"
        class="shrink-0 text-xs font-semibold text-white/30 transition hover:text-white/70"
      >
        {{ actionLabel || 'See all' }}
        <i aria-hidden="true" class="pi pi-chevron-left mr-1 text-[10px]" />
      </RouterLink>
      <slot name="action" />
    </div>
    <slot />
  </section>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  eyebrow?: string
  actionLabel?: string
  actionTo?: string
}>()
</script>
