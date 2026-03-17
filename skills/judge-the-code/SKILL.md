---
name: judge-the-code
description: >-
  帮助人类在 AI 大量生成代码的时代，保持对代码的 Judgment 和 Taste。
  包含三个渐进式 skill：understand-repo（建立结构认知）、
  philosophy-extractor（提炼设计哲学）、demon-hunter（发现安全漏洞与设计陷阱）。
  TRIGGER when: 用户想理解、评估、学习一个代码库，或想 review AI 生成的代码。
origin: judge-the-code
version: 0.4.1
---

# judge-the-code

> AI 让代码能跑。但能跑不等于好。

这个 skill 套件帮你保持对代码的判断力。

## 三层工作流

```
understand-repo  →  philosophy-extractor  →  demon-hunter（规划中）
"这个项目长什么样"    "哪里设计得好，为什么"      "哪里有恶魔"
     结构层                 欣赏层                   判断层
```

### 第一步：建立结构认知

```
/understand-repo .
/understand-repo ~/projects/some-repo
```

5 个并行 Agent 分析技术栈、架构、入口、依赖、开发环境。
输出 `UNDERSTANDING.md` + 渐进式导览模式。

### 第二步：提炼设计哲学

```
/philosophy-extractor .
```

4 个并行 Agent 分析命名风格、错误处理、测试取向、架构决策。
输出 `PHILOSOPHY.md`，每条决策打标签：🔮 精妙 / ✅ 合理 / ⚠️ 存疑 / ❌ 反模式。

### 第三步：猎杀恶魔（规划中）

```
/demon-hunter .
```

结合 semgrep / trivy 等工具 + Claude 语义分析，发现安全漏洞、性能隐患、设计陷阱。

---

## 技术准则

所有 Agent 在做任何分析决策时，遵循以下优先级：

```
如实 > 速度 > 节省 token
```

### 如实（Evidence-based）— 最高优先级

- 每个结论必须有**文件路径 + 行号**作为证据
- 观察不到的东西，不猜测，直接说"未发现"
- 不确定时降级表述（"可能" / "倾向于"），而不是给出确定性结论
- 宁可输出"无法判断"，也不输出无根据的观点

### 速度（Speed）— 次优先

- 所有 Agent **并行执行**，不串行等待
- 优先用单次 `grep` 命令覆盖全库，而不是多次读取单个文件
- 减少串行工具调用数量

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
