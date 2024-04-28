//go:build !windows
// +build !windows

/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/

package ssh

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

// Connect to the ssh server.
func Connect(profile *Profile) {
	cmd := exec.Command("ssh", "-p", strconv.Itoa(profile.Port))
	passwordSSH := profile.Password != ""
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

	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Println("✗ Failed to start ssh:", err)
		os.Exit(0)
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	if runtime.GOOS != "windows" {
		//
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		go func() {
			for range ch {
				if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
					log.Printf("error resizing pty: %s", err)
				}
			}
		}()
		ch <- syscall.SIGWINCH                        // Initial resize.
		defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.
	}

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("✗ Failed to set stdin in raw mode:", err)
		os.Exit(0)
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
			fmt.Println("✗ Failed to read pty:", err)
			os.Exit(0)
		}
		if bytes.HasPrefix(buf[:n], []byte("The authenticity of host ")) {
			// Write yes to the pty.
			_, err = ptmx.Write([]byte("yes\n"))
			if err != nil {
				fmt.Println("✗ Failed to write yes to pty:", err)
				os.Exit(0)
			}
			continue
		}
		if bytes.HasSuffix(buf[:n], []byte("password: ")) {
			// Write the password to the pty.
			_, err = ptmx.Write([]byte(profile.Password + "\n"))
			if err != nil {
				fmt.Println("✗ Failed to write password to pty:", err)
				os.Exit(0)
			}
		}
		break
	}

	_, _ = io.Copy(os.Stdout, ptmx)
}
