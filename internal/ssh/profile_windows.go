//go:build windows
// +build windows

/*
Copyright © 2024 Pone Ding <poneding@gmail.com>
*/
package ssh

import (
	"fmt"
	"os"
)

type Profile struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
	PrivateKey string `json:"privateKey"`
	Current    bool   `json:"current"`
}

type Options func(*Profile)

func NewProfile(name, host string, options ...Options) *Profile {
	p := &Profile{
		Name: name,
		Host: host,
	}

	for _, option := range options {
		option(p)
	}

	if p.Port == 0 {
		p.Port = 22
	}

	if p.Host == "" {
		fmt.Println("✗ Host is required for a profile")
		os.Exit(0)
	}

	return p
}

func WithPort(port int) Options {
	return func(p *Profile) {
		p.Port = port
	}
}

func WithUser(user string) Options {
	return func(p *Profile) {
		p.User = user
	}
}

func WithPrivateKey(privateKey string) Options {
	return func(p *Profile) {
		p.PrivateKey = privateKey
	}
}
