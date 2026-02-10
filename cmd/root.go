/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var appVersion, appCommit, appDate string

// SetVersionInfo sets the version info from main (populated by ldflags)
func SetVersionInfo(version, commit, date string) {
	appVersion = version
	appCommit = commit
	appDate = date
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "repo-lister",
	Short:   "A CLI tool to manage container images using Kubernetes credentials",
	Version: "dev",
	Long: `repo-lister is a CLI tool that helps you manage container images across registries
using Kubernetes credentials for authentication.

Features:
  - list:  List image tags from a registry
  - copy:  Copy/retag images between registries
  - pull:  Pull images from registry to local storage
  - push:  Push images from local storage to registry

All commands use Kubernetes secrets for registry authentication, making it easy
to work with private registries in your cluster.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", appVersion, appCommit, appDate)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
