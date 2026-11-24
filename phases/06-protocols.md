# Phase 6 — Protocols

Days 55–70. TLS 1.3 in detail, then Noise, Signal, and MLS. Each one is a different answer to "who shares keys with whom, and what happens when a key leaks later."

## Day 55 — What broke before TLS 1.3

**Goal.** Attach each famous attack to a property, so the TLS 1.3 design has a reason.

**Read or watch.**

- Boneh, [authenticated encryption slides](https://crypto.stanford.edu/~dabo/courses/OnlineCrypto/slides/07-authenc-v2-annotated.pdf), the TLS 1.2 case study.
- Cloudflare, [Padding oracles and the decline of CBC-mode cipher suites](https://blog.cloudflare.com/padding-oracles-and-the-decline-of-cbc-mode-ciphersuites/).
- [ROBOT](https://robotattack.org/), the summary.

**Do.**

- Make a table with one row each for BEAST, CRIME, POODLE, Lucky13, and ROBOT. Columns: the property that failed (confidentiality of a byte, compression leak, padding oracle, RSA padding oracle), and the TLS feature that removed it (or did not).
- It is fine if one cell says "I need the primary writeup." Fill it on day 74 or 76 rather than inventing the mechanism.

**Reflect.**

- Which of those attacks still matters to a system that negotiates only TLS 1.3?

## Day 56 — What TLS 1.3 set out to do

**Goal.** List the design cuts: fewer round trips, fewer algorithms, no RSA key transport, no CBC.

**Read or watch.**

- [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446), section 1 and section 1.2.
- Cloudflare, [RFC 8446, aka TLS 1.3](https://blog.cloudflare.com/rfc-8446-aka-tls-1-3/).
- Computerphile, [TLS Handshake Explained](https://www.youtube.com/watch?v=86cQJ0MMses). The video is a sketch. The RFC is the spec. Use the video to see the flight, then close it.

**Do.**

- From RFC 8446 section 1.2, copy the list of properties the handshake is supposed to provide. Do not paraphrase yet.
- Under each property, name the piece you already know: AEAD, ephemeral DH, a signature over the transcript, HKDF.

**Reflect.**

- "Forward secrecy" is one of those properties. Which long-term key can leak tomorrow without opening today's recording?

## Day 57 — The picture

**Goal.** Be able to point at every flight in a TLS 1.3 handshake.

**Read or watch.**

- [The Illustrated TLS 1.3 Connection](https://tls13.xargs.org/). Do this with the browser, not as a PDF skim. Click the bytes.

**Do.**

- On paper, draw ClientHello, ServerHello, EncryptedExtensions, Certificate, CertificateVerify, Finished, in order. Mark the first flight that is encrypted.
- Find the key share and the signature algorithm on the illustrated page. Those are the two choices day 46's checklist was about.

**Reflect.**

- The client sends a key share before it knows which group the server prefers. What does the client send to make that work, and what retry happens if it guessed wrong?

## Day 58 — How to read RFC 8446

**Goal.** Find a message definition and the transcript hash without getting lost.

**Read or watch.**

- [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446), section 2 (protocol overview) and section 4 (handshake protocol), through the ClientHello structure.

**Do.**

- Locate `ClientHello.legacy_version` and the `supported_versions` extension. Write down the fixed value of `legacy_version`.
- Locate the transcript hash. You need to know that Finished and the signatures cover a hash of the handshake, not a hash of the record ciphertext.
- Day 60's lab parses the start of this structure from a live Go handshake. Note the field offsets you will check: content type, legacy record version, handshake type, legacy_version.

**Reflect.**

- Why is the transcript hash a better thing to sign than the raw bytes on the wire, once records can be padded or retransmitted?

## Day 59 — The key schedule

**Goal.** Trace one secret from the DH shared secret to the application traffic key.

**Read or watch.**

- [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446), section 7.1. Read it next to your day-30 notes on `HKDF-Expand-Label`.
- The illustrated handshake's key-schedule panel, if [tls13.xargs.org](https://tls13.xargs.org/) shows one for the secrets you are tracing. The RFC wins on names.

**Do.**

- List the labels you actually care about: `derived`, `c hs traffic`, `s hs traffic`, `c ap traffic`, `s ap traffic`. One sentence each on what the key encrypts.
- Find the early secret. Leave it for day 61. Do not pretend 0-RTT is free.

**Reflect.**

- The handshake traffic keys and the application traffic keys are different. What does that buy you when the handshake messages and the data have different lifetimes?

## Day 60 — A local handshake and a key log

**Goal.** Produce a TLS 1.3 ClientHello, parse the legacy versions, and see the NSS key log labels.

**Read or watch.**

- [labs/tls13](../labs/tls13/handshake.go).
- RFC 8446, appendix D.4, on the middlebox compatibility mode. This is why the record version is 0x0301 while `legacy_version` inside ClientHello is 0x0303.

**Do.**

- Run `go test ./labs/tls13/`.
- Read the captured ClientHello assertion in the test. Confirm it matches what you wrote down on day 58.
- The key log contains traffic secrets. The lab keeps it in memory. Write the rule for a real `SSLKEYLOGFILE`: it is a private key for those sessions.

**Reflect.**

- The lab trusts its own CA through `RootCAs`. What would `InsecureSkipVerify` have skipped that this test still checks?

## Day 61 — 0-RTT

**Goal.** Say what 0-RTT is allowed to protect, and why replay is the whole problem.

**Read or watch.**

- [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446), section 2.3, and section 8 (especially the 0-RTT discussion).
- Appendix E.5 if you want the attacker's view of a replay.

**Do.**

- List three requests that must not be answered from 0-RTT data: anything that is not idempotent. Be concrete. A POST that charges a card. A request that rotates a key. A request that deletes a row.
- List one request that is a plausible fit: a GET of a public object, and even that wants a careful read of the RFC's warnings.

**Reflect.**

- Early data is encrypted to keys derived from a previous session. Which property from day 56 is weaker for those bytes?

## Day 62 — Records, padding, and the weird ClientHello

**Goal.** Read the record layer, and understand why the first flight looks like TLS 1.2 from a distance.

**Read or watch.**

- [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446), section 5, and appendix D.
- Your day-60 capture. You already have the bytes.

**Do.**

- From section 5, write the record header fields and the content types you will see: handshake, application data, alert.
- Note that after the handshake starts encrypting, records are type `application_data` on the outside, with the real type inside. Middleboxes are why.

**Reflect.**

- Padding exists so traffic analysis has less to chew on. What attack from day 55 was a reason to stop compressing plaintext instead?

## Day 63 — crypto/tls, the knobs that matter

**Goal.** Know which `tls.Config` fields still do something in TLS 1.3, and which are 1.2 leftovers.

**Read or watch.**

- [`crypto/tls`](https://pkg.go.dev/crypto/tls), type `Config`. Read `MinVersion`, `MaxVersion`, `CurvePreferences`, `ClientAuth`, `RootCAs`, `Certificates`, `KeyLogWriter`, `CipherSuites`, and `InsecureSkipVerify`.

**Do.**

- Write a server `Config` you would accept in review: `MinVersion: tls.VersionTLS13`, a certificate, no `InsecureSkipVerify`, and no custom cipher-suite list unless you can say which 1.2 suites you still need.
- Read the docs for `CipherSuites` far enough to see that TLS 1.3 does not use that list the way TLS 1.2 did. Go orders those suites itself. Filippo's [cipher suite note](https://words.filippo.io/cipher-suite-ordering/) is the design rationale.

**Reflect.**

- `CurvePreferences` is where a post-quantum hybrid will show up. What do you want the default to be, before you read day 87?

## Day 64 — Noise

**Goal.** Read a handshake pattern as a one-line score, and know what the tokens mean.

**Read or watch.**

- [The Noise Protocol Framework](https://noiseprotocol.org/noise.html), sections 1 through 7. Patterns, tokens (`e`, `s`, `ee`, `es`, `se`, `ss`), and the split into transport keys.
- Latacora, [Factoring the Noise protocol matrix](https://www.latacora.com/blog/2018/07/18/factoring-the-noise/).

**Do.**

- Decode `Noise_IK` and `Noise_XX` from the pattern language. For each, say whether the initiator knows the responder's static key before the first message.
- Write the rule the Latacora post is arguing: if you do not control both ends and you do not have a reason to leave TLS, you stay in TLS.

**Reflect.**

- Noise gives you the handshake and the transport. It does not give you certificates. What did you spend phase 5 building that a Noise deployment still has to invent?

## Day 65 — WireGuard as a concrete Noise system

**Goal.** See one deployed pattern, including the identity key and the replay window.

**Read or watch.**

- [WireGuard paper](https://www.wireguard.com/papers/wireguard.pdf), sections 1 through 5. The Noise pattern is `Noise_IK`.

**Do.**

- From the paper, write where the static keys live, what is in the first handshake message, and how a peer is named (a public key, not a certificate).
- Note the handshake rate limit. That is an engineering control on an unauthenticated first message. Connect it to day 38.

**Reflect.**

- WireGuard's identity is the static Curve25519 key. What does revocation mean in a system with no certificate and no expiry?

## Day 66 — X3DH

**Goal.** See how Signal sets up a conversation when one side is offline.

**Read or watch.**

- [The X3DH specification](https://signal.org/docs/specifications/x3dh/). Read sections 1 through 4. The signed prekey and the one-time prekey are the moving parts.

**Do.**

- Draw the keys each party holds before the first message: identity key, signed prekey, one-time prekey.
- Write which of those can be fetched from a server without the recipient being online, and which signature binds the signed prekey to the identity key.

**Reflect.**

- A one-time prekey that is reused is a different DH. What property degrades?

## Day 67 — The Double Ratchet

**Goal.** Explain the two ratchets, and what each one does after a key leaks.

**Read or watch.**

- [The Double Ratchet specification](https://signal.org/docs/specifications/doubleratchet/), sections 1 through 3. The Diffie-Hellman ratchet and the symmetric-key ratchet.

**Do.**

- In two sentences: the symmetric ratchet moves forward with every message so that a stolen message key does not open the next message. The DH ratchet moves when the parties exchange new ephemeral keys so that a stolen chain key does not open future epochs.
- Find the header that carries the new ratchet public key. That header is what makes the next DH step possible.

**Reflect.**

- This is stronger than TLS's "forward secrecy for the session" story. What user experience pays for it? (Out-of-order messages, and a lot of key state.)

## Day 68 — MLS

**Goal.** See why pairwise ratchets do not scale to a large group, and what a tree buys.

**Read or watch.**

- [RFC 9420](https://www.rfc-editor.org/rfc/rfc9420), sections 1 through 5. Stop after the tree and the commit/welcome sketch if section 5 gets too deep. The overview in sections 2 and 3 is the requirement.
- [Messaging Layer Security](https://messaginglayersecurity.rocks/), the introduction.

**Do.**

- Write the cost of a naive group: a sender encrypts to N members, or shares one key that cannot be rotated when one member leaves. MLS is the alternative.
- From the RFC, define epoch, commit, and welcome in one sentence each.

**Reflect.**

- Removing a member has to change a secret that member knew. Which object in the tree is that secret derived from?

## Day 69 — Four protocols, four jobs

**Goal.** Pick TLS, Noise, Signal, or MLS for a concrete system without treating them as interchangeable.

**Read or watch.**

- Your notes from days 56, 64, 67, and 68. No new spec today.
- [Real World Crypto 2025 program](https://rwc.iacr.org/2025/program.php), only far enough to see that deployed-protocol talks are the literature for this phase. Watch one talk later in the week if a title matches a system you operate. The [IACR channel](https://www.youtube.com/@TheIACR) is where the recordings go.

**Do.**

- For each of these systems, name the protocol you would start from and the reason: a public HTTPS API, a VPN between hosts you both administer, a two-person messenger with offline delivery, a group chat with a few hundred members and removals.
- Write the thing each choice does not give you: TLS does not give you post-compromise security of a long-running chat. Noise does not give you WebPKI. Signal's ratchet is not a group tree. MLS is not a substitute for a CA.

**Reflect.**

- Which of the four do you operate today without being able to name the pattern or the RFC section? That is the one to reread.

## Day 70 — A one-page security review

**Goal.** Practice the review you will be asked for, on a protocol you already touched.

**Read or watch.**

- Nothing new. Bring RFC 8446 or the Noise spec, whichever you can find a section in quickly.

**Do.**

- Write one page, using this outline, about TLS 1.3 as Go implements it, or about WireGuard if that is closer to your work:

  1. Assets and the attacker you are willing to consider.
  2. The protocol and the section numbers you actually read.
  3. Keys: what is long-term, what is ephemeral, what is derived.
  4. One misuse that breaks it (nonce reuse, skipped authentication, 0-RTT on a charge request, a logged key log).
  5. What you would ask an implementer to test.

- Put the page in your log. It is the sample you will turn into a talk on day 99.

**Reflect.**

- Which section did you cite from memory and then fail to find? That is the one to reread before phase 7.
