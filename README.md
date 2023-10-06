# htb-cli
![Coverage](https://img.shields.io/badge/Coverage-13.6%25-red)

![Workflows (main)](https://github.com/GoToolSharing/htb-cli/actions/workflows/go.yml/badge.svg?branch=main)
![Workflows (dev)](https://github.com/GoToolSharing/htb-cli/actions/workflows/go.yml/badge.svg?branch=dev)
![GitHub go.mod Go version (main)](https://img.shields.io/github/go-mod/go-version/GoToolSharing/htb-cli/main)
![GitHub go.mod Go version (dev)](https://img.shields.io/github/go-mod/go-version/GoToolSharing/htb-cli/dev)
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

![Helper](/assets/helper.png)

## Start

![Start machine](/assets/start.png)

## Stop

![Stop machine](/assets/stop.png)

## Reset

![Reset machine](/assets/reset.png)

## Submit

### Submit machine flag

![Submit machine flag](/assets/submit_machine.png)

### Submit challenge flag

![Submit challenge flag](/assets/submit_challenge.png)

## Info

![Info active machine](/assets/info_active.png)

![Info machines](/assets/info_machines.png)

## Status

![Status](/assets/status.png)