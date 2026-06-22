// Package assets embeds the agent-integration files
package assets

import _ "embed"

// OpencodeSkill is the jumper skill for opencode, written to
// ~/.config/opencode/skills/jumper/SKILL.md by `jumper agent setup opencode`.
//
//go:embed opencode/skill/SKILL.md
var OpencodeSkill string

// OpencodePlugin is the jumper-inject opencode plugin, written to
// ~/.config/opencode/plugins/jumper-inject.ts by `jumper agent setup opencode`.
//
//go:embed opencode/plugins/jumper-inject.ts
var OpencodePlugin string
