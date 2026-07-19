package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/sheet-platform/backend/internal/dto/request"
	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请求参数错误"))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请求参数错误"))
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(resp))
}

func (h *AuthHandler) SendCode(c *gin.Context) {
	var req request.SendCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请求参数错误"))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请求参数错误"))
		return
	}

	if err := h.authService.SendVerificationCode(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success("验证码已发送"))
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	profile, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "获取用户信息失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(profile))
}

func (h *AuthHandler) RegisterRoutes(public, auth, admin *gin.RouterGroup) {
	public.POST("/auth/login", h.Login)
	public.POST("/auth/register", h.Register)
	public.POST("/auth/send-code", h.SendCode)
	auth.GET("/profile", h.Profile)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请求参数错误"))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "请求参数错误"))
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(resp))
}
