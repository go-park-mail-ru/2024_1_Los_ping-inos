package requests

import (
	"strings"
	"testing"
	"time"

	"main.go/internal/types"
)

func TestCreateCSRFToken(t *testing.T) {
	SID := "session123"
	UID := types.UserID(12345)
	tokenExpTime := time.Now().Add(1 * time.Hour).Unix()

	token, err := CreateCSRFToken(SID, UID, tokenExpTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token == "" {
		t.Error("expected a non-empty token")
	}

	// Split the token to check its parts
	tokenParts := strings.Split(token, ":")
	if len(tokenParts) != 2 {
		t.Errorf("expected token to have 2 parts, got %d", len(tokenParts))
	}
}

func TestCheckCSRFToken(t *testing.T) {
	SID := "session123"
	UID := types.UserID(12345)
	tokenExpTime := time.Now().Add(1 * time.Hour).Unix()

	validToken, err := CreateCSRFToken(SID, UID, tokenExpTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test valid token
	valid, err := CheckCSRFToken(SID, UID, validToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !valid {
		t.Error("expected token to be valid")
	}

	// Test expired token
	expiredTokenExpTime := time.Now().Add(-1 * time.Hour).Unix()
	expiredToken, err := CreateCSRFToken(SID, UID, expiredTokenExpTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	valid, err = CheckCSRFToken(SID, UID, expiredToken)
	if err == nil {
		t.Error("expected error for expired token")
	}
	if valid {
		t.Error("expected token to be invalid due to expiration")
	}

	// Test invalid token format
	invalidToken := "invalid_token_format"
	valid, err = CheckCSRFToken(SID, UID, invalidToken)
	if err == nil {
		t.Error("expected error for invalid token format")
	}
	if valid {
		t.Error("expected token to be invalid due to bad format")
	}

	// Test token with wrong MAC
	tokenParts := strings.Split(validToken, ":")
	wrongToken := tokenParts[0][:len(tokenParts[0])-1] + "x:" + tokenParts[1]

	valid, err = CheckCSRFToken(SID, UID, wrongToken)
	if err == nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valid {
		t.Error("expected token to be invalid due to wrong MAC")
	}
}
