package utility

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// PullImage pulls an image from a registry and saves it to a local tar file
func PullImage(imageRef string, outputPath string, secretName string, namespace string) error {
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

	fmt.Printf("Pulling image %s...\n", imageRef)

	// Fetch the image from registry
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(kc))
	if err != nil {
		return fmt.Errorf("failed to fetch image from registry: %w", err)
	}

	fmt.Printf("Saving image to %s...\n", outputPath)

	// Write image to tar file
	err = tarball.WriteToFile(outputPath, ref, img)
	if err != nil {
		return fmt.Errorf("failed to write image to tar file: %w", err)
	}

	fmt.Printf("✓ Successfully pulled image to %s\n", outputPath)
	return nil
}
