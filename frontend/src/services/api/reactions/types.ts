import { z } from 'zod'

export const ReactionSchema = z.object({
  id: z.string(),
  target_id: z.string(),
  target_type: z.string(),
  type: z.string(),
  user_id: z.string(),
  created_at: z.string(),
})

export const CountsResponseSchema = z.object({
  like: z.number(),
  love: z.number(),
  dislike: z.number(),
  total: z.number(),
})

export const ReactRequestSchema = z.object({
  target_id: z.string(),
  target_type: z.enum(['track', 'album', 'playlist', 'artist', 'comment', 'club_discussion', 'club_discussion_reply']),
  type: z.enum(['like', 'love', 'dislike']),
})

export type Reaction = z.infer<typeof ReactionSchema>
export type CountsResponse = z.infer<typeof CountsResponseSchema>
export type ReactRequest = z.infer<typeof ReactRequestSchema>
