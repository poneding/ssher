/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package util

import (
	"fmt"
	"os"
)

// UserHomeDirOrDie returns the user's home directory.
// If the home directory cannot be determined, the program will exit.
func UserHomeDirOrDie() string {
	path, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("✗ Error:", err)
		os.Exit(0)
	}
	return path
}
