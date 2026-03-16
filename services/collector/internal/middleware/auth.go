package middleware

import (
	"crypto/subtle"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(validKeys []string) gin.HandlerFunc {
	keySet := make(map[string]bool, len(validKeys))
	for _, k := range validKeys {
		keySet[k] = true
	}

	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing X-API-Key header",
			})
			return
		}

		if !constantTimeContains(validKeys, key) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid API key",
			})
			return
		}

		c.Next()
	}
}

// constantTimeContains prevents timing attacks by using constant-time comparison.
func constantTimeContains(keys []string, candidate string) bool {
	for _, k := range keys {
		if subtle.ConstantTimeCompare([]byte(k), []byte(candidate)) == 1 {
			return true
		}
	}
	return false
}
