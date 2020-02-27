// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

var language string
var outputJSON bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:    "list",
	Short:  "List Samples",
	Hidden: true,
	Long: `Lists the available samples. Checks online if newer sample index
	is available`,
	Run: func(cmd *cobra.Command, args []string) {

		if language == "" {
			for _, l := range getAggregator().GetLanguages() {
				fmt.Printf("%s\n", l)
			}
			os.Exit(1)
		}

		if getAggregator().Samples[language] == nil {
			fmt.Printf("Invalid language provided, available languages: %v\n", getAggregator().GetLanguages())
			os.Exit(1)
		}

		if outputJSON {
			fmt.Printf("%s\n", prettyPrint(getAggregator().Samples[language]))
			return
		}

		for _, s := range getAggregator().Samples[language] {
			fmt.Printf("%s:\n\t%s\n", s.Fields.Name, s.Fields.Description)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&language, "output", "o", "", "specific language samples you want to list")
	listCmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "output as JSON")
}
