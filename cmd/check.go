// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
/*
	usage:
	oneapi-cli check  --deps="mkl,tbb"

	TODO: --html flag

	https://software.intel.com/en-us/oneapi

*/

package cmd

import (
	"fmt"
	"os"

	"github.com/intel/oneapi-cli/pkg/deps"
	"github.com/spf13/cobra"
)

var depsParam []string
var root string

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:    "check",
	Short:  "check dependencies",
	Hidden: true,
	Long:   `check dependencies, returns an error/retrieve-it- message if dependencies are absent`,
	Run: func(cmd *cobra.Command, args []string) {

		if root == "" {
			//Find the oneAPI root
			var err error
			root, err = deps.GetOneAPIRoot()
			if err != nil {
				fmt.Println(err) //Failed to find the Env, may be unset.
				os.Exit(-1)
			}
		}

		//Check the deps at the found root.
		msg, errCode := deps.CheckDeps(depsParam, root)
		if errCode != 0 {
			fmt.Println(msg)
			os.Exit(errCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringSliceVarP(&depsParam, "deps", "", nil, "comma seperated dependency array")
	checkCmd.Flags().StringVar(&root, "oneapi-root", "", "(optional) path to oneAPI root, default attempts to use environment varible ONEAPI_ROOT")
}
