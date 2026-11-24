# Phase 3 — Hashes, MACs, and KDFs

Days 23–32. Collision resistance is not a MAC. A password hash is not a way to turn a DH secret into an AES key.

## Day 23 — What a hash has to resist

**Goal.** Separate preimage, second preimage, and collision, and name length extension.

**Read or watch.**

- Boneh, [collision-resistance slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/06-collision-resistance-v2-annotated.pdf).
- HAC [chapter 9](https://cacr.uwaterloo.ca/hac/about/chap9.pdf), sections 9.1 and 9.2. Definitions only.
- Boneh and Shoup, chapter 8, "Message integrity from collision resistant hashing," in the [graduate book](https://toc.cryptobook.us/book.pdf). The Merkle–Damgård section is the one that explains length extension.

**Do.**

- With `crypto/sha256`, hash a message and hash that message with extra bytes. You will not find a collision. Write down what a collision would be, in terms of the function, so the exercise is the definition and not a search.
- Look up "length extension" in the Boneh slides and write one sentence: who can compute `H(M || pad || extra)` without knowing `M`, and from what input.

**Reflect.**

- Why does length extension make `SHA-256(key || message)` a bad MAC?

## Day 24 — SHA-256 and SHA-3

**Goal.** Know which hash you are being offered, and what a sponge changes.

**Read or watch.**

- [FIPS 180-4](https://csrc.nist.gov/pubs/fips/180-4/upd1/final), the introduction and the table of hash sizes. SHA-1 is in this document. That is not a recommendation.
- [FIPS 202](https://csrc.nist.gov/pubs/fips/202/final), sections 1 through 3. SHA-3 and SHAKE. You want the sponge picture and the output lengths, not the Keccak round.
- [`crypto/sha256`](https://pkg.go.dev/crypto/sha256) and [`crypto/sha3`](https://pkg.go.dev/crypto/sha3).

**Do.**

- Hash the same 32-byte string with SHA-256 and SHA3-256. Confirm the digests differ and both are 32 bytes.
- Note the SHAKE functions on the `sha3` page. Day 31 is the day you use an XOF on purpose. Not today.

**Reflect.**

- A protocol that says "SHA-2" has not picked a function. What do you write instead?

## Day 25 — HMAC

**Goal.** Know why HMAC is the construction, and check a tag in constant time.

**Read or watch.**

- [RFC 2104](https://www.rfc-editor.org/rfc/rfc2104). It is short. Read it all.
- Boneh, [integrity slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/05-integrity-v2-annotated.pdf), the HMAC section, and the timing-attack section.
- [`crypto/hmac`](https://pkg.go.dev/crypto/hmac) and [`crypto/subtle`](https://pkg.go.dev/crypto/subtle).

**Do.**

- Compute HMAC-SHA-256 of a message. Verify it with `hmac.Equal`, not `bytes.Equal`.
- Change one message byte and confirm the verify fails.

**Reflect.**

- `subtle.ConstantTimeCompare` returns 1 or 0, not a bool. Why is that the signature?

## Day 26 — MAC, signature, AEAD tag

**Goal.** Stop using these words as synonyms.

**Read or watch.**

- Boneh's integrity slides from day 25, the opening definitions.
- [RFC 5116](https://www.rfc-editor.org/rfc/rfc5116), section 2.1, for what an AEAD tag covers: ciphertext and additional data, under a nonce.

**Do.**

- Make a three-row table. Columns: who holds the verification key, whether the tag hides the message, whether a nonce is required. Rows: HMAC, Ed25519 (you will build it on day 43; for today the row is "public verification key, no confidentiality"), AES-GCM tag.
- Fill the table from the docs you have already read. Leave a cell blank rather than guess. Check it on day 43.

**Reflect.**

- An AEAD tag is not a signature you can show a third party. What key would that third party need, and why is that the wrong shape?

## Day 27 — Password hashes are not KDFs for protocol keys

**Goal.** Put Argon2 and HKDF in different boxes.

**Read or watch.**

- [RFC 9106](https://www.rfc-editor.org/rfc/rfc9106), sections 1 through 4. Argon2id. Read the parameter guidance. Do not tune memory by folklore.
- Filippo Valsorda, [Cryptographic Right Answers](https://words.filippo.io/right-answers/), the password section.
- The odds-and-ends key-derivation slides: [Boneh, week 4](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/08-odds-and-ends-v2-annotated.pdf).

**Do.**

- Write two sentences you could drop into a design review. One says when Argon2id is the right function (a human password). One says when it is the wrong function (deriving AES keys from a DH shared secret or from a high-entropy IKM).
- Name HKDF as the second function. You read the spec tomorrow.

**Reflect.**

- A slow hash is a feature against guessing. Why is that feature a bug inside a handshake?

## Day 28 — HKDF

**Goal.** Read extract and expand as two different steps with two different inputs.

**Read or watch.**

- [RFC 5869](https://www.rfc-editor.org/rfc/rfc5869). Sections 2 and 3, plus test case 1 in the appendix. This is the spec the lab checks.
- Filippo Valsorda, [The FIPS compliance of HKDF](https://words.filippo.io/fips-hkdf/). Read it for the extract-then-expand split and for how NIST documents talk about it. The FIPS status is the part you will reuse on day 94.

**Do.**

- On paper, label the four inputs: IKM, salt, info, L. Write which step consumes each.
- Read the note in [labs/hkdf](../labs/hkdf/hkdf.go) about a nil salt versus an empty slice. HMAC's zero-padding is why they often match. The RFC's words are still "if not provided, HashLen zeros."

**Reflect.**

- TLS 1.3 derives a lot of keys with different `info` strings from related secrets. What goes wrong if two of those strings are equal?

## Day 29 — HKDF by hand

**Goal.** Implement extract and expand, and match both the RFC vector and `golang.org/x/crypto/hkdf`.

**Read or watch.**

- RFC 5869, appendix A, test case 1, again.
- [`golang.org/x/crypto/hkdf`](https://pkg.go.dev/golang.org/x/crypto/hkdf) and [`crypto/hkdf`](https://pkg.go.dev/crypto/hkdf). The standard-library package is what you call from new Go. The lab compares against `x/crypto` because that is the check this path asked for.

**Do.**

- Read [labs/hkdf/hkdf.go](../labs/hkdf/hkdf.go) only after you have sketched `Extract` and `Expand` yourself.
- Run `go test ./labs/hkdf/`.
- The standard library is the implementation you ship. Delete the scratch copy.

**Reflect.**

- Which test failed first when you had a bug: the RFC vector, or the comparison with `x/crypto`? If neither failed, name the line you were least sure of.

## Day 30 — Labels

**Goal.** See domain separation as a string in the KDF, not as a comment.

**Read or watch.**

- [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446), section 7.1, only the `HKDF-Expand-Label` definition and the list of labels. You are not learning the handshake today. You are learning why the labels exist. The handshake is days 55–63.
- RFC 5869, section 3.2, on the `info` field.

**Do.**

- Using your day-29 function or `crypto/hkdf`, derive two 32-byte keys from the same IKM and salt, with info strings `tls13 key` and `tls13 iv`. Confirm they differ.
- Derive again with the same info string and confirm they match. That match is the bug when two callers share a label.

**Reflect.**

- A label that includes a protocol version is there for a migration. What do you lose if the version is only in a comment?

## Day 31 — SHAKE, when you need a variable output

**Goal.** Know when an XOF is the tool, and when HKDF already did the job.

**Read or watch.**

- [FIPS 202](https://csrc.nist.gov/pubs/fips/202/final), the SHAKE128 and SHAKE256 sections.
- [FIPS 205](https://csrc.nist.gov/pubs/fips/205/final), the introduction only. SLH-DSA uses SHAKE. You are not implementing it. You are seeing why the hash phase shows up again in the post-quantum phase.

**Do.**

- With `crypto/sha3`, read 16 bytes and then 64 bytes from SHAKE256 of the same input. Confirm the 16-byte prefix matches.
- Write one line: "I will not use SHAKE as a password hash." Cross-reference day 27.

**Reflect.**

- HKDF-Expand and SHAKE both produce "more bytes." What input does HKDF take that SHAKE, called on a bare string, does not?

## Day 32 — The Go surface

**Goal.** Leave with a short list of packages and the one comparison function you will accept in review.

**Read or watch.**

- [`crypto/sha256`](https://pkg.go.dev/crypto/sha256), [`crypto/sha3`](https://pkg.go.dev/crypto/sha3), [`crypto/hmac`](https://pkg.go.dev/crypto/hmac), [`crypto/hkdf`](https://pkg.go.dev/crypto/hkdf), [`crypto/subtle`](https://pkg.go.dev/crypto/subtle).

**Do.**

- Write the import paths you want to see in a pull request, and the ones you want a comment on: a hand-rolled HMAC, `==` on a tag, SHA-1, MD5.
- `crypto/md5` and `crypto/sha1` are still in the tree because the world still verifies old signatures. Using them for a new MAC is a different act. Say so in your log.

**Reflect.**

- Which function from this phase is the one Vault Transit is calling when it derives a key from a context? You will check that against the Transit docs on day 91.
