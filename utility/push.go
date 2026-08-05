package utility

import (
	"fmt"
	"os"
	"strings"

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
	if strings.TrimSpace(sourceImageRef) == "" {
		return fmt.Errorf("source image is required")
	}

	if strings.TrimSpace(destinationImageRef) == "" {
		return fmt.Errorf("destination image is required")
	}

	// Parse source image reference
	srcRef, err := name.ParseReference(sourceImageRef)
	if err != nil {
		return fmt.Errorf("failed to parse source image reference '%s': %w", sourceImageRef, err)
	}

	// Parse destination image reference before touching local daemon state.
	dstRef, err := name.ParseReference(destinationImageRef)
	if err != nil {
		return fmt.Errorf("failed to parse destination image reference '%s': %w", destinationImageRef, err)
	}

	fmt.Printf("Loading local image %s from Docker daemon...\n", sourceImageRef)

	img, err := loadLocalImageFn(srcRef)
	if err != nil {
		return fmt.Errorf("local image '%s' not found in Docker daemon. Tag the image first, then retry: %w", sourceImageRef, err)
	}

	// Create keychain after source image resolution so local-source errors fail first.
	kc, err := createKeychainFn(namespace, secretName)
	if err != nil {
		return fmt.Errorf("failed to create keychain: %w", err)
	}

	reporter := newPushProgressReporter(os.Stdout, isInteractiveTerminal())
	reporter.Start(100)
	reporter.Update(20, "preparing image")
	reporter.Update(40, "resolving auth")
	reporter.Update(70, "pushing layers")

	// Push image to registry
	err = remoteWriteFn(dstRef, img, remote.WithAuthFromKeychain(kc))
	if err != nil {
		mappedErr := HandleRegistryError(err, "pushing image to", destinationImageRef)
		reporter.Fail(mappedErr)
		return mappedErr
	}

	reporter.Update(100, "finalizing")
	reporter.Complete(fmt.Sprintf("Successfully pushed image to %s", destinationImageRef))
	return nil
}

func loadLocalImage(src name.Reference) (v1.Image, error) {
	return daemon.Image(src)
}
