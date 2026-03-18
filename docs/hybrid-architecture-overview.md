# 全局工具化决策：哪些 Skill 需要工具，哪些不需要

> 决策时间: 2026-03-18
> 适用范围: judge-the-code 全部 4 个子 Skill

---

## 核心原则

**凡是有 ground truth（客观真相）的任务，用工具；凡是需要 judgment（主观判断）的任务，用 LLM。**

judge-the-code 的四层工作流中，每一层的本质不同：

| 层 | Skill | 任务本质 | 是否有客观真相 | 结论 |
|----|-------|---------|:----------:|------|
| 结构层 | code-explore | 提取事实（依赖是什么、代码量多少） | ✅ 是 | **工具 + LLM 混合** |
| 欣赏层 | design-lens | 审美判断（设计好不好、为什么） | ❌ 否 | **纯 LLM** |
| 判断层 | demon-hunter | 扫描漏洞（CVE、密钥、注入） | ✅ 是 | **工具 + LLM 混合** |
| 经济层 | token-optimize | 经验推理（这个调用值不值、会不会爆） | ❌ 否 | **纯 LLM** |

---

## 各 Skill 决策详情

### code-explore — 工具 + LLM 混合

**引入的工具：** `syft`（依赖分析）+ `scc`（代码物理分布）

**理由：** 依赖树和代码行数是客观事实，LLM 读 80 行配置文件来猜测既不准确也浪费 Token。工具可以 100% 准确地提取这些数据，LLM 只需基于确定数据做总结和画图。

**不引入工具的 Agent：** Agent 1（技术栈）、Agent 3（入口追踪）、Agent 5（开发环境）。这三个维度当前的 grep/glob + LLM 推理方案已经足够可靠。

详见 → `docs/code-explore-tool-selection.md`

---

### demon-hunter — 工具 + LLM 混合

**引入的工具：** `bearer`（SAST 代码漏洞）+ `trivy`（依赖 CVE + 密钥）+ `gitleaks`（Git 历史密钥）

**理由：** CVE 编号、密钥泄漏、注入漏洞都有客观判定标准（要么存在要么不存在），且需要依赖 CVE 数据库等 LLM 不具备的外部知识。工具的确定性扫描结果是报告的骨架，LLM 负责解释"为什么危险"和"怎么修"。

详见 → `docs/demon-hunter-tool-selection.md`

---

### design-lens — 纯 LLM

**不引入任何工具。**

design-lens 的 4 个 Agent 分析的是命名风格、错误处理取向、测试信仰、架构决策。这些任务的本质是**主观审美判断**——"这个命名揭示了一个务实的工程师"、"这个抽象层级过深属于过度设计"——没有任何确定性工具能输出这种洞察。

**曾考虑但否决的方向：**

| 候选思路 | 否决理由 |
|---------|---------|
| 圈复杂度工具（`gocyclo`, `radon`） | 只输出数字，"15 算不算高"取决于上下文，最终还是 LLM 判断。为一个数字引入二进制不值得。 |
| 代码风格检查器（`eslint`, `golint`） | 检查"是否符合规则"，而 design-lens 关注"规则背后的设计哲学是什么"。层次完全不同。 |

---

### token-optimize — 纯 LLM

**不引入任何工具。**

token-optimize 的 3 个 Phase（找调用点 → 推演隐患 → 优化建议）全部是经验推理任务。LLM API 的调用特征词（`messages.create`、`chat.completions`）极度明显，grep 准确率接近 100%；推演"如果数组无限增长会怎样"和建议"降级为 Haiku"都是纯语境推理，没有工具能代劳。

**曾考虑但否决的方向：**

| 候选思路 | 否决理由 |
|---------|---------|
| `tiktoken`（精确 Token 计数） | 只有 Python 库，无独立二进制。且关注的不是精确到个位的 token 数，而是"这个模式会不会爆炸"——量级判断 LLM 自己能做。 |
| `ast-grep`（精确定位 API 调用） | LLM API 调用的特征词极度明显，grep 准确率接近 100%，不需要 AST 级别精度。 |

---

## 总览

```
code-explore   →  syft + scc                        （客观事实提取）
demon-hunter   →  bearer + trivy + gitleaks          （客观漏洞扫描）
design-lens    →  纯 LLM                             （主观审美判断）
token-optimize →  纯 LLM                             （经验推理）
```

**判断力这件事，终究还是要靠大脑。工具只负责收集大脑需要的证据。**
