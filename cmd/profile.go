/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// profileCmd represents the profile file operations command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "ssh profile file operations.",
	Long:  `ssh profile file operations.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
