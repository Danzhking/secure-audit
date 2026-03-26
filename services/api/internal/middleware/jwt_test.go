package middleware

import (
	"testing"
	"time"
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
	if claims.Sub != "admin" {
		t.Errorf("expected sub 'admin', got '%s'", claims.Sub)
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
	if err.Error() != "срок действия токена истёк" {
		t.Errorf("expected 'срок действия токена истёк', got '%s'", err.Error())
	}
}

func TestInvalidSignature(t *testing.T) {
	token, _ := GenerateToken(testSecret, "admin", "admin", time.Hour)
	_, err := validateToken(token, "wrong-secret-key-not-matching!!")
	if err == nil {
		t.Fatal("expected error for invalid signature")
	}
	if err.Error() != "неверная подпись токена" {
		t.Errorf("expected 'неверная подпись токена', got '%s'", err.Error())
	}
}

func TestMalformedToken(t *testing.T) {
	_, err := validateToken("not.a.valid.jwt.token", testSecret)
	if err == nil {
		t.Fatal("expected error for malformed token")
	}
}
