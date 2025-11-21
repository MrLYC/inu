# Add Web API

## Why
当前 Inu 仅支持 CLI 命令行接口，限制了其集成和使用场景。许多用户希望：
- 通过 HTTP API 集成到现有系统中（如日志处理管道、数据处理服务）
- 提供 Web 服务供多个客户端共享（避免每个客户端都配置 LLM credentials）
- 支持远程调用和微服务架构
- 减少 CLI 启动开销，通过常驻进程提升性能

添加 Web API 将使 Inu 更加灵活和易于集成，同时保持 CLI 接口向后兼容。

## What Changes
- 添加 `inu web` 子命令，启动 HTTP API 服务器
- 使用 Gin 框架提供 RESTful API
- 将现有的 `anonymize` 和 `restore` 功能封装为 HTTP 端点
- 提供基本的身份认证机制（admin token）
- 支持可配置的监听地址和端口

**新增命令**:
```bash
inu web [--addr 127.0.0.1:8080] [--admin-user admin] [--admin-token <token>]
```

**新增 API 端点**:
- `POST /api/v1/anonymize` - 脱敏文本
- `POST /api/v1/restore` - 还原文本
- `GET /health` - 健康检查

## Impact
- **影响的 specs**: 需要新增 `web-api` spec
- **影响的代码**:
  - 新增 `cmd/inu/commands/web.go` - Web 命令实现
  - 新增 `pkg/web/` - Web 服务器和 API handlers
    - `pkg/web/server.go` - Gin 服务器初始化
    - `pkg/web/handlers/anonymize.go` - 脱敏 API handler
    - `pkg/web/handlers/restore.go` - 还原 API handler
    - `pkg/web/middleware/auth.go` - 身份认证中间件
  - 修改 `cmd/inu/main.go` - 注册 web 命令
- **依赖变更**: 
  - 新增 `github.com/gin-gonic/gin` (Web 框架)
- **不破坏现有功能**: CLI 命令保持完全兼容
