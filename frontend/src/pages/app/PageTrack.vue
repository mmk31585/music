<template>
  <div :key="String(route.params.id)" class="relative mx-auto min-h-screen pb-36">
    <!-- Ambient background derived from cover art -->
    <div
      class="pointer-events-none fixed inset-0 transition-all duration-1000"
      :style="ambientBg"
      aria-hidden="true"
    />
    <div
      class="pointer-events-none fixed inset-0 opacity-[0.015]"
      style="background-image: repeating-radial-gradient(circle at 50% 50%, transparent 0, transparent 2px, rgba(255,255,255,0.04) 2px, rgba(255,255,255,0.04) 3px); background-size: 6px 6px;"
      aria-hidden="true"
    />

    <div class="relative z-10 mx-auto w-full max-w-7xl px-4 pt-4 md:px-6 lg:px-8">
      <!-- Back button -->
      <button
        type="button"
        class="mb-6 inline-flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm text-white/50 transition hover:bg-white/6 hover:text-white"
        @click="goBack"
      >
        <i aria-hidden="true" class="pi pi-arrow-left text-xs" />
        Back
      </button>

      <!-- Loading -->
      <div v-if="loading" class="space-y-6">
        <div class="flex flex-col gap-10 md:flex-row">
          <SkeletonLoader variant="card" class="h-85 w-85 shrink-0" />
          <div class="flex-1 space-y-4">
            <SkeletonLoader variant="lines" :lines="1" class="max-w-sm" />
            <SkeletonLoader variant="lines" :lines="1" class="max-w-xs" />
            <SkeletonLoader variant="lines" :lines="2" class="max-w-md" />
            <div class="mt-8 flex gap-4">
              <SkeletonLoader variant="card" class="h-12 w-32 rounded-full" />
              <SkeletonLoader variant="card" class="h-12 w-32 rounded-full" />
            </div>
          </div>
        </div>
        <div class="space-y-2">
          <SkeletonLoader v-for="i in 5" :key="i" variant="track" />
        </div>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="flex flex-col items-center gap-4 py-24 text-center">
        <div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-white/4">
          <i aria-hidden="true" class="pi pi-exclamation-circle text-3xl text-slate-500" />
        </div>
        <h2 class="text-xl font-bold text-white">Track not found</h2>
        <p class="text-sm text-slate-400">This track may have been removed or the link is invalid.</p>
        <RouterLink to="/" class="mt-2 text-sm font-medium text-spotify underline underline-offset-2">
          Go home
        </RouterLink>
      </div>

      <!-- Content -->
      <template v-else-if="track">
        <!-- ════════════════════════════════════════ -->
        <!-- HERO SECTION                           -->
        <!-- ════════════════════════════════════════ -->
        <div class="flex flex-col gap-10 md:flex-row md:items-end">
          <!-- Left: Cover Art -->
          <div class="group shrink-0">
            <div
              class="relative h-75 w-75 overflow-hidden rounded-2xl bg-white/6 shadow-2xl ring-1 ring-white/10 transition-all duration-500 md:h-85 md:w-85"
              :style="coverGlowStyle"
            >
              <img
                v-if="track.cover_url"
                :src="track.cover_url"
                :alt="track.title"
                loading="eager"
                class="h-full w-full select-none object-cover"
                @error="onImgError"
              />
              <div v-else class="flex h-full items-center justify-center">
                <i aria-hidden="true" class="pi pi-compact-disc text-5xl text-slate-500" />
              </div>

              <!-- Playing indicator overlay -->
              <div
                v-if="isPlaying"
                class="absolute inset-0 flex items-center justify-center bg-black/20 backdrop-blur-[1px]"
              >
                <div class="flex items-end gap-0.75">
                  <span class="now-playing-bar h-5 w-0.75 rounded-full bg-white" />
                  <span class="now-playing-bar h-8 w-0.75 rounded-full bg-white" style="animation-delay: 0.15s" />
                  <span class="now-playing-bar h-6 w-0.75 rounded-full bg-white" style="animation-delay: 0.3s" />
                  <span class="now-playing-bar h-4 w-0.75 rounded-full bg-white" style="animation-delay: 0.45s" />
                  <span class="now-playing-bar h-7 w-0.75 rounded-full bg-white" style="animation-delay: 0.2s" />
                </div>
              </div>

              <!-- Admin edit overlay -->
              <RouterLink
                v-if="isAdmin"
                :to="`/admin/catalog/tracks/${trackId}`"
                class="absolute top-3 right-3 flex h-8 w-8 items-center justify-center rounded-full bg-black/50 text-white opacity-0 backdrop-blur-xs transition-opacity hover:bg-black/70 group-hover:opacity-100"
                title="Edit track in admin"
                @click.stop
              >
                <i aria-hidden="true" class="pi pi-pencil text-xs" />
              </RouterLink>
            </div>
            <!-- Sleeve frame accent -->
            <div
              class="pointer-events-none absolute inset-0 rounded-2xl ring-1 ring-inset ring-white/4"
              aria-hidden="true"
            />
          </div>

          <!-- Right: Track Info + Controls -->
          <div class="flex-1">
            <!-- Eyebrow + Admin badge -->
            <div class="flex items-center gap-3">
              <p class="text-[10px] font-bold tracking-[0.35em] text-white/30 uppercase">Track</p>
              <span
                v-if="isAdmin"
                class="rounded-full bg-amber-500/15 px-2 py-0.5 text-[9px] font-bold text-amber-400"
              >Admin</span>
            </div>

            <!-- Title -->
            <div class="mt-3 flex items-start gap-3">
              <h1 class="text-3xl font-black leading-[1.1] text-white md:text-4xl lg:text-5xl">
                {{ track.title }}
              </h1>
              <span
                v-if="track.explicit"
                class="mt-1.5 inline-flex h-6 w-6 shrink-0 items-center justify-center rounded bg-white/15 text-[11px] font-bold tracking-wide text-white"
                title="Explicit"
              >E</span>
            </div>

            <!-- Artists (all) -->
            <div class="mt-4 flex flex-wrap items-center gap-x-2 gap-y-1.5 text-sm">
              <template v-for="(a, i) in trackArtists" :key="String(a.artistId || i)">
                <RouterLink
                  :to="`/artist/${a.artistId}`"
                  class="inline-flex items-center gap-1.5 font-bold text-white underline underline-offset-4 decoration-white/20 transition hover:text-spotify hover:decoration-[#1db954]"
                >
                  {{ a.name }}
                </RouterLink>
                <span
                  v-if="a.role && !['main', 'primary'].includes(a.role as string)"
                  class="rounded-full bg-purple-500/15 px-1.5 py-0.5 text-[9px] font-medium text-purple-400"
                >feat.</span>
                <span v-if="i < trackArtists.length - 1" class="text-white/30">,</span>
              </template>
              <span v-if="track.album_title" class="text-white/30">&middot;</span>
              <RouterLink
                v-if="track.album_title"
                :to="`/album/${track.album_id}`"
                class="text-white/70 transition hover:text-white"
              >
                {{ track.album_title }}
              </RouterLink>
            </div>

            <!-- Genres -->
            <div
              v-if="genreList.length > 0"
              class="mt-4 flex flex-wrap items-center gap-1.5"
            >
              <span
                v-for="g in genreList"
                :key="g.id"
                class="rounded-full border border-white/10 bg-white/6 px-3 py-0.5 text-xs font-medium text-slate-300 transition hover:border-spotify/30 hover:text-white"
              >
                {{ g.name }}
              </span>
            </div>

            <!-- Play count -->
            <p v-if="track.play_count > 0" class="mt-3 text-xs text-slate-500">
              {{ formatPlayCount(track.play_count) }} plays
            </p>

            <!-- Main Controls (prev / play / next) -->
            <div class="mt-8 flex items-center gap-4">
              <button
                type="button"
                aria-label="Previous track"
                class="flex h-10 w-10 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
                @click="player.playPrevious"
              >
                <i aria-hidden="true" class="pi pi-step-backward text-lg" />
              </button>

              <button
                type="button"
                aria-label="Toggle play"
                class="relative flex h-14 w-14 items-center justify-center rounded-full bg-white text-black shadow-2xl transition-all hover:scale-105 hover:bg-spotify hover:text-white active:scale-95 md:h-16 md:w-16"
                @click="togglePlay"
              >
                <i
                  :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'"
                  class="ml-0.5 text-xl md:text-2xl"
                />
              </button>

              <button
                type="button"
                aria-label="Next track"
                class="flex h-10 w-10 items-center justify-center rounded-full text-slate-400 transition hover:bg-white/10 hover:text-white"
                @click="player.playNext"
              >
                <i aria-hidden="true" class="pi pi-step-forward text-lg" />
              </button>
            </div>

            <!-- Progress bar -->
            <div class="mt-5 flex w-full max-w-md items-center gap-3">
              <span class="min-w-10 text-right text-xs font-medium text-slate-500 tabular-nums">
                {{ formatTime(player.currentTime.value) }}
              </span>
              <input
                type="range"
                min="0"
                max="100"
                step="0.1"
                aria-label="Seek"
                class="player-range flex-1"
                :style="{ '--range-progress': `${Number(player.progressPercent.value || 0)}%` }"
                :value="player.progressPercent.value"
                @input="onSeek"
              />
              <span class="min-w-10 text-xs font-medium text-slate-500 tabular-nums">
                {{ formatTime(player.duration.value || track.duration_seconds || 0) }}
              </span>
            </div>

            <!-- Action Buttons -->
            <div class="mt-6 flex flex-wrap items-center gap-3">
              <button
                type="button"
                class="glow-green inline-flex items-center gap-2.5 rounded-full bg-spotify px-7 py-2.5 text-sm font-bold text-black transition hover:scale-105 hover:bg-spotify-hover"
                @click="togglePlay"
              >
                <i aria-hidden="true" :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" />
                {{ isPlaying ? 'Pause' : 'Play' }}
              </button>

              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-full border border-white/15 bg-white/4 px-5 py-2.5 text-sm font-bold transition"
                :class="isLiked ? 'border-spotify/30 text-spotify' : 'text-white/80 hover:border-white/30 hover:bg-white/8 hover:text-white'"
                @click="toggleLike"
              >
                <i aria-hidden="true" :class="isLiked ? 'pi pi-heart-fill' : 'pi pi-heart'" />
                {{ isLiked ? 'Liked' : 'Like' }}
              </button>

              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/3 px-5 py-2.5 text-sm font-bold text-white/60 transition hover:bg-white/8 hover:text-white"
                @click="showAddToPlaylist = true"
              >
                <i aria-hidden="true" class="pi pi-plus" />
                Playlist
              </button>

              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/3 px-5 py-2.5 text-sm font-bold text-white/60 transition hover:bg-white/8 hover:text-white"
                @click="shareTrack"
              >
                <i aria-hidden="true" class="pi pi-share-alt" />
                Share
              </button>

              <!-- Admin edit button -->
              <RouterLink
                v-if="isAdmin"
                :to="`/admin/catalog/tracks/${trackId}`"
                class="inline-flex items-center gap-2 rounded-full border border-amber-500/20 bg-amber-500/10 px-5 py-2.5 text-sm font-bold text-amber-400 transition hover:bg-amber-500/20"
              >
                <i aria-hidden="true" class="pi pi-pencil text-sm" />
                Edit
              </RouterLink>

              <!-- Create edit (always shown) -->
              <button
                type="button"
                class="inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/3 px-5 py-2.5 text-sm font-bold text-white/60 transition hover:bg-white/8 hover:text-white"
                @click="goToCreateEdit"
              >
                <i aria-hidden="true" class="pi pi-video text-sm" />
                Edit Video
              </button>
            </div>

            <!-- Secondary Controls -->
            <div class="mt-6 flex items-center gap-6">
              <button
                type="button"
                class="flex items-center gap-2 text-sm font-medium transition"
                :class="player.shuffleMode ? 'text-spotify' : 'text-slate-400 hover:text-white'"
                @click="player.toggleShuffle"
              >
                <i class="pi pi-sort-alt text-lg" />
                Shuffle
              </button>

              <button
                type="button"
                class="relative flex items-center gap-2 text-sm font-medium transition"
                :class="player.repeatMode.value !== 'off' ? 'text-spotify' : 'text-slate-400 hover:text-white'"
                @click="player.toggleRepeat"
              >
                <i aria-hidden="true" class="pi pi-refresh text-lg" />
                <span
                  v-if="player.repeatMode.value === 'one'"
                  class="absolute -top-1 -right-3 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-spotify text-[8px] font-bold text-black"
                >1</span>
                Repeat
              </button>

              <button
                type="button"
                class="flex items-center gap-2 text-sm font-medium text-slate-400 transition hover:text-white"
                @click="showQueue = true"
              >
                <i aria-hidden="true" class="pi pi-list text-lg" />
                Queue
              </button>
            </div>
          </div>
        </div>

        <!-- ════════════════════════════════════════ -->
        <!-- MUSIC VISUALIZER                       -->
        <!-- ════════════════════════════════════════ -->
        <section class="mx-auto mt-16 w-full max-w-2xl">
          <div class="relative h-28 overflow-hidden rounded-2xl border border-white/4 bg-white/2 backdrop-blur-xs">
            <div class="flex h-full items-end justify-center gap-0.5 px-4 pb-3">
              <div
                v-for="i in 96"
                :key="i"
                class="visualizer-bar w-[2.5px] rounded-t-full"
                :class="isPlaying ? 'bg-white/20' : 'bg-white/4'"
                :style="{
                  height: isPlaying ? `${getBarHeight(i)}%` : '8%',
                  animationDelay: isPlaying ? `${i * 0.025}s` : '0s',
                }"
              />
            </div>

            <!-- Center play indicator -->
            <div
              class="absolute inset-0 flex cursor-pointer items-center justify-center bg-black/10 opacity-0 transition hover:opacity-100"
              @click="togglePlay"
            >
              <div
                class="flex h-12 w-12 items-center justify-center rounded-full bg-spotify/90 text-black shadow-xl backdrop-blur-xs transition-transform hover:scale-110"
              >
                <i aria-hidden="true" :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-lg" />
              </div>
            </div>
          </div>
        </section>

        <!-- ════════════════════════════════════════ -->
        <!-- MUSIC VIDEOS SECTION                   -->
        <!-- ════════════════════════════════════════ -->
        <section v-if="trackVideos.length > 0" class="mt-16">
          <div class="relative mb-6">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-white/6" />
            </div>
            <div class="relative flex justify-between items-center">
              <span class="bg-surface-base px-4 text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase">
                Music Videos
              </span>
              <RouterLink
                to="/videos"
                class="bg-surface-base px-4 text-xs font-medium text-spotify transition hover:text-spotify-hover"
              >
                Browse all &rarr;
              </RouterLink>
            </div>
          </div>
          <!-- Mobile: horizontal snap scroll | Desktop: grid -->
          <div class="flex gap-3 overflow-x-auto px-1 pb-2 snap-x snap-mandatory scrollbar-hidden sm:grid sm:grid-cols-3 sm:overflow-visible sm:px-0 sm:pb-0 sm:snap-none md:grid-cols-4 lg:grid-cols-5">
            <div
              v-for="v in trackVideos"
              :key="String(v.id)"
              role="button"
              tabindex="0"
              class="group w-[45vw] shrink-0 snap-start cursor-pointer overflow-hidden rounded-xl bg-white/5 transition-all duration-150 hover:bg-white/8 hover:shadow-lg hover:shadow-black/20 sm:w-auto sm:shrink"
              @click="openVideo(v)"
              @keydown.enter="openVideo(v)"
              @keydown.space.prevent="openVideo(v)"
            >
              <div class="relative aspect-9/16 w-full overflow-hidden">
                <img
                  v-if="v.thumbnail_url || v.thumbnail_path || v.track_cover_url"
                  :src="(v.thumbnail_url || v.thumbnail_path || v.track_cover_url) ?? undefined"
                  :alt="v.title"
                  class="h-full w-full object-cover transition-transform duration-150 group-hover:scale-[1.03]"
                  loading="lazy"
                  @error="($event.target as HTMLImageElement).style.display='none'"
                />
                <div v-else class="flex h-full items-center justify-center bg-white/3">
                  <i aria-hidden="true" class="pi pi-video text-2xl text-slate-500" />
                </div>

                <!-- Overlay gradient -->
                <div class="pointer-events-none absolute inset-0 bg-linear-to-t from-black/70 via-transparent to-transparent" />

                <!-- Type badge -->
                <div
                  class="absolute top-2 right-2 rounded-full px-2 py-0.5 text-[10px] font-bold backdrop-blur-xs"
                  :class="v.type === 'official_mv' ? 'bg-spotify/20 text-spotify' : 'bg-blue-500/20 text-blue-400'"
                >
                  {{ v.type === 'official_mv' ? 'MV' : 'Edit' }}
                </div>

                <!-- Bottom stats -->
                <div class="absolute right-2 bottom-2 left-2 flex items-center justify-between">
                  <div class="flex items-center gap-2 text-[11px] font-medium text-white/80">
                    <span class="flex items-center gap-1">
                      <i aria-hidden="true" class="pi pi-eye text-[10px]" />
                      {{ formatCount(v.view_count) }}
                    </span>
                    <span class="flex items-center gap-1">
                      <i aria-hidden="true" class="pi pi-heart text-[10px]" />
                      {{ formatCount(v.like_count) }}
                    </span>
                  </div>
                  <div class="flex h-7 w-7 items-center justify-center rounded-full bg-black/50 text-white opacity-0 backdrop-blur-xs transition-opacity group-hover:opacity-100">
                    <i aria-hidden="true" class="pi pi-play-fill text-xs" />
                  </div>
                </div>
              </div>
              <div class="p-2.5">
                <p class="truncate text-xs font-semibold text-white/90">{{ v.title }}</p>
                <p v-if="v.uploader" class="mt-0.5 truncate text-[10px] text-white/40">
                  {{ v.uploader.username }}
                </p>
              </div>
            </div>
          </div>
        </section>

        <!-- ════════════════════════════════════════ -->
        <!-- LYRICS SECTION                         -->
        <!-- ════════════════════════════════════════ -->
        <section class="mt-16">
          <div class="relative mb-6">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-white/6" />
            </div>
            <div class="relative flex justify-center">
              <span class="bg-surface-base px-4 text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase">
                Lyrics
              </span>
            </div>
          </div>

          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-bold text-white">Lyrics</h2>
            <button
              v-if="lyrics && lyrics.content"
              type="button"
              class="flex items-center gap-1.5 rounded-full bg-white/10 px-3 py-1.5 text-xs font-bold transition hover:bg-white/15"
              :class="karaokeActive ? 'bg-spotify/15 text-spotify' : 'text-white/60'"
              @click="karaokeActive = !karaokeActive"
            >
              <i aria-hidden="true" class="pi pi-mic text-[10px]" />
              Karaoke
            </button>
          </div>

          <div
            v-if="karaokeActive && lyrics?.content"
            class="overflow-hidden rounded-2xl border border-white/6 bg-white/3 backdrop-blur-xl"
            style="height: 400px"
          >
            <KaraokeLyrics
              :content="lyrics.content"
              :type="lyrics.type || 'plain'"
              :language="lyrics.language || 'en'"
              :current-time="player.currentTime.value"
              :loading="loading"
              :karaoke="true"
              @seek="player.seek"
            />
          </div>
          <div
            v-else
            class="rounded-2xl border border-white/6 bg-white/3 p-6 backdrop-blur-xl"
          >
            <LyricsDisplay :lyrics="lyrics" :loading="loading" :error="lyrics! && loading!" />
          </div>
        </section>

        <!-- ════════════════════════════════════════ -->
        <!-- CREDITS SECTION                        -->
        <!-- ════════════════════════════════════════ -->
        <section v-if="trackArtists.length > 0 || trackCredits.length > 0" class="mt-16">
          <div class="relative mb-8">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-white/6" />
            </div>
            <div class="relative flex justify-center">
              <span class="bg-surface-base px-4 text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase">
                Credits
              </span>
            </div>
          </div>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <!-- Artists -->
            <div
              v-for="(a, i) in trackArtists"
              :key="String(a.artistId || i)"
              class="group flex items-center gap-4 rounded-2xl border border-white/4 bg-white/2 px-5 py-4 transition hover:border-white/8 hover:bg-white/4"
            >
              <div
                class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-linear-to-br from-white/8 to-white/2 text-sm font-bold text-white/70 ring-1 ring-white/4"
              >
                {{ String(a.name).charAt(0).toUpperCase() }}
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-2">
                  <RouterLink
                    :to="`/artist/${a.artistId}`"
                    class="block truncate text-sm font-semibold text-white transition group-hover:text-spotify"
                  >
                    {{ a.name }}
                  </RouterLink>
                  <span
                    v-if="a.role && !['main', 'primary'].includes(a.role as string)"
                    class="rounded-full bg-purple-500/15 px-1.5 py-0.5 text-[9px] font-medium text-purple-400"
                  >feat.</span>
                </div>
                <p class="mt-0.5 text-xs text-white/30">
                  {{ a.role && !['main', 'primary'].includes(a.role as string) ? a.role : 'Main artist' }}
                </p>
              </div>
            </div>

            <!-- Credits -->
            <div
              v-for="c in trackCredits"
              :key="String(c.id)"
              class="group flex items-center gap-4 rounded-2xl border border-white/4 bg-white/2 px-5 py-4 transition hover:border-white/8 hover:bg-white/4"
            >
              <div
                class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-linear-to-br from-white/8 to-white/2 text-sm font-bold text-white/70 ring-1 ring-white/4"
              >
                {{ String(c.artistName).charAt(0).toUpperCase() }}
              </div>
              <div class="min-w-0">
                <p class="truncate text-sm font-semibold text-white">{{ c.artistName }}</p>
                <p class="mt-0.5 text-xs text-white/40 capitalize">{{ c.creditType }}</p>
              </div>
            </div>
          </div>
        </section>

        <!-- ════════════════════════════════════════ -->
        <!-- YOU MIGHT LIKE                         -->
        <!-- ════════════════════════════════════════ -->
        <section v-if="similarTracks.length" class="mt-16" aria-live="polite">
          <div class="relative mb-8">
            <div class="absolute inset-0 flex items-center">
              <div class="w-full border-t border-white/6" />
            </div>
            <div class="relative flex justify-center">
              <span class="bg-surface-base px-4 text-[10px] font-bold tracking-[0.3em] text-white/20 uppercase">
                You might like
              </span>
            </div>
          </div>
          <div class="space-y-1">
            <div
              v-for="(st, index) in similarTracks"
              :key="st.id"
              class="group grid grid-cols-[48px_1fr_auto] items-center gap-4 rounded-xl px-3 py-2.5 transition-all duration-200 hover:bg-white/8"
            >
              <!-- Play button -->
              <button
                type="button"
                aria-label="Play track"
                class="relative flex h-11 w-11 items-center justify-center overflow-hidden rounded-xl bg-white/10 text-white transition-all duration-200 hover:scale-105 hover:bg-spotify hover:text-black"
                :class="{ 'bg-spotify text-black': isCurrentSimilarTrack(st) }"
                @click="playSimilar(st, index)"
              >
                <img
                  v-if="st.cover_url"
                  :src="st.cover_url"
                  :alt="st.title"
                  class="absolute inset-0 h-full w-full object-cover opacity-60 transition group-hover:opacity-35"
                  loading="lazy"
                  @error="onImgError"
                />
                <span class="relative z-10 flex items-center justify-center">
                  <template v-if="isCurrentSimilarTrack(st) && player.isPlaying.value">
                    <span class="flex h-4 items-end gap-0.5" aria-label="Playing">
                      <span class="similar-eq-bar h-2" />
                      <span class="similar-eq-bar animation-delay-150 h-4" />
                      <span class="similar-eq-bar animation-delay-300 h-3" />
                    </span>
                  </template>
                  <i v-else aria-hidden="true" class="pi pi-play text-sm" />
                </span>
              </button>

              <div class="min-w-0">
                <RouterLink
                  :to="`/track/${st.id}`"
                  class="truncate text-sm font-semibold transition hover:underline"
                  :class="isCurrentSimilarTrack(st) ? 'text-spotify' : 'text-white'"
                  @click.stop
                >
                  {{ st.title }}
                </RouterLink>
                <p class="mt-0.5 truncate text-xs text-slate-400">{{ st.artist_name }}</p>
              </div>

              <div class="flex items-center gap-4 text-xs text-slate-400">
                <span class="tabular-nums">{{ formatTime(st.duration_seconds) }}</span>
              </div>
            </div>
          </div>
        </section>
      </template>
    </div>

    <!-- Queue Panel -->
    <QueuePanel v-model:visible="showQueue" />

    <!-- Add to Playlist Dialog -->
    <AddToPlaylistDialog
      v-model:visible="showAddToPlaylist"
      :track-id="trackId"
      :track-title="track?.title ?? ''"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { SkeletonLoader } from '@/components/common'
import { useTrack } from '@/composables/catalog/useTrack'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { buildPlaybackTrack, mapToPlaybackTracks } from '@/factories/playbackTrack'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { useSocialShare } from '@/composables/social'
import { useVideoApi } from '@/services/api/video'
import { useUserAuthStore } from '@/stores'
import { onImgError } from '@/utils/helpers'
import { formatCount } from '@/utils/number'
import { LyricsDisplay, QueuePanel, KaraokeLyrics, AddToPlaylistDialog } from '@/components/music'
import type { VideoItem } from '@/services/api/video/types'

const route = useRoute()
const router = useRouter()
const trackId = String(route.params.id)

const {
  track,
  similarTracks,
  lyrics,
  trackArtists: _trackArtists,
  trackCredits: _trackCredits,
  genreList,
  isLiked,
  loading,
  error,
  fetchTrack,
  toggleLike,
} = useTrack(trackId)

const trackArtists = computed(() => _trackArtists.value as unknown as Record<string, unknown>[])
const trackCredits = computed(() => _trackCredits.value as unknown as Record<string, unknown>[])
const player = usePlayer()
const playerApi = usePlayerApi()
const videoApi = useVideoApi()
const authStore = useUserAuthStore()
const showQueue = ref(false)
const showAddToPlaylist = ref(false)
const karaokeActive = ref(false)
const trackVideos = ref<VideoItem[]>([])

const isAdmin = computed(() => authStore.isAdmin)
const isPlaying = player.isPlaying

// ── Ambient background from cover art ──
const coverUrl = computed(() => track.value?.cover_url || null)
const { palette } = useAlbumColors(coverUrl)
const accentColor = computed(() => palette.value.vibrant || '#1db954')

const ambientBg = computed(() => {
  if (!coverUrl.value) return { background: '#0A0A0F' }
  const c = accentColor.value
  return {
    background: `
      radial-gradient(ellipse 80% 50% at 50% 0%, ${c}1A 0%, transparent 70%),
      radial-gradient(ellipse 60% 40% at 100% 100%, ${c}0D 0%, transparent 50%),
      #0A0A0F
    `,
  }
})

const coverGlowStyle = computed(() => {
  if (!coverUrl.value) return {}
  const c = accentColor.value
  return {
    boxShadow: `0 0 40px ${c}40, 0 0 80px ${c}20, 0 0 120px ${c}10`,
    transition: 'box-shadow 0.6s ease',
  }
})

// ── Navigation ──
function isCurrentSimilarTrack(st: { id: string }): boolean {
  return player.currentTrack.value?.id === String(st.id)
}

function goBack() {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}

function goToCreateEdit() {
  router.push({ name: 'app.create-edit', query: { trackId } })
}

function openVideo(v: VideoItem) {
  router.push(`/music-video/${v.id}`)
}

// ── Playback ──
function togglePlay() {
  if (!track.value) return
  const pb = buildPlaybackTrack({
    id: trackId,
    title: track.value.title,
    artist_name: track.value.artist_name || 'Unknown',
    album_title: track.value.album_title || null,
    cover_url: track.value.cover_url || null,
    duration_seconds: track.value.duration_seconds ?? null,
    audio_url: track.value.audio_url || undefined,
  })

  if (isPlaying.value && player.currentTrack.value?.id === trackId) {
    player.pause()
  } else {
    player.playTrack(pb)
  }
}

function playSimilar(
  st: {
    id: string
    title: string
    artist_name?: string | null
    cover_url?: string | null
    duration_seconds?: number | null
  },
  index: number,
) {
  const queue = mapToPlaybackTracks(similarTracks.value.map((t) => ({
    ...t,
    audio_url: t.audio_url || undefined,
  })))

  player.setQueueAndPlay(queue, index)
}

function onSeek(event: Event) {
  const target = event.target as HTMLInputElement
  player.seekPercent(Number(target.value))
}

// ── Share ──
const { copyLink } = useSocialShare()
function shareTrack() {
  if (!track.value) return
  copyLink({
    id: trackId,
    title: track.value.title,
    type: 'track',
    artistName: track.value.artist_name,
    coverUrl: track.value.cover_url,
  })
}

// ── Formatters ──
function formatTime(seconds?: number | null) {
  if (!seconds && seconds !== 0) return '0:00'
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

function formatPlayCount(count: number) {
  if (count >= 1_000_000) return `${(count / 1_000_000).toFixed(1)}M`
  if (count >= 1_000) return `${(count / 1_000).toFixed(1)}K`
  return String(count)
}

// ── Load track videos ──
async function fetchTrackVideos() {
  try {
    const res = await videoApi.getTrackVideos(trackId)
    trackVideos.value = (res?.items || []).slice(0, 10)
  } catch {
    trackVideos.value = []
  }
}

// ── Visualizer bars ──
const barHeights = Array.from({ length: 96 }, (_, i) => {
  const center = 48
  const distance = Math.abs(i - center) / center
  return 15 + Math.cos(distance * Math.PI * 0.5) * 55 + Math.random() * 10
})

function getBarHeight(i: number) {
  return barHeights[i] ?? 20
}

onMounted(() => {
  fetchTrack()
  fetchTrackVideos()
})
</script>

<style scoped>
/* ── Now Playing bars (on cover art overlay) ── */
.now-playing-bar {
  animation: now-playing 1s ease-in-out infinite alternate;
}
@keyframes now-playing {
  0% { transform: scaleY(0.3); opacity: 0.5; }
  100% { transform: scaleY(1); opacity: 1; }
}

/* ── Visualizer bars ── */
.visualizer-bar {
  animation: viz-pulse 1.2s ease-in-out infinite alternate;
}
.visualizer-bar:nth-child(even) {
  animation-delay: 0.1s;
}
.visualizer-bar:nth-child(3n) {
  animation-delay: 0.25s;
}
.visualizer-bar:nth-child(5n+2) {
  animation-delay: 0.05s;
}
.visualizer-bar:nth-child(7n+4) {
  animation-delay: 0.35s;
}
/* ── Similar tracks equalizer bars ── */
.similar-eq-bar {
  width: 3px;
  border-radius: 999px;
  background: currentColor;
  animation: similar-equalizer 850ms ease-in-out infinite alternate;
  will-change: transform, opacity;
}
.animation-delay-150 {
  animation-delay: 150ms;
}
.animation-delay-300 {
  animation-delay: 300ms;
}
@keyframes similar-equalizer {
  from { transform: scaleY(0.45); opacity: 0.6; }
  to { transform: scaleY(1); opacity: 1; }
}

@keyframes viz-pulse {
  0% { transform: scaleY(0.4); opacity: 0.3; }
  100% { transform: scaleY(1); opacity: 0.7; }
}

/* ── Player range slider ── */
.player-range {
  --range-progress: 0%;
  width: 100%;
  height: 18px;
  cursor: pointer;
  appearance: none;
  background: transparent;
}

.player-range::-webkit-slider-runnable-track {
  height: 4px;
  border-radius: 999px;
  background: linear-gradient(
    to right,
    #1db954 0%,
    #1db954 var(--range-progress),
    rgba(255, 255, 255, 0.12) var(--range-progress),
    rgba(255, 255, 255, 0.12) 100%
  );
}

.player-range::-webkit-slider-thumb {
  width: 14px;
  height: 14px;
  margin-top: -5px;
  border-radius: 999px;
  appearance: none;
  background: #fff;
  box-shadow: 0 0 16px rgba(29, 185, 84, 0.6);
  opacity: 0;
  transition:
    opacity 160ms ease,
    transform 160ms ease;
}

.player-range:hover::-webkit-slider-thumb,
.player-range:active::-webkit-slider-thumb {
  opacity: 1;
}

.player-range:active::-webkit-slider-thumb {
  transform: scale(1.25);
}

.player-range::-moz-range-track {
  height: 4px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.12);
}

.player-range::-moz-range-progress {
  height: 4px;
  border-radius: 999px;
  background: #1db954;
}

.player-range::-moz-range-thumb {
  width: 14px;
  height: 14px;
  border: 0;
  border-radius: 999px;
  background: #fff;
  box-shadow: 0 0 16px rgba(29, 185, 84, 0.6);
  opacity: 0;
  transition:
    opacity 160ms ease,
    transform 160ms ease;
}

.player-range:hover::-moz-range-thumb,
.player-range:active::-moz-range-thumb {
  opacity: 1;
}

/* ── Horizontal scroll carousel (mobile) ── */
.scrollbar-hidden {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.scrollbar-hidden::-webkit-scrollbar {
  display: none;
}
</style>
