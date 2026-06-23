<template>
  <div class="min-h-screen">
    <!-- ═══ LOADING ═══ -->
    <div v-if="loading" class="mx-auto max-w-6xl px-4 pt-4 md:px-6 md:pt-8">
      <div class="mx-auto aspect-video max-w-5xl animate-pulse rounded-2xl bg-white/[0.03]" />
      <div class="mt-6 flex flex-col gap-6 lg:flex-row">
        <div class="flex-1 space-y-4">
          <div class="h-8 w-3/4 animate-pulse rounded-lg bg-white/[0.06]" />
          <div class="h-4 w-1/3 animate-pulse rounded bg-white/[0.04]" />
        </div>
        <div class="w-full lg:w-80">
          <div class="grid grid-cols-2 gap-3">
            <div v-for="i in 4" :key="i" class="animate-pulse">
              <div class="aspect-video w-full rounded-xl bg-white/[0.04]" />
              <div class="mt-2 space-y-1.5 px-1">
                <div class="h-3 w-3/4 rounded bg-white/[0.06]" />
                <div class="h-2.5 w-1/2 rounded bg-white/[0.04]" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ═══ ERROR ═══ -->
    <div v-else-if="error" class="flex flex-col items-center justify-center px-4 py-32 text-center">
      <div class="mb-6 flex h-20 w-20 items-center justify-center rounded-full border border-amber-500/20 bg-amber-500/5">
        <i class="pi pi-video text-3xl text-amber-400/60" />
      </div>
      <h2 class="mb-2 text-2xl font-black text-white/60">Video not found</h2>
      <p class="mb-8 text-sm text-white/30">This music video doesn't exist or is still being processed.</p>
      <div class="flex gap-3">
        <button
          class="rounded-xl border border-white/[0.08] bg-white/[0.04] px-5 py-2.5 text-sm font-semibold text-white/60 transition hover:bg-white/[0.08] hover:text-white"
          @click="goBack"
        >
          <i class="pi pi-arrow-right ml-2" />
          Go Back
        </button>
        <button
          class="rounded-xl bg-amber-500/10 px-5 py-2.5 text-sm font-semibold text-amber-400 transition hover:bg-amber-500/20"
          @click="router.push('/music-videos')"
        >
          <i class="pi pi-th-large ml-2" />
          Browse Videos
        </button>
      </div>
    </div>

    <!-- ═══ VIDEO CONTENT ═══ -->
    <div v-else-if="video">

      <!-- ══ PLAYER ══ -->
      <div class="mx-auto max-w-5xl px-0 md:px-6">
        <div
          ref="playerContainerRef"
          class="group relative w-full overflow-hidden bg-black md:rounded-2xl"
          :class="isFullscreen ? 'fixed inset-0 z-50 !rounded-none' : 'aspect-video'"
          @mousemove="showControls"
          @mouseleave="startHideTimer"
        >
        <video
          ref="videoRef"
          :src="videoSrc"
          :poster="videoPoster"
          class="h-full w-full cursor-pointer object-contain"
          :class="{ 'cursor-pointer': !isPlaying }"
          :playsinline="true"
          preload="metadata"
          @timeupdate="onTimeUpdate"
          @loadedmetadata="onLoaded"
          @ended="onEnded"
          @waiting="buffering = true"
          @canplay="buffering = false"
          @play="onPlay"
          @pause="onPause"
          @error="onVideoError"
          @click.prevent="togglePlay"
        />

        <!-- Big play -->
        <Transition name="fade">
          <div
            v-if="!isPlaying && !buffering && !videoError"
            class="absolute inset-0 flex cursor-pointer items-center justify-center bg-black/30"
            @click="togglePlay"
          >
            <div class="flex h-20 w-20 items-center justify-center rounded-full bg-white/5 text-white shadow-2xl backdrop-blur-2xl transition-all hover:scale-110 hover:bg-white/10 md:h-24 md:w-24">
              <i class="pi pi-play-fill ml-1 text-4xl md:text-5xl" />
            </div>
          </div>
        </Transition>

        <!-- Buffering -->
        <Transition name="fade">
          <div
            v-if="buffering"
            class="absolute inset-0 flex items-center justify-center bg-black/40"
          >
            <div class="flex flex-col items-center gap-3">
              <i class="pi pi-spin pi-spinner text-2xl text-white/40" />
              <span class="text-xs tracking-widest text-white/30 uppercase">Buffering</span>
            </div>
          </div>
        </Transition>

        <!-- Video error -->
        <Transition name="fade">
          <div
            v-if="videoError"
            class="absolute inset-0 flex flex-col items-center justify-center bg-black/70 gap-4"
          >
            <div class="flex h-16 w-16 items-center justify-center rounded-full border border-red-500/20 bg-red-500/5">
              <i class="pi pi-exclamation-triangle text-xl text-red-400/50" />
            </div>
            <p class="text-sm text-white/40 text-center max-w-xs">This video can't be played right now.</p>
            <button
              class="rounded-xl bg-white/10 px-4 py-2 text-xs font-semibold text-white/70 hover:bg-white/20"
              @click="retryVideo"
            >
              <i class="pi pi-refresh ml-1.5" /> Retry
            </button>
          </div>
        </Transition>

        <!-- ══ Controls overlay ══ -->
        <Transition name="fade">
          <div
            v-if="(controlsVisible || !isPlaying) && !videoError"
            class="absolute inset-0 flex flex-col justify-end bg-gradient-to-t from-black/90 via-black/10 to-transparent"
          >
            <div class="px-4 pb-1">
              <div class="group relative cursor-pointer py-2" @click="seekTo">
                <div class="h-1 overflow-hidden rounded-full bg-white/15 transition-all group-hover:h-1.5">
                  <div
                    class="h-full rounded-full bg-gradient-to-r from-amber-400 to-amber-500 transition-all"
                    :style="{ width: progressPercent + '%' }"
                  />
                </div>
              </div>
            </div>
            <div class="flex items-center gap-3 px-4 pb-4">
              <button
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-white/80 hover:bg-white/10"
                @click="togglePlay"
                aria-label="Play / Pause"
              >
                <i :class="isPlaying ? 'pi pi-pause-fill' : 'pi pi-play-fill'" class="text-lg" />
              </button>
              <span class="min-w-[4.5rem] text-xs font-medium text-white/60 tabular-nums">
                {{ formatTime(currentTime) }} / {{ formatTime(duration) }}
              </span>
              <div class="hidden items-center gap-1.5 sm:flex">
                <button
                  class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 hover:bg-white/10 hover:text-white"
                  @click="toggleMute"
                  :aria-label="isMuted ? 'Unmute' : 'Mute'"
                >
                  <i :class="isMuted ? 'pi pi-volume-off' : 'pi pi-volume-up'" class="text-sm" />
                </button>
                <div class="relative h-1 w-20 cursor-pointer overflow-hidden rounded-full bg-white/10" @click="seekVolume">
                  <div class="h-full rounded-full bg-white/50 transition-all" :style="{ width: volumePercent + '%' }" />
                </div>
              </div>
              <div class="flex-1" />
              <button
                class="flex h-8 w-8 items-center justify-center rounded-full transition"
                :class="isLiked ? 'text-red-400' : 'text-white/40 hover:bg-white/10 hover:text-white'"
                @click="toggleLike"
                :aria-label="isLiked ? 'Unlike video' : 'Like video'"
              >
                <i :class="isLiked ? 'pi pi-heart-fill' : 'pi pi-heart'" class="text-sm" />
              </button>
              <button
                class="flex h-8 w-8 items-center justify-center rounded-full text-white/40 hover:bg-white/10 hover:text-white"
                @click="toggleFullscreen"
                :aria-label="isFullscreen ? 'Exit fullscreen' : 'Fullscreen'"
              >
                <i :class="isFullscreen ? 'pi pi-window-minimize' : 'pi pi-window-maximize'" class="text-sm" />
              </button>
            </div>
          </div>
        </Transition>
      </div>
      </div>

      <!-- ══ BELOW PLAYER ══ -->
      <div class="mx-auto max-w-6xl px-4 pt-6 pb-20 md:px-6">
        <div class="flex flex-col gap-6 lg:flex-row">

          <!-- ── Main: Video Info ── -->
          <div class="min-w-0 flex-1 space-y-5">
            <div>
              <h1 class="text-xl font-black text-white md:text-2xl">{{ video.title }}</h1>
              <div class="mt-2 flex flex-wrap items-center gap-x-4 gap-y-1.5 text-sm text-white/40">
                <span>{{ formatCount(video.view_count) }} views</span>
                <span class="h-1 w-1 rounded-full bg-white/10" />
                <span>{{ formatCount(video.like_count) }} likes</span>
                <span class="h-1 w-1 rounded-full bg-white/10" />
                <span>{{ timeAgo(video.created_at) }}</span>
                <span
                  class="rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase tracking-wider"
                  :class="video.type === 'official_mv'
                    ? 'bg-amber-500/10 text-amber-400 border border-amber-500/15'
                    : 'bg-blue-500/10 text-blue-400 border border-blue-500/15'"
                >
                  {{ video.type === 'official_mv' ? 'Official MV' : 'Fan Edit' }}
                </span>
              </div>
              <p v-if="video.description" class="mt-3 text-sm leading-relaxed text-white/30 line-clamp-2">
                {{ video.description }}
              </p>
              <div class="mt-3 flex items-center gap-2">
                <button
                  class="flex items-center gap-2 rounded-lg border border-white/[0.06] bg-white/[0.03] px-3 py-1.5 text-xs font-semibold text-white/40 hover:border-white/10 hover:bg-white/[0.06] hover:text-white/70 transition"
                  @click="shareVideo"
                >
                  <i class="pi pi-share-alt" /> Share
                </button>
                <RouterLink
                  v-if="video.track_id"
                  :to="`/track/${video.track_id}`"
                  class="flex items-center gap-2 rounded-lg border border-white/[0.06] bg-white/[0.03] px-3 py-1.5 text-xs font-semibold text-white/40 hover:border-white/10 hover:bg-white/[0.06] hover:text-white/70 transition"
                >
                  <i class="pi pi-music" /> View Track
                </RouterLink>
              </div>
            </div>

            <!-- Track Card (compact) -->
            <div v-if="trackInfo" class="flex items-center gap-4 rounded-2xl border border-white/[0.06] bg-white/[0.02] p-4">
              <div class="relative h-14 w-14 shrink-0 overflow-hidden rounded-xl bg-white/10 shadow-lg">
                <img v-if="trackInfo.cover_url" :src="trackInfo.cover_url" :alt="trackInfo.title" class="h-full w-full object-cover" />
                <div v-else class="flex h-full items-center justify-center"><i class="pi pi-music text-lg text-slate-500" /></div>
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-[10px] font-bold tracking-wider text-white/30 uppercase">Track</p>
                <p class="truncate text-sm font-semibold text-white">{{ trackInfo.title }}</p>
                <p class="truncate text-xs text-white/40">{{ trackInfo.artist_name || 'Unknown artist' }}</p>
              </div>
              <button
                type="button"
                class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-amber-500 text-black shadow-lg shadow-amber-500/20 transition-all hover:scale-105 hover:bg-amber-400 active:scale-95"
                :title="isTrackPlaying ? 'Now Playing' : 'Play Track'"
                @click="playTrack"
              >
                <i :class="isTrackPlaying ? 'pi pi-check' : 'pi pi-play-fill'" class="ml-0.5 text-sm" />
              </button>
            </div>

            <!-- Comments -->
            <div>
              <h3 class="mb-3 text-xs font-bold tracking-wider text-white/40 uppercase">
                Comments <span v-if="commentsTotal > 0" class="font-normal text-white/20">({{ commentsTotal }})</span>
              </h3>
              <div v-if="auth.isAuthenticated" class="mb-4">
                <textarea
                  v-model="commentContent"
                  placeholder="Write a comment…"
                  class="w-full resize-none rounded-xl border border-white/[0.06] bg-white/[0.03] px-4 py-3 text-sm text-white/80 placeholder:text-white/20 outline-none focus:border-white/20 focus:bg-white/[0.06]"
                  rows="2"
                  maxlength="1000"
                  @keydown.ctrl.enter="postComment"
                />
                <div class="mt-1.5 flex items-center justify-between">
                  <span class="text-[10px] text-white/20">{{ commentContent.length }}/1000</span>
                  <Button label="Post" size="small" :loading="submittingComment" :disabled="!commentContent.trim() || submittingComment"
                    class="!rounded-xl !bg-white/10 !text-white hover:!bg-white/20 !text-xs !py-1.5 !px-4" @click="postComment" />
                </div>
              </div>
              <p v-else class="mb-4 text-sm text-white/30">
                <RouterLink to="/login" class="text-white/60 underline hover:text-white">Log in</RouterLink> to comment.
              </p>
              <div v-if="loadingComments" class="space-y-3">
                <div v-for="i in 3" :key="i" class="flex animate-pulse gap-3">
                  <div class="h-8 w-8 rounded-full bg-white/[0.06]" />
                  <div class="flex-1 space-y-2"><div class="h-3 w-24 rounded bg-white/[0.06]" /><div class="h-4 w-3/4 rounded bg-white/[0.04]" /></div>
                </div>
              </div>
              <div v-else-if="comments.length" class="space-y-4">
                <div v-for="cm in comments" :key="cm.id" class="flex gap-3">
                  <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-white/[0.06] text-xs font-bold text-white/40">{{ cm.author.charAt(0).toUpperCase() }}</div>
                  <div class="min-w-0 flex-1">
                    <div class="flex items-center gap-2"><span class="text-xs font-semibold text-white/70">{{ cm.author }}</span><span class="text-[10px] text-white/20">{{ timeAgo(cm.created_at) }}</span></div>
                    <p class="mt-0.5 text-sm leading-relaxed text-white/60 break-words">{{ cm.content }}</p>
                  </div>
                </div>
              </div>
              <div v-else class="py-8 text-center text-sm text-white/20"><i class="pi pi-comment mb-2 block text-xl text-white/10" />No comments yet.</div>
            </div>
          </div>

          <!-- ── Sidebar: Related Videos ── -->
          <div class="w-full lg:w-80 shrink-0">
            <h3 class="mb-3 text-xs font-bold tracking-wider text-white/40 uppercase">
              Related <span v-if="relatedVideos.length" class="font-normal text-white/20">({{ relatedVideos.length }})</span>
            </h3>
            <div class="grid grid-cols-2 gap-3 lg:grid-cols-1">
              <div
                v-for="rel in relatedVideos"
                :key="String(rel.id)"
                class="group flex cursor-pointer gap-3 rounded-xl p-2 transition hover:bg-white/[0.04]"
                @click="openVideo(rel)"
              >
                <div class="aspect-video w-40 shrink-0 overflow-hidden rounded-lg bg-white/10">
                  <img
                    v-if="rel.thumbnail_url || rel.thumbnail_path || rel.track_cover_url"
                    :src="(rel.thumbnail_url || rel.thumbnail_path || rel.track_cover_url) ?? undefined"
                    :alt="rel.title"
                    class="h-full w-full object-cover transition-transform group-hover:scale-[1.03]"
                    loading="lazy"
                  />
                  <div v-else class="flex h-full items-center justify-center"><i class="pi pi-video text-slate-500" /></div>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="line-clamp-2 text-sm font-medium leading-snug text-white/90">{{ rel.title }}</p>
                  <p class="mt-0.5 truncate text-xs text-white/40">{{ formatCount(rel.view_count) }} views</p>
                </div>
              </div>
            </div>
            <div v-if="!relatedVideos.length" class="py-8 text-center text-sm text-white/20">
              <i class="pi pi-video mb-2 block text-xl text-white/10" />No related videos
            </div>
          </div>
        </div>
      </div>
    </div>

    <AddToPlaylistDialog
      v-if="trackInfo?.id"
      v-model:visible="showAddToPlaylist"
      :track-id="String(trackInfo.id)"
      :track-title="trackInfo.title"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import { useRequest } from '@/composables/useRequest'
import { useVideoApi } from '@/services/api/video'
import { useTracksApi } from '@/services/api/catalog/tracks'
import { useLibraryApi } from '@/services/api/library'
import { usePlayer } from '@/composables/player'
import { useUserAuthStore } from '@/stores'
import { AddToPlaylistDialog } from '@/components/music'
import { formatCount } from '@/utils/number'
import { useSocialShare } from '@/composables/social'
import type { VideoItem } from '@/services/api/video/types'
import type { Track } from '@/services/api/catalog/tracks'

const route = useRoute()
const router = useRouter()
const toast = useToast()
const videoApi = useVideoApi()
const tracksApi = useTracksApi()
const libraryApi = useLibraryApi()
const player = usePlayer()
const auth = useUserAuthStore()
const { copyLink } = useSocialShare()

const videoId = computed(() => route.params.videoId as string)

// ── Data ──
const video = ref<VideoItem | null>(null)
const relatedVideos = ref<VideoItem[]>([])
const loading = ref(true)
const error = ref(false)

const trackInfo = ref<Pick<Track, 'id' | 'title' | 'artist_name' | 'cover_url' | 'duration_seconds'> | null>(null)
const loadingTrack = ref(false)

const isTrackPlaying = computed(() => {
  if (!trackInfo.value) return false
  return String(player.currentTrack.value?.id) === String(trackInfo.value.id) && player.isPlaying.value
})

const trackLiked = ref(false)
const showAddToPlaylist = ref(false)

// ── Comments ──
interface CommentItem { id: string; author: string; content: string; created_at: string }
const comments = ref<CommentItem[]>([])
const commentsTotal = ref(0)
const commentContent = ref('')
const submittingComment = ref(false)
const loadingComments = ref(false)

// ── Player ──
const videoRef = ref<HTMLVideoElement | null>(null)
const playerContainerRef = ref<HTMLDivElement | null>(null)
const isPlaying = ref(false)
const isMuted = ref(false)
const isFullscreen = ref(false)
const currentTime = ref(0)
const buffering = ref(false)
const controlsVisible = ref(true)
const volume = ref(1)
const isLiked = ref(false)
const videoError = ref(false)
let controlsTimer: ReturnType<typeof setTimeout> | null = null
let hideTimer: ReturnType<typeof setTimeout> | null = null

const duration = computed(() => (video.value ? video.value.duration_ms / 1000 : 0))
const progressPercent = computed(() => { if (!duration.value) return 0; return (currentTime.value / duration.value) * 100 })
const volumePercent = computed(() => volume.value * 100)

const videoSrc = computed(() => {
  if (!video.value) return undefined
  if (video.value.final_video_url) return video.value.final_video_url
  return `/api/v1/videos/${video.value.id}/stream`
})

const videoPoster = computed(() => {
  if (!video.value) return undefined
  return video.value.thumbnail_url || video.value.thumbnail_path || video.value.track_cover_url || (trackInfo.value?.cover_url || undefined) || undefined
})

function formatTime(s: number): string { if (!s || !isFinite(s)) return '0:00'; return `${Math.floor(s / 60)}:${Math.floor(s % 60).toString().padStart(2, '0')}` }
function timeAgo(d: string): string { const diff = Math.floor((Date.now() - new Date(d).getTime()) / 1000); if (diff < 60) return 'just now'; if (diff < 3600) return `${Math.floor(diff / 60)}m ago`; if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`; return `${Math.floor(diff / 86400)}d ago` }

async function fetchVideo() {
  loading.value = true; error.value = false
  try {
    const data = await useRequest<VideoItem>(`/videos/${videoId.value}`, { method: 'GET' })
    video.value = data; isLiked.value = data.is_liked || false
    fetchRelated(); fetchComments()
    if (data.track_id) fetchTrackInfo(data.track_id)
    else if (data.track?.id) fetchTrackInfo(String(data.track.id))
  } catch { error.value = true } finally { loading.value = false }
}

async function fetchTrackInfo(trackId: string) {
  loadingTrack.value = true
  try {
    const track = await tracksApi.getTrack(trackId, { silent: true } as any)
    trackInfo.value = { id: track.id, title: track.title, artist_name: (track as any).artist_name || null, cover_url: (track as any).cover_url || null, duration_seconds: (track as any).duration_seconds || null }
  } catch { /* silent */ } finally { loadingTrack.value = false }
}

async function fetchRelated() {
  const tid = video.value?.track_id || video.value?.track?.id
  if (!tid) return
  try {
    const res = await useRequest<{ items: VideoItem[] }>(`/tracks/${tid}/videos`, { method: 'GET' })
    relatedVideos.value = (res?.items || []).filter(v => String(v.id) !== videoId.value)
  } catch { relatedVideos.value = [] }
}

function togglePlay() {
  if (videoError.value) return; const el = videoRef.value; if (!el) return
  el.paused ? el.play().catch(() => {}) : el.pause()
}
function onPlay() { isPlaying.value = true; startHideTimer() }
function onPause() { isPlaying.value = false; cancelHideTimer(); controlsVisible.value = true }
function onVideoError() { videoError.value = true; isPlaying.value = false; buffering.value = false }
function retryVideo() { videoError.value = false; buffering.value = true; videoRef.value?.load() }
function toggleMute() { const el = videoRef.value; if (!el) return; el.muted = !el.muted; isMuted.value = el.muted }
function seekTo(e: MouseEvent) { const el = videoRef.value; if (!el || !duration.value) return; const r = (e.currentTarget as HTMLElement).getBoundingClientRect(); el.currentTime = ((e.clientX - r.left) / r.width) * duration.value }
function seekVolume(e: MouseEvent) { const el = videoRef.value; if (!el) return; const pct = Math.max(0, Math.min(1, (e.clientX - (e.currentTarget as HTMLElement).getBoundingClientRect().left) / (e.currentTarget as HTMLElement).offsetWidth)); volume.value = pct; el.volume = pct; if (pct > 0 && el.muted) { el.muted = false; isMuted.value = false } }
function toggleFullscreen() { const c = playerContainerRef.value; if (!c) return; !document.fullscreenElement ? c.requestFullscreen().catch(() => {}) : document.exitFullscreen().catch(() => {}) }
function onTimeUpdate() { const el = videoRef.value; if (el) currentTime.value = el.currentTime }
function onLoaded() { buffering.value = false; videoError.value = false; const el = videoRef.value; if (el) el.volume = volume.value }
function onEnded() { isPlaying.value = false; controlsVisible.value = true }
function showControls() { controlsVisible.value = true; cancelHideTimer(); if (isPlaying.value) controlsTimer = setTimeout(() => { controlsVisible.value = false }, 3000) }
function startHideTimer() { cancelHideTimer(); controlsTimer = setTimeout(() => { if (isPlaying.value) controlsVisible.value = false }, 3000) }
function cancelHideTimer() { if (controlsTimer) { clearTimeout(controlsTimer); controlsTimer = null } }

async function toggleLike() {
  if (!video.value) return; const prev = isLiked.value; isLiked.value = !isLiked.value
  try { isLiked.value ? await videoApi.likeVideo(String(video.value.id)) : await videoApi.unlikeVideo(String(video.value.id)) }
  catch { isLiked.value = prev; toast.add({ severity: 'error', summary: 'Failed to update like', life: 3000 }) }
}
function playTrack() { const tid = trackInfo.value?.id || video.value?.track_id || video.value?.track?.id; if (tid) player.playTrackById(String(tid)) }
function shareVideo() { if (!video.value) return; copyLink({ id: String(video.value.id), title: video.value.title, type: 'video' }) }
function openVideo(rel: VideoItem) { if (String(rel.id) === videoId.value) return; router.push(`/music-video/${rel.id}`) }
function goBack() { router.back() }

async function fetchComments() {
  if (!video.value?.id) return; loadingComments.value = true
  try { const r = await useRequest<{ items: CommentItem[]; total: number }>(`/videos/${video.value.id}/comments?limit=20&offset=0`, { method: 'GET' }); comments.value = r?.items || []; commentsTotal.value = r?.total || 0 }
  catch { comments.value = [] } finally { loadingComments.value = false }
}
async function postComment() {
  if (!video.value?.id || !commentContent.value.trim()) return; submittingComment.value = true
  try { await useRequest(`/videos/${video.value.id}/comments`, { method: 'POST', data: { content: commentContent.value.trim() } }, { silent: true }); commentContent.value = ''; toast.add({ severity: 'success', summary: 'Comment posted', life: 3000 }); fetchComments() }
  catch { toast.add({ severity: 'error', summary: 'Failed to post comment', life: 3000 }) } finally { submittingComment.value = false }
}

function onFullscreenChange() { isFullscreen.value = !!document.fullscreenElement }
function onKeydown(e: KeyboardEvent) {
  if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return
  switch (e.key) {
    case ' ': case 'k': e.preventDefault(); togglePlay(); break
    case 'f': toggleFullscreen(); break
    case 'm': toggleMute(); break
    case 'ArrowLeft': if (videoRef.value) videoRef.value.currentTime = Math.max(0, (videoRef.value.currentTime || 0) - 5); break
    case 'ArrowRight': if (videoRef.value) videoRef.value.currentTime = Math.min(duration.value, (videoRef.value.currentTime || 0) + 5); break
  }
}

onMounted(() => { fetchVideo(); document.addEventListener('fullscreenchange', onFullscreenChange); document.addEventListener('keydown', onKeydown) })
onUnmounted(() => { document.removeEventListener('fullscreenchange', onFullscreenChange); document.removeEventListener('keydown', onKeydown); cancelHideTimer(); if (hideTimer) clearTimeout(hideTimer) })
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.line-clamp-2 { display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
</style>
