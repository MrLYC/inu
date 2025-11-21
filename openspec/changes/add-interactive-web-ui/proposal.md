# 提案: 添加交互式 Web 界面

## 为什么

目前 `inu web` 仅提供 RESTful API，没有用户界面。用户必须通过 curl 或 Postman 等 API 客户端与匿名化服务交互。`inu interactive` 命令提供了更好的工作流程，具有跨多个还原周期的实体内存，但仅限于 CLI 环境。

**问题:**
- Web 模式缺少可视化界面，限制了可访问性
- 不熟悉 API 的用户无法轻松使用 Web 服务
- 交互式工作流（一次匿名化，多次还原）在 Web 模式下不可用
- 还原过程中缺少实体映射的可视化反馈

## 做什么

为 `inu web` 添加单页 Web 界面，在浏览器界面中复制交互式命令工作流。

**核心功能:**

1. **匿名化视图**
   - 实体类型选择器（从 `--entity-types` 参数填充 + 自定义输入）
   - 用于输入待匿名化文本的文本区域
   - 显示匿名化结果的输出面板
   - 带加载状态的"匿名化"按钮
   - "切换到还原模式"按钮（匿名化后出现）

2. **还原视图**
   - 实体映射显示（显示占位符 → 原始值对）
   - 只读的匿名化文本面板（左侧）
   - 可编辑的输入文本区域（右侧）用于外部处理结果
   - "还原"按钮用于反匿名化当前文本
   - "返回匿名化"按钮返回第一个视图

3. **状态管理**
   - 使用 sessionStorage 进行实体的客户端存储
   - 实体在同一浏览器会话的视图切换中持久化
   - 无需服务器端会话（无状态 API 设计）

**实现方式:**
- 原生 HTML/CSS/JavaScript（无框架依赖）
- 通过 Gin 的静态文件处理器提供静态文件
- 复用现有的 `/api/v1/anonymize` 和 `/api/v1/restore` 端点
- 新路由: `GET /` 提供 UI 主页（无需认证）
- 新路由: `GET /api/v1/config` 提供实体类型配置（可选）

## 影响

**规范:**
- **web-api**: 为 UI 路由、静态文件服务和前端功能添加了 ADDED 要求

**代码:**
- `pkg/web/server.go`: 添加 `GET /` 路由和静态文件处理器
- `pkg/web/static/`: 用于 HTML/CSS/JS 文件的新目录
- `cmd/inu/commands/web.go`: 无需更改（实体类型已可配置）

**依赖项:**
- 无（原生 JS，无新的 Go 依赖）

**破坏性变更:**
- 无（API 端点保持不变）

**文档:**
- 使用 Web UI 更新 README
- 添加 Web 界面截图

## 考虑的替代方案

1. **React/Vue SPA**: 更复杂，需要构建工具，对于简单 UI 来说过于复杂
2. **服务器渲染模板**: 需要服务器端状态管理，破坏无状态 API 设计
3. **独立的 UI 服务器**: 额外的部署复杂性，对于单团队项目来说不必要

## 待解决问题

1. UI 路由是否需要认证？（建议：GET / 无需认证，保持 API 认证）
2. 实体类型是否应该在运行时可配置？（建议：v1 仅支持 CLI 参数）
3. 是否应该支持多个并发会话？（建议：不需要，客户端状态足够）

## Impact

**Specs:**
- **web-api**: ADDED requirements for UI routes, static file serving, and frontend functionality

**Code:**
- `pkg/web/server.go`: Add `GET /` route and static file handler
- `pkg/web/static/`: New directory for HTML/CSS/JS files
- `cmd/inu/commands/web.go`: No changes needed (entity types already configurable)

**Dependencies:**
- None (vanilla JS, no new Go dependencies)

**Breaking Changes:**
- None (API endpoints remain unchanged)

**Documentation:**
- Update README with web UI usage
- Add screenshots of web interface

## Alternatives Considered

1. **React/Vue SPA**: More complex, requires build tooling, overkill for simple UI
2. **Server-rendered templates**: Requires state management on server, breaks stateless API design
3. **Separate UI server**: Additional deployment complexity, unnecessary for single-team project

## Open Questions

1. Should UI routes require authentication? (Recommend: no auth for GET /, keep API auth)
2. Should entity types be configurable at runtime? (Recommend: CLI flag only for v1)
3. Should we support multiple concurrent sessions? (Recommend: not needed, client-side state sufficient)
