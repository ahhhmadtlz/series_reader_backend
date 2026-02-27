package service

import (
	"crypto/sha256"
	"fmt"
)

// hashToken returns a SHA-256 hex string of the token.
// We store the hash, never the raw token.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", sum)
}