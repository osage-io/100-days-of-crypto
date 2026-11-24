# Phase 2 — Symmetric cryptography

Days 11–22. AES, modes, AEAD, ChaCha20-Poly1305. You already know the names. The point is the failure mode of each one.

## Day 11 — AES, without implementing AES

**Goal.** Describe AES as a 128-bit permutation with a key schedule, and know where the spec says that.

**Read or watch.**

- [FIPS 197](https://csrc.nist.gov/pubs/fips/197/final), sections 1 through 5. Stop before you start memorizing the S-box.
- Boneh, [block-cipher slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/03-block-v2-annotated.pdf).
- Computerphile, [AES Explained](https://www.youtube.com/watch?v=O4xNJsjtN6E). Use it for the picture of the state grid. The spec wins if they disagree.
- Paar's [playlist](https://www.youtube.com/playlist?list=PL6N5qY2nvvJE8X75VkXglSrVhLv1tVcfy), the AES lecture.

**Do.**

- Draw one round: SubBytes, ShiftRows, MixColumns, AddRoundKey. Note that MixColumns is missing in the last round, because FIPS 197 says so.
- From `crypto/aes`, encrypt one 16-byte block. Do not write a mode.

**Reflect.**

- AES is a permutation of 128-bit blocks under a key. Why is that fact alone not an encryption scheme for a message?

## Day 12 — ECB, CBC, and CTR

**Goal.** Know which mode leaks what.

**Read or watch.**

- Boneh, [using block ciphers](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/04-using-block-v2-annotated.pdf). One-time key, then many-time key.
- [NIST SP 800-38A](https://csrc.nist.gov/pubs/sp/800/38/a/final). Read the ECB, CBC, and CTR specifications. Skip the appendices.

**Do.**

- Encrypt 48 identical bytes with AES-ECB and with AES-CTR, using `crypto/cipher`. Print the ciphertext in 16-byte rows.
- Keep the program. Day 73 uses the ECB half.

**Reflect.**

- CTR decryption flips a ciphertext bit into the same plaintext bit. What does a tamperer get from that, before any MAC enters the picture?

## Day 13 — Padding

**Goal.** Treat padding as an oracle waiting to happen, not as a formatting detail.

**Read or watch.**

- Boneh, [authenticated encryption slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/07-authenc-v2-annotated.pdf), the CBC padding section. You will come back to the attack on day 74. Today, just see how the padding is defined and how a decryptor checks it.
- [SP 800-38A](https://csrc.nist.gov/pubs/sp/800/38/a/final), the CBC section, for the rule that padding is outside the mode itself.

**Do.**

- Implement PKCS#7 pad and unpad for a 16-byte block in a scratch file. Make unpad return an error on a bad pad.
- Count how many different error strings you were tempted to return. One is already too many if an attacker can see it.

**Reflect.**

- CTR does not need padding. What did you pay for that, in nonce handling?

## Day 14 — Unauthenticated encryption fails open

**Goal.** State an active attack against a CPA-secure mode that has no MAC.

**Read or watch.**

- Boneh, [authenticated encryption slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/07-authenc-v2-annotated.pdf), "active attacks on CPA-secure encryption."
- [Latacora, 2018](https://www.latacora.com/blog/cryptographic-right-answers/), the "Encrypting data" section. Notice the date. The recommendation has moved. The failure mode of "CTR or CBC alone" has not.

**Do.**

- Take yesterday's CTR ciphertext and flip one byte. Decrypt. Describe the plaintext you got, including the bytes you did not flip.
- Write the attack in one sentence that names the attacker's goal: change the message, not read it.

**Reflect.**

- Confidentiality and integrity are different games. Which game did AES-CTR, used alone, never claim to win?

## Day 15 — AEAD, and the order of MAC and encrypt

**Goal.** Know why encrypt-then-MAC is the construction, and what AEAD packages for you.

**Read or watch.**

- [RFC 5116](https://www.rfc-editor.org/rfc/rfc5116), sections 1 and 2. This is the interface: a nonce, additional data, plaintext in, ciphertext including the tag out.
- Boneh and Shoup, chapter 9, "Authenticated encryption," in the [graduate book](https://toc.cryptobook.us/book.pdf). Read the constructions. Skip the proofs if you can state encrypt-then-MAC versus MAC-then-encrypt without them.
- Boneh's slides from day 14, the "constructions from ciphers and MACs" section.

**Do.**

- In three lines, write MAC-then-encrypt, encrypt-and-MAC, and encrypt-then-MAC. Circle the one TLS 1.3 and Go's AEAD APIs give you.
- Additional data is not encrypted and not optional to authenticate. Name one header you would put there (a version, a key id, a record number).

**Reflect.**

- Why is "we HMAC the plaintext, then encrypt plaintext || tag with CBC" the combination that kept producing oracles?

## Day 16 — AES-GCM

**Goal.** See GCM as CTR plus GHASH, and see the nonce rule as a hard limit.

**Read or watch.**

- [NIST SP 800-38D](https://csrc.nist.gov/pubs/sp/800/38/d/final). Read section 8 on uniqueness of the IV, and the security bounds. You do not need to implement GHASH.
- Adam Langley, [AEADs: getting better at symmetric cryptography](https://www.imperialviolet.org/2015/05/16/aeads.html).

**Do.**

- With `cipher.NewGCM`, seal two messages under one key and two different 12-byte nonces. Confirm both open.
- Read the Go docs for [`crypto/cipher`](https://pkg.go.dev/crypto/cipher) on `NewGCM` and `Seal`. Note the nonce size and the overhead.

**Reflect.**

- SP 800-38D does not say "try not to reuse the nonce." What does it say happens to the authenticator key when you do?

## Day 17 — Nonce reuse, on purpose

**Goal.** Recover plaintext from two AES-GCM messages that shared a nonce, and watch the tags still verify.

**Read or watch.**

- SP 800-38D, the uniqueness requirement, again. One page.
- The package comment in [labs/aesgcm](../labs/aesgcm/nonce.go).

**Do.**

- Before reading the lab, predict what `ct1 XOR ct2` equals when the nonce and key match.
- Then run `go test ./labs/aesgcm/` and read the test. The lab is a demo. It is not a sealing API.

**Reflect.**

- The attack in the lab recovers plaintext when one message is known. The forbidden attack on GHASH forges tags without that. You do not have to implement the forgery. What is the operational rule that makes both attacks impossible?

## Day 18 — ChaCha20

**Goal.** Read a stream cipher spec the way you read FIPS 197: state, rounds, and the nonce.

**Read or watch.**

- [RFC 8439](https://www.rfc-editor.org/rfc/rfc8439), sections 2.1 through 2.4. Run the test vector in section 2.3.2 by hand far enough to see the state layout, or check it with a few lines of Go against `golang.org/x/crypto/chacha20`. Comparing to the RFC vector is the point. Shipping your own ChaCha20 is not.
- Paar's playlist lecture on stream ciphers, if the AES lecture made the contrast useful: [playlist](https://www.youtube.com/playlist?list=PL6N5qY2nvvJE8X75VkXglSrVhLv1tVcfy).

**Do.**

- Write down the ChaCha20 inputs: 256-bit key, 96-bit nonce, 32-bit counter. Compare that nonce to the AES-GCM nonce from day 16.
- Note the sentence in the RFC about unique nonces. It is the same rule, with a different failure.

**Reflect.**

- ChaCha20 alone is still a stream cipher. What is missing, in terms of day 10's games?

## Day 19 — Poly1305

**Goal.** See a one-time authenticator: the key is used for one message.

**Read or watch.**

- [RFC 8439](https://www.rfc-editor.org/rfc/rfc8439), section 2.5, including the test vector in 2.5.2.
- Boneh and Shoup, chapter 7, "Message integrity from universal hashing," in the [graduate book](https://toc.cryptobook.us/book.pdf). Read the idea of a one-time MAC. Poly1305 is the deployed example.

**Do.**

- From the RFC, list the inputs to Poly1305: the 256-bit one-time key (r, s) and the message.
- Find the sentence that says what happens if the key is reused. Copy it into your log.

**Reflect.**

- AEAD constructions derive a fresh Poly1305 key per nonce. What breaks if two messages share that key?

## Day 20 — ChaCha20-Poly1305

**Goal.** Treat the pair as one AEAD, with the nonce and additional-data rules written down.

**Read or watch.**

- [RFC 8439](https://www.rfc-editor.org/rfc/rfc8439), section 2.6 through 2.8, and section 4.
- [`golang.org/x/crypto/chacha20poly1305`](https://pkg.go.dev/golang.org/x/crypto/chacha20poly1305).

**Do.**

- Seal and open one message with `chacha20poly1305.New`. Use a 12-byte nonce.
- Pass additional data on seal and the wrong additional data on open. The open must fail.

**Reflect.**

- The RFC's AEAD is what Go's `cipher.AEAD` already matched. Which arguments of `Seal` are the nonce, the plaintext, and the additional data?

## Day 21 — Limits

**Goal.** Know the message and nonce limits that make "AES-256-GCM" a conditional recommendation.

**Read or watch.**

- SP 800-38D, the security bounds you skipped on day 16. The number you want is how many messages you may encrypt under one key with random 96-bit nonces.
- Filippo Valsorda, [(Post-Quantum) Cryptographic Right Answers](https://words.filippo.io/right-answers/), the "Encrypting data" section. Read the nonce-budget paragraph.
- Filippo Valsorda, [cipher suite ordering](https://words.filippo.io/cipher-suite-ordering/), for why Go prefers ChaCha20-Poly1305 when AES hardware is missing.

**Do.**

- Write the rule you will actually use: random 96-bit nonces are acceptable under a key only up to a message budget you can name; past that, use an extended-nonce construction (XChaCha20-Poly1305 or XAES-256-GCM) or a counter you persist.
- Key commitment: GCM can authenticate the same ciphertext under two different keys. You do not need the attack. Note that "the tag checked" does not bind the key. Filippo's post is the pointer.

**Reflect.**

- Vault Transit's default key type is `aes256-gcm96`. Which limit from today applies to a busy Transit key, and what does the Transit documentation tell an operator to do about it? (You will read that doc on day 91. Guess now, check then.)

## Day 22 — Pick one, in Go

**Goal.** Leave the phase with a default you can defend, and two APIs you will not call.

**Read or watch.**

- [`crypto/cipher`](https://pkg.go.dev/crypto/cipher), `NewGCM` and `AEAD`.
- [`chacha20poly1305`](https://pkg.go.dev/golang.org/x/crypto/chacha20poly1305).
- Latacora's [2018 answers](https://www.latacora.com/blog/cryptographic-right-answers/) and Filippo's [current answers](https://words.filippo.io/right-answers/), side by side, only the symmetric-encryption section.

**Do.**

- Write a 15-line wrapper, in a scratch file, that seals with AES-256-GCM and a random 12-byte nonce, and returns `nonce || ciphertext`. Document the message budget in a comment.
- List the calls you will reject in review: ECB, CBC without an encrypt-then-MAC construction you can point at, and a GCM nonce shorter than 12 bytes.

**Reflect.**

- What did you believe about "AES is secure" on day 1 that you would not write in a design review now?
