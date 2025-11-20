# Implementation Tasks

## 1. 项目结构重组
- [x] 1.1 创建 `cmd/inu/` 目录
- [x] 1.2 将 `main.go` 移动到 `cmd/inu/main.go`
- [x] 1.3 创建 `pkg/anonymizer/` 目录
- [x] 1.4 提取核心逻辑到 `pkg/anonymizer/anonymizer.go`
- [x] 1.5 提取实体定义到 `pkg/anonymizer/entity.go`
- [x] 1.6 提取 LLM 客户端到 `pkg/anonymizer/llm.go`
- [x] 1.7 更新 `cmd/inu/main.go` 使用新的包结构
- [x] 1.8 删除根目录旧的 `main.go`

## 2. 构建系统实现
- [x] 2.1 创建 `Makefile`，实现以下 targets：
  - [x] `build`: 编译当前平台二进制文件
  - [x] `build-all`: 交叉编译所有平台
  - [x] `test`: 运行测试
  - [x] `lint`: 代码检查
  - [x] `clean`: 清理编译产物
  - [x] `help`: 显示帮助信息
- [x] 2.2 实现版本信息注入（version, commit, build time）
- [x] 2.3 配置输出目录为 `bin/`
- [x] 2.4 配置目标平台列表（linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64）
- [x] 2.5 创建 `.golangci.yml` 配置文件

## 3. 更新 .gitignore
- [x] 3.1 添加 `bin/` 目录到忽略列表
- [x] 3.2 添加常见 Go 编译产物（`*.exe`, `*.test`, `*.out`）
- [x] 3.3 添加 IDE 相关文件（`.idea/`, `.vscode/`, `*.swp`）
- [x] 3.4 添加 macOS 系统文件（`.DS_Store`）

## 4. CI 工作流实现
- [x] 4.1 创建 `.github/workflows/ci.yml`
- [x] 4.2 配置触发条件（push to main, pull_request）
- [x] 4.3 配置 Go 版本（1.24.x）
- [x] 4.4 添加依赖缓存
- [x] 4.5 添加测试步骤（`go test ./...`）
- [x] 4.6 添加 lint 步骤（golangci-lint）
- [x] 4.7 添加编译验证步骤

## 5. Release 工作流实现
- [x] 5.1 创建 `.github/workflows/release.yml`
- [x] 5.2 配置触发条件（push tags `v*.*.*`）
- [x] 5.3 配置多平台编译矩阵
- [x] 5.4 实现产物打包（tar.gz 和 zip）
- [x] 5.5 生成 SHA256 校验和文件
- [x] 5.6 配置 GitHub Release 创建
- [x] 5.7 配置产物上传

## 6. 文档完善
- [x] 6.1 创建 `README.md`，包含：
  - [x] 项目简介
  - [x] 安装说明
  - [x] 使用示例
  - [x] 环境变量配置
  - [x] 构建说明
- [x] 6.2 创建或更新 `LICENSE` 文件（Apache 2.0）
- [x] 6.3 在 README 中添加 CI 状态徽章

## 7. 测试验证
- [x] 7.1 本地测试 `make build` - 编译成功
- [x] 7.2 本地测试 `make build-all` - 所有平台编译成功（5个平台）
- [x] 7.3 本地测试 `make test` - 通过（当前无测试文件）
- [x] 7.4 本地测试 `make lint` - 正确提示需要安装 golangci-lint
- [ ] 7.5 推送代码验证 CI workflow
- [ ] 7.6 创建测试 tag 验证 release workflow
- [x] 7.7 验证编译产物可正常运行 - bin/inu 运行正常
