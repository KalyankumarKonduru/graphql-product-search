package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

// Role defines the authorization level of a user.
type Role string

const (
	RoleAnonymous Role = "anonymous"
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
)

// User represents an authenticated user extracted from a JWT token.
type User struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  Role   `json:"role"`
}

type contextKey string

const userContextKey contextKey = "auth_user"

// JWTHeader is the standard JWT header for HMAC-SHA256.
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// JWTClaims represents the payload of a JWT token.
type JWTClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Iat   int64  `json:"iat"`
	Exp   int64  `json:"exp"`
	Iss   string `json:"iss"`
}

// Authenticator provides JWT-based authentication and API key validation.
type Authenticator struct {
	jwtSecret []byte
	apiKeys   map[string]bool
	issuer    string
}

// NewAuthenticator creates a new Authenticator with the given secret and API keys.
func NewAuthenticator(jwtSecret string, apiKeys []string) *Authenticator {
	keyMap := make(map[string]bool, len(apiKeys))
	for _, key := range apiKeys {
		if key != "" {
			keyMap[key] = true
		}
	}

	return &Authenticator{
		jwtSecret: []byte(jwtSecret),
		apiKeys:   keyMap,
		issuer:    "graphql-product-search",
	}
}

// GenerateToken creates a new JWT token for the given user.
func (a *Authenticator) GenerateToken(user *User, duration time.Duration) (string, error) {
	now := time.Now()

	claims := JWTClaims{
		Sub:   user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  string(user.Role),
		Iat:   now.Unix(),
		Exp:   now.Add(duration).Unix(),
		Iss:   a.issuer,
	}

	header := JWTHeader{Alg: "HS256", Typ: "JWT"}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signingInput := headerEncoded + "." + claimsEncoded
	signature := a.sign([]byte(signingInput))
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	return signingInput + "." + signatureEncoded, nil
}

// ValidateToken verifies and parses a JWT token string.
func (a *Authenticator) ValidateToken(tokenStr string) (*User, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	// Verify signature
	signingInput := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid token signature encoding")
	}

	expectedSig := a.sign([]byte(signingInput))
	if !hmac.Equal(signature, expectedSig) {
		return nil, errors.New("invalid token signature")
	}

	// Decode claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token claims encoding")
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, errors.New("invalid token claims")
	}

	// Check expiration
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}

	// Check issuer
	if claims.Iss != a.issuer {
		return nil, errors.New("invalid token issuer")
	}

	return &User{
		ID:    claims.Sub,
		Email: claims.Email,
		Name:  claims.Name,
		Role:  Role(claims.Role),
	}, nil
}

func (a *Authenticator) sign(data []byte) []byte {
	mac := hmac.New(sha256.New, a.jwtSecret)
	mac.Write(data)
	return mac.Sum(nil)
}

// ValidateAPIKey checks if the given API key is valid.
func (a *Authenticator) ValidateAPIKey(key string) bool {
	return a.apiKeys[key]
}

// Middleware extracts authentication credentials from the request
// and adds the user to the context. It supports both JWT Bearer tokens
// and API key authentication. Unauthenticated requests are allowed
// through as anonymous users - authorization is enforced at the
// resolver level.
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := &User{Role: RoleAnonymous}

		// Check for Bearer token
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				if u, err := a.ValidateToken(token); err == nil {
					user = u
				}
			}
		}

		// Check for API key (used by service-to-service communication)
		if apiKey := r.Header.Get("X-API-Key"); apiKey != "" && user.Role == RoleAnonymous {
			if a.ValidateAPIKey(apiKey) {
				user = &User{
					ID:   "service",
					Name: "API Service",
					Role: RoleAdmin,
				}
			}
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromContext retrieves the authenticated user from the context.
func UserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value(userContextKey).(*User); ok {
		return user
	}
	return &User{Role: RoleAnonymous}
}

// RequireRole checks if the current user has the required role.
func RequireRole(ctx context.Context, roles ...Role) bool {
	user := UserFromContext(ctx)
	for _, role := range roles {
		if user.Role == role {
			return true
		}
	}
	return false
}
