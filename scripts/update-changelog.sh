#!/bin/bash
set -e

LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo 'v0.0.0')
echo "Generating changelog entries since $LAST_TAG..."

# Generate changelog entry for unreleased changes
echo "## [Unreleased]" > /tmp/new-changes.md
echo "" >> /tmp/new-changes.md

# Added section for features
if git log "$LAST_TAG"..HEAD --grep="^feat" --pretty=format:%s | grep -q "."; then
    echo "### Added" >> /tmp/new-changes.md
    git log "$LAST_TAG"..HEAD --grep="^feat" --pretty=format:"- %s" | sed 's/feat: //' >> /tmp/new-changes.md
    echo "" >> /tmp/new-changes.md
fi

# Fixed section for bugs
if git log "$LAST_TAG"..HEAD --grep="^fix" --pretty=format:%s | grep -q "."; then
    echo "### Fixed" >> /tmp/new-changes.md
    git log "$LAST_TAG"..HEAD --grep="^fix" --pretty=format:"- %s" | sed 's/fix: //' >> /tmp/new-changes.md
    echo "" >> /tmp/new-changes.md
fi

# Changed section for refactoring, etc.
if git log "$LAST_TAG"..HEAD --grep="^refactor\|^perf\|^style" --pretty=format:%s | grep -q "."; then
    echo "### Changed" >> /tmp/new-changes.md
    git log "$LAST_TAG"..HEAD --grep="^refactor\|^perf\|^style" --pretty=format:"- %s" | sed 's/refactor: //' | sed 's/perf: //' | sed 's/style: //' >> /tmp/new-changes.md
    echo "" >> /tmp/new-changes.md
fi

# Check if CHANGELOG.md exists, create if not
if [ ! -f CHANGELOG.md ]; then
    echo "# Changelog" > CHANGELOG.md
    echo "" >> CHANGELOG.md
    echo "All notable changes to this project will be documented in this file." >> CHANGELOG.md
    echo "" >> CHANGELOG.md
    echo "The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)," >> CHANGELOG.md
    echo "and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)." >> CHANGELOG.md
    echo "" >> CHANGELOG.md
    echo "## [Unreleased]" >> CHANGELOG.md
    echo "" >> CHANGELOG.md
fi

# Use a different approach for replacing the Unreleased section
# Create a temporary file with the new content
cat CHANGELOG.md > /tmp/changelog.tmp
# Replace the content between ## [Unreleased] and the next section
awk '
    BEGIN { unreleased=0; printed=0 }
    /^## \[Unreleased\]/ {
        print $0;
        system("cat /tmp/new-changes.md");
        unreleased=1;
        printed=1;
        next;
    }
    /^## \[/ {
        if (unreleased && !printed) {
        system("cat /tmp/new-changes.md");
        printed=1;
        }
        unreleased=0;
        print $0;
        next;
    }
    {
        if (!unreleased) print $0;
    }
    ' /tmp/changelog.tmp > CHANGELOG.md

echo "Changelog updated with new entries."