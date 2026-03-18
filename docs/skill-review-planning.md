# skill-review：Skill/Prompt 工程项目的专属审查

> 记录时间: 2026-03-18
> 状态: 规划阶段，待 code-explore 混合架构升级完成后启动

---

## 背景

随着 AI Agent 生态的发展，越来越多的项目不再是传统的源代码项目，而是 **Skill 项目**——由自然语言 Prompt、Agent 定义文件、执行流编排组成的工程产物（如 `~/.agents/skills/` 下的所有 skill）。

当前 judge-the-code 的全部 4 个子 skill（code-explore / design-lens / demon-hunter / token-optimize）都是围绕**源代码**设计的。它们对 Skill 项目**几乎完全不适用**：

| 子 Skill | 对 Skill 项目的适用性 | 原因 |
|---------|:------------------:|------|
| code-explore | ❌ | 找 `package.json`、入口函数？Skill 项目没有这些。`syft` 和 `scc` 在 markdown 上毫无意义。 |
| design-lens | ❌ | 分析命名、错误处理、抽象层级？Skill 项目是自然语言，没有函数和变量。 |
| demon-hunter | ❌ | `bearer`/`trivy`/`gitleaks` 扫 markdown？扫不出任何东西。Skill 的安全隐患是 Prompt Injection，不是 CVE。 |
| token-optimize | ❌ | 找 `messages.create`？Skill 本身就是 Prompt，它不调用 API，它被 API 消费。 |

**根本原因：源代码和自然语言指令是两种完全不同的工程产物，质量维度不同，审查方法不同。**

---

## Skill 项目的质量维度

| 维度 | 具体审查问题 | 对标传统代码中的什么 |
|------|------------|-----------------|
| **Prompt 清晰度** | 指令是否有歧义？不同智商的模型会不会理解出不同结果？是否存在可被忽略的关键步骤？ | 命名与可读性 |
| **执行流设计** | Phase 分层是否合理？有没有死路？Agent 之间有没有信息断层？ | 架构设计 |
| **输入验证** | Phase 0 是否覆盖了所有异常输入？路径为空、文件不存在、权限不足？ | 边界检查 |
| **Agent 编排** | 并行的 Agent 是否真正独立？有没有隐式依赖导致串行才正确？ | 并发设计 |
| **容错与降级** | 某个 Agent 失败了怎么办？工具缺失时有没有 fallback？ | 错误处理 |
| **Token 自身效率** | Skill 本身的 Prompt 是否臃肿？能不能更精简地表达同样的指令？ | 代码体积 |
| **安全边界** | 有没有 Prompt Injection 风险？是否给了 Agent 过宽的文件系统访问权限？用户输入是否与系统指令混在一起？ | 安全审计 |
| **模型兼容性** | 指令是否过度依赖某个模型的特性（如 Claude 的 XML 标签）？换一个模型会不会崩？ | 可移植性 |

---

## 架构方向

不改造现有 4 个子 skill，另起一个独立的子 skill `skill-review`：

```
judge-the-code
├── 传统源代码审查路径
│   ├── code-explore    (结构层)
│   ├── design-lens     (欣赏层)
│   ├── demon-hunter    (判断层)
│   └── token-optimize  (经济层)
│
└── Skill/Prompt 工程审查路径（新）
    └── skill-review    (专属审查)
```

两条路径共享 judge-the-code 的品牌和 dashboard 基础设施，但分析维度完全独立。

---

## 开放问题（待后续设计时回答）

1. `skill-review` 是否也需要确定性工具？还是纯 LLM 就够？（初步判断：大部分维度是语义判断，可能纯 LLM 为主。但 Token 效率维度或许可以用 `tiktoken` 类工具做精确计算。）
2. 如何自动区分"这是一个源代码项目"还是"这是一个 Skill 项目"？（可能的信号：根目录或子目录存在 `SKILL.md` 文件。）
3. Dashboard 如何呈现？新增一个独立 Tab，还是根据项目类型自动切换展示？
