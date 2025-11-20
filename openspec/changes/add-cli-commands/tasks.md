# Implementation Tasks

## 1. 依赖管理
- [ ] 1.1 添加 CLI 框架依赖：`github.com/spf13/cobra`
- [ ] 1.2 添加配置管理依赖：`github.com/spf13/viper`
- [ ] 1.3 运行 `go mod tidy` 更新依赖

## 2. 创建 CLI 基础结构
- [ ] 2.1 创建 `pkg/cli/` 目录
- [ ] 2.2 实现 `pkg/cli/input.go`：输入处理函数
  - [ ] `ReadInput(file, content string, stdin io.Reader) (string, error)` - 按优先级读取输入
- [ ] 2.3 实现 `pkg/cli/output.go`：输出处理函数
  - [ ] `WriteOutput(content string, print bool, outputFile string) error` - 处理输出
  - [ ] `WriteEntities(entities []*Entity, print bool, outputFile string) error` - 输出实体信息
- [ ] 2.4 实现 `pkg/cli/entities.go`：实体 YAML 序列化/反序列化
  - [ ] 定义 YAML 结构体（包含 entities 数组）
  - [ ] `SaveEntitiesToYAML(entities []*Entity, file string) error`
  - [ ] `LoadEntitiesFromYAML(file string) ([]*Entity, error)` - 使用 viper 读取 YAML

## 3. 实现命令结构
- [ ] 3.1 创建 `cmd/inu/commands/` 目录
- [ ] 3.2 实现 `cmd/inu/commands/anonymize.go`：匿名化命令
  - [ ] 定义命令和所有 flags（file, content, entity-types, print, print-entities, output, output-entities）
  - [ ] 实现命令执行逻辑
  - [ ] 添加参数验证（至少有一种输入）
  - [ ] 集成 `pkg/anonymizer` 进行匿名化
  - [ ] 处理输出（文本和实体）
- [ ] 3.3 实现 `cmd/inu/commands/restore.go`：还原命令
  - [ ] 定义命令和所有 flags（file, content, entities, print, output）
  - [ ] 实现命令执行逻辑
  - [ ] 添加参数验证（必须有 entities 参数）
  - [ ] 集成 `pkg/anonymizer` 进行还原
  - [ ] 处理输出

## 4. 重写主入口
- [ ] 4.1 重写 `cmd/inu/main.go`
  - [ ] 初始化 cobra root command
  - [ ] 配置 app 元信息（Use, Short, Long, Version）
  - [ ] 注册 `anonymize` 和 `restore` 子命令
  - [ ] 添加全局 flags（如果需要）
  - [ ] 实现版本信息注入（从编译时 ldflags）
- [ ] 4.2 删除旧的 demo 代码

## 5. 错误处理和用户体验
- [ ] 5.1 实现环境变量检查和友好错误提示
  - [ ] 在 LLM 初始化失败时给出配置提示
- [ ] 5.2 实现文件操作错误处理
  - [ ] 文件不存在错误
  - [ ] 权限错误
  - [ ] YAML 解析错误
- [ ] 5.3 添加进度提示（可选）
  - [ ] 在 stderr 输出 "Processing..." 消息

## 6. 测试
- [ ] 6.1 创建 `pkg/cli/input_test.go`：输入处理单元测试
- [ ] 6.2 创建 `pkg/cli/output_test.go`：输出处理单元测试
- [ ] 6.3 创建 `pkg/cli/entities_test.go`：实体序列化单元测试
- [ ] 6.4 手动测试所有命令和参数组合
  - [ ] `inu --help`
  - [ ] `inu --version`
  - [ ] `inu anonymize` 各种参数组合
  - [ ] `inu restore` 各种参数组合

## 7. 文档更新
- [ ] 7.1 更新 `README.md`
  - [ ] 添加 CLI 使用示例
  - [ ] 更新"快速开始"部分
  - [ ] 添加命令参考文档
- [ ] 7.2 更新 `openspec/project.md`
  - [ ] 更新项目目的描述（从 demo 到可用工具）
  - [ ] 添加新增依赖说明

## 8. 验证和构建
- [ ] 8.1 运行所有测试：`make test`
- [ ] 8.2 运行代码检查：`make lint`（如果已安装 golangci-lint）
- [ ] 8.3 构建二进制文件：`make build`
- [ ] 8.4 验证编译的二进制文件可以正常运行
- [ ] 8.5 测试实际使用场景（端到端）
  - [ ] 创建测试文件，执行匿名化，保存实体
  - [ ] 使用保存的实体还原文本
  - [ ] 验证还原结果与原文一致
