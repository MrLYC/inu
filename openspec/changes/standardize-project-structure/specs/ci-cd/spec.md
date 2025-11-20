# CI/CD Specification

## ADDED Requirements

### Requirement: 持续集成工作流
系统 SHALL 提供 GitHub Actions 工作流以在 Pull Request 时自动运行测试和代码检查。

#### Scenario: PR 触发自动测试
- **WHEN** 创建或更新 Pull Request
- **THEN** GitHub Actions 应该自动运行以下检查：
  - Go 版本使用 1.24.x
  - 执行 `go test ./...`
  - 执行 `golangci-lint run`
  - 所有检查通过后显示绿色状态

#### Scenario: 推送到 main 分支触发检查
- **WHEN** 推送代码到 main 分支
- **THEN** GitHub Actions 应该自动运行测试和 lint 检查

#### Scenario: 多平台编译验证
- **WHEN** CI 运行时
- **THEN** 应该验证所有目标平台的编译是否成功（不上传产物）

### Requirement: 自动发布工作流
系统 SHALL 提供 GitHub Actions 工作流以在创建 release tag 时自动构建和发布。

#### Scenario: 创建 release tag 触发发布
- **WHEN** 推送格式为 `v*.*.*` 的 Git tag（如 v1.0.0）
- **THEN** GitHub Actions 应该：
  1. 为所有目标平台编译二进制文件
  2. 为每个平台创建压缩包（.tar.gz 或 .zip）
  3. 计算 SHA256 校验和
  4. 创建 GitHub Release 并上传所有产物

#### Scenario: Release 产物命名规范
- **WHEN** 创建 release 产物
- **THEN** 产物命名应该遵循以下格式：
  - Linux/macOS: `inu-${VERSION}-${OS}-${ARCH}.tar.gz`
  - Windows: `inu-${VERSION}-windows-${ARCH}.zip`
  - 校验和文件: `inu-${VERSION}-checksums.txt`

#### Scenario: Release 自动生成说明
- **WHEN** 创建 GitHub Release
- **THEN** Release 应该包含：
  - 版本号作为标题
  - 自动生成的变更日志（基于 commits）
  - 安装说明
  - 所有平台的下载链接和 SHA256 校验和

### Requirement: 工作流配置
系统 SHALL 使用合理的 GitHub Actions 配置以确保效率和安全。

#### Scenario: 使用最新稳定版本的 actions
- **WHEN** 配置 GitHub Actions workflow
- **THEN** 应该使用以下 actions 的最新稳定版本：
  - `actions/checkout@v4`
  - `actions/setup-go@v5`
  - `golangci/golangci-lint-action@v6`
  - `softprops/action-gh-release@v2`

#### Scenario: 最小化权限原则
- **WHEN** 配置 workflow 权限
- **THEN** 应该只授予必要的权限：
  - CI workflow: `contents: read`
  - Release workflow: `contents: write`

#### Scenario: 依赖缓存优化
- **WHEN** 运行 CI/CD workflow
- **THEN** 应该缓存 Go modules 以加速构建
