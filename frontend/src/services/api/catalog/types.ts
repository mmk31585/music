import { z } from 'zod'

const IdSchema = z.union([z.string(), z.number()])

export const ArtistSchema = z
  .object({
    id: IdSchema,
    name: z.string(),
    image_url: z.string().optional().nullable(),
    imageUrl: z.string().optional().nullable(),
  })
  .transform((artist) => ({
    id: artist.id,
    name: artist.name,
    image_url: artist.image_url ?? artist.imageUrl ?? null,
  }))

export const AlbumSchema = z
  .object({
    id: IdSchema,
    title: z.string(),
    cover_url: z.string().optional().nullable(),
    coverUrl: z.string().optional().nullable(),
    artist_id: IdSchema.optional().nullable(),
    artistId: IdSchema.optional().nullable(),
    artist_name: z.string().optional().nullable(),
  })
  .transform((album) => ({
    id: album.id,
    title: album.title,
    cover_url: album.cover_url ?? album.coverUrl ?? null,
    artist_id: album.artist_id ?? album.artistId ?? null,
    artist_name: album.artist_name ?? null,
  }))

export const GenreSchema = z.object({
  id: IdSchema,
  name: z.string(),
})

export const TrackSchema = z
  .object({
    id: IdSchema,
    title: z.string(),
    duration_seconds: z.number().optional().nullable(),
    durationSeconds: z.number().optional().nullable(),
    audio_url: z.string().optional().nullable(),
    audioUrl: z.string().optional().nullable(),
    cover_url: z.string().optional().nullable(),
    coverUrl: z.string().optional().nullable(),
    artist_id: IdSchema.optional().nullable(),
    artistId: IdSchema.optional().nullable(),
    album_id: IdSchema.optional().nullable(),
    albumId: IdSchema.optional().nullable(),
    genre_id: IdSchema.optional().nullable(),
    artist_name: z.string().optional().nullable(),
    album_title: z.string().optional().nullable(),
    genres: z.array(GenreSchema).optional().nullable(),
  })
  .transform((track) => ({
    id: track.id,
    title: track.title,
    duration_seconds: track.duration_seconds ?? track.durationSeconds ?? null,
    audio_url: track.audio_url ?? track.audioUrl ?? null,
    cover_url: track.cover_url ?? track.coverUrl ?? null,
    artist_id: track.artist_id ?? track.artistId ?? null,
    album_id: track.album_id ?? track.albumId ?? null,
    genre_id: track.genre_id ?? track.genres?.[0]?.id ?? null,
    artist_name: track.artist_name ?? null,
    album_title: track.album_title ?? null,
    genres: track.genres ?? [],
  }))

export type Artist = z.infer<typeof ArtistSchema>
export type Album = z.infer<typeof AlbumSchema>
export type Genre = z.infer<typeof GenreSchema>
export type Track = z.infer<typeof TrackSchema>
