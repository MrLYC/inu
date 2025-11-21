# 设计: 交互式 Web 界面

## 架构

### 前端技术
**决策: 原生 HTML/CSS/JavaScript**

**理由:**
- 无需构建工具（更简单的开发和部署）
- 无框架依赖（更小的包体积，更快的加载）
- 足以应对 UI 复杂度（两个视图，简单的状态管理）
- 小团队更易维护

**权衡:**
- 手动 DOM 操作（相比声明式框架）
- 无组件复用性（对于小型 UI 可接受）
- 相比框架 DSL 更冗长

### 状态管理
**决策: 客户端 sessionStorage**

**理由:**
- 实体仅在单个浏览器会话中需要
- 无服务器端会话复杂性
- 维持无状态 API 设计
- 会话结束时自动清理

- 会话结束时自动清理

**状态模式:**
```javascript
{
  "entities": {
    "[PERSON_1]": "张三",
    "[ORG_1]": "阿里巴巴",
    // ... 更多映射
  },
  "anonymizedText": "...",
  "entityTypes": ["PERSON", "ORG", "EMAIL", "PHONE", "ADDRESS"]
}
```

**权衡:**
- 关闭标签页时数据丢失（对于匿名化工作流可接受）
- 无法跨标签页共享（单用户工具不需要）
- 浏览器存储限制（约 5MB，足够用于实体映射）

### UI 布局

#### 视图 1: 匿名化模式
```
┌────────────────────────────────────────────────┐
│ 实体类型: [PERSON▼] [+ 添加自定义类型]         │
├────────────────────┬───────────────────────────┤
│                    │                           │
│ 输入文本           │ 匿名化输出                 │
│ (可编辑)           │ (只读，显示结果)           │
│                    │                           │
│                    │                           │
├────────────────────┴───────────────────────────┤
│ [匿名化] [切换到还原模式]                      │
└────────────────────────────────────────────────┘
```

#### 视图 2: 还原模式
```
┌────────────────────────────────────────────────┐
│ 实体映射:                                      │
│ [PERSON_1] → 张三  [ORG_1] → 阿里巴巴  ...      │
├────────────────────┬───────────────────────────┤
│ 匿名化文本         │ 还原输入                   │
│ (只读)             │ (可编辑，粘贴到此处)       │
│                    │                           │
│                    │                           │
├────────────────────┴───────────────────────────┤
│ [还原] [返回匿名化]                            │
└────────────────────────────────────────────────┘
```

### API 集成

**现有端点（无变更）:**
- `POST /api/v1/anonymize`: 请求体 `{"text": "...", "entity_types": ["PERSON", ...]}`
- `POST /api/v1/restore`: 请求体 `{"text": "...", "entities": {"[PERSON_1]": "张三", ...}}`

**新端点:**
- `GET /`: 提供 `index.html`（无需认证）
- `GET /static/*`: 提供 CSS/JS 文件（无需认证）
- `GET /api/v1/config`: 返回从 CLI 参数获取的 `{"entity_types": ["PERSON", ...]}`（可选，v2 版本）

**认证:**
- UI 路由（`GET /`、`GET /static/*`）: 无需认证（公开访问）
- API 路由（`POST /api/v1/*`）: 需要 Basic Auth（现有行为）
- 前端使用凭据提示处理 401 响应

### 文件结构
```
pkg/web/
├── server.go           # 添加 GET / 路由和静态文件处理器
├── static/
│   ├── index.html      # 单页应用入口点
│   ├── styles.css      # UI 样式（响应式、简洁设计）
│   └── app.js          # 视图切换、API 调用、状态管理
└── handlers/
    └── config.go       # （可选）GET /api/v1/config 处理器
```

### 用户流程

**匿名化流程:**
1. 用户在浏览器中打开 `http://localhost:8080/`
2. UI 加载视图 1（匿名化模式）
3. 用户从下拉菜单选择实体类型或添加自定义类型
4. 用户在左侧面板输入文本
5. 用户点击"匿名化" → 出现加载动画
6. JavaScript 使用凭据调用 `POST /api/v1/anonymize`
7. 成功后：右侧面板显示匿名化文本，"切换到还原模式"按钮出现
8. JavaScript 将实体和匿名化文本存储到 sessionStorage

**还原流程:**
1. 用户点击"切换到还原模式"
2. UI 切换到视图 2（还原模式）
3. 顶部显示实体映射，左侧面板显示匿名化文本（只读）
4. 用户将外部处理的文本粘贴到右侧面板
5. 用户点击"还原" → 出现加载动画
6. JavaScript 使用 sessionStorage 中的实体调用 `POST /api/v1/restore`
7. 成功后：右侧面板更新为还原后的文本
8. 用户可以多次点击"还原"（实体保留）
9. 用户点击"返回匿名化"返回视图 1

**Trade-offs:**
- Data lost on tab close (acceptable for anonymization workflow)
- Cannot share across tabs (not needed for single-user tool)
- Browser storage limits (~5MB, sufficient for entity mappings)

### UI Layout

#### View 1: Anonymize Mode
```
┌────────────────────────────────────────────────┐
│ Entity Types: [PERSON▼] [+ Add Custom Type]   │
├────────────────────┬───────────────────────────┤
│                    │                           │
│ Input Text         │ Anonymized Output         │
│ (editable)         │ (read-only, shows result) │
│                    │                           │
│                    │                           │
├────────────────────┴───────────────────────────┤
│ [Anonymize] [Switch to Restore Mode]          │
└────────────────────────────────────────────────┘
```

#### View 2: Restore Mode
```
┌────────────────────────────────────────────────┐
│ Entity Mappings:                               │
│ [PERSON_1] → 张三  [ORG_1] → 阿里巴巴  ...      │
├────────────────────┬───────────────────────────┤
│ Anonymized Text    │ Input for Restoration     │
│ (read-only)        │ (editable, paste here)    │
│                    │                           │
│                    │                           │
├────────────────────┴───────────────────────────┤
│ [Restore] [Back to Anonymize]                  │
└────────────────────────────────────────────────┘
```

### API Integration

**Existing Endpoints (no changes):**
- `POST /api/v1/anonymize`: Request body `{"text": "...", "entity_types": ["PERSON", ...]}`
- `POST /api/v1/restore`: Request body `{"text": "...", "entities": {"[PERSON_1]": "张三", ...}}`

**New Endpoints:**
- `GET /`: Serve `index.html` (no auth required)
- `GET /static/*`: Serve CSS/JS files (no auth required)
- `GET /api/v1/config`: Return `{"entity_types": ["PERSON", ...]}` from CLI flag (optional, for v2)

**Authentication:**
- UI routes (`GET /`, `GET /static/*`): No authentication (public access)
- API routes (`POST /api/v1/*`): Basic Auth required (existing behavior)
- Frontend handles 401 responses with credential prompt

### File Structure
```
pkg/web/
├── server.go           # Add GET / route and static file handler
├── static/
│   ├── index.html      # Single-page application entry point
│   ├── styles.css      # UI styling (responsive, clean design)
│   └── app.js          # View switching, API calls, state management
└── handlers/
    └── config.go       # (Optional) GET /api/v1/config handler
```

### User Flow

**Anonymize Flow:**
1. User opens `http://localhost:8080/` in browser
2. UI loads with View 1 (Anonymize Mode)
3. User selects entity types from dropdown or adds custom types
4. User inputs text in left panel
5. User clicks "Anonymize" → Loading spinner appears
6. JavaScript calls `POST /api/v1/anonymize` with credentials
7. On success: Right panel shows anonymized text, "Switch to Restore Mode" button appears
8. JavaScript stores entities and anonymized text in sessionStorage

**Restore Flow:**
1. User clicks "Switch to Restore Mode"
2. UI transitions to View 2 (Restore Mode)
3. Entity mappings displayed at top, anonymized text in left panel (read-only)
4. User pastes externally processed text into right panel
5. User clicks "Restore" → Loading spinner appears
6. JavaScript calls `POST /api/v1/restore` with entities from sessionStorage
7. On success: Right panel updates with restored text
8. User can click "Restore" multiple times (entities preserved)
9. 用户点击"返回匿名化"返回视图 1

### 错误处理

**客户端:**
- 401 未授权 → 提示输入用户名/密码，使用 Basic Auth 重试
- 400 错误请求 → 显示 API 响应的错误消息
- 500 服务器错误 → 显示通用错误，建议检查日志
- 网络故障 → 显示"无法连接到服务器"消息

**服务器端:**
- 无效的实体类型 → 返回 400 及错误详情
- 空文本 → 返回 400 "文本不能为空"
- LLM 错误 → 返回 500 及清理后的错误消息（无敏感数据）

### 响应式设计

**桌面端（>768px）:**
- 左右分栏布局（左：输入，右：输出）
- 实体选择器在顶部（横向）

**移动端（<768px）:**
- 堆叠布局（输入在上，输出在下）
- 实体选择器为下拉菜单（纵向）
- 可滚动的实体映射

### 性能考虑

- 懒加载：仅在需要时加载视图组件
- 防抖：输入后等待 500ms 再启用"匿名化"按钮（防止误点击）
- 缓存：浏览器使用版本控制缓存静态资源（CSS/JS）
- 压缩：生产环境压缩 CSS/JS（未来优化）

### 安全考虑

**认证:**
- Basic Auth 凭据存储在内存中（不存 localStorage）以防止 XSS 攻击
- 生产环境建议使用 HTTPS（防止凭据嗅探）

**XSS 防护:**
- 使用 `textContent` 而不是 `innerHTML` 显示用户输入
- 显示前清理实体映射（转义 HTML 实体）

**CSRF:**
- 不适用（无状态 API，无 cookies/sessions）
- Authorization 头中的 Basic Auth

### 测试策略

**手动测试:**
- 浏览器兼容性：Chrome、Firefox、Safari、Edge
- 响应式设计：桌面、平板、移动视口
- 用户流程：匿名化 → 还原 → 返回，多次还原周期
- 错误场景：401、400、500、网络故障

**自动化测试（未来）:**
- JavaScript 函数单元测试（状态管理、API 调用）
- Playwright E2E 测试（完整用户流程）

### 未来增强功能（超出范围）

1. **多会话**: 允许保存/加载实体映射为文件
2. **实体编辑**: 在 UI 中手动编辑实体映射
3. **批量处理**: 上传多个文件进行匿名化
4. **暗黑模式**: 切换 UI 主题
5. **国际化**: 多语言支持（英文、中文）
6. **WebSocket**: 实时匿名化流（用于大文本）
