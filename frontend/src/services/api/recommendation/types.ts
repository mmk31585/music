import { z } from 'zod'

export const RecommendationTrackSchema = z.object({
  id: z.string(),
  title: z.string(),
  artist_id: z.string().nullable().optional(),
  artist_name: z.string().nullable().optional(),
  album_id: z.string().nullable().optional(),
  album_title: z.string().nullable().optional(),
  genre: z.string().nullable().optional(),
  cover_url: z.string().nullable().optional(),
  audio_url: z.string().nullable().optional(),
  duration_seconds: z.number().nullable().optional(),
  score: z.number().nullable().optional(),
})

export type RecommendationTrack = z.infer<typeof RecommendationTrackSchema>

export const RecommendationResponseSchema = z.object({
  type: z.string(),
  items: z.array(RecommendationTrackSchema).catch([]),
  limit: z.number(),
})

export type RecommendationResponse = z.infer<typeof RecommendationResponseSchema>
