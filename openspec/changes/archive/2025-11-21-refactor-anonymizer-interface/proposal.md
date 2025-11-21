# Refactor Anonymizer to Interface

## Why
当前 `Anonymizer` 是一个具体的 struct 实现，直接依赖 LLM 和特定的响应格式（`<<<PAIR>>>`）。这种设计存在以下问题：

**现有问题**：
1. **扩展性差**：无法轻松添加其他脱敏策略（如基于规则、基于字典等）
2. **测试困难**：Web handlers 和 CLI 直接依赖具体实现，难以进行单元测试
3. **命名不清晰**：当前实现紧密耦合了 `<<<PAIR>>>` 格式，但名称未体现这一特点
4. **违反依赖倒置原则**：高层模块（handlers, commands）依赖低层实现细节

**改进后的效果**：
- 接口抽象：`Anonymizer` 作为接口定义核心能力
- 实现分离：`HasHidePair` 作为基于 `<<<PAIR>>>` 格式的具体实现
- 易于扩展：可添加其他实现（如 `RuleBasedAnonymizer`）
- 测试友好：handlers 依赖接口，易于 mock

## What Changes
- **重构** `Anonymizer` 为接口类型，定义核心方法签名
- **重命名** 当前实现为 `HasHidePair` struct
- **更新** 构造函数：`NewHashHidePair(llm)` 创建具体实现
- **调整** 所有调用方：CLI commands, Web server, handlers
- **保持** 接口兼容：方法签名不变，只是类型层级调整

**接口定义**：
```go
// Anonymizer 定义文本脱敏的核心接口
type Anonymizer interface {
    AnonymizeText(ctx context.Context, types []string, text string) (string, []*Entity, error)
    AnonymizeTextStream(ctx context.Context, types []string, text string, writer io.Writer) ([]*Entity, error)
    RestoreText(ctx context.Context, entities []*Entity, text string) (string, error)
}

// HasHidePair 是基于 <<<PAIR>>> 格式的 Anonymizer 实现
type HasHidePair struct {
    anonymizeTemplate *prompt.DefaultChatTemplate
    llm               model.BaseChatModel
}
```

**调用方变更**：
- `anonymizer.New(llm)` → `anonymizer.NewHashHidePair(llm)`
- 返回类型：`*anonymizer.Anonymizer` → `anonymizer.Anonymizer` (接口)
- Web Server 存储：`*anonymizer.Anonymizer` → `anonymizer.Anonymizer`
- Handlers 已使用接口，无需修改（已正确设计）

## Impact
- **影响的 specs**: 无需修改 specs（接口行为不变）
- **影响的代码**:
  - `pkg/anonymizer/anonymizer.go` - 接口定义 + struct 重命名
  - `cmd/inu/commands/anonymize.go` - 调用 `NewHashHidePair`
  - `cmd/inu/commands/restore.go` - 调用 `NewHashHidePair`
  - `cmd/inu/commands/web.go` - 调用 `NewHashHidePair`
  - `pkg/web/server.go` - 存储接口类型
  - `README.md` - 示例代码更新
- **破坏性变更**: 否
  - 对外 API 行为完全一致
  - 只是内部重构，不影响功能
  - 测试继续通过（mock 已经使用接口模式）
- **向后兼容性**: 完全兼容
  - 方法签名不变
  - 行为不变
  - 只是类型系统重构

## Implementation Notes
1. 在 `anonymizer.go` 顶部定义 `Anonymizer` 接口
2. 重命名 `type Anonymizer struct` → `type HasHidePair struct`
3. 更新所有方法接收者：`(a *Anonymizer)` → `(h *HasHidePair)`
4. 重命名构造函数：`New(llm)` → `NewHashHidePair(llm)`
5. 返回类型改为接口：`*Anonymizer` → `Anonymizer`
6. 更新所有调用方（commands, web server）
7. 运行完整测试确保兼容性
