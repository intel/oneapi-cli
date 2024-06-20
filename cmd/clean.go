// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
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
		if err := os.RemoveAll(filepath.Join(baseFilePath, aggregator.AggregatorLocalAPILevel)); err != nil {
			fmt.Println("Failed to clean sample cache.")
			fmt.Printf("%s \n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
