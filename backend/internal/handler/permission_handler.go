package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/service"
)

type PermissionHandler struct {
	permService *service.PermissionService
}

func NewPermissionHandler(permService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permService: permService}
}

func parseIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return 0, false
	}
	return uint(id), true
}

// ===== 权限模块 =====

func (h *PermissionHandler) Tree(c *gin.Context) {
	tree, err := h.permService.BuildTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取权限树失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(tree))
}

func (h *PermissionHandler) CreateModule(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required,max=100"`
		ParentID *uint  `json:"parent_id"`
		Path     string `json:"path" binding:"max=200"`
		Icon     string `json:"icon" binding:"max=50"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}

	module := &model.PermissionModule{
		Name:     req.Name,
		ParentID: req.ParentID,
		Path:     req.Path,
		Icon:     req.Icon,
	}

	if err := h.permService.CreateModule(module); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(module))
}

func (h *PermissionHandler) UpdateModule(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}

	var req struct {
		Name     string `json:"name" binding:"required,max=100"`
		ParentID *uint  `json:"parent_id"`
		Path     string `json:"path" binding:"max=200"`
		Icon     string `json:"icon" binding:"max=50"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}

	if req.ParentID != nil && *req.ParentID == id {
		c.JSON(http.StatusBadRequest, response.Error(400, "上级模块不能是自身"))
		return
	}

	module, err := h.permService.GetModule(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "模块不存在"))
		return
	}
	module.Name = req.Name
	module.ParentID = req.ParentID
	module.Path = req.Path
	module.Icon = req.Icon

	if err := h.permService.UpdateModule(module); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(module))
}

func (h *PermissionHandler) DeleteModule(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.permService.DeleteModule(id); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *PermissionHandler) GetModulePermissions(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	ids, err := h.permService.GetModulePermissionIDs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取模块权限失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(ids))
}

func (h *PermissionHandler) AssignModulePermissions(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req struct {
		PermIDs []uint `json:"perm_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	if err := h.permService.SetModulePermissions(id, req.PermIDs); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

// ===== 权限 =====

func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	perms, err := h.permService.ListPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取权限列表失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(perms))
}

type permissionReq struct {
	Name     string `json:"name" binding:"required,max=100"`
	Code     string `json:"code" binding:"required,max=100"`
	URL      string `json:"url" binding:"max=200"`
	Method   string `json:"method" binding:"omitempty,oneof=GET POST PUT DELETE"`
	ModuleID *uint  `json:"module_id"`
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req permissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	perm := &model.Permission{
		Name:   req.Name,
		Code:   req.Code,
		URL:    req.URL,
		Method: req.Method,
	}
	if err := h.permService.CreatePermission(perm, req.ModuleID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(perm))
}

func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req permissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	perm := &model.Permission{
		ID:     id,
		Name:   req.Name,
		Code:   req.Code,
		URL:    req.URL,
		Method: req.Method,
	}
	if err := h.permService.UpdatePermission(perm, req.ModuleID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(perm))
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.permService.DeletePermission(id); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

// ===== 角色 =====

func (h *PermissionHandler) ListRoles(c *gin.Context) {
	roles, err := h.permService.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取角色列表失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(roles))
}

type roleReq struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=200"`
}

func (h *PermissionHandler) CreateRole(c *gin.Context) {
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	role, err := h.permService.CreateRole(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(role))
}

func (h *PermissionHandler) UpdateRole(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req roleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	role, err := h.permService.UpdateRole(id, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(role))
}

func (h *PermissionHandler) DeleteRole(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.permService.DeleteRole(id); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *PermissionHandler) GetRoleModules(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	ids, err := h.permService.GetRoleModuleIDs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取角色模块失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(ids))
}

func (h *PermissionHandler) AssignRoleModules(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req struct {
		ModuleIDs []uint `json:"module_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	if err := h.permService.SetRoleModules(id, req.ModuleIDs); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}
