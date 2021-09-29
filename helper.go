package rand

import "encoding/hex"

// New returns 16 random bytes.
//
// The result is not printable, you can use encoding/hex or encoding/base64 to print it.
func New() (r [16]byte) {
	Read(r[:])
	return
}

// NewHex returns 32 random hex-encoded bytes.
func NewHex() []byte {
	buf := New()
	result := make([]byte, hex.EncodedLen(len(buf)))
	hex.Encode(result, buf[:])
	return result
}
