# Agent 3 — 设计层恶魔分析（Claude 语义分析）

> 归属：`demon-hunter` skill，Phase 2 并行分析之一。
> 工具扫不到的设计层面问题，由 Claude 负责。

## 任务

发现工具无法检测的隐性风险：性能隐患、隐性耦合、陷阱 API、债务炸弹。

## 步骤

### 1. 读取已有上下文

优先复用已有分析，避免重复工作：
```bash
# 检查是否有 design-lens 的输出
ls .judge-the-code/philosophy.md 2>/dev/null && echo "FOUND" || echo "MISSING"
ls .judge-the-code/understanding.md 2>/dev/null && echo "FOUND" || echo "MISSING"
```
- 若有 `.judge-the-code/philosophy.md`：读取"隐含设计原则"和"存疑决策"，作为判断基准
- 若有 `.judge-the-code/understanding.md`：读取架构图和核心模块列表

### 2. 性能时间炸弹检测

```bash
# N+1 查询模式：循环内有数据库调用
grep -rn "for\|forEach\|\.map\|while" \
  src/ --include="*.ts" --include="*.py" --include="*.go" \
  -A 5 | grep -B 3 "findOne\|findById\|query\|SELECT\|\.find(" \
  | grep -v "node_modules\|\.test\." | head -30

# 无超时的外部调用
grep -rn "fetch(\|axios\.\|http\.get\|requests\.get" \
  src/ --include="*.ts" --include="*.py" \
  | grep -v "timeout\|AbortSignal\|node_modules\|\.test\." | head -20

# 大对象序列化（潜在内存炸弹）
grep -rn "JSON\.stringify\|json\.dumps\|json\.marshal" \
  src/ --include="*.ts" --include="*.py" --include="*.go" \
  | grep -v "node_modules\|\.test\." | head -20
```

### 3. 隐性耦合检测

```bash
# 全局状态依赖
grep -rn "^let \|^var \|global\.\|window\.\|process\.env" \
  src/ --include="*.ts" --include="*.js" \
  | grep -v "node_modules\|\.test\.\|\.d\.ts" | head -20

# 直接读取 process.env（分散而非集中配置）
grep -rn "process\.env\." \
  src/ --include="*.ts" --include="*.js" \
  | grep -v "node_modules\|\.test\.\|config\." \
  | wc -l
# 若 > 10 处分散读取 → 配置管理混乱隐患
```

### 4. 哲学破坏检测（依赖 .judge-the-code/philosophy.md）

若 `.judge-the-code/philosophy.md` 存在：
- 找出"隐含设计原则"中的核心规则
- 用 grep 验证是否存在违反这些原则的代码
- 示例：原则是"边界处统一校验"，则检查是否有绕过校验层直接操作数据的路径

若不存在：跳过此步，注明"建议先运行 /design-lens 获取更深层分析"

### 5. 陷阱 API 检测

```bash
# 容易误用的异步操作（忘记 await）
grep -rn "async function\|async (" \
  src/ --include="*.ts" -A 10 \
  | grep -v "await\|return\|node_modules\|\.test\." \
  | grep "^\s*[a-zA-Z].*(" | head -20

# 静默失败的 Promise（未处理的 rejection）
grep -rn "\.then(" \
  src/ --include="*.ts" --include="*.js" \
  | grep -v "\.catch\|node_modules\|\.test\." | head -15
```

## 输出格式

```
## ⚡ 性能时间炸弹

#### N+1 查询隐患 — `src/services/VideoService.ts:145`
- **是什么**：forEach 循环内调用 findById，N 个视频 = N 次数据库查询
- **触发条件**：列表页加载，视频数量越多越慢
- **影响**：100 个视频 → 100+ 次 DB 查询，高并发时必崩
- **修复方向**：使用 findByIds([...]) 批量查询

#### 无超时的外部 API 调用 — `src/lib/ai-client.ts:89`
- **是什么**：fetch() 无 timeout，第三方 API 挂起会导致请求永久 hang
- **修复方向**：添加 AbortSignal.timeout(30000)

## 🔗 隐性耦合

#### 配置分散读取 — 全库 23 处直接读 process.env
- **是什么**：配置散落各处，无统一入口
- **风险**：环境变量命名不一致，缺少校验，上线前难以检查完整性
- **修复方向**：集中到 src/config.ts，启动时统一校验

## 💣 哲学破坏（基于 .judge-the-code/philosophy.md）

#### 违反"边界处统一校验"原则 — `src/api/notes.ts:67`
- **是什么**：直接接受用户输入写入数据库，跳过了 validation 层
- **项目原则**：其他所有 API 都经过 zod schema 校验
- **风险**：输入不一致可能导致数据污染

## 🪤 陷阱 API

#### Promise 静默失败 — `src/lib/translation.ts:34`
- **是什么**：.then() 链无 .catch()，失败时静默忽略
- **触发条件**：翻译 API 超时或报错
- **表现**：用户看到空白内容，无错误提示，难以排查
```

## 完成

输出：
```
✅ Agent 3 完成 — 发现 [N] 个设计层隐患
```
