package x25519

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestBothSidesAgree(t *testing.T) {
	alice, err := Generate(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	bob, err := Generate(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ab, err := alice.Shared(bob.Public())
	if err != nil {
		t.Fatal(err)
	}
	ba, err := bob.Shared(alice.Public())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ab, ba) {
		t.Fatal("shared secrets differ")
	}
	if len(ab) != 32 {
		t.Fatalf("shared secret length %d", len(ab))
	}
	if bytes.Equal(ab, alice.Public()) || bytes.Equal(ab, bob.Public()) {
		t.Fatal("shared secret matched a public key")
	}

	carol, err := Generate(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	ac, err := alice.Shared(carol.Public())
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(ab, ac) {
		t.Fatal("a third party derived the same secret")
	}
}

func TestRejectsBadPublic(t *testing.T) {
	alice, err := Generate(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := alice.Shared(make([]byte, 32)); err == nil {
		t.Fatal("all-zero public key was accepted")
	}
	if _, err := alice.Shared([]byte{1, 2, 3}); err == nil {
		t.Fatal("short public key was accepted")
	}
}
