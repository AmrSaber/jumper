// Package assets embeds the agent-integration files
package assets

import _ "embed"

// OpencodePlugin is the jumper-inject opencode plugin, written to
// ~/.config/opencode/plugins/jumper-inject.ts by `jumper agent setup opencode`.
//
//go:embed opencode/plugins/jumper-inject.ts
var OpencodePlugin string
