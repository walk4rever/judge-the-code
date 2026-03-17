# Agent 2 — 错误处理取向分析师

> 归属：`design-lens` skill，Phase 1 并行分析之一。

## 任务

通过错误处理方式，揭示代码库对"意外"的态度。

## 步骤

1. **提取错误处理代码**（带上下文，单次命令覆盖全库）：

   ```bash
   # 主要错误处理模式（含前1行后3行上下文）
   grep -rn "catch\|throw new\|panic(\|except \|raise \|if err " \
     src/ --include="*.ts" --include="*.go" --include="*.py" \
     -A 3 -B 1 \
     | grep -v "node_modules\|\.test\.\|\.spec\.\|dist/" | head -80
   ```

   ```bash
   # 专项：静默忽略检测（高危模式）
   grep -rn "catch\s*[({][^}]*}\s*$\|catch\s*(.*)\s*{}" \
     src/ --include="*.ts" --include="*.js" \
     | grep -v "node_modules\|\.test\." | head -20
   ```

   ```bash
   # 专项：Go 风格错误忽略
   grep -rn "_ = err\|_ , err\b" \
     . --include="*.go" | grep -v "_test.go" | head -10
   ```

2. 分析**错误处理风格**：

   | 风格 | 特征 | 透露的取向 |
   |------|------|-----------|
   | 快速失败 | `panic`, `assert`, `throw` 用于不可恢复错误 | 防御性强，宁可崩溃也不带病运行 |
   | 防御包裹 | 大量 `try/catch`，返回默认值 | 容错优先，可能掩盖问题 |
   | 显式传播 | `Result<T, E>`, `if err != nil` | 强制调用方处理，类型安全 |
   | 静默忽略 | `catch (e) {}`, `_ = err` | 可能的隐患，也可能是有意为之 |
   | 统一处理 | 全局 error handler，middleware 层统一捕获 | 架构意识强 |

3. **统计错误处理密度**（用数字说话）：

   ```bash
   # 统计 try/catch 总数
   grep -rn "\btry\b\|\bcatch\b" \
     src/ --include="*.ts" --include="*.js" \
     | grep -v "node_modules\|\.test\." | wc -l

   # 统计 throw 总数
   grep -rn "\bthrow\b" \
     src/ --include="*.ts" --include="*.js" \
     | grep -v "node_modules\|\.test\." | wc -l
   ```

4. 检查**错误信息质量**：

   ```bash
   # 错误消息是否携带上下文（包含冒号分隔的描述）
   grep -rn "throw new.*\".*:.*\"\|new Error(.*:.*)" \
     src/ --include="*.ts" | grep -v "node_modules\|\.test\." | head -10

   # 是否有自定义错误类型
   grep -rn "extends Error\|class.*Error\b" \
     src/ --include="*.ts" | grep -v "node_modules\|\.test\." | head -10
   ```

5. 检查**边界保护意识**：

   ```bash
   # 输入校验层
   grep -rn "\.parse(\|\.safeParse(\|\.validate(\|z\.\|joi\.\|yup\." \
     src/ --include="*.ts" | grep -v "node_modules\|\.test\." | head -20

   # 空值保护
   grep -rn "??\|!\.\|null check\|undefined check\|\bif.*null\b\|\bif.*undefined\b" \
     src/ --include="*.ts" | grep -v "node_modules\|\.test\." | wc -l
   ```

## 输出格式

```
## 错误处理取向

### 主导风格
[快速失败 / 防御包裹 / 显式传播 / 混合]
**证据**：
- `src/services/UserService.ts:45` — throw 用于不合法的业务状态（快速失败）
- `src/api/handler.ts:12` — 全局 catch 统一处理（架构意识）
**结论**：...

### 错误处理密度
- try/catch 块：[N] 处
- throw 语句：[M] 处
- 静默忽略（catch {}）：[K] 处 → [评价]

### 错误信息质量
[信息丰富 / 一般 / 过于简略]
**证据**：...

### 边界保护
[严格 / 一般 / 宽松]
**证据**：
- 使用 zod 在 API 层统一校验 → 边界清晰
- service 层未做二次校验 → 信任上游，有一定风险
**结论**：...

### 值得关注的错误处理决策
- [具体例子 + 分析]
```

## 完成

输出：
```
✅ Agent 2/4 完成 — 错误处理取向已分析
```
