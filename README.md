# awsctx

A lightweight CLI tool for switching between AWS profiles interactively.

## Features

- List all AWS profiles from `~/.aws/config` with name, region, and account ID
- Show current active profile (marked with `*`)
- Async account ID loading via STS with a spinner
- Switch profiles interactively with keyboard navigation
- Saves selected profile to `~/.awsctx` for shell integration

## Demo

```
AWS Profiles

>  * dev                  ap-northeast-2       111122223333
     staging              us-east-1            444455556666
     prod                 ap-northeast-2       777788889999

↑↓/jk: move  enter: select  q: quit
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

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Select profile |
| `q` / `Ctrl+C` | Quit |

## Requirements

- Go 1.21+
- AWS credentials configured in `~/.aws/config`
- AWS credentials with `sts:GetCallerIdentity` permission (for account ID lookup)

## How it works

1. Parses `~/.aws/config` to collect all `[profile ...]` sections
2. Fetches account IDs asynchronously via `sts:GetCallerIdentity` for each profile
3. On profile selection, writes the profile name to `~/.awsctx`
4. The shell function re-exports `AWS_PROFILE` after each switch
