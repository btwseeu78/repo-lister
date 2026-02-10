package utility

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// PushImage pushes an image from a local tar file to a registry
func PushImage(imageRef string, sourcePath string, secretName string, namespace string) error {
	// Create keychain
	kc, err := CreateKeychain(namespace, secretName)
	if err != nil {
		return fmt.Errorf("failed to create keychain: %w", err)
	}

	// Parse image reference
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("failed to parse image reference '%s': %w", imageRef, err)
	}

	fmt.Printf("Loading image from %s...\n", sourcePath)

	// Load image from tar file
	img, err := tarball.ImageFromPath(sourcePath, nil)
	if err != nil {
		return fmt.Errorf("failed to load image from tar file: %w", err)
	}

	fmt.Printf("Pushing image to %s...\n", imageRef)

	// Push image to registry
	err = remote.Write(ref, img, remote.WithAuthFromKeychain(kc))
	if err != nil {
		return fmt.Errorf("failed to push image to registry: %w", err)
	}

	fmt.Printf("✓ Successfully pushed image to %s\n", imageRef)
	return nil
}
