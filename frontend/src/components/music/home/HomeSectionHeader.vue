<template>
  <div class="mb-4 flex items-end justify-between gap-4">
    <div>
      <h2 class="text-xl font-bold text-white md:text-2xl" style="letter-spacing: -0.02em">
        {{ title }}
      </h2>
      <p v-if="subtitle" class="mt-0.5 text-sm text-white/60">
        {{ subtitle }}
      </p>
    </div>
    <RouterLink
      v-if="seeAllRoute"
      :to="seeAllRoute"
      class="shrink-0 whitespace-nowrap text-sm font-medium text-white/60 transition hover:text-white hover:underline"
    >
      {{ seeAllLabel || 'مشاهده همه' }}
    </RouterLink>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
  seeAllRoute?: string
  seeAllLabel?: string
}>()
</script>
