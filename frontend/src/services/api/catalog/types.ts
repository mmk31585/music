import { z } from 'zod'

const IdSchema = z.union([z.string(), z.number()])

export const ArtistSchema = z
  .object({
    id: IdSchema,
    name: z.string(),
    bio: z.string().optional().nullable(),

    image_url: z.string().optional().nullable(),

    is_verified: z.boolean().optional().nullable(),

    monthly_listeners: z.number().optional().nullable(),
  })
  .transform((artist) => ({
    id: artist.id,
    name: artist.name,
    bio: artist.bio ?? null,
    image_url: artist.image_url ,
    is_verified: artist.is_verified ,
    monthly_listeners: artist.monthly_listeners ,
  }))

export const AlbumSchema = z
  .object({
    id: IdSchema,
    title: z.string(),
    cover_url: z.string().optional().nullable(),
    artist_id: IdSchema.optional().nullable(),
    artist_name: z.string().optional().nullable(),
  })
  .transform((album) => ({
    id: album.id,
    title: album.title,
    cover_url: album.cover_url ,
    artist_id: album.artist_id ,
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
    audio_url: z.string().optional().nullable(),
    cover_url: z.string().optional().nullable(),
    artist_id: IdSchema.optional().nullable(),
    album_id: IdSchema.optional().nullable(),
    genre_id: IdSchema.optional().nullable(),
    artist_name: z.string().optional().nullable(),
    album_title: z.string().optional().nullable(),
    genres: z.array(GenreSchema).optional().nullable(),
  })
  .transform((track) => ({
    id: track.id,
    title: track.title,
    duration_seconds: track.duration_seconds ,
    audio_url: track.audio_url ,
    cover_url: track.cover_url ,
    artist_id: track.artist_id ,
    album_id: track.album_id ,
    genre_id: track.genre_id ?? track.genres?.[0]?.id ?? null,
    artist_name: track.artist_name ?? null,
    album_title: track.album_title ?? null,
    genres: track.genres ?? [],
  }))

export type Artist = z.infer<typeof ArtistSchema>
export type Album = z.infer<typeof AlbumSchema>
export type Genre = z.infer<typeof GenreSchema>
export type Track = z.infer<typeof TrackSchema>
