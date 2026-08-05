package utility

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func TestPushImageFailsWhenSourceReferenceInvalid(t *testing.T) {
	err := PushImage(":::invalid:::", "registry.io/team/app:1.0.0", "", "default")
	if err == nil || !strings.Contains(err.Error(), "failed to parse source image reference") {
		t.Fatalf("expected source parse error, got %v", err)
	}
}

func TestPushImageBehaviorMatrix(t *testing.T) {
	originalCreateKeychainFn := createKeychainFn
	originalRemoteWriteFn := remoteWriteFn
	originalLoadLocalImageFn := loadLocalImageFn
	t.Cleanup(func() {
		createKeychainFn = originalCreateKeychainFn
		remoteWriteFn = originalRemoteWriteFn
		loadLocalImageFn = originalLoadLocalImageFn
	})

	cases := []struct {
		name                 string
		source               string
		destination          string
		secret               string
		loadErr              error
		wantErrSubstr        string
		wantLoadCalls        int
		wantKeychainCalls    int
		wantRemoteWriteCalls int
	}{
		{name: "missing source", source: "", destination: "registry.io/team/app:1.0.0", secret: "", wantErrSubstr: "source image is required", wantLoadCalls: 0, wantKeychainCalls: 0, wantRemoteWriteCalls: 0},
		{name: "missing destination", source: "local/app:dev", destination: "", secret: "", wantErrSubstr: "destination image is required", wantLoadCalls: 0, wantKeychainCalls: 0, wantRemoteWriteCalls: 0},
		{name: "invalid destination", source: "local/app:dev", destination: ":::bad:::", secret: "", wantErrSubstr: "failed to parse destination image reference", wantLoadCalls: 1, wantKeychainCalls: 0, wantRemoteWriteCalls: 0},
		{name: "local source missing", source: "local/not-found:dev", destination: "registry.io/team/app:1.0.0", secret: "", loadErr: errors.New("daemon image not found"), wantErrSubstr: "local image", wantLoadCalls: 1, wantKeychainCalls: 0, wantRemoteWriteCalls: 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			loadCalls := 0
			keychainCalls := 0
			remoteWriteCalls := 0

			loadLocalImageFn = func(_ name.Reference) (v1.Image, error) {
				loadCalls++
				if tc.loadErr != nil {
					return nil, tc.loadErr
				}
				return empty.Image, nil
			}

			createKeychainFn = func(_, _ string) (authn.Keychain, error) {
				keychainCalls++
				return authn.DefaultKeychain, nil
			}

			remoteWriteFn = func(_ name.Reference, _ v1.Image, _ ...remote.Option) error {
				remoteWriteCalls++
				return nil
			}

			err := PushImage(tc.source, tc.destination, tc.secret, "default")
			if err == nil {
				t.Fatalf("expected error")
			}

			if !strings.Contains(err.Error(), tc.wantErrSubstr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErrSubstr, err)
			}

			if loadCalls != tc.wantLoadCalls {
				t.Fatalf("expected loadLocalImageFn calls %d, got %d", tc.wantLoadCalls, loadCalls)
			}

			if keychainCalls != tc.wantKeychainCalls {
				t.Fatalf("expected createKeychainFn calls %d, got %d", tc.wantKeychainCalls, keychainCalls)
			}

			if remoteWriteCalls != tc.wantRemoteWriteCalls {
				t.Fatalf("expected remoteWriteFn calls %d, got %d", tc.wantRemoteWriteCalls, remoteWriteCalls)
			}
		})
	}
}

func TestPushImageLocalSourceFailFastBeforeKeychainCreation(t *testing.T) {
	originalCreateKeychainFn := createKeychainFn
	originalRemoteWriteFn := remoteWriteFn
	originalLoadLocalImageFn := loadLocalImageFn
	t.Cleanup(func() {
		createKeychainFn = originalCreateKeychainFn
		remoteWriteFn = originalRemoteWriteFn
		loadLocalImageFn = originalLoadLocalImageFn
	})

	keychainCalled := false
	remoteWriteCalled := false

	loadLocalImageFn = func(_ name.Reference) (v1.Image, error) {
		return nil, errors.New("daemon image not found")
	}

	createKeychainFn = func(_, _ string) (authn.Keychain, error) {
		keychainCalled = true
		return authn.DefaultKeychain, nil
	}

	remoteWriteFn = func(_ name.Reference, _ v1.Image, _ ...remote.Option) error {
		remoteWriteCalled = true
		return nil
	}

	err := PushImage("local/not-found:dev", "registry.io/team/app:1.0.0", "regcred", "default")
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), "local image") {
		t.Fatalf("expected local image error before keychain creation, got %v", err)
	}

	if strings.Contains(err.Error(), "failed to create keychain") {
		t.Fatalf("expected local-source fail-fast before keychain creation, got %v", err)
	}

	if keychainCalled {
		t.Fatalf("expected keychain creation not to be called on local-source failure")
	}

	if remoteWriteCalled {
		t.Fatalf("expected remote write not to be called on local-source failure")
	}
}

func TestPushImageUsesDefaultKeychainWhenSecretEmpty(t *testing.T) {
	originalCreateKeychainFn := createKeychainFn
	originalRemoteWriteFn := remoteWriteFn
	originalLoadLocalImageFn := loadLocalImageFn
	t.Cleanup(func() {
		createKeychainFn = originalCreateKeychainFn
		remoteWriteFn = originalRemoteWriteFn
		loadLocalImageFn = originalLoadLocalImageFn
	})

	createCalled := false
	writeCalled := false

	createKeychainFn = func(namespace, secretName string) (authn.Keychain, error) {
		createCalled = true
		if namespace != "default" {
			t.Fatalf("expected namespace default, got %q", namespace)
		}
		if secretName != "" {
			t.Fatalf("expected empty secret name for default keychain path, got %q", secretName)
		}
		return authn.DefaultKeychain, nil
	}

	loadLocalImageFn = func(_ name.Reference) (v1.Image, error) {
		return empty.Image, nil
	}

	remoteWriteFn = func(_ name.Reference, _ v1.Image, _ ...remote.Option) error {
		writeCalled = true
		return nil
	}

	err := PushImage("local/app:dev", "registry.io/team/app:1.0.0", "", "default")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !createCalled {
		t.Fatalf("expected keychain factory to be called")
	}

	if !writeCalled {
		t.Fatalf("expected remote write to be called")
	}
}

func TestPushImageWrapsRegistryWriteError(t *testing.T) {
	originalCreateKeychainFn := createKeychainFn
	originalRemoteWriteFn := remoteWriteFn
	originalLoadLocalImageFn := loadLocalImageFn
	t.Cleanup(func() {
		createKeychainFn = originalCreateKeychainFn
		remoteWriteFn = originalRemoteWriteFn
		loadLocalImageFn = originalLoadLocalImageFn
	})

	createKeychainFn = func(_, _ string) (authn.Keychain, error) {
		return authn.DefaultKeychain, nil
	}

	loadLocalImageFn = func(_ name.Reference) (v1.Image, error) {
		return empty.Image, nil
	}

	remoteWriteFn = func(_ name.Reference, _ v1.Image, _ ...remote.Option) error {
		return errors.New("unauthorized: authentication required")
	}

	err := PushImage("local/app:dev", "registry.io/team/app:1.0.0", "", "default")
	if err == nil {
		t.Fatalf("expected error")
	}

	if !strings.Contains(err.Error(), "authentication failed while pushing image to") {
		t.Fatalf("expected registry auth error mapping, got %v", err)
	}
}
