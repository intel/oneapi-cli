// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
package deps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestGetOneAPIRoot(t *testing.T) {

	os.Unsetenv(rootEnvKey) //Clear Env, just incase
	_, err := GetOneAPIRoot()
	if err == nil {
		t.Errorf("should have failed to find the env key, we unset it for test")
	}

	const testValue = "/opt/inteltestval/"
	os.Setenv(rootEnvKey, testValue)
	val, err := GetOneAPIRoot()
	if err != nil {
		t.Error(err)
	}
	if val != testValue {
		t.Errorf("unexpected env value receieved")
	}
	os.Unsetenv(rootEnvKey) //Clear Env, just incase
}

func setupTestRoot(t *testing.T, testDeps []string) (root string) {
	t.Helper()
	root, err := ioutil.TempDir("", "depscheck")
	if err != nil {
		t.Error(err)

	}
	for _, k := range testDeps {
		err := os.MkdirAll(filepath.Join(root, k), os.ModePerm)
		if err != nil {
			t.Error(err)
		}
	}
	return root
}

func TestCheckDeps(t *testing.T) {

	testingGold := []string{"cheese", "milk"}

	root := setupTestRoot(t, testingGold)

	msg, errCode := CheckDeps(testingGold, root)
	if errCode > 0 {
		t.Errorf("Golden test failed found missing %s", msg)
	}

}

func TestGenerateMessage(t *testing.T) {
	missing := []string{"foo"}
	msg := GenerateMessage(missing)
	fmt.Println(msg)

	if !strings.Contains(msg, "(foo)") {
		t.Errorf("foo wasn't reported")
	}
}

func TestCheckCompilerDeps(t *testing.T) {
	testingGold := []string{"cheese", "milk"}

	root := setupTestRoot(t, testingGold)
	deps := []string{"compiler|gomer"}
	_, errCode := checkCompilerDeps(deps, root)
	if errCode == 0 {
		t.Errorf("gomer compiler should never have been found")
	}
}

/*
func TestReadSomeJSONAndTheJSON(t *testing.T) {

	var mapping []compDir
	err := readSomeJSON(filepath.Join("..", "..", "json", "compmapping.json"), &mapping)
	if err != nil {
		t.Errorf("failure loading compmapping.json - %s", err)
	}

	var sweetComps []suiteComponent
	err = readSomeJSON(filepath.Join("..", "..", "json", "suite-components.json"), &sweetComps)
	if err != nil {
		t.Errorf("failure loading suite-components.json - %s", err)
	}

	var suites []suite
	err = readSomeJSON(filepath.Join("..", "..", "json", "suites.json"), &suites)
	if err != nil {
		t.Errorf("failure loading suites.json - %s", err)
	}

	//not worrying about invalid-components-hosts or components.json right now.
}
*/

func TestParseSomeJSONAndTheJSON(t *testing.T) {
	var mapping []compDir
	err := parseSomeJSON(compmappingJSON, &mapping)
	if err != nil {
		t.Errorf("failure parsing  Compmappingjson - %s", err)
	}
	if len(mapping) == 0 {
		t.Errorf("compmappingJSON empty after parsing")
	}

	var sweetComps []suiteComponent
	err = parseSomeJSON(sweetComponentsJSON, &sweetComps)
	if err != nil {
		t.Errorf("failure parsing  sweetComponentsJSON - %s", err)
	}
	if len(sweetComps) == 0 {
		t.Errorf("sweetComponentsJSON empty after parsing")
	}

	var suites []suite
	err = parseSomeJSON(suitesJSON, &suites)
	if err != nil {
		t.Errorf("failure parsing  suitesJSON - %s", err)
	}
	if len(suites) == 0 {
		t.Errorf("suitesJSON empty after parsing")
	}
}

func TestContains(t *testing.T) {
	haystack := []string{"marvel", "crunch", "america", "picard"}
	needle := "picard"
	youshouldlivesolong := "ryker"
	wat := ""

	if !contains(haystack, needle) {
		t.Errorf("%s not found", needle)
	}

	if contains(haystack, youshouldlivesolong) {
		t.Errorf("captain %s?? Never", youshouldlivesolong)
	}

	if contains(haystack, wat) {
		t.Errorf("Alfred North Whitehead hates you")
	}
}

func TestMapStringArr(t *testing.T) {
	haystack := []string{"marvel", "crunch", "america", "picard"}
	fuzz := mapStringArr(haystack, func(a string) string {
		return a[0:1]
	})
	expectation := []string{"m", "c", "a", "p"}
	if !reflect.DeepEqual(fuzz, expectation) {
		t.Errorf("string mapping not that hard, maybe consider career change?")
	}

	emptyArr := []string{}
	ress := mapStringArr(emptyArr, func(a string) string {
		return a
	})
	if len(ress) > 0 {
		t.Errorf("string mapping not that hard, maybe consider career change?")
	}
	//amazing enuogh, DeepEqual returns the wrong value.  These are both empty, it says not equal.
	// if !reflect.DeepEqual(emptyArr, ress) {
	// 	t.Errorf("string mapping not that hard, maybe consider career change? %s  %s", emptyArr, ress)
	// }

}

func findDirectory(mapping []compDir, componentId string) string {
	for _, compDir := range mapping {
		if compDir.ComponentId == componentId {
			return compDir.Dir
		}
	}
	return ""
}

func TestFindDir(t *testing.T) {
	var mapping []compDir
	compMapErr := parseSomeJSON(compmappingJSON, &mapping)

	dir := findDirectory(mapping, "intel_advisor")
	if dir == "" {
		t.Errorf("findDirectory failed, unable to locate advisor in compmappingJSON")
	}

	dir = findDirectory(mapping, "xanadu")
	if dir != "" {
		t.Errorf("In Xanadu did Kubla Khan a stately pleasure dome decree")
	}

	if compMapErr != nil {
		t.Errorf("failure parsing  the JSON constants compmappingJSON: %s", compMapErr)
	}
}

func TestIntegrityOfJSONConstants(t *testing.T) {
	// we have a test above that makes sure the constants (suitesJSON, sweetComponentJSON, etc)
	// are valid JSON strings.
	// this test makes sure that every entry in sweetComponentsJSON has a matching component id in componentsMapping
	// extra entries in componentMapping is ok, there may be vestigial components there. Doesn't matter.
	// and that every entry in sweetComponentJSON has a suite in suitesJSON, and vice versa.

	var mapping []compDir
	compMapErr := parseSomeJSON(compmappingJSON, &mapping)

	var sweetComps []suiteComponent
	sweetCompErr := parseSomeJSON(sweetComponentsJSON, &sweetComps)

	var suites []suite
	suitesErr := parseSomeJSON(suitesJSON, &suites)

	if compMapErr != nil || sweetCompErr != nil || suitesErr != nil {
		t.Errorf("failure parsing  the JSON constants compmappingJSON: %s - sweetComponentsJSON: %s - suitesJSON: %s", compMapErr, sweetCompErr, suitesErr)
	}

	//run through sweetComps looking for id that is not in compDir map
	var dir string
	for _, suiteComp := range sweetComps {
		dir = findDirectory(mapping, suiteComp.ComponentId)
		if dir == "" {
			t.Errorf("unable to locate directory in compmapping JSON for component declared in sweetComponentsJSON: %s", suiteComp.ComponentId)
		}
	}

	//run through sweetComps looking for suite that is not in suitesJSON
	var slug string
	for _, sweetComp := range sweetComps {
		slug = findSlug(suites, sweetComp.SuiteId)
		if slug == "" {
			t.Errorf("unable to locate suite in suitesJSON for entry declared in sweetComponentsJSON: %s", sweetComp.SuiteId)
		}
	}

	//run through suites looking for suite that is not in sweetComponentsJSON
	for _, suiteEntry := range suites {
		match := false
		for _, sweetComp := range sweetComps {
			if sweetComp.SuiteId == suiteEntry.SuiteId {
				match = true
				break
			}
		}
		if !match {
			t.Errorf("unable to locate suite in sweetComponentsJSON for entry declared in sutiesJSON: %s", suiteEntry.SuiteId)
		}
	}

}

func TestSeparatethSheepsGoats(t *testing.T) {
	//violates Deuteronomy 6:16 .  Invoke at your peril.
	sheep := []string{"Dolly", "Montauciel", "Methuselina", "Lance-Corporal-Derby-XXX"}
	goats := []string{"Goat|Billy", "Goat|Pan", "Goat|Nanny", "Goat|Rudy-Giuliani"}

	components, _ := separatethSheepsGoats(sheep)
	if !reflect.DeepEqual(components, sheep) {
		t.Errorf("repent!")
	}

	_, special := separatethSheepsGoats(goats)
	if !reflect.DeepEqual(special, goats) {
		t.Errorf("repent!")
	}

	herd := append(sheep, goats...)

	components, special = separatethSheepsGoats(herd)
	if !reflect.DeepEqual(components, sheep) {
		t.Errorf("repent!")
	}
	if !reflect.DeepEqual(special, goats) {
		t.Errorf("repent!")
	}
}

func TestSimplifyMsgErrCode(t *testing.T) {
	msg, errCode := simplifyMsgErrCode("yams", 0, "hams", 0)
	if errCode != 0 {
		t.Errorf("errCode should be 0")
	}

	msg, errCode = simplifyMsgErrCode("yams", 1, "hams", 0)
	if msg != "yams" {
		t.Errorf("should have yams")
	}

	msg, errCode = simplifyMsgErrCode("yams", 0, "hams", 1)
	if msg != "hams" {
		t.Errorf("should have hams")
	}

	msg, errCode = simplifyMsgErrCode("yams", 1, "hams", 1)
	if msg != "yams\nhams" {
		t.Errorf("should have yams and hams")
	}
}

func TestParseDep(t *testing.T) {
	regEx := regexp.MustCompile(pkgReg)
	pkg, url := parseDep(regEx, "pkg|mraa|http://www.intel.com")
	if pkg != "mraa" {
		t.Errorf("pkg should have parsed to mraa, %s", pkg)
	}
	if url != "http://www.intel.com" {
		t.Errorf("url should have parsed to http://www.intel.com, %s", url)
	}

	//url optional
	pkg, url = parseDep(regEx, "pkg|npm")
	if pkg != "npm" {
		t.Errorf("pkg should have parsed to npm, %s", pkg)
	}

	regEx = regexp.MustCompile(compilerReg)
	compiler, _ := parseDep(regEx, "compiler|icc")
	if compiler != "icc" {
		t.Errorf("compiler should have parsed to icc, %s", compiler)
	}
}

func TestFileExists(t *testing.T) {
	no := fileExists("")
	if no {
		t.Errorf("why does '' exist?")
	}
	yes := fileExists(".")
	if !yes {
		t.Errorf("why doesn't '.' exist?")
	}

}
