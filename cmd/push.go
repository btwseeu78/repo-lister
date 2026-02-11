package cmd

import (
	"os"
	"repo-lister/utility"

	"github.com/spf13/cobra"
)

var (
	pushImage     string
	pushSource    string
	pushSecret    string
	pushNamespace string
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push an image from local storage to registry",
	Long: `Push a container image from a local tar file to a registry.

This command loads an image from a tar archive and uploads it to a registry
using Kubernetes credentials for authentication. The tar file is typically
created using the pull command or docker save.

The push operation is useful for:
  - Uploading locally modified images
  - Migrating images from pull command to another registry
  - Restoring backed up images
  - Publishing images to private registries`,
	Example: `  # Push an image from a tar file
  repo-lister push \
    --image linuxarpan/testpush:v2.0.0 \
    --source /tmp/my-image.tar \
    --secret regcred \
    --namespace default

  # Push to a private registry
  repo-lister push \
    --image myregistry.io/app:latest \
    --source ./backup/app-latest.tar \
    --secret registry-cred`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the PushImage function from the utility package
		err := utility.PushImage(pushImage, pushSource, pushSecret, pushNamespace)
		if err != nil {
			cmd.PrintErrln("Error pushing image:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	// Define flags for the push command
	pushCmd.Flags().StringVarP(&pushImage, "image", "i", "", "Destination image reference (e.g., registry.io/image:tag) (required)")
	pushCmd.Flags().StringVarP(&pushSource, "source", "f", "", "Source tar file path (required)")
	pushCmd.Flags().StringVarP(&pushSecret, "secret", "s", "", "Kubernetes secret name for registry authentication (required)")
	pushCmd.Flags().StringVarP(&pushNamespace, "namespace", "n", "default", "Kubernetes namespace where the secret is located")

	// Mark required flags
	_ = pushCmd.MarkFlagRequired("image")
	_ = pushCmd.MarkFlagRequired("source")
	_ = pushCmd.MarkFlagRequired("secret")
}
