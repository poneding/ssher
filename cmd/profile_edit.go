/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// profileEditCmd represents the edit command
var profileEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the ssh profile file.",
	Long:  `Edit the ssh profile file. The default editor is vim or vi on Unix-like systems, and notepad on Windows.`,
	Run: func(cmd *cobra.Command, args []string) {
		runProfileEdit()
	},
}

func init() {
	profileCmd.AddCommand(profileEditCmd)
}

func runProfileEdit() {
	file := viper.ConfigFileUsed()
	if file == "" {
		fmt.Println("✗ No profile file found.")
		os.Exit(0)
	}
	switch runtime.GOOS {
	case "windows":
		// use notepad to edit
		if err := exec.Command("notepad", file).Run(); err != nil {
			fmt.Printf("✗ Failed to open %s with notepad: %v\n", file, err)
		}
	default:
		// use vi or vim to edit
		textEditor := "vim" // default to vim
		if _, err := exec.LookPath("vim"); err != nil {
			textEditor = "vi"
		}
		cmd := exec.Command(textEditor, file)
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Printf("✗ Failed to open %s with vi: %v\n", file, err)
		}
	}
	os.Exit(0)
}
