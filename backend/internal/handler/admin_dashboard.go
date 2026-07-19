package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
)

func (h *AdminHandler) registerDashboardRoutes(admin *gin.RouterGroup) {
	admin.GET("/dashboard", h.Dashboard)
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	stats, err := h.adminService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取统计数据失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(stats))
}
