// Package mlkem uses the standard library ML-KEM-768 implementation.
//
// crypto/mlkem has been in the standard library since Go 1.24. Most programs
// should not call it directly. TLS already negotiates a hybrid. This package
// is the smallest loop that shows encapsulation, decapsulation, and the
// implicit rejection FIPS 203 requires.
package mlkem

import "crypto/mlkem"

// Exchange768 generates a decapsulation key, encapsulates to its public key
// from a second handle, and decapsulates. Both shared secrets match when
// the ciphertext is intact.
func Exchange768() (sharedSender, sharedReceiver, ciphertext []byte, err error) {
	dk, err := mlkem.GenerateKey768()
	if err != nil {
		return nil, nil, nil, err
	}
	ek, err := mlkem.NewEncapsulationKey768(dk.EncapsulationKey().Bytes())
	if err != nil {
		return nil, nil, nil, err
	}
	sharedSender, ciphertext = ek.Encapsulate()
	sharedReceiver, err = dk.Decapsulate(ciphertext)
	if err != nil {
		return nil, nil, nil, err
	}
	return sharedSender, sharedReceiver, ciphertext, nil
}

// Decapsulate768 decapsulates ciphertext with seed, which is the 64-byte
// d || z seed from DecapsulationKey768.Bytes.
func Decapsulate768(seed, ciphertext []byte) ([]byte, error) {
	dk, err := mlkem.NewDecapsulationKey768(seed)
	if err != nil {
		return nil, err
	}
	return dk.Decapsulate(ciphertext)
}
