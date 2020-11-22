package main

func init() {
	prerequisites = [][]string{
		{"Homebrew", "brew"},
		{"Git", "git"},
		{"JDK", "javac"},
		{"Android Studio", "android-studio"},
		{"Flutter", "flutter"},
	}
}

func (i installer) downloadHomebrew() {
	execute("\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)\"")
}

func (i installer) downloadGit() {
	execute("brew install git")
}
