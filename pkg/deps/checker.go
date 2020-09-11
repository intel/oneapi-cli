// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
package deps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	rootEnvKey = "ONEAPI_ROOT"

	baseURL = "https://software.intel.com/en-us/oneapi/"

	formatStr = `The following tools are needed to build this sample but are not locally installed: (%s)
You may continue and view the sample without the prerequisites. To install the missing prerequisites, visit:
%s%s
`

	pkgReg      = "pkg\\|([^|]*)\\|?(.*)"
	compilerReg = "compiler\\|(.*)"
)

//CheckDeps as
func CheckDeps(dependencies []string, root string) (msg string, errCode int) {
	//dependencies are both "normal" component dependencies ( ["mkl", "vtune"])
	//and "special" dependencies  ( ["pkg|mraa", "compiler|icc"])

	componentDeps, specialDeps := separatethSheepsGoats(dependencies)

	componentMsg, componentErrCode := checkComponentDeps(componentDeps, root)
	specialMsg, specialErrCode := checkSpecialDeps(specialDeps, root)

	//now gather results and return
	msg, errCode = simplifyMsgErrCode(specialMsg, specialErrCode, componentMsg, componentErrCode)
	return msg, errCode
}
func simplifyMsgErrCode(msg1 string, errCode1 int, msg2 string, errCode2 int) (msg string, errCode int) {
	//takes a two pairs of messages and error codes and returns their concatenation (or whatever is appropriate)
	msg = ""
	errCode = 0
	divider := ""

	if errCode1 != 0 {
		msg = msg1
		errCode = errCode1
		divider = "\n"
	}
	if errCode2 != 0 {
		msg = msg + divider + msg2
		errCode = errCode2
	}
	return msg, errCode
}

func checkComponentDeps(dependencies []string, root string) (msg string, errCode int) {

	var missing []string
	for _, k := range dependencies {
		if _, err := os.Stat(filepath.Join(root, k)); os.IsNotExist(err) {
			//does NOT EXIST
			missing = append(missing, k)
		}
	}
	// Something was missing, get a message
	if len(missing) > 0 {
		msg := GenerateMessage(missing)
		return msg, -1
	}
	return "", 0
}

func checkPackageDeps(packageDeps []string, root string) (msg string, errCode int) {
	// deps = pckg|<package-name>|url
	// 1. check for pkg-config
	// 1.F   if not: message returned says "this sample requires <package-name> which we unable to verify.  Be sure it is installed. <url>."
	// 1.T   if so: call    `pkg-config --exists <package-name>`
	// 1.T.T if exists - OK, no message
	// 1.T.F if not: "this sample requires <package-name> which is not installed. You can obtain it here: <url>"

	//errCode -1 no package    , -2 no pkg-config

	msg = ""
	errCode = 0
	divider := ""

	//0. setup regex that will parse dependency
	regEx := regexp.MustCompile(pkgReg)

	//1.
	_, pkgErr := exec.LookPath("pkg-config")

	for _, dep := range packageDeps {
		pkg, url := parseDep(regEx, dep)
		err := pkgErr
		if err == nil {
			cmd := exec.Command("pkg-config", "--exists", pkg)
			err = cmd.Run()
			if err == nil {
				//we are good to go.
			} else {
				msg = msg + divider + fmt.Sprintf("this sample requires %s which is not installed. To obtain: %s", pkg, url)
				divider = "\n"
				errCode = -1
			}
		} else {
			msg = msg + divider + fmt.Sprintf("this sample requires %s which we are unable to verify. Please make sure it is installed. %s", pkg, url)
			divider = "\n"
			errCode = -2
		}
	}

	return msg, errCode
}

func parseDep(re *regexp.Regexp, dep string) (pkg string, url string) {
	//given a regex and a string, parses it out to package and url.
	// this parser can be used for both pkg| and compiler| , use constants pkgReg or compilerReg as first arg
	// example: parseDep(pgkReg, "pkg|mraa|www.intel.com") => "mraa", "www.intel.com"
	//          parseDep(compilerReg, "compiler|icc") => "icc", ""
	match := re.FindStringSubmatch(dep)
	pkg = ""

	if len(match) > 1 {
		pkg = match[1]
		url = fmt.Sprintf("Search for 'install %s' for help.", pkg)
	}
	if len(match) > 2 {
		url = match[2]
	}
	return pkg, url
}

func checkCompilerDeps(compilerDeps []string, root string) (msg string, errCode int) {

	msg = ""
	errCode = 0
	var missing []string

	winCompilers := map[string]string{
		"icc":     "windows/bin/intel64/icl.exe",
		"fortran": "windows/bin/intel64/ifort.exe",
		"dpcpp":   "windows/bin/dpcpp.exe",
		"icpc":    "windows/bin/intel64/icpc.exe",
		"icx":     "windows/bin/icx.exe",
		"icpcx":   "windows/bin/icpcx.exe",
	}
	linCompilers := map[string]string{
		"icc":     "linux/bin/intel64/icc",
		"fortran": "linux/bin/intel64/ifort",
		"dpcpp":   "linux/bin/dpcpp",
		"icpc":    "linux/bin/intel64/icpc",
		"icx":     "linux/bin/icx",
		"icpcx":   "linux/bin/icpcx",
	}
	macCompilers := map[string]string{
		"icpc":  "mac/bin/intel64/icpc",
		"icc":   "mac/bin/intel64/icc",
		"ifort": "mac/bin/intel64/ifort",
		"icx":   "mac/bin/icx",
		"icpcx": "mac/bin/icpcx",
	}

	compilerRoot := GetCompilerRoot(root)

	//0. setup regex that will parse dependency
	regEx := regexp.MustCompile(compilerReg)

	for _, dep := range compilerDeps {
		compiler, _ := parseDep(regEx, dep)
		var pathTail string

		switch runtime.GOOS {
		case "linux":
			pathTail = linCompilers[compiler]
		case "windows":
			pathTail = winCompilers[compiler]
		case "darwin":
			pathTail = macCompilers[compiler]
		default:
			msg = "Cannot check Compiler, unsupported OS"
			errCode = 01
			return msg, errCode
		}

		fullPath := filepath.Join(compilerRoot, pathTail)
		if pathTail == "" || !fileExists(fullPath) {
			missing = append(missing, compiler)
		}
	}
	if len(missing) > 0 {
		msg = GenerateMessage(missing)
		errCode = 01
	}
	return msg, errCode
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func checkSpecialDeps(specialDependencies []string, root string) (msg string, errCode int) {
	packageDeps, remainingDeps := separatethSheepsGoatsRhematosC(specialDependencies, func(dep string) bool { return strings.HasPrefix(dep, "pkg|") })
	compilerDeps, remainingDeps := separatethSheepsGoatsRhematosC(remainingDeps, func(dep string) bool { return strings.HasPrefix(dep, "compiler|") })

	packageMsg, packageErrCode := checkPackageDeps(packageDeps, root)
	compilerMsg, compilerErrCode := checkCompilerDeps(compilerDeps, root)
	msg, errCode = simplifyMsgErrCode(packageMsg, packageErrCode, compilerMsg, compilerErrCode)
	return msg, errCode
}

//GetOneAPIRoot gets the root the OneAPI installation
//based on the ONEAPI_ROOT
func GetOneAPIRoot() (path string, err error) {
	root, ok := os.LookupEnv(rootEnvKey)
	if !ok {
		return "", fmt.Errorf("%s not defined.  Be sure to run oneapi environment script ( source setvars.sh )", rootEnvKey)
	}
	return root, nil
}

// CMPLR_ROOT was, for awhile, always defined in the environment. But no longer.
func GetCompilerRoot(root string) (compilerRoot string) {
	return filepath.Join(root, "compiler", "latest")
}

type suiteComponent struct {
	SuiteId     string `json:"suiteId"`
	ComponentId string `json:"componentId"`
	Primary     bool   `json:"primary"`
}

type compDir struct {
	Dir         string `json:"dir"`
	ComponentId string `json:"componentId"`
}

type suite struct {
	SuiteId     string `json:"id"`
	Label       string `json:"label"`
	UrlSlug     string `json:"urlSlug"`
	BaseToolkit string `json:"baseToolkit"`
}

func GenerateMessage(missing []string) string {

	//1.0 read in compmapping.json
	//1.1 - translate missing to id list
	//2.0 read in suite-components.json
	//2.1 expand id-list to toolkits
	//3.0 count/find most frequent toolkit (if any)
	//4.0 read suites.json
	//4.1 get slug from toolkit

	//5. return message with url.

	//baseURL := "https://software.intel.com/en-us/oneapi"   // this is captured in format String
	var slug string
	fallbackMsg := fmt.Sprintf(formatStr, strings.Join(missing, " "), baseURL, slug)

	//1  read in compmapping.json
	var mapping []compDir
	//err := readSomeJSON(filepath.Join("json", "compmapping.json"), &mapping)
	err := parseSomeJSON(compmappingJSON, &mapping)
	if err != nil {
		return fallbackMsg
	}
	//1.1  translate missing to id list
	idList := mapStringArr(missing, func(miss string) string {
		for _, v := range mapping {
			if v.Dir == miss {
				return v.ComponentId
			}
		}
		return ""
	})

	//2.0 read in suite-components.json
	var sweetComps []suiteComponent
	//err = readSomeJSON(filepath.Join("json", "suite-components.json"), &sweetComps)
	err = parseSomeJSON(sweetComponentsJSON, &sweetComps)
	if err != nil {
		return fallbackMsg
	}
	//2.1 expand id-list to toolkits

	var matchedSuite []string

	for _, suiteComp := range sweetComps {
		if contains(idList, suiteComp.ComponentId) {
			matchedSuite = append(matchedSuite, suiteComp.SuiteId)
		}
	}

	if len(matchedSuite) > 1 {
		//4.0 read suites.json
		var suites []suite
		//err = readSomeJSON(filepath.Join("json", "suites.json"), &suites)
		err = parseSomeJSON(suitesJSON, &suites)
		if err != nil {
			return fallbackMsg
		}

		var composed string

		for _, suite := range matchedSuite {
			slug := findSlug(suites, suite)

			if len(composed) == 0 {
				composed = fmt.Sprintf(formatStr, strings.Join(missing, " "), baseURL, slug)
			} else {
				composed = fmt.Sprintf("%sor %s%s", composed, baseURL, slug)
			}
		}
		return composed
	}

	return fallbackMsg
}

func readSomeJSON(path string, something interface{}) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	json.Unmarshal(byteValue, something)
	return nil
}

//keeping the same interface as the file reading function for now.
func parseSomeJSON(str string, something interface{}) error {
	json.Unmarshal([]byte(str), something)
	return nil
}

func mapStringArr(arr []string, f func(a string) string) []string {
	var result []string
	for _, v := range arr {
		result = append(result, f(v))
	}
	return result
}

func contains(arr []string, needle string) bool {
	for _, v := range arr {
		if needle == v {
			return true
		}
	}
	return false
}

func findSlug(suites []suite, maxSuite string) string {
	for _, aSuite := range suites {
		if aSuite.SuiteId == maxSuite {
			return aSuite.UrlSlug
		}
	}
	//not worth dealing with the error
	return ""
}

func separatethSheepsGoats(dependencies []string) ([]string, []string) {
	// the dependencies consist of "normal" component dependencies ( ["mkl", "vtune"])  i.e. "sheep"
	// and "special" dependencies ( ["pkg|mraa", "compiler|icc"])  i.e. "goats"
	// as foretold, they are separateth into two groups.
	return separatethSheepsGoatsRhematosC(dependencies, func(dep string) bool { return !strings.Contains(dep, "|") })
}

func separatethSheepsGoatsRhematosC(dependencies []string, predicate func(a string) bool) ([]string, []string) {
	// And before him shall be gathered all nations, and he shall separate them one from another as a shepherd separateth the sheep from the goats.
	// this function takes a divine predicate and uses that to split a string array in twain
	var sheep []string
	var goats []string
	for _, dep := range dependencies {
		if predicate(dep) {
			sheep = append(sheep, dep)
		} else {
			goats = append(goats, dep)
		}
	}
	return sheep, goats
}
