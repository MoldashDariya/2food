package generate

import (
	"crypto/rand"
	"math/big"
)

const numericChars = "0123456789"

// Code generates a random numeric code of the given length.
func Code(length int) string {
	code := make([]byte, length)
	max := big.NewInt(int64(len(numericChars)))

	for i := range code {
		n, _ := rand.Int(rand.Reader, max)
		code[i] = numericChars[n.Int64()]
	}

	return string(code)
}
