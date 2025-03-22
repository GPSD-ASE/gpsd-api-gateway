#!/bin/bash
set -e

# Get the latest tag or default to v0.0.0 if none exists
LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo 'v0.0.0')
TODAY=$(date +%Y-%m-%d)

ORG=$(echo $REMOTE_URL | rev | cut -d'/' -f2 | rev)
REPO=$(echo $REMOTE_URL | rev | cut -d'/' -f1 | rev | sed 's/\.git$//')
REPO_URL="https://github.com/$ORG/$REPO"

echo "Using repository URL: $REPO_URL"

echo "Generating changelog entries since $LATEST_TAG..."

# Create a new temporary changelog
cat > CHANGELOG.new << EOF
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

EOF

# Function to get commits with their hash
get_commits_with_hash() {
    git log --pretty=format:"%h %s" $LATEST_TAG..HEAD | grep -E "$1" | sed "s/^\\([a-f0-9]\\+\\) $2: \\(.*\\)/- \\2 ([\\1](${REPO_URL}\\/commit\\/\\1))/"
}

# Get commits by type
FEATURES=$(get_commits_with_hash "^[a-f0-9]+ feat:" "feat")
FIXES=$(get_commits_with_hash "^[a-f0-9]+ fix:" "fix")
BREAKING=$(get_commits_with_hash "^[a-f0-9]+ BREAKING CHANGE:" "BREAKING CHANGE")

# Only add sections if they have content
if [ ! -z "$FEATURES" ]; then
    echo -e "\n### Added" >> CHANGELOG.new
    echo "$FEATURES" >> CHANGELOG.new
    echo "" >> CHANGELOG.new  # Add a newline after section
fi

if [ ! -z "$FIXES" ]; then
    echo -e "\n### Fixed" >> CHANGELOG.new
    echo "$FIXES" >> CHANGELOG.new
    echo "" >> CHANGELOG.
    
    new  # Add a newline after section
fi

if [ ! -z "$BREAKING" ]; then
    echo -e "\n### Breaking Changes" >> CHANGELOG.new
    echo "$BREAKING" >> CHANGELOG.new
    echo "" >> CHANGELOG.new  # Add a newline after section
fi

# If no changes were found, add a placeholder
if [ -z "$FEATURES" ] && [ -z "$FIXES" ] && [ -z "$BREAKING" ]; then
    echo -e "\n### Maintenance" >> CHANGELOG.new
    echo "- Minor updates and improvements" >> CHANGELOG.new
    echo "" >> CHANGELOG.new  # Add a newline after section
fi

# Replace the changelog
mv CHANGELOG.new CHANGELOG.md
echo "Changelog updated with new entries."