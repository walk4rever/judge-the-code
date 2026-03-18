# code-explore 工具化改造决策 (混合架构)

> 决策时间: 2026-03-18
> 适用组件: code-explore (P1)

---

## 背景与痛点

当前的 `code-explore` 包含 5 个并行 Agent。它们的核心机制是通过提示词指导大模型使用 `grep`、`glob` 和读取前 80 行等方式去**猜测**项目的语言、架构和入口。

**当前纯 Prompt 模式的致命缺陷：**
1. **可靠性极低**：在复杂的 Monorepo、多语言混合、或者使用了别名（`@/utils`）的项目中，大模型的正则匹配和目录直觉经常失效。
2. **Context 浪费严重**：让大模型去翻 `package.json` 前 80 行来找框架，既慢又贵，且极易截断关键信息。
3. **能力天花板低**：纯 LLM 无法进行精准的 AST（抽象语法树）解析，找入口文件永远是靠“碰运气”看文件名。

**演进方向：混合架构 (Hybrid Architecture)**
像 `demon-hunter` 猎杀安全漏洞一样，`code-explore` 必须走向 **“工具主导收集，LLM 主导提炼”**。底层运行零依赖的 C/Go/Rust 确定性二进制工具生成 JSON，LLM 只负责将 JSON 翻译为人类可读的架构说明。

---

## 确定性工具选型矩阵

为了彻底解决上述痛点，我们选择以下 5 个全球顶尖的开源静态分析工具，去平替掉大模型盲目的探测行为。

### 核心原则：零侵入与零依赖
与 `demon-hunter` 的策略完全一致：**必须是单文件静态二进制（Static Binaries），支持 Mac/Linux/Win，无需用户系统安装 Go/Rust/Node 等任何环境，不污染用户的 PATH。**

| 替换对象 | 选定工具 | 语言 | 业界地位 / 理由 | 核心能力 |
|---------|---------|------|----------------|---------|
| **Agent 1** (技术栈探测) | [**enry**](https://github.com/go-enry/go-enry) | Go | GitHub 官方语言识别引擎 Linguist 的 Go 高性能版本。 | 毫秒级输出项目中各种语言/文件的绝对物理占比。不再靠猜 `package.json`。 |
| **Agent 2** (架构与代码分布) | [**scc**](https://github.com/boyter/scc) | Go | 世界上最快的代码行数与物理结构统计器。 | 穿透所有目录，精确吐出"哪个目录包含多少行真实业务逻辑"，帮大模型建立上帝视角，不被垃圾文件干扰。 |
| **Agent 3** (入口与时序追踪) | [**ast-grep**](https://github.com/ast-grep/ast-grep) | Rust | 爆火的跨语言 AST（语法树）结构化搜索神器。 | 不再用正则找 `listen()`。可以直接按语法树定义找"真正的启动函数"，100% 准确，无视注释和死代码。 |
| **Agent 4** (依赖树与生态) | [**syft**](https://github.com/anchore/syft) | Go | 全球最权威的 SBOM (软件物料清单) 生成器，Docker 官方采用。 | 一键输出完美的依赖解析树 JSON（穿透 node_modules/vendor）。 |
| **Agent 5** (开发环境解析) | [**yq**](https://github.com/mikefarah/yq) | Go | YAML 处理领域的绝对霸主。 | 确定性提取 `docker-compose.yml` 中的 GPU、环境变量依赖等深层嵌套配置，不再用易错的正则。 |

---

## 改造路径与隔离执行策略

改造将分为两个阶段，以 `enry` 和 `scc` 最先落地：

### 1. 隔离分发层 (`setup` 脚本)
复用 `demon-hunter/setup` 的逻辑。在 `skills/judge-the-code/code-explore/setup` 脚本中：
- 识别宿主架构（macOS/Linux, arm64/amd64）。
- 通过 `curl` 下载 `enry`, `scc`, `ast-grep` 等的 Release 压缩包。
- 解压后扔进 `skills/judge-the-code/bin/` 中。

### 2. 混合执行层
不再让 Agent 直接跑 `bash`。而是让 `code-explore` 启动时，由 `bash` 预先执行：
```bash
# 获取极度准确的技术栈与代码量 JSON
{SKILL_DIR}/bin/enry > .judge-the-code/state/lang-dist.json
{SKILL_DIR}/bin/scc -f json > .judge-the-code/state/arch-dist.json
```
**Agent 的 Prompt 演进：**
*以前：* "请使用 grep 查找 `package.json` 前 80 行并猜框架..."
*现在：* "请读取 `lang-dist.json` 和 `arch-dist.json`。由于这是底层工具 100% 准确的扫描结果，你只需总结它的工程特点，并挑选代码量最大的前 3 个目录作为核心模块绘制 Mermaid 架构图。"

---

## 预期收益

1. **Token 成本暴降**：LLM 不再需要读取数万字符的源码和长配置。
2. **速度提升 10 倍**：底层工具的扫描都在毫秒级别。
3. **彻底消除幻觉**：依赖图谱、入口位置、代码占比成为绝对客观的真理，不会因为大模型降级（如切换至 Haiku 或本地模型）而产生分析降级。