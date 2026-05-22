package tenant

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

// KMSClient defines the interface for decrypting sensitive tenant metadata
type KMSClient interface {
	Decrypt(ciphertext string) (string, error)
}

// SymmetricKMS is a local implementation of KMSClient using AES-GCM
type SymmetricKMS struct {
	masterKey []byte
}

func NewSymmetricKMS(key string) (*SymmetricKMS, error) {
	k, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid master key encoding: %w", err)
	}
	if len(k) != 32 {
		return nil, errors.New("master key must be 32 bytes (base64 encoded)")
	}
	return &SymmetricKMS{masterKey: k}, nil
}

func (s *SymmetricKMS) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.masterKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, encryptedMessage := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
