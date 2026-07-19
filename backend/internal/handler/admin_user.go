package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
)

func (h *AdminHandler) registerUserRoutes(admin *gin.RouterGroup) {
	admin.GET("/users", h.ListUsers)
	admin.POST("/users/:id/ban", h.BanUser)
	admin.POST("/users/:id/unban", h.UnbanUser)
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
