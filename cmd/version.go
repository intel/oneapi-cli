// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version set during build via `go build -ldflags -X "-X main.Version=version"`
var version string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the CLI version information",
	Long:  `Show the CLI version information`,
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			version = "devel"
		}
		fmt.Printf("%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
