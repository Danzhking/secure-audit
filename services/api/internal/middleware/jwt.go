package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Claims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
}

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := validateToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user", claims.Sub)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func GenerateToken(secret, sub, role string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Sub:  sub,
		Role: role,
		Exp:  now.Add(ttl).Unix(),
		Iat:  now.Unix(),
	}

	header := base64url([]byte(`{"alg":"HS256","typ":"JWT"}`))

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payload := base64url(claimsJSON)

	sig := sign(header+"."+payload, secret)
	return header + "." + payload + "." + sig, nil
}

func validateToken(token, secret string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errInvalid("malformed token")
	}

	expectedSig := sign(parts[0]+"."+parts[1], secret)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, errInvalid("invalid signature")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errInvalid("invalid payload encoding")
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, errInvalid("invalid payload")
	}

	if time.Now().Unix() > claims.Exp {
		return nil, errInvalid("token expired")
	}

	return &claims, nil
}

func sign(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64url(h.Sum(nil))
}

func base64url(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

type tokenError struct{ msg string }

func (e *tokenError) Error() string  { return e.msg }
func errInvalid(msg string) error    { return &tokenError{msg: msg} }
