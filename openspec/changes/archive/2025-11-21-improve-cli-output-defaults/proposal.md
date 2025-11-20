# Improve CLI Output Defaults

## Why
当前 CLI 的输出行为需要显式指定 `--print` 才会输出结果到 stdout，这对于管道操作和脚本使用不够友好：

**当前问题**：
1. 用户必须记得加 `--print` 才能看到输出，否则默认静默
2. `--print-entities` 是一个独立的标志，容易忘记或误用
3. 不符合 Unix 管道哲学：默认应该输出到 stdout，需要静默时才加参数
4. 实体信息输出到 stdout 会干扰管道操作

**改进后的行为**：
- 默认输出结果到 stdout（符合 Unix 工具惯例）
- 实体信息输出到 stderr（日志性质，不干扰管道）
- 需要静默时使用 `--no-print` 参数
- 简化用户体验，减少认知负担

## What Changes
- **删除** `anonymize` 命令的 `--print` 和 `--print-entities` 参数
- **删除** `restore` 命令的 `--print` 参数
- **添加** `--no-print` 参数（适用于两个命令）
- **修改** 默认行为：始终输出到 stdout，除非指定 `--no-print` 或 `--output`
- **修改** `anonymize` 命令：默认将实体信息输出到 stderr
- **修改** `--no-print` 参数：同时禁用 stdout 和 stderr 的输出

**使用场景对比**：

| 场景 | 旧命令 | 新命令 |
|------|--------|--------|
| 查看结果 | `inu anonymize -f input.txt --print` | `inu anonymize -f input.txt` |
| 管道使用 | `echo "text" \| inu anonymize --print \| grep ...` | `echo "text" \| inu anonymize \| grep ...` |
| 保存到文件 | `inu anonymize -f input.txt -o output.txt` | `inu anonymize -f input.txt -o output.txt` |
| 查看实体 | `inu anonymize -f input.txt --print --print-entities` | `inu anonymize -f input.txt 2>&1` |
| 只保存不显示 | `inu anonymize -f input.txt -o output.txt` | `inu anonymize -f input.txt -o output.txt --no-print` |

## Impact
- **影响的 specs**: `cli` spec 需要更新输出行为的场景
- **影响的代码**:
  - `cmd/inu/commands/anonymize.go` - 移除 `--print`/`--print-entities`，添加 `--no-print`
  - `cmd/inu/commands/restore.go` - 移除 `--print`，添加 `--no-print`
  - `pkg/cli/output.go` - 修改 `WriteOutput` 函数签名和默认行为
  - 相关测试文件需要更新
- **破坏性变更**: 是
  - 现有使用 `--print` 的脚本需要移除该参数
  - 现有使用 `--print-entities` 的脚本需要改用 stderr 重定向
  - 但新行为更符合 Unix 惯例，迁移成本低
- **向后兼容策略**: 
  - 可以在文档中提供迁移指南
  - 考虑在过渡期保留 `--print` 作为废弃参数（no-op）并发出警告

## Migration Guide
### 用户脚本迁移

**Scenario 1: 打印输出**
```bash
# 旧方式
inu anonymize --file input.txt --print

# 新方式（移除 --print）
inu anonymize --file input.txt
```

**Scenario 2: 查看实体**
```bash
# 旧方式
inu anonymize --file input.txt --print --print-entities

# 新方式（实体自动输出到 stderr）
inu anonymize --file input.txt 2>&1
# 或只看实体
inu anonymize --file input.txt 2>&1 1>/dev/null
```

**Scenario 3: 静默保存**
```bash
# 旧方式（已经静默）
inu anonymize --file input.txt --output result.txt --output-entities entities.yaml

# 新方式（需要明确禁用输出）
inu anonymize --file input.txt --output result.txt --output-entities entities.yaml --no-print
```

**Scenario 4: 管道使用**
```bash
# 旧方式
cat input.txt | inu anonymize --print | inu restore --entities e.yaml --print

# 新方式（更简洁）
cat input.txt | inu anonymize | inu restore --entities e.yaml
```
