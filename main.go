/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import "repo-lister/cmd"

// Set via ldflags by GoReleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
