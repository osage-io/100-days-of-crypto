# 100 Days of Crypto

A free, one-day-at-a-time path from cryptographic foundations to the work a security engineer actually has to defend: TLS 1.3, public-key infrastructure, post-quantum migration, and key management.

This path is written for an experienced security engineer who is not trying to become a mathematician. You are expected to read a spec, write Go, and already know what a threat model is. The new material is the constructions, the assumptions they rest on, and the mistakes that make them fail.

## How to use it

One day is about an hour.

1. Read the day in [`phases/`](phases/01-foundations.md). Read the pages it names. Do not open a second browser full of extra posts.
2. Do the task. On a lab day, try it yourself before you read the code under [`labs/`](labs/README.md).
3. Answer the reflection in a sentence. That sentence is your log entry.
4. Stop. The next day is tomorrow.

Miss a day when you need to. Mark it in [`LOG.md`](LOG.md) and continue. Do not restart at day 1, and do not skip a lab day without coming back to it.

## Philosophy

Cryptography fails in small, specific ways. A nonce reused under AES-GCM does not "weaken" the message. It publishes the XOR of the plaintexts, and the authentication tag still verifies. A textbook RSA implementation does not become safe because the primes are "pretty large."

The order of the path follows that. You meet a primitive, you use it from Go's standard library, and you break a misuse on purpose before you trust it inside TLS, Noise, Signal, or Vault. Post-quantum work comes after you can say which classical assumption a quantum computer actually removes. Key management comes after you know what the key is for.

Nothing here is a production library. `labs/toyrsa` is textbook RSA with no padding. `labs/aesgcm` reuses a nonce on purpose. `labs/caesar` and `labs/xorbreak` recover plaintexts. They are marked in the source. Do not copy them into a service.

## Phases

| Phase | Days | File |
| --- | --- | --- |
| Foundations | 1–10 | [phases/01-foundations.md](phases/01-foundations.md) |
| Symmetric cryptography | 11–22 | [phases/02-symmetric.md](phases/02-symmetric.md) |
| Hashes, MACs, and KDFs | 23–32 | [phases/03-hashes-macs-kdfs.md](phases/03-hashes-macs-kdfs.md) |
| Public key | 33–46 | [phases/04-public-key.md](phases/04-public-key.md) |
| PKI and X.509 | 47–54 | [phases/05-pki.md](phases/05-pki.md) |
| Protocols | 55–70 | [phases/06-protocols.md](phases/06-protocols.md) |
| Attacks and cryptanalysis | 71–80 | [phases/07-attacks.md](phases/07-attacks.md) |
| Post-quantum | 81–89 | [phases/08-post-quantum.md](phases/08-post-quantum.md) |
| Applied key management | 90–94 | [phases/09-key-management.md](phases/09-key-management.md) |
| Giving back | 95–100 | [phases/10-giving-back.md](phases/10-giving-back.md) |

The reading list, with the same links, is in [`RESOURCES.md`](RESOURCES.md).

## #100DaysOfCrypto

When you finish a day, publish one short note. A gist, a blog post, or a few sentences in public is enough.

- Include the day number and the tag `#100DaysOfCrypto`.
- Say what you did, and one thing that surprised you.
- Link the day file or the resource you actually used.
- Do not put a Cryptopals or CryptoHack solution in the title. If you publish a writeup, put the answer behind a clear spoiler warning.
- A missed day is a blank row, not a failure.

Copy [`LOG.md`](LOG.md) and fill it in as you go.

## Labs

The labs are Go. You need **Go 1.26 or later**. [`crypto/mlkem`](https://pkg.go.dev/crypto/mlkem) has been in the standard library since Go 1.24. This repo also depends on `golang.org/x/crypto`, and the current module requires Go 1.26.

```bash
go test ./...
go vet ./...
```

The Vault lab's live test talks to a real dev server only when you ask it to:

```bash
vault server -dev
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=<the root token the dev server prints>
go test ./labs/vaulttransit/ -run TestLiveVault
```

`go test ./...` does not need Vault. It uses an in-process stand-in for the Transit HTTP API.

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md). Resource suggestions use the [issue template](.github/ISSUE_TEMPLATE/resource-suggestion.yml). New resources have to be free, and the link has to load.

## Licenses

- Code under `labs/` is [MIT](LICENSE).
- Written content (the phase notes, this README, and the other prose) is [CC BY 4.0](LICENSE-CONTENT).

You can reuse the writing with attribution. You can reuse the code under the MIT terms. Keep the toy-code warnings attached if you republish a lab.
