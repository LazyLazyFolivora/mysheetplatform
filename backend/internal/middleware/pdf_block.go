package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BlockPDF 阻止直接访问 PDF 文件，只能通过鉴权下载接口获取
func BlockPDF() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/uploads/sheets_pdf/") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
