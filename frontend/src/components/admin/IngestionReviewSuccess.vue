<template>
  <div class="py-12 text-center">
    <div class="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-green-500/20">
      <CheckCircle aria-hidden="true" class="text-3xl text-green-400"  />
    </div>
    <h2 class="mb-2 text-2xl font-bold text-white">Published to Catalog</h2>
    <p class="mb-8 text-slate-400">{{ filename }} has been published successfully.</p>

    <div class="mx-auto mb-8 grid max-w-md gap-4">
      <a
        v-if="publishResult.artistId"
        :href="`/admin/catalog/artists/${publishResult.artistId}`"
        class="flex items-center gap-3 rounded-lg border border-white/6 bg-white/3 p-4 text-left transition-colors hover:bg-white/6"
      >
        <img
          v-if="finalMetadata.artists[0]?.imageUrl"
          :src="finalMetadata.artists[0].imageUrl"
          alt=""
          class="h-10 w-10 shrink-0 rounded-full object-cover"
          @error="($event.target as HTMLImageElement).style.display='none'"
        />
        <User aria-hidden="true" v-else class="text-xl text-emerald-400"  />
        <div>
          <p class="text-sm font-medium text-white">View Artist in Catalog</p>
          <p class="text-xs text-slate-500">{{ finalMetadata.artists[0]?.name || 'Artist' }}</p>
        </div>
      </a>

      <a
        v-if="publishResult.albumId"
        :href="`/admin/catalog/albums/${publishResult.albumId}`"
        class="flex items-center gap-3 rounded-lg border border-white/6 bg-white/3 p-4 text-left transition-colors hover:bg-white/6"
      >
        <img
          v-if="finalMetadata.album.coverUrl"
          :src="finalMetadata.album.coverUrl"
          alt=""
          class="h-10 w-10 shrink-0 rounded object-cover"
          @error="($event.target as HTMLImageElement).style.display='none'"
        />
        <Book aria-hidden="true" v-else class="text-xl text-emerald-400"  />
        <div>
          <p class="text-sm font-medium text-white">View Album in Catalog</p>
          <p class="text-xs text-slate-500">{{ finalMetadata.album.title }}</p>
        </div>
      </a>

      <a
        v-if="publishResult.trackId"
        :href="`/admin/catalog/tracks/${publishResult.trackId}`"
        class="flex items-center gap-3 rounded-lg border border-white/6 bg-white/3 p-4 text-left transition-colors hover:bg-white/6"
      >
        <img
          v-if="finalMetadata.track.coverUrl"
          :src="finalMetadata.track.coverUrl"
          alt=""
          class="h-10 w-10 shrink-0 rounded object-cover"
          @error="($event.target as HTMLImageElement).style.display='none'"
        />
        <Music aria-hidden="true" v-else class="text-xl text-emerald-400"  />
        <div>
          <p class="text-sm font-medium text-white">View Track in Catalog</p>
          <p class="text-xs text-slate-500">{{ finalMetadata.track.title }}</p>
        </div>
      </a>

      <a
        :href="publishResult.audioUrl"
        target="_blank"
        class="flex items-center gap-3 rounded-lg border border-white/6 bg-white/3 p-4 text-left transition-colors hover:bg-white/6"
      >
        <ExternalLink aria-hidden="true" class="text-xl text-emerald-400"  />
        <div>
          <p class="text-sm font-medium text-white">Audio File URL</p>
          <p class="truncate text-xs text-slate-500">{{ publishResult.audioUrl }}</p>
        </div>
      </a>
    </div>

    <div class="flex justify-center gap-3">
      <Button label="Back to Drafts" icon="pi pi-arrow-left" severity="secondary" @click="goBack" />
      <Button label="Upload Another" icon="pi pi-upload" @click="uploadAnother" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Book, CheckCircle, ExternalLink, Music, User } from 'lucide-vue-next'
import { useRouter } from 'vue-router'
import type { FinalizeResult, SaveFinalMetadataRequest } from '@/services/api/ingestion/types'

const props = defineProps<{
  filename: string
  publishResult: FinalizeResult
  finalMetadata: SaveFinalMetadataRequest
}>()

const router = useRouter()

const emit = defineEmits<{
  back: []
}>()

function goBack() {
  emit('back')
}

function uploadAnother() {
  router.push({ name: 'admin.ingestion' })
}
</script>
