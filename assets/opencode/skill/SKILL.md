---
name: jumper
description: Use when managing jumper bookmarks — marking, deleting, renaming, or pruning directory bookmarks.
---

## What this skill is for

jumper is directory-bookmark CLI — short names mapped to directories. Reading bookmarks is handled without this skill (see below), so this skill exists only for the **management (write) commands**, which are destructive enough to warrant their own guidance and should never run unprompted.

## Getter commands — you rarely need these

`jumper list` and `jumper resolve <name>[/subpath]` read bookmarks. You normally do **not** need to call them to get a path — the current list (name → path) is already injected every turn in the `<jumper-bookmarks>` block, so read from there. Reach for `jumper resolve` only when you need a path that isn't directly in the block, e.g. a subpath (`jumper resolve <name>/<subpath>`) for scripting.

## Management commands — only on explicit user request

Never invoke these proactively. Use only when the user explicitly asks to manage their bookmarks.

| Command                      | Aliases               | What it does                                                    |
| ---------------------------- | --------------------- | --------------------------------------------------------------- |
| `jumper mark [name] [dir]`   |                       | Bookmark a directory; defaults to current dir and its base name |
| `jumper delete [name\|path]` | `del`, `rm`, `unmark` | Delete by name or path; no arg = current directory              |
| `jumper rename <old> <new>`  | `mv`                  | Rename a bookmark (case-insensitive match on old name)          |
| `jumper prune`               | `clean`               | Remove bookmarks whose paths no longer exist                    |

## Gotchas

- Bookmark names cannot start with `.`, `~`, or `/`; `-` is reserved.
- `mark` silently overwrites an existing bookmark of the same name.
- `delete` with a path argument (starts with `/`, `~`, or `.`) removes **all** bookmarks pointing to that path — use name form to be precise.
- The injected `<jumper-bookmarks>` block (and `jumper list`) may show paths that have since been deleted or moved; `jumper prune` cleans these up.
