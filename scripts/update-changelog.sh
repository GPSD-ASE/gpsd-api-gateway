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

# Replace the Unreleased section in CHANGELOG.md
sed -i "" -e '/## \[Unreleased\]/,/^## \[/c\
## [Unreleased]\
\
' CHANGELOG.md

# Insert the new changes after the Unreleased header
awk '
/^## \[Unreleased\]/ {
print; 
system("cat /tmp/new-changes.md"); 
next;
}
{print}
' CHANGELOG.md > /tmp/CHANGELOG.md.new
mv /tmp/CHANGELOG.md.new CHANGELOG.md

echo "Changelog updated with new entries."