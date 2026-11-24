// Package xorbreak recovers plaintexts from reused XOR keystreams.
//
// The attacks are Cryptopals set 1, challenges 2 and 3, plus a repeating-key
// break. They are here so you can see why a stream cipher nonce must never
// repeat. This is not a cipher implementation to ship.
package xorbreak

import "errors"

// XOR returns a XOR b. Both slices must be the same length.
func XOR(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, errors.New("xorbreak: length mismatch")
	}
	out := make([]byte, len(a))
	for i := range a {
		out[i] = a[i] ^ b[i]
	}
	return out, nil
}

// SingleByteXOR searches for the one-byte key that makes ct look most like English.
func SingleByteXOR(ct []byte) (key byte, plain []byte) {
	best := -1e9
	for k := 0; k < 256; k++ {
		cand := make([]byte, len(ct))
		for i, c := range ct {
			cand[i] = c ^ byte(k)
		}
		if s := scoreEnglish(cand); s > best {
			best = s
			key = byte(k)
			plain = cand
		}
	}
	return key, plain
}

// RepeatingKeyXOR recovers a short repeating key. It tries each length
// from 2 to 16, breaks each column as a single-byte XOR, and keeps the
// plaintext that looks most like English. Cryptopals challenge 6 uses
// Hamming distance to guess the length first. Either check works. The
// English score is the one that rejects a wrong length.
func RepeatingKeyXOR(ct []byte) (key, plain []byte) {
	limit := 16
	if limit > len(ct)/8 {
		limit = len(ct) / 8
	}
	if limit < 2 {
		limit = 2
	}
	best := -1e18
	for size := 2; size <= limit; size++ {
		k, p := breakSize(ct, size)
		if s := scoreEnglish(p); s > best {
			best = s
			key, plain = k, p
		}
	}
	return key, plain
}

func breakSize(ct []byte, size int) (key, plain []byte) {
	key = make([]byte, size)
	cols := make([][]byte, size)
	for i, b := range ct {
		cols[i%size] = append(cols[i%size], b)
	}
	plain = make([]byte, len(ct))
	for i, col := range cols {
		k, dec := SingleByteXOR(col)
		key[i] = k
		for n, b := range dec {
			plain[i+n*size] = b
		}
	}
	return key, plain
}

// DetectECB reports whether ct contains a repeated block. Cryptopals set 1
// challenge 8 is this test and nothing more.
func DetectECB(ct []byte, blockSize int) bool {
	if blockSize <= 0 || len(ct) < blockSize*2 || len(ct)%blockSize != 0 {
		return false
	}
	seen := make(map[string]struct{}, len(ct)/blockSize)
	for i := 0; i < len(ct); i += blockSize {
		block := string(ct[i : i+blockSize])
		if _, ok := seen[block]; ok {
			return true
		}
		seen[block] = struct{}{}
	}
	return false
}

func scoreEnglish(b []byte) float64 {
	var s float64
	for _, c := range b {
		switch {
		case c >= 'a' && c <= 'z':
			s += englishWeight[c-'a']
		case c >= 'A' && c <= 'Z':
			s += englishWeight[c-'A'] * 0.6
		case c == ' ':
			s += 0.18
		case c == '\'' || c == ',' || c == '.' || c == '!' || c == '?' || c == ';':
			s += 0.02
		case c == '\n' || c == '\t':
			s += 0.01
		case c >= 32 && c < 127:
			s -= 0.05
		default:
			s -= 1.5
		}
	}
	return s
}

var englishWeight = [26]float64{
	0.08167, 0.01492, 0.02782, 0.04253, 0.12702, 0.02228, 0.02015,
	0.06094, 0.06966, 0.00153, 0.00772, 0.04025, 0.02406, 0.06749,
	0.07507, 0.01929, 0.00095, 0.05987, 0.06327, 0.09056, 0.02758,
	0.00978, 0.02360, 0.00150, 0.01974, 0.00074,
}
