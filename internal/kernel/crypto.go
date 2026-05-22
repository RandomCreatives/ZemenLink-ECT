package kernel

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
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

	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}

	return &EscrowCrypto{EscrowPublicKey: rsaPub}, nil
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
