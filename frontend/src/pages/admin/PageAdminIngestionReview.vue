<!-- TODO MEDIUM: Uses surface-* / primary-* PrimeVue classes. Should use slate-* / white/* admin dark theme. -->
<template>
  <div class="mx-auto w-full max-w-5xl px-4 py-6 md:px-6 lg:px-8">
    <AdminSectionHeader
      eyebrow="Content"
      title="Review Ingestion Draft"
      :description="`${draftDetail?.originalFilename || ''} — ${formatDuration(draftDetail?.durationSeconds)}`"
    >
      <template #actions>
        <Button
          v-if="!enriching && !enrichment?.enrichment_attempted && draftDetail?.status === 'review'"
          label="Enrich Now"
          icon="pi pi-magic"
          severity="warn"
          size="small"
          @click="triggerEnrich"
        />
        <Button
          v-if="enriching"
          label="Enriching..."
          icon="pi pi-spin pi-spinner"
          severity="warn"
          size="small"
          disabled
        />
        <Button
          label="Back to List"
          icon="pi pi-arrow-left"
          severity="secondary"
          text
          @click="goBack"
        />
      </template>
    </AdminSectionHeader>

    <div v-if="loading" class="py-20 text-center">
      <i class="pi pi-spin pi-spinner text-3xl text-surface-400"></i>
      <p class="mt-4 text-surface-500">Loading draft...</p>
    </div>

    <div v-else-if="error" class="mb-8">
      <Message severity="error" :closable="false">{{ error }}</Message>
    </div>

    <template v-else-if="published && publishResult">
      <div class="py-12 text-center">
        <div class="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-green-500/20">
          <i class="pi pi-check-circle text-3xl text-green-400"></i>
        </div>
        <h2 class="mb-2 text-2xl font-bold text-white">Published to Catalog</h2>
        <p class="mb-8 text-surface-400">{{ draftDetail?.originalFilename }} has been published successfully.</p>
        <div class="mx-auto mb-8 grid max-w-md gap-4">
          <a
            v-if="publishResult.artistId"
            :href="`/admin/catalog/artists/${publishResult.artistId}`"
            class="flex items-center gap-3 rounded-lg border border-surface-700 bg-surface-800 p-4 text-left transition-colors hover:bg-surface-700"
          >
            <img
              v-if="finalMetadata.artist.imageUrl"
              :src="finalMetadata.artist.imageUrl"
              alt=""
              class="h-10 w-10 flex-shrink-0 rounded-full object-cover"
              @error="($event.target as HTMLImageElement).style.display='none'"
            />
            <i v-else class="pi pi-user text-xl text-primary"></i>
            <div>
              <p class="text-sm font-medium text-white">View Artist in Catalog</p>
              <p class="text-xs text-surface-500">{{ finalMetadata.artist.name }}</p>
            </div>
          </a>
          <a
            v-if="publishResult.albumId"
            :href="`/admin/catalog/albums/${publishResult.albumId}`"
            class="flex items-center gap-3 rounded-lg border border-surface-700 bg-surface-800 p-4 text-left transition-colors hover:bg-surface-700"
          >
            <img
              v-if="finalMetadata.album.coverUrl"
              :src="finalMetadata.album.coverUrl"
              alt=""
              class="h-10 w-10 flex-shrink-0 rounded object-cover"
              @error="($event.target as HTMLImageElement).style.display='none'"
            />
            <i v-else class="pi pi-book text-xl text-primary"></i>
            <div>
              <p class="text-sm font-medium text-white">View Album in Catalog</p>
              <p class="text-xs text-surface-500">{{ finalMetadata.album.title }}</p>
            </div>
          </a>
          <a
            v-if="publishResult.trackId"
            :href="`/admin/catalog/tracks/${publishResult.trackId}`"
            class="flex items-center gap-3 rounded-lg border border-surface-700 bg-surface-800 p-4 text-left transition-colors hover:bg-surface-700"
          >
            <img
              v-if="finalMetadata.track.coverUrl"
              :src="finalMetadata.track.coverUrl"
              alt=""
              class="h-10 w-10 flex-shrink-0 rounded object-cover"
              @error="($event.target as HTMLImageElement).style.display='none'"
            />
            <i v-else class="pi pi-music text-xl text-primary"></i>
            <div>
              <p class="text-sm font-medium text-white">View Track in Catalog</p>
              <p class="text-xs text-surface-500">{{ finalMetadata.track.title }}</p>
            </div>
          </a>
          <a
            :href="publishResult.audioUrl"
            target="_blank"
            class="flex items-center gap-3 rounded-lg border border-surface-700 bg-surface-800 p-4 text-left transition-colors hover:bg-surface-700"
          >
            <i class="pi pi-external-link text-xl text-primary"></i>
            <div>
              <p class="text-sm font-medium text-white">Audio File URL</p>
              <p class="truncate text-xs text-surface-500">{{ publishResult.audioUrl }}</p>
            </div>
          </a>
        </div>
        <div class="flex justify-center gap-3">
          <Button label="Back to Drafts" icon="pi pi-arrow-left" severity="secondary" @click="goBack" />
          <Button label="Upload Another" icon="pi pi-upload" @click="router.push({ name: 'admin.ingestion' })" />
        </div>
      </div>
    </template>

    <template v-else>
      <div class="mb-8">
        <div class="flex items-center gap-2">
          <div
            v-for="(s, i) in steps"
            :key="i"
            class="flex items-center"
          >
            <div
              class="flex cursor-pointer items-center gap-2 rounded-full px-3 py-1.5 text-xs font-medium transition-colors"
              :class="stepperClass(i)"
              @click="step = i"
            >
              <span
                class="flex h-5 w-5 items-center justify-center rounded-full text-[10px] font-bold"
                :class="stepIconClass(i)"
              >
                <i v-if="i < step" class="pi pi-check"></i>
                <span v-else>{{ i + 1 }}</span>
              </span>
              {{ s.label }}
            </div>
            <i v-if="i < steps.length - 1" class="pi pi-chevron-right mx-2 text-xs text-surface-500"></i>
          </div>
        </div>
      </div>

      <div v-if="enriching" class="mb-8">
        <Message severity="info" :closable="false">
          <i class="pi pi-spin pi-spinner mr-2"></i>Enrichment in progress — covers, bio, and suggestions will appear once complete.
        </Message>
      </div>

      <div v-if="!enrichment?.enrichment_attempted && !enriching && draftDetail?.status === 'review'" class="mb-8">
        <Message severity="warn" :closable="false">
          <i class="pi pi-exclamation-triangle mr-2"></i>No enrichment data available. Click "Enrich Now" above or fill in the fields manually.
        </Message>
      </div>

      <Transition name="fade" mode="out-in">
        <div :key="step">
          <!-- Step 1: Artist -->
          <div v-if="step === 0">
            <div class="mb-6">
              <h2 class="text-lg font-semibold text-white">Artist</h2>
              <p class="mt-1 text-sm text-surface-400">Confirm or edit the artist information.</p>
            </div>

            <div class="mb-6 grid gap-6 md:grid-cols-3">
              <div
                v-if="enrichment?.spotify?.artistImageUrl"
                class="flex-shrink-0"
              >
                <p class="mb-2 text-xs font-medium text-surface-500">
                  Suggested Image
                  <i class="pi pi-spotify ml-1 text-green-400 text-[10px]"></i>
                </p>
                <img
                  :src="enrichment.spotify.artistImageUrl"
                  alt="Artist"
                  class="h-32 w-32 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
                  :class="finalMetadata.artist.imageUrl === enrichment.spotify.artistImageUrl ? 'ring-green-500' : 'ring-purple-500 hover:ring-green-400'"
                  @click="finalMetadata.artist.imageUrl = enrichment.spotify.artistImageUrl"
                  title="Click to use this image"
                />
              </div>
              <div v-if="finalMetadata.artist.imageUrl && finalMetadata.artist.imageUrl !== enrichment?.spotify?.artistImageUrl" class="flex-shrink-0">
                <p class="mb-2 text-xs font-medium text-surface-500">Selected Image</p>
                <img
                  :src="finalMetadata.artist.imageUrl"
                  alt="Selected"
                  class="h-32 w-32 rounded-lg object-cover shadow-md ring-2 ring-green-500"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
              </div>
              <div v-if="!finalMetadata.artist.imageUrl" class="flex flex-shrink-0 flex-col items-center justify-center gap-2">
                <p class="mb-1 text-xs font-medium text-surface-500">No Image Yet</p>
                <label class="flex h-24 w-24 cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-surface-600 bg-surface-800 text-surface-400 transition-colors hover:border-primary hover:text-primary">
                  <i class="pi pi-upload text-xl"></i>
                  <span class="mt-1 text-[10px]">Upload</span>
                  <input type="file" accept="image/jpeg,image/png,image/webp" class="hidden" @change="uploadArtistImage" />
                </label>
              </div>
              <div v-if="enrichment?.lastfm?.artistBio" class="md:col-span-1">
                <p class="mb-2 text-xs font-medium text-surface-500">
                  Last.fm Bio
                  <i class="pi pi-star-fill ml-1 text-yellow-400 text-[10px]"></i>
                </p>
                <p class="max-h-24 overflow-y-auto rounded bg-surface-800 p-3 text-xs leading-relaxed text-surface-300">
                  {{ truncateBio(enrichment.lastfm.artistBio) }}
                </p>
              </div>
            </div>

            <div class="mb-6">
              <p class="mb-3 text-sm font-medium text-surface-400">Artist Action</p>
              <div class="flex flex-wrap gap-3">
                <Button
                  :severity="finalMetadata.artist.action === 'create' ? 'primary' : 'secondary'"
                  :outlined="finalMetadata.artist.action !== 'create'"
                  size="small"
                  @click="finalMetadata.artist.action = 'create'; finalMetadata.artist.existingId = undefined"
                >
                  <i class="pi pi-plus mr-1"></i>Create New
                </Button>
                <Button
                  :severity="finalMetadata.artist.action === 'link' ? 'primary' : 'secondary'"
                  :outlined="finalMetadata.artist.action !== 'link'"
                  size="small"
                  @click="finalMetadata.artist.action = 'link'"
                >
                  <i class="pi pi-link mr-1"></i>Use Existing
                </Button>
              </div>
            </div>

            <div v-if="finalMetadata.artist.action === 'link'" class="mb-6">
              <label class="mb-2 block text-xs font-medium text-surface-400">Search Existing Artists</label>
              <div class="relative">
                <IconField>
                  <InputIcon><i class="pi pi-search"></i></InputIcon>
                  <InputText
                    v-model="artistSearchQuery"
                    placeholder="Type artist name..."
                    class="w-full"
                    @update:model-value="debouncedArtistSearch"
                  />
                </IconField>
              </div>
              <div v-if="artistSearchResults.length > 0" class="mt-2 max-h-48 overflow-y-auto rounded border border-surface-700 bg-surface-800">
                <div
                  v-for="a in artistSearchResults"
                  :key="a.id"
                  class="flex cursor-pointer items-center gap-3 px-3 py-2 text-sm transition-colors hover:bg-surface-700"
                  :class="{ 'bg-primary/20': selectedArtistId === a.id }"
                  @click="selectExistingArtist(a)"
                >
                  <img
                    v-if="a.imageUrl"
                    :src="a.imageUrl"
                    class="h-8 w-8 rounded-full object-cover"
                  />
                  <div v-else class="flex h-8 w-8 items-center justify-center rounded-full bg-surface-600 text-xs">
                    <i class="pi pi-user"></i>
                  </div>
                  <div>
                    <p class="font-medium text-white">{{ a.name }}</p>
                    <p v-if="a.country" class="text-[11px] text-surface-500">{{ a.country }}</p>
                  </div>
                </div>
              </div>
              <p v-else-if="artistSearchQuery && !artistSearching" class="mt-2 text-xs text-surface-500">
                No existing artists found.
              </p>
            </div>

            <div class="space-y-4 rounded-lg border border-surface-700 bg-surface-800 p-5">
              <div class="grid gap-4 md:grid-cols-2">
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">
                    Name
                    <template v-if="findSuggestion('artist')">
                      <Badge
                        :value="sourceLabel(findSuggestion('artist')!.source)"
                        :severity="sourceSeverity(findSuggestion('artist')!.source)"
                        size="small"
                      />
                      <Badge
                        :value="findSuggestion('artist')!.confidence"
                        :severity="confidenceSeverity(findSuggestion('artist')!.confidence)"
                        size="small"
                      />
                    </template>
                  </label>
                  <InputText
                    v-model="finalMetadata.artist.name"
                    dir="auto"
                    class="w-full"
                    placeholder="Artist name"
                  />
                </div>
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">Country (ISO 3166-1 alpha-2)</label>
                  <InputText
                    v-model="finalMetadata.artist.country"
                    maxlength="2"
                    class="w-20 uppercase"
                    placeholder="IR"
                  />
                </div>
              </div>
              <div>
                <label class="mb-1 block text-xs font-medium text-surface-400">Bio</label>
                <Textarea
                  v-model="finalMetadata.artist.bio"
                  dir="auto"
                  :auto-resize="true"
                  class="w-full"
                  rows="3"
                  placeholder="Artist biography"
                />
              </div>
              <div>
                <label class="mb-1 block text-xs font-medium text-surface-400">Image URL</label>
                <InputText
                  v-model="finalMetadata.artist.imageUrl"
                  class="w-full font-mono text-xs"
                  placeholder="https://..."
                />
              </div>
              <div v-if="enrichment?.musicbrainz?.artistMbid">
                <label class="mb-1 block text-xs font-medium text-surface-400">MusicBrainz MBID</label>
                <InputText
                  :model-value="enrichment.musicbrainz.artistMbid"
                  class="w-full font-mono text-xs text-surface-400"
                  disabled
                />
              </div>
            </div>
          </div>

          <!-- Step 2: Album -->
          <div v-if="step === 1">
            <div class="mb-6">
              <h2 class="text-lg font-semibold text-white">Album</h2>
              <p class="mt-1 text-sm text-surface-400">Confirm or edit the album information.</p>
            </div>

            <div class="mb-6 flex flex-wrap gap-6">
              <div v-if="embeddedCover">
                <p class="mb-2 text-xs font-medium text-surface-500">Embedded Cover</p>
                <img
                  :src="embeddedCover"
                  alt="Embedded"
                  class="h-32 w-32 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
                  :class="finalMetadata.album.coverUrl === embeddedCover ? 'ring-green-500' : 'ring-transparent hover:ring-surface-400'"
                  @click="finalMetadata.album.coverUrl = embeddedCover"
                />
              </div>
              <div v-if="suggestedAlbumCover">
                <p class="mb-2 text-xs font-medium text-surface-500">Spotify</p>
                <img
                  :src="suggestedAlbumCover"
                  alt="Spotify"
                  class="h-32 w-32 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
                  :class="finalMetadata.album.coverUrl === suggestedAlbumCover ? 'ring-green-500' : 'ring-transparent hover:ring-surface-400'"
                  @click="finalMetadata.album.coverUrl = suggestedAlbumCover"
                />
              </div>
              <div v-if="finalMetadata.album.coverUrl && finalMetadata.album.coverUrl !== embeddedCover && finalMetadata.album.coverUrl !== suggestedAlbumCover" class="flex-shrink-0">
                <p class="mb-2 text-xs font-medium text-surface-500">Selected</p>
                <img
                  :src="finalMetadata.album.coverUrl"
                  alt="Selected"
                  class="h-32 w-32 rounded-lg object-cover shadow-md ring-2 ring-green-500"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
              </div>
              <div v-if="!finalMetadata.album.coverUrl" class="flex flex-shrink-0 flex-col items-center justify-center gap-2">
                <p class="mb-1 text-xs font-medium text-surface-500">Upload Cover</p>
                <label class="flex h-24 w-24 cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-surface-600 bg-surface-800 text-surface-400 transition-colors hover:border-primary hover:text-primary">
                  <i class="pi pi-upload text-xl"></i>
                  <span class="mt-1 text-[10px]">Upload</span>
                  <input type="file" accept="image/jpeg,image/png,image/webp" class="hidden" @change="uploadAlbumCover" />
                </label>
              </div>
            </div>

            <div class="mb-6">
              <p class="mb-3 text-sm font-medium text-surface-400">Album Action</p>
              <div class="flex flex-wrap gap-3">
                <Button
                  :severity="finalMetadata.album.action === 'create' ? 'primary' : 'secondary'"
                  :outlined="finalMetadata.album.action !== 'create'"
                  size="small"
                  @click="finalMetadata.album.action = 'create'; finalMetadata.album.existingId = undefined"
                >
                  <i class="pi pi-plus mr-1"></i>Create New
                </Button>
                <Button
                  :severity="finalMetadata.album.action === 'link' ? 'primary' : 'secondary'"
                  :outlined="finalMetadata.album.action !== 'link'"
                  size="small"
                  @click="finalMetadata.album.action = 'link'"
                >
                  <i class="pi pi-link mr-1"></i>Use Existing
                </Button>
              </div>
            </div>

            <div v-if="finalMetadata.album.action === 'link'" class="mb-6">
              <label class="mb-2 block text-xs font-medium text-surface-400">Search Existing Albums</label>
              <div class="relative">
                <IconField>
                  <InputIcon><i class="pi pi-search"></i></InputIcon>
                  <InputText
                    v-model="albumSearchQuery"
                    placeholder="Type album title..."
                    class="w-full"
                    @update:model-value="debouncedAlbumSearch"
                  />
                </IconField>
              </div>
              <div v-if="albumSearchResults.length > 0" class="mt-2 max-h-48 overflow-y-auto rounded border border-surface-700 bg-surface-800">
                <div
                  v-for="a in albumSearchResults"
                  :key="a.id"
                  class="flex cursor-pointer items-center gap-3 px-3 py-2 text-sm transition-colors hover:bg-surface-700"
                  :class="{ 'bg-primary/20': selectedAlbumId === a.id }"
                  @click="selectExistingAlbum(a)"
                >
                  <img
                    v-if="a.coverUrl"
                    :src="a.coverUrl"
                    class="h-8 w-8 rounded object-cover"
                  />
                  <div v-else class="flex h-8 w-8 items-center justify-center rounded bg-surface-600 text-xs">
                    <i class="pi pi-image"></i>
                  </div>
                  <div>
                    <p class="font-medium text-white">{{ a.title }}</p>
                    <p v-if="a.artistName" class="text-[11px] text-surface-500">{{ a.artistName }}</p>
                  </div>
                </div>
              </div>
              <p v-else-if="albumSearchQuery && !albumSearching" class="mt-2 text-xs text-surface-500">
                No existing albums found.
              </p>
            </div>

            <div class="space-y-4 rounded-lg border border-surface-700 bg-surface-800 p-5">
              <div class="grid gap-4 md:grid-cols-2">
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">
                    Title
                    <template v-if="findSuggestion('album')">
                      <Badge
                        :value="sourceLabel(findSuggestion('album')!.source)"
                        :severity="sourceSeverity(findSuggestion('album')!.source)"
                        size="small"
                      />
                      <Badge
                        :value="findSuggestion('album')!.confidence"
                        :severity="confidenceSeverity(findSuggestion('album')!.confidence)"
                        size="small"
                      />
                    </template>
                  </label>
                  <InputText
                    v-model="finalMetadata.album.title"
                    dir="auto"
                    class="w-full"
                    placeholder="Album title"
                  />
                </div>
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">
                    Release Year
                    <template v-if="findSuggestion('year')">
                      <Badge
                        :value="sourceLabel(findSuggestion('year')!.source)"
                        :severity="sourceSeverity(findSuggestion('year')!.source)"
                        size="small"
                      />
                      <Badge
                        :value="findSuggestion('year')!.confidence"
                        :severity="confidenceSeverity(findSuggestion('year')!.confidence)"
                        size="small"
                      />
                    </template>
                  </label>
                  <InputNumber
                    v-model="finalMetadata.album.releaseYear"
                    :use-grouping="false"
                    class="w-full"
                    placeholder="2024"
                  />
                </div>
              </div>
              <div class="grid gap-4 md:grid-cols-2">
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">Genre</label>
                  <InputText
                    v-model="finalMetadata.album.genre"
                    dir="auto"
                    class="w-full"
                    placeholder="Genre"
                  />
                </div>
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">Cover URL</label>
                  <InputText
                    v-model="finalMetadata.album.coverUrl"
                    class="w-full font-mono text-xs"
                    placeholder="https://..."
                  />
                </div>
              </div>
              <div v-if="enrichment?.musicbrainz?.albumMbid">
                <label class="mb-1 block text-xs font-medium text-surface-400">MusicBrainz Release ID</label>
                <InputText
                  :model-value="enrichment.musicbrainz.albumMbid"
                  class="w-full font-mono text-xs text-surface-400"
                  disabled
                />
              </div>
            </div>
          </div>

          <!-- Step 3: Track -->
          <div v-if="step === 2">
            <div class="mb-6">
              <h2 class="text-lg font-semibold text-white">Track</h2>
              <p class="mt-1 text-sm text-surface-400">Review and edit the track details.</p>
            </div>

            <div class="mb-6 flex flex-wrap gap-6">
              <div v-if="embeddedCover">
                <p class="mb-2 text-xs font-medium text-surface-500">Embedded Cover</p>
                <img
                  :src="embeddedCover"
                  alt="Embedded"
                  class="h-24 w-24 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
                  :class="finalMetadata.track.coverUrl === embeddedCover ? 'ring-green-500' : 'ring-transparent hover:ring-surface-400'"
                  @click="finalMetadata.track.coverUrl = embeddedCover"
                />
              </div>
              <div v-if="suggestedAlbumCover">
                <p class="mb-2 text-xs font-medium text-surface-500">
                  Spotify
                  <i class="pi pi-spotify ml-1 text-green-400 text-[10px]"></i>
                </p>
                <img
                  :src="suggestedAlbumCover"
                  alt="Spotify"
                  class="h-24 w-24 cursor-pointer rounded-lg object-cover shadow-md ring-2 transition-all"
                  :class="finalMetadata.track.coverUrl === suggestedAlbumCover ? 'ring-green-500' : 'ring-transparent hover:ring-surface-400'"
                  @click="finalMetadata.track.coverUrl = suggestedAlbumCover"
                />
              </div>
              <div v-if="finalMetadata.track.coverUrl && finalMetadata.track.coverUrl !== embeddedCover && finalMetadata.track.coverUrl !== suggestedAlbumCover" class="flex-shrink-0">
                <p class="mb-2 text-xs font-medium text-surface-500">Selected</p>
                <img
                  :src="finalMetadata.track.coverUrl"
                  alt="Selected"
                  class="h-24 w-24 rounded-lg object-cover shadow-md ring-2 ring-green-500"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
              </div>
            </div>

            <div class="space-y-4 rounded-lg border border-surface-700 bg-surface-800 p-5">
              <div class="grid gap-4 md:grid-cols-2">
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">
                    Title
                    <template v-if="findSuggestion('title')">
                      <Badge
                        :value="sourceLabel(findSuggestion('title')!.source)"
                        :severity="sourceSeverity(findSuggestion('title')!.source)"
                        size="small"
                      />
                      <Badge
                        :value="findSuggestion('title')!.confidence"
                        :severity="confidenceSeverity(findSuggestion('title')!.confidence)"
                        size="small"
                      />
                    </template>
                  </label>
                  <InputText
                    v-model="finalMetadata.track.title"
                    dir="auto"
                    class="w-full"
                    placeholder="Track title"
                  />
                </div>
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">Track Number</label>
                  <InputNumber
                    v-model="finalMetadata.track.trackNumber"
                    :use-grouping="false"
                    :min="1"
                    class="w-full"
                    placeholder="1"
                  />
                </div>
              </div>

              <div class="grid gap-4 md:grid-cols-3">
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">Duration (seconds)</label>
                  <InputNumber
                    :model-value="finalMetadata.track.durationSeconds"
                    class="w-full text-surface-400"
                    disabled
                  />
                </div>
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">Genre</label>
                  <InputText
                    v-model="finalMetadata.track.genre"
                    dir="auto"
                    class="w-full"
                    placeholder="Genre"
                  />
                </div>
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">
                    Explicit
                    <i class="pi pi-info-circle ml-1 text-surface-500 text-[10px]"></i>
                  </label>
                  <div class="flex items-center gap-2 pt-1">
                    <ToggleSwitch v-model="finalMetadata.track.explicit" />
                    <span class="text-xs text-surface-400">{{ finalMetadata.track.explicit ? 'Yes' : 'No' }}</span>
                  </div>
                </div>
              </div>

              <div>
                <div class="mb-1 flex items-center justify-between">
                  <label class="text-xs font-medium text-surface-400">Lyrics</label>
                  <div class="flex items-center gap-2">
                    <Badge
                      v-if="findSuggestion('lyrics_type')"
                      :value="findSuggestion('lyrics_type')!.value === 'lrc' ? 'Synced' : 'Plain'"
                      :severity="findSuggestion('lyrics_type')!.value === 'lrc' ? 'success' : 'info'"
                      size="small"
                    />
                    <Badge
                      v-if="findSuggestion('lyrics')"
                      :value="sourceLabel(findSuggestion('lyrics')!.source)"
                      :severity="sourceSeverity(findSuggestion('lyrics')!.source)"
                      size="small"
                    />
                    <Button
                      icon="pi pi-arrows-alt"
                      severity="secondary"
                      text
                      size="small"
                      :pt="{ root: { class: 'h-6 w-6' } }"
                      @click="lyricsExpanded = !lyricsExpanded"
                    />
                  </div>
                </div>
                <Textarea
                  v-model="finalMetadata.track.lyrics"
                  dir="auto"
                  :auto-resize="true"
                  class="w-full font-mono text-xs leading-relaxed"
                  :class="{ 'min-h-64': lyricsExpanded, 'min-h-24': !lyricsExpanded }"
                  placeholder="Lyrics..."
                />
                <p v-if="findSuggestion('lyrics_type')?.value === 'lrc'" class="mt-1 text-[11px] text-green-400">
                  <i class="pi pi-check-circle mr-1"></i>Synced LRC lyrics detected — timestamps will sync with playback.
                </p>
              </div>

              <div class="grid gap-4 md:grid-cols-2">
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">Cover URL</label>
                  <InputText
                    v-model="finalMetadata.track.coverUrl"
                    class="w-full font-mono text-xs"
                    placeholder="https://..."
                  />
                </div>
                <div>
                  <label class="mb-1 block text-xs font-medium text-surface-400">Spotify Preview URL</label>
                  <InputText
                    v-model="finalMetadata.track.spotifyPreviewUrl"
                    class="w-full font-mono text-xs text-surface-400"
                    :disabled="!enrichment?.spotify?.previewUrl"
                    :placeholder="enrichment?.spotify?.previewUrl || 'Not available'"
                  />
                </div>
              </div>
            </div>
          </div>

          <div v-if="step === 3">
            <div class="mb-6">
              <h2 class="text-lg font-semibold text-white">Confirm & Publish</h2>
              <p class="mt-1 text-sm text-surface-400">Review your selections before publishing to the catalog.</p>
            </div>

            <div class="mb-6 grid gap-6 md:grid-cols-3">
              <div class="rounded-lg border border-surface-700 bg-surface-800 p-4">
                <h3 class="mb-3 flex items-center gap-2 text-sm font-semibold text-white">
                  <i class="pi pi-user text-primary"></i>Artist
                </h3>
                <img
                  v-if="finalMetadata.artist.imageUrl"
                  :src="finalMetadata.artist.imageUrl"
                  alt=""
                  class="mb-2 h-24 w-24 rounded-lg object-cover shadow-sm"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
                <div class="space-y-1 text-xs">
                  <p><span class="text-surface-500">Action:</span> {{ finalMetadata.artist.action === 'link' ? 'Link existing' : 'Create new' }}</p>
                  <p><span class="text-surface-500">Name:</span> {{ finalMetadata.artist.name || '(empty)' }}</p>
                  <p v-if="finalMetadata.artist.country"><span class="text-surface-500">Country:</span> {{ finalMetadata.artist.country }}</p>
                </div>
              </div>
              <div class="rounded-lg border border-surface-700 bg-surface-800 p-4">
                <h3 class="mb-3 flex items-center gap-2 text-sm font-semibold text-white">
                  <i class="pi pi-book text-primary"></i>Album
                </h3>
                <img
                  v-if="finalMetadata.album.coverUrl"
                  :src="finalMetadata.album.coverUrl"
                  alt=""
                  class="mb-2 h-24 w-24 rounded-lg object-cover shadow-sm"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
                <div class="space-y-1 text-xs">
                  <p><span class="text-surface-500">Action:</span> {{ finalMetadata.album.action === 'link' ? 'Link existing' : 'Create new' }}</p>
                  <p><span class="text-surface-500">Title:</span> {{ finalMetadata.album.title || '(empty)' }}</p>
                  <p v-if="finalMetadata.album.releaseYear"><span class="text-surface-500">Year:</span> {{ finalMetadata.album.releaseYear }}</p>
                </div>
              </div>
              <div class="rounded-lg border border-surface-700 bg-surface-800 p-4">
                <h3 class="mb-3 flex items-center gap-2 text-sm font-semibold text-white">
                  <i class="pi pi-music text-primary"></i>Track
                </h3>
                <img
                  v-if="finalMetadata.track.coverUrl"
                  :src="finalMetadata.track.coverUrl"
                  alt=""
                  class="mb-2 h-24 w-24 rounded-lg object-cover shadow-sm"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
                <div class="space-y-1 text-xs">
                  <p><span class="text-surface-500">Title:</span> {{ finalMetadata.track.title || '(empty)' }}</p>
                  <p><span class="text-surface-500">Duration:</span> {{ formatDuration(finalMetadata.track.durationSeconds) }}</p>
                  <p><span class="text-surface-500">Explicit:</span> {{ finalMetadata.track.explicit ? 'Yes' : 'No' }}</p>
                </div>
              </div>
            </div>

            <Message v-if="publishingError" severity="error" :closable="false" class="mb-4">
              {{ publishingError }}
            </Message>
          </div>
        </div>
      </Transition>

      <div class="mt-8 flex items-center justify-between border-t border-surface-700 pt-6">
        <div>
          <Button
            v-if="step > 0"
            label="Back"
            icon="pi pi-chevron-left"
            severity="secondary"
            text
            @click="step--"
          />
        </div>
        <div class="flex items-center gap-3">
          <Button
            v-if="step < 3"
            label="Next"
            icon="pi pi-chevron-right"
            icon-pos="right"
            :disabled="!canProceed"
            @click="step++"
          />
          <Button
            v-if="step === 3"
            label="Publish to Catalog"
            icon="pi pi-check"
            :loading="publishing"
            :disabled="!isValid"
            @click="publish"
          />
          <Button
            label="Reject"
            icon="pi pi-times"
            severity="danger"
            text
            :loading="rejecting"
            @click="confirmReject"
          />
        </div>
      </div>
    </template>

    <Dialog
      v-model:visible="rejectDialogVisible"
      header="Reject Draft"
      :modal="true"
      class="w-full max-w-md"
    >
      <div class="space-y-4">
        <p class="text-sm text-surface-400">Provide a reason for rejecting this draft (optional).</p>
        <Textarea
          v-model="rejectReason"
          dir="auto"
          :auto-resize="true"
          class="w-full"
          rows="3"
          placeholder="Reason..."
        />
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" text @click="rejectDialogVisible = false" />
        <Button label="Confirm Reject" severity="danger" :loading="rejecting" @click="rejectDraft" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, onBeforeUnmount, watch } from 'vue'
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router'
import AdminSectionHeader from '@/components/admin/AdminSectionHeader.vue'
import { useIngestionApi } from '@/services/api/ingestion/routes'
import type {
  DraftDetailResponse,
  EnrichmentResult,
  EnrichedSuggestion,
  SaveFinalMetadataRequest,
  ArtistSearchResult,
  AlbumSearchResult,
  FinalizeResult,
} from '@/services/api/ingestion/types'
import { useToast } from 'primevue/usetoast'

// TODO MEDIUM: This component is 1135 lines — too large. Split into ArtistStep, AlbumStep, TrackStep, ConfirmStep sub-components.
defineOptions({ name: 'PageAdminIngestionReview' })

const route = useRoute()
const router = useRouter()
const toast = useToast()
const ingestionApi = useIngestionApi()

const draftId = route.params.id as string

const loading = ref(true)
const error = ref<string | null>(null)
const publishing = ref(false)
const publishingError = ref<string | null>(null)
const rejecting = ref(false)
const published = ref(false)
const publishResult = ref<FinalizeResult | null>(null)
const rejectDialogVisible = ref(false)
const rejectReason = ref('')
const lyricsExpanded = ref(false)

const enriching = ref(false)

const draftDetail = ref<DraftDetailResponse | null>(null)
const enrichment = ref<EnrichmentResult | null>(null)
const step = ref(0)

const steps = [
  { label: 'Artist', key: 'artist' },
  { label: 'Album', key: 'album' },
  { label: 'Track', key: 'track' },
  { label: 'Confirm', key: 'confirm' },
]

const hasUnsavedChanges = ref(false)

const finalMetadata = reactive<SaveFinalMetadataRequest>({
  artist: { action: 'create', name: '', bio: '', imageUrl: '', country: '', musicbrainzMbid: '' },
  album: { action: 'create', title: '', releaseYear: undefined, genre: '', coverUrl: '', musicbrainzReleaseId: '' },
  track: { title: '', trackNumber: undefined, durationSeconds: 0, genre: '', lyrics: '', explicit: false, spotifyPreviewUrl: '', coverUrl: '' },
})

const embeddedCover = computed(() => {
  return draftDetail.value?.assets?.find(a => a.assetType === 'cover')?.url || null
})

const suggestedAlbumCover = computed(() => {
  return enrichment.value?.spotify?.albumCoverUrl || null
})

const suggestedArtistImage = computed(() => {
  return enrichment.value?.spotify?.artistImageUrl || null
})

const selectedArtistId = ref<string | undefined>()
const selectedAlbumId = ref<string | undefined>()
const artistSearchQuery = ref('')
const albumSearchQuery = ref('')
const artistSearchResults = ref<ArtistSearchResult[]>([])
const albumSearchResults = ref<AlbumSearchResult[]>([])
const artistSearching = ref(false)
const albumSearching = ref(false)
let artistSearchTimer: ReturnType<typeof setTimeout> | null = null
let albumSearchTimer: ReturnType<typeof setTimeout> | null = null

const canProceed = computed(() => {
  if (step.value === 0) return finalMetadata.artist.name.trim().length > 0
  if (step.value === 1) return finalMetadata.album.title.trim().length > 0
  if (step.value === 2) return finalMetadata.track.title.trim().length > 0
  return true
})

const isValid = computed(() => {
  return (
    finalMetadata.artist.name.trim().length > 0 &&
    finalMetadata.album.title.trim().length > 0 &&
    finalMetadata.track.title.trim().length > 0
  )
})

function findSuggestion(field: string): EnrichedSuggestion | undefined {
  return enrichment.value?.suggestions?.find(s => s.field === field)
}

function sourceLabel(source: string): string {
  const map: Record<string, string> = { file: 'Embedded', musicbrainz: 'MB', lastfm: 'Last.fm', spotify: 'Spotify' }
  return map[source] || source
}

function sourceSeverity(source: string) {
  const map: Record<string, string> = { file: 'info', musicbrainz: 'warn', lastfm: 'help', spotify: 'success' }
  return map[source] || undefined as any
}

function confidenceSeverity(c: string) {
  const map: Record<string, string> = { exact_match: 'success', fuzzy: 'warn', fallback: 'danger' }
  return map[c] || undefined as any
}

function stepperClass(i: number): string {
  if (i === step.value) return 'bg-primary/20 text-primary'
  if (i < step.value) return 'bg-green-500/10 text-green-400'
  return 'text-surface-500 hover:text-surface-300'
}

function stepIconClass(i: number): string {
  if (i === step.value) return 'bg-primary text-white'
  if (i < step.value) return 'bg-green-500 text-white'
  return 'bg-surface-600 text-surface-400'
}

function truncateBio(bio: string): string {
  if (!bio) return ''
  return bio.length > 300 ? bio.slice(0, 300) + '...' : bio
}

function formatDuration(seconds?: number): string {
  if (!seconds) return ''
  const m = Math.floor(seconds / 60)
  const s = Math.round(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

function selectExistingArtist(a: ArtistSearchResult) {
  selectedArtistId.value = a.id
  finalMetadata.artist.existingId = a.id
  finalMetadata.artist.name = a.name
  finalMetadata.artist.country = a.country || ''
  finalMetadata.artist.imageUrl = a.imageUrl || ''
  finalMetadata.artist.bio = a.bio || ''
  hasUnsavedChanges.value = true
}

function selectExistingAlbum(a: AlbumSearchResult) {
  selectedAlbumId.value = a.id
  finalMetadata.album.existingId = a.id
  finalMetadata.album.title = a.title
  finalMetadata.album.coverUrl = a.coverUrl || ''
  if (a.releaseYear) finalMetadata.album.releaseYear = a.releaseYear
  hasUnsavedChanges.value = true
}

function debouncedArtistSearch() {
  if (artistSearchTimer) clearTimeout(artistSearchTimer)
  const q = artistSearchQuery.value.trim()
  if (q.length < 2) {
    artistSearchResults.value = []
    return
  }
  artistSearchTimer = setTimeout(async () => {
    artistSearching.value = true
    try {
      const res = await ingestionApi.searchArtists(q)
      artistSearchResults.value = Array.isArray(res) ? res : []
    } catch {
      artistSearchResults.value = []
    } finally {
      artistSearching.value = false
    }
  }, 300)
}

function debouncedAlbumSearch() {
  if (albumSearchTimer) clearTimeout(albumSearchTimer)
  const q = albumSearchQuery.value.trim()
  if (q.length < 2) {
    albumSearchResults.value = []
    return
  }
  albumSearchTimer = setTimeout(async () => {
    albumSearching.value = true
    try {
      const res = await ingestionApi.searchAlbums(q)
      albumSearchResults.value = Array.isArray(res) ? res : []
    } catch {
      albumSearchResults.value = []
    } finally {
      albumSearching.value = false
    }
  }, 300)
}

async function loadDraft() {
  loading.value = true
  error.value = null
  try {
    const detail = await ingestionApi.getDraftDetail(draftId)
    if (!detail) {
      error.value = 'Draft not found.'
      return
    }
    draftDetail.value = detail
    enrichment.value = detail.enrichedMetadata || null
    if (detail.status === 'enriching') {
      enriching.value = true
      pollEnrichment()
    }
    prefillForm(detail)
  } catch (err: any) {
    error.value = err?.message || 'Failed to load draft.'
  } finally {
    loading.value = false
  }
}

function pollEnrichment() {
  const maxAttempts = 30
  let attempts = 0
  const timer = setInterval(async () => {
    attempts++
    try {
      const detail = await ingestionApi.getDraftDetail(draftId)
      draftDetail.value = detail
      if (detail && detail.status !== 'enriching') {
        clearInterval(timer)
        enriching.value = false
        enrichment.value = detail.enrichedMetadata || null
        prefillForm(detail)
      }
    } catch {
      // ignore polling errors
    }
    if (attempts >= maxAttempts) {
      clearInterval(timer)
      enriching.value = false
    }
  }, 2000)
}

async function triggerEnrich() {
  enriching.value = true
  try {
    await ingestionApi.enrichDraft(draftId)
    pollEnrichment()
  } catch {
    enriching.value = false
    toast.add({ severity: 'error', summary: 'Failed to start enrichment', life: 3000 })
  }
}

function prefillForm(detail: DraftDetailResponse) {
  const tags = detail.extractedMetadata
  const sug = enrichment.value?.suggestions || []
  const mb = enrichment.value?.musicbrainz
  const lfm = enrichment.value?.lastfm
  const spot = enrichment.value?.spotify

  const sugMap = new Map<string, any>()
  for (const s of sug) {
    sugMap.set(s.field, s.value)
  }

  // Artist
  finalMetadata.artist.name = sugMap.get('artist') || tags?.artist || ''
  finalMetadata.artist.bio = lfm?.artistBio || sugMap.get('artist_bio') || ''
  finalMetadata.artist.imageUrl = spot?.artistImageUrl || sugMap.get('artist_image_url') || ''
  if (mb?.artistMbid) finalMetadata.artist.musicbrainzMbid = mb.artistMbid

  // Album
  finalMetadata.album.title = sugMap.get('album') || tags?.album || ''
  finalMetadata.album.releaseYear = sugMap.get('year') || (tags?.year || undefined) as any
  finalMetadata.album.genre = sugMap.get('genre') || tags?.genre || ''
  finalMetadata.album.coverUrl = spot?.albumCoverUrl || sugMap.get('album_cover_url') || embeddedCover.value || ''
  if (mb?.albumMbid) finalMetadata.album.musicbrainzReleaseId = mb.albumMbid

  // Track
  finalMetadata.track.title = sugMap.get('title') || tags?.title || ''
  finalMetadata.track.trackNumber = tags?.trackNumber || undefined
  finalMetadata.track.durationSeconds = Math.round(detail.durationSeconds || tags?.duration || 0)
  finalMetadata.track.genre = sugMap.get('genre') || tags?.genre || ''
  finalMetadata.track.lyrics = sugMap.get('lyrics') || tags?.lyrics || ''
  finalMetadata.track.explicit = false
  finalMetadata.track.spotifyPreviewUrl = spot?.previewUrl || ''
  finalMetadata.track.coverUrl = embeddedCover.value || spot?.albumCoverUrl || ''
}

function goBack() {
  if (hasUnsavedChanges.value) {
    if (!window.confirm('You have unsaved changes. Are you sure you want to leave?')) return
  }
  router.push({ name: 'admin.ingestion' })
}

async function uploadArtistImage(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  try {
    const res = await ingestionApi.uploadDraftImage(draftId, 'artist', file)
    if (res?.url) {
      finalMetadata.artist.imageUrl = res.url
      hasUnsavedChanges.value = true
      toast.add({ severity: 'success', summary: 'Image uploaded', detail: 'Artist image uploaded successfully.', life: 3000 })
    }
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Upload failed', detail: err?.message || 'Failed to upload image.', life: 5000 })
  }
  target.value = ''
}

async function uploadAlbumCover(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  try {
    const res = await ingestionApi.uploadDraftImage(draftId, 'album', file)
    if (res?.url) {
      finalMetadata.album.coverUrl = res.url
      hasUnsavedChanges.value = true
      toast.add({ severity: 'success', summary: 'Cover uploaded', detail: 'Album cover uploaded successfully.', life: 3000 })
    }
  } catch (err: any) {
    toast.add({ severity: 'error', summary: 'Upload failed', detail: err?.message || 'Failed to upload image.', life: 5000 })
  }
  target.value = ''
}

async function publish() {
  publishing.value = true
  publishingError.value = null
  try {
    await ingestionApi.saveFinalMetadata(draftId, { ...finalMetadata })
    const result = await ingestionApi.finalizeDraft(draftId)
    publishResult.value = result
    hasUnsavedChanges.value = false
    published.value = true
    toast.add({ severity: 'success', summary: 'Draft published', detail: 'The draft has been published to the catalog.', life: 5000 })
  } catch (err: any) {
    publishingError.value = err?.message || 'Failed to publish draft.'
  } finally {
    publishing.value = false
  }
}

function confirmReject() {
  rejectDialogVisible.value = true
}

async function rejectDraft() {
  rejecting.value = true
  try {
    await ingestionApi.rejectDraft(draftId, rejectReason.value)
    hasUnsavedChanges.value = false
    rejectDialogVisible.value = false
    toast.add({ severity: 'info', summary: 'Draft rejected', detail: 'The draft has been rejected.', life: 4000 })
    setTimeout(() => router.push({ name: 'admin.ingestion' }), 1500)
  } catch {
    // error handled by useRequest
  } finally {
    rejecting.value = false
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (rejectDialogVisible.value || loading.value || published.value) return
  const tag = (e.target as HTMLElement)?.tagName
  const isInput = tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT'
  if (e.key === 'Escape') {
    if (!isInput && step.value > 0) {
      e.preventDefault()
      step.value--
    }
  } else if (e.key === 'Enter') {
    if (!isInput) {
      e.preventDefault()
      if (step.value < 3 && canProceed.value) {
        step.value++
      } else if (step.value === 3 && isValid.value) {
        publish()
      }
    }
  }
}

// TODO MEDIUM: finalMetadata is reactive (not a ref). Ensure deep watch on reactive tracks changes properly.
watch(finalMetadata, () => {
  hasUnsavedChanges.value = true
}, { deep: true })

onMounted(() => {
  loadDraft()
  window.addEventListener('keydown', handleKeydown)
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onUnmounted(() => {
  if (artistSearchTimer) clearTimeout(artistSearchTimer)
  if (albumSearchTimer) clearTimeout(albumSearchTimer)
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

onBeforeRouteLeave((_to, _from, next) => {
  if (hasUnsavedChanges.value) {
    const answer = window.confirm('You have unsaved changes. Leave anyway?')
    if (!answer) {
      next(false)
      return
    }
  }
  next()
})

function handleBeforeUnload(e: BeforeUnloadEvent) {
  if (hasUnsavedChanges.value) {
    e.preventDefault()
  }
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
