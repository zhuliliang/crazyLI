# Cron Setup for Monitoring & Cleanup

Add the following entries via `crontab -e` (adjust paths if repo moves):

```
# Every 10 minutes: agent monitoring loop
*/10 * * * * cd /Users/openclaw/.openclaw/workspace/Igloo/crazyLI && \
  ./.clawdbot/scripts/check_agents.sh >> .clawdbot/logs/check.log 2>&1

# Daily at 03:00: cleanup finished/failed tasks
0 3 * * * cd /Users/openclaw/.openclaw/workspace/Igloo/crazyLI && \
  ./.clawdbot/scripts/cleanup_tasks.sh >> .clawdbot/logs/cleanup.log 2>&1
```

**Tips**
- Run `crontab -l` to verify installation.
- Logs live in `.clawdbot/logs/`; review them daily.
- If using launchd instead of cron, convert the commands into plist jobs referencing the same scripts.
