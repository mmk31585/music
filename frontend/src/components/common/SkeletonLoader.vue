<template>
  <div :class="containerCls">
    <!-- Card variant -->
    <template v-if="variant === 'card'">
      <div class="animate-pulse rounded-2xl bg-white/[0.04] p-4">
        <div class="mb-3 aspect-square overflow-hidden rounded-xl bg-white/10 shimmer-base" />
        <div class="space-y-2">
          <div class="h-4 w-3/4 rounded bg-white/10 shimmer-base" />
          <div class="h-3 w-1/2 rounded bg-white/10 shimmer-base" />
        </div>
      </div>
    </template>

    <!-- Track row variant -->
    <template v-else-if="variant === 'track'">
      <div class="flex items-center gap-3 rounded-xl px-3 py-2">
        <div class="h-11 w-11 shrink-0 animate-pulse rounded-xl bg-white/10 shimmer-base" />
        <div class="min-w-0 flex-1 space-y-2">
          <div class="h-4 w-2/3 rounded bg-white/10 shimmer-base" />
          <div class="h-3 w-1/3 rounded bg-white/10 shimmer-base" />
        </div>
        <div class="h-3 w-10 rounded bg-white/10 shimmer-base" />
      </div>
    </template>

    <!-- Hero variant -->
    <template v-else-if="variant === 'hero'">
      <div class="flex flex-col items-center gap-6 py-12 md:flex-row md:items-end md:gap-8">
        <div class="h-56 w-56 animate-pulse rounded-full bg-white/10 shimmer-base md:h-64 md:w-64" />
        <div class="flex-1 space-y-3">
          <div class="mx-auto h-5 w-24 rounded bg-white/10 shimmer-base md:mx-0" />
          <div class="mx-auto h-8 w-48 rounded bg-white/10 shimmer-base md:mx-0 md:w-64" />
          <div class="mx-auto h-4 w-32 rounded bg-white/10 shimmer-base md:mx-0" />
        </div>
      </div>
    </template>

    <!-- Lines variant (default) -->
    <template v-else>
      <div :class="['space-y-2', cls]">
        <div
          v-for="i in lines"
          :key="i"
          class="animate-pulse rounded bg-white/10 shimmer-base"
          :style="{
            width: typeof width === 'string' ? width : Array.isArray(width) ? width[i - 1] || width[0] : '100%',
            height: typeof height === 'string' ? height : Array.isArray(height) ? height[i - 1] || height[0] : '16px',
          }"
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type SkeletonVariant = 'lines' | 'card' | 'track' | 'hero'

const props = withDefaults(
  defineProps<{
    variant?: SkeletonVariant
    lines?: number
    width?: string | string[]
    height?: string | string[]
    class?: string
  }>(),
  {
    variant: 'lines',
    lines: 3,
    width: '100%',
    height: '16px',
    class: '',
  },
)

const containerCls = computed(() => props.class || undefined)
const cls = computed(() => props.class || undefined)
</script>

<style scoped>
@keyframes shimmer-slide {
  0% {
    background-position: -200% 0;
  }
  100% {
    background-position: 200% 0;
  }
}

.shimmer-base {
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.03) 0%,
    rgba(255, 255, 255, 0.08) 50%,
    rgba(255, 255, 255, 0.03) 100%
  );
  background-size: 200% 100%;
  animation: shimmer-slide 1.5s ease-in-out infinite;
}
</style>
