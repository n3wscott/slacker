# slacker

`slacker` is a tool for interacting to slack via a one-off api call. Intended to be used with automation.

[![GoDoc](https://godoc.org/github.com/n3wscott/prbuilder?status.svg)](https://godoc.org/github.com/n3wscott/slacker)
[![Go Report Card](https://goreportcard.com/badge/n3wscott/prbuilder)](https://goreportcard.com/report/n3wscott/slacker)

_Work in progress._

## Installation

`slacker` can be installed via:

```shell
go get github.com/n3wscott/slacker/cmd/slacker
```

To update your installation:

```shell
go get -u github.com/n3wscott/slacker/cmd/slacker
```

### Auth

A slackbot oauth token needs to exist in a file called `token` at `/var/bindings/slackbot`.

This token can be produced on the slack api oauth page for bots.

## Usage

`slacker` has two commands, `get` and `send`

```shell
Interact with slack from a command line.

Usage:
  slacker [command]

Available Commands:
  get         Get a list of channels or direct messages.
  help        Help about any command
  send        Send a message or response to a channel or thread.

Flags:
  -h, --help   help for slacker
```

### Get

Get a list of channels or direct messages.

```shell
Usage:
  slacker get [flags]
  slacker get [command]

Available Commands:
  channel     Get a list of channels.
  dm          Get a list of direct message channels.

Flags:
  -h, --help   help for get
```

### Send

Send a message or response to a channel or thread.

```shell
Usage:
  slacker send [flags]

Flags:
  -h, --help              help for send
      --id string         Channel ID to use.
      --json              Output as JSON.
      --name string       Channel Name to use, will resolve to a Channel ID.
      --reaction          Message is treated as a reaction to be added.
      --remove-reaction   Message is treated as a reaction to be removed.
      --thread string     Unique identifier of a thread's parent message.
```
