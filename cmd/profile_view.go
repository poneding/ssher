/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// profileViewCmd represents the view command
var profileViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View the ssh profile file.",
	Long:  `View the ssh profile file.`,
	Run: func(cmd *cobra.Command, args []string) {
		runProfileView()
	},
}

func init() {
	profileCmd.AddCommand(profileViewCmd)
}

func runProfileView() {
	file := viper.ConfigFileUsed()
	if file == "" {
		fmt.Println("✗ No profile file found.")
		os.Exit(0)
	}
	b, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("✗ Failed to read %s: %v\n", file, err)
		os.Exit(0)
	}
	fmt.Println(string(b))
}
