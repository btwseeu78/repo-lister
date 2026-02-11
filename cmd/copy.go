package cmd

import (
	"os"
	"repo-lister/utility"

	"github.com/spf13/cobra"
)

var (
	copySource          string
	copyDestination     string
	copySourceSecret    string
	copyDestSecret      string
	copySourceNamespace string
	copyDestNamespace   string
	copyShowProgress    bool
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy/retag an image from one registry to another",
	Long: `Copy an image from a source registry to a destination registry without using local storage.

This command streams the image directly from source to destination, supporting:
  - Retagging within the same registry
  - Copying between different registries
  - Using different credentials for source and destination
  - Multi-architecture images (image indexes)

The copy operation is efficient as it doesn't require local disk storage for the image.`,
	Example: `  # Copy from public source to private destination
  repo-lister copy \
    --source docker.io/library/nginx:latest \
    --destination myregistry.io/nginx:latest \
    --dest-secret regcred

  # Copy image with different tag in same registry
  repo-lister copy \
    --source myregistry.io/app:v1.0.0 \
    --destination myregistry.io/app:v2.0.0 \
    --source-secret regcred \
    --dest-secret regcred

  # Copy between different registries with different credentials
  repo-lister copy \
    --source gcr.io/project/app:latest \
    --destination registry.io/team/app:latest \
    --source-secret gcr-secret \
    --dest-secret registry-secret \
    --source-namespace kube-system \
    --dest-namespace default \
    --progress

  # Copy with progress indicator
  repo-lister copy \
    --source linuxarpan/testpush:v1.0.0 \
    --destination linuxarpan/testpush:v2.0.0 \
    --source-secret regcred \
    --dest-secret regcred \
    --progress`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the CopyImage function from the utility package
		err := utility.CopyImage(
			copySource,
			copyDestination,
			copySourceSecret,
			copyDestSecret,
			copySourceNamespace,
			copyDestNamespace,
			copyShowProgress,
		)
		if err != nil {
			cmd.PrintErrln("Error copying image:", err)
			os.Exit(1)
		}

		if !copyShowProgress {
			cmd.Printf("Successfully copied %s to %s\n", copySource, copyDestination)
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// Define flags for the copy command
	copyCmd.Flags().StringVarP(&copySource, "source", "s", "", "Source image reference (e.g., registry.io/image:tag) (required)")
	copyCmd.Flags().StringVarP(&copyDestination, "destination", "d", "", "Destination image reference (e.g., registry.io/image:newtag) (required)")
	copyCmd.Flags().StringVar(&copySourceSecret, "source-secret", "", "Kubernetes secret name for source registry authentication (optional for public registries)")
	copyCmd.Flags().StringVar(&copyDestSecret, "dest-secret", "", "Kubernetes secret name for destination registry authentication (optional for public registries)")
	copyCmd.Flags().StringVar(&copySourceNamespace, "source-namespace", "default", "Kubernetes namespace for source secret")
	copyCmd.Flags().StringVar(&copyDestNamespace, "dest-namespace", "default", "Kubernetes namespace for destination secret")
	copyCmd.Flags().BoolVarP(&copyShowProgress, "progress", "p", false, "Show progress during copy operation")

	// Mark required flags
	_ = copyCmd.MarkFlagRequired("source")
	_ = copyCmd.MarkFlagRequired("destination")
}
