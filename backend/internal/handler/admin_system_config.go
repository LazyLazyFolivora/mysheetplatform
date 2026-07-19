package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/model"
)

func (h *AdminHandler) registerSystemConfigRoutes(admin *gin.RouterGroup) {
	admin.GET("/system-config", h.ListSystemConfigs)
	admin.POST("/system-config", h.CreateSystemConfig)
	admin.PUT("/system-config", h.UpdateSystemConfig)
	admin.DELETE("/system-config/:key", h.DeleteSystemConfig)
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
