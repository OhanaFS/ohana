package util

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomHex returns a random hex string of the given number of bytes.
// Note that n is the number of bytes, not the number of hex characters. To get
// a string of M characters, pass n = M/2.
func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
