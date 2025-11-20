# Implementation Tasks

## 1. 项目配置和依赖
- [ ] 1.1 添加 Gin 框架依赖到 `go.mod`
- [ ] 1.2 更新 `openspec/project.md` 添加 Web API 相关技术栈说明

## 2. Web 服务器核心实现
- [ ] 2.1 创建 `pkg/web/server.go`
  - [ ] 实现 `Server` 结构体（持有 Anonymizer 实例和配置）
  - [ ] 实现 `NewServer(anonymizer, config)` 构造函数
  - [ ] 实现 `Start()` 方法启动 Gin 服务器
  - [ ] 实现 `Stop()` 方法优雅关闭
  - [ ] 配置 Gin 路由和中间件
- [ ] 2.2 创建 `pkg/web/config.go`
  - [ ] 定义 `Config` 结构体（Addr, AdminUser, AdminToken）
  - [ ] 实现配置验证逻辑

## 3. 身份认证中间件
- [ ] 3.1 创建 `pkg/web/middleware/auth.go`
  - [ ] 实现 `BasicAuth(adminUser, adminToken)` 中间件
  - [ ] 支持 HTTP Basic Authentication
  - [ ] 验证失败返回 401 Unauthorized
  - [ ] 添加日志记录认证失败事件

## 4. API Handlers 实现
- [ ] 4.1 创建 `pkg/web/handlers/anonymize.go`
  - [ ] 定义请求结构体 `AnonymizeRequest{Text, EntityTypes}`
  - [ ] 定义响应结构体 `AnonymizeResponse{AnonymizedText, Entities}`
  - [ ] 实现 `AnonymizeHandler(anonymizer)` 处理函数
  - [ ] 处理输入验证（空文本、无效类型）
  - [ ] 调用 `anonymizer.AnonymizeText()`
  - [ ] 处理 LLM 错误并返回适当的 HTTP 状态码
- [ ] 4.2 创建 `pkg/web/handlers/restore.go`
  - [ ] 定义请求结构体 `RestoreRequest{AnonymizedText, Entities}`
  - [ ] 定义响应结构体 `RestoreResponse{RestoredText}`
  - [ ] 实现 `RestoreHandler(anonymizer)` 处理函数
  - [ ] 处理输入验证
  - [ ] 调用 `anonymizer.RestoreText()`
  - [ ] 错误处理
- [ ] 4.3 创建 `pkg/web/handlers/health.go`
  - [ ] 实现 `HealthHandler()` - 返回 200 OK 和简单状态信息

## 5. Web 命令实现
- [ ] 5.1 创建 `cmd/inu/commands/web.go`
  - [ ] 定义命令行标志变量（webAddr, webAdminUser, webAdminToken）
  - [ ] 实现 `NewWebCmd()` 创建 Cobra 命令
  - [ ] 实现 `runWeb()` 执行函数
    - [ ] 检查环境变量（LLM credentials）
    - [ ] 初始化 LLM 和 Anonymizer
    - [ ] 创建 Web Server
    - [ ] 启动服务器并监听信号（SIGINT, SIGTERM）
    - [ ] 优雅关闭
- [ ] 5.2 修改 `cmd/inu/main.go`
  - [ ] 导入 web 命令
  - [ ] 注册到 root 命令

## 6. 错误处理和日志
- [ ] 6.1 统一错误响应格式
  - [ ] 定义 `ErrorResponse{Error, Message, Code}` 结构
  - [ ] 创建辅助函数 `respondError(c, statusCode, message)`
- [ ] 6.2 添加请求日志中间件
  - [ ] 记录请求方法、路径、状态码、耗时
  - [ ] 记录错误详情（调试用）

## 7. 测试实现
- [ ] 7.1 创建 `pkg/web/handlers/anonymize_test.go`
  - [ ] 测试正常匿名化请求
  - [ ] 测试空文本错误
  - [ ] 测试 LLM 错误处理
  - [ ] 测试无效 JSON 请求
- [ ] 7.2 创建 `pkg/web/handlers/restore_test.go`
  - [ ] 测试正常还原请求
  - [ ] 测试空实体列表
  - [ ] 测试错误处理
- [ ] 7.3 创建 `pkg/web/middleware/auth_test.go`
  - [ ] 测试正确的认证信息
  - [ ] 测试错误的认证信息
  - [ ] 测试缺少认证头
- [ ] 7.4 创建 `pkg/web/server_test.go`
  - [ ] 测试服务器启动和关闭
  - [ ] 集成测试：完整 API 请求流程

## 8. 文档和示例
- [ ] 8.1 更新 `README.md`
  - [ ] 添加 Web API 使用说明
  - [ ] 添加 API 端点文档
  - [ ] 添加 curl 示例
- [ ] 8.2 创建 API 示例脚本（可选）
  - [ ] `examples/api_anonymize.sh` - 匿名化示例
  - [ ] `examples/api_restore.sh` - 还原示例

## 9. 集成和验证
- [ ] 9.1 运行所有测试：`make test`
- [ ] 9.2 手动测试 Web 服务器
  - [ ] 启动服务器：`inu web --admin-token test123`
  - [ ] 测试健康检查：`curl http://localhost:8080/health`
  - [ ] 测试匿名化 API（带认证）
  - [ ] 测试还原 API
  - [ ] 测试无认证访问（应返回 401）
- [ ] 9.3 验证与现有 CLI 命令的兼容性
- [ ] 9.4 验证构建和发布流程
