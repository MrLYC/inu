# 设计：将静态资源嵌入二进制文件

## 概述

使用 Go 1.16+ 的 `embed` 包将 Web UI 静态资源（HTML/CSS/JS）在编译时嵌入到二进制文件中，实现单文件分发和零依赖部署。

## 技术架构

### 当前架构

```
pkg/web/server.go
├── setupRoutes()
│   ├── ui.GET("/", func) -> c.File("pkg/web/static/index.html")
│   └── ui.Static("/static", "pkg/web/static")
└── 依赖外部文件系统中的 pkg/web/static/ 目录
```

**问题**：
- 运行时需要 `pkg/web/static/` 目录存在
- 二进制文件无法独立运行
- 路径依赖可能导致部署问题

### 目标架构

```
pkg/web/server.go
├── //go:embed static/*
├── var staticFS embed.FS
├── setupRoutes()
│   ├── ui.GET("/", serveIndexHTML) -> 从嵌入 FS 读取
│   └── ui.GET("/static/*", gin.WrapH(http.FileServer(http.FS(staticFS))))
└── 使用嵌入的文件系统，无需外部文件
```

**优势**：
- 编译时将静态文件打包进二进制
- 运行时从内存中的嵌入 FS 读取
- 单文件分发，零外部依赖

## 实现细节

### 1. 使用 embed 包

```go
package web

import (
	"embed"
	"io/fs"
	"net/http"
	// ... 其他导入
)

//go:embed static/*
var staticFS embed.FS
```

**说明**：
- `//go:embed` 是编译器指令，不是普通注释
- `static/*` 匹配 `pkg/web/static/` 下的所有文件
- `embed.FS` 实现了 `fs.FS` 接口，可用于文件读取

### 2. 创建 HTTP 文件系统

```go
// 获取 static 子目录的文件系统
staticSubFS, err := fs.Sub(staticFS, "static")
if err != nil {
	return nil, fmt.Errorf("failed to create sub filesystem: %w", err)
}

// 创建 HTTP 文件服务器
httpFS := http.FS(staticSubFS)
```

**说明**：
- `fs.Sub` 用于获取子目录，避免 URL 中包含 "static/static/"
- `http.FS` 将 `fs.FS` 适配为 `http.FileSystem`

### 3. 修改路由配置

#### 首页路由

```go
ui.GET("/", func(c *gin.Context) {
	// 从嵌入 FS 读取 index.html
	data, err := staticSubFS.Open("index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to load page")
		return
	}
	defer data.Close()
	
	c.DataFromReader(http.StatusOK, -1, "text/html; charset=utf-8", data, nil)
})
```

或使用更简洁的方式：

```go
ui.GET("/", func(c *gin.Context) {
	c.FileFromFS("index.html", httpFS)
})
```

#### 静态文件路由

```go
// 方案 1：使用 Gin 的静态文件处理
ui.StaticFS("/static", httpFS)

// 方案 2：使用标准库 http.FileServer
ui.GET("/static/*filepath", func(c *gin.Context) {
	http.FileServer(httpFS).ServeHTTP(c.Writer, c.Request)
})
```

**推荐方案 1**：更符合 Gin 的使用习惯。

### 4. 完整修改示例

```go
// setupRoutes configures all HTTP routes and middleware
func (s *Server) setupRoutes() {
	// 创建嵌入的静态文件系统
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("Failed to create static filesystem: %v", err)
	}
	httpFS := http.FS(staticSubFS)

	// Determine if auth is enabled
	authEnabled := s.config.IsAuthEnabled()

	// UI routes (auth required if enabled)
	ui := s.engine.Group("/")
	if authEnabled {
		ui.Use(middleware.BasicAuth(s.config.AdminUser, s.config.AdminToken))
	}
	{
		// 首页
		ui.GET("/", func(c *gin.Context) {
			c.FileFromFS("index.html", httpFS)
		})
		// 静态资源
		ui.StaticFS("/static", httpFS)
	}

	// ... 其他路由保持不变
}
```

## 构建和测试

### 构建流程

1. **无需修改 Makefile**：Go 编译器自动处理 `//go:embed`
2. **编译**：`make build` 或 `go build`
3. **验证**：检查二进制文件大小增加（约 10-20KB）

### 测试验证

#### 1. 单元测试

```go
func TestStaticFilesEmbedded(t *testing.T) {
	// 验证 staticFS 可访问
	entries, err := staticFS.ReadDir("static")
	assert.NoError(t, err)
	assert.NotEmpty(t, entries)
	
	// 验证关键文件存在
	expectedFiles := []string{"index.html", "app.js", "styles.css"}
	for _, filename := range expectedFiles {
		_, err := staticFS.ReadFile("static/" + filename)
		assert.NoError(t, err, "File should be embedded: %s", filename)
	}
}
```

#### 2. 集成测试

```bash
# 编译二进制
make build

# 删除静态文件目录（验证嵌入生效）
mv pkg/web/static pkg/web/static.bak

# 启动服务器
./bin/inu web --admin-token test

# 测试 UI 访问
curl -u admin:test http://localhost:8080/
curl -u admin:test http://localhost:8080/static/app.js

# 恢复静态文件
mv pkg/web/static.bak pkg/web/static
```

#### 3. 手动测试

1. 编译二进制文件
2. 将二进制文件复制到空目录（无静态文件）
3. 运行 `./inu web --admin-token test`
4. 浏览器访问 `http://localhost:8080/`
5. 验证 UI 完整显示和功能正常

## 边界情况处理

### 1. 文件不存在

嵌入的文件系统在编译时确定，不会出现文件不存在的情况。如果编译时文件缺失，编译会失败。

### 2. 开发模式

在开发时修改静态文件后：
1. 重新编译：`make build`
2. 重启服务器
3. 静态文件自动更新

### 3. 大文件处理

当前静态文件总大小约 10-20KB，对二进制文件影响可忽略。如果未来需要嵌入大文件：
- 考虑压缩静态资源
- 或使用外部 CDN 存储大型资源

## 兼容性

### Go 版本要求

- **最低版本**：Go 1.16（引入 `embed` 包）
- **当前项目**：Go 1.24.4 ✅

### 平台兼容性

- ✅ Linux (amd64, arm64)
- ✅ macOS (amd64, arm64)  
- ✅ Windows (amd64)

所有平台的编译和运行行为一致。

## 性能考虑

### 内存使用

- 嵌入的文件在编译时存储在二进制的只读数据段
- 运行时通过内存映射访问，无需全部加载到堆内存
- 性能影响可忽略

### 响应速度

- 从内存读取比磁盘 I/O 更快
- HTTP 缓存机制仍然有效
- 用户体验无差异或更好

## 回滚方案

如果嵌入方案出现问题，可以快速回滚：

```go
// 回滚到外部文件系统（保留原有代码）
ui.GET("/", func(c *gin.Context) {
	c.File("pkg/web/static/index.html")
})
ui.Static("/static", "pkg/web/static")
```

## 参考资料

- [Go embed 包文档](https://pkg.go.dev/embed)
- [Go 1.16 Release Notes - Embed](https://go.dev/doc/go1.16#library-embed)
- [Gin 静态文件服务](https://gin-gonic.com/docs/examples/serving-static-files/)
