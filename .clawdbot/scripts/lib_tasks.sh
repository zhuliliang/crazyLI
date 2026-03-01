#!/usr/bin/env bash
# Helper functions for manipulating the task registry JSON.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
source "${ROOT_DIR}/.clawdbot/config.sh"

ensure_registry() {
  if [[ ! -f "${TASK_REGISTRY}" ]]; then
    echo "[]" >"${TASK_REGISTRY}"
  fi
}

add_task() {
  local id="$1"
  local description="$2"
  local agent="$3"
  local branch="$4"
  local worktree="$5"
  local session="$6"

  ensure_registry
  python3 - "$@" <<'PY'
import json, sys, time, pathlib
id, description, agent, branch, worktree, session = sys.argv[1:7]
registry_path = pathlib.Path(sys.argv[7])
entry = {
    "id": id,
    "description": description,
    "agent": agent,
    "branch": branch,
    "worktree": worktree,
    "tmuxSession": session,
    "status": "running",
    "respawns": 0,
    "startedAt": int(time.time() * 1000),
}
data = json.loads(registry_path.read_text())
existing = [t for t in data if t["id"] == id]
if existing:
    raise SystemExit(f"Task {id} already exists in registry.")
data.append(entry)
registry_path.write_text(json.dumps(data, indent=2))
PY
}

update_task_field() {
  local id="$1"
  local field="$2"
  local value="$3"
  ensure_registry
  python3 - "$@" <<'PY'
import json, sys, pathlib
id, field, value, registry = sys.argv[1:5]
path = pathlib.Path(registry)
data = json.loads(path.read_text())
updated = False
for task in data:
    if task.get("id") == id:
        task[field] = value
        updated = True
        break
if not updated:
    raise SystemExit(f"Task {id} not found.")
path.write_text(json.dumps(data, indent=2))
PY
}

update_task_status() {
  local id="$1"
  local status="$2"
  update_task_field "$id" "status" "$status"
  if [[ "$status" == "done" ]]; then
    local completed=$(python3 - <<'PY'
import time
print(int(time.time() * 1000))
PY
)
    update_task_field "$id" "completedAt" "$completed"
  fi
}

increment_respawns() {
  local id="$1"
  ensure_registry
  python3 - "$@" <<'PY'
import json, sys, pathlib
id, registry = sys.argv[1:3]
path = pathlib.Path(registry)
data = json.loads(path.read_text())
for task in data:
    if task.get("id") == id:
        task["respawns"] = task.get("respawns", 0) + 1
        path.write_text(json.dumps(data, indent=2))
        break
else:
    raise SystemExit(f"Task {id} not found.")
PY
}

remove_task() {
  local id="$1"
  ensure_registry
  python3 - "$@" <<'PY'
import json, sys, pathlib
id, registry = sys.argv[1:3]
path = pathlib.Path(registry)
data = json.loads(path.read_text())
filtered = [t for t in data if t.get("id") != id]
path.write_text(json.dumps(filtered, indent=2))
PY
}

list_tasks() {
  ensure_registry
  cat "${TASK_REGISTRY}"
}
