# htb-cli
![Coverage](https://img.shields.io/badge/Coverage-55.2%25-yellow)

![Workflows (main)](https://github.com/GoToolSharing/htb-cli/actions/workflows/go.yml/badge.svg?branch=main)
![Workflows (dev)](https://github.com/GoToolSharing/htb-cli/actions/workflows/go.yml/badge.svg?branch=dev)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/GoToolSharing/htb-cli)
![GitHub release](https://img.shields.io/github/v/release/GoToolSharing/htb-cli)
![GitHub Repo stars](https://img.shields.io/github/stars/GoToolSharing/htb-cli)

<div>
  <img alt="current version" src="https://img.shields.io/badge/linux-supported-success">
  <img alt="current version" src="https://img.shields.io/badge/windows-supported-success">
  <img alt="current version" src="https://img.shields.io/badge/mac-supported-success">
  <br>
  <img alt="amd64" src="https://img.shields.io/badge/amd64%20(x86__64)-supported-success">
  <img alt="arm64" src="https://img.shields.io/badge/arm64%20(aarch64)-supported-success">
</div>

## Installation

`go install github.com/GoToolSharing/htb-cli@latest`

## Configuration

You must add a Hackthebox **App token** in the **HTB_TOKEN** environment variable (zshrc maybe).
API Token can be find here : https://app.hackthebox.com/profile/settings => `Create App Token`

```
export HTB_TOKEN=eyJ...
```

## Helper

```
This software, engineered using the Go programming language, serves to streamline and automate various tasks for the HackTheBox platform, enhancing user efficiency and productivity.

Usage:
  htb-cli [command]

Available Commands:
  active      Catalogue of active machines
  help        Help about any command
  info        Showcase detailed machine information
  reset       Reset a machine
  start       Start a machine
  status      Displays the status of HackTheBox servers
  stop        Stop the current machine
  submit      Submit credentials (User and Root Flags)

Flags:
  -h, --help           help for htb-cli
  -p, --proxy string   Configure a URL for an HTTP proxy
  -v, --verbose        Verbose mode

Use "htb-cli [command] --help" for more information about a command.
```

## Start

```
❯ htb-cli start -m Blue
? The following machine was found : Blue Yes
Machine deployed to lab.
```

## Stop

```
❯ htb-cli stop
Machine terminated.
```

## Reset

```
❯ htb-cli reset
CozyHosting will be reset in 1 minute.
```

## Submit

This command allows to submit the user flag and the root flag of active and retired machines. The first argument is the flag and the second the difficulty /10.

### Submit machine flag
```
❯ htb-cli submit -m SteamCloud -f flag4testing -d 3
? The following machine was found : SteamCloud Yes
SteamCloud user is now owned.
```

### Submit challenge flag
```
❯ htb-cli submit -c Phonebook -f flag4testing -d 3
? The following challenge was found : Phonebook Yes
Incorrect flag
```

## Info

```
❯ htb-cli info
? Do you want to check for active machine ? Yes
Name   |OS        |Active   |Difficulty   |Stars   |IP            |Status            |Release
Blue   |Windows   |0        |Easy         |4.5     |10.10.10.40   |✅ User - ✅ Root   |2017-07-28
```

```
❯ htb-cli info -m Zip -m pilgrimage
? Do you want to check for active machine ? No
? The following machine was found : Zipping Yes
Name      |OS      |Active   |Difficulty   |Stars   |FirstUserBlood   |FirstRootBlood   |Status            |Release
Zipping   |Linux   |✅        |Medium       |4.1     |0H 15M 9S        |1H 12M 28S       |❌ User - ❌ Root   |2023-08-26
? The following machine was found : Pilgrimage Yes
Pilgrimage   |Linux   |✅    |Easy   |4.5   |0H 17M 0S   |0H 20M 33S   |✅ User - ✅ Root   |2023-06-24
```

## Status

```
❯ htb-cli status
All Systems Operational
```
