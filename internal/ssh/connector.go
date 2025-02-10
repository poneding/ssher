//go:build !windows
// +build !windows

/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/

package ssh

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/creack/pty"
	"github.com/poneding/ssher/internal/output"
	"golang.org/x/term"
)

// Connect to the server.
func Connect(server *Server) {
	cmd := exec.Command("ssh", "-p", strconv.Itoa(server.Port))
	passwordSSH := server.Password != ""
	if server.IdentityFile != "" {
		cmd.Args = append(cmd.Args, "-i", server.IdentityFile)
	}

	if server.User != "" {
		cmd.Args = append(cmd.Args, server.User+"@"+server.Host)
	} else {
		cmd.Args = append(cmd.Args, server.Host)
	}

	ptmx, err := pty.Start(cmd)
	if err != nil {
		output.Fatal("Failed to start ssh: %s", err)
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("failed to resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		output.Fatal("Failed to set stdin in raw mode: %s", err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()

	for i := 0; passwordSSH && i < 3; i++ {
		// Read the pty until wait for password.
		buf := make([]byte, 1024*4)
		n, err := ptmx.Read(buf)
		if err != nil {
			output.Fatal("Failed to read pty: %s", err)
		}
		if bytes.HasPrefix(buf[:n], []byte("The authenticity of host ")) {
			// Write yes to the pty.
			_, err = ptmx.Write([]byte("yes\n"))
			if err != nil {
				output.Fatal("Failed to write yes to pty: %s", err)
			}
			continue
		}
		if bytes.HasSuffix(buf[:n], []byte("password: ")) {
			// Write the password to the pty.
			decodedPwd, err := base64.StdEncoding.DecodeString(server.Password)
			if err != nil {
				output.Fatal("Failed to decode password: %s", err)
			}
			_, err = ptmx.Write(append(decodedPwd, []byte("\n")...))
			if err != nil {
				output.Fatal("Failed to write password to pty: %s", err)
			}
			break
		}
	}
	_, _ = io.Copy(os.Stdout, ptmx)
}
