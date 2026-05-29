import { z } from 'zod'

export const PlaybackTrackSchema = z.object({
  id: z.string(),
  title: z.string(),
  artistName: z.string().default('Unknown artist'),
  albumTitle: z.string().nullable().optional(),
  coverUrl: z.string().nullable().optional(),
  durationSeconds: z.number().nullable().optional(),
  streamUrl: z.string(),
})

export type PlaybackTrack = z.infer<typeof PlaybackTrackSchema>
