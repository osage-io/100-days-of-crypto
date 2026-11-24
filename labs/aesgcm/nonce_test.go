package aesgcm

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestNonceReuseRecoversPlaintext(t *testing.T) {
	key := make([]byte, 32)
	nonce := make([]byte, 12)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	if _, err := rand.Read(nonce); err != nil {
		t.Fatal(err)
	}
	known := []byte("session token = root")
	secret := []byte("session token = admin-key")
	aad := []byte("header")

	sealedKnown, err := Seal(key, nonce, known, aad)
	if err != nil {
		t.Fatal(err)
	}
	sealedSecret, err := Seal(key, nonce, secret, aad)
	if err != nil {
		t.Fatal(err)
	}

	// Both tags still verify. GCM does not notice the reuse.
	opened, err := Open(key, nonce, sealedSecret, aad)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(opened, secret) {
		t.Fatalf("open got %q", opened)
	}

	got, err := RecoverPlaintext(sealedKnown, sealedSecret, known)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, secret[:len(known)]) {
		t.Fatalf("recovered %q", got)
	}
}

func TestFreshNonceDoesNotOpenUnderTheOther(t *testing.T) {
	key := make([]byte, 16)
	n1 := make([]byte, 12)
	n2 := make([]byte, 12)
	if _, err := rand.Read(key); err != nil {
		t.Fatal(err)
	}
	if _, err := rand.Read(n1); err != nil {
		t.Fatal(err)
	}
	if _, err := rand.Read(n2); err != nil {
		t.Fatal(err)
	}
	sealed, err := Seal(key, n1, []byte("one"), nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Open(key, n2, sealed, nil); err == nil {
		t.Fatal("a ciphertext opened under a different nonce")
	}
}
