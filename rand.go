package rand

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"
)

const (
	randSaltLen            = 45   // see Rand(), 6+4+45==55<56, Best performance for md5
	randSaltUpdateInterval = 3600 // seconds
)

var (
	randSalt     = make([]byte, randSaltLen)
	randSequence = Uint32()

	mutex                  sync.Mutex
	randSaltLastUpdateTime int64 = -randSaltUpdateInterval
)

// Rand returns 16-byte raw random bytes.
// It is not printable, you can use encoding/hex or encoding/base64 to print it.
func Rand() (rd [16]byte) {
	timeNow := time.Now()
	timeNowUnix := timeNow.Unix()

	if timeNowUnix >= atomic.LoadInt64(&randSaltLastUpdateTime)+randSaltUpdateInterval {
		mutex.Lock() // Lock
		if timeNowUnix >= randSaltLastUpdateTime+randSaltUpdateInterval {
			randSaltLastUpdateTime = timeNowUnix
			mutex.Unlock() // Unlock

			Read(randSalt)
			copy(rd[:], randSalt)
			return
		}
		mutex.Unlock() // Unlock
	}
	sequence := atomic.AddUint32(&randSequence, 1)

	var src [6 + 4 + randSaltLen]byte // 6+4+45==55
	timeNowUnixNano := timeNow.UnixNano()
	src[0] = byte(timeNowUnixNano >> 40)
	src[1] = byte(timeNowUnixNano >> 32)
	src[2] = byte(timeNowUnixNano >> 24)
	src[3] = byte(timeNowUnixNano >> 16)
	src[4] = byte(timeNowUnixNano >> 8)
	src[5] = byte(timeNowUnixNano)
	src[6] = byte(sequence >> 24)
	src[7] = byte(sequence >> 16)
	src[8] = byte(sequence >> 8)
	src[9] = byte(sequence)
	copy(src[10:], randSalt)

	return md5.Sum(src[:])
}

// RandHex returns 32-byte hex-encoded bytes.
func RandHex() (rd []byte) {
	rdx := Rand()
	rd = make([]byte, hex.EncodedLen(len(rdx)))
	hex.Encode(rd, rdx[:])
	return
}
