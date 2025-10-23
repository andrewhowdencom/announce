# ruf

A (vibe coded) application to make calls.

## What it does

This application is a CLI tool to send calls to different platforms. Currently, it supports Slack.

## Usage

To see a list of all available commands and flags, run:

```bash
ruf --help
```

## Getting it

You can download the latest version of the application from the [GitHub Releases page](https://github.com/andrewhowdencom/ruf/releases).

## Development

This application has been almost entirely "vibe coded" with Google Jules & Gemini.

## Task Runner

This project uses [Taskfile](https://taskfile.dev/) as a task runner for common development tasks. To use it, you'll first need to install it.

### Installation

You can install Taskfile with the following command:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

### Usage

To see a list of all available tasks, run:

```bash
task --list
```

You can then run any task with `task <task-name>`.

## Running as a Service

This project includes an example `systemd` unit file that can be used to run the application as a user-level service.

To install it, copy the file to `~/.config/systemd/user/`:

```bash
mkdir -p ~/.config/systemd/user
cp examples/ruf.service ~/.config/systemd/user/
```