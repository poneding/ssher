[![Go Report Card](https://goreportcard.com/badge/github.com/poneding/ssher)](https://goreportcard.com/report/github.com/poneding/ssher)
[![GitHub release](https://img.shields.io/github/v/release/poneding/ssher)](https://img.shields.io/github/v/release/poneding/ssher)
[![GitHub license](https://img.shields.io/github/license/poneding/ssher)](https://img.shields.io/github/license/poneding/ssher)
[![GitHub stars](https://img.shields.io/github/stars/poneding/ssher)](https://img.shields.io/github/stars/poneding/ssher)
[![GitHub forks](https://img.shields.io/github/forks/poneding/ssher)](https://img.shields.io/github/forks/poneding/ssher)

# ssher

中文 | [ENG](README_en-US.md)

ssher 是一款轻量的 SSH Profile 管理命令行工具，让你可以更方便的连接到你的服务器。

由于是使用 Go 语言开发，所以支持多平台，包括 Linux、MacOS 和 Windows。

## 🔍 预览

![ssher](https://images.poneding.com/2024/04/202404260925762.gif)

## ⚙️ 安装

### 直接 Go install 安装

如果本地已经安装了 Go 环境，可以直接使用 `go install` 安装：

```bash
go install github.com/poneding/ssher@latest
```

### 二进制文件

MacOS & Linux 安装，参考以下命令：

```bash
# MacOS
sudo wget https://ghproxy.ketches.cn/https://github.com/poneding/ssher/releases/download/v0.1.0/ssher_0.1.0_darwin_arm64 -O /user/local/bin/ssher && sudo chmod +x /user/local/bin/ssher

# Linux
sudo wget https://ghproxy.ketches.cn/https://github.com/poneding/ssher/releases/download/v0.1.0/ssher_0.1.0_linux_amd64 -O /user/local/bin/ssher && sudo chmod +x /user/local/bin/ssher
```

> 注意：下载前确认你的系统是 `arm64` 还是 `amd64`，下载对应的二进制文件。

Windows 安装，参考以下步骤：

首先下载 `ssher.exe` 文件：

```bash
# 下载 .exe 文件
wget https://ghproxy.ketches.cn/https://github.com/poneding/ssher/releases/download/v0.1.0/ssher_0.1.0_windows_amd64.exe
```

下载完成后，将 `ssher.exe` 文件路径添加到环境变量中，或者将其放到一个已经添加到环境变量的路径下。

### 浏览器下载

[👉🏻 发布下载](https://github.com/poneding/ssher/releases)，国内网络访问可能受阻。

### 源码编译

需要[安装 Go 环境](https://golang.google.cn/doc/install)，然后执行以下命令：

```bash
git clone https://github.com/poneding/ssher.git
cd ssher
go build -o ssher main.go
```

## 🛠️ 使用

### 快速开始

```bash
ssher
```

输入以上命令，会进入交互模式，在交互模式中你可以使用通过 `↓ ↑ → ←` 键选择你要连接的 SSH Profile 或者添加新的 SSH Profile。

### SSH 连接管理操作

```bash
# 查看 ssh 连接列表
ssher list

# 添加 ssh 连接
ssher add

# 删除 ssh 连接
ssher remove
```

### SSH Profile 文件操作

Profile 文件中存储了你的服务器信息，你可以通过以下命令查看和编辑 profile 文件：

```bash
# 查看 profile
ssher profile view

# 编辑 profile
ssher profile edit
```

### 版本和升级

```bash
# 查看版本
ssher version

# 升级
ssher upgrade
```

### 命令自动补全

```bash
# 将补全脚本写入到 ~/.bashrc 或者 ~/.zshrc 中
# bash
echo 'source <(ssher completion bash)' >> ~/.bashrc
source ~/.bashrc

# zsh
echo 'source <(ssher completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

## ⭐️ Stars

[![Stargazers over time](https://starchart.cc/poneding/ssher.svg?variant=adaptive)](https://starchart.cc/poneding/ssher)

如果您觉得这个项目不错，欢迎给我一个 Star ⭐️，你的支持是我最大的动力。
