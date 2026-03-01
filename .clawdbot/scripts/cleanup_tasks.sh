#!/usr/bin/env bash
# Removes finished worktrees, kills tmux sessions, prunes registry entries.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
source "${ROOT_DIR}/.clawdbot/config.sh"
source "${SCRIPT_DIR}/lib_tasks.sh"

ensure_registry

readarray -t TASKS < <(jq -c '.[]' "$TASK_REGISTRY")
if [[ ${#TASKS[@]} -eq 0 ]]; then
  echo "No tasks to clean up."
  exit 0
fi

for entry in "${TASKS[@]}"; do
  id=$(jq -r '.id' <<<"$entry")
  status=$(jq -r '.status' <<<"$entry")
  session=$(jq -r '.tmuxSession' <<<"$entry")
  worktree=$(jq -r '.worktree' <<<"$entry")
  branch=$(jq -r '.branch' <<<"$entry")

  if [[ "$status" != "done" && "$status" != "failed" ]]; then
    continue
  fi

  echo "Cleaning task $id ..."

  if tmux has-session -t "$session" 2>/dev/null; then
    tmux kill-session -t "$session"
    echo "  - tmux session killed"
  fi

  if [[ -d "$worktree" ]]; then
    git -C "$ROOT_DIR" worktree remove "$worktree" --force || true
    echo "  - worktree removed"
  fi

  if git -C "$ROOT_DIR" branch --list "$branch" >/dev/null 2>&1; then
    git -C "$ROOT_DIR" branch -D "$branch" || true
    echo "  - branch $branch deleted"
  fi

  remove_task "$id" "$TASK_REGISTRY"
  echo "  - registry entry removed"

done
