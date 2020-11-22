package main

func init() {
	prerequisites = map[string]string{
		"Homebrew": "brew",
		"Git":      "git",
		"JDK":      "javac",
		// TODO:
		"Android Studio": "android-studio",
		"Flutter":        "flutter",
	}
}

func (i installer) downloadHomebrew() {}
func (i installer) downloadGit()      {}
