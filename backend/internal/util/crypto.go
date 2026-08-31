package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/google/uuid"
	"golang.org/x/crypto/curve25519"
)

type RealityKeypair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

// GenerateRealityKeypair generates X25519 keypair for VLESS Reality
func GenerateRealityKeypair() (*RealityKeypair, error) {
	var privateKey [32]byte
	if _, err := io.ReadFull(rand.Reader, privateKey[:]); err != nil {
		return nil, err
	}

	var publicKey [32]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	return &RealityKeypair{
		PrivateKey: base64.RawURLEncoding.EncodeToString(privateKey[:]),
		PublicKey:  base64.RawURLEncoding.EncodeToString(publicKey[:]),
	}, nil
}

// GenerateShortID generates random hex string for Reality short_id (e.g. 8 chars)
func GenerateShortID() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "01234567"
	}
	return hex.EncodeToString(bytes)
}

// GenerateUUID generates standard UUID v4
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateRandomPassword generates a secure random alphanumeric string
func GenerateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return uuid.New().String()[:length]
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}
