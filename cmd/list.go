package cmd

import (
	"os"
	"repo-lister/utility"

	"github.com/spf13/cobra"
)

var listImageName, listImageFilter, listSecretName, listNamespace string
var listLimit int

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List image tags from a container registry",
	Long: `List all tags for a container image from a registry using Kubernetes credentials.

This command uses Kubernetes secrets to authenticate with container registries
and lists available tags for the specified image. Results can be filtered using
regex patterns and limited to a specific number of results.`,
	Example: `  # List tags from a public registry (no secret needed)
  repo-lister list --image linuxarpan/testpush --limit 5

  # List tags from a private registry with secret
  repo-lister list --image myregistry.io/app --secret registry-cred --namespace default --limit 5

  # List tags with a filter
  repo-lister list --image myregistry.io/app --secret registry-cred --filter "v[0-9]+.*"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the ListImage function from the utility package
		tags, err := utility.ListImage(listImageName, listImageFilter, listSecretName, listNamespace, listLimit)
		if err != nil {
			cmd.PrintErrln("Error listing image tags:", err)
			os.Exit(1)
		}
		for _, tag := range tags {
			cmd.Println(tag)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Define flags for the list command
	listCmd.Flags().StringVarP(&listImageName, "image", "i", "", "Image name to list tags for (required)")
	listCmd.Flags().StringVarP(&listImageFilter, "filter", "f", ".*", "Regex filter to apply to image tags")
	listCmd.Flags().StringVarP(&listSecretName, "secret", "s", "", "Kubernetes secret name for registry authentication (optional for public registries)")
	listCmd.Flags().StringVarP(&listNamespace, "namespace", "n", "default", "Kubernetes namespace where the secret is located")
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 5, "Maximum number of tags to return")

	// Mark required flags
	listCmd.MarkFlagRequired("image")
}
