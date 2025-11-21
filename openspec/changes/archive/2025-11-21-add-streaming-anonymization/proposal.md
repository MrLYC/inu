# Add Streaming Anonymization

## Why
当前 Anonymizer 的 `AnonymizeText` 方法采用批量处理模式，需要等待 LLM 完整生成响应后才能返回结果。这在处理大文本时存在以下问题：

**现有问题**：
1. 用户体验差：处理长文本时缺乏实时反馈，用户只能等待
2. 内存占用高：需要在内存中缓存完整的响应
3. 响应延迟长：首字节时间（TTFB）等于整个生成时间
4. 无法提前终止：即使用户只需要部分结果也要等待全部生成

**改进后的效果**：
- 实时输出：LLM 生成每个 token 时立即输出，用户获得即时反馈
- 降低内存：流式处理避免大量数据积累
- 更快响应：首字节延迟显著降低
- 保持兼容：现有 API 继续可用，基于流式实现

## What Changes
- **添加** `AnonymizeTextStream` 方法：接受 `io.Writer` 参数，流式输出脱敏文本
- **重构** `AnonymizeText` 方法：基于 `AnonymizeTextStream` 实现，保持向后兼容
- **修改** CLI `anonymize` 命令：使用流式 API，实时输出到 stdout
- **保持** Web API 不变：暂不支持流式响应（未来可选 SSE）

**方法签名对比**：

```go
// 新增：流式方法
func (a *Anonymizer) AnonymizeTextStream(
    ctx context.Context,
    types []string,
    text string,
    writer io.Writer,
) ([]*Entity, error)

// 重构：批量方法（基于流式实现）
func (a *Anonymizer) AnonymizeText(
    ctx context.Context,
    types []string,
    text string,
) (string, []*Entity, error)
```

**流式输出特性**：
- LLM 生成的每个 token 立即写入 writer
- 实体映射在完成后通过返回值提供
- 遇到错误立即中断并返回 error
- Writer 错误和 LLM 错误都会导致函数终止

## Impact
- **影响的 specs**: `cli` spec 需要添加流式输出场景
- **影响的代码**:
  - `pkg/anonymizer/anonymizer.go` - 添加 `AnonymizeTextStream`，重构 `AnonymizeText`
  - `cmd/inu/commands/anonymize.go` - 使用流式 API
  - 测试文件需要添加流式测试
- **破坏性变更**: 否
  - 现有 `AnonymizeText` API 完全兼容
  - CLI 输出行为不变（仍然输出到 stdout）
  - 只是内部实现从批量改为流式
- **性能提升**:
  - 首字节时间：从 O(全文生成时间) 降低到 O(首个token时间)
  - 内存占用：从 O(响应大小) 降低到 O(token大小)
  - 用户体验：大文本处理时显著改善

## Implementation Notes
1. 使用 `llm.Stream()` 代替 `llm.Generate()`
2. 从 `StreamReader` 读取 token 并立即写入 writer
3. 仍需解析 `<<<PAIR>>>` 标记来分离文本和映射
4. `AnonymizeText` 使用 `bytes.Buffer` 作为 writer 调用流式方法
5. 错误处理：任何阶段的错误都立即返回
