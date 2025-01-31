package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Encode(value []byte, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))

	if _, err := h.Write(value); err != nil {
		return "", fmt.Errorf("hash.encode: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
