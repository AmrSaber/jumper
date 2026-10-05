// jumper-inject gives OpenCode a fresh bookmark list without adding it to
// session history. The context hook also explains how agents use the list.
//
// Opt-out: JUMPER_INJECT=false|0|no suppresses both injections.

import { Plugin } from '@opencode/plugin';

const INJECT_ENABLED = !['false', '0', 'no'].includes((process.env.JUMPER_INJECT ?? '').toLowerCase());
const BLOCK_OPEN = '<jumper-bookmarks>';
const BLOCK_CLOSE = '</jumper-bookmarks>';
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

export default Plugin.define({
  id: 'jumper-inject',
  async setup(ctx) {
    await ctx.session.hook('context', async (event) => {
      if (!INJECT_ENABLED) return;

      event.system.push({ type: 'text', text: SYSTEM_PRIMER });

      const bookmarks = await renderBookmarks();
      if (bookmarks) event.system.push({ type: 'text', text: bookmarks });
    });
  },
});
