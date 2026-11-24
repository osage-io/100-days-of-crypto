// Package x25519 performs an X25519 key agreement with crypto/ecdh.
//
// The shared secret is not a key. Run it through HKDF, with a label, before
// you encrypt anything. This package only does the Diffie-Hellman step.
package x25519

import (
	"crypto/ecdh"
	"crypto/rand"
	"io"
)

// Key is one side of an X25519 exchange.
type Key struct {
	priv *ecdh.PrivateKey
}

// Generate draws a fresh X25519 private key.
func Generate(random io.Reader) (*Key, error) {
	if random == nil {
		random = rand.Reader
	}
	priv, err := ecdh.X25519().GenerateKey(random)
	if err != nil {
		return nil, err
	}
	return &Key{priv: priv}, nil
}

// Public returns the 32-byte public key.
func (k *Key) Public() []byte {
	return k.priv.PublicKey().Bytes()
}

// Shared computes the raw X25519 shared secret with a peer public key.
func (k *Key) Shared(peerPublic []byte) ([]byte, error) {
	pub, err := ecdh.X25519().NewPublicKey(peerPublic)
	if err != nil {
		return nil, err
	}
	return k.priv.ECDH(pub)
}
