# Design: Web API Implementation

## Overview
本设计文档描述如何为 Inu 添加基于 Gin 框架的 Web API，将现有的文本匿名化和还原功能封装为 RESTful API 端点。

## Architecture

### Component Structure
```
inu
├── cmd/inu/
│   ├── main.go (注册 web 命令)
│   └── commands/
│       ├── anonymize.go (现有 CLI)
│       ├── restore.go (现有 CLI)
│       └── web.go (新增 Web 命令)
└── pkg/
    ├── anonymizer/ (现有核心逻辑)
    │   ├── anonymizer.go
    │   └── llm.go
    ├── cli/ (现有 CLI 工具)
    └── web/ (新增 Web 层)
        ├── server.go (服务器实例和路由)
        ├── config.go (配置结构)
        ├── handlers/
        │   ├── anonymize.go (匿名化 API)
        │   ├── restore.go (还原 API)
        │   └── health.go (健康检查)
        └── middleware/
            └── auth.go (身份认证)
```

### Request Flow
```
Client Request
    ↓
Gin Router
    ↓
Auth Middleware (验证 admin token)
    ↓
Handler (anonymize/restore)
    ↓
pkg/anonymizer (复用核心逻辑)
    ↓
LLM API (OpenAI)
    ↓
Response to Client
```

## Key Design Decisions

### 1. Web 框架选择：Gin
**决策**: 使用 `github.com/gin-gonic/gin`

**理由**:
- 高性能、轻量级
- 广泛使用、社区活跃
- 内置丰富的中间件和工具
- API 设计清晰、易于测试
- 与 CloudWeGo 生态兼容性好

**替代方案考虑**:
- **Echo**: 性能相近，但生态略小
- **Fiber**: 高性能但 API 设计与标准库差异大
- **标准库 net/http**: 需要更多手动路由和中间件实现

### 2. 身份认证方案：HTTP Basic Auth
**决策**: 使用 HTTP Basic Authentication

**理由**:
- 简单直接，适合管理员 API
- 无需额外状态管理（无 session/cookie）
- 易于测试和调试
- 适合内部服务和受信任网络

**安全考虑**:
- 建议部署时使用 HTTPS（TLS）
- token 应使用强随机值
- 未来可扩展支持 JWT/OAuth2

**替代方案**:
- **JWT**: 更复杂，适合多用户场景（当前不需要）
- **API Key in Header**: 与 Basic Auth 类似，但缺少标准
- **无认证**: 不安全，不适合生产环境

### 3. API 设计模式：RESTful
**决策**: 采用 RESTful API 设计

**端点设计**:
```
POST /api/v1/anonymize
POST /api/v1/restore
GET  /health
```

**理由**:
- 清晰的资源和操作映射
- 标准的 HTTP 方法语义
- 易于理解和集成
- 支持 API 版本控制（/v1）

**请求/响应格式**: JSON
- 广泛支持
- 易于调试
- Gin 内置优秀的 JSON 处理

### 4. 错误处理策略
**HTTP 状态码映射**:
- `200 OK`: 成功处理
- `400 Bad Request`: 输入验证失败（空文本、无效格式）
- `401 Unauthorized`: 认证失败
- `500 Internal Server Error`: LLM 调用失败或内部错误

**错误响应格式**:
```json
{
  "error": "invalid_input",
  "message": "Text cannot be empty",
  "code": 400
}
```

### 5. 服务器生命周期管理
**决策**: 支持优雅关闭（Graceful Shutdown）

**实现**:
- 监听 `SIGINT` 和 `SIGTERM` 信号
- 使用 `context.WithTimeout` 设置关闭超时（默认 5 秒）
- 等待正在处理的请求完成
- 关闭数据库连接、LLM 客户端等资源

**理由**:
- 避免请求中断导致数据不一致
- 云环境友好（K8s, Docker）
- 生产环境最佳实践

## API Specification

### POST /api/v1/anonymize
匿名化文本中的敏感信息。

**Request**:
```json
{
  "text": "张三的电话是 13800138000",
  "entity_types": ["个人信息"]  // 可选，默认使用全部类型
}
```

**Response** (200 OK):
```json
{
  "anonymized_text": "<个人信息[0].姓名.全名>的电话是<个人信息[1].电话.号码>",
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

**Error Response** (400 Bad Request):
```json
{
  "error": "invalid_input",
  "message": "Text cannot be empty",
  "code": 400
}
```

### POST /api/v1/restore
还原匿名化的文本。

**Request**:
```json
{
  "anonymized_text": "<个人信息[0].姓名.全名>的电话是<个人信息[1].电话.号码>",
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

**Response** (200 OK):
```json
{
  "restored_text": "张三的电话是 13800138000"
}
```

**Error Response** (400 Bad Request):
```json
{
  "error": "invalid_input",
  "message": "Entities cannot be empty",
  "code": 400
}
```

### GET /health
健康检查端点（无需认证）。

**Response** (200 OK):
```json
{
  "status": "ok",
  "version": "v0.1.0"
}
```

## Configuration

### Command Line Flags
```bash
inu web \
  --addr "127.0.0.1:8080" \      # 监听地址和端口
  --admin-user "admin" \          # 管理员用户名
  --admin-token "secret123"       # 管理员密码（token）
```

**默认值**:
- `--addr`: `127.0.0.1:8080` (本地监听，安全默认值)
- `--admin-user`: `admin`
- `--admin-token`: 无默认值（必须指定）

**环境变量支持**:
- 复用现有的 LLM 配置：`OPENAI_API_KEY`, `OPENAI_MODEL_NAME`, `OPENAI_BASE_URL`
- 可选支持：`INU_WEB_ADDR`, `INU_ADMIN_USER`, `INU_ADMIN_TOKEN`

## Testing Strategy

### Unit Tests
- **Handlers**: 使用 `httptest` 测试 HTTP handlers
  - Mock `anonymizer.Anonymizer` 接口
  - 测试输入验证、错误处理、响应格式
- **Middleware**: 测试认证中间件的各种场景
- **Server**: 测试服务器启动和关闭逻辑

### Integration Tests
- 启动完整的 Web 服务器
- 使用真实的 HTTP 客户端发送请求
- Mock LLM API（使用 `pkg/anonymizer/mock_llm_test.go`）
- 验证端到端流程

### Manual Testing
```bash
# 1. 启动服务器
inu web --admin-token test123

# 2. 健康检查（无需认证）
curl http://localhost:8080/health

# 3. 匿名化（需要认证）
curl -X POST http://localhost:8080/api/v1/anonymize \
  -u admin:test123 \
  -H "Content-Type: application/json" \
  -d '{"text": "张三的电话是 13800138000"}'

# 4. 还原（需要认证）
curl -X POST http://localhost:8080/api/v1/restore \
  -u admin:test123 \
  -H "Content-Type: application/json" \
  -d '{
    "anonymized_text": "<个人信息[0].姓名.全名>的电话是<个人信息[1].电话.号码>",
    "entities": [...]
  }'
```

## Performance Considerations

### LLM 连接池
- 复用单个 `anonymizer.Anonymizer` 实例
- LLM 客户端在服务器启动时初始化一次
- 避免每次请求重新创建连接

### 并发处理
- Gin 默认使用 goroutine 池处理请求
- 无需额外配置，自动支持并发

### 超时控制
- 建议为 LLM API 调用设置超时（如 30 秒）
- HTTP 请求超时通过 context 传递到 LLM 层

## Security Considerations

### 1. 认证强制
- 所有 `/api/*` 端点必须经过认证
- `/health` 端点可公开访问（健康检查需要）

### 2. HTTPS 建议
- 生产环境必须使用 HTTPS
- 可通过反向代理（Nginx, Caddy）提供 TLS

### 3. 速率限制（未来扩展）
- 考虑添加限流中间件防止滥用
- 可基于 IP 或用户限制请求频率

### 4. 输入验证
- 限制请求 body 大小（默认 Gin 限制为 8MB）
- 验证文本长度（避免超大输入导致 LLM 超时）

## Future Enhancements

### 1. 异步处理
- 对于大文本，支持异步处理模式
- 返回任务 ID，客户端轮询结果

### 2. 批量处理
- 支持一次请求处理多个文本
- `POST /api/v1/batch/anonymize`

### 3. WebSocket 支持
- 实时流式处理
- 适合大文件和交互式场景

### 4. 多用户支持
- 支持多个 API token
- 用户级别的配额和限流

### 5. Metrics 和监控
- 添加 Prometheus metrics 端点
- 记录请求量、延迟、错误率

## Migration Path
本变更完全向后兼容：
- 现有 CLI 命令 (`anonymize`, `restore`) 保持不变
- Web API 作为新功能添加
- 可以在同一二进制文件中同时支持 CLI 和 Web 模式
