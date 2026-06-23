<template>
  <Teleport to="body">
    <Transition name="fullscreen">
      <div v-if="isOpen" class="fixed inset-0 z-[9999] flex flex-col" :style="{ ...dynamicBg, transition: 'background 0.8s cubic-bezier(0.19, 1, 0.22, 1)' }">

        <!-- Background blur layer with crossfade -->
        <div class="pointer-events-none absolute -inset-5 z-0 scale-110">
          <img
            v-if="currentTrack?.coverUrl"
            :src="currentTrack?.coverUrl"
            class="h-full w-full object-cover opacity-40 md:opacity-50 transition-all duration-1000"
            style="filter: blur(100px) saturate(2)"
          />
        </div>
        <div class="pointer-events-none absolute inset-0 z-[1] bg-gradient-to-b from-black/70 via-black/30 to-black/90" />

        <!-- ── Top bar ── -->
        <div class="relative z-10 flex shrink-0 items-center justify-between px-4 pt-3 md:px-8 md:pt-5" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
          <button class="flex h-10 w-10 items-center justify-center rounded-full text-white/50 transition-all hover:text-white hover:bg-white/10" aria-label="Close" @click="close">
            <i aria-hidden="true" class="pi pi-chevron-down text-xl" />
          </button>
          <p class="text-[10px] font-semibold tracking-[0.2em] text-white/30 uppercase">Now Playing</p>
          <div class="flex items-center gap-1">
            <!-- Toggle: lyrics/wave vs queue -->
            <button
              class="flex h-10 w-10 items-center justify-center rounded-full transition-all"
              :class="rightPanelTab === 'queue' ? 'bg-white/15 text-white' : 'text-white/50 hover:text-white hover:bg-white/10'"
              :aria-label="rightPanelTab === 'queue' ? 'Show lyrics' : 'Show queue'"
              @click="rightPanelTab = rightPanelTab === 'queue' ? 'lyrics' : 'queue'"
            >
              <i aria-hidden="true" class="pi pi-list text-lg" />
            </button>
            <button class="flex h-10 w-10 items-center justify-center rounded-full text-white/50 transition-all hover:text-white hover:bg-white/10" aria-label="Close" @click="emit('toggle-queue')">
              <i aria-hidden="true" class="pi pi-external-link text-sm" />
            </button>
          </div>
        </div>

        <!-- ── Desktop: side-by-side (md+) ── -->
        <div class="relative z-10 hidden md:flex flex-1 min-h-0 px-8 pb-6">
          <div class="flex flex-1 gap-8 min-h-0 w-full max-w-7xl mx-auto">

            <!-- LEFT: Cover + info + controls (slightly narrower to make room for queue preview) -->
            <div class="flex flex-col items-center gap-4 shrink-0 w-[380px] lg:w-[420px] justify-center pb-12">

              <!-- Album art -->
              <div class="relative">
                <div
                  class="w-[240px] lg:w-[280px] aspect-square overflow-hidden rounded-3xl shadow-2xl transition-all duration-700"
                  :class="isPlaying ? 'scale-100 cover-glow' : 'scale-95 opacity-80'"
                >
                  <img
                    v-if="currentTrack?.coverUrl"
                    :src="currentTrack?.coverUrl"
                    :alt="currentTrack?.title"
                    class="h-full w-full object-cover"
                  />
                  <div
                    v-else
                    class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#a855f7]/30"
                  >
                    <i aria-hidden="true" class="pi pi-music text-5xl text-white/30" />
                  </div>
                </div>
                <div
                  v-if="isPlaying"
                  class="pointer-events-none absolute -inset-3 animate-pulse rounded-full border-2"
                  :style="{ borderColor: `${palette.vibrant}33` }"
                />
              </div>

              <!-- Track info -->
              <div class="w-full max-w-sm text-center">
                <div class="flex items-center justify-center gap-3">
                  <div class="min-w-0">
                    <h2 class="text-xl lg:text-2xl font-bold text-white truncate max-w-70">{{ currentTrack?.title }}</h2>
                    <p class="text-sm lg:text-base text-white/50 mt-0.5 truncate max-w-70">{{ currentTrack?.artistName }}</p>
                  </div>

                </div>
              </div>

              <!-- Progress bar (dynamic color) with seek preview -->
              <div class="w-full max-w-sm">
                <div
                  role="slider"
                  tabindex="0"
                  ref="progressRef"
                  class="relative flex h-8 cursor-pointer items-center group"
                  aria-label="Seek"
                  aria-valuemin="0"
                  :aria-valuemax="duration"
                  :aria-valuenow="currentTime"
                  @click="seek"
                  @mousedown="startDrag"
                  @keydown.enter="seek"
                  @keydown.space.prevent="seek"
                  @mousemove="onSeekHover"
                  @mouseleave="seekHoverTime = null"
                >
                  <div class="absolute inset-x-0 h-1 rounded-full transition-all duration-150 group-hover:h-1.5" :class="isPlaying ? 'bg-white/15' : 'bg-white/10'" />
                  <div v-if="isPlaying" class="absolute left-0 h-1 rounded-full opacity-25 blur-[3px] transition-all duration-150" :style="{ width: progressPercent + '%', background: progressColor }" />
                  <div class="absolute left-0 h-1 rounded-full transition-all duration-150 group-hover:h-1.5" :style="{ width: progressPercent + '%', background: progressColor, boxShadow: `0 0 8px ${progressColor}66` }" />
                  <div class="absolute top-1/2 -translate-y-1/2 size-4 rounded-full shadow-xl opacity-0 group-hover:opacity-100 scale-0 group-hover:scale-100 transition-all duration-300 ease-spring" :style="{ left: `calc(${progressPercent}% - 8px)`, background: progressColor, boxShadow: `0 0 0 3px ${progressColor}22, 0 4px 12px rgba(0,0,0,0.5)` }" />
                  <!-- Seek preview tooltip -->
                  <div
                    v-if="seekHoverTime !== null && seekHoverPos !== null"
                    class="absolute -top-7 z-20 rounded-md bg-[#1a1a1a] px-2 py-1 text-[11px] font-mono tabular-nums text-white shadow-xl border border-white/10 transition-opacity duration-100 pointer-events-none"
                    :style="{ left: `${seekHoverPos}px`, transform: 'translateX(-50%)' }"
                  >
                    {{ formatTime(seekHoverTime) }}
                  </div>
                </div>
                <div class="flex justify-between mt-1">
                  <span class="text-[11px] text-white/40 font-mono tabular-nums">{{ formatTime(currentTime) }}</span>
                  <span class="text-[11px] text-white/40 font-mono tabular-nums">{{ formatTime(duration) }}</span>
                </div>
              </div>

              <!-- Playback controls -->
              <div class="flex items-center justify-center gap-1 w-full max-w-sm">
                <div class="flex shrink-0 items-center gap-1">
                  <button class="flex h-9 w-9 items-center justify-center transition-all hover:scale-110" :class="liked ? 'text-[#f472b6]' : 'text-white/40 hover:text-white'" :aria-label="liked ? 'Unlike' : 'Like'" @click="toggleLike">
                    <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-base" />
                  </button>
                  <button class="flex h-9 w-9 items-center justify-center transition-all hover:scale-110 text-white/40 hover:text-white" aria-label="Add to playlist" @click="showAddToPlaylist = true">
                    <i aria-hidden="true" class="pi pi-list-plus text-base" />
                  </button>
                </div>
                <div class="relative">
                  <button
                    type="button"
                    class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200 hover:scale-110"
                    :class="shuffleMode !== 'off' ? 'text-[#a855f7]' : 'text-white/40 hover:text-white'"
                    :aria-label="shuffleMode === 'queue' ? 'Shuffle queue' : shuffleMode === 'catalog' ? 'Random catalog tracks' : shuffleMode === 'similar' ? 'Similar tracks' : 'Shuffle off'"
                    data-shuffle-btn
                    @click="openShuffleMenu($event)"
                  >
                    <i aria-hidden="true" class="pi pi-sort-alt text-sm" />
                    <span v-if="shuffleMode !== 'off'" class="absolute -top-0.5 -right-0.5 flex h-3 w-3 items-center justify-center rounded-full bg-[#a855f7] text-[7px] font-bold text-white">{{ shuffleMode === 'queue' ? 'Q' : shuffleMode === 'catalog' ? 'R' : 'S' }}</span>
                  </button>
                </div>
                <button class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all duration-200 hover:text-white hover:scale-110 disabled:opacity-20" :disabled="!hasPrevious" aria-label="Previous track" @click="playPrevious">
                  <i aria-hidden="true" class="pi pi-step-backward text-lg" />
                </button>
                <button class="flex h-14 w-14 items-center justify-center rounded-full shadow-2xl transition-all duration-200 active:scale-95 hover:scale-105 disabled:opacity-40" :style="{ background: progressColor }" :disabled="!currentTrack || isLoadingTrack" :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'" @click="togglePlay">
                  <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-xl text-white" />
                  <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause' : 'pi pi-play'" class="text-xl ml-1 text-white" />
                </button>
                <button class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all duration-200 hover:text-white hover:scale-110 disabled:opacity-20" :disabled="!hasNext" aria-label="Next track" @click="playNext">
                  <i aria-hidden="true" class="pi pi-step-forward text-lg" />
                </button>
                <div class="relative">
                  <button
                    type="button"
                    class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200 hover:scale-110"
                    :class="repeatMode !== 'off' ? 'text-[#f472b6]' : 'text-white/40 hover:text-white'"
                    :aria-label="repeatMode === 'off' ? 'Repeat off' : repeatMode === 'all' ? 'Repeat all' : 'Repeat one'"
                    data-repeat-btn
                    @click="openRepeatMenu($event)"
                  >
                    <i aria-hidden="true" class="pi pi-refresh text-sm" />
                    <span v-if="repeatMode === 'one'" class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-[#f472b6] text-[7px] font-bold text-black">1</span>
                  </button>
                </div>
              </div>

              <!-- Volume -->
              <div class="flex items-center gap-2 w-full max-w-65 mx-auto">
                <button class="flex h-8 w-8 items-center justify-center text-white/40 transition-colors hover:text-white shrink-0" :aria-label="muted ? 'Unmute' : 'Mute'" @click="toggleMute">
                  <i aria-hidden="true" :class="muted ? 'pi pi-volume-off' : 'pi pi-volume-up'" class="text-xs" />
                </button>
                <div class="relative flex-1 flex items-center group/vol h-4" dir="ltr">
                  <div class="absolute inset-x-0 h-0.5 rounded-full bg-white/15" />
                  <div class="absolute left-0 h-0.5 rounded-full transition-all" :style="{ width: `${muted ? 0 : Number(volume) * 100}%`, background: progressColor }" />
                  <input type="range" min="0" max="1" step="0.01" class="absolute inset-0 w-full cursor-pointer opacity-0 z-10"  :value="muted ? 0 : volume" @input="onVolume" />
                </div>
                <button class="flex h-8 w-8 items-center justify-center text-white/40 transition-colors hover:text-white shrink-0" aria-label="Toggle lyrics" @click="emit('toggle-lyrics')">
                  <i aria-hidden="true" class="pi pi-align-left text-xs" />
                </button>
              </div>
            </div>

            <!-- RIGHT: Toggle between Lyrics/Wave and Queue (desktop only) -->
            <div class="flex-1 min-h-0 flex flex-col md:pl-4">
              <!-- Panel tab bar -->
              <div class="flex items-center gap-3 mb-3 shrink-0">
                <button
                  class="text-xs font-bold tracking-wider uppercase transition"
                  :class="rightPanelTab === 'lyrics' ? 'text-white' : 'text-white/30 hover:text-white/60'"
                  @click="rightPanelTab = 'lyrics'"
                >
                  <i aria-hidden="true" class="pi pi-align-left me-1.5 text-[10px]" />
                  Lyrics
                </button>
                <button
                  v-if="queuePreview.length > 0"
                  class="text-xs font-bold tracking-wider uppercase transition"
                  :class="rightPanelTab === 'queue' ? 'text-white' : 'text-white/30 hover:text-white/60'"
                  @click="rightPanelTab = 'queue'"
                >
                  <i aria-hidden="true" class="pi pi-list me-1.5 text-[10px]" />
                  Up Next
                </button>
              </div>

              <div class="flex-1 min-h-0 overflow-hidden relative">

                <!-- ── LYRICS / WAVE tab ── -->
                <div v-show="rightPanelTab === 'lyrics'" class="h-full">
                  <!-- Lyrics loading shimmer -->
                  <div v-if="lyricsLoading" class="flex h-full flex-col items-center justify-center gap-3">
                    <div class="shimmer h-4 w-48 rounded bg-white/[0.06]" />
                    <div class="shimmer h-4 w-36 rounded bg-white/[0.04]" />
                    <div class="shimmer h-4 w-40 rounded bg-white/[0.03]" />
                  </div>

                  <!-- Synced lyrics -->
                  <SyncedLyrics
                    v-else-if="parsedLines.length > 0"
                    :lines="parsedLines"
                    :current-time="currentTime"
                    :duration="duration"
                    :active-color="palette.vibrant"
                    :muted-color="palette.muted"
                    :language="lyricsLanguage"
                    @seek="seekTo"
                  />

                  <!-- Simple wave visualizer when no lyrics -->
                  <div v-else class="flex h-full flex-col items-center justify-center gap-6">
                    <div class="wave-visualizer flex items-end gap-[3px] h-16" :style="{ '--wave-color': palette.vibrant }">
                      <span v-for="i in 32" :key="i" class="wave-bar w-[3px] rounded-full" :style="{ animationDelay: `${i * 0.08}s`, height: `${20 + Math.sin(i * 0.7) * 30 + 30}%`, background: palette.vibrant || '#1db954' }" />
                    </div>
                    <p class="text-xs text-white/25">No synced lyrics available</p>
                  </div>
                </div>

                <!-- ── QUEUE tab ── -->
                <div v-show="rightPanelTab === 'queue'" class="h-full flex flex-col px-1">
                  <div v-if="queuePreview.length > 0" class="flex-1 space-y-1 overflow-y-auto scroll-thin">
                    <div
                      v-for="(track, i) in queuePreview" :key="track.id || i"
                      class="group flex items-center gap-3 rounded-xl p-2 transition hover:bg-white/[0.06] cursor-pointer"
                      @click="playQueueItem(i)"
                    >
                      <div class="h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10">
                        <img
                          v-if="track.coverUrl"
                          :src="track.coverUrl"
                          :alt="track.title"
                          class="h-full w-full object-cover"
                          loading="lazy"
                        />
                        <div v-else class="flex h-full items-center justify-center">
                          <i aria-hidden="true" class="pi pi-music text-xs text-white/30" />
                        </div>
                      </div>
                      <div class="min-w-0 flex-1">
                        <p class="truncate text-sm font-medium text-white group-hover:text-[#1db954] transition-colors">
                          {{ track.title }}
                        </p>
                        <p class="truncate text-xs text-white/40">{{ track.artistName }}</p>
                      </div>
                      <span class="text-[10px] text-white/25 tabular-nums shrink-0">
                        {{ formatTime(track.durationSeconds ?? 0) }}
                      </span>
                    </div>
                  </div>
                  <div v-else class="flex h-full items-center justify-center">
                    <p class="text-xs text-white/20">Queue is empty</p>
                  </div>
                </div>

              </div>
            </div>

          </div>
        </div>

        <!-- ── Mobile: single column stacked (<md) ── -->
        <div
          class="relative z-10 flex md:hidden flex-1 min-h-0 flex-col"
          ref="mobileContainerRef"
          @touchstart="onSwipeStart"
          @touchmove.prevent="onSwipeMove"
          @touchend="onSwipeEnd"
          @touchcancel="onSwipeEnd"
        >
          <!-- Drag handle for pull-to-close -->
          <div class="flex justify-center pt-2 pb-0 shrink-0">
            <div
              class="h-1 w-10 rounded-full transition-all"
              :style="{ background: swipeProgress > 0.3 ? `${palette.vibrant}66` : 'rgba(255,255,255,0.2)' }"
            />
          </div>

          <!-- Scrollable middle area: cover → lyrics → controls; flex-col justify-center for vertical centering -->
          <div class="flex-1 overflow-y-auto px-4 pb-4 scroll-smooth flex flex-col justify-center" :style="{ transform: `translateY(${swipeProgress * 80}px)`, opacity: 1 - swipeProgress, transition: isSwiping ? 'none' : 'transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.4s ease' }" style="scrollbar-width: thin;">

            <!-- Album art -->
            <div class="flex flex-col items-center pt-4">
              <div class="relative">
                <div
                  class="w-[min(200px,45vw)] aspect-square overflow-hidden rounded-2xl shadow-2xl transition-all duration-700"
                  :class="isPlaying ? 'scale-100 cover-glow' : 'scale-95 opacity-80'"
                >
                  <img v-if="currentTrack?.coverUrl" :src="currentTrack?.coverUrl" :alt="currentTrack?.title" class="h-full w-full object-cover" />
                  <div v-else class="flex h-full w-full items-center justify-center bg-gradient-to-br from-[#1db954]/30 to-[#a855f7]/30">
                    <i aria-hidden="true" class="pi pi-music text-4xl text-white/30" />
                  </div>
                </div>
                <div v-if="isPlaying" class="pointer-events-none absolute -inset-3 animate-pulse rounded-full border-2" :style="{ borderColor: `${palette.vibrant}33` }" />
              </div>
            </div>

            <!-- Track info -->
            <div class="w-full text-center mt-4">
              <div class="flex items-center justify-center gap-2">
                <div class="min-w-0">
                  <h2 class="text-base font-bold text-white truncate max-w-[220px]">{{ currentTrack?.title }}</h2>
                  <p class="text-sm text-white/50 mt-0.5 truncate max-w-[220px]">{{ currentTrack?.artistName }}</p>
                </div>
                <div class="flex shrink-0 items-center gap-0.5">
                  <button class="flex h-8 w-8 items-center justify-center transition-all hover:scale-110" :class="liked ? 'text-[#f472b6]' : 'text-white/40 hover:text-white'" :aria-label="liked ? 'Unlike' : 'Like'" @click="toggleLike">
                    <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-sm" />
                  </button>
                </div>
              </div>
            </div>

            <!-- ── Mobile lyrics (between cover and controls) ── -->
            <div class="w-full mt-4 px-2 min-h-[60px] flex items-center justify-center">
              <div v-if="lyricsLoading" class="shimmer h-5 w-56 rounded bg-white/[0.06]" />
              <template v-else-if="activeLineText">
                <div class="transition-all duration-700 ease-out text-center">
                  <p class="text-lg font-bold leading-relaxed transition-all duration-500" :style="{ color: palette.vibrant, textShadow: `0 0 40px ${palette.vibrant}44` }">{{ activeLineText }}</p>
                </div>
                <div v-if="parsedLines.length > 1" class="absolute bottom-1 left-1/2 -translate-x-1/2 flex items-center gap-1">
                  <span v-for="dot in Math.min(parsedLines.length, 7)" :key="dot" class="h-1 rounded-full transition-all duration-300" :class="dot - 1 === activeLineIdx ? 'w-3 opacity-80' : 'w-1 opacity-20'" :style="{ background: dot - 1 === activeLineIdx ? palette.vibrant : '#fff' }" />
                </div>
              </template>
              <template v-else>
                <p class="text-sm text-white/15 text-center">No synced lyrics</p>
              </template>
            </div>

            <!-- Progress bar (mobile) -->
            <div class="w-full mt-3 px-2">
              <div
                role="slider"
                tabindex="0"
                ref="progressRefMobile"
                class="relative flex h-8 cursor-pointer items-center group"
                aria-label="Seek"
                aria-valuemin="0"
                :aria-valuemax="duration"
                :aria-valuenow="currentTime"
                @click="seek"
                @touchstart.prevent="startTouchDrag"
                @keydown.enter="seek"
                @keydown.space.prevent="seek"
                @mousemove="onSeekHoverMobile"
                @mouseleave="mobileSeekHoverTime = null"
              >
                <div class="absolute inset-x-0 h-1 rounded-full bg-white/12" />
                <div class="absolute left-0 h-1 rounded-full" :style="{ width: progressPercent + '%', background: progressColor, boxShadow: `0 0 8px ${progressColor}66` }" />
                <div class="absolute top-1/2 -translate-y-1/2 size-4 rounded-full shadow-xl" :style="{ left: `calc(${progressPercent}% - 8px)`, background: progressColor, boxShadow: `0 0 0 3px ${progressColor}22` }" />
                <!-- Seek preview tooltip (mobile) -->
                <div
                  v-if="mobileSeekHoverTime !== null && mobileSeekHoverPos !== null"
                  class="absolute -top-8 z-20 rounded-md bg-[#1a1a1a] px-2 py-1 text-[10px] font-mono tabular-nums text-white shadow-xl border border-white/10 transition-opacity duration-100 pointer-events-none"
                  :style="{ left: `${mobileSeekHoverPos}px`, transform: 'translateX(-50%)' }"
                >
                  {{ formatTime(mobileSeekHoverTime) }}
                </div>
              </div>
              <div class="flex justify-between mt-1 px-0.5">
                <span class="text-[10px] text-white/40 font-mono tabular-nums">{{ formatTime(currentTime) }}</span>
                <span class="text-[10px] text-white/40 font-mono tabular-nums">{{ formatTime(duration) }}</span>
              </div>
            </div>

            <!-- Playback controls (mobile) -->
            <div class="flex items-center justify-center gap-2 mt-3">
              <button data-shuffle-btn class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all hover:text-white active:scale-90" :class="shuffleMode !== 'off' ? 'text-[#a855f7]' : ''" aria-label="Shuffle" @click="openShuffleMenu($event)">
                <i aria-hidden="true" class="pi pi-sort-alt text-sm" />
              </button>
              <button class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all hover:text-white active:scale-90" :disabled="!hasPrevious" aria-label="Previous track" @click="playPrevious">
                <i aria-hidden="true" class="pi pi-step-backward text-lg" />
              </button>
              <button class="flex h-12 w-12 items-center justify-center rounded-full shadow-2xl transition-all active:scale-90" :style="{ background: progressColor }" :disabled="!currentTrack || isLoadingTrack" :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'" @click="togglePlay">
                <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-lg text-white" />
                <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-lg ml-0.5 text-white" />
              </button>
              <button class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all hover:text-white active:scale-90" :disabled="!hasNext" aria-label="Next track" @click="playNext">
                <i aria-hidden="true" class="pi pi-step-forward text-lg" />
              </button>
              <button data-repeat-btn class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all hover:text-white active:scale-90" :class="repeatMode !== 'off' ? 'text-[#f472b6]' : ''" aria-label="Repeat" @click="openRepeatMenu($event)">
                <i aria-hidden="true" class="pi pi-refresh text-sm" />
              </button>
            </div>

            <!-- Volume row -->
            <div class="flex items-center gap-2 w-full max-w-[240px] mx-auto mt-3 mb-4">
              <button class="flex h-7 w-7 items-center justify-center text-white/40 hover:text-white shrink-0" :aria-label="muted ? 'Unmute' : 'Mute'" @click="toggleMute">
                <i aria-hidden="true" :class="muted ? 'pi pi-volume-off' : 'pi pi-volume-up'" class="text-[10px]" />
              </button>
              <div class="relative flex-1 flex items-center h-3">
                <div class="absolute inset-x-0 h-0.5 rounded-full bg-white/15" />
                <div class="absolute left-0 h-0.5 rounded-full" :style="{ width: `${muted ? 0 : Number(volume) * 100}%`, background: progressColor }" />
                <input type="range" min="0" max="1" step="0.01" class="absolute inset-0 w-full cursor-pointer opacity-0 z-10" :value="muted ? 0 : volume" @input="onVolume" />
              </div>
            </div>

          </div>
        </div>

        <!-- ── Shared shuffle / repeat popups (rendered here, positioned by JS) ── -->
        <div class="relative z-50">
          <Transition name="fade">
            <div v-if="showShuffleMenu" class="shuffle-menu fixed z-[9999] min-w-[150px] rounded-xl border border-white/10 bg-[#121212] p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl" :style="shuffleMenuPos" style="backdrop-filter: blur(24px);">
              <button v-for="mode in shuffleModes" :key="mode.value" type="button" class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/[0.08]" :class="shuffleMode === mode.value ? 'text-[#a855f7] bg-white/[0.06]' : 'text-slate-400 hover:text-white'" @click="setShuffleMode(mode.value)">
                <i aria-hidden="true" :class="mode.icon" class="text-sm" />
                <span class="flex-1 text-left">{{ mode.label }}</span>
                <span v-if="shuffleMode === mode.value" class="h-2 w-2 rounded-full bg-[#a855f7]" />
              </button>
            </div>
          </Transition>
          <Transition name="fade">
            <div v-if="showRepeatMenu" class="repeat-menu fixed z-[9999] min-w-[130px] rounded-xl border border-white/10 bg-[#121212] p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl" :style="repeatMenuPos" style="backdrop-filter: blur(24px);">
              <button v-for="mode in repeatModes" :key="mode.value" type="button" class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/[0.08]" :class="repeatMode === mode.value ? 'text-[#f472b6] bg-white/[0.06]' : 'text-slate-400 hover:text-white'" @click="setRepeatMode(mode.value)">
                <i aria-hidden="true" :class="mode.icon" class="text-sm" />
                <span class="flex-1 text-left">{{ mode.label }}</span>
                <span v-if="repeatMode === mode.value" class="h-2 w-2 rounded-full bg-[#f472b6]" />
              </button>
            </div>
          </Transition>
        </div>

      </div>
    </Transition>
  </Teleport>

  <AddToPlaylistDialog
    :visible="showAddToPlaylist"
    :track-id="currentTrack?.id || ''"
    :track-title="currentTrack?.title"
    @update:visible="showAddToPlaylist = $event"
  />
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, onUnmounted, ref, watch } from 'vue'
import { usePlayerControls, useTrackLike } from '@/composables/player'
import { usePlayerStore } from '@/stores/player'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { useLyricsApi } from '@/services/api/lyrics'
import { parseLRCLines, parsePlainLines } from '@/composables/lyrics'
import type { ParsedLine } from '@/composables/lyrics'
import AddToPlaylistDialog from './AddToPlaylistDialog.vue'
import SyncedLyrics from './SyncedLyrics.vue'

const props = withDefaults(defineProps<{
  visible: boolean
  initialTab?: string
}>(), {
  initialTab: 'now-playing',
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'toggle-queue': []
  'toggle-lyrics': []
}>()

const isOpen = ref(props.visible)
watch(() => props.visible, (v) => { isOpen.value = v })
watch(isOpen, (v) => { emit('update:visible', v) })

const pc = usePlayerControls()

const {
  currentTrack,
  isPlaying,
  isBuffering,
  isLoadingTrack,
  currentTime,
  duration,
  volume,
  muted,
  hasNext,
  hasPrevious,
  shuffleMode,
  repeatMode,
  togglePlayPause,
  toggleMute,
  setVolume,
  playNext,
  playPrevious,
  toggleShuffle,
  toggleRepeat,
} = pc

const showShuffleMenu = ref(false)
const showRepeatMenu = ref(false)
const shuffleMenuPos = ref({ top: '0px', left: '0px' })
const repeatMenuPos = ref({ top: '0px', left: '0px' })
const shuffleModes = [
  { value: 'off' as const, label: 'Off', icon: 'pi pi-ban' },
  { value: 'queue' as const, label: 'Shuffle Queue', icon: 'pi pi-sort-alt' },
  { value: 'catalog' as const, label: 'Random Catalog', icon: 'pi pi-globe' },
  { value: 'similar' as const, label: 'Similar Tracks', icon: 'pi pi-star' },
]
const repeatModes = [
  { value: 'off' as const, label: 'No Repeat', icon: 'pi pi-refresh' },
  { value: 'all' as const, label: 'Repeat All', icon: 'pi pi-sync' },
  { value: 'one' as const, label: 'Repeat One', icon: 'pi pi-undo' },
]
function setShuffleMode(mode: 'off' | 'queue' | 'catalog' | 'similar') {
  pc.setShuffleMode(mode)
  showShuffleMenu.value = false
}
function setRepeatMode(mode: 'off' | 'all' | 'one') {
  usePlayerStore().repeatMode = mode
  showRepeatMenu.value = false
}
function openShuffleMenu(e: MouseEvent) {
  const btn = e.currentTarget as HTMLElement
  const rect = btn.getBoundingClientRect()
  shuffleMenuPos.value = {
    top: `${rect.top - 8}px`,
    left: `${rect.left + rect.width / 2 - 75}px`,
  }
  showShuffleMenu.value = !showShuffleMenu.value
}
function openRepeatMenu(e: MouseEvent) {
  const btn = e.currentTarget as HTMLElement
  const rect = btn.getBoundingClientRect()
  repeatMenuPos.value = {
    top: `${rect.top - 8}px`,
    left: `${rect.left + rect.width / 2 - 65}px`,
  }
  showRepeatMenu.value = !showRepeatMenu.value
}
onMounted(() => {
  document.addEventListener('click', handleOutsideClick)
})
onUnmounted(() => {
  document.removeEventListener('click', handleOutsideClick)
})
function handleOutsideClick(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (showShuffleMenu.value) {
    const btn = document.querySelector('[data-shuffle-btn]')
    if (btn && !btn.contains(target) && !target.closest('.shuffle-menu')) {
      showShuffleMenu.value = false
    }
  }
  if (showRepeatMenu.value) {
    const btn = document.querySelector('[data-repeat-btn]')
    if (btn && !btn.contains(target) && !target.closest('.repeat-menu')) {
      showRepeatMenu.value = false
    }
  }
}

// ─── Seek hover preview (desktop) ──────────────────────────────
const seekHoverTime = ref<number | null>(null)
const seekHoverPos = ref<number | null>(null)

function onSeekHover(e: MouseEvent) {
  const el = progressRef.value
  if (!el || !duration.value) return
  const rect = el.getBoundingClientRect()
  const x = e.clientX - rect.left
  const ratio = Math.max(0, Math.min(1, x / rect.width))
  seekHoverTime.value = ratio * duration.value
  seekHoverPos.value = x
}

// ─── Seek hover preview (mobile) ───────────────────────────────
const mobileSeekHoverTime = ref<number | null>(null)
const mobileSeekHoverPos = ref<number | null>(null)

function onSeekHoverMobile(e: MouseEvent) {
  const el = progressRefMobile.value
  if (!el || !duration.value) return
  const rect = el.getBoundingClientRect()
  const x = e.clientX - rect.left
  const ratio = Math.max(0, Math.min(1, x / rect.width))
  mobileSeekHoverTime.value = ratio * duration.value
  mobileSeekHoverPos.value = x
}

const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const progressColor = computed(() => {
  return coverUrl.value && palette.value.vibrant ? palette.value.vibrant : '#1db954'
})

const progressRef = ref<HTMLElement>()
const progressRefMobile = ref<HTMLElement>()

const progressPercent = computed(() =>
  duration.value ? (currentTime.value / duration.value) * 100 : 0
)

const dynamicBg = computed(() => {
  const p = palette.value
  return {
    background: coverUrl.value
      ? `radial-gradient(ellipse 80% 60% at 50% 0%, ${p.vibrant}33 0%, transparent 70%), radial-gradient(ellipse 60% 40% at 100% 100%, ${p.muted}44 0%, transparent 60%), ${p.dark}`
      : '#08080A',
  }
})

function seekTo(seconds: number) {
  if (!duration.value) return
  const clamped = Math.max(0, Math.min(seconds, duration.value))
  pc.seek(clamped)
}

function getProgressEl() {
  return progressRef.value || progressRefMobile.value
}

function seek(e: MouseEvent | Touch | KeyboardEvent) {
  const el = getProgressEl()
  if (!el || !duration.value) return
  const rect = el.getBoundingClientRect()
  const clientX = 'clientX' in e ? e.clientX : 0
  const ratio = (clientX - rect.left) / rect.width
  seekTo(ratio * duration.value)
}

function startDrag() {
  const move = (ev: MouseEvent) => seek(ev)
  const up = () => {
    window.removeEventListener('mousemove', move)
    window.removeEventListener('mouseup', up)
  }
  window.addEventListener('mousemove', move)
  window.addEventListener('mouseup', up)
}

function startTouchDrag(e: TouchEvent) {
  const el = getProgressEl()
  const touch = e.touches[0]
  if (!touch || !el || !duration.value) return
  const rect = el.getBoundingClientRect()
  const ratio = (touch.clientX - rect.left) / rect.width
  seekTo(ratio * duration.value)

  const move = (ev: TouchEvent) => {
    const t = ev.touches[0]
    if (!t) return
    const r = (t.clientX - rect.left) / rect.width
    seekTo(r * duration.value)
  }
  const up = () => {
    window.removeEventListener('touchmove', move)
    window.removeEventListener('touchend', up)
  }
  window.addEventListener('touchmove', move)
  window.addEventListener('touchend', up)
}

function formatTime(s: number) {
  if (!s || !isFinite(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

function togglePlay() {
  togglePlayPause()
}

function close() {
  isOpen.value = false
}

function onVolume(e: Event) {
  setVolume(Number((e.target as HTMLInputElement).value))
}

const rootEl = ref<HTMLElement | null>(null)

onMounted(async () => {
  await nextTick()
  rootEl.value?.focus()
})

onBeforeUnmount(() => {})

// ─── Right panel tab: 'lyrics' | 'queue' ─────────────────────
const rightPanelTab = ref<'lyrics' | 'queue'>('lyrics')

// ─── Queue preview (Up Next) ──────────────────────────────────
const playerStore = usePlayerStore()

const queuePreview = computed(() => {
  const q = playerStore.queue
  if (!currentTrack.value || q.length === 0) return []
  const idx = q.findIndex((t) => t.id === currentTrack.value?.id)
  if (idx < 0) return q.slice(0, 5)
  return q.slice(idx + 1, idx + 6)
})

function playQueueItem(previewIndex: number) {
  const q = playerStore.queue
  if (!currentTrack.value) return
  const currentIdx = q.findIndex((t) => t.id === currentTrack.value?.id)
  if (currentIdx < 0) return
  const targetIdx = currentIdx + 1 + previewIndex
  const target = q[targetIdx]
  if (target) {
    playerStore.playTrack(target)
  }
}

// ─── Mobile: pull-to-close swipe gesture ──────────────────────
const swipeStartY = ref(0)
const swipeProgress = ref(0)
const isSwiping = ref(false)
const mobileContainerRef = ref<HTMLElement | null>(null)

function onSwipeStart(e: TouchEvent) {
  // Only start swipe if scrolled to top or nearly there
  const container = mobileContainerRef.value
  if (container && container.scrollTop > 10) return
  swipeStartY.value = e.touches[0]!.clientY
  isSwiping.value = true
}

function onSwipeMove(e: TouchEvent) {
  if (!isSwiping.value) return
  const delta = e.touches[0]!.clientY - swipeStartY.value
  swipeProgress.value = Math.max(0, Math.min(1, delta / 200))
}

function onSwipeEnd() {
  if (isSwiping.value && swipeProgress.value > 0.4) {
    close()
  }
  swipeProgress.value = 0
  isSwiping.value = false
}

const showAddToPlaylist = ref(false)

const trackId = computed(() => currentTrack.value?.id)
const { liked, toggleLike } = useTrackLike(trackId)

// ─── Active line for mobile single-line display ──────────────
const activeLineIdx = computed(() => {
  const t = currentTime.value
  const lines = parsedLines.value
  for (let i = lines.length - 1; i >= 0; i--) {
    if (t >= lines[i]!.timeSeconds) return i
  }
  if (lines.length > 0 && duration.value > 0 && lines[0]!.timeSeconds === 0) {
    const progress = t / duration.value
    return Math.min(Math.floor(progress * lines.length), lines.length - 1)
  }
  return -1
})

const activeLineText = computed(() => {
  const idx = activeLineIdx.value
  return idx >= 0 ? parsedLines.value[idx]?.text ?? null : null
})

// ─── Synced lyrics ─────────────────────────────────────────────
const lyricsApi = useLyricsApi()
const lyricsContent = ref<string | null>(null)
const lyricsType = ref<'lrc' | 'plain'>('plain')
const lyricsLanguage = ref<string>('en')
const lyricsLoading = ref(false)
const noLyrics = ref(false)

const parsedLines = ref<ParsedLine[]>([])

watch(() => currentTrack.value?.id, async (id) => {
  if (!id) {
    lyricsContent.value = null
    parsedLines.value = []
    noLyrics.value = false
    return
  }
  lyricsLoading.value = true
  noLyrics.value = false
  try {
    const data = await lyricsApi.getTrackLyrics(id, undefined, { silent: true })
    if (data?.content) {
      lyricsContent.value = data.content
      lyricsType.value = data.type === 'lrc' ? 'lrc' : 'plain'
      lyricsLanguage.value = data.language || 'en'
    } else {
      lyricsContent.value = null
      noLyrics.value = true
      lyricsLanguage.value = 'en'
    }
  } catch {
    noLyrics.value = true
    lyricsContent.value = null
    lyricsLanguage.value = 'en'
  } finally {
    lyricsLoading.value = false
  }
}, { immediate: false })

watch([() => lyricsContent.value, () => lyricsType.value], () => {
  if (!lyricsContent.value) {
    parsedLines.value = []
    return
  }
  parsedLines.value = lyricsType.value === 'lrc'
    ? parseLRCLines(lyricsContent.value)
    : parsePlainLines(lyricsContent.value)
}, { immediate: true })
</script>

<style scoped>
.fullscreen-enter-active,
.fullscreen-leave-active {
  transition: all 0.5s cubic-bezier(0.19, 1, 0.22, 1);
}
.fullscreen-enter-from,
.fullscreen-leave-to {
  opacity: 0;
  transform: translateY(100%);
}

@keyframes pulse {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 0.15; }
}
.animate-pulse {
  animation: pulse 2s ease-in-out infinite;
}

@keyframes cover-glow {
  0%, 100% { box-shadow: 0 0 60px rgba(168, 85, 247, 0.15), 0 0 120px rgba(29, 185, 84, 0.08); }
  50% { box-shadow: 0 0 80px rgba(168, 85, 247, 0.3), 0 0 160px rgba(29, 185, 84, 0.15); }
}
.cover-glow {
  animation: cover-glow 3s ease-in-out infinite;
}

/* Shimmer loading */
@keyframes shimmer-pulse {
  0% { opacity: 0.06; }
  50% { opacity: 0.15; }
  100% { opacity: 0.06; }
}
.shimmer {
  animation: shimmer-pulse 1.5s ease-in-out infinite;
}

/* Spring easing for thumb dot */
.ease-spring {
  transition-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1);
}

/* Wave visualizer bars */
@keyframes wave-bounce {
  0%, 100% { transform: scaleY(0.5); opacity: 0.3; }
  50% { transform: scaleY(1); opacity: 0.8; }
}
.wave-bar {
  animation: wave-bounce 0.8s ease-in-out infinite alternate;
  transform-origin: bottom;
}

@media (prefers-reduced-motion: reduce) {
  .wave-bar { animation: none; opacity: 0.4; }
}

/* Custom scrollbar for lyrics area */
.scroll-thin {
  scrollbar-width: thin;
  scrollbar-color: rgba(255,255,255,0.06) transparent;
}

@media (prefers-reduced-motion: reduce) {
  .fullscreen-enter-active,
  .fullscreen-leave-active { transition: none; }
  .fullscreen-enter-from,
  .fullscreen-leave-to { transform: none; opacity: 1; }
  .cover-glow { animation: none; }
  .shimmer { animation: none; }
}
</style>
