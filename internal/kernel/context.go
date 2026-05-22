package kernel

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type TenantContext struct {
	TenantID        string
	DB              *sqlx.DB
	FeatureFlags    map[string]bool
	ComplianceLevel string
}

type contextKey string

const tenantContextKey contextKey = "tenantContext"

func WithTenantContext(ctx context.Context, tc *TenantContext) context.Context {
	return context.WithValue(ctx, tenantContextKey, tc)
}

func GetTenantContext(ctx context.Context) (*TenantContext, bool) {
	tc, ok := ctx.Value(tenantContextKey).(*TenantContext)
	return tc, ok
}
