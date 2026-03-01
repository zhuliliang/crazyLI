# CI & Auto-Reviewer Plan

## CI Pipeline (GitHub Actions)
- Defined in `.github/workflows/ci.yml`.
- Jobs:
  1. `lint-test`: checkout → Node 20 → `pnpm install` → `pnpm lint` / `pnpm test` / `pnpm run e2e`.
  2. `ui-screenshot-check`: if PR body包含 `UI Changes: yes`，检查 PR 描述里是否出现 "Screenshot" 字样，否则失败。
- TODO: replace placeholder commands with actual scripts once codebase exists (e.g., `pnpm lint` → `pnpm run lint:ci`).

## PR Template
- `.github/PULL_REQUEST_TEMPLATE.md` 强制作者勾选测试状态、标记 UI changes，并贴截图链接。
- Zoe 在生成 prompt 时应提醒 agent 勾选对应项。

## Auto Reviewers
1. **Codex Reviewer**
   - Command placeholder: `codex review --model gpt-5.3-codex --pr <number>`.
   - Hook: run via GitHub Actions job or manual `check_agents.sh` hook when PR ready.
2. **Gemini Reviewer**
   - Use `gemini_cli.sh --model gemini-2.5-pro -f prompts/review-<pr>.md` to analyze `gh pr view -p` output.
3. **Claude Reviewer**
   - Placeholder CLI: `claude review --model claude-opus-4.5 --pr <number>`.

## Integration Points
- Extend `check_agents.sh` later to detect PR ready state → call reviewer scripts → post comments via `gh pr comment`.
- For now, manual step: run reviewers via CLI once PR passes CI, record results in `active-tasks.json`.

## Definition of Done (Recap)
- CI green (lint/types/tests/E2E)
- All three reviewers mark pass/acceptable
- Screenshots provided when UI changes occur
- Telegram notification sent → human merges
