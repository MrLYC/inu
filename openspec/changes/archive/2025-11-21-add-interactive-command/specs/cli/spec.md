# cli Specification Delta

## ADDED Requirements

### Requirement: 交互式命令
系统 SHALL 提供 `interactive` 子命令来执行交互式的脱敏和还原流程。

#### Scenario: 基本交互式流程
- **WHEN** 用户执行 `inu interactive -f input.txt`
- **THEN** 系统应该：
  1. 读取并脱敏 input.txt 的内容
  2. 将脱敏文本输出到 stdout
  3. 在 stderr 显示详细使用提示（如何输入、如何结束）
  4. 等待用户从 stdin 输入处理后的文本
  5. 用户按 Ctrl+D 后，使用内存中的实体信息还原文本
  6. 将还原后的文本输出到 stdout
  7. 继续等待下一次输入（循环）

#### Scenario: 从命令行参数输入
- **WHEN** 用户执行 `inu interactive -c "张三的电话是 13800138000"`
- **THEN** 系统应该使用提供的文本内容执行交互式流程

#### Scenario: 指定实体类型
- **WHEN** 用户执行 `inu interactive -c "张三在 ABC 公司工作" --entity-types "个人信息"`
- **THEN** 系统应该只识别和脱敏 "个人信息" 类型的实体
- **AND** 忽略其他类型（如 "组织机构"）

#### Scenario: 自定义分隔符
- **WHEN** 用户执行 `inu interactive -f input.txt --delimiter "END"`
- **THEN** 系统应该在提示中说明使用 END 作为分隔符
- **AND** 当用户输入单独的 "END" 行时，触发还原处理
- **AND** 分隔符行本身不包含在处理文本中
- **AND** 处理完成后继续等待下一次输入

#### Scenario: 多次输入处理
- **WHEN** 用户在交互模式下输入多次处理后的文本
- **THEN** 系统应该：
  1. 第一次输入 + 分隔符/EOF → 还原并输出
  2. 提示 "Ready for next input..."
  3. 第二次输入 + 分隔符/EOF → 再次还原并输出
  4. 持续循环直到用户 Ctrl+C 退出

#### Scenario: 详细提示信息
- **WHEN** 命令启动并完成脱敏
- **THEN** stderr 应该输出类似：
  ```
  === Anonymization Complete ===
  The text above has been anonymized.
  You can now:
  1. Copy the anonymized text
  2. Process it externally (e.g., paste to ChatGPT)
  3. Paste the processed text back here
  4. Press Ctrl+D (or type 'END' if --delimiter is set) to restore

  Waiting for your input...
  ```

#### Scenario: 禁用详细提示信息
- **WHEN** 用户执行 `inu interactive -f input.txt --no-prompt`
- **THEN** 系统应该只输出简短提示："Waiting for input..."
- **AND** 不显示详细的使用说明
- **AND** 仍然等待 stdin 输入

#### Scenario: 空输入处理
- **WHEN** 用户在等待输入阶段直接输入分隔符或 EOF 而不输入任何内容
- **THEN** 系统应该忽略本次输入
- **AND** 继续等待下一次输入
- **AND** 不报错

#### Scenario: 用户中断（Ctrl+C）
- **GIVEN** 命令正在等待用户输入
- **WHEN** 用户按 Ctrl+C
- **THEN** 系统应该立即退出
- **AND** 退出状态码为非零（中断）

### Requirement: 输出流控制
系统 SHALL 正确分离 stdout 和 stderr 输出，支持重定向。

#### Scenario: 分离输出流
- **WHEN** 用户执行 `inu interactive -f input.txt > output.txt 2> prompts.log`
- **THEN** 脱敏文本和还原后文本应该写入 output.txt
- **AND** 提示信息应该写入 prompts.log
- **AND** 命令仍然等待 stdin 输入（终端交互）

### Requirement: 交互式命令错误处理和用户体验
系统 SHALL 在交互式命令中提供清晰的错误信息和友好的交互提示。

#### Scenario: 环境变量未配置
- **WHEN** 用户执行 `inu interactive -f input.txt` 但未设置 OPENAI_API_KEY
- **THEN** 系统应该显示友好的错误信息
- **AND** 说明需要配置的环境变量
- **AND** 退出并显示非零状态码

#### Scenario: LLM API 调用失败
- **WHEN** 脱敏过程中 LLM API 调用失败
- **THEN** 系统应该显示错误："Error: Failed to anonymize text: <具体错误>"
- **AND** 不进入等待输入阶段
- **AND** 退出并显示非零状态码

#### Scenario: 无输入内容
- **WHEN** 用户执行 `inu interactive`（无 --file, --content, 无 stdin）
- **THEN** 系统应该返回错误："Error: No input provided. Use --file, --content, or stdin."
- **AND** 显示使用示例

#### Scenario: 还原阶段错误处理
- **WHEN** 还原过程中发生错误（如占位符格式损坏）
- **THEN** 系统应该尽力还原能匹配的部分
- **AND** 在 stderr 输出警告："Warning: Some placeholders could not be restored"
- **AND** 仍然输出部分还原的文本到 stdout

#### Scenario: 交互提示清晰性
- **WHEN** 命令等待用户输入时
- **THEN** 提示信息应该清楚说明：
  - 当前正在等待输入
  - 如何结束输入（Ctrl+D）
  - 可以粘贴多行文本
- **AND** 提示应该输出到 stderr，不干扰 stdout

### Requirement: 命令帮助和文档
系统 SHALL 提供完整的命令帮助信息和使用示例。

#### Scenario: 显示 pipe 命令帮助
- **WHEN** 用户执行 `inu interactive --help`
- **THEN** 系统应该显示：
  - 命令简短描述
  - 详细描述和使用场景
  - 所有可用标志及说明
  - 使用示例（至少 2-3 个）
  - 与 anonymize/restore 命令的对比

#### Scenario: 帮助文本包含典型用例
- **WHEN** 用户查看 `inu interactive --help`
- **THEN** 帮助文本应该包含：
  - 基本交互式使用示例
  - 与 LLM 工具集成示例
  - 还原模式使用示例
  - 输出重定向示例

#### Scenario: 错误提示包含帮助链接
- **WHEN** 命令因参数错误失败
- **THEN** 错误信息应该提示："Run 'inu interactive --help' for usage information."
