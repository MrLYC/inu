# 设计：修复静态资源路由 301 重定向循环

## 问题诊断

### 症状
用户访问 `http://localhost:8080/` 时，浏览器显示不断的 301 重定向，页面无法加载。

### 根因分析

**当前实现**（`pkg/web/server.go:67-85`）：
```go
staticSubFS, err := fs.Sub(staticFS, "static")
if err != nil {
    log.Fatalf("Failed to create static filesystem: %v", err)
}
httpFS := http.FS(staticSubFS)

ui := s.engine.Group("/")
if authEnabled {
    ui.Use(middleware.BasicAuth(s.config.AdminUser, s.config.AdminToken))
}
{
    ui.GET("/", func(c *gin.Context) {
        c.FileFromFS("index.html", httpFS)
    })
    ui.StaticFS("/static", httpFS)  // ❌ 问题所在
}
```

**为什么 `StaticFS` 会导致 301 重定向？**

Gin 的 `StaticFS(relativePath string, fs http.FileSystem)` 内部实现：

1. 对于路径 `/static`（不带尾部斜杠），Gin 会检查文件系统中是否存在 "static" 文件或目录
2. 如果是目录，Gin 会发送 301 重定向到 `/static/`（带尾部斜杠）
3. 对于 `/static/`（带尾部斜杠），Gin 尝试返回目录索引

但在我们的情况下：
- `httpFS` 的根目录是原 `static/` 目录（通过 `fs.Sub` 提取）
- 当访问 `/static` 时，Gin 在 `httpFS` 根目录查找名为 "static" 的条目
- 如果找不到，或者找到但处理不当，可能导致重定向循环

**实际问题**：
`fs.Sub(staticFS, "static")` 返回的文件系统根目录包含 `index.html`、`app.js`、`styles.css`，但 **不包含** 名为 "static" 的目录。因此：
- 访问 `/static` → Gin 查找 "static" → 未找到 → 可能触发意外行为
- 访问 `/static/app.js` → Gin 查找 "app.js" → 找到 → 正常返回

但如果 Gin 的路由匹配逻辑在处理 `/static` 前缀时与 `/` 路由冲突，可能导致重定向循环。

### 验证假设

让我们检查 Gin 的 `StaticFS` 源码行为（基于 Gin v1.9+）：

```go
// Gin 内部实现（简化）
func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
    if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
        panic("URL parameters cannot be used when serving a static folder")
    }
    handler := group.createStaticHandler(relativePath, fs)
    urlPattern := path.Join(relativePath, "/*filepath")
    
    // Register GET and HEAD handlers
    group.GET(urlPattern, handler)
    group.HEAD(urlPattern, handler)
    return group.returnObj()
}
```

关键点：
- `StaticFS("/static", fs)` 实际注册的路由是 `/static/*filepath`
- 这意味着 `/static`（不带尾部斜杠）**不会**被这个路由匹配
- Gin 可能会有默认的重定向行为处理 `/static` → `/static/`

## 解决方案设计

### 方案 1：使用 `http.FileServer` 和自定义处理器（推荐）

**实现**：
```go
func (s *Server) setupRoutes() {
    staticSubFS, err := fs.Sub(staticFS, "static")
    if err != nil {
        log.Fatalf("Failed to create static filesystem: %v", err)
    }
    httpFS := http.FS(staticSubFS)
    
    authEnabled := s.config.IsAuthEnabled()
    
    // UI routes
    ui := s.engine.Group("/")
    if authEnabled {
        ui.Use(middleware.BasicAuth(s.config.AdminUser, s.config.AdminToken))
    }
    {
        // 首页
        ui.GET("/", func(c *gin.Context) {
            c.FileFromFS("index.html", httpFS)
        })
        
        // 静态资源：使用 http.FileServer
        fileServer := http.FileServer(httpFS)
        ui.GET("/static/*filepath", func(c *gin.Context) {
            // 去除 /static 前缀，保留文件路径
            c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/static")
            if c.Request.URL.Path == "" {
                c.Request.URL.Path = "/"
            }
            fileServer.ServeHTTP(c.Writer, c.Request)
        })
    }
    
    // ... 其他路由
}
```

**优点**：
- ✅ 使用标准 `http.FileServer`，行为可预测
- ✅ 避免 Gin `StaticFS` 的路径匹配问题
- ✅ 保留文件服务器的完整功能（ETag、Range 请求等）
- ✅ 与认证中间件兼容

**缺点**：
- 需要手动处理路径前缀

### 方案 2：直接注册静态文件路由

**实现**：
```go
ui.GET("/static/app.js", func(c *gin.Context) {
    c.FileFromFS("app.js", httpFS)
})
ui.GET("/static/styles.css", func(c *gin.Context) {
    c.FileFromFS("styles.css", httpFS)
})
```

**优点**：
- ✅ 最直接，完全控制路由
- ✅ 避免任何路径匹配问题

**缺点**：
- ❌ 需要为每个静态文件注册路由（不可扩展）
- ❌ 添加新静态文件时需要修改代码

### 方案 3：使用 `NoRoute` 处理 fallback

**实现**：
```go
ui.GET("/", func(c *gin.Context) {
    c.FileFromFS("index.html", httpFS)
})

// NoRoute 处理所有未匹配的路由
s.engine.NoRoute(func(c *gin.Context) {
    path := c.Request.URL.Path
    if strings.HasPrefix(path, "/static/") {
        filename := strings.TrimPrefix(path, "/static/")
        c.FileFromFS(filename, httpFS)
    } else {
        c.JSON(404, gin.H{"error": "not found"})
    }
})
```

**优点**：
- ✅ 灵活，处理所有静态资源

**缺点**：
- ❌ 影响其他路由的 404 处理
- ❌ 认证中间件可能不生效（NoRoute 不在 ui 组内）

## 推荐实现

**选择方案 1**，原因：
1. 使用标准库 `http.FileServer`，可靠且功能完整
2. 路径处理清晰，易于理解和维护
3. 与现有认证中间件无缝集成
4. 性能优秀（支持 ETag、Range 请求等）

## 实现步骤

1. **修改 `setupRoutes()` 方法**
   - 替换 `ui.StaticFS("/static", httpFS)` 
   - 使用 `http.FileServer` + 自定义路由处理器

2. **添加 import**
   - 添加 `strings` 包（如果尚未导入）

3. **更新测试**
   - 在 `server_test.go` 中添加路由测试
   - 测试 `/static/app.js`、`/static/styles.css` 的访问
   - 测试 `/static` 和 `/static/` 的行为

4. **集成测试**
   - 启动服务器，使用 curl 验证所有路径
   - 使用浏览器访问，确认无重定向循环

## 回退计划

如果方案 1 仍有问题，可以临时回退到方案 2（逐个注册文件），虽然不优雅但能快速解决问题。

长期来看，如果静态资源数量增加，可以考虑：
- 使用构建工具生成路由注册代码
- 或使用嵌入式文件系统的目录遍历功能动态注册路由
