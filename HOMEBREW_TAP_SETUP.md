# Homebrew Tap Release Guide

This document explains how to set up and use Homebrew tap for releasing repo-lister.

## Prerequisites

- GitHub account with repo-lister repository
- GoReleaser configured (already done in `.goreleaser.yml`)
- GitHub Actions enabled (already configured)

## Setup Steps

### 1. Create Homebrew Tap Repository

Create a new public repository named `homebrew-tap`:

**Using GitHub CLI:**
```bash
gh repo create homebrew-tap --public --description "Homebrew formulae for repo-lister"
```

**Or via GitHub Web:**
1. Go to https://github.com/new
2. Repository name: `homebrew-tap`
3. Description: "Homebrew formulae for repo-lister"
4. Make it **Public** (required for Homebrew taps)
5. Don't initialize with README
6. Click "Create repository"

### 2. Verify GoReleaser Configuration

The `.goreleaser.yml` already has Homebrew tap configured:

```yaml
brews:
  - name: repo-lister
    repository:
      owner: "{{ .Env.GITHUB_REPOSITORY_OWNER }}"
      name: homebrew-tap
      token: "{{ .Env.GITHUB_TOKEN }}"
    folder: Formula
    homepage: https://github.com/{{ .Env.GITHUB_REPOSITORY_OWNER }}/repo-lister
    description: Container image management tool using Kubernetes credentials
    license: MIT
```

This configuration automatically:
- Generates Homebrew formula
- Pushes formula to `homebrew-tap` repository
- Updates formula on each new release

### 3. Create Your First Release

```bash
# Ensure all changes are committed
git add .
git commit -m "feat: add Homebrew tap support"
git push origin master

# Create and push a version tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

**What Happens Automatically:**
1. GitHub Actions workflow triggers (`.github/workflows/release.yml`)
2. GoReleaser builds Linux binaries for:
   - amd64 (x86_64)
   - arm64
   - armv7
   - armv6
3. Creates GitHub Release with all binaries
4. Generates Homebrew formula
5. Pushes formula to `homebrew-tap/Formula/repo-lister.rb`
6. Creates DEB and RPM packages

### 4. Verify the Release

Check that everything worked:

```bash
# 1. Check GitHub Releases
# Visit: https://github.com/your-username/repo-lister/releases

# 2. Check Homebrew Tap Repository
# Visit: https://github.com/your-username/homebrew-tap
# Should have: Formula/repo-lister.rb

# 3. View the generated formula
gh repo view your-username/homebrew-tap
```

### 5. Test Installation

Test the Homebrew installation:

```bash
# Add your tap
brew tap your-username/tap

# Check available formulae
brew search repo-lister

# Install repo-lister
brew install repo-lister

# Verify installation
repo-lister --help

# Check installed version
repo-lister --version
```

## User Installation Instructions

Share these instructions with users:

```bash
# One-time setup: Add the tap
brew tap your-username/tap

# Install repo-lister
brew install repo-lister

# Or install directly without tapping
brew install your-username/tap/repo-lister

# Upgrade to latest version
brew upgrade repo-lister

# Uninstall
brew uninstall repo-lister
```

## Updating Releases

To release a new version:

```bash
# Make your changes
git add .
git commit -m "feat: add new feature"
git push origin master

# Tag new version (use semantic versioning)
git tag -a v0.2.0 -m "Release v0.2.0: Added new features"
git push origin v0.2.0

# GoReleaser automatically:
# - Builds new binaries
# - Creates new GitHub release
# - Updates Homebrew formula with new version
```

Users update with:
```bash
brew upgrade repo-lister
```

## Troubleshooting

### Issue: Formula Not Created

**Problem:** Homebrew formula wasn't pushed to homebrew-tap repository.

**Solutions:**

1. **Check repository exists:**
   ```bash
   gh repo view your-username/homebrew-tap
   ```

2. **Verify GITHUB_TOKEN permissions:**
   - Go to repo Settings → Actions → General
   - Under "Workflow permissions", ensure:
     - ✅ Read and write permissions
     - ✅ Allow GitHub Actions to create and approve pull requests

3. **Check release logs:**
   ```bash
   gh run list --workflow=release.yml
   gh run view <run-id> --log
   ```

4. **Manual token (if needed):**
   If default GITHUB_TOKEN doesn't work:
   ```bash
   # Create Personal Access Token with 'repo' scope
   # Add to repository secrets:
   gh secret set HOMEBREW_TAP_TOKEN --body "ghp_your_token_here"
   ```

   Then update `.goreleaser.yml`:
   ```yaml
   brews:
     - name: repo-lister
       repository:
         token: "{{ .Env.HOMEBREW_TAP_TOKEN }}"
   ```

### Issue: Formula Fails to Install

**Problem:** `brew install repo-lister` fails.

**Check:**

1. **Formula syntax:**
   ```bash
   brew audit --strict your-username/tap/repo-lister
   ```

2. **Test formula locally:**
   ```bash
   brew install --build-from-source your-username/tap/repo-lister
   ```

3. **Check formula file:**
   ```bash
   brew cat your-username/tap/repo-lister
   ```

### Issue: Architecture Not Supported

**Problem:** Formula installs but binary doesn't work.

**Note:** We only build for Linux:
- ✅ Linux amd64
- ✅ Linux arm64
- ✅ Linux armv7
- ✅ Linux armv6
- ❌ macOS (not built)
- ❌ Windows (not built)

Homebrew on Linux (Linuxbrew) users can install it.
macOS users should download binaries manually from GitHub releases.

## Advanced Configuration

### Adding Dependencies

If your tool requires other packages:

```yaml
brews:
  - name: repo-lister
    dependencies:
      - kubernetes-cli      # Required dependency
      - name: docker
        type: optional      # Optional dependency
```

### Custom Install Steps

For complex installations:

```yaml
brews:
  - name: repo-lister
    install: |
      bin.install "repo-lister"

      # Install completions
      bash_completion.install "completions/repo-lister.bash" => "repo-lister"
      zsh_completion.install "completions/_repo-lister"
      fish_completion.install "completions/repo-lister.fish"

      # Install man pages
      man1.install "man/repo-lister.1"
```

### Multiple Versions

To support multiple major versions:

```yaml
brews:
  - name: repo-lister
    # Installs as 'repo-lister'

  - name: repo-lister@2
    # For version 2.x, installs as 'repo-lister@2'
    # Users: brew install repo-lister@2
```

## Automation Workflow

```
Developer:
├─ 1. Make changes
├─ 2. Commit & push
├─ 3. Create git tag (v0.x.0)
└─ 4. Push tag

GitHub Actions:
├─ 5. Triggers release workflow
├─ 6. GoReleaser builds binaries
├─ 7. Creates GitHub Release
└─ 8. Updates Homebrew formula

User:
├─ 9. brew tap user/tap (one-time)
├─ 10. brew install repo-lister
└─ 11. brew upgrade repo-lister (for updates)
```

## Resources

- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [GoReleaser Homebrew Documentation](https://goreleaser.com/customization/homebrew/)
- [GitHub Actions Permissions](https://docs.github.com/en/actions/security-guides/automatic-token-authentication)

## Quick Reference

```bash
# Create homebrew-tap repo
gh repo create homebrew-tap --public

# Release new version
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# Users install
brew tap your-username/tap
brew install repo-lister

# View formula
gh repo view your-username/homebrew-tap

# Check release status
gh run list --workflow=release.yml
gh run view --web
```

## Example: Complete First Release

```bash
# Step 1: Create homebrew-tap repository
gh repo create homebrew-tap --public --description "Homebrew tap for repo-lister"

# Step 2: Ensure everything is committed
cd repo-lister
git status
git add .
git commit -m "chore: prepare for first release"
git push origin master

# Step 3: Create first release
git tag -a v0.1.0 -m "Initial release v0.1.0"
git push origin v0.1.0

# Step 4: Wait for GitHub Actions (check progress)
gh run list --workflow=release.yml
gh run watch

# Step 5: Verify release created
gh release view v0.1.0

# Step 6: Check homebrew formula
gh repo view your-username/homebrew-tap

# Step 7: Test installation
brew tap your-username/tap
brew install repo-lister
repo-lister --help

# Success! 🎉
```

## Distribution Summary

After setup, your tool will be available via:

1. **Homebrew** (Primary method for Linux users with brew)
   ```bash
   brew install your-username/tap/repo-lister
   ```

2. **Direct Binary Download** (All Linux architectures)
   - GitHub Releases page
   - wget/curl from release assets

3. **Package Managers**
   - DEB packages (Debian/Ubuntu)
   - RPM packages (RHEL/CentOS/Fedora)

4. **Build from Source** (Go developers)
   ```bash
   go install github.com/your-username/repo-lister@latest
   ```
