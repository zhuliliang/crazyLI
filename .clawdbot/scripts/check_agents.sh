#!/usr/bin/env bash
# Monitoring loop similar to the article: verifies tmux sessions, PRs, CI, reviewers.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
source "${ROOT_DIR}/.clawdbot/config.sh"
source "${SCRIPT_DIR}/lib_tasks.sh"

require_cli() {
  for bin in "$@"; do
    if ! command -v "$bin" >/dev/null 2>&1; then
      echo "Warning: $bin not found in PATH." >&2
    fi
  done
}

require_cli tmux gh jq

ensure_registry

TASK_COUNT=$(jq 'length' "$TASK_REGISTRY")
if [[ "$TASK_COUNT" -eq 0 ]]; then
  echo "No active tasks."
  exit 0
fi

while IFS= read -r entry; do
  id=$(jq -r '.id' <<<"$entry")
  agent=$(jq -r '.agent' <<<"$entry")
  worktree=$(jq -r '.worktree' <<<"$entry")
  session=$(jq -r '.tmuxSession' <<<"$entry")
  branch=$(jq -r '.branch' <<<"$entry")
  respawns=$(jq -r '.respawns // 0' <<<"$entry")
  status=$(jq -r '.status' <<<"$entry")

  echo "→ Checking task ${id} (${agent})"

  if [[ "$status" == "done" ]]; then
    echo "  - Already done. Skipping."
    continue
  fi

  # tmux session alive?
  if ! tmux has-session -t "$session" 2>/dev/null; then
    echo "  - tmux session missing."
    if (( respawns < MAX_RESPAWNS )); then
      echo "  - Marking for respawn."
      increment_respawns "$id" "$TASK_REGISTRY"
      update_task_status "$id" "needs-respawn"
    else
      echo "  - Max respawns reached. Marking as failed."
      update_task_status "$id" "failed"
      [[ -x "$NOTIFY_SCRIPT" ]] && "$NOTIFY_SCRIPT" "Task ${id} failed: tmux session gone."
    fi
    continue
  fi

  # PR status via gh (looks for open PR targeting branch)
  pr_json=$(gh pr list --state open --head "$branch" --json number,title,state 2>/dev/null || echo "[]")
  pr_count=$(jq 'length' <<<"$pr_json")
  if (( pr_count > 0 )); then
    pr_number=$(jq -r '.[0].number' <<<"$pr_json")
    echo "  - PR #${pr_number} is open."
    update_task_field "$id" "pr" "$pr_number"
  else
    echo "  - No PR open yet."
  fi

  # CI check (requires gh)
  if [[ -n "${pr_number:-}" ]]; then
    status_checks=$(gh pr checks "$pr_number" --json status 2>/dev/null || echo '{}')
    ci_state=$(jq -r '.status // empty' <<<"$status_checks")
    [[ -n "$ci_state" ]] && update_task_field "$id" "ci" "$ci_state"
    echo "  - CI status: ${ci_state:-unknown}" | sed 's/\n/ /g'

    if [[ "$ci_state" == "SUCCESS" ]]; then
      echo "  - CI passed. You can soon mark task done after reviews."
    fi
  fi

done < <(jq -c '.[]' "$TASK_REGISTRY")
