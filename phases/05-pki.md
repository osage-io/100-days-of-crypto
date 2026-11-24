# Phase 5 — PKI and X.509

Days 47–54. A certificate is a signed document. Path validation is a program with a lot of branches. Most outages and most forgery bugs are in those branches, not in RSA.

## Day 47 — What is in a certificate

**Goal.** Read a certificate as a structure, not as a PEM blob.

**Read or watch.**

- [RFC 5280](https://www.rfc-editor.org/rfc/rfc5280), section 4, through the list of fields and the extension overview in 4.2.
- [`crypto/x509`](https://pkg.go.dev/crypto/x509), the `Certificate` struct.

**Do.**

- Generate a self-signed certificate in a scratch file, the same shape as [labs/tls13](../labs/tls13/handshake.go) `newCert`. Print `Subject`, `DNSNames`, `NotAfter`, `KeyUsage`, and `ExtKeyUsage`.
- Decode the PEM with `x509.ParseCertificate` and confirm you read back what you set.

**Reflect.**

- The common name is not the name verifiers should use for TLS. Which field is?

## Day 48 — Path validation and the extensions that bite

**Goal.** Know which extensions change the answer of "is this chain valid?"

**Read or watch.**

- [RFC 5280](https://www.rfc-editor.org/rfc/rfc5280), section 6. Read 6.1. The algorithm is tedious and that is the point. Then read the extension processing for key usage, extended key usage, basic constraints, and name constraints.
- The `VerifyOptions` docs on [`crypto/x509`](https://pkg.go.dev/crypto/x509).

**Do.**

- Build a two-certificate chain in scratch code: a CA and a leaf. Set `IsCA` and `MaxPathLen` only on the CA. Verify the leaf with `Roots` set to the CA.
- Issue a second leaf with `ExtKeyUsage` of code signing and try to verify it as a server. Record the error.

**Reflect.**

- Name constraints are how a public CA can delegate a private namespace without delegating the whole web. Why do they have to be marked critical to mean what you think they mean?

## Day 49 — Trust anchors and chain building

**Goal.** See that "the chain" is an output of path building, not a file the peer gets to dictate.

**Read or watch.**

- RFC 5280, section 6.1, the trust-anchor input, again.
- [`crypto/x509`](https://pkg.go.dev/crypto/x509), `SystemCertPool` and `Verify`.

**Do.**

- Print the number of certificates in `SystemCertPool` on your machine, or note the error if the environment has none. A container often has a smaller pool than a laptop. That difference is an incident waiting for a customer.
- Verify your day-48 leaf once with the right root and once with an empty pool. Keep both errors.

**Reflect.**

- A server that sends the root in the handshake has not saved the client a decision. Who chooses the trust anchor?

## Day 50 — Revocation

**Goal.** Explain CRL, OCSP, stapling, and short-lived certificates as different answers to the same stale-key problem.

**Read or watch.**

- [RFC 5280](https://www.rfc-editor.org/rfc/rfc5280), section 5, CRLs. Read the reason codes.
- [RFC 6960](https://www.rfc-editor.org/rfc/rfc6960), sections 1 through 3. OCSP. Section 2.2 is the request. You want the privacy and soft-fail problems more than the ASN.1.
- Let's Encrypt, [How It Works](https://letsencrypt.org/how-it-works/), the lifetime paragraph. Short lifetime is a revocation strategy.

**Do.**

- Write the failure of "fail open when the OCSP responder is down" in one sentence, from the client's point of view.
- Write the failure of "fail closed" in one sentence, from the operator's point of view.
- That pair is why stapling and short-lived certificates exist. You do not have to build either today.

**Reflect.**

- A private PKI with 90-day leaves and no CRL is a choice. What incident makes that choice wrong?

## Day 51 — Certificate Transparency

**Goal.** Know what a signed certificate timestamp is evidence of, and what it is not.

**Read or watch.**

- [Certificate Transparency](https://certificate.transparency.dev/). Read the introduction.
- [RFC 6962](https://www.rfc-editor.org/rfc/rfc6962), sections 1 and 3. The log is an append-only record of certificates. It does not, by itself, say the certificate was authorized.

**Do.**

- Take a public site you operate or a well-known site, and find its CT entries from the documentation on certificate.transparency.dev (the how-to-look-up section, or the log list). If you cannot reach a lookup tool from this environment, write the steps you would use and the field you would check: the logged certificate's subject and the issuer.
- Write one sentence an incident responder can use: "CT shows the certificate was logged. It does not show the requester was us."

**Reflect.**

- Why is CT more useful after a mis-issuance than before?

## Day 52 — WebPKI and ACME

**Goal.** See public issuance as a protocol with a proof of control, not as a form you submit.

**Read or watch.**

- Let's Encrypt, [How It Works](https://letsencrypt.org/how-it-works/).
- [RFC 8555](https://www.rfc-editor.org/rfc/rfc8555), sections 1 through 3, and section 8 on the challenges. HTTP-01 and DNS-01 are the ones you will meet.

**Do.**

- Draw the ACME actors: account key, order, authorization, challenge, certificate. The account key is not the certificate key.
- List the challenge type you would use for a server that is not reachable on port 80, and the one you would use when the DNS zone is the thing you control.

**Reflect.**

- An ACME client that stores the account key next to the certificate key has collapsed two credentials. What can an attacker do with each?

## Day 53 — Private PKI

**Goal.** Decide what a private CA is for, and what name constraints have to say about it.

**Read or watch.**

- RFC 5280, the name-constraints extension, section 4.2.1.10.
- Vault's [PKI secrets engine](https://developer.hashicorp.com/vault/docs/secrets/pki) overview. Read the introduction only. Day 92 is the working day. Today you want the shape: a mount that issues certificates, with a CA it holds.

**Do.**

- Write a one-page policy, half a page is enough, for a private CA you might actually run. Include: what names it may sign, the maximum leaf lifetime, whether the CA key is online, and who can ask for a certificate.
- Compare that policy to public WebPKI. The difference is the point of a private PKI.

**Reflect.**

- An online intermediate is convenient and it is a second root in practice. What control makes that acceptable?

## Day 54 — crypto/x509, the footguns

**Goal.** Know the Go calls that disable the checks you just spent a week learning.

**Read or watch.**

- [`crypto/x509`](https://pkg.go.dev/crypto/x509), `Certificate.Verify`, `VerifyOptions`, and `SystemCertPool`.
- [`crypto/tls`](https://pkg.go.dev/crypto/tls), `Config.InsecureSkipVerify` and `Config.RootCAs`. Read the docs. The lab on day 60 does not set `InsecureSkipVerify`. It installs the lab CA in `RootCAs`.

**Do.**

- In the scratch program from day 47, call `Verify` with `DNSName` set to a name that is not in the certificate. Record the error.
- Grep a repo you maintain for `InsecureSkipVerify`. If you find it, write the condition under which it is set. "Tests" is only an answer if the test does not talk to a network you do not control.

**Reflect.**

- `InsecureSkipVerify: true` does not mean "I checked the fingerprint somewhere else" unless the code actually checks a fingerprint somewhere else. Does the call site you found do that?
