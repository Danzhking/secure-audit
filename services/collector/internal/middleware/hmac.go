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

// HMACVerify validates the HMAC-SHA256 signature of the request body.
// The client must compute HMAC-SHA256(body, secret) and send it
// in the X-Signature header as a hex-encoded string.
func HMACVerify(secret string) gin.HandlerFunc {
	secretBytes := []byte(secret)

	return func(c *gin.Context) {
		signature := c.GetHeader("X-Signature")
		if signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing X-Signature header",
			})
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "failed to read request body",
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		expectedSig, err := hex.DecodeString(signature)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid signature format (expected hex)",
			})
			return
		}

		mac := hmac.New(sha256.New, secretBytes)
		mac.Write(body)
		computedSig := mac.Sum(nil)

		if !hmac.Equal(computedSig, expectedSig) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid HMAC signature",
			})
			return
		}

		c.Next()
	}
}
