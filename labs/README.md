# Labs

Small Go programs. `go test ./...` from the repo root runs all of them. Read the package comment before you copy anything.

| Lab | Day | What it shows |
| --- | --- | --- |
| [caesar](caesar/caesar.go) | 2 | Frequency analysis breaks Caesar and Vigenère. Not a cipher to use. |
| [xorbreak](xorbreak/xor.go) | 72 | Cryptopals set 1: fixed XOR, single-byte XOR, repeating-key XOR, ECB detection. |
| [aesgcm](aesgcm/nonce.go) | 17 | AES-GCM nonce reuse. The tags still verify. Demo only. |
| [hkdf](hkdf/hkdf.go) | 29 | HKDF extract and expand, checked against `golang.org/x/crypto/hkdf` and RFC 5869. |
| [toyrsa](toyrsa/rsa.go) | 34 | Textbook RSA. No padding, small primes. Not for production. |
| [x25519](x25519/exchange.go) | 42 | X25519 with `crypto/ecdh`. The shared secret still needs a KDF. |
| [tls13](tls13/handshake.go) | 60 | A local TLS 1.3 handshake, the ClientHello bytes, and an NSS key log. |
| [mlkem](mlkem/exchange.go) | 84 | ML-KEM-768 from `crypto/mlkem`, including implicit rejection. |
| [vaulttransit](vaulttransit/client.go) | 91 | Transit encrypt, decrypt, rotate, and rewrap. Tests use a stand-in, not Vault. |

Go 1.26 or later. See the root [README](../README.md).
