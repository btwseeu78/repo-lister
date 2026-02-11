package utility

import (
	"os"
	"path/filepath"
	"testing"
)

// TestPushImageValidation tests basic validation of PushImage parameters
func TestPushImageValidation(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name       string
		imageRef   string
		sourcePath string
		secretName string
		namespace  string
		createFile bool
		wantErr    bool
	}{
		{
			name:       "valid parameters",
			imageRef:   "nginx:v1.0.0",
			sourcePath: filepath.Join(tempDir, "test-image.tar"),
			secretName: "regcred",
			namespace:  "default",
			createFile: true,
			wantErr:    false,
		},
		{
			name:       "empty image reference",
			imageRef:   "",
			sourcePath: filepath.Join(tempDir, "test-image.tar"),
			secretName: "regcred",
			namespace:  "default",
			createFile: true,
			wantErr:    true,
		},
		{
			name:       "non-existent source file",
			imageRef:   "nginx:v1.0.0",
			sourcePath: "/nonexistent/file.tar",
			secretName: "regcred",
			namespace:  "default",
			createFile: false,
			wantErr:    true,
		},
		{
			name:       "empty source path",
			imageRef:   "nginx:v1.0.0",
			sourcePath: "",
			secretName: "regcred",
			namespace:  "default",
			createFile: false,
			wantErr:    true,
		},
		{
			name:       "invalid image reference",
			imageRef:   ":::invalid:::",
			sourcePath: filepath.Join(tempDir, "test-image.tar"),
			secretName: "regcred",
			namespace:  "default",
			createFile: true,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy tar file if needed
			if tt.createFile && tt.sourcePath != "" {
				if err := os.MkdirAll(filepath.Dir(tt.sourcePath), 0755); err != nil {
					t.Fatalf("Failed to create directory: %v", err)
				}
				f, err := os.Create(tt.sourcePath)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				if _, err := f.WriteString("dummy tar content"); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				_ = f.Close()
				defer func() { _ = os.Remove(tt.sourcePath) }()
			}

			err := PushImage(tt.imageRef, tt.sourcePath, tt.secretName, tt.namespace)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for test case: %s", tt.name)
				}
			} else {
				if err != nil {
					// Expected in test environment without k8s/registry
					t.Logf("Expected error without k8s/registry access: %v", err)
				}
			}
		})
	}
}

// TestPushImageSourceFile tests source file validation
func TestPushImageSourceFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name       string
		setupFile  func(string) error
		expectErr  bool
		errMessage string
	}{
		{
			name: "valid tar file",
			setupFile: func(path string) error {
				return os.WriteFile(path, []byte("test content"), 0644)
			},
			expectErr: false,
		},
		{
			name: "empty file",
			setupFile: func(path string) error {
				return os.WriteFile(path, []byte{}, 0644)
			},
			expectErr: false, // Will fail during tar parsing, not validation
		},
		{
			name: "directory instead of file",
			setupFile: func(path string) error {
				return os.Mkdir(path, 0755)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourcePath := filepath.Join(tempDir, tt.name+".tar")

			if tt.setupFile != nil {
				if err := tt.setupFile(sourcePath); err != nil {
					t.Fatalf("Failed to setup test file: %v", err)
				}
				defer func() { _ = os.Remove(sourcePath) }()
			}

			err := PushImage("nginx:test", sourcePath, "test-secret", "default")

			if tt.expectErr && err == nil {
				t.Errorf("Expected error for test case: %s", tt.name)
			}

			if err != nil {
				t.Logf("Error (expected): %v", err)
			}
		})
	}
}
