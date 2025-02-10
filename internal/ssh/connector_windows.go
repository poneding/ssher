//go:build windows
// +build windows

/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/

package ssh

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/poneding/ssher/internal/output"
)

// Connect to the ssh server.
func Connect(server *Server) {
	cmd := exec.Command("ssh", "-p", strconv.Itoa(server.Port))
	if server.IdentityFile != "" {
		cmd.Args = append(cmd.Args, "-i", server.IdentityFile)
	}

	if server.User != "" {
		cmd.Args = append(cmd.Args, server.User+"@"+server.Host)
	} else {
		cmd.Args = append(cmd.Args, server.Host)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		output.Fatal("Failed to run ssh:", err)
	}
	cmd.Wait()
}
