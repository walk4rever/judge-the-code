# Agent 2 — Secrets & 依赖漏洞解读（trivy + gitleaks）

> 归属：`demon-hunter` skill，Phase 2 并行分析之一。
> 输入：trivy JSON + gitleaks JSON 扫描结果

## 任务

解读依赖漏洞和密钥泄漏，区分真实风险和噪音。

## Secrets 分析（trivy secrets + gitleaks）

1. 合并两个工具的 secrets findings，去重（同一文件+行号）

2. 对每条 finding 判断**真实性**：
   ```bash
   # 读取泄漏位置附近代码确认
   sed -n 'X,Yp' <file>
   ```
   - 若是 `.env.example` / `test/fixture` / 注释中的示例 → 排除，注明原因
   - 若是真实密钥格式（有熵值、特定前缀如 `sk-`/`ghp_`/`AKIA`）→ 保留为高危

3. 区分**当前代码**（trivy）vs **git history**（gitleaks）：
   - git history 中的泄漏即使已删除也需要处理（需 rotate key）

## 依赖漏洞分析（trivy SCA）

1. 只关注 **CRITICAL / HIGH** 级别 CVE

2. 对每个漏洞补充**可利用性判断**：
   - 该依赖是 direct dependency 还是 transitive？
   - 项目是否实际使用了漏洞涉及的功能？
   - 是否有已发布的 fix 版本？

3. 识别**批量可升级**的情况（多个 CVE 可通过升级同一个包解决）

## 输出格式

```
## 🔑 密钥泄漏

### 当前代码中
#### AWS Access Key — `config/deploy.ts:23`
- **类型**：AWS_ACCESS_KEY_ID（AKIA 前缀确认）
- **状态**：当前文件存在，需立即 rotate
- **发现工具**：trivy + gitleaks

### Git History 中（已从代码删除，但仍需 rotate）
- `src/config.ts` — commit a3f7c21（2025-11-03）泄漏 Stripe Secret Key

## 📦 依赖漏洞（CVE）

### Critical（需立即处理）
#### CVE-2024-XXXX — `lodash@4.17.15`
- **影响**：原型链污染，可导致 RCE
- **是否可利用**：项目使用了 `_.merge()`，属于漏洞路径 ⚠️
- **修复**：升级到 lodash@4.17.21（已有 patch）

### High
[同上格式]

### 批量升级建议
- 升级 `express` 4.18 → 4.21 可修复 3 个 High CVE
```

## 完成

输出：
```
✅ Agent 2 完成 — 发现 [N] 个密钥风险，[M] 个高危 CVE
```
