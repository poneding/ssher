/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ssher/internal/ssh"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new server.",
	Long:  `Add a new server, you will be prompted to input the server name, host, port, user, password and private key file path.`,
	Run: func(cmd *cobra.Command, args []string) {
		runAddServer()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAddServer() {
	server := ssh.AddFormPrompt()
	ssh.AddServer(server)
}
