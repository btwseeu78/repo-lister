package utility

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// CopyImage copies an image from source to destination registry without local storage
// It supports different source and destination secrets for cross-registry copying
func CopyImage(
	sourceImage string,
	destImage string,
	sourceSecret string,
	destSecret string,
	sourceNamespace string,
	destNamespace string,
	showProgress bool,
) error {
	// Validate that source and destination are different
	if sourceImage == destImage {
		return fmt.Errorf("source and destination images are identical: %s", sourceImage)
	}

	// Create source keychain
	sourceKC, err := CreateKeychain(sourceNamespace, sourceSecret)
	if err != nil {
		return fmt.Errorf("failed to create source keychain: %w", err)
	}

	// Create destination keychain
	destKC, err := CreateKeychain(destNamespace, destSecret)
	if err != nil {
		return fmt.Errorf("failed to create destination keychain: %w", err)
	}

	// Parse source image reference
	srcRef, err := name.ParseReference(sourceImage)
	if err != nil {
		return fmt.Errorf("failed to parse source image reference '%s': %w", sourceImage, err)
	}

	// Parse destination image reference
	dstRef, err := name.ParseReference(destImage)
	if err != nil {
		return fmt.Errorf("failed to parse destination image reference '%s': %w", destImage, err)
	}

	if showProgress {
		fmt.Printf("Copying image...\n")
		fmt.Printf("  Source:      %s\n", sourceImage)
		fmt.Printf("  Destination: %s\n", destImage)
		fmt.Println()
	}

	// Fetch image descriptor from source
	if showProgress {
		fmt.Println("Fetching image from source registry...")
	}

	desc, err := remote.Get(srcRef, remote.WithAuthFromKeychain(sourceKC))
	if err != nil {
		return HandleRegistryError(err, "fetching image from source", sourceImage)
	}

	// Check if it's an image index (multi-arch) or regular image
	if showProgress {
		fmt.Println("Analyzing image type...")
	}

	// Create progress channel if needed
	var updates chan v1.Update
	if showProgress {
		updates = make(chan v1.Update, 100)
		go printProgress(updates)
		// Note: remote.Write/WriteIndex will close the channel
	}

	// Try to get as image
	img, err := desc.Image()
	if err == nil {
		// It's a regular image
		if showProgress {
			fmt.Println("Copying image layers to destination registry...")
			err = remote.Write(dstRef, img,
				remote.WithAuthFromKeychain(destKC),
				remote.WithProgress(updates))
		} else {
			err = remote.Write(dstRef, img, remote.WithAuthFromKeychain(destKC))
		}

		if err != nil {
			return HandleRegistryError(err, "writing image to destination", destImage)
		}
	} else {
		// Try as image index (multi-arch)
		idx, err := desc.ImageIndex()
		if err != nil {
			return fmt.Errorf("failed to process image (not a valid image or image index): %w", err)
		}

		if showProgress {
			fmt.Println("Copying image index (multi-arch) to destination registry...")
			err = remote.WriteIndex(dstRef, idx,
				remote.WithAuthFromKeychain(destKC),
				remote.WithProgress(updates))
		} else {
			err = remote.WriteIndex(dstRef, idx, remote.WithAuthFromKeychain(destKC))
		}

		if err != nil {
			return HandleRegistryError(err, "writing image index to destination", destImage)
		}
	}

	if showProgress {
		fmt.Printf("\n✓ Successfully copied image to %s\n", destImage)
	}

	return nil
}

// printProgress reads progress updates from a channel and prints them
func printProgress(updates <-chan v1.Update) {
	for update := range updates {
		if update.Total > 0 {
			percent := float64(update.Complete) / float64(update.Total) * 100
			fmt.Printf("\rProgress: %d/%d bytes (%.1f%%)", update.Complete, update.Total, percent)
		}
	}
}
