package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, exists := c.Get(ContextUserKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		claims, ok := v.(*Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user context",
			})
			return
		}

		if claims.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "bạn đéo phải người nắm quyền sinh quyền sát",
			})
			return
		}

		c.Next()
	}
}
