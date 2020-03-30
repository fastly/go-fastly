### How to cut a new release for the go-fastly HTTP client
1. Merge all PRs for the release into the master branch
1. Ensure all tests are passing
1. Update the [CHANGELOG.md](https://github.com/fastly/go-fastly/blob/master/CHANGELOG.md):
	* The format is based on https://keepachangelog.com/en/1.0.0/
	* This project adheres to https://semver.org/spec/v2.0.0.html
	* **Only add relevant changes**
	* Please keep the Unreleased header/link at the top
	* Please don’t forget to update the header links at the bottom
1. Bump the project version in [fastly/client.go](https://github.com/fastly/go-fastly/blob/master/fastly/client.go)
1. Send PR for the `CHANGELOG` and `client.go` changes.
1. Checkout and update the master branch:
	* `git checkout master`
	* `git pull`
1. Create a  new tag for master:
	* `git tag vx.x.x`
1. Push the new tag:
	* `git push upstream vx.x.x`
1. On GitHub, go to the main page of the repository
1. Under the repository name, click [Releases](https://github.com/fastly/go-fastly/releases)
1. Click [Draft a new release](https://github.com/fastly/go-fastly/releases/new)
	* Select the new tag for the tag version
	* Use the format: `vx.x.x - yyyy-mm-dd` for the release title
	* Use what you have added to the CHANGELOG for the release description
1. Click Publish release
