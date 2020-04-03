// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/intel/oneapi-cli/pkg/aggregator"
	"github.com/intel/oneapi-cli/pkg/ui"
	"github.com/spf13/cobra"
)

//SamplesEndpointDefault default samples endpoint
const SamplesEndpointDefault = "https://iotdk.intel.com/samples-iss"

//2021.1-beta05/

//SampleLatestKey the location from the default path that points to a "latest version"
const SampleLatestKey = "latest"

//LocalStorageDefault the default path root where the local cache is kept
const LocalStorageDefault = ".oneapi-cli"

var baseURL string
var baseFilePath string
var cAggregator *aggregator.Aggregator
var defaultLanguages = []string{"cpp", "python"}
var enabledLanguages []string
var userHome string
var ignoreOS bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oneapi-cli",
	Short: "oneapi-cli a tool to fetch samples",
	Long: `oneapi-cli is tool for fetching samples. It intends to be used either
	interactively or called from another tool`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Connecting to online Sample Aggregator, this may take some time based on network conditions\n")
		app, err := ui.NewCLI(getAggregator(), userHome)
		if err != nil {
			log.Fatal(err)
		}
		app.Show()

	},
}

func getAggregator() *aggregator.Aggregator {
	if cAggregator == nil {
		var err error
		cAggregator, err = aggregator.NewAggregator(baseURL, baseFilePath, enabledLanguages, ignoreOS)
		if err != nil && err != aggregator.ErrCacheLock {
			//Most errors we are going to find are network related :/
			fmt.Printf("Failed to fetch sample index, this *may* be your network/proxy environment.\nYou might try setting http_proxy in your environment, for example:\n")
			fmt.Printf("\tLinux: export http_proxy=http://your.proxy:8080\n")
			fmt.Printf("\tWindows: set http_proxy=http://your.proxy:8080\n")
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		if err == aggregator.ErrCacheLock {
			fmt.Printf("Local Sample cache is corrupt! Please clean the cache and retry!\n")
			fmt.Printf("\toneapi-cli clean\n")
			fmt.Printf("\toneapi-cli\n")
			os.Exit(1)
		}
	}
	return cAggregator

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	var err error
	userHome, err = os.UserHomeDir()
	if err != nil {
		fmt.Printf("Unable to locate Home Directory - %v\n", err)
		os.Exit(1)
	}
	defaultBaseFilePath := filepath.Join(userHome, LocalStorageDefault)

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cmd.yaml)")
	rootCmd.PersistentFlags().StringVarP(&baseURL, "url", "u", getVersionInfo(), "URL of remote sample aggregator")
	rootCmd.PersistentFlags().StringVarP(&baseFilePath, "directory", "d", defaultBaseFilePath, "location to store local oneapi samples cache")
	rootCmd.PersistentFlags().StringSliceVarP(&enabledLanguages, "languages", "l", defaultLanguages, "enabled languages")
	rootCmd.PersistentFlags().BoolVar(&ignoreOS, "ignore-os", false, "ignore Host-OS based filtering when showing/outputting samples")

}

//looks at the command bin path and looks for "version.txt" which
//points to which sample version to look at. If it cant find it it
//returns the Latestkey const
func getVersionInfo() string {
	bin, err := os.Executable()
	if err != nil {
		return fmt.Sprintf("%s/%s/", SamplesEndpointDefault, SampleLatestKey)
	}
	versionPath := filepath.Join(filepath.Dir(bin), "version.txt")

	file, err := os.Open(versionPath)
	if err != nil {
		return fmt.Sprintf("%s/%s/", SamplesEndpointDefault, SampleLatestKey)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("%s/%s/", SamplesEndpointDefault, SampleLatestKey)
	}
	return fmt.Sprintf("%s/%s/", SamplesEndpointDefault, scanner.Text())
}
