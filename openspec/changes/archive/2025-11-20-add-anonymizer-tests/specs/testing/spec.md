# testing Specification Delta

## ADDED Requirements

### Requirement: Anonymizer 单元测试
系统 SHALL 为 `pkg/anonymizer` 包提供完整的单元测试覆盖，使用 mock 避免真实 LLM API 调用。

#### Scenario: 测试脱敏单个实体
- **WHEN** 调用 `AnonymizeText` 脱敏包含一个实体的文本
- **THEN** 系统应该返回脱敏后的文本和对应的实体信息

#### Scenario: 测试脱敏多个实体
- **WHEN** 调用 `AnonymizeText` 脱敏包含多种实体类型的文本
- **THEN** 系统应该正确识别并脱敏所有实体，返回完整的实体列表

#### Scenario: 测试脱敏空文本
- **WHEN** 调用 `AnonymizeText` 时输入为空字符串
- **THEN** 系统应该正确处理并返回空结果或适当的响应

#### Scenario: 测试脱敏无匹配实体
- **WHEN** LLM 返回无匹配实体（"None"）
- **THEN** 系统应该返回原文本和空实体列表，不报错

#### Scenario: 测试 LLM API 调用失败
- **WHEN** LLM API 调用失败（模拟网络错误或认证失败）
- **THEN** 系统应该返回明确的错误信息，包含失败原因

#### Scenario: 测试响应格式错误
- **WHEN** LLM 返回的响应格式不正确（缺少 `<<<PAIR>>>` 分隔符）
- **THEN** 系统应该返回格式错误，不应崩溃

#### Scenario: 测试 JSON 解析失败
- **WHEN** LLM 返回的 mapping 部分不是有效的 JSON
- **THEN** 系统应该返回 JSON 解析错误，不应崩溃

#### Scenario: 测试实体 key 格式错误
- **WHEN** LLM 返回的实体 key 不符合 `<EntityType[ID].Category.Detail>` 格式
- **THEN** 系统应该返回格式错误，指出具体的错误 key

### Requirement: 文本还原测试
系统 SHALL 测试 `RestoreText` 方法的文本还原功能。

#### Scenario: 测试还原单个实体
- **WHEN** 调用 `RestoreText` 还原包含一个实体占位符的文本
- **THEN** 系统应该将占位符替换为原始值，返回还原后的文本

#### Scenario: 测试还原多个实体
- **WHEN** 调用 `RestoreText` 还原包含多个实体占位符的文本
- **THEN** 系统应该正确替换所有占位符，返回完整还原的文本

#### Scenario: 测试还原空文本
- **WHEN** 调用 `RestoreText` 时输入为空字符串
- **THEN** 系统应该返回空字符串，不报错

#### Scenario: 测试还原空实体列表
- **WHEN** 调用 `RestoreText` 时实体列表为空
- **THEN** 系统应该返回原文本不变，不报错

#### Scenario: 测试还原无匹配占位符
- **WHEN** 文本中不包含任何实体占位符
- **THEN** 系统应该返回原文本不变，不报错

#### Scenario: 测试实体无 Values
- **WHEN** 实体的 Values 数组为空
- **THEN** 系统应该跳过该实体，继续处理其他实体

#### Scenario: 测试完整往返
- **WHEN** 对文本执行脱敏后再执行还原
- **THEN** 还原后的文本应该与原文本完全一致

### Requirement: Mock LLM 实现
系统 SHALL 提供 `model.BaseChatModel` 接口的 mock 实现用于测试。

#### Scenario: Mock 返回预定义响应
- **WHEN** 测试调用 mock LLM 的 `Generate` 方法
- **THEN** mock 应该返回预先配置的响应内容，无需真实 API 调用

#### Scenario: Mock 返回错误
- **WHEN** 测试需要模拟 LLM 调用失败
- **THEN** mock 应该返回预先配置的错误，模拟真实失败场景

#### Scenario: Mock 响应格式可配置
- **WHEN** 测试需要不同格式的 LLM 响应
- **THEN** mock 应该支持灵活配置响应内容（脱敏文本、实体映射）

#### Scenario: Mock 不需要环境变量
- **WHEN** 运行使用 mock 的测试
- **THEN** 测试应该在未设置 `OPENAI_API_KEY` 等环境变量的情况下成功运行

### Requirement: 测试覆盖率和质量
系统 SHALL 确保测试的覆盖率和质量达到标准。

#### Scenario: 覆盖核心方法
- **WHEN** 运行测试套件
- **THEN** 应该覆盖 `AnonymizeText`、`RestoreText`、`New` 等所有核心方法

#### Scenario: 测试覆盖率要求
- **WHEN** 运行 `go test -cover`
- **THEN** `pkg/anonymizer` 包的测试覆盖率应该 ≥ 80%

#### Scenario: 测试执行速度
- **WHEN** 运行完整的 anonymizer 测试套件
- **THEN** 所有测试应该在 1 秒内完成（无网络依赖）

#### Scenario: 测试稳定性
- **WHEN** 在 CI 环境中多次运行测试
- **THEN** 测试结果应该一致且稳定，无随机失败

#### Scenario: 测试独立性
- **WHEN** 运行单个测试或完整测试套件
- **THEN** 每个测试应该独立运行，不依赖其他测试的状态或执行顺序
