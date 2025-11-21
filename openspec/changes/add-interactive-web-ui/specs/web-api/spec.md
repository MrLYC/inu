# Web API Spec Deltas

## ADDED Requirements

### Requirement: 交互式 Web 界面
Web 服务器 SHALL 提供单页应用 (SPA) 界面,允许用户通过浏览器进行匿名化和还原操作。UI 路由需要 Basic Auth 认证（与 API 使用相同的认证机制）。

#### 场景: 加载主页面
```
请求: GET /
请求头: Authorization: Basic YWRtaW46c2VjcmV0
响应: 200 OK
Content-Type: text/html
返回 index.html 文件
```

#### 场景: 加载主页面未授权
```
请求: GET /
请求头: 无 Authorization 或凭据错误
响应: 401 Unauthorized
WWW-Authenticate: Basic realm="Inu Web UI"
```

#### 场景: 加载静态资源
```
请求: GET /static/app.js
请求头: Authorization: Basic YWRtaW46c2VjcmV0
响应: 200 OK
Content-Type: application/javascript
返回 JavaScript 文件

请求: GET /static/styles.css
请求头: Authorization: Basic YWRtaW46c2VjcmV0
响应: 200 OK
Content-Type: text/css
返回 CSS 文件
```

#### 场景: 静态文件不存在
```
请求: GET /static/nonexistent.js
请求头: Authorization: Basic YWRtaW46c2VjcmV0
响应: 404 Not Found
```

### Requirement: 配置端点
Web 服务器 SHALL 提供配置端点,返回服务器启动时的实体类型配置。

#### 场景: 获取实体类型配置
```
请求: GET /api/v1/config
请求头: Authorization: Basic YWRtaW46c2VjcmV0
响应: 200 OK
Content-Type: application/json

{
  "entity_types": ["PERSON", "ORG", "EMAIL", "PHONE", "ADDRESS"]
}

说明: entity_types 从 --entity-types CLI 参数获取
```

#### 场景: 获取配置未授权
```
请求: GET /api/v1/config
请求头: 无 Authorization 或凭据错误
响应: 401 Unauthorized
```

### Requirement: 匿名化视图
UI SHALL 提供匿名化视图,包含实体类型选择器、输入文本区域、输出显示区域和操作按钮。

#### 场景: 选择预定义实体类型
```
UI状态:
- 页面加载时调用 GET /api/v1/config 获取实体类型
- 实体类型下拉菜单显示从配置端点获取的类型
- 默认选中第一个类型

用户操作:
- 从下拉菜单选择 PERSON 类型
- 可多选 (PERSON, ORG, EMAIL)

预期:
- 选中的类型用于后续匿名化请求的 entity_types 字段
```

#### 场景: 添加自定义实体类型
```
UI状态:
- 实体类型选择器包含 "+ 添加自定义类型" 按钮

用户操作:
- 点击 "+ 添加自定义类型" 按钮
- 输入自定义类型名称 "ID_CARD"
- 确认添加

预期:
- "ID_CARD" 出现在类型选择列表中
- 可以选择该自定义类型用于匿名化
```

#### 场景: 输入并匿名化文本
```
UI状态:
- 左侧输入框为空
- 右侧输出区域显示 "等待匿名化结果..."
- "匿名化" 按钮可点击
- "切换到还原模式" 按钮不可见

用户操作:
- 在左侧输入框输入: "张三在阿里巴巴工作,邮箱是zhangsan@example.com"
- 选择实体类型: PERSON, ORG, EMAIL
- 点击 "匿名化" 按钮

预期:
- "匿名化" 按钮显示加载动画 (禁用状态)
- JavaScript 调用 POST /api/v1/anonymize
- 成功后右侧显示: "[PERSON_1]在[ORG_1]工作,邮箱是[EMAIL_1]"
- "匿名化" 按钮恢复可点击状态
- "切换到还原模式" 按钮出现
- 实体映射和匿名化文本存入 sessionStorage
```

#### 场景: 匿名化失败 - 401 未授权
```
UI状态:
- 用户输入文本后点击 "匿名化"

后端响应:
- 401 Unauthorized (未提供 Basic Auth 或凭据错误)

预期:
- UI 显示登录弹窗,提示输入用户名和密码
- 用户输入凭据后重试请求
- 凭据存储在内存中 (不存入 localStorage)

说明: 
- 由于 UI 路由也需要认证,用户在访问 GET / 时就会被要求登录
- 浏览器会记住凭据并自动在后续 API 请求中使用
- 此场景处理会话过期或凭据变更的情况
```

#### 场景: 匿名化失败 - 400 错误请求
```
UI状态:
- 用户输入空文本后点击 "匿名化"

后端响应:
- 400 Bad Request
- {"error": "Text cannot be empty"}

预期:
- 右侧输出区域显示详细错误信息: "错误: Text cannot be empty"
- "匿名化" 按钮恢复可点击状态
```

#### 场景: 匿名化失败 - 500 服务器错误
```
UI状态:
- 用户输入文本后点击 "匿名化"

后端响应:
- 500 Internal Server Error
- {"error": "LLM service unavailable"}

预期:
- 右侧输出区域显示详细错误信息: "服务器错误: LLM service unavailable"
- "匿名化" 按钮恢复可点击状态
- 建议用户检查服务器日志

说明: 保留详细错误消息以便调试,因为工具主要用于内部/开发环境
```

### Requirement: 还原视图
UI SHALL 提供还原视图,显示实体映射,包含只读的匿名化文本、可编辑的输入区域和还原按钮。

#### 场景: 切换到还原模式
```
UI状态:
- 用户在匿名化视图完成匿名化
- "切换到还原模式" 按钮可见

用户操作:
- 点击 "切换到还原模式" 按钮

预期:
- UI 切换到还原视图
- 顶部显示实体映射: "[PERSON_1] → 张三  [ORG_1] → 阿里巴巴  [EMAIL_1] → zhangsan@example.com"
- 左侧只读区域显示匿名化文本: "[PERSON_1]在[ORG_1]工作,邮箱是[EMAIL_1]"
- 右侧可编辑区域为空
- "还原" 和 "返回匿名化" 按钮可见
```

#### 场景: 输入并还原文本
```
UI状态:
- 用户在还原视图

用户操作:
- 在右侧输入框粘贴外部处理后的文本: "[PERSON_1]的新邮箱是[EMAIL_1]"
- 点击 "还原" 按钮

预期:
- "还原" 按钮显示加载动画 (禁用状态)
- JavaScript 调用 POST /api/v1/restore,使用 sessionStorage 中的实体映射
- 成功后右侧更新为: "张三的新邮箱是zhangsan@example.com"
- "还原" 按钮恢复可点击状态
```

#### 场景: 多次还原
```
UI状态:
- 用户完成一次还原操作

用户操作:
- 在右侧输入框清空并粘贴新文本: "联系[PERSON_1]和[ORG_1]"
- 点击 "还原" 按钮

预期:
- 再次调用 POST /api/v1/restore
- 右侧更新为: "联系张三和阿里巴巴"
- 实体映射保持不变 (从首次匿名化获取)
- 可以无限次还原不同文本
```

#### 场景: 返回匿名化视图
```
UI状态:
- 用户在还原视图

用户操作:
- 点击 "返回匿名化" 按钮

预期:
- UI 切换回匿名化视图
- 左侧输入框保留之前的原始文本 (提供更好的用户体验)
- 右侧输出区域保留匿名化结果
- 实体映射仍在 sessionStorage 中 (未清除)
- "切换到还原模式" 按钮仍可见
```

#### 场景: 还原失败 - 400 错误请求
```
UI状态:
- 用户在还原视图输入空文本后点击 "还原"

后端响应:
- 400 Bad Request
- {"error": "Text cannot be empty"}

预期:
- 右侧输入区域显示错误信息: "错误: Text cannot be empty"
- "还原" 按钮恢复可点击状态
```

### Requirement: 前端状态管理
前端 SHALL 使用 sessionStorage 存储实体映射、匿名化文本和原始文本,实现视图间状态共享。

#### 场景: 存储匿名化结果
```
时机: POST /api/v1/anonymize 成功返回

数据结构:
{
  "entities": {
    "[PERSON_1]": "张三",
    "[ORG_1]": "阿里巴巴",
    "[EMAIL_1]": "zhangsan@example.com"
  },
  "anonymizedText": "[PERSON_1]在[ORG_1]工作,邮箱是[EMAIL_1]",
  "originalText": "张三在阿里巴巴工作,邮箱是zhangsan@example.com",
  "entityTypes": ["PERSON", "ORG", "EMAIL"]
}

存储位置: sessionStorage.setItem("inuState", JSON.stringify(data))

预期:
- 数据在浏览器标签页关闭前保持
- 刷新页面后数据仍存在
- 关闭标签页后数据自动清除
```

#### 场景: 恢复状态
```
时机: 页面刷新或用户返回

用户操作:
- 用户完成匿名化后刷新页面

预期:
- JavaScript 从 sessionStorage.getItem("inuState") 读取数据
- 如果数据存在:
  - 左侧输入框显示原始文本
  - 右侧显示匿名化结果
  - "切换到还原模式" 按钮可见
  - 实体类型选择器恢复之前选择
- 如果数据不存在:
  - UI 恢复初始状态 (空输入框, 等待匿名化)
```

#### 场景: 清除状态
```
时机: 用户开始新的匿名化操作

用户操作:
- 用户在匿名化视图点击 "匿名化" 按钮 (新文本)

预期:
- sessionStorage.removeItem("inuState") 清除旧数据
- 新的匿名化结果覆盖 sessionStorage
- 之前的实体映射被丢弃
```

### Requirement: 响应式设计
UI SHALL 适配不同屏幕尺寸,提供桌面和移动端优化布局。

#### 场景: 桌面端布局 (>768px)
```
设备: 1920x1080 桌面浏览器

预期:
- 匿名化视图: 左右分栏 (输入在左, 输出在右)
- 还原视图: 左右分栏 (匿名化文本在左, 输入在右)
- 实体类型选择器: 横向排列
- 实体映射显示: 横向滚动条 (超出时)
```

#### 场景: 移动端布局 (<768px)
```
设备: 375x667 iPhone SE

预期:
- 匿名化视图: 上下堆叠 (输入在上, 输出在下)
- 还原视图: 上下堆叠 (匿名化文本在上, 输入在下)
- 实体类型选择器: 下拉菜单 (节省空间)
- 实体映射显示: 纵向滚动
- 按钮: 全宽显示 (易于点击)
```

### Requirement: 错误处理和用户反馈
UI SHALL 提供清晰的详细错误提示和加载状态,改善用户体验。保留详细错误消息以便内部调试。

#### 场景: 网络连接失败
```
UI状态: 用户点击 "匿名化" 按钮

后端状态: 服务器不可达 (网络断开或服务器停止)

预期:
- JavaScript fetch() 抛出 TypeError 或 Network Error
- 右侧输出区域显示详细错误: "无法连接到服务器,请检查网络或服务器状态"
- "匿名化" 按钮恢复可点击状态
```

#### 场景: 加载动画
```
时机: 调用 API 期间 (anonymize 或 restore)

预期:
- 按钮文本替换为加载图标 (旋转动画)
- 按钮禁用 (防止重复点击)
- 输出区域显示 "处理中..." 提示
- API 返回后恢复正常状态
```

#### 场景: 空输入提示
```
UI状态: 用户未输入任何文本

用户操作: 点击 "匿名化" 或 "还原" 按钮

预期:
- 不调用 API (前端验证)
- 输入框边框高亮红色
- 显示提示信息: "请输入要处理的文本"
```
