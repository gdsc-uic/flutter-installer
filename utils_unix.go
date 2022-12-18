//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func execute(shcmd string) {
	cmd := exec.Command("/bin/sh", "-c", shcmd)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Run(); err != nil {
		panicIfErr(err)
	}

	fmt.Println(stdBuffer.String())
}
