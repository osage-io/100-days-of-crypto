# Phase 4 — Public key

Days 33–46. RSA, Diffie-Hellman, elliptic curves, Ed25519, X25519. The arithmetic from phase 1 is the whole trick. The engineering is the padding, the cofactor, and the nonce.

## Day 33 — Textbook RSA

**Goal.** Compute `c = m^e mod n` and say exactly why that sentence is not a scheme.

**Read or watch.**

- Boneh, [trapdoor slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/11-pubkey-trapdoor-annotated.pdf), through "the RSA trapdoor permutation."
- HAC [chapter 8](https://cacr.uwaterloo.ca/hac/about/chap8.pdf), section 8.2. Read the textbook scheme and the list of attacks on it. You are not implementing the attacks today.

**Do.**

- With two primes from day 6 or from `crypto/rand`, compute one encryption and decryption with `math/big`.
- Confirm that encrypting the same `m` twice returns the same `c`.

**Reflect.**

- Deterministic public-key encryption fails IND-CPA immediately. What is the attacker's two queries?

## Day 34 — Toy RSA, marked as a toy

**Goal.** Have a working textbook implementation in front of you, with the warning in the package comment, and then stop using it.

**Read or watch.**

- The comment at the top of [labs/toyrsa](../labs/toyrsa/rsa.go).
- [RFC 8017](https://www.rfc-editor.org/rfc/rfc8017), section 7, the overview of the schemes. You want the names: OAEP for encryption, PSS for signatures. The details are tomorrow.

**Do.**

- Run `go test ./labs/toyrsa/`.
- Read `Generate`. Find the line that refuses to proceed when `gcd(e, φ(n)) ≠ 1`. That is day 4.
- Do not raise the bit size and call it done. The missing padding is the bug that remains at 4096 bits.

**Reflect.**

- The lab keeps `P` and `Q` on the struct so you can see them. What should a real key type do with those factors after `D` is computed? (Chinese Remainder Theorem speedups still need them. Logs do not.)

## Day 35 — PKCS #1 padding

**Goal.** Distinguish v1.5, OAEP, and PSS, and know which one has a padding oracle with a famous name.

**Read or watch.**

- [RFC 8017](https://www.rfc-editor.org/rfc/rfc8017), sections 7.1 (OAEP), 7.2 (v1.5 encryption), 8.1 (PSS), and 8.2 (v1.5 signatures). Read the formatting figures. Skip the ASN.1 appendices.
- Boneh's trapdoor slides, the PKCS #1 section.
- [ROBOT](https://robotattack.org/), the summary of Bleichenbacher's attack. The attack day is day 75. Today you only need to see that v1.5 decryption errors became an oracle.

**Do.**

- In `crypto/rsa`, find `EncryptOAEP`, `SignPSS`, and `SignPKCS1v15`. Write down which RFC section each one implements.
- Note that Go still exposes v1.5 signatures because verifiers still meet them. New encryption should be OAEP if you are encrypting with RSA at all. Day 45 argues for not encrypting with RSA.

**Reflect.**

- A signature padding and an encryption padding fail differently. Which one lets an attacker forge, and which one lets an attacker decrypt, when the check is wrong?

## Day 36 — RSA in practice

**Goal.** Know the sizes, the exponent, and the reason RSA is a bad bulk cipher.

**Read or watch.**

- Boneh's trapdoor slides, "RSA in practice" and "is RSA a one-way function?"
- HAC [chapter 8](https://cacr.uwaterloo.ca/hac/about/chap8.pdf), the implementation notes on CRT decryption if you want the speedup. Optional.
- [FIPS 186-5](https://csrc.nist.gov/pubs/fips/186-5/final), the RSA section's allowed lengths. Read the table, not the standard cover to cover.

**Do.**

- Generate a 2048-bit key with `rsa.GenerateKey` and time it. Compare that, informally, to the X25519 generation you will do on day 42.
- Write the rule: 2048 is a floor for legacy RSA, not a target for a new protocol. The public exponent you want to see is 65537.

**Reflect.**

- CRT decryption uses `p` and `q`. A fault in one half is a key-recovery attack. You do not need the math today. What operational fact follows: private-key operations belong where?

## Day 37 — Diffie-Hellman

**Goal.** Compute the protocol from the assumption you wrote on day 8.

**Read or watch.**

- Boneh, [key-exchange slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/09-basic-key-exchange-annotated.pdf).
- Computerphile, [Secret Key Exchange (Diffie-Hellman)](https://www.youtube.com/watch?v=NmM9HA2MQGI). The video uses a prime modulus. The move to elliptic curves is day 39.
- HAC [chapter 12](https://cacr.uwaterloo.ca/hac/about/chap12.pdf), section 12.6, the basic DH protocol.

**Do.**

- With `math/big`, pick a small prime `p` and generator `g`, choose `a` and `b`, and check `g^{ab} ≡ g^{ba} (mod p)`.
- The shared secret is `g^{ab} mod p`, not the AES key. Write "HKDF next" next to it.

**Reflect.**

- Static DH and ephemeral DH give an attacker different recordings. Which one is "harvest now, decrypt later" actually about?

## Day 38 — DH parameters that are not safe

**Goal.** Recognize a bad group instead of treating "it is DH" as a review comment.

**Read or watch.**

- [RFC 7919](https://www.rfc-editor.org/rfc/rfc7919), sections 1 through 3. Named groups exist because ad hoc primes were a bad idea. TLS 1.3 then dropped finite-field DH. The history is the lesson.
- HAC [chapter 8](https://cacr.uwaterloo.ca/hac/about/chap8.pdf) is the wrong chapter for this; stay with chapter 12's notes on group order if you want a reference, or Boneh's "intractable problems" slides from day 8.
- Boneh and Shoup, chapter 10, the Diffie-Hellman section of the [graduate book](https://toc.cryptobook.us/book.pdf), on active attacks and authenticated key exchange. Read far enough to see that unauthenticated DH does not name the peer.

**Do.**

- Write three failure names in your own words: small subgroup, a smooth group order, and a static key reused with no authentication.
- For each, one sentence on what the attacker learns.

**Reflect.**

- TLS 1.3 does not negotiate custom DH primes. What class of bug does that remove from a review?

## Day 39 — Elliptic curves, as a group

**Goal.** Be able to say "point addition" without pretending you derived the group law.

**Read or watch.**

- Computerphile, [Elliptic Curves](https://www.youtube.com/watch?v=NF1pwjL9-DE).
- Boneh and Shoup, chapter 15, "Elliptic curve cryptography and pairings," in the [graduate book](https://toc.cryptobook.us/book.pdf). Read the opening until the discrete-log problem is restated on points. Skip pairings.
- Paar's playlist, the elliptic-curve lectures: [playlist](https://www.youtube.com/playlist?list=PL6N5qY2nvvJE8X75VkXglSrVhLv1tVcfy).

**Do.**

- Rewrite your day-8 card. Replace `g^a mod p` with `[a]P`. The hard problem is the same shape.
- Do not implement point addition. The standard library already did, on curves you should actually use.

**Reflect.**

- Shor's algorithm is aimed at this problem and at factoring. You will make that precise on day 81. For today: which of your day-8 assumptions does an elliptic curve not remove?

## Day 40 — Curve25519's design constraints

**Goal.** Read a curve recommendation as a list of implementation traps, not as a brand.

**Read or watch.**

- [SafeCurves](https://safecurves.cr.yp.to/). Read the criteria list: ladder, completeness, indistinguishability of points, and the twist. You are collecting words you will meet in reviews.
- [RFC 7748](https://www.rfc-editor.org/rfc/rfc7748), section 1, the motivation.

**Do.**

- From SafeCurves, write down one criterion NIST P-256 is marked badly on, and one sentence on whether that criterion matters for ephemeral Diffie-Hellman in TLS. Filippo's [right answers](https://words.filippo.io/right-answers/) has the practical version of this argument in the key-exchange section. Read that paragraph.
- You are allowed to conclude "use X25519 anyway." You are not allowed to conclude it without the criterion in front of you.

**Reflect.**

- A review comment of "P-256 is backdoored" is usually confused about which attack. What is the attack SafeCurves is actually scoring?

## Day 41 — X25519

**Goal.** Read the Diffie-Hellman function as specified: clamping, the u-coordinate, and the all-zero output.

**Read or watch.**

- [RFC 7748](https://www.rfc-editor.org/rfc/rfc7748), sections 5 and 6. The X25519 test vectors are in section 6.1.
- [`crypto/ecdh`](https://pkg.go.dev/crypto/ecdh).

**Do.**

- Call `ecdh.X25519().GenerateKey` and print the public key length. It is 32 bytes.
- Read the RFC paragraph on checking for the all-zero shared secret. Go's `ECDH` rejects low-order peer keys. The lab test on day 42 locks that in.

**Reflect.**

- The shared secret is not uniform in the way a fresh AES key is. What function sits between `ECDH` and `cipher.NewGCM`?

## Day 42 — X25519 in Go

**Goal.** Run an exchange, confirm both sides match, and reject a garbage public key.

**Read or watch.**

- [labs/x25519](../labs/x25519/exchange.go).
- RFC 7748, the security considerations, section 7.

**Do.**

- Run `go test ./labs/x25519/`.
- Add a scratch call that feeds the shared secret to HKDF from day 29 with the label `day42`. Do not seal with the raw secret.

**Reflect.**

- The lab refuses the all-zero public key. Name one more check a protocol still has to do that X25519 will not do for you. (Authentication of the peer.)

## Day 43 — Ed25519

**Goal.** Separate the signature scheme from the Diffie-Hellman function that shares a curve family.

**Read or watch.**

- [RFC 8032](https://www.rfc-editor.org/rfc/rfc8032), sections 1 through 5. The signing procedure is section 5.1. Read the clamping and the fact that the signature is deterministic.
- [ed25519.cr.yp.to](https://ed25519.cr.yp.to/), the short introduction.
- [`crypto/ed25519`](https://pkg.go.dev/crypto/ed25519).

**Do.**

- Sign a message and verify it. Flip one signature byte and confirm verify fails.
- Verify with the wrong public key and confirm that fails too.
- Update the table from day 26 if a cell was blank.

**Reflect.**

- Ed25519's nonce is derived from the private key and the message. What failure of ECDSA does that choice remove? Tomorrow's reading answers it. Write your guess first.

## Day 44 — ECDSA and the nonce

**Goal.** See a signature scheme where a repeated or biased nonce is a private-key leak.

**Read or watch.**

- [FIPS 186-5](https://csrc.nist.gov/pubs/fips/186-5/final), the ECDSA generation section, far enough to see `k` and the warning about it.
- [RFC 6979](https://www.rfc-editor.org/rfc/rfc6979), sections 1 and 2. Deterministic ECDSA is the mitigation. Ed25519 baked the same idea in from the start.
- [`crypto/ecdsa`](https://pkg.go.dev/crypto/ecdsa). Go's `SignASN1` is deterministic as of the version notes on that page. Read the docs rather than assuming an old blog post.

**Do.**

- Sign the same message twice with `ecdsa.SignASN1` and the same key. Record whether the signatures match. Whatever you observe, write the sentence from the docs that explains it.
- State the attack in one line: if the same `k` signs two messages, the private key is an equation away. You do not have to solve it.

**Reflect.**

- PS3 and a long list of Bitcoin incidents were this bug. What is the review check: "is the nonce generated where, and is it ever logged?"

## Day 45 — KEM-DEM, and HPKE

**Goal.** Stop encrypting messages with a public key. Encapsulate a secret, then use an AEAD.

**Read or watch.**

- [RFC 9180](https://www.rfc-editor.org/rfc/rfc9180), sections 1 through 4. HPKE is the standardized version of "KEM, then KDF, then AEAD."
- Boneh and Shoup, chapter 11, the hybrid-encryption discussion in the [graduate book](https://toc.cryptobook.us/book.pdf).

**Do.**

- Draw four boxes: recipient public key, KEM encapsulation, HKDF, AEAD seal. Label the wire messages: encapsulated key, and ciphertext.
- This is the shape of ML-KEM on day 83, with a different KEM. Write that in the margin.

**Reflect.**

- RSA-OAEP can play the KEM role for a small symmetric key. Why is that still a poor default for a new design in 2026? (Day 81 is the full answer. "Harvest now" is the short one.)

## Day 46 — The Go public-key shelf

**Goal.** Know which package to import, and which function you do not want in a new protocol.

**Read or watch.**

- [`crypto/rsa`](https://pkg.go.dev/crypto/rsa), [`crypto/ecdsa`](https://pkg.go.dev/crypto/ecdsa), [`crypto/ed25519`](https://pkg.go.dev/crypto/ed25519), [`crypto/ecdh`](https://pkg.go.dev/crypto/ecdh).

**Do.**

- Write a review checklist of six lines. Include: OAEP rather than v1.5 encryption, PSS rather than a new v1.5 deployment, X25519 or P-256 via `crypto/ecdh` rather than a hand-rolled ladder, Ed25519 when you want a signature that does not take a nonce from the caller, HKDF on every shared secret, and no textbook RSA.
- Point each line at the package you just opened.

**Reflect.**

- Which of these packages will TLS 1.3 call on your behalf, so that application code should not call them at all?
