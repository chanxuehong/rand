package rand

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"time"
)

const (
	randSaltLen            = 45   // see Rand(), 8+2+45 == 55, Better performance for md5
	randSaltUpdateInterval = 3600 // seconds
)

var (
	randSalt = make([]byte, randSaltLen)

	randMutex              sync.Mutex
	randSaltLastUpdateTime int64  = -randSaltUpdateInterval
	randClockSequence      uint32 = Uint32()
)

// Rand returns 16-byte raw and not printable random bytes.
// You can use encoding/hex or encoding/base64 to print it.
func Rand() (ret [16]byte) {
	timeNow := time.Now()
	timeNowUnix := timeNow.Unix()

	randMutex.Lock() // Lock
	if timeNowUnix >= randSaltLastUpdateTime+randSaltUpdateInterval {
		randSaltLastUpdateTime = timeNowUnix
		Read(randSalt)

		randMutex.Unlock() // Unlock
		copy(ret[:], randSalt)
		return
	}
	randClockSequence++
	sequence := randClockSequence
	randMutex.Unlock() // Unlock

	var src [8 + 2 + randSaltLen]byte // 8+2+45 == 55
	timeNowUnixNano := timeNow.UnixNano()
	src[0] = byte(timeNowUnixNano >> 56)
	src[1] = byte(timeNowUnixNano >> 48)
	src[2] = byte(timeNowUnixNano >> 40)
	src[3] = byte(timeNowUnixNano >> 32)
	src[4] = byte(timeNowUnixNano >> 24)
	src[5] = byte(timeNowUnixNano >> 16)
	src[6] = byte(timeNowUnixNano >> 8)
	src[7] = byte(timeNowUnixNano)
	src[8] = byte(sequence >> 8)
	src[9] = byte(sequence)
	copy(src[10:], randSalt)

	ret = md5.Sum(src[:])
	return
}

// RandHex returns 32-byte hex-encoded bytes.
func RandHex() (ret []byte) {
	rd := Rand()
	ret = make([]byte, hex.EncodedLen(len(rd)))
	hex.Encode(ret, rd[:])
	return
}
