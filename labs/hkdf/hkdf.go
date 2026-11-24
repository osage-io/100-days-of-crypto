// Package hkdf implements HKDF as specified in RFC 5869, using crypto/hmac.
//
// The tests check it against golang.org/x/crypto/hkdf. Use that package, or
// crypto/hkdf in the standard library, in real code. This file exists so the
// extract and expand steps are visible.
package hkdf

import (
	"crypto/hmac"
	"errors"
	"hash"
)

// Extract returns HMAC(salt, secret). A nil or empty salt becomes HashLen
// zero bytes, which is what RFC 5869 means by "salt not provided".
//
// golang.org/x/crypto/hkdf rewrites only a nil salt. An empty slice is
// handed to HMAC as an empty key, and HMAC zero-pads that key out to the
// block size, so the digest matches anyway.
func Extract(h func() hash.Hash, secret, salt []byte) []byte {
	if len(salt) == 0 {
		salt = make([]byte, h().Size())
	}
	mac := hmac.New(h, salt)
	mac.Write(secret)
	return mac.Sum(nil)
}

// Expand returns the first length bytes of the HKDF expand step.
func Expand(h func() hash.Hash, pseudorandomKey, info []byte, length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("hkdf: length must be positive")
	}
	hashLen := h().Size()
	if length > 255*hashLen {
		return nil, errors.New("hkdf: length exceeds 255 * HashLen")
	}
	n := (length + hashLen - 1) / hashLen
	var prev []byte
	out := make([]byte, 0, n*hashLen)
	for i := 1; i <= n; i++ {
		mac := hmac.New(h, pseudorandomKey)
		mac.Write(prev)
		mac.Write(info)
		mac.Write([]byte{byte(i)})
		prev = mac.Sum(nil)
		out = append(out, prev...)
	}
	return out[:length], nil
}

// Derive is Extract followed by Expand.
func Derive(h func() hash.Hash, secret, salt, info []byte, length int) ([]byte, error) {
	return Expand(h, Extract(h, secret, salt), info, length)
}
