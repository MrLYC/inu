# Implementation Tasks

## 1. 依赖管理
- [x] 1.1 添加 CLI 框架依赖：`github.com/spf13/cobra`
- [x] 1.2 添加配置管理依赖：`github.com/spf13/viper`
- [x] 1.3 运行 `go mod tidy` 更新依赖

## 2. 创建 CLI 基础结构
- [x] 2.1 创建 `pkg/cli/` 目录
- [x] 2.2 实现 `pkg/cli/input.go`：输入处理函数
  - [x] `ReadInput(file, content string, stdin io.Reader) (string, error)` - 按优先级读取输入
- [x] 2.3 实现 `pkg/cli/output.go`：输出处理函数
  - [x] `WriteOutput(content string, print bool, outputFile string) error` - 处理输出
  - [x] `WriteEntities(entities []*Entity, print bool, outputFile string) error` - 输出实体信息
- [x] 2.4 实现 `pkg/cli/entities.go`：实体 YAML 序列化/反序列化
  - [x] 定义 YAML 结构体（包含 entities 数组）
  - [x] `SaveEntitiesToYAML(entities []*Entity, file string) error`
  - [x] `LoadEntitiesFromYAML(file string) ([]*Entity, error)` - 使用 viper 读取 YAML

## 3. 实现命令结构
- [x] 3.1 创建 `cmd/inu/commands/` 目录
- [x] 3.2 实现 `cmd/inu/commands/anonymize.go`：脱敏命令
  - [x] 定义命令和所有 flags（file, content, entity-types, print, print-entities, output, output-entities）
  - [x] 实现命令执行逻辑
  - [x] 添加参数验证（至少有一种输入）
  - [x] 集成 `pkg/anonymizer` 进行脱敏
  - [x] 处理输出（文本和实体）
- [x] 3.3 实现 `cmd/inu/commands/restore.go`：还原命令
  - [x] 定义命令和所有 flags（file, content, entities, print, output）
  - [x] 实现命令执行逻辑
  - [x] 添加参数验证（必须有 entities 参数）
  - [x] 集成 `pkg/anonymizer` 进行还原
  - [x] 处理输出

## 4. 重写主入口
- [x] 4.1 重写 `cmd/inu/main.go`
  - [x] 初始化 cobra root command
  - [x] 配置 app 元信息（Use, Short, Long, Version）
  - [x] 注册 `anonymize` 和 `restore` 子命令
  - [x] 添加全局 flags（如果需要）
  - [x] 实现版本信息注入（从编译时 ldflags）
- [x] 4.2 删除旧的 demo 代码

## 5. 错误处理和用户体验
- [x] 5.1 实现环境变量检查和友好错误提示
  - [x] 在 LLM 初始化失败时给出配置提示
- [x] 5.2 实现文件操作错误处理
  - [x] 文件不存在错误
  - [x] 权限错误
  - [x] YAML 解析错误
- [x] 5.3 添加进度提示（可选）
  - [x] 在 stderr 输出 "Processing..." 消息

## 6. 测试
- [x] 6.1 创建 `pkg/cli/input_test.go`：输入处理单元测试
- [x] 6.2 创建 `pkg/cli/output_test.go`：输出处理单元测试
- [x] 6.3 创建 `pkg/cli/entities_test.go`：实体序列化单元测试
- [x] 6.4 手动测试所有命令和参数组合
  - [x] `inu --help`
  - [x] `inu --version`
  - [x] `inu anonymize` 各种参数组合
  - [x] `inu restore` 各种参数组合

## 7. 文档更新
- [x] 7.1 更新 `README.md`
  - [x] 添加 CLI 使用示例
  - [x] 更新"快速开始"部分
  - [x] 添加命令参考文档
- [x] 7.2 更新 `openspec/project.md`
  - [x] 更新项目目的描述（从 demo 到可用工具）
  - [x] 添加新增依赖说明

## 8. 验证和构建
- [x] 8.1 运行所有测试：`make test`
- [x] 8.2 运行代码检查：`make lint`（如果已安装 golangci-lint）
- [x] 8.3 构建二进制文件：`make build`
- [x] 8.4 验证编译的二进制文件可以正常运行
- [x] 8.5 测试实际使用场景（端到端）
  - [x] 创建测试文件，执行脱敏，保存实体
  - [x] 使用保存的实体还原文本
  - [x] 验证还原结果与原文一致
  - 注：已创建 `test_e2e.sh` 脚本用于端到端测试，需要配置 OpenAI API 密钥后执行
