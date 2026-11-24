package mlkem

import (
	"bytes"
	"crypto/mlkem"
	"testing"
)

func TestExchange768(t *testing.T) {
	send, recv, ct, err := Exchange768()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(send, recv) {
		t.Fatal("shared secrets differ")
	}
	if len(send) != mlkem.SharedKeySize {
		t.Fatalf("shared key size %d", len(send))
	}
	if len(ct) != mlkem.CiphertextSize768 {
		t.Fatalf("ciphertext size %d", len(ct))
	}
}

func TestImplicitRejection(t *testing.T) {
	dk, err := mlkem.GenerateKey768()
	if err != nil {
		t.Fatal(err)
	}
	ek, err := mlkem.NewEncapsulationKey768(dk.EncapsulationKey().Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if len(ek.Bytes()) != mlkem.EncapsulationKeySize768 {
		t.Fatalf("encapsulation key size %d", len(ek.Bytes()))
	}
	shared, ct := ek.Encapsulate()
	damaged := bytes.Clone(ct)
	damaged[0] ^= 0x01

	// FIPS 203: a ciphertext of the right length never returns an error.
	// A damaged one decapsulates to a different shared key.
	got, err := dk.Decapsulate(damaged)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(got, shared) {
		t.Fatal("damaged ciphertext reproduced the sender's shared key")
	}

	if _, err := dk.Decapsulate(ct[:len(ct)-1]); err == nil {
		t.Fatal("a short ciphertext was accepted")
	}
}

func TestSeedRoundTrip(t *testing.T) {
	dk, err := mlkem.GenerateKey768()
	if err != nil {
		t.Fatal(err)
	}
	_, ct := dk.EncapsulationKey().Encapsulate()
	shared, err := dk.Decapsulate(ct)
	if err != nil {
		t.Fatal(err)
	}
	again, err := Decapsulate768(dk.Bytes(), ct)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(shared, again) {
		t.Fatal("seed did not reproduce the decapsulation key")
	}
	if len(dk.Bytes()) != mlkem.SeedSize {
		t.Fatalf("seed size %d", len(dk.Bytes()))
	}
}
