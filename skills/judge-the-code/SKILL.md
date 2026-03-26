---
name: judge-the-code
description: >-
  帮助人类在 AI 大量生成代码的时代，保持对代码的 Judgment 和 Taste。
  这是总入口 skill：自动在 code-explore、design-lens、demon-hunter、
  token-optimize、skill-review 之间编排与分流。TRIGGER when:
  用户想理解、评估、学习一个代码库，或想 review AI 生成的代码，
  且没有明确指定只跑某个子 skill。
origin: judge-the-code
---

# judge-the-code

> AI 让代码能跑。但能跑不等于好。

这是 `judge-the-code` 的总入口。对最终用户，默认只需要记一个命令：

```bash
/judge-the-code .
/judge-the-code /path/to/project
```

根 skill 负责判断目标类型，并自动选择合适的分析链路。子 skill 继续保留，但属于内部执行单元和高级用法，不应成为默认 onboarding。

## 入口职责

### 默认入口

```
/judge-the-code <项目路径>
```

如果用户没有给路径，默认按当前目录处理。

在执行任何分析前，先检查当前这份 skill 目录是否完成 setup。

- 若未完成 setup：立即停止，不进入后续分析
- 明确提示用户先在当前 `judge-the-code` 目录执行 `./setup`
- 不要让用户先看到“missing executable”这类底层报错

### 自然语言触发

以下表达都应优先触发本根 skill，而不是要求用户手动记子命令：

- 帮我完整审查这个仓库
- 看看这个项目值不值得学
- 分析这个 repo 的结构、设计和风险
- review 一下这份 AI 生成的代码
- 这个 Skill 项目有没有什么隐患

### 分流规则

收到 `/judge-the-code <path>` 后，先判断目标类型，再选择执行链路：

1. **Skill / Prompt 工程项目**
   - 特征：存在 `SKILL.md`、agent 定义、prompt 编排文件，且缺少传统应用入口
   - 执行：`skill-review`

2. **传统代码库**
   - 执行：`code-explore → design-lens → demon-hunter → token-optimize`

3. **Hybrid 项目**
   - 同时包含代码与 Skill 编排
   - 默认执行：传统代码库主链路
   - 若用户明确要求审查 Prompt/Skill 质量，追加 `skill-review`

## 代码库主链路

```
code-explore  →  design-lens  →  demon-hunter  →  token-optimize
"这个项目长什么样"    "哪里设计得好，为什么"      "哪里有恶魔"        "哪里在烧钱"
     结构层                 欣赏层                   判断层                 经济层
```

根 skill 应按以下顺序组织执行，而不是把选择权丢给用户：

- `code-explore`：建立结构认知
- `design-lens`：提炼设计哲学
- `demon-hunter`：发现安全与设计隐患
- `token-optimize`：发现 Token 浪费与上下文污染

各阶段输出统一落到 `.judge-the-code/`。

## Skill / Prompt 路径

若目标更像 Skill/Prompt 工程项目，则走：

```
/skill-review .
```

用于审查 Skill/Prompt 工程项目：Prompt 清晰度、执行流设计、Agent 编排、容错降级、安全边界与模型兼容性。
输出 `.judge-the-code/skill-review.md`。

## 高级用法

只有在用户明确要求某个维度，或你确认只需要局部分析时，才直接使用子 skill：

- `/code-explore .`
- `/design-lens .`
- `/demon-hunter .`
- `/token-optimize .`
- `/skill-review .`

若用户只是泛泛地说“帮我看看这个项目”，默认仍应优先使用 `/judge-the-code .` 这个单入口。

---

## 技术准则

所有 Agent 在做任何分析决策时，遵循以下优先级：

```
如实 > 速度 > 节省 token
```

### 输出语言

**跟用户的对话语言走**，与项目本身的语言无关：
- 用户用中文提问 → 所有报告、注释、说明用中文
- 用户用英文提问 → 所有报告、注释、说明用英文
- 混合语言 → 以用户最后一条消息的语言为准

### 如实（Evidence-based）— 最高优先级

- 每个结论必须有**文件路径 + 行号**作为证据
- 观察不到的东西，不猜测，直接说"未发现"
- 不确定时降级表述（"可能" / "倾向于"），而不是给出确定性结论
- 宁可输出"无法判断"，也不输出无根据的观点

### 速度（Speed）— 次优先

速度指用户从发出命令到看到有用输出的**整体感知耗时**，包含三层：

- **LLM 推理**：token 少 → 推理快；并行 Agent 而非串行
- **工具执行**：grep/bash 命令本身要快；未来混合架构（demon-hunter 调用 semgrep/trivy/Go CLI 等）中，优先选择启动快、增量扫描的工具，避免全量深度扫描
- **感知延迟**：尽早输出中间进度（如进度 banner），让用户知道任务在推进，而不是沉默等待

### 节省 token（Efficiency）— 最低优先

- 使用 `grep` + `sed` 精准提取，而非 `Read` 整个文件
- 用 `wc -l` 统计数量，而非把所有内容返回再手数
- **冲突规则**：如果节省 token 会导致证据不足、结论不可靠，宁可多用 token

---

## 适用场景

- **评估一个库要不要引入** — 不只看功能，还看坑
- **Review AI 生成的代码** — 验证没有破坏设计哲学，没有埋雷
- **学习优秀项目的设计** — 带着批判性眼光，找到真正值得偷的东西
- **接手陌生代码库** — 快速建立判断力，不只是走马观花
