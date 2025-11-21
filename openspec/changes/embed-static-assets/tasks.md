# 任务：将静态资源嵌入二进制文件

## 1. 准备和验证
- [ ] 1.1 确认 Go 版本 >= 1.16（embed 包支持）
- [ ] 1.2 阅读 embed 包文档和最佳实践
- [ ] 1.3 检查当前静态文件结构和大小

## 2. 核心实现
- [ ] 2.1 修改 `pkg/web/server.go`
  - [ ] 导入 `embed` 和 `io/fs` 包
  - [ ] 添加 `//go:embed static/*` 指令
  - [ ] 声明 `var staticFS embed.FS`
  
- [ ] 2.2 修改 `setupRoutes()` 函数
  - [ ] 创建静态文件子系统：`fs.Sub(staticFS, "static")`
  - [ ] 创建 HTTP 文件系统：`http.FS(staticSubFS)`
  - [ ] 修改首页路由使用 `c.FileFromFS("index.html", httpFS)`
  - [ ] 修改静态资源路由使用 `ui.StaticFS("/static", httpFS)`

## 3. 测试验证
- [ ] 3.1 编写单元测试 `pkg/web/server_test.go`
  - [ ] `TestStaticFilesEmbedded` - 验证文件嵌入成功
  - [ ] 测试读取 index.html
  - [ ] 测试读取 app.js
  - [ ] 测试读取 styles.css
  
- [ ] 3.2 编译测试
  - [ ] 运行 `make build` 验证编译成功
  - [ ] 检查二进制文件大小增加（约 10-20KB）
  - [ ] 运行 `make test` 确保所有测试通过
  
- [ ] 3.3 集成测试
  - [ ] 编译二进制文件
  - [ ] 备份并删除 `pkg/web/static/` 目录
  - [ ] 启动 Web 服务器：`./bin/inu web --admin-token test`
  - [ ] 使用 curl 测试 `GET /`
  - [ ] 使用 curl 测试 `GET /static/app.js`
  - [ ] 使用 curl 测试 `GET /static/styles.css`
  - [ ] 恢复 `pkg/web/static/` 目录

- [ ] 3.4 手动功能测试
  - [ ] 将二进制文件复制到空目录
  - [ ] 启动服务器
  - [ ] 浏览器访问 `http://localhost:8080/`
  - [ ] 验证 UI 完整显示（HTML、CSS、JS 都加载）
  - [ ] 测试脱敏功能正常
  - [ ] 测试还原功能正常
  - [ ] 测试实体映射显示正常

## 4. 多平台构建验证
- [ ] 4.1 构建所有平台二进制
  - [ ] 运行 `make build-all`
  - [ ] 验证所有平台编译成功
  
- [ ] 4.2 平台独立性测试
  - [ ] 在 Linux 上测试二进制（如有环境）
  - [ ] 在 macOS 上测试二进制
  - [ ] 在 Windows 上测试二进制（如有环境）

## 5. 文档更新
- [ ] 5.1 更新 `README.md`
  - [ ] 说明二进制文件已包含 Web UI
  - [ ] 更新部署说明：单文件分发
  - [ ] 移除关于静态文件的手动配置说明（如有）
  
- [ ] 5.2 更新 `openspec/project.md`
  - [ ] 在 "Architecture Patterns" 中添加静态资源嵌入说明
  - [ ] 更新部署相关描述

## 6. CI/CD 验证
- [ ] 6.1 推送代码触发 CI
  - [ ] 验证 CI 构建成功
  - [ ] 验证测试全部通过
  - [ ] 检查构建产物大小

## 7. 清理和优化
- [ ] 7.1 代码审查
  - [ ] 检查错误处理是否完善
  - [ ] 确保代码注释清晰
  - [ ] 验证日志输出合理
  
- [ ] 7.2 性能验证
  - [ ] 对比嵌入前后的响应时间
  - [ ] 验证内存使用无明显增加
  - [ ] 确认静态文件缓存正常工作

## 8. 最终验证
- [ ] 8.1 端到端测试
  - [ ] 从 GitHub Release 下载二进制（或本地构建）
  - [ ] 在干净环境中运行
  - [ ] 完整测试所有 Web UI 功能
  - [ ] 验证用户体验符合预期
  
- [ ] 8.2 回归测试
  - [ ] 运行完整测试套件：`make test`
  - [ ] 测试 CLI 命令不受影响
  - [ ] 测试 API 端点功能正常

## 9. 准备归档
- [ ] 9.1 更新变更文档
  - [ ] 标记所有任务为完成
  - [ ] 更新 proposal.md 状态
  - [ ] 准备变更总结

- [ ] 9.2 OpenSpec 验证
  - [ ] 运行 `openspec validate embed-static-assets --strict`
  - [ ] 修复所有验证错误
  - [ ] 准备归档
