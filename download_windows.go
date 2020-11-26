package main

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
)

func init() {
	prerequisites = [][]string{
		{"Git", "git"},
		{"JDK", "javac"},
		{"Visual Studio Code", "code"},
		{"Android Studio", filepath.Join("C:\\", "Program Files", "Android", "Android Studio", "bin", "studio64.exe")},
		{"Flutter", "flutter"},
	}
}

func getGitWindowsLatestRelease(arch string) (string, string, error) {
	rel, _, err := ghClient.Repositories.GetLatestRelease(context.Background(), "git-for-windows", "git")
	panicIfErr(err)

	installerPrefix := "64-bit.exe"
	if strings.HasSuffix(arch, "86") {
		installerPrefix = "32-bit.exe"
	}

	for _, asset := range rel.Assets {
		if strings.HasSuffix(*asset.Name, installerPrefix) {
			return *asset.Name, *asset.BrowserDownloadURL, nil
		}
	}

	return "", "", errors.New("no git installer matched for this system")
}

func (i installer) downloadGit() {
	_, installerURL, err := getGitWindowsLatestRelease(i.arch)
	panicIfErr(err)
	downloadFile(installerURL)
}

func (i installer) downloadJDK() {
	i.downloadJDKFromMirror()
}
