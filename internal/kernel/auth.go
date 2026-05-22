package kernel

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

// TenantResolver is an interface to fetch tenant details during authentication
type TenantResolver interface {
	GetDB(ctx context.Context, tenantID string) (*sqlx.DB, error)
	GetTenant(ctx context.Context, tenantID string) (interface{}, error) // Returns a tenant object with compliance info
}

type FeatureProvider interface {
	GetFlagsForTenant(tenantID string, complianceLevel string) map[string]bool
}

// AuthMiddleware validates JWT and populates TenantContext
func AuthMiddleware(secret []byte, resolver TenantResolver, features FeatureProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid claims", http.StatusUnauthorized)
				return
			}

			tenantID, ok := claims["tenant_id"].(string)
			if !ok {
				http.Error(w, "tenant_id missing in token", http.StatusUnauthorized)
				return
			}

			userID, _ := claims["user_id"].(string) // Optional for now, but good practice

			// Resolve tenant-specific resources
			db, err := resolver.GetDB(r.Context(), tenantID)
			if err != nil {
				http.Error(w, "Failed to resolve tenant resources", http.StatusInternalServerError)
				return
			}

			// Resolve tenant metadata
			tenantObj, err := resolver.GetTenant(r.Context(), tenantID)
			if err != nil {
				http.Error(w, "Failed to resolve tenant metadata", http.StatusInternalServerError)
				return
			}

			// We need a way to access fields from the tenant object (e.g., via interface or type assertion)
			// For this example, we assume it has a ComplianceLevel. In production, use a shared struct.
			compliance := "standard"
			if t, ok := tenantObj.(interface{ GetCompliance() string }); ok {
				compliance = t.GetCompliance()
			}

			role, _ := claims["role"].(string)

			tc := &TenantContext{
				TenantID:        tenantID,
				UserID:          userID,
				DB:              db,
				ComplianceLevel: compliance,
				FeatureFlags:    features.GetFlagsForTenant(tenantID, compliance),
			}

			// Add role to context
			ctx := WithUserRole(r.Context(), role)
			ctx = WithTenantContext(ctx, tc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GenerateTestToken is a helper for testing
func GenerateTestToken(secret []byte, tenantID, userID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tenant_id": tenantID,
		"user_id":   userID,
		"role":      role,
		"exp":       time.Now().Add(time.Hour).Unix(),
	})
	return token.SignedString(secret)
}
