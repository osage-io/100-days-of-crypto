package xorbreak

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestFixedXOR(t *testing.T) {
	// Cryptopals set 1 challenge 2.
	a := mustHex(t, "1c0111001f010100061a024b53535009181c")
	b := mustHex(t, "686974207468652062756c6c277320657965")
	got, err := XOR(a, b)
	if err != nil {
		t.Fatal(err)
	}
	want := mustHex(t, "746865206b696420646f6e277420706c6179")
	if !bytes.Equal(got, want) {
		t.Fatalf("got %x", got)
	}
}

func TestSingleByteCryptopals(t *testing.T) {
	// Cryptopals set 1 challenge 3. Try it before reading this expected plaintext.
	ct := mustHex(t, "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")
	key, plain := SingleByteXOR(ct)
	if key != 0x58 {
		t.Fatalf("key = %#x", key)
	}
	if string(plain) != "Cooking MC's like a pound of bacon" {
		t.Fatalf("plain = %q", plain)
	}
}

func TestRepeatingKey(t *testing.T) {
	plain := []byte("A repeating XOR key is a stream cipher with a short period. " +
		"Hamming distance between blocks finds that period when the plaintext is English. " +
		"After the columns split, each column is the single-byte problem again. " +
		"Use a real nonce, or do not build a stream cipher at all.")
	key := []byte("CRYPTO")
	ct := repeatXOR(plain, key)
	gotKey, got := RepeatingKeyXOR(ct)
	if !bytes.Equal(gotKey, key) {
		t.Fatalf("key = %q, want %q", gotKey, key)
	}
	if !bytes.Equal(got, plain) {
		t.Fatalf("plain = %q", got)
	}
}

func TestDetectECB(t *testing.T) {
	block := bytes.Repeat([]byte("YELLOW SUBMARINE"), 4)
	if !DetectECB(block, 16) {
		t.Fatal("expected a repeated block")
	}
	mixed := []byte("0123456789abcdef00112233445566778899aabbccddeeff")
	if DetectECB(mixed, 16) {
		t.Fatal("distinct blocks were flagged as ECB")
	}
}

func repeatXOR(plain, key []byte) []byte {
	out := make([]byte, len(plain))
	for i, b := range plain {
		out[i] = b ^ key[i%len(key)]
	}
	return out
}

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
