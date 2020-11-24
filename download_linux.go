package main

import (
	"fmt"
	"strings"
)

func init() {
	prerequisites = [][]string{
		{"Git", "git"},
		{"JDK", "javac"},
		{"Android Studio", "android-studio"},
		{"Flutter", "flutter"},
	}
}

func (i installer) downloadGit() {
	args := []string{"sudo"}

	switch i.platformFamily {
	case "ubuntu":
		args = append(args, "apt-get", "-y", "install", "git")
	case "debian":
		args = append(args, "apt-get", "-y", "install", "git")
	case "fedora":
		// > fedora 22
		args = append(args, "dnf", "install", "git")
	case "gentoo":
		args = append(args, "emerge", "--ask", "--verbose", "dev-vcs/git")
	case "arch":
		args = append(args, "pacman", "-S", "git")
	case "suse":
		args = append(args, "zypper", "install", "git")
	case "alpine":
		args = append(args, "apk", "add", "git")
	default:
		panic(fmt.Sprintf("no git installer matched for \"%s\"", i.platformFamily))
	}

	fmt.Printf("Executing %s...\n", strings.Join(args, " "))
	execute(strings.Join(args, " "))
}

func (i installer) downloadJDK() {
	args := []string{"sudo"}

	if i.platformFamily == "ubuntu" || i.platformFamily == "debian" {
		args = append(args, "apt-get", "-y", "install", "openjdk-13-jdk")
		fmt.Printf("Executing %s...\n", strings.Join(args, " "))
		execute(strings.Join(args, " "))
		return
	}

	i.downloadJDKFromMirror()
	// panic(fmt.Sprintf("no java installer matched for \"%s\"", i.platformFamily))
}
