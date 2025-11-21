# Refactor Anonymizer to Interface - Design

## Context
当前 `Anonymizer` 是一个具体的 struct，包含 LLM 客户端和 prompt 模板，使用 `<<<PAIR>>>` 格式分隔脱敏文本和实体映射。虽然 Web handlers 已经定义了局部接口（好的设计），但核心包仍然暴露具体类型，限制了扩展性。

**当前架构**：
```
cmd/inu/commands/*.go
    ↓ (依赖具体类型)
pkg/anonymizer.Anonymizer (struct)
    ↓
CloudWeGo Eino LLM
```

**目标架构**：
```
cmd/inu/commands/*.go
pkg/web/server.go
    ↓ (依赖接口)
pkg/anonymizer.Anonymizer (interface)
    ↑ (实现)
pkg/anonymizer.HasHidePair (struct)
    ↓
CloudWeGo Eino LLM
```

## Goals / Non-Goals

**Goals**：
- 定义清晰的 `Anonymizer` 接口作为稳定契约
- 将当前实现重命名为 `HasHidePair`，体现其特定格式
- 所有调用方依赖接口而非具体实现
- 保持 100% 向后兼容（行为不变）
- 为未来扩展其他实现铺路

**Non-Goals**：
- 不添加新的脱敏策略（本次仅重构）
- 不修改 `<<<PAIR>>>` 格式或解析逻辑
- 不改变 Web API 或 CLI 的行为
- 不修改测试用例的逻辑（只调整构造调用）

## Decisions

### Decision 1: 接口位置和命名
**选择**：在 `pkg/anonymizer/anonymizer.go` 顶部定义 `Anonymizer` 接口

**理由**：
- 接口是包的核心抽象，应放在最显眼位置
- Go 约定：接口定义在使用方包，但此处作为核心能力定义在提供方包也合理
- 命名保持 `Anonymizer` 体现核心领域概念

**备选方案及拒绝理由**：
- 单独文件 `interface.go`：过度分割，接口简单不需要
- 放在调用方包：会导致循环依赖问题

### Decision 2: 实现类命名
**选择**：`HasHidePair` 替代原 `Anonymizer` struct

**命名依据**：
- **Has** 前缀：表示"具有...特性"的实现
- **HidePair**：隐藏分隔对（`<<<PAIR>>>`），体现实现细节
- 与接口名区分，避免混淆

**备选方案及拒绝理由**：
- `PairBasedAnonymizer`：冗长，不符合 Go 简洁风格
- `DefaultAnonymizer`：未体现实现特点
- `LLMAnonymizer`：过于宽泛，未来可能有其他 LLM 实现

### Decision 3: 接口方法定义
**选择**：接口包含三个核心方法

```go
type Anonymizer interface {
    AnonymizeText(ctx context.Context, types []string, text string) (string, []*Entity, error)
    AnonymizeTextStream(ctx context.Context, types []string, text string, writer io.Writer) ([]*Entity, error)
    RestoreText(ctx context.Context, entities []*Entity, text string) (string, error)
}
```

**理由**：
- 完整覆盖当前所有公开方法
- `AnonymizeText` 和 `AnonymizeTextStream` 是核心能力
- `RestoreText` 是必备的逆操作
- 内部辅助方法（如 `createAnonymizeMessages`）不暴露

**不包含的方法**：
- `createAnonymizeMessages`：内部实现细节
- `parseAnonymizeResponse`：包级私有函数
- `parseAnonymizeEntities`：包级私有函数

### Decision 4: 构造函数返回类型
**选择**：返回接口类型 `Anonymizer`

```go
func NewHashHidePair(chatModel model.BaseChatModel) (Anonymizer, error) {
    return &HasHidePair{...}, nil
}
```

**理由**：
- 调用方依赖接口，不关心具体实现
- 符合依赖倒置原则
- 便于未来替换实现或添加工厂方法

**备选方案及拒绝理由**：
- 返回 `*HasHidePair`：暴露具体类型，违背接口设计初衷
- 继续用 `New` 命名：不清楚返回的是哪个实现

### Decision 5: Web Server 存储类型
**选择**：存储接口而非指针

```go
type Server struct {
    anonymizer anonymizer.Anonymizer  // 接口类型
}
```

**理由**：
- 接口本身就是引用语义，无需指针
- 符合 Go 接口使用惯例
- 更灵活，可存储任何实现

## Risks / Trade-offs

### Risk 1: 命名混淆
**风险**：`HasHidePair` 名称可能不直观

**缓解**：
- 添加清晰的文档注释说明含义
- 在构造函数注释中强调 `<<<PAIR>>>` 格式
- 代码审查时验证命名可读性

### Risk 2: 接口膨胀
**风险**：未来可能添加更多方法导致接口过大

**缓解**：
- 当前只包含核心三个方法，已足够精简
- 如需扩展可考虑接口组合（如 `Anonymizer` + `Streamer`）
- Go 接口是隐式实现，扩展不强制修改已有代码

### Trade-off: 抽象 vs 简单性
**取舍**：引入接口增加了一层抽象

**权衡考量**：
- **成本**：代码多一层类型定义，初学者需理解接口概念
- **收益**：扩展性显著提升，测试更友好，架构更清晰
- **结论**：收益远大于成本，符合长期演进需求

## Migration Plan

### 阶段 1: 定义接口和重命名（本提案）
- 添加 `Anonymizer` 接口定义
- 重命名 struct 为 `HasHidePair`
- 更新所有调用方
- 确保所有测试通过

### 阶段 2: 扩展能力（未来可选）
- 添加基于规则的实现：`RuleBasedAnonymizer`
- 添加基于字典的实现：`DictBasedAnonymizer`
- 添加组合实现：`HybridAnonymizer`

### 阶段 3: 工厂模式（未来可选）
- 添加工厂函数根据配置创建实现
- 支持运行时切换策略

## Open Questions

1. **是否需要更多接口方法？**
   - 当前：只有三个核心方法
   - 可选：添加配置方法（如 `SetEntityTypes`）
   - 决策：暂不需要，保持接口精简

2. **是否需要接口版本化？**
   - 当前：单一接口版本
   - 可选：`AnonymizerV2` 支持破坏性变更
   - 决策：暂不需要，接口稳定

3. **是否需要包级别的默认实例？**
   - 当前：每次显式创建
   - 可选：`var Default Anonymizer` 全局实例
   - 决策：不需要，避免全局状态

## Implementation Notes

### 接口定义建议
```go
// Anonymizer 定义文本敏感信息脱敏的核心接口。
// 实现此接口的类型应当能够：
//  1. 将原始文本中的敏感实体替换为占位符
//  2. 记录实体映射关系以支持还原
//  3. 支持流式输出以改善用户体验
type Anonymizer interface {
    // AnonymizeText 批量脱敏文本，返回完整结果
    AnonymizeText(ctx context.Context, types []string, text string) (string, []*Entity, error)
    
    // AnonymizeTextStream 流式脱敏文本，实时写入 writer
    AnonymizeTextStream(ctx context.Context, types []string, text string, writer io.Writer) ([]*Entity, error)
    
    // RestoreText 使用实体映射还原脱敏文本
    RestoreText(ctx context.Context, entities []*Entity, text string) (string, error)
}
```

### HasHidePair 注释建议
```go
// HasHidePair 是基于 <<<PAIR>>> 分隔符格式的 Anonymizer 实现。
// 它使用 LLM 生成脱敏文本和实体映射，响应格式为：
//   <脱敏文本>
//   <<<PAIR>>>
//   <JSON 映射>
// 此实现支持流式输出，在遇到 <<<PAIR>>> 标记前将 token 实时写入输出。
type HasHidePair struct {
    anonymizeTemplate *prompt.DefaultChatTemplate
    llm               model.BaseChatModel
}
```

### 测试注意事项
- 所有测试中的 `New()` 改为 `NewHashHidePair()`
- Mock 测试无需修改（已使用接口模式）
- 集成测试验证接口行为一致性
