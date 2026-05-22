package kernel

import (
	"context"
	"fmt"
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

	// 1. Generate/Retrieve the Symmetric Key for this message
	// For this PoC, we assume the message already contains the encrypted symmetric key for the recipient.
	// We now encrypt that SAME key with the Corporate Escrow Public Key.

	// Placeholder for actual encryption logic
	escrowKeyEncrypted := "escrow-encrypted-key-placeholder"

	// 2. Persist the Escrow Key
	query := `
		INSERT INTO message_keys (message_id, recipient_id, key_encrypted, is_escrow)
		VALUES ($1, NULL, $2, TRUE)
	`
	_, err := tc.DB.ExecContext(ctx, query, msg.ID, escrowKeyEncrypted)
	return err
}

func PersistenceInterceptor(ctx context.Context, msg *Message) error {
	tc, ok := GetTenantContext(ctx)
	if !ok {
		return fmt.Errorf("tenant context required for persistence")
	}

	query := `
		INSERT INTO messages (id, chat_id, sender_id, content_encrypted, content_type, metadata, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tc.DB.ExecContext(ctx, query,
		msg.ID,
		msg.Metadata["chat_id"],
		tc.UserID,
		msg.Content,
		msg.Metadata["content_type"],
		msg.Metadata, // Persist all metadata as JSONB
		"sent",
	)
	return err
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
