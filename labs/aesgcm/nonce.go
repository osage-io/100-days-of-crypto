// Package aesgcm shows what happens when an AES-GCM nonce is reused.
//
// DEMO ONLY. Nonce reuse breaks confidentiality of GCM immediately: the
// two ciphertexts XOR to the two plaintexts. The authentication tags still
// verify, so the API will not save you. Do not copy this into a design.
// A production caller uses crypto/cipher AEAD with a nonce that never repeats
// under the same key.
package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

// Seal encrypts plaintext with AES-GCM. The nonce must be 12 bytes.
// The same (key, nonce) pair must not be used twice. This function will
// not stop you. That is the bug the tests demonstrate.
func Seal(key, nonce, plaintext, additionalData []byte) ([]byte, error) {
	aead, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	if len(nonce) != aead.NonceSize() {
		return nil, fmt.Errorf("aesgcm: nonce must be %d bytes", aead.NonceSize())
	}
	return aead.Seal(nil, nonce, plaintext, additionalData), nil
}

// Open decrypts a Seal result. A reused nonce still opens cleanly.
func Open(key, nonce, sealed, additionalData []byte) ([]byte, error) {
	aead, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	if len(nonce) != aead.NonceSize() {
		return nil, fmt.Errorf("aesgcm: nonce must be %d bytes", aead.NonceSize())
	}
	return aead.Open(nil, nonce, sealed, additionalData)
}

// RecoverPlaintext uses one known plaintext and two sealed messages that
// shared a nonce. It returns the other plaintext, truncated to the length
// of known. Tags are ignored: the keystream is in the ciphertext bytes.
func RecoverPlaintext(sealed, otherSealed, known []byte) ([]byte, error) {
	ct1, err := body(sealed)
	if err != nil {
		return nil, err
	}
	ct2, err := body(otherSealed)
	if err != nil {
		return nil, err
	}
	if len(known) > len(ct1) || len(known) > len(ct2) {
		return nil, errors.New("aesgcm: known plaintext longer than ciphertext")
	}
	out := make([]byte, len(known))
	for i := range known {
		out[i] = ct1[i] ^ ct2[i] ^ known[i]
	}
	return out, nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func body(sealed []byte) ([]byte, error) {
	const tag = 16
	if len(sealed) < tag {
		return nil, errors.New("aesgcm: ciphertext shorter than a tag")
	}
	return sealed[:len(sealed)-tag], nil
}
