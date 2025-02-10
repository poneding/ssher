/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ssher/internal/ssh"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a server.",
	Long:    "Remove a server, you will be prompted to select a server to remove or just input the server name after `--name` or `-n",
	Aliases: []string{"rm"},
	Run: func(cmd *cobra.Command, args []string) {
		runRemove(args)
	},
	ValidArgsFunction: completeWithServers,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(args []string) {
	var servers []*ssh.Server
	if len(args) == 0 {
		servers = append(servers, ssh.SelectPrompt(ssh.RemovePromptLable))
	} else {
		for _, arg := range args {
			server := ssh.GetServer(arg)
			if server != nil {
				servers = append(servers, server)
			}
		}
	}

	for _, dst := range servers {
		ssh.RemoveServer(dst)
	}
}
