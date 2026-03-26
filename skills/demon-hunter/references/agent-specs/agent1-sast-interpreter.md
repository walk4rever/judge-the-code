# Agent 1 — SAST 结果解读（bearer）

> 归属：`demon-hunter` skill，Phase 2 并行分析之一。
> 输入：bearer JSON 扫描结果（已过滤 critical/high）

## 任务

解读 bearer 的代码安全扫描结果，结合项目上下文判断真实风险。

## 步骤

1. 读取传入的 bearer JSON 片段，提取每条 finding 的：
   - `rule_id` / `title` — 漏洞类型
   - `filename` + `line_number` — 位置
   - `description` — 技术描述
   - `severity` — critical/high/medium

2. 对每条 finding 做**上下文判断**（不要照单全收）：
   - 读取 `filename:line_number` 附近 ±10 行（用 `sed -n 'X,Yp' file`）
   - 判断是否为**真实威胁**还是**误报**（测试文件/mock 数据/注释中的示例）
   - 若是测试文件中的 finding → 降级为低优先，注明"测试代码"

3. 对真实 finding 补充**业务影响说明**：
   - 这个漏洞在什么条件下会被触发？
   - 影响范围：单个端点 / 整个认证层 / 数据存储？
   - 修复难度：一行 / 架构改动？

4. 按严重程度分组输出（critical → high → medium）

## 输出格式

```
## 🔴 代码安全漏洞（SAST）

### Critical

#### [漏洞类型] — `src/api/users.ts:45`
- **是什么**：SQL 拼接注入，用户输入直接进入查询字符串
- **触发条件**：任何调用 `/api/users?filter=` 的请求
- **影响范围**：整个用户数据表可被读取/篡改
- **修复方向**：使用参数化查询或 ORM 的 .where() 方法
- **证据**：`bearer rule: sql-injection`

### High
[同上格式]

### 误报（已排除）
- `tests/mock.ts:12` — 测试文件中的示例代码，非生产路径
```

## 完成

输出：
```
✅ Agent 1 完成 — SAST 分析完成，发现 [N] 个真实漏洞
```
