// +build linux

package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func init() {
	prerequisites = map[string]string{
		"Git":            "git",
		"JDK":            "javac",
		"Android Studio": "android-studio",
		"Flutter":        "flutter",
	}
}

func (i installer) downloadGit() {
	command := ""
	args := []string{}

	switch i.platformFamily {
	case "debian":
		command = "apt-get"
		args = append(args, "install", "git")
	case "fedora":
		// > fedora 22
		command = "dnf"
		args = append(args, "install", "git")
	case "gentoo":
		command = "emerge"
		args = append(args, "--ask", "--verbose", "dev-vcs/git")
	case "arch":
		command = "pacman"
		args = append(args, "-S", "git")
	case "suse":
		command = "zypper"
		args = append(args, "install", "git")
	case "alpine":
		command = "apk"
		args = append(args, "add", "git")
	default:
		panic(fmt.Sprintf("linux platform \"%s\" not supported", i.platformFamily))
	}

	fmt.Printf("Executing %s %s\n", command, strings.Join(args, " "))
	cmd := exec.Command("sudo "+command, args...)
	stdout, err := cmd.StdoutPipe()
	panicIfErr(err)
	cmd.Start()
	buf := bufio.NewReader(stdout)
	num := 1
	for {
		line, _, _ := buf.ReadLine()
		fmt.Println(string(line))
		num++
		if num > 3 {
			break
		}
	}
}
