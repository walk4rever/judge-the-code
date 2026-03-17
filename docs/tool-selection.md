# demon-hunter 工具选型决策

> 决策时间: 2026-03-17
> 适用组件: demon-hunter (P3)

---

## 工具选型

### 核心三件套

| 工具 | 类别 | 检测范围 | 付费情况 | 市场地位 |
|------|------|---------|---------|---------|
| **semgrep** | SAST | 代码安全漏洞（注入/XSS/越权/硬编码密钥）| CLI 完全免费；Pro $25/user/mo（托管规则+dashboard，我们不需要）| 开发者 SAST 首选；Dropbox/Snowflake/Trail of Bits 在用 |
| **trivy** | SCA + Secrets + IaC | 依赖 CVE（npm/pip/cargo/go/maven...）+ 硬编码密钥（当前文件）+ 容器/IaC 配置错误 | 完全免费，Aqua Security 出品 | SCA 领域覆盖最全的免费工具 |
| **gitleaks** | Secrets | Git history 中的密钥泄漏（trivy 不扫历史）| 完全免费 | 600+ 检测规则，git history 扫描最佳 |

### 未来可加（增强覆盖，暂不实现）

| 工具 | 理由 | 时机 |
|------|------|------|
| **osv-scanner** | Google 出品，OSV 数据库（与 GitHub Dependabot 同源），CVE 更新更快；与 trivy 互补降低漏报 | trivy 漏报问题被用户反馈时 |

---

## 覆盖分析

| Demon 类型 | 工具 | 覆盖情况 |
|-----------|------|---------|
| 代码安全漏洞（注入/XSS/越权）| semgrep | ✅ 规则匹配，全语言 |
| 依赖 CVE | trivy | ✅ 全生态 |
| 硬编码密钥（当前代码）| trivy + gitleaks | ✅ 双重覆盖 |
| Git history 密钥泄漏 | gitleaks | ✅ 唯一能扫历史的 |
| 容器/IaC 配置错误 | trivy | ✅ |
| 性能隐患（N+1/无索引）| Claude | ✅ grep 模式匹配 |
| 隐性耦合/设计陷阱 | Claude | ✅ 语义分析 |
| 哲学破坏 | Claude | ✅ 读 PHILOSOPHY.md |
| **深层语义漏洞（跨函数数据流）** | — | ⚠️ semgrep 是规则匹配，无法做数据流分析；CodeQL 能做但需要 GitHub Enterprise，不引入 |
| **许可证合规** | — | ⚠️ trivy 有基础支持，非重点 |

---

## 用户环境多样性分析

demon-hunter 需要在以下所有环境下可用：

| 维度 | 需要覆盖的情况 |
|------|-------------|
| **操作系统** | macOS (Intel x86_64 / Apple Silicon arm64) · Linux (Ubuntu/Alpine/Fedora...) · WSL2 |
| **架构** | x86_64(amd64) · arm64/aarch64 · 未来可能的 arm/v7 |
| **包管理器** | brew(macOS only) · apt/yum/apk · 或完全没有 |
| **Python 环境** | 有 Python ≥ 3.9 · 有 uv · 有 pipx · 无任何 Python 工具 |
| **网络** | 直连 · 企业代理 · 离线/受限 |
| **权限** | 完整权限 · 无 sudo（企业机器常见）|

**核心矛盾**：

- `trivy` / `gitleaks` 是 **Go 单二进制**，覆盖所有 OS+arch，直接下载，问题最小
- `semgrep` **无独立二进制**（仅 pip 分发），是最大的环境多样性风险

**semgrep 在各环境的可用性**：

| 安装方式 | macOS | Linux | WSL2 | 无 Python | 无网络 |
|---------|:-----:|:-----:|:----:|:---------:|:------:|
| `uvx semgrep` | ✅ 若有 uv | ✅ 若有 uv | ✅ 若有 uv | ✅ uv 自带 Python | ❌ |
| `pip install` | ✅ | ✅ | ✅ | ❌ | ❌ |
| `brew install` | ✅ | ❌ | ❌ | ✅ | ❌ |
| `docker run` | ✅ 若有 Docker | ✅ | ✅ | ✅ | ❌ |
| **skill 内置 venv** | ✅ 需 Python | ✅ 需 Python | ✅ 需 Python | ❌ | ❌ |

**结论**：没有任何一种方式能覆盖所有环境。必须使用**有序 fallback 链**。

---

## 隔离执行策略

### 核心原则

**不侵入用户的系统级环境。** 所有工具运行在 skill 自管理的隔离目录内。

### 方案

参考 gstack 的 `browse/dist/browse` 模式——skill 目录内自带二进制。

```
~/.agents/skills/judge-the-code/
└── demon-hunter/
    ├── SKILL.md
    ├── setup              ← 一次性安装脚本（类似 gstack/setup）
    └── bin/               ← skill 自管理的工具目录，不写入系统 PATH
        ├── trivy          ← Go 单二进制，直接从 GitHub Releases 下载
        ├── gitleaks       ← Go 单二进制，直接从 GitHub Releases 下载
        └── semgrep        ← 见下方说明
```

### 各工具的隔离方式

**trivy 和 gitleaks**（Go 单二进制）：
- 直接从 GitHub Releases 下载到 `SKILL_DIR/bin/`
- 无任何依赖，无系统污染
- 按需下载，首次运行时触发

```bash
# trivy
curl -sSfL "https://github.com/aquasecurity/trivy/releases/download/v${VER}/trivy_${VER}_${OS}_${ARCH}.tar.gz" \
  | tar -xz -C "$SKILL_DIR/bin" trivy

# gitleaks
curl -sSfL "https://github.com/gitleaks/gitleaks/releases/download/v${VER}/gitleaks_${VER}_${os}_${arch}.tar.gz" \
  | tar -xz -C "$SKILL_DIR/bin" gitleaks
```

**semgrep**（Python 工具，无独立二进制分发）：

优先级检测，按顺序尝试：

```
优先级 1：uvx semgrep        — 若 uv 已安装，最优选择
                               uv 自动创建隔离临时环境，无系统污染
                               运行后环境自动清理
                               uv 是现代 Python 工具链标准，越来越普及

优先级 2：SKILL_DIR/.venv    — skill 内置 venv（uv venv 创建）
          uv venv $SKILL_DIR/.venv && uv pip install semgrep
          完全隔离，不影响系统 Python 和其他项目

优先级 3：提示安装            — 以上都不可用时，给出安装建议
          "需要安装 uv（推荐）或 semgrep：brew install uv"
```

### trivy / gitleaks 的二进制下载策略

Go 单二进制，按 OS + arch 自动选择正确包：

```bash
OS=$(uname -s)    # Darwin | Linux
ARCH=$(uname -m)  # x86_64 | arm64 | aarch64
[ "$ARCH" = "aarch64" ] && ARCH="arm64"   # Linux arm64 归一化

# trivy 命名规则: trivy_VERSION_macOS-64bit.tar.gz / trivy_VERSION_Linux-ARM64.tar.gz
# gitleaks 命名规则: gitleaks_VERSION_darwin_arm64.tar.gz / gitleaks_VERSION_linux_x64.tar.gz
```

覆盖矩阵：

| OS | arch | trivy 包名 | gitleaks 包名 |
|----|------|-----------|--------------|
| macOS | x86_64 | `macOS-64bit` | `darwin_x64` |
| macOS | arm64 | `macOS-ARM64` | `darwin_arm64` |
| Linux | x86_64 | `Linux-64bit` | `linux_x64` |
| Linux | arm64 | `Linux-ARM64` | `linux_arm64` |

### semgrep 的有序 fallback 链

```
检测顺序（优先使用隔离方式，最后才是系统级安装）：

1. SKILL_DIR/.venv/bin/semgrep  → 之前已在 skill 内创建过 venv（最优：完全隔离，无需网络）
2. uvx semgrep                  → uv 已安装（优秀：隔离，自动管理）
3. pipx run semgrep             → pipx 已安装（良好：隔离）
4. docker run semgrep/semgrep   → Docker 已安装（可用：完全隔离）
5. python3 -m semgrep           → semgrep 已用 pip 装过（可用：系统级）
6. semgrep                      → 系统 PATH 中有（可用）
7. → 触发安装引导               → 根据检测到的可用工具给出建议
```

### 安装引导（fallback 最终兜底）

当 semgrep 完全不可用时，给出环境感知的建议而非一条命令：

```
⚠️ 未找到 semgrep，SAST 扫描将跳过。

检测到你的环境：
  ✅ uv 已安装  →  推荐：一次性运行 `uv tool install semgrep`
  ✅ Python 3.x →  备选：`pip install --user semgrep`
  ❌ brew 未安装

其余扫描（trivy + gitleaks）继续执行。
```

### 整体执行入口（find-tools 脚本）

```bash
#!/usr/bin/env bash
# SKILL_DIR/bin/find-tools
# 输出各工具的可执行路径，供 SKILL.md 引用

SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$SKILL_DIR/bin"

# trivy
if [ -x "$BIN/trivy" ]; then echo "TRIVY=$BIN/trivy"
elif command -v trivy &>/dev/null; then echo "TRIVY=$(command -v trivy)"
else echo "TRIVY=MISSING"; fi

# gitleaks
if [ -x "$BIN/gitleaks" ]; then echo "GITLEAKS=$BIN/gitleaks"
elif command -v gitleaks &>/dev/null; then echo "GITLEAKS=$(command -v gitleaks)"
else echo "GITLEAKS=MISSING"; fi

# semgrep（有序 fallback）
if [ -x "$SKILL_DIR/.venv/bin/semgrep" ]; then echo "SEMGREP=$SKILL_DIR/.venv/bin/semgrep"
elif command -v uvx &>/dev/null; then echo "SEMGREP=uvx semgrep"
elif command -v pipx &>/dev/null; then echo "SEMGREP=pipx run semgrep"
elif command -v docker &>/dev/null; then echo "SEMGREP=docker run --rm -v \$(pwd):/src semgrep/semgrep"
elif command -v semgrep &>/dev/null; then echo "SEMGREP=$(command -v semgrep)"
else echo "SEMGREP=MISSING"; fi
```

### 缺失工具的处理原则

**不因某个工具缺失而整体失败**，而是：

1. 输出哪些维度被跳过及原因
2. 给出针对当前环境的安装建议
3. 已可用的工具继续扫描，输出部分结果

```
🔍 demon-hunter 扫描报告

✅ 依赖 CVE 扫描 (trivy)      — 完成，发现 3 个高危漏洞
✅ Git 密钥扫描 (gitleaks)    — 完成，未发现泄漏
⚠️ 代码漏洞扫描 (semgrep)    — 已跳过（未安装）
   → 安装建议：uv tool install semgrep
✅ 设计陷阱分析 (Claude)      — 完成，发现 2 个隐患
```

---

## 不选择的方案及理由

| 方案 | 否决理由 |
|------|---------|
| 系统级 `brew install` / `pip install` | 侵入用户系统环境，违反隔离原则 |
| Docker | 用户不一定装了 Docker，镜像拉取慢，overhead 大 |
| Snyk | 免费版有调用次数限制（200次/月），不可靠 |
| CodeQL | 需要 GitHub Enterprise 或 Actions 环境，无法本地独立运行 |
| GitGuardian | 付费，且需要账号，适合团队而非个人开发者 |
| SonarQube | 重量级，需要独立服务，不适合本地 CLI 使用 |
| pipx | 没有 uv 普及，且 uv 已包含 pipx 的功能 |
