// jumper-inject gives OpenCode a fresh bookmark list without adding it to
// session history. The context hook also explains how agents use the list.
//
// Opt-out: JUMPER_INJECT=false|0|no suppresses both injections.

const INJECT_ENABLED = !['false', '0', 'no'].includes((process.env.JUMPER_INJECT ?? '').toLowerCase());
const BLOCK_OPEN = '<jumper-bookmarks>';
const BLOCK_CLOSE = '</jumper-bookmarks>';
const SKILL_CONTENT = `---
name: jumper
description: Use when managing jumper bookmarks — marking, deleting, renaming, or pruning directory bookmarks.
---

## What this skill is for

jumper is directory-bookmark CLI — short names mapped to directories. Reading bookmarks is handled without this skill (see below), so this skill exists only for the **management (write) commands**, which are destructive enough to warrant their own guidance and should never run unprompted.

## Getter commands — you rarely need these

\`jumper list\` and \`jumper resolve <name>[/subpath]\` read bookmarks. You normally do **not** need to call them to get a path — the current list (name → path) is already injected every turn in the \`<jumper-bookmarks>\` block, so read from there. Reach for \`jumper resolve\` only when you need a path that isn't directly in the block, for example, a subpath (\`jumper resolve <name>/<subpath>\`) for scripting.

## Management commands — only on explicit user request

Never invoke these proactively. Use only when the user explicitly asks to manage their bookmarks.

| Command                      | Aliases               | What it does                                                    |
| ---------------------------- | --------------------- | --------------------------------------------------------------- |
| \`jumper mark [name] [dir]\`   |                       | Bookmark a directory; defaults to current dir and its base name |
| \`jumper delete [name\|path]\` | \`del\`, \`rm\`, \`unmark\` | Delete by name or path; no arg = current directory              |
| \`jumper rename <old> <new>\`  | \`mv\`                  | Rename a bookmark (case-insensitive match on old name)          |
| \`jumper prune\`               | \`clean\`               | Remove bookmarks whose paths no longer exist                    |

## Gotchas

- Bookmark names cannot start with \`.\`, \`~\`, or \`/\`; \`-\` is reserved.
- \`mark\` silently overwrites an existing bookmark of the same name.
- \`delete\` with a path argument (starts with \`/\`, \`~\`, or \`.\`) removes **all** bookmarks pointing to that path — use name form to be precise.
- The injected \`<jumper-bookmarks>\` block (and \`jumper list\`) might show paths that have since been deleted or moved; \`jumper prune\` cleans these up.
`;
const SYSTEM_PRIMER = [
  "jumper is the user's directory-bookmark CLI: short names mapped to directories.",
  `The current bookmarks are injected each turn in a \`${BLOCK_OPEN}\` block (name: path).`,
  'When the user refers to a project or directory by name, use that path directly instead of asking.',
].join('\n');

async function renderBookmarks(): Promise<string> {
  try {
    const result = await Bun.$`jumper list -o json`.quiet().nothrow();
    if (result.exitCode !== 0) return '';

    const bookmarks = JSON.parse(result.stdout.toString());
    if (!Array.isArray(bookmarks) || bookmarks.length === 0) return '';

    const lines = bookmarks
      .filter(
        (bookmark: { title?: unknown; path?: unknown }) =>
          typeof bookmark.title === 'string' && typeof bookmark.path === 'string',
      )
      .map((bookmark: { title: string; path: string }) => `- ${bookmark.title}: ${bookmark.path}`);
    return lines.length ? [BLOCK_OPEN, ...lines, BLOCK_CLOSE].join('\n') : '';
  } catch {
    return '';
  }
}

export default {
  id: 'jumper-inject',
  async setup(ctx) {
    await ctx.skill.transform((editor) => {
      editor.add({
        id: 'jumper',
        name: 'jumper',
        description: 'Use when managing jumper bookmarks — marking, deleting, renaming, or pruning directory bookmarks.',
        path: new URL(import.meta.url).pathname,
        content: SKILL_CONTENT,
      });
    });

    await ctx.session.hook('context', async (event) => {
      if (!INJECT_ENABLED) return;

      event.system.push({ type: 'text', text: SYSTEM_PRIMER });

      const bookmarks = await renderBookmarks();
      if (bookmarks) event.system.push({ type: 'text', text: bookmarks });
    });
  },
};
