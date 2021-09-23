# Release Process

1. Merge all PRs intended for the release.
2. Rebase latest remote main branch locally (`git pull --rebase origin main`).
3. Ensure all analysis checks and tests are passing (`make all`).
4. Open a new PR to update CHANGELOG ([example](https://github.com/fastly/go-fastly/pull/272))<sup>[1](#note1),[2](#note2),[3](#note3)</sup>.
5. Merge CHANGELOG.
6. Rebase latest remote main branch locally (`git pull --rebase origin main`).
7. Tag a new release (`tag=vX.Y.Z && git tag -s $tag -m "$tag" && git push origin $tag`).
8. Copy/paste CHANGELOG into a new [draft release](https://github.com/fastly/go-fastly/releases)<sup>[4](#note4)</sup>.
9. Publish draft release.
10. Communicate the release in the relevant Slack channels<sup>[5](#note5)</sup>.

## Footnotes

1. <a name="note1"></a>We utilize [semantic versioning](https://semver.org/) and only include relevant/significant changes within the CHANGELOG.
2. <a name="note2"></a>Also bump `ProjectVersion` in `fastly/client.go`.
3. <a name="note3"></a>If a major version change, then update references to the version in `go.mod` and `README.md`.
4. <a name="note4"></a>Use the format: `vX.Y.Z - yyyy-mm-dd` for the release title.
5. <a name="note5"></a>Fastly make internal announcements in the Slack channels: `#api-clients`, `#ecp-languages`.
