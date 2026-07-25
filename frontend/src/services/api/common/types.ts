import { z } from 'zod'
import { UserSchema } from '@/services/api/auth/types'
import type { ApiResponseProps } from '@/plugins/client/types.ts'

export type ApiClientError = Omit<ApiResponseProps, 'data'> & {
  errors?: Record<string, string>
}

export const TimestampsSchema = z.object({
  created_at: z.string().optional().nullable(),
  updated_at: z.string().optional().nullable(),
  deleted_at: z.string().optional().nullable(),
})

export type TimestampsProps = z.infer<typeof TimestampsSchema>

export const UserstampsSchema = z.object({
  created_by: z
    .lazy(() => UserSchema)
    .optional()
    .nullable(),
  updated_by: z
    .lazy(() => UserSchema)
    .optional()
    .nullable(),
  deleted_by: z
    .lazy(() => UserSchema)
    .optional()
    .nullable(),
})

export type UserstampsProps = z.infer<typeof UserstampsSchema>

export const StampsSchema = TimestampsSchema.extend(UserstampsSchema.shape)

export type StampsProps = z.infer<typeof StampsSchema>
