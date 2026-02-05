package auth

import (
	"context"
	"testing"
	"time"
)

func TestAuthenticator_GenerateAndValidateToken(t *testing.T) {
	auth := NewAuthenticator("test-secret-key-12345", nil)

	user := &User{
		ID:    "user-1",
		Email: "test@example.com",
		Name:  "Test User",
		Role:  RoleUser,
	}

	token, err := auth.GenerateToken(user, 1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	// Validate the token
	validated, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if validated.ID != user.ID {
		t.Errorf("expected user ID '%s', got '%s'", user.ID, validated.ID)
	}
	if validated.Email != user.Email {
		t.Errorf("expected email '%s', got '%s'", user.Email, validated.Email)
	}
	if validated.Role != user.Role {
		t.Errorf("expected role '%s', got '%s'", user.Role, validated.Role)
	}
}

func TestAuthenticator_ExpiredToken(t *testing.T) {
	auth := NewAuthenticator("test-secret-key-12345", nil)

	user := &User{
		ID:   "user-1",
		Role: RoleUser,
	}

	// Generate a token that's already expired
	token, err := auth.GenerateToken(user, -1*time.Hour)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = auth.ValidateToken(token)
	if err == nil {
		t.Error("expected error for expired token")
	}
}

func TestAuthenticator_InvalidSignature(t *testing.T) {
	auth1 := NewAuthenticator("secret-1", nil)
	auth2 := NewAuthenticator("secret-2", nil)

	user := &User{
		ID:   "user-1",
		Role: RoleUser,
	}

	token, _ := auth1.GenerateToken(user, 1*time.Hour)

	// Try to validate with different secret
	_, err := auth2.ValidateToken(token)
	if err == nil {
		t.Error("expected error for token signed with different secret")
	}
}

func TestAuthenticator_InvalidTokenFormat(t *testing.T) {
	auth := NewAuthenticator("test-secret", nil)

	testCases := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"single part", "abc"},
		{"two parts", "abc.def"},
		{"invalid base64", "abc.def.ghi"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := auth.ValidateToken(tc.token)
			if err == nil {
				t.Errorf("expected error for token: %s", tc.token)
			}
		})
	}
}

func TestAuthenticator_APIKeyValidation(t *testing.T) {
	keys := []string{"key-1", "key-2", "key-3"}
	auth := NewAuthenticator("secret", keys)

	t.Run("valid API key", func(t *testing.T) {
		if !auth.ValidateAPIKey("key-1") {
			t.Error("expected key-1 to be valid")
		}
	})

	t.Run("invalid API key", func(t *testing.T) {
		if auth.ValidateAPIKey("key-invalid") {
			t.Error("expected key-invalid to be invalid")
		}
	})

	t.Run("empty API key", func(t *testing.T) {
		if auth.ValidateAPIKey("") {
			t.Error("expected empty key to be invalid")
		}
	})
}

func TestUserFromContext(t *testing.T) {
	t.Run("returns user from context", func(t *testing.T) {
		user := &User{ID: "user-1", Role: RoleAdmin}
		ctx := context.WithValue(context.Background(), userContextKey, user)

		result := UserFromContext(ctx)
		if result.ID != "user-1" {
			t.Errorf("expected user ID 'user-1', got '%s'", result.ID)
		}
		if result.Role != RoleAdmin {
			t.Errorf("expected role 'admin', got '%s'", result.Role)
		}
	})

	t.Run("returns anonymous for empty context", func(t *testing.T) {
		result := UserFromContext(context.Background())
		if result.Role != RoleAnonymous {
			t.Errorf("expected anonymous role, got '%s'", result.Role)
		}
	})
}

func TestRequireRole(t *testing.T) {
	adminCtx := context.WithValue(context.Background(), userContextKey, &User{Role: RoleAdmin})
	userCtx := context.WithValue(context.Background(), userContextKey, &User{Role: RoleUser})
	anonCtx := context.Background()

	t.Run("admin has admin access", func(t *testing.T) {
		if !RequireRole(adminCtx, RoleAdmin) {
			t.Error("admin should have admin access")
		}
	})

	t.Run("user does not have admin access", func(t *testing.T) {
		if RequireRole(userCtx, RoleAdmin) {
			t.Error("user should not have admin access")
		}
	})

	t.Run("anonymous does not have user access", func(t *testing.T) {
		if RequireRole(anonCtx, RoleUser) {
			t.Error("anonymous should not have user access")
		}
	})

	t.Run("multiple roles", func(t *testing.T) {
		if !RequireRole(userCtx, RoleAdmin, RoleUser) {
			t.Error("user should match when RoleUser is in the list")
		}
	})
}
