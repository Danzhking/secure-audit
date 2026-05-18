package middleware

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key-for-unit-tests-32ch"

func TestGenerateAndValidateToken(t *testing.T) {
	token, err := GenerateToken(testSecret, "admin", "admin", time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}

	claims, err := validateToken(token, testSecret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}
	if claims.Subject != "admin" {
		t.Errorf("expected sub 'admin', got '%s'", claims.Subject)
	}
	if claims.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", claims.Role)
	}
}

func TestExpiredToken(t *testing.T) {
	token, _ := GenerateToken(testSecret, "admin", "admin", -time.Hour)
	_, err := validateToken(token, testSecret)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Errorf("expected ErrTokenExpired, got: %v", err)
	}
}

func TestInvalidSignature(t *testing.T) {
	token, _ := GenerateToken(testSecret, "admin", "admin", time.Hour)
	_, err := validateToken(token, "wrong-secret-key-not-matching!!")
	if err == nil {
		t.Fatal("expected error for invalid signature")
	}
	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		t.Errorf("expected ErrTokenSignatureInvalid, got: %v", err)
	}
}

func TestMalformedToken(t *testing.T) {
	_, err := validateToken("not.a.valid.jwt.token", testSecret)
	if err == nil {
		t.Fatal("expected error for malformed token")
	}
}
