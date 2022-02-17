// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package browser

import (
	exec "golang.org/x/sys/execabs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

//OpenBrowser opens the url passed
func OpenBrowser(url string) error {
	r := strings.NewReplacer("&", "^&")
	rundll32 := filepath.Join(os.Getenv("SystemRoot"), "System32", "rundll32.exe")

	cmd := exec.Command(rundll32, "url.dll,FileProtocolHandler", r.Replace(url))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()

}
