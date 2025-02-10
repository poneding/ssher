/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"os"
	"runtime"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/poneding/ssher/internal/output"
	"github.com/poneding/ssher/internal/ssh"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all servers.",
	Long:    `List all servers.`,
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		runListServers()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runListServers() {
	servers := ssh.GetServers()
	if len(servers) == 0 {
		output.Note("No servers found, add a new server with `ssher add`.")
		return
	}

	var paddingLen int
	for _, s := range servers {
		if len(s.Name) > paddingLen {
			paddingLen = len(s.Name)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if runtime.GOOS != "windows" {
		t.AppendHeader(table.Row{"", "NAME", "HOST", "PORT", "USER", "PASSWORD", "IDENTITY_FILE"})
	} else {
		t.AppendHeader(table.Row{"", "NAME", "HOST", "PORT", "USER", "IDENTITY_FILE"})
	}

	current := func(s *ssh.Server) string {
		if s.Current {
			return "✦"
		}
		return " "
	}
	password := func(s *ssh.Server) string {
		if s.Password != "" {
			return "******"
		}
		return ""
	}

	for _, s := range servers {
		if runtime.GOOS != "windows" {
			t.AppendRow(table.Row{current(s), s.Name, s.Host, s.Port, s.User, password(s), s.IdentityFile})
		} else {
			t.AppendRow(table.Row{current(s), s.Name, s.Host, s.Port, s.User, s.IdentityFile})
		}
	}
	t.SetStyle(table.StyleRounded)
	t.Render()
}
