# learn-from-github

帮助普通技术人快速理解陌生 GitHub 项目的 Claude Code Skill 集合。

## 核心 Skill：`/understand-repo`

通过 **5 个并行 Agent** 分析代码库，30 分钟内建立完整的项目心智模型。

### 使用方式

在 Claude Code 中：

```
/understand-repo .
/understand-repo /path/to/cloned/repo
```

或直接描述：
```
帮我理解一下这个项目 ~/projects/some-repo
分析一下 /path/to/project 这个仓库
```

### 5 个并行分析维度

| Agent | 分析内容 | 输出 |
|-------|---------|------|
| Stack Detector | 语言、框架、运行时 | 技术栈一览 |
| Architecture Mapper | 目录结构、设计模式 | 架构图 + 模块说明 |
| Entry Point Tracer | 主入口、启动流程 | 核心执行路径 |
| Dependency Analyst | 依赖库 + 使用意图 | 关键依赖解析 |
| Dev Setup Guide | 本地运行所需全部信息 | 快速启动手册 |

### 输出

- **`UNDERSTANDING.md`** — 保存到分析的项目根目录
- **交互式答疑** — 分析完成后进入专家陪伴模式，可直接提问

## 安装

```bash
# 安装 skill 到 Claude Code
cp -r skills/understand-repo ~/.claude/skills/understand-repo
```

## 目录结构

```
learn-from-github/
└── skills/
    └── understand-repo/
        └── SKILL.md    ← 核心 skill 定义
```
