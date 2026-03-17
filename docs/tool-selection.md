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

### 执行入口

```bash
# SKILL_DIR/bin/find-tools（参考 gstack/browse/bin/find-browse）
SKILL_DIR=$(find ~/.agents/skills ~/.claude/skills \
  -name "SKILL.md" -path "*/demon-hunter/*" 2>/dev/null \
  | head -1 | xargs dirname)

TRIVY="$SKILL_DIR/bin/trivy"
GITLEAKS="$SKILL_DIR/bin/gitleaks"

# semgrep 按优先级
if command -v uvx &>/dev/null; then
  SEMGREP="uvx semgrep"
elif [ -x "$SKILL_DIR/.venv/bin/semgrep" ]; then
  SEMGREP="$SKILL_DIR/.venv/bin/semgrep"
else
  SEMGREP=""  # 触发安装提示
fi
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
