# Quick Release Checklist

## One-Time Setup

- [ ] Create `homebrew-tap` repository on GitHub (must be public)
- [ ] Verify `.goreleaser.yml` has brew configuration
- [ ] Ensure GitHub Actions has write permissions

## For Each Release

### 1. Prepare Release
```bash
# Commit all changes
git add .
git commit -m "feat: your changes"
git push origin master
```

### 2. Create Version Tag
```bash
# Use semantic versioning: vMAJOR.MINOR.PATCH
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

### 3. Monitor Release
```bash
# Watch GitHub Actions progress
gh run watch

# Or view in browser
gh run view --web
```

### 4. Verify Release
```bash
# Check release created
gh release view v0.1.0

# Check homebrew formula updated
gh repo view your-username/homebrew-tap
```

### 5. Test Installation
```bash
# Test brew installation
brew tap your-username/tap
brew install repo-lister
repo-lister --help
```

## Release Checklist

- [ ] All tests passing (`go test ./...`)
- [ ] Build successful (`go build .`)
- [ ] Version bumped appropriately
- [ ] Changelog updated (auto-generated)
- [ ] Tag created and pushed
- [ ] GitHub Actions workflow completed
- [ ] Release assets uploaded
- [ ] Homebrew formula updated
- [ ] Installation tested

## Quick Commands

```bash
# Current version
git describe --tags --abbrev=0

# Create next patch version
git tag -a v0.1.1 -m "Release v0.1.1"

# Create next minor version
git tag -a v0.2.0 -m "Release v0.2.0"

# Create next major version
git tag -a v1.0.0 -m "Release v1.0.0"

# Push tag
git push origin <tag-name>

# Delete tag (if mistake)
git tag -d v0.1.0
git push origin :refs/tags/v0.1.0
```

## User Installation

Share with users:
```bash
brew tap your-username/tap
brew install repo-lister
```

## Troubleshooting

If release fails:
1. Check workflow logs: `gh run view --web`
2. Verify homebrew-tap repo exists and is public
3. Ensure GITHUB_TOKEN has write permissions
4. Check `.goreleaser.yml` syntax
5. See full guide: `docs/HOMEBREW_TAP_SETUP.md`
