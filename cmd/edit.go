/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ssher/internal/ssh"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a server.",
	Long:  `Edit a server, you will be prompted to input the server name, host, port, user, password and private key file path.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runEditServer(args)
	},
	ValidArgsFunction: completeWithServers,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEditServer(args []string) {
	var server *ssh.Server
	if len(args) == 0 {
		server = ssh.SelectPrompt(ssh.EditPromptLable)
	} else {
		server = ssh.GetServer(args[0])
	}

	if server != nil {
		newServer := ssh.EditFormPrompt(server)
		ssh.UpdateServer(server.Name, newServer)
	}
}
