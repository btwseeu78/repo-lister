package utility

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateK8sClient creates a Kubernetes client
// It first tries to use in-cluster config, then falls back to kubeconfig
func CreateK8sClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kube config
		kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes client config: %w", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return clientset, nil
}

// CreateKeychain creates a k8schain keychain for the given namespace and secret
// This keychain can be used to authenticate with container registries using k8s secrets
func CreateKeychain(namespace, secretName string) (authn.Keychain, error) {
	clientset, err := CreateK8sClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	kc, err := k8schain.New(ctx, clientset, k8schain.Options{
		Namespace:        namespace,
		ImagePullSecrets: []string{secretName},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes keychain: %w", err)
	}

	return kc, nil
}
