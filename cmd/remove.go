/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/poneding/ssher/internal/ssh"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a ssh profile.",
	Long:    "Remove a ssh profile, you will be prompted to select a profile to remove or just input the profile name after `--name` or `-n",
	Aliases: []string{"rm"},
	Run: func(cmd *cobra.Command, args []string) {
		runRemove()
	},
}

var nameSSHProfileRemoved string

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringVarP(&nameSSHProfileRemoved, "name", "n", "", "ssh profile to remove (ssh name)")
	removeCmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		profiles := ssh.GetProfiles()
		var completions []string
		for _, p := range profiles {
			if p.Name == toComplete || p.Host == toComplete {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			if p.Name != "" {
				completions = append(completions, p.Name)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})
}

func runRemove() {
	var profile *ssh.Profile
	if nameSSHProfileRemoved != "" {
		fmt.Printf("✓ SSH profile to remove: %s\n", nameSSHProfileRemoved)
		profile = ssh.GetProfile(nameSSHProfileRemoved)
		if profile == nil {
			fmt.Printf("✗ Profile %s not found\n", nameSSHProfileRemoved)
			os.Exit(0)
		}
	} else {
		profile = ssh.SelectPrompt(ssh.RemovePromptLable)
	}

	ssh.RemoveProfile(profile)
}
