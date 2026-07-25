package permissions

import (
	"context"
	"sync"
	"time"
)

const (
	cacheTTL      = 5 * time.Minute
	cleanupEvery  = 10 * time.Minute
)

type ServiceInterface interface {
	// Permissions
	ListPermissions(ctx context.Context) ([]Permission, error)

	// Roles
	ListRoles(ctx context.Context) ([]Role, error)
	GetRole(ctx context.Context, slug string) (Role, error)
	CreateRole(ctx context.Context, role Role) (Role, error)
	UpdateRole(ctx context.Context, id int64, updates map[string]any) error
	DeleteRole(ctx context.Context, id int64) error

	// Role permissions
	GetRolePermissions(ctx context.Context, roleSlug string) ([]string, error)
	SetRolePermissions(ctx context.Context, roleSlug string, permSlugs []string) error

	// User role management
	AssignUserRole(ctx context.Context, userID string, roleSlug string, assignedBy *string) error
	RemoveUserRole(ctx context.Context, userID string, roleSlug string) error

	// User permission overrides
	GetUserOverrides(ctx context.Context, userID string) ([]UserPermissionOverride, error)
	GrantUserPermission(ctx context.Context, userID string, permSlug string, reason string, grantedBy *string) error
	RevokeUserPermission(ctx context.Context, userID string, permSlug string, grantedBy *string) error

	// Access resolution
	ResolveUserAccess(ctx context.Context, userID string) (UserAccess, error)
	HasPermission(ctx context.Context, userID string, permSlug string) (bool, error)
	HasAnyPermission(ctx context.Context, userID string, permSlugs ...string) (bool, error)

	// Cache management
	InvalidateUser(userID string)
	InvalidateAll()
}

type Service struct {
	repo  RepositoryInterface
	cache sync.Map // map[string]*PermissionCacheEntry
}

func NewService(repo RepositoryInterface) *Service {
	s := &Service{repo: repo}
	go s.cleanupLoop()
	return s
}

func (s *Service) ListPermissions(ctx context.Context) ([]Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *Service) ListRoles(ctx context.Context) ([]Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *Service) GetRole(ctx context.Context, slug string) (Role, error) {
	return s.repo.GetRoleBySlug(ctx, slug)
}

func (s *Service) CreateRole(ctx context.Context, role Role) (Role, error) {
	return s.repo.CreateRole(ctx, role)
}

func (s *Service) UpdateRole(ctx context.Context, id int64, updates map[string]any) error {
	return s.repo.UpdateRole(ctx, id, updates)
}

func (s *Service) DeleteRole(ctx context.Context, id int64) error {
	return s.repo.DeleteRole(ctx, id)
}

func (s *Service) GetRolePermissions(ctx context.Context, roleSlug string) ([]string, error) {
	role, err := s.repo.GetRoleBySlug(ctx, roleSlug)
	if err != nil {
		return nil, err
	}
	return s.repo.GetRolePermissions(ctx, role.ID)
}

func (s *Service) SetRolePermissions(ctx context.Context, roleSlug string, permSlugs []string) error {
	role, err := s.repo.GetRoleBySlug(ctx, roleSlug)
	if err != nil {
		return err
	}
	return s.repo.SetRolePermissions(ctx, role.ID, permSlugs)
}

func (s *Service) AssignUserRole(ctx context.Context, userID string, roleSlug string, assignedBy *string) error {
	return s.repo.AssignUserRole(ctx, userID, roleSlug, assignedBy)
}

func (s *Service) RemoveUserRole(ctx context.Context, userID string, roleSlug string) error {
	return s.repo.RemoveUserRole(ctx, userID, roleSlug)
}

func (s *Service) GetUserOverrides(ctx context.Context, userID string) ([]UserPermissionOverride, error) {
	return s.repo.GetUserOverrides(ctx, userID)
}

func (s *Service) GrantUserPermission(ctx context.Context, userID string, permSlug string, reason string, grantedBy *string) error {
	return s.repo.UpsertUserOverride(ctx, userID, permSlug, true, reason, grantedBy)
}

func (s *Service) RevokeUserPermission(ctx context.Context, userID string, permSlug string, grantedBy *string) error {
	return s.repo.UpsertUserOverride(ctx, userID, permSlug, false, "revoked", grantedBy)
}

func (s *Service) ResolveUserAccess(ctx context.Context, userID string) (UserAccess, error) {
	// Check cache
	if cached, ok := s.cache.Load(userID); ok {
		entry := cached.(*PermissionCacheEntry)
		if time.Now().Before(entry.ExpiresAt) {
			return entry.Access, nil
		}
		s.cache.Delete(userID)
	}

	access, err := s.repo.ResolveUserAccess(ctx, userID)
	if err != nil {
		return UserAccess{}, err
	}

	// Store in cache
	s.cache.Store(userID, &PermissionCacheEntry{
		Access:    access,
		ExpiresAt: time.Now().Add(cacheTTL),
	})

	return access, nil
}

func (s *Service) HasPermission(ctx context.Context, userID string, permSlug string) (bool, error) {
	access, err := s.ResolveUserAccess(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, p := range access.Permissions {
		if p == permSlug {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) HasAnyPermission(ctx context.Context, userID string, permSlugs ...string) (bool, error) {
	access, err := s.ResolveUserAccess(ctx, userID)
	if err != nil {
		return false, err
	}
	permSet := make(map[string]struct{}, len(access.Permissions))
	for _, p := range access.Permissions {
		permSet[p] = struct{}{}
	}
	for _, slug := range permSlugs {
		if _, ok := permSet[slug]; ok {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) InvalidateUser(userID string) {
	s.cache.Delete(userID)
}

func (s *Service) InvalidateAll() {
	s.cache.Range(func(key, _ any) bool {
		s.cache.Delete(key)
		return true
	})
}

func (s *Service) cleanupLoop() {
	ticker := time.NewTicker(cleanupEvery)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.cache.Range(func(key, value any) bool {
			entry := value.(*PermissionCacheEntry)
			if now.After(entry.ExpiresAt) {
				s.cache.Delete(key)
			}
			return true
		})
	}
}
