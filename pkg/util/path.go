/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package util

import (
	"os"

	"github.com/poneding/ssher/internal/output"
)

// UserHomeDirOrDie returns the user's home directory.
// If the home directory cannot be determined, the program will exit.
func UserHomeDirOrDie() string {
	path, err := os.UserHomeDir()
	if err != nil {
		output.Fatal("Failed to get user home directory: %s", err)
	}
	return path
}
