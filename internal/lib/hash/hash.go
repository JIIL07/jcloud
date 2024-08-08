package jhash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(src string) string {
	hasher := sha256.New()
	hasher.Write([]byte(src))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Compare(hashed, new string) bool {
	return hashed == Hash(new)
}
