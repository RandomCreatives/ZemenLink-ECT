package kernel

import (
	"context"
)

type Message struct {
	ID      string
	Content string
	Metadata map[string]interface{}
}

type Interceptor func(ctx context.Context, msg *Message) error

type Pipeline struct {
	interceptors []Interceptor
}

func NewPipeline(interceptors ...Interceptor) *Pipeline {
	return &Pipeline{interceptors: interceptors}
}

// NewMessagePipeline is a factory for the standard messaging interceptors
func NewMessagePipeline() *Pipeline {
	return NewPipeline(
		EncryptionInterceptor,
		EscrowInterceptor,
		PersistenceInterceptor,
		BroadcastInterceptor,
	)
}

// RealTimeProvider is an interface for broadcasting messages
type RealTimeProvider interface {
	BroadcastToTenant(ctx context.Context, tenantID string, payload interface{}) error
}

var GlobalRTProvider RealTimeProvider

// Interceptor implementations (Stubs for now)

func EncryptionInterceptor(ctx context.Context, msg *Message) error {
	// TODO: Signal Protocol implementation
	return nil
}

func EscrowInterceptor(ctx context.Context, msg *Message) error {
	tc, ok := GetTenantContext(ctx)
	if !ok || !tc.FeatureFlags["escrow"] {
		return nil
	}
	// TODO: Multi-encrypt symmetric key for escrow
	return nil
}

func PersistenceInterceptor(ctx context.Context, msg *Message) error {
	// TODO: Final DB write logic if not handled by handler
	return nil
}

func BroadcastInterceptor(ctx context.Context, msg *Message) error {
	tc, ok := GetTenantContext(ctx)
	if !ok || GlobalRTProvider == nil {
		return nil
	}
	return GlobalRTProvider.BroadcastToTenant(ctx, tc.TenantID, msg)
}

func (p *Pipeline) Execute(ctx context.Context, msg *Message) error {
	for _, interceptor := range p.interceptors {
		if err := interceptor(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
