package hkdf

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"testing"

	xhkdf "golang.org/x/crypto/hkdf"
)

func TestRFC5869Case1(t *testing.T) {
	// RFC 5869, appendix A, test case 1.
	ikm := bytes.Repeat([]byte{0x0b}, 22)
	salt := mustHex(t, "000102030405060708090a0b0c")
	info := mustHex(t, "f0f1f2f3f4f5f6f7f8f9")
	wantPRK := mustHex(t, "077709362c2e32df0ddc3f0dc47bba6390b6c73bb50f9c3122ec844ad7c2b3e5")
	wantOKM := mustHex(t, "3cb25f25faacd57a90434f64d0362f2a2d2d0a90cf1a5a4c5db02d56ecc4c5bf34007208d5b887185865")

	prk := Extract(sha256.New, ikm, salt)
	if !bytes.Equal(prk, wantPRK) {
		t.Fatalf("PRK = %x", prk)
	}
	okm, err := Expand(sha256.New, prk, info, 42)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(okm, wantOKM) {
		t.Fatalf("OKM = %x", okm)
	}

	if !bytes.Equal(prk, xhkdf.Extract(sha256.New, ikm, salt)) {
		t.Fatal("Extract disagreed with golang.org/x/crypto/hkdf")
	}
	if !bytes.Equal(okm, readX(t, ikm, salt, info, 42)) {
		t.Fatal("Expand disagreed with golang.org/x/crypto/hkdf")
	}
}

func TestNilSaltMatchesXCrypto(t *testing.T) {
	secret := []byte("ikm")
	info := []byte("tls13 info")
	got, err := Derive(sha256.New, secret, nil, info, 32)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, readX(t, secret, nil, info, 32)) {
		t.Fatal("nil salt disagreed with golang.org/x/crypto/hkdf")
	}
}

func TestEmptySaltMatchesXCrypto(t *testing.T) {
	secret := []byte("ikm")
	info := []byte("info")
	got, err := Derive(sha256.New, secret, []byte{}, info, 16)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, readX(t, secret, []byte{}, info, 16)) {
		t.Fatal("empty salt disagreed with golang.org/x/crypto/hkdf")
	}
	if !bytes.Equal(got, readX(t, secret, nil, info, 16)) {
		t.Fatal("empty salt disagreed with a nil salt")
	}
}

func readX(t *testing.T, secret, salt, info []byte, n int) []byte {
	t.Helper()
	buf := make([]byte, n)
	if _, err := io.ReadFull(xhkdf.New(sha256.New, secret, salt, info), buf); err != nil {
		t.Fatal(err)
	}
	return buf
}

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
