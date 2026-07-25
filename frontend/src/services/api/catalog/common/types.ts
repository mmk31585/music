import { z } from 'zod'

// Shared ID schema used across all catalog entities
export const IdSchema = z.union([z.string(), z.number()])

export type Id = z.infer<typeof IdSchema>

// Common pagination types
export interface PaginationParams {
  page?: number
  limit?: number
  sort?: string
  order?: 'asc' | 'desc'
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  totalPages: number
}
