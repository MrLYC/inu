# Implementation Tasks

## 1. Define Anonymizer Interface
- [ ] 1.1 在 `pkg/anonymizer/anonymizer.go` 文件顶部添加 `Anonymizer` 接口定义
  - 定义 `AnonymizeText` 方法签名
  - 定义 `AnonymizeTextStream` 方法签名
  - 定义 `RestoreText` 方法签名
  - 添加接口文档注释说明设计意图

## 2. Rename Struct to HasHidePair
- [ ] 2.1 重命名 struct 定义
  - `type Anonymizer struct` → `type HasHidePair struct`
  - 保持字段不变（`anonymizeTemplate`, `llm`）
  - 添加注释说明 HasHidePair 是基于 `<<<PAIR>>>` 格式的实现
- [ ] 2.2 更新所有方法接收者
  - `(a *Anonymizer)` → `(h *HasHidePair)`
  - 方法体内的 `a.` 引用改为 `h.`
  - 包括：`createAnonymizeMessages`, `AnonymizeTextStream`, `AnonymizeText`, `RestoreText`

## 3. Rename Constructor Function
- [ ] 3.1 重命名构造函数
  - `func New(chatModel)` → `func NewHashHidePair(chatModel)`
  - 返回类型：`(*Anonymizer, error)` → `(Anonymizer, error)` (接口类型)
  - 实际返回：`&HasHidePair{...}` (实现了接口)
- [ ] 3.2 更新构造函数文档
  - 说明返回的是 Anonymizer 接口实现
  - 强调基于 `<<<PAIR>>>` 分隔格式

## 4. Update CLI Commands
- [ ] 4.1 更新 `cmd/inu/commands/anonymize.go`
  - 导入保持不变
  - 调用 `anonymizer.NewHashHidePair(llm)` 替代 `anonymizer.New(llm)`
  - 变量类型保持 `anon` (Go 接口无需显式声明类型)
- [ ] 4.2 更新 `cmd/inu/commands/restore.go`
  - 调用 `anonymizer.NewHashHidePair(llm)` 替代 `anonymizer.New(llm)`
- [ ] 4.3 更新 `cmd/inu/commands/web.go`
  - 调用 `anonymizer.NewHashHidePair(llm)` 替代 `anonymizer.New(llm)`

## 5. Update Web Server
- [ ] 5.1 更新 `pkg/web/server.go`
  - 字段类型：`anonymizer *anonymizer.Anonymizer` → `anonymizer anonymizer.Anonymizer` (接口，去掉指针)
  - 构造函数参数：`anon *anonymizer.Anonymizer` → `anon anonymizer.Anonymizer`
  - 字段赋值保持：`anonymizer: anon`
- [ ] 5.2 验证 handlers 无需修改
  - `handlers/anonymize.go` 已使用接口 `Anonymizer` interface
  - `handlers/restore.go` 已使用接口 `Restorer` interface
  - Mock 测试已正确使用接口模式

## 6. Update Documentation
- [ ] 6.1 更新 `README.md`
  - 示例代码：`anonymizer.New(llm)` → `anonymizer.NewHashHidePair(llm)`
  - 添加接口说明（如果需要）
- [ ] 6.2 更新代码注释
  - 在接口定义处说明设计意图
  - 在 HasHidePair struct 处说明实现细节

## 7. Testing and Validation
- [ ] 7.1 运行单元测试
  - `go test ./pkg/anonymizer -v` 确保核心逻辑测试通过
  - 测试文件中的 `New()` 调用改为 `NewHashHidePair()`
- [ ] 7.2 运行完整测试套件
  - `go test ./...` 确保所有包测试通过
  - Web API 测试、CLI 测试都应通过
- [ ] 7.3 编译验证
  - `make build` 确保编译成功
  - 检查是否有未更新的引用
- [ ] 7.4 手动功能测试
  - CLI anonymize 命令功能正常
  - CLI restore 命令功能正常
  - Web API 功能正常

## Dependencies
- Task 2 depends on 1 (接口定义后才能重命名 struct)
- Task 3 depends on 2 (struct 重命名后才能更新构造函数)
- Task 4,5 depends on 3 (构造函数准备好后才能更新调用方)
- Task 7 depends on all above (所有代码更新后才能测试)

## Validation Criteria
- [ ] 所有单元测试通过
- [ ] 所有集成测试通过
- [ ] 代码编译无错误无警告
- [ ] CLI 功能手动测试通过
- [ ] Web API 功能手动测试通过
- [ ] 代码审查：命名清晰、接口设计合理

## Rollback Plan
如果出现问题，可以简单回退所有命名修改：
1. `NewHashHidePair` → `New`
2. `HasHidePair` → `Anonymizer` (struct)
3. 移除接口定义
4. 恢复指针类型引用
