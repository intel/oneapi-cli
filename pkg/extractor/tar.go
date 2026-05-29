// Copyright 2019-2026 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
package extractor

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// secure join to prevent path traversal
func safeJoin(baseDir, name string) (string, error) {
	base := filepath.Clean(baseDir)
	cleanName := filepath.Clean(name)

	// Reject absolute paths
	if filepath.IsAbs(cleanName) {
		return "", fmt.Errorf("absolute paths not allowed: %s", name)
	}

	target := filepath.Join(base, cleanName)

	// Ensure target is within base
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path traversal detected: %s", name)
	}

	return target, nil
}

// ExtractTarGz extracts a tar.gz to the destination
func ExtractTarGz(sourcetb string, out string) error {

	// Ensure Output exists
	if err := os.MkdirAll(out, 0750); err != nil {
		return err
	}

	tbz, err := os.Open(sourcetb)
	if err != nil {
		return err
	}
	defer tbz.Close()

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
			return nil
		case err != nil:
			return err
		}

		// PTK0006726
		target, err := safeJoin(out, hdr.Name)
		if err != nil {
			return err
		}

		switch hdr.Typeflag {

		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil {
				return err
			}

		case tar.TypeReg:
			// Ensure parent dir exists
			if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
				return err
			}

			// Safer file write (truncate if exists)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}

			if err := f.Close(); err != nil {
				return err
			}

		// 🔒 Block symlinks & hardlinks (critical)
		case tar.TypeSymlink, tar.TypeLink:
			return fmt.Errorf("unsupported file type (link): %s", hdr.Name)

		default:
			// Optional: ignore or fail
			return fmt.Errorf("unsupported file type: %v in %s", hdr.Typeflag, hdr.Name)
		}
	}
}
