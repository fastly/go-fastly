### How to cut a new release for the go-fastly HTTP client
1. Merge all PRs for the release into the master branch
2. Ensure all tests are passing
3. Update the [CHANGELOG.md](https://github.com/fastly/go-fastly/blob/master/CHANGELOG.md):
	* The format is based on https://keepachangelog.com/en/1.0.0/
	* This project adheres to https://semver.org/spec/v2.0.0.html
	* Only add relevant changes
	* Please keep the Unreleased header/link at the top
	* Please don’t forget to update the header links at the bottom
4. Bump the project version in fastly/client.go
5. Checkout and update the master branch:
	* $ git checkout master
	* $ git pull upstream master
6. Create a  new tag for master:
	* $ git tag vx.x.x
7. Push the new tag:
	* $ git push upstream vx.x.x
8. On GitHub, go to the main page of the repository
9. Under the repository name, click Releases
10. Click Draft a new release
	* Select the new tag for the tag version
	* Use the format: vx.x.x - yyyy-mm-dd for the release title
	* Use what you have added to the CHANGELOG for the release description
11. Click Publish release
