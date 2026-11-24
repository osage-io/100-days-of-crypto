# Phase 8 — Post-quantum

Days 81–89. You do not need a quantum-computing course. You need to know which assumptions die, which NIST standards replaced them, where hybrids sit in TLS, and what CNSA 2.0 actually says.

## Day 81 — What a cryptographically relevant quantum computer changes

**Goal.** Separate Shor's threat from Grover's, in terms of key sizes you already use.

**Read or watch.**

- Boneh and Shoup, chapter 17, "Post-quantum cryptography from lattices," the introduction, in the [graduate book](https://toc.cryptobook.us/book.pdf). Read until the classical assumptions are named. Skip the lattice math on the first pass.
- NIST's [post-quantum project](https://csrc.nist.gov/projects/post-quantum-cryptography), the overview.
- Filippo Valsorda, [Cryptographic Right Answers](https://words.filippo.io/right-answers/), the key-exchange section, and Latacora's [post-quantum update](https://www.latacora.com/blog/post-quantum-cryptographic-right-answers/). The 2018 Latacora post is not the advice for this phase.

**Do.**

- Two columns. Shor: factoring and discrete log, so RSA, finite-field DH, and elliptic-curve DH and signatures. Grover: a square-root speedup against unstructured search, which is why AES-256 and SHA-384 show up in migration tables instead of "AES is dead."
- Write "harvest now, decrypt later" as a sentence about recorded handshakes, not about data that was encrypted under a symmetric key you still control.

**Reflect.**

- Which of your day-1 systems sends a classical key exchange across a network an adversary can record?

## Day 82 — What "standardized" means here

**Goal.** Read a FIPS cover page and know what it does and does not claim.

**Read or watch.**

- [FIPS 203](https://csrc.nist.gov/pubs/fips/203/final), [FIPS 204](https://csrc.nist.gov/pubs/fips/204/final), and [FIPS 205](https://csrc.nist.gov/pubs/fips/205/final). Read the three introductions and the parameter tables. Stop.
- The NIST project page from day 81, the timeline of the competition, far enough to see that these three are the standards and that the names changed: Kyber to ML-KEM, Dilithium to ML-DSA, SPHINCS+ to SLH-DSA.

**Do.**

- Make a table: standard, old name, job (KEM or signature), parameter sets. ML-KEM-512/768/1024. ML-DSA-44/65/87. The SLH-DSA sets are in the FIPS 205 table. Copy the small/fast distinction from the document rather than from memory.
- Note the date on each FIPS. CNSA 2.0's 2022 text still says "Kyber" and "Dilithium" and "TBD" for the specification. Day 88 is about reading that document as a 2022 document.

**Reflect.**

- A FIPS names an algorithm. A CMVP certificate names an implementation. Which one does a procurement form that says "FIPS validated" actually ask for? Day 94 returns to this.

## Day 83 — ML-KEM

**Goal.** Use ML-KEM as a KEM: encapsulation key, ciphertext, shared secret. Not as a block cipher.

**Read or watch.**

- [FIPS 203](https://csrc.nist.gov/pubs/fips/203/final), sections 1 through 3, and the section that states implicit rejection. You want the interface and the failure behavior. You do not need to implement NTT.
- The [Kyber site](https://pq-crystals.org/kyber/index.shtml), the short description, for the history of the name. The FIPS is the spec.
- [`crypto/mlkem`](https://pkg.go.dev/crypto/mlkem).

**Do.**

- From the Go docs, copy the sizes for ML-KEM-768: encapsulation key, ciphertext, shared key, seed. Check them against the constants in the package. The lab asserts them tomorrow.
- Write the sentence from the docs: a ciphertext of the wrong length is an error; a ciphertext of the right length that was not produced by Encapsulate still returns a shared key, and it will not match the sender's.

**Reflect.**

- Implicit rejection removes an oracle. Compare it, in one sentence, to the padding oracle on day 74.

## Day 84 — ML-KEM in Go

**Goal.** Run encapsulation and decapsulation, damage a ciphertext, and watch the shared key change without an error.

**Read or watch.**

- [labs/mlkem](../labs/mlkem/exchange.go).
- The [Go 1.24 release notes](https://go.dev/doc/go1.24), the `crypto/mlkem` paragraph. The API you are calling is that one. This repo builds with Go 1.26 or later because of `golang.org/x/crypto`.

**Do.**

- Run `go test ./labs/mlkem/`.
- Read `TestImplicitRejection` and connect it to yesterday's sentence.
- Read `TestSeedRoundTrip`. The 64-byte seed is the private key you would store. The expanded decapsulation key is a cache.

**Reflect.**

- Most applications should not call `crypto/mlkem`. What should they call, if the goal is a recorded TLS session that survives a later quantum computer? Check your answer on day 87.

## Day 85 — ML-DSA

**Goal.** Know what a Dilithium signature is for, and how big it is, without implementing it.

**Read or watch.**

- [FIPS 204](https://csrc.nist.gov/pubs/fips/204/final), the introduction and the parameter table.
- The [Dilithium site](https://pq-crystals.org/dilithium/index.shtml), the sizes table, as a check on the FIPS names.

**Do.**

- Compare an Ed25519 signature (64 bytes) to the ML-DSA sizes in the FIPS table. Write the parameter set a CNSA reader will ask about. The 2022 CNSA text says Level V. In FIPS 204, check which parameter set is category 5 before you write the name down. Do not trust this sentence over the table.
- Note that Go's standard library, as of the 1.24 notes, standardized the KEM first. Do not invent an import path for ML-DSA if `go doc` does not have one. Write "not in the standard library yet" if that is what you see, and name the library you would evaluate instead on day 97.

**Reflect.**

- A signature that is a few kilobytes changes certificates, handshakes, and CT logs. Which of those does your job touch?

## Day 86 — SLH-DSA, and the stateful schemes CNSA already wanted

**Goal.** Separate stateless hash signatures from stateful ones, and see why state is dangerous.

**Read or watch.**

- [FIPS 205](https://csrc.nist.gov/pubs/fips/205/final), the introduction. SLH-DSA is stateless.
- [SPHINCS+](https://sphincs.org/), the one-page description.
- [NIST SP 800-208](https://csrc.nist.gov/pubs/sp/800/208/final), the introduction. LMS and XMSS. The word you are looking for is "state."

**Do.**

- Two rows. SLH-DSA: stateless, larger signatures, no "how many times have I signed" counter. LMS/XMSS: smaller, and a repeated state value is a key-recovery event. Copy the warning from SP 800-208 rather than softening it.
- CNSA 2.0 picked the stateful schemes for software and firmware signing, with the state requirement attached. You will read that choice in the advisory on day 88. Write the question you want the advisory to answer: why a stateful scheme for firmware, and not for general certificates?

**Reflect.**

- A firmware signing key that lives in an HSM is one answer to "where does the state live?" What does the HSM still not know? (Which build should have been signed.)

## Day 87 — Hybrid TLS

**Goal.** Read the TLS hybrid as a concatenation of two shared secrets, and know the group names.

**Read or watch.**

- [RFC 9954](https://www.rfc-editor.org/rfc/rfc9954), sections 1 and 2. The construction. It does not pick the algorithms.
- [RFC 10024](https://www.rfc-editor.org/rfc/rfc10024), sections 1 through 3. X25519MLKEM768, SecP256r1MLKEM768, SecP384r1MLKEM1024.
- Filippo Valsorda, [The 2025 Go Cryptography State of the Union](https://words.filippo.io/2025-state/), the post-quantum key-exchange section.
- [`crypto/tls`](https://pkg.go.dev/crypto/tls), the `CurvePreferences` docs, and search the page for MLKEM. Prefer the docs over a blog if they disagree about the default.

**Do.**

- Write the hybrid rule in one sentence: the session is secure if at least one of the two component secrets is. Name the failure that breaks X25519MLKEM768 (a break of X25519 and a break of ML-KEM-768).
- From RFC 10024, note which shared secret is concatenated first for X25519MLKEM768, and why the P-256 hybrid is ordered differently. The RFC's FIPS discussion is the reason. Read it. Do not summarize from this paragraph.

**Reflect.**

- A configuration that "turns off the hybrid until the standard is finished" is now a decision to stay with a classical key exchange. What recorded traffic does that decision expose?

## Day 88 — CNSA 2.0

**Goal.** Read the advisory as a procurement and engineering document, including the dates and the 2022 names.

**Read or watch.**

- NSA, [Commercial National Security Algorithm Suite 2.0](https://media.defense.gov/2025/May/30/2003728741/-1/-1/0/CSA_CNSA_2.0_ALGORITHMS.PDF) (PP-22-1338, September 2022). Read the algorithm tables and the timing section. The PDF is the September 2022 text, including where it still says Kyber, Dilithium, and "TBD."
- [RFC 9151](https://www.rfc-editor.org/rfc/rfc9151) is the CNSA profile for TLS as cited by that advisory. It is a CNSA 1.0 profile. The advisory says so.

**Do.**

- Copy the timing bullets as written: software and firmware signing (begin immediately, prefer by 2025, exclusive by 2030), web and cloud (prefer by 2025, exclusive by 2033), traditional networking (2026 / 2030), operating systems (2027 / 2033), niche equipment (2030 / 2033), custom and legacy (2033). The overall NSS line in the advisory is 2035. Check each number against the PDF. Several "by 2025" dates are already behind you. The document has not been rewritten by this repo.
- Copy the symmetric requirements: AES-256, SHA-384 or SHA-512.
- Copy the footnote on hybrids: a hybrid may be how you deploy, and CNSA 2.0 algorithms still have to be selected at the exclusive date. CNSA 1.0 alone is not the end state.
- Translate Level V with the FIPS tables from day 82 and day 85, and write the translation next to the quote, not inside it.

**Reflect.**

- CNSA 2.0 is about National Security Systems. Which of your customers actually cite it, and which of them cite it for systems it does not cover?

## Day 89 — Inventory, then CBOM

**Goal.** Find where cryptography is hiding before you plan a migration.

**Read or watch.**

- [CycloneDX CBOM](https://cyclonedx.org/capabilities/cbom/). Read the description of a cryptographic asset: algorithm, certificate, protocol, related material.
- The [Authoritative Guide to CBOM](https://cyclonedx.org/guides/OWASP_CycloneDX-Authoritative-Guide-to-CBOM-en.pdf), sections 1 and 2. The point of the format is the dependency between an application and an algorithm, not a logo.

**Do.**

- Pick one service you can describe. List ten cryptographic assets without a scanner: TLS version and groups, certificate and issuer, signing keys, HMAC keys, password hash, KMS or Vault key names, SSH host keys, any hardcoded curve, any token format. Mark each classical public-key asset as "recorded traffic is a harvest risk" or "not."
- Put the list in your log. A CBOM is this list in a schema. The list comes first.

**Reflect.**

- Which asset on the list is not in source control and not in a key inventory you could hand to someone else?
