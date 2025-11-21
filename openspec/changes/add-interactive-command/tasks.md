# Implementation Tasks

## 1. 设计验证和准备
- [x] 1.1 验证设计方案的可行性
  - [x] 测试 stdin 交互行为（EOF 和分隔符检测）
  - [x] 测试多次输入循环
  - [x] 验证 stderr 和 stdout 分离输出

## 2. 核心功能实现
- [x] 2.1 创建 `cmd/inu/commands/interactive.go`
  - [x] 定义命令行标志变量
    - [x] `interactiveFile`, `interactiveContent` (输入)
    - [x] `interactiveEntityTypes` (实体类型)
    - [x] `interactiveDelimiter` (自定义分隔符)
    - [x] `interactiveNoPrompt` (禁用详细提示)
  - [x] 实现 `NewInteractiveCmd()` 创建 Cobra 命令
  - [x] 实现 `runInteractive()` 主执行函数

## 3. 交互循环实现
- [x] 3.1 实现匿名化阶段
  - [x] 读取原始文本输入（复用 `cli.ReadInput`）
  - [x] 检查环境变量（LLM credentials）
  - [x] 初始化 LLM 和 Anonymizer
  - [x] 调用 `AnonymizeText()` 匿名化
  - [x] 输出匿名化文本到 stdout
  - [x] 在内存中保存实体信息
  
- [x] 3.2 实现 `printPrompt()` 辅助函数
  - [x] 根据 `--no-prompt` 标志决定详细程度
  - [x] 详细模式：输出完整使用说明
  - [x] 简短模式：仅输出 "Waiting for input..."
  - [x] 提示中包含分隔符信息（如果使用）
  
- [x] 3.3 实现 `readUntilDelimiter()` 辅助函数
  - [x] 使用 `bufio.Scanner` 逐行读取 stdin
  - [x] 检测 EOF 或自定义分隔符
  - [x] 分隔符行本身不包含在返回内容中
  - [x] 返回读取的完整文本（去除分隔符行）
  
- [x] 3.4 实现循环还原逻辑
  - [x] 在循环中调用 `readUntilDelimiter()`
  - [x] 空输入时跳过（继续等待下一次）
  - [x] 调用 `RestoreText()` 使用内存中的实体还原
  - [x] 输出还原后的文本到 stdout
  - [x] 输出 "Ready for next input..." 提示
  - [x] 继续循环直到用户 Ctrl+C

## 4. 命令注册
- [x] 4.1 修改 `cmd/inu/main.go`
  - [x] 导入 interactive 命令
  - [x] 在 root 命令中注册 `NewInteractiveCmd()`

## 5. 错误处理和验证
- [x] 5.1 添加输入验证
  - [x] 验证原始文本非空
  - [x] 优雅处理 Ctrl+C 中断
- [x] 5.2 添加错误提示
  - [x] 友好的错误信息
  - [x] 使用场景示例
- [x] 5.3 处理边界情况
  - [x] 用户输入 Ctrl+C 中断
  - [x] stdin 读取错误
  - [x] LLM API 调用失败

## 6. 测试实现
- [x] 6.1 创建 `cmd/inu/commands/interactive_test.go`
  - [x] 测试命令创建和参数解析
  - [x] 测试标志验证逻辑
- [x] 6.2 单元测试
  - [x] `TestDelimiterLogic_EOF` - 测试 EOF 检测
  - [x] `TestDelimiterLogic_Custom` - 测试自定义分隔符
  - [x] `TestNewInteractiveCmd` - 测试命令创建
- [x] 6.3 集成测试（基本测试覆盖）
  - [x] 命令行为逻辑测试（通过单元测试覆盖）

## 7. 文档更新
- [x] 7.1 更新 `README.md`
  - [x] 添加 `inu interactive` 命令说明
  - [x] 添加使用场景和示例
  - [x] 添加与 ChatGPT 交互的最佳实践
- [x] 7.2 更新命令帮助文本
  - [x] 编写清晰的 Short 和 Long 描述
  - [x] 为每个标志添加详细说明
  - [x] 添加使用示例到帮助文本
