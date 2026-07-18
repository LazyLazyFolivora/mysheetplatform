package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/service"
)

func RequirePermission(permService *service.PermissionService, moduleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(401, "未登录"))
			return
		}

		allowed, err := permService.CheckAccess(userID.(uint), moduleCode)
		if err != nil || !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error(403, "权限不足"))
			return
		}

		c.Next()
	}
}
