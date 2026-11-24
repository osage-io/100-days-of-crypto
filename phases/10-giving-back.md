# Phase 10 — Giving back

Days 95–100. The point of the hundred days is that someone else can use what you learned. Pick work that matches how these projects actually take changes.

## Day 95 — Go's crypto packages

**Goal.** Find one real change you are qualified to make, and stop before you mail it half-finished.

**Read or watch.**

- [Contribution guide](https://go.dev/doc/contribute).
- [CONTRIBUTING.md](https://github.com/golang/go/blob/master/CONTRIBUTING.md) in the main repo.
- [Go security policy](https://go.dev/doc/security/policy). A bug in `crypto/tls` is not a GitHub issue if the policy says to mail it.
- Filippo Valsorda, [Post-quantum Cryptography for the Go Ecosystem](https://words.filippo.io/mlkem768/). Read it as a worked example of a change written to be reviewed: a small API, tests, and a reason.

**Do.**

- Clone `golang/go` or browse `src/crypto` and `src/crypto/tls` on GitHub. Find one of: a doc comment that assumes the reader already knows a term you did not know on day 1, a test that does not cover a failure you hit in a lab, or an open issue that is a documentation bug rather than a design change.
- Write the change in a private branch if it is real. Do not open the GitHub issue or the CL unless the contributing guide's steps are done, including the mailing-list or Gerrit step the guide actually names. The log entry is the proposal, not the merge.

**Reflect.**

- Which part of the guide is easy to skip and would waste a reviewer's time? (Usually: a change with no test, or a security report filed in public.)

## Day 96 — BoringSSL and OpenSSL, starting with the docs

**Goal.** See two C libraries the way an outsider is allowed to help: accurate comments and reproducible builds, not a new cipher.

**Read or watch.**

- [BoringSSL README](https://boringssl.googlesource.com/boringssl/+/refs/heads/master/README.md). Read the "don't use this unless you are Chrome or Android or you know why" warning, and the build instructions far enough to see the style of the project.
- [OpenSSL CONTRIBUTING](https://github.com/openssl/openssl/blob/master/CONTRIBUTING.md) and the [community page](https://www.openssl.org/community/).

**Do.**

- Pick one function you now understand (HKDF, AES-GCM, X25519) and find it in one of the two trees. Read the comment and the test name. Write one sentence the comment should contain if a reader has only the RFC.
- If that sentence is already there, write that down. A review that says "this is fine, here is the line" is a finished day.

**Reflect.**

- BoringSSL's README tells most applications to use something else. What should a Go service use, given days 22 and 63?

## Day 97 — liboqs and Open Quantum Safe

**Goal.** Find the post-quantum implementation work that is not inside the Go standard library yet.

**Read or watch.**

- [Open Quantum Safe](https://openquantumsafe.org/).
- [liboqs](https://github.com/open-quantum-safe/liboqs) and its [CONTRIBUTING](https://github.com/open-quantum-safe/liboqs/blob/main/CONTRIBUTING.md).

**Do.**

- From the liboqs README, list the algorithms they build that you met in phase 8: ML-KEM, ML-DSA, SLH-DSA. Note which ones `go doc` still does not offer, if that is still true on your toolchain.
- Read the contributing file's test instructions. Run nothing you cannot build. The day is done when you can say where a documentation fix would go, or where a test vector from a FIPS lives in the tree.

**Reflect.**

- A binding in a language you like is useful only if the C library's contract is the one you think it is. Which FIPS section would you diff the binding against?

## Day 98 — Teach without spoiling

**Goal.** Help the next person through Cryptopals or CryptoHack without handing them the answer.

**Read or watch.**

- [Cryptopals](https://cryptopals.com/sets/1), the set you found hardest.
- [CryptoHack](https://cryptohack.org/). Create an account if you want the exercises. The site is free.

**Do.**

- Write a hint page, not a solution, for one challenge you solved. A good hint names the tool ("this is a padding oracle," "the key is one byte," "read RFC 8439 section 2.5") and not the plaintext.
- If you publish it, put any actual solution below a spoiler warning, and do not file it as a pull request against this repo.
- Offer the hint to one person who is a week behind you, or leave it where a reader of your #100DaysOfCrypto log can find it.

**Reflect.**

- Which hint would have saved you an hour on day 72 or 78, and which hint would have deleted the hour you actually needed?

## Day 99 — Write it down, then say it

**Goal.** Turn one lab into a short post and a short talk.

**Read or watch.**

- One talk from the [IACR channel](https://www.youtube.com/@TheIACR), chosen from the [Real World Crypto 2025 program](https://rwc.iacr.org/2025/program.php). Watch the first ten minutes for the shape: a concrete system, a constraint, a result. You are not copying their topic.
- Your day-70 review.

**Do.**

- Outline an 800-word post. Use one of these, unless your day-70 page is stronger:

  1. AES-GCM nonce reuse: what the API returns, what `ct1 XOR ct2` is, what you want a reviewer to grep for.
  2. The TLS 1.3 ClientHello you parsed: why `legacy_version` is 0x0303 and the record says 0x0301.
  3. ML-KEM implicit rejection, in contrast to a padding oracle.

- Outline a 10-minute talk from the same notes. Six slides is enough: the bug, the spec sentence, the Go call, the test, the review rule, the thing you will not claim.
- Publish the post if you can. A private draft still counts for the log. The tag is `#100DaysOfCrypto`.

**Reflect.**

- What did you cut because you could not point at a spec or a test? Leave it cut.

## Day 100 — Responsible disclosure, and the log

**Goal.** Write a vulnerability report you would be willing to receive, and close the hundred days honestly.

**Read or watch.**

- [Go security policy](https://go.dev/doc/security/policy). Note the address `security@golang.org` and the request to stay off the public tracker until the policy says otherwise.
- [RFC 9116](https://www.rfc-editor.org/rfc/rfc9116), `security.txt`. Read the field list.
- OWASP, [Vulnerability Disclosure Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Vulnerability_Disclosure_Cheat_Sheet.html).

**Do.**

- Write a practice report whose subject is [labs/toyrsa](../labs/toyrsa/rsa.go). This is not a real product bug and you will not mail it to anyone. Use the sections a real report needs:

  1. Summary, one sentence, including "textbook RSA, no padding, teaching code."
  2. Affected component, with the path.
  3. What an attacker who used this as if it were `crypto/rsa` could do. Deterministic ciphertexts and no padding. Do not dress it up.
  4. How you confirmed it. The test name.
  5. The fix: delete it and call `crypto/rsa` or, better, stop encrypting messages with RSA. Point at day 45.
  6. What you want the maintainer to do, and what you will not do (no public exploit, no demand for a CVE for a labeled toy).

- Then write the retrospective in `LOG.md`: ten lines, one per phase, each line a thing you can now explain that you could not on day 1. Publish it with `#100DaysOfCrypto`.

**Reflect.**

- Who will read the next version of your day-70 review, and what do you owe them besides the link?
