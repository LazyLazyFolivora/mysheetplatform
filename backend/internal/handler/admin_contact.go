package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/model"
)

func (h *AdminHandler) registerContactRoutes(admin *gin.RouterGroup) {
	admin.GET("/contact-messages", h.ListMessages)
	admin.POST("/contact-messages/:id/read", h.MarkMessageRead)
	admin.DELETE("/contact-messages/:id", h.DeleteMessage)
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
