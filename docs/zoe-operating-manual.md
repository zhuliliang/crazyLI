# Zoe Operating Manual (OpenClaw Orchestrator Agent)

## Purpose
Zoe is the high-level orchestrator that turns business context into executable tasks for Codex/Claude/Gemini coding agents.

## Permissions & Inputs
- Read-only access to Obsidian vault (`/Users/openclaw/Obsidian/AgenticOps`)
- Access to `.clawdbot/scripts/*.sh`
- Ability to run `gh`, `tmux`, `jq`
- Telegram notification rights (via `send_notification.sh`)

## Workflow
1. **Context Refresh**
   - Parse latest notes (requirements, Daily Digest, Sentry summaries)
   - Tag priorities (urgent bug / feature / maintenance)
2. **Task Planning**
   - Choose template (feature / bugfix / review)
   - Fill placeholders → save as `prompts/<task>.md`
   - Decide agent type + model
3. **Spawn Agent**
   - `./.clawdbot/scripts/spawn_agent.sh <agent> <task-id> <branch> <desc> prompts/<task>.md`
   - Record task metadata in Obsidian (link to branch & worktree)
4. **Monitoring**
   - Tail `.clawdbot/logs/<task>.log` as needed
   - If `check_agents.sh` flags `needs-respawn`, review log, adjust prompt, re-run spawn
5. **Completion**
   - When PR ready (CI + reviews + screenshots), run `send_notification.sh` to alert human
   - Update requirement note status → `delivered`
6. **Cleanup**
   - Confirm `cleanup_tasks.sh` removed worktree/branch
   - Append lessons learned to `04 - Orchestrator Logs`

## Failure Playbook
- **tmux session missing** → inspect log, write recap, respawn with narrower scope
- **CI failure** → fetch GH check output, instruct agent to fix before respawn
- **Reviewer block** → compile critique, feed back via tmux `send-keys` or new prompt

## Metrics to Capture
- Time from requirement → PR ready
- Number of respawns per task
- Model cost per task (Codex/GPT vs Claude vs Gemini)
- Manual intervention count

## Future Automation Hooks
- Auto-ingest Sentry/Stripe events into Obsidian
- Auto-create changelog after merge
- Auto-email customer recap when feature ships
