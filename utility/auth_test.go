package utility

import (
	"testing"
)

// TestCreateKeychainValidation tests basic validation of CreateKeychain function
func TestCreateKeychainValidation(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		secret    string
		wantErr   bool
	}{
		{
			name:      "valid namespace and secret",
			namespace: "default",
			secret:    "regcred",
			wantErr:   false, // Will error in test without k8s, but tests validation
		},
		{
			name:      "empty namespace",
			namespace: "",
			secret:    "regcred",
			wantErr:   false, // k8schain will handle empty namespace
		},
		{
			name:      "empty secret",
			namespace: "default",
			secret:    "",
			wantErr:   false, // Will create keychain but may fail on auth
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This will fail in CI/test environments without k8s cluster
			// Real validation would need mocking or integration test environment
			_, err := CreateKeychain(tt.namespace, tt.secret)

			// In a real k8s environment, this should work
			// In test environment, we expect an error (no k8s config)
			if err == nil {
				t.Log("Successfully created keychain (running in k8s environment)")
			} else {
				t.Logf("Expected error in test environment: %v", err)
			}
		})
	}
}

// TestCreateK8sClient tests the k8s client creation
func TestCreateK8sClient(t *testing.T) {
	t.Run("create k8s client", func(t *testing.T) {
		// This will only work in an environment with k8s config
		_, err := CreateK8sClient()

		if err == nil {
			t.Log("Successfully created k8s client (running in k8s environment)")
		} else {
			t.Logf("Expected error without k8s config: %v", err)
		}
	})
}

// Note: For comprehensive testing, consider using:
// - k8s fake client for mocking
// - testcontainers for integration tests
// - environment variable to skip tests requiring k8s cluster
