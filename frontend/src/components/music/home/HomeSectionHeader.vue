<template>
  <div class="mb-5 flex items-end justify-between gap-4">
    <div>
      <p v-if="eyebrow" class="text-[10px] font-bold tracking-[0.3em] text-white/30 uppercase">
        {{ eyebrow }}
      </p>
      <h2 class="mt-1 text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">
        {{ title }}
      </h2>
      <p v-if="subtitle" class="mt-1 text-sm text-white/50">
        {{ subtitle }}
      </p>
    </div>
    <RouterLink
      v-if="seeAllRoute"
      :to="seeAllRoute"
      class="shrink-0 whitespace-nowrap text-xs font-medium text-white/30 transition hover:text-white"
    >
      {{ seeAllLabel || 'See all' }}
      <i aria-hidden="true" class="pi pi-chevron-left ml-0.5 text-[10px]" />
    </RouterLink>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  eyebrow?: string
  subtitle?: string
  seeAllRoute?: string
  seeAllLabel?: string
}>()
</script>
