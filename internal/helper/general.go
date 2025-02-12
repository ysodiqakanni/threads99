package helper

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

func GenerateUsernameFromEmail(email string) string {
	// Remove everything after @ and any dots
	username := strings.Split(email, "@")[0]
	username = strings.ReplaceAll(username, ".", "")

	// Take first 6 characters (or less if email is shorter)
	length := min(6, len(username))
	prefix := username[:length]

	// Generate 2 random bytes (4 hex characters)
	randomBytes := make([]byte, 2)
	rand.Read(randomBytes)
	suffix := hex.EncodeToString(randomBytes)

	// Combine prefix and suffix, ensure total length <= 10
	result := prefix + suffix
	if len(result) > 10 {
		result = result[:10]
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
