package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/model"
	"github.com/sheet-platform/backend/internal/pkg"
)

func Auth(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(401, "未登录"))
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := pkg.ParseJWT(cfg.JWT.Secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(401, "登录已过期，请重新登录"))
			return
		}

		// 校验账号状态，防止被封禁用户在 token 过期前继续访问
		var user model.User
		if err := db.Select("status").First(&user, claims.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(401, "用户不存在"))
			return
		}
		if user.Status != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error(403, "账号已被禁用"))
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_id", claims.RoleID)
		c.Next()
	}
}

func AdminAuth(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(401, "未登录"))
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := pkg.ParseJWT(cfg.JWT.Secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(401, "登录已过期"))
			return
		}

		// 以数据库为准校验管理员身份，避免伪造/过期的声明绕过权限
		var user model.User
		if err := db.Select("is_admin", "status").First(&user, claims.UserID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error(401, "用户不存在"))
			return
		}
		if !user.IsAdmin || user.Status != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error(403, "没有管理员权限"))
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", true)
		c.Next()
	}
}
