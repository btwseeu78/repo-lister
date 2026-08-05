package utility

import (
	"strings"
	"testing"
)

func TestPushImageFailsWhenSourceReferenceInvalid(t *testing.T) {
	err := PushImage(":::invalid:::", "registry.io/team/app:1.0.0", "", "default")
	if err == nil || !strings.Contains(err.Error(), "failed to parse source image reference") {
		t.Fatalf("expected source parse error, got %v", err)
	}
}

func TestPushImageFailsWhenLocalSourceMissing(t *testing.T) {
	err := PushImage("local/not-found:dev", "registry.io/team/app:1.0.0", "", "default")
	if err == nil || !strings.Contains(err.Error(), "local image") {
		t.Fatalf("expected local image missing error, got %v", err)
	}
}

func TestPushImageLocalSourceFailFastBeforeKeychainCreation(t *testing.T) {
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
}
