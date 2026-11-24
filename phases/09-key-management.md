# Phase 9 — Applied key management

Days 90–94. A key has a lifecycle. Vault Transit and Vault PKI are two different lifecycles. An HSM does not pick the algorithm. FIPS 140-3 is a claim about a module, not about your protocol.

## Day 90 — Lifecycle

**Goal.** Write generation, storage, use, rotation, and destruction as separate steps with owners.

**Read or watch.**

- [NIST SP 800-57 Part 1 Rev. 5](https://csrc.nist.gov/pubs/sp/800/57/pt1/r5/final), section 1 and the key-state discussion in section 8. Read for the states. Skip the enormous tables on first pass.
- HAC [chapter 13](https://cacr.uwaterloo.ca/hac/about/chap13.pdf), sections 13.1 and 13.2, as the shorter engineering chapter.

**Do.**

- Take one key from your day-89 inventory. Fill a table: who generates it, where it lives, which API uses it, what rotation means (new version, or new certificate, or new session), and how you know the old one is gone.
- If "rotation" is "we restart the process," write that down. It is a real answer and usually a bad one.

**Reflect.**

- Destruction is the step people skip. What does "gone" mean for a key that was in a replica, a backup, and a core dump?

## Day 91 — Vault Transit

**Goal.** Encrypt, rotate, and rewrap through the Transit HTTP API, and read the ciphertext prefix correctly.

**Read or watch.**

- Vault, [Transit secrets engine](https://developer.hashicorp.com/vault/docs/secrets/transit). Do the encrypt, decrypt, rotate, and rewrap steps in the docs.
- The [Transit HTTP API](https://developer.hashicorp.com/vault/api-docs/secret/transit), the encrypt and rewrap endpoints. Plaintext is base64 because it is not text.
- [labs/vaulttransit](../labs/vaulttransit/client.go). The tests speak to a stand-in. The prefix `vault:v<keyVersion>:` matches Transit's version template. The bytes after the prefix in the stand-in are not Vault's blob. The client does not care. It sends the string back.

**Do.**

- Run `go test ./labs/vaulttransit/`.
- Against a dev server, if you have one:

  ```bash
  vault server -dev
  export VAULT_ADDR=http://127.0.0.1:8200
  export VAULT_TOKEN=<the root token the dev server prints>
  go test ./labs/vaulttransit/ -run TestLiveVault
  ```

- After rotate, decrypt a ciphertext minted on version 1. Then rewrap it and decrypt the result. The plaintext should not appear in the rewrap response. Check that in the test.
- Read the Transit doc's note on the AES-GCM message limit and rotation. Compare it to day 21.

**Reflect.**

- Rewrap is a privilege that does not return plaintext. Why is that a better grant for a background job than `decrypt`?

## Day 92 — Vault PKI is not Transit

**Goal.** Issue a certificate from a PKI mount and say why that is a different product from encrypting a column.

**Read or watch.**

- Vault, [PKI secrets engine](https://developer.hashicorp.com/vault/docs/secrets/pki). Read the setup for a root and an intermediate, and the issue endpoint. You do not need every tune flag.
- Your day-53 policy.

**Do.**

- On a dev server, or on paper if you cannot run Vault: enable `pki`, generate a root, issue one leaf with a short TTL and a name constraint or an allowed domain that matches the policy you wrote on day 53.
- Write the difference in two sentences. Transit: the key never leaves, the ciphertext does, rotation rewraps ciphertext. PKI: the CA key signs, the leaf private key is handed to a caller or generated with the certificate, rotation is issuance and revocation.

**Reflect.**

- Putting a TLS key in Transit and asking Transit to sign the handshake is the wrong shape. What would you issue from PKI instead?

## Day 93 — HSMs and PKCS #11

**Goal.** Know what an HSM actually holds, and what a PKCS #11 session is asking it to do.

**Read or watch.**

- [PKCS #11 3.1](https://docs.oasis-open.org/pkcs11/pkcs11-spec/v3.1/os/pkcs11-spec-v3.1-os.html), the overview and the object model: token, session, object, attribute. You are not memorizing every function code.
- Vault, [seal concepts](https://developer.hashicorp.com/vault/docs/concepts/seal). The seal is how Vault protects its own master key. It is not the Transit key your application asked to encrypt a row.

**Do.**

- Draw three keys that people conflate: the seal key (unwraps Vault), a Transit key (encrypts application data), a PKI CA key (signs certificates). For each, say whether an HSM-backed seal replaces it. Only the seal is replaced by "auto-unseal with an HSM." The others move only if you configured them to.
- Write one thing an HSM does not check: whether you should have signed this firmware, whether this nonce has been used, whether this certificate name is the one the customer asked for.

**Reflect.**

- "It is in an HSM" is not an answer to a padding-oracle bug in the application. Which day-74 fix is still yours?

## Day 94 — FIPS 140-3

**Goal.** Read the level numbers as a buyer, and read Go's FIPS mode as a build flag with a narrow claim.

**Read or watch.**

- [FIPS 140-3](https://csrc.nist.gov/pubs/fips/140-3/final), the introduction and the section that defines the four levels. Levels are about the module boundary and the physical and role requirements. They are not a ranking of algorithms.
- The [CMVP project](https://csrc.nist.gov/projects/cryptographic-module-validation-program), the "what is validated" explanation.
- Go, [FIPS 140-3 compliance](https://go.dev/doc/security/fips140). Read what `GOFIPS140` and `GODEBUG=fips140` change, and the paragraph that says the Go team does not claim your regulatory problem is solved.
- Filippo Valsorda, [The FIPS compliance of HKDF](https://words.filippo.io/fips-hkdf/), the part about which documents list HKDF. You read it on day 28. Read the conclusion again with a module in mind.

**Do.**

- Write three sentences you can say to a customer. (1) FIPS 140-3 level N is a statement about a cryptographic module. (2) A validated module can still be used with the wrong algorithm, the wrong nonce, and the wrong key lifetime. (3) Go can build a binary that sticks to the approved set. That binary is not, by itself, a completed CMVP certificate. Check sentence 3 against the Go doc before you repeat it.
- From the Go doc, note one algorithm TLS will refuse in FIPS-only mode. That refusal is why a "just turn on FIPS" flag changes handshakes.

**Reflect.**

- CNSA's advisory told you not to assess NSS against "FIPS-validated" alone. Where did it say the approval has to come from instead? Check day 88's notes on the enforcement section before you answer from memory.
