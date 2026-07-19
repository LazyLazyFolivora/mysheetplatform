package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/service"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) RegisterRoutes(public, auth, admin *gin.RouterGroup) {
	public.GET("/friend-links", h.ListPublicFriendLinks)
	public.POST("/contact-messages", h.CreateContactMessage)
	h.registerDashboardRoutes(admin)
	h.registerUserRoutes(admin)
	h.registerContactRoutes(admin)
	h.registerFriendLinkRoutes(admin)
	h.registerIPBanRoutes(admin)
	h.registerSystemConfigRoutes(admin)
}
