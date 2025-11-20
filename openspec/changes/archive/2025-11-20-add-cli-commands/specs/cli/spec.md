# CLI Specification

## ADDED Requirements

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
系统 SHALL 提供 `anonymize` 子命令来匿名化文本中的敏感信息。

#### Scenario: 从标准输入读取并打印到标准输出
- **WHEN** 用户执行 `echo "张三的电话是 13800138000" | inu anonymize --print`
- **THEN** 系统应该读取标准输入，匿名化文本，并将结果打印到标准输出

#### Scenario: 从文件读取内容
- **WHEN** 用户执行 `inu anonymize --file input.txt --print`
- **THEN** 系统应该读取 input.txt 文件的内容进行匿名化

#### Scenario: 从命令行参数读取内容
- **WHEN** 用户执行 `inu anonymize --content "张三的电话是 13800138000" --print`
- **THEN** 系统应该使用提供的内容字符串进行匿名化

#### Scenario: 输入优先级
- **WHEN** 用户同时指定 `--file`、`--content` 和标准输入
- **THEN** 系统应该按优先级使用：`--file` > `--content` > 标准输入

#### Scenario: 指定实体类型
- **WHEN** 用户执行 `inu anonymize --entity-types "个人信息,业务信息" --content "张三" --print`
- **THEN** 系统应该只识别和匿名化指定的实体类型

#### Scenario: 使用默认实体类型
- **WHEN** 用户执行 `inu anonymize` 而不指定 `--entity-types`
- **THEN** 系统应该使用默认实体类型列表：["个人信息", "业务信息", "资产信息", "账户信息", "位置数据", "文档名称", "组织机构", "岗位称谓"]

#### Scenario: 输出匿名化文本到文件
- **WHEN** 用户执行 `inu anonymize --file input.txt --output result.txt`
- **THEN** 系统应该将匿名化后的文本写入 result.txt 文件

#### Scenario: 同时输出到终端和文件
- **WHEN** 用户执行 `inu anonymize --content "text" --print --output result.txt`
- **THEN** 系统应该同时将匿名化文本打印到终端并写入文件

#### Scenario: 打印实体信息（简化格式）
- **WHEN** 用户执行 `inu anonymize --content "张三的电话是 13800138000" --print-entities`
- **THEN** 系统应该以简化格式打印实体信息到终端：
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
- **WHEN** 用户执行 `inu anonymize --print` 且无标准输入、无 `--file`、无 `--content`
- **THEN** 系统应该退出并显示错误：需要提供输入内容

#### Scenario: API 调用失败时报错
- **WHEN** 匿名化过程中 LLM API 调用失败（网络错误、认证失败等）
- **THEN** 系统应该退出并显示清晰的错误信息，包括失败原因

### Requirement: 还原命令
系统 SHALL 提供 `restore` 子命令来还原匿名化的文本。

#### Scenario: 从标准输入读取匿名文本并还原
- **WHEN** 用户执行 `echo "<个人信息[0].姓名.张三>" | inu restore --entities entities.yaml --print`
- **THEN** 系统应该读取标准输入和实体文件，还原文本并打印

#### Scenario: 从文件读取匿名文本
- **WHEN** 用户执行 `inu restore --file anonymized.txt --entities entities.yaml --print`
- **THEN** 系统应该读取文件内容进行还原

#### Scenario: 从命令行参数读取匿名文本
- **WHEN** 用户执行 `inu restore --content "<个人信息[0].姓名.张三>" --entities entities.yaml --print`
- **THEN** 系统应该还原提供的内容字符串

#### Scenario: 输入优先级
- **WHEN** 用户同时指定 `--file`、`--content` 和标准输入
- **THEN** 系统应该按优先级使用：`--file` > `--content` > 标准输入

#### Scenario: 输出还原文本到文件
- **WHEN** 用户执行 `inu restore --file anonymized.txt --entities entities.yaml --output restored.txt`
- **THEN** 系统应该将还原后的文本写入文件

#### Scenario: 同时输出到终端和文件
- **WHEN** 用户执行 `inu restore --content "text" --entities entities.yaml --print --output restored.txt`
- **THEN** 系统应该同时打印到终端并写入文件

#### Scenario: 缺少实体文件时报错
- **WHEN** 用户执行 `inu restore --content "text" --print` 但未指定 `--entities`
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
