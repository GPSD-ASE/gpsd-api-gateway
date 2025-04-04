### Implementing CI/CD workflow
@author: Kaaviya Ramkumar

@email: prkaaviya17@gmail.com

@date: 22.03.2025

## Steps

1. Create a .commitlintrc.json file at the repository root.
```json
{
  "extends": ["@commitlint/config-conventional"],
  "rules": {
    "type-enum": [2, "always", [
      "feat", "fix", "docs", "style", "refactor", "perf", "test", "chore", "ci", "build", "revert"
    ]]
  }
}
```

2. Install commitlint packages (if not already available).
```bash
npm install --save-dev @commitlint/cli @commitlint/config-conventional
```

3. Create a CHANGELOG.md at the repository root (see CHANGELOG.md).

4. Create a scripts directory and add these files.

- scripts/bump-version.sh
- scripts/update-changelog.sh


5. Make the scripts executable.
```bash
chmod +x scripts/bump-version.sh scripts/update-changelog.sh
```

4.  Setup GitHub Actions workflow with .github/workflows/release.yml.

5. Update Makefile to include release tasks.

6. Add secrets to current GitHub repository.

7. Push to main and check the GitHub Actions tab to see the workflow run

It should update the changelog
Bump the version
Build and push a new Docker image
Update the Helm chart
