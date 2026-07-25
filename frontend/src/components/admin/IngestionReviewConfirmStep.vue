<template>
  <div>
    <div class="mb-6">
      <h2 class="text-lg font-semibold text-white">Confirm & Publish</h2>
      <p class="mt-1 text-sm text-slate-400">Review your selections before publishing to the catalog.</p>
    </div>

    <div class="mb-6 grid gap-6 md:grid-cols-3">
      <!-- Artists -->
      <div class="rounded-lg border border-white/6 bg-white/3 p-4">
        <h3 class="mb-3 flex items-center gap-2 text-sm font-semibold text-white">
          <Users aria-hidden="true" class="text-emerald-400"  />Artists
        </h3>
        <div class="space-y-2">
          <div
            v-for="(a, i) in finalMetadata.artists"
            :key="i"
            class="flex items-center gap-2 rounded-lg bg-white/4 p-2"
          >
            <div
              v-if="a.imageUrl"
              class="h-8 w-8 shrink-0 overflow-hidden rounded-full"
            >
              <img
                :src="a.imageUrl"
                :alt="a.name"
                class="h-full w-full object-cover"
                @error="($event.target as HTMLImageElement).style.display='none'"
              />
            </div>
            <div v-else class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-white/6 text-xs text-slate-400">
              <User aria-hidden="true" class=""  />
            </div>
            <div class="min-w-0">
              <p class="truncate text-xs font-medium text-white">{{ a.name || '(empty)' }}</p>
              <div class="flex items-center gap-1.5">
                <span
                  class="rounded-full px-1.5 py-0.5 text-[9px] font-bold uppercase tracking-wider"
                  :class="i === 0
                    ? 'bg-emerald-500/15 text-emerald-400'
                    : 'bg-purple-500/15 text-purple-400'"
                >{{ i === 0 ? 'Primary' : 'Featured' }}</span>
                <span class="text-[10px] text-slate-500">
                  {{ a.action === 'link' ? 'Existing' : 'New' }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Album -->
      <div class="rounded-lg border border-white/6 bg-white/3 p-4">
        <h3 class="mb-3 flex items-center gap-2 text-sm font-semibold text-white">
          <Book aria-hidden="true" class="text-emerald-400"  />Album
        </h3>
        <img
          v-if="finalMetadata.album.coverUrl"
          :src="finalMetadata.album.coverUrl"
          alt=""
          class="mb-2 h-24 w-24 rounded-lg object-cover shadow-xs"
          @error="($event.target as HTMLImageElement).style.display='none'"
        />
        <div class="space-y-1 text-xs">
          <p><span class="text-slate-500">Action:</span> {{ finalMetadata.album.action === 'link' ? 'Link existing' : 'Create new' }}</p>
          <p><span class="text-slate-500">Title:</span> {{ finalMetadata.album.title || '(empty)' }}</p>
          <p v-if="finalMetadata.album.releaseYear"><span class="text-slate-500">Year:</span> {{ finalMetadata.album.releaseYear }}</p>
        </div>
      </div>

      <!-- Track -->
      <div class="rounded-lg border border-white/6 bg-white/3 p-4">
        <h3 class="mb-3 flex items-center gap-2 text-sm font-semibold text-white">
          <Music aria-hidden="true" class="text-emerald-400"  />Track
        </h3>
        <img
          v-if="finalMetadata.track.coverUrl"
          :src="finalMetadata.track.coverUrl"
          alt=""
          class="mb-2 h-24 w-24 rounded-lg object-cover shadow-xs"
          @error="($event.target as HTMLImageElement).style.display='none'"
        />
        <div class="space-y-1 text-xs">
          <p><span class="text-slate-500">Title:</span> {{ finalMetadata.track.title || '(empty)' }}</p>
          <p><span class="text-slate-500">Duration:</span> {{ formatDuration(finalMetadata.track.durationSeconds) }}</p>
          <p><span class="text-slate-500">Explicit:</span> {{ finalMetadata.track.explicit ? 'Yes' : 'No' }}</p>
        </div>
      </div>
    </div>

    <Message v-if="publishingError" severity="error" :closable="false" class="mb-4">
      {{ publishingError }}
    </Message>
  </div>
</template>

<script setup lang="ts">
import { Book, Music, User, Users } from 'lucide-vue-next'
import type { SaveFinalMetadataRequest } from '@/services/api/ingestion/types'
import { formatDuration } from '@/utils/format'

defineProps<{
  finalMetadata: SaveFinalMetadataRequest
  publishingError: string | null
}>()
</script>
