# Peek-a-Queue

Peek-a-Queue is a tiny, cheeky TUI tool that sneaks a peek at your AWS SQS queues. It lists queues, shows per-queue attributes, and lets you poke around with a keyboard. Kinda like window-shopping for messages, without the commitment to AWS console.

Think of it as a queue inspector that wears a CLI hat.

## Why use this?
Because clicking through the AWS Console to find your ApproximateNumberOfMessages is a form of digital penance nobody deserves. Peek-a-Queue gives you a bird’s-eye view of your message traffic without the heavy lifting.

## Features
- Zero-Touch Inspection: View Messages Available and Messages In-Flight without accidentally consuming or deleting them.

- Periodic background refreshes: keep your data fresh, so you can watch your queue drain (or pile up) in real-time.

- Filtering: Instantly filter through hundreds of queues to find that one specific dead-letter queue (DLQ) that's currently on fire.

- Keyboard-First: Because your mouse is for Slack, not for infrastructure.

- TUI Magic: Built with the Bubble Tea framework, making it smoother than a buttered slide.

## Requirements
- Go 1.20+ (or whatever your project uses)
- AWS credentials available via environment or `~/.aws/credentials` (standard AWS SDK config)

## Install

### From source:

```bash
cd /path/to/peek-a-queue
go build -o peek-a-queue .
# or run directly
go run .
```

### Global Install:

```bash
go install github.com/nicknoonan/peek-a-queue@latest
# then run
peek-a-queue
```

## Quick Start

1. Make sure your AWS credentials are set:
```bash
export AWS_PROFILE=default
export AWS_REGION=us-east-1

# or
aws login

# idk you've probably got a special way for loading aws creds that sdk can use. just do that.
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

Pagination and filtering example:

![Peek-a-Queue pagination and filtering](assets/pickaqueue.gif)

Refresh example:
![Peek-a-Queue refreshing example](assets/refreshexample.gif)


## Contributing
Bugs, features, and parade floats welcome. Open an issue or send a PR. Keep it small, keep it kind, and maybe include a gif.

## Happy queue peeking! 👀🚀