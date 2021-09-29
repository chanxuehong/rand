package rand

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math"
	mathrand "math/rand"
	"time"
)

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

// Read reads/generates len(p) random bytes and writes them into p.
func Read(p []byte) {
	if len(p) <= 0 {
		return
	}
	if _, err := cryptorand.Read(p); err != nil {
		mathrand.Read(p)
	}
}

// Int31 returns a non-negative random 31-bit integer as an int32.
func Int31() int32 {
	return int32(Uint32() & math.MaxInt32)
}

// Int63 returns a non-negative random 63-bit integer as an int64.
func Int63() int64 {
	return int64(Uint64() & math.MaxInt64)
}

// Uint32 returns a random 32-bit integer as an uint32.
func Uint32() uint32 {
	var x [4]byte
	Read(x[:])
	return binary.BigEndian.Uint32(x[:])
}

// Uint64 returns a random 64-bit integer as an uint64.
func Uint64() uint64 {
	var x [8]byte
	Read(x[:])
	return binary.BigEndian.Uint64(x[:])
}
