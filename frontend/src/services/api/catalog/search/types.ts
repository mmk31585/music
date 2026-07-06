import { z } from 'zod'
import { IdSchema } from '../common'

/**
 * Search-specific schemas match the backend search module's camelCase DTOs.
 * The catalog API uses snake_case, but the search module uses camelCase.
 * These are kept separate to avoid confusion.
 */

const SearchTrackArtistSchema = z.object({
  artistId: z.string().optional(),
  name: z.string(),
  slug: z.string().optional(),
  role: z.string().optional(),
  position: z.number().optional(),
})

const SearchGenreSchema = z.object({
  id: z.string(),
  name: z.string(),
  slug: z.string().optional(),
})

const SearchTrackSchema = z.object({
  id: z.string(),
  title: z.string(),
  slug: z.string().optional(),
  artistId: z.string().nullable().optional(),
  albumId: z.string().nullable().optional(),
  coverUrl: z.string().nullable().optional(),
  audioUrl: z.string().nullable().optional(),
  durationSeconds: z.number().default(0),
  trackNumber: z.number().nullable().optional(),
  explicit: z.boolean().default(false),
  playCount: z.number().default(0),
  isPublic: z.boolean().optional().default(true),
  createdAt: z.string().optional(),
  artists: z.array(SearchTrackArtistSchema).optional(),
  genres: z.array(SearchGenreSchema).optional(),
})

const SearchAlbumSchema = z.object({
  id: z.string(),
  title: z.string(),
  slug: z.string().optional(),
  artistId: z.string().nullable().optional(),
  coverUrl: z.string().nullable().optional(),
  releaseDate: z.string().nullable().optional(),
  artists: z.array(SearchTrackArtistSchema).optional(),
})

const SearchArtistSchema = z.object({
  id: z.string(),
  name: z.string(),
  slug: z.string().optional(),
  bio: z.string().nullable().optional(),
  imageUrl: z.string().nullable().optional(),
  isVerified: z.boolean().default(false),
  monthlyListeners: z.number().default(0),
})

const SearchPlaylistSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string().nullable().optional(),
  coverUrl: z.string().nullable().optional(),
  userId: z.string().nullable().optional(),
  isPublic: z.boolean().nullable().optional(),
})

export const SearchResultSchema = z.object({
  query: z.string().optional(),
  tracks: z.array(SearchTrackSchema).optional().default([]),
  albums: z.array(SearchAlbumSchema).optional().default([]),
  artists: z.array(SearchArtistSchema).optional().default([]),
  playlists: z.array(SearchPlaylistSchema).optional().default([]),
})

export type SearchResult = z.infer<typeof SearchResultSchema>

export interface SearchParams {
  query: string
  type?: 'all' | 'artists' | 'albums' | 'tracks'
  limit?: number
}
