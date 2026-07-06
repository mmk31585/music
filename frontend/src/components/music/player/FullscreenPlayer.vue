<template>
  <Teleport to="body">
    <Transition name="fullscreen">
      <div
        v-if="isOpen"
        role="dialog"
        aria-modal="true"
        aria-label="Now Playing"
        class="fixed inset-0 z-9999 flex flex-col"
        :style="{ ...dynamicBg, transition: 'background 0.8s cubic-bezier(0.19, 1, 0.22, 1)' }"
      >

        <!-- Background blur layer with crossfade -->
        <div class="pointer-events-none absolute -inset-5 z-0 scale-110">
          <img
            v-if="currentTrack?.coverUrl"
            :src="currentTrack?.coverUrl"
            class="h-full w-full object-cover opacity-40 md:opacity-50 transition-all duration-1000"
            style="filter: blur(100px) saturate(2)"
          />
        </div>
        <div class="pointer-events-none absolute inset-0 z-1 bg-linear-to-b from-black/70 via-black/30 to-black/90" />

        <!-- ── Top bar ── -->
        <div class="relative z-10 flex shrink-0 items-center justify-between px-4 pt-3 md:px-8 md:pt-5" style="padding-top: max(0.75rem, env(safe-area-inset-top))">
          <button class="flex h-10 w-10 items-center justify-center rounded-full text-white/50 transition-all hover:text-white hover:bg-white/10" aria-label="Close" @click="close">
            <i aria-hidden="true" class="pi pi-chevron-down text-xl" />
          </button>
          <p class="text-[10px] font-semibold tracking-[0.2em] text-white/30 uppercase">Now Playing</p>
          <div class="flex items-center gap-1">
            <!-- Lyrics toggle -->
            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full transition-all"
              :class="activeTab === 'lyrics' ? 'bg-white/15 text-white' : 'text-white/50 hover:text-white hover:bg-white/10'"
              aria-label="Show lyrics"
              @click="activeTab = 'lyrics'"
            >
              <i aria-hidden="true" class="pi pi-align-left text-lg" />
            </button>
            <!-- Queue toggle -->
            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full transition-all"
              :class="activeTab === 'queue' ? 'bg-white/15 text-white' : 'text-white/50 hover:text-white hover:bg-white/10'"
              aria-label="Show queue"
              @click="activeTab = 'queue'"
            >
              <i aria-hidden="true" class="pi pi-list text-lg" />
            </button>
            <!-- Pop out to Picture-in-Picture (mini player) -->
            <button
              type="button"
              class="flex h-10 w-10 items-center justify-center rounded-full transition-all"
              :class="isPiPActive ? 'bg-white/15 text-white' : 'text-white/50 hover:text-white hover:bg-white/10'"
              aria-label="Pop out player"
              @click="togglePiP"
            >
              <i aria-hidden="true" class="pi pi-external-link text-sm" />
            </button>
          </div>
        </div>

        <!-- ── Desktop: side-by-side (md+) ── -->
        <div class="relative z-10 hidden md:flex flex-1 min-h-0 px-8 pb-6">
          <div
            class="flex flex-1 min-h-0 w-full mx-auto transition-all duration-300"
            :class="activeTab === 'queue' ? 'max-w-5xl justify-center gap-12' : 'max-w-7xl gap-8'"
          >

            <!-- LEFT: Cover + info + controls -->
            <div
              class="flex flex-col items-center gap-4 shrink-0 justify-center pb-12 transition-all duration-300"
              :class="activeTab === 'queue' ? 'w-95 lg:w-105' : 'w-95 lg:w-105'"
            >

              <!-- Album art -->
              <div class="relative">
                <div
                  class="w-60 lg:w-70 aspect-square overflow-hidden rounded-3xl shadow-2xl transition-all duration-700"
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
                    class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30"
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
                    <h2 class="text-xl lg:text-2xl font-bold truncate max-w-70" :style="{ color: palette.vibrant || '#fff' }">{{ currentTrack?.title }}</h2>
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
                  @keydown.arrow-left.prevent="seekTo(currentTime - 5)"
                  @keydown.arrow-right.prevent="seekTo(currentTime + 5)"
                  @keydown.home.prevent="seekTo(0)"
                  @keydown.end.prevent="seekTo(duration)"
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
                  <button class="flex h-9 w-9 items-center justify-center transition-all hover:scale-110" :class="liked ? 'text-aurora-pink' : 'text-white/40 hover:text-white'" :aria-label="liked ? 'Unlike' : 'Like'" @click="toggleLike">
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
                    :class="shuffleMode !== 'off' ? 'text-aurora-purple' : 'text-white/40 hover:text-white'"
                    :aria-label="shuffleMode === 'queue' ? 'Shuffle queue' : shuffleMode === 'catalog' ? 'Random catalog tracks' : shuffleMode === 'similar' ? 'Similar tracks' : 'Shuffle off'"
                    data-shuffle-btn
                    @click="openShuffleMenu($event)"
                  >
                    <i aria-hidden="true" class="pi pi-sort-alt text-sm" />
                    <span v-if="shuffleMode !== 'off'" class="absolute -top-0.5 -right-0.5 flex h-3 w-3 items-center justify-center rounded-full bg-aurora-purple text-[7px] font-bold text-white">{{ shuffleMode === 'queue' ? 'Q' : shuffleMode === 'catalog' ? 'R' : 'S' }}</span>
                  </button>
                </div>
                <button class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all duration-200 hover:text-white hover:scale-110 disabled:opacity-20" :disabled="!hasPrevious" aria-label="Previous track" @click="playPrevious">
                  <i aria-hidden="true" class="pi pi-step-backward text-lg" />
                </button>
                <button class="flex h-14 w-14 items-center justify-center rounded-full shadow-2xl transition-all duration-200 active:scale-95 hover:scale-105 disabled:opacity-40" :style="{ background: progressColor }" :disabled="!currentTrack || isLoadingTrack" :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'" @click="togglePlay">
                  <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-xl text-white" />
                  <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause' : 'pi pi-play'" class="text-xl text-white" />
                </button>
                <button class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all duration-200 hover:text-white hover:scale-110 disabled:opacity-20" :disabled="!hasNext" aria-label="Next track" @click="playNext">
                  <i aria-hidden="true" class="pi pi-step-forward text-lg" />
                </button>
                <div class="relative">
                  <button
                    type="button"
                    class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200 hover:scale-110"
                    :class="repeatMode !== 'off' ? 'text-aurora-pink' : 'text-white/40 hover:text-white'"
                    :aria-label="repeatMode === 'off' ? 'Repeat off' : repeatMode === 'all' ? 'Repeat all' : 'Repeat one'"
                    data-repeat-btn
                    @click="openRepeatMenu($event)"
                  >
                    <i aria-hidden="true" class="pi pi-refresh text-sm" />
                    <span v-if="repeatMode === 'one'" class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-aurora-pink text-[7px] font-bold text-black">1</span>
                  </button>
                </div>
              </div>

              <!-- Volume -->
              <div class="flex items-center gap-2 w-full max-w-65 mx-auto">
                <button class="flex h-8 w-8 items-center justify-center text-white/40 transition-colors hover:text-white shrink-0" :aria-label="muted ? 'Unmute' : 'Mute'" @click="toggleMute">
                  <i aria-hidden="true" :class="muted ? 'pi pi-volume-off' : 'pi pi-volume-up'" class="text-xs" />
                </button>
                <div class="relative flex-1 flex items-center group/vol h-4">
                  <div class="absolute inset-x-0 h-0.5 rounded-full bg-white/15" />
                  <div class="absolute left-0 h-0.5 rounded-full transition-all" :style="{ width: `${muted ? 0 : Number(volume) * 100}%`, background: progressColor }" />
                  <Slider :model-value="muted ? 0 : volume" @update:model-value="onVolume" :min="0" :max="1" :step="0.01" class="w-full z-10" aria-label="Volume" />
                </div>
                <button class="flex h-8 w-8 items-center justify-center text-white/40 transition-colors hover:text-white shrink-0" aria-label="Toggle lyrics" @click="activeTab = activeTab === 'lyrics' ? 'now-playing' : 'lyrics'">
                  <i aria-hidden="true" class="pi pi-align-left text-xs" />
                </button>
              </div>
            </div>

            <!-- RIGHT: Lyrics / Queue panel (desktop only) -->
            <div
              class="min-h-0 flex flex-col md:pl-6 transition-all duration-300"
              :class="activeTab === 'lyrics' ? 'flex-1' : 'w-85 shrink-0'"
            >

              <!-- ── LYRICS / WAVE ── -->
              <div v-if="activeTab === 'lyrics'" class="flex-1 min-h-0 flex flex-col">
                <div v-if="lyricsLoading" class="flex h-full flex-col items-center justify-center gap-3">
                  <div class="shimmer h-4 w-48 rounded bg-white/6" />
                  <div class="shimmer h-4 w-36 rounded bg-white/4" />
                  <div class="shimmer h-4 w-40 rounded bg-white/3" />
                </div>
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
                <div v-else class="flex h-full flex-col items-center justify-center gap-6">
                  <div class="wave-visualizer flex items-end gap-0.75 h-16" :style="{ '--wave-color': palette.vibrant }">
                    <span v-for="i in 32" :key="i" class="wave-bar w-0.75 rounded-full" :style="{ animationDelay: `${i * 0.08}s`, height: `${20 + Math.sin(i * 0.7) * 30 + 30}%`, background: palette.vibrant || '#1db954' }" />
                  </div>
                  <p class="text-xs text-white/25">No synced lyrics available</p>
                </div>
              </div>

              <!-- ── QUEUE ── -->
              <div v-else class="flex-1 min-h-0 flex flex-col">
                <div v-if="playerStore.queue.length > 0" class="flex-1 overflow-y-auto scroll-thin pr-1">
                  <!-- Up next badge -->
                  <div class="flex items-center gap-2 px-1 pb-3 text-[10px] font-semibold tracking-wider text-white/30 uppercase">
                    <i aria-hidden="true" class="pi pi-arrow-down text-[9px]" />
                    Queue
                    <span class="h-3.5 w-3.5 rounded-full bg-white/8 flex items-center justify-center text-[8px] font-bold text-white/40">{{ playerStore.queue.length }}</span>
                    <i v-if="shuffleMode !== 'off'" aria-hidden="true" class="pi pi-sort-alt text-[9px] text-aurora-purple ml-auto" title="Shuffle is on — actual next track is from shuffle order, not queue order" />
                  </div>
                  <draggable
                    :list="localQueue"
                    item-key="id"
                    handle=".drag-handle"
                    animation="200"
                    ghost-class="opacity-30"
                    class="space-y-1"
                    @end="onReorder"
                  >
                    <template #item="{ element: track, index }">
                      <div
                        class="group flex items-center gap-3 rounded-xl px-3 py-2.5 transition-all duration-200 cursor-pointer"
                        :class="index === 0 ? 'bg-white/6 ring-1 ring-white/8' : 'hover:bg-white/4'"
                        @click="playQueueItem(index)"
                      >
                        <span class="drag-handle flex w-5 items-center justify-center cursor-grab active:cursor-grabbing text-white/20 hover:text-white/60 transition-colors me-3">
                          <i aria-hidden="true" class="pi pi-bars text-xs" />
                        </span>
                        <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10 shadow-md ring-1 ring-white/6">
                          <img
                            v-if="track.coverUrl"
                            :src="track.coverUrl"
                            :alt="track.title"
                            class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-110"
                            loading="lazy"
                            @error="onImgError"
                          />
                          <div v-else class="flex h-full items-center justify-center">
                            <i aria-hidden="true" class="pi pi-music text-xs text-white/30" />
                          </div>
                          <div class="absolute inset-0 flex items-center justify-center bg-black/0 transition-all duration-200 group-hover:bg-black/30">
                            <i aria-hidden="true" class="pi pi-play-fill text-white text-sm opacity-0 transition-all duration-200 group-hover:opacity-100 drop-shadow-lg" />
                          </div>
                        </div>
                        <div class="min-w-0 flex-1">
                          <p class="truncate text-sm font-medium text-white/90 group-hover:text-white transition-colors">
                            {{ track.title }}
                          </p>
                          <p class="truncate text-xs text-white/40">{{ track.artistName }}</p>
                        </div>
                        <span class="text-[10px] font-mono tabular-nums shrink-0 text-white/25 group-hover:text-white/50 transition-colors">
                          {{ formatTime(track.durationSeconds ?? 0) }}
                        </span>
                        <button
                          type="button"
                          class="flex h-7 w-7 items-center justify-center rounded-full text-white/15 opacity-0 group-hover:opacity-100 transition-all hover:bg-white/10 hover:text-white/60 active:scale-90"
                          aria-label="Remove from queue"
                          @click.stop="removeFromQueue(index)"
                        >
                          <i aria-hidden="true" class="pi pi-times text-[10px]" />
                        </button>
                      </div>
                    </template>
                  </draggable>
                </div>
                <div v-else class="flex h-full items-center justify-center">
                  <div class="flex flex-col items-center gap-2">
                    <i aria-hidden="true" class="pi pi-list text-lg text-white/20" />
                    <p class="text-xs text-white/30">Queue is empty</p>
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
                  <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
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
                  <h2 class="text-base font-bold truncate max-w-55" :style="{ color: palette.vibrant || '#fff' }">{{ currentTrack?.title }}</h2>
                  <p class="text-sm text-white/50 mt-0.5 truncate max-w-55">{{ currentTrack?.artistName }}</p>
                </div>
                <div class="flex shrink-0 items-center gap-0.5">
                  <button class="flex h-8 w-8 items-center justify-center transition-all hover:scale-110" :class="liked ? 'text-aurora-pink' : 'text-white/40 hover:text-white'" :aria-label="liked ? 'Unlike' : 'Like'" @click="toggleLike">
                    <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-sm" />
                  </button>
                </div>
              </div>
            </div>

            <!-- ── Mobile lyrics (between cover and controls) ── -->
            <div class="w-full mt-4 px-2 min-h-15 flex items-center justify-center">
              <div v-if="lyricsLoading" class="shimmer h-5 w-56 rounded bg-white/6" />
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
                @keydown.arrow-left.prevent="seekTo(currentTime - 5)"
                @keydown.arrow-right.prevent="seekTo(currentTime + 5)"
                @keydown.home.prevent="seekTo(0)"
                @keydown.end.prevent="seekTo(duration)"
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
              <button data-shuffle-btn class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all hover:text-white active:scale-90" :class="shuffleMode !== 'off' ? 'text-aurora-purple' : ''" aria-label="Shuffle" @click="openShuffleMenu($event)">
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
              <button data-repeat-btn class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 transition-all hover:text-white active:scale-90" :class="repeatMode !== 'off' ? 'text-aurora-pink' : ''" aria-label="Repeat" @click="openRepeatMenu($event)">
                <i aria-hidden="true" class="pi pi-refresh text-sm" />
              </button>
            </div>

            <!-- Volume row -->
            <div class="flex items-center gap-2 w-full max-w-60 mx-auto mt-3 mb-4">
              <button class="flex h-7 w-7 items-center justify-center text-white/40 hover:text-white shrink-0" :aria-label="muted ? 'Unmute' : 'Mute'" @click="toggleMute">
                <i aria-hidden="true" :class="muted ? 'pi pi-volume-off' : 'pi pi-volume-up'" class="text-[10px]" />
              </button>
              <div class="relative flex-1 flex items-center h-3">
                <div class="absolute inset-x-0 h-0.5 rounded-full bg-white/15" />
                <div class="absolute left-0 h-0.5 rounded-full" :style="{ width: `${muted ? 0 : Number(volume) * 100}%`, background: progressColor }" />
                <Slider :model-value="muted ? 0 : volume" @update:model-value="onVolume" :min="0" :max="1" :step="0.01" class="absolute inset-0 w-full z-10" aria-label="Volume" />
              </div>
            </div>

          </div>
        </div>

        <!-- ── Shared shuffle / repeat popups (rendered here, positioned by JS) ── -->
        <div class="relative z-50">
          <Transition name="fade">
            <div
              v-if="showShuffleMenu"
              ref="shuffleMenuRef"
              role="menu"
              aria-label="Select shuffle mode"
              class="shuffle-menu fixed z-9999 min-w-37.5 rounded-xl border border-white/10 bg-surface-raised p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl"
              :style="shuffleMenuPos"
              style="backdrop-filter: blur(24px);"
              @keydown="onMenuKeydown($event, 'shuffle')"
            >
              <button v-for="(mode, idx) in shuffleModes" :key="mode.value" type="button" :ref="(el) => { if (el) shuffleItemRefs[idx] = el as HTMLElement }" role="menuitem" :tabindex="idx === 0 ? 0 : -1" class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/8" :class="shuffleMode === mode.value ? 'text-aurora-purple bg-white/6' : 'text-slate-400 hover:text-white'" @click="setShuffleMode(mode.value)">
                <i aria-hidden="true" :class="mode.icon" class="text-sm" />
                <span class="flex-1 text-left">{{ mode.label }}</span>
                <span v-if="shuffleMode === mode.value" class="h-2 w-2 rounded-full bg-aurora-purple" />
              </button>
            </div>
          </Transition>
          <Transition name="fade">
            <div
              v-if="showRepeatMenu"
              ref="repeatMenuRef"
              role="menu"
              aria-label="Select repeat mode"
              class="repeat-menu fixed z-9999 min-w-32.5 rounded-xl border border-white/10 bg-surface-raised p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl"
              :style="repeatMenuPos"
              style="backdrop-filter: blur(24px);"
              @keydown="onMenuKeydown($event, 'repeat')"
            >
              <button v-for="(mode, idx) in repeatModes" :key="mode.value" type="button" :ref="(el) => { if (el) repeatItemRefs[idx] = el as HTMLElement }" role="menuitem" :tabindex="idx === 0 ? 0 : -1" class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/8" :class="repeatMode === mode.value ? 'text-aurora-pink bg-white/6' : 'text-slate-400 hover:text-white'" @click="setRepeatMode(mode.value)">
                <i aria-hidden="true" :class="mode.icon" class="text-sm" />
                <span class="flex-1 text-left">{{ mode.label }}</span>
                <span v-if="repeatMode === mode.value" class="h-2 w-2 rounded-full bg-aurora-pink" />
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
import { usePlayerPiPController } from '@/composables/usePlayerPiPController'
import { useLyricsApi } from '@/services/api/lyrics'
import { parseLRCLines, parsePlainLines } from '@/composables/lyrics'
import type { ParsedLine } from '@/composables/lyrics'
import { onImgError } from '@/utils/helpers'
import { queueManager } from '@/services/player/queue-manager'
import draggable from 'vuedraggable'
import AddToPlaylistDialog from './AddToPlaylistDialog.vue'
import SyncedLyrics from './SyncedLyrics.vue'
import type { PlaybackTrack } from '@/services/api/player'

const props = withDefaults(defineProps<{
  visible: boolean
  initialTab?: string
}>(), {
  initialTab: 'now-playing',
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()

const isOpen = ref(props.visible)
watch(() => props.visible, (v) => { isOpen.value = v })
watch(isOpen, (v) => { emit('update:visible', v) })

// ─── Picture-in-Picture ──────────────────────────────────────
const pip = usePlayerPiPController()
const isPiPActive = pip.isOpen

async function togglePiP() {
  if (pip.isOpen.value) {
    pip.close()
  } else {
    try {
      await pip.open()
    } catch {
      // PiP not supported or user denied
    }
  }
}

// Auto-open PiP when leaving the page while fullscreen is open & playing
function onVisibilityChange() {
  if (document.visibilityState === 'hidden' && isOpen.value && isPlaying.value) {
    pip.open().catch(() => {})
  }
}

onMounted(() => {
  document.addEventListener('visibilitychange', onVisibilityChange)
})
onBeforeUnmount(() => {
  document.removeEventListener('visibilitychange', onVisibilityChange)
})

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
} = pc

const showShuffleMenu = ref(false)
const showRepeatMenu = ref(false)
const shuffleMenuRef = ref<HTMLElement | null>(null)
const shuffleItemRefs = ref<HTMLElement[]>([])
const repeatMenuRef = ref<HTMLElement | null>(null)
const repeatItemRefs = ref<HTMLElement[]>([])
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

// Auto-focus first menu item when menu opens
watch(showShuffleMenu, async (v) => {
  if (v) await nextTick(); shuffleItemRefs.value[0]?.focus()
})
watch(showRepeatMenu, async (v) => {
  if (v) await nextTick(); repeatItemRefs.value[0]?.focus()
})

function onMenuKeydown(e: KeyboardEvent, type: 'shuffle' | 'repeat') {
  const items = type === 'shuffle' ? shuffleItemRefs.value : repeatItemRefs.value
  const currentIdx = items.findIndex((el) => el === document.activeElement)
  switch (e.key) {
    case 'ArrowDown': { e.preventDefault(); const next = (currentIdx + 1) % items.length; items[next]?.focus(); break }
    case 'ArrowUp': { e.preventDefault(); const prev = (currentIdx - 1 + items.length) % items.length; items[prev]?.focus(); break }
    case 'Home': { e.preventDefault(); items[0]?.focus(); break }
    case 'End': { e.preventDefault(); items[items.length - 1]?.focus(); break }
    case 'Escape':
    case 'Tab': {
      if (type === 'shuffle') showShuffleMenu.value = false
      else showRepeatMenu.value = false
      break
    }
  }
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
    if (!btn || (!btn.contains(target) && !target.closest('.shuffle-menu'))) {
      showShuffleMenu.value = false
    }
  }
  if (showRepeatMenu.value) {
    const btn = document.querySelector('[data-repeat-btn]')
    if (!btn || (!btn.contains(target) && !target.closest('.repeat-menu'))) {
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
  if (!isFinite(s)) return '0:00'
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

function onVolume(val: number | number[]) {
  setVolume(typeof val === 'number' ? val : val[0] ?? 0)
}

const rootEl = ref<HTMLElement | null>(null)

onMounted(async () => {
  await nextTick()
  rootEl.value?.focus()
})

onBeforeUnmount(() => {})

// ─── Right panel tab: 'lyrics' | 'queue' ─────────────────────
const activeTab = ref('lyrics')

// Local queue copy for drag-and-drop reordering
const playerStore = usePlayerStore()
const localQueue = ref<PlaybackTrack[]>([])

// Keep localQueue in sync with the store queue.
// Uses deep watch on the full array ref to catch all mutations (push, splice, replace).
watch(() => [...playerStore.queue], (synced) => {
  localQueue.value = synced
}, { immediate: true })

function onReorder(event: { oldIndex: number; newIndex: number }) {
  queueManager.reorderQueue(event.oldIndex, event.newIndex)
  playerStore.queue = queueManager.all()
}

function playQueueItem(index: number) {
  const target = localQueue.value[index]
  if (target) {
    playerStore.playTrack(target)
  }
}

function removeFromQueue(index: number) {
  const q = [...playerStore.queue]
  q.splice(index, 1)
  playerStore.queue = q
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

/** Shared lyrics fetcher used by both track-change and open events */
async function fetchLyrics(id: string | number | undefined) {
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
}

watch(() => currentTrack.value?.id, async (id) => {
  await fetchLyrics(id)
}, { immediate: false })

// Fetch lyrics when the fullscreen player opens (even for the same track)
watch(isOpen, async (open) => {
  if (open && currentTrack.value?.id) {
    await fetchLyrics(currentTrack.value.id)
  }
})

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
/* ── PrimeVue Slider overrides for FullscreenPlayer ────────── */
:deep(.p-slider) {
  background: transparent !important;
}
:deep(.p-slider-range) {
  background: transparent !important;
}
/* Keep handle visible and styled for the dynamic bg.
   PrimeVue sets transform: translateX(-50%) by default to center
   the handle horizontally; we ADD vertical centering without
   removing horizontal centering. */
:deep(.p-slider-handle) {
  width: 14px !important;
  height: 14px !important;
  border-radius: 50% !important;
  background: #fff !important;
  box-shadow: 0 0 8px rgba(255, 255, 255, 0.3), 0 2px 8px rgba(0, 0, 0, 0.4) !important;
  border: none !important;
  margin-top: 0 !important;
  margin-left: 0 !important;
  top: 50% !important;
  transform: translate(-50%, -50%) !important;
  opacity: 1;
  transition: opacity 0.15s ease;
}
:deep(.p-slider:hover .p-slider-handle) {
  opacity: 1;
  transform: translate(-50%, -50%) scale(1.2);
}
:deep(.p-slider:active .p-slider-handle) {
  opacity: 1;
}
:deep(.p-slider-handle:focus-visible) {
  opacity: 1;
  box-shadow: 0 0 0 3px rgba(29, 185, 84, 0.4), 0 0 12px rgba(255, 255, 255, 0.3) !important;
}

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
