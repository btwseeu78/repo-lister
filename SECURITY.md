# Security Policy

## Binary Signatures

All release binaries and checksums are signed using [cosign](https://github.com/sigstore/cosign) with keyless signing (Sigstore). This provides cryptographic verification that binaries haven't been tampered with.

### Verifying Signatures

**Install cosign:**
```bash
# macOS
brew install cosign

# Linux
wget "https://github.com/sigstore/cosign/releases/latest/download/cosign-linux-amd64"
sudo mv cosign-linux-amd64 /usr/local/bin/cosign
sudo chmod +x /usr/local/bin/cosign
```

**Verify a release:**
```bash
# Download the binary, checksum, signature, and certificate
VERSION=v0.1.0
curl -LO "https://github.com/btwseeu78/repo-lister/releases/download/${VERSION}/checksums.txt"
curl -LO "https://github.com/btwseeu78/repo-lister/releases/download/${VERSION}/checksums.txt.pem"
curl -LO "https://github.com/btwseeu78/repo-lister/releases/download/${VERSION}/checksums.txt.sig"

# Verify the signature
cosign verify-blob \
  --certificate checksums.txt.pem \
  --signature checksums.txt.sig \
  --certificate-identity-regexp="https://github.com/btwseeu78/repo-lister/.*" \
  --certificate-oidc-issuer="https://token.actions.githubusercontent.com" \
  checksums.txt

# If verification succeeds, verify the binary checksum
curl -LO "https://github.com/btwseeu78/repo-lister/releases/download/${VERSION}/repo-lister_${VERSION#v}_darwin_arm64.tar.gz"
shasum -a 256 -c checksums.txt --ignore-missing
```

### What Gets Signed

- `checksums.txt` (SHA256 hashes of all release artifacts)
- Signature is created using GitHub Actions OIDC with Sigstore
- Certificate and signature are published alongside each release

### Trust Model

Signatures are created using [keyless signing](https://docs.sigstore.dev/cosign/keyless/) with GitHub Actions OIDC. This means:

- No private keys to manage or leak
- Signatures are tied to the GitHub Actions workflow execution
- Verified against Sigstore's transparency log (Rekor)
- Certificate proves the build came from the official `btwseeu78/repo-lister` repository

### Reporting Security Issues

If you discover a security vulnerability, please email security concerns to the repository maintainers or open a private security advisory on GitHub.

**Do not report security vulnerabilities through public GitHub issues.**

## macOS Gatekeeper

**Note:** Cosign signatures verify authenticity but don't bypass macOS Gatekeeper warnings. On first run:

```bash
# Option 1: Remove quarantine attribute
xattr -d com.apple.quarantine /path/to/repo-lister

# Option 2: Right-click → Open (one-time approval)
```

For a fully signed and notarized macOS experience without warnings, an Apple Developer account ($99/year) would be required.
