import { z } from 'zod'

export const UserSchema = z.object({
  id: z.union([z.string(), z.number()]),
  email: z.string().email(),
  username: z.string().optional().nullable(),
  displayName: z.string().optional().nullable(),
  name: z.string().optional().nullable(),
  role: z.string().optional().nullable(),
})

export const LoginPayloadSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
})

export const RegisterPayloadSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
  username: z.string().min(2).optional(),
  displayName: z.string().min(2).optional(),
  name: z.string().min(2).optional(),
})

export const LogoutPayloadSchema = z.object({
  refreshToken: z.string().optional(),
})

const AuthDataSchema = z.object({
  user: UserSchema,
  access_token: z.string(),
  refresh_token: z.string().optional(),
  token: z.string().optional(),
  token_type: z.string().optional(),
  expires_in: z.number().optional(),
})

/**
 * Supports both shapes:
 *
 * 1. Full backend response:
 * {
 *   success,
 *   message,
 *   data: { user, access_token, refresh_token }
 * }
 *
 * 2. Unwrapped request-factory response:
 * {
 *   user,
 *   access_token,
 *   refresh_token
 * }
 */
export const AuthResponseSchema = z
  .union([
    z.object({
      success: z.boolean().optional(),
      message: z.string().optional(),
      data: AuthDataSchema,
    }),
    AuthDataSchema,
  ])
  .transform((value) => {
    if ('data' in value) {
      return value.data
    }

    return value
  })

/**
 * Supports both shapes:
 *
 * 1. Full backend response:
 * {
 *   success,
 *   message,
 *   data: { user }
 * }
 *
 * 2. Unwrapped request-factory response:
 * {
 *   user
 * }
 *
 * 3. Direct user:
 * {
 *   id,
 *   email,
 *   ...
 * }
 */
export const CurrentUserResponseSchema = z
  .union([
    z.object({
      success: z.boolean().optional(),
      message: z.string().optional(),
      data: z.object({
        user: UserSchema,
      }),
    }),
    z.object({
      user: UserSchema,
    }),
    UserSchema,
  ])
  .transform((value) => {
    if ('data' in value) {
      return value.data.user
    }

    if ('user' in value) {
      return value.user
    }

    return value
  })

export type User = z.infer<typeof UserSchema>
export type LoginPayload = z.infer<typeof LoginPayloadSchema>
export type RegisterPayload = z.infer<typeof RegisterPayloadSchema>
export type LogoutPayload = z.infer<typeof LogoutPayloadSchema>

export type AuthResponse = z.infer<typeof AuthResponseSchema>
export type CurrentUserResponse = z.infer<typeof CurrentUserResponseSchema>
