package utility

import (
	"context"
	"fmt"
	"regexp"
	"sort"

	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ListImage(imageName string, imageFilter string, secretName string, namespace string, limit int) ([]string, error) {

	// Create KeycHain to be used with docker
	// Create a Kubernetes client
	var config *rest.Config
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Println("Error creating Kubernetes client config:", err)
			return nil, err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		return nil, err
	}

	// Create a Kubernetes keychain using the specific secret
	ctx := context.Background()
	kc, err := k8schain.New(ctx, clientset, k8schain.Options{
		Namespace:        namespace, // Change this to the namespace where your secret is located
		ImagePullSecrets: []string{secretName},
	})
	if err != nil {
		fmt.Println("Error creating Kubernetes keychain:", err)
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
	fmt.Println("limt", limit)
	// Sort tags in descending order
	sort.Sort(sort.Reverse(sort.StringSlice(filteredTags)))
	if limit > 0 && limit < len(filteredTags) {
		filteredTags = filteredTags[:limit]
	}
	return filteredTags, nil

}
