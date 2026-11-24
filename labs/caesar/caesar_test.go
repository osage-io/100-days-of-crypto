package caesar

import "testing"

const prose = "" +
	"Security engineers already know how to threat model a system. What is new " +
	"here is the habit of counting letters. A substitution cipher keeps every " +
	"pattern of the language and only renames the symbols. English still has " +
	"too many e letters, too many t letters, and too many short words. Chi squared " +
	"against those frequencies is enough to undo a Caesar shift. A Vigenere key " +
	"is the same attack once you split the ciphertext into columns that were " +
	"shifted by a single letter. None of this is modern cryptography. It is the " +
	"warmup that makes the later failure modes feel familiar: the construction " +
	"leaked a statistic, and the statistic was enough."

func TestCaesarRoundTripAndBreak(t *testing.T) {
	const shift = 11
	ct := Encrypt(prose, shift)
	if ct == prose {
		t.Fatal("encryption did not change the text")
	}
	if got := Decrypt(ct, shift); got != prose {
		t.Fatalf("decrypt mismatch")
	}
	plain, found := Break(ct)
	if found != shift {
		t.Fatalf("shift = %d, want %d", found, shift)
	}
	if plain != prose {
		t.Fatalf("broken plaintext mismatch")
	}
}

func TestCaesarLeavesNonLetters(t *testing.T) {
	got := Encrypt("Go 1.26!", 1)
	if got != "Hp 1.26!" {
		t.Fatalf("got %q", got)
	}
}

func TestVigenereRoundTripAndBreak(t *testing.T) {
	const key = "CRYPTO"
	ct := EncryptVigenere(prose, key)
	if DecryptVigenere(ct, key) != prose {
		t.Fatal("decrypt mismatch")
	}
	plain, found := BreakVigenere(ct)
	if plain != prose {
		t.Fatalf("plaintext mismatch\n got: %s\nwant: %s", plain, prose)
	}
	if found != key {
		t.Fatalf("key = %q, want %q", found, key)
	}
}
