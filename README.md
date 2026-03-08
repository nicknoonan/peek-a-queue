# Peek-a-Queue

Peek-a-Queue is a tiny, cheeky TUI tool that sneaks a peek at your AWS SQS queues. It lists queues, shows per-queue attributes, and lets you poke around with a keyboard. Kinda like window-shopping for messages, without the commitment to AWS console.

Think of it as a queue inspector that wears a CLI hat.

## Features
- Lists all SQS queues in your account
- Loads and displays queue attributes on demand
- Keyboard-driven UI with filtering, pagination, and a help menu
- Periodic background refreshes of visible queue attributes

## Requirements
- Go 1.20+ (or whatever your project uses)
- AWS credentials available via environment or `~/.aws/credentials` (standard AWS SDK config)

## Install

From source (recommended during development):

```bash
cd /path/to/peek-a-queue
go build -o peek-a-queue .
# or run directly
go run .
```

Install globally with `go install` (replace with your module path):

```bash
go install github.com/nicknoonan/peek-a-queue@latest
# then run
peek-a-queue
```

## Quick Start

1. Make sure your AWS credentials are set:
```bash
export AWS_PROFILE=default   # on Windows PowerShell: $env:AWS_PROFILE="default"
export AWS_REGION=us-east-1  # or set in your profile
```

2. Run the TUI:
```bash
peek-a-queue
```

3. Keyboard cheatsheet (because memorizing is optional, but fun):
- s — toggle spinner
- T — toggle title bar (also toggles filter)
- S — toggle status bar
- P — toggle pagination
- H — toggle help menu
- Enter — refresh selected queue's attributes
- r — refresh attributes for the visible page


## Screenshots

Drop your screenshots into the repo and replace these placeholders:

![Peek-a-Queue main screen - placeholder](assets/pickaqueue.gif)
*(Replace with a screenshot of the main UI)*


## Contributing
Bugs, features, and parade floats welcome. Open an issue or send a PR. Keep it small, keep it kind, and maybe include a gif.

## Happy queue peeking! 👀🚀