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
	Short: "Add a ssh profile.",
	Long:  `Add a ssh profile, you will be prompted to input the profile name, host, port, user, password and private key file path.`,
	Run: func(cmd *cobra.Command, args []string) {
		profile := ssh.FormPrompt()
		ssh.AddProfile(profile)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
