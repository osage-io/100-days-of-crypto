# Phase 7 — Attacks and cryptanalysis

Days 71–80. The job is to name the oracle, the assumption it violates, and the fix. Implementing a famous attack is useful once. Collecting writeups is not the goal.

Cryptopals is the exercise set: [set 1](https://cryptopals.com/sets/1), [set 2](https://cryptopals.com/sets/2), [set 3](https://cryptopals.com/sets/3). Do the challenges. Do not publish a solution in the title of a #100DaysOfCrypto post. This repo's `labs/xorbreak` covers the small set-1 cases so the tests have something to run. It is not a walkthrough of the later sets.

## Day 71 — Oracles

**Goal.** Describe an attack as an oracle the attacker is allowed to query.

**Read or watch.**

- Boneh and Shoup, the chosen-ciphertext sections of chapter 2 and the start of chapter 9, in the [graduate book](https://toc.cryptobook.us/book.pdf). You read the definitions on day 10. Today, read them as a description of a network service.
- Your day-10 note on IND-CCA.

**Do.**

- Pick a service you know (a decrypt API, a padding error, a MAC verify, a login). Write the oracle: the input the attacker chooses, the output they see, and the bit they are not supposed to learn.
- If the output is only "success" or "failure," it is still an oracle. Timing is another output. Day 76.

**Reflect.**

- Which of your oracles is available to anyone on the internet, and which one requires a stolen session?

## Day 72 — XOR, on purpose

**Goal.** Finish the set-1 XOR challenges against known answers, in Go.

**Read or watch.**

- [Cryptopals set 1](https://cryptopals.com/sets/1), challenges 2 and 3. Try them before you open the lab.
- The two-time pad is the same fact as day 17: a reused keystream cancels.

**Do.**

- Run `go test ./labs/xorbreak/` after you have a version of your own that passes challenge 2 and recovers challenge 3's plaintext.
- The repeating-key test in the lab is the set-1 repeating-key exercise on a short English plaintext. Challenge 6's file is longer. Do that one from the Cryptopals site, in a private scratch file.

**Reflect.**

- Single-byte XOR falls over because English is not uniform. What would the attack do against random plaintext?

## Day 73 — ECB is a mode, and a bug

**Goal.** Detect ECB, and say why "the blocks are encrypted" was never the claim.

**Read or watch.**

- Cryptopals set 1, challenge 8, on the [set 1 page](https://cryptopals.com/sets/1).
- Your day-12 program that printed ECB ciphertext in 16-byte rows.

**Do.**

- `DetectECB` in [labs/xorbreak](../labs/xorbreak/xor.go) is the whole of challenge 8. Read it after you have written the detection yourself.
- Describe the byte-at-a-time ECB decryption game in one paragraph, without implementing it: the attacker appends bytes until a block boundary, and uses the oracle as an equality test. That is Cryptopals challenge 12. Do it if you have the hour. The description is the requirement.

**Reflect.**

- Why does AES-GCM not have this attack, and what attack does it have instead when the nonce repeats?

## Day 74 — Padding oracles

**Goal.** Explain Vaudenay's attack as a decrypt oracle that returns a padding error.

**Read or watch.**

- Boneh, [authenticated encryption slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/07-authenc-v2-annotated.pdf), "CBC padding attacks."
- Cloudflare, [Padding oracles and the decline of CBC-mode cipher suites](https://blog.cloudflare.com/padding-oracles-and-the-decline-of-cbc-mode-ciphersuites/).

**Do.**

- Write the attacker loop in pseudocode: flip a byte in the previous block, ask "does this pad?", and use the answer to solve for one plaintext byte. Stop when the pseudocode is correct. You do not have to run it against a server.
- Connect this to day 13's error strings. One distinct padding error is the oracle.

**Reflect.**

- Encrypt-then-MAC, checked before any padding parse, closes the oracle. Where in the decrypt function does that check have to sit?

## Day 75 — Bleichenbacher and ROBOT

**Goal.** See an RSA v1.5 padding oracle as the same shape as yesterday, on a different primitive.

**Read or watch.**

- [ROBOT](https://robotattack.org/). Read the attack summary and the mitigations.
- [RFC 8017](https://www.rfc-editor.org/rfc/rfc8017), section 7.2.2, the decryption checks. Those checks, if they return different errors, are the oracle.

**Do.**

- Write the parallel with day 74 in four lines: the secret, the oracle, the query the attacker sends, the fix.
- The fix in TLS 1.3 is not a more careful error message. It is removing RSA key transport. Say that explicitly.

**Reflect.**

- A library that "mitigates Bleichenbacher" by making every error look the same still has to do that in constant time. What did yesterday's padding attack become when the error message was unified but the timing was not?

## Day 76 — Timing

**Goal.** Treat a fast path and a slow path as an output bit.

**Read or watch.**

- Boneh's integrity slides, the timing section: [slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/05-integrity-v2-annotated.pdf).
- [Lucky Thirteen](https://www.isg.rhul.ac.uk/tls/Lucky13.html), the page's own summary.
- [`crypto/subtle`](https://pkg.go.dev/crypto/subtle), `ConstantTimeCompare` and `ConstantTimeCopy`.

**Do.**

- Find one comparison in a repo you work on that checks a MAC, a token, or a password with `==` or `bytes.Equal`. If you cannot find one, write the bad Go and the replacement side by side in your log.
- Lucky13 is the CBC padding oracle when the error messages were already unified. Add that sentence to your day-74 note.

**Reflect.**

- Constant-time code is a local property of the secret-dependent branches. Name one secret in a handshake that is allowed to influence a branch (a public certificate) and one that is not (a MAC key).

## Day 77 — Nonce reuse, both shapes

**Goal.** Put the GCM failure and the ECDSA failure on the same page.

**Read or watch.**

- Your day-17 test and your day-44 note.
- [SP 800-38D](https://csrc.nist.gov/pubs/sp/800/38/d/final), the uniqueness requirement.
- [RFC 6979](https://www.rfc-editor.org/rfc/rfc6979), section 1, the list of what goes wrong when `k` repeats or is biased.

**Do.**

- Two columns. GCM: what the attacker needs (two ciphertexts, or one known plaintext), what they get (plaintext XOR, and a path to tag forgery). ECDSA: what the attacker needs (two signatures with the same `k`), what they get (the private key).
- Write the operational control for each: a nonce counter or a random nonce with a budget you can name; a deterministic signature scheme or a `k` that never leaves the signing function.

**Reflect.**

- Which of these can a log line cause? (A logged ECDSA nonce. A logged GCM nonce is less interesting than a repeated one. A logged key log from day 60 is a third column. Add it.)

## Day 78 — Cryptopals set 2

**Goal.** Do the CBC padding-oracle challenges yourself, with yesterday's pseudocode.

**Read or watch.**

- [Cryptopals set 2](https://cryptopals.com/sets/2). Challenges 9 through 17 are the block. Challenge 16 is the CBC bitflip. Challenge 17 is the padding oracle.

**Do.**

- Solve at least challenge 9 (PKCS#7), challenge 10 (CBC), and challenge 16 or 17. Keep the code in a private directory.
- Do not open a pull request with the solution. A #100DaysOfCrypto post can say "I did set 2 challenge 17" and stop there.

**Reflect.**

- Which challenge was hard because of the crypto, and which was hard because of the encoding? Be honest. Both happen in real bugs.

## Day 79 — Other oracles: compression and errors

**Goal.** Add CRIME/BREACH and error-message oracles to the list, at the level of the leak.

**Read or watch.**

- The CRIME row of your day-55 table. Fill it now if it was blank. Cloudflare's padding-oracle post from day 74 states the compression lesson in the same family of attacks: do not hand the attacker a size channel on a secret.
- [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446), appendix E, the section on compression and the section you already used for 0-RTT. Read the compression warning.

**Do.**

- Write the BREACH/CRIME pattern in one sentence: the attacker chooses a prefix, the compressed length of secret-plus-prefix gets shorter when the guess was right, the length is visible.
- Write one error-oracle pattern from your own experience: a login form, a decrypt API, or a certificate error that distinguishes two cases a caller should not be able to tell apart.

**Reflect.**

- TLS 1.3 removed compression. An application that compresses a secret next to attacker-controlled bytes has put it back. Where in your stack does that happen?

## Day 80 — Read one attack all the way through

**Goal.** Practice extracting the broken assumption from a primary source.

**Read or watch.**

- Pick one: the [ROBOT](https://robotattack.org/) writeup, the [Lucky Thirteen](https://www.isg.rhul.ac.uk/tls/Lucky13.html) page, or one attack talk from the [Real World Crypto 2025 program](https://rwc.iacr.org/2025/program.php) whose recording is on the [IACR channel](https://www.youtube.com/@TheIACR).

**Do.**

- One page, four headings:

  1. The assumption the designer made.
  2. The oracle the attacker actually had.
  3. What the attacker computed.
  4. The fix, and whether the fix is in the primitive, the protocol, or the implementation.

- This is the same outline as day 70. You should be faster at it now.

**Reflect.**

- If you had to teach this attack in fifteen minutes, which of the four headings would you cut, and which one is the one people usually skip?
