// Package caesar breaks classical substitution ciphers.
//
// This is practice for frequency analysis. Caesar and Vigenère are not
// encryption. Do not use them for anything but this exercise.
package caesar

import "math"

// English letter frequencies, A through Z.
var english = [26]float64{
	0.08167, 0.01492, 0.02782, 0.04253, 0.12702, 0.02228, 0.02015,
	0.06094, 0.06966, 0.00153, 0.00772, 0.04025, 0.02406, 0.06749,
	0.07507, 0.01929, 0.00095, 0.05987, 0.06327, 0.09056, 0.02758,
	0.00978, 0.02360, 0.00150, 0.01974, 0.00074,
}

// Encrypt shifts English letters by shift positions. Other bytes are copied.
func Encrypt(plain string, shift int) string {
	return crypt(plain, shift)
}

// Decrypt reverses Encrypt.
func Decrypt(cipher string, shift int) string {
	return crypt(cipher, -shift)
}

// Break recovers a Caesar shift by chi-squared fit against English letter
// frequencies. It returns the plaintext and the shift that was used to encrypt.
func Break(cipher string) (plain string, shift int) {
	best := math.Inf(1)
	bestShift := 0
	for s := 0; s < 26; s++ {
		candidate := Decrypt(cipher, s)
		score := chiSquare(candidate)
		if score < best {
			best = score
			bestShift = s
			plain = candidate
		}
	}
	return plain, bestShift
}

// EncryptVigenere encrypts with a repeating alphabetic key. The key is
// interpreted case-insensitively. A non-letter in the key is skipped.
func EncryptVigenere(plain, key string) string {
	return vigenere(plain, key, 1)
}

// DecryptVigenere reverses EncryptVigenere.
func DecryptVigenere(cipher, key string) string {
	return vigenere(cipher, key, -1)
}

// BreakVigenere recovers plaintext from an English Vigenère ciphertext.
// The key length is the shortest length whose columns look like English.
func BreakVigenere(cipher string) (plain, key string) {
	letters := letterBytes(cipher)
	if len(letters) < 2 {
		return cipher, ""
	}
	maxK := 16
	if maxK > len(letters)/4 {
		maxK = len(letters) / 4
	}
	if maxK < 1 {
		maxK = 1
	}
	k := 1
	for n := 1; n <= maxK; n++ {
		if avgIOC(letters, n) >= 0.055 {
			k = n
			break
		}
	}
	if avgIOC(letters, k) < 0.055 {
		best := -1.0
		for n := 1; n <= maxK; n++ {
			score := avgIOC(letters, n)
			if score > best {
				best = score
				k = n
			}
		}
	}

	shifts := make([]int, k)
	cols := make([][]byte, k)
	for i, b := range letters {
		cols[i%k] = append(cols[i%k], b)
	}
	keyLetters := make([]byte, k)
	for i, col := range cols {
		_, shift := Break(string(col))
		shifts[i] = shift
		keyLetters[i] = 'A' + byte(shift)
	}
	return DecryptVigenere(cipher, string(keyLetters)), string(keyLetters)
}

func crypt(s string, shift int) string {
	shift = ((shift % 26) + 26) % 26
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch {
		case b >= 'a' && b <= 'z':
			out[i] = 'a' + (b-'a'+byte(shift))%26
		case b >= 'A' && b <= 'Z':
			out[i] = 'A' + (b-'A'+byte(shift))%26
		default:
			out[i] = b
		}
	}
	return string(out)
}

func vigenere(s, key string, dir int) string {
	shifts := keyShifts(key)
	if len(shifts) == 0 {
		return s
	}
	out := make([]byte, len(s))
	ki := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		if !isLetter(b) {
			out[i] = b
			continue
		}
		shift := shifts[ki%len(shifts)] * dir
		ki++
		out[i] = crypt(string([]byte{b}), shift)[0]
	}
	return string(out)
}

func keyShifts(key string) []int {
	var shifts []int
	for i := 0; i < len(key); i++ {
		b := key[i]
		switch {
		case b >= 'a' && b <= 'z':
			shifts = append(shifts, int(b-'a'))
		case b >= 'A' && b <= 'Z':
			shifts = append(shifts, int(b-'A'))
		}
	}
	return shifts
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func letterBytes(s string) []byte {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if isLetter(s[i]) {
			out = append(out, s[i])
		}
	}
	return out
}

func chiSquare(s string) float64 {
	var counts [26]int
	n := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch {
		case b >= 'a' && b <= 'z':
			counts[b-'a']++
			n++
		case b >= 'A' && b <= 'Z':
			counts[b-'A']++
			n++
		}
	}
	if n == 0 {
		return math.Inf(1)
	}
	var chi float64
	fn := float64(n)
	for i, exp := range english {
		expected := exp * fn
		diff := float64(counts[i]) - expected
		chi += diff * diff / expected
	}
	return chi
}

func avgIOC(letters []byte, k int) float64 {
	if k < 1 {
		return 0
	}
	var sum float64
	var used int
	for i := 0; i < k; i++ {
		var col []byte
		for j := i; j < len(letters); j += k {
			col = append(col, letters[j])
		}
		if len(col) < 2 {
			continue
		}
		sum += ioc(col)
		used++
	}
	if used == 0 {
		return 0
	}
	return sum / float64(used)
}

func ioc(col []byte) float64 {
	var counts [26]int
	n := 0
	for _, b := range col {
		switch {
		case b >= 'a' && b <= 'z':
			counts[b-'a']++
			n++
		case b >= 'A' && b <= 'Z':
			counts[b-'A']++
			n++
		}
	}
	if n < 2 {
		return 0
	}
	var num float64
	for _, c := range counts {
		num += float64(c * (c - 1))
	}
	return num / float64(n*(n-1))
}
