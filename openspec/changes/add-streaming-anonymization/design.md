# Streaming Anonymization Design

## Context
当前 Anonymizer 使用 CloudWeGo Eino 的 `BaseChatModel.Generate()` 方法，该方法等待 LLM 完整生成响应后一次性返回。Eino 同时提供了 `BaseChatModel.Stream()` 方法，返回 `*schema.StreamReader[*schema.Message]`，支持逐 token 读取。

**技术约束**：
- LLM 响应格式：`<匿名化文本>\n<<<PAIR>>>\n<JSON映射>`
- 需要完整读取才能解析实体映射
- CLI 默认输出到 stdout，需要实时显示进度
- 测试使用 mock LLM，需要支持流式 mock

## Goals / Non-Goals

**Goals**：
- 提供流式匿名化 API，降低 TTFB 和内存占用
- 保持现有 API 完全兼容
- CLI 实时输出匿名化文本，改善用户体验
- 支持流式测试验证

**Non-Goals**：
- Web API 流式响应（SSE/chunked）- 未来可选
- 并发流式处理多个请求
- 流式处理超大文件（>100MB）
- 自定义输出格式或协议

## Decisions

### Decision 1: 流式方法签名
**选择**：新增 `AnonymizeTextStream(ctx, types, text, writer)` 方法

**理由**：
- `io.Writer` 是 Go 标准接口，灵活且可测试
- 独立方法保持向后兼容，降低风险
- 调用方可选择批量或流式

**备选方案及拒绝理由**：
- 修改现有方法签名：破坏性变更，影响 Web API
- 使用回调函数：不符合 Go 惯用模式，难以测试
- 返回 channel：增加复杂度，需要额外 goroutine 管理

### Decision 2: 响应解析策略
**选择**：边读取边输出文本，完成后解析实体映射

**流程**：
```
1. 调用 llm.Stream() 获取 StreamReader
2. 循环读取 token：
   a. 写入 writer（实时输出）
   b. 追加到 buffer（用于解析映射）
   c. 检测错误立即返回
3. 读取完成后，从 buffer 解析 <<<PAIR>>> 和 JSON
4. 返回实体列表
```

**理由**：
- 无法提前知道 `<<<PAIR>>>` 的位置
- 实体映射必须完整才能解析
- 双写（writer + buffer）开销可接受

### Decision 3: AnonymizeText 重构
**选择**：使用 `bytes.Buffer` 调用 `AnonymizeTextStream`

```go
func (a *Anonymizer) AnonymizeText(ctx context.Context, types []string, text string) (string, []*Entity, error) {
    var buf bytes.Buffer
    entities, err := a.AnonymizeTextStream(ctx, types, text, &buf)
    if err != nil {
        return "", nil, err
    }
    return buf.String(), entities, nil
}
```

**理由**：
- 复用流式逻辑，避免重复代码
- `bytes.Buffer` 实现 `io.Writer`，无额外依赖
- 性能影响可忽略（相比 LLM 调用）

### Decision 4: 错误处理
**选择**：任何错误立即中断并返回

**错误类型**：
1. LLM 流错误：`StreamReader.Recv()` 返回 error
2. Writer 错误：`writer.Write()` 返回 error
3. 解析错误：JSON unmarshal 或格式错误

**处理策略**：
- 所有错误立即返回，不重试
- 已写入 writer 的数据无法撤回（流式特性）
- 错误信息包含上下文（使用 eris.Wrap）

### Decision 5: 流式 Mock
**选择**：在 mock 中实现 `Stream()` 返回预设 token 序列

```go
type mockStreamReader struct {
    tokens []string
    index  int
}

func (m *mockStreamReader) Recv() (*schema.Message, error) {
    if m.index >= len(m.tokens) {
        return nil, io.EOF
    }
    token := m.tokens[m.index]
    m.index++
    return &schema.Message{Content: token}, nil
}
```

**理由**：
- 可控测试：预设 token 序列验证流式逻辑
- 覆盖边界：空响应、错误、大响应
- 与真实 StreamReader 行为一致

## Risks / Trade-offs

### Risk 1: 响应格式变化
**风险**：LLM 改变 `<<<PAIR>>>` 格式导致解析失败

**缓解**：
- 提示词明确要求固定格式
- 添加格式验证测试
- 错误信息清晰指示解析失败位置

### Risk 2: Writer 阻塞
**风险**：Writer 阻塞（如管道满）导致整体延迟

**缓解**：
- 文档说明调用方需提供非阻塞 writer
- CLI 使用 os.Stdout（操作系统缓冲）
- 未来可考虑带缓冲的 writer

### Trade-off: 内存双写
**取舍**：流式输出的同时需要 buffer 存储全文

**影响**：
- 内存占用：O(响应大小) - 与批量模式相同
- 但首字节延迟显著降低
- 用户体验改善大于内存成本

## Migration Plan

### 阶段 1: 添加流式 API（本提案）
- 实现 `AnonymizeTextStream`
- 重构 `AnonymizeText`
- CLI 采用流式
- 添加流式测试

### 阶段 2: Web API 流式支持（未来）
- 添加 SSE endpoint: `POST /api/v1/anonymize/stream`
- 客户端示例和文档
- 性能基准测试

### 阶段 3: 优化（可选）
- 实现真正的流式解析（无需 buffer 全文）
- 支持分段输出实体映射
- 并发处理优化

## Open Questions

1. **是否需要进度回调？**
   - 当前：无进度通知，只有流式输出
   - 可选：添加 `ProgressFunc(percent float64)` 参数
   - 决策：暂不需要，stderr 已有进度信息

2. **是否支持取消？**
   - 当前：依赖 context 取消
   - 可选：显式 `Cancel()` 方法
   - 决策：context 足够，符合 Go 惯例

3. **大文本分段处理？**
   - 当前：整文本传给 LLM
   - 可选：自动分段 + 合并
   - 决策：未来优化，先实现基础流式
