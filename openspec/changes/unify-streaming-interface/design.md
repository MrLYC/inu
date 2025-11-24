# 设计文档: 统一流式接口

## Context

### 当前架构

```go
type Anonymizer interface {
    // 方式 1: 返回完整字符串(内部包装 AnonymizeTextStream)
    AnonymizeText(ctx context.Context, entityTypes []string, text string) (string, []*Entity, error)

    // 方式 2: 流式输出到 writer
    AnonymizeTextStream(ctx context.Context, entityTypes []string, text string, writer io.Writer) ([]*Entity, error)

    // 问题: 只返回完整字符串,无法告知哪些占位符未还原
    RestoreText(ctx context.Context, entities []*Entity, text string) (string, error)
}
```

**当前 RestoreText 实现** (pkg/anonymizer/anonymizer.go:249-271):
```go
func (a *anonymizer) RestoreText(ctx context.Context, entities []*Entity, text string) (string, error) {
    placeholderMap := make(map[string]string)
    for _, entity := range entities {
        if len(entity.Values) == 0 {
            continue
        }
        placeholderMap[entity.Key] = entity.Values[0]
    }

    result := placeholderRegex.ReplaceAllStringFunc(text, func(match string) string {
        normalized := normalizePlaceholder(match)
        if original, ok := placeholderMap[normalized]; ok {
            return original
        }
        return match // 保持占位符不变,但无法告知调用方
    })

    return result, nil
}
```

**核心问题**: ReplaceAllStringFunc 中遇到未匹配的占位符时只能返回原样,调用方无法知道哪些占位符未还原。

### 使用场景

**CLI restore 命令** (cmd/inu/commands/restore.go:99):
```go
result, err := anon.RestoreText(ctx, entities, input)
if err != nil {
    return err
}
cli.WriteOutput(result, noPrint, output)
// 问题: 无法向用户显示哪些占位符未还原
```

**CLI interactive 命令** (cmd/inu/commands/interactive.go:129-134):
```go
restoredText, err := anon.RestoreText(ctx, entities, processedText)
if err != nil {
    fmt.Fprintln(os.Stderr, "Warning: Some placeholders could not be restored")
    restoredText = processedText
}
fmt.Println(restoredText)
// 问题: 只能显示通用警告,无法告知具体哪些占位符失败
```

**Web API restore handler** (pkg/web/handlers/restore.go:54):
```go
restoredText, err := a.restorer.RestoreText(c.Request.Context(), entities, text)
if err != nil {
    return err
}
c.JSON(http.StatusOK, RestoreResponse{RestoredText: restoredText})
// 问题: 无法在响应中包含失败信息
```

## Goals

### 主要目标

1. **统一接口设计**
   - 所有方法都使用 writer 参数(流式输出)
   - 删除冗余方法(AnonymizeText)
   - 简化方法命名(AnonymizeTextStream → Anonymize)

2. **增强错误报告**
   - RestoreText 返回未还原的占位符列表
   - CLI 命令显示详细的失败信息
   - Web API 可选择性地包含失败信息

3. **保持性能**
   - 流式输出避免大字符串拷贝
   - 不增加额外的内存开销
   - 不影响现有的批处理性能

### 非目标

- ❌ 不改变占位符格式规范
- ❌ 不修改实体识别逻辑
- ❌ 不改变 Web API 的 HTTP 接口
- ✅ ~~不提供详细的失败原因~~ → **已采纳**: 区分 "not_found" vs "empty_values"

## Decisions

### 决策 1: 新接口签名设计

**选择**:
```go
type Anonymizer interface {
    Anonymize(ctx context.Context, entityTypes []string, text string, writer io.Writer) ([]*Entity, error)
    RestoreText(ctx context.Context, entities []*Entity, text string, writer io.Writer) ([]RestoreFailure, error)
}

type RestoreFailure struct {
    Placeholder string // normalized placeholder
    Reason      string // "not_found" | "empty_values"
}
```

**理由**:
- ✅ 接口一致性: 两个方法都使用 writer
- ✅ 简洁性: 删除 AnonymizeText,重命名 AnonymizeTextStream
- ✅ 透明性: RestoreText 返回失败的占位符列表

**考虑的替代方案**:

**方案 A: 返回详细的错误对象**
```go
type RestoreError struct {
    Placeholder string
    Reason      string // "not_found" | "no_values" | "format_error"
}

RestoreText(...) ([]RestoreError, error)
```
❌ 拒绝理由: 复杂度高,大多数情况下只需要知道哪些占位符失败

**方案 B: 保留 AnonymizeText 作为便捷方法**
```go
AnonymizeText(ctx, types, text) (string, []*Entity, error) {
    var buf bytes.Buffer
    entities, err := Anonymize(ctx, types, text, &buf)
    return buf.String(), entities, err
}
```
❌ 拒绝理由: API 膨胀,用户困惑(应该用哪个?)

**方案 C: 添加新方法而非修改现有方法**
```go
AnonymizeStream(...)  // 新方法
RestoreTextStream(...) // 新方法
AnonymizeText(...)    // 保留旧方法
RestoreText(...)      // 保留旧方法
```
❌ 拒绝理由: API 膨胀严重,维护成本高

### 决策 2: 未还原占位符的返回格式

**选择**: 返回 `[]RestoreFailure` (包含占位符和失败原因)

**定义**:
```go
type RestoreFailure struct {
    Placeholder string // normalized placeholder (e.g., "<个人信息[1].姓名.全名>")
    Reason      string // "not_found" | "empty_values"
}
```

**理由**:
- ✅ **区分失败原因**: 用户知道是"实体不存在"还是"实体无值"
- ✅ **调试价值高**: "not_found" → 检查实体文件; "empty_values" → 检查 LLM 输出
- ✅ **实现简单**: 只需多维护一个 `emptyKeys` map,代码增加不到 10 行
- ✅ **用户体验更好**: CLI 可以显示更详细的警告信息

**示例**:
```go
// entities 包含:
// {Key: "<个人信息[0].姓名.全名>", Values: ["Alice"]}
// {Key: "<个人信息[1].姓名.全名>", Values: []}  ← 无值
// 缺少 <个人信息[2].姓名.全名>

// text: "{{个人信息[0].姓名.全名}} met {{个人信息[1].姓名.全名}} and {{个人信息[2].姓名.全名}}"

failures, err := anon.RestoreText(ctx, entities, text, writer)
// failures = []RestoreFailure{
//     {Placeholder: "<个人信息[1].姓名.全名>", Reason: "empty_values"},
//     {Placeholder: "<个人信息[2].姓名.全名>", Reason: "not_found"},
// }
// writer 输出: "Alice met {{个人信息[1].姓名.全名}} and {{个人信息[2].姓名.全名}}"
```

**CLI 输出效果**:
```
Alice met {{个人信息[1].姓名.全名}} and {{个人信息[2].姓名.全名}}
Warning: 2 placeholder(s) could not be restored:
  - <个人信息[1].姓名.全名> (entity has no values)
  - <个人信息[2].姓名.全名> (not found in entities file)
```

**考虑的替代方案**:

**方案 A: 只返回 []string (不区分原因)**
```go
RestoreText(...) ([]string, error)
// 返回: ["<个人信息[1].姓名.全名>", "<个人信息[2].姓名.全名>"]
```
❌ 拒绝理由: 用户无法知道如何修复问题

**方案 B: 返回 []*Entity (失败的实体对象)**
```go
RestoreText(...) ([]*Entity, error)
```
❌ 拒绝理由:
- text 中可能有占位符但 entities 中没有对应实体
- 无法返回"不存在的实体"

### 决策 3: CLI Warning 输出方式

**选择**: 使用 stderr 输出警告,格式化为列表

**示例输出**:
```
Alice met Bob and {{PERSON_2}}
Warning: 1 placeholder(s) could not be restored:
  - PERSON_2
```

**理由**:
- ✅ 不影响主输出(stdout 仍然是纯文本)
- ✅ 可以被管道忽略(`inu restore ... 2>/dev/null`)
- ✅ 清晰的视觉反馈

**实现**:
```go
// 使用 cli 包的 WarningMessage 方法(输出到 stderr)
if len(unrestoredPlaceholders) > 0 {
    cli.WarningMessage("Warning: %d placeholder(s) could not be restored:", len(unrestoredPlaceholders))
    for _, placeholder := range unrestoredPlaceholders {
        cli.WarningMessage("  - %s", placeholder)
    }
}
```

**考虑的替代方案**:

**方案 A: 在输出文本中内联标记**
```
Alice met Bob and [FAILED:{{PERSON_2}}]
```
❌ 拒绝理由: 污染输出内容,破坏文本格式

**方案 B: 使用颜色高亮**
```
Alice met Bob and \033[31m{{PERSON_2}}\033[0m
```
❌ 拒绝理由:
- 管道到文件时有转义符
- 不是所有终端支持颜色
- stderr 警告已经足够清晰

**方案 C: 保存到单独的错误文件**
```
inu restore ... --error-log restore-errors.txt
```
❌ 拒绝理由: 增加参数复杂度,大多数情况不需要

### 决策 4: 向后兼容策略

**选择**: 接受破坏性变更,提供清晰的迁移指南

**理由**:
- ✅ 这是内部 API,不影响最终用户
- ✅ CLI 命令行为保持兼容(只是增加了警告)
- ✅ Web API 端点保持不变
- ✅ SDK 集成者可以通过文档快速迁移

**破坏性变更清单**:
1. `AnonymizeText()` 方法被删除
2. `AnonymizeTextStream()` 重命名为 `Anonymize()`
3. `RestoreText()` 签名变化:
   - 添加 `writer io.Writer` 参数
   - 返回值从 `(string, error)` 变为 `([]string, error)`

**缓解措施**:
- 在 proposal.md 中包含详细的迁移指南
- 所有测试用例更新确保行为一致
- Web handlers 使用 bytes.Buffer 包装以保持响应格式不变

**考虑的替代方案**:

**方案 A: 保留旧方法 + deprecated 标记**
```go
// Deprecated: Use Anonymize instead
func (a *anonymizer) AnonymizeText(...) (string, []*Entity, error)

// Deprecated: Use Anonymize instead
func (a *anonymizer) AnonymizeTextStream(...) ([]*Entity, error)
```
❌ 拒绝理由:
- 维护两套实现
- 用户仍然困惑(应该用哪个?)
- 最终还是要删除,不如一次性改

**方案 B: 添加新接口,保留旧接口**
```go
type AnonymizerV2 interface {
    Anonymize(...)
    RestoreText(...) ([]string, error)
}

type Anonymizer interface { // 保留旧接口
    AnonymizeText(...)
    AnonymizeTextStream(...)
    RestoreText(...) (string, error)
}
```
❌ 拒绝理由: 接口膨胀,用户需要选择使用哪个版本

### 决策 5: Web API 适配策略

**选择**: 保持 HTTP 响应格式不变,内部使用 bytes.Buffer 适配

**RestoreResponse 结构** (可选扩展):
```go
type RestoreResponse struct {
    RestoredText          string   `json:"restored_text"`
    UnrestoredPlaceholders []string `json:"unrestored_placeholders,omitempty"` // 可选字段
}
```

**实现**:
```go
var buf bytes.Buffer
unrestored, err := a.restorer.RestoreText(c.Request.Context(), entities, text, &buf)
if err != nil {
    return err
}

c.JSON(http.StatusOK, RestoreResponse{
    RestoredText:          buf.String(),
    UnrestoredPlaceholders: unrestored, // 如果为空则被省略
})
```

**理由**:
- ✅ 向后兼容: 旧客户端忽略新字段
- ✅ 可扩展: 新客户端可以处理失败信息
- ✅ 性能: bytes.Buffer 开销极小

**考虑的替代方案**:

**方案 A: 添加新的 API 端点**
```
POST /v2/restore  // 新端点,返回详细信息
POST /restore     // 旧端点,保持不变
```
❌ 拒绝理由: API 版本膨胀,维护成本高

**方案 B: 使用错误码表示部分失败**
```go
// HTTP 207 Multi-Status
c.JSON(http.StatusMultiStatus, RestoreResponse{
    RestoredText: buf.String(),
    Errors: []string{"PERSON_2 not found"},
})
```
❌ 拒绝理由:
- 部分失败不是错误(正常业务场景)
- HTTP 207 语义不清晰

## Implementation Details

### 核心逻辑修改

**RestoreText 新实现** (伪代码):
```go
func (a *anonymizer) RestoreText(ctx context.Context, entities []*Entity, text string, writer io.Writer) ([]RestoreFailure, error) {
    // 1. 构建两个映射
    entityMap := make(map[string]string)  // 有值的实体
    emptyKeys := make(map[string]bool)    // 无值的实体

    for _, entity := range entities {
        normalizedKey := normalizePlaceholder(entity.Key)
        if len(entity.Values) == 0 {
            emptyKeys[normalizedKey] = true
        } else {
            entityMap[normalizedKey] = entity.Values[0]
        }
    }

    // 2. 收集失败信息
    var failures []RestoreFailure
    seenFailures := make(map[string]bool)  // 去重

    // 3. 流式替换并输出
    lastIndex := 0
    matches := placeholderRegex.FindAllStringIndex(text, -1)

    for _, match := range matches {
        // 写入前面的文本
        if _, err := writer.Write([]byte(text[lastIndex:match[0]])); err != nil {
            return nil, err
        }

        // 处理占位符
        placeholder := text[match[0]:match[1]]
        normalizedKey := normalizePlaceholder(placeholder)

        if value, exists := entityMap[normalizedKey]; exists {
            // 还原成功
            if _, err := writer.Write([]byte(value)); err != nil {
                return nil, err
            }
        } else {
            // 还原失败,保留占位符
            if _, err := writer.Write([]byte(placeholder)); err != nil {
                return nil, err
            }

            // 记录失败原因
            if !seenFailures[normalizedKey] {
                reason := "not_found"
                if emptyKeys[normalizedKey] {
                    reason := "empty_values"
                }
                failures = append(failures, RestoreFailure{
                    Placeholder: normalizedKey,
                    Reason:      reason,
                })
                seenFailures[normalizedKey] = true
            }
        }

        lastIndex = match[1]
    }

    // 写入剩余文本
    if _, err := writer.Write([]byte(text[lastIndex:])); err != nil {
        return nil, err
    }

    return failures, nil
}
```

**关键变化**:
1. 使用 `FindAllStringIndex` 而非 `ReplaceAllStringFunc` (更灵活)
2. 流式写入到 writer(避免字符串拼接)
3. 维护两个 map: `entityMap` (有值) 和 `emptyKeys` (无值)
4. 收集失败占位符并区分原因: "not_found" vs "empty_values"
5. 返回 `[]RestoreFailure` 而非完整字符串

### 测试策略

**单元测试** (pkg/anonymizer/anonymizer_test.go):
```go
func TestRestoreText_WithUnrestoredPlaceholders(t *testing.T) {
    entities := []*Entity{
        {Key: "<个人信息[0].姓名.全名>", Values: []string{"Alice"}},
    }
    text := "{{个人信息[0].姓名.全名}} met {{个人信息[1].姓名.全名}}"

    var buf bytes.Buffer
    failures, err := anon.RestoreText(ctx, entities, text, &buf)

    assert.NoError(t, err)
    assert.Equal(t, "Alice met {{个人信息[1].姓名.全名}}", buf.String())
    assert.Len(t, failures, 1)
    assert.Equal(t, "<个人信息[1].姓名.全名>", failures[0].Placeholder)
    assert.Equal(t, "not_found", failures[0].Reason)
}

func TestRestoreText_EmptyValues(t *testing.T) {
    entities := []*Entity{
        {Key: "<个人信息[0].姓名.全名>", Values: []string{"Alice"}},
        {Key: "<个人信息[1].姓名.全名>", Values: []string{}},  // 无值
    }
    text := "{{个人信息[0].姓名.全名}} met {{个人信息[1].姓名.全名}}"

    var buf bytes.Buffer
    failures, err := anon.RestoreText(ctx, entities, text, &buf)

    assert.NoError(t, err)
    assert.Equal(t, "Alice met {{个人信息[1].姓名.全名}}", buf.String())
    assert.Len(t, failures, 1)
    assert.Equal(t, "<个人信息[1].姓名.全名>", failures[0].Placeholder)
    assert.Equal(t, "empty_values", failures[0].Reason)
}

func TestRestoreText_AllRestored(t *testing.T) {
    entities := []*Entity{
        {Key: "<个人信息[0].姓名.全名>", Values: []string{"Alice"}},
        {Key: "<个人信息[1].姓名.全名>", Values: []string{"Bob"}},
    }
    text := "{{个人信息[0].姓名.全名}} met {{个人信息[1].姓名.全名}}"

    var buf bytes.Buffer
    failures, err := anon.RestoreText(ctx, entities, text, &buf)

    assert.NoError(t, err)
    assert.Equal(t, "Alice met Bob", buf.String())
    assert.Empty(t, failures)
}
```

**集成测试** (CLI 命令):
```go
func TestRestoreCommand_WithWarnings(t *testing.T) {
    // 准备测试文件
    // ...

    stdout, stderr := captureOutput(func() {
        cmd := newRestoreCommand()
        cmd.SetArgs([]string{"--entities", entitiesFile, inputFile})
        cmd.Execute()
    })

    assert.Contains(t, stdout, "Alice met {{个人信息[1].姓名.全名}}")
    assert.Contains(t, stderr, "Warning: 1 placeholder(s) could not be restored")
    assert.Contains(t, stderr, "- <个人信息[1].姓名.全名> (not found in entities file)")
}

func TestRestoreCommand_EmptyValuesWarning(t *testing.T) {
    // entities 包含无值实体
    // ...

    stdout, stderr := captureOutput(func() {
        cmd := newRestoreCommand()
        cmd.Execute()
    })

    assert.Contains(t, stderr, "(entity has no values)")
}
```

### 迁移步骤

**阶段 1: 核心接口修改**
1. 修改 `pkg/anonymizer/anonymizer.go` 接口定义
2. 实现新的 `RestoreText` 逻辑
3. 删除 `AnonymizeText` 方法
4. 重命名 `AnonymizeTextStream` → `Anonymize`
5. 更新所有单元测试

**阶段 2: CLI 命令适配**
1. 修改 `cmd/inu/commands/anonymize.go`
2. 修改 `cmd/inu/commands/restore.go` (添加 warning 输出)
3. 修改 `cmd/inu/commands/interactive.go` (添加 warning 输出)
4. 更新集成测试

**阶段 3: Web API 适配**
1. 修改 `pkg/web/handlers/anonymize.go` (使用 bytes.Buffer)
2. 修改 `pkg/web/handlers/restore.go` (添加 unrestored_placeholders 字段)
3. 更新 handler 测试

**阶段 4: 文档更新**
1. 更新 README.md 示例
2. 更新 API 文档(如果有)
3. 更新 spec deltas

## Trade-offs

### 优点
- ✅ 接口一致性: 统一使用 writer 模式
- ✅ 用户体验: 清晰的错误反馈
- ✅ 性能: 流式输出避免大字符串拷贝
- ✅ 可维护性: 减少冗余方法

### 缺点
- ❌ 破坏性变更: 需要更新所有调用代码
- ❌ 测试工作量: 大量测试需要更新
- ❌ 迁移成本: SDK 集成者需要适配(如果有外部用户)

### 权衡决策
**接受破坏性变更的理由**:
1. 这是内部 API,外部用户影响有限
2. 长期收益(一致性、可维护性)大于短期迁移成本
3. 清晰的迁移指南降低迁移难度
4. 最终用户(CLI/Web)体验保持兼容

## Open Questions

### Q1: 是否需要在 Web API 响应中包含 unrestored_placeholders?
**答案**: 是,作为可选字段(omitempty)
- 向后兼容旧客户端
- 新客户端可以利用这个信息改进 UI(如高亮显示未还原占位符)

### Q2: 是否需要区分"占位符不存在"和"实体没有值"?
**答案**: 暂不需要
- 对用户来说结果相同(无法还原)
- 可以作为未来增强(如果有明确需求)
- 当前设计不阻碍未来扩展

### Q3: Interactive 命令是否需要交互式地提示用户修复失败的占位符?
**答案**: 暂不需要
- 增加复杂度
- 大多数情况下用户会直接检查实体文件
- 可以作为未来增强(如果有明确需求)

### Q4: 是否需要提供--strict模式(有未还原占位符时返回错误)?
**答案**: 可选,作为后续 PR
- 有些用户可能希望这种行为(CI/CD 场景)
- 当前设计不阻碍未来添加
- 不在本次变更范围内

## Success Metrics

**功能完整性**:
- ✅ 所有单元测试通过
- ✅ 所有集成测试通过
- ✅ CLI 命令行为符合预期
- ✅ Web API 响应格式兼容

**用户体验**:
- ✅ CLI 用户可以看到详细的失败占位符
- ✅ 警告信息清晰易懂
- ✅ 不影响正常工作流(警告不阻塞)

**代码质量**:
- ✅ 测试覆盖率不降低
- ✅ 接口设计简洁一致
- ✅ 文档完整准确

**性能**:
- ✅ 流式输出性能与之前一致
- ✅ 无额外的内存开销
- ✅ 大文件处理能力不下降
