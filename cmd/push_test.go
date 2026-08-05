package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func executePush(t *testing.T, args ...string) error {
	t.Helper()

	pushSourceImage = ""
	pushDestinationImage = ""
	pushSecret = ""
	pushNamespace = "default"

	rootCmd.SetArgs(append([]string{"push"}, args...))
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	return rootCmd.Execute()
}

func TestPushCommandRequiresSourceAndDestinationImage(t *testing.T) {
	err := executePush(t, "--destination-image", "registry.io/team/app:1.0.0")
	if err == nil {
		t.Fatalf("expected missing source-image error")
	}

	if !strings.Contains(err.Error(), `required flag(s) "source-image" not set`) {
		t.Fatalf("expected missing source-image error, got: %v", err)
	}
}

func TestPushCommandWiresSourceDestinationSecretAndNamespace(t *testing.T) {
	originalPushImageFn := pushImageFn
	t.Cleanup(func() {
		pushImageFn = originalPushImageFn
	})

	var capturedSource string
	var capturedDestination string
	var capturedSecret string
	var capturedNamespace string

	pushImageFn = func(sourceImage, destinationImage, secretName, namespace string) error {
		capturedSource = sourceImage
		capturedDestination = destinationImage
		capturedSecret = secretName
		capturedNamespace = namespace
		return nil
	}

	err := executePush(t,
		"--source-image", "docker.io/library/nginx:1.25.0",
		"--destination-image", "registry.io/team/nginx:1.25.0",
		"--secret", "regcred",
		"--namespace", "images",
	)
	if err != nil {
		t.Fatalf("expected command execution to succeed, got: %v", err)
	}

	if capturedSource != "docker.io/library/nginx:1.25.0" {
		t.Fatalf("unexpected source image: %s", capturedSource)
	}

	if capturedDestination != "registry.io/team/nginx:1.25.0" {
		t.Fatalf("unexpected destination image: %s", capturedDestination)
	}

	if capturedSecret != "regcred" {
		t.Fatalf("unexpected secret name: %s", capturedSecret)
	}

	if capturedNamespace != "images" {
		t.Fatalf("unexpected namespace: %s", capturedNamespace)
	}
}

func TestPushCommandDefaultsNamespace(t *testing.T) {
	originalPushImageFn := pushImageFn
	t.Cleanup(func() {
		pushImageFn = originalPushImageFn
	})

	var capturedNamespace string

	pushImageFn = func(_, _, _, namespace string) error {
		capturedNamespace = namespace
		return nil
	}

	err := executePush(t,
		"--source-image", "docker.io/library/busybox:1.36",
		"--destination-image", "registry.io/team/busybox:1.36",
	)
	if err != nil {
		t.Fatalf("expected command execution to succeed, got: %v", err)
	}

	if capturedNamespace != "default" {
		t.Fatalf("expected default namespace to be 'default', got: %s", capturedNamespace)
	}
}

func TestPushCommandReturnsRunError(t *testing.T) {
	originalPushImageFn := pushImageFn
	t.Cleanup(func() {
		pushImageFn = originalPushImageFn
	})

	pushImageFn = func(_, _, _, _ string) error {
		return fmt.Errorf("push failed")
	}

	err := executePush(t,
		"--source-image", "docker.io/library/redis:7",
		"--destination-image", "registry.io/team/redis:7",
	)
	if err == nil {
		t.Fatalf("expected command execution to fail")
	}

	if !strings.Contains(err.Error(), "push failed") {
		t.Fatalf("expected push error to bubble up, got: %v", err)
	}
}
