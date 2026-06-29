import { z } from 'zod'

const optionalNumber = z.preprocess((value) => {
  if (value === null || value === undefined || value === '') return undefined

  const numberValue = Number(value)

  return Number.isFinite(numberValue) ? numberValue : undefined
}, z.number().optional())

const optionalBoolean = z.preprocess((value) => {
  if (value === null || value === undefined || value === '') return undefined
  if (typeof value === 'boolean') return value
  if (value === 'true') return true
  if (value === 'false') return false

  return value
}, z.boolean().optional())

export const UploadResponseSchema = z.object({
  mediaId: z.union([z.string(), z.number()]).optional(),
  url: z.string().optional(),
  file_url: z.string().optional(),
  fileUrl: z.string().optional(),
  path: z.string().optional(),

  filename: z.string().optional(),
  fileName: z.string().optional(),
  original_name: z.string().optional(),
  originalName: z.string().optional(),

  mime_type: z.string().optional(),
  mimeType: z.string().optional(),

  size: optionalNumber,

  durationSeconds: optionalNumber,
  duration_seconds: optionalNumber,

  duplicate: optionalBoolean,
})

export type UploadResponse = z.infer<typeof UploadResponseSchema>

export interface Media {
  id: string | number
  objectKey: string
  originalFilename?: string | null
  fileSize?: number | null
  mimeType?: string | null
  mediaType?: string | null
  publicUrl?: string | null
  createdAt?: string | null
}
