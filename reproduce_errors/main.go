package main

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func main() {
	tests := []struct {
		name string
		repo string
	}{
		{"Valid public repo", "index.docker.io/library/alpine"},
		{"Non-existent repo", "index.docker.io/library/this-repo-should-not-exist-12345"},
		{"Invalid host", "nonexistent.registry.local/repo"},
		{"Private repo (unauth)", "gcr.io/google-containers/pause"}, // Might work if public, let's try a clearly private one or just interpret the result
	}

	for _, tt := range tests {
		fmt.Printf("Testing: %s (%s)\n", tt.name, tt.repo)
		ref, err := name.NewRepository(tt.repo)
		if err != nil {
			fmt.Printf("  Error parsing repo: %v\n", err)
			continue
		}

		tags, err := remote.List(ref)
		if err != nil {
			fmt.Printf("  Error listing tags: %v\n", err)
		} else {
			fmt.Printf("  Success! Found %d tags\n", len(tags))
		}
		fmt.Println("---------------------------------------------------")
	}
}
