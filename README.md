# Jumper - Directory Bookmarks for Your Shell

[![Latest Release](https://img.shields.io/github/v/release/AmrSaber/jumper?logo=github)](https://github.com/AmrSaber/jumper/releases/latest)
[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)](https://go.dev/)

Jumper lets you bookmark directories and jump to them by name from anywhere in your shell. No more `cd ../../some/deeply/nested/project` — just `jump proj`.

**Inspired by [zshmarks](https://github.com/jocelynmallon/zshmarks)**, Jumper is a standalone Go binary with shell integration, case-insensitive subdirectory tab-completion, and Table/JSON/YAML output.

## Table of Contents

- [Installation](#installation)
- [Shell Setup](#shell-setup)
- [Quick Start](#quick-start)
- [Tab Completion](#tab-completion)
- [Data Storage](#data-storage)

## Installation

**Homebrew (macOS and Linux):**

```bash
brew install AmrSaber/tap/jumper
```

**Go:**

```bash
go install github.com/AmrSaber/jumper@latest
```

## Shell Setup

Jumper provides a shell integration script that defines the `jump` function (which handles the `cd` for you). Add one line to your shell's rc file:

```bash
eval "$(jumper init)"
```

Jumper auto-detects your shell from `$SHELL`. You can also pass it explicitly:

```bash
eval "$(jumper init zsh)"   # ~/.zshrc
eval "$(jumper init bash)"  # ~/.bashrc
```

Then reload your shell:

```bash
source ~/.zshrc   # or ~/.bashrc
```

> Only `jump` is a shell function — all other commands (`jumper mark`, `jumper delete`, etc.) are invoked directly.

## Quick Start

```bash
# Bookmark the current directory using its base name
cd ~/Projects/my-app
jumper mark           # saves as "my-app"

# Bookmark the current directory with an explicit name
jumper mark my-app

# Bookmark a specific directory without cd-ing into it
jumper mark my-app ~/Projects/my-app

# Jump to a bookmark from anywhere
jump my-app

# Jump into a subdirectory directly
jump my-app/src/components

# List all bookmarks
jumper list
# ┌──────────┬───────────────────┐
# │ TITLE    │ PATH              │
# ├──────────┼───────────────────┤
# │ dotfiles │ ~/.dotfiles       │
# │ my-app   │ ~/Projects/my-app │
# └──────────┴───────────────────┘

# Rename a bookmark
jumper rename my-app app

# Delete a bookmark by name
jumper delete dotfiles

# Delete multiple bookmarks by name
jumper delete dotfiles my-app

# Delete all bookmarks pointing to a specific directory
jumper delete ~/Projects/my-app

# Delete all bookmarks pointing to the current directory
jumper delete           # or: jumper delete .

# Use `jumper help` to get the full help
```

## Tab Completion

The `jump` function supports tab completion for bookmark names and subdirectories.

- **Bookmark names** are matched case-insensitively — typing `proj` matches `Projects`
- **A trailing `/` is appended** to each completion so you can keep chaining without extra keystrokes
- **Subdirectory expansion** works at any depth, also case-insensitively

```
jump pro<Tab>         → jump Projects/
jump Projects/<Tab>   → jump Projects/src/
                        jump Projects/docs/
jump Projects/sr<Tab> → jump Projects/src/
```

## Data Storage

Bookmarks are stored as a JSON file at your system's standard data location:

- **Linux**: `$XDG_DATA_HOME/jumper/bookmarks.json` (defaults to `~/.local/share/jumper/bookmarks.json`)
- **macOS**: `~/Library/Application Support/jumper/bookmarks.json`

The file is a plain JSON array — easy to inspect, back up, or sync with your dotfiles.

```json
[
  { "title": "dotfiles", "path": "/home/amr/.dotfiles" },
  { "title": "proj",     "path": "/home/amr/Projects/my-app" }
]
```

---

## Contributing

1. Fork the repository at [github.com/AmrSaber/jumper](https://github.com/AmrSaber/jumper)
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes
4. Open a Pull Request

Bug reports and feature requests: **https://github.com/AmrSaber/jumper/issues**
