# Project Context

## Purpose
Inu 是一个基于 AI 大模型的文本敏感信息脱敏工具。它能够识别文本中的个人信息、联系方式、组织名称等敏感实体，将其替换为占位符，并支持后续还原。主要应用场景包括：
- 数据脱敏处理
- 日志安全存储
- 测试数据生成
- 隐私保护

未来计划扩展为包含 CLI 和 Web 两种使用模式的完整工具。

## Tech Stack
- **语言**: Go 1.24.4
- **AI 框架**: CloudWeGo Eino（用于 LLM 交互）
- **LLM 提供商**: OpenAI API（支持自定义 base URL）
- **CLI 框架**: Cobra（命令行解析）+ Viper（配置管理）
- **Web 框架**: Gin（HTTP API 服务器）
- **构建工具**: Make, Go toolchain
- **CI/CD**: GitHub Actions
- **版本管理**: Git tags (vX.Y.Z)

## Project Conventions

### Code Style
- 使用 `gofmt` 标准格式化
- 使用 `golangci-lint` 进行代码质量检查
- 包命名使用小写单词，避免下划线
- 导出的标识符使用驼峰命名法
- 错误处理使用 `github.com/rotisserie/eris` 进行错误包装

### Architecture Patterns
- **项目结构**:
  - `cmd/inu/`: CLI 入口点
  - `cmd/inu/commands/`: CLI 子命令实现（anonymize, restore, web）
  - `pkg/anonymizer/`: 核心业务逻辑
  - `pkg/cli/`: CLI 工具函数（输入输出、实体管理）
  - `pkg/web/`: Web API 服务器（handlers, middleware, server）
  - `bin/`: 编译产物输出目录（不提交到版本控制）
  - `openspec/`: OpenSpec 规范和变更提案
- **CLI 架构**: 使用 Cobra 构建子命令结构，Viper 处理配置文件
- **Web 架构**: 使用 Gin 框架提供 RESTful API，支持 HTTP Basic Auth
- **依赖注入**: 构造函数模式（如 `NewHas`）
- **错误处理**: 统一使用 eris 进行错误包装和追踪

### Testing Strategy
- 单元测试覆盖核心逻辑
- 使用 `go test ./...` 运行所有测试
- CI 中自动运行测试确保质量

### Git Workflow
- **主分支**: `main`
- **标签格式**: `vX.Y.Z`（语义化版本）
- **提交约定**: 清晰描述变更内容
- **Release 流程**: 
  1. 推送 `v*.*.*` 格式的 Git tag
  2. GitHub Actions 自动构建多平台二进制文件
  3. 自动创建 GitHub Release 并上传产物

## Domain Context
### 敏感信息类型
默认支持识别和脱敏的实体类型包括：
- 个人信息：姓名、身份证号、电话号码等
- 业务信息：业务数据、客户信息等
- 资产信息：财产、资源信息等
- 账户信息：银行账号、信用卡号等
- 位置数据：地址、地理位置等
- 文档名称：文件名、文档标题等
- 组织机构：公司名称、机构名称等
- 岗位称谓：职位、头衔等

用户可以通过 CLI 参数自定义实体类型。

### 实体格式
脱敏后的占位符格式为：`<EntityType[ID].Category.Detail>`
- 示例：`<个人信息[0].姓名.全名>`

### CLI 命令
- `inu anonymize`: 脱敏文本
  - 输入：`--file` / `--content` / stdin（优先级递减）
  - 输出：`--print` 和/或 `--output`
  - 实体：`--output-entities` 保存到 YAML 文件
- `inu restore`: 还原文本
  - 输入：`--file` / `--content` / stdin
  - 实体：`--entities` (必需)
  - 输出：`--print` 和/或 `--output`
- `inu web`: 启动 Web API 服务器
  - 配置：`--addr` (监听地址), `--admin-user`, `--admin-token`
  - API 端点：
    - `GET /health` - 健康检查（无需认证）
    - `POST /api/v1/anonymize` - 脱敏文本（需要认证）
    - `POST /api/v1/restore` - 还原文本（需要认证）

## Important Constraints
- 依赖外部 LLM API（OpenAI 或兼容服务）
- 需要配置环境变量：`OPENAI_API_KEY`, `OPENAI_MODEL_NAME`, `OPENAI_BASE_URL`
- 脱敏质量依赖于 LLM 的能力

## External Dependencies
- **CloudWeGo Eino**: AI 工具链框架
- **OpenAI API**: 大语言模型服务（可自定义 endpoint）
- **Cobra**: CLI 命令行框架
- **Viper**: 配置和 YAML 文件管理
- **Eris**: Go 错误处理增强
- **Gin**: Web 框架（HTTP API 服务器）
