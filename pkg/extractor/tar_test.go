// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
package extractor

import (
	"crypto/sha512"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

//testdata/golden.tar.gz contents
const testalpha = "testalpha.txt"
const testbeta = "testbeta.txt"
const testdir = "dirtest"
const testalphaSum = "10b8eefa145e6f3ff612197247765c0fa15788874b2f483879c8528383489d97f4ac18daed45334341d55342c94602976b401c88e879a9d8fd950c69040e29e2"
const testbetaSum = "e73842277cc6739947522cc2cac9dc150524d2d6683fe54b79da6a098a2cffc8a9498583c2f17d7897db36fa1980d6d0e729f4989780e0be38ef3684210c5d99"

func setupGoldTemp(t *testing.T) string {
	t.Helper()
	dir, err := ioutil.TempDir("", "extar")
	if err != nil {
		t.Error(err)
	}
	return dir
}

func cleanupTemp(t *testing.T, path string) {
	t.Helper()
	os.RemoveAll(path)
}

//localhash returns hash,  error
func localHash(t *testing.T, path string) string {
	t.Helper()
	hasher := sha512.New()
	local, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}

	_, err = hasher.Write(local)
	if err != nil {
		t.Error(err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func TestExtractTarGz(t *testing.T) {

	golden := filepath.Join("testdata", "golden.tar.gz")

	tempPath := setupGoldTemp(t)
	defer cleanupTemp(t, tempPath)
	output := filepath.Join(tempPath, "gold")

	err := ExtractTarGz(golden, output)
	if err != nil {
		t.Error(err) // Failed to run against golden tar
	}

	//test the output of alpha was good
	extractedAlphaSum := localHash(t, filepath.Join(output, testalpha))
	if extractedAlphaSum != testalphaSum {
		t.Errorf("outputted alpha file does not match expected output sha512 "+
			"golden: %s and we just read: %s", testalphaSum, extractedAlphaSum)
	}

	betaCombinedPath := filepath.Join(output, testdir, testbeta)
	//test the output of beta was good, this also tests the directory was made!
	extractedBetaSum := localHash(t, betaCombinedPath)
	if extractedBetaSum != testbetaSum {
		t.Errorf("outputted beta file does not match expected output sha512 "+
			"golden: %s and we just read: %s", testbetaSum, extractedBetaSum)
	}

}

func TestOKTarGz(t *testing.T) {

	oktar := filepath.Join("testdata", "ok.tar.gz")

	tempPath := setupGoldTemp(t)
	defer cleanupTemp(t, tempPath)
	output := filepath.Join(tempPath, "ok")

	err := ExtractTarGz(oktar, output)
	if err != nil {
		t.Error(err) // Failed to run against ok-ish tar
	}
}
