<template>
  <!-- DESKTOP BAR -->
  <Transition name="bar-slide">
    <div
      v-if="currentTrack"
      ref="barRef"
      role="contentinfo"
      aria-label="Music player"
      class="fixed inset-x-0 bottom-0 z-50 hidden md:block select-none"
    >
      <div
        class="relative flex flex-col overflow-visible transition-all duration-300 ease-out"
        :class="collapsed ? 'pb-0' : ''"
        style="background: #08080A; backdrop-filter: blur(32px); -webkit-backdrop-filter: blur(32px);"
      >
        <!-- Background blur gets its own overflow clip so it doesn't clip popup menus -->
        <div class="pointer-events-none absolute inset-0 overflow-hidden rounded-[inherit]">
          <div
            v-if="currentTrack.coverUrl"
            class="absolute inset-0 scale-110"
            aria-hidden="true"
          >
            <img
              :src="currentTrack.coverUrl"
              alt=""
              class="h-full w-full object-cover transition-opacity duration-300 ease-out"
              :class="collapsed ? 'opacity-[0.15]' : 'opacity-[0.35]'"
              style="filter: blur(60px) saturate(1.5)"
            />
          </div>
          <div class="pointer-events-none absolute inset-0 bg-linear-to-t from-black/80 via-black/50 to-transparent" />
        </div>

        <!-- ── MINI COLLAPSED BAR ── -->
        <!-- Thin progress line at top -->
        <div class="relative z-10 h-0.5 bg-white/5">
          <div class="h-full rounded-full transition-[width] duration-100" :style="progressStyle" />
        </div>
        <div v-if="collapsed" class="relative z-10 flex items-center gap-3 px-4 h-14">
          <div role="button" tabindex="0" aria-label="Open fullscreen player" class="relative shrink-0 cursor-pointer" @click="emit('toggle-fullscreen')" @keydown.enter="emit('toggle-fullscreen')" @keydown.space.prevent="emit('toggle-fullscreen')">
            <div class="h-10 w-10 overflow-hidden rounded-xl shadow-lg ring-1 ring-white/10">
              <img
                v-if="currentTrack.coverUrl"
                :src="currentTrack.coverUrl"
                :alt="currentTrack.title"
                class="h-full w-full object-cover"
                loading="lazy"
              />
              <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
                <i aria-hidden="true" class="pi pi-headphones text-sm text-slate-500" />
              </div>
            </div>
          </div>

          <div role="button" tabindex="0" aria-label="Open fullscreen player" class="min-w-0 flex-1 cursor-pointer" @click="emit('toggle-fullscreen')" @keydown.enter="emit('toggle-fullscreen')" @keydown.space.prevent="emit('toggle-fullscreen')">
            <p class="truncate text-sm font-bold text-white leading-tight">{{ currentTrack.title }}</p>
            <p class="truncate text-xs text-white/60 leading-tight">{{ currentTrack.artistName }}</p>
          </div>

          <div class="flex items-center gap-2">
            <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full shadow-lg disabled:opacity-40 hover:scale-110 active:scale-90 transition-all duration-200"
                :style="{ background: progressColor }"
                :disabled="!currentTrack || isLoadingTrack"
                :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'"
                @click="isPlaying ? togglePlayPause() : proceed()"
              >
                <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm text-white" />
                <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause' : 'pi pi-play'" class="ml-0.5 text-sm text-white" />
              </button>
            </GuestPlayGate>

            <button
              type="button"
              class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
              :aria-label="collapsed ? 'Expand player' : 'Collapse player'"
              @click="toggleCollapsed"
            >
              <i aria-hidden="true" :class="collapsed ? 'pi pi-chevron-up' : 'pi pi-chevron-down'" class="text-sm" />
            </button>
          </div>
        </div>

        <!-- ── FULL EXPANDED BAR ── -->
        <Transition name="expand">
          <div v-if="!collapsed" class="overflow-hidden">
          <div class="relative z-10 flex items-center gap-4 px-6 pt-3">
            <div class="flex min-w-0 w-[25%] items-center gap-3">
              <div role="button" tabindex="0" aria-label="Open fullscreen player" class="relative shrink-0 cursor-pointer" @click="emit('toggle-fullscreen')" @keydown.enter="emit('toggle-fullscreen')" @keydown.space.prevent="emit('toggle-fullscreen')">
                <div
                  class="h-14 w-14 overflow-hidden rounded-[18px] shadow-[0_16px_32px_rgba(0,0,0,0.5)] ring-1 ring-white/10 transition-all duration-700"
                  :class="isPlaying ? 'scale-100' : 'scale-95 opacity-80'"
                >
                  <img
                    v-if="currentTrack.coverUrl"
                    :src="currentTrack.coverUrl"
                    :alt="currentTrack.title"
                    class="h-full w-full object-cover"
                    loading="lazy"
                  />
                  <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
                    <i aria-hidden="true" class="pi pi-headphones text-lg text-slate-500" />
                  </div>
                </div>
              </div>

              <div class="min-w-0 flex-1">
                <div class="flex items-center gap-2">
                  <p class="truncate text-sm font-bold text-white leading-tight">{{ currentTrack.title }}</p>
                  <button
                    type="button"
                    class="shrink-0 flex items-center justify-center transition-all hover:scale-110 active:scale-90"
                    :class="liked ? 'text-aurora-pink' : 'text-white/30 hover:text-white/60'"
                    :aria-label="liked ? 'Unlike' : 'Like'"
                    @click.stop="toggleLike"
                  >
                    <i aria-hidden="true" :class="liked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-xs" />
                  </button>
                  <button
                    type="button"
                    class="shrink-0 flex items-center justify-center transition-all hover:scale-110 active:scale-90 text-white/30 hover:text-white/60"
                    aria-label="Add to playlist"
                    @click.stop="showAddToPlaylist = true"
                  >
                    <i aria-hidden="true" class="pi pi-plus text-xs" />
                  </button>
                </div>
                <p class="truncate text-xs text-white/50 mt-0.5 leading-tight">{{ currentTrack.artistName }}</p>
                <div v-if="isPlaying" class="flex items-center gap-1 mt-1">
                  <span class="size-1 rounded-full bg-white/60 animate-bounce" style="animation-delay: 0ms" />
                  <span class="size-1 rounded-full bg-white/60 animate-bounce" style="animation-delay: 150ms" />
                  <span class="size-1 rounded-full bg-white/60 animate-bounce" style="animation-delay: 300ms" />
                </div>
              </div>
            </div>

            <div class="flex flex-1 items-center justify-center gap-1.5">
              <div class="relative">
                <button
                  ref="shuffleBtnRef"
                  type="button"
                  class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200"
                  :class="shuffleMode !== 'off' ? 'text-aurora-purple' : 'text-white/40 hover:text-white hover:bg-white/10'"
                  :disabled="!currentTrack"
                  :aria-label="shuffleMode === 'queue' ? 'Shuffle queue' : shuffleMode === 'catalog' ? 'Random catalog tracks' : shuffleMode === 'similar' ? 'Similar tracks' : 'Shuffle off'"
                  @click="showShuffleMenu = !showShuffleMenu"
                >
                  <i aria-hidden="true" class="pi pi-sort-alt text-sm" />
                  <span
                    v-if="shuffleMode !== 'off'"
                    class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-aurora-purple text-[8px] font-bold text-white"
                  >{{ shuffleMode === 'queue' ? 'Q' : shuffleMode === 'catalog' ? 'R' : 'S' }}</span>
                </button>

                <!-- Shuffle mode selector popup -->
                <Transition name="fade">
                  <div
                    v-if="showShuffleMenu"
                    ref="shuffleMenuRef"
                    role="menu"
                    aria-label="Select shuffle mode"
                    class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 z-60 min-w-37.5 rounded-xl border border-white/10 bg-surface-raised p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl"
                    style="backdrop-filter: blur(24px);"
                    @keydown="onMenuKeydown($event, 'shuffle')"
                  >
                    <button
                      v-for="(mode, idx) in shuffleModes"
                      :key="mode.value"
                      :ref="(el) => { if (el) shuffleItemRefs[idx] = el as HTMLElement }"
                      type="button"
                      role="menuitem"
                      :tabindex="idx === 0 ? 0 : -1"
                      class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/8"
                      :class="shuffleMode === mode.value ? 'text-aurora-purple bg-white/6' : 'text-slate-400 hover:text-white'"
                      @click="setShuffleMode(mode.value)"
                    >
                      <i aria-hidden="true" :class="mode.icon" class="text-sm" />
                      <span class="flex-1 text-left">{{ mode.label }}</span>
                      <span v-if="shuffleMode === mode.value" class="h-2 w-2 rounded-full bg-aurora-purple" />
                    </button>
                  </div>
                </Transition>
              </div>

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white disabled:opacity-25 hover:bg-white/10 transition-all duration-200"
                :disabled="!hasPrevious"
                aria-label="Previous track"
                @click="playPrevious"
              >
                <i aria-hidden="true" class="pi pi-step-backward text-base" />
              </button>

              <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
                <button
                  type="button"
                  class="flex h-11 w-11 items-center justify-center rounded-full shadow-xl disabled:opacity-40 hover:scale-110 active:scale-95 transition-all duration-200"
                  :style="{ background: progressColor }"
                  :disabled="!currentTrack || isLoadingTrack"
                  :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'"
                @click="isPlaying ? togglePlayPause() : proceed()"
                >
                  <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-base text-white" />
                  <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause' : 'pi pi-play'" class="ml-0.5 text-base text-white" />
                </button>
              </GuestPlayGate>

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/50 hover:text-white disabled:opacity-25 hover:bg-white/10 transition-all duration-200"
                :disabled="!hasNext"
                aria-label="Next track"
                @click="playNext"
              >
                <i aria-hidden="true" class="pi pi-step-forward text-base" />
              </button>

              <div class="relative">
                <button
                  ref="repeatBtnRef"
                  type="button"
                  class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200"
                  :class="repeatMode !== 'off' ? 'text-aurora-pink' : 'text-white/40 hover:text-white hover:bg-white/10'"
                  :disabled="!currentTrack"
                  :aria-label="repeatMode === 'off' ? 'Repeat off' : repeatMode === 'all' ? 'Repeat all' : 'Repeat one'"
                  @click="showRepeatMenu = !showRepeatMenu"
                >
                  <i aria-hidden="true" class="pi pi-refresh text-sm" />
                  <span
                    v-if="repeatMode === 'one'"
                    class="absolute -top-0.5 -right-0.5 flex h-4 w-4 items-center justify-center rounded-full bg-aurora-pink text-[8px] font-bold text-black"
                  >1</span>
                </button>

                <!-- Repeat mode selector popup -->
                <Transition name="fade">
                  <div
                    v-if="showRepeatMenu"
                    ref="repeatMenuRef"
                    role="menu"
                    aria-label="Select repeat mode"
                    class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 z-[60] min-w-[130px] rounded-xl border border-white/10 bg-surface-raised p-1.5 shadow-[0_8px_32px_rgba(0,0,0,0.6)] backdrop-blur-2xl"
                    style="backdrop-filter: blur(24px);"
                    @keydown="onMenuKeydown($event, 'repeat')"
                  >
                    <button
                      v-for="(mode, idx) in repeatModes"
                      :key="mode.value"
                      :ref="(el) => { if (el) repeatItemRefs[idx] = el as HTMLElement }"
                      type="button"
                      role="menuitem"
                      :tabindex="idx === 0 ? 0 : -1"
                      class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-xs font-bold transition hover:bg-white/8"
                      :class="repeatMode === mode.value ? 'text-aurora-pink bg-white/6' : 'text-slate-400 hover:text-white'"
                      @click="setRepeatMode(mode.value)"
                    >
                      <i aria-hidden="true" :class="mode.icon" class="text-sm" />
                      <span class="flex-1 text-left">{{ mode.label }}</span>
                      <span v-if="repeatMode === mode.value" class="h-2 w-2 rounded-full bg-aurora-pink" />
                    </button>
                  </div>
                </Transition>
              </div>
            </div>

            <div class="flex w-[25%] items-center justify-end gap-1">
              <div class="relative">
                <button
                  ref="queueBtnRef"
                  type="button"
                  class="relative flex h-9 w-9 items-center justify-center rounded-full transition-all duration-200"
                  :class="upcomingCount > 0 ? 'text-white hover:bg-white/10' : 'text-white/40 hover:text-white hover:bg-white/10'"
                  :disabled="!currentTrack"
                  :aria-label="`Queue — ${upcomingCount} upcoming`"
                  @click="showQueuePreview = !showQueuePreview"
                >
                  <i aria-hidden="true" class="pi pi-list text-sm" />
                  <span
                    v-if="upcomingCount > 0"
                    class="absolute -top-0.5 -right-0.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-spotify text-[8px] font-bold text-black"
                  >{{ upcomingCount > 9 ? '9+' : upcomingCount }}</span>
                </button>

                <!-- Up Next mini-preview popup -->
                <Transition name="fade">
                  <div
                    v-if="showQueuePreview"
                    class="absolute bottom-full right-0 mb-2 z-60 w-72 origin-bottom-right rounded-2xl border border-white/8 p-2 shadow-[0_12px_48px_rgba(0,0,0,0.7)] backdrop-blur-2xl"
                    style="backdrop-filter: blur(32px); -webkit-backdrop-filter: blur(32px); background: rgba(10, 10, 12, 0.94);"
                  >
                    <!-- Now Playing -->
                    <div class="mb-2 px-2 pt-1">
                      <p class="text-[10px] font-semibold tracking-wider text-white/30 uppercase">Now Playing</p>
                      <div class="mt-1.5 flex items-center gap-2.5">
                        <div class="h-9 w-9 shrink-0 overflow-hidden rounded-lg bg-white/10 ring-1 ring-white/6">
                          <img
                            v-if="currentTrack?.coverUrl"
                            :src="currentTrack.coverUrl"
                            :alt="currentTrack.title"
                            class="h-full w-full object-cover"
                            loading="lazy"
                          />
                          <div v-else class="flex h-full items-center justify-center">
                            <i aria-hidden="true" class="pi pi-music text-xs text-white/30" />
                          </div>
                        </div>
                        <div class="min-w-0 flex-1">
                          <p class="truncate text-sm font-bold text-white/90">{{ currentTrack?.title }}</p>
                          <p class="truncate text-xs text-white/40">{{ currentTrack?.artistName }}</p>
                        </div>
                        <i aria-hidden="true" class="pi pi-waveform text-sm text-spotify" />
                      </div>
                    </div>

                    <div class="mx-2 my-1.5 border-t border-white/6" />

                    <!-- Next Up -->
                    <div v-if="nextTrack" class="px-2 pb-2">
                      <div class="flex items-center gap-2 text-[10px] font-semibold tracking-wider text-white/30 uppercase mb-2">
                        <i aria-hidden="true" class="pi pi-arrow-down text-[9px]" />
                        Up Next
                        <span v-if="upcomingCount > 1" class="h-3.5 w-3.5 rounded-full bg-white/8 flex items-center justify-center text-[8px] font-bold text-white/40">{{ upcomingCount }}</span>
                        <i v-if="shuffleMode !== 'off'" aria-hidden="true" class="pi pi-sort-alt text-[9px] text-aurora-purple ml-auto" title="Shuffle is on — next track from shuffle order" />
                      </div>
                      <div
                        class="group flex items-center gap-2.5 rounded-xl px-2.5 py-2 transition-all duration-200 cursor-pointer ring-1 ring-white/6 bg-white/6 hover:bg-white/10"
                        @click="playNextTrack"
                      >
                        <div class="relative h-10 w-10 shrink-0 overflow-hidden rounded-lg bg-white/10 shadow-sm ring-1 ring-white/6">
                          <img
                            v-if="nextTrack.coverUrl"
                            :src="nextTrack.coverUrl"
                            :alt="nextTrack.title"
                            class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-110"
                            loading="lazy"
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
                            {{ nextTrack.title }}
                          </p>
                          <p class="truncate text-xs text-white/40">{{ nextTrack.artistName }}</p>
                        </div>
                        <span class="text-[10px] font-mono tabular-nums shrink-0 text-white/25 group-hover:text-white/50 transition-colors">
                          {{ formatTime(nextTrack.durationSeconds ?? 0) }}
                        </span>
                      </div>
                    </div>
                    <!-- Shuffle catalog/similar: can't predict next track -->
                    <div v-else-if="shuffleMode === 'catalog' || shuffleMode === 'similar'" class="px-2 pb-2">
                      <div class="flex items-center gap-2 text-[10px] font-semibold tracking-wider text-white/30 uppercase mb-2">
                        <i aria-hidden="true" class="pi pi-sort-alt text-[9px] text-aurora-purple" />
                        Up Next
                        <span class="text-[8px] text-aurora-purple/60 font-normal">— random track</span>
                      </div>
                      <div
                        class="flex items-center gap-2.5 rounded-xl px-2.5 py-2 transition-all duration-200 cursor-pointer ring-1 ring-white/6 bg-white/6 hover:bg-white/10 group"
                        @click="playNextTrack"
                      >
                        <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-aurora-purple/10 ring-1 ring-aurora-purple/20">
                          <i aria-hidden="true" class="pi pi-shuffle text-sm text-aurora-purple" />
                        </div>
                        <div class="min-w-0 flex-1">
                          <p class="text-sm font-medium text-white/70 group-hover:text-white transition-colors">Skip to random track</p>
                          <p class="text-xs text-white/30">Next track selected from {{ shuffleMode === 'catalog' ? 'full catalog' : 'similar tracks' }}</p>
                        </div>
                      </div>
                    </div>
                    <div v-else class="flex flex-col items-center gap-1.5 px-2 pb-3 pt-2 text-center">
                      <i aria-hidden="true" class="pi pi-list text-lg text-white/20" />
                      <p class="text-xs text-white/30">No upcoming tracks in queue</p>
                    </div>

                    <!-- View full queue -->
                    <button
                      type="button"
                      class="flex w-full items-center justify-center gap-1.5 rounded-lg px-3 py-2 text-xs font-medium text-white/40 transition-all hover:bg-white/6 hover:text-white/70"
                      @click="openFullQueue"
                    >
                      View full queue
                      <i aria-hidden="true" class="pi pi-arrow-right text-[10px]" />
                    </button>
                  </div>
                </Transition>
              </div>

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :disabled="!currentTrack"
                aria-label="Toggle lyrics"
                @click="emit('toggle-lyrics')"
              >
                <i aria-hidden="true" class="pi pi-align-left text-sm" />
              </button>

              <div class="mx-1 h-6 w-px bg-white/10" />

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :aria-label="muted ? 'Unmute' : 'Mute'"
                @click="toggleMute"
              >
                <i aria-hidden="true" :class="volumeIcon" class="text-sm" />
              </button>
            <div class="w-20">
              <Slider
                :model-value="muted ? 0 : volume"
                @update:model-value="onVolume"
                :min="0"
                :max="1"
                :step="0.01"
                aria-label="Volume"
              />
            </div>

              <div class="mx-1 h-6 w-px bg-white/10" />

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :disabled="!currentTrack"
                aria-label="Fullscreen"
                @click="emit('toggle-fullscreen')"
              >
                <i aria-hidden="true" class="pi pi-arrow-up-right-and-arrow-down-left-from-center text-sm" />
              </button>

              <button
                type="button"
                class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                :aria-label="collapsed ? 'Expand player' : 'Collapse player'"
                @click="toggleCollapsed"
              >
                <i aria-hidden="true" :class="collapsed ? 'pi pi-chevron-up' : 'pi pi-chevron-down'" class="text-sm" />
              </button>

              <div class="relative overflow-menu-container">
                <button
                  type="button"
                  class="flex h-9 w-9 items-center justify-center rounded-full text-white/40 hover:text-white hover:bg-white/10 transition-all duration-200"
                  :disabled="!currentTrack"
                  aria-label="More options"
                  @click.stop="showOverflow = !showOverflow"
                >
                  <i aria-hidden="true" class="pi pi-ellipsis-h text-sm" />
                </button>
                <PlayerOverflowMenu v-if="showOverflow" @close="showOverflow = false" @add-to-playlist="showAddToPlaylist = true" @toggle-pip="handleTogglePiP" />
              </div>
            </div>
          </div>

          <div class="relative z-10 flex items-center gap-3 px-6 pb-3 pt-1">
            <span class="text-[11px] text-white/50 font-mono tabular-nums w-10 text-right">{{ formatTime(currentTime) }}</span>
            <div
              role="button"
              tabindex="0"
              class="group relative flex-1 h-1.5 cursor-pointer rounded-full transition-all duration-150"
              :class="isPlaying ? 'bg-white/15' : 'bg-white/10'"
              dir="ltr"
              @click="onSeekClick"
              @keydown.enter="onSeekClick"
              @keydown.space.prevent="onSeekClick"
              @keydown.arrow-left.prevent="seekRelative(-5)"
              @keydown.arrow-right.prevent="seekRelative(5)"
              @keydown.home.prevent="seekPercent(0)"
              @keydown.end.prevent="seekPercent(100)"
            >
              <!-- Base track glow -->
              <div
                v-if="isPlaying"
                class="absolute inset-0 rounded-full opacity-20 blur-xs transition-opacity duration-500"
                :style="{ background: progressColor }"
              />
              <!-- Fill -->
              <div
                class="relative h-full rounded-full"
                :class="isPlaying ? 'progress-bar-fill' : ''"
                :style="progressStyle"
              />
              <!-- Thumb dot -->
              <div
                class="absolute top-1/2 -translate-y-1/2 size-4 rounded-full shadow-xl opacity-0 group-hover:opacity-100 scale-0 group-hover:scale-100 transition-all duration-300 ease-spring"
                :style="{
                  left: `calc(${progressPercent}% - 8px)`,
                  background: progressColor,
                  boxShadow: `0 0 0 3px ${progressColor}22, 0 4px 12px rgba(0,0,0,0.5)`,
                }"
              />
            </div>
            <span class="text-[11px] text-white/50 font-mono tabular-nums w-10">{{ formatTime(duration) }}</span>
          </div>

          <div
            v-if="playbackError"
            class="relative z-10 flex items-center justify-center gap-2 bg-red-500/10 px-5 py-1.5 text-xs text-red-400"
          >
            <i aria-hidden="true" class="pi pi-exclamation-circle text-xs" />
            <span>{{ playbackError }}</span>
          </div>
        </div>
        </Transition>
      </div>
    </div>
  </Transition>

  <!-- MOBILE BAR - floats above bottom nav -->
  <Transition name="bar-slide">
    <div
      v-if="currentTrack"
      key="mobile-bar"
      role="contentinfo"
      aria-label="Music player"
      class="fixed inset-x-0 z-[45] md:hidden select-none"
      style="bottom: calc(4.5rem + env(safe-area-inset-bottom, 0px))"
    >
      <button
        type="button"
        class="relative mx-3 flex w-[calc(100%-1.5rem)] items-center gap-3 overflow-hidden rounded-2xl text-left shadow-2xl transition-all duration-500 active:scale-[0.98]"
        :style="{ background: coverUrl ? `${palette.dark}dd` : 'rgba(8, 8, 10, 0.92)', backdropFilter: 'blur(24px)', WebkitBackdropFilter: 'blur(24px)' }"
        @click="emit('toggle-fullscreen')"
      >
        <div
          v-if="currentTrack.coverUrl"
          class="absolute inset-0 scale-110 bg-cover bg-center blur-2xl opacity-[0.25] transition-all duration-700"
          :style="{ backgroundImage: `url(${currentTrack.coverUrl})` }"
          aria-hidden="true"
        />
        <div class="absolute inset-0 bg-linear-to-r from-black/80 via-black/50 to-black/80" />

        <div class="relative shrink-0">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-lg ring-1 ring-white/10">
            <img
              v-if="currentTrack.coverUrl"
              :src="currentTrack.coverUrl"
              :alt="currentTrack.title"
              class="h-full w-full object-cover"
              loading="lazy"
            />
            <div v-else class="flex h-full w-full items-center justify-center bg-linear-to-br from-spotify/30 to-aurora-purple/30">
              <i aria-hidden="true" class="pi pi-headphones text-sm text-slate-500" />
            </div>
          </div>
        </div>

        <div class="relative min-w-0 flex-1">
          <p class="truncate text-sm font-bold text-white">{{ currentTrack.title }}</p>
          <p class="truncate text-xs text-white/50">{{ currentTrack.artistName }}</p>
        </div>

        <div class="relative flex items-center gap-0.5 pr-1" @click.stop>
          <GuestPlayGate action="play" @proceed="togglePlayPause" #default="{ proceed }">
            <button
              type="button"
              class="flex h-9 w-9 items-center justify-center rounded-full bg-white/80 text-black shadow-lg disabled:opacity-40 active:scale-90 transition-transform"
              :disabled="!currentTrack || isLoadingTrack"
              :aria-label="isLoadingTrack || isBuffering ? 'Loading' : isPlaying ? 'Pause' : 'Play'"
              @click="isPlaying ? togglePlayPause() : proceed()"
            >
              <i aria-hidden="true" v-if="isLoadingTrack || isBuffering" class="pi pi-spin pi-spinner text-sm" />
              <i aria-hidden="true" v-else :class="isPlaying ? 'pi pi-pause' : 'pi pi-play'" class="ml-0.5 text-sm" />
            </button>
          </GuestPlayGate>
          <button
            type="button"
            class="flex h-9 w-9 items-center justify-center rounded-full text-white/60 hover:text-white active:scale-90 transition-transform"
            :disabled="!hasNext"
            aria-label="Next track"
            @click.stop="playNext"
          >
            <i aria-hidden="true" class="pi pi-step-forward text-sm" />
          </button>
        </div>
      </button>

      <div
        v-if="playbackError"
        class="mx-3 mt-1 flex items-center justify-center gap-2 rounded-xl bg-red-500/10 px-4 py-1.5 text-xs text-red-400"
      >
        <i aria-hidden="true" class="pi pi-exclamation-circle text-xs" />
        <span>{{ playbackError }}</span>
      </div>
    </div>
  </Transition>

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
import { usePlayerPiPController } from '@/composables/usePlayerPiPController'
import { useAlbumColors } from '@/composables/useAlbumColors'
import AddToPlaylistDialog from './AddToPlaylistDialog.vue'
import PlayerOverflowMenu from './PlayerOverflowMenu.vue'
import GuestPlayGate from '@/components/common/GuestPlayGate.vue'
import type { PlaybackTrack } from '@/services/api/player'

const emit = defineEmits<{
  'toggle-queue': []
  'toggle-fullscreen': []
  'toggle-lyrics': []
  'toggle-mobile-sheet': []
}>()

const showOverflow = ref(false)
const showAddToPlaylist = ref(false)
const collapsed = ref(localStorage.getItem('player-bar-collapsed') === 'true')
const showShuffleMenu = ref(false)
const shuffleBtnRef = ref<HTMLElement | null>(null)
const shuffleMenuRef = ref<HTMLElement | null>(null)
const shuffleItemRefs = ref<HTMLElement[]>([])
const showRepeatMenu = ref(false)
const repeatBtnRef = ref<HTMLElement | null>(null)
const repeatMenuRef = ref<HTMLElement | null>(null)
const repeatItemRefs = ref<HTMLElement[]>([])
const showQueuePreview = ref(false)
const queueBtnRef = ref<HTMLElement | null>(null)
const showFullQueue = ref(false)

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
    case 'ArrowDown': {
      e.preventDefault()
      const next = (currentIdx + 1) % items.length
      items[next]?.focus()
      break
    }
    case 'ArrowUp': {
      e.preventDefault()
      const prev = (currentIdx - 1 + items.length) % items.length
      items[prev]?.focus()
      break
    }
    case 'Home': {
      e.preventDefault()
      items[0]?.focus()
      break
    }
    case 'End': {
      e.preventDefault()
      items[items.length - 1]?.focus()
      break
    }
    case 'Escape':
    case 'Tab': {
      if (type === 'shuffle') showShuffleMenu.value = false
      else showRepeatMenu.value = false
      // Restore focus to trigger button
      const trigger = type === 'shuffle' ? shuffleBtnRef.value : repeatBtnRef.value
      trigger?.focus()
      break
    }
  }
}

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

// Close popup menus on outside click
onMounted(() => {
  document.addEventListener('click', handleOutsideClick)
})
onUnmounted(() => {
  document.removeEventListener('click', handleOutsideClick)
})
function handleOutsideClick(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (showShuffleMenu.value && shuffleBtnRef.value && !shuffleBtnRef.value.contains(target)) {
    showShuffleMenu.value = false
  }
  if (showRepeatMenu.value && repeatBtnRef.value && !repeatBtnRef.value.contains(target)) {
    showRepeatMenu.value = false
  }
  if (showQueuePreview.value && queueBtnRef.value && !queueBtnRef.value.contains(target)) {
    showQueuePreview.value = false
  }
}

function toggleCollapsed() {
  collapsed.value = !collapsed.value
  localStorage.setItem('player-bar-collapsed', String(collapsed.value))
  window.dispatchEvent(new CustomEvent('playerbar-collapse', { detail: collapsed.value }))
}

const pip = usePlayerPiPController()

const pc = usePlayerControls()

const {
  currentTrack,
  isPlaying,
  isBuffering,
  isLoadingTrack,
  volume,
  muted,
  currentTime,
  duration,
  error: playbackError,
  hasNext,
  hasPrevious,
  shuffleMode,
  repeatMode,
  volumeIcon,
  togglePlayPause,
  toggleMute,
  setVolume,
  progressPercent,
  seekPercent,
  playNext,
  playPrevious,
} = pc

const trackId = computed(() => currentTrack.value?.id)
const { liked, toggleLike } = useTrackLike(trackId)

// Queue next-up preview — uses store directly to ensure reactivity on queue changes
const playerStore = usePlayerStore()
const queueTracks = computed(() => playerStore.queue as PlaybackTrack[])
const currentTrackIndex = computed(() => {
  if (!currentTrack.value) return -1
  return queueTracks.value.findIndex((t) => t.id === currentTrack.value?.id)
})
/** The actual next track that will play (accounting for shuffle order). */
const nextTrack = computed(() => {
  // Use the engine's shuffle-aware peek when available
  const engineNext = playerStore.nextUpTrack
  if (engineNext) return engineNext

  // Fallback: sequential next from queue (for catalog/similar shuffle)
  const idx = currentTrackIndex.value
  if (idx < 0) return queueTracks.value[0] ?? null
  return queueTracks.value[idx + 1] ?? null
})
const upcomingCount = computed(() => {
  const idx = currentTrackIndex.value
  if (idx < 0) return queueTracks.value.length

  if (shuffleMode.value === 'catalog' || shuffleMode.value === 'similar') {
    // Can't predict count for random-api-based shuffle modes
    return queueTracks.value.length - idx - 1
  }

  return Math.max(0, queueTracks.value.length - idx - 1)
})

function playNextTrack() {
  playNext()
  showQueuePreview.value = false
}

function openFullQueue() {
  showQueuePreview.value = false
  emit('toggle-queue')
}

const coverUrl = computed(() => currentTrack.value?.coverUrl || null)
const { palette } = useAlbumColors(coverUrl)

const progressColor = computed(() => {
  return coverUrl.value && palette.value.vibrant ? palette.value.vibrant : '#1db954'
})

const progressGlow = computed(() => {
  const c = progressColor.value
  return `0 0 8px ${c}66, 0 0 20px ${c}33`
})

const progressStyle = computed(() => ({
  width: `${progressPercent.value}%`,
  background: progressColor.value,
  boxShadow: progressGlow.value,
  transition: 'width 100ms linear, background 0.5s ease, box-shadow 0.3s ease',
}))

function formatTime(s: number) {
  if (!isFinite(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

const barRef = ref<HTMLElement | null>(null)

function onVolume(val: number | number[]) {
  pc.setVolume(typeof val === 'number' ? val : val[0] ?? 0)
}

function onSeekClick(e: MouseEvent | KeyboardEvent) {
  const bar = e.currentTarget as HTMLElement
  const rect = bar.getBoundingClientRect()
  const pct = (((e as MouseEvent).clientX - rect.left) / rect.width) * 100
  seekPercent(pct)
}

function seekRelative(seconds: number) {
  const newTime = Math.max(0, Math.min(duration.value, currentTime.value + seconds))
  const pct = duration.value > 0 ? (newTime / duration.value) * 100 : 0
  seekPercent(pct)
}

function handleTogglePiP() {
  pip.toggle().catch(() => {
    // PiP not supported or user denied — silently fail
  })
}

function onOverflowClickOutside(e: MouseEvent) {
  if (!(e.target as HTMLElement).closest('.overflow-menu-container')) {
    showOverflow.value = false
  }
}

onMounted(() => document.addEventListener('click', onOverflowClickOutside))
onBeforeUnmount(() => document.removeEventListener('click', onOverflowClickOutside))
</script>

<style scoped>
.bar-slide-enter-active {
  transition: transform 400ms cubic-bezier(0.19, 1, 0.22, 1);
}
.bar-slide-leave-active {
  transition: transform 300ms ease-in;
}
.bar-slide-enter-from,
.bar-slide-leave-to {
  transform: translateY(100%);
}

/* ── Expand/Collapse Transition ────────────────── */
.expand-enter-active {
  transition: opacity 300ms ease, transform 300ms ease;
}
.expand-leave-active {
  transition: opacity 200ms ease, transform 200ms ease;
}
.expand-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}
.expand-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}
.animate-bounce {
  animation: bounce 0.6s ease-in-out infinite;
}

/* ── Progress Bar ─────────────────────────────── */

/* Subtle pulse glow on the progress fill when playing */
@keyframes progress-glow {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}
.progress-bar-fill::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  background: inherit;
  filter: blur(6px);
  opacity: 0.35;
  animation: progress-glow 2s ease-in-out infinite;
  z-index: -1;
}

/* Ease-spring utility for the thumb dot */
.ease-spring {
  transition-timing-function: cubic-bezier(0.34, 1.56, 0.64, 1);
}

/* Time labels hover highlight */
.group:hover .time-highlight {
  color: rgba(255, 255, 255, 0.7);
}

@media (prefers-reduced-motion: reduce) {
  .bar-slide-enter-active,
  .bar-slide-leave-active,
  .expand-enter-active,
  .expand-leave-active { transition: none; }
  .bar-slide-enter-from,
  .bar-slide-leave-to,
  .expand-enter-from,
  .expand-leave-to,
  .expand-enter-to,
  .expand-leave-from { opacity: 1; transform: none; }
  .animate-bounce { animation: none; }
  .progress-bar-fill::after { animation: none; opacity: 0; }
  .transition-opacity { transition: none; }
}
</style>
