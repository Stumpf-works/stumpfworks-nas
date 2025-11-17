// Revision: 2025-11-17 | Author: Claude | Version: 1.0.0
package users

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Stumpf-works/stumpfworks-nas/internal/config"
)

// setupTestConfig creates a test configuration
func setupTestConfig() {
	config.GlobalConfig = &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "test-secret-key-for-jwt-testing-minimum-32-characters-long",
			JWTExpirationHours: 24,
			JWTRefreshHours:    168,
		},
	}
}

// TestGenerateToken tests JWT token generation
func TestGenerateToken(t *testing.T) {
	setupTestConfig()

	testUser := &User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}

	token, err := GenerateToken(testUser)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	// Verify token structure (should have 3 parts: header.payload.signature)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Errorf("Expected token to have 3 parts, got %d", len(parts))
	}
}

// TestValidateToken tests JWT token validation
func TestValidateToken(t *testing.T) {
	setupTestConfig()

	testUser := &User{
		ID:       42,
		Username: "alice",
		Role:     "admin",
	}

	// Generate a valid token
	token, err := GenerateToken(testUser)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Verify claims
	if claims.UserID != testUser.ID {
		t.Errorf("Expected UserID %d, got %d", testUser.ID, claims.UserID)
	}
	if claims.Username != testUser.Username {
		t.Errorf("Expected Username %s, got %s", testUser.Username, claims.Username)
	}
	if claims.Role != testUser.Role {
		t.Errorf("Expected Role %s, got %s", testUser.Role, claims.Role)
	}
	if claims.Issuer != "stumpfworks-nas" {
		t.Errorf("Expected Issuer 'stumpfworks-nas', got %s", claims.Issuer)
	}
}

// TestValidateToken_InvalidToken tests validation of invalid tokens
func TestValidateToken_InvalidToken(t *testing.T) {
	setupTestConfig()

	tests := []struct {
		name        string
		token       string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Empty token",
			token:       "",
			shouldError: true,
			errorMsg:    "failed to parse token",
		},
		{
			name:        "Malformed token",
			token:       "not.a.valid.jwt.token",
			shouldError: true,
			errorMsg:    "failed to parse token",
		},
		{
			name:        "Token with wrong signature",
			token:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImZha2UiLCJyb2xlIjoidXNlciJ9.wrong_signature",
			shouldError: true,
			errorMsg:    "failed to parse token",
		},
		{
			name:        "Random string",
			token:       "thisisnotajwttoken",
			shouldError: true,
			errorMsg:    "failed to parse token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for token '%s', but got none. Claims: %+v", tt.token, claims)
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', but got: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

// TestValidateToken_ExpiredToken tests validation of expired tokens
func TestValidateToken_ExpiredToken(t *testing.T) {
	setupTestConfig()

	// Create an expired token manually
	expirationTime := time.Now().Add(-1 * time.Hour) // Expired 1 hour ago

	claims := &Claims{
		UserID:   1,
		Username: "testuser",
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			Issuer:    "stumpfworks-nas",
			Subject:   "1",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GlobalConfig.Auth.JWTSecret))
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	// Try to validate the expired token
	_, err = ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for expired token, but got none")
	}
	if !strings.Contains(err.Error(), "failed to parse token") {
		t.Errorf("Expected error about expired token, got: %s", err.Error())
	}
}

// TestValidateToken_WrongSigningMethod tests rejection of tokens with wrong signing method
func TestValidateToken_WrongSigningMethod(t *testing.T) {
	setupTestConfig()

	// Create a token with RS256 instead of HS256
	claims := &Claims{
		UserID:   1,
		Username: "testuser",
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "stumpfworks-nas",
			Subject:   "1",
		},
	}

	// Note: This will fail to sign without RSA keys, but the test demonstrates the concept
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("Failed to create token with none algorithm: %v", err)
	}

	// Try to validate the token
	_, err = ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for wrong signing method, but got none")
	}
}

// TestGenerateRefreshToken tests refresh token generation
func TestGenerateRefreshToken(t *testing.T) {
	setupTestConfig()

	testUser := &User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}

	refreshToken, err := GenerateRefreshToken(testUser)
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	if refreshToken == "" {
		t.Fatal("Generated refresh token is empty")
	}

	// Validate the refresh token
	claims, err := ValidateToken(refreshToken)
	if err != nil {
		t.Fatalf("Failed to validate refresh token: %v", err)
	}

	// Verify it's a valid token with correct claims
	if claims.UserID != testUser.ID {
		t.Errorf("Expected UserID %d, got %d", testUser.ID, claims.UserID)
	}

	// Verify expiration is longer than regular token
	// Refresh token should expire in 168 hours (7 days)
	expectedExpiration := time.Now().Add(time.Hour * time.Duration(config.GlobalConfig.Auth.JWTRefreshHours))
	if claims.ExpiresAt.Time.Before(expectedExpiration.Add(-1 * time.Minute)) {
		t.Error("Refresh token expiration is shorter than expected")
	}
}

// TestTokenRolePermissions tests that different roles are correctly encoded
func TestTokenRolePermissions(t *testing.T) {
	setupTestConfig()

	roles := []string{"admin", "user", "guest"}

	for _, role := range roles {
		t.Run("Role_"+role, func(t *testing.T) {
			testUser := &User{
				ID:       1,
				Username: "testuser",
				Role:     role,
			}

			token, err := GenerateToken(testUser)
			if err != nil {
				t.Fatalf("Failed to generate token for role %s: %v", role, err)
			}

			claims, err := ValidateToken(token)
			if err != nil {
				t.Fatalf("Failed to validate token for role %s: %v", role, err)
			}

			if claims.Role != role {
				t.Errorf("Expected role %s, got %s", role, claims.Role)
			}
		})
	}
}

// TestTokenWithDifferentSecrets tests that tokens signed with different secrets fail validation
func TestTokenWithDifferentSecrets(t *testing.T) {
	// Create token with one secret
	config.GlobalConfig = &config.Config{
		Auth: config.AuthConfig{
			JWTSecret:          "original-secret-key-for-testing-must-be-32-chars-long",
			JWTExpirationHours: 24,
			JWTRefreshHours:    168,
		},
	}

	testUser := &User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}

	token, err := GenerateToken(testUser)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Change the secret
	config.GlobalConfig.Auth.JWTSecret = "different-secret-key-for-testing-must-be-32-chars"

	// Try to validate with different secret
	_, err = ValidateToken(token)
	if err == nil {
		t.Error("Expected validation to fail with different secret, but it succeeded")
	}
	if !strings.Contains(err.Error(), "failed to parse token") {
		t.Errorf("Expected signature verification error, got: %s", err.Error())
	}
}

// BenchmarkGenerateToken benchmarks token generation performance
func BenchmarkGenerateToken(b *testing.B) {
	setupTestConfig()

	testUser := &User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateToken(testUser)
	}
}

// BenchmarkValidateToken benchmarks token validation performance
func BenchmarkValidateToken(b *testing.B) {
	setupTestConfig()

	testUser := &User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}

	token, _ := GenerateToken(testUser)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ValidateToken(token)
	}
}
