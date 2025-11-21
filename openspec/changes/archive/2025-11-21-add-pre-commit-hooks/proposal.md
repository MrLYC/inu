# 提案：添加 Pre-commit Hooks 配置

## Why

### 问题
当前项目使用 GitHub Actions 在 CI 阶段运行 golangci-lint 检查，但存在以下问题：

1. **反馈延迟**：开发者需要推送代码到 GitHub 后才能发现 lint 问题，导致：
   - 修复时需要额外的提交
   - CI/CD 管道浪费资源在失败的构建上
   - 开发迭代效率降低

2. **本地验证缺失**：没有统一的本地代码质量检查机制，导致：
   - 不同开发者可能使用不同的工具或配置
   - 容易遗漏代码格式化和 lint 检查
   - 代码质量不一致

3. **现有 lint 问题**：golangci-lint action 在 CI 中失败，需要：
   - 识别和修复现有的代码质量问题
   - 确保修复后的代码符合 lint 标准
   - 防止未来引入新的问题

### 影响范围
- **严重性**: 🟡 Medium - 影响开发效率和代码质量
- **受影响用户**: 所有项目贡献者
- **受影响流程**:
  - 开发工作流（本地开发）
  - CI/CD 流程（代码质量检查）

### 目标
1. **添加 pre-commit hooks 配置**：
   - 在本地提交前自动运行代码格式化和 lint 检查
   - 阻止不符合质量标准的代码提交
   - 提供即时反馈，提高开发效率

2. **修复现有 lint 问题**：
   - 识别所有 golangci-lint 报告的问题
   - 修复代码使其通过所有 lint 检查
   - 确保 CI 中的 lint job 能够通过

3. **改进开发者体验**：
   - 提供简单的设置文档
   - 自动化代码质量检查流程
   - 减少 CI 失败和修复周期

## What Changes

### 新增文件
- `.pre-commit-config.yaml`: Pre-commit hooks 配置文件
  - gofmt hook: 自动格式化 Go 代码
  - goimports hook: 自动整理导入语句
  - golangci-lint hook: 运行 lint 检查
  - go test hook: 可选的测试运行（可配置）

### 受影响的组件
- **开发工具链**：
  - 添加 pre-commit 框架
  - 配置 Go 相关 hooks
  - 更新 mise.toml（如果需要添加工具）

- **文档**：
  - README.md: 添加 pre-commit 安装和使用说明
  - 可能添加 CONTRIBUTING.md: 开发者贡献指南

- **代码库**：
  - 修复所有现有的 golangci-lint 问题
  - 确保代码通过所有检查

### 技术方案

#### Pre-commit 配置
使用 [pre-commit](https://pre-commit.com/) 框架，它支持：
- 多语言 hooks
- 自动安装依赖
- 本地和 CI 中一致的检查

**配置示例**：
```yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files

  - repo: https://github.com/golangci/golangci-lint
    rev: v1.56.0
    hooks:
      - id: golangci-lint
        args: ['--timeout=5m']

  - repo: local
    hooks:
      - id: gofmt
        name: gofmt
        entry: gofmt -w
        language: system
        files: \.go$

      - id: goimports
        name: goimports
        entry: goimports -w -local github.com/mrlyc/inu
        language: system
        files: \.go$
```

#### Lint 问题修复策略
1. 运行 `golangci-lint run --timeout=5m` 识别所有问题
2. 按优先级修复：
   - **高优先级**：类型错误、未使用的导入、明显的 bug
   - **中优先级**：代码格式、命名规范
   - **低优先级**：代码风格建议
3. 如果某些问题需要抑制，在 `.golangci.yml` 中添加明确的注释说明原因

#### 开发者工作流
1. **首次设置**：
   ```bash
   # 安装 pre-commit
   pip install pre-commit  # 或 brew install pre-commit

   # 安装 hooks
   pre-commit install
   ```

2. **日常使用**：
   ```bash
   git add .
   git commit -m "..."  # 自动运行 hooks
   ```

3. **手动运行**（可选）：
   ```bash
   pre-commit run --all-files  # 检查所有文件
   ```

4. **跳过 hooks**（紧急情况）：
   ```bash
   git commit --no-verify -m "..."
   ```

## Benefits

### 开发效率
- ✅ 本地即时反馈，无需等待 CI
- ✅ 减少因 lint 失败导致的额外提交
- ✅ 统一的代码质量标准

### CI/CD 优化
- ✅ 减少因代码质量问题导致的 CI 失败
- ✅ 节省 CI 资源和时间
- ✅ 更快的 PR 合并周期

### 代码质量
- ✅ 防止低质量代码进入代码库
- ✅ 一致的代码风格和格式
- ✅ 及早发现潜在问题

## Risks

### 技术风险
- **低风险**：Pre-commit 是成熟的工具，广泛使用
- **依赖管理**：需要开发者本地安装 pre-commit（可通过文档说明）
- **性能影响**：Hooks 可能增加提交时间（通常 <30 秒，可接受）

### 迁移风险
- **现有代码修复**：可能需要较大的修复工作量
- **向后兼容**：不影响现有的 CI/CD 流程
- **开发者适应**：需要文档和沟通，确保团队理解新流程

### 缓解措施
- 提供清晰的安装和使用文档
- 在 README 中突出说明 pre-commit 的价值
- 支持 `--no-verify` 选项用于紧急情况
- 分阶段修复现有 lint 问题（如果太多）

## Alternatives Considered

### 1. 仅依赖 CI 检查
**优点**：
- 无需本地配置
- 统一的检查环境

**缺点**：
- 反馈延迟
- CI 资源浪费
- 开发效率低

**结论**：不推荐，当前问题无法解决

### 2. 使用 Git hooks 脚本
**优点**：
- 无需额外依赖
- 轻量级

**缺点**：
- 需要手动管理脚本
- 跨平台兼容性差
- 难以版本控制和分发

**结论**：不推荐，pre-commit 提供更好的抽象

### 3. 使用 IDE 插件
**优点**：
- IDE 集成体验好
- 实时反馈

**缺点**：
- 不同开发者使用不同 IDE
- 无法强制执行
- 不一致的配置

**结论**：可作为补充，但不能替代 pre-commit

## Affected Specs
- `ci-cd`：添加 Pre-commit Hooks 相关需求

## References
- Pre-commit 官方文档: https://pre-commit.com/
- golangci-lint pre-commit: https://github.com/golangci/golangci-lint/tree/master/.github/hooks
- GitHub Actions golangci-lint: https://github.com/golangci/golangci-lint-action
