# Stage 2 Report – Zoe Orchestrator Setup

Date: 2026-02-28

## Deliverables
1. **Prompt Templates** (`.clawdbot/templates/`)
   - `feature-prompt.md`: includes context, scope, DoD checklist.
   - `bugfix-prompt.md`: incident metadata, repro steps, SLA notes.
   - `review-prompt.md`: structured review checklist/output format.
2. **Prompt Storage Guide** (`prompts/README.md`)
   - Naming guidance + required metadata for each prompt file.
3. **Zoe Operating Manual** (`docs/zoe-operating-manual.md`)
   - Permissions, workflow, failure playbook, metrics, future hooks.

## Pending Items / Constraints
- Demo run of `spawn_agent.sh` is deferred because actual coding agents (Codex/Claude CLI) are not yet configured and the repository contains no code to operate on. Recommendation: perform the dry-run once code + CLI binaries are in place to avoid misleading failures.
- Zoe agent in OpenClaw still needs to be instantiated (with vault permissions). Manual describes the expected behavior.

Stage 2 artifacts are ready for review. Once CLI tooling and repo contents exist, we can proceed to Stage 3 (monitoring automation) and actually bind Zoe into OpenClaw.
