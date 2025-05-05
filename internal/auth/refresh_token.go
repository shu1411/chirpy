package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// Read never returns an error,
// no error handling necessary
func MakeRefreshToken() string {
	data := make([]byte, 32)
	rand.Read(data)
	return hex.EncodeToString(data)
}
