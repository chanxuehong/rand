package rand

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	saltLen            = 45   // see New(), 6+4+45==55<56, Best performance for md5
	saltUpdateInterval = 3600 // seconds
)

var (
	_mutex                   sync.Mutex
	_salt                    unsafe.Pointer
	_saltLastUpdateTimestamp int64 = -saltUpdateInterval
	_sequence                      = Uint32()
)

func init() {
	salt := make([]byte, saltLen)
	Read(salt)
	storeSalt(salt)
}

func storeSalt(salt []byte) {
	atomic.StorePointer(&_salt, unsafe.Pointer(&salt))
}

func loadSalt() []byte {
	p := atomic.LoadPointer(&_salt)
	if p == nil {
		return nil
	}
	return *(*[]byte)(p)
}

// New returns 16-byte raw random bytes.
// It is not printable, you can use encoding/hex or encoding/base64 to print it.
func New() (rd [16]byte) {
	timeNow := time.Now()
	timeNowUnix := timeNow.Unix()

	if timeNowUnix >= atomic.LoadInt64(&_saltLastUpdateTimestamp)+saltUpdateInterval {
		_mutex.Lock() // Lock
		if timeNowUnix >= atomic.LoadInt64(&_saltLastUpdateTimestamp)+saltUpdateInterval {
			salt := make([]byte, saltLen)
			Read(salt)
			storeSalt(salt)
			atomic.StoreInt64(&_saltLastUpdateTimestamp, timeNowUnix)
			copy(rd[:], salt)
			_mutex.Unlock() // Unlock
			return
		}
		_mutex.Unlock() // Unlock
	}

	timeNowUnixNano := timeNow.UnixNano()
	sequence := atomic.AddUint32(&_sequence, 1)
	var src [6 + 4 + saltLen]byte // 6+4+45==55
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
	copy(src[10:], loadSalt())

	return md5.Sum(src[:])
}

// NewHex returns 32-byte hex-encoded bytes.
func NewHex() (rd []byte) {
	rdx := New()
	rd = make([]byte, hex.EncodedLen(len(rdx)))
	hex.Encode(rd, rdx[:])
	return
}
