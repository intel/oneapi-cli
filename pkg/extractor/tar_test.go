// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
package extractor

import (
	"crypto/sha256"
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
const testalphaSum = "178f2ea6fbb119e0623edba9baf9a8e3e2607c6791622ad0ab149bfe9573479a"
const testbetaSum = "c1e75e9b1a795c7c4da7180c6a812779e2e6957e47581719b697b880aae07251"

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
	hasher := sha256.New()
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
		t.Errorf("outputted alpha file does not match expected output sha256 "+
			"golden: %s and we just read: %s", testalphaSum, extractedAlphaSum)
	}

	betaCombinedPath := filepath.Join(output, testdir, testbeta)
	//test the output of beta was good, this also tests the directory was made!
	extractedBetaSum := localHash(t, betaCombinedPath)
	if extractedBetaSum != testbetaSum {
		t.Errorf("outputted beta file does not match expected output sha256 "+
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
