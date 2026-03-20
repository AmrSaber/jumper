Output static/demo.mp4
Output static/demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 900
Set Height 450
Set Theme "Dracula"
Set Padding 10

# ── Hidden setup ────────────────────────────────────────────────────────────
Hide

# Point jumper at a throw-away data dir so real bookmarks are untouched
Type "export XDG_DATA_HOME=$(mktemp -d)"
Enter
Sleep 200ms

# Create demo directory tree
Type "mkdir -p ~/demo-dirs/my-app/src ~/demo-dirs/.dotfiles ~/demo-dirs/scripts/deploy"
Enter
Sleep 200ms

# Load shell integration (defines `jump` and tab completion)
# Explicit "bash" because VHS may not set $SHELL correctly
Type `eval "$(jumper init bash)"`
Enter
Sleep 200ms

# Start from home so the first visible `cd` is intentional
Type "cd ~"
Enter
Sleep 200ms

Ctrl+L

Show
# ── Visible demo ────────────────────────────────────────────────────────────
Type "# First let's bookmark some directories"
Enter
Sleep 1s

# --- Mark: from inside the directory (implicit name from base name) ---
Type "cd ~/demo-dirs/my-app"
Enter

Type "jumper mark"
Enter
Sleep 2s

Enter

# --- Mark: explicit name ---
Type "cd ~/demo-dirs/.dotfiles"
Enter

Type "jumper mark dotfiles"
Enter
Sleep 2s

Enter

# --- Mark: alias for a nested path ---
Type "jumper mark deploy ~/demo-dirs/scripts/deploy"
Enter
Sleep 2s

Enter

# --- Mark: non-existing directory ---
Type "jumper mark ghost ~/demo-dirs/does-not-exist"
Enter
Sleep 2s

Enter

# --- List all bookmarks ---
Type "jumper list"
Enter
Sleep 4s

Ctrl+L

# --- Jump to another bookmark ---
Type "# Jump to marks"
Enter

Type "cd && pwd"
Enter
Sleep 500ms

Enter

Type "jump dotfiles"
Enter

Type "pwd"
Enter
Sleep 2s

Ctrl+L

# --- Tab completion: type partial name, Tab completes it ---
Type "# Auto-completion"
Enter
Sleep 500ms

Type "jump app"
Sleep 500ms
Tab
Sleep 500ms

Enter

Type "pwd"
Enter
Sleep 2s

Ctrl+L

# --- Subdirectory traversal ---
Type "# Nested directories auto-completion"
Enter
Sleep 500ms

Type "jump app"
Sleep 500ms
Tab
Sleep 500ms

Type "s"
Tab
Sleep 500ms

Enter

Type "pwd"
Enter
Sleep 2s

Ctrl+L

# --- Rename ---
Type "jumper rename my-app app"
Enter

Enter

# Show list to confirm the rename visually
Type "jumper list"
Enter
Sleep 2s

Enter

Type "jump app"
Enter

Type "pwd"
Enter
Sleep 5s

Ctrl+L

# --- Delete multiple by name ---
Type "jumper delete deploy ghost"
Enter
Sleep 2s

Enter

# --- Final list showing clean state ---
Type "jumper list"
Enter
Sleep 5s

# ── Hidden teardown ─────────────────────────────────────────────────────────
Hide

# Clean up
Type "rm -rf $XDG_DATA_HOME ~/demo-dirs"
Enter
Sleep 100ms
