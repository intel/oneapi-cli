// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package aggregator

// Sample Type
type Sample struct {
	Path   string `json:"path"`
	SHA    string `json:"sha"`
	Fields Fields `json:"example"`
}

// Fields type (nested struct in sample type)
type Fields struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Categories   []string `json:"categories"`
	Author       string   `json:"author"`
	Date         string   `json:"date"`
	Tag          string   `json:"tag"`
	Dependencies []string `json:"dependencies"`
	OS           []string `json:"os"`
	ReadmeURI    string   `json:"sample_readme_uri"`
	TargetDevice []string `json:"targetDevice"`
	Builder      []string `json:"builder"`
	Toolchain    []string `json:"toolchain"`

	//Not Parsing these out atm
	ProjectOptions   []interface{} `json:"projectOptions"`
	MakeVariables    map[string]interface{}
	IndexerVariables map[string]interface{}
}
