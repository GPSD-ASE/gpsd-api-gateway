#!/bin/bash
set -e

# Get the latest tag or default to v0.0.0 if none exists
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo 'v0.0.0')
TODAY=$(date +%Y-%m-%d)

echo "Generating changelog entries since $LATEST_TAG..."

# Create a new temporary changelog
cat > CHANGELOG.new << EOF
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

EOF

# Get commits since the last tag
COMMITS=$(git log --pretty=format:"- %s" $LATEST_TAG..HEAD | grep -E "^- (feat|fix|BREAKING CHANGE):")

# Group commits by type
FEATURES=$(echo "$COMMITS" | grep -E "^- feat:" | sed 's/^- feat: //')
FIXES=$(echo "$COMMITS" | grep -E "^- fix:" | sed 's/^- fix: //')
BREAKING=$(echo "$COMMITS" | grep -E "^- BREAKING CHANGE:" | sed 's/^- BREAKING CHANGE: //')

# Only add sections if they have content
if [ ! -z "$FEATURES" ]; then
    echo -e "\n### Added" >> CHANGELOG.new
    echo "$FEATURES" >> CHANGELOG.new
fi

if [ ! -z "$FIXES" ]; then
    echo -e "\n### Fixed" >> CHANGELOG.new
    echo "$FIXES" >> CHANGELOG.new
fi

if [ ! -z "$BREAKING" ]; then
    echo -e "\n### Breaking Changes" >> CHANGELOG.new
    echo "$BREAKING" >> CHANGELOG.new
fi

# Replace the changelog
mv CHANGELOG.new CHANGELOG.md
echo "Changelog updated with new entries."