package token

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

// Generate returns a cryptographically secure opaque token for a client cookie.
func Generate() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate random token: %w", err)
	}

	return hex.EncodeToString(buf), nil
}

// Hash returns the SHA-256 representation stored in the database.
func Hash(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// Verify checks a raw client token against the stored SHA-256 hash.
func Verify(raw, storedHash string) bool {
	return subtle.ConstantTimeCompare(
		[]byte(Hash(raw)),
		[]byte(storedHash),
	) == 1
}
