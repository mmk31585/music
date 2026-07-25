import { computed, type Ref } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayer } from '@/composables/player'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'
import { useAppToast } from '@/composables/useAppToast'
import { useAlbumColors } from '@/composables/useAlbumColors'
import { buildPlaybackTrack } from '@/factories/playbackTrack'
import type { ContextMenuSection, ContextMenuHeader } from '@/types/context-menu'

/**
 * Generic track shape accepted by the context menu composable.
 * Covers all the field variants used across the codebase.
 */
export interface TrackContextItem {
  id: string | number
  title?: string | null
  artistName?: string | null
  artist_name?: string | null
  artist_id?: string | number | null
  artist?: { name?: string | null; id?: string | number } | null
  artists?: { name?: string | null; id?: string | number }[] | null
  coverUrl?: string | null
  cover_url?: string | null
  cover?: string | null
  album?: {
    title?: string | null
    coverUrl?: string | null
    cover_url?: string | null
    id?: string | number
  } | null
  album_id?: string | number | null
  albumTitle?: string | null
  album_title?: string | null
  durationSeconds?: number | null
  duration_seconds?: number | null
  duration?: number | null
}

interface UseTrackContextMenuOptions {
  /** If provided, used for "Add to queue" badge count. */
  queue?: Ref<TrackContextItem[]>
  /** Callback to open radio mode from the menu. */
  openRadio?: (trackId: string, seedLabel?: string) => void
  /** Callback to open the add-to-playlist picker. */
  openAddToPlaylist?: () => void
  /** Callback fired when the menu is closed. */
  onClose?: () => void
}

// ── Helpers ──────────────────────────────────────────────────────

function resolveArtistId(t: TrackContextItem): string | number | null {
  return t.artist_id ?? t.artist?.id ?? t.artists?.[0]?.id ?? null
}

function resolveAlbumId(t: TrackContextItem): string | number | null {
  return t.album_id ?? t.album?.id ?? null
}

function resolveArtistName(t: TrackContextItem): string {
  return (
    t.artistName ||
    t.artist_name ||
    t.artist?.name ||
    (Array.isArray(t.artists) && t.artists[0]?.name) ||
    'Unknown artist'
  )
}

function resolveTitle(t: TrackContextItem): string {
  return t.title || 'Untitled'
}

function resolveCoverUrl(t: TrackContextItem): string | null {
  return t.coverUrl || t.cover_url || t.cover || t.album?.coverUrl || t.album?.cover_url || null
}

function resolveAlbumTitle(t: TrackContextItem): string | null {
  return t.albumTitle || t.album_title || t.album?.title || null
}

function resolveDuration(t: TrackContextItem): number | null {
  return t.durationSeconds ?? t.duration_seconds ?? t.duration ?? null
}

// ── Composable ───────────────────────────────────────────────────

/**
 * Builds context menu sections and header for a track.
 *
 * Usage:
 * ```vue
 * <script setup>
 * import { useTrackContextMenu } from '@/composables/useTrackContextMenu'
 * const { sections, header, accentColor } = useTrackContextMenu(trackRef, {
 *   openRadio: (id, label) => openRadio(id, label),
 * })
 * </script>
 *
 * <template>
 *   <ContextMenu v-model:visible="menuVisible" :sections="sections" :header="header" :accent-color="accentColor" />
 * </template>
 * ```
 */
export function useTrackContextMenu(
  trackRef: Ref<TrackContextItem | null | undefined>,
  options?: UseTrackContextMenuOptions,
) {
  const router = useRouter()
  const player = usePlayer()
  const playerApi = usePlayerApi()
  const toast = useAppToast()

  // ── playback track builder ────────────────────────────────────

  function toPlaybackTrack(t: TrackContextItem): PlaybackTrack {
    return buildPlaybackTrack({
      id: String(t.id),
      title: resolveTitle(t),
      artist_name: resolveArtistName(t),
      album_title: resolveAlbumTitle(t),
      cover_url: resolveCoverUrl(t),
      duration_seconds: resolveDuration(t),
      streamUrl: playerApi.getTrackStreamUrl(String(t.id)),
    })
  }

  // ── album colors for dynamic accent ───────────────────────────

  const coverUrl = computed(() => trackRef.value ? resolveCoverUrl(trackRef.value) : null)
  const { palette } = useAlbumColors(coverUrl)

  // ── individual actions ────────────────────────────────────────

  function playNow() {
    const t = trackRef.value
    if (!t) return
    player.toggleTrack(toPlaybackTrack(t))
    options?.onClose?.()
  }

  function playNext() {
    const t = trackRef.value
    if (!t) return
    const pt = toPlaybackTrack(t)
    if (player.queue.value.length > 0) {
      const items = [...player.queue.value]
      items.splice(0, 0, pt)
      player.updateQueue(items)
    } else {
      player.updateQueue([pt])
    }
    toast.success(`"${pt.title}" will play next`)
  }

  function addToQueue() {
    const t = trackRef.value
    if (!t) return
    const pt = toPlaybackTrack(t)
    player.updateQueue([...player.queue.value, pt])
    toast.success(`"${pt.title}" added to queue`)
  }

  function goToTrack() {
    const t = trackRef.value
    if (!t) return
    router.push(`/track/${t.id}`)
    options?.onClose?.()
  }

  function goToArtist() {
    const t = trackRef.value
    if (!t) return
    const artistId = resolveArtistId(t)
    if (artistId) router.push(`/artist/${artistId}`)
    options?.onClose?.()
  }

  function goToAlbum() {
    const t = trackRef.value
    if (!t) return
    const albumId = resolveAlbumId(t)
    if (albumId) router.push(`/album/${albumId}`)
    options?.onClose?.()
  }

  function copyLink() {
    const t = trackRef.value
    if (!t) return
    const url = `${window.location.origin}/track/${t.id}`
    navigator.clipboard?.writeText(url).then(() => {
      toast.success(`"${resolveTitle(t)}" link copied to clipboard`)
    }).catch(() => {
      toast.error('Could not copy link')
    })
  }

  function shareTrack() {
    const t = trackRef.value
    if (!t) return
    const url = `${window.location.origin}/track/${t.id}`
    if (navigator.share) {
      navigator.share({ title: resolveTitle(t), url }).catch(() => {})
    } else {
      copyLink()
    }
  }

  function startRadio() {
    const t = trackRef.value
    if (!t) return
    options?.openRadio?.(String(t.id), `${resolveTitle(t)} • ${resolveArtistName(t)}`)
    options?.onClose?.()
  }

  function reportIssue() {
    toast.info('Issue reporting will be available in a future update')
    options?.onClose?.()
  }

  // ── sections (reactive) ───────────────────────────────────────

  const sections = computed<ContextMenuSection[]>(() => {
    const t = trackRef.value
    if (!t) return []

    const hasArtist = resolveArtistId(t) !== null
    const hasAlbum = resolveAlbumId(t) !== null
    const queueCount = options?.queue?.value?.length ?? 0

    return [
      // ── Playback ─────────────────────────────────────────
      {
        id: 'playback',
        label: 'PLAYBACK',
        items: [
          {
            id: 'play-now',
            label: 'Play Now',
            icon: 'Play',
            action: playNow,
          },
          {
            id: 'play-next',
            label: 'Play Next',
            icon: 'SkipForward',
            action: playNext,
          },
          {
            id: 'add-to-queue',
            label: 'Add to Queue',
            icon: 'ListMusic',
            badge: queueCount > 0 ? `${queueCount} songs` : undefined,
            action: addToQueue,
          },
          {
            id: 'start-radio',
            label: 'Start Radio',
            icon: 'Radio',
            separator: true,
            action: startRadio,
          },
        ],
      },

      // ── Library ──────────────────────────────────────────
      {
        id: 'library',
        label: 'LIBRARY',
        items: [
          {
            id: 'add-to-playlist',
            label: 'Add to Playlist',
            icon: 'PlusCircle',
            separator: true,
            action: options?.openAddToPlaylist ?? (() => {
              toast.info('Playlist picker coming soon')
            }),
          },
        ],
      },

      // ── Go To ────────────────────────────────────────────
      {
        id: 'navigate',
        label: 'GO TO',
        items: [
          {
            id: 'go-to-track',
            label: 'Go to Track',
            icon: 'Music2',
            action: goToTrack,
          },
          ...(hasArtist
            ? [{
                id: 'go-to-artist',
                label: 'Go to Artist',
                icon: 'UserRound',
                action: goToArtist,
              }]
            : []),
          ...(hasAlbum
            ? [{
                id: 'go-to-album',
                label: 'Go to Album',
                icon: 'Disc3',
                action: goToAlbum,
              }]
            : []),
        ],
      },

      // ── Share ────────────────────────────────────────────
      {
        id: 'share',
        label: 'SHARE',
        items: [
          {
            id: 'copy-link',
            label: 'Copy Link',
            icon: 'Link2',
            shortcut: 'Ctrl+C',
            separator: true,
            action: copyLink,
          },
          {
            id: 'share',
            label: 'Share',
            icon: 'Share2',
            action: shareTrack,
          },
        ],
      },

      // ── Advanced ─────────────────────────────────────────
      {
        id: 'advanced',
        label: 'ADVANCED',
        items: [
          {
            id: 'track-info',
            label: 'Track Info',
            icon: 'Info',
            shortcut: 'I',
            separator: true,
            action: goToTrack,
          },
          {
            id: 'report',
            label: 'Report',
            icon: 'Flag',
            danger: true,
            action: reportIssue,
          },
        ],
      },
    ]
  })

  // ── header (reactive) ─────────────────────────────────────────

  const header = computed<ContextMenuHeader | undefined>(() => {
    const t = trackRef.value
    if (!t) return undefined

    const trackId = String(t.id)
    const isPlaying = player.currentTrack.value?.id === trackId
      && player.isPlaying.value

    return {
      coverUrl: resolveCoverUrl(t) ?? undefined,
      title: resolveTitle(t),
      artistName: resolveArtistName(t),
      artistId: resolveArtistId(t) ?? undefined,
      albumName: resolveAlbumTitle(t) ?? undefined,
      albumId: resolveAlbumId(t) ?? undefined,
      duration: resolveDuration(t) ?? undefined,
      isPlaying,
    }
  })

  return {
    /** Ordered sections of menu actions (reactive). */
    sections,
    /** Premium track header for the menu. */
    header,
    /** Dynamic album accent color for hover/focus tints. */
    accentColor: computed(() => palette.value.vibrant),
    /** Build a PlaybackTrack from the context item. */
    toPlaybackTrack,
    /** Individual actions callable outside the menu. */
    actions: {
      playNow,
      playNext,
      addToQueue,
      goToTrack,
      goToArtist,
      goToAlbum,
      copyLink,
      shareTrack,
      startRadio,
    },
  }
}
