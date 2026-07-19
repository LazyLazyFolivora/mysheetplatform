package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/model"
)

func (h *AdminHandler) registerIPBanRoutes(admin *gin.RouterGroup) {
	admin.GET("/banned-ip", h.ListBannedIPs)
	admin.POST("/banned-ip", h.CreateBannedIP)
	admin.DELETE("/banned-ip/:id", h.DeleteBannedIP)
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
