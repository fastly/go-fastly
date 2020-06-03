#!/bin/bash
set -e

if ! command -v github_changelog_generator > /dev/null; then
	echo "No github_changelog_generator in \$PATH, install via 'gem install github_changelog_generator'."
	exit 1
fi

if [ -z "$CHANGELOG_GITHUB_TOKEN" ]; then
	printf "\nWARNING: No \$CHANGELOG_GITHUB_TOKEN in environment, set one to avoid hitting rate limit.\n\n"
fi

if [ -z "$SEMVER_TAG" ]; then
	echo "You must set \$SEMVER_TAG to your desired release semver version."
	exit 1
fi

github_changelog_generator -u fastly -p go-fastly \
	--future-release $SEMVER_TAG \
	--no-pr-wo-labels \
	--no-author \
	--base CHANGELOG_HISTORY.md \
	--since "v1.14.0" \
	--enhancement-label "**Enhancements:**" \
	--bugs-label "**Bug fixes:**" \
	--release-url "https://github.com/fastly/go-fastly/releases/tag/%s" \
	--exclude-labels documentation \
	--exclude-tags-regex "v.*-.*"
