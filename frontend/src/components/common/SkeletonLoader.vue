<template>
  <div :class="containerCls">
    <!-- Card variant -->
    <template v-if="variant === 'card'">
      <div class="rounded-2xl bg-white/4 p-4">
        <Skeleton class="mb-3 aspect-square w-full rounded-xl" />
        <div class="space-y-2">
          <Skeleton class="h-4 w-3/4 rounded" />
          <Skeleton class="h-3 w-1/2 rounded" />
        </div>
      </div>
    </template>

    <!-- Track row variant -->
    <template v-else-if="variant === 'track'">
      <div class="flex items-center gap-3 rounded-xl px-3 py-2">
        <Skeleton class="h-11 w-11 shrink-0 rounded-xl" />
        <div class="min-w-0 flex-1 space-y-2">
          <Skeleton class="h-4 w-2/3 rounded" />
          <Skeleton class="h-3 w-1/3 rounded" />
        </div>
        <Skeleton class="h-3 w-10 rounded" />
      </div>
    </template>

    <!-- Hero variant -->
    <template v-else-if="variant === 'hero'">
      <div class="flex flex-col items-center gap-6 py-12 md:flex-row md:items-end md:gap-8">
        <Skeleton shape="circle" class="h-56 w-56 md:h-64 md:w-64" />
        <div class="flex-1 space-y-3">
          <Skeleton class="mx-auto h-5 w-24 rounded md:mx-0" />
          <Skeleton class="mx-auto h-8 w-48 rounded md:mx-0 md:w-64" />
          <Skeleton class="mx-auto h-4 w-32 rounded md:mx-0" />
        </div>
      </div>
    </template>

    <!-- Lines variant (default) -->
    <template v-else>
      <div :class="['space-y-2', cls]">
        <Skeleton
          v-for="i in lines"
          :key="i"
          class="rounded-sm"
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
