// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package aggregator

import (
	"crypto/sha512"
	"io/ioutil"
	"os"
)

//FileExists helper function for checking a file exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

//localhash returns hash,  error
func localHash(path string) ([]byte, error) {
	hasher := sha512.New()
	local, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	_, err = hasher.Write(local)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}
