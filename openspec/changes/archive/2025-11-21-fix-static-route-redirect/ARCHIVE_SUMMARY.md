# 归档总结：fix-static-route-redirect

**归档日期**: 2025-11-21
**状态**: ✅ 已完成并归档

## 变更概述

修复 Web UI 访问时的 301 重定向循环问题，该问题在 embed-static-assets 实施后出现，导致用户无法正常使用 Web UI。

## 实施成果

### 核心修改
- **pkg/web/server.go**: 重构静态文件服务实现（~40 行修改）
  - 移除使用 `ui.StaticFS` 和 `c.FileFromFS` 的实现
  - 使用 `embed.FS.ReadFile` 直接读取文件内容
  - 使用 `c.Data()` 返回内容，手动设置 Content-Type
  - 禁用 Gin 的自动重定向功能（`RedirectTrailingSlash` 和 `RedirectFixedPath`）
  - 实现路径处理逻辑：去除前导斜杠，根据扩展名设置正确的 Content-Type

### 测试覆盖
- **pkg/web/server_test.go**: 新增单元测试
  - `TestStaticRoutes`: 验证静态文件路由（4 个测试用例）
    - GET /static/app.js → 200 + JavaScript 内容
    - GET /static/styles.css → 200 + CSS 内容
    - GET /static/index.html → 200 + HTML 内容
    - GET /static/nonexistent.js → 404
  - `TestHomePageRoute`: 验证首页路由
    - GET / → 200 + index.html 内容
  - 创建 `mockAnonymizer` 辅助测试
- **测试结果**: 11/11 测试通过（新增 2 个测试函数，5 个测试用例）

### 问题根因分析

**原始问题**：
使用 Gin 的 `StaticFS` 和 `c.FileFromFS` 方法处理嵌入文件系统时：
```go
staticSubFS, _ := fs.Sub(staticFS, "static")
httpFS := http.FS(staticSubFS)
ui.StaticFS("/static", httpFS)
```

这导致：
1. `c.FileFromFS` 内部调用 `http.ServeFile`，在处理文件路径时会检查是否为目录
2. 如果路径被误判为目录，会返回 301 重定向到 `./`
3. 浏览器收到 `Location: ./` 后持续重定向，形成死循环

**最终方案**：
- 不使用 `fs.Sub` 创建子文件系统（避免根目录被视为目录）
- 不使用 `ui.StaticFS` 或 `c.FileFromFS`（避免自动重定向逻辑）
- 直接使用 `staticFS.ReadFile("static/" + filepath)` 读取文件
- 使用 `c.Data()` 手动设置 Content-Type 和状态码
- 禁用 Gin 的 `RedirectTrailingSlash` 和 `RedirectFixedPath`

### 实现代码

```go
// 禁用自动重定向
engine.RedirectTrailingSlash = false
engine.RedirectFixedPath = false

// 首页路由
ui.GET("/", func(c *gin.Context) {
    data, err := staticFS.ReadFile("static/index.html")
    if err != nil {
        c.String(http.StatusNotFound, "File not found")
        return
    }
    c.Data(http.StatusOK, "text/html; charset=utf-8", data)
})

// 静态文件路由
ui.GET("/static/*filepath", func(c *gin.Context) {
    filepath := strings.TrimPrefix(c.Param("filepath"), "/")
    data, err := staticFS.ReadFile("static/" + filepath)
    if err != nil {
        c.String(http.StatusNotFound, "File not found")
        return
    }

    // 根据扩展名设置 Content-Type
    contentType := "application/octet-stream"
    if strings.HasSuffix(filepath, ".html") {
        contentType = "text/html; charset=utf-8"
    } else if strings.HasSuffix(filepath, ".css") {
        contentType = "text/css; charset=utf-8"
    } else if strings.HasSuffix(filepath, ".js") {
        contentType = "application/javascript; charset=utf-8"
    }

    c.Data(http.StatusOK, contentType, data)
})
```

## 规格变更

### web-api Spec

**MODIFIED Requirements**:
- "静态资源嵌入式服务" - 添加 5 个新场景：
  - 访问首页不产生重定向循环
  - 访问静态资源路径不产生重定向
  - 访问 /static 前缀路径（无尾部斜杠）
  - 静态文件路由与 API 路由不冲突
  - 使用标准 HTTP 文件服务器特性

**ADDED Requirements**:
- "静态文件路由实现规范" - 新增 3 个场景：
  - 直接读取并返回静态资源
  - 路径前缀正确映射
  - 认证中间件正确应用

## 技术细节

### 为什么不使用 http.FileServer

尝试过使用标准库的 `http.FileServer` 配合 `http.StripPrefix`：
```go
fileServer := http.StripPrefix("/static", http.FileServer(httpFS))
ui.GET("/static/*filepath", gin.WrapH(fileServer))
```

但这仍然产生 301 重定向。原因：
- `http.FileServer` 在遇到目录时会自动添加尾部斜杠并重定向
- 与 Gin 的路由参数 `*filepath` 交互时，路径解析出现问题
- 即使禁用 Gin 的重定向，`http.FileServer` 内部仍会产生重定向

### 优势与约束

**优势**：
1. **完全控制**: 手动管理文件读取和响应，避免框架黑盒行为
2. **零重定向**: 直接返回内容，不依赖任何可能产生重定向的机制
3. **简单明确**: 代码逻辑清晰，易于理解和调试
4. **兼容性好**: 不依赖特定框架版本的行为

**约束**：
- 需要手动设置 Content-Type（当前实现支持 .html、.css、.js）
- 不支持 Range 请求和 ETag（可在未来添加）
- 每次请求都读取文件（Go 的 embed.FS 已经缓存在内存中，性能影响可忽略）

## 验证检查清单

- [x] 代码实现完成
- [x] 单元测试通过（11/11）
- [x] 修复了 301 重定向循环
- [x] 首页 `/` 正常返回 index.html
- [x] 静态资源 `/static/*` 正常返回文件
- [x] 404 错误正确处理
- [x] 文档更新完成
- [x] Spec deltas 应用到主 spec
- [x] 所有测试保持通过

## 后续建议

1. **性能优化**: 考虑添加 ETag 支持以启用浏览器缓存
2. **Content-Type**: 扩展支持更多文件类型（如 .png、.svg、.ico）
3. **压缩**: 考虑添加 gzip 压缩支持
4. **监控**: 在生产环境中监控静态资源加载性能
5. **CI/CD**: 添加集成测试验证浏览器访问无重定向

## 相关变更

- **前置变更**: `embed-static-assets` (2025-11-21)
  - 引入了 embed 嵌入静态资源功能
  - 使用 `ui.StaticFS` 实现，导致了本次修复的问题

---

**归档方式**: 手动归档（OpenSpec CLI 未安装）
**验证方式**: 运行单元测试 + 手动测试
**归档路径**: `openspec/changes/archive/2025-11-21-fix-static-route-redirect/`
