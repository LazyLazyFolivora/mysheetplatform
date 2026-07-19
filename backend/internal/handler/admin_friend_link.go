package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/service"
)

func (h *AdminHandler) registerFriendLinkRoutes(admin *gin.RouterGroup) {
	admin.GET("/friend-links", h.ListFriendLinks)
	admin.POST("/friend-links", h.CreateFriendLink)
	admin.PUT("/friend-links/:id", h.UpdateFriendLink)
	admin.DELETE("/friend-links/:id", h.DeleteFriendLink)
}

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

func (h *AdminHandler) CreateFriendLink(c *gin.Context) {
	var req service.CreateFriendLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	link, err := h.adminService.CreateFriendLink(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "创建失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(link))
}

func (h *AdminHandler) UpdateFriendLink(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req service.UpdateFriendLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}
	if err := h.adminService.UpdateFriendLink(uint(id), req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "更新失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}

func (h *AdminHandler) DeleteFriendLink(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.adminService.DeleteFriendLink(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(nil))
}
