package permissions

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	apperrors "music/internal/common/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInterface interface {
	// Permissions
	ListPermissions(ctx context.Context) ([]Permission, error)
	GetPermissionBySlug(ctx context.Context, slug string) (Permission, error)

	// Roles
	ListRoles(ctx context.Context) ([]Role, error)
	GetRoleBySlug(ctx context.Context, slug string) (Role, error)
	GetRoleByID(ctx context.Context, id int64) (Role, error)
	CreateRole(ctx context.Context, role Role) (Role, error)
	UpdateRole(ctx context.Context, id int64, updates map[string]any) error
	DeleteRole(ctx context.Context, id int64) error

	// Role ↔ Permission
	GetRolePermissions(ctx context.Context, roleID int64) ([]string, error)
	SetRolePermissions(ctx context.Context, roleID int64, permissionSlugs []string) error

	// User ↔ Role
	GetUserRoles(ctx context.Context, userID string) ([]UserRole, error)
	AssignUserRole(ctx context.Context, userID string, roleSlug string, assignedBy *string) error
	RemoveUserRole(ctx context.Context, userID string, roleSlug string) error

	// User Permission Overrides
	GetUserOverrides(ctx context.Context, userID string) ([]UserPermissionOverride, error)
	UpsertUserOverride(ctx context.Context, userID string, permSlug string, granted bool, reason string, grantedBy *string) error
	RemoveUserOverride(ctx context.Context, userID string, permSlug string) error

	// Resolved access
	ResolveUserAccess(ctx context.Context, userID string) (UserAccess, error)
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListPermissions(ctx context.Context) ([]Permission, error) {
	query := `
		SELECT id, slug, category, label, description, is_active, created_at, updated_at
		FROM permissions
		WHERE is_active = TRUE
		ORDER BY category, slug
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Slug, &p.Category, &p.Label, &p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, nil
}

func (r *Repository) GetPermissionBySlug(ctx context.Context, slug string) (Permission, error) {
	query := `
		SELECT id, slug, category, label, description, is_active, created_at, updated_at
		FROM permissions
		WHERE slug = $1
	`

	var p Permission
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.Slug, &p.Category, &p.Label, &p.Description, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Permission{}, apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "permission not found", nil)
		}
		return Permission{}, err
	}
	return p, nil
}

func (r *Repository) ListRoles(ctx context.Context) ([]Role, error) {
	query := `
		SELECT id, slug, label, description, hierarchy_level, is_system_role, is_active, created_at, updated_at
		FROM roles
		WHERE is_active = TRUE
		ORDER BY hierarchy_level DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Slug, &role.Label, &role.Description, &role.HierarchyLevel, &role.IsSystemRole, &role.IsActive, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *Repository) GetRoleBySlug(ctx context.Context, slug string) (Role, error) {
	query := `
		SELECT id, slug, label, description, hierarchy_level, is_system_role, is_active, created_at, updated_at
		FROM roles
		WHERE slug = $1
	`

	var role Role
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&role.ID, &role.Slug, &role.Label, &role.Description, &role.HierarchyLevel, &role.IsSystemRole, &role.IsActive, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Role{}, apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "role not found", nil)
		}
		return Role{}, err
	}
	return role, nil
}

func (r *Repository) GetRoleByID(ctx context.Context, id int64) (Role, error) {
	query := `
		SELECT id, slug, label, description, hierarchy_level, is_system_role, is_active, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	var role Role
	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID, &role.Slug, &role.Label, &role.Description, &role.HierarchyLevel, &role.IsSystemRole, &role.IsActive, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Role{}, apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "role not found", nil)
		}
		return Role{}, err
	}
	return role, nil
}

func (r *Repository) CreateRole(ctx context.Context, role Role) (Role, error) {
	query := `
		INSERT INTO roles (slug, label, description, hierarchy_level, is_system_role, is_active)
		VALUES ($1, $2, $3, $4, FALSE, TRUE)
		RETURNING id, slug, label, description, hierarchy_level, is_system_role, is_active, created_at, updated_at
	`

	var created Role
	err := r.db.QueryRow(ctx, query,
		role.Slug, role.Label, role.Description, role.HierarchyLevel,
	).Scan(
		&created.ID, &created.Slug, &created.Label, &created.Description,
		&created.HierarchyLevel, &created.IsSystemRole, &created.IsActive,
		&created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return Role{}, err
	}
	return created, nil
}

func (r *Repository) UpdateRole(ctx context.Context, id int64, updates map[string]any) error {
	allowed := map[string]string{
		"label":       "label",
		"description": "description",
		"hierarchy_level": "hierarchy_level",
		"is_active":   "is_active",
	}

	setClauses := make([]string, 0, len(updates))
	args := make([]any, 0, len(updates)+1)
	argIdx := 1

	for key, value := range updates {
		col, ok := allowed[key]
		if !ok {
			return apperrors.New(http.StatusBadRequest, apperrors.CodeBadRequest, "unknown role field: "+key, nil)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, argIdx))
		args = append(args, value)
		argIdx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE roles SET %s WHERE id = $%d AND is_system_role = FALSE", strings.Join(setClauses, ", "), argIdx)
	args = append(args, id)

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "role not found or is a system role", nil)
	}
	return nil
}

func (r *Repository) DeleteRole(ctx context.Context, id int64) error {
	query := `DELETE FROM roles WHERE id = $1 AND is_system_role = FALSE`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(http.StatusNotFound, apperrors.CodeNotFound, "role not found or is a system role", nil)
	}
	return nil
}

func (r *Repository) GetRolePermissions(ctx context.Context, roleID int64) ([]string, error) {
	query := `
		SELECT p.slug
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.category, p.slug
	`

	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		slugs = append(slugs, slug)
	}
	return slugs, nil
}

func (r *Repository) SetRolePermissions(ctx context.Context, roleID int64, permissionSlugs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Clear existing
	if _, err := tx.Exec(ctx, "DELETE FROM role_permissions WHERE role_id = $1", roleID); err != nil {
		return err
	}

	// Insert new
	if len(permissionSlugs) > 0 {
		for i, slug := range permissionSlugs {
			query := fmt.Sprintf(`
				INSERT INTO role_permissions (role_id, permission_id)
				SELECT $1, id FROM permissions WHERE slug = $%d
			`, i+2)
			args := append([]any{roleID}, any(slug))
			if _, err := tx.Exec(ctx, query, args...); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetUserRoles(ctx context.Context, userID string) ([]UserRole, error) {
	query := `
		SELECT ur.id, ur.user_id, ur.role_id, r.slug AS role_slug,
		       ur.assigned_by, ur.assigned_at, ur.expires_at, ur.is_active
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND ur.is_active = TRUE
		  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		ORDER BY r.hierarchy_level DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRoles []UserRole
	for rows.Next() {
		var ur UserRole
		if err := rows.Scan(&ur.ID, &ur.UserID, &ur.RoleID, &ur.RoleSlug, &ur.AssignedBy, &ur.AssignedAt, &ur.ExpiresAt, &ur.IsActive); err != nil {
			return nil, err
		}
		userRoles = append(userRoles, ur)
	}
	return userRoles, nil
}

func (r *Repository) AssignUserRole(ctx context.Context, userID string, roleSlug string, assignedBy *string) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_by)
		SELECT $1, id, $3 FROM roles WHERE slug = $2
		ON CONFLICT (user_id, role_id) DO UPDATE SET is_active = TRUE, assigned_at = NOW()
	`
	_, err := r.db.Exec(ctx, query, userID, roleSlug, assignedBy)
	return err
}

func (r *Repository) RemoveUserRole(ctx context.Context, userID string, roleSlug string) error {
	query := `
		UPDATE user_roles SET is_active = FALSE
		WHERE user_id = $1 AND role_id = (SELECT id FROM roles WHERE slug = $2)
	`
	_, err := r.db.Exec(ctx, query, userID, roleSlug)
	return err
}

func (r *Repository) GetUserOverrides(ctx context.Context, userID string) ([]UserPermissionOverride, error) {
	query := `
		SELECT upo.id, upo.user_id, upo.permission_id, p.slug AS perm_slug,
		       upo.granted, upo.reason, upo.granted_by, upo.granted_at, upo.expires_at
		FROM user_permission_overrides upo
		JOIN permissions p ON p.id = upo.permission_id
		WHERE upo.user_id = $1
		  AND (upo.expires_at IS NULL OR upo.expires_at > NOW())
		ORDER BY p.category, p.slug
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	overrides := make([]UserPermissionOverride, 0)
	for rows.Next() {
		var o UserPermissionOverride
		if err := rows.Scan(&o.ID, &o.UserID, &o.PermissionID, &o.PermSlug, &o.Granted, &o.Reason, &o.GrantedBy, &o.GrantedAt, &o.ExpiresAt); err != nil {
			return nil, err
		}
		overrides = append(overrides, o)
	}
	return overrides, nil
}

func (r *Repository) UpsertUserOverride(ctx context.Context, userID string, permSlug string, granted bool, reason string, grantedBy *string) error {
	query := `
		INSERT INTO user_permission_overrides (user_id, permission_id, granted, reason, granted_by)
		SELECT $1, id, $3, $4, $5 FROM permissions WHERE slug = $2
		ON CONFLICT (user_id, permission_id) DO UPDATE
		SET granted = $3, reason = $4, granted_by = $5, granted_at = NOW()
	`
	_, err := r.db.Exec(ctx, query, userID, permSlug, granted, reason, grantedBy)
	return err
}

func (r *Repository) RemoveUserOverride(ctx context.Context, userID string, permSlug string) error {
	query := `
		DELETE FROM user_permission_overrides
		WHERE user_id = $1 AND permission_id = (SELECT id FROM permissions WHERE slug = $2)
	`
	_, err := r.db.Exec(ctx, query, userID, permSlug)
	return err
}

func (r *Repository) ResolveUserAccess(ctx context.Context, userID string) (UserAccess, error) {
	// Get highest active role
	roleQuery := `
		SELECT r.slug, r.hierarchy_level
		FROM user_roles ur
		JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND ur.is_active = TRUE
		  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		ORDER BY r.hierarchy_level DESC
		LIMIT 1
	`

	var access UserAccess
	access.UserID = userID

	err := r.db.QueryRow(ctx, roleQuery, userID).Scan(&access.RoleSlug, &access.RoleLevel)
	if err != nil {
		if err == pgx.ErrNoRows {
			// No role assigned; default to listener
			access.RoleSlug = "listener"
			access.RoleLevel = 10
		} else {
			return UserAccess{}, err
		}
	}

	// Get role permissions
	permQuery := `
		SELECT DISTINCT p.slug
		FROM user_roles ur
		JOIN role_permissions rp ON rp.role_id = ur.role_id
		JOIN permissions p ON p.id = rp.permission_id
		WHERE ur.user_id = $1 AND ur.is_active = TRUE
		  AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
		  AND p.is_active = TRUE
	`

	rows, err := r.db.Query(ctx, permQuery, userID)
	if err != nil {
		return UserAccess{}, err
	}
	defer rows.Close()

	permSet := make(map[string]struct{})
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return UserAccess{}, err
		}
		permSet[slug] = struct{}{}
	}

	// Apply overrides
	overrideQuery := `
		SELECT p.slug, upo.granted
		FROM user_permission_overrides upo
		JOIN permissions p ON p.id = upo.permission_id
		WHERE upo.user_id = $1
		  AND (upo.expires_at IS NULL OR upo.expires_at > NOW())
	`

	orows, err := r.db.Query(ctx, overrideQuery, userID)
	if err != nil {
		return UserAccess{}, err
	}
	defer orows.Close()

	for orows.Next() {
		var slug string
		var granted bool
		if err := orows.Scan(&slug, &granted); err != nil {
			return UserAccess{}, err
		}
		if granted {
			permSet[slug] = struct{}{}
		} else {
			delete(permSet, slug)
		}
	}

	access.Permissions = make([]string, 0, len(permSet))
	for slug := range permSet {
		access.Permissions = append(access.Permissions, slug)
	}

	return access, nil
}
