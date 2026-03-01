#!/usr/bin/env bash
# Create a dedicated git worktree, install deps, and launch a coding agent in tmux.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
source "${ROOT_DIR}/.clawdbot/config.sh"
source "${SCRIPT_DIR}/lib_tasks.sh"

usage() {
  cat <<'EOF'
Usage: spawn_agent.sh <agent> <task-id> <branch-name> <description> <prompt-file>

agent: codex | claude | gemini
branch-name: feature branch to create from MAIN_BRANCH
prompt-file: text/markdown file containing the agent prompt
EOF
}

if [[ $# -lt 5 ]]; then
  usage
  exit 1
fi

AGENT="$1"
TASK_ID="$2"
BRANCH="$3"
DESCRIPTION="$4"
PROMPT_FILE="$5"

if [[ ! -f "$PROMPT_FILE" ]]; then
  echo "Prompt file not found: $PROMPT_FILE" >&2
  exit 1
fi

mkdir -p "$WORKTREES_DIR"
WORKTREE_DIR="${WORKTREES_DIR}/${TASK_ID}"
SESSION_NAME="${AGENT}-${TASK_ID}"

if git -C "$ROOT_DIR" worktree list | grep -q "${WORKTREE_DIR}"; then
  echo "Worktree already exists: ${WORKTREE_DIR}"
else
  git -C "$ROOT_DIR" fetch origin "$MAIN_BRANCH"
  git -C "$ROOT_DIR" worktree add -B "$BRANCH" "$WORKTREE_DIR" "origin/${MAIN_BRANCH}"
fi

pushd "$WORKTREE_DIR" >/dev/null
if [[ -f package.json ]]; then
  case "$PACKAGE_MANAGER" in
    pnpm) pnpm install ;;
    npm) npm install ;;
    yarn) yarn install ;;
    bun) bun install ;;
    *) echo "Unsupported PACKAGE_MANAGER=$PACKAGE_MANAGER" >&2; exit 1 ;;
  esac
fi
popd >/dev/null

if tmux has-session -t "$SESSION_NAME" 2>/dev/null; then
  echo "tmux session already exists: $SESSION_NAME" >&2
  exit 1
fi

mkdir -p "${ROOT_DIR}/.clawdbot/logs"
LOG_FILE="${ROOT_DIR}/.clawdbot/logs/${TASK_ID}.log"

declare AGENT_CMD
case "$AGENT" in
  codex)
    AGENT_CMD="cd ${WORKTREE_DIR} && ${CODEX_BIN} --model ${CODEX_MODEL} -c \"model_reasoning_effort=${CODEX_REASONING}\" --dangerously-bypass-approvals-and-sandbox \"$(cat "$PROMPT_FILE")\""
    ;;
  claude)
    AGENT_CMD="cd ${WORKTREE_DIR} && ${CLAUDE_BIN} --model ${CLAUDE_MODEL} --dangerously-skip-permissions -p \"$(cat "$PROMPT_FILE")\""
    ;;
  gemini)
    AGENT_CMD="cd ${WORKTREE_DIR} && ${GEMINI_BIN} --model ${GEMINI_MODEL} -p \"$(cat "$PROMPT_FILE")\""
    ;;
  *)
    echo "Unknown agent: $AGENT" >&2
    exit 1
    ;;
endcase

# Launch tmux session running the agent command, logging output via script(1)
tmux new-session -d -s "$SESSION_NAME" "bash -lc 'set -euo pipefail; ${AGENT_CMD} | tee \"${LOG_FILE}\"'"

add_task "$TASK_ID" "$DESCRIPTION" "$AGENT" "$BRANCH" "$WORKTREE_DIR" "$SESSION_NAME" "$TASK_REGISTRY"

echo "Spawned $AGENT agent in session $SESSION_NAME (worktree: $WORKTREE_DIR)"
