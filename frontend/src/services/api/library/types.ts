import { z } from 'zod'

const NullableString = z.string().nullable().optional()
const NullableNumber = z.number().nullable().optional()

export const LibraryTrackSchema = z.object({
  track_id: z.string(),
  title: z.string(),
  artist_id: NullableString,
  artist_name: NullableString,
  album_id: NullableString,
  album_title: NullableString,
  cover_url: NullableString,
  audio_url: NullableString,
  duration_seconds: NullableNumber,
  added_at: z.string(),
  last_played_at: z.string().nullable().optional(),
})

export const LibraryAlbumSchema = z.object({
  album_id: z.string(),
  title: z.string(),
  artist_id: NullableString,
  artist_name: NullableString,
  cover_url: NullableString,
  release_date: z.string().nullable().optional(),
  added_at: z.string(),
})

export const LibraryArtistSchema = z.object({
  artist_id: z.string(),
  name: z.string(),
  cover_url: NullableString,
  followed_at: z.string(),
})

export const PlayHistoryItemSchema = z.object({
  track_id: z.string(),
  title: z.string(),
  artist_id: NullableString,
  artist_name: NullableString,
  album_id: NullableString,
  album_title: NullableString,
  cover_url: NullableString,
  audio_url: NullableString,
  duration_seconds: NullableNumber,
  added_at: z.string(),
  last_played_at: z.string().nullable().optional(),
})

export const LibraryTrackLikePayloadSchema = z.object({
  track_id: z.string(),
})

export const LibraryAlbumLikePayloadSchema = z.object({
  album_id: z.string(),
})

export const LibraryArtistFollowPayloadSchema = z.object({
  artist_id: z.string(),
})

export const AddPlayHistoryPayloadSchema = z.object({
  track_id: z.string(),
  duration: z.number().min(0).optional(),
  completed: z.boolean().optional(),
})

export type LibraryTrack = z.infer<typeof LibraryTrackSchema>
export type LibraryAlbum = z.infer<typeof LibraryAlbumSchema>
export type LibraryArtist = z.infer<typeof LibraryArtistSchema>
export type PlayHistoryItem = z.infer<typeof PlayHistoryItemSchema>

export type LibraryTrackLikePayload = z.infer<typeof LibraryTrackLikePayloadSchema>
export type LibraryAlbumLikePayload = z.infer<typeof LibraryAlbumLikePayloadSchema>
export type LibraryArtistFollowPayload = z.infer<typeof LibraryArtistFollowPayloadSchema>
export type AddPlayHistoryPayload = z.infer<typeof AddPlayHistoryPayloadSchema>
