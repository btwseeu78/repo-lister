## Usage
```sh
go run main.go -i linuxarpan/testpush -s regcred -l 2
```
## Local Uses
 build the tool with go build ., copt to your local bin.

## Requirements

- Connectivity to cluster
- kubeconfig
- Access to get secret

## Flags
```sh
    --imageName (string): Name of the image.
    --imageFilter (string): Filter for the image.
    --secretName (string): Name of the secret.
    --limit (int): Limit for the number of results.
```