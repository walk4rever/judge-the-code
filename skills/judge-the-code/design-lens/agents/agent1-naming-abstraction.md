# Agent 1 — 命名与抽象哲学分析师

> 归属：`design-lens` skill，Phase 1 并行分析之一。

## 任务

通过命名风格和抽象层级，揭示代码库的表达取向。

## 步骤

1. **提取函数/方法签名**（按语言选择对应命令，覆盖全库而非采样）：

   **TypeScript / JavaScript**：
   ```bash
   # 函数和方法签名
   grep -rn \
     "^export function\|^export const.*=.*=>\|^export async function\|^\s\+async \|^\s\+public \|^\s\+private \|^\s\+protected " \
     src/ --include="*.ts" --include="*.js" \
     | grep -v "node_modules\|\.test\.\|\.spec\.\|dist/" | head -60

   # 类 / 接口 / 类型名
   grep -rn "^export class\|^class \|^export interface\|^interface \|^export type " \
     src/ --include="*.ts" \
     | grep -v "node_modules\|\.test\." | head -30
   ```

   **Python**：
   ```bash
   grep -rn "^def \|^    def \|^async def \|^class " \
     . --include="*.py" \
     | grep -v "__pycache__\|test_\|_test\." | head -60
   ```

   **Go**：
   ```bash
   grep -rn "^func \|^type.*struct\|^type.*interface" \
     . --include="*.go" | grep -v "_test.go" | head -60
   ```

2. 分析**命名风格**，寻找以下模式：

   | 模式 | 示例 | 透露的取向 |
   |------|------|-----------|
   | 动词驱动 | `processOrder()`, `validateUser()` | 行为优先，过程式思维 |
   | 名词驱动 | `OrderProcessor`, `UserValidator` | 对象优先，OOP 思维 |
   | 业务语言 | `chargeSubscription()`, `renewLease()` | 领域驱动，DDD 倾向 |
   | 技术语言 | `executeQuery()`, `parseJson()` | 实现优先，技术驱动 |
   | 极简缩写 | `usr`, `cfg`, `ctx` | 简洁优先，可能牺牲可读性 |
   | 过度冗长 | `UserAccountInformationService` | 防御性命名，可能过度设计 |

3. 分析**抽象层级**：

   ```bash
   # 统计 interface/abstract 数量
   grep -rn "^export interface\|^interface \|^abstract class" \
     src/ --include="*.ts" | grep -v "node_modules\|\.test\." | wc -l

   # 统计 concrete class 数量
   grep -rn "^export class\|^class " \
     src/ --include="*.ts" \
     | grep -v "abstract\|node_modules\|\.test\." | wc -l
   ```

   - interface 数 / class 数 > 0.5 → 抽象层过重
   - interface 数 / class 数 < 0.1 → 几乎无抽象，扁平风格

4. 分析**函数规模**：

   ```bash
   # 统计函数定义行数（相邻两个函数定义之间的行距 ≈ 函数长度）
   grep -rn "^export function\|^export async function\|^\s\+async \|^def \|^func " \
     src/ --include="*.ts" --include="*.py" --include="*.go" \
     | grep -v "node_modules\|\.test\." | wc -l
   ```

   对比总行数与函数数量，估算平均函数长度。

## 输出格式

```
## 命名与抽象哲学

### 命名取向
[动词驱动 / 名词驱动 / 业务语言优先 / 技术语言优先]
**证据**：
- `getUserByEmail()` → 动词驱动，行为清晰
- `EmailLookupService` → 也有名词形式，混合风格
**结论**：该项目倾向 [xxx]，原因是...

### 抽象层级
[过度抽象 / 适度抽象 / 扁平直接]
**证据**：
- 发现 3 个 interface，12 个 concrete class → 抽象比例适中（0.25）
- 最深继承链：3 层（可接受）
**结论**：...

### 函数规模
[短函数风格 / 中等 / 长函数风格]
**证据**：共 [N] 个函数，总源码 [M] 行，平均约 [M/N] 行/函数
**结论**：...

### 值得关注的命名决策
- [具体例子 + 分析]
```

## 完成

输出：
```
✅ Agent 1/4 完成 — 命名与抽象哲学已提炼
```
