# 提案：修复静态资源路由 301 重定向循环

## Why

### 问题
在 embed-static-assets 变更实施后，访问 Web UI 时出现 301 重定向循环问题。用户访问主页时浏览器会不断重定向，导致无法正常使用 Web UI。

**根本原因**：
在 `pkg/web/server.go` 的 `setupRoutes()` 中，使用了不正确的静态文件服务配置：

```go
staticSubFS, _ := fs.Sub(staticFS, "static")
httpFS := http.FS(staticSubFS)
ui.StaticFS("/static", httpFS)
```

这导致 Gin 的 `StaticFS` 在处理 `/static/*` 路径时行为异常。Gin 的 `StaticFS` 方法期望文件系统的根目录直接映射到 URL 路径，但这里的配置可能导致路径解析问题。

### 影响范围
- **严重性**: 🔴 Critical - Web UI 完全无法使用
- **受影响用户**: 所有使用 Web UI 的用户
- **受影响功能**: Web UI 首页、静态资源加载（CSS、JS）

### 目标
修复静态资源路由配置，确保：
1. 访问 `/` 能正常返回 `index.html`
2. 访问 `/static/app.js`、`/static/styles.css` 等能正常返回对应文件
3. 不出现 301 重定向循环
4. 保持嵌入式文件系统的使用（不回退到外部文件）

## What Changes

### 受影响的组件
- `pkg/web/server.go`：修复 `setupRoutes()` 方法中的静态文件路由配置

### 技术方案

**问题分析**：
Gin 的 `StaticFS(relativePath string, fs http.FileSystem)` 方法将 URL 路径 `relativePath` 映射到文件系统 `fs` 的根目录。当前代码：

```go
staticSubFS, _ := fs.Sub(staticFS, "static")  // 从 embed.FS 中提取 static/ 子目录
httpFS := http.FS(staticSubFS)                // 转换为 http.FileSystem
ui.StaticFS("/static", httpFS)                // 映射 /static -> httpFS 的根（即原 static/ 目录）
```

这看起来是正确的，但可能存在以下问题之一：
1. Gin 的 `StaticFS` 处理尾部斜杠时的重定向逻辑
2. `http.FS` 包装后的路径解析问题
3. 与认证中间件的交互问题

**解决方案选项**：

**选项 1（推荐）**：使用 `Static` 方法的文件服务器适配器
```go
// 不使用 StaticFS，而是手动创建文件服务器
fileServer := http.FileServer(httpFS)
ui.GET("/static/*filepath", func(c *gin.Context) {
    c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/static")
    fileServer.ServeHTTP(c.Writer, c.Request)
})
```

**选项 2**：使用 Gin 的 `StaticFileFS` 逐个注册文件
```go
ui.GET("/static/app.js", func(c *gin.Context) {
    c.FileFromFS("app.js", httpFS)
})
ui.GET("/static/styles.css", func(c *gin.Context) {
    c.FileFromFS("styles.css", httpFS)
})
```

**选项 3**：调整 `StaticFS` 的路径映射
```go
// 不使用 fs.Sub，直接从 staticFS 的 "static" 前缀访问
httpFS := http.FS(staticFS)
ui.GET("/static/*filepath", func(c *gin.Context) {
    // 手动处理路径：/static/app.js -> static/app.js
    path := c.Param("filepath")
    c.FileFromFS("static"+path, httpFS)
})
```

**推荐方案**：选项 1，因为它：
- 使用标准的 `http.FileServer`，行为可预测
- 保留文件系统的完整功能（目录列表、ETag 等）
- 与认证中间件兼容性好
- 代码简洁明了

### 验证方法
1. **单元测试**：更新 `pkg/web/server_test.go`，添加路由测试
2. **集成测试**：启动服务器，使用 curl 测试各路径
3. **浏览器测试**：在浏览器中访问 Web UI，确认无重定向循环

## Risks

### 技术风险
- **低风险**：修改仅涉及路由配置，不影响核心业务逻辑
- **向后兼容**：URL 路径保持不变（`/` 和 `/static/*`）

### 依赖关系
- **无阻塞依赖**：可以立即实施
- **后续验证**：需要确认所有静态资源（HTML、CSS、JS）都能正常加载

## Affected Specs
- `web-api`：修复 "静态资源嵌入式服务" 需求的实现

## References
- Gin StaticFS 文档: https://pkg.go.dev/github.com/gin-gonic/gin#RouterGroup.StaticFS
- Go http.FileServer 文档: https://pkg.go.dev/net/http#FileServer
- 相关 issue: embed-static-assets 实施后的路由问题
