import { z } from 'zod'
import { ArtistSchema } from '../artists'
import { AlbumSchema } from '../albums'
import { TrackSchema } from '../tracks'

export const SearchResultSchema = z.object({
  artists: z.array(ArtistSchema).optional().default([]),
  albums: z.array(AlbumSchema).optional().default([]),
  tracks: z.array(TrackSchema).optional().default([]),
})

export type SearchResult = z.infer<typeof SearchResultSchema>

export interface SearchParams {
  query: string
  type?: 'all' | 'artists' | 'albums' | 'tracks'
  limit?: number
}
