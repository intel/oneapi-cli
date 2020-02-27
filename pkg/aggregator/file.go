// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package aggregator

import (
	"crypto/sha256"
	"io/ioutil"
	"os"
)

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

//localhash returns hash,  error
func localHash(path string) ([]byte, error) {
	hasher := sha256.New()
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
