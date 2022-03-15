// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package browser

import (
	exec "golang.org/x/sys/execabs"
)

//OpenBrowser opens the url passed
func OpenBrowser(url string) error {
	cmd := exec.Command("open", url)
	return cmd.Run()
}
