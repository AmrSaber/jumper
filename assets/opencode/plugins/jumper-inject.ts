// jumper-inject — opencode plugin: teach the agent about jumper once (system
// prompt) and feed it the live bookmark list every turn (messages), so it
// always has the current name→path map without running `jumper list` itself.
//
// Two hooks, deliberately split:
//
//   experimental.chat.system.transform — push a tiny STATIC explainer into the
//     system prompt: what jumper is, how to use it at a high level, and that the
//     live directory list is injected each turn. This is what makes the injected
//     list intelligible even on turns where the agent hasn't loaded the jumper
//     skill (skills are only surfaced as a name + description until invoked).
//
//   experimental.chat.messages.transform — push the actual bookmark list. This
//     hook fires when opencode builds the message list for EACH outgoing LLM
//     request and mutates only that per-request payload — it is never written
//     back into stored session history. So the model sees exactly ONE fresh copy
//     of the list per turn (regenerated from `jumper list` each time) and nothing
//     accumulates. (A system-prompt injection would instead be cached per session
//     and go stale mid-session — hence the list lives here, not in system.)
//
// Why in-place mutation: opencode passes the same `messages` array reference on
// to model-message conversion; reassigning `output.messages = …` is a silent
// no-op (see anomalyco/opencode#25754). We therefore push onto an existing
// message's `parts` array in place.
//
// Best-effort throughout — if jumper is missing, errors, or returns nothing, we
// inject nothing and never disrupt the request.
//
// Opt-out: JUMPER_INJECT=false|0|no suppresses both injections.

import type { Plugin } from '@opencode-ai/plugin';

const INJECT_ENABLED = !['false', '0', 'no'].includes((process.env.JUMPER_INJECT ?? '').toLowerCase());

const BLOCK_OPEN = '<jumper-bookmarks>';
const BLOCK_CLOSE = '</jumper-bookmarks>';

// Static, high-level primer injected into the system prompt. Keeps the dynamic
// list (below) free of framing and ensures the agent understands the
// `<jumper-bookmarks>` block even without the jumper skill loaded.
const SYSTEM_PRIMER = [
  "jumper is the user's directory-bookmark CLI: short names mapped to directories.",
  `The current bookmarks are injected each turn in a \`${BLOCK_OPEN}\` block (name: path).`,
  'When the user refers to a project/dir by name, use that path directly instead of asking.',
].join('\n');

export const JumperInject: Plugin = async ({ $ }) => {
  // Fork `jumper list -o json` and render it as a compact `name: path` block.
  // Returns '' on any failure or when there are no bookmarks, so callers can
  // simply skip injection.
  async function renderBookmarks(): Promise<string> {
    try {
      const res = await $`jumper list -o json`.quiet().nothrow();
      if (res.exitCode !== 0) return '';

      const parsed = JSON.parse(res.stdout.toString());
      if (!Array.isArray(parsed) || parsed.length === 0) return '';

      const lines = parsed.filter((b: any) => b?.title && b?.path).map((b: any) => `- ${b.title}: ${b.path}`);
      if (lines.length === 0) return '';

      return [BLOCK_OPEN, ...lines, BLOCK_CLOSE].join('\n');
    } catch {
      return '';
    }
  }

  return {
    // Static primer — explains jumper once, in the system prompt.
    'experimental.chat.system.transform': async (_input, output) => {
      if (!INJECT_ENABLED) return;
      output.system.push(SYSTEM_PRIMER);
    },

    // Dynamic list — fresh bookmarks appended to the latest user message.
    'experimental.chat.messages.transform': async (_input, output) => {
      if (!INJECT_ENABLED) return;
      try {
        const messages = output.messages;
        if (!Array.isArray(messages) || messages.length === 0) return;

        // Anchor to the last user message — the most recent thing the model
        // reads — so the bookmark block sits right next to current intent.
        const target = [...messages].reverse().find((m) => m?.info?.role === 'user');
        if (!target) return;

        const block = await renderBookmarks();
        if (!block) return;

        // Append in place (reassigning output.messages is a no-op in opencode).
        target.parts.push({
          id: `jumper-inject-${Date.now()}`,
          sessionID: target.info.sessionID,
          messageID: target.info.id,
          type: 'text',
          text: block,
          synthetic: true,
        } as any);
      } catch {
        // best-effort — never disrupt the request
      }
    },
  };
};

export default JumperInject;
