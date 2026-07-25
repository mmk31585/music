<template>
  <div>
    <div class="mb-6">
      <h2 class="text-lg font-semibold text-white">Album</h2>
      <p class="mt-1 text-sm text-slate-400">Confirm or edit the album information.</p>
    </div>

    <div class="mb-6 flex flex-wrap gap-6">
      <div v-if="ctx.embeddedCover.value">
        <p class="mb-2 text-xs font-medium text-slate-500">Embedded Cover</p>
        <img
          :src="ctx.embeddedCover.value"
          alt="Embedded"
          class="h-32 w-32 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
          :class="ctx.finalMetadata.album.coverUrl === ctx.embeddedCover.value ? 'ring-green-500' : 'ring-transparent hover:ring-white/20'"
          @click="ctx.finalMetadata.album.coverUrl = ctx.embeddedCover.value!"
        />
      </div>
      <div v-if="ctx.suggestedAlbumCover.value">
<p class="mb-2 text-xs font-medium text-slate-500 flex items-center gap-1">
          <Music2 aria-hidden="true" class="text-green-400"></Music2>Spotify
        </p>
        <img
          :src="ctx.suggestedAlbumCover.value"
          alt="Spotify"
          class="h-32 w-32 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
          :class="ctx.finalMetadata.album.coverUrl === ctx.suggestedAlbumCover.value ? 'ring-green-500' : 'ring-transparent hover:ring-white/20'"
          @click="ctx.finalMetadata.album.coverUrl = ctx.suggestedAlbumCover.value!"
        />
      </div>
      <div v-if="ctx.suggestedLastfmAlbumCover.value" class="shrink-0">
<p class="mb-2 text-xs font-medium text-slate-500 flex items-center gap-1">
            <Star aria-hidden="true" class="text-yellow-400"></Star>Last.fm
          </p>
        <img
          :src="ctx.suggestedLastfmAlbumCover.value"
          alt="Last.fm"
          class="h-32 w-32 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
          :class="ctx.finalMetadata.album.coverUrl === ctx.suggestedLastfmAlbumCover.value ? 'ring-green-500' : 'ring-transparent hover:ring-white/20'"
          @click="ctx.finalMetadata.album.coverUrl = ctx.suggestedLastfmAlbumCover.value!"
        />
      </div>
      <div v-if="ctx.enrichment.value?.ml?.albumCoverUrl" class="shrink-0">
        <p class="mb-2 text-xs font-medium text-slate-500 flex items-center gap-1">
          <Bot aria-hidden="true" class="text-purple-400"></Bot>ML Server
        </p>
        <img
          :src="ctx.enrichment.value.ml.albumCoverUrl"
          alt="ML"
          class="h-32 w-32 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
          :class="ctx.finalMetadata.album.coverUrl === ctx.enrichment.value.ml.albumCoverUrl ? 'ring-green-500' : 'ring-transparent hover:ring-white/20'"
          @click="ctx.finalMetadata.album.coverUrl = ctx.enrichment.value!.ml!.albumCoverUrl"
        />
      </div>
      <div v-if="ctx.finalMetadata.album.coverUrl && ctx.finalMetadata.album.coverUrl !== ctx.embeddedCover.value && ctx.finalMetadata.album.coverUrl !== ctx.suggestedAlbumCover.value && ctx.finalMetadata.album.coverUrl !== ctx.suggestedLastfmAlbumCover.value && ctx.finalMetadata.album.coverUrl !== ctx.enrichment.value?.ml?.albumCoverUrl" class="shrink-0">
        <p class="mb-2 text-xs font-medium text-slate-500">Selected</p>
        <img
          :src="ctx.finalMetadata.album.coverUrl"
          alt="Selected"
          class="h-32 w-32 rounded-lg object-cover shadow-md ring-2 ring-green-500"
          @error="($event.target as HTMLImageElement).style.display='none'"
        />
      </div>
      <div v-if="!ctx.finalMetadata.album.coverUrl" class="flex shrink-0 flex-col items-center justify-center gap-2">
        <p class="mb-1 text-xs font-medium text-slate-500">Upload Cover</p>
        <label class="flex h-24 w-24 cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-white/6 bg-white/3 text-slate-400 transition-colors hover:border-emerald-500/50 hover:text-emerald-400">
          <Upload aria-hidden="true" class="text-xl"></Upload>
          <span class="mt-1 text-[10px]">Upload</span>
          <input type="file" accept="image/jpeg,image/png,image/webp" class="hidden" @change="ctx.uploadAlbumCover" />
        </label>
      </div>
    </div>

    <!-- Fetch album data -->
    <div class="mb-6">
      <Button
        :loading="ctx.enrichingAlbum.value"
        :disabled="ctx.enriching.value"
        label="Fetch Album Data"
        icon="pi pi-cloud-download"
        severity="help"
        size="small"
        outlined
        @click="ctx.fetchSection('album')"
      />
      <p class="mt-1.5 text-[11px] text-slate-500">
        Load album data from the track's embedded metadata (ID3 tags).
      </p>
    </div>

    <div class="mb-6">
      <p class="mb-3 text-sm font-medium text-slate-400">Album Action</p>
      <div class="flex flex-wrap gap-3">
        <Button
          :severity="ctx.finalMetadata.album.action === 'create' ? 'primary' : 'secondary'"
          :outlined="ctx.finalMetadata.album.action !== 'create'"
          size="small"
          @click="ctx.finalMetadata.album.action = 'create'; ctx.finalMetadata.album.existingId = undefined"
        >
          <Plus aria-hidden="true" class="mr-1"></Plus>Create New
        </Button>
        <Button
          :severity="ctx.finalMetadata.album.action === 'link' ? 'primary' : 'secondary'"
          :outlined="ctx.finalMetadata.album.action !== 'link'"
          size="small"
          @click="ctx.finalMetadata.album.action = 'link'"
        >
          <Link aria-hidden="true" class="mr-1"></Link>Use Existing
        </Button>
      </div>
    </div>

    <div v-if="ctx.finalMetadata.album.action === 'link'" class="mb-6">
      <label class="mb-2 block text-xs font-medium text-slate-400">Search Existing Albums</label>
      <div class="relative">
        <IconField>
          <InputIcon><Search aria-hidden="true" class=""></Search></InputIcon>
          <InputText
            :model-value="ctx.albumSearchQuery.value"
            @update:model-value="(v: string | undefined) => { if (v !== undefined) { ctx.albumSearchQuery.value = v; ctx.debouncedAlbumSearch() } }"
            placeholder="Type album title..."
            class="w-full"
          />
        </IconField>
      </div>
      <div v-if="ctx.albumSearchResults.value.length > 0" class="mt-2 max-h-48 overflow-y-auto rounded border border-white/6 bg-white/3">
        <div
          v-for="a in ctx.albumSearchResults.value"
          :key="a.id"
          class="flex cursor-pointer items-center gap-3 px-3 py-2 text-sm transition-colors hover:bg-white/6"
          :class="{ 'bg-emerald-500/20': ctx.selectedAlbumId.value === a.id }"
          @click="ctx.selectExistingAlbum(a)"
        >
          <img v-if="a.coverUrl" :src="a.coverUrl" class="h-8 w-8 rounded object-cover" />
          <div v-else class="flex h-8 w-8 items-center justify-center rounded bg-white/6 text-xs">
            <Image aria-hidden="true" class=""></Image>
          </div>
          <div>
            <p class="font-medium text-white">{{ a.title }}</p>
            <p v-if="a.artistName" class="text-[11px] text-slate-500">{{ a.artistName }}</p>
          </div>
        </div>
      </div>
      <p v-else-if="ctx.albumSearchQuery.value && !ctx.albumSearching.value" class="mt-2 text-xs text-slate-500">
        No existing albums found.
      </p>
    </div>

    <div class="space-y-4 rounded-lg border border-white/6 bg-white/3 p-5">
      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">
            Title
            <template v-if="ctx.findSuggestion('album')">
              <Badge :value="ctx.sourceLabel(ctx.findSuggestion('album')!.source)" :severity="ctx.sourceSeverity(ctx.findSuggestion('album')!.source)" size="small" />
              <Badge :value="ctx.findSuggestion('album')!.confidence" :severity="ctx.confidenceSeverity(ctx.findSuggestion('album')!.confidence)" size="small" />
            </template>
          </label>
          <InputText v-model="ctx.finalMetadata.album.title" dir="auto" class="w-full" placeholder="Album title" />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">
            Release Year
            <template v-if="ctx.findSuggestion('year')">
              <Badge :value="ctx.sourceLabel(ctx.findSuggestion('year')!.source)" :severity="ctx.sourceSeverity(ctx.findSuggestion('year')!.source)" size="small" />
              <Badge :value="ctx.findSuggestion('year')!.confidence" :severity="ctx.confidenceSeverity(ctx.findSuggestion('year')!.confidence)" size="small" />
            </template>
          </label>
          <InputNumber v-model="ctx.finalMetadata.album.releaseYear" :use-grouping="false" class="w-full" placeholder="2024" />
        </div>
      </div>
      <div class="grid gap-4 md:grid-cols-2">
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">Genre</label>
          <InputText v-model="ctx.finalMetadata.album.genre" dir="auto" class="w-full" placeholder="Genre" />
        </div>
        <div>
          <label class="mb-1 block text-xs font-medium text-slate-400">Cover URL</label>
          <InputText v-model="ctx.finalMetadata.album.coverUrl" class="w-full font-mono text-xs" placeholder="https://..." />
        </div>
      </div>
      <div v-if="ctx.enrichment.value?.musicbrainz?.albumMbid">
        <label class="mb-1 block text-xs font-medium text-slate-400">MusicBrainz Release ID</label>
        <div class="flex items-center gap-2">
          <InputText :model-value="ctx.enrichment.value.musicbrainz.albumMbid" class="flex-1 font-mono text-xs text-slate-400" disabled />
          <Button
            v-if="ctx.finalMetadata.album.musicbrainzReleaseId !== ctx.enrichment.value.musicbrainz.albumMbid"
            label="Apply"
            icon="pi pi-check"
            size="small"
            severity="success"
            outlined
            @click="ctx.finalMetadata.album.musicbrainzReleaseId = ctx.enrichment.value!.musicbrainz!.albumMbid; ctx.hasUnsavedChanges.value = true"
          />
          <CheckCircle aria-hidden="true" v-else class="text-green-400" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Bot, CheckCircle, Image, Link, Music2, Plus, Search, Star, Upload } from 'lucide-vue-next'
import { useIngestionReviewContext } from '@/composables/useIngestionReviewContext'

const ctx = useIngestionReviewContext()
</script>
