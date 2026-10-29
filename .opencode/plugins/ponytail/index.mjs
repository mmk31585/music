// ponytail — OpenCode V2 plugin.
// Injects the ponytail ruleset into every chat's system prompt at the active intensity,
// and persists /ponytail mode switches.

import { Plugin } from "@opencode/plugin";
import fs from 'fs';
import os from 'os';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

// OpenCode has no flag-file convention of its own; keep mode beside its config.
const statePath = path.join(
  process.env.XDG_CONFIG_HOME || path.join(os.homedir(), '.config'),
  'opencode',
  '.ponytail-active',
);

const VALID_MODES = ['lite', 'full', 'ultra', 'off'];
const DEFAULT_MODE = 'full';

function readMode() {
  try {
    const mode = fs.readFileSync(statePath, 'utf8').trim();
    return VALID_MODES.includes(mode) ? mode : DEFAULT_MODE;
  } catch (e) {
    return DEFAULT_MODE;
  }
}

function writeMode(mode) {
  fs.mkdirSync(path.dirname(statePath), { recursive: true });
  fs.writeFileSync(statePath, mode);
}

const PONYTAIL_INSTRUCTIONS = {
  lite: `
## Ponytail (lite)
Build what's asked, but name the lazier alternative in one line. User picks.
`,
  full: `
## Ponytail (full)
The ladder enforced. Stdlib and native first. Shortest diff, shortest explanation. Default.
Stop at the first rung that holds:
1. Does this need to exist at all? (YAGNI)
2. Already in this codebase? Reuse it.
3. Stdlib does it? Use it.
4. Native platform feature covers it? Use it.
5. Already-installed dependency solves it? Use it.
6. Can it be one line? One line.
7. Only then: minimum code that works.

No unrequested abstractions, no boilerplate, no scaffolding. Deletion over addition.
Mark deliberate simplifications with a // ponytail: comment.
Output: Code first, then at most three short lines.
  `,
  ultra: `
## Ponytail (ultra)
YAGNI extremist. Deletion before addition. Ship the one-liner and challenge the rest.
No cache until a profiler says so. Hand-rolled TTL cache class is a bug farm with a hit rate.
  `,
  off: ''
};

export default Plugin.define({
  id: "ponytail",
  async setup(ctx) {
    const ponytailSkillsDir = path.resolve(__dirname, '../skills');

    // Register skills directory so opencode discovers ponytail skills.
    await ctx.skill.transform((editor) => {
      // The skills.paths config handles this; transforms are for runtime edits.
      // We'll just ensure the skills are reloadable.
    });

    // Reload skills to pick up the ponytail skills directory
    await ctx.skill.reload();

    // Append the ruleset to the system prompt before every model request (context hook)
    await ctx.session.hook("context", (event) => {
      const mode = readMode();
      if (mode === 'off') return;
      const instructions = PONYTAIL_INSTRUCTIONS[mode] || PONYTAIL_INSTRUCTIONS.full;
      event.system.push({ type: "text", text: instructions });
    });

    // Register the /ponytail command to persist mode
    await ctx.command.transform((editor) => {
      editor.add({
        name: "ponytail",
        description: "Switch ponytail intensity level (lite/full/ultra/off)",
        execute: async ({ sessionID, prompt }) => {
          const args = prompt.text?.replace(/^\/ponytail\s*/, '').trim() || '';
          const mode = VALID_MODES.includes(args) ? args : DEFAULT_MODE;
          writeMode(mode);
          
          await ctx.session.prompt({
            sessionID,
            text: `Ponytail mode set to: ${mode}`,
            delivery: "steer",
          });
        },
      });
    });

    // Reload commands
    await ctx.command.reload();
  },
});