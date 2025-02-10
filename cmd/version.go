/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ssher/internal/output"
	"github.com/spf13/cobra"
)

const version = "v1.1.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "ssher version.",
	Long:  `ssher version`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Done("Version: %s", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
