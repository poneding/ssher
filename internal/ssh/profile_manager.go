/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package ssh

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type config struct {
	Profiles []*Profile `json:"profiles"`
}

var c *config

// readProfiles reads the profiles from the config file
func readProfiles() {
	if c == nil {
		c = new(config)
		err := viper.Unmarshal(c)
		if err != nil {
			fmt.Printf("✗ Unable to decode into struct %v", err)
			os.Exit(0)
		}
	}
}

// AddProfile adds a new ssh profile
func AddProfile(profile *Profile) {
	defer func() {
		SaveProfiles()
		fmt.Printf("✔ A ssh profile added: %s\t%s\n", profile.Name, profile.Host)
	}()

	readProfiles()

	if profile == nil {
		fmt.Println("✗ ssh profile is required")
		os.Exit(0)
	}

	if profile.Name == "" {
		fmt.Println("✗ ssh profile name is required for a profile")
		os.Exit(0)
	}

	if profile.Host == "" {
		fmt.Println("✗ ssh profile host is required for a profile")
		os.Exit(0)
	}

	if profile.Port == 0 {
		profile.Port = 22
	}

	if len(c.Profiles) == 0 {
		profile.Current = true
	}
	c.Profiles = append(c.Profiles, profile)
}

// GetProfiles returns all the profiles
func GetProfiles() []*Profile {
	readProfiles()
	return c.Profiles
}

// GetProfile returns a profile by name
func GetProfile(name string) *Profile {
	readProfiles()
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			return c.Profiles[i]
		}
	}
	return nil
}

// GetCurrentProfile returns the current profile
func GetCurrentProfile() *Profile {
	readProfiles()
	for i := range c.Profiles {
		if c.Profiles[i].Current {
			return c.Profiles[i]
		}
	}
	return nil
}

// RemoveProfile removes a profile
func RemoveProfile(profile *Profile) {
	defer func() {
		SaveProfiles()
		fmt.Printf("✔ A ssh profile removed: %s\t%s\n", profile.Name, profile.Host)
	}()

	readProfiles()

	for i := range c.Profiles {
		if c.Profiles[i].Name == profile.Name {
			c.Profiles = append(c.Profiles[:i], c.Profiles[i+1:]...)
			return
		}
	}
}

// CreateProfileFile creates ssh profile file
func CreateProfileFile(file string) {
	f, err := os.Create(file)
	if err != nil {
		fmt.Printf("✗ Unable to create profile file, %v", err)
		os.Exit(0)
	}
	defer f.Close()
}

// SaveProfiles save the profiles to the profile file
func SaveProfiles() {
	readProfiles()
	viper.Set("profiles", c.Profiles)

	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("✗ Unable to write to profile file, %v", err)
		os.Exit(0)
	}
}
