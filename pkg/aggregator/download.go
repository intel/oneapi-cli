// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package aggregator

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-ieproxy"
)

func init() {
	http.DefaultTransport.(*http.Transport).Proxy = ieproxy.GetProxyFunc()
}

//downloadFileDirect Fetchs URL into local file
func downloadFileDirect(path string, url string) error {

	// Get the data
	c := &http.Client{
		Timeout: HTTPTimeout * time.Second,
	}
	// Get the data
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP-%v on %s", resp.StatusCode, url)
	}

	//Ensure Directory for local path of language exists

	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func sha256URL(url string) ([]byte, []byte, error) {

	c := &http.Client{
		Timeout: HTTPTimeout * time.Second,
	}
	// Get the data
	resp, err := c.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("HTTP-%v on %s", resp.StatusCode, url)
	}
	// Write the body to file
	hasher := sha256.New()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	_, err = hasher.Write(body)
	if err != nil {
		return nil, nil, nil
	}

	return hasher.Sum(nil), body, nil
}

func checkURL(URL string) (*url.URL, error) {
	return url.ParseRequestURI(URL)
}
