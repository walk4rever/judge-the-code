# judge-the-code

> 帮助人类在 AI 大量生成代码的时代，保持对代码的 **Judgment** 和 **Taste**。

[English](README.md) | 中文

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## 为什么需要这个工具

以前，写代码和理解代码是同一个动作。你写了什么，你就理解什么。

现在，AI 在写代码。写和懂被解耦了。

**AI 让代码能跑。但能跑不等于好。**

AI 生成的代码可能：
- 引入你没意识到的安全漏洞
- 破坏项目原有的设计哲学
- 埋下在 10 万用户时才爆的性能炸弹
- 无声地燃烧你的 LLM 预算——把 Token 费用变成黑洞
- 用了"有效"的捷径，制造下一个人踩不完的坑

发现这些，需要人真正理解代码库的 DNA——它的设计取向、历史决策、在乎什么。

这个理解，不能靠 lint，不能靠测试，**只能靠人的判断力**。

`judge-the-code` 是帮你维持这个判断力的工具。

---

## 两件事

```
Taste（欣赏力）                    Judgment（判断力）
──────────────────────────────────────────────────────
这里的设计很精妙，为什么？           这里有个坑，小心
这个抽象层级恰到好处               这个模式看起来干净，但会爆
这是一个值得学习的决策             这里有个隐性安全漏洞
这个 API 让错误用法很难发生         这个假设在高并发下会失效
                                  这个 Prompt 在烧 10 倍不该烧的 Token

让人看见代码的好                   让人看见代码的恶
```

---

## 架构

```
工具找问题（确定性）  +  Claude 解释问题（语义性）

Skill 层（Claude）          Tool 层（Go 二进制）
─────────────────────────────────────────────
code-explore               ← 分析目录结构、scc、syft
design-lens                ← 采样源码文件
demon-hunter  ←────────────── bearer / trivy / gitleaks
token-optimize             ← 静态分析 LLM 调用点
• 解读扫描结果                  确定性扫描，CVE 数据库
• 结合项目上下文判断             单二进制，setup 一行装好
• 解释为什么危险
• 给出修复建议
```

## Skills

| 组件 | 形态 | 作用 | 状态 |
|------|------|------|------|
| `code-explore` | Skill + 工具 | 建立代码库全局认知（结构、技术栈、入口、依赖）| ✅ 可用 |
| `design-lens` | Skill | 提取设计哲学与关键决策，找到值得学习和质疑的地方 | ✅ 可用 |
| `demon-hunter` | Skill + 工具 | 发现安全漏洞、依赖 CVE、密钥泄漏、性能隐患、设计陷阱 | ✅ 可用 |
| `token-optimize` | Skill | 发现 LLM 集成中的 Token 浪费——钱包黑洞、注意力污染、无意义的上下文膨胀 | ✅ 可用 |
| `skill-review` | Skill | 审查 Skill/Prompt 工程项目的质量——指令清晰度、Agent 编排、注入风险 | 🚧 MVP |

四个组件构成完整工作流：

```
code-explore  →  design-lens  →  demon-hunter  →  token-optimize
"这个项目       "哪里设计得好，   "哪里有恶魔"      "哪里在烧钱"
 长什么样"       为什么"
    结构层          欣赏层          判断层            经济层
```

### Skill/Prompt 路径：`skill-review`

随着 AI Agent 生态的发展，越来越多的项目不再是传统源代码，而是 **Skill 项目**：自然语言 Prompt、Agent 定义、执行流编排。现有工具（lint、SAST 扫描器、依赖审计）对这些完全无用。

`skill-review` 把 judge-the-code 的理念带到这个新前沿：

- **Prompt 清晰度** — 指令是否有歧义？低智商模型会不会理解出不同结果？
- **执行流设计** — Phase 分层是否合理？有没有死路或信息断层？
- **Agent 编排** — 并行 Agent 是否真正独立？有没有隐式串行依赖？
- **容错与降级** — Agent 失败或工具缺失时有没有 fallback？
- **安全边界** — Prompt Injection 风险？文件系统访问权限是否过宽？
- **模型兼容性** — 是否过度依赖某个模型的特性？

---

## 安装

先把 `judge-the-code` 目录放到你自己的 agent skills 加载路径里，然后在这个目录内部执行 setup。

```bash
cd /path/to/judge-the-code
./setup
```

`setup` 是原地初始化：

- 只为当前这份 `judge-the-code` 目录准备工具
- 不会帮你复制到任何全局 skills 路径

> 升级提示：用新版本替换当前目录后，重新执行 `./setup`。

## 使用方式

默认只需要记一个入口：

```bash
/judge-the-code .
/judge-the-code /path/to/project
```

根 skill 会自动分流：

- 代码仓库：`code-explore -> design-lens -> demon-hunter -> token-optimize`
- Skill/Prompt 项目：`skill-review`
- Hybrid 项目：自动识别并选择合适路径

自然语言也应该落到同一个总入口，例如：

- “帮我完整审查这个仓库”
- “帮我理解这个代码库并找风险”
- “审查一下这份 AI 生成的项目”
- “看看这个 Skill 项目有没有 Prompt 和编排问题”

### 高级用法

只有在你明确只想跑某个维度时，才直接调用子 skill：

```bash
/code-explore .       # 结构理解与上手
/design-lens .        # 设计哲学
/demon-hunter .       # 安全与隐患扫描
/token-optimize .     # LLM / Token 成本审查
/skill-review .       # Skill / Prompt 工程审查
```

### 非对话 CLI 入口

```bash
./bin/judge-the-code .
```

`judge-the-code` 是面向混合仓库的非对话统一入口：

- hybrid/skill：自动执行 `skill-review`
- hybrid/code：默认执行完整 baseline
- 全部产物统一落在 `TARGET/.judge-the-code/`

`run-judge` 作为底层实现脚本保留，用于兼容已有调用。

### Dashboard

```bash
./bin/view .
```

自动生成 `.judge-the-code/summary.html` 并在浏览器打开。

### 输出文件

```
.judge-the-code/
├── code-explore.md     ← code-explore 报告
├── design-lens.md      ← design-lens 报告
├── demon-hunter.md     ← demon-hunter 报告
├── token-optimize.md   ← token-optimize 报告
├── skill-review.md     ← skill-review 报告
├── summary.html        ← 可视化总览
└── state/              ← skill 内部状态（不用管）
```

---

## 适用场景

- **评估一个库要不要引入** — 不只看功能，还看坑
- **学习优秀项目的设计** — 带着批判性眼光，找到真正值得偷的东西
- **Review AI 生成的代码** — 验证没有破坏设计哲学，没有埋雷
- **接手陌生代码库** — 快速建立判断力，不只是走马观花
- **审计 LLM 集成成本** — 找到 Token 浪费、上下文膨胀、无意义的烧钱点

---

## 路线图

| 里程碑 | 说明 | 状态 |
|--------|------|------|
| code-explore | 代码库结构分析，Mermaid 可视化 | ✅ 已发布 |
| design-lens | 设计哲学提取与决策考古 | ✅ 已发布 |
| demon-hunter | 安全扫描（bearer + trivy + gitleaks）+ 语义分析 | ✅ 已发布 |
| token-optimize | LLM Token 浪费检测与优化建议 | ✅ 已发布 |
| code-explore 混合架构 | 确定性工具（scc + syft）用于架构与依赖分析 | 🚧 进行中 |
| skill-review | Skill/Prompt 工程项目的质量审查 | 🚧 MVP |

---

## 协议

[MIT](LICENSE) — 自由使用、修改和分发。
