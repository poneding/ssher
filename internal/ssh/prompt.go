/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package ssh

import (
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/poneding/ssher/internal/output"
	"github.com/poneding/ssher/pkg/util"
)

const (
	ConnectPromptLable = "Select a server to connect"
	RemovePromptLable  = "Select a server to remove"
	EditPromptLable    = "Select a server to edit"
)

const (
	AddServerPromptOperation = "Add a new server"
	ExitPromptOperation      = "Exit"
)

func Comfirm(label string) bool {
	templates := &promptui.SelectTemplates{
		Label:    promptui.Styler(promptui.FGYellow)("# {{ . }}?"),
		Active:   promptui.Styler(promptui.FGCyan, promptui.FGUnderline)("➤ {{ . }}"),
		Inactive: promptui.Styler(promptui.FGFaint)("  {{ . }}"),
	}
	prompt := promptui.Select{
		Label:        label,
		Items:        []string{"No", "Yes"},
		Templates:    templates,
		Size:         4,
		HideSelected: true,
	}
	_, obj, err := prompt.Run()
	if err != nil {
		output.Fatal("Prompt failed %v", err)
	}

	return obj == "Yes"
}

// SelectPrompt prompts the user to select a server, and returns the selected server.
func SelectPrompt(label string) (selected *Server) {
	servers := GetServers()
	cursorPos := 0
	for i, s := range servers {
		if s.Current {
			cursorPos = i
		}
	}

	prompt := promptui.Select{
		Label: label,
		Items: append(servers, &Server{Name: AddServerPromptOperation, Emoji: "+"}, &Server{Name: ExitPromptOperation, Emoji: "✗"}),
		Searcher: func(input string, index int) bool {
			if index < 0 || index >= len(servers) {
				return false
			}

			current := servers[index]
			if strings.Contains(strings.ToLower(current.Name), strings.ToLower(input)) ||
				strings.Contains(strings.ToLower(current.Host), strings.ToLower(input)) {
				return true
			}

			return false
		},
		HideSelected: true,
		CursorPos:    cursorPos,
		Templates: &promptui.SelectTemplates{
			Label:    promptui.Styler(promptui.FGYellow)("# {{ . }}:"),
			Active:   promptui.Styler(promptui.FGCyan, promptui.FGUnderline)("➤ {{ .Emoji }} {{ .Name }}"),
			Inactive: promptui.Styler(promptui.FGFaint)("  {{ .Emoji }} {{ .Name }}"),
			Details: `{{if .Host}}
---------- Server ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Host:" | faint }}	{{ .Host }}
{{ "Port:" | faint }}	{{ .Port }}
{{ "User:" | faint }}	{{ .User }}
{{ "IdentityFile:" | faint }}	{{ .IdentityFile }}{{end}}`,
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		fmt.Println("✗ Prompt failed:", err)
		os.Exit(0)
	}

	// AddServerPromptOperation selected
	if index == len(servers) {
		server := AddFormPrompt()
		AddServer(server)
		os.Exit(0)
	}

	// ExitPromptOperation selected
	if index == len(servers)+1 {
		os.Exit(0)
	}

	if label == RemovePromptLable {
		if !Comfirm("Are you sure to remove this server") {
			os.Exit(0)
		}
	}

	selected = servers[index]
	return
}

func FormPrompt(server *Server) *Server {
	servers := GetServers()
	var newServer = new(Server)

	output.Note("Add a new server, fill in the following fields (q to quit):")

	var promptLabels = []string{"Name(*)", "Host(*)", "Port(*)", "User(*)", "Password", "IdentityFile"}
	m := make(map[string]string, len(promptLabels))

	for _, label := range promptLabels {
		prompt := promptui.Prompt{
			Label: promptui.Styler(promptui.FGBlue)(label),
		}
		switch label {
		case "Name(*)":
			prompt.Default = server.Name
			prompt.AllowEdit = true
			prompt.Validate = func(input string) error {
				input = strings.TrimSpace(input)
				if input == "" {
					return fmt.Errorf("Server name is required")
				}

				for _, s := range servers {
					if s.Name == input && s.Name != server.Name {
						return fmt.Errorf("Server name %s already exists", input)
					}
				}
				return nil
			}
		case "Host(*)":
			prompt.Default = server.Host
			prompt.AllowEdit = true
			prompt.Validate = func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("Server host is required")
				}
				return nil
			}
		case "Port(*)":
			prompt.Default = strconv.Itoa(server.Port)
			prompt.AllowEdit = true
			prompt.Validate = func(input string) error {
				input = strings.TrimSpace(input)
				if input == "q" {
					os.Exit(0)
				}
				_, err := strconv.Atoi(input)
				return err
			}
		case "User(*)":
			prompt.Default = server.User
			prompt.AllowEdit = true
			prompt.Validate = func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("Server user is required")
				}
				return nil
			}
		case "Password":
			if runtime.GOOS == "windows" {
				continue
			}
			prompt.Mask = '*'
		case "IdentityFile":
			prompt.Default = server.IdentityFile
			prompt.AllowEdit = true
			prompt.Validate = func(input string) error {
				input = strings.TrimSpace(input)
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
		m[label] = strings.TrimSpace(value)
	}

	newServer.Name = m["Name(*)"]
	newServer.Host = m["Host(*)"]
	newServer.Port, _ = strconv.Atoi(m["Port(*)"])
	newServer.User = m["User(*)"]
	if runtime.GOOS != "windows" {
		newServer.Password = base64.StdEncoding.EncodeToString([]byte(m["Password"]))
	}
	newServer.IdentityFile = m["IdentityFile"]

	return newServer
}

// AddFormPrompt prompts the user to fill in the server
// fields, and returns the new server.
func AddFormPrompt() *Server {
	defaultServer := &Server{
		Port:         22,
		User:         "root",
		IdentityFile: "~/.ssh/id_rsa",
	}

	return FormPrompt(defaultServer)
}

// EditFormPrompt prompts the user to modify the server
// fields, and returns the new server.
func EditFormPrompt(server *Server) *Server {
	return FormPrompt(server)
}
