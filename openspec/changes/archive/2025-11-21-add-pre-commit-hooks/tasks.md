# 任务：添加 Pre-commit Hooks 配置

## 1. 准备工作

### 1.1 识别现有 Lint 问题
- [ ] 确保本地安装了 golangci-lint
  ```bash
  # 检查安装
  golangci-lint --version

  # 如未安装，选择一种方式安装：
  brew install golangci-lint
  # 或
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

- [ ] 运行 golangci-lint 识别所有问题
  ```bash
  golangci-lint run --timeout=5m > lint-issues.txt 2>&1
  cat lint-issues.txt
  ```

- [ ] 分析问题并分类：
  - 记录问题总数
  - 按类型分组（格式、未使用、错误处理等）
  - 确定优先级（P0: 必须修复，P1: 应该修复，P2: 可以修复）

### 1.2 准备修复环境
- [ ] 确保安装了 goimports
  ```bash
  go install golang.org/x/tools/cmd/goimports@latest
  ```

- [ ] 创建新分支进行修复
  ```bash
  git checkout -b add-pre-commit-hooks
  ```

- [ ] 备份当前代码（可选）
  ```bash
  git stash
  ```

## 2. 修复现有代码质量问题

### 2.1 自动修复格式问题
- [ ] 运行 gofmt 格式化所有 Go 文件
  ```bash
  gofmt -w .
  ```

- [ ] 运行 goimports 整理导入
  ```bash
  goimports -w -local github.com/mrlyc/inu .
  ```

- [ ] 运行 golangci-lint 自动修复
  ```bash
  golangci-lint run --fix --timeout=5m
  ```

- [ ] 查看修改
  ```bash
  git diff
  ```

### 2.2 手动修复剩余问题
- [ ] 再次运行 golangci-lint 检查
  ```bash
  golangci-lint run --timeout=5m
  ```

- [ ] 对于每个剩余问题：
  - [ ] 理解问题根因
  - [ ] 手动修复代码
  - [ ] 确保修复不破坏功能
  - [ ] 提交修复（小批量提交，方便审查）

- [ ] 处理特殊情况：
  - 如果某些问题不应修复，在 `.golangci.yml` 中添加排除规则
  - 添加注释说明为什么排除

### 2.3 验证修复结果
- [ ] 运行所有测试
  ```bash
  go test ./... -v -race -cover
  ```

- [ ] 确保所有测试通过
- [ ] 验证 golangci-lint 无错误
  ```bash
  golangci-lint run --timeout=5m
  ```

- [ ] 验证代码编译
  ```bash
  make build
  ```

- [ ] 提交所有修复
  ```bash
  git add .
  git commit -m "fix: resolve all golangci-lint issues"
  ```

## 3. 创建 Pre-commit 配置

### 3.1 创建配置文件
- [ ] 在项目根目录创建 `.pre-commit-config.yaml`
  ```bash
  touch .pre-commit-config.yaml
  ```

- [ ] 添加配置内容（参考 design.md 中的完整配置）
  - [ ] 添加 pre-commit-hooks repo（基础文件检查）
  - [ ] 添加 local hooks（Go 工具）
  - [ ] 配置 gofmt hook
  - [ ] 配置 goimports hook
  - [ ] 配置 golangci-lint hook

### 3.2 测试配置
- [ ] 安装 pre-commit（如未安装）
  ```bash
  pip install pre-commit
  # 或
  brew install pre-commit
  ```

- [ ] 安装 hooks
  ```bash
  pre-commit install
  ```

- [ ] 测试所有 hooks
  ```bash
  pre-commit run --all-files
  ```

- [ ] 验证所有 hooks 通过
- [ ] 如果失败，调试并修复配置

### 3.3 调整配置（如需要）
- [ ] 根据测试结果调整超时时间
- [ ] 调整 hook 参数
- [ ] 确保 hook 执行时间合理（< 1 分钟）

## 4. 更新项目文档

### 4.1 更新 README.md
- [ ] 添加 "Development Setup" 或 "开发环境设置" 章节
- [ ] 包含以下内容：
  - [ ] Pre-commit 工具介绍
  - [ ] 安装 pre-commit 的多种方式
  - [ ] 安装 goimports 说明
  - [ ] 安装 golangci-lint 说明
  - [ ] 安装 hooks 的命令
  - [ ] 使用示例
  - [ ] 常见问题解答

### 4.2 更新 mise.toml（可选）
- [ ] 考虑是否添加工具到 mise.toml
  ```toml
  [tools]
  go = "1.24"
  golangci-lint = "latest"
  # pre-commit = "latest"  # 可选
  ```

- [ ] 测试 mise 配置
  ```bash
  mise install
  mise exec -- golangci-lint --version
  ```

### 4.3 创建 CONTRIBUTING.md（可选）
- [ ] 如果项目需要贡献指南，创建此文件
- [ ] 包含：
  - [ ] 代码风格要求
  - [ ] 提交流程
  - [ ] Pre-commit hooks 使用
  - [ ] 测试要求
  - [ ] PR 流程

## 5. 测试集成

### 5.1 本地测试
- [ ] 创建测试提交验证 hooks
  ```bash
  # 修改一个文件
  echo "// test" >> pkg/anonymizer/anonymizer.go
  git add pkg/anonymizer/anonymizer.go
  git commit -m "test: verify pre-commit hooks"
  ```

- [ ] 验证 hooks 自动运行
- [ ] 验证格式问题被自动修复
- [ ] 恢复测试修改
  ```bash
  git reset --soft HEAD~1
  git restore pkg/anonymizer/anonymizer.go
  ```

### 5.2 测试失败场景
- [ ] 故意引入格式问题测试 hook 能否捕获
  ```bash
  # 添加未使用的导入
  echo 'import "fmt"' | cat - pkg/anonymizer/anonymizer.go > temp && mv temp pkg/anonymizer/anonymizer.go
  git add pkg/anonymizer/anonymizer.go
  git commit -m "test: should fail"
  ```

- [ ] 验证 commit 被阻止
- [ ] 查看错误信息是否清晰
- [ ] 恢复修改
  ```bash
  git restore pkg/anonymizer/anonymizer.go
  ```

### 5.3 测试跳过功能
- [ ] 测试 `--no-verify` 选项
  ```bash
  git commit --no-verify -m "test: skip hooks"
  ```

- [ ] 验证提交成功（跳过 hooks）
- [ ] 重置提交
  ```bash
  git reset HEAD~1
  ```

## 6. CI 验证

### 6.1 推送到远程分支
- [ ] 确保所有本地测试通过
- [ ] 推送分支到 GitHub
  ```bash
  git push origin add-pre-commit-hooks
  ```

### 6.2 创建 Pull Request
- [ ] 在 GitHub 创建 PR
- [ ] 填写 PR 描述，说明变更内容：
  - 添加 pre-commit 配置
  - 修复所有 golangci-lint 问题
  - 更新开发文档

### 6.3 验证 CI 检查
- [ ] 等待 CI 运行完成
- [ ] 验证所有 CI jobs 通过：
  - [ ] Test job 通过
  - [ ] Lint job 通过（关键！）
  - [ ] Build matrix job 通过

- [ ] 如果 CI 失败：
  - [ ] 查看失败日志
  - [ ] 在本地重现问题
  - [ ] 修复并重新推送
  - [ ] 重复验证

## 7. 团队通知和培训

### 7.1 准备通知
- [ ] 准备通知文档/邮件，包括：
  - [ ] 变更说明
  - [ ] 安装步骤
  - [ ] 使用指南
  - [ ] 常见问题
  - [ ] 支持联系方式

### 7.2 发布变更
- [ ] 合并 PR 到 main 分支
- [ ] 通知团队成员
- [ ] 分享安装和使用指南
- [ ] 提供支持和答疑

### 7.3 收集反馈
- [ ] 监控团队采纳情况
- [ ] 收集使用反馈
- [ ] 记录常见问题
- [ ] 根据反馈优化配置和文档

## 8. 后续改进

### 8.1 监控效果
- [ ] 监控 PR 中 lint 失败率（应该降低）
- [ ] 观察代码质量指标
- [ ] 收集开发者体验反馈

### 8.2 可选增强
- [ ] 考虑添加更多 hooks：
  - [ ] go test（可选）
  - [ ] go vet（可选）
  - [ ] 安全检查（gosec）

- [ ] 考虑优化性能：
  - [ ] 使用缓存加速 golangci-lint
  - [ ] 调整 hook 并行执行

- [ ] 考虑在 CI 中集成 pre-commit
  ```yaml
  - name: Run pre-commit
    uses: pre-commit/action@v3.0.0
  ```

### 8.3 文档维护
- [ ] 根据实际使用更新文档
- [ ] 添加疑难解答章节
- [ ] 维护最佳实践指南

---

## 检查清单

### 完成标准
- [x] 所有 golangci-lint 问题已修复
- [x] `.pre-commit-config.yaml` 配置正确
- [x] 所有测试通过
- [x] CI 中 lint job 通过
- [x] README.md 包含开发设置说明
- [x] 本地测试验证 hooks 正常工作
- [x] PR 已创建并通过审查

### 验证命令
```bash
# 验证 golangci-lint
golangci-lint run --timeout=5m

# 验证测试
go test ./... -v -race -cover

# 验证构建
make build

# 验证 pre-commit
pre-commit run --all-files

# 验证格式
gofmt -l .
goimports -l -local github.com/mrlyc/inu .
```

### 成功指标
- golangci-lint 零错误零警告
- 所有单元测试通过
- CI 中所有 jobs 通过
- Pre-commit hooks 能够捕获常见问题
- 开发者反馈积极
