<template>
  <div>
    <div class="mb-6">
      <h2 class="text-lg font-semibold text-white">Track</h2>
      <p class="mt-1 text-sm text-slate-400">Review and edit the track details.</p>
    </div>

    <!-- Artists summary -->
    <div class="mb-4 rounded-lg border border-white/6 bg-white/3 p-3">
      <div class="flex flex-wrap items-center gap-2">
        <span class="text-xs font-medium text-slate-500">Artists:</span>
        <template v-for="(a, i) in ctx.finalMetadata.artists" :key="i">
          <span class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs"
            :class="i === 0 ? 'bg-emerald-500/15 text-emerald-300' : 'bg-purple-500/15 text-purple-300'">
            <User aria-hidden="true" class="text-[9px]" />
            {{ a.name || 'Unnamed' }}
          </span>
          <span v-if="i < ctx.finalMetadata.artists.length - 1" class="text-white/20">,</span>
        </template>
      </div>
    </div>

    <!-- Cover image picker -->
    <div class="mb-6 flex flex-wrap gap-6">
      <div v-if="ctx.embeddedCover.value">
        <p class="mb-2 text-xs font-medium text-slate-500">Embedded Cover</p>
        <img
          :src="ctx.embeddedCover.value"
          alt="Embedded"
          class="h-24 w-24 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
          :class="ctx.finalMetadata.track.coverUrl === ctx.embeddedCover.value ? 'ring-green-500' : 'ring-transparent hover:ring-white/20'"
          @click="ctx.finalMetadata.track.coverUrl = ctx.embeddedCover.value!"
        />
      </div>
      <div v-if="ctx.suggestedAlbumCover.value">
        <p class="mb-2 text-xs font-medium text-slate-500 flex items-center gap-1">
          <Music2 aria-hidden="true" class="text-green-400"></Music2>Spotify
        </p>
        <img
          :src="ctx.suggestedAlbumCover.value"
          alt="Spotify"
          class="h-24 w-24 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
          :class="ctx.finalMetadata.track.coverUrl === ctx.suggestedAlbumCover.value ? 'ring-green-500' : 'ring-transparent hover:ring-white/20'"
          @click="ctx.finalMetadata.track.coverUrl = ctx.suggestedAlbumCover.value!"
        />
      </div>
      <div v-if="ctx.suggestedLastfmAlbumCover.value">
        <p class="mb-2 text-xs font-medium text-slate-500 flex items-center gap-1">
          <Star aria-hidden="true" class="text-yellow-400"></Star>Last.fm
        </p>
        <img
          :src="ctx.suggestedLastfmAlbumCover.value"
          alt="Last.fm"
          class="h-24 w-24 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
          :class="ctx.finalMetadata.track.coverUrl === ctx.suggestedLastfmAlbumCover.value ? 'ring-green-500' : 'ring-transparent hover:ring-white/20'"
          @click="ctx.finalMetadata.track.coverUrl = ctx.suggestedLastfmAlbumCover.value!"
        />
      </div>
      <div v-if="ctx.enrichment.value?.ml?.albumCoverUrl">
        <p class="mb-2 text-xs font-medium text-slate-500 flex items-center gap-1">
          <Bot aria-hidden="true" class="text-purple-400"></Bot>ML Server
        </p>
        <img
          :src="ctx.enrichment.value.ml.albumCoverUrl"
          alt="ML"
          class="h-24 w-24 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
          :class="ctx.finalMetadata.track.coverUrl === ctx.enrichment.value.ml.albumCoverUrl ? 'ring-green-500' : 'ring-transparent hover:ring-white/20'"
          @click="ctx.finalMetadata.track.coverUrl = ctx.enrichment.value!.ml!.albumCoverUrl"
        />
      </div>
      <div v-if="ctx.finalMetadata.track.coverUrl && ctx.finalMetadata.track.coverUrl !== ctx.embeddedCover.value && ctx.finalMetadata.track.coverUrl !== ctx.suggestedAlbumCover.value && ctx.finalMetadata.track.coverUrl !== ctx.suggestedLastfmAlbumCover.value && ctx.finalMetadata.track.coverUrl !== ctx.enrichment.value?.ml?.albumCoverUrl" class="shrink-0">
        <p class="mb-2 text-xs font-medium text-slate-500">Selected</p>
        <img
          :src="ctx.finalMetadata.track.coverUrl"
          alt="Selected"
          class="h-24 w-24 rounded-lg object-cover shadow-md ring-2 ring-green-500"
          @error="($event.target as HTMLImageElement).style.display='none'"
        />
      </div>
    </div>

    <!-- Fetch track data -->
    <div class="mb-6">
      <Button
        :loading="ctx.enrichingTrack.value"
        :disabled="ctx.enriching.value"
        label="Fetch Lyrics &amp; Cover"
        icon="pi pi-cloud-download"
        severity="help"
        size="small"
        outlined
        @click="ctx.fetchSection('track')"
      />
      <p class="mt-1.5 text-[11px] text-slate-500">
        Load track data from the track's embedded metadata (ID3 tags).
      </p>
    </div>

    <!-- Audio Preview Player -->
    <AdminAudioPreview v-if="ctx.trackAudioUrl.value" :src="ctx.trackAudioUrl.value" />

    <div class="space-y-4 rounded-lg border border-white/6 bg-white/3 p-5">
      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">
            Title
            <template v-if="ctx.findSuggestion('title')">
              <Badge :value="ctx.sourceLabel(ctx.findSuggestion('title')!.source)" :severity="ctx.sourceSeverity(ctx.findSuggestion('title')!.source)" size="small" />
              <Badge :value="ctx.findSuggestion('title')!.confidence" :severity="ctx.confidenceSeverity(ctx.findSuggestion('title')!.confidence)" size="small" />
            </template>
          </label>
          <InputText v-model="ctx.finalMetadata.track.title" dir="auto" class="w-full" placeholder="Track title" />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">Track Number</label>
          <InputNumber v-model="ctx.finalMetadata.track.trackNumber" :use-grouping="false" :min="1" class="w-full" placeholder="1" />
        </div>
      </div>

      <div class="grid gap-4 md:grid-cols-3">
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">Duration (seconds)</label>
          <InputNumber :model-value="ctx.finalMetadata.track.durationSeconds" class="w-full text-slate-400" disabled />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">Genre</label>
          <InputText v-model="ctx.finalMetadata.track.genre" dir="auto" class="w-full" placeholder="Genre" />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">
            Explicit
            <Info aria-hidden="true" class="ml-1 text-slate-500 text-[10px]"></Info>
          </label>
          <div class="flex items-center gap-2 pt-1">
            <ToggleSwitch v-model="ctx.finalMetadata.track.explicit" />
            <span class="text-xs text-slate-400">{{ ctx.finalMetadata.track.explicit ? 'Yes' : 'No' }}</span>
          </div>
        </div>
      </div>

      <div>
        <div class="mb-1 flex items-center justify-between">
          <label class="text-xs font-medium text-slate-400">Lyrics</label>
          <div class="flex items-center gap-2">
            <Badge
              v-if="ctx.findSuggestion('lyrics_type')"
              :value="ctx.findSuggestion('lyrics_type')!.value === 'lrc' ? 'Synced' : 'Plain'"
              :severity="ctx.findSuggestion('lyrics_type')!.value === 'lrc' ? 'success' : 'info'"
              size="small"
            />
            <Badge
              v-if="ctx.findSuggestion('lyrics')"
              :value="ctx.sourceLabel(ctx.findSuggestion('lyrics')!.source)"
              :severity="ctx.sourceSeverity(ctx.findSuggestion('lyrics')!.source)"
              size="small"
            />
            <Button
              v-if="ctx.trackAudioUrl.value && ctx.lyricsType.value === 'lrc' && ctx.finalMetadata.track.lyrics"
              :icon="ctx.showSyncedLyrics.value ? 'pi pi-eye-slash' : 'pi pi-eye'"
              :label="ctx.showSyncedLyrics.value ? 'Hide' : 'Karaoke'"
              severity="help"
              size="small"
              @click="ctx.showSyncedLyrics.value = !ctx.showSyncedLyrics.value"
            />
            <Button
              icon="pi pi-arrows-alt"
              severity="secondary"
              text
              size="small"
              :pt="{ root: { class: 'h-6 w-6' } }"
              @click="ctx.lyricsExpanded.value = !ctx.lyricsExpanded.value"
            />
          </div>
        </div>

        <!-- Live synced lyrics display -->
        <div
          v-if="ctx.showSyncedLyrics.value && ctx.trackAudioUrl.value && parsedLyrics.length > 0"
          ref="syncedLyricsContainer"
          class="mb-3 max-h-64 overflow-y-auto rounded-lg bg-black/60 p-4 scrollbar-none"
        >
          <div
            v-for="(line, idx) in parsedLyrics"
            :key="idx"
            :ref="(el: any) => { if (el) syncedLyricLineRefs[idx] = el as HTMLElement }"
            class="cursor-pointer px-2 py-1.5 text-center text-sm leading-relaxed transition-all duration-300"
            :class="{
              'scale-105 font-bold text-white': idx === activeLyricLine,
              'text-white/15': activeLyricLine >= 0 && idx < activeLyricLine,
              'text-white/30': idx > activeLyricLine,
            }"
            @click="ctx.seekAudio(line.timeSeconds)"
          >
            {{ line.text }}
          </div>
        </div>

        <Textarea
          v-model="ctx.finalMetadata.track.lyrics"
          dir="auto"
          :auto-resize="true"
          class="w-full font-mono text-xs leading-relaxed"
          :class="{ 'min-h-64': ctx.lyricsExpanded.value, 'min-h-24': !ctx.lyricsExpanded.value }"
          placeholder="Lyrics..."
        />
      </div>

      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">Cover URL</label>
          <InputText v-model="ctx.finalMetadata.track.coverUrl" class="w-full font-mono text-xs" placeholder="https://..." />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">Spotify Preview URL</label>
          <InputText
            v-model="ctx.finalMetadata.track.spotifyPreviewUrl"
            class="w-full font-mono text-xs text-slate-400"
            :disabled="!!ctx.enrichment.value?.spotify?.previewUrl"
            :placeholder="ctx.enrichment.value?.spotify?.previewUrl || 'Not available'"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Bot, Info, Music2, Star, User } from 'lucide-vue-next'
import { ref, computed, watch, onUnmounted } from 'vue'
import { useIngestionReviewContext } from '@/composables/useIngestionReviewContext'
import AdminAudioPreview from '@/components/admin/AdminAudioPreview.vue'
import { parseLRCLines, parsePlainLines } from '@/composables/lyrics'
import type { ParsedLine } from '@/composables/lyrics'

const ctx = useIngestionReviewContext()

let parsedLyrics: ParsedLine[] = []
const syncedLyricsContainer = ref<HTMLElement | null>(null)
const syncedLyricLineRefs = ref<HTMLElement[]>([])
let lyricsScrollTimer: ReturnType<typeof setTimeout> | null = null

function rebuildParsedLyrics() {
  const text = ctx.finalMetadata.track.lyrics
  if (!text) {
    parsedLyrics = []
    return
  }
  parsedLyrics = ctx.lyricsType.value === 'lrc' ? parseLRCLines(text) : parsePlainLines(text)
}

watch(
  () => ctx.finalMetadata.track.lyrics,
  () => { rebuildParsedLyrics() },
  { immediate: true },
)

const activeLyricLine = computed(() => {
  const t = ctx.audioCurrentTime.value
  for (let i = parsedLyrics.length - 1; i >= 0; i--) {
    if (t >= parsedLyrics[i]!.timeSeconds) return i
  }
  return -1
})

watch(activeLyricLine, (idx) => {
  if (lyricsScrollTimer) clearTimeout(lyricsScrollTimer)
  lyricsScrollTimer = setTimeout(() => {
    if (idx < 0 || !syncedLyricsContainer.value) return
    const target = syncedLyricLineRefs.value[idx]
    target?.scrollIntoView({ block: 'center', behavior: 'smooth' })
  }, 80)
})

onUnmounted(() => {
  if (lyricsScrollTimer) clearTimeout(lyricsScrollTimer)
})
</script>
