package kernel

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type TenantContext struct {
	TenantID        string
	UserID          string
	DB              *sqlx.DB
	FeatureFlags    map[string]bool
	ComplianceLevel string
	EscrowPublicKey string
}

type contextKey string

const (
	tenantContextKey contextKey = "tenantContext"
	userRoleKey      contextKey = "userRole"
)

func WithTenantContext(ctx context.Context, tc *TenantContext) context.Context {
	return context.WithValue(ctx, tenantContextKey, tc)
}

func GetTenantContext(ctx context.Context) (*TenantContext, bool) {
	tc, ok := ctx.Value(tenantContextKey).(*TenantContext)
	return tc, ok
}

func WithUserRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, userRoleKey, role)
}

func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(userRoleKey).(string)
	return role, ok
}
