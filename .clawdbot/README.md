# .clawdbot Orchestrator

This directory hosts the local orchestration layer (“Zoe”) for running specialized coding agents (Codex, Claude Code, Gemini) against this repository.

## Layout

- `config.sh` – shared configuration (repo path, default branch, package manager, CLI binaries, etc.)
- `active-tasks.json` – source of truth for every running agent (tmux session, worktree, PR status, etc.)
- `scripts/` – helper scripts launched manually or by cron: spawn agents, monitor health, clean up worktrees, notify via Telegram, etc.
- `templates/` *(optional later)* – prompt and task templates.

## Key Scripts

| Script | Purpose |
| --- | --- |
| `scripts/spawn_agent.sh` | Creates a worktree, installs deps, launches Codex/Claude/Gemini in a tmux session, and records the task registry entry. |
| `scripts/check_agents.sh` | Monitoring loop (cron every 10 min recommended). Verifies tmux session is alive, checks for PR + CI status through `gh`, flags tasks that need respawn or alert. |
| `scripts/cleanup_tasks.sh` | Kills finished tmux sessions, removes worktrees/branches, prunes registry entries. Run daily. |
| `scripts/send_notification.sh` | Placeholder notification hook (Telegram-ready). Called when tasks fail or finish. |
| `scripts/lib_tasks.sh` | JSON helper library used by other scripts to mutate `active-tasks.json`. |
| `scripts/gemini_cli.sh` + `gemini_cli.py` | Local Gemini CLI wrapper used by `spawn_agent.sh` when `agent=gemini`. |

## Gemini CLI Setup

1. Install the dependency (once per machine):
   ```bash
   pip install google-generativeai
   ```
2. Export your API key (add to shell profile):
   ```bash
   export GOOGLE_API_KEY="<your_gemini_key>"
   ```
3. Optionally override the default model in `config.sh` (`GEMINI_MODEL`).

After这些步骤，`spawn_agent.sh gemini ...` 会通过 `.clawdbot/scripts/gemini_cli.sh` 调用 Google Gemini。

## Example Workflow

1. Write a prompt file (e.g. `prompts/feat-template.md`) containing full context.
2. Spawn an agent:
   ```bash
   cd /Users/openclaw/dev/openclaw/Igloo/crazyLI
   ./.clawdbot/scripts/spawn_agent.sh codex feat-custom-templates feat/custom-templates "Custom template feature" prompts/feat-template.md
   ```
3. Cron (every 10 min):
   ```cron
   */10 * * * * cd /Users/openclaw/dev/openclaw/Igloo/crazyLI && ./.clawdbot/scripts/check_agents.sh >> .clawdbot/logs/check.log 2>&1
   ```
4. When CI + reviewers pass, mark task done / notify.
5. Daily cleanup:
   ```cron
   0 3 * * * cd /Users/openclaw/dev/openclaw/Igloo/crazyLI && ./.clawdbot/scripts/cleanup_tasks.sh >> .clawdbot/logs/cleanup.log 2>&1
   ```

The orchestrator keeps all business context in OpenClaw; these scripts focus on repeatable local automation, matching the workflow described in Elvis Sun’s Twitter thread.
