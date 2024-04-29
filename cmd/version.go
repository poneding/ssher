/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v1.0.2"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "ssher version.",
	Long:  `ssher version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("✓ Version: %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
