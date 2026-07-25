import { z } from 'zod'

export const ReputationSummarySchema = z.object({
  userId: z.string(),
  trustScore: z.number(),
  tier: z.string(),
  tierLabel: z.string(),
  acceptedContributions: z.number(),
  totalContributions: z.number(),
  uploadSlots: z.number(),
  autoPublish: z.boolean(),
  canReview: z.boolean(),
})

export type ReputationSummary = z.infer<typeof ReputationSummarySchema>

export const TrustTierSchema = z.object({
  id: z.number(),
  slug: z.string(),
  label: z.string(),
  description: z.string().optional(),
  minScore: z.number(),
  minAccepted: z.number(),
  uploadSlots: z.number(),
  autoPublish: z.boolean(),
  canReview: z.boolean(),
  hierarchyLevel: z.number(),
})

export type TrustTier = z.infer<typeof TrustTierSchema>
