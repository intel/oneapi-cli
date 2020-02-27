// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
package extractor

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

//ExtractTarGz extracts a tar.gz to the destination
func ExtractTarGz(sourcetb string, out string) error {

	//Ensure Output exists
	if err := os.MkdirAll(out, 0750); err != nil {
		return err
	}

	tbz, err := os.Open(sourcetb)
	if err != nil {
		return err
	}

	gzr, err := gzip.NewReader(tbz)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()

		switch {

		case err == io.EOF:
			return nil // return when no more files, good path

		case err != nil:
			return err
		}

		// the target location where the dir/file should be created
		target := filepath.Join(out, hdr.Name)

		// check the file type, are we a directory for example
		switch hdr.Typeflag {

		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil {
				return err
			}

		// we have a file, create it with the stored attr from the header
		case tar.TypeReg:
			//Sometimes the file can come before its directory listing, or it never has one :S
			if !fileExists(filepath.Dir(target)) {
				if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
					return err
				}
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}

			// Store into destination
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		}
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
