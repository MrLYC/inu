# Inu 🐕

[![CI](https://github.com/MrLYC/inu/actions/workflows/ci.yml/badge.svg)](https://github.com/MrLYC/inu/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/go-1.24.4-blue.svg)](https://golang.org/dl/)

**Inu** 是一个基于 AI 大模型的文本敏感信息匿名化工具。它能够智能识别文本中的个人信息、联系方式、组织名称等敏感实体，将其替换为可追溯的占位符，并支持后续还原。

## ✨ 特性

- 🤖 **AI 驱动**：基于大语言模型（LLM）的智能实体识别
- 🔒 **安全可靠**：敏感信息完全匿名化，保护隐私
- 🔄 **可逆转换**：支持将匿名化文本还原为原始内容
- 🎯 **精准识别**：支持多种实体类型（人名、联系方式、地址、IP 等）
- 🌐 **灵活配置**：支持自定义 LLM API endpoint（兼容 OpenAI API）

## 📦 安装

### 从 Release 下载

从 [GitHub Releases](https://github.com/MrLYC/inu/releases) 下载适合你平台的预编译二进制文件：

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

## 🚀 快速开始

### 配置环境变量

Inu 需要连接到 OpenAI API（或兼容的服务）。请设置以下环境变量：

```bash
export OPENAI_API_KEY="your-api-key"
export OPENAI_MODEL_NAME="gpt-4"
export OPENAI_BASE_URL="https://api.openai.com/v1"  # 可选，默认为 OpenAI
```

### 基本使用

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
    
    // 创建匿名化器
    anon, err := anonymizer.New(llm)
    if err != nil {
        log.Fatal(err)
    }
    
    // 匿名化文本
    text := "张三的身份证号是 110101199001011234，他的电话号码是 13800138000。"
    types := []string{"人名", "联系方式", "职务", "密码", "组织", "地址", "文件", "账号", "网址", "IP"}
    
    result, entities, err := anon.AnonymizeText(ctx, types, text)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("匿名化结果: %s", result)
    // 输出: <人名[0].姓名.张三> 的身份证号是 <联系方式[0].身份证.110101199001011234>...
    
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

Inu 可以识别和匿名化以下类型的敏感信息：

- **人名**：个人姓名
- **联系方式**：电话号码、邮箱地址、身份证号等
- **职务**：职位、头衔
- **密码**：密码、密钥
- **组织**：公司、机构名称
- **地址**：物理地址、邮寄地址
- **文件**：文件名、路径
- **账号**：用户名、账户ID
- **网址**：URL、域名
- **IP**：IP 地址

## 🛠️ 开发

### 项目结构

```
inu/
├── cmd/inu/           # CLI 入口
├── pkg/anonymizer/    # 核心匿名化逻辑
├── bin/               # 编译产物（不提交）
├── openspec/          # OpenSpec 规范和变更提案
└── .github/           # GitHub Actions workflows
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

需要安装 [golangci-lint](https://golangci-lint.run/usage/install/)：

```bash
make lint
```

## 📋 路线图

- [ ] CLI 命令行工具（交互式和批处理模式）
- [ ] Web 界面
- [ ] 支持更多 LLM 提供商
- [ ] 批量文件处理
- [ ] 自定义实体类型
- [ ] 配置文件支持

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
- [Eris](https://github.com/rotisserie/eris) - Go 错误处理库

## 📬 联系方式

- GitHub Issues: [https://github.com/MrLYC/inu/issues](https://github.com/MrLYC/inu/issues)
- Author: [@MrLYC](https://github.com/MrLYC)

---

Made with ❤️ by [MrLYC](https://github.com/MrLYC)
