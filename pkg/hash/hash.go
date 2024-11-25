package jhash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(src []byte) string {
	hash := sha256.New()
	hash.Write(src)
	return hex.EncodeToString(hash.Sum(nil))
}

func Compare(hashed, new []byte) bool {
	return string(hashed) == Hash(new)
}
