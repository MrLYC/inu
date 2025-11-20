# cli Specification

## Purpose
定义 Inu CLI 命令行接口的输出行为规范和命令接口。
## Requirements

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

### Requirement: 命令行接口
系统 SHALL 提供命令行接口来执行文本匿名化和还原操作。

#### Scenario: 显示帮助信息
- **WHEN** 用户执行 `inu --help` 或 `inu -h`
- **THEN** 系统应该显示可用命令和全局选项的帮助信息

#### Scenario: 显示版本信息
- **WHEN** 用户执行 `inu --version` 或 `inu -v`
- **THEN** 系统应该显示版本号、Git commit hash 和构建时间

#### Scenario: 显示子命令帮助
- **WHEN** 用户执行 `inu anonymize --help` 或 `inu restore --help`
- **THEN** 系统应该显示该子命令的详细使用说明和参数列表

### Requirement: 匿名化命令
系统 SHALL 提供 `anonymize` 子命令来匿名化文本中的敏感信息，使用流式输出改善用户体验。

#### Scenario: 从标准输入读取并流式输出到标准输出
- **WHEN** 用户执行 `echo "张三的电话是 13800138000" | inu anonymize`
- **THEN** 系统应该读取标准输入，流式生成匿名化文本
- **AND** 实时输出到标准输出（逐 token）
- **AND** 在流式输出完成后，实体信息输出到 stderr

#### Scenario: 从文件读取并流式输出
- **WHEN** 用户执行 `inu anonymize --file input.txt`
- **THEN** 系统应该读取 input.txt 文件的内容
- **AND** 流式生成并输出匿名化文本到标准输出
- **AND** 用户可以实时看到输出进度

#### Scenario: 从命令行参数读取内容
- **WHEN** 用户执行 `inu anonymize --content "张三的电话是 13800138000"`
- **THEN** 系统应该使用提供的内容字符串进行匿名化并输出到标准输出

#### Scenario: 输入优先级
- **WHEN** 用户同时指定 `--file`、`--content` 和标准输入
- **THEN** 系统应该按优先级使用：`--file` > `--content` > 标准输入

#### Scenario: 指定实体类型
- **WHEN** 用户执行 `inu anonymize --entity-types "个人信息,业务信息" --content "张三"`
- **THEN** 系统应该只识别和匿名化指定的实体类型并输出到标准输出

#### Scenario: 使用默认实体类型
- **WHEN** 用户执行 `inu anonymize` 而不指定 `--entity-types`
- **THEN** 系统应该使用默认实体类型列表：["个人信息", "业务信息", "资产信息", "账户信息", "位置数据", "文档名称", "组织机构", "岗位称谓"]

#### Scenario: 输出匿名化文本到文件
- **WHEN** 用户执行 `inu anonymize --file input.txt --output result.txt`
- **THEN** 系统应该流式生成匿名化文本
- **AND** 同时写入 result.txt 和 stdout
- **AND** 两个输出目标都是流式写入

#### Scenario: 流式输出到管道
- **WHEN** 用户执行 `inu anonymize --file input.txt | grep "个人信息"`
- **THEN** 系统应该将 token 实时传递给管道下游命令
- **AND** 下游命令可以立即开始处理
- **AND** 不会因为上游缓冲导致延迟

#### Scenario: 只输出到文件不显示
- **WHEN** 用户执行 `inu anonymize --content "text" --no-print --output result.txt`
- **THEN** 系统应该只将匿名化文本写入文件，不输出到标准输出

#### Scenario: 实体信息输出到 stderr（默认）
- **WHEN** 用户执行 `inu anonymize --content "张三的电话是 13800138000"`
- **THEN** 系统应该将实体信息输出到 stderr：
  ```
  <个人信息[0].姓名.张三>: 张三
  <个人信息[1].电话.13800138000>: 13800138000
  ```

#### Scenario: 输出实体信息到 YAML 文件
- **WHEN** 用户执行 `inu anonymize --file input.txt --output-entities entities.yaml`
- **THEN** 系统应该将实体信息以 YAML 格式写入文件：
  ```yaml
  entities:
    - key: "<个人信息[0].姓名.张三>"
      type: "个人信息"
      id: "0"
      category: "姓名"
      detail: "张三"
      values:
        - "张三"
  ```

#### Scenario: 无输入内容时报错
- **WHEN** 用户执行 `inu anonymize` 且无标准输入、无 `--file`、无 `--content`
- **THEN** 系统应该退出并显示错误：需要提供输入内容

#### Scenario: API 调用失败时报错
- **WHEN** 匿名化过程中 LLM API 调用失败（网络错误、认证失败等）
- **THEN** 系统应该退出并显示清晰的错误信息，包括失败原因

#### Scenario: 流式输出被中断
- **WHEN** 用户在流式输出过程中按 Ctrl+C
- **THEN** 系统应该立即停止 LLM 请求
- **AND** 已输出的部分保留在输出中
- **AND** 实体信息可能不完整（未完全解析）

#### Scenario: 流式输出遇到错误
- **WHEN** 流式输出过程中 LLM 返回错误（如网络中断）
- **THEN** 系统应该立即停止输出
- **AND** 显示错误信息到 stderr
- **AND** 已输出的部分文本保留
- **AND** 返回非零退出码

### Requirement: 还原命令
系统 SHALL 提供 `restore` 子命令来还原匿名化的文本。

#### Scenario: 从标准输入读取匿名文本并还原
- **WHEN** 用户执行 `echo "<个人信息[0].姓名.张三>" | inu restore --entities entities.yaml`
- **THEN** 系统应该读取标准输入和实体文件，还原文本并输出到标准输出（默认行为）

#### Scenario: 从文件读取匿名文本
- **WHEN** 用户执行 `inu restore --file anonymized.txt --entities entities.yaml`
- **THEN** 系统应该读取文件内容进行还原并输出到标准输出

#### Scenario: 从命令行参数读取匿名文本
- **WHEN** 用户执行 `inu restore --content "<个人信息[0].姓名.张三>" --entities entities.yaml`
- **THEN** 系统应该还原提供的内容字符串并输出到标准输出

#### Scenario: 输入优先级
- **WHEN** 用户同时指定 `--file`、`--content` 和标准输入
- **THEN** 系统应该按优先级使用：`--file` > `--content` > 标准输入

#### Scenario: 输出还原文本到文件
- **WHEN** 用户执行 `inu restore --file anonymized.txt --entities entities.yaml --output restored.txt`
- **THEN** 系统应该将还原后的文本写入文件
- **AND** 同时输出到标准输出（默认行为）

#### Scenario: 只输出到文件不显示
- **WHEN** 用户执行 `inu restore --content "text" --entities entities.yaml --no-print --output restored.txt`
- **THEN** 系统应该只将还原后的文本写入文件，不输出到标准输出

#### Scenario: 缺少实体文件时报错
- **WHEN** 用户执行 `inu restore --content "text"` 但未指定 `--entities`
- **THEN** 系统应该退出并显示错误：需要提供实体文件

#### Scenario: 实体文件格式错误时报错
- **WHEN** 用户提供的 `--entities` 文件不是有效的 YAML 格式
- **THEN** 系统应该退出并显示清晰的解析错误信息

### Requirement: 错误处理和用户体验
系统 SHALL 提供清晰的错误信息和良好的用户体验。

#### Scenario: 环境变量未配置时提示
- **WHEN** 用户执行命令但未设置 `OPENAI_API_KEY` 等必需环境变量
- **THEN** 系统应该显示友好的错误信息，说明需要配置的环境变量和示例

#### Scenario: 文件不存在时报错
- **WHEN** 用户指定的 `--file` 或 `--entities` 文件不存在
- **THEN** 系统应该退出并显示文件不存在的错误

#### Scenario: 写入文件权限不足时报错
- **WHEN** 用户指定的 `--output` 或 `--output-entities` 路径无写入权限
- **THEN** 系统应该退出并显示权限错误

#### Scenario: 显示进度指示
- **WHEN** 处理大文件或 API 调用耗时较长
- **THEN** 系统应该在 stderr 显示进度提示（如 "Processing..." 或进度条）

## Implementation Notes

### Streaming Anonymization
- 流式输出基于 CloudWeGo Eino 的 `StreamReader` 接口
- 内部使用 `llm.Stream()` 替代 `llm.Generate()`
- 实体映射在完整响应后解析（`<<<PAIR>>>` 分隔符）
- 使用 `io.MultiWriter` 支持同时输出到多个目标
- 在 `<<<PAIR>>>` 标记前：所有 token 实时写入 writer
- 在 `<<<PAIR>>>` 标记后：token 用于实体 JSON 解析，不写入 writer
- Fallback 机制：Stream() 失败时自动使用 Generate() 并解析完整响应