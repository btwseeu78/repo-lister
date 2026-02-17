package utility

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// normalizeImageName ensures the image name is a valid registry repository reference.
// For Docker Hub short names (e.g. "linuxarpan/testpush"), it prepends "docker.io/".
func normalizeImageName(imageName string) string {
	// If the image already contains a registry domain (has a dot or localhost), return as-is
	parts := strings.SplitN(imageName, "/", 2)
	if len(parts) == 2 && (strings.Contains(parts[0], ".") || strings.Contains(parts[0], ":") || parts[0] == "localhost") {
		return imageName
	}
	// Docker Hub short name — prepend docker.io
	return "docker.io/" + imageName
}

// ListImage lists tags from a container registry, with optional filtering and sorting by semver.
// secretName is optional — if empty, anonymous/public access is used.
func ListImage(imageName string, imageFilter string, secretName string, namespace string, limit int) ([]string, error) {

	// Create keychain using shared authentication (anonymous if no secret)
	kc, err := CreateKeychain(namespace, secretName)
	if err != nil {
		return nil, fmt.Errorf("error creating keychain: %w", err)
	}

	// Normalize the image name for proper registry resolution
	repoName := normalizeImageName(imageName)

	// Parse the repository name
	repo, err := name.NewRepository(repoName)
	if err != nil {
		return nil, fmt.Errorf("error parsing repository name '%s': %w", repoName, err)
	}

	// List all tags in the repository
	var filteredTags []string
	tags, err := remote.List(repo, remote.WithAuthFromKeychain(kc))
	if err != nil {
		return nil, HandleRegistryError(err, "listing tags for", repoName)
	}

	// Handle empty repository case
	if len(tags) == 0 {
		return nil, fmt.Errorf("repository '%s' is empty (no tags found)", repoName)
	}

	// Apply regex filter if provided
	if imageFilter != "" {
		regex, err := regexp.Compile(imageFilter)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex '%s': %w", imageFilter, err)
		}
		for _, tag := range tags {
			if regex.MatchString(tag) {
				filteredTags = append(filteredTags, tag)
			}
		}
	} else {
		filteredTags = tags
	}

	// Separate semver and non-semver tags
	var semverTags []semver.Version
	var nonSemverTags []string
	for _, tag := range filteredTags {
		v, err := semver.ParseTolerant(tag)
		if err == nil {
			semverTags = append(semverTags, v)
		} else {
			nonSemverTags = append(nonSemverTags, tag)
		}
	}

	// Sort semver tags descending (newest first)
	sort.Slice(semverTags, func(i, j int) bool {
		return semverTags[i].GT(semverTags[j])
	})

	// Convert sorted semver tags back to string, then append non-semver tags
	sortedTags := make([]string, 0, len(semverTags)+len(nonSemverTags))
	for _, v := range semverTags {
		sortedTags = append(sortedTags, v.String())
	}
	sortedTags = append(sortedTags, nonSemverTags...)

	if limit > 0 && limit < len(sortedTags) {
		sortedTags = sortedTags[:limit]
	}
	return sortedTags, nil
}
