# Logging & Alerting Policy

1. **Locations**
   - `.clawdbot/logs/check.log` – cron output for monitoring loop
   - `.clawdbot/logs/cleanup.log` – daily cleanup summary
   - `.clawdbot/logs/<task>.log` – per-agent tmux output streamed by `spawn_agent.sh`

2. **Rotation**
   - Weekly manual rotation: `cp check.log check-2026-03-01.log && : > check.log`
   - Consider installing `logrotate` if logs exceed 100 MB.

3. **Alerting Rules**
   - Any `failed` or `needs-respawn` status → Telegram ping (already wired via `send_notification.sh` inside scripts)
   - `check_agents.sh` to add additional alerts as new failure modes emerge (CI failure, reviewer blockers, etc.)

4. **Audit Trail**
   - Keep `.clawdbot/logs/<task>.log` until feature merges + retro complete, then delete to protect storage.
