# 统一流式接口设计

## Why

### 问题
当前 `Anonymizer` 接口设计存在以下问题:

1. **RestoreText 缺乏透明度**:
   - 当部分占位符无法还原时,用户无法知道哪些实体还原失败
   - 只能通过比较输出中是否还有占位符来推断失败,但无法确定具体原因
   - Interactive 和 Restore 命令只能简单提示"部分占位符无法还原",没有详细信息

2. **接口不一致**:
   - `AnonymizeTextStream` 使用 writer 参数实现流式输出
   - `AnonymizeText` 实际上是 `AnonymizeTextStream` 的包装器
   - `RestoreText` 返回完整字符串,没有流式输出选项
   - Web API 和 CLI 命令需要将完整结果加载到内存

3. **命名混淆**:
   - `AnonymizeText` 和 `AnonymizeTextStream` 功能重复
   - 用户需要理解两个方法的差异
   - 大多数情况下应该使用流式版本以获得更好的性能

### 影响范围
- **严重性**: 🟡 Medium - 影响用户体验和 API 一致性
- **受影响用户**: CLI 用户、Web API 用户、SDK 集成者
- **受影响模块**:
  - `pkg/anonymizer` - 核心接口
  - `cmd/inu/commands` - CLI 命令(anonymize, restore, interactive)
  - `pkg/web/handlers` - Web API handlers
  - 测试代码

## What Changes

### 1. 统一接口命名
- **删除** `AnonymizeText()` 方法
- **重命名** `AnonymizeTextStream()` → `Anonymize()`
- **修改** `RestoreText()` 签名:
  - 添加 `writer io.Writer` 参数(输出还原后的文本)
  - 返回 `[]string` (还原失败的占位符列表)
  - 完整签名: `Anonymize(ctx, types, text, writer) ([]*Entity, error)`
  - 完整签名: `RestoreText(ctx, entities, text, writer) ([]string, error)`

### 2. 新接口定义
```go
type Anonymizer interface {
    // Anonymize 脱敏文本并流式输出到 writer
    // 返回识别到的实体列表和可能的错误
    Anonymize(ctx context.Context, types []string, text string, writer io.Writer) ([]*Entity, error)

    // RestoreText 还原文本并流式输出到 writer
    // 返回无法还原的占位符列表(如果为空则全部还原成功)
    RestoreText(ctx context.Context, entities []*Entity, text string, writer io.Writer) ([]string, error)
}
```

### 3. 命令适配

#### anonymize 命令
```go
// 旧代码
entities, err := anon.AnonymizeTextStream(ctx, types, input, writer)

// 新代码
entities, err := anon.Anonymize(ctx, types, input, writer)
```

#### restore 命令
```go
// 旧代码
result, err := anon.RestoreText(ctx, entities, input)
if err != nil {
    return err
}
cli.WriteOutput(result, noPrint, output)

// 新代码
unrestoredPlaceholders, err := anon.RestoreText(ctx, entities, input, writer)
if err != nil {
    return err
}

// 显示警告
if len(unrestoredPlaceholders) > 0 {
    cli.WarningMessage("Warning: %d placeholder(s) could not be restored:", len(unrestoredPlaceholders))
    for _, placeholder := range unrestoredPlaceholders {
        cli.WarningMessage("  - %s", placeholder)
    }
}
```

#### interactive 命令
```go
// 旧代码
entities, err := anon.AnonymizeTextStream(ctx, types, input, os.Stdout)
// ...
restoredText, err := anon.RestoreText(ctx, entities, processedText)
if err != nil {
    fmt.Fprintln(os.Stderr, "Warning: Some placeholders could not be restored")
    restoredText = processedText
}
fmt.Println(restoredText)

// 新代码
entities, err := anon.Anonymize(ctx, types, input, os.Stdout)
// ...
unrestored, err := anon.RestoreText(ctx, entities, processedText, os.Stdout)
if len(unrestored) > 0 {
    cli.WarningMessage("\nWarning: %d placeholder(s) could not be restored:", len(unrestored))
    for _, p := range unrestored {
        cli.WarningMessage("  - %s", p)
    }
}
```

### 4. Web API 适配

#### Anonymize Handler
```go
// 旧代码
anonymizedText, entities, err := anon.AnonymizeText(c.Request.Context(), entityTypes, req.Text)

// 新代码
var buf bytes.Buffer
entities, err := anon.Anonymize(c.Request.Context(), entityTypes, req.Text, &buf)
anonymizedText := buf.String()
```

#### Restore Handler
保持当前行为(不显示未还原占位符),因为:
- Web UI 可视化更好处理(高亮显示未还原占位符)
- 可以在 response 中添加可选的 `unrestored_placeholders` 字段

### 5. 向后兼容性

**破坏性变更**:
- ✅ `AnonymizeText()` 方法被移除
- ✅ `AnonymizeTextStream()` 重命名为 `Anonymize()`
- ✅ `RestoreText()` 签名变化(添加 writer 参数,返回值变化)

**缓解策略**:
- 这是内部 API,不影响最终用户
- Web API 端点保持不变
- CLI 命令行为保持兼容(只是增加了警告信息)

## Benefits

### 用户体验改进
- ✅ 清晰的错误反馈: 用户知道哪些占位符无法还原
- ✅ 便于调试: 可以针对性地检查实体文件或文本内容
- ✅ 更好的交互式体验: 即时了解还原状态

### API 一致性
- ✅ 统一的流式接口设计
- ✅ 方法命名更简洁(Anonymize vs AnonymizeText/AnonymizeTextStream)
- ✅ 相同的参数模式(都使用 writer)

### 性能优化
- ✅ 避免不必要的字符串拷贝
- ✅ 支持大文件处理(流式输出)
- ✅ 减少内存占用

## Risks

### 技术风险
- **低风险**: 变更范围明确,测试覆盖充分
- **接口稳定性**: 内部 API 变更不影响外部用户
- **测试工作量**: 需要更新所有相关测试用例

### 迁移风险
- **影响范围**: pkg/anonymizer, cmd/inu/commands, pkg/web/handlers
- **测试覆盖**: 单元测试 + 集成测试确保行为一致
- **回滚成本**: 低(Git revert 即可)

### 缓解措施
- 完整的测试覆盖(单元测试 + 集成测试)
- 详细的迁移文档
- 分阶段提交(接口变更 → CLI 适配 → Web 适配 → 测试)

## Alternatives Considered

### 1. 保留 AnonymizeText 作为便捷方法
**优点**:
- 向后兼容
- 对简单用例更友好

**缺点**:
- API 膨胀
- 用户需要理解两个方法的差异
- 维护成本增加

**结论**: 不采纳,保持接口简洁

### 2. RestoreText 返回详细的错误对象
**示例**:
```go
type RestoreError struct {
    Placeholder string
    Reason      string // "not_found" | "no_values" | "format_error"
}

RestoreText(...) (string, []RestoreError, error)
```

**优点**:
- 提供更详细的错误信息

**缺点**:
- 复杂度增加
- 大多数情况下只需要知道哪些占位符失败即可

**结论**: 暂不采纳,可以作为未来增强

### 3. 添加新方法而非修改现有方法
**示例**:
```go
AnonymizeStream() // 新方法
AnonymizeText()   // 保留旧方法
RestoreTextStream() // 新方法
RestoreText()      // 保留旧方法
```

**优点**:
- 完全向后兼容

**缺点**:
- API 膨胀严重
- 用户困惑(应该用哪个?)
- 维护成本高

**结论**: 不采纳,清晰胜于兼容

## Affected Specs
- `cli` - 更新 anonymize, restore, interactive 命令行为
- `web-api` - 可选:添加 unrestored_placeholders 字段到 RestoreResponse

## Migration Guide

### 对于 SDK 集成者

**场景 1: 使用 AnonymizeText**
```go
// 旧代码
result, entities, err := anon.AnonymizeText(ctx, types, text)

// 新代码
var buf bytes.Buffer
entities, err := anon.Anonymize(ctx, types, text, &buf)
result := buf.String()
```

**场景 2: 使用 AnonymizeTextStream**
```go
// 旧代码
entities, err := anon.AnonymizeTextStream(ctx, types, text, writer)

// 新代码
entities, err := anon.Anonymize(ctx, types, text, writer)
```

**场景 3: 使用 RestoreText**
```go
// 旧代码
restored, err := anon.RestoreText(ctx, entities, text)
if err != nil {
    return err
}
fmt.Println(restored)

// 新代码
var buf bytes.Buffer
unrestored, err := anon.RestoreText(ctx, entities, text, &buf)
if err != nil {
    return err
}
restored := buf.String()
if len(unrestored) > 0 {
    log.Printf("Warning: %d placeholders not restored: %v", len(unrestored), unrestored)
}
fmt.Println(restored)
```

### 对于最终用户

**CLI 用户**: 无需变更,命令行为保持一致,只是会看到更详细的警告信息

**Web UI 用户**: 无需变更,API 响应格式保持兼容
