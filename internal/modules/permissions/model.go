package permissions

import "time"

type Permission struct {
	ID          int64     `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	Category    string    `json:"category" db:"category"`
	Label       string    `json:"label" db:"label"`
	Description string    `json:"description,omitempty" db:"description"`
	IsActive    bool      `json:"isActive" db:"is_active"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type Role struct {
	ID              int64     `json:"id" db:"id"`
	Slug            string    `json:"slug" db:"slug"`
	Label           string    `json:"label" db:"label"`
	Description     string    `json:"description,omitempty" db:"description"`
	HierarchyLevel  int       `json:"hierarchyLevel" db:"hierarchy_level"`
	IsSystemRole    bool      `json:"isSystemRole" db:"is_system_role"`
	IsActive        bool      `json:"isActive" db:"is_active"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type UserRole struct {
	ID         int64      `json:"id" db:"id"`
	UserID     string     `json:"userId" db:"user_id"`
	RoleID     int64      `json:"roleId" db:"role_id"`
	RoleSlug   string     `json:"roleSlug" db:"role_slug"`
	AssignedBy *string    `json:"assignedBy,omitempty" db:"assigned_by"`
	AssignedAt time.Time  `json:"assignedAt" db:"assigned_at"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
	IsActive   bool       `json:"isActive" db:"is_active"`
}

type UserPermissionOverride struct {
	ID           int64      `json:"id" db:"id"`
	UserID       string     `json:"userId" db:"user_id"`
	PermissionID int64      `json:"permissionId" db:"permission_id"`
	PermSlug     string     `json:"permSlug" db:"perm_slug"`
	Granted      bool       `json:"granted" db:"granted"`
	Reason       string     `json:"reason,omitempty" db:"reason"`
	GrantedBy    *string    `json:"grantedBy,omitempty" db:"granted_by"`
	GrantedAt    time.Time  `json:"grantedAt" db:"granted_at"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty" db:"expires_at"`
}

// RoleWithPermissions is a role enriched with its permission slugs.
type RoleWithPermissions struct {
	Role
	Permissions []string `json:"permissions"`
}

// UserAccess is the full resolved access set for a user.
type UserAccess struct {
	UserID      string   `json:"userId"`
	RoleSlug    string   `json:"roleSlug"`
	RoleLevel   int      `json:"roleLevel"`
	Permissions []string `json:"permissions"`
}

// PermissionCacheEntry is a cached user access record.
type PermissionCacheEntry struct {
	Access    UserAccess
	ExpiresAt time.Time
}
