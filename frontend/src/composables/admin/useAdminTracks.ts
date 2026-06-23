import { ref } from 'vue'
import {
  useTracksApi,
  type Track,
  type TrackCreatePayload,
  type TrackUpdatePayload,
  type TrackArtistRequest,
} from '@/services/api/catalog/tracks'
import type { UploadResponse } from '@/services/api/media'
import { useLyricsApi } from '@/services/api/lyrics'

export type TrackCreditPayload = {
  artist_id: string | number
  role: string
  position?: number
}

export type TrackFormPayload = {
  title: string
  artist_id?: string | number
  album_id?: string | number | null
  album_artist_id?: string | number | null
  duration_seconds?: number | null
  audio_url?: string | null
  cover_url?: string | null

  artist_ids?: Array<string | number>
  featured_artist_ids?: Array<string | number>
  credits?: TrackCreditPayload[]
  artists?: TrackCreditPayload[]
  genre_ids?: Array<string | number>

  track_number?: number | null
  disc_number?: number | null
  year?: number | null
  composer?: string | null
  album_artist?: string | null
  explicit?: boolean | null
  isrc?: string | null
  language?: string | null
  release_date?: string | null
  label?: string | null

  lyrics?: string | null
  lyrics_language?: string | null
  lyrics_type?: 'plain' | 'synced' | null

  audioFile?: File | null
  coverFile?: File | null
}

function buildArtistsFromPayload(payload: TrackFormPayload): TrackArtistRequest[] {
  if (payload.credits?.length) {
    return payload.credits.map((credit, index) => ({
      artistId: credit.artist_id,
      role: credit.role,
      position: index,
    }))
  }

  if (payload.artists?.length) {
    return payload.artists.map((artist, index) => ({
      artistId: artist.artist_id,
      role: artist.role || (index === 0 ? 'primary' : 'featured'),
      position: artist.position ?? index,
    }))
  }

  const primaryArtistId = payload.artist_ids?.[0] ?? payload.artist_id

  if (!primaryArtistId) {
    throw new Error('artist_id is required')
  }

  return [
    {
      artistId: primaryArtistId,
      role: 'primary',
      position: 0,
    },
  ]
}

function toCreatePayload(payload: TrackFormPayload): TrackCreatePayload {
  const primaryArtistId = payload.artist_ids?.[0] ?? payload.artist_id

  if (!primaryArtistId) {
    throw new Error('artist_id is required')
  }

  return {
    title: payload.title,
    artistId: primaryArtistId,
    artists: buildArtistsFromPayload(payload),
    albumId: payload.album_id ?? null,
    durationSeconds: payload.duration_seconds ?? null,
    audioUrl: payload.audio_url ?? null,
    coverUrl: payload.cover_url ?? null,
    genreIds: payload.genre_ids ?? [],
    trackNumber: payload.track_number ?? null,
    explicit: payload.explicit ?? false,
    isPublic: true,
  }
}


function toUpdatePayload(payload: TrackFormPayload): TrackUpdatePayload {
  const primaryArtistId = payload.artist_ids?.[0] ?? payload.artist_id

  return {
    title: payload.title,
    artistId: primaryArtistId ?? null,
    artists: buildArtistsFromPayload(payload),
    albumId: payload.album_id ?? null,
    durationSeconds: payload.duration_seconds ?? null,
    audioUrl: payload.audio_url ?? null,
    coverUrl: payload.cover_url ?? null,
    genreIds: payload.genre_ids ?? [],
    trackNumber: payload.track_number ?? null,
    explicit: payload.explicit ?? false,
    isPublic: true,
  }
}

function getUploadedAudioUrl(uploaded: UploadResponse): string | null {
  return uploaded.url ?? uploaded.file_url ?? uploaded.fileUrl ?? null
}

function getUploadedCoverUrl(uploaded: UploadResponse): string | null {
  return uploaded.url ?? uploaded.file_url ?? uploaded.fileUrl ?? null
}

export function useAdminTracks() {
  const {
    adminGetTracks,
    adminCreateTrack,
    adminUploadTrackAudio,
    adminUploadTrackCover,
    adminUpdateTrack,
    adminDeleteTrack,
  } = useTracksApi()

  const { getTrackLyrics, getLyricsByTrackID, adminCreateLyrics, adminUpdateLyrics, adminDeleteLyrics } = useLyricsApi()

  const tracks = ref<Track[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const deleting = ref(false)
  const error = ref<any>(null)

  async function fetchTracks() {
    loading.value = true
    error.value = null

    try {
      const response = await adminGetTracks()
      tracks.value = response
      return response
    } catch (err) {
      error.value = err
      throw err
    } finally {
      loading.value = false
    }
  }

  async function syncTrackLyrics(trackId: string | number, payload: TrackFormPayload) {
    const content = payload.lyrics?.trim()

    // If no content, delete all existing lyrics and bail
    if (!content) {
      try {
        const allLyrics = await getLyricsByTrackID(trackId, { silent: true })
        for (const l of allLyrics) {
          if (l?.id) await adminDeleteLyrics(l.id)
        }
      } catch {
        // no lyrics to delete — fine
      }
      return
    }

    const language = payload.lyrics_language || 'en'
    const type = payload.lyrics_type || 'plain'

    // Try to find existing lyrics with the SAME language first (update path)
    try {
      const existing = await getTrackLyrics(trackId, language, { silent: true })
      if (existing?.id) {
        await adminUpdateLyrics(existing.id, {
          track_id: trackId,
          content,
          language,
          type,
        })
        return
      }
    } catch {
      // no existing lyrics with this language — will delete old ones and create
    }

    // Language changed (or first-time create): delete any old lyrics first,
    // then create fresh ones to avoid duplicates
    try {
      const allLyrics = await getLyricsByTrackID(trackId, { silent: true })
      for (const l of allLyrics) {
        if (l?.id) await adminDeleteLyrics(l.id)
      }
    } catch {
      // no old lyrics to delete
    }

    await adminCreateLyrics({
      track_id: trackId,
      content,
      language,
      type,
    })
  }

  async function createTrack(payload: TrackFormPayload) {
    saving.value = true
    error.value = null

    try {
      let createPayload = toCreatePayload(payload)

      if (payload.audioFile) {
        const uploadedAudio = await adminUploadTrackAudio(payload.audioFile)
        createPayload = {
          ...createPayload,
          audioUrl: getUploadedAudioUrl(uploadedAudio),
        }
      }

      if (payload.coverFile) {
        const uploadedCover = await adminUploadTrackCover(payload.coverFile)
        createPayload = {
          ...createPayload,
          coverUrl: getUploadedCoverUrl(uploadedCover),
        }
      }

      const created = await adminCreateTrack(createPayload)
      await syncTrackLyrics(created.id, payload)

      tracks.value = [created, ...tracks.value]
      return created
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function updateTrack(id: string | number, payload: TrackFormPayload) {
    saving.value = true
    error.value = null

    try {
      let updatePayload = toUpdatePayload(payload)

      if (payload.audioFile) {
        const uploadedAudio = await adminUploadTrackAudio(payload.audioFile)
        updatePayload = {
          ...updatePayload,
          audioUrl: getUploadedAudioUrl(uploadedAudio),
        }
      }

      if (payload.coverFile) {
        const uploadedCover = await adminUploadTrackCover(payload.coverFile)
        updatePayload = {
          ...updatePayload,
          coverUrl: getUploadedCoverUrl(uploadedCover),
        }
      }

      const updated = await adminUpdateTrack(id, updatePayload)
      await syncTrackLyrics(id, payload)

      tracks.value = tracks.value.map((track) =>
        String(track.id) === String(id) ? updated : track,
      )

      return updated
    } catch (err) {
      error.value = err
      throw err
    } finally {
      saving.value = false
    }
  }

  async function deleteTrack(id: string | number) {
    deleting.value = true
    error.value = null

    try {
      await adminDeleteTrack(id)
      tracks.value = tracks.value.filter((track) => String(track.id) !== String(id))
    } catch (err) {
      error.value = err
      throw err
    } finally {
      deleting.value = false
    }
  }

  function resetTracks() {
    tracks.value = []
    error.value = null
  }

  return {
    tracks,
    loading,
    saving,
    deleting,
    error,
    fetchTracks,
    createTrack,
    updateTrack,
    deleteTrack,
    resetTracks,
  }
}
