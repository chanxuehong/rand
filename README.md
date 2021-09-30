# golang rand package

[![GoPkg Widget]][GoPkg]

Package rand provides random number/bytes related functions.

The Rand type is similar to math/rand.Rand, except that it prefers to use crypto/rand to implement functions.

## NOTE

Because we use crypto/rand to implement functions, the performance is not as good as math/rand. If you care about
performance, please use it carefully.

[GoPkg]: https://pkg.go.dev/github.com/chanxuehong/rand

[GoPkg Widget]: https://pkg.go.dev/badge/github.com/chanxuehong/rand.svg
