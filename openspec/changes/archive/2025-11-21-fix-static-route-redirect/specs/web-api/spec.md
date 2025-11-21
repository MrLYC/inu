# web-api 规格变更

## MODIFIED Requirements

### Requirement: 静态资源嵌入式服务（修改实现）
系统 SHALL 使用嵌入的文件系统提供 Web UI 静态资源，无需依赖外部文件。实现应避免路由冲突和重定向循环。

#### Scenario: 访问首页不产生重定向循环
- **GIVEN** Web 服务器正在运行
- **AND** 静态资源已通过 embed 嵌入
- **WHEN** 客户端发送 `GET /` 请求
- **THEN** 服务器应该直接返回 200 OK 和 `index.html` 内容
- **AND** 不应该发送任何 301 或 302 重定向
- **AND** 响应头应该包含 `Content-Type: text/html`

#### Scenario: 访问静态资源路径不产生重定向
- **GIVEN** Web 服务器正在运行
- **WHEN** 客户端请求 `GET /static/app.js`
- **THEN** 服务器应该直接返回 200 OK 和 JavaScript 内容
- **AND** 不应该发送 301 重定向到 `/static/app.js/` 或其他路径
- **AND** Content-Type 应该是 `application/javascript` 或 `text/javascript`

#### Scenario: 访问 /static 前缀路径（无尾部斜杠）
- **WHEN** 客户端请求 `GET /static` （不带尾部斜杠）
- **THEN** 服务器可以返回 404 Not Found 或重定向到 `/static/`
- **AND** 如果重定向，应该是 301 Moved Permanently 到 `/static/`
- **AND** 重定向不应该形成循环（即 `/static/` 不应再次重定向）

#### Scenario: 静态文件路由与 API 路由不冲突
- **GIVEN** 系统同时提供静态资源和 API 端点
- **WHEN** 客户端请求 `/api/v1/anonymize`
- **THEN** 请求应该路由到 API 处理器，而不是静态文件处理器
- **AND** 返回 JSON 响应，而非 HTML 或 404

#### Scenario: 使用标准 HTTP 文件服务器特性
- **WHEN** 客户端请求静态资源
- **THEN** 响应应该包含标准的文件服务器特性：
  - `Content-Type` 根据文件扩展名正确设置
  - `ETag` 或 `Last-Modified` 用于缓存验证
  - 支持 `If-None-Match` / `If-Modified-Since` 条件请求
  - 支持 `Range` 请求（部分内容）

## ADDED Requirements

### Requirement: 静态文件路由实现规范
系统 SHALL 使用 `http.FileServer` 或等效机制来服务嵌入的静态资源，确保路由行为符合标准 HTTP 语义。

#### Scenario: 使用 http.FileServer 提供静态资源
- **GIVEN** 静态资源嵌入在 `embed.FS` 中
- **WHEN** 系统初始化静态文件路由
- **THEN** 应该使用 `http.FileServer` 或兼容的文件服务器
- **AND** 文件服务器应该正确处理路径前缀（如 `/static`）
- **AND** 文件服务器应该支持标准 HTTP 特性（ETag、Range 等）

#### Scenario: 路径前缀正确映射
- **GIVEN** 静态资源位于嵌入的 `static/` 目录
- **WHEN** 客户端请求 `/static/app.js`
- **THEN** 服务器应该从嵌入的文件系统中读取 `app.js`（而非 `static/app.js`）
- **AND** 路径映射应该透明且一致

#### Scenario: 认证中间件正确应用
- **GIVEN** Web UI 启用了 Basic Auth
- **WHEN** 客户端请求静态资源 `/static/app.js` 不带认证头
- **THEN** 服务器应该返回 401 Unauthorized
- **WHEN** 客户端提供正确的认证信息
- **THEN** 服务器应该返回 200 OK 和文件内容

## REMOVED Requirements

无移除的需求。

---

**实现注意事项**：
- 避免使用可能导致路径混淆的框架方法（如 Gin 的 `StaticFS`，如果其行为不符合预期）
- 优先使用标准库 `http.FileServer` 或明确的路径处理逻辑
- 确保路由注册顺序正确，避免通配符路由覆盖具体路由
- 静态文件路由应该在 API 路由之前或之后注册，以避免冲突

**测试要求**：
- 必须包含单元测试验证静态文件路由返回正确的状态码和内容
- 必须包含集成测试在浏览器中验证无重定向循环
- 必须测试认证中间件对静态资源的影响
