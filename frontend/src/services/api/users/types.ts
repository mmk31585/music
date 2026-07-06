import z from 'zod'

export const UserPropsSchema = z.object({
  id: z.string(),
  full_name: z.string(),
  username: z.string(),
  avatar_url: z.string().nullable().optional(),
  created_at: z.string(),
  updated_at: z.string().nullable(),
})

export type UserProps = z.infer<typeof UserPropsSchema>
