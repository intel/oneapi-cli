// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"fmt"
	"os"

	"github.com/intel/oneapi-cli/pkg/aggregator"
	"github.com/intel/oneapi-cli/pkg/extractor"
	"github.com/spf13/cobra"
)

var sampleLang string

// listCmd represents the list command
var createCmd = &cobra.Command{
	Use:    "create",
	Short:  "Create Sample",
	Hidden: true,
	Long: `Creates the sample based on the passed in path

	i.e. oneapi-cli create -s cpp my/long/path/from/index/json /tmp/mynewproject`,
	Run: func(cmd *cobra.Command, args []string) {

		//Arg 0 being sample
		//arg 1 being where to create the sample. Complete path

		if len(args) != 2 || args[0] == "" || args[1] == "" {
			fmt.Println("Please pass both a sample and where you want it extracted to")
			os.Exit(1)
		}

		tarPath, err := aggregator.GetTarBall(baseFilePath, baseURL, sampleLang, args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		err = extractor.ExtractTarGz(tarPath, args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&sampleLang, "sampleLangauge", "s", "cpp", "specific language of the samples you want to create")
}
