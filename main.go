package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/Delta456/box-cli-maker/v2"
	"github.com/cavaliercoder/grab"
	"github.com/cheggaaa/pb/v3"
	"github.com/google/go-github/v32/github"
	"github.com/shirou/gopsutil/v3/cpu"
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
	version = "1.0"
)

var (
	downloadFolder = ""
	ghClient       = github.NewClient(nil)
	dlClient       = grab.NewClient()
	// map[Program Name]program (e.g map[Git]git)
	prerequisites = map[string]string{}
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
	resp := dlClient.Do(req)
	t := time.NewTicker(200 * time.Millisecond)
	downloadBar := pb.Full.Start64(resp.Size)
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

	fmt.Printf("Download complete.\n")
}

func (i installer) downloadJDK() {
	installerURL := ""

	fmt.Println(" not present.\nDownloading installer for JDK...")
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

	if len(installerURL) == 0 {
		fmt.Println("jdk installer not found")
	}

	downloadFile(installerURL)
}

func (i installer) downloadAndroidStudio() {

}

func (i installer) downloadFlutter() {

}

func main() {
	sysInfo := strings.Builder{}
	wd, err := os.Getwd()
	panicIfErr(err)
	downloadFolder = path.Join(wd, "dsc-flutter-installer_downloads")
	if _, err := os.Stat(downloadFolder); os.IsNotExist(err) {
		os.Mkdir(downloadFolder, os.ModeAppend)
	}

	systemInfoBox := box.New(box.Config{Px: 2, Py: 0, Type: "Single", ContentAlign: "Left", Color: "Green", TitlePos: "Inside"})
	// Print system environment
	// OS
	// CPU Details
	// RAM Info
	// Memory and free space
	hostStat, _ := host.Info()
	cpuStats, _ := cpu.Info()
	virtMemInfo, _ := mem.VirtualMemory()
	ram := virtMemInfo.Total
	// diskStat, _ := disk.Usage("/")

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

	sysInfo.WriteString(fmt.Sprintf("%-8v", "CPU:"))
	for cpuName, cpuCores := range cpus {
		sysInfo.WriteString(fmt.Sprintf("%20v\n", fmt.Sprintf("%s x %v", cpuName, cpuCores)))
	}

	sysInfo.WriteString(fmt.Sprintf("%-8v%v", "RAM:", byteCountIEC(ram)))
	systemInfoBox.Print("System Info", sysInfo.String())

	inst := installer{
		os:             hostStat.OS,
		platformFamily: hostStat.Platform,
		arch:           hostStat.KernelArch,
	}

	for programName, execName := range prerequisites {
		fmt.Printf("Checking %s...", programName)
		if isProgramInstalled := programExists(execName); isProgramInstalled {
			fmt.Println(" present")
			continue
		}

		fmt.Printf(" not present.\nDownloading installer for %s...\n", programName)
		switch programName {
		case "Git":
			inst.downloadGit()
		case "JDK":
			inst.downloadJDK()
		case "Android Studio":
			inst.downloadAndroidStudio()
		case "Flutter":
			inst.downloadFlutter()
		default:
			if inst.os != "darwin" && programName != "Homebrew" {
				continue
			}
			inst.downloadHomebrew()
		}
	}
}
