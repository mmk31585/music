<template>
  <div>
    <div class="mb-6">
      <h2 class="text-lg font-semibold text-white">Artists</h2>
      <p class="mt-1 text-sm text-slate-400">
        Add all artists credited on this track.
        <span class="text-slate-500">The first artist is the primary; additional artists are featured.</span>
      </p>
    </div>

    <!-- Fetch artist data -->
    <div class="mb-6">
      <Button
        :loading="ctx.enrichingArtist.value"
        :disabled="ctx.enriching.value"
        label="Fetch Artist Data"
        icon="pi pi-cloud-download"
        severity="help"
        size="small"
        outlined
        @click="ctx.fetchSection('artist')"
      />
      <p class="mt-1.5 text-[11px] text-slate-500">
        Load artist data from the track's embedded metadata (ID3 tags).
      </p>
    </div>

    <!-- Artist cards loop -->
    <div class="space-y-4">
      <div
        v-for="(artist, aIdx) in ctx.finalMetadata.artists"
        :key="aIdx"
        class="overflow-hidden rounded-xl border transition-all duration-200"
        :class="aIdx === 0
          ? 'border-emerald-500/30 bg-emerald-500/4'
          : 'border-white/6 bg-white/3'"
      >
        <!-- Card header -->
        <div class="flex items-center justify-between px-5 py-3">
          <div class="flex items-center gap-3">
            <span
              class="flex h-7 w-7 items-center justify-center rounded-full text-[11px] font-bold"
              :class="aIdx === 0
                ? 'bg-emerald-500/20 text-emerald-400'
                : 'bg-white/6 text-slate-400'"
            >{{ aIdx + 1 }}</span>
            <div>
              <span class="text-sm font-semibold text-white">{{ artist.name || 'Unnamed artist' }}</span>
              <div class="flex items-center gap-2 mt-0.5">
                <span
                  class="rounded-full px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider"
                  :class="aIdx === 0
                    ? 'bg-emerald-500/15 text-emerald-400'
                    : 'bg-purple-500/15 text-purple-400'"
                >
                  {{ aIdx === 0 ? 'Primary' : 'Featured' }}
                </span>
                <span
                  v-if="artist.action === 'link'"
                  class="rounded-full bg-sky-500/15 px-2 py-0.5 text-[10px] font-medium text-sky-400"
                >Existing</span>
                <span
                  v-else
                  class="rounded-full bg-amber-500/15 px-2 py-0.5 text-[10px] font-medium text-amber-400"
                >New</span>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-1">
            <button
              v-if="aIdx > 0"
              type="button"
              :disabled="aIdx <= 1"
              class="flex h-7 w-7 items-center justify-center rounded-lg text-slate-500 transition hover:bg-white/6 hover:text-white disabled:opacity-30"
              @click="ctx.moveFeaturedArtist(aIdx, -1)"
              title="Move up"
            >
              <ChevronUp aria-hidden="true" class="text-xs"></ChevronUp>
            </button>
            <button
              v-if="aIdx > 0 && aIdx < ctx.finalMetadata.artists.length - 1"
              type="button"
              class="flex h-7 w-7 items-center justify-center rounded-lg text-slate-500 transition hover:bg-white/6 hover:text-white"
              @click="ctx.moveFeaturedArtist(aIdx, 1)"
              title="Move down"
            >
              <ChevronDown aria-hidden="true" class="text-xs"></ChevronDown>
            </button>
            <button
              v-if="aIdx > 0"
              type="button"
              class="flex h-7 w-7 items-center justify-center rounded-lg text-red-400/50 transition hover:bg-red-500/10 hover:text-red-400"
              @click="ctx.removeFeaturedArtist(aIdx)"
              title="Remove artist"
            >
              <Trash2 aria-hidden="true" class="text-xs"></Trash2>
            </button>
          </div>
        </div>

        <!-- Card body -->
        <div class="border-t border-white/4 px-5 py-4">
          <!-- Artist image picker -->
          <div class="mb-4">
            <p class="mb-2 text-xs font-medium text-slate-500">
              {{ aIdx === 0 ? 'Primary Artist Image' : 'Artist Image' }}
            </p>
            <div class="flex flex-wrap gap-3">
              <div v-if="ctx.enrichment.value?.spotify?.artistImageUrl" class="flex flex-col items-center gap-1">
                <p class="text-[10px] font-medium text-slate-500 flex items-center gap-1">
                  <Music2 aria-hidden="true" class="text-green-400"></Music2>
                </p>
                <img
                  :src="ctx.enrichment.value.spotify.artistImageUrl"
                  alt="Spotify"
                  class="h-16 w-16 cursor-pointer rounded-lg object-cover shadow-xs ring-2 transition-all hover:opacity-80"
                  :class="artist.imageUrl === ctx.enrichment.value.spotify.artistImageUrl ? 'ring-green-500' : 'ring-transparent'"
                  @click="artist.imageUrl = ctx.enrichment.value!.spotify!.artistImageUrl"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
              </div>
              <div v-if="ctx.enrichment.value?.lastfm?.artistImageUrl" class="flex flex-col items-center gap-1">
                <p class="text-[10px] font-medium text-slate-500 flex items-center gap-1">
                  <Star aria-hidden="true" class="text-yellow-400"></Star>
                </p>
                <img
                  :src="ctx.enrichment.value.lastfm.artistImageUrl"
                  alt="Last.fm"
                  class="h-16 w-16 cursor-pointer rounded-lg object-cover shadow-xs ring-2 transition-all hover:opacity-80"
                  :class="artist.imageUrl === ctx.enrichment.value.lastfm.artistImageUrl ? 'ring-green-500' : 'ring-transparent'"
                  @click="artist.imageUrl = ctx.enrichment.value!.lastfm!.artistImageUrl"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
              </div>
              <div v-if="ctx.enrichment.value?.ml?.artistImageUrl" class="flex flex-col items-center gap-1">
                <p class="text-[10px] font-medium text-slate-500 flex items-center gap-1">
                  <Bot aria-hidden="true" class="text-purple-400"></Bot>
                </p>
                <img
                  :src="ctx.enrichment.value.ml.artistImageUrl"
                  alt="ML"
                  class="h-16 w-16 cursor-pointer rounded-lg object-cover shadow-xs ring-2 transition-all hover:opacity-80"
                  :class="artist.imageUrl === ctx.enrichment.value.ml.artistImageUrl ? 'ring-green-500' : 'ring-transparent'"
                  @click="artist.imageUrl = ctx.enrichment.value!.ml!.artistImageUrl"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
              </div>
              <div v-if="artist.imageUrl && ![ctx.enrichment.value?.spotify?.artistImageUrl, ctx.enrichment.value?.lastfm?.artistImageUrl, ctx.enrichment.value?.ml?.artistImageUrl].includes(artist.imageUrl)" class="flex flex-col items-center gap-1">
                <p class="text-[10px] font-medium text-green-400 flex items-center gap-1">
                  <CheckCircle aria-hidden="true" class=""></CheckCircle>
                </p>
                <img
                  :src="artist.imageUrl"
                  alt="Selected"
                  class="h-16 w-16 rounded-lg object-cover shadow-xs ring-2 ring-green-500"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
              </div>
              <label class="flex h-16 w-16 cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-white/6 bg-white/3 text-slate-500 transition-colors hover:border-emerald-500/50 hover:text-emerald-400">
                <Upload aria-hidden="true" class="text-sm"></Upload>
                <span class="mt-0.5 text-[9px]">Upload</span>
                <input type="file" accept="image/jpeg,image/png,image/webp" class="hidden" @change="ctx.uploadArtistImage($event, aIdx)" />
              </label>
            </div>
          </div>

          <!-- Action (create vs link) -->
          <div class="mb-4">
            <p class="mb-2 text-xs font-medium text-slate-500">Action</p>
            <div class="flex flex-wrap gap-2">
              <Button
                :severity="artist.action === 'create' ? 'primary' : 'secondary'"
                :outlined="artist.action !== 'create'"
                size="small"
                @click="artist.action = 'create'; artist.existingId = undefined"
              >
                <Plus aria-hidden="true" class="mr-1"></Plus>Create New
              </Button>
              <Button
                :severity="artist.action === 'link' ? 'primary' : 'secondary'"
                :outlined="artist.action !== 'link'"
                size="small"
                @click="artist.action = 'link'"
              >
                <Link aria-hidden="true" class="mr-1"></Link>Use Existing
              </Button>
            </div>
          </div>

          <!-- Search existing (when linking) -->
          <div v-if="artist.action === 'link'" class="mb-4">
            <label class="mb-1.5 block text-xs font-medium text-slate-500">Search Existing Artists</label>
            <div class="relative">
              <IconField>
                <InputIcon><Search aria-hidden="true" class=""></Search></InputIcon>
                <InputText
                  :model-value="ctx.artistSearchQuery.value"
                  @update:model-value="(v: string | undefined) => { if (v !== undefined) { ctx.artistSearchQuery.value = v; ctx.debouncedArtistSearch() } }"
                  placeholder="Type artist name..."
                  class="w-full"
                />
              </IconField>
            </div>
            <div v-if="ctx.artistSearchResults.value.length > 0" class="mt-1 max-h-36 overflow-y-auto rounded border border-white/6 bg-white/3">
              <div
                v-for="a in ctx.artistSearchResults.value"
                :key="a.id"
                class="flex cursor-pointer items-center gap-3 px-3 py-2 text-sm transition-colors hover:bg-white/6"
                :class="{ 'bg-emerald-500/20': ctx.selectedArtistId.value === a.id }"
                @click="ctx.selectExistingArtist(a, aIdx)"
              >
                <img v-if="a.imageUrl" :src="a.imageUrl" class="h-7 w-7 rounded-full object-cover" />
                <div v-else class="flex h-7 w-7 items-center justify-center rounded-full bg-white/6 text-xs">
                  <User aria-hidden="true" class=""></User>
                </div>
                <div>
                  <p class="font-medium text-white">{{ a.name }}</p>
                  <p v-if="a.country" class="text-[11px] text-slate-500">{{ a.country }}</p>
                </div>
              </div>
            </div>
            <p v-else-if="ctx.artistSearchQuery.value && !ctx.artistSearching.value" class="mt-1 text-xs text-slate-500">No existing artists found.</p>
          </div>

          <!-- Fields grid -->
          <div class="grid gap-3 md:grid-cols-2">
            <div>
              <label class="mb-1 block text-xs font-medium text-slate-500">
                Name
                <template v-if="aIdx === 0 && ctx.findSuggestion('artist')">
                  <Badge :value="ctx.sourceLabel(ctx.findSuggestion('artist')!.source)" :severity="ctx.sourceSeverity(ctx.findSuggestion('artist')!.source)" size="small" />
                </template>
              </label>
              <InputText
                :model-value="artist.name"
                @update:model-value="(v: string | undefined) => { if (v !== undefined) { artist.name = v; ctx.hasUnsavedChanges.value = true } }"
                dir="auto"
                class="w-full"
                :placeholder="aIdx === 0 ? 'Primary artist name' : 'Featured artist name'"
              />
            </div>
            <div>
              <label class="mb-1 block text-xs font-medium text-slate-500">Country</label>
              <InputText
                :model-value="artist.country"
                @update:model-value="(v: string | undefined) => { if (v !== undefined) { artist.country = v; ctx.hasUnsavedChanges.value = true } }"
                maxlength="2"
                class="w-20 uppercase"
                placeholder="IR"
              />
            </div>
          </div>
          <div class="mt-3">
            <label class="mb-1 block text-xs font-medium text-slate-500">Bio</label>
            <Textarea
              :model-value="artist.bio"
              @update:model-value="(v: string | undefined) => { if (v !== undefined) { artist.bio = v; ctx.hasUnsavedChanges.value = true } }"
              dir="auto"
              :auto-resize="true"
              class="w-full"
              rows="2"
              placeholder="Artist biography"
            />
          </div>
          <div class="mt-3">
            <label class="mb-1 block text-xs font-medium text-slate-500">Image URL</label>
            <InputText
              :model-value="artist.imageUrl"
              @update:model-value="(v: string | undefined) => { if (v !== undefined) { artist.imageUrl = v; ctx.hasUnsavedChanges.value = true } }"
              class="w-full font-mono text-xs"
              placeholder="https://..."
            />
          </div>
          <div v-if="aIdx === 0 && ctx.enrichment.value?.musicbrainz?.artistMbid" class="mt-3">
            <label class="mb-1 block text-xs font-medium text-slate-500">MusicBrainz MBID</label>
            <div class="flex items-center gap-2">
              <InputText
                :model-value="ctx.enrichment.value.musicbrainz.artistMbid"
                class="flex-1 font-mono text-xs text-slate-400"
                disabled
              />
              <Button
                v-if="artist.musicbrainzMbid !== ctx.enrichment.value.musicbrainz.artistMbid"
                label="Apply"
                icon="pi pi-check"
                size="small"
                severity="success"
                outlined
                @click="artist.musicbrainzMbid = ctx.enrichment.value!.musicbrainz!.artistMbid; ctx.hasUnsavedChanges.value = true"
              />
              <CheckCircle aria-hidden="true" v-else class="text-green-400" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Featured Artist button -->
    <div class="mt-4">
      <Button
        label="Add Featured Artist"
        icon="pi pi-plus"
        severity="secondary"
        size="small"
        outlined
        @click="ctx.addFeaturedArtist"
      />
      <p class="mt-1.5 text-[11px] text-slate-500">
        Add featured artists for collaborations (feat., ft., &amp;).
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Bot, CheckCircle, ChevronDown, ChevronUp, Link, Music2, Plus, Search, Star, Trash2, Upload, User } from 'lucide-vue-next'
import { useIngestionReviewContext } from '@/composables/useIngestionReviewContext'

const ctx = useIngestionReviewContext()
</script>
