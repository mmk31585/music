<template>
  <div class="relative mx-auto min-h-screen w-full pb-36">
    <!-- ── Ambient Aurora Background ── -->
    <div class="aurora-bg pointer-events-none fixed inset-0" aria-hidden="true">
      <div class="aurora-spot-1" />
      <div class="aurora-spot-2" />
    </div>

    <div class="relative z-10">
      <!-- ════════════════════════════════════════ -->
      <!-- HERO — Your Music Mind                   -->
      <!-- ════════════════════════════════════════ -->
      <section class="relative overflow-hidden px-4 pt-8 md:px-6 lg:px-8">
        <div class="mx-auto max-w-7xl">
          <div
            class="relative rounded-3xl border border-white/6 bg-linear-to-br from-spotify/10 via-aurora-purple/5 to-black/40 p-8 backdrop-blur-2xl md:p-14"
          >
            <!-- Aurora overlay -->
            <div class="pointer-events-none absolute inset-0 overflow-hidden rounded-3xl">
              <div
                class="absolute -top-1/2 -right-1/4 h-96 w-96 rounded-full bg-spotify/10 blur-[120px] animate-aurora-drift"
              />
              <div
                class="absolute -bottom-1/2 -left-1/4 h-80 w-80 rounded-full bg-aurora-purple/10 blur-[100px] animate-aurora-drift-2"
              />
            </div>

            <div class="relative">
              <!-- AI Eyebrow -->
              <div class="mb-4 flex items-center gap-2">
                <span
                  class="inline-flex items-center gap-1.5 rounded-full border border-spotify/20 bg-spotify/10 px-3 py-1 text-[10px] font-bold tracking-[0.2em] text-spotify uppercase"
                >
                  <i aria-hidden="true" class="pi pi-sparkles text-[10px]" />
                  AI-Powered
                </span>
                <span class="text-[10px] font-medium text-white/40">
                  Updated just now
                </span>
              </div>

              <!-- Headline -->
              <h1
                class="text-gradient from-white via-spotify to-aurora-purple bg-clip-text text-5xl font-black leading-tight text-transparent md:text-7xl"
              >
                Your Music Mind
              </h1>

              <!-- AI Insight -->
              <p class="mt-4 max-w-2xl text-base leading-relaxed text-white/70 md:text-lg">
                <template v-if="!loadingStats && listeningStats">
                  You've listened to
                  <strong class="text-white">{{ listeningStats.total_minutes_listened }} minutes</strong>
                  this month across
                  <strong class="text-white">{{ listeningStats.unique_artists_count }} artists</strong>.
                  Your discovery score is
                  <strong class="text-spotify">{{ discoveryScore }}</strong>.
                  <template v-if="dominantMoodLabel">
                    We're sensing a <strong class="text-white">{{ dominantMoodLabel }}</strong> vibe today.
                  </template>
                </template>
                <template v-else>
                  Discover tracks curated by AI based on your unique taste and listening journey.
                </template>
              </p>

              <!-- Quick Stats -->
              <div class="mt-6 flex flex-wrap items-center gap-3">
                <div
                  v-if="!loadingStats && listeningStats"
                  class="flex items-center gap-1.5 rounded-full border border-white/6 bg-white/4 px-3 py-1.5 text-xs text-white/60"
                >
                  <i aria-hidden="true" class="pi pi-clock text-[10px]" />
                  {{ listeningStats.total_minutes_listened }} min
                </div>
                <div
                  v-if="!loadingStats && listeningStats"
                  class="flex items-center gap-1.5 rounded-full border border-white/6 bg-white/4 px-3 py-1.5 text-xs text-white/60"
                >
                  <i aria-hidden="true" class="pi pi-users text-[10px]" />
                  {{ listeningStats.unique_artists_count }} artists
                </div>
                <div
                  v-if="!loadingStats && listeningStats"
                  class="flex items-center gap-1.5 rounded-full border border-spotify/20 bg-spotify/10 px-3 py-1.5 text-xs text-spotify"
                >
                  <i aria-hidden="true" class="pi pi-chart-line text-[10px]" />
                  {{ discoveryScore }} discovery
                </div>
                <div
                  class="flex items-center gap-1.5 rounded-full border border-white/6 bg-white/4 px-3 py-1.5 text-xs text-white/60"
                >
                  <i aria-hidden="true" class="pi pi-star text-[10px]" />
                  {{ forYouTracks.length + popularTracks.length }} recommendations
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ════════════════════════════════════════ -->
      <!-- AI MOOD QUICK-PICK                       -->
      <!-- ════════════════════════════════════════ -->
      <section class="mt-10 px-4 md:px-6 lg:px-8">
        <div class="mx-auto max-w-7xl">
          <div class="mb-4 flex items-center justify-between">
            <div>
              <p class="text-[10px] font-bold tracking-[0.25em] text-spotify uppercase">
                Instant Mood Match
              </p>
              <h2 class="mt-1 text-lg font-black text-white">How are you feeling?</h2>
            </div>
            <RouterLink
              to="/ai/mood-explorer"
              class="text-xs text-white/40 transition hover:text-white"
            >
              Explore all moods
              <i aria-hidden="true" class="pi pi-arrow-right ml-1 text-[10px]" />
            </RouterLink>
          </div>

          <div
            class="no-scrollbar flex gap-2 overflow-x-auto pb-2"
            role="group"
            aria-label="Mood quick-pick"
          >
            <button
              v-for="mood in moodOptions"
              :key="mood.value"
              type="button"
              class="reveal-up group relative flex shrink-0 items-center gap-2 rounded-full border px-4 py-2.5 text-xs font-bold transition-all duration-300"
              :class="[
                selectedMood === mood.value
                  ? 'border-spotify/40 bg-spotify/15 text-white shadow-lg shadow-spotify/10'
                  : 'border-white/6 bg-white/3 text-white/60 hover:border-white/12 hover:bg-white/6 hover:text-white',
              ]"
              :style="{ transitionDelay: `${moodOptions.indexOf(mood) * 50}ms` }"
              @click="handleMoodPick(mood.value)"
            >
              <!-- Mood gradient pill background -->
              <div
                class="pointer-events-none absolute inset-0 rounded-full opacity-0 transition-opacity duration-300 group-hover:opacity-100"
                :class="getMoodGradient(mood.value)"
              />
              <span class="relative z-10 flex items-center gap-2">
                <i aria-hidden="true" :class="mood.icon" class="text-sm" />
                {{ mood.label }}
              </span>
              <span
                v-if="selectedMood === mood.value"
                class="flex h-4 w-4 items-center justify-center rounded-full bg-spotify text-[8px] text-black"
              >
                <i aria-hidden="true" class="pi pi-check" />
              </span>
            </button>
          </div>

          <!-- Inline AI Mood Results -->
          <div v-if="loadingMoodPlaylist" class="mt-4 flex items-center gap-3 rounded-2xl bg-white/2 px-6 py-4">
            <i aria-hidden="true" class="pi pi-spin pi-spinner text-spotify" />
            <span class="text-sm text-white/60">AI is curating tracks for your mood...</span>
          </div>
          <div
            v-else-if="moodPlaylist && moodPlaylist.tracks.length"
            class="mt-4 overflow-hidden rounded-2xl border border-white/6 bg-white/2"
          >
            <div class="flex items-center justify-between px-4 py-3">
              <div class="flex items-center gap-2">
                <i aria-hidden="true" class="pi pi-sparkles text-xs text-spotify" />
                <span class="text-xs font-bold text-white">{{ moodPlaylist.name }}</span>
              </div>
              <button
                type="button"
                class="text-[10px] text-white/40 transition hover:text-white"
                @click="selectedMood = ''; moodPlaylist = null"
              >
                Dismiss
              </button>
            </div>
            <div class="border-t border-white/6">
              <div
                v-for="(track, index) in moodPlaylist.tracks.slice(0, 5)"
                :key="track.id"
                role="button"
                tabindex="0"
                class="group flex cursor-pointer items-center gap-3 px-4 py-2 transition hover:bg-white/4"
                @click="playTrack(track, index)"
                @keydown.enter="playTrack(track, index)"
              >
                <div class="relative h-9 w-9 shrink-0 overflow-hidden rounded-lg bg-white/10">
                  <img
                    v-if="track.cover_url"
                    :src="track.cover_url"
                    :alt="track.title"
                    loading="lazy"
                    class="h-full w-full object-cover"
                    @error="onImgError"
                  />
                  <div v-else class="flex h-full items-center justify-center">
                    <i aria-hidden="true" class="pi pi-music text-[10px] text-white/30" />
                  </div>
                  <div
                    class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100"
                  >
                    <i aria-hidden="true" class="pi pi-play-fill text-[10px] text-white" />
                  </div>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="truncate text-sm font-medium text-white">{{ track.title }}</p>
                  <p class="truncate text-xs text-white/40">{{ track.artist || 'Unknown' }}</p>
                </div>
                <span class="shrink-0 text-xs text-white/30 tabular-nums">
                  {{ formatDuration(track.duration) }}
                </span>
              </div>
              <div
                v-if="moodPlaylist.tracks.length > 5"
                class="border-t border-white/4 px-4 py-2 text-center"
              >
                <RouterLink
                  to="/ai/mood-explorer"
                  class="text-xs text-white/30 transition hover:text-white"
                >
                  View all {{ moodPlaylist.tracks.length }} tracks
                </RouterLink>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ════════════════════════════════════════ -->
      <!-- CATEGORY CARDS — Mood Gradients          -->
      <!-- ════════════════════════════════════════ -->
      <section class="mt-10 px-4 md:px-6 lg:px-8">
        <div class="mx-auto grid max-w-7xl gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <RouterLink
            v-for="card in categoryCards"
            :key="card.to"
            :to="card.to"
            class="group relative overflow-hidden rounded-2xl border border-white/6 p-6 transition-all duration-500 hover:-translate-y-1 hover:border-white/12"
          >
            <!-- Hover gradient overlay -->
            <div
              class="pointer-events-none absolute inset-0 bg-linear-to-br opacity-0 transition-opacity duration-500 group-hover:opacity-100"
              :class="card.gradient"
            />
            <!-- Card content -->
            <div class="relative">
              <div
                class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl text-lg backdrop-blur-xs"
                :class="card.iconBg"
              >
                <i aria-hidden="true" :class="card.icon" class="relative z-10" />
              </div>
              <h3 class="text-lg font-black text-white">{{ card.title }}</h3>
              <p class="mt-1.5 text-sm leading-relaxed text-white/50">
                {{ card.description }}
              </p>
            </div>
            <!-- Arrow indicator -->
            <div
              class="absolute right-5 bottom-5 flex h-8 w-8 items-center justify-center rounded-full border border-white/6 text-xs text-white/30 transition-all duration-300 group-hover:border-white/15 group-hover:text-white/70"
            >
              <i aria-hidden="true" class="pi pi-arrow-right" />
            </div>
          </RouterLink>
        </div>
      </section>

      <!-- ════════════════════════════════════════ -->
      <!-- PERSONALIZED — AI For You Picks          -->
      <!-- ════════════════════════════════════════ -->
      <section v-if="forYouTracks.length" class="mt-14 px-4 md:px-6 lg:px-8">
        <div class="mx-auto max-w-7xl">
          <div class="mb-5 flex items-end justify-between">
            <div>
              <p class="text-[10px] font-bold tracking-[0.25em] text-spotify uppercase">
                AI Curated
              </p>
              <h2 class="mt-1 text-2xl font-black text-white">Made for You</h2>
            </div>
            <RouterLink
              to="/recommendations/for-you"
              class="text-xs text-white/40 transition hover:text-white"
            >
              View all
            </RouterLink>
          </div>

          <!-- Loading -->
          <div v-if="loadingForYou" class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <div v-for="i in 5" :key="i" class="space-y-3">
              <div class="shimmer aspect-square rounded-2xl" />
              <div class="shimmer h-4 w-3/4 rounded-lg" />
              <div class="shimmer h-3 w-1/2 rounded-lg" />
            </div>
          </div>

          <!-- Track Grid -->
          <div v-else class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <div
              v-for="(track, idx) in forYouTracks"
              :key="track.id"
              role="button"
              tabindex="0"
              class="group cursor-pointer"
              :style="{ animationDelay: `${idx * 80}ms` }"
              @click="playTrack(track, idx)"
              @keydown.enter="playTrack(track, idx)"
            >
              <div
                class="relative mb-3 aspect-square overflow-hidden rounded-2xl bg-white/5 shadow-lg ring-1 ring-white/6 transition-all duration-500 group-hover:scale-[1.02] group-hover:ring-spotify/30"
              >
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-music text-2xl text-white/20" />
                </div>
                <!-- Play overlay -->
                <div
                  class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 backdrop-blur-xs transition-all duration-300 group-hover:opacity-100"
                >
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-full bg-spotify/90 text-black shadow-xl shadow-spotify/20 transition-transform duration-300 group-hover:scale-105"
                  >
                    <i aria-hidden="true" class="pi pi-play-fill text-lg" />
                  </div>
                </div>
                <!-- AI Reason Chip -->
                <div
                  v-if="track.genre"
                  class="absolute top-2 left-2 rounded-full bg-black/60 px-2 py-0.5 text-[9px] font-medium text-white/70 backdrop-blur-xs"
                >
                  {{ track.genre }}
                </div>
              </div>
              <p class="truncate text-sm font-semibold text-white">{{ track.title }}</p>
              <p class="truncate text-xs text-white/40">{{ track.artist_name || 'Unknown' }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- ════════════════════════════════════════ -->
      <!-- SIMILAR — Because You Listened To        -->
      <!-- ════════════════════════════════════════ -->
      <section
        v-if="personalizedTracks.length"
        class="mt-14 px-4 md:px-6 lg:px-8"
      >
        <div class="mx-auto max-w-7xl">
          <div class="mb-5 flex items-end justify-between">
            <div>
              <p class="text-[10px] font-bold tracking-[0.25em] text-aurora-purple uppercase">
                AI Similarity
              </p>
              <h2 class="mt-1 text-2xl font-black text-white">Because you listened to...</h2>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
            <div
              v-for="(track, idx) in personalizedTracks.slice(0, 5)"
              :key="track.id"
              role="button"
              tabindex="0"
              class="group cursor-pointer"
              :style="{ animationDelay: `${idx * 80}ms` }"
              @click="playTrack(track, idx)"
              @keydown.enter="playTrack(track, idx)"
            >
              <div
                class="relative mb-3 aspect-square overflow-hidden rounded-2xl bg-white/5 shadow-lg ring-1 ring-white/6 transition-all duration-500 group-hover:scale-[1.02] group-hover:ring-aurora-purple/30"
              >
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-music text-2xl text-white/20" />
                </div>
                <div
                  class="absolute inset-0 flex items-center justify-center bg-black/40 opacity-0 backdrop-blur-xs transition-all duration-300 group-hover:opacity-100"
                >
                  <div
                    class="flex h-12 w-12 items-center justify-center rounded-full bg-aurora-purple/90 text-white shadow-xl shadow-aurora-purple/20 transition-transform duration-300 group-hover:scale-105"
                  >
                    <i aria-hidden="true" class="pi pi-play-fill text-lg" />
                  </div>
                </div>
                <div
                  class="absolute top-2 left-2 rounded-full bg-aurora-purple/30 px-2 py-0.5 text-[9px] font-medium text-white/80 backdrop-blur-xs"
                >
                  <i aria-hidden="true" class="pi pi-bolt mr-0.5 text-[8px]" />
                  Similar
                </div>
              </div>
              <p class="truncate text-sm font-semibold text-white">{{ track.title }}</p>
              <p class="truncate text-xs text-white/40">{{ track.artist_name || 'Unknown' }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- ════════════════════════════════════════ -->
      <!-- DISCOVER WEEKLY — AI Weekly Playlist     -->
      <!-- ════════════════════════════════════════ -->
      <section
        v-if="!loadingDiscoverWeekly && discoverWeekly"
        class="mt-14 px-4 md:px-6 lg:px-8"
      >
        <div class="mx-auto max-w-7xl">
          <div
            class="relative overflow-hidden rounded-2xl border border-white/6 bg-linear-to-br from-aurora-purple/10 via-spotify/5 to-black/40 p-8 backdrop-blur-2xl md:p-10"
          >
            <div class="pointer-events-none absolute inset-0 overflow-hidden rounded-2xl">
              <div
                class="absolute -top-1/3 -left-1/4 h-64 w-64 rounded-full bg-aurora-purple/10 blur-[100px]"
              />
              <div
                class="absolute -bottom-1/3 -right-1/4 h-64 w-64 rounded-full bg-spotify/10 blur-[100px]"
              />
            </div>

            <div class="relative flex flex-col gap-6 md:flex-row md:items-center md:justify-between">
              <div class="flex-1">
                <div class="mb-3 flex items-center gap-2">
                  <span
                    class="inline-flex items-center gap-1.5 rounded-full border border-aurora-purple/20 bg-aurora-purple/10 px-3 py-1 text-[10px] font-bold tracking-[0.2em] text-aurora-purple uppercase"
                  >
                    <i aria-hidden="true" class="pi pi-calendar text-[10px]" />
                    Weekly
                  </span>
                </div>
                <h2 class="text-2xl font-black text-white md:text-3xl">
                  Discover Weekly
                </h2>
                <p class="mt-2 max-w-lg text-sm text-white/60">
                  AI-curated playlist updated every Monday.
                  {{ discoverWeekly.playlist.track_count }} fresh tracks based on your taste.
                </p>
              </div>
              <RouterLink
                to="/recommendations/for-you"
                class="inline-flex items-center gap-2 rounded-full bg-white/10 px-6 py-3 text-sm font-bold text-white transition hover:bg-white/20"
              >
                <i aria-hidden="true" class="pi pi-play" />
                Listen now
              </RouterLink>
            </div>

            <!-- Preview of discover weekly tracks -->
            <div class="relative mt-6">
              <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
                <div
                  v-for="(track, idx) in discoverWeekly.tracks.slice(0, 5)"
                  :key="track.id"
                  role="button"
                  tabindex="0"
                  class="group cursor-pointer"
                  @click="playTrack(track, idx)"
                  @keydown.enter="playTrack(track, idx)"
                >
                  <div
                    class="relative mb-2 aspect-square overflow-hidden rounded-xl bg-white/5 ring-1 ring-white/6 transition-all duration-300 group-hover:ring-aurora-purple/30"
                  >
                    <img
                      v-if="track.cover_url"
                      :src="track.cover_url"
                      :alt="track.title"
                      loading="lazy"
                      class="h-full w-full object-cover transition duration-300 group-hover:scale-105"
                      @error="onImgError"
                    />
                    <div v-else class="flex h-full items-center justify-center">
                      <i aria-hidden="true" class="pi pi-music text-xl text-white/20" />
                    </div>
                    <div
                      class="absolute inset-0 flex items-center justify-center bg-black/30 opacity-0 transition group-hover:opacity-100"
                    >
                      <div
                        class="flex h-10 w-10 items-center justify-center rounded-full bg-aurora-purple/80 text-white"
                      >
                        <i aria-hidden="true" class="pi pi-play-fill text-sm" />
                      </div>
                    </div>
                  </div>
                  <p class="truncate text-xs font-medium text-white/80">{{ track.title }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ════════════════════════════════════════ -->
      <!-- POPULAR — Trending Now                   -->
      <!-- ════════════════════════════════════════ -->
      <section class="mt-14 px-4 md:px-6 lg:px-8">
        <div class="mx-auto max-w-7xl">
          <div class="mb-5 flex items-end justify-between">
            <div>
              <p class="text-[10px] font-bold tracking-[0.25em] text-white/30 uppercase">
                Trending
              </p>
              <h2 class="mt-1 text-2xl font-black text-white">Popular Now</h2>
            </div>
            <RouterLink
              to="/recommendations/popular"
              class="text-xs text-white/40 transition hover:text-white"
            >
              View all
            </RouterLink>
          </div>

          <div v-if="loadingPopular" class="space-y-2">
            <div v-for="i in 5" :key="i" class="shimmer h-16 rounded-2xl" />
          </div>
          <div
            v-else-if="popularTracks.length === 0"
            class="rounded-2xl border border-dashed border-white/6 px-6 py-12 text-center"
          >
            <p class="text-sm text-white/40">No popular tracks yet.</p>
          </div>
          <div
            v-else
            class="overflow-hidden rounded-2xl border border-white/6 bg-white/2 backdrop-blur-xs"
          >
            <div
              v-for="(track, index) in popularTracks"
              :key="track.id"
              role="button"
              tabindex="0"
              :style="{ animationDelay: `${index * 60}ms` }"
              class="group flex cursor-pointer items-center gap-4 px-4 py-3 transition hover:bg-white/4"
              @click="playTrack(track, index)"
              @keydown.enter="playTrack(track, index)"
            >
              <!-- Number / Play icon -->
              <span class="flex w-7 items-center justify-center">
                <span class="text-sm font-bold text-white/20 group-hover:hidden">{{ index + 1 }}</span>
                <i
                  aria-hidden="true"
                  class="pi pi-play-fill hidden text-sm text-white group-hover:block"
                />
              </span>

              <!-- Cover -->
              <div
                class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/5 ring-1 ring-white/6"
              >
                <img
                  v-if="track.cover_url"
                  :src="track.cover_url"
                  :alt="track.title"
                  loading="lazy"
                  class="h-full w-full object-cover"
                  @error="onImgError"
                />
                <div v-else class="flex h-full items-center justify-center">
                  <i aria-hidden="true" class="pi pi-music text-xs text-white/20" />
                </div>
              </div>

              <!-- Info -->
              <div class="min-w-0 flex-1">
                <p class="truncate font-semibold text-white">{{ track.title }}</p>
                <p class="truncate text-sm text-white/40">
                  {{ track.artist_name || 'Unknown artist' }}
                </p>
              </div>

              <!-- Genre badge -->
              <span
                v-if="track.genre"
                class="hidden rounded-full bg-white/6 px-2.5 py-1 text-[10px] font-medium text-white/40 md:block"
              >
                {{ track.genre }}
              </span>

              <!-- Score indicator -->
              <span
                v-if="track.score"
                class="hidden items-center gap-1 text-xs text-white/30 sm:flex"
              >
                <i aria-hidden="true" class="pi pi-chart-line text-[10px]" />
                {{ Math.round(track.score) }}
              </span>

              <!-- Duration -->
              <span class="shrink-0 text-xs text-white/30 tabular-nums">
                {{ formatDuration(track.duration_seconds) }}
              </span>
            </div>
          </div>
        </div>
      </section>

      <!-- ════════════════════════════════════════ -->
      <!-- BOTTOM SPACER for player                 -->
      <!-- ════════════════════════════════════════ -->
      <div class="h-8" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAIRecommendations } from '@/composables/useAIRecommendations'
import { usePlayer } from '@/composables/player'
import { usePlayerApi } from '@/services/api/player'
import { MOOD_OPTIONS } from '@/services/api/ai/types'
import type { RecommendationTrack } from '@/services/api/recommendation/types'
import type { AITrackItem } from '@/services/api/ai/types'
import { onImgError } from '@/utils/helpers'

const router = useRouter()
const player = usePlayer()
const playerApi = usePlayerApi()

const {
  popularTracks,
  forYouTracks,
  personalizedTracks,
  discoverWeekly,
  listeningStats,
  moodPlaylist,
  loadingPopular,
  loadingForYou,
  loadingDiscoverWeekly,
  loadingStats,
  loadingMoodPlaylist,
  fetchPopular,
  fetchForYou,
  fetchDiscoverWeekly,
  fetchListeningStats,
  fetchPersonalized,
  fetchMoodPlaylist,
  formatDuration,
  getMoodGradient,
  dominantMood,
  discoveryScore,
} = useAIRecommendations()

const selectedMood = ref('')

const moodOptions = [...MOOD_OPTIONS]

const dominantMoodLabel = computed(() => {
  const mood = moodOptions.find((m) => m.value === dominantMood.value)
  return mood?.label || null
})

const categoryCards = [
  {
    title: 'For You',
    description: 'AI-curated picks based on your unique taste and listening history.',
    to: '/recommendations/for-you',
    icon: 'pi pi-heart',
    iconBg: 'bg-pink-500/20 text-pink-400',
    gradient: 'from-pink-500/10 via-rose-500/5 to-transparent',
  },
  {
    title: 'Mood Explorer',
    description: 'Browse tracks by mood — energy, valence, tempo, and more.',
    to: '/ai/mood-explorer',
    icon: 'pi pi-magic',
    iconBg: 'bg-blue-500/20 text-blue-400',
    gradient: 'from-blue-500/10 via-cyan-500/5 to-transparent',
  },
  {
    title: 'AI Playlist',
    description: 'Generate a custom playlist from a prompt, mood, or activity.',
    to: '/ai/playlist-generator',
    icon: 'pi pi-sparkles',
    iconBg: 'bg-amber-500/20 text-amber-400',
    gradient: 'from-amber-500/10 via-orange-500/5 to-transparent',
  },
  {
    title: 'Popular',
    description: 'The most played tracks trending across the catalog right now.',
    to: '/recommendations/popular',
    icon: 'pi pi-chart-line',
    iconBg: 'bg-spotify/20 text-spotify',
    gradient: 'from-spotify/10 via-emerald-500/5 to-transparent',
  },
  {
    title: 'Best Tracks',
    description: 'Curated high-quality tracks for the best listening session.',
    to: '/recommendations/best',
    icon: 'pi pi-star',
    iconBg: 'bg-yellow-400/20 text-yellow-300',
    gradient: 'from-yellow-400/10 via-amber-500/5 to-transparent',
  },
  {
    title: 'Recent',
    description: 'Freshly added tracks to expand your musical horizons.',
    to: '/recommendations/recent',
    icon: 'pi pi-clock',
    iconBg: 'bg-sky-400/20 text-sky-300',
    gradient: 'from-sky-400/10 via-blue-500/5 to-transparent',
  },
]

async function handleMoodPick(mood: string) {
  if (selectedMood.value === mood) {
    selectedMood.value = ''
    moodPlaylist.value = null
    return
  }
  selectedMood.value = mood
  await fetchMoodPlaylist(mood)
}

function playTrack(track: RecommendationTrack | Record<string, unknown>, index: number) {
  const allTracks = [...forYouTracks.value, ...popularTracks.value, ...personalizedTracks.value]
  const source = track.id ? forYouTracks.value.find(t => t.id === track.id)
    ? forYouTracks.value
    : popularTracks.value.find(t => t.id === track.id)
      ? popularTracks.value
      : personalizedTracks.value
    : allTracks

  const queue = source.map((t: any) => ({
    id: String(t.id),
    title: String(t.title || ''),
    artistName: String(t.artist_name || t.artist || 'Unknown'),
    albumTitle: String(t.album_title || t.album || ''),
    coverUrl: String(t.cover_url || ''),
    durationSeconds: t.duration_seconds ?? t.duration ?? null,
    streamUrl: playerApi.getTrackStreamUrl(String(t.id)),
  }))
  const startIdx = queue.findIndex((t) => t.id === String(track.id || track.id))
  if (startIdx >= 0) {
    player.setQueueAndPlay(queue, startIdx)
  }
}

onMounted(async () => {
  await Promise.all([
    fetchPopular(8),
    fetchForYou(5),
    fetchPersonalized(5),
  ])
  fetchDiscoverWeekly()
  fetchListeningStats('month')
})
</script>
