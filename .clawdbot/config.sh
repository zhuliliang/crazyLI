#!/usr/bin/env bash
# Global configuration for the local orchestrator ("Zoe").

set -euo pipefail

# Absolute path to the repository root.
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Default branch name used when creating worktrees.
MAIN_BRANCH="main"

# Package manager used for dependency installs inside worktrees.
PACKAGE_MANAGER="pnpm"  # change to npm/yarn/bun if needed

# Binary names or paths for each coding agent CLI.
CODEX_BIN="codex"
CLAUDE_BIN="claude"
GEMINI_BIN="${REPO_ROOT}/.clawdbot/scripts/gemini_cli.sh"

# Default models.
CODEX_MODEL="gpt-5.3-codex"
CLAUDE_MODEL="claude-opus-4.5"
GEMINI_MODEL="gemini-2.5-pro"

# Reasoning effort presets.
CODEX_REASONING="high"
CLAUDE_REASONING="medium"
GEMINI_REASONING="medium"

# Where worktrees will be created. Each task gets its own directory here.
WORKTREES_DIR="${REPO_ROOT}/../worktrees"

# JSON registry path (shared across scripts).
TASK_REGISTRY="${REPO_ROOT}/.clawdbot/active-tasks.json"

# Optional: path to a script that sends Telegram/Signal notifications
NOTIFY_SCRIPT="${REPO_ROOT}/.clawdbot/scripts/send_notification.sh"

# Max auto-respawn attempts per task when monitoring fails.
MAX_RESPAWNS=3
