# Contributing

This repo is a public learning path. Changes should make a day clearer, fix a lab, or add a resource that a reader can open for free.

## Licenses

- Code under `labs/` is MIT. See [LICENSE](LICENSE).
- Prose is CC BY 4.0. See [LICENSE-CONTENT](LICENSE-CONTENT).

By sending a pull request you agree that your code is contributed under MIT and your writing under CC BY 4.0.

## Resources

Use the resource-suggestion issue template, or a pull request that edits the day and [`RESOURCES.md`](RESOURCES.md) together.

A resource belongs here when all of these are true:

- You opened the link and it loaded.
- The whole resource is free. A paid certificate, a paywalled book, or a free preview of the first chapter does not count.
- It teaches the construction, the failure, or the operational decision for a specific day.
- It is not a Cryptopals or CryptoHack solution writeup.

If you are unsure the link is stable, say so in the pull request.

## Labs

Labs are Go, under `labs/`, with tests.

- `go test ./...` and `go vet ./...` pass.
- Toy code stays marked as not for production, in the package comment, next to the dangerous function.
- Do not add a Cryptopals solution beyond the small reference labs already in the tree. Point readers at the challenge and let them solve it.
- Prefer the standard library. Use `golang.org/x/crypto` only when the standard library does not have the primitive, or when the lab is explicitly a comparison against it.

## Conduct

[`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) applies to issues, pull requests, and reviews.

## Reports about the labs

The labs are teaching code. `labs/toyrsa` is broken on purpose. Please do not file a product vulnerability report against it.

If a comment or a test would teach someone to ship a real bug, open an ordinary issue and say which day it affects.
