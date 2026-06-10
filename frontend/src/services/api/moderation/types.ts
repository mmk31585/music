import { z } from 'zod'

export const ContentReportSchema = z.object({
  id: z.string(),
  reporter_id: z.string(),
  target_id: z.string(),
  target_type: z.string(),
  reason: z.string(),
  description: z.string().nullable().optional(),
  status: z.string(),
  moderator_id: z.string().nullable().optional(),
  resolved_at: z.string().nullable().optional(),
  created_at: z.string(),
  resolution_note: z.string().nullable().optional(),
})

export const ContentFlagSchema = z.object({
  id: z.string(),
  target_id: z.string(),
  target_type: z.string(),
  flag_type: z.string(),
  flagged_at: z.string(),
  expires_at: z.string().nullable().optional(),
})

export const ModerationActionSchema = z.object({
  id: z.string(),
  report_id: z.string().nullable().optional(),
  moderator_id: z.string(),
  action: z.string(),
  target_id: z.string(),
  target_type: z.string(),
  previous_status: z.string().nullable().optional(),
  new_status: z.string().nullable().optional(),
  note: z.string().nullable().optional(),
  metadata: z.record(z.string(), z.any()).nullable().optional(),
  created_at: z.string(),
})

export const ModerationStatsSchema = z.object({
  total_reports: z.number(),
  pending_reports: z.number(),
  resolved_today: z.number(),
  flagged_content: z.number(),
  unique_reporters: z.number(),
  avg_resolution_hours: z.number(),
  by_reason: z.record(z.string(), z.number()).nullable().optional(),
  by_target_type: z.record(z.string(), z.number()).nullable().optional(),
})

export type ContentReport = z.infer<typeof ContentReportSchema>
export type ContentFlag = z.infer<typeof ContentFlagSchema>
export type ModerationAction = z.infer<typeof ModerationActionSchema>
export type ModerationStats = z.infer<typeof ModerationStatsSchema>
