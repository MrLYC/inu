# Design: Improve CLI Output Defaults

## Overview
本设计重构 CLI 命令的输出行为，使其更符合 Unix 工具的惯例和管道操作的最佳实践。

## Design Decisions

### 1. 默认输出到 stdout
**决策**: 移除 `--print` 参数，默认将结果输出到 stdout

**理由**:
- **Unix 哲学**: 命令行工具默认应该输出到 stdout，这是 grep、sed、awk 等所有标准工具的行为
- **管道友好**: 用户可以直接将输出传递给其他命令，无需记得加 `--print`
- **认知负担**: 减少用户需要记忆的参数，简化使用体验
- **常见实践**: 大多数 CLI 工具（如 jq, yq, base64）都默认输出结果

**替代方案考虑**:
- 保留 `--print` 作为默认行为：但这违背了 Unix 惯例
- 使用环境变量控制：过于复杂，不够直观

### 2. 实体信息输出到 stderr
**决策**: 默认将实体信息输出到 stderr，而不是 stdout

**理由**:
- **流分离**: stderr 用于日志/诊断信息，stdout 用于主要输出
- **管道安全**: 实体信息不会干扰管道中的数据流
- **标准实践**: 类似于 `curl -v`（详细信息到 stderr）、`git`（状态信息到 stderr）
- **灵活控制**: 用户可以分别重定向 stdout 和 stderr

**行为**:
```bash
# stdout: 脱敏后的文本
# stderr: 实体信息（如果有）
inu anonymize -f input.txt

# 只要文本
inu anonymize -f input.txt 2>/dev/null

# 只要实体信息
inu anonymize -f input.txt 1>/dev/null

# 都要（合并到一起）
inu anonymize -f input.txt 2>&1
```

### 3. 新增 --no-print 参数
**决策**: 添加 `--no-print` 参数来禁用所有输出（stdout 和 stderr）

**理由**:
- **明确意图**: 用户可以明确表示"我只想保存到文件，不要输出"
- **完全静默**: 适合脚本场景，只关心副作用（保存文件）
- **一致性**: 控制所有输出（包括实体信息），避免部分静默的混乱

**使用场景**:
- 批处理脚本：只需要文件输出，不需要终端输出
- 后台任务：避免产生不必要的日志
- 性能优化：跳过格式化和输出操作

### 4. 输出逻辑矩阵

| 场景 | stdout | stderr | 说明 |
|------|--------|--------|------|
| 默认行为 | 主输出 | 实体信息 | 用户友好，信息完整 |
| `--no-print` | 无 | 无 | 完全静默 |
| `--output file` | 主输出 | 实体信息 | 同时输出到文件和终端 |
| `--output + --no-print` | 无 | 无 | 只保存文件 |
| 进度信息 | - | 始终输出 | `ProgressMessage()` 不受 `--no-print` 影响 |

### 5. 进度信息的处理
**决策**: 进度信息（"Initializing LLM...", "Anonymization complete"）继续输出到 stderr，不受 `--no-print` 影响

**理由**:
- 这些是诊断/状态信息，不是主要输出
- 用户通常希望看到进度，即使禁用了主输出
- 如果用户真的想完全静默，可以用 `2>/dev/null`

**替代方案**: 
- 也被 `--no-print` 禁用：但这会让用户失去长时间运行任务的反馈
- 添加 `--quiet` 参数来控制进度信息：增加复杂度，暂不必要

## Implementation Plan

### Phase 1: 修改输出函数
修改 `pkg/cli/output.go`:
- `WriteOutput(content, print, outputFile)` → `WriteOutput(content, noPrint, outputFile)`
- 默认输出到 stdout，除非 `noPrint=true` 或只有 `outputFile`
- 添加 `WriteEntitiesToStderr(entities, noPrint)` 函数

### Phase 2: 更新命令
修改 `cmd/inu/commands/anonymize.go`:
- 移除 `--print` 和 `--print-entities` 标志
- 添加 `--no-print` 标志
- 调用新的输出函数

修改 `cmd/inu/commands/restore.go`:
- 移除 `--print` 标志
- 添加 `--no-print` 标志
- 调用新的输出函数

### Phase 3: 更新测试
- 更新 `output_test.go` 的测试用例
- 确保管道行为正确（通过集成测试）

### Phase 4: 更新文档和规范
- 更新 `openspec/specs/cli/spec.md`
- 更新 `README.md` 的使用示例
- 添加迁移指南到文档

## Backward Compatibility

### 破坏性变更
1. 移除 `--print` 参数
2. 移除 `--print-entities` 参数
3. 默认行为改变（从静默到输出）

### 迁移策略
**Option 1: 直接移除（推荐）**
- 干净利落，符合语义化版本的 major bump
- 在 v1.0.0 之前进行此变更
- 提供清晰的迁移文档

**Option 2: 废弃警告**
- 保留 `--print` 作为 no-op（不做任何事）
- 发出废弃警告："--print is deprecated and has no effect (output is now enabled by default)"
- 在下一个 major 版本移除
- 缺点：增加代码复杂度，延迟清理时间

**推荐**: 采用 Option 1，因为项目还处于早期阶段（v0.x.x），用户基数小，迁移成本低。

## Testing Strategy

### Unit Tests
- `TestWriteOutput_DefaultBehavior` - 测试默认输出到 stdout
- `TestWriteOutput_NoPrint` - 测试 `--no-print` 禁用输出
- `TestWriteOutput_FileOnly` - 测试只输出到文件
- `TestWriteEntitiesToStderr_Default` - 测试实体输出到 stderr
- `TestWriteEntitiesToStderr_NoPrint` - 测试 `--no-print` 禁用实体输出

### Integration Tests
```bash
# 测试默认输出
output=$(echo "张三" | ./inu anonymize)
[[ -n "$output" ]] || exit 1

# 测试 --no-print
output=$(echo "张三" | ./inu anonymize --no-print)
[[ -z "$output" ]] || exit 1

# 测试实体输出到 stderr
stderr=$(echo "张三" | ./inu anonymize 2>&1 1>/dev/null)
[[ "$stderr" =~ "个人信息" ]] || exit 1

# 测试管道
output=$(echo "张三的电话是 13800138000" | ./inu anonymize -e /tmp/e.yaml | ./inu restore -e /tmp/e.yaml)
[[ "$output" == "张三的电话是 13800138000" ]] || exit 1
```

## Risks and Mitigations

### Risk 1: 现有用户脚本中断
**影响**: 高  
**概率**: 中  
**缓解**: 
- 提供详细的迁移指南和示例
- 在 Release Notes 中突出显示
- 考虑在项目达到 v1.0.0 之前进行此变更

### Risk 2: 用户不想看到默认输出
**影响**: 低  
**概率**: 低  
**缓解**: 
- 提供 `--no-print` 参数
- 文档中说明如何重定向到 `/dev/null`

### Risk 3: 实体信息输出到 stderr 被忽略
**影响**: 低  
**概率**: 中  
**缓解**: 
- 在文档中清楚说明 stderr 的用途
- 提供重定向示例（`2>&1`）
