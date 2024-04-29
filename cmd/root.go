/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/poneding/ssher/internal/ssh"
	"github.com/poneding/ssher/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssher",
	Short: "ssher is a lightweight ssh profile cli manager.",
	Long:  `ssher is a lightweight ssh profile cli manager.`,
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

var nameSSHProfileConnect string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ssher.yaml)")

	rootCmd.Flags().StringVarP(&nameSSHProfileConnect, "name", "n", "", "ssh profile to connect. (ssh name)")
	rootCmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		profiles := ssh.GetProfiles()
		var completions []string
		for _, p := range profiles {
			if p.Name == toComplete {
				return []string{}, cobra.ShellCompDirectiveNoFileComp
			}
			if p.Name != "" {
				completions = append(completions, p.Name)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".ssher" (without extension).
		viper.AddConfigPath(util.UserHomeDirOrDie())
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ssher")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); errors.As(err, &viper.ConfigFileNotFoundError{}) {
		file := viper.ConfigFileUsed()
		if file == "" {

			file = util.UserHomeDirOrDie() + "/.ssher.yaml"
		}
		ssh.CreateProfileFile(file)
	}
}

func runConnect() {
	var profile *ssh.Profile

	if nameSSHProfileConnect != "" {
		fmt.Printf("✓ SSH profile to connect: %s\n", nameSSHProfileConnect)
		profile = ssh.GetProfile(nameSSHProfileConnect)
		if profile == nil {
			fmt.Println("✗ No such profile found.")
			os.Exit(0)
		}
	} else {
		profile = ssh.SelectPrompt(ssh.ConnectPromptLable)
	}

	// reset current profile
	if current := ssh.GetCurrentProfile(); current != nil {
		current.Current = false
	}
	profile.Current = true
	ssh.SaveProfiles()

	// ssh connect
	ssh.Connect(profile)
}
