package kernel

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// EscrowCrypto handles public-key encryption for corporate escrow
type EscrowCrypto struct {
	EscrowPublicKey *rsa.PublicKey
}

func NewEscrowCrypto(publicKeyPEM []byte) (*EscrowCrypto, error) {
	if len(publicKeyPEM) == 0 {
		// Fallback for PoC: generate a dummy key if none provided
		// IN PRODUCTION: This MUST be loaded from KMS or per-tenant configuration
		priv, _ := rsa.GenerateKey(rand.Reader, 2048)
		return &EscrowCrypto{EscrowPublicKey: &priv.PublicKey}, nil
	}

	// Logic to decode PEM public key would go here
	return &EscrowCrypto{}, nil
}

func (c *EscrowCrypto) EncryptForKeyEscrow(plaintextKey string) (string, error) {
	if plaintextKey == "" {
		return "", fmt.Errorf("plaintext key is empty")
	}

	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		c.EscrowPublicKey,
		[]byte(plaintextKey),
		nil,
	)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
