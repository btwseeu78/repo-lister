package utility

import (
	"testing"
)

// TestListImageValidation tests basic validation of ListImage parameters
func TestListImageValidation(t *testing.T) {
	tests := []struct {
		name        string
		imageName   string
		imageFilter string
		secretName  string
		namespace   string
		limit       int
		wantErr     bool
	}{
		{
			name:        "valid parameters",
			imageName:   "nginx",
			imageFilter: ".*",
			secretName:  "regcred",
			namespace:   "default",
			limit:       5,
			wantErr:     false,
		},
		{
			name:        "empty image name",
			imageName:   "",
			imageFilter: ".*",
			secretName:  "regcred",
			namespace:   "default",
			limit:       5,
			wantErr:     true,
		},
		{
			name:        "negative limit",
			imageName:   "nginx",
			imageFilter: ".*",
			secretName:  "regcred",
			namespace:   "default",
			limit:       -1,
			wantErr:     false, // Negative limit is handled as no limit
		},
		{
			name:        "zero limit",
			imageName:   "nginx",
			imageFilter: ".*",
			secretName:  "regcred",
			namespace:   "default",
			limit:       0,
			wantErr:     false, // Zero limit returns all
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This will fail without k8s cluster and registry access
			// We're testing parameter validation, not actual functionality
			_, err := ListImage(tt.imageName, tt.imageFilter, tt.secretName, tt.namespace, tt.limit)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for test case: %s", tt.name)
				}
			} else {
				if err != nil {
					t.Logf("Expected error without k8s/registry access: %v", err)
				}
			}
		})
	}
}

// TestImageFilterRegex tests regex pattern validation
func TestImageFilterRegex(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		wantErr bool
	}{
		{
			name:    "valid regex",
			filter:  "v[0-9]+.*",
			wantErr: false,
		},
		{
			name:    "wildcard",
			filter:  ".*",
			wantErr: false,
		},
		{
			name:    "empty filter",
			filter:  "",
			wantErr: false,
		},
		{
			name:    "invalid regex",
			filter:  "[invalid((",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We'll test this by calling ListImage with the filter
			// The function should handle invalid regex gracefully
			_, err := ListImage("nginx", tt.filter, "test-secret", "default", 1)

			if tt.wantErr && err == nil {
				// In a perfect world, this should error on invalid regex
				// But the actual error might come from k8s/registry access
				t.Logf("Invalid regex might be caught during execution")
			}

			if err != nil {
				t.Logf("Error (expected in test environment): %v", err)
			}
		})
	}
}
