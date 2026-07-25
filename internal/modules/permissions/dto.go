package permissions

// --- Role DTOs ---

type CreateRoleRequest struct {
	Slug           string `json:"slug" binding:"required,min=2,max=50"`
	Label          string `json:"label" binding:"required,min=2,max=100"`
	Description    string `json:"description"`
	HierarchyLevel int    `json:"hierarchyLevel" binding:"min=0,max=9999"`
}

type UpdateRoleRequest struct {
	Label          *string `json:"label"`
	Description    *string `json:"description"`
	HierarchyLevel *int    `json:"hierarchyLevel" binding:"omitempty,min=0,max=9999"`
	IsActive       *bool   `json:"isActive"`
}

type SetRolePermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required,min=1"`
}

// --- User Role DTOs ---

type AssignUserRoleRequest struct {
	RoleSlug string `json:"roleSlug" binding:"required"`
}

type UserRoleResponse struct {
	UserID     string `json:"userId"`
	RoleSlug   string `json:"roleSlug"`
	AssignedBy string `json:"assignedBy,omitempty"`
}

// --- User Permission Override DTOs ---

type GrantPermissionRequest struct {
	PermSlug string `json:"permSlug" binding:"required"`
	Reason   string `json:"reason" binding:"required,min=3"`
}

type RevokePermissionRequest struct {
	PermSlug string `json:"permSlug" binding:"required"`
}

// --- Access Resolution DTOs ---

type UserAccessResponse struct {
	UserID      string   `json:"userId"`
	RoleSlug    string   `json:"roleSlug"`
	RoleLevel   int      `json:"roleLevel"`
	Permissions []string `json:"permissions"`
}

type CheckPermissionRequest struct {
	PermSlug string `json:"permSlug" binding:"required"`
}

type CheckPermissionResponse struct {
	HasPermission bool   `json:"hasPermission"`
	PermSlug      string `json:"permSlug"`
}

// --- Permission DTOs ---

type PermissionResponse struct {
	ID          int64  `json:"id"`
	Slug        string `json:"slug"`
	Category    string `json:"category"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}
