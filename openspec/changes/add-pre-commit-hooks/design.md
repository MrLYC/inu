# 设计文档：添加 Pre-commit Hooks 配置

## Context

当前项目在 CI 阶段使用 golangci-lint 进行代码质量检查，但开发者在本地提交代码时没有自动化的质量检查机制。这导致：
1. CI 经常因 lint 问题失败
2. 开发者需要多次提交来修复简单的格式问题
3. 代码审查时花费时间在格式和基础质量问题上

通过引入 pre-commit hooks，可以在本地提交前自动检查代码质量，提供即时反馈。

## Goals / Non-Goals

### Goals
- 在本地 git commit 时自动运行代码质量检查
- 修复当前所有 golangci-lint 检查失败的问题
- 提供简单的设置和使用流程
- 保持与现有 CI 检查的一致性

### Non-Goals
- 不强制开发者必须使用 pre-commit（可以 --no-verify 跳过）
- 不替代 CI 检查（CI 仍然作为最后的守门人）
- 不引入新的 lint 规则（使用现有 .golangci.yml）
- 不要求复杂的本地环境配置

## Decisions

### Decision 1: 使用 Pre-commit 框架

**选择**: 使用 [pre-commit](https://pre-commit.com/) 框架

**理由**:
- ✅ 业界标准工具，成熟稳定
- ✅ 支持多语言和多种 hooks
- ✅ 自动管理 hook 依赖和版本
- ✅ 配置简单，易于维护
- ✅ 大量现成的 hooks 可用

**替代方案考虑**:
1. **自定义 Git hooks 脚本**
   - ❌ 需要手动管理跨平台兼容性
   - ❌ 难以版本控制和分发
   - ❌ 每个开发者需要手动复制脚本

2. **Husky (Node.js 工具)**
   - ❌ 需要 Node.js 依赖（项目是 Go）
   - ❌ 不是 Go 生态的标准工具
   - ❌ 额外的技术栈

3. **Lefthook**
   - ✅ Go 编写，性能好
   - ❌ 社区生态较小
   - ❌ 现成 hooks 较少
   - ❌ 配置相对复杂

**结论**: Pre-commit 提供最好的易用性和生态支持

### Decision 2: Pre-commit Hooks 配置

**选择的 Hooks**:

1. **基础文件检查** (pre-commit/pre-commit-hooks):
   ```yaml
   - trailing-whitespace  # 移除行尾空格
   - end-of-file-fixer    # 确保文件以换行结束
   - check-yaml          # 验证 YAML 语法
   - check-added-large-files  # 防止提交大文件
   - check-merge-conflict     # 检查未解决的合并冲突
   ```

2. **Go 代码格式化**:
   ```yaml
   - gofmt    # 标准 Go 格式化
   - goimports  # 整理导入语句
   ```

3. **Go Lint 检查**:
   ```yaml
   - golangci-lint  # 使用项目的 .golangci.yml 配置
   ```

**配置策略**:
- 使用项目现有的 `.golangci.yml` 配置（保持一致性）
- goimports 使用 `-local github.com/mrlyc/inu` 参数
- 设置合理的超时时间（5分钟）

**不包含的检查**:
- ❌ 单元测试：测试可能耗时较长，放在 CI 中运行
- ❌ 构建检查：编译在 CI 中验证更合适
- ❌ 安全扫描：可选，未来可添加

### Decision 3: 修复现有 Lint 问题的策略

**方法**:

1. **识别问题**:
   ```bash
   golangci-lint run --timeout=5m > lint-issues.txt
   ```

2. **分类和优先级**:
   - **P0 (必须修复)**:
     - 编译错误
     - 未使用的导入
     - 明显的逻辑错误
     - 类型错误

   - **P1 (应该修复)**:
     - 代码格式问题
     - 命名规范问题
     - 未使用的变量

   - **P2 (可以修复)**:
     - 代码风格建议
     - 可选的优化建议

3. **修复方式**:
   - 自动修复：
     ```bash
     gofmt -w .
     goimports -w -local github.com/mrlyc/inu .
     golangci-lint run --fix
     ```

   - 手动修复：对于自动修复无法处理的问题

   - 抑制（最后手段）：
     - 在 `.golangci.yml` 中添加 `issues.exclude-rules`
     - 必须包含注释说明原因

4. **验证**:
   ```bash
   golangci-lint run --timeout=5m
   go test ./...
   ```

### Decision 4: 开发者工作流

**安装流程**:

1. **文档位置**: README.md 中添加 "Development Setup" 章节

2. **安装命令**:
   ```bash
   # 方式 1: pip
   pip install pre-commit

   # 方式 2: brew (macOS)
   brew install pre-commit

   # 方式 3: mise (如果项目使用)
   mise use -g pre-commit@latest

   # 安装 hooks
   pre-commit install
   ```

3. **首次运行**:
   ```bash
   # 可选：检查所有文件
   pre-commit run --all-files
   ```

**日常使用**:
- 正常提交：`git commit -m "..."`（自动运行 hooks）
- 跳过检查：`git commit --no-verify -m "..."`（紧急情况）
- 手动运行：`pre-commit run --all-files`

**失败处理**:
- Hook 失败时会显示错误信息
- 修复问题后重新 `git add` 和 `git commit`
- 如果是格式问题，hook 可能已自动修复，只需重新 add

### Decision 5: CI 集成

**策略**: 保持现有 CI 配置不变

**理由**:
- Pre-commit hooks 是本地开发的辅助工具
- CI 仍然作为代码质量的最终守门人
- 允许开发者在特殊情况下跳过 hooks（--no-verify）
- CI 确保所有代码（包括跳过 hooks 的）都被检查

**可选增强**（未来考虑）:
```yaml
# 在 CI 中也运行 pre-commit（可选）
- name: Run pre-commit
  uses: pre-commit/action@v3.0.0
```

## Technical Details

### .pre-commit-config.yaml 完整配置

```yaml
# Pre-commit hooks configuration
# See https://pre-commit.com for more information

repos:
  # General file checks
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
        name: Remove trailing whitespace
      - id: end-of-file-fixer
        name: Fix end of files
      - id: check-yaml
        name: Check YAML syntax
        args: ['--unsafe']  # Allow custom YAML tags
      - id: check-added-large-files
        name: Check for large files
        args: ['--maxkb=1000']
      - id: check-merge-conflict
        name: Check for merge conflicts

  # Go-specific hooks
  - repo: local
    hooks:
      - id: gofmt
        name: Run gofmt
        entry: gofmt -w
        language: system
        files: \.go$
        pass_filenames: true

      - id: goimports
        name: Run goimports
        entry: goimports
        args: [-w, -local, github.com/mrlyc/inu]
        language: system
        files: \.go$
        pass_filenames: true

      - id: golangci-lint
        name: Run golangci-lint
        entry: golangci-lint run
        args: ['--timeout=5m', '--fix']
        language: system
        files: \.go$
        pass_filenames: false
```

### 工具依赖

**必需工具**:
- `pre-commit`: Hook 管理框架
- `go`: Go 工具链（已有）
- `gofmt`: 随 Go 一起安装
- `goimports`: 需要安装
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  ```
- `golangci-lint`: 需要安装
  ```bash
  # 方式 1: brew
  brew install golangci-lint

  # 方式 2: go install
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

  # 方式 3: mise
  mise use -g golangci-lint@latest
  ```

**版本要求**:
- pre-commit: >= 3.0.0
- goimports: latest
- golangci-lint: >= 1.55.0

### 性能考虑

**预期执行时间**:
- gofmt: < 1 秒
- goimports: < 2 秒
- golangci-lint: 5-30 秒（取决于变更文件数量）
- 总计: 通常 < 30 秒

**优化策略**:
- golangci-lint 只检查变更的文件（默认行为）
- 使用 `--fix` 自动修复简单问题
- 设置合理的超时时间

**开发者可选操作**:
```bash
# 跳过 hooks（紧急情况）
git commit --no-verify

# 只运行特定 hook
pre-commit run golangci-lint --files file.go

# 临时禁用某个 hook
SKIP=golangci-lint git commit -m "..."
```

## Migration Plan

### Phase 1: 准备阶段
1. ✅ 创建 `.pre-commit-config.yaml`
2. ✅ 更新 README.md（添加开发设置说明）
3. ✅ 可选：创建 CONTRIBUTING.md

### Phase 2: 修复现有问题
1. 运行 golangci-lint 识别所有问题
2. 自动修复可修复的问题
3. 手动修复剩余问题
4. 验证所有测试通过
5. 提交修复

### Phase 3: 部署
1. 合并 pre-commit 配置到主分支
2. 通知团队成员安装 pre-commit
3. 更新项目文档

### Phase 4: 验证
1. 监控 CI 中的 lint job 通过率
2. 收集开发者反馈
3. 根据需要调整配置

### Rollback Plan
如果出现问题：
1. 开发者可以使用 `--no-verify` 跳过 hooks
2. 可以删除 `.pre-commit-config.yaml` 回退
3. CI 检查不受影响，保证代码质量底线

## Risks / Trade-offs

### Risks

1. **开发者采纳率**
   - 风险：部分开发者可能不安装 pre-commit
   - 缓解：清晰的文档 + CI 作为后备

2. **首次运行耗时**
   - 风险：首次 commit 需要下载 hooks
   - 缓解：文档中建议先运行 `pre-commit run --all-files`

3. **工具版本不一致**
   - 风险：不同开发者的 golangci-lint 版本可能不同
   - 缓解：在 mise.toml 中固定版本（可选）

### Trade-offs

**选择**: 本地 hooks + CI 检查
**代价**:
- 开发者需要额外安装工具
- 提交时间略有增加（~30秒）

**收益**:
- 即时反馈，减少 CI 失败
- 更好的代码质量
- 减少代码审查负担

## Open Questions

1. **是否在 mise.toml 中添加工具？**
   - 选项 A: 添加 pre-commit, golangci-lint
   - 选项 B: 仅在文档中说明安装方法
   - 建议：选项 A，统一开发环境

2. **是否在 CI 中也运行 pre-commit？**
   - 选项 A: 添加 pre-commit CI job
   - 选项 B: 保持现有 golangci-lint job
   - 建议：选项 B，避免重复

3. **是否添加更多 hooks？**
   - 可选：go test, go vet, staticcheck
   - 建议：暂不添加，保持简单快速

## References

- Pre-commit 官方文档: https://pre-commit.com/
- Pre-commit hooks 列表: https://pre-commit.com/hooks.html
- golangci-lint 配置: https://golangci-lint.run/usage/configuration/
- Go code review comments: https://github.com/golang/go/wiki/CodeReviewComments
