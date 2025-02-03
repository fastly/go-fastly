# Release Process

## Prerequisites

For security we sign tags. To be able to sign tags you need to tell Git which key you would like to use. Please follow these
[steps](https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key) to
tell Git about your signing key.

## Steps

1. Merge all PRs intended for the release.
2. Rebase latest remote main branch locally: `git pull --rebase origin main`
3. Ensure all analysis checks and tests are passing: `make all`
4. Open a new PR to update CHANGELOG ([example](https://github.com/fastly/go-fastly/pull/272)).
   - We utilize [semantic versioning](https://semver.org/) and only include relevant/significant changes within the CHANGELOG.
   - Also bump `ProjectVersion` in `fastly/client.go`.
   - If a major version change, then update references to the version in `go.mod` and `README.md` (also in code example tests, `./fastly/example_*_test.go`).
5. Merge CHANGELOG.
6. Rebase latest remote main branch locally: `git pull --rebase origin main`
7. Create a new signed tag (replace `{{remote}}` with the remote pointing to the official repository i.e. `origin` or `upstream` depending on your Git workflow): `tag=vX.Y.Z && git tag -s $tag -m $tag && git push {{remote}} $tag`
8. Copy/paste CHANGELOG into a new [draft release](https://github.com/fastly/go-fastly/releases).
   - Use the format: `vX.Y.Z - yyyy-mm-dd` for the release title.
9. Publish draft release.
