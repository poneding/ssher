/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package ssh

import (
	"bytes"
	"os"

	"github.com/poneding/ssher/internal/output"
	"github.com/poneding/ssher/pkg/util"
	"gopkg.in/yaml.v2"
)

type config struct {
	Servers []*Server `yaml:"servers"`
}

var c *config

func ConfigFile() string {
	return util.UserHomeDirOrDie() + "/.ssher.yaml"
}

// readServers reads the servers from the configuration file
func readServers() {
	if c == nil {
		configContent, err := os.ReadFile(ConfigFile())
		if err != nil {
			output.Fatal("Unable to read configuration file, %v", err)
		}

		c = new(config)
		if err := yaml.Unmarshal(configContent, &c); err != nil {
			output.Fatal("Unable to unmarshal configuration file, %v", err)
		}

		for _, s := range c.Servers {
			s.Emoji = " "
			if s.Current {
				s.Emoji = "✦"
			}
		}
	}
}

// AddServer adds a new server
func AddServer(server *Server) {
	defer func() {
		SaveConfigurationFile()
		output.Done("Server %s(%s) added.", server.Name, server.Host)
	}()

	readServers()

	if server == nil {
		output.Fatal("Server is required.")
	}

	if server.Name == "" {
		output.Fatal("Server name is required.")
	}

	if server.Host == "" {
		output.Fatal("Server host is required.")
	}

	if server.Port == 0 {
		server.Port = 22
	}

	if len(c.Servers) == 0 {
		server.Current = true
	}
	c.Servers = append(c.Servers, server)
}

// GetServers returns all the servers
func GetServers() []*Server {
	readServers()
	return c.Servers
}

// GetServer returns a server by name
func GetServer(name string) *Server {
	readServers()
	for i := range c.Servers {
		if c.Servers[i].Name == name {
			return c.Servers[i]
		}
	}
	return nil
}

// GetCurrentServer returns the current server
func GetCurrentServer() *Server {
	readServers()
	for i := range c.Servers {
		if c.Servers[i].Current {
			return c.Servers[i]
		}
	}
	return nil
}

// RemoveServer removes a server by name
func RemoveServer(server *Server) {
	defer func() {
		SaveConfigurationFile()
		output.Done("Server %s(%s) removed.", server.Name, server.Host)
	}()

	readServers()

	for i := range c.Servers {
		if c.Servers[i].Name == server.Name {
			c.Servers = append(c.Servers[:i], c.Servers[i+1:]...)
			return
		}
	}
}

// UpdateServer update a server by name
func UpdateServer(name string, newServer *Server) {
	defer func() {
		SaveConfigurationFile()
		output.Done("Server %s(%s) updated.", newServer.Name, newServer.Host)
	}()

	readServers()

	for i := range c.Servers {
		if c.Servers[i].Name == name {
			c.Servers[i] = newServer
			return
		}
	}
}

// TryCreateConfigurationFile creates server file
func TryCreateConfigurationFile() {
	configFile := ConfigFile()
	if _, err := os.Stat(configFile); err != nil && os.IsNotExist(err) {
		f, err := os.Create(ConfigFile())
		if err != nil {
			output.Fatal("Unable to create configuration file, %v", err)
		}
		defer f.Close()
	}
}

// SaveConfigurationFile save the servers to the configuration file
func SaveConfigurationFile() {
	readServers()

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	// encoder.SetIndent(0)
	if err := encoder.Encode(c); err != nil {
		output.Fatal("Unable to marshal configuration file, %v", err)
	}

	configContent := buf.Bytes()
	if err := os.WriteFile(ConfigFile(), configContent, 0644); err != nil {
		output.Fatal("Unable to write to configuration file, %v", err)
	}
}
