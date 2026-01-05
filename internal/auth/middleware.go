package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextUserKey = "user"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(ContextUserKey, claims)
		c.Next()
	}
}
func humanFileSize(size int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
	)

	switch {
	case size >= MB:
		return fmt.Sprintf("%d MB", size/MB)
	case size >= KB:
		return fmt.Sprintf("%d KB", size/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func LimitUploadSize(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			maxSize,
		)

		if err := c.Request.ParseMultipartForm(maxSize); err != nil {
			c.AbortWithStatusJSON(
				http.StatusRequestEntityTooLarge,
				gin.H{
					"error": fmt.Sprintf(
						"File vượt quá dung lượng cho phép (%s)",
						humanFileSize(maxSize),
					),
				},
			)
			return
		}

		c.Next()
	}
}
