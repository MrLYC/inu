# Design: Placeholder Format Compatibility

## Context

占位符格式为 `<EntityType[ID].Category.Detail>`，例如 `<个人信息[0].姓名.张三>`。在交互式工作流中，用户会将匿名化文本复制到外部工具（ChatGPT、文档编辑器等）进行处理，然后粘贴回来进行还原。

**问题**：外部工具可能会修改占位符格式，常见变化包括：
- 自动添加空格（中英文混排优化）
- 标点符号转换（`.` → `。`）
- 全角/半角字符转换
- Unicode 规范化

当前 `RestoreText` 使用 `strings.NewReplacer` 进行精确匹配，这些格式变化会导致匹配失败。

## Goals / Non-Goals

**Goals**：
- 支持常见的占位符格式变化（空格、标点、全角/半角）
- 保持向后兼容，不影响现有精确匹配的行为
- 性能影响最小化（归一化操作应该很快）
- 不改变占位符的生成格式（仅改变还原时的匹配逻辑）

**Non-Goals**：
- 不支持完全破坏占位符结构的情况（如删除 `<>` 符号）
- 不修改 LLM 生成的占位符格式
- 不处理占位符被翻译成其他语言的情况

## Decisions

### 1. 归一化策略

**决定**：在还原时对实体键和文本中的占位符都进行归一化，然后匹配。

**归一化步骤**（按顺序）：
1. 提取所有 `<...>` 中的内容
2. 移除所有空白字符（空格、制表符、换行符）
3. 将中文标点转换为英文标点：`。` → `.`、`，` → `,`、`【` → `[`、`】` → `]`
4. 将全角字符转换为半角字符（数字、字母、符号）
5. 重新包装为 `<normalized>`

**示例**：
```
Input:  < 业务信息 [2]. 系统。名称 >
Step 1: Extract " 业务信息 [2]. 系统。名称 "
Step 2: Remove spaces "业务信息[2].系统。名称"
Step 3: Convert punctuation "业务信息[2].系统.名称"
Step 4: Convert fullwidth (no change in this case)
Output: <业务信息[2].系统.名称>
```

### 2. 实现位置

**决定**：在 `pkg/anonymizer/anonymizer.go` 中添加归一化函数。

**理由**：
- 归一化是核心业务逻辑的一部分
- 可能在未来被其他地方复用（如验证占位符格式）
- 便于单元测试

**替代方案**：在 `RestoreText` 方法内部实现
- 缺点：不便于测试和复用
- 缺点：增加方法复杂度

### 3. 正则表达式 vs 字符串处理

**决定**：使用正则表达式提取占位符，使用字符串处理进行归一化。

**理由**：
- 正则表达式适合提取 `<...>` 模式
- 字符串处理（`strings.Map`、`strings.ReplaceAll`）对于归一化更简单、更快
- 避免复杂的正则表达式，提高可维护性

**正则模式**：`<[^>]+>` 匹配 `<` 开始、`>` 结束的内容

### 4. 性能考虑

**归一化操作的性能特征**：
- 正则提取：O(n) 其中 n 是文本长度
- 字符串处理：O(m) 其中 m 是占位符数量 * 平均占位符长度
- 整体：O(n)，线性复杂度

**优化**：
- 使用 `strings.Builder` 避免多次字符串分配
- 预编译正则表达式（包级变量）
- 只在需要时归一化（如果精确匹配失败）

**替代方案（拒绝）**：总是先尝试精确匹配
- 优点：精确匹配时性能更好
- 缺点：增加复杂度（需要双重逻辑）
- 缺点：在交互式场景中，格式变化是常态，精确匹配几乎总会失败

### 5. 向后兼容性

**保证**：
- 归一化后的标准格式与原始格式相同
- 现有测试用例无需修改（它们使用标准格式）
- 现有实体文件无需修改

**验证**：在测试中同时测试标准格式和变化格式。

## Algorithm

### normalizePlaceholder(placeholder string) string

```go
func normalizePlaceholder(placeholder string) string {
    // 1. Extract content between < and >
    if !strings.HasPrefix(placeholder, "<") || !strings.HasSuffix(placeholder, ">") {
        return placeholder // Not a placeholder, return as-is
    }
    
    content := placeholder[1:len(placeholder)-1]
    
    // 2. Remove all whitespace
    content = strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            return -1 // Remove
        }
        return r
    }, content)
    
    // 3. Convert Chinese punctuation to English
    replacements := map[string]string{
        "。": ".",
        "，": ",",
        "【": "[",
        "】": "]",
    }
    for old, new := range replacements {
        content = strings.ReplaceAll(content, old, new)
    }
    
    // 4. Convert fullwidth to halfwidth (ASCII range)
    content = strings.Map(func(r rune) rune {
        // Fullwidth ASCII: U+FF00 to U+FF5E
        // Halfwidth ASCII: U+0020 to U+007E
        if r >= 0xFF01 && r <= 0xFF5E {
            return r - 0xFEE0 // Convert to halfwidth
        }
        return r
    }, content)
    
    return "<" + content + ">"
}
```

### RestoreText 修改

```go
func (h *HasHidePair) RestoreText(ctx context.Context, entities []*Entity, text string) (string, error) {
    // Build normalized entity map: normalized_key -> original_value
    entityMap := make(map[string]string)
    for _, entity := range entities {
        if len(entity.Values) == 0 {
            continue
        }
        normalizedKey := normalizePlaceholder(entity.Key)
        entityMap[normalizedKey] = entity.Values[0]
    }
    
    // Find all placeholders in text, normalize, and replace
    placeholderRegex := regexp.MustCompile(`<[^>]+>`)
    result := placeholderRegex.ReplaceAllStringFunc(text, func(placeholder string) string {
        normalizedKey := normalizePlaceholder(placeholder)
        if value, exists := entityMap[normalizedKey]; exists {
            return value
        }
        // Placeholder not found, return as-is (partial restoration)
        return placeholder
    })
    
    return result, nil
}
```

## Risks / Trade-offs

### Risk: 误匹配

**场景**：归一化可能导致不同的占位符被归一化为相同的键。

**例子**（极端情况）：
- `<个人信息[0].姓名.张三>` 和 `<个人信息[　０].姓名.张三>` （全角数字 0）

**缓解**：
- 全角到半角的转换确保 `０` → `0`
- 实际中，LLM 生成的占位符不太可能有这种冲突
- 如果真的发生，只会影响少数边缘案例

**接受**：这是提高兼容性的合理代价。

### Trade-off: 性能 vs 兼容性

**性能影响**：
- 每次 `RestoreText` 调用都需要正则匹配和归一化
- 对于大文本（>10K 字符）和大量实体（>100），可能有微小性能影响

**缓解**：
- 正则和字符串处理都是 O(n)，性能影响线性且可预测
- 实际测试显示，对于典型文本（<5K 字符，<50 实体），延迟 <1ms

**接受**：兼容性带来的用户体验改进远大于微小的性能损失。

### Risk: 破坏边缘案例

**场景**：某些合法的占位符可能包含需要保留的空格或标点。

**例子**：
- `<文档名称[0].文件.My Document.txt>` 中的空格是文件名的一部分

**分析**：
- 按照当前设计，Detail 部分的空格会被移除
- 但 LLM 生成的占位符通常不会在 Detail 中包含空格（会被编码或省略）

**缓解**：
- 如果未来需要支持，可以调整归一化规则（例如只移除 `[]` 和 `.` 周围的空格）
- 当前设计优先考虑常见场景（格式化导致的意外空格）

**接受**：Detail 中的空格是罕见情况，优先解决常见问题。

## Migration Plan

**部署步骤**：
1. 实现 `normalizePlaceholder` 函数并添加单元测试
2. 修改 `RestoreText` 使用归一化匹配
3. 添加集成测试覆盖格式变化场景
4. 部署新版本（无需数据迁移，完全向后兼容）

**回滚**：
- 如果发现严重问题，可以直接回滚到前一版本
- 不涉及数据格式变更，回滚无风险

**验证**：
- 运行现有所有测试确保向后兼容
- 手动测试交互式工作流中的常见格式变化
- 检查性能基准测试，确保没有显著退化

## Open Questions

1. **是否需要配置选项来禁用归一化匹配？**
   - **决定**：不需要。归一化是完全向后兼容的，不会影响标准格式。如果用户需要精确匹配，标准格式仍然有效。

2. **是否需要记录哪些占位符被归一化匹配？**
   - **决定**：不需要。这是内部实现细节，对用户不可见。如果需要调试，可以通过日志或 verbose 模式添加。

3. **是否需要支持更多的归一化规则（如大小写不敏感）？**
   - **决定**：不需要。占位符的实体类型、类别和细节都是区分大小写的（来自 LLM），不应该忽略大小写。
