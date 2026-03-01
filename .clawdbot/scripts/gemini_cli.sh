#!/usr/bin/env bash
# Thin wrapper so the orchestrator can call "gemini" like other CLIs.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PYTHON_BIN="python3"
CLI_PY="${SCRIPT_DIR}/gemini_cli.py"

exec "${PYTHON_BIN}" "${CLI_PY}" "$@"
