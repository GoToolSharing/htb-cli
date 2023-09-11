```
This software, engineered using the Go programming language, serves to streamline and automate various tasks for the HackTheBox platform, enhancing user efficiency and productivity.

Usage:
  htb-cli [command]

Available Commands:
  active      Catalogue of active machines
  submit      Submit credentials (User and Root Flags)
  help        Help about any command
  info        Showcase detailed machine information
  reset       Reset a machine - [WIP]
  start       Start a machine
  stop        Stop the current machine

Flags:
  -h, --help           help for htb-cli
  -p, --proxy string   Configure a URL for an HTTP proxy
  -v, --verbose        Verbose mode

Use "htb-cli [command] --help" for more information about a command.
```

## Installation

`go install github.com/GoToolSharing/htb-cli@latest`

## Configuration

You must add a Hackthebox **App token** in the **HTB_TOKEN** environment variable (zshrc maybe).
API Token can be find here : https://app.hackthebox.com/profile/settings => Create App Token

```
export HTB_TOKEN=eyJ...
```

## Start

```
> htb-cli start -m Flight
Machine deployed to lab.
```

## Stop

```
> htb-cli stop
Machine terminated.
```

## Reset

```
> htb-cli reset -m Flight
Machine terminated.
```

## Submit

This command allows to submit the user flag and the root flag of active and retired machines. The first argument is the flag and the second the difficulty /10.

```
> htb-cli submit -f flag4testing -d 3

SteamCloud user is now owned.
```

## Info

By default the command shows the active machine.

```
> htb-cli info -m Zipper -m Sau

Name     |OS      |Active   |Difficulty   |Stars   |FirstUserBlood   |FirstRootBlood   |Status            |Release
Zipper   |Linux   |0        |Hard         |4.5     |1H 57M 59S       |2H 8M 18S        |❌ User - ❌ Root   |2018-10-20
Sau      |Linux   |1        |Easy         |4.6     |0H 8M 39S        |0H 11M 40S       |✅ User - ✅ Root   |2023-07-08
```