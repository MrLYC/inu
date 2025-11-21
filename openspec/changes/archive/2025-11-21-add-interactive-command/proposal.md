# Add Interactive Command

## Why
当前 Inu 的 `anonymize` 和 `restore` 命令是独立运行的，用户需要：
1. 运行 `anonymize` 保存实体文件
2. 手动处理脱敏文本（如使用 ChatGPT 总结、翻译等）
3. 运行 `restore` 还原处理后的文本

这个流程存在以下问题：
- 需要管理中间文件（entities.yaml）
- 多步骤操作容易出错
- 不适合交互式场景（复制粘贴到外部工具处理）

添加交互式命令可以：
- **简化工作流**：一个命令完成整个流程
- **无需中间文件**：实体信息在内存中保持
- **支持多次处理**：可以多次输入不同的处理结果
- **提升用户体验**：适合与 ChatGPT 等外部工具配合使用

**典型使用场景**：
```bash
# 场景 1: 与 ChatGPT 交互
$ inu interactive -f sensitive-report.txt
[复制脱敏文本到 ChatGPT 请求总结]
[复制 ChatGPT 的总结粘贴回终端]
[得到还原后的总结]
[继续请求翻译...]

# 场景 2: 自定义分隔符
$ inu interactive -c "张三在 ABC 公司工作" --delimiter "END"
[处理文本]
END
[得到还原文本]
```

## What Changes
- 添加 `inu interactive` 子命令，提供交互式脱敏和还原流程
- 支持单一交互模式：脱敏 → 循环等待输入 → 还原 → 继续等待
- 实体信息保存在内存中，不输出到文件或标准流
- 支持自定义输入分隔符（类似 heredoc）

**新增命令**:
```bash
# 基本交互模式
inu interactive [--file <input>] [--content <text>] [--entity-types <types>]

# 自定义分隔符
inu interactive -f input.txt --delimiter "END"

# 精简提示
inu interactive -c "文本" --no-prompt
```

**工作流程**:
```
1. 读取原始文本（从 --file / --content / stdin）
2. 调用 LLM 进行脱敏
3. 输出脱敏文本到 stdout
4. 在 stderr 显示详细使用提示
5. 循环：
   a. 等待从 stdin 读取处理后的文本
   b. 遇到分隔符或 EOF 时处理
   c. 使用内存中的实体信息还原文本
   d. 输出还原后的文本到 stdout
   e. 继续等待下一次输入
6. 用户 Ctrl+C 退出
```

## Impact
- **影响的 specs**: 需要修改 `cli` spec，添加 interactive 命令相关需求
- **影响的代码**:
  - 新增 `cmd/inu/commands/interactive.go` - Interactive 命令实现
  - 修改 `cmd/inu/main.go` - 注册 interactive 命令
- **依赖变更**: 无（复用现有 anonymizer 和 CLI 工具）
- **不破坏现有功能**: 
  - `anonymize` 和 `restore` 命令保持不变
  - 新增独立的 `interactive` 命令
