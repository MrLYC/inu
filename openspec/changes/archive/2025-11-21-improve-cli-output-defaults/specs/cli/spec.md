# cli Specification

## Purpose
定义 Inu CLI 命令行接口的输出行为规范。

## MODIFIED Requirements

### Requirement: 命令输出行为
系统 SHALL 默认将主要输出写入 stdout，将诊断信息写入 stderr。

#### Scenario: 默认输出到 stdout - anonymize
- **WHEN** 用户执行 `inu anonymize --file input.txt`
- **THEN** 系统应该将匿名化后的文本输出到 stdout
- **AND** 不需要指定 `--print` 参数

#### Scenario: 默认输出到 stdout - restore
- **WHEN** 用户执行 `inu restore --file input.txt --entities e.yaml`
- **THEN** 系统应该将还原后的文本输出到 stdout
- **AND** 不需要指定 `--print` 参数

#### Scenario: 实体信息输出到 stderr
- **WHEN** 用户执行 `inu anonymize --file input.txt`
- **THEN** 系统应该将识别到的实体信息输出到 stderr
- **AND** 格式为 `key: values`（每行一个实体）
- **AND** 不干扰 stdout 的主要输出

#### Scenario: 禁用所有输出
- **WHEN** 用户执行 `inu anonymize --file input.txt --no-print`
- **THEN** 系统不应该向 stdout 输出匿名化文本
- **AND** 系统不应该向 stderr 输出实体信息
- **AND** 进度信息仍然输出到 stderr

#### Scenario: 同时输出到文件和终端
- **WHEN** 用户执行 `inu anonymize --file input.txt --output result.txt`
- **THEN** 系统应该将匿名化文本写入 result.txt
- **AND** 同时将匿名化文本输出到 stdout
- **AND** 将实体信息输出到 stderr

#### Scenario: 只输出到文件
- **WHEN** 用户执行 `inu anonymize --file input.txt --output result.txt --no-print`
- **THEN** 系统应该将匿名化文本写入 result.txt
- **AND** 不向 stdout 输出任何内容
- **AND** 不向 stderr 输出实体信息
- **AND** 进度信息仍然输出到 stderr

#### Scenario: 管道操作 - 只传递主输出
- **WHEN** 用户执行 `echo "张三" | inu anonymize | grep "个人信息"`
- **THEN** 系统应该将匿名化文本传递给下一个命令
- **AND** 实体信息不会干扰管道数据流（因为在 stderr）

#### Scenario: 管道操作 - 合并输出流
- **WHEN** 用户执行 `echo "张三" | inu anonymize 2>&1 | grep "个人信息"`
- **THEN** 系统应该将 stdout 和 stderr 合并传递给下一个命令
- **AND** 用户可以在管道中处理实体信息

#### Scenario: 重定向 - 分离输出流
- **WHEN** 用户执行 `inu anonymize -f input.txt > output.txt 2> entities.log`
- **THEN** 系统应该将匿名化文本写入 output.txt
- **AND** 将实体信息和进度信息写入 entities.log

#### Scenario: 重定向 - 只要主输出
- **WHEN** 用户执行 `inu anonymize -f input.txt 2>/dev/null`
- **THEN** 系统应该将匿名化文本输出到 stdout
- **AND** 丢弃所有 stderr 输出（实体信息、进度信息）

#### Scenario: 重定向 - 只要实体信息
- **WHEN** 用户执行 `inu anonymize -f input.txt 1>/dev/null`
- **THEN** 系统应该将实体信息输出到 stderr
- **AND** 丢弃 stdout 输出（匿名化文本）

### Requirement: 进度信息输出
系统 SHALL 将进度和状态信息输出到 stderr，不受 `--no-print` 影响。

#### Scenario: 显示进度信息
- **WHEN** 用户执行 `inu anonymize --file input.txt`
- **THEN** 系统应该在 stderr 输出：
  - "Initializing LLM client..."
  - "Anonymizing text..."
  - "Anonymization complete"
- **AND** 这些信息不会干扰 stdout 的主要输出

#### Scenario: 进度信息不受 --no-print 影响
- **WHEN** 用户执行 `inu anonymize --file input.txt --no-print`
- **THEN** 系统仍然应该在 stderr 输出进度信息
- **AND** 用户可以通过 `2>/dev/null` 单独禁用进度信息

### Requirement: 命令参数
系统 SHALL 提供 `--no-print` 参数来控制输出行为。

#### Scenario: --no-print 参数存在于 anonymize 命令
- **WHEN** 用户执行 `inu anonymize --help`
- **THEN** 帮助信息应该包含 `--no-print` 参数说明
- **AND** 不应该包含 `--print` 参数
- **AND** 不应该包含 `--print-entities` 参数

#### Scenario: --no-print 参数存在于 restore 命令
- **WHEN** 用户执行 `inu restore --help`
- **THEN** 帮助信息应该包含 `--no-print` 参数说明
- **AND** 不应该包含 `--print` 参数

## REMOVED Requirements

### Requirement: --print 参数（已移除）
~~系统 SHALL 提供 `--print` 参数来控制是否输出到 stdout。~~

**移除原因**: 默认输出到 stdout 更符合 Unix 工具惯例，无需显式参数。

#### Scenario: ~~使用 --print 输出~~（已移除）
~~- **WHEN** 用户执行 `inu anonymize --file input.txt --print`~~
~~- **THEN** 系统应该将匿名化文本输出到 stdout~~

### Requirement: --print-entities 参数（已移除）
~~系统 SHALL 提供 `--print-entities` 参数来控制是否输出实体信息。~~

**移除原因**: 实体信息默认输出到 stderr，更符合日志和诊断信息的惯例。

#### Scenario: ~~使用 --print-entities 输出实体~~（已移除）
~~- **WHEN** 用户执行 `inu anonymize --file input.txt --print-entities`~~
~~- **THEN** 系统应该将实体信息输出到 stdout~~
