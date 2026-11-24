// Package toyrsa is textbook RSA with no padding.
//
// NOT FOR PRODUCTION.
// The primes are whatever size you ask for, encryption is deterministic,
// and there is no OAEP or PSS. It exists so the arithmetic is visible.
// Real code uses crypto/rsa, and usually does not encrypt a message with
// RSA at all. Encrypt a key with a KEM, then use an AEAD.
package toyrsa

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

// PrivateKey is a textbook RSA key. P and Q are kept so you can see CRT
// inputs. Nothing here is padded.
type PrivateKey struct {
	N, E, D *big.Int
	P, Q    *big.Int
}

// Generate builds a key with two primes of the given bit length and e = 65537.
// bits is the size of each prime, not of N. 256 is plenty for the tests and
// still useless against a real attacker.
func Generate(random io.Reader, bits int) (*PrivateKey, error) {
	if bits < 32 {
		return nil, errors.New("toyrsa: refuse primes under 32 bits; the result would be even easier to mistake for a key")
	}
	if random == nil {
		random = rand.Reader
	}
	e := big.NewInt(65537)
	one := big.NewInt(1)
	for {
		p, err := rand.Prime(random, bits)
		if err != nil {
			return nil, err
		}
		q, err := rand.Prime(random, bits)
		if err != nil {
			return nil, err
		}
		if p.Cmp(q) == 0 {
			continue
		}
		p1 := new(big.Int).Sub(p, one)
		q1 := new(big.Int).Sub(q, one)
		phi := new(big.Int).Mul(p1, q1)
		if new(big.Int).GCD(nil, nil, e, phi).Cmp(one) != 0 {
			continue
		}
		d := new(big.Int).ModInverse(e, phi)
		if d == nil {
			continue
		}
		return &PrivateKey{
			N: new(big.Int).Mul(p, q),
			E: new(big.Int).Set(e),
			D: d,
			P: p,
			Q: q,
		}, nil
	}
}

// Encrypt computes m^e mod N. m must satisfy 0 <= m < N.
// The same m always produces the same ciphertext.
func (k *PrivateKey) Encrypt(m *big.Int) (*big.Int, error) {
	if m == nil || m.Sign() < 0 || m.Cmp(k.N) >= 0 {
		return nil, errors.New("toyrsa: message representative out of range")
	}
	return new(big.Int).Exp(m, k.E, k.N), nil
}

// Decrypt computes c^d mod N.
func (k *PrivateKey) Decrypt(c *big.Int) *big.Int {
	return new(big.Int).Exp(c, k.D, k.N)
}

// FromBytes maps a short byte string into a message representative.
func (k *PrivateKey) FromBytes(b []byte) (*big.Int, error) {
	m := new(big.Int).SetBytes(b)
	if len(b) == 0 || m.Cmp(k.N) >= 0 {
		return nil, errors.New("toyrsa: message does not fit under N")
	}
	return m, nil
}
