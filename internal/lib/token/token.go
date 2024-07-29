package jtoken

import (
	"crypto/rand"
	"encoding/hex"
)

var GlobalToken string

func Generate(length int) string {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	GlobalToken = hex.EncodeToString(bytes)
	return GlobalToken[:length]
}
