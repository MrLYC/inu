# 实施任务清单

## Phase 1: 核心接口修改
- [ ] **1.1** 修改 `pkg/anonymizer/anonymizer.go` 中的接口定义
  - [ ] 删除 `AnonymizeText()` 方法
  - [ ] 重命名 `AnonymizeTextStream()` → `Anonymize()`
  - [ ] 修改 `RestoreText()` 签名: 添加 `writer io.Writer` 参数,返回 `([]string, error)`

- [ ] **1.2** 实现新的 `RestoreText()` 逻辑
  - [ ] 使用 `FindAllStringIndex` 替代 `ReplaceAllStringFunc`
  - [ ] 实现流式写入到 writer
  - [ ] 收集未还原的占位符(去重)
  - [ ] 返回未还原占位符列表

- [ ] **1.3** 更新 `pkg/anonymizer/anonymizer_test.go`
  - [ ] 更新 `TestRestoreText` 系列测试
    - [ ] 测试全部还原成功的情况(返回空列表)
    - [ ] 测试部分还原失败的情况(返回失败列表)
    - [ ] 测试无实体的情况(保留所有占位符)
  - [ ] 更新所有使用 `AnonymizeText` 的测试 → `Anonymize`
  - [ ] 更新所有使用 `AnonymizeTextStream` 的测试 → `Anonymize`
  - [ ] 验证测试覆盖率不降低

## Phase 2: CLI 命令适配

- [ ] **2.1** 修改 `cmd/inu/commands/anonymize.go`
  - [ ] 将 `AnonymizeTextStream` 调用改为 `Anonymize`
  - [ ] 验证行为一致性

- [ ] **2.2** 修改 `cmd/inu/commands/restore.go`
  - [ ] 修改 `RestoreText` 调用:添加 writer 参数
  - [ ] 处理返回的未还原占位符列表
  - [ ] 实现警告输出逻辑:
    ```go
    if len(unrestoredPlaceholders) > 0 {
        cli.WarningMessage("Warning: %d placeholder(s) could not be restored:", len(unrestoredPlaceholders))
        for _, placeholder := range unrestoredPlaceholders {
            cli.WarningMessage("  - %s", placeholder)
        }
    }
    ```
  - [ ] 更新 `getInput()` 方法以返回 writer

- [ ] **2.3** 修改 `cmd/inu/commands/interactive.go`
  - [ ] 将 `AnonymizeTextStream` 调用改为 `Anonymize`
  - [ ] 修改 `RestoreText` 调用:添加 writer 参数
  - [ ] 实现警告输出逻辑(同 2.2)
  - [ ] 移除旧的错误处理逻辑 (lines 131-134)

- [ ] **2.4** 更新 CLI 命令测试
  - [ ] 更新 `cmd/inu/commands/interactive_test.go`
  - [ ] 添加测试:验证警告输出格式
  - [ ] 添加测试:验证部分还原失败的场景
  - [ ] 验证所有集成测试通过

## Phase 3: Web API 适配

- [ ] **3.1** 修改 `pkg/web/handlers/anonymize.go`
  - [ ] 更新接口定义: `AnonymizeText` → `Anonymize`
  - [ ] 使用 `bytes.Buffer` 包装调用:
    ```go
    var buf bytes.Buffer
    entities, err := anon.Anonymize(c.Request.Context(), entityTypes, req.Text, &buf)
    anonymizedText := buf.String()
    ```
  - [ ] 验证 HTTP 响应格式不变

- [ ] **3.2** 修改 `pkg/web/handlers/restore.go`
  - [ ] 使用 `bytes.Buffer` 包装调用:
    ```go
    var buf bytes.Buffer
    unrestored, err := a.restorer.RestoreText(c.Request.Context(), entities, text, &buf)
    ```
  - [ ] 扩展 `RestoreResponse` 结构:
    ```go
    type RestoreResponse struct {
        RestoredText          string   `json:"restored_text"`
        UnrestoredPlaceholders []string `json:"unrestored_placeholders,omitempty"`
    }
    ```
  - [ ] 在响应中包含未还原占位符(如果非空)

- [ ] **3.3** 更新 Web handlers 测试
  - [ ] 更新 `pkg/web/handlers/anonymize_test.go`
  - [ ] 更新 `pkg/web/handlers/restore_test.go`
    - [ ] 测试 `unrestored_placeholders` 字段存在
    - [ ] 测试全部还原时该字段省略
    - [ ] 验证 HTTP 响应格式向后兼容
  - [ ] 更新 mock 定义 (`pkg/web/handlers/mock_test.go`)

## Phase 4: Spec 更新

- [ ] **4.1** 创建 CLI spec delta
  - [ ] 创建 `openspec/changes/unify-streaming-interface/specs/cli/spec.md`
  - [ ] 标记 MODIFIED: Anonymize 命令(方法重命名)
  - [ ] 标记 MODIFIED: Restore 命令(添加警告输出行为)
  - [ ] 标记 MODIFIED: Interactive 命令(添加警告输出行为)

- [ ] **4.2** (可选)创建 Web API spec delta
  - [ ] 创建 `openspec/changes/unify-streaming-interface/specs/web-api/spec.md`
  - [ ] 标记 MODIFIED: POST /restore 响应格式(添加 unrestored_placeholders 字段)

## Phase 5: 文档和验证

- [ ] **5.1** 更新项目文档
  - [ ] 检查 README.md 是否有接口使用示例(如果有则更新)
  - [ ] 检查是否有 API 文档需要更新

- [ ] **5.2** 运行完整测试套件
  - [ ] 运行单元测试: `go test ./pkg/...`
  - [ ] 运行集成测试: `go test ./cmd/...`
  - [ ] 验证测试覆盖率: `go test -cover ./...`

- [ ] **5.3** 手动测试
  - [ ] 测试 `inu anonymize` 命令(验证方法重命名正常工作)
  - [ ] 测试 `inu restore` 命令:
    - [ ] 全部还原成功的场景(无警告)
    - [ ] 部分还原失败的场景(显示警告和占位符列表)
  - [ ] 测试 `inu interactive` 命令:
    - [ ] 脱敏流程正常工作
    - [ ] 还原流程显示警告(如果有失败)
  - [ ] 测试 Web API:
    - [ ] POST /anonymize 正常工作
    - [ ] POST /restore 返回正确格式(包含 unrestored_placeholders)

- [ ] **5.4** 验证 OpenSpec 合规性
  - [ ] 运行 `openspec validate unify-streaming-interface --strict`
  - [ ] 修复所有验证错误
  - [ ] 确保所有 spec deltas 正确引用现有 requirements

## Phase 6: 代码审查和提交

- [ ] **6.1** 自我代码审查
  - [ ] 检查所有变更符合 Go 编码规范
  - [ ] 验证错误处理完整
  - [ ] 验证日志输出格式统一
  - [ ] 检查是否有遗漏的 TODO 注释

- [ ] **6.2** Git 提交策略
  - [ ] Commit 1: 核心接口修改(pkg/anonymizer)
  - [ ] Commit 2: CLI 命令适配(cmd/inu/commands)
  - [ ] Commit 3: Web API 适配(pkg/web/handlers)
  - [ ] Commit 4: Spec 更新和文档
  - [ ] 每个 commit 包含对应的测试更新

- [ ] **6.3** 最终验证
  - [ ] 所有测试通过 ✅
  - [ ] OpenSpec 验证通过 ✅
  - [ ] 手动测试通过 ✅
  - [ ] 代码审查通过 ✅

## Dependencies

**任务依赖关系**:
- Phase 2 依赖 Phase 1 完成(CLI 命令依赖核心接口)
- Phase 3 依赖 Phase 1 完成(Web API 依赖核心接口)
- Phase 4 依赖 Phase 2 和 Phase 3 完成(spec deltas 需要反映实现细节)
- Phase 5 依赖 Phase 1-4 完成(验证所有变更)
- Phase 6 依赖 Phase 5 完成(确保质量后提交)

**外部依赖**:
- 无新的外部库依赖
- Go 版本要求不变
- 测试工具不变

## Risks and Mitigation

**风险 1: RestoreText 性能下降**
- **缓解**: 使用 benchmark 测试验证性能
- **预期**: 流式写入性能应该与字符串拼接持平或更好

**风险 2: 测试覆盖不全**
- **缓解**:
  - 运行 `go test -cover` 确保覆盖率不降低
  - 添加边界情况测试(空输入、超大输入等)

**风险 3: Web API 响应格式破坏向后兼容**
- **缓解**:
  - 使用 `omitempty` 标签使新字段可选
  - 添加集成测试验证旧客户端兼容性

**风险 4: CLI 警告输出影响脚本解析**
- **缓解**:
  - 警告输出到 stderr(不影响 stdout)
  - 保持主输出格式不变
  - 文档说明如何抑制警告(`2>/dev/null`)

## Estimated Timeline

- **Phase 1**: 2-3 hours (核心逻辑修改)
- **Phase 2**: 2-3 hours (CLI 命令适配)
- **Phase 3**: 1-2 hours (Web API 适配)
- **Phase 4**: 1 hour (Spec 更新)
- **Phase 5**: 1-2 hours (验证和文档)
- **Phase 6**: 1 hour (代码审查和提交)

**总计**: 8-12 hours

## Success Criteria

**必须满足**:
- [ ] 所有单元测试通过
- [ ] 所有集成测试通过
- [ ] OpenSpec 验证通过(strict 模式)
- [ ] CLI 命令行为符合 spec 描述
- [ ] Web API 响应向后兼容

**可选但推荐**:
- [ ] 测试覆盖率 ≥ 80%
- [ ] Benchmark 性能与之前持平
- [ ] 所有 edge cases 有对应测试
- [ ] 代码审查无 critical 问题

## Notes

**重要提醒**:
1. 先运行测试再提交代码
2. 每个 phase 完成后运行相关测试确保未破坏现有功能
3. Spec deltas 需要准确引用现有 requirements(使用 line numbers)
4. 警告输出格式需要与现有的 `cli.WarningMessage` 风格一致
5. Web API 新字段使用 `omitempty` 确保向后兼容

**可选增强**(不在本次变更范围):
- 提供 `--strict` 模式(有未还原占位符时返回错误)
- Interactive 命令提供交互式修复失败占位符的功能
- 区分"占位符不存在"和"实体没有值"两种失败原因
- 提供颜色高亮显示未还原占位符(CLI)
