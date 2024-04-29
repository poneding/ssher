//go:build !windows
// +build !windows

/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package ssh

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/poneding/ssher/pkg/util"
)

type PromptLable string

const (
	ConnectPromptLable PromptLable = "Select a ssh profile to connect: "
	RemovePromptLable  PromptLable = "Select a ssh profile to remove: "
)

type PromptOperation string

const (
	AddSSHProfilePromptOperation PromptOperation = "+ Add a new ssh profile"
	ExitPromptOperation          PromptOperation = "✗ Exit"
)

// SelectPrompt prompts the user to select a ssh profile, and returns the selected profile.
func SelectPrompt(label PromptLable) (selected *Profile) {
	profiles := GetProfiles()
	cursorPos := 0
	for i, p := range profiles {
		if p.Current {
			cursorPos = i
		}
	}

	prompt := promptui.Select{
		// Label: label,
		Items: append(profiles, &Profile{Name: string(AddSSHProfilePromptOperation)}, &Profile{Name: string(ExitPromptOperation)}),
		Searcher: func(input string, index int) bool {
			if index < 0 || index >= len(profiles) {
				return false
			}

			current := profiles[index]
			if strings.Contains(strings.ToLower(current.Name), strings.ToLower(input)) ||
				strings.Contains(strings.ToLower(current.Host), strings.ToLower(input)) {
				return true
			}

			return false
		},
		CursorPos: cursorPos,
		Templates: &promptui.SelectTemplates{
			Label:    string(label),
			Active:   promptui.Styler(promptui.FGRed, promptui.FGBold)("➜  {{if .Current}}{{ \"* \" }}{{end}}{{ .Name }}\t{{ .Host }}"),
			Inactive: promptui.Styler(promptui.FGCyan)("  {{if .Current}}{{ \"* \" }}{{end}}{{ .Name }}\t{{ .Host }}"),
			Selected: promptui.Styler(promptui.FGCyan)("✔ Seletcd: {{ .Name }}\t{{ .Host }}"),
			Details: `
---------- Profile ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Host:" | faint }}	{{ .Host }}
{{ "Port:" | faint }}	{{ .Port }}
{{ "User:" | faint }}	{{ .User }}
{{ "PrivateKey:" | faint }}	{{ .PrivateKey }}`,
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println("✗ Prompt failed:", err)
		os.Exit(0)
	}

	// AddSSHProfilePromptOperation selected
	if index == len(profiles) {
		profile := FormPrompt()
		AddProfile(profile)
		os.Exit(0)
	}

	// ExitPromptOperation selected
	if index == len(profiles)+1 {
		os.Exit(0)
	}

	if label == RemovePromptLable {
		prompt := promptui.Prompt{
			Label:   promptui.Styler(promptui.FGRed)("⚠ Are you sure to remove this profile? [y/n]"),
			Default: "n",
		}
		c, err := prompt.Run()
		if err != nil {
			fmt.Println("✗ Prompt failed:", err)
			os.Exit(0)
		}
		if strings.ToLower(c) != "y" {
			os.Exit(0)
		}
	}

	selected = profiles[index]
	return
}

// FormPrompt prompts the user to fill in the ssh profile fields, and returns the new profile.
func FormPrompt() *Profile {
	profiles := GetProfiles()
	var newProfile = new(Profile)

	fmt.Println("Add a new ssh profile, fill in the following fields (q to quit):")

	var promptLabels = []string{"Name(*)", "Host(*)", "Port", "User", "Password", "PrivateKey"}
	m := make(map[string]string, len(promptLabels))

	for _, p := range promptLabels {
		prompt := promptui.Prompt{
			Label: p,
		}
		switch p {
		case "Name":
			prompt.Validate = func(input string) error {
				input = strings.Trim(input, " ")
				if input == "" {
					return fmt.Errorf("ssh profile %s is required", p)
				}

				for _, p := range profiles {
					if p.Name == input {
						return fmt.Errorf("ssh profile name %s already exists", input)
					}
				}
				return nil
			}
		case "Host":
			prompt.Validate = func(input string) error {
				input = strings.Trim(input, " ")
				if input == "" {
					return fmt.Errorf("ssh profile %s is required", p)
				}
				return nil
			}
		case "Password":
			prompt.Mask = '*'
		case "Port":
			prompt.Default = " 22"
			prompt.Validate = func(input string) error {
				input = strings.Trim(input, " ")
				if input == "q" {
					os.Exit(0)
				}
				_, err := strconv.Atoi(input)
				return err
			}
		case "PrivateKey":
			prompt.Validate = func(input string) error {
				input = strings.Trim(input, " ")
				if input == "" {
					return nil
				}
				if input == "q" {
					os.Exit(0)
				}

				// compatible with home path with ~, e.g. ~/.ssh/id_rsa
				if strings.HasPrefix(input, "~") {
					input = strings.Replace(input, "~", util.UserHomeDirOrDie(), 1)
				}
				_, err := os.Stat(input)
				return err
			}
		}
		value, _ := prompt.Run()
		if value == "q" {
			os.Exit(0)
		}
		m[p] = strings.Trim(value, " ")
	}

	newProfile.Name = m["Name(*)"]
	newProfile.Host = m["Host(*)"]
	newProfile.Port, _ = strconv.Atoi(m["Port"])
	newProfile.User = m["User"]
	newProfile.Password = m["Password"]
	newProfile.PrivateKey = m["PrivateKey"]

	return newProfile
}
