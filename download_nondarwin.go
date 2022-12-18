//go:build !darwin
// +build !darwin

package main

func (i installer) downloadHomebrew() {
	panic("error: only for macos users only.")
}

func (i installer) downloadAndroidStudio() {
	dlType := "install"
	ext := "exe"
	os := i.os

	if i.os == "linux" {
		ext = "tar.gz"
		dlType = "ide-zips"
	}

	androidStudioVersion := "2021.2.1.16"
	filename := fmt.Sprintf("android-studio-ide-%s-%s.%s", androidStudioVersion, os, ext)
	downloadFile(fmt.Sprintf("https://redirector.gvt1.com/edgedl/android/studio/%s/%s/%s", androidStudioVersion, dlType, filename))
}
