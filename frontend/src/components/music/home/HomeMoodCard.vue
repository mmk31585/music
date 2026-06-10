<template>
  <button
    type="button"
    class="group relative flex h-24 w-full cursor-pointer items-end overflow-hidden rounded-2xl p-4 text-right transition-all duration-200 ease-spring hover:scale-[1.03] hover:brightness-110 focus-visible:ring-2 focus-visible:ring-[#1db954] focus-visible:ring-offset-2 focus-visible:outline-none md:h-28"
    :style="{ background: gradient }"
    @click="$emit('click', moodId)"
  >
    <span
      class="absolute top-3 right-3 text-2xl transition-all duration-300 group-hover:scale-125 group-hover:-rotate-6 md:top-4 md:right-4 md:text-3xl"
    >
      {{ icon }}
    </span>
    <div class="absolute inset-0 rounded-2xl border-2 border-white/0 transition-all duration-300 group-hover:border-white/10" />
    <p class="relative z-10 text-base font-bold text-white drop-shadow-xl md:text-lg">
      {{ label }}
    </p>
  </button>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  icon: string
  gradient: string
  moodId: string
}>()

defineEmits<{
  click: [moodId: string]
}>()
</script>
