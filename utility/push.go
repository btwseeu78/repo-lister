package utility

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

var createKeychainFn = CreateKeychain
var remoteWriteFn = remote.Write
var loadLocalImageFn = loadLocalImage

// PushImage pushes an image from the local Docker daemon to a registry.
func PushImage(sourceImageRef string, destinationImageRef string, secretName string, namespace string) error {
	// Parse source image reference
	srcRef, err := name.ParseReference(sourceImageRef)
	if err != nil {
		return fmt.Errorf("failed to parse source image reference '%s': %w", sourceImageRef, err)
	}

	fmt.Printf("Loading local image %s from Docker daemon...\n", sourceImageRef)

	img, err := loadLocalImageFn(srcRef)
	if err != nil {
		return fmt.Errorf("local image '%s' not found in Docker daemon. Tag the image first, then retry: %w", sourceImageRef, err)
	}

	// Parse destination image reference
	dstRef, err := name.ParseReference(destinationImageRef)
	if err != nil {
		return fmt.Errorf("failed to parse destination image reference '%s': %w", destinationImageRef, err)
	}

	// Create keychain after source image resolution so local-source errors fail first.
	kc, err := createKeychainFn(namespace, secretName)
	if err != nil {
		return fmt.Errorf("failed to create keychain: %w", err)
	}

	fmt.Printf("Pushing image to %s...\n", destinationImageRef)

	// Push image to registry
	err = remoteWriteFn(dstRef, img, remote.WithAuthFromKeychain(kc))
	if err != nil {
		return HandleRegistryError(err, "pushing image to", destinationImageRef)
	}

	fmt.Printf("✓ Successfully pushed image to %s\n", destinationImageRef)
	return nil
}

func loadLocalImage(src name.Reference) (v1.Image, error) {
	return daemon.Image(src)
}
