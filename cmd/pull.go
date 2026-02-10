package cmd

import (
	"os"
	"repo-lister/utility"

	"github.com/spf13/cobra"
)

var (
	pullImage     string
	pullOutput    string
	pullSecret    string
	pullNamespace string
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull an image from registry to local storage",
	Long: `Pull a container image from a registry and save it to a local tar file.

This command downloads an image from a registry using Kubernetes credentials
and saves it as a tar archive. The tar file can later be used with the push
command to upload to a different registry, or can be loaded into a local
Docker daemon.

The pull operation is useful for:
  - Backing up images locally
  - Transferring images to air-gapped environments
  - Inspecting image contents offline
  - Migrating images between registries (using pull + push)`,
	Example: `  # Pull an image to a tar file
  repo-lister pull \
    --image linuxarpan/testpush:v1.0.0 \
    --output /tmp/my-image.tar \
    --secret regcred \
    --namespace default

  # Pull from a private registry
  repo-lister pull \
    --image myregistry.io/app:latest \
    --output ./backup/app-latest.tar \
    --secret registry-cred`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the PullImage function from the utility package
		err := utility.PullImage(pullImage, pullOutput, pullSecret, pullNamespace)
		if err != nil {
			cmd.PrintErrln("Error pulling image:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Define flags for the pull command
	pullCmd.Flags().StringVarP(&pullImage, "image", "i", "", "Image reference to pull (e.g., registry.io/image:tag) (required)")
	pullCmd.Flags().StringVarP(&pullOutput, "output", "o", "", "Output path for the tar file (required)")
	pullCmd.Flags().StringVarP(&pullSecret, "secret", "s", "", "Kubernetes secret name for registry authentication (required)")
	pullCmd.Flags().StringVarP(&pullNamespace, "namespace", "n", "default", "Kubernetes namespace where the secret is located")

	// Mark required flags
	pullCmd.MarkFlagRequired("image")
	pullCmd.MarkFlagRequired("output")
	pullCmd.MarkFlagRequired("secret")
}
