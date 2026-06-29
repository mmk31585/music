<template>
  <Teleport to="body">
    <div
      ref="containerRef"
      class="fixed inset-0 z-[9999] flex flex-col bg-black"
      dir="ltr"
      role="dialog"
      aria-label="Video player"
      @keydown="onKeyDown"
      tabindex="0"
    >
      <!-- Top bar (fades out after 2s idle) -->
      <Transition name="fade">
        <div
          v-if="showTopBar"
          class="pointer-events-none absolute top-0 right-0 left-0 z-10 flex items-center justify-between bg-linear-to-b from-black/60 to-transparent px-4 pb-6 pt-4"
        >
          <button
            type="button"
            class="pointer-events-auto flex h-9 w-9 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-xs transition hover:bg-white/20"
            aria-label="Close"
            @click="handleClose"
          >
            <i aria-hidden="true" class="pi pi-chevron-right text-lg" />
          </button>

          <button
            type="button"
            class="pointer-events-auto flex h-9 w-9 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-xs transition hover:bg-white/20"
            aria-label="More options"
            @click="showOverflow = !showOverflow"
          >
            <i aria-hidden="true" class="pi pi-ellipsis-v text-sm" />
          </button>
        </div>
      </Transition>

      <!-- Overflow menu -->
      <Transition name="fade">
        <div
          v-if="showOverflow"
          class="absolute top-16 left-4 z-20 min-w-44 overflow-hidden rounded-2xl border border-white/10 bg-surface-raised shadow-2xl shadow-black/40 backdrop-blur-2xl"
        >
          <button
            type="button"
            class="flex w-full items-center gap-3 px-4 py-3 text-sm font-medium text-white/80 transition hover:bg-white/10"
            @click="shareOnTelegram"
          >
            <i aria-hidden="true" class="pi pi-telegram text-base" />
            Share on Telegram
          </button>
          <button
            type="button"
            class="flex w-full items-center gap-3 px-4 py-3 text-sm font-medium text-white/80 transition hover:bg-white/10"
            @click="goToTrack"
          >
            <i aria-hidden="true" class="pi pi-music text-base" />
            Go to track page
          </button>
          <button
            type="button"
            class="flex w-full items-center gap-3 px-4 py-3 text-sm font-medium text-red-400/80 transition hover:bg-white/10"
            @click="showOverflow = false"
          >
            <i aria-hidden="true" class="pi pi-flag text-base" />
            Report
          </button>
        </div>
      </Transition>

      <!-- Desktop side navigation arrows -->
      <button
        v-if="hasPrevious && !isTransitioning"
        type="button"
        class="absolute top-1/2 left-4 z-20 hidden h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full bg-black/50 text-white backdrop-blur-xs transition hover:bg-white/20 md:flex"
        aria-label="Previous video"
        @click="goToPrevious"
      >
        <i aria-hidden="true" class="pi pi-chevron-up text-xl" />
      </button>
      <button
        v-if="hasNext && !isTransitioning"
        type="button"
        class="absolute top-1/2 right-4 z-20 hidden h-12 w-12 -translate-y-1/2 items-center justify-center rounded-full bg-black/50 text-white backdrop-blur-xs transition hover:bg-white/20 md:flex"
        aria-label="Next video"
        @click="goToNext"
      >
        <i aria-hidden="true" class="pi pi-chevron-down text-xl" />
      </button>

      <!-- Swipeable video area -->
      <div
        ref="swipeRef"
        class="relative flex flex-1 touch-none overflow-hidden outline-hidden"
      >
        <!-- Current video -->
        <div
          class="absolute inset-0 flex items-center justify-center"
          :class="{ 'transition-transform': reducedMotion! }"
          :style="{ transform: `translateY(${currentTranslateY}%)`, transitionDuration: transitionDuration }"
        >
          <video
            ref="currentVideoRef"
            :src="currentSrc"
            :poster="currentPoster"
            class="h-full w-full object-contain"
            :muted="isMuted"
            playsinline
            loop="false"
            preload="auto"
            @loadedmetadata="onVideoLoaded"
            @ended="onEnded"
            @error="onVideoError"
          />

          <!-- Tap zone: left tap = prev, right tap = next, center = toggle mute -->
          <div class="absolute inset-0 hidden md:flex">
            <div class="flex-1 cursor-w-resize" @click="goToPrevious" />
            <div class="flex cursor-pointer items-center justify-center" @click="toggleMute">
              <div class="h-16 w-16" />
            </div>
            <div class="flex-1 cursor-e-resize" @click="goToNext" />
          </div>
        </div>

        <!-- Next video (prefetched) -->
        <div
          v-if="nextVideo"
          class="absolute inset-0 flex items-center justify-center"
          :class="{ 'transition-transform': reducedMotion! }"
          :style="{ transform: `translateY(${nextTranslateY}%)`, transitionDuration: transitionDuration }"
        >
          <video
            ref="nextVideoRef"
            :src="nextSrc"
            class="h-full w-full object-contain"
            muted
            playsinline
            preload="auto"
          />
        </div>

        <!-- Swipe hint -->
        <Transition name="fade">
          <div
            v-if="showSwipeHint && hasNext"
            class="absolute bottom-20 left-1/2 z-20 -translate-x-1/2 animate-bounce"
          >
            <div class="flex flex-col items-center gap-1 rounded-full bg-black/40 px-4 py-2 backdrop-blur-xs">
              <i aria-hidden="true" class="pi pi-chevron-down text-sm text-white/60" />
              <span class="text-[10px] font-medium text-white/40">Swipe for more</span>
            </div>
          </div>
        </Transition>
      </div>

      <!-- Right side actions -->
      <div class="absolute bottom-28 left-0 z-10 flex flex-col items-center gap-5 px-3">
        <!-- Like button -->
        <div class="flex flex-col items-center gap-1">
          <button
            type="button"
            class="like-btn flex h-11 w-11 items-center justify-center rounded-full text-lg backdrop-blur-xs transition-all duration-150"
            :class="isLiked ? 'bg-spotify/20 text-spotify' : 'bg-white/10 text-white hover:bg-white/20'"
            :style="{ transform: `scale(${likeScale})` }"
            @click.stop="toggleLike"
            :aria-label="isLiked ? 'Unlike' : 'Like'"
          >
            <i aria-hidden="true" :class="isLiked ? 'pi pi-heart-fill' : 'pi pi-heart'" />
          </button>
          <span class="text-[10px] font-bold text-white/70">{{ formatCount(currentVideo?.like_count || 0) }}</span>
        </div>

        <!-- Comment button (placeholder) -->
        <div class="flex flex-col items-center gap-1">
          <button
            type="button"
            class="flex h-11 w-11 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-xs transition hover:bg-white/20"
            aria-label="Comments"
            @click.stop="showComments = !showComments"
          >
            <i aria-hidden="true" class="pi pi-comment text-lg" />
          </button>
          <span class="text-[10px] font-bold text-white/70">{{ commentCount }}</span>
        </div>

        <!-- Track indicator -->
        <div
          v-if="currentVideo?.track"
          class="flex flex-col items-center gap-1"
        >
          <div
            class="flex h-11 w-11 items-center justify-center overflow-hidden rounded-full bg-white/10"
          >
            <AppImage
              :src="currentVideo.track.cover_url"
              :alt="currentVideo.track.title"
              class="h-full w-full object-cover"
              fallback-icon="pi pi-music"
              icon-size="1rem"
            />
          </div>
        </div>
      </div>

      <!-- Mute/unmute indicator -->
      <Transition name="fade">
        <div
          v-if="showMuteIndicator"
          class="pointer-events-none absolute top-1/2 left-1/2 z-20 -translate-x-1/2 -translate-y-1/2"
        >
          <div class="flex h-16 w-16 items-center justify-center rounded-full bg-black/50 backdrop-blur-xs">
            <i aria-hidden="true" :class="isMuted ? 'pi pi-volume-off' : 'pi pi-volume-up'" class="text-2xl text-white" />
          </div>
        </div>
      </Transition>

      <!-- Bottom strip -->
      <div class="absolute right-0 bottom-0 left-0 z-10 space-y-2 bg-linear-to-t from-black/80 via-black/40 to-transparent px-4 pb-4 pt-12">
        <!-- Uploader info -->
        <div v-if="currentVideo?.uploader" class="flex items-center gap-2">
          <div class="flex h-7 w-7 shrink-0 items-center justify-center overflow-hidden rounded-full bg-white/10">
            <AppImage
              :src="currentVideo.uploader.avatar_url"
              :alt="currentVideo.uploader.username"
              class="h-full w-full object-cover"
              fallback-icon="pi pi-user"
              icon-size="0.625rem"
            />
          </div>
          <span class="text-xs font-semibold text-white">{{ currentVideo.uploader.username }}</span>
          <span v-if="currentVideo.type === 'official_mv'" class="rounded-sm bg-spotify/20 px-1.5 py-0.5 text-[9px] font-bold text-spotify">MV</span>
        </div>

        <!-- Title -->
        <p class="line-clamp-2 text-sm font-medium text-white/90">
          {{ currentVideo?.title || '' }}
        </p>

        <!-- NowPlayingStrip -->
        <NowPlayingStrip
          v-if="currentVideo?.track"
          :track="currentVideo.track"
          @play="handlePlayTrack"
        />
      </div>

      <!-- Comments panel (slide-in) -->
      <Transition name="slide-up">
        <div
          v-if="showComments"
          class="absolute inset-0 z-30 flex flex-col bg-surface-base/95 backdrop-blur-xl"
        >
          <div class="flex items-center justify-between border-b border-white/6 px-4 py-3">
            <h3 class="text-sm font-bold text-white">Comments</h3>
            <button
              type="button"
              class="flex h-8 w-8 items-center justify-center rounded-full bg-white/10 text-white/60 transition hover:bg-white/20 hover:text-white"
              aria-label="Close comments"
              @click="showComments = false"
            >
              <i aria-hidden="true" class="pi pi-times text-sm" />
            </button>
          </div>

          <!-- Comments list -->
          <div class="flex-1 overflow-y-auto px-4 py-3">
            <div v-if="comments.length === 0" class="flex flex-col items-center justify-center py-16 text-center">
              <div class="flex h-12 w-12 items-center justify-center rounded-full bg-white/5">
                <i aria-hidden="true" class="pi pi-comment text-slate-500" />
              </div>
              <p class="mt-3 text-sm font-medium text-white/50">No comments yet</p>
              <p class="mt-1 text-xs text-white/30">Be the first to share your thoughts!</p>
            </div>
            <div v-else class="space-y-4">
              <div
                v-for="comment in comments"
                :key="String(comment.id)"
                class="flex gap-3"
              >
                <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-white/10 text-xs font-bold text-white">
                  {{ comment.author?.charAt(0).toUpperCase() || '?' }}
                </div>
                <div class="min-w-0 flex-1">
                  <div class="flex items-baseline gap-2">
                    <span class="text-xs font-semibold text-white/90">{{ comment.author || 'Anonymous' }}</span>
                    <span class="text-[10px] text-white/30">{{ timeAgo(comment.created_at) }}</span>
                  </div>
                  <p class="mt-0.5 text-sm leading-relaxed text-white/80">{{ comment.content }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Comment input -->
          <div class="border-t border-white/6 px-4 py-3">
            <div class="flex items-center gap-2">
              <input
                v-model="newComment"
                type="text"
                placeholder="Write a comment..."
                class="flex-1 rounded-xl border border-white/8 bg-white/3 px-4 py-2.5 text-sm text-white placeholder:text-slate-600 outline-hidden transition focus:border-spotify/40 focus:bg-white/6"
                @keydown.enter="submitComment"
              />
              <button
                type="button"
                class="flex h-10 w-10 items-center justify-center rounded-xl bg-spotify text-black transition hover:bg-spotify/90 disabled:opacity-30"
                :disabled="!newComment.trim() || submittingComment"
                @click="submitComment"
              >
                <i v-if="!submittingComment" class="pi pi-send text-sm" />
                <i v-else class="pi pi-spin pi-spinner text-sm" />
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import type { VideoItem, TrackSummary } from '@/services/api/video/types'
import AppImage from '@/components/common/AppImage.vue'
import NowPlayingStrip from '@/components/video/NowPlayingStrip.vue'
import { useVerticalSwipe } from '@/composables/useVerticalSwipe'
import { usePlayerStore } from '@/stores/player'
import { useVideoApi } from '@/services/api/video'
import { useRequest } from '@/composables/useRequest'
import { formatCount } from '@/utils/number'

interface VideoComment {
  id: string | number
  author?: string
  content: string
  created_at: string
}

const props = defineProps<{
  videos: VideoItem[]
  startIndex: number
}>()

const emit = defineEmits<{
  close: []
}>()

const router = useRouter()
const playerStore = usePlayerStore()
const videoApi = useVideoApi()

// ── State ─────────────────────────────────────────────────────────────

const containerRef = ref<HTMLElement | null>(null)
const swipeRef = ref<HTMLElement | null>(null)
const currentVideoRef = ref<HTMLVideoElement | null>(null)
const nextVideoRef = ref<HTMLVideoElement | null>(null)

const currentIndex = ref(props.startIndex)
const translateY = ref(0)
const showTopBar = ref(true)
const showOverflow = ref(false)
const showSwipeHint = ref(true)
const isMuted = ref(true)
const showMuteIndicator = ref(false)
const isTransitioning = ref(false)

const isLiked = ref(false)
const likeScale = ref(1)
const viewed = ref<Set<string>>(new Set())

// Comments
const showComments = ref(false)
const comments = ref<VideoComment[]>([])
const newComment = ref('')
const submittingComment = ref(false)
const commentCount = ref(0)
let commentsAbort: AbortController | null = null

// Animation
const reducedMotion = ref(false)
const transitionDuration = '250ms'

// ── Computed ──────────────────────────────────────────────────────────

const currentVideo = computed(() => props.videos[currentIndex.value] || null)
const nextVideo = computed(() => {
  const idx = currentIndex.value + 1
  return idx < props.videos.length ? props.videos[idx] : null
})
const hasNext = computed(() => currentIndex.value < props.videos.length - 1)
const hasPrevious = computed(() => currentIndex.value > 0)

const currentTranslateY = computed(() => -translateY.value)
const nextTranslateY = computed(() => 100 - translateY.value)

const currentSrc = computed(() => {
  const v = currentVideo.value
  if (!v) return ''
  return v.final_video_url || `/api/v1/videos/${v.id}/stream`
})
const currentPoster = computed(() => {
  const v = currentVideo.value
  if (!v) return undefined
  return v.thumbnail_url || v.thumbnail_path || v.track_cover_url || undefined
})
const nextSrc = computed(() => {
  const v = nextVideo.value
  if (!v) return ''
  return v.final_video_url || `/api/v1/videos/${v.id}/stream`
})

// ── Lifecycle ─────────────────────────────────────────────────────────

onMounted(() => {
  reducedMotion.value = window.matchMedia('(prefers-reduced-motion: reduce)').matches

  nextTick(() => { playCurrentVideo() })

  // Focus container for keyboard events
  containerRef.value?.focus()

  let topBarTimer = setTimeout(() => {
    showTopBar.value = false
  }, 2000)

  function resetTopBar() {
    showTopBar.value = true
    clearTimeout(topBarTimer)
    topBarTimer = setTimeout(() => {
      showTopBar.value = false
    }, 2000)
  }

  document.addEventListener('touchstart', resetTopBar, { passive: true })
  document.addEventListener('mousemove', resetTopBar, { passive: true })

  // Hide swipe hint after 5s
  setTimeout(() => { showSwipeHint.value = false }, 5000)

  countView()

  onUnmounted(() => {
    document.removeEventListener('touchstart', resetTopBar)
    document.removeEventListener('mousemove', resetTopBar)
    clearTimeout(topBarTimer)
    commentsAbort?.abort()
  })
})

watch(currentIndex, () => {
  countView()
  showOverflow.value = false
  // Reload comments for the new video
  if (showComments.value) loadComments()
})

// ── Keyboard Navigation ───────────────────────────────────────────────

function onKeyDown(e: KeyboardEvent) {
  if (showComments.value) {
    if (e.key === 'Escape') showComments.value = false
    return
  }
  switch (e.key) {
    case 'ArrowUp':
    case 'ArrowLeft':
      e.preventDefault()
      goToPrevious()
      break
    case 'ArrowDown':
    case 'ArrowRight':
      e.preventDefault()
      goToNext()
      break
    case 'Escape':
      handleClose()
      break
    case ' ':
    case 'm':
      e.preventDefault()
      toggleMute()
      break
  }
}

// ── Vertical Swipe ────────────────────────────────────────────────────

useVerticalSwipe({
  element: swipeRef as any,
  threshold: 60,
  onSwipeUp: () => {
    if (!hasNext.value) return
    if (reducedMotion.value) {
      goToNext()
      return
    }
    isTransitioning.value = true
    const animDuration = 250
    const startTime = performance.now()
    function animate(time: number) {
      const elapsed = time - startTime
      const progress = Math.min(1, elapsed / animDuration)
      translateY.value = -100 * progress
      if (progress < 1) {
        requestAnimationFrame(animate)
      } else {
        translateY.value = 0
        isTransitioning.value = false
        goToNext()
      }
    }
    requestAnimationFrame(animate)
  },
  onSwipeDown: () => {
    if (!hasPrevious.value) return
    if (reducedMotion.value) {
      goToPrevious()
      return
    }
    isTransitioning.value = true
    const animDuration = 250
    const startTime = performance.now()
    function animate(time: number) {
      const elapsed = time - startTime
      const progress = Math.min(1, elapsed / animDuration)
      translateY.value = 100 * progress
      if (progress < 1) {
        requestAnimationFrame(animate)
      } else {
        translateY.value = 0
        isTransitioning.value = false
        goToPrevious()
      }
    }
    requestAnimationFrame(animate)
  },
})

// ── Video playback ────────────────────────────────────────────────────

function playCurrentVideo() {
  const video = currentVideoRef.value
  if (!video) return
  video.currentTime = 0
  video.muted = isMuted.value
  video.play().catch(() => {})
}

function onVideoLoaded() {
  const video = currentVideoRef.value
  if (!video) return
  video.muted = isMuted.value
  video.play().catch(() => {})
}

function onEnded() {
  // Auto-advance to next on end
  if (hasNext.value) goToNext()
}

function onVideoError() {
  // Skip to next on error
  if (hasNext.value) goToNext()
}

function toggleMute() {
  isMuted.value = isMuted.value!
  if (currentVideoRef.value) {
    currentVideoRef.value.muted = isMuted.value
  }
  showMuteIndicator.value = true
  setTimeout(() => { showMuteIndicator.value = false }, 1000)
}

// ── Navigation ────────────────────────────────────────────────────────

function goToNext() {
  if (!hasNext.value) return
  isTransitioning.value = true
  currentIndex.value++
  nextTick(() => {
    playCurrentVideo()
    isTransitioning.value = false
  })
}

function goToPrevious() {
  if (!hasPrevious.value) return
  isTransitioning.value = true
  currentIndex.value--
  nextTick(() => {
    playCurrentVideo()
    isTransitioning.value = false
  })
}

function handleClose() {
  currentVideoRef.value?.pause()
  emit('close')
}

// ── Like ──────────────────────────────────────────────────────────────

async function toggleLike() {
  const video = currentVideo.value
  if (!video) return

  const previousState = isLiked.value
  const previousCount = video.like_count
  isLiked.value = isLiked.value!
  video.like_count += isLiked.value ? 1 : -1
  if (video.like_count < 0) video.like_count = 0

  likeScale.value = 1.1
  setTimeout(() => { likeScale.value = 1 }, 150)

  try {
    if (isLiked.value) {
      await videoApi.likeVideo(String(video.id))
    } else {
      await videoApi.unlikeVideo(String(video.id))
    }
  } catch {
    isLiked.value = previousState
    video.like_count = previousCount
  }
}

// ── View count ────────────────────────────────────────────────────────

function countView() {
  const video = currentVideo.value
  if (!video || viewed.value.has(String(video.id))) return
  viewed.value.add(String(video.id))
  videoApi.viewVideo(String(video.id)).catch(() => {})
}

// ── Actions ───────────────────────────────────────────────────────────

function handlePlayTrack(track: TrackSummary) {
  if (!track) return
  playerStore.playTrackById(String(track.id))
}

function shareOnTelegram() {
  const video = currentVideo.value
  if (!video?.track?.id) return
  const url = `${window.location.origin}/track/${String(video.track.id)}`
  window.open(`https://t.me/share/url?url=${encodeURIComponent(url)}&text=${encodeURIComponent(video.title)}`, '_blank')
  showOverflow.value = false
}

function goToTrack() {
  const video = currentVideo.value
  if (!video?.track?.id) return
  currentVideoRef.value?.pause()
  router.push(`/track/${String(video.track.id)}`)
  emit('close')
}

// ── Comments ──────────────────────────────────────────────────────────

async function loadComments() {
  const video = currentVideo.value
  if (!video) return
  commentsAbort?.abort()
  commentsAbort = new AbortController()
  try {
    const res = await useRequest<{ items: VideoComment[]; total: number }>(
      `/videos/${video.id}/comments`,
      { method: 'GET', signal: commentsAbort.signal } as any,
      { silent: true },
    )
    comments.value = res?.items || []
    commentCount.value = res?.total || comments.value.length
  } catch {
    // Silent fail
  }
}

async function submitComment() {
  const video = currentVideo.value
  if (!video || !newComment.value.trim() || submittingComment.value) return
  submittingComment.value = true
  try {
    const result = await useRequest<VideoComment>(
      `/videos/${video.id}/comments`,
      {
        method: 'POST',
        data: { content: newComment.value.trim() },
      },
      { silent: true },
    )
    if (result) {
      comments.value.unshift(result as VideoComment)
      commentCount.value++
      newComment.value = ''
    }
  } catch {
    // Silent fail
  } finally {
    submittingComment.value = false
  }
}

function timeAgo(dateStr: string): string {
  const now = Date.now()
  const then = new Date(dateStr).getTime()
  const diff = Math.floor((now - then) / 1000)
  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`
  return `${Math.floor(diff / 86400)}d`
}

// Load initial like state
watch(currentVideo, (v) => {
  if (v) isLiked.value = v.is_liked || false
}, { immediate: true })
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 200ms ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.like-btn {
  will-change: transform;
}

/* Slide-up panel for comments */
.slide-up-enter-active {
  transition: transform 0.3s ease-out;
}
.slide-up-leave-active {
  transition: transform 0.2s ease-in;
}
.slide-up-enter-from {
  transform: translateY(100%);
}
.slide-up-leave-to {
  transform: translateY(100%);
}

@media (prefers-reduced-motion: reduce) {
  .transition-transform {
    transition-duration: 0s !important;
  }
  .slide-up-enter-active,
  .slide-up-leave-active {
    transition: none;
  }
  .slide-up-enter-from,
  .slide-up-leave-to {
    transform: none;
  }
}
</style>
