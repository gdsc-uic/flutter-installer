// +build darwin linux

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/mattn/go-isatty"
)

func init() {
	if (runtime.GOOS == "linux" || runtime.GOOS == "darwin") && !isatty.IsTerminal(os.Stderr.Fd()) {
		execPath, _ := os.Executable()

		// TODO: debian-based distros only. not sure to other distros
		if runtime.GOOS == "linux" {
			execute(fmt.Sprintf("x-terminal-emulator -e \"%s\"", execPath))
		} else {
			execute(fmt.Sprintf("open -a Terminal.app %s", execPath))
		}

		os.Exit(0)
	}
}
