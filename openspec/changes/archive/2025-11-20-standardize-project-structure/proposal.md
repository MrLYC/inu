# Standardize Project Structure

## Why
当前项目缺乏标准化的构建和发布流程，导致编译产物管理混乱、无法自动化发布。需要建立规范的项目结构以支持多平台编译、自动化测试和 GitHub Release 发布。

## What Changes
- 添加 `Makefile` 支持 `make build`、`make test`、`make clean`、`make lint` 等标准命令
- 配置编译输出到 `bin/` 目录，支持多平台交叉编译（linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64）
- 创建 GitHub Actions workflow 实现：
  - PR 时自动运行测试和 lint
  - 创建 release tag (vX.Y.Z) 时自动构建多平台二进制文件并发布到 GitHub Release
- 重组代码结构：将 `main.go` 移至 `cmd/inu/main.go`，核心逻辑提取到 `pkg/` 目录
- 完善 `.gitignore` 忽略编译产物和临时文件
- 添加项目基础文档（README.md, LICENSE）

## Impact
- 影响的 specs: 新增 `build-system` 和 `ci-cd` 两个 capability
- 影响的代码:
  - 项目根目录：新增 `Makefile`
  - `.github/workflows/`：新增 `ci.yml` 和 `release.yml`
  - 代码结构重组：`main.go` → `cmd/inu/main.go`，核心逻辑 → `pkg/anonymizer/`
  - `.gitignore`：添加编译产物忽略规则
  - 新增文档文件
