package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HMACVerify проверяет HMAC-SHA256 подпись тела запроса.
// Клиент вычисляет HMAC-SHA256(body, secret) и передаёт её
// в заголовке X-Signature в виде hex-строки.
func HMACVerify(secret string) gin.HandlerFunc {
	secretBytes := []byte(secret)

	return func(c *gin.Context) {
		signature := c.GetHeader("X-Signature")
		if signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "отсутствует заголовок X-Signature",
			})
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "не удалось прочитать тело запроса",
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		expectedSig, err := hex.DecodeString(signature)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "неверный формат подписи (ожидается hex)",
			})
			return
		}

		mac := hmac.New(sha256.New, secretBytes)
		mac.Write(body)
		computedSig := mac.Sum(nil)

		if !hmac.Equal(computedSig, expectedSig) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "неверная HMAC-подпись",
			})
			return
		}

		c.Next()
	}
}
