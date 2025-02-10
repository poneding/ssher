/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package ssh

import (
	"encoding/base64"

	"github.com/poneding/ssher/internal/output"
)

type Server struct {
	Name         string `yaml:"name"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password,omitempty"`
	IdentityFile string `yaml:"identity_file,omitempty"`
	Current      bool   `yaml:"current,omitempty"`
	Emoji        string `yaml:"-"`
}

type Options func(*Server)

func NewServer(name, host string, options ...Options) *Server {
	s := &Server{
		Name: name,
		Host: host,
	}

	for _, option := range options {
		option(s)
	}

	if s.Port == 0 {
		s.Port = 22
	}

	if s.Host == "" {
		output.Fatal("Server host is required")
	}

	return s
}

func WithPort(port int) Options {
	return func(s *Server) {
		s.Port = port
	}
}

func WithUser(user string) Options {
	return func(s *Server) {
		s.User = user
	}
}

func WithPassword(password string) Options {
	return func(s *Server) {
		s.Password = base64.StdEncoding.EncodeToString([]byte(password))
	}
}

func WithIdentityFile(identityFile string) Options {
	return func(s *Server) {
		s.IdentityFile = identityFile
	}
}
