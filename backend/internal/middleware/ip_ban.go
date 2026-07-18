package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/model"
)

func IPBan(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := clientIP(c)

		var banned model.BannedIp
		err := db.Where("ip = ?", ip).First(&banned).Error
		if err != nil {
			c.Next()
			return
		}

		if banned.ExpiresAt != nil && banned.ExpiresAt.Before(time.Now()) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, response.Error(403, "该IP已被封禁"))
	}
}

func clientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return ip
}
