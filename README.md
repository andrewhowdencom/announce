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
| `source.urls` | A list of URLs to fetch calls from. Remote (`https://...`), local (`file://...`) and git (`git://...`) URLs are supported. See the Git Sources section for more information. |
| `slack.app_token` | The Slack app token to use for sending calls. |
| `git.tokens` | A map of git providers to personal access tokens. Currently, only `github.com` is supported. |

### Example

```yaml
source:
  urls:
    - "https://example.com/announcements.yaml"
    - "file:///path/to/local/announcements.yaml"
    - "git://github.com/andrewhowdencom/ruf-example-announcements/tree/main/example.yaml"

slack:
  app:
    token: ""

git:
  tokens:
    github.com: "YOUR_GITHUB_TOKEN"
```

### Git Sources

The application supports fetching calls from Git repositories. The URL format is:

`git://<repository>/tree/<refspec>/<file-path>`

For example:

`git://github.com/andrewhowdencom/ruf-example-announcements/tree/main/example.yaml`

### Slack Configuration

To use the Slack integration, you'll need to create a Slack app and install it in your workspace. The app will need the following permissions:

- `channels:read`: To list public channels.
- `groups:read`: To list private channels.
- `chat:write`: To send messages.

## Call Format

The application expects the source YAML files to contain a top-level `calls` list. Optionally, a `campaign` can be specified. If a campaign is not specified, it will be derived from the filename.

```yaml
campaign:
  id: "my-campaign"
  name: "My Campaign"
calls:
- id: "unique-id-1"
  author: "author@example.com"
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

## Listing Sent Calls

When you list the sent calls, you will see the following statuses:

| Status | Description |
| --- | --- |
| `sent` | The call has been successfully sent. |
| `deleted` | The call has been sent and then subsequently deleted. |

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
