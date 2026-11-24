# Phase 1 — Foundations

Days 1–10. Enough number theory and security language to read a spec. Not a math degree.

The free books for this phase:

- Dan Boneh and Victor Shoup, [A Graduate Course in Applied Cryptography](https://toc.cryptobook.us/book.pdf) (version 0.6). The [contents page](https://toc.cryptobook.us/) is the map.
- Victor Shoup, [A Computational Introduction to Number Theory and Algebra](https://www.shoup.net/ntb/ntb-v2.pdf). Read the sections named below. Do not read the whole book this week.
- [Handbook of Applied Cryptography](https://cacr.uwaterloo.ca/hac/), free chapter PDFs. Use it as a reference, not as a novel.
- Boneh's [Stanford course page](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/) hosts the slide PDFs. The Coursera video links on that page were dead when this path was written. Use the slides.

## Day 1 — What this path is for

**Goal.** Decide what "cryptography expert" means for your job, and set up the repo.

**Read or watch.**

- The root [README](../README.md), including the log rules.
- [Cryptographic Right Answers (Latacora, 2018)](https://www.latacora.com/blog/cryptographic-right-answers/). Read it as a map of decisions, not as current advice. The post-quantum update is day 81.
- The preface of [Boneh and Shoup](https://toc.cryptobook.us/book.pdf), through "Intended audience".

**Do.**

- Install Go 1.26 or later. Run `go test ./...` from a clone of this repo.
- Copy [LOG.md](../LOG.md) somewhere you will actually edit.
- Write five lines: which protocols you already operate (TLS, SSH, IPsec, Vault, a KMS), and which of those you could not explain the key schedule for.

**Reflect.**

- Which of those five lines is the one your customers ask you about?

## Day 2 — One-time pad, XOR, and letter counting

**Goal.** Separate perfect secrecy from computational security, and see a statistic leak.

**Read or watch.**

- Boneh, [stream-cipher slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/02-stream-v2-annotated.pdf). The one-time pad and "what is a secure cipher" sections.
- Christof Paar's [lecture playlist](https://www.youtube.com/playlist?list=PL6N5qY2nvvJE8X75VkXglSrVhLv1tVcfy), lecture 1, and the [course page](https://www.cryptography-textbook.com/online-course/). The Springer textbook is not free. The lectures and the [first-edition slide zip](https://www.cryptography-textbook.com/files/slides_1st.zip) are.
- HAC [chapter 1](https://cacr.uwaterloo.ca/hac/about/chap1.pdf), sections 1.1 through 1.4, as a glossary pass.

**Do.**

- Before opening the source, break a Caesar ciphertext and a Vigenère ciphertext you encrypt yourself.
- Then read [labs/caesar](../labs/caesar/caesar.go) and run `go test ./labs/caesar/`.
- Convince yourself the one-time pad is only perfect when the key is uniform, as long as the message, and never reused. Two-time pad is day 72.

**Reflect.**

- Frequency analysis did not "crack the algorithm." What statistic did the cipher fail to hide?

## Day 3 — Modular arithmetic

**Goal.** Compute in Z/nZ without hand-waving the inverse.

**Read or watch.**

- Shoup, chapter 2, "Congruences," in the [number theory book](https://www.shoup.net/ntb/ntb-v2.pdf). Stop at the end of the basic congruence rules.
- Boneh, [number-theory slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/10-numth-annotated.pdf), the notation section.
- HAC [chapter 2](https://cacr.uwaterloo.ca/hac/about/chap2.pdf), section 2.4. Use it when a spec says "mod" and you want the precise convention.

**Do.**

- In a short Go program using `math/big`, compute `(a+b) mod n`, `(a*b) mod n`, and a modular inverse for a pair you pick. Print the check `a * inv ≡ 1 (mod n)`.
- Note when the inverse does not exist. That case comes back as "gcd(e, φ(n)) ≠ 1" on day 34.

**Reflect.**

- Why is `int` overflow the wrong tool for this, even when the numbers look small?

## Day 4 — Euclid and the inverse

**Goal.** Compute gcd and modular inverses yourself, once.

**Read or watch.**

- Shoup, chapter 4, "Euclid's algorithm," in the [same book](https://www.shoup.net/ntb/ntb-v2.pdf). The extended algorithm is the part you need.
- HAC [chapter 2](https://cacr.uwaterloo.ca/hac/about/chap2.pdf), section 2.4.3, if you want the algorithm written as numbered steps.

**Do.**

- Implement extended Euclid for `big.Int`. Check it against `new(big.Int).ModInverse`.
- Do not leave this function in a package you might import later. A scratch file is the right place.

**Reflect.**

- `ModInverse` returning nil is a normal answer. What input produces it?

## Day 5 — Groups, rings, fields

**Goal.** Learn the words specs use, and stop there.

**Read or watch.**

- Shoup, the openings of chapter 6 ("Abelian groups") and chapter 7 ("Rings") in the [number theory book](https://www.shoup.net/ntb/ntb-v2.pdf). Read until the definitions sit still. Skip the exercises.
- Boneh and Shoup, the start of chapter 10, "Public key tools," in the [graduate book](https://toc.cryptobook.us/book.pdf). This is where the same words show up next to Diffie-Hellman.

**Do.**

- Write one sentence each for group, ring, and field, using only examples you already have: integers mod n, and integers mod a prime.
- Mark which of those are fields.

**Reflect.**

- AES and ML-KEM both do arithmetic in a field that is not the integers. What do you need from a field that you do not need from a ring?

## Day 6 — Fermat, Euler, and why RSA exponentiation works

**Goal.** See why `m^{ed} ≡ m (mod n)` for textbook RSA, and where that proof stops.

**Read or watch.**

- Boneh, [number-theory slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/10-numth-annotated.pdf), Fermat, Euler, and modular e-th roots.
- HAC [chapter 2](https://cacr.uwaterloo.ca/hac/about/chap2.pdf), Euler's theorem and the definition of φ(n).

**Do.**

- Pick two small primes, compute n and φ(n), pick e = 3 or 65537 if it is coprime to φ(n), and compute d. Check `m^{ed} ≡ m (mod n)` for one m smaller than n.
- This is the whole of day 34's lab, without the warnings. Read those warnings before you keep the code.

**Reflect.**

- The congruence does not say the encryption is safe. What did you not prove?

## Day 7 — Primes

**Goal.** Know what "probable prime" means when Go generates a key.

**Read or watch.**

- Shoup, the opening of chapter 5 ("The distribution of primes") and the opening of chapter 10 ("Probabilistic primality testing") in the [number theory book](https://www.shoup.net/ntb/ntb-v2.pdf).
- The `Prime` function in the [crypto/rand docs](https://pkg.go.dev/crypto/rand).

**Do.**

- Read the Go source for `rand.Prime` (`src/crypto/rand/util.go` in the Go tree, or `go doc crypto/rand`). Note the primality test it calls.
- Generate a 256-bit prime and time it. Then read the comment in [labs/toyrsa](../labs/toyrsa/rsa.go) so the size does not impress you tomorrow.

**Reflect.**

- A primality test that can be wrong with a tiny probability is how real keys are made. What would you need before you called that a bug?

## Day 8 — Discrete log and Diffie-Hellman

**Goal.** Name the assumption a DH group is making, separately from the protocol.

**Read or watch.**

- Boneh, [key-exchange slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/09-basic-key-exchange-annotated.pdf), the Diffie-Hellman section.
- Boneh, [number-theory slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/10-numth-annotated.pdf), "intractable problems."
- HAC [chapter 3](https://cacr.uwaterloo.ca/hac/about/chap3.pdf), the discrete logarithm problem. Skim. You want the problem statement, not the algorithms.

**Do.**

- Write the three names on a card: discrete log, computational Diffie-Hellman, decisional Diffie-Hellman. One sentence each, saying what the attacker is given and what they must produce.
- You will use the card on day 37.

**Reflect.**

- Which of the three does a passive eavesdropper on textbook DH have to solve to read the shared secret?

## Day 9 — Polynomials and finite fields

**Goal.** Recognize a byte as a field element when AES or ML-KEM says so.

**Read or watch.**

- Shoup, the opening of chapter 17 ("Polynomial arithmetic") and chapter 19 ("Finite fields") in the [number theory book](https://www.shoup.net/ntb/ntb-v2.pdf). Stop before the implementation algorithms.
- FIPS 197, section 4, from the [NIST publication page](https://csrc.nist.gov/pubs/fips/197/final). This is the AES field. You only need to see that the bytes are elements of GF(2^8).

**Do.**

- In one paragraph, explain why AES can multiply bytes and always get a byte back, and why "mod 256" is not that multiplication.
- Do not implement the AES S-box.

**Reflect.**

- Where have you already seen a reduction polynomial without calling it that? (Hint: CRC.)

## Day 10 — Security games

**Goal.** Read an IND-CPA or IND-CCA definition without translating it into vibes.

**Read or watch.**

- Boneh and Shoup, chapter 2, "Encryption," in the [graduate book](https://toc.cryptobook.us/book.pdf). Read the security definitions. Skip the long proofs on a first pass. The preface tells you that is allowed.
- Mike Rosulek, [The Joy of Cryptography](https://joyofcryptography.com/pdf/book.pdf). Read the preface and the first chapter that writes a security goal as a game between a challenger and an adversary.
- [Crypto 101](https://www.crypto101.io/) is the shorter companion if the games feel abstract. It is CC BY-NC. Use it as a second explanation, not as the definition.

**Do.**

- Write the IND-CPA game in your own words: what the attacker is allowed to ask, what the challenger flips, and what "advantage" is.
- Then write one sentence on how IND-CCA adds a decryption oracle, and why the challenge ciphertext is excluded.

**Reflect.**

- A padding-oracle attack is an attacker who has which kind of oracle? Keep the answer for day 74.
