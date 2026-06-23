import z from 'zod'

export const UserPropsSchema = z.object({
  id: z.number().int(),
  first_name: z.string().optional(),
  last_name: z.string().optional(),
  full_name: z.string().optional(),
  username: z.string(),
  created_at: z.date().optional(),
  updated_at: z.date().optional(),
})

export type UserProps = z.infer<typeof UserPropsSchema>
