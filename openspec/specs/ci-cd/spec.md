# ci-cd Specification

## Purpose
TBD - created by archiving change standardize-project-structure. Update Purpose after archive.
## Requirements
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

#### Scenario: Lint 检查与本地一致
- **GIVEN** 开发者在本地运行了 pre-commit hooks
- **AND** 所有 hooks 都通过
- **WHEN** 推送代码到 GitHub 触发 CI
- **THEN** CI 中的 golangci-lint 检查应该也能通过
- **AND** 使用与本地相同的 `.golangci.yml` 配置
- **AND** 产生一致的检查结果

---

**实现注意事项**：

1. **Pre-commit 配置文件格式**：
   - 使用 YAML 格式
   - repos 列表包含所有 hook 源
   - 每个 hook 指定 id, name, entry, args 等属性

2. **Hooks 执行顺序**：
   - 先运行快速的格式化 hooks（gofmt, goimports）
   - 再运行较慢的检查 hooks（golangci-lint）
   - 失败时立即停止，不继续执行后续 hooks

3. **自动修复行为**：
   - gofmt 和 goimports 应该自动修复格式问题
   - golangci-lint 使用 `--fix` 参数自动修复可修复的问题
   - 修复后文件需要重新 add 并提交

4. **性能优化**：
   - golangci-lint 只检查变更的文件（默认行为）
   - 使用 `pass_filenames: false` 可以让 lint 检查整个项目
   - 设置合理的超时时间避免长时间等待

5. **开发者体验**：
   - 提供清晰的错误信息
   - 支持 `--no-verify` 跳过检查（紧急情况）
   - 文档中说明如何调试 hook 失败

6. **CI 集成**：
   - CI 保持作为最终的质量守门人
   - 即使开发者跳过本地 hooks，CI 仍然会检查
   - 不在 CI 中重复运行 pre-commit（避免重复）

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

### Requirement: Pre-commit Hooks 配置
系统 SHALL 提供 pre-commit hooks 配置，在本地提交前自动运行代码质量检查。

#### Scenario: 配置文件存在
- **GIVEN** 开发者克隆了项目仓库
- **WHEN** 查看项目根目录
- **THEN** 应该存在 `.pre-commit-config.yaml` 文件
- **AND** 配置文件应该包含以下 hooks：
  - 基础文件检查（trailing-whitespace, end-of-file-fixer, check-yaml 等）
  - Go 代码格式化（gofmt, goimports）
  - Go 代码检查（golangci-lint）

#### Scenario: 安装 pre-commit hooks
- **GIVEN** 开发者已安装 pre-commit 工具
- **WHEN** 在项目根目录运行 `pre-commit install`
- **THEN** 系统应该成功安装 git hooks
- **AND** 输出消息确认安装成功
- **AND** `.git/hooks/pre-commit` 文件应该被创建

#### Scenario: 提交时自动运行 hooks
- **GIVEN** pre-commit hooks 已安装
- **AND** 开发者修改了 Go 代码文件
- **WHEN** 执行 `git commit`
- **THEN** 系统应该自动运行以下检查：
  - 移除行尾空格
  - 确保文件以换行结束
  - 运行 gofmt 格式化代码
  - 运行 goimports 整理导入
  - 运行 golangci-lint 检查代码质量
- **AND** 如果所有检查通过，允许提交
- **AND** 如果任何检查失败，阻止提交并显示错误信息

#### Scenario: Hooks 自动修复问题
- **GIVEN** pre-commit hooks 已安装
- **AND** Go 代码有格式问题（如缺少换行符、行尾空格）
- **WHEN** 执行 `git commit`
- **THEN** hooks 应该自动修复这些问题
- **AND** 提示开发者文件已被修改
- **AND** 开发者需要重新 `git add` 并再次提交

#### Scenario: 跳过 hooks 检查
- **GIVEN** pre-commit hooks 已安装
- **WHEN** 开发者执行 `git commit --no-verify -m "message"`
- **THEN** 系统应该跳过所有 pre-commit hooks
- **AND** 直接创建提交
- **AND** 不运行任何代码质量检查

#### Scenario: 手动运行所有 hooks
- **GIVEN** pre-commit hooks 已安装
- **WHEN** 开发者执行 `pre-commit run --all-files`
- **THEN** 系统应该检查所有项目文件（不仅是 staged 文件）
- **AND** 报告所有发现的问题
- **AND** 如果启用了自动修复，自动修复可修复的问题

#### Scenario: 手动运行特定 hook
- **GIVEN** pre-commit hooks 已安装
- **WHEN** 开发者执行 `pre-commit run golangci-lint`
- **THEN** 系统应该只运行 golangci-lint hook
- **AND** 跳过其他 hooks
- **AND** 报告 golangci-lint 的检查结果

### Requirement: 本地代码质量工具
系统 SHALL 确保开发者能够在本地运行与 CI 相同的代码质量检查。

#### Scenario: golangci-lint 配置一致性
- **GIVEN** 项目中存在 `.golangci.yml` 配置文件
- **WHEN** 开发者在本地运行 `golangci-lint run`
- **THEN** 应该使用与 CI 相同的配置
- **AND** 检查相同的 linters
- **AND** 使用相同的超时设置（5分钟）
- **AND** 产生与 CI 一致的结果

#### Scenario: goimports 配置本地前缀
- **GIVEN** 项目使用 `github.com/mrlyc/inu` 作为模块路径
- **WHEN** 运行 goimports
- **THEN** 应该使用 `-local github.com/mrlyc/inu` 参数
- **AND** 本地包导入应该分组在第三方包之后
- **AND** 导入顺序应该一致

#### Scenario: 文档说明工具安装
- **GIVEN** 开发者阅读项目 README
- **WHEN** 查看开发环境设置章节
- **THEN** 应该包含以下工具的安装说明：
  - pre-commit 框架
  - goimports
  - golangci-lint
- **AND** 提供多种安装方式（pip, brew, mise）
- **AND** 包含验证安装的命令

### Requirement: 代码质量标准
系统 SHALL 确保所有提交的代码符合项目定义的质量标准。

#### Scenario: 所有代码通过 golangci-lint
- **GIVEN** 项目配置了 golangci-lint
- **WHEN** 运行 `golangci-lint run --timeout=5m`
- **THEN** 应该没有任何错误或警告
- **AND** 所有 Go 文件应该通过以下检查：
  - gofmt（格式化）
  - goimports（导入整理）
  - govet（静态分析）
  - errcheck（错误检查）
  - staticcheck（静态检查）
  - unused（未使用代码）
  - gosimple（简化建议）
  - ineffassign（无效赋值）
  - typecheck（类型检查）

#### Scenario: 代码格式一致性
- **GIVEN** 项目中的所有 Go 文件
- **WHEN** 运行 `gofmt -l .`
- **THEN** 不应该输出任何文件名
- **AND** 所有文件应该已经格式化
- **WHEN** 运行 `goimports -l -local github.com/mrlyc/inu .`
- **THEN** 不应该输出任何文件名
- **AND** 所有导入语句应该已经整理

#### Scenario: 没有未使用的导入
- **GIVEN** 项目中的所有 Go 文件
- **WHEN** 运行代码检查
- **THEN** 不应该存在未使用的导入语句
- **AND** 不应该存在未使用的变量（除非有合理的 `_ =` 赋值）

#### Scenario: YAML 文件语法正确
- **GIVEN** 项目中的配置文件（如 .golangci.yml, .github/workflows/*.yml）
- **WHEN** pre-commit hooks 检查 YAML 语法
- **THEN** 所有 YAML 文件应该语法正确
- **AND** 能够被正确解析
