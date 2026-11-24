package toyrsa

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestRoundTripIsDeterministicAndUnpadded(t *testing.T) {
	k, err := Generate(rand.Reader, 256)
	if err != nil {
		t.Fatal(err)
	}
	m, err := k.FromBytes([]byte("not for production"))
	if err != nil {
		t.Fatal(err)
	}
	c1, err := k.Encrypt(m)
	if err != nil {
		t.Fatal(err)
	}
	c2, err := k.Encrypt(m)
	if err != nil {
		t.Fatal(err)
	}
	if c1.Cmp(c2) != 0 {
		t.Fatal("textbook RSA was not deterministic")
	}
	if k.Decrypt(c1).Cmp(m) != 0 {
		t.Fatal("decrypt mismatch")
	}
}

func TestRejectsMessageOutsideRange(t *testing.T) {
	k, err := Generate(rand.Reader, 64)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := k.Encrypt(k.N); err == nil {
		t.Fatal("encrypt accepted m == N")
	}
	if _, err := k.Encrypt(new(big.Int).Add(k.N, big.NewInt(1))); err == nil {
		t.Fatal("encrypt accepted m > N")
	}
}
