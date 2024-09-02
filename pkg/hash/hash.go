package jhash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(src string) string {
	hash := sha256.New()
	hash.Write([]byte(src))
	return hex.EncodeToString(hash.Sum(nil))
}

func Compare(hashed, new string) bool {
	return hashed == Hash(new)
}
