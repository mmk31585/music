import { z } from 'zod'

export const MoodTagSchema = z.object({
  track_id: z.string(),
  name: z.string(),
  confidence: z.number(),
})

export const MoodResponseSchema = z.object({
  track_id: z.string(),
  mood_tags: z.array(MoodTagSchema).optional(),
  energy: z.number(),
  valence: z.number(),
  tempo: z.number(),
  danceability: z.number(),
  acousticness: z.number(),
  instrumentalness: z.number(),
  liveness: z.number(),
  speechiness: z.number(),
})

export type MoodResponse = z.infer<typeof MoodResponseSchema>

export const AITrackItemSchema = z.object({
  id: z.string(),
  title: z.string(),
  artist: z.string().optional(),
  cover_url: z.string().optional(),
  duration: z.number().optional(),
})

export const AIPlaylistResponseSchema = z.object({
  id: z.string(),
  name: z.string(),
  description: z.string(),
  tracks: z.array(AITrackItemSchema),
  generated_at: z.string(),
})

export type AIPlaylistResponse = z.infer<typeof AIPlaylistResponseSchema>

export const TrackMetaSchema = z.object({
  id: z.string(),
  title: z.string(),
  artist: z.string().optional(),
  album: z.string().optional(),
  genre: z.string().optional(),
  duration: z.number().optional(),
})

export type TrackMeta = z.infer<typeof TrackMetaSchema>

export interface GeneratePlaylistPayload {
  prompt: string
  mood?: string
  activity?: string
  seed_track_id?: string
  genre?: string
  limit?: number
  exclude_ids?: string[]
}

export interface AnalyzeMoodPayload {
  track_id: string
}

export interface EmbeddingPayload {
  track_ids: string[]
}

export const MOOD_OPTIONS = [
  { label: 'Energetic', value: 'energetic', icon: 'pi pi-bolt' },
  { label: 'Happy', value: 'happy', icon: 'pi pi-sun' },
  { label: 'Chill', value: 'chill', icon: 'pi pi-cloud' },
  { label: 'Calm', value: 'calm', icon: 'pi pi-moon' },
  { label: 'Sad', value: 'sad', icon: 'pi pi-flag' },
  { label: 'Focus', value: 'focus', icon: 'pi pi-eye' },
  { label: 'Romantic', value: 'romantic', icon: 'pi pi-heart' },
  { label: 'Intense', value: 'intense', icon: 'pi pi-fire' },
  { label: 'Confident', value: 'confident', icon: 'pi pi-star' },
  { label: 'Sleep', value: 'sleep', icon: 'pi pi-moon' },
] as const

export const ACTIVITY_OPTIONS = [
  { label: 'Workout', value: 'workout', icon: 'pi pi-bolt' },
  { label: 'Studying', value: 'studying', icon: 'pi pi-book' },
  { label: 'Party', value: 'party', icon: 'pi pi-users' },
  { label: 'Commuting', value: 'commuting', icon: 'pi pi-car' },
  { label: 'Cooking', value: 'cooking', icon: 'pi pi-sparkles' },
  { label: 'Gaming', value: 'gaming', icon: 'pi pi-play' },
  { label: 'Relaxing', value: 'relaxing', icon: 'pi pi-inbox' },
  { label: 'Running', value: 'running', icon: 'pi pi-arrow-right' },
] as const
