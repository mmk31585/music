import { computed, type Ref } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayer } from '@/composables/player'
import { usePlayerApi, type PlaybackTrack } from '@/services/api/player'
import { useToast } from 'primevue/usetoast'
import { buildPlaybackTrack } from '@/factories/playbackTrack'

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
  /** If provided, these tracks are used for "Play next" / "Add to queue" bulk actions. */
  queue?: Ref<TrackContextItem[]>
}

/**
 * Builds PrimeVue ContextMenu items for a track.
 *
 * Usage in a component:
 * ```vue
 * <script setup>
 * import { useTrackContextMenu } from '@/composables/useTrackContextMenu'
 * const { model, show } = useTrackContextMenu(trackRef)
 * </script>
 *
 * <template>
 *   <ContextMenu :model="model" ref="ctxRef" />
 *   <div @contextmenu="show($event)">...</div>
 * </template>
 * ```
 */
export function useTrackContextMenu(
  trackRef: Ref<TrackContextItem | null | undefined>,
  _options?: UseTrackContextMenuOptions,
) {
  const router = useRouter()
  const player = usePlayer()
  const playerApi = usePlayerApi()
  const toast = useToast()

  // ── helpers ──────────────────────────────────────────────────────
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

  function toPlaybackTrack(t: TrackContextItem): PlaybackTrack {
    return buildPlaybackTrack({
      id: String(t.id),
      title: resolveTitle(t),
      artist_name: resolveArtistName(t),
      album_title: t.albumTitle || t.album_title || t.album?.title || null,
      cover_url: t.coverUrl || t.cover_url || t.cover || t.album?.coverUrl || t.album?.cover_url || null,
      duration_seconds: t.durationSeconds ?? t.duration_seconds ?? t.duration ?? null,
      streamUrl: playerApi.getTrackStreamUrl(String(t.id)),
    })
  }

  // ── actions ──────────────────────────────────────────────────────
  function playNow() {
    const t = trackRef.value
    if (!t) return
    player.toggleTrack(toPlaybackTrack(t))
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
    toast.add({ severity: 'success', summary: 'Added', detail: `"${pt.title}" will play next`, life: 2000 })
  }

  function addToQueue() {
    const t = trackRef.value
    if (!t) return
    const pt = toPlaybackTrack(t)
    player.updateQueue([...player.queue.value, pt])
    toast.add({ severity: 'success', summary: 'Queued', detail: `"${pt.title}" added to queue`, life: 2000 })
  }

  function goToTrack() {
    const t = trackRef.value
    if (!t) return
    router.push(`/track/${t.id}`)
  }

  function goToArtist() {
    const t = trackRef.value
    if (!t) return
    const artistId = resolveArtistId(t)
    if (artistId) router.push(`/artist/${artistId}`)
  }

  function goToAlbum() {
    const t = trackRef.value
    if (!t) return
    const albumId = resolveAlbumId(t)
    if (albumId) router.push(`/album/${albumId}`)
  }

  function shareTrack() {
    const t = trackRef.value
    if (!t) return
    const url = `${window.location.origin}/track/${t.id}`
    navigator.clipboard?.writeText(url).then(() => {
      toast.add({ severity: 'success', summary: 'Link copied', detail: `"${resolveTitle(t)}" link copied to clipboard`, life: 2500 })
    }).catch(() => {
      toast.add({ severity: 'error', summary: 'Failed', detail: 'Could not copy link', life: 3000 })
    })
  }

  // ── menu model (PrimeVue 4 MenuItem[]) ──────────────────────────
  const model = computed(() => {
    const t = trackRef.value
    const hasArtist = t ? resolveArtistId(t) !== null : false
    const hasAlbum = t ? resolveAlbumId(t) !== null : false

    return [
      {
        label: 'Play Now',
        icon: 'pi pi-play',
        command: playNow,
      },
      {
        label: 'Play Next',
        icon: 'pi pi-step-forward',
        command: playNext,
      },
      {
        label: 'Add to Queue',
        icon: 'pi pi-list',
        command: addToQueue,
      },
      { separator: true },
      {
        label: 'Go to Track',
        icon: 'pi pi-music',
        command: goToTrack,
      },
      ...(hasArtist
        ? [
            {
              label: 'Go to Artist',
              icon: 'pi pi-user',
              command: goToArtist,
            },
          ]
        : []),
      ...(hasAlbum
        ? [
            {
              label: 'Go to Album',
              icon: 'pi pi-book',
              command: goToAlbum,
            },
          ]
        : []),
      { separator: true },
      {
        label: 'Share Track',
        icon: 'pi pi-share-alt',
        command: shareTrack,
      },
    ]
  })

  /**
   * Call from a @contextmenu.prevent handler to show the menu at cursor position.
   */
  function _show(event: MouseEvent) {
    // The ContextMenu component's `show` method is called via ref in the consuming component.
    // This helper just prevents default and is a no-op placeholder — the real wiring
    // happens in the component template with `ctxRef.show(event)`.
    event.preventDefault()
  }

  return {
    /** PrimeVue ContextMenu model (reactive array of MenuItems). */
    model,
    /** Helper to build a PlaybackTrack from the current context track. */
    toPlaybackTrack,
    /** Individual actions if you need to call them outside the menu. */
    actions: {
      playNow,
      playNext,
      addToQueue,
      goToTrack,
      goToArtist,
      goToAlbum,
      shareTrack,
    },
  }
}
