/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/poneding/ssher/internal/output"
	"github.com/spf13/cobra"
)

const (
	GH_REPO              = "github.com/poneding/ssher"
	GH_ADDR              = "https://" + GH_REPO
	GH_RELEASE_ADDR_BASE = GH_ADDR + "/releases"
	GH_API_ADDR_BASE     = "https://api.github.com/repos/poneding/ssher"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade ssher to the latest version or specific version.",
	Long:  `Upgrade ssher to the latest version or specific version.`,
	Run: func(cmd *cobra.Command, args []string) {
		runUpgrade()
	},
}

var targetVersion string

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringVarP(&targetVersion, "version", "v", "", "upgrade to specific version")
}

type getReleaseResult struct {
	TagName string `json:"tag_name"`
}

func runUpgrade() {
	if targetVersion == "" {
		targetVersion = "latest"
	}

	if targetVersion == version && targetVersion != "latest" {
		output.Done("Current version: %s, upgradtion ignored.", version)
		return
	}

	// remote version check
	suffix := func() string {
		if targetVersion == "latest" {
			return targetVersion
		} else {
			return "tags/" + targetVersion
		}
	}()
	r, err := http.Get(fmt.Sprintf("%s/releases/%s", GH_API_ADDR_BASE, suffix))
	if err != nil {
		output.Fatal("Failed to get version: %s", err)
	}
	if r.StatusCode != http.StatusOK {
		output.Fatal("Failed to get version: %s", r.Status)
	}
	defer r.Body.Close()

	result := getReleaseResult{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		output.Fatal("Failed to get version: %s", err)
	}
	if result.TagName == "" {
		output.Fatal("Failed to get version: empty")
	}

	targetVersion = result.TagName
	if targetVersion == version {
		output.Done("Current version: %s, you are up to date.", version)
		return
	}

	output.Note("Version %s is available, upgrading now...", targetVersion)

	upgrade()
}

// upgrade ssher to target version, go install is recommended
func upgrade() {
	_, err := exec.LookPath("go")
	if err == nil {
		output.Note("Upgrading by go install...")
		upgradeByGoInstall(targetVersion)
	} else {
		output.Note("Downloading ssher from github release...")
		upgradeByDownload(targetVersion)
	}
	output.Done("Upgraded to version: %s", targetVersion)
}

// upgradeByGoInstall upgrade by go install command
func upgradeByGoInstall(v string) {
	if targetVersion == "" {
		targetVersion = "latest"
	}

	cmd := exec.Command("go", "install", GH_REPO+"@"+v)
	err := cmd.Run()
	if err != nil {
		output.Fatal("Upgrade failed by go install: %s", err)
	}
}

// upgradeByDownload download binary from github release and replace the current binary
func upgradeByDownload(v string) {
	if targetVersion == "" {
		targetVersion = "latest"
	}

	var ext string
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	targetFile := fmt.Sprintf("ssher_%s_%s_%s%s", strings.Trim(v, "v"), runtime.GOOS, runtime.GOARCH, ext)

	// download from github release. eg: https://github.com/poneding/ssher/releases/download/v1.1.0/ssher_1.1.0_linux_amd64
	r, err := http.Get(fmt.Sprintf("%s/download/%s/%s", GH_RELEASE_ADDR_BASE, v, targetFile))
	if err != nil || r.StatusCode != http.StatusOK {
		output.Fatal("Failed to download: %s", err)
	}
	defer r.Body.Close()

	// get ssher path
	ssherPath, err := os.Executable()
	if err != nil {
		output.Fatal("Failed to get ssher path: %s", err)
	}

	binaryData, err := io.ReadAll(r.Body)
	if err != nil {
		output.Fatal("Failed to read all: %s", err)
	}

	// write binaryData to targetFile and set permission
	err = os.WriteFile(targetFile, binaryData, 0755)
	if err != nil {
		output.Fatal("Failed to write file: %s", err)
	}

	// mv targetFile to ssherPath
	err = os.Rename(targetFile, ssherPath)
	if err != nil {
		output.Fatal("Failed to rename: %s", err)
	}
}
