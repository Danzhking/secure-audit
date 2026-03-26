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
				"error": "отсутствует заголовок X-API-Key",
			})
			return
		}

		if !constantTimeContains(validKeys, key) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "недействительный API-ключ",
			})
			return
		}

		c.Next()
	}
}

// constantTimeContains сравнивает ключи за постоянное время (защита от timing-атак).
func constantTimeContains(keys []string, candidate string) bool {
	for _, k := range keys {
		if subtle.ConstantTimeCompare([]byte(k), []byte(candidate)) == 1 {
			return true
		}
	}
	return false
}
