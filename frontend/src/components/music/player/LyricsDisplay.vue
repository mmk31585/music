<template>
  <div class="space-y-4">
    <div v-if="loading" class="space-y-3">
      <div
        v-for="i in 8"
        :key="i"
        class="h-6 w-full animate-pulse rounded bg-white/[0.06]"
        :style="{ width: `${60 + Math.random() * 30}%` }"
      />
    </div>

    <div
      v-else-if="error"
      class="flex flex-col items-center gap-3 rounded-2xl border border-white/[0.06] bg-white/[0.03] px-6 py-12 text-center"
    >
      <i aria-hidden="true" class="pi pi-align-left text-2xl text-slate-500" />
      <p class="text-sm text-slate-400">No lyrics available</p>
      <button
        v-if="onAddLyrics"
        type="button"
        class="text-xs font-medium text-[#1db954] underline underline-offset-2 transition hover:text-[#1ed760]"
        @click="onAddLyrics"
      >
        Add lyrics
      </button>
    </div>

    <div v-else-if="lyrics" class="space-y-1">
      <div class="mb-4 flex items-center gap-2">
        <span class="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-slate-400">
          {{ lyrics.language || 'en' }}
        </span>
        <span class="rounded-full bg-white/10 px-3 py-1 text-xs font-medium text-slate-400">
          {{ lyrics.type || 'plain' }}
        </span>
      </div>

      <div
        class="text-sm leading-relaxed whitespace-pre-line text-slate-300"
        :class="{ 'text-center': centered }"
        :dir="isRtl ? 'rtl' : 'ltr'"
      >
        {{ lyrics.content }}
      </div>
    </div>

    <div
      v-else
      class="flex flex-col items-center gap-3 rounded-2xl border border-white/[0.06] bg-white/[0.03] px-6 py-12 text-center"
    >
      <i aria-hidden="true" class="pi pi-align-left text-2xl text-slate-500" />
      <p class="text-sm text-slate-400">No lyrics available</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Lyrics } from '@/services/api/lyrics'

const props = defineProps<{
  lyrics: Lyrics | null
  loading?: boolean
  error?: boolean
  centered?: boolean
  onAddLyrics?: () => void
}>()

const isRtl = computed(() => {
  const lang = props.lyrics?.language?.toLowerCase()
  return lang === 'fa' || lang === 'far' || lang?.startsWith('fa-') || lang?.startsWith('fa_')
})
</script>
