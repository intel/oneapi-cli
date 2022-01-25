// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause

package browser

import (
	exec "golang.org/x/sys/execabs"
)

//OpenBrowser opens the url passed
func OpenBrowser(url string) error {
	cmd := exec.Command("xdg-open", url)
	//Maybe we need to disown the process, but we should be better
	//than just calling run.
	return cmd.Start()

}
