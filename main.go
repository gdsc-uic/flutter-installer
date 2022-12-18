package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/Delta456/box-cli-maker/v2"
	"github.com/cavaliergopher/grab/v3"
	"github.com/cheggaaa/pb/v3"
	"github.com/google/go-github/v32/github"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// DSC Flutter Installer
// Flow:
// - Check system environment
//   - NOTE: check disk space
// - Install homebrew (??) (for mac)
// - Check and install Git (if not present)
// - Check and install Java (if not present)
// - Check and install Android studio (if not present)
//   - TODO: Additional setup
// - Check and install Flutter (if not present)

type installer struct {
	os             string
	platformFamily string
	arch           string
}

const (
	version        = "1.0.1"
	flutterVersion = "3.3.10"
)

var (
	downloadFolder = ""
	ghClient       = github.NewClient(nil)
	dlClient       = grab.NewClient()
	// map[Program Name]program (e.g map[Git]git)
	prerequisites = [][]string{}
)

// https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func byteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getRegexMatches(r *regexp.Regexp, txt string) map[string]string {
	matches := r.FindStringSubmatch(txt)[1:]
	keys := r.SubexpNames()[1:]
	matchMap := map[string]string{}
	for i := range matches {
		matchMap[keys[i]] = matches[i]
	}
	return matchMap
}

func programExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func downloadFile(url string) {
	fmt.Printf("Downloading %s\n", url)
	req, _ := grab.NewRequest(downloadFolder, url)
	req.Filename = filepath.Join(downloadFolder, filepath.Base(url))
	resp := dlClient.Do(req)
	t := time.NewTicker(100 * time.Millisecond)
	downloadBar := pb.Full.Start64(resp.Size())
	downloadBar.Set(pb.SIBytesPrefix, true)
	downloadBar.Set(pb.Bytes, true)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			downloadBar.SetCurrent(resp.BytesComplete())
		case <-resp.Done:
			downloadBar.Finish()
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		panic(fmt.Sprintf("Download failed: %v\n", err))
	}

	fmt.Printf("Download for \"%s\" was complete.\n", filepath.Base(url))
}

func (i installer) downloadJDKFromMirror() {
	installerURL := ""
	// TODO:
	rel, _, err := ghClient.Repositories.GetReleaseByTag(context.Background(), "AdoptOpenJDK", "openjdk15-binaries", "jdk-15.0.1+9")
	panicIfErr(err)

	assetRegexp, err := regexp.Compile(`(?:OpenJDK15U-jdk_(?P<Arch>\w+?)_(?P<OS>.*?)_.*?)(?P<Format>.msi|.pkg|.tar.gz)`)
	panicIfErr(err)

	for _, asset := range rel.Assets {
		if !assetRegexp.MatchString(*asset.Name) {
			continue
		}
		matches := getRegexMatches(assetRegexp, *asset.Name)
		if matches["OS"] != i.os {
			continue
		}
		if (i.arch == matches["Arch"]) || (i.arch == "x86_64" && matches["Arch"] == "x64") || strings.HasPrefix(matches["Arch"], "32") || (i.arch == "arm64" && matches["Arch"] == "aarch64") {
			installerURL = *asset.BrowserDownloadURL
			break
		}
	}

	if len(installerURL) > 0 {
		downloadFile(installerURL)
	}

	panic(fmt.Sprintf("jdk installer not found for \"%s\"", i.os))
}

func (i installer) downloadVscode() {
	systemName := i.os

	if i.os == "windows" {
		systemName = "win32"

		switch i.arch {
		case "x86_64":
			systemName += "-x64"
		case "arm":
			systemName += "-arm64"
		}

		systemName += "-user"
	} else if i.os == "linux" {
		if i.platformFamily == "ubuntu" || i.platformFamily == "debian" {
			systemName += "-deb"
		} else {
			systemName += "-rpm"
		}

		switch i.arch {
		case "arm":
			systemName += "armhf"
		case "arm64":
			systemName += "arm64"
		default:
			systemName += "x64"
		}
	}

	downloadLink := fmt.Sprintf("https://update.code.visualstudio.com/latest/%s/stable", systemName)
	downloadFile(downloadLink)
}

func (i installer) downloadFlutter() {
	os := i.os
	ext := "zip"
	if os == "darwin" {
		os = "macos"
	} else if os == "linux" {
		ext = "tar.xz"
	}

	filename := fmt.Sprintf("flutter_%s_%s-stable.%s", os, flutterVersion, ext)
	downloadLink := fmt.Sprintf("https://storage.googleapis.com/flutter_infra_release/releases/stable/%s/%s", os, filename)
	downloadFile(downloadLink)
}

func main() {
	sysInfo, diskInfo := strings.Builder{}, strings.Builder{}
	wd, err := os.Getwd()
	panicIfErr(err)
	downloadFolder = filepath.Join(wd, "dsc-flutter-installer_downloads")
	if _, err := os.Stat(downloadFolder); os.IsNotExist(err) {
		os.Mkdir(downloadFolder, 0777)
	}

	warningBox := box.New(box.Config{Px: 2, Py: 0, Type: "Double", ContentAlign: "Center", Color: "Yellow", TitlePos: "Inside"})
	systemInfoBox := box.New(box.Config{Px: 2, Py: 0, Type: "Single", ContentAlign: "Left", Color: "Green", TitlePos: "Top"})
	// Print system environment
	// OS
	// CPU Details
	// RAM Info
	// Memory and free space
	hostStat, _ := host.Info()
	cpuStats, _ := cpu.Info()
	virtMemInfo, _ := mem.VirtualMemory()
	ram := virtMemInfo.Total
	partitions, _ := disk.Partitions(false)

	sysInfo.WriteString(fmt.Sprintf("%-8v%s %v\n", "OS:", hostStat.OS, hostStat.Platform))
	sysInfo.WriteString(fmt.Sprintf("%-8v%v\n", "Arch:", hostStat.KernelArch))

	// print the cpus used
	cpus := map[string]int{}
	for _, cpuStat := range cpuStats {
		if _, cpuExists := cpus[cpuStat.ModelName]; cpuExists {
			cpus[cpuStat.ModelName]++
		} else {
			cpus[cpuStat.ModelName] = 1
		}
	}

	if len(cpus) > 0 {
		sysInfo.WriteString(fmt.Sprintf("%-8v", "CPU:"))
		for cpuName, cpuCores := range cpus {
			sysInfo.WriteString(fmt.Sprintf("%20v\n", fmt.Sprintf("%s x %v", cpuName, cpuCores)))
		}
	}

	sysInfo.WriteString(fmt.Sprintf("%-8v%v", "RAM:", byteCountIEC(ram)))
	diskInfo.WriteString("Disks:\n")
	// print disks
	for i, part := range partitions {
		if strings.HasPrefix(part.Mountpoint, "/boot") || strings.HasPrefix(part.Mountpoint, "/snap") {
			continue
		}
		diskStat, _ := disk.Usage(part.Mountpoint)
		diskInfo.WriteString(fmt.Sprintf("        %-25v %v/%v (%.2f%% used)", part.Mountpoint, byteCountIEC(diskStat.Used), byteCountIEC(diskStat.Total), diskStat.UsedPercent))
		if i+2 < len(partitions) {
			diskInfo.WriteString("\n")
		}
	}

	sysInfo.WriteString("\n" + diskInfo.String())
	systemInfoBox.Println("System Info", sysInfo.String())
	warningBox.Println("WARNING!", fmt.Sprintf("Your downloads will be placed at:\n%s\n\nIf you wish to put the downloads into another directory,\nplease move the installer to its desired destination and run it again.", downloadFolder))
	fmt.Print("Would you like to proceed? (Y/N): ")
	var choice string
	_, _ = fmt.Scanln(&choice)
	if strings.ToLower(choice) != "y" {
		os.Exit(1)
	}

	inst := installer{
		os:             hostStat.OS,
		platformFamily: hostStat.Platform,
		arch:           hostStat.KernelArch,
	}

	for _, requi := range prerequisites {
		programName, execName := requi[0], requi[1]
		if programName == "Android Studio" && runtime.GOOS == "windows" {
			var decision string
			fmt.Print("Do you have Android Studio installed? (Y/N): ")
			_, _ = fmt.Scanln(&decision)
			if strings.ToLower(decision) == "y" {
				var customStudioFolder string

				fmt.Print("Enter Android Studio installation folder (leave blank if using default folder): ")
				_, err = fmt.Scanln(&customStudioFolder)
				if err == nil {
					execName = filepath.Join(customStudioFolder, "bin", "studio64.exe")
				}
			}
		}

		fmt.Printf("Checking %s...", programName)
		time.Sleep(2 * time.Second)
		if isProgramInstalled := programExists(execName); isProgramInstalled {
			fmt.Println(" already installed.")
			continue
		}

		fmt.Printf(" not installed.\nDownloading installer for %s...\n", programName)
		switch programName {
		case "Git":
			inst.downloadGit()
		case "JDK":
			inst.downloadJDK()
		case "Android Studio":
			inst.downloadAndroidStudio()
		case "Flutter":
			inst.downloadFlutter()
		case "Visual Studio Code":
			inst.downloadVscode()
		case "Homebrew":
			inst.downloadHomebrew()
		}
	}
	time.Sleep(2 * time.Second)

	fmt.Println("Download has finished! Press ENTER to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
