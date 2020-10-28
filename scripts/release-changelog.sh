#!/bin/bash
set -e

prev_tag="$(source scripts/tags.sh; previous_tag)"

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
  --no-pr-wo-labels \
  --no-author \
  --no-issues \
  --enhancement-label "**Enhancements:**" \
  --bugs-label "**Bug fixes:**" \
  --release-url "https://github.com/fastly/go-fastly/releases/tag/%s" \
  --exclude-labels documentation \
  --exclude-tags-regex "v.*-.*" \
  --output RELEASE_CHANGELOG.md \
  --since-tag $prev_tag
