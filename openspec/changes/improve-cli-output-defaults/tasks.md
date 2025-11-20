# Implementation Tasks

## 1. 修改输出函数
- [ ] 1.1 重构 `pkg/cli/output.go`
  - [ ] 修改 `WriteOutput` 函数签名：`(content, print, outputFile)` → `(content, noPrint, outputFile)`
  - [ ] 默认输出到 stdout（当 `noPrint=false` 且 `outputFile=""` 时）
  - [ ] 当 `noPrint=true` 时禁用 stdout 输出
  - [ ] 当 `outputFile` 存在时写入文件（不管 `noPrint` 值）
  - [ ] 添加 `WriteEntitiesToStderr(entities, noPrint)` 函数
  - [ ] 实体信息默认输出到 stderr，除非 `noPrint=true`
- [ ] 1.2 更新 `output_test.go` 测试
  - [ ] 修改 `TestWriteOutput_PrintOnly` → `TestWriteOutput_DefaultBehavior`
  - [ ] 添加 `TestWriteOutput_NoPrint` 测试
  - [ ] 修改其他相关测试用例
  - [ ] 添加 `TestWriteEntitiesToStderr_Default` 测试
  - [ ] 添加 `TestWriteEntitiesToStderr_NoPrint` 测试

## 2. 更新 anonymize 命令
- [ ] 2.1 修改 `cmd/inu/commands/anonymize.go`
  - [ ] 移除 `anonymizePrint` 变量
  - [ ] 移除 `anonymizePrintEntities` 变量
  - [ ] 添加 `anonymizeNoPrint` 变量
  - [ ] 移除 `--print` 标志（`flags.BoolVarP(&anonymizePrint, "print", "p", ...)`）
  - [ ] 移除 `--print-entities` 标志
  - [ ] 添加 `--no-print` 标志
  - [ ] 修改 `runAnonymize` 调用 `WriteOutput(result, anonymizeNoPrint, anonymizeOutput)`
  - [ ] 修改输出实体信息的逻辑，调用 `WriteEntitiesToStderr(entities, anonymizeNoPrint)`
  - [ ] 移除 `if anonymizePrintEntities` 条件判断

## 3. 更新 restore 命令
- [ ] 3.1 修改 `cmd/inu/commands/restore.go`
  - [ ] 移除 `restorePrint` 变量
  - [ ] 添加 `restoreNoPrint` 变量
  - [ ] 移除 `--print` 标志
  - [ ] 添加 `--no-print` 标志
  - [ ] 修改 `runRestore` 调用 `WriteOutput(result, restoreNoPrint, restoreOutput)`

## 4. 更新规范文档
- [ ] 4.1 更新 `openspec/specs/cli/spec.md`
  - [ ] 复制 `changes/improve-cli-output-defaults/specs/cli/spec.md` 的内容
  - [ ] 确保 MODIFIED 和 REMOVED 部分正确标注

## 5. 更新用户文档
- [ ] 5.1 更新 `README.md`
  - [ ] 修改 "命令行使用" 部分的示例，移除 `--print`
  - [ ] 添加 `--no-print` 使用示例
  - [ ] 添加实体信息重定向示例（`2>&1`, `2>/dev/null` 等）
  - [ ] 添加管道使用示例
- [ ] 5.2 创建迁移指南
  - [ ] 在 README 或单独文档中添加 "从 v0.x 迁移" 章节
  - [ ] 说明参数变更和行为变更
  - [ ] 提供旧命令到新命令的对照表

## 6. 集成测试
- [ ] 6.1 添加集成测试脚本（可选）
  - [ ] 测试默认输出行为
  - [ ] 测试 `--no-print` 行为
  - [ ] 测试实体信息输出到 stderr
  - [ ] 测试管道操作
  - [ ] 测试重定向操作
- [ ] 6.2 手动验证
  - [ ] 验证 `inu anonymize` 默认输出
  - [ ] 验证 `inu restore` 默认输出
  - [ ] 验证管道：`echo "text" | inu anonymize | inu restore -e e.yaml`
  - [ ] 验证 stderr 输出：`inu anonymize 2>&1 | grep "个人信息"`
  - [ ] 验证 `--no-print`：`inu anonymize --no-print` 无输出
  - [ ] 验证 `--output + --no-print`：只写文件不输出

## 7. 验证和清理
- [ ] 7.1 运行所有测试：`go test ./...`
- [ ] 7.2 验证构建：`go build -o bin/inu ./cmd/inu`
- [ ] 7.3 检查帮助信息：`inu anonymize --help` 和 `inu restore --help`
- [ ] 7.4 搜索代码中是否还有遗漏的 `--print` 引用
- [ ] 7.5 更新 CHANGELOG（如果有）

## 任务依赖关系

```
1.1 (修改输出函数)
  ↓
1.2 (更新输出函数测试)
  ↓
2.1 + 3.1 (并行更新两个命令)
  ↓
6.1 (集成测试)
  ↓
6.2 (手动验证)
  ↓
4.1 + 5.1 + 5.2 (并行更新文档)
  ↓
7.1-7.5 (最终验证)
```

## 预估工作量
- **核心实现**: 2-3 小时（步骤 1-3）
- **测试和验证**: 1-2 小时（步骤 6-7）
- **文档更新**: 1 小时（步骤 4-5）
- **总计**: 4-6 小时

## 风险和注意事项
- ⚠️ 破坏性变更：确保在发布时更新版本号（major bump）
- ⚠️ 测试覆盖：特别注意 stderr 输出的测试（需要捕获 stderr）
- ⚠️ 文档完整性：提供清晰的迁移指南，避免用户困惑
- ⚠️ 向后兼容：考虑是否需要在过渡期保留废弃警告
