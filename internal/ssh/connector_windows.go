//go:build windows
// +build windows

/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/

package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Connect to the ssh server.
func Connect(profile *Profile) {
	cmd := exec.Command("ssh", "-p", strconv.Itoa(profile.Port))
	if profile.PrivateKey != "" {
		if profile.PrivateKey != "" {
			cmd.Args = append(cmd.Args, "-i", profile.PrivateKey)
		}
	}

	if profile.User != "" {
		cmd.Args = append(cmd.Args, profile.User+"@"+profile.Host)
	} else {
		cmd.Args = append(cmd.Args, profile.Host)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("✗ Failed to run ssh:", err)
		os.Exit(0)
	}
	cmd.Wait()
}
