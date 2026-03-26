# demon-hunter 工具选型决策

> 决策时间: 2026-03-17
> 适用组件: demon-hunter (P3)

---

## 工具选型

### 核心三件套（全部 Go 单二进制，零环境依赖）

| 工具 | 类别 | 检测范围 | 付费情况 | 选择理由 |
|------|------|---------|---------|---------|
| **bearer** | SAST | 代码安全漏洞（注入/XSS/越权）+ PII 泄漏检测，**数据流分析** | 完全免费开源 | 单 Go 二进制；数据流分析比 semgrep 更深；覆盖 JS/TS/Python/Go/Java/Ruby/PHP |
| **trivy** | SCA + Secrets + IaC | 依赖 CVE（npm/pip/cargo/go/maven...）+ 硬编码密钥（当前文件）+ 容器/IaC 配置错误 | 完全免费，Aqua Security 出品 | SCA 领域覆盖最全的免费工具 |
| **gitleaks** | Secrets | Git history 中的密钥泄漏（trivy 不扫历史）| 完全免费 | 600+ 检测规则，git history 扫描最佳 |

### 为什么用 Bearer 而不是 semgrep

| 维度 | semgrep | bearer |
|------|---------|--------|
| 技术深度 | 规则匹配（Pattern matching）| **数据流分析（Taint analysis）** |
| 跨函数漏洞 | ❌ 无法追踪 | ✅ 追踪污点传播路径 |
| 分发方式 | pip only，**无独立二进制** | **单 Go 二进制** |
| 环境依赖 | 需要 Python 3.9+ | 零依赖 |
| 免费程度 | CLI 免费，Pro 付费 | 完全免费开源 |
| 成熟度 | 14k stars，5年+ | 2.6k stars，2021年起 |

**结论**：semgrep 的最大痛点（无独立二进制 → 环境多样性问题）和技术局限（无数据流分析）促使我们选择 Bearer。三件套全部是 Go 单二进制，从根本上解决了环境多样性问题。

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

**三件套全部是 Go 单二进制**，覆盖所有主流环境：

| 工具 | macOS x86 | macOS arm64 | Linux x86 | Linux arm64 | WSL2 | 无网络 |
|------|:---------:|:-----------:|:---------:|:-----------:|:----:|:------:|
| bearer | ✅ | ✅ | ✅ | ✅ | ✅ | 需预下载 |
| trivy | ✅ | ✅ | ✅ | ✅ | ✅ | 需预下载 |
| gitleaks | ✅ | ✅ | ✅ | ✅ | ✅ | 需预下载 |

**结论**：环境多样性问题通过"全部选 Go 单二进制"从根本上解决。

---

## 隔离执行策略

### 核心原则

**不侵入用户的系统级环境。** 所有工具运行在 skill 自管理的隔离目录内。

### 方案

参考 gstack 的 `browse/dist/browse` 模式——skill 目录内自带二进制。

```
judge-the-code/
└── skills/
    └── demon-hunter/
        ├── SKILL.md
        ├── setup              ← 一次性安装脚本（类似 gstack/setup）
        └── bin/               ← skill 自管理的工具目录，不写入系统 PATH
            ├── trivy          ← Go 单二进制，直接从 GitHub Releases 下载
            ├── gitleaks       ← Go 单二进制，直接从 GitHub Releases 下载
            └── semgrep        ← 见下方说明
```

### 各工具的隔离方式

三件套全部是 **Go 单二进制**，安装方式完全一致：

- 下载到 `SKILL_DIR/bin/`，不写入系统 PATH
- 无任何运行时依赖
- 首次运行时由 `setup` 脚本自动下载，后续直接使用

```bash
# bearer
curl -sSfL "https://github.com/Bearer/bearer/releases/download/v${VER}/bearer_${VER}_${os}_${arch}.tar.gz" \
  | tar -xz -C "$SKILL_DIR/bin" bearer

# trivy
curl -sSfL "https://github.com/aquasecurity/trivy/releases/download/v${VER}/trivy_${VER}_${OS}_${ARCH}.tar.gz" \
  | tar -xz -C "$SKILL_DIR/bin" trivy

# gitleaks
curl -sSfL "https://github.com/gitleaks/gitleaks/releases/download/v${VER}/gitleaks_${VER}_${os}_${arch}.tar.gz" \
  | tar -xz -C "$SKILL_DIR/bin" gitleaks
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

三件套统一逻辑：优先用 `SKILL_DIR/bin/`，其次用系统 PATH，缺失则标记 MISSING。

```bash
#!/usr/bin/env bash
# SKILL_DIR/bin/find-tools
SKILL_DIR="$(cd "$(dirname "$0")/.." && pwd)"
BIN="$SKILL_DIR/bin"

for tool in bearer trivy gitleaks; do
  if [ -x "$BIN/$tool" ]; then
    echo "${tool^^}=$BIN/$tool"
  elif command -v "$tool" &>/dev/null; then
    echo "${tool^^}=$(command -v $tool)"
  else
    echo "${tool^^}=MISSING"
  fi
done
```

### 缺失工具的处理原则

**不因某个工具缺失而整体失败**，而是：

1. 输出哪些维度被跳过及原因
2. 给出安装建议（setup 脚本路径）
3. 已可用的工具继续扫描，输出部分结果

```
🔍 demon-hunter 扫描报告

✅ 代码漏洞扫描 (bearer)      — 完成，发现 2 个高危漏洞
✅ 依赖 CVE 扫描 (trivy)      — 完成，发现 3 个高危漏洞
⚠️ Git 密钥扫描 (gitleaks)   — 已跳过（未安装）
   → 运行 setup 脚本安装：{SKILL_DIR}/setup
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
