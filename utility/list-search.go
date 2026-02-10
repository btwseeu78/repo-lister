package utility

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/blang/semver"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// ListImage lists tags from a container registry, with optional filtering and sorting by semver.
func ListImage(imageName string, imageFilter string, secretName string, namespace string, limit int) ([]string, error) {

	// Create keychain using shared authentication
	kc, err := CreateKeychain(namespace, secretName)
	if err != nil {
		fmt.Println("Error creating keychain:", err)
		return nil, err
	}

	// Define the repository to list tags from
	repoName := imageName

	// Parse the repository name
	repo, err := name.NewRepository(repoName)
	if err != nil {
		fmt.Println("Error parsing repository name:", err)
		return nil, err
	}
	var filteredTags []string
	// List all tags in the repository
	tags, err := remote.List(repo, remote.WithAuthFromKeychain(kc))
	if err != nil {
		fmt.Println("Error listing tags:", err)
		return nil, err
	}
	var regex *regexp.Regexp
	if imageFilter != "" {
		regex, err = regexp.Compile(imageFilter)
		if err != nil {
			fmt.Println("Error compiling regex:", err)
			return nil, err
		}

		// Filter tags using the regex pattern

		for _, tag := range tags {
			if regex.MatchString(tag) {
				filteredTags = append(filteredTags, tag)
			}
		}
	} else {
		filteredTags = tags
	}

	var semverTags []semver.Version
	for _, tag := range filteredTags {
		v, err := semver.ParseTolerant(tag)
		if err == nil {
			semverTags = append(semverTags, v)
		}
	}

	sort.Slice(semverTags, func(i, j int) bool {
		return semverTags[i].GT(semverTags[j])
	})

	// Convert sorted semver tags back to string
	sortedTags := make([]string, len(semverTags))
	for i, v := range semverTags {
		sortedTags[i] = v.String()
	}
	if limit > 0 && limit < len(sortedTags) {
		sortedTags = sortedTags[:limit]
	}
	return sortedTags, nil

}
