# Add CLI Commands

## Why
当前 Inu 只是一个硬编码的 demo 程序，无法被用户实际使用。需要实现完整的命令行界面（CLI），使其成为一个可用的工具，支持从多种输入源读取内容、灵活配置实体类型、将结果输出到终端或文件。

## What Changes
- 实现两个 CLI 子命令：
  - `inu anonymize`: 匿名化文本中的敏感信息
  - `inu restore`: 还原匿名化的文本
- 支持多种输入方式：文件（`--file`）、命令行参数（`--content`）、标准输入
- 支持多种输出方式：终端打印（`--print`）、文件输出（`--output`）
- 支持实体信息的输出和读取（YAML 格式）
- 使用 cobra + viper 实现 CLI 框架和配置管理
- 保留原有的核心匿名化逻辑，只改造入口层

## Impact
- 影响的 specs: 新增 `cli` capability
- 影响的代码:
  - `cmd/inu/main.go`: 完全重写，从 demo 改为 CLI 命令
  - 新增 `cmd/inu/commands/`: 命令实现（anonymize.go, restore.go）
  - 可选：新增 `pkg/cli/`: CLI 辅助函数（输入/输出处理）
  - `pkg/anonymizer/`: 保持不变，无需修改
- 依赖变更:
  - 新增 CLI 框架：github.com/spf13/cobra
  - 新增配置管理：github.com/spf13/viper（处理 YAML 配置）
