package cmd

import (
	"repo-lister/utility"

	"github.com/spf13/cobra"
)

var (
	pushSourceImage      string
	pushDestinationImage string
	pushSecret           string
	pushNamespace        string
	pushImageFn          = utility.PushImage
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push an image from local daemon to registry",
	Long: `Push a container image from your local daemon to a registry.

This command reads an image that already exists locally and pushes it to a
destination registry image reference using Kubernetes credentials when provided.

The push operation is useful for:
  - Publishing locally built images
  - Mirroring images to another registry
  - Promoting tagged images across environments`,
	Example: `  # Push an image from local daemon to a destination reference
  repo-lister push \
    --source-image linuxarpan/testpush:v2.0.0 \
    --destination-image registry.io/team/testpush:v2.0.0 \
    --secret regcred \
    --namespace default

  # Push to a private registry
  repo-lister push \
	    --source-image app:latest \
	    --destination-image myregistry.io/app:latest`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return pushImageFn(pushSourceImage, pushDestinationImage, pushSecret, pushNamespace)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Define flags for the push command
	pushCmd.Flags().StringVar(&pushSourceImage, "source-image", "", "Local daemon source image reference (required)")
	pushCmd.Flags().StringVar(&pushDestinationImage, "destination-image", "", "Destination registry image reference (required)")
	pushCmd.Flags().StringVarP(&pushSecret, "secret", "s", "", "Kubernetes secret name for registry authentication")
	pushCmd.Flags().StringVarP(&pushNamespace, "namespace", "n", "default", "Kubernetes namespace where the secret is located")

	// Mark required flags
	_ = pushCmd.MarkFlagRequired("source-image")
	_ = pushCmd.MarkFlagRequired("destination-image")
}
