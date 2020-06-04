### How to cut a new release for the go-fastly HTTP client
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html); therefore first determine the appropriate version tag based on the change set. If in doubt discuss with the team via Slack before releasing.

1. Merge all PRs intended for the release into the `master` branch
1. Checkout and update the master branch and ensure all tests are passing:
    * `git checkout master`
    * `git pull`
    * `make test`
1. Update the [`CHANGELOG.md`](https://github.com/fastly/go-fastly/blob/master/CHANGELOG.md):
    * Apply necessary labels (`enchancement`, `bug`, `documentation` etc) to all PRs intended for the release that you wish to appear in the `CHANGELOG.md`
    * **Only add labels for relevant changes**
    * `git checkout -b vX.Y.Z` where `vX.Y.Z` is your target version tag
    * `CHANGELOG_GITHUB_TOKEN=xxxx SEMVER_TAG=vX.Y.Z make changelog`
    * `git add CHANGELOG.md && git commit -m "vX.Y.Z"`
1. Bump the project version in fastly/client.go
1. Send PR for the `CHANGELOG.md` and `client.go` changes.
1. Once approved and merged, checkout and update the `master` branch:
    * `git checkout master`
    * `git pull`
1. Create a new tag for `master`:
    * `git tag -s vX.Y.Z -m "vX.Y.Z"`
1. Push the new tag:
    * `git push upstream vX.Y.Z`
1. Under the repository name, click [Releases](https://github.com/fastly/go-fastly/releases)
1. Click [Draft a new release](https://github.com/fastly/go-fastly/releases/new)
	  * Select the new tag for the tag version
	  * Use the format: `vX.Y.Z - yyyy-mm-dd` for the release title
    * Run the following and paste the output in the release description.
       * `make release-changelog`
       * `cat RELEASE_CHANGELOG.md | pbcopy && rm -rf RELEASE_CHANGELOG.md`
1. Click Publish release
1. Celebrate :tada:
