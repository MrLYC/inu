# Design Document: Anonymizer Unit Tests with Mocking

## Problem Statement
`pkg/anonymizer` 包依赖外部 LLM API（通过 CloudWeGo Eino 的 `model.BaseChatModel` 接口），导致编写单元测试时面临以下挑战：
1. **外部依赖**：测试需要真实的 API key 和网络连接
2. **不确定性**：LLM 响应可能不稳定，导致测试结果不可预测
3. **成本问题**：频繁测试会产生 API 调用费用
4. **速度问题**：网络请求会大幅降低测试执行速度

需要设计一个 mock 策略，在不依赖真实 LLM 的情况下测试核心逻辑。

## Solution Overview
使用 Go 的接口特性实现 `model.BaseChatModel` 的 mock 版本，通过依赖注入在测试中替换真实实现。

### Key Design Decisions

#### 1. Mock 实现方式
**选择：手写 mock 结构体**
- **理由**：
  - `model.BaseChatModel` 接口相对简单（主要是 `Generate` 方法）
  - 避免引入额外依赖（如 gomock、testify/mock）
  - 更容易控制和理解 mock 行为
  - 符合项目"最小依赖"原则

**替代方案（未采用）**：
- gomock：需要代码生成，增加复杂度
- testify/mock：额外依赖，对简单接口过度设计

#### 2. Mock 响应格式
**设计**：mock 需要模拟 LLM 返回的特定格式
```
<anonymized_text>
<<<PAIR>>>
{"<EntityType[ID].Category.Detail>": ["original_value"]}
```

**实现**：
```go
type mockChatModel struct {
    response      string  // 完整响应内容
    responseError error   // 模拟错误
}

// 辅助函数简化响应构造
func newMockAnonymizeResponse(anonymizedText string, mapping map[string][]string) string
```

#### 3. 测试数据策略
**选择：每个测试用例显式定义输入和期望输出**
- **理由**：
  - 测试意图清晰
  - 易于调试失败的测试
  - 避免测试之间的隐式依赖

**示例**：
```go
tests := []struct {
    name           string
    mockResponse   string
    inputTypes     []string
    inputText      string
    expectText     string
    expectEntities int
    expectError    bool
}{
    {
        name: "single entity anonymization",
        mockResponse: "mock response...",
        // ...
    },
}
```

#### 4. 覆盖范围
**重点测试**：
1. **核心功能**：AnonymizeText, RestoreText
2. **边界情况**：空输入、无匹配、错误格式
3. **错误处理**：API 失败、解析错误

**不测试**：
- LLM 本身的行为（外部依赖）
- 真实的网络交互（集成测试范围）

## Implementation Details

### Mock LLM 结构
```go
type mockChatModel struct {
    response      *schema.Message // 模拟成功响应
    responseError error            // 模拟失败
}

func (m *mockChatModel) Generate(ctx context.Context, messages []*schema.Message, opts ...model.Option) (*schema.Message, error) {
    if m.responseError != nil {
        return nil, m.responseError
    }
    return m.response, nil
}
```

### 响应构造辅助
```go
// 构造标准脱敏响应
func newMockAnonymizeResponse(anonymizedText string, mapping map[string][]string) *schema.Message {
    mappingJSON, _ := json.Marshal(mapping)
    content := fmt.Sprintf("%s\n<<<PAIR>>>\n%s", anonymizedText, mappingJSON)
    return &schema.Message{Content: content}
}

// 构造错误响应
func newMockError(errMsg string) error {
    return errors.New(errMsg)
}
```

### 测试模式示例
```go
func TestAnonymizeText_Success(t *testing.T) {
    // 准备 mock
    mockLLM := &mockChatModel{
        response: newMockAnonymizeResponse(
            "<个人信息[0].姓名.全名> lives in Beijing",
            map[string][]string{
                "<个人信息[0].姓名.全名>": {"张三"},
            },
        ),
    }

    // 创建脱敏器
    anon, _ := New(mockLLM)

    // 执行测试
    result, entities, err := anon.AnonymizeText(ctx, []string{"个人信息"}, "张三 lives in Beijing")

    // 验证结果
    assert.NoError(t, err)
    assert.Equal(t, "<个人信息[0].姓名.全名> lives in Beijing", result)
    assert.Len(t, entities, 1)
}
```

## Trade-offs

### 优势
1. **快速可靠**：无网络依赖，测试秒级完成
2. **可重复**：每次运行结果一致
3. **零成本**：不产生 API 调用费用
4. **易维护**：mock 逻辑简单，与接口同步更新容易

### 限制
1. **不测试真实集成**：无法发现 LLM API 变化导致的问题
2. **假设响应格式**：如果 LLM 响应格式变化，测试可能无法捕获
3. **需要手动同步**：接口变更时需要更新 mock

### 缓解措施
- 保留现有的手动测试文件（`test_e2e.sh`）用于集成测试
- 在 CI 中添加可选的集成测试步骤（需要 API key）
- 定期人工验证真实 LLM 交互

## Alternatives Considered

### 方案 A：使用 gomock
- **优点**：自动生成，类型安全
- **缺点**：增加构建复杂度，代码生成步骤
- **决策**：不采用，接口太简单不值得

### 方案 B：使用 testify/mock
- **优点**：功能丰富，API 友好
- **缺点**：外部依赖，对简单场景过度
- **决策**：不采用，手写更清晰

### 方案 C：使用真实 API 的测试模式
- **优点**：测试最真实
- **缺点**：需要 API key，不稳定，慢
- **决策**：仅用于集成测试，不用于单元测试

## Success Criteria
1. 所有核心函数有测试覆盖（AnonymizeText, RestoreText, New）
2. 测试覆盖率 ≥ 80%
3. 测试在无网络环境下可运行
4. 测试执行时间 < 1 秒
5. CI 中测试稳定通过
