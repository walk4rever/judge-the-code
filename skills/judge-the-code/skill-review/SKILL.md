---
name: skill-review
description: >-
  审查 Skill/Prompt 工程项目质量，覆盖 Prompt 清晰度、执行流设计、Agent 编排、
  容错降级、安全边界和模型兼容性。输出 .judge-the-code/skill-review.md 报告。
  TRIGGER when: 用户要 review Skill 项目、Prompt 工程、agent 编排，或询问
  "这个 SKILL.md 写得好吗"、"有没有 prompt 注入风险"、"并行 agent 设计合理吗"。
origin: judge-the-code
---

# Skill Review — Skill/Prompt 工程审查

> 专门审查自然语言指令工程，而不是传统源代码。

## 调用方式

```bash
/skill-review .
/skill-review /path/to/skill-project
```

## 执行流程

### Phase 0：输入验证

```bash
TARGET="{用户提供的路径}"

if [ -z "$TARGET" ]; then
  echo "EMPTY"
elif [ ! -e "$TARGET" ]; then
  echo "NOT_FOUND"
elif [ -f "$TARGET" ]; then
  echo "IS_FILE"
elif [ -d "$TARGET" ]; then
  echo "OK"
fi
```

- `EMPTY`：停止，提示 `/skill-review .`
- `NOT_FOUND`：停止，提示路径不存在
- `IS_FILE`：停止，提示需要目录路径
- `OK`：继续

### Phase 1：结构发现（确定审查对象）

重点识别以下文件：
- `SKILL.md`
- `agents/*.md`
- `setup`
- `bin/*`

```bash
find "$TARGET" -maxdepth 4 -type f \
  \( -name "SKILL.md" -o -path "*/agents/*.md" -o -name "setup" -o -path "*/bin/*" \) \
  | head -200
```

若未发现 `SKILL.md`，提示“该目录可能不是 Skill 项目”，但允许继续做弱审查。

### Phase 2：规则扫描（确定性）

优先使用确定性扫描器：

```bash
if [ -x "{SKILL_DIR}/bin/lint-skill-review" ]; then
  "{SKILL_DIR}/bin/lint-skill-review" "{TARGET}" "{TARGET}/.judge-the-code/skill-review-findings.json"
else
  echo "NO_LINTER"
fi
```

- 若输出 `NO_LINTER`：使用 grep fallback 做最小规则检查。
- 若扫描器可用：读取 `.judge-the-code/skill-review-findings.json`，将 findings 作为 Phase 3 的高置信证据输入。

规则覆盖：
- 是否有 Phase 0 输入校验（`EMPTY/NOT_FOUND/IS_FILE/OK`）
- 是否有 fallback 逻辑（工具缺失/Agent 失败时）
- 是否存在危险命令（`rm -rf`, `git reset --hard`, `DROP TABLE`）
- 是否声明 Prompt Injection 边界（系统指令与用户输入隔离）
- 并行 Agent 是否存在隐式串行依赖
- Prompt 是否过长（高 token 冗余风险）

### Phase 3：4 Agent 并行语义审查

按以下规格并行执行：

| Agent | 规格文件 | 维度 |
|------|---------|------|
| Agent 1 | `agents/agent1-clarity.md` | Prompt 清晰度 |
| Agent 2 | `agents/agent2-flow.md` | 执行流设计 |
| Agent 3 | `agents/agent3-safety.md` | 安全边界与注入风险 |
| Agent 4 | `agents/agent4-orchestration.md` | 编排/容错/兼容性 |

### Phase 4：综合输出

优先使用一键流水线（自动生成报告并更新 dashboard）：

```bash
if [ -x "{SKILL_DIR}/bin/run-skill-review" ]; then
  "{SKILL_DIR}/bin/run-skill-review" "{TARGET}"
  exit 0
fi
```

若一键脚本不可用，则手动生成 `.judge-the-code/skill-review.md`，结构如下：

```markdown
# [项目名] — Skill Review 报告

> 生成时间: [date] | 分析工具: skill-review

## 总评
- 综合评分: [0-100]
- 等级: [A/B/C/D]
- 结论: [2-3 句]

## 八维评分
| 维度 | 分数 | 结论 |
|------|------|------|
| Prompt 清晰度 | 0-10 | ... |
| 执行流设计 | 0-10 | ... |
| 输入验证 | 0-10 | ... |
| Agent 编排 | 0-10 | ... |
| 容错与降级 | 0-10 | ... |
| Token 效率 | 0-10 | ... |
| 安全边界 | 0-10 | ... |
| 模型兼容性 | 0-10 | ... |

## Findings（按严重级别）
### Critical
### High
### Medium
### Low

## Top 5 Fixes（按 ROI）
1. ...
2. ...
3. ...
4. ...
5. ...
```

保存完成后：

```bash
find ~/.agents/skills ~/.claude/skills -name "SKILL.md" -path "*/skill-review/*" 2>/dev/null \
  | head -1 | xargs dirname
```

若存在历史脚本，追加本次运行快照：

```bash
if [ -x "{SKILL_DIR}/bin/update-skill-review-history" ]; then
  "{SKILL_DIR}/bin/update-skill-review-history" "{TARGET}" \
    "{TARGET}/.judge-the-code/skill-review.md" \
    "{TARGET}/.judge-the-code/skill-review-findings.json"
fi
```

若存在 `bin/view`，执行：

```bash
"{SKILL_DIR}/bin/view" .
```

## 质量标准

- 每个结论必须有路径与行号证据
- 没证据时写“未发现/无法判断”，不猜测
- 建议必须可执行，避免空泛措辞
