[![Go Report Card](https://goreportcard.com/badge/github.com/poneding/ssher)](https://goreportcard.com/report/github.com/poneding/ssher)
[![GitHub release](https://img.shields.io/github/v/release/poneding/ssher)](https://img.shields.io/github/v/release/poneding/ssher)
[![GitHub license](https://img.shields.io/github/license/poneding/ssher)](https://img.shields.io/github/license/poneding/ssher)
[![GitHub stars](https://img.shields.io/github/stars/poneding/ssher)](https://img.shields.io/github/stars/poneding/ssher)
[![GitHub forks](https://img.shields.io/github/forks/poneding/ssher)](https://img.shields.io/github/forks/poneding/ssher)

ENG | [中文](README.md)

# ssher

ssher is an easy-to-use command line tool for connecting to remote servers.

ssher developed with Go language, so it supports multiple platforms, including Linux, MacOS, and Windows.

## 🔍 Preview

![ssher](https://github.com/poneding/images/blob/master/2024/04/202404260925762.gif?raw=true)

## ⚙️ Installation

### Go install

If you have already installed the Go environment locally, you can directly use `go install` to install:

```bash
go install github.com/poneding/ssher@latest
```

### Wget binary

MacOS & Linux installation, refer to the following commands:

```bash
# MacOS
sudo wget https://github.com/poneding/ssher/releases/download/v1.1.0/ssher_1.1.0_darwin_arm64 -O /user/local/bin/ssher && sudo chmod +x /user/local/bin/ssher

# Linux
sudo wget https://github.com/poneding/ssher/releases/download/v1.1.0/ssher_1.1.0_linux_amd64 -O /user/local/bin/ssher && sudo chmod +x /user/local/bin/ssher
```

> Note: Before downloading, make sure your system is `arm64` or `amd64`, and download the corresponding binary file.
\
Windows installation, refer to the following steps:

Download the `ssher.exe` file first:

```bash
# Download .exe file
wget https://github.com/poneding/ssher/releases/download/v1.1.0/ssher_1.1.0_windows_amd64.exe
```

Add the `ssher.exe` file path to the environment variable after download done, or put it in a path that has already been added to the environment variable.

### Download from browser

[👉🏻 GitHub Releases](https://github.com/poneding/ssher/releases)

### Compile from source

[Go environment](https://go.dev/doc/install) is required, then execute the following command:

```bash
git clone https://github.com/poneding/ssher.git
cd ssher
go build -o ssher main.go
```

## 🛠️ Usage

### Get started

```bash
ssher
ssher -s <server>
```

Execute the above command, you will enter the interactive mode, where you can use the `↓ ↑ → ←` keys to select the server you want to connect to or add a new server.

### Server management operation

```bash
# List all servers(list, ls)
ssher list

# Add a new server
ssher add

# Remove a server(remove, rm)
ssher remove
ssher remove <server_a> <server_b>

# Edit a server
ssher edit
ssher edit <server>
```

### Version and upgrade

```bash
# Check version
ssher version

# Upgrade
ssher upgrade
```

### Auto-completion

```bash
# Append the completion script to ~/.bashrc or ~/.zshrc
# bash
echo 'source <(ssher completion bash)' >> ~/.bashrc
source ~/.bashrc

# zsh
echo 'source <(ssher completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

## ⭐️ Stars

[![Stargazers over time](https://starchart.cc/poneding/ssher.svg?variant=adaptive)](https://starchart.cc/poneding/ssher)

Welcome to give me a Star ⭐️ if this project is helpful to you, your support is my greatest motivation.
