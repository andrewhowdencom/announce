# ruf

A (vibe coded) application to make calls.

## What it does

This application is a CLI tool to send calls to different platforms. Currently, it supports Slack.

## Usage

To see a list of all available commands and flags, run:

```bash
ruf --help
```

## Configuration

The application is configured using a YAML file located at `$XDG_CONFIG_HOME/ruf/config.yaml`. The following configuration options are available:

| Name | Description |
| --- | --- |
| `source.urls` | A list of URLs to fetch calls from. Both remote (`https://...`) and local (`file://...`) URLs are supported. |
| `slack.app_token` | The Slack app token to use for sending calls. |

### Example

```yaml
source:
  urls:
    - "https://example.com/announcements.yaml"
    - "file:///path/to/local/announcements.yaml"

slack:
  app:
    token: ""
```

## Call Format

The application expects the source YAML files to contain a list of calls with the following format:

```yaml
- id: "unique-id-1"
  subject: "Hello!"
  content: "Hello, world!"
  destinations:
    - type: "slack"
      to:
        - "C1234567890"
  scheduled_at: "2025-01-01T12:00:00Z"
- id: "unique-id-2"
  subject: "Recurring hello!"
  content: "Hello, recurring world!"
  destinations:
    - type: "slack"
      to:
        - "C1234567890"
  cron: "0 * * * *"
  recurring: true
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
