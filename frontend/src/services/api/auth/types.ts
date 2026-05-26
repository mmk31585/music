import { z } from 'zod'

export const UserSchema = z.object({
  id: z.union([z.string(), z.number()]),
  email: z.email(),
  name: z.string().nullable().optional(),
  username: z.string().nullable().optional(),
  role: z.string().nullable().optional(),
})

export type User = z.infer<typeof UserSchema>

export const LoginPayloadSchema = z.object({
  email: z.email(),
  password: z.string().min(1),
})

export type LoginPayload = z.infer<typeof LoginPayloadSchema>

export const RegisterPayloadSchema = z.object({
  name: z.string().min(2).optional(),
  email: z.email(),
  password: z.string().min(6),
})

export type RegisterPayload = z.infer<typeof RegisterPayloadSchema>

/**
 * This is intentionally flexible because Go backends often return:
 * {
 *   user,
 *   access_token,
 *   refresh_token
 * }
 *
 * or:
 * {
 *   data: { user, access_token, refresh_token }
 * }
 *
 * If your backend response is different, only update this schema.
 */
export const AuthResponseSchema = z.object({
  user: UserSchema,
  access_token: z.string(),
  refresh_token: z.string(),
  token: z.string().optional(),

  data: z
    .object({
      user: UserSchema.optional(),
      access_token: z.string().optional(),
      refresh_token: z.string().optional(),
      token: z.string().optional(),
    })
    .optional(),

  message: z.string().optional(),
  success: z.boolean().optional(),
})

export type AuthResponse = z.infer<typeof AuthResponseSchema>
