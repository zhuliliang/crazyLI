# Obsidian Vault Setup for Agent Orchestration

Use this file as the seed note inside your Obsidian vault. It mirrors the Stage 0 preparation tasks in我们的执行计划。

## 1. Vault Structure

推荐顶层文件夹：

- `01 - Requirements`：客户请求、会议纪要、支持票据
- `02 - Engineering Context`：架构决策、系统图、依赖清单
- `03 - Operations`：Sentry 摘要、Stripe/CRM 快照、操作手册
- `04 - Orchestrator Logs`：Zoe prompt 历史、agent 结果、复盘
- `Templates`：prompt 模板、PR 检查表、review 指南

## 2. 命名规范

`YYYY-MM-DD - 简短主题`，如 `2026-02-28 - Customer Template Request`。

## 3. Frontmatter 示例

```yaml
---
type: requirement
source: customer_call
status: queued
priority: high
linked_task: feat-custom-templates
---
```

## 4. 每日流程

1. 客户会议后立即更新/创建 requirement note。
2. 每日结束写一条 `Daily Digest` 总结关键事件。
3. Zoe 读取 `Daily Digest` + 未完成需求决定次日任务。

## 5. 同步提示

- Vault 放在本地，例如 `/Users/openclaw/Obsidian/AgenticOps`，确保 Zoe 可访问。
- 如需版本控制，可把 Vault 作为 git submodule。

设置完该结构后，Zoe 就能系统性地解析需求→生成 prompt。
