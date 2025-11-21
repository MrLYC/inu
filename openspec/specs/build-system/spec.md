# build-system Specification

## Purpose
TBD - created by archiving change standardize-project-structure. Update Purpose after archive.
## Requirements
### Requirement: Makefile Build Commands
系统 SHALL 提供标准的 Makefile 构建命令来管理项目的编译、测试和清理操作。

#### Scenario: 执行 make build 编译项目
- **WHEN** 开发者在项目根目录执行 `make build`
- **THEN** 系统应该编译当前平台的二进制文件到 `bin/inu` 或 `bin/inu.exe`（Windows）

#### Scenario: 执行 make build-all 交叉编译
- **WHEN** 开发者执行 `make build-all`
- **THEN** 系统应该为所有目标平台（linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64）编译二进制文件到 `bin/` 目录

#### Scenario: 执行 make test 运行测试
- **WHEN** 开发者执行 `make test`
- **THEN** 系统应该运行所有 Go 测试用例并输出结果

#### Scenario: 执行 make clean 清理产物
- **WHEN** 开发者执行 `make clean`
- **THEN** 系统应该删除 `bin/` 目录下的所有编译产物

#### Scenario: 执行 make lint 代码检查
- **WHEN** 开发者执行 `make lint`
- **THEN** 系统应该运行 golangci-lint 进行代码质量检查

### Requirement: 多平台编译支持
系统 SHALL 支持为多个目标平台交叉编译二进制文件。

#### Scenario: 指定目标平台编译
- **WHEN** 设置 GOOS 和 GOARCH 环境变量后执行编译
- **THEN** 系统应该生成对应平台的二进制文件，命名格式为 `inu-${GOOS}-${GOARCH}`（Windows 追加 `.exe`）

#### Scenario: 生成所有平台二进制文件
- **WHEN** 执行全平台编译命令
- **THEN** 系统应该在 `bin/` 目录生成以下文件：
  - `inu-linux-amd64`
  - `inu-linux-arm64`
  - `inu-darwin-amd64`
  - `inu-darwin-arm64`
  - `inu-windows-amd64.exe`

### Requirement: 版本信息注入
系统 SHALL 在编译时注入版本信息到二进制文件中。

#### Scenario: 编译时注入 Git 版本信息
- **WHEN** 执行 `make build` 或 `make build-all`
- **THEN** 编译的二进制文件应该包含版本号、Git commit hash 和构建时间信息

#### Scenario: 从 Git tag 获取版本号
- **WHEN** 当前 commit 有 Git tag（格式为 vX.Y.Z）
- **THEN** 版本号应该使用该 tag，否则使用 `dev`

### Requirement: 项目结构组织
系统 SHALL 采用标准的 Go 项目结构组织代码。

#### Scenario: 入口文件位置
- **WHEN** 查看项目结构
- **THEN** main 函数应该位于 `cmd/inu/main.go`

#### Scenario: 核心逻辑分离
- **WHEN** 查看项目结构
- **THEN** 核心业务逻辑应该位于 `pkg/anonymizer/` 目录，包括：
  - `anonymizer.go`：脱敏核心逻辑
  - `entity.go`：实体定义
  - `llm.go`：LLM 客户端封装

### Requirement: 静态资源编译时嵌入
系统 SHALL 在编译时将 Web UI 静态资源自动嵌入到二进制文件中，无需额外的构建步骤或工具。

#### Scenario: 编译自动包含静态资源
- **GIVEN** 存在 `pkg/web/static/` 目录包含 HTML/CSS/JS 文件
- **WHEN** 执行 `make build` 或 `go build ./cmd/inu`
- **THEN** 生成的二进制文件应该包含所有静态资源
- **AND** 二进制文件可以独立运行，无需外部静态文件目录

#### Scenario: 静态文件修改后重新编译
- **GIVEN** 修改了 `pkg/web/static/` 中的任何文件
- **WHEN** 重新执行 `make build`
- **THEN** 新生成的二进制文件应该包含更新后的静态资源
- **AND** 运行新二进制时应该加载新版本的静态文件

#### Scenario: 多平台编译包含静态资源
- **WHEN** 执行 `make build-all` 进行交叉编译
- **THEN** 所有平台的二进制文件应该都包含静态资源
- **AND** 每个平台的二进制文件都可以独立运行 Web UI

#### Scenario: 编译失败当静态文件缺失
- **GIVEN** `pkg/web/static/` 目录不存在或缺少关键文件
- **WHEN** 执行 `make build`
- **THEN** 编译应该失败并提示静态文件缺失错误
- **AND** 错误信息应该明确指出缺失的文件或目录

### Requirement: 二进制文件大小管理
系统 SHALL 合理管理嵌入静态资源后的二进制文件大小。

#### Scenario: 嵌入静态资源后的文件大小
- **GIVEN** 静态资源总大小约为 10-20KB
- **WHEN** 编译包含嵌入资源的二进制文件
- **THEN** 二进制文件大小增加应该接近静态资源大小（±10%）
- **AND** 文件大小增加对分发和部署无显著影响

#### Scenario: 构建输出显示文件大小
- **WHEN** 执行 `make build` 或 `make build-all`
- **THEN** 构建完成后应该显示二进制文件的大小信息
- **AND** 开发者可以验证嵌入资源的影响
