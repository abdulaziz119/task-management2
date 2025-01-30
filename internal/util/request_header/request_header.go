package request_header

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetTokenFromHeader(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if header == "" {
		return ""
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
