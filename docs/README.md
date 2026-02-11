# repo-lister

A CLI tool to manage container images across registries using Kubernetes credentials.

## Features

- **list** - List image tags from a container registry
- **copy** - Copy/retag images between registries without local storage
- **pull** - Pull images from registry to local tar files
- **push** - Push images from local tar files to registry

All commands use Kubernetes secrets for registry authentication, making it easy to work with private registries in your cluster.

## Security

🔒 **All releases are cryptographically signed** using [cosign](https://github.com/sigstore/cosign) with keyless signing (Sigstore). See [SECURITY.md](SECURITY.md) for verification instructions.

## Installation

### Homebrew (macOS & Linux)

```sh
brew install --cask btwseeu78/repo-lister/repo-lister
```

### Download Pre-built Binary

Download the latest release for your platform from the [releases page](https://github.com/btwseeu78/repo-lister/releases).

```sh
# macOS (Apple Silicon)
curl -LO https://github.com/btwseeu78/repo-lister/releases/latest/download/repo-lister_<version>_darwin_arm64.tar.gz
tar -xzf repo-lister_<version>_darwin_arm64.tar.gz
sudo mv repo-lister /usr/local/bin/

# Linux (amd64)
curl -LO https://github.com/btwseeu78/repo-lister/releases/latest/download/repo-lister_<version>_linux_x86_64.tar.gz
tar -xzf repo-lister_<version>_linux_x86_64.tar.gz
sudo mv repo-lister /usr/local/bin/

# Linux (arm64)
curl -LO https://github.com/btwseeu78/repo-lister/releases/latest/download/repo-lister_<version>_linux_arm64.tar.gz
tar -xzf repo-lister_<version>_linux_arm64.tar.gz
sudo mv repo-lister /usr/local/bin/
```

### Package Manager

**Debian/Ubuntu:**
```sh
curl -LO https://github.com/btwseeu78/repo-lister/releases/latest/download/repo-lister_<version>_amd64.deb
sudo dpkg -i repo-lister_<version>_amd64.deb
```

**RHEL/CentOS/Fedora:**
```sh
curl -LO https://github.com/btwseeu78/repo-lister/releases/latest/download/repo-lister_<version>_x86_64.rpm
sudo rpm -i repo-lister_<version>_x86_64.rpm
```

### Build from Source

```sh
git clone https://github.com/your-github-username/repo-lister.git
cd repo-lister
go build .
sudo mv repo-lister /usr/local/bin/
```

## Requirements

- Connectivity to Kubernetes cluster
- Valid kubeconfig (or in-cluster configuration)
- Access to get secrets in the specified namespace
- Kubernetes secrets of type `kubernetes.io/dockerconfigjson`

## Commands

### 1. List - List image tags

List all available tags for a container image from a registry.

```sh
repo-lister list \
  --image <image-name> \
  --secret <secret-name> \
  --namespace <namespace> \
  --limit <number>
```

**Flags:**
- `-i, --image` - Image name to list tags for (required)
- `-s, --secret` - Kubernetes secret name for authentication (required)
- `-n, --namespace` - Kubernetes namespace where secret is located (default: "default")
- `-f, --filter` - Regex filter to apply to image tags (default: ".*")
- `-l, --limit` - Maximum number of tags to return (default: 5)

**Examples:**

```sh
# List latest 5 tags
repo-lister list \
  --image linuxarpan/testpush \
  --secret regcred \
  --namespace default \
  --limit 5

# List tags with filter
repo-lister list \
  --image myregistry.io/app \
  --secret registry-cred \
  --filter "v[0-9]+.*" \
  --limit 10
```

### 2. Copy - Copy/retag images between registries

Copy an image from source to destination registry without using local disk storage. Supports different credentials for source and destination.

```sh
repo-lister copy \
  --source <source-image:tag> \
  --destination <dest-image:tag> \
  --source-secret <secret> \
  --dest-secret <secret> \
  --source-namespace <namespace> \
  --dest-namespace <namespace> \
  --progress
```

**Flags:**
- `-s, --source` - Source image reference (required)
- `-d, --destination` - Destination image reference (required)
- `--source-secret` - Secret for source registry (required)
- `--dest-secret` - Secret for destination registry (required)
- `--source-namespace` - Namespace for source secret (default: "default")
- `--dest-namespace` - Namespace for destination secret (default: "default")
- `-p, --progress` - Show progress during copy operation

**Examples:**

```sh
# Copy with new tag in same registry
repo-lister copy \
  --source myregistry.io/app:v1.0.0 \
  --destination myregistry.io/app:v2.0.0 \
  --source-secret regcred \
  --dest-secret regcred

# Copy between different registries with progress
repo-lister copy \
  --source gcr.io/project/app:latest \
  --destination registry.io/team/app:latest \
  --source-secret gcr-secret \
  --dest-secret registry-secret \
  --progress

# Copy with different namespaces
repo-lister copy \
  --source linuxarpan/testpush:v1.0.0 \
  --destination myregistry.io/testpush:v1.0.0 \
  --source-secret dockerhub-cred \
  --dest-secret registry-cred \
  --source-namespace kube-system \
  --dest-namespace default
```

### 3. Pull - Pull image to local storage

Pull a container image from a registry and save it to a local tar file.

```sh
repo-lister pull \
  --image <image:tag> \
  --output <path.tar> \
  --secret <secret> \
  --namespace <namespace>
```

**Flags:**
- `-i, --image` - Image reference to pull (required)
- `-o, --output` - Output path for tar file (required)
- `-s, --secret` - Kubernetes secret for authentication (required)
- `-n, --namespace` - Namespace where secret is located (default: "default")

**Examples:**

```sh
# Pull image to tar file
repo-lister pull \
  --image linuxarpan/testpush:v1.0.0 \
  --output /tmp/test-image.tar \
  --secret regcred \
  --namespace default

# Pull from private registry
repo-lister pull \
  --image myregistry.io/app:latest \
  --output ./backup/app-latest.tar \
  --secret registry-cred
```

### 4. Push - Push image from local storage

Push a container image from a local tar file to a registry.

```sh
repo-lister push \
  --image <image:tag> \
  --source <path.tar> \
  --secret <secret> \
  --namespace <namespace>
```

**Flags:**
- `-i, --image` - Destination image reference (required)
- `-f, --source` - Source tar file path (required)
- `-s, --secret` - Kubernetes secret for authentication (required)
- `-n, --namespace` - Namespace where secret is located (default: "default")

**Examples:**

```sh
# Push image from tar file
repo-lister push \
  --image linuxarpan/testpush:v2.0.0 \
  --source /tmp/test-image.tar \
  --secret regcred \
  --namespace default

# Push to private registry
repo-lister push \
  --image myregistry.io/app:latest \
  --source ./backup/app-latest.tar \
  --secret registry-cred
```

## Common Workflows

### Workflow 1: Retag an image in the same registry

```sh
repo-lister copy \
  --source myregistry.io/app:latest \
  --destination myregistry.io/app:v1.0.0 \
  --source-secret regcred \
  --dest-secret regcred
```

### Workflow 2: Migrate image between registries

```sh
repo-lister copy \
  --source gcr.io/project/app:latest \
  --destination myregistry.io/app:latest \
  --source-secret gcr-cred \
  --dest-secret registry-cred \
  --progress
```

### Workflow 3: Backup and restore images

```sh
# Backup: Pull image to local tar
repo-lister pull \
  --image myregistry.io/app:v1.0.0 \
  --output ./backup/app-v1.0.0.tar \
  --secret regcred

# Restore: Push image from local tar
repo-lister push \
  --image myregistry.io/app:v1.0.0 \
  --source ./backup/app-v1.0.0.tar \
  --secret regcred
```

### Workflow 4: Air-gapped environment transfer

```sh
# On connected machine: Pull image
repo-lister pull \
  --image gcr.io/project/app:latest \
  --output app-latest.tar \
  --secret gcr-cred

# Transfer app-latest.tar to air-gapped environment

# On air-gapped machine: Push to local registry
repo-lister push \
  --image local-registry.io/app:latest \
  --source app-latest.tar \
  --secret local-cred
```

## Creating Kubernetes Secrets

To use repo-lister, you need Kubernetes secrets with registry credentials:

### Docker Hub

```sh
kubectl create secret docker-registry dockerhub-cred \
  --docker-server=docker.io \
  --docker-username=<username> \
  --docker-password=<password> \
  --namespace=default
```

### Private Registry

```sh
kubectl create secret docker-registry registry-cred \
  --docker-server=myregistry.io \
  --docker-username=<username> \
  --docker-password=<password> \
  --namespace=default
```

### Google Container Registry (GCR)

```sh
kubectl create secret docker-registry gcr-cred \
  --docker-server=gcr.io \
  --docker-username=_json_key \
  --docker-password="$(cat key.json)" \
  --namespace=default
```

## Troubleshooting

### Authentication errors

- Verify the secret exists: `kubectl get secret <secret-name> -n <namespace>`
- Verify secret type: `kubectl get secret <secret-name> -n <namespace> -o yaml`
- Ensure secret is type `kubernetes.io/dockerconfigjson`
- Check RBAC permissions to read secrets

### Image not found

- Verify image name and tag are correct
- Check registry URL format (include registry domain for non-Docker Hub images)
- Ensure credentials have access to the image

### Network issues

- Verify connectivity to registry from cluster
- Check if registry requires VPN or special network configuration
- Verify firewall rules allow registry access

## License

See LICENSE file for details.
