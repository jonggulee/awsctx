# awsctx

A CLI tool for switching between AWS profiles.

## Features

- List all AWS profiles from `~/.aws/config` with account IDs
- Show current active profile
- Switch profiles interactively with keyboard navigation

## Demo

```
AWS Profiles:

>   dev                            111122223333
  * staging                        444455556666
    prod                           777788889999

↑↓: move  enter: select  q: quit
```

## Installation

```bash
git clone https://github.com/jonggulee/awsctx.git
cd awsctx
go build -o awsctx .
sudo mv awsctx /usr/local/bin/awsctx
```

Add to your `~/.zshrc` (or `~/.bashrc`):

```zsh
export AWS_PROFILE=$(cat ~/.awsctx 2>/dev/null)

function awsctx() {
    /usr/local/bin/awsctx
    export AWS_PROFILE=$(cat ~/.awsctx 2>/dev/null)
}
```

Then reload your shell:

```bash
source ~/.zshrc
```

## Usage

```bash
awsctx
```

Use `↑↓` or `j/k` to navigate, `Enter` to select, `q` to quit.

## Requirements

- Go 1.21+
- AWS credentials configured in `~/.aws/config`
