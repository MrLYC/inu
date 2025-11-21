# cli Specification

## ADDED Requirements

### Requirement: 流式脱敏输出
系统 SHALL 在 LLM 生成脱敏文本时实时输出到 stdout，提供即时反馈。

#### Scenario: 实时流式输出脱敏文本
- **WHEN** 用户执行 `inu anonymize --file large_input.txt`
- **THEN** 系统应该在 LLM 生成每个 token 时立即输出到 stdout
- **AND** 用户可以实时看到脱敏进度
- **AND** 不需要等待全部生成完成

#### Scenario: 流式输出到文件
- **WHEN** 用户执行 `inu anonymize --file input.txt --output result.txt`
- **THEN** 系统应该将 LLM 生成的 token 实时写入 result.txt
- **AND** 同时输出到 stdout（默认行为）
- **AND** 文件内容逐步增长而非等待完成后一次性写入

#### Scenario: 流式输出到管道
- **WHEN** 用户执行 `inu anonymize --file input.txt | grep "个人信息"`
- **THEN** 系统应该将 token 实时传递给管道下游命令
- **AND** 下游命令可以立即开始处理
- **AND** 不会因为上游缓冲导致延迟

#### Scenario: 仅流式输出到文件（禁用 stdout）
- **WHEN** 用户执行 `inu anonymize --file input.txt --output result.txt --no-print`
- **THEN** 系统应该将 token 实时写入 result.txt
- **AND** 不向 stdout 输出
- **AND** 仍然保持流式写入行为

#### Scenario: 流式输出被中断
- **WHEN** 用户在流式输出过程中按 Ctrl+C
- **THEN** 系统应该立即停止 LLM 请求
- **AND** 已输出的部分保留在输出中
- **AND** 实体信息可能不完整（未完全解析）

#### Scenario: 流式输出遇到错误
- **WHEN** 流式输出过程中 LLM 返回错误（如网络中断）
- **THEN** 系统应该立即停止输出
- **AND** 显示错误信息到 stderr
- **AND** 已输出的部分文本保留
- **AND** 返回非零退出码

#### Scenario: Writer 错误中断流式输出
- **WHEN** 流式输出到文件时磁盘空间不足
- **THEN** 系统应该立即停止处理
- **AND** 显示清晰的错误信息（磁盘空间不足）
- **AND** 返回非零退出码

### Requirement: 流式输出性能特征
系统 SHALL 通过流式输出降低首字节延迟和改善用户体验。

#### Scenario: 首字节时间显著降低
- **WHEN** 处理大文本（>1KB）进行脱敏
- **THEN** 首字节输出时间应该远小于完整生成时间
- **AND** 用户在秒级内看到首批输出（而非等待数十秒）

#### Scenario: 内存占用保持合理
- **WHEN** 使用流式输出处理大文本
- **THEN** 内存占用应该与批量模式接近
- **AND** 不应因流式处理显著增加内存（允许缓冲开销）

#### Scenario: 向后兼容批量模式
- **WHEN** 使用现有的 CLI 命令和参数
- **THEN** 输出结果应该与之前版本一致
- **AND** 实体信息格式不变
- **AND** 仅实现方式从批量改为流式

## MODIFIED Requirements

### Requirement: 脱敏命令
系统 SHALL 提供 `anonymize` 子命令来脱敏文本中的敏感信息，使用流式输出改善用户体验。

#### Scenario: 从标准输入读取并流式输出到标准输出
- **WHEN** 用户执行 `echo "张三的电话是 13800138000" | inu anonymize`
- **THEN** 系统应该读取标准输入，流式生成脱敏文本
- **AND** 实时输出到标准输出（逐 token）
- **AND** 在流式输出完成后，实体信息输出到 stderr

#### Scenario: 从文件读取并流式输出
- **WHEN** 用户执行 `inu anonymize --file input.txt`
- **THEN** 系统应该读取 input.txt 文件的内容
- **AND** 流式生成并输出脱敏文本到标准输出
- **AND** 用户可以实时看到输出进度

#### Scenario: 流式输出同时保存到文件
- **WHEN** 用户执行 `inu anonymize --file input.txt --output result.txt`
- **THEN** 系统应该流式生成脱敏文本
- **AND** 同时写入 result.txt 和 stdout
- **AND** 两个输出目标都是流式写入

## Implementation Notes
- 流式输出基于 CloudWeGo Eino 的 `StreamReader` 接口
- 内部使用 `llm.Stream()` 替代 `llm.Generate()`
- 实体映射仍需等待完整响应后解析
- 使用 `io.MultiWriter` 支持同时输出到多个目标
