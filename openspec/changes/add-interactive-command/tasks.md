# Implementation Tasks

## 1. 设计验证和准备
- [ ] 1.1 验证设计方案的可行性
  - [ ] 测试 stdin 交互行为（EOF 和分隔符检测）
  - [ ] 测试多次输入循环
  - [ ] 验证 stderr 和 stdout 分离输出

## 2. 核心功能实现
- [ ] 2.1 创建 `cmd/inu/commands/interactive.go`
  - [ ] 定义命令行标志变量
    - [ ] `interactiveFile`, `interactiveContent` (输入)
    - [ ] `interactiveEntityTypes` (实体类型)
    - [ ] `interactiveDelimiter` (自定义分隔符)
    - [ ] `interactiveNoPrompt` (禁用详细提示)
  - [ ] 实现 `NewInteractiveCmd()` 创建 Cobra 命令
  - [ ] 实现 `runInteractive()` 主执行函数

## 3. 交互循环实现
- [ ] 3.1 实现匿名化阶段
  - [ ] 读取原始文本输入（复用 `cli.ReadInput`）
  - [ ] 检查环境变量（LLM credentials）
  - [ ] 初始化 LLM 和 Anonymizer
  - [ ] 调用 `AnonymizeText()` 匿名化
  - [ ] 输出匿名化文本到 stdout
  - [ ] 在内存中保存实体信息
  
- [ ] 3.2 实现 `printPrompt()` 辅助函数
  - [ ] 根据 `--no-prompt` 标志决定详细程度
  - [ ] 详细模式：输出完整使用说明
  - [ ] 简短模式：仅输出 "Waiting for input..."
  - [ ] 提示中包含分隔符信息（如果使用）
  
- [ ] 3.3 实现 `readUntilDelimiter()` 辅助函数
  - [ ] 使用 `bufio.Scanner` 逐行读取 stdin
  - [ ] 检测 EOF 或自定义分隔符
  - [ ] 分隔符行本身不包含在返回内容中
  - [ ] 返回读取的完整文本（去除分隔符行）
  
- [ ] 3.4 实现循环还原逻辑
  - [ ] 在循环中调用 `readUntilDelimiter()`
  - [ ] 空输入时跳过（继续等待下一次）
  - [ ] 调用 `RestoreText()` 使用内存中的实体还原
  - [ ] 输出还原后的文本到 stdout
  - [ ] 输出 "Ready for next input..." 提示
  - [ ] 继续循环直到用户 Ctrl+C

## 4. 命令注册
- [ ] 4.1 修改 `cmd/inu/main.go`
  - [ ] 导入 interactive 命令
  - [ ] 在 root 命令中注册 `NewInteractiveCmd()`

## 5. 错误处理和验证
- [ ] 5.1 添加输入验证
  - [ ] 验证原始文本非空
  - [ ] 优雅处理 Ctrl+C 中断
- [ ] 5.2 添加错误提示
  - [ ] 友好的错误信息
  - [ ] 使用场景示例
- [ ] 5.3 处理边界情况
  - [ ] 用户输入 Ctrl+C 中断
  - [ ] stdin 读取错误
  - [ ] LLM API 调用失败

## 6. 测试实现
- [ ] 6.1 创建 `cmd/inu/commands/interactive_test.go`
  - [ ] 测试命令创建和参数解析
  - [ ] 测试标志验证逻辑
- [ ] 6.2 单元测试
  - [ ] `TestReadUntilDelimiter_EOF` - 测试 EOF 检测
  - [ ] `TestReadUntilDelimiter_CustomDelimiter` - 测试自定义分隔符
  - [ ] `TestPrintPrompt_Detailed` - 测试详细提示
  - [ ] `TestPrintPrompt_Concise` - 测试简短提示
- [ ] 6.3 集成测试
  - [ ] `TestInteractiveCommand_SingleInput` - 测试单次输入流程（mock LLM 和 stdin）
  - [ ] `TestInteractiveCommand_MultipleInputs` - 测试多次输入循环
  - [ ] `TestInteractiveCommand_EmptyInput` - 测试空输入跳过
  - [ ] `TestInteractiveCommand_CustomDelimiter` - 测试自定义分隔符

## 7. 文档更新
- [ ] 7.1 更新 `README.md`
  - [ ] 添加 `inu interactive` 命令说明
  - [ ] 添加使用场景和示例
  - [ ] 添加与 ChatGPT 交互的最佳实践
- [ ] 8.2 更新命令帮助文本
  - [ ] 编写清晰的 Short 和 Long 描述
  - [ ] 为每个标志添加详细说明
  - [ ] 添加使用示例到帮助文本
- [ ] 8.3 创建用户指南（可选）
  - [ ] 交互式使用教程
  - [ ] 管道集成示例
  - [ ] 故障排查指南

## 9. 手动测试和验证
- [ ] 9.1 基本功能测试
  - [ ] 测试从文件输入：`inu interactive -f test.txt`
  - [ ] 测试从内容输入：`inu interactive -c "张三"`
  - [ ] 测试从 stdin 输入：`echo "张三" | inu interactive`
  - [ ] 验证匿名化文本输出正确
  - [ ] 验证实体信息输出到 stderr
  - [ ] 验证交互等待和还原功能
- [ ] 9.2 模式测试
  - [ ] 测试嵌入实体模式：`inu interactive --embed-entities`
  - [ ] 测试还原模式：`inu interactive --restore-mode --entities e.yaml`
  - [ ] 测试自定义分隔符
  - [ ] 测试禁用提示：`--no-prompt`
- [ ] 9.3 错误场景测试
  - [ ] 测试无输入：`inu interactive` （应报错）
  - [ ] 测试还原模式缺少 --entities（应报错）
  - [ ] 测试无效实体文件格式
  - [ ] 测试处理后文本为空
  - [ ] 测试 LLM API 失败
- [ ] 9.4 实际使用场景测试
  - [ ] 测试与 ChatGPT 交互（手动复制粘贴）
  - [ ] 测试与本地 LLM 工具集成
  - [ ] 测试重定向输出：`inu interactive -f input.txt > output.txt 2> entities.json`
- [ ] 9.5 性能和体验测试
  - [ ] 测试大文本处理（1MB+）
  - [ ] 验证提示信息清晰度
  - [ ] 验证错误信息友好性

## 10. 集成和发布准备
- [ ] 10.1 运行完整测试套件：`make test`
- [ ] 10.2 验证构建成功：`make build`
- [ ] 10.3 检查代码格式：`gofmt`, `golangci-lint`
- [ ] 10.4 更新 CHANGELOG（如有）
- [ ] 10.5 验证与现有命令的兼容性
  - [ ] 确认 `anonymize` 命令不受影响
  - [ ] 确认 `restore` 命令不受影响
  - [ ] 确认 `web` 命令不受影响
