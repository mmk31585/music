<template>
  <div class="flex flex-col items-center justify-center py-16 text-center">
    <div
      class="mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-white/4 ring-1 ring-white/6"
    >
      <i aria-hidden="true" :class="[icon, 'text-2xl text-slate-500']" />
    </div>

    <h3 class="text-base font-semibold text-white">{{ title }}</h3>
    <p class="mt-1.5 max-w-sm text-sm text-slate-500">{{ description }}</p>

    <div v-if="$slots.action" class="mt-5">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    icon?: string
    title: string
    description?: string
  }>(),
  {
    icon: 'pi pi-inbox',
    description: '',
  },
)
</script>
