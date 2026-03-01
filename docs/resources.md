# Reference Materials for Agentic Orchestration

## Core Articles
- Elvis Sun, "OpenClaw + Codex/Claude Agent Swarm: The One-Person Dev Team" (X, Feb 2026)
  - 已拷贝关键结构：两层上下文、tmux 监控、自动 reviewer。

## Internal Docs
- `.clawdbot/README.md`：本地 orchestrator 指南
- `docs/obsidian-setup.md`：Vault 结构
- （待建）`docs/agent-orchestration.md`：端到端流程（阶段 4 输出）

## Credential Inventory
- `~/.openclaw/.env`
  - `OPENAI_API_KEY`
  - `GOOGLE_API_KEY`
- Claude/Codex/Telegram key 位置：待补充

## Next Additions
- 会议纪要模版
- PR 模板
- 自动 reviewer CLI 使用手册
- `docs/zoe-operating-manual.md`：Zoe 运行手册
- `.clawdbot/templates/*`：prompt 模板（feature / bugfix / review）
- `.github/workflows/ci.yml`：CI baseline
- `.github/PULL_REQUEST_TEMPLATE.md`：PR 要求模板
- `docs/ci-review-plan.md`：CI & 审查策略
