import { z } from 'zod'

export const PermissionSchema = z.object({
  id: z.number(),
  slug: z.string(),
  category: z.string(),
  label: z.string(),
  description: z.string().optional(),
})

export type Permission = z.infer<typeof PermissionSchema>

export const RoleSchema = z.object({
  id: z.number(),
  slug: z.string(),
  label: z.string(),
  description: z.string().optional(),
  hierarchyLevel: z.number(),
  isSystemRole: z.boolean(),
  isActive: z.boolean(),
  createdAt: z.string(),
  updatedAt: z.string().optional().default(''),
})

export type Role = z.infer<typeof RoleSchema>

export const RoleWithPermissionsSchema = RoleSchema.extend({
  permissions: z.array(z.string()).optional().default([]),
})

export type RoleWithPermissions = z.infer<typeof RoleWithPermissionsSchema>

export const UserAccessSchema = z.object({
  userId: z.string(),
  roleSlug: z.string(),
  roleLevel: z.number(),
  permissions: z.array(z.string()),
})

export type UserAccess = z.infer<typeof UserAccessSchema>

export const CheckPermissionResponseSchema = z.object({
  hasPermission: z.boolean(),
  permSlug: z.string(),
})

export type CheckPermissionResponse = z.infer<typeof CheckPermissionResponseSchema>

export const UserPermissionOverrideSchema = z.object({
  id: z.number(),
  userId: z.string(),
  permissionId: z.number(),
  permSlug: z.string(),
  granted: z.boolean(),
  reason: z.string().optional(),
  grantedBy: z.string().optional(),
  grantedAt: z.string(),
  expiresAt: z.string().optional(),
})

export type UserPermissionOverride = z.infer<typeof UserPermissionOverrideSchema>
