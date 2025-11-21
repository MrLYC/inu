# 任务：修复静态资源路由 301 重定向循环

## 1. 问题调查和验证
- [x] 1.1 重现 301 重定向问题
- [x] 1.2 分析 `pkg/web/server.go` 中的路由配置
- [x] 1.3 研究 Gin `StaticFS` 的行为和限制
- [x] 1.4 确定根本原因（`StaticFS` 路径匹配问题）

## 2. 实现修复
- [ ] 2.1 修改 `pkg/web/server.go` 的 `setupRoutes()` 方法
  - [ ] 移除 `ui.StaticFS("/static", httpFS)` 行
  - [ ] 添加 `import "strings"` (如果缺失)
  - [ ] 实现基于 `http.FileServer` 的静态文件处理器：
    ```go
    fileServer := http.FileServer(httpFS)
    ui.GET("/static/*filepath", func(c *gin.Context) {
        c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/static")
        if c.Request.URL.Path == "" {
            c.Request.URL.Path = "/"
        }
        fileServer.ServeHTTP(c.Writer, c.Request)
    })
    ```

## 3. 测试验证
- [ ] 3.1 更新 `pkg/web/server_test.go`
  - [ ] 添加测试：`TestStaticRoutes`
    - 测试 `/static/app.js` 返回 200 和正确的 Content-Type
    - 测试 `/static/styles.css` 返回 200 和正确的 Content-Type
    - 测试 `/static/nonexistent.js` 返回 404
  - [ ] 添加测试：`TestHomePageRoute`
    - 测试 `/` 返回 200 和 `index.html` 内容
  
- [ ] 3.2 运行单元测试
  - [ ] 执行 `go test ./pkg/web -v`
  - [ ] 确保所有测试通过

- [ ] 3.3 手动集成测试
  - [ ] 编译项目：`make build`
  - [ ] 启动服务器：`./bin/inu web --admin-token test`
  - [ ] 测试首页：`curl -u admin:test http://localhost:8080/`
    - 预期：返回 HTML 内容，无 301 重定向
  - [ ] 测试 JS 文件：`curl -u admin:test http://localhost:8080/static/app.js`
    - 预期：返回 JavaScript 内容
  - [ ] 测试 CSS 文件：`curl -u admin:test http://localhost:8080/static/styles.css`
    - 预期：返回 CSS 内容
  - [ ] 测试不存在的文件：`curl -u admin:test http://localhost:8080/static/nonexistent.js`
    - 预期：返回 404

- [ ] 3.4 浏览器测试
  - [ ] 在浏览器中访问 `http://localhost:8080/`
  - [ ] 输入认证信息（admin / test）
  - [ ] 确认页面正常加载，无重定向循环
  - [ ] 打开浏览器开发者工具，检查网络请求
    - 确认 `/static/app.js` 返回 200
    - 确认 `/static/styles.css` 返回 200
    - 确认无 301 重定向
  - [ ] 测试 Web UI 功能（脱敏、还原）正常工作

## 4. 代码审查
- [ ] 4.1 检查代码风格
  - [ ] 使用 `gofmt` 格式化代码
  - [ ] 确保导入语句有序
  - [ ] 添加必要的注释

- [ ] 4.2 错误处理
  - [ ] 确认 `fs.Sub` 的错误处理保留
  - [ ] 确认 `FileServer` 的错误由 HTTP 层处理

## 5. 文档更新
- [ ] 5.1 更新相关注释
  - [ ] 在 `setupRoutes()` 中添加注释说明静态文件路由的实现方式

- [ ] 5.2 检查是否需要更新用户文档
  - [ ] README.md 无需更新（用户行为无变化）
  - [ ] 内部文档无需更新（实现细节变更）

## 6. 回归测试
- [ ] 6.1 运行完整测试套件
  - [ ] 执行 `make test`
  - [ ] 确保所有测试通过（预期 103+ 测试）

- [ ] 6.2 验证其他路由未受影响
  - [ ] 测试 `/health` 端点
  - [ ] 测试 `/api/v1/anonymize` 端点
  - [ ] 测试 `/api/v1/restore` 端点
  - [ ] 测试 `/api/v1/config` 端点

## 7. 性能验证
- [ ] 7.1 对比修复前后的响应时间
  - [ ] 使用 `ab` 或 `hey` 工具进行简单压测
  - [ ] 确认性能无明显退化

- [ ] 7.2 验证 ETag 和缓存头
  - [ ] 确认静态资源响应包含 ETag
  - [ ] 测试条件请求（If-None-Match）返回 304

## 8. 准备提交
- [ ] 8.1 整理提交
  - [ ] 创建清晰的 commit message
  - [ ] 引用此提案 ID：fix-static-route-redirect

- [ ] 8.2 更新变更文档
  - [ ] 标记所有任务为完成
  - [ ] 更新 proposal.md 状态

## 9. 验证和归档准备
- [ ] 9.1 最终验证
  - [ ] 在干净环境中编译和测试
  - [ ] 确认问题已完全解决

- [ ] 9.2 准备归档
  - [ ] 运行 `openspec validate fix-static-route-redirect --strict`
  - [ ] 修复所有验证错误
  - [ ] 准备归档到 archive/
