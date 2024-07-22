package cloudfiles

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func ComparePassword(enteredPassword, storedHash string) bool {
	return HashPassword(enteredPassword) == storedHash
}
