# Inu 🐕

[![CI](https://github.com/MrLYC/inu/actions/workflows/ci.yml/badge.svg)](https://github.com/MrLYC/inu/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/go-1.24.4-blue.svg)](https://golang.org/dl/)

**Inu** 是一个基于 AI 大模型的文本敏感信息脱敏工具。它能够智能识别文本中的个人信息、联系方式、组织名称等敏感实体，将其替换为可追溯的占位符，并支持后续还原。

## ✨ 特性

- 🤖 **AI 驱动**：基于大语言模型（LLM）的智能实体识别
- 🔒 **安全可靠**：敏感信息完全脱敏，保护隐私
- 🔄 **可逆转换**：支持将脱敏文本还原为原始内容
- 🎯 **精准识别**：支持多种实体类型（人名、联系方式、地址、IP 等）
- 🌐 **灵活配置**：支持自定义 LLM API endpoint（兼容 OpenAI API）
- 🖥️ **CLI + Web API**：同时支持命令行工具和 HTTP API 服务

## 📦 安装

### 从 Release 下载

从 [GitHub Releases](https://github.com/MrLYC/inu/releases) 下载适合你平台的预编译二进制文件。

**单文件分发**：二进制文件已包含完整的 Web UI，无需额外的静态文件或配置，下载即可使用。

**Linux (amd64)**
```bash
curl -LO https://github.com/MrLYC/inu/releases/latest/download/inu-linux-amd64.tar.gz
tar xzf inu-linux-amd64.tar.gz
sudo mv inu /usr/local/bin/
```

**macOS (Apple Silicon)**
```bash
curl -LO https://github.com/MrLYC/inu/releases/latest/download/inu-darwin-arm64.tar.gz
tar xzf inu-darwin-arm64.tar.gz
sudo mv inu /usr/local/bin/
```

**macOS (Intel)**
```bash
curl -LO https://github.com/MrLYC/inu/releases/latest/download/inu-darwin-amd64.tar.gz
tar xzf inu-darwin-amd64.tar.gz
sudo mv inu /usr/local/bin/
```

### 从源码编译

要求：Go 1.24 或更高版本

```bash
git clone https://github.com/MrLYC/inu.git
cd inu
make build
```

编译后的二进制文件位于 `bin/inu`。

## �️ 开发环境配置

如果你想参与 Inu 的开发或贡献代码，需要配置本地开发环境。

### 安装 Pre-commit

Pre-commit 会在每次 git commit 前自动运行代码格式化和质量检查，确保代码符合项目标准。

**安装 pre-commit**：
```bash
# 方式 1: 使用 pip
pip3 install pre-commit

# 方式 2: 使用 Homebrew (macOS)
brew install pre-commit

# 方式 3: 使用 mise (如果项目使用 mise)
mise use -g pre-commit@latest
```

**安装 Git hooks**：
```bash
cd inu
pre-commit install
```

### 安装代码质量工具

**安装 goimports**（整理 Go 导入语句）：
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

**安装 golangci-lint**（代码静态检查）：
```bash
# 方式 1: 使用 Homebrew (macOS)
brew install golangci-lint

# 方式 2: 使用 go install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 方式 3: 使用 mise
mise use -g golangci-lint@latest
```

### 使用 Pre-commit Hooks

安装完成后，每次 `git commit` 时会自动运行以下检查：
- 移除行尾空格
- 确保文件以换行结束
- 检查 YAML 语法
- 运行 `gofmt` 格式化代码
- 运行 `goimports` 整理导入
- 运行 `golangci-lint` 进行代码质量检查

**手动运行所有检查**：
```bash
pre-commit run --all-files
```

**跳过 hooks（紧急情况）**：
```bash
git commit --no-verify -m "urgent fix"
```

### 疑难解答

如果遇到 `goimports` 或 `golangci-lint` 命令找不到的错误：

1. **检查 GOPATH**：
   ```bash
   echo $GOPATH
   # 应该输出类似 /Users/username/.go
   ```

2. **确认工具路径**：
   ```bash
   ls $HOME/.go/bin
   # 应该看到 goimports 和 golangci-lint
   ```

3. **添加到 PATH**（如果需要）：
   ```bash
   export PATH="$HOME/.go/bin:$PATH"
   # 或将上面的命令添加到 ~/.zshrc 或 ~/.bashrc
   ```

4. **验证安装**：
   ```bash
   goimports -version
   golangci-lint --version
   ```

## �🚀 快速开始

### 配置环境变量

Inu 需要连接到 OpenAI API（或兼容的服务）。请设置以下环境变量：

```bash
export OPENAI_API_KEY="your-api-key"
export OPENAI_MODEL_NAME="gpt-4"
export OPENAI_BASE_URL="https://api.openai.com/v1"  # 可选，默认为 OpenAI
```

### 命令行使用

#### 脱敏文本

从文件读取并输出到标准输出：
```bash
inu anonymize --file input.txt --output-entities entities.yaml
```

保存到文件（默认也输出到标准输出）：
```bash
inu anonymize --file input.txt --output anonymized.txt --output-entities entities.yaml
```

只保存到文件，不打印到标准输出：
```bash
inu anonymize --file input.txt --no-print --output anonymized.txt --output-entities entities.yaml
```

从命令行参数读取：
```bash
inu anonymize --content "张三的电话是 13800138000"
```

从标准输入（默认输出到标准输出和标准错误）：
```bash
echo "李四住在北京市朝阳区" | inu anonymize --output-entities entities.yaml
```

使用管道（entities 信息输出到 stderr）：
```bash
cat input.txt | inu anonymize --output-entities entities.yaml > anonymized.txt 2> entities.log
```

指定实体类型：
```bash
inu anonymize --file input.txt --entity-types "个人信息,业务信息,资产信息"
```

#### 还原文本

从文件读取并输出到标准输出：
```bash
inu restore --file anonymized.txt --entities entities.yaml
```

保存到文件（默认也输出到标准输出）：
```bash
inu restore --file anonymized.txt --entities entities.yaml --output restored.txt
```

只保存到文件，不打印到标准输出：
```bash
inu restore --file anonymized.txt --entities entities.yaml --no-print --output restored.txt
```

使用管道：
```bash
cat anonymized.txt | inu restore --entities entities.yaml > restored.txt
```

#### 交互式工作流

`interactive` 命令提供了一个便捷的交互式流程，特别适合与 ChatGPT 等外部工具配合使用：

**基本用法**
```bash
inu interactive --file sensitive.txt

# 工作流程：
# 1. 命令输出脱敏文本（带分隔线标识）
# 2. 复制脱敏文本，粘贴到 ChatGPT 进行处理（如总结、翻译）
# 3. 复制 ChatGPT 的回复
# 4. 粘贴回终端
# 5. 按 Ctrl+D (Unix) 或 Ctrl+Z (Windows)
# 6. 命令输出还原后的文本（带分隔线标识）
# 7. 可以继续步骤 3-6 进行多次处理
```

**从命令行参数输入**
```bash
inu interactive -c "张三在 ABC 公司工作"
```

**精简提示模式**
```bash
inu interactive -f input.txt --no-prompt

# 减少提示信息，适合脚本化使用
```

**指定实体类型**
```bash
inu interactive -f input.txt --entity-types "个人信息,业务信息"

# 只脱敏指定类型的实体
```

**显示效果示例**：

```
$ inu interactive -c "张三的电话是 13800138000"

============================================================
ANONYMIZED TEXT:
============================================================
<个人信息[0].姓名.全名>的电话是<个人信息[1].电话.号码>
============================================================

------------------------------------------------------------
✅ Anonymization Complete
------------------------------------------------------------
Next steps:
  1. Copy the anonymized text above
  2. Process it externally (e.g., paste to ChatGPT)
  3. Paste the processed text below
  4. Press Ctrl+D (Unix) or Ctrl+Z (Windows) to restore
------------------------------------------------------------

📝 Paste your processed text here:
[此处粘贴处理后的文本]
^D

============================================================
RESTORED TEXT:
============================================================
张三的电话是 13800138000
============================================================

📝 Ready for next input (Ctrl+D to restore, Ctrl+C to exit)
```

**典型工作流示例**：

1. **与 ChatGPT 进行文档总结**
   ```bash
   $ inu interactive -f confidential-report.txt
   # [看到脱敏文本，有清晰的分隔线]
   # [复制到 ChatGPT: "请总结这份报告"]
   # [粘贴 ChatGPT 的总结]
   # [按 Ctrl+D]
   # [得到还原后的总结，有清晰的分隔线]
   ```

2. **多次处理流程**
   ```bash
   $ inu interactive -f report.txt
   # 第一次：获取总结
   [粘贴到 ChatGPT 获取总结]
   ^D
   [得到还原的总结]

   # 第二次：翻译总结
   [粘贴总结到 ChatGPT 获取翻译]
   ^D
   [得到还原的翻译]

   # 按 Ctrl+C 退出
   ```

**优势**：
- 🔒 **保护隐私**：敏感信息不会泄露给 ChatGPT
- 🔄 **支持多次处理**：一次脱敏，多次使用不同 AI 服务
- 🎯 **简单直观**：无需管理中间文件
- 💾 **实体在内存**：整个流程在一个进程中完成
- 📊 **清晰显示**：使用分隔线明确区分输入输出内容

#### 从旧版本迁移

**⚠️ Breaking Changes in v0.2.0**

从 v0.2.0 开始，CLI 输出行为已更改，遵循 Unix 标准约定：

**旧版本（v0.1.x）**：
- 默认不输出到 stdout，需要使用 `--print` 参数
- 使用 `--print-entities` 输出 entities

**新版本（v0.2.0+）**：
- **默认输出到 stdout**（标准输出）
- **entities 默认输出到 stderr**（标准错误）
- 使用 `--no-print` 参数来**抑制** stdout 输出
- 移除了 `--print` 和 `--print-entities` 参数

**迁移示例**：

```bash
# 旧版本：
inu anonymize -f input.txt -o output.txt --print --print-entities

# 新版本（等效）：
inu anonymize -f input.txt -o output.txt
# 输出到 stdout 和 output.txt，entities 到 stderr

# 如果只想要文件输出（旧版本的默认行为）：
inu anonymize -f input.txt -o output.txt --no-print
```

这个改变使 `inu` 更符合 Unix 哲学，更容易在管道中使用：
```bash
# 现在可以直接这样使用：
cat input.txt | inu anonymize | tee anonymized.txt

# entities 可以重定向到日志文件：
cat input.txt | inu anonymize 2> entities.log > anonymized.txt
```

### Web API 使用

#### 启动 Web 服务器

**需要认证的方式（推荐）：**
```bash
inu web --admin-token your-secret-token
```

**不需要认证的方式（仅开发环境）：**
```bash
inu web
```

⚠️ **警告**：不使用 `--admin-token` 参数时，服务器将运行在无认证模式下，任何人都可以访问和使用 API。这仅适用于本地开发环境，生产环境中务必启用认证！

使用自定义地址、端口和实体类型：
```bash
inu web --addr 0.0.0.0:9090 \
  --admin-user admin \
  --admin-token your-secret-token \
  --entity-types "PERSON,ORG,EMAIL,PHONE,ADDRESS"
```

服务器启动后，可以通过 Web 界面或 HTTP API 进行脱敏和还原操作。

**部署说明**：
- 二进制文件包含完整的 Web UI，无需部署额外的静态文件
- 可以将单个二进制文件复制到任何目录直接运行
- 适合容器化部署和离线环境

#### Web 界面

打开浏览器访问 `http://localhost:8080/`。

**认证说明：**
- 如果启动服务器时设置了 `--admin-token`，访问时需要输入用户名和密码
- 如果未设置 `--admin-token`，无需认证即可直接访问

**功能特性：**
- 🎨 **双视图模式**：脱敏视图和还原视图
- 🔄 **实时处理**：即时脱敏和还原文本
- 💾 **会话状态**：自动保存实体映射到浏览器会话
- 📱 **响应式设计**：支持桌面端和移动端
- 🎯 **自定义实体类型**：支持添加自定义实体类型
- 🔒 **安全认证**：可选的 Basic Auth 认证保护

**使用流程：**

1. **脱敏文本**
   - 选择要识别的实体类型（支持多选）
   - 在左侧输入框输入原始文本
   - 点击"脱敏"按钮
   - 在右侧查看脱敏结果
   - 实体映射自动保存到浏览器会话

2. **还原文本**
   - 点击"切换到还原模式"
   - 查看顶部的实体映射
   - 在右侧输入框粘贴需要还原的文本
   - 点击"还原"按钮
   - 查看还原后的结果
   - 可以多次还原不同的文本（实体保留）

3. **多次处理**
   - 一次脱敏，可以多次还原不同文本
   - 刷新页面后状态自动恢复（使用 sessionStorage）
   - 关闭标签页后数据自动清除

#### API 端点

**健康检查（无需认证）**
```bash
curl http://localhost:8080/health
```

响应：
```json
{
  "status": "ok",
  "version": "v0.1.0"
}
```

**获取配置（需要认证）**
```bash
curl http://localhost:8080/api/v1/config \
  -u admin:your-secret-token
```

响应：
```json
{
  "entity_types": ["PERSON", "ORG", "EMAIL", "PHONE", "ADDRESS"]
}
```

**脱敏文本（需要认证）**
```bash
curl -X POST http://localhost:8080/api/v1/anonymize \
  -u admin:your-secret-token \
  -H "Content-Type: application/json" \
  -d '{
    "text": "张三的电话是 13800138000"
  }'
```

响应：
```json
{
  "anonymized_text": "<个人信息[0].姓名.全名>的电话是 <个人信息[1].电话.号码>",
  "entities": [
    {
      "key": "<个人信息[0].姓名.全名>",
      "type": "个人信息",
      "id": "0",
      "category": "姓名",
      "detail": "张三",
      "values": ["张三"]
    },
    {
      "key": "<个人信息[1].电话.号码>",
      "type": "个人信息",
      "id": "1",
      "category": "电话",
      "detail": "13800138000",
      "values": ["13800138000"]
    }
  ]
}
```

**指定实体类型**
```bash
curl -X POST http://localhost:8080/api/v1/anonymize \
  -u admin:your-secret-token \
  -H "Content-Type: application/json" \
  -d '{
    "text": "张三在 ABC 公司工作",
    "entity_types": ["个人信息"]
  }'
```

**还原文本（需要认证）**
```bash
curl -X POST http://localhost:8080/api/v1/restore \
  -u admin:your-secret-token \
  -H "Content-Type: application/json" \
  -d '{
    "anonymized_text": "<个人信息[0].姓名.全名>的电话是 <个人信息[1].电话.号码>",
    "entities": [
      {
        "key": "<个人信息[0].姓名.全名>",
        "values": ["张三"]
      },
      {
        "key": "<个人信息[1].电话.号码>",
        "values": ["13800138000"]
      }
    ]
  }'
```

响应：
```json
{
  "restored_text": "张三的电话是 13800138000"
}
```

#### 身份认证

**Web 界面和 API 端点**的认证是可选的，取决于启动服务器时是否设置了 `--admin-token`。

**启用认证：**
```bash
inu web --admin-token your-secret-token
```

- **Web 界面**：浏览器会弹出认证对话框，输入用户名和密码
- **API 调用**：使用 `-u username:password` 或设置 `Authorization` 头

```bash
# 方式 1：使用 -u 参数
curl -u admin:your-secret-token http://localhost:8080/api/v1/anonymize ...

# 方式 2：使用 Authorization 头
curl -H "Authorization: Basic $(echo -n 'admin:your-secret-token' | base64)" \
  http://localhost:8080/api/v1/anonymize ...
```

**禁用认证（不推荐）：**
```bash
inu web
```

所有端点无需认证即可访问。⚠️ 仅用于本地开发环境！

**注意**：生产环境中建议同时使用 HTTPS 和认证保护。

### 编程接口

```go
package main

import (
    "context"
    "log"

    "github.com/mrlyc/inu/pkg/anonymizer"
)

func main() {
    ctx := context.Background()

    // 创建 LLM 客户端
    llm, err := anonymizer.CreateOpenAIChatModel(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // 创建脱敏器
    anon, err := anonymizer.NewHashHidePair(llm)
    if err != nil {
        log.Fatal(err)
    }

    // 脱敏文本
    text := "张三的身份证号是 110101199001011234，他的电话号码是 13800138000。"
    types := []string{"个人信息", "业务信息", "资产信息", "账户信息", "位置数据", "文档名称", "组织机构", "岗位称谓"}

    result, entities, err := anon.AnonymizeText(ctx, types, text)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("脱敏结果: %s", result)
    // 输出: <个人信息[0].姓名.全名> 的身份证号是 <个人信息[1].身份证.110101199001011234>...

    // 还原文本
    restored, err := anon.RestoreText(ctx, entities, result)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("还原结果: %s", restored)
    // 输出: 张三的身份证号是 110101199001011234，他的电话号码是 13800138000。
}
```

## 📖 支持的实体类型

Inu 默认识别以下类型的敏感信息：

- **个人信息**：姓名、身份证号、电话号码等
- **业务信息**：业务数据、客户信息等
- **资产信息**：财产、资源信息等
- **账户信息**：银行账号、信用卡号等
- **位置数据**：地址、地理位置等
- **文档名称**：文件名、文档标题等
- **组织机构**：公司名称、机构名称等
- **岗位称谓**：职位、头衔等

你也可以通过 `--entity-types` 参数自定义要识别的实体类型。

## 🛠️ 开发

### 开发环境设置

本项目使用 [pre-commit](https://pre-commit.com/) 来保证代码质量，在提交代码前自动运行格式化和 lint 检查。

#### 安装 pre-commit

**macOS / Linux (推荐使用 mise):**
```bash
mise install  # 如果项目配置了 mise.toml
```

**macOS (使用 Homebrew):**
```bash
brew install pre-commit
```

**通用方式 (使用 pip):**
```bash
pip install pre-commit
# 或
pip3 install pre-commit
```

#### 安装 Go 工具

确保安装了以下 Go 代码质量工具：

```bash
# 安装 goimports（导入排序）
go install golang.org/x/tools/cmd/goimports@latest

# 安装 golangci-lint（代码检查）
# macOS:
brew install golangci-lint

# Linux / Windows:
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

#### 初始化 pre-commit hooks

在项目根目录执行：

```bash
pre-commit install
```

这会在本地 `.git/hooks/` 目录中安装 Git hooks。之后每次 `git commit` 时，hooks 会自动运行。

#### 使用方法

**正常提交（自动运行 hooks）：**
```bash
git add .
git commit -m "your commit message"
# pre-commit 会自动运行，如果有问题会阻止提交
```

**手动运行所有 hooks：**
```bash
pre-commit run --all-files
```

**跳过 hooks（不推荐）：**
```bash
git commit --no-verify -m "your message"
```

**只运行特定 hook：**
```bash
pre-commit run gofmt --all-files
pre-commit run golangci-lint --all-files
```

#### Hooks 说明

pre-commit 会运行以下检查：

- **文件检查**
  - 去除行尾空白字符
  - 确保文件以换行符结尾
  - 检查 YAML 语法
  - 检查大文件（超过 1MB）
  - 检查是否有未解决的合并冲突

- **Go 代码检查**
  - `gofmt`: 自动格式化 Go 代码
  - `goimports`: 整理和优化导入语句
  - `golangci-lint`: 运行 lint 检查（支持自动修复）

如果 hook 自动修复了代码，你需要重新 `git add` 并再次提交。

#### 故障排除

**Hook 运行失败：**
```bash
# 查看详细错误信息
pre-commit run --all-files --verbose

# 清除缓存并重试
pre-commit clean
pre-commit run --all-files
```

**跳过特定文件：**

编辑 `.pre-commit-config.yaml`，在对应 hook 中添加 `exclude` 参数：
```yaml
- id: gofmt
  exclude: ^vendor/|^.openspec/
```

**更新 hooks 版本：**
```bash
pre-commit autoupdate
```

### 项目结构

```
inu/
├── cmd/inu/               # CLI 入口
│   └── commands/          # CLI 子命令（anonymize, restore, web）
├── pkg/
│   ├── anonymizer/        # 核心脱敏逻辑
│   ├── cli/               # CLI 工具函数（输入输出、实体管理）
│   └── web/               # Web API 服务器和 UI
│       ├── handlers/      # HTTP handlers（anonymize, restore, health, config）
│       ├── middleware/    # 认证中间件
│       └── static/        # Web UI 静态文件（HTML, CSS, JS）
├── bin/                   # 编译产物（不提交）
├── openspec/              # OpenSpec 规范和变更提案
├── .github/               # GitHub Actions workflows
└── .pre-commit-config.yaml  # Pre-commit hooks 配置
```

### 构建命令

```bash
make help           # 显示所有可用命令
make build          # 编译当前平台二进制文件
make build-all      # 交叉编译所有平台
make test           # 运行测试
make lint           # 代码检查
make clean          # 清理编译产物
```

### 测试

```bash
go test ./...
```

### 代码检查

**使用 pre-commit（推荐）：**
```bash
pre-commit run --all-files
```

**直接使用 golangci-lint：**
```bash
golangci-lint run --timeout=5m
# 或使用 Makefile
make lint
```

## 📋 路线图

- [x] 核心脱敏和还原功能
- [x] CLI 命令行工具（`inu anonymize` / `inu restore`）
- [x] 多种输入方式（文件、命令行参数、标准输入）
- [x] 实体 YAML 文件管理
- [x] CI/CD 自动化构建和发布
- [x] Web API 服务（`inu web`）
- [x] HTTP 身份认证
- [x] Web 界面（交互式脱敏和还原）
- [ ] 支持更多 LLM 提供商
- [ ] 批量文件处理
- [ ] 更丰富的配置文件支持
- [ ] 插件系统

## 🤝 贡献

欢迎贡献！请查看 [OpenSpec 规范](openspec/) 了解项目的开发流程。

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 Apache License 2.0 许可证。详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [CloudWeGo Eino](https://github.com/cloudwego/eino) - AI 工具链框架
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Viper](https://github.com/spf13/viper) - 配置管理
- [Eris](https://github.com/rotisserie/eris) - Go 错误处理库

## 📬 联系方式

- GitHub Issues: [https://github.com/MrLYC/inu/issues](https://github.com/MrLYC/inu/issues)
- Author: [@MrLYC](https://github.com/MrLYC)

---

Made with ❤️ by [MrLYC](https://github.com/MrLYC)
