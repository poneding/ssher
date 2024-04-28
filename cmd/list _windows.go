//go:build windows
// +build windows

/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/poneding/ssher/internal/ssh"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all ssh profiles.",
	Long:    `List all ssh profiles.`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		runList()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList() {
	profiles := ssh.GetProfiles()
	if len(profiles) == 0 {
		fmt.Println("No profiles found, add a new ssh profile with `ssher add`.")
		return
	}

	var paddingLen int
	for _, p := range profiles {
		if len(p.Name) > paddingLen {
			paddingLen = len(p.Name)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"", "Name", "Host", "Port", "User", "Private Key"})

	current := func(p *ssh.Profile) string {
		if p.Current {
			return "*"
		}
		return " "
	}

	for _, p := range profiles {
		t.AppendRow(table.Row{current(p), p.Name, p.Host, p.Port, p.User, p.PrivateKey})
	}
	t.Render()
}
