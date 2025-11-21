# Add Anonymizer Tests

## Why
`pkg/anonymizer` 包是 Inu 的核心业务逻辑，负责文本脱敏和还原功能，但目前完全没有单元测试覆盖。这导致：
- 代码变更时无法快速验证核心逻辑的正确性
- 重构和优化存在较高风险
- CI/CD 中无法捕获脱敏逻辑的回归问题
- 新贡献者难以理解预期行为

由于脱敏功能依赖外部 LLM API，测试需要通过 mock 避免实际 API 调用，确保：
- 测试可以在无网络/无 API key 的环境中运行
- 测试执行速度快且稳定
- 不产生外部 API 调用成本

## What Changes
- 为 `pkg/anonymizer` 包添加完整的单元测试覆盖
- 实现 LLM 接口的 mock，避免真实 API 调用
- 测试核心方法：`AnonymizeText`、`RestoreText`、`New`
- 测试边界情况：空输入、错误格式、无实体匹配等
- 使用 Go 标准库的接口实现 mock（无需引入第三方 mock 框架）

## Impact
- 影响的 specs: 无（纯测试增强，不改变现有功能）
- 影响的代码:
  - 新增 `pkg/anonymizer/anonymizer_test.go`: 核心脱敏逻辑测试
  - 新增 `pkg/anonymizer/mock_llm_test.go`: LLM 接口 mock 实现
  - 可选：重构 `pkg/anonymizer/anonymizer.go` 以提升可测试性（仅在必要时）
- 依赖变更: 无（使用 Go 标准测试库）
