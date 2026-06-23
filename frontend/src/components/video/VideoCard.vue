<template>
  <div
    role="button"
    tabindex="0"
    class="group relative cursor-pointer overflow-hidden rounded-2xl bg-white/5 transition-all duration-150"
    :class="hoverClass"
    @click="$emit('open')"
    @keydown.enter="$emit('open')"
    @keydown.space.prevent="$emit('open')"
  >
    <!-- Thumbnail -->
    <div class="relative aspect-9/16 w-full overflow-hidden">
      <AppImage
        :src="video.thumbnail_url || video.thumbnail_path || video.track_cover_url"
        :alt="video.title"
        class="h-full w-full object-cover transition-transform duration-150 group-hover:scale-[1.03]"
        fallback-icon="pi pi-video"
        icon-size="2rem"
      />

      <!-- Overlay gradient -->
      <div class="pointer-events-none absolute inset-0 bg-linear-to-t from-black/70 via-transparent to-transparent" />

      <!-- Status badge (for processing/failed) -->
      <div
        v-if="video.status && video.status !== 'ready'"
        class="absolute top-2 left-2 z-10 flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold backdrop-blur-xs"
        :class="statusBadgeClass"
      >
        <i v-if="video.status === 'processing'" class="pi pi-spin pi-spinner text-[10px]" />
        <i v-else-if="video.status === 'failed'" class="pi pi-exclamation-circle text-[10px]" />
        <span>{{ statusLabel }}</span>
      </div>

      <!-- Type badge -->
      <div
        class="absolute top-2 right-2 flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-bold backdrop-blur-xs"
        :class="typeBadgeClass"
      >
        <i aria-hidden="true" :class="typeIcon" class="text-[10px]" />
        <span>{{ video.type === 'official_mv' ? 'MV' : 'Edit' }}</span>
      </div>

      <!-- Bottom stats -->
      <div class="absolute right-2 bottom-2 left-2 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <span class="flex items-center gap-1 text-[11px] font-medium text-white/80">
            <i aria-hidden="true" class="pi pi-eye text-[10px]" />
            {{ formatCount(video.view_count) }}
          </span>
          <span class="flex items-center gap-1 text-[11px] font-medium text-white/80">
            <i aria-hidden="true" class="pi pi-heart text-[10px]" />
            {{ formatCount(video.like_count) }}
          </span>
        </div>

        <div
          class="flex h-7 w-7 items-center justify-center rounded-full bg-black/50 text-white opacity-0 backdrop-blur-xs transition-opacity duration-150 group-hover:opacity-100"
        >
          <i aria-hidden="true" class="pi pi-play-fill text-xs" />
        </div>
      </div>
    </div>

    <!-- Title (below thumbnail) -->
    <div class="p-2.5">
      <p class="truncate text-xs font-semibold text-white/90">{{ video.title }}</p>
      <p v-if="video.uploader" class="mt-0.5 truncate text-[10px] text-white/40">
        {{ video.uploader.username }}
      </p>
    </div>

    <!-- Loading skeleton (shown while image loads) -->
    <div
      v-if="!loaded"
      class="absolute inset-0 rounded-2xl bg-white/5"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { VideoItem } from '@/services/api/video/types'
import AppImage from '@/components/common/AppImage.vue'
import { formatCount } from '@/utils/number'

const props = defineProps<{
  video: VideoItem
}>()

defineEmits<{
  open: []
}>()

const loaded = ref(false)

const hoverClass = computed(() => {
  return 'hover:bg-white/8 hover:shadow-lg hover:shadow-black/20'
})

const typeBadgeClass = computed(() => {
  if (props.video.type === 'official_mv') {
    return 'bg-spotify/20 text-spotify'
  }
  return 'bg-blue-500/20 text-blue-400'
})

const typeIcon = computed(() => {
  if (props.video.type === 'official_mv') {
    return 'pi pi-play-circle'
  }
  return 'pi pi-pencil'
})

const statusLabel = computed(() => {
  if (props.video.status === 'processing') return 'Processing'
  if (props.video.status === 'failed') return 'Failed'
  return ''
})

const statusBadgeClass = computed(() => {
  if (props.video.status === 'processing') return 'bg-yellow-500/20 text-yellow-400'
  if (props.video.status === 'failed') return 'bg-red-500/20 text-red-400'
  return ''
})
</script>
