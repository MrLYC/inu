# Project Context

## Purpose
Inu 是一个基于 AI 大模型的文本敏感信息匿名化工具。它能够识别文本中的个人信息、联系方式、组织名称等敏感实体，将其替换为占位符，并支持后续还原。主要应用场景包括：
- 数据脱敏处理
- 日志安全存储
- 测试数据生成
- 隐私保护

未来计划扩展为包含 CLI 和 Web 两种使用模式的完整工具。

## Tech Stack
- **语言**: Go 1.24.4
- **AI 框架**: CloudWeGo Eino（用于 LLM 交互）
- **LLM 提供商**: OpenAI API（支持自定义 base URL）
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
  - `pkg/anonymizer/`: 核心业务逻辑
  - `bin/`: 编译产物输出目录（不提交到版本控制）
  - `openspec/`: OpenSpec 规范和变更提案
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
支持识别和匿名化的实体类型包括：
- 人名
- 联系方式（电话、邮箱等）
- 职务
- 密码
- 组织名称
- 地址
- 文件名
- 账号
- 网址
- IP 地址

### 实体格式
匿名化后的占位符格式为：`<EntityType[ID].Category.Detail>`
- 示例：`<人名[0].姓名.张三>`

## Important Constraints
- 依赖外部 LLM API（OpenAI 或兼容服务）
- 需要配置环境变量：`OPENAI_API_KEY`, `OPENAI_MODEL_NAME`, `OPENAI_BASE_URL`
- 匿名化质量依赖于 LLM 的能力

## External Dependencies
- **CloudWeGo Eino**: AI 工具链框架
- **OpenAI API**: 大语言模型服务（可自定义 endpoint）
