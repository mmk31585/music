package permissions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	apperrors "music/internal/common/errors"
	"music/internal/common/response"
	"music/internal/modules/auth"
)

type Handler struct {
	svc ServiceInterface
}

func NewHandler(svc ServiceInterface) *Handler {
	return &Handler{svc: svc}
}

// --- Permission Endpoints ---

func (h *Handler) ListPermissions(c *gin.Context) error {
	perms, err := h.svc.ListPermissions(c.Request.Context())
	if err != nil {
		return err
	}

	result := make([]PermissionResponse, len(perms))
	for i, p := range perms {
		result[i] = PermissionResponse{
			ID:          p.ID,
			Slug:        p.Slug,
			Category:    p.Category,
			Label:       p.Label,
			Description: p.Description,
		}
	}

	response.Success(c, http.StatusOK, "permissions retrieved", result)
	return nil
}

// --- Role Endpoints ---

func (h *Handler) ListRoles(c *gin.Context) error {
	roles, err := h.svc.ListRoles(c.Request.Context())
	if err != nil {
		return err
	}
	response.Success(c, http.StatusOK, "roles retrieved", roles)
	return nil
}

func (h *Handler) GetRole(c *gin.Context) error {
	slug := c.Param("slug")
	role, err := h.svc.GetRole(c.Request.Context(), slug)
	if err != nil {
		return err
	}

	perms, err := h.svc.GetRolePermissions(c.Request.Context(), slug)
	if err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "role retrieved", RoleWithPermissions{
		Role:        role,
		Permissions: perms,
	})
	return nil
}

func (h *Handler) CreateRole(c *gin.Context) error {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	role := Role{
		Slug:            req.Slug,
		Label:           req.Label,
		Description:     req.Description,
		HierarchyLevel:  req.HierarchyLevel,
	}

	created, err := h.svc.CreateRole(c.Request.Context(), role)
	if err != nil {
		return err
	}

	response.Success(c, http.StatusCreated, "role created", created)
	return nil
}

func (h *Handler) UpdateRole(c *gin.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeBadRequest, "invalid role ID", nil)
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	updates := make(map[string]any)
	if req.Label != nil {
		updates["label"] = *req.Label
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.HierarchyLevel != nil {
		updates["hierarchy_level"] = *req.HierarchyLevel
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.svc.UpdateRole(c.Request.Context(), id, updates); err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "role updated", struct{}{})
	return nil
}

func (h *Handler) DeleteRole(c *gin.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeBadRequest, "invalid role ID", nil)
	}

	if err := h.svc.DeleteRole(c.Request.Context(), id); err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "role deleted", struct{}{})
	return nil
}

func (h *Handler) SetRolePermissions(c *gin.Context) error {
	slug := c.Param("slug")
	var req SetRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	if err := h.svc.SetRolePermissions(c.Request.Context(), slug, req.Permissions); err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "role permissions updated", struct{}{})
	return nil
}

// --- User Role Endpoints ---

func (h *Handler) AssignUserRole(c *gin.Context) error {
	userID := c.Param("id")
	var req AssignUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	assignedBy := auth.UserIDFromContext(c)
	if err := h.svc.AssignUserRole(c.Request.Context(), userID, req.RoleSlug, &assignedBy); err != nil {
		return err
	}

	// Invalidate cache
	h.svc.InvalidateUser(userID)

	response.Success(c, http.StatusOK, "role assigned", UserRoleResponse{
		UserID:   userID,
		RoleSlug: req.RoleSlug,
	})
	return nil
}

func (h *Handler) RemoveUserRole(c *gin.Context) error {
	userID := c.Param("id")
	roleSlug := c.Param("roleSlug")

	if err := h.svc.RemoveUserRole(c.Request.Context(), userID, roleSlug); err != nil {
		return err
	}

	h.svc.InvalidateUser(userID)

	response.Success(c, http.StatusOK, "role removed", struct{}{})
	return nil
}

// --- User Permission Override Endpoints ---

func (h *Handler) GetUserOverrides(c *gin.Context) error {
	userID := c.Param("id")
	overrides, err := h.svc.GetUserOverrides(c.Request.Context(), userID)
	if err != nil {
		return err
	}
	response.Success(c, http.StatusOK, "user overrides retrieved", overrides)
	return nil
}

func (h *Handler) GrantUserPermission(c *gin.Context) error {
	userID := c.Param("id")
	var req GrantPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	grantedBy := auth.UserIDFromContext(c)
	if err := h.svc.GrantUserPermission(c.Request.Context(), userID, req.PermSlug, req.Reason, &grantedBy); err != nil {
		return err
	}

	h.svc.InvalidateUser(userID)

	response.Success(c, http.StatusOK, "permission granted", struct{}{})
	return nil
}

func (h *Handler) RevokeUserPermission(c *gin.Context) error {
	userID := c.Param("id")
	var req RevokePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	grantedBy := auth.UserIDFromContext(c)
	if err := h.svc.RevokeUserPermission(c.Request.Context(), userID, req.PermSlug, &grantedBy); err != nil {
		return err
	}

	h.svc.InvalidateUser(userID)

	response.Success(c, http.StatusOK, "permission revoked", struct{}{})
	return nil
}

// --- Access Resolution Endpoints ---

func (h *Handler) GetUserAccess(c *gin.Context) error {
	userID := c.Param("id")
	access, err := h.svc.ResolveUserAccess(c.Request.Context(), userID)
	if err != nil {
		return err
	}
	response.Success(c, http.StatusOK, "user access resolved", UserAccessResponse(access))
	return nil
}

func (h *Handler) CheckPermission(c *gin.Context) error {
	userID := c.Param("id")
	var req CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return apperrors.New(http.StatusBadRequest, apperrors.CodeValidation, "invalid request", err.Error())
	}

	ok, err := h.svc.HasPermission(c.Request.Context(), userID, req.PermSlug)
	if err != nil {
		return err
	}

	response.Success(c, http.StatusOK, "permission check completed", CheckPermissionResponse{
		HasPermission: ok,
		PermSlug:      req.PermSlug,
	})
	return nil
}

func (h *Handler) GetMyAccess(c *gin.Context) error {
	userID := auth.UserIDFromContext(c)
	if userID == "" {
		return apperrors.New(http.StatusUnauthorized, apperrors.CodeUnauthorized, "authentication required", nil)
	}

	access, err := h.svc.ResolveUserAccess(c.Request.Context(), userID)
	if err != nil {
		return err
	}
	response.Success(c, http.StatusOK, "your access resolved", UserAccessResponse(access))
	return nil
}
