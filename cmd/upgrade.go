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

	"github.com/spf13/cobra"
)

const (
	GH_REPO                    = "github.com/poneding/ssher"
	GH_PROXY                   = "https://ghproxy.ketches.cn/"
	GH_ADDR                    = "https://" + GH_REPO
	GH_RELEASE_ADDR_BASE       = GH_ADDR + "/releases"
	GH_RELEASE_PROXY_ADDR_BASE = GH_PROXY + GH_RELEASE_ADDR_BASE
	GH_API_PROXY_ADDR_BASE     = GH_PROXY + "https://api.github.com/repos/poneding/ssher"
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

	if targetVersion == version {
		fmt.Printf("✓ Current version: %s, upgradtion ignored.\n", version)
	}

	// remote version check
	suffix := func() string {
		if targetVersion == "latest" {
			return targetVersion
		} else {
			return "tags/" + targetVersion
		}
	}()
	r, err := http.Get(fmt.Sprintf("%s/releases/%s", GH_API_PROXY_ADDR_BASE, suffix))
	if err != nil {
		fmt.Println("✗ Failed to get version:", err)
		os.Exit(0)
	}
	if r.StatusCode != http.StatusOK {
		fmt.Println("✗ Failed to get version:", r.Status)
		os.Exit(0)
	}
	defer r.Body.Close()

	result := getReleaseResult{}
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		fmt.Println("✗ Failed to get version:", err)
		os.Exit(0)
	}
	if result.TagName == "" {
		fmt.Println("✗ Failed to get version: empty")
		os.Exit(0)
	}

	targetVersion = result.TagName
	fmt.Printf("✓ Version %s is available, upgrading now...\n", targetVersion)

	upgrade()
}

// upgrade ssher to target version, go install is recommended
func upgrade() {
	_, err := exec.LookPath("go")
	if err == nil {
		fmt.Println("✓ Upgrading by go install...")
		upgradeByGoInstall(targetVersion)
	} else {
		fmt.Println("✓ Downloading ssher from github release...")
		upgradeByDownload(targetVersion)
	}
	fmt.Println("✓ Upgraded to version:", targetVersion)
}

// upgradeByGoInstall upgrade by go install command
func upgradeByGoInstall(v string) {
	if targetVersion == "" {
		targetVersion = "latest"
	}

	cmd := exec.Command("go", "install", GH_REPO+"@"+v)
	err := cmd.Run()
	if err != nil {
		fmt.Println("✗ Upgrade failed by go install:", err)
		os.Exit(0)
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

	// download from github release. eg: https://github.com/poneding/ssher/releases/download/v0.1.0/ssher_0.1.0_linux_amd64
	r, err := http.Get(fmt.Sprintf("%s/download/%s/%s", GH_RELEASE_PROXY_ADDR_BASE, v, targetFile))
	if err != nil || r.StatusCode != http.StatusOK {
		fmt.Println("✗ Failed to download:", err)
		os.Exit(0)
	}
	defer r.Body.Close()

	// get ssher path
	ssherPath, err := os.Executable()
	if err != nil {
		fmt.Println("✗ Failed to get ssher path:", err)
		os.Exit(0)
	}

	binaryData, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("✗ Failed to read all:", err)
		os.Exit(0)
	}

	// write binaryData to targetFile and set permission
	err = os.WriteFile(targetFile, binaryData, 0755)
	if err != nil {
		fmt.Println("✗ Failed to write file:", err)
		os.Exit(0)
	}

	// mv targetFile to ssherPath
	err = os.Rename(targetFile, ssherPath)
	if err != nil {
		fmt.Println("✗ Failed to rename:", err)
		os.Exit(0)
	}
}
