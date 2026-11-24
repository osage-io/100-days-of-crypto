# Resources

Every link below was opened while this repo was put together. If one dies, open an issue with the [resource template](.github/ISSUE_TEMPLATE/resource-suggestion.yml) rather than replacing it with a guess.

The phase files say which section to read. This list is the shelf.

## Books and courses

- Dan Boneh and Victor Shoup, [A Graduate Course in Applied Cryptography](https://toc.cryptobook.us/) ([PDF, version 0.6](https://toc.cryptobook.us/book.pdf)).
- Boneh's [Stanford course page](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/) and the slide PDFs linked from it. Slide files used in the path:
  - [Stream ciphers](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/02-stream-v2-annotated.pdf)
  - [Block ciphers](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/03-block-v2-annotated.pdf)
  - [Using block ciphers](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/04-using-block-v2-annotated.pdf)
  - [Message integrity](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/05-integrity-v2-annotated.pdf)
  - [Collision resistance](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/06-collision-resistance-v2-annotated.pdf)
  - [Authenticated encryption](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/07-authenc-v2-annotated.pdf)
  - [Odds and ends, including key derivation](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/08-odds-and-ends-v2-annotated.pdf)
  - [Basic key exchange](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/09-basic-key-exchange-annotated.pdf)
  - [Number theory](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/10-numth-annotated.pdf)
  - [Trapdoor permutations](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/11-pubkey-trapdoor-annotated.pdf)
- Victor Shoup, [A Computational Introduction to Number Theory and Algebra](https://www.shoup.net/ntb/ntb-v2.pdf).
- [Handbook of Applied Cryptography](https://cacr.uwaterloo.ca/hac/). Chapters cited in the path:
  - [1](https://cacr.uwaterloo.ca/hac/about/chap1.pdf), [2](https://cacr.uwaterloo.ca/hac/about/chap2.pdf), [3](https://cacr.uwaterloo.ca/hac/about/chap3.pdf), [8](https://cacr.uwaterloo.ca/hac/about/chap8.pdf), [9](https://cacr.uwaterloo.ca/hac/about/chap9.pdf), [12](https://cacr.uwaterloo.ca/hac/about/chap12.pdf), [13](https://cacr.uwaterloo.ca/hac/about/chap13.pdf).
- Christof Paar, [lecture playlist](https://www.youtube.com/playlist?list=PL6N5qY2nvvJE8X75VkXglSrVhLv1tVcfy), [course page](https://www.cryptography-textbook.com/online-course/), [first-edition slides](https://www.cryptography-textbook.com/files/slides_1st.zip).
- Mike Rosulek, [The Joy of Cryptography](https://joyofcryptography.com/pdf/book.pdf).
- Laurens Van Houtven, [Crypto 101](https://www.crypto101.io/). CC BY-NC.

## What to use

- Latacora, [Cryptographic Right Answers (2018)](https://www.latacora.com/blog/cryptographic-right-answers/). Historical map. Not the current choice list.
- Latacora, [Cryptographic Right Answers: Post Quantum Edition](https://www.latacora.com/blog/post-quantum-cryptographic-right-answers/).
- Filippo Valsorda, [(Post-Quantum) Cryptographic Right Answers](https://words.filippo.io/right-answers/).
- Filippo Valsorda, [cipher suite ordering](https://words.filippo.io/cipher-suite-ordering/), [HKDF and FIPS](https://words.filippo.io/fips-hkdf/), [2025 Go cryptography state of the union](https://words.filippo.io/2025-state/), [ML-KEM for Go](https://words.filippo.io/mlkem768/).
- Adam Langley, [AEADs](https://www.imperialviolet.org/2015/05/16/aeads.html).

## Specifications

- [FIPS 197](https://csrc.nist.gov/pubs/fips/197/final), AES.
- [FIPS 180-4](https://csrc.nist.gov/pubs/fips/180-4/upd1/final), SHA-1 and SHA-2. SHA-1 stays in the document. Do not take that as permission.
- [FIPS 202](https://csrc.nist.gov/pubs/fips/202/final), SHA-3 and SHAKE.
- [FIPS 186-5](https://csrc.nist.gov/pubs/fips/186-5/final), digital signatures.
- [FIPS 140-3](https://csrc.nist.gov/pubs/fips/140-3/final).
- [FIPS 203](https://csrc.nist.gov/pubs/fips/203/final), ML-KEM.
- [FIPS 204](https://csrc.nist.gov/pubs/fips/204/final), ML-DSA.
- [FIPS 205](https://csrc.nist.gov/pubs/fips/205/final), SLH-DSA.
- [SP 800-38A](https://csrc.nist.gov/pubs/sp/800/38/a/final), block cipher modes.
- [SP 800-38D](https://csrc.nist.gov/pubs/sp/800/38/d/final), GCM.
- [SP 800-57 Part 1 Rev. 5](https://csrc.nist.gov/pubs/sp/800/57/pt1/r5/final), key management.
- [SP 800-208](https://csrc.nist.gov/pubs/sp/800/208/final), LMS and XMSS.
- [NIST post-quantum project](https://csrc.nist.gov/projects/post-quantum-cryptography).
- [CMVP](https://csrc.nist.gov/projects/cryptographic-module-validation-program).
- [RFC 2104](https://www.rfc-editor.org/rfc/rfc2104), HMAC.
- [RFC 5116](https://www.rfc-editor.org/rfc/rfc5116), AEAD interface.
- [RFC 5280](https://www.rfc-editor.org/rfc/rfc5280), X.509 path validation.
- [RFC 5869](https://www.rfc-editor.org/rfc/rfc5869), HKDF.
- [RFC 6960](https://www.rfc-editor.org/rfc/rfc6960), OCSP.
- [RFC 6962](https://www.rfc-editor.org/rfc/rfc6962), Certificate Transparency.
- [RFC 6979](https://www.rfc-editor.org/rfc/rfc6979), deterministic ECDSA.
- [RFC 7748](https://www.rfc-editor.org/rfc/rfc7748), X25519 and X448.
- [RFC 7919](https://www.rfc-editor.org/rfc/rfc7919), named finite-field groups.
- [RFC 8017](https://www.rfc-editor.org/rfc/rfc8017), PKCS #1.
- [RFC 8032](https://www.rfc-editor.org/rfc/rfc8032), EdDSA.
- [RFC 8439](https://www.rfc-editor.org/rfc/rfc8439), ChaCha20-Poly1305.
- [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446), TLS 1.3.
- [RFC 8555](https://www.rfc-editor.org/rfc/rfc8555), ACME.
- [RFC 9106](https://www.rfc-editor.org/rfc/rfc9106), Argon2.
- [RFC 9116](https://www.rfc-editor.org/rfc/rfc9116), security.txt.
- [RFC 9151](https://www.rfc-editor.org/rfc/rfc9151), CNSA suite profile for TLS. This is the CNSA 1.0 profile.
- [RFC 9180](https://www.rfc-editor.org/rfc/rfc9180), HPKE.
- [RFC 9420](https://www.rfc-editor.org/rfc/rfc9420), MLS.
- [RFC 9954](https://www.rfc-editor.org/rfc/rfc9954), hybrid key exchange in TLS 1.3.
- [RFC 10024](https://www.rfc-editor.org/rfc/rfc10024), X25519MLKEM768 and the NIST-curve hybrids.
- [SafeCurves](https://safecurves.cr.yp.to/).
- [Ed25519](https://ed25519.cr.yp.to/).
- [Kyber](https://pq-crystals.org/kyber/index.shtml), [Dilithium](https://pq-crystals.org/dilithium/index.shtml), [SPHINCS+](https://sphincs.org/).
- NSA, [CNSA 2.0](https://media.defense.gov/2025/May/30/2003728741/-1/-1/0/CSA_CNSA_2.0_ALGORITHMS.PDF) (PP-22-1338, September 2022).

## Protocols and PKI

- [The Illustrated TLS 1.3 Connection](https://tls13.xargs.org/).
- Cloudflare, [TLS 1.3](https://blog.cloudflare.com/rfc-8446-aka-tls-1-3/) and [padding oracles](https://blog.cloudflare.com/padding-oracles-and-the-decline-of-cbc-mode-ciphersuites/).
- [The Noise Protocol Framework](https://noiseprotocol.org/noise.html).
- Latacora, [Factoring the Noise protocol matrix](https://www.latacora.com/blog/2018/07/18/factoring-the-noise/).
- [WireGuard paper](https://www.wireguard.com/papers/wireguard.pdf).
- [X3DH](https://signal.org/docs/specifications/x3dh/) and the [Double Ratchet](https://signal.org/docs/specifications/doubleratchet/).
- [Messaging Layer Security](https://messaginglayersecurity.rocks/).
- Let's Encrypt, [How It Works](https://letsencrypt.org/how-it-works/).
- [Certificate Transparency](https://certificate.transparency.dev/).
- [ROBOT](https://robotattack.org/).
- [Lucky Thirteen](https://www.isg.rhul.ac.uk/tls/Lucky13.html).

## Go

- [`crypto/cipher`](https://pkg.go.dev/crypto/cipher), [`crypto/hmac`](https://pkg.go.dev/crypto/hmac), [`crypto/sha256`](https://pkg.go.dev/crypto/sha256), [`crypto/sha3`](https://pkg.go.dev/crypto/sha3), [`crypto/hkdf`](https://pkg.go.dev/crypto/hkdf), [`crypto/subtle`](https://pkg.go.dev/crypto/subtle), [`crypto/rand`](https://pkg.go.dev/crypto/rand).
- [`crypto/rsa`](https://pkg.go.dev/crypto/rsa), [`crypto/ecdsa`](https://pkg.go.dev/crypto/ecdsa), [`crypto/ed25519`](https://pkg.go.dev/crypto/ed25519), [`crypto/ecdh`](https://pkg.go.dev/crypto/ecdh).
- [`crypto/x509`](https://pkg.go.dev/crypto/x509), [`crypto/tls`](https://pkg.go.dev/crypto/tls), [`crypto/mlkem`](https://pkg.go.dev/crypto/mlkem).
- [`golang.org/x/crypto/hkdf`](https://pkg.go.dev/golang.org/x/crypto/hkdf), [`golang.org/x/crypto/chacha20poly1305`](https://pkg.go.dev/golang.org/x/crypto/chacha20poly1305).
- [Go 1.24 release notes](https://go.dev/doc/go1.24), [FIPS 140-3 in Go](https://go.dev/doc/security/fips140), [contribution guide](https://go.dev/doc/contribute), [security policy](https://go.dev/doc/security/policy), [CONTRIBUTING.md](https://github.com/golang/go/blob/master/CONTRIBUTING.md).

## Key management

- Vault [Transit](https://developer.hashicorp.com/vault/docs/secrets/transit), [Transit HTTP API](https://developer.hashicorp.com/vault/api-docs/secret/transit), [PKI](https://developer.hashicorp.com/vault/docs/secrets/pki), [seal](https://developer.hashicorp.com/vault/docs/concepts/seal).
- [PKCS #11 3.1](https://docs.oasis-open.org/pkcs11/pkcs11-spec/v3.1/os/pkcs11-spec-v3.1-os.html).
- [CycloneDX CBOM](https://cyclonedx.org/capabilities/cbom/) and the [Authoritative Guide to CBOM](https://cyclonedx.org/guides/OWASP_CycloneDX-Authoritative-Guide-to-CBOM-en.pdf).

## Practice and talks

- [Cryptopals](https://cryptopals.com/), [set 1](https://cryptopals.com/sets/1), [set 2](https://cryptopals.com/sets/2), [set 3](https://cryptopals.com/sets/3).
- [CryptoHack](https://cryptohack.org/).
- Computerphile: [AES](https://www.youtube.com/watch?v=O4xNJsjtN6E), [Diffie-Hellman](https://www.youtube.com/watch?v=NmM9HA2MQGI), [elliptic curves](https://www.youtube.com/watch?v=NF1pwjL9-DE), [TLS handshake](https://www.youtube.com/watch?v=86cQJ0MMses).
- [Real World Crypto 2025 program](https://rwc.iacr.org/2025/program.php) and the [IACR YouTube channel](https://www.youtube.com/@TheIACR).

## Giving back

- [BoringSSL README](https://boringssl.googlesource.com/boringssl/+/refs/heads/master/README.md).
- [OpenSSL CONTRIBUTING](https://github.com/openssl/openssl/blob/master/CONTRIBUTING.md), [OpenSSL community](https://www.openssl.org/community/).
- [Open Quantum Safe](https://openquantumsafe.org/), [liboqs](https://github.com/open-quantum-safe/liboqs), [liboqs CONTRIBUTING](https://github.com/open-quantum-safe/liboqs/blob/main/CONTRIBUTING.md).
- OWASP, [Vulnerability Disclosure Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Vulnerability_Disclosure_Cheat_Sheet.html).

## Left out on purpose

- *Understanding Cryptography* (Paar, Pelzl, Güneysu), the Springer book. The lecture videos and the first-edition slide archive above are free. The book is not.
- Coursera Cryptography I as a free audit. Boneh's Stanford page is the free slide archive. The Coursera lecture URLs linked from that page returned 404, and Coursera's enrollment text no longer promises a free full course. The certificate is paid. The Boneh-Shoup book is the free textbook for the same material.
- A Let's Encrypt page at `/docs/revocation/`. That path returned 404. Revocation in the path uses RFC 5280 and RFC 6960.
- The PKCS #11 "base" URL under the v3.1 directory. That path returned 404. The v3.1 specification link above is the one that loaded.
