// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"os"
	"path/filepath"

	"github.com/intel/oneapi-cli/pkg/aggregator"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean Sample Cache",
	Long:  `Removes local Sample Cache`,
	Run: func(cmd *cobra.Command, args []string) {
		os.RemoveAll(filepath.Join(baseFilePath, aggregator.AggregatorLocalAPILevel))
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
