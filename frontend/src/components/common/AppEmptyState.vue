<template>
  <div
    class="flex flex-col items-center gap-3 py-16 text-center"
    :class="[
      bordered ? 'rounded-3xl border border-white/10 bg-black/20 px-6' : '',
      variant === 'error' ? 'border-red-500/20 bg-red-500/3' : '',
      variant === 'offline' ? 'border-amber-500/20 bg-amber-500/3' : '',
    ]"
  >
    <div
      v-if="icon"
      class="flex h-16 w-16 items-center justify-center"
      :class="[
        bordered ? 'rounded-full bg-white/10' : 'rounded-2xl bg-white/4 ring-1 ring-white/6',
        variant === 'error' ? 'text-red-400' : '',
        variant === 'offline' ? 'text-amber-400' : '',
      ]"
    >
      <i aria-hidden="true" :class="[icon, iconClass]" />
    </div>
    <h3
      class="text-white"
      :class="[
        bordered ? 'mt-5 text-xl font-black' : 'text-base font-semibold',
        variant === 'error' ? 'text-red-300' : '',
      ]"
    >
      {{ title }}
    </h3>
    <p v-if="description" class="max-w-sm text-sm" :class="bordered ? 'mt-2 text-slate-400' : 'text-slate-500'">
      {{ description }}
    </p>
    <div v-if="$slots.action" class="mt-2">
      <slot name="action" />
    </div>
    <slot v-else />
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  icon?: string
  title?: string
  description?: string
  iconClass?: string
  bordered?: boolean
  variant?: 'empty' | 'error' | 'offline'
}>(), {
  icon: 'pi pi-inbox',
  iconClass: 'text-2xl text-slate-500',
  bordered: false,
  variant: 'empty',
})
</script>
