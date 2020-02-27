// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package browser

import (
	"os/exec"
	"strings"
	"syscall"
)

//OpenBrowser opens the url passed
func OpenBrowser(url string) error {
	r := strings.NewReplacer("&", "^&")
	cmd := exec.Command("cmd", "/c", "start", r.Replace(url))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()

}
