package main

func init() {
	prerequisites = [][]string{
		{"Homebrew", "brew"},
		{"Git", "git"},
		{"JDK", "javac"},
		{"Visual Studio Code", "code"},
		{"Android Studio", "/Applications/Android\\ Studio.app"},
		{"Flutter", "flutter"},
	}
}

func (i installer) downloadHomebrew() {
	execute("\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)\"")
}

func (i installer) downloadGit() {
	execute("brew install git")
}

func (i installer) downloadJDK() {
	// tap adoptopenjdk
	execute("brew tap AdoptOpenJDK/openjdk")
	execute("brew cask install adoptopenjdk15")
}

func (i installer) downloadAndroidStudio() {
	execute("brew install android-studio")
}
