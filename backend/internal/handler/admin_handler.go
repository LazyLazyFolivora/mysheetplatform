package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/service"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	stats, err := h.adminService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取统计数据失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(stats))
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, total, err := h.adminService.ListUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取用户列表失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(response.PageResponse{List: users, Total: total, Page: page, PageSize: pageSize}))
}

func (h *AdminHandler) BanUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.adminService.BanUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "操作失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *AdminHandler) UnbanUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.adminService.UnbanUser(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "操作失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *AdminHandler) ListMessages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	msgs, total, err := h.adminService.ListContactMessages(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取消息列表失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(response.PageResponse{List: msgs, Total: total, Page: page, PageSize: pageSize}))
}

func (h *AdminHandler) MarkMessageRead(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.adminService.MarkMessageRead(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "操作失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *AdminHandler) DeleteMessage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.adminService.DeleteMessage(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "操作失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *AdminHandler) CreateContactMessage(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required,max=100"`
		Email   string `json:"email" binding:"omitempty,max=100"`
		Subject string `json:"subject" binding:"omitempty,max=200"`
		Content string `json:"content" binding:"required,max=2000"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请填写姓名和留言内容"))
		return
	}

	msg := &model.ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Content: req.Content,
	}
	if err := h.adminService.CreateContactMessage(msg); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "提交失败，请稍后重试"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

// ListPublicFriendLinks 前台友情链接，只返回启用的
func (h *AdminHandler) ListPublicFriendLinks(c *gin.Context) {
	links, err := h.adminService.ListEnabledFriendLinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取友情链接失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(links))
}

func (h *AdminHandler) ListFriendLinks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	links, total, err := h.adminService.ListFriendLinks(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取友情链接失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(response.PageResponse{List: links, Total: total, Page: page, PageSize: pageSize}))
}

type friendLinkReq struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	URL         string `json:"url" binding:"required,url,max=500"`
	Description string `json:"description" binding:"max=500"`
	Logo        string `json:"logo" binding:"max=500"`
	IsEnabled   *bool  `json:"is_enabled"`
	SortOrder   int    `json:"sort_order" binding:"min=0"`
}

func (h *AdminHandler) CreateFriendLink(c *gin.Context) {
	var req friendLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	enabled := true
	if req.IsEnabled != nil {
		enabled = *req.IsEnabled
	}
	link := model.FriendLink{
		Name:        req.Name,
		URL:         req.URL,
		Description: req.Description,
		Logo:        req.Logo,
		IsEnabled:   enabled,
		SortOrder:   req.SortOrder,
	}
	if err := h.adminService.CreateFriendLink(&link); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "创建失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(link))
}

func (h *AdminHandler) UpdateFriendLink(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req friendLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	link, err := h.adminService.GetFriendLink(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "链接不存在"))
		return
	}
	link.Name = req.Name
	link.URL = req.URL
	link.Description = req.Description
	link.Logo = req.Logo
	if req.IsEnabled != nil {
		link.IsEnabled = *req.IsEnabled
	}
	link.SortOrder = req.SortOrder
	if err := h.adminService.UpdateFriendLink(link); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "更新失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(link))
}

func (h *AdminHandler) DeleteFriendLink(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.adminService.DeleteFriendLink(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *AdminHandler) ListBannedIPs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	ips, total, err := h.adminService.ListBannedIPs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取封禁IP列表失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(response.PageResponse{List: ips, Total: total, Page: page, PageSize: pageSize}))
}

func (h *AdminHandler) CreateBannedIP(c *gin.Context) {
	var ip model.BannedIp
	if err := c.ShouldBindJSON(&ip); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	if err := h.adminService.CreateBannedIP(&ip); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "操作失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(ip))
}

func (h *AdminHandler) DeleteBannedIP(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.adminService.DeleteBannedIP(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *AdminHandler) ListSystemConfigs(c *gin.Context) {
	configs, err := h.adminService.ListSystemConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取系统配置失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(configs))
}

func (h *AdminHandler) UpdateSystemConfig(c *gin.Context) {
	var req struct {
		Key    string `json:"key" binding:"required,max=100"`
		Value  string `json:"value"`
		Remark string `json:"remark" binding:"max=200"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	if err := h.adminService.UpdateSystemConfig(req.Key, req.Value, req.Remark); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "更新失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *AdminHandler) CreateSystemConfig(c *gin.Context) {
	var req struct {
		Key    string `json:"key" binding:"required,min=1,max=100"`
		Value  string `json:"value" binding:"required"`
		Remark string `json:"remark" binding:"max=200"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	cfg := model.SystemConfig{ConfigKey: req.Key, ConfigVal: req.Value, Remark: req.Remark}
	if err := h.adminService.CreateSystemConfig(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(cfg))
}

func (h *AdminHandler) DeleteSystemConfig(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	if err := h.adminService.DeleteSystemConfig(key); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}
