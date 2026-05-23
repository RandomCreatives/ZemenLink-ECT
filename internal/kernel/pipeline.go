package kernel

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

type Message struct {
	ID           string
	Content      string
	RecipientID  string                 // Target user for 1:1, or NULL/empty for groups
	EncryptedKey string                 // The symmetric key encrypted for the recipient
	Metadata     map[string]interface{} // Includes chat_id, content_type, extra
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

// GlobalEscrowCrypto is a placeholder for tenant-specific crypto services
var GlobalEscrowCrypto interface {
	EncryptForKeyEscrow(plaintextKey string) (string, error)
}

// Interceptor implementations (Stubs for now)

func EncryptionInterceptor(ctx context.Context, msg *Message) error {
	// SIMULATED E2EE
	// IN PRODUCTION: Full implementation of the Signal Protocol (Double Ratchet)
	// For this foundation, we simulate the "encrypted payload" format using Base64.
	if msg.Content != "" && !strings.HasPrefix(msg.Content, "E2EE:") {
		msg.Content = fmt.Sprintf("E2EE:v1:%s", base64.StdEncoding.EncodeToString([]byte(msg.Content)))
	}
	return nil
}

func EscrowInterceptor(ctx context.Context, msg *Message) error {
	tc, ok := GetTenantContext(ctx)
	if !ok || !tc.FeatureFlags["escrow"] {
		return nil
	}

	if msg.EncryptedKey == "" {
		return nil
	}

	// Determine the escrow crypto service
	var crypto *EscrowCrypto
	if tc.EscrowPublicKey != "" {
		var err error
		crypto, err = NewEscrowCrypto([]byte(tc.EscrowPublicKey))
		if err != nil {
			log.Printf("Warning: failed to load tenant escrow key: %v", err)
		}
	}

	escrowKeyEncrypted := "escrow-encrypted-key-placeholder"
	if crypto != nil {
		var err error
		// Simulate multi-encryption
		escrowKeyEncrypted, err = crypto.EncryptForKeyEscrow("simulated-symmetric-key")
		if err != nil {
			return fmt.Errorf("escrow encryption failed: %w", err)
		}
	} else if GlobalEscrowCrypto != nil {
		var err error
		// Fallback to global key if set
		escrowKeyEncrypted, err = GlobalEscrowCrypto.(interface {
			EncryptForKeyEscrow(string) (string, error)
		}).EncryptForKeyEscrow("simulated-symmetric-key")
		if err != nil {
			return fmt.Errorf("escrow encryption failed: %w", err)
		}
	}

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

	tx, err := tc.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert Message
	query := `
		INSERT INTO messages (id, chat_id, sender_id, content_encrypted, content_type, metadata, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.ExecContext(ctx, query,
		msg.ID,
		msg.Metadata["chat_id"],
		tc.UserID,
		msg.Content,
		msg.Metadata["content_type"],
		msg.Metadata,
		"sent",
	)
	if err != nil {
		return err
	}

	// 2. Insert Recipient Key (if provided)
	if msg.RecipientID != "" && msg.EncryptedKey != "" {
		keyQuery := `
			INSERT INTO message_keys (message_id, recipient_id, key_encrypted, is_escrow)
			VALUES ($1, $2, $3, FALSE)
		`
		_, err = tx.ExecContext(ctx, keyQuery, msg.ID, msg.RecipientID, msg.EncryptedKey)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
