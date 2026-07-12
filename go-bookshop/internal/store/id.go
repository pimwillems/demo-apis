package store

import (
	"crypto/rand"
	"math/big"
)

const idAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func newShortID(prefix string) string {
	const length = 10
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(idAlphabet))))
		if err != nil {
			// crypto/rand should never fail; fall back to '0' on the
			// off chance it does so callers always get a valid id.
			b[i] = idAlphabet[0]
			continue
		}
		b[i] = idAlphabet[n.Int64()]
	}
	return prefix + string(b)
}

func newOrderID() string    { return newShortID("ord-") }
func newCustomerID() string { return newShortID("cust-") }
