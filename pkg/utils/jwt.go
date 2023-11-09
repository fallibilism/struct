package utils

import (
	"crypto/sha1"
	"encoding/hex"

)

// to hash id using *sha1*
func GenHash(id string) (hash string) {
	hasher := sha1.New()
	hasher.Write([]byte(id))
	hash = hex.EncodeToString(hasher.Sum(nil))

	return hash
}
