/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/poneding/ssher/internal/output"
	"github.com/poneding/ssher/internal/ssh"
	"github.com/spf13/cobra"
)

// var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssher",
	Short: "ssher is an easy-to-use command line tool for connecting to remote servers.",
	Long:  `ssher is an easy-to-use command line tool for connecting to remote servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		runConnect()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("✗ Error:", err)
		os.Exit(0)
	}
}

var targetServer string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&targetServer, "server", "s", "", "target server to connect. (ssh name)")
	rootCmd.RegisterFlagCompletionFunc("server", completeWithServers)
}

func initConfig() {
	ssh.TryCreateConfigurationFile()
}

func runConnect() {
	var server *ssh.Server

	if targetServer != "" {
		server = ssh.GetServer(targetServer)
		if server == nil {
			output.Fatal("No server found.")
		}
	} else {
		server = ssh.SelectPrompt(ssh.ConnectPromptLable)
	}

	// reset current server
	if current := ssh.GetCurrentServer(); current != nil {
		current.Current = false
	}
	server.Current = true
	ssh.SaveConfigurationFile()

	// ssh connect
	ssh.Connect(server)
}

func completeWithServers(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	servers := ssh.GetServers()

	var completions []string
	for _, s := range servers {
		completions = append(completions, fmt.Sprintf("%s\t[%s] %s:%d - %s", s.Name, s.Emoji, s.Host, s.Port, s.User))
	}

	return completions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}
