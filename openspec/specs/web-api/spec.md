# web-api Specification

## Purpose
定义 Inu 的 Web API 功能规范，提供基于 HTTP 的文本脱敏和还原服务。

## ADDED Requirements

### Requirement: Web 服务器命令
系统 SHALL 提供 `web` 子命令来启动 HTTP API 服务器。

#### Scenario: 使用默认配置启动服务器
- **GIVEN** 用户已配置必需的环境变量（OPENAI_API_KEY, OPENAI_MODEL_NAME）
- **WHEN** 用户执行 `inu web --admin-token secret123`
- **THEN** 系统应该在 `127.0.0.1:8080` 启动 Web 服务器
- **AND** 使用默认管理员用户名 `admin`
- **AND** 使用提供的 token `secret123` 进行认证

#### Scenario: 自定义监听地址和端口
- **WHEN** 用户执行 `inu web --addr 0.0.0.0:9090 --admin-token secret123`
- **THEN** 系统应该在 `0.0.0.0:9090` 启动服务器
- **AND** 服务器可从外部网络访问

#### Scenario: 自定义管理员用户名
- **WHEN** 用户执行 `inu web --admin-user root --admin-token secret123`
- **THEN** 系统应该使用 `root` 作为管理员用户名进行认证

#### Scenario: 缺少 admin-token 时报错
- **WHEN** 用户执行 `inu web` 而不提供 `--admin-token`
- **THEN** 系统应该退出并显示错误：需要指定 --admin-token

#### Scenario: 优雅关闭服务器
- **GIVEN** Web 服务器正在运行
- **WHEN** 用户发送 SIGINT (Ctrl+C) 或 SIGTERM 信号
- **THEN** 系统应该等待正在处理的请求完成（最多 5 秒）
- **AND** 关闭服务器并释放资源
- **AND** 输出关闭完成的日志信息

#### Scenario: 显示启动信息
- **WHEN** Web 服务器成功启动
- **THEN** 系统应该输出：
  - 监听地址和端口
  - API 版本
  - 管理员用户名（不显示 token）
  - 可用的 API 端点列表
  - Web UI 已内置的说明（无需配置静态文件目录）

#### Scenario: 环境变量未配置时报错
- **WHEN** 用户执行 `inu web --admin-token secret123` 但未设置 OPENAI_API_KEY
- **THEN** 系统应该显示友好的错误信息，说明需要配置的环境变量

### Requirement: 脱敏 API 端点
系统 SHALL 提供 `POST /api/v1/anonymize` 端点来脱敏文本。

#### Scenario: 成功脱敏单个实体
- **GIVEN** 客户端已通过 HTTP Basic Auth 认证
- **WHEN** 客户端发送 POST 请求到 `/api/v1/anonymize`：
  ```json
  {
    "text": "张三的电话是 13800138000"
  }
  ```
- **THEN** 系统应该返回 200 OK 和脱敏结果：
  ```json
  {
    "anonymized_text": "<个人信息[0].姓名.全名>的电话是<个人信息[1].电话.号码>",
    "entities": [
      {
        "key": "<个人信息[0].姓名.全名>",
        "type": "个人信息",
        "id": "0",
        "category": "姓名",
        "detail": "张三",
        "values": ["张三"]
      },
      {
        "key": "<个人信息[1].电话.号码>",
        "type": "个人信息",
        "id": "1",
        "category": "电话",
        "detail": "13800138000",
        "values": ["13800138000"]
      }
    ]
  }
  ```

#### Scenario: 指定实体类型
- **WHEN** 客户端发送请求并指定 `entity_types`：
  ```json
  {
    "text": "张三在 ABC 公司工作",
    "entity_types": ["个人信息"]
  }
  ```
- **THEN** 系统应该只识别和脱敏 "个人信息" 类型的实体
- **AND** 忽略 "组织机构" 类型（ABC 公司）

#### Scenario: 使用默认实体类型
- **WHEN** 客户端发送请求不指定 `entity_types`
- **THEN** 系统应该使用默认实体类型列表：["个人信息", "业务信息", "资产信息", "账户信息", "位置数据", "文档名称", "组织机构", "岗位称谓"]

#### Scenario: 空文本输入
- **WHEN** 客户端发送空文本：
  ```json
  {
    "text": ""
  }
  ```
- **THEN** 系统应该返回 400 Bad Request：
  ```json
  {
    "error": "invalid_input",
    "message": "Text cannot be empty",
    "code": 400
  }
  ```

#### Scenario: 缺少认证信息
- **WHEN** 客户端发送请求不带 Authorization 头
- **THEN** 系统应该返回 401 Unauthorized：
  ```json
  {
    "error": "unauthorized",
    "message": "Authentication required",
    "code": 401
  }
  ```

#### Scenario: 认证信息错误
- **WHEN** 客户端发送错误的用户名或密码
- **THEN** 系统应该返回 401 Unauthorized

#### Scenario: 无效的 JSON 格式
- **WHEN** 客户端发送格式错误的 JSON
- **THEN** 系统应该返回 400 Bad Request 和 JSON 解析错误信息

#### Scenario: LLM API 调用失败
- **WHEN** 脱敏过程中 LLM API 调用失败
- **THEN** 系统应该返回 500 Internal Server Error：
  ```json
  {
    "error": "llm_error",
    "message": "Failed to call LLM API: connection timeout",
    "code": 500
  }
  ```

#### Scenario: 无实体匹配
- **WHEN** LLM 在文本中未找到任何指定类型的实体
- **THEN** 系统应该返回 200 OK：
  ```json
  {
    "anonymized_text": "This is plain text",
    "entities": []
  }
  ```

### Requirement: 还原 API 端点
系统 SHALL 提供 `POST /api/v1/restore` 端点来还原脱敏文本。

#### Scenario: 成功还原文本
- **GIVEN** 客户端已通过认证
- **WHEN** 客户端发送 POST 请求到 `/api/v1/restore`：
  ```json
  {
    "anonymized_text": "<个人信息[0].姓名.全名>的电话是<个人信息[1].电话.号码>",
    "entities": [
      {
        "key": "<个人信息[0].姓名.全名>",
        "values": ["张三"]
      },
      {
        "key": "<个人信息[1].电话.号码>",
        "values": ["13800138000"]
      }
    ]
  }
  ```
- **THEN** 系统应该返回 200 OK：
  ```json
  {
    "restored_text": "张三的电话是 13800138000"
  }
  ```

#### Scenario: 空实体列表
- **WHEN** 客户端发送空的 `entities` 数组：
  ```json
  {
    "anonymized_text": "Some text",
    "entities": []
  }
  ```
- **THEN** 系统应该返回 200 OK 并返回原文本（无需还原）：
  ```json
  {
    "restored_text": "Some text"
  }
  ```

#### Scenario: 空脱敏文本
- **WHEN** 客户端发送空的 `anonymized_text`
- **THEN** 系统应该返回 400 Bad Request：
  ```json
  {
    "error": "invalid_input",
    "message": "Anonymized text cannot be empty",
    "code": 400
  }
  ```

#### Scenario: 缺少 entities 字段
- **WHEN** 客户端请求中不包含 `entities` 字段
- **THEN** 系统应该返回 400 Bad Request

#### Scenario: 实体无 values
- **WHEN** 某个实体的 `values` 数组为空：
  ```json
  {
    "entities": [
      {
        "key": "<个人信息[0].姓名.全名>",
        "values": []
      }
    ]
  }
  ```
- **THEN** 系统应该跳过该实体，不进行替换

#### Scenario: 文本中无匹配的占位符
- **WHEN** `anonymized_text` 中不包含任何实体 key
- **THEN** 系统应该返回 200 OK 和原文本（无需还原）

### Requirement: 健康检查端点
系统 SHALL 提供 `GET /health` 端点用于健康检查。

#### Scenario: 健康检查成功
- **WHEN** 客户端发送 GET 请求到 `/health`（无需认证）
- **THEN** 系统应该返回 200 OK：
  ```json
  {
    "status": "ok",
    "version": "v0.1.0"
  }
  ```

#### Scenario: 健康检查不需要认证
- **WHEN** 客户端不带认证信息访问 `/health`
- **THEN** 系统应该仍然返回 200 OK

### Requirement: 请求日志和错误处理
系统 SHALL 记录所有 API 请求和错误信息。

#### Scenario: 记录成功请求
- **WHEN** API 请求成功处理
- **THEN** 系统应该在日志中记录：
  - 请求方法（POST/GET）
  - 请求路径
  - HTTP 状态码（200）
  - 处理耗时
  - 客户端 IP（可选）

#### Scenario: 记录失败请求
- **WHEN** API 请求失败（4xx 或 5xx）
- **THEN** 系统应该在日志中记录：
  - 错误类型
  - 错误详细信息
  - 请求上下文

#### Scenario: 记录认证失败
- **WHEN** 客户端认证失败
- **THEN** 系统应该记录认证失败事件：
  - 尝试的用户名
  - 客户端 IP
  - 时间戳

#### Scenario: 敏感信息不记录
- **WHEN** 记录请求日志
- **THEN** 系统不应记录：
  - Admin token
  - 完整的请求 body（可能包含敏感文本）
  - LLM API key
- **AND** 仅记录元数据（路径、状态码、耗时）
## Requirements
### Requirement: 交互式 Web 界面
Web 服务器 SHALL 提供单页应用 (SPA) 界面,允许用户通过浏览器进行脱敏和还原操作。UI 路由需要 Basic Auth 认证（与 API 使用相同的认证机制）。

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

### Requirement: 脱敏视图
UI SHALL 提供脱敏视图,包含实体类型选择器、输入文本区域、输出显示区域和操作按钮。

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
- 选中的类型用于后续脱敏请求的 entity_types 字段
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
- 可以选择该自定义类型用于脱敏
```

#### 场景: 输入并脱敏文本
```
UI状态:
- 左侧输入框为空
- 右侧输出区域显示 "等待脱敏结果..."
- "脱敏" 按钮可点击
- "切换到还原模式" 按钮不可见

用户操作:
- 在左侧输入框输入: "张三在阿里巴巴工作,邮箱是zhangsan@example.com"
- 选择实体类型: PERSON, ORG, EMAIL
- 点击 "脱敏" 按钮

预期:
- "脱敏" 按钮显示加载动画 (禁用状态)
- JavaScript 调用 POST /api/v1/anonymize
- 成功后右侧显示: "[PERSON_1]在[ORG_1]工作,邮箱是[EMAIL_1]"
- "脱敏" 按钮恢复可点击状态
- "切换到还原模式" 按钮出现
- 实体映射和脱敏文本存入 sessionStorage
```

#### 场景: 脱敏失败 - 401 未授权
```
UI状态:
- 用户输入文本后点击 "脱敏"

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

#### 场景: 脱敏失败 - 400 错误请求
```
UI状态:
- 用户输入空文本后点击 "脱敏"

后端响应:
- 400 Bad Request
- {"error": "Text cannot be empty"}

预期:
- 右侧输出区域显示详细错误信息: "错误: Text cannot be empty"
- "脱敏" 按钮恢复可点击状态
```

#### 场景: 脱敏失败 - 500 服务器错误
```
UI状态:
- 用户输入文本后点击 "脱敏"

后端响应:
- 500 Internal Server Error
- {"error": "LLM service unavailable"}

预期:
- 右侧输出区域显示详细错误信息: "服务器错误: LLM service unavailable"
- "脱敏" 按钮恢复可点击状态
- 建议用户检查服务器日志

说明: 保留详细错误消息以便调试,因为工具主要用于内部/开发环境
```

### Requirement: 静态资源嵌入式服务
系统 SHALL 使用嵌入的文件系统提供 Web UI 静态资源，无需依赖外部文件。

#### Scenario: 访问首页从嵌入资源加载
- **GIVEN** Web 服务器正在运行
- **AND** `pkg/web/static/` 目录中的文件已在编译时嵌入
- **WHEN** 客户端发送 `GET /` 请求
- **THEN** 服务器应该从嵌入的文件系统读取 `index.html`
- **AND** 返回 200 OK 和 HTML 内容
- **AND** Content-Type 应该是 `text/html; charset=utf-8`
- **AND** 无需访问外部文件系统

#### Scenario: 访问静态资源文件
- **WHEN** 客户端请求 `GET /static/app.js`
- **THEN** 服务器应该从嵌入的文件系统返回 `app.js` 内容
- **AND** Content-Type 应该是 `application/javascript`
- **WHEN** 客户端请求 `GET /static/styles.css`
- **THEN** 服务器应该从嵌入的文件系统返回 `styles.css` 内容
- **AND** Content-Type 应该是 `text/css`

#### Scenario: 静态资源不存在返回 404
- **WHEN** 客户端请求不存在的静态资源 `GET /static/nonexistent.js`
- **THEN** 服务器应该返回 404 Not Found
- **AND** 响应体应该包含友好的错误信息

#### Scenario: 嵌入资源无需外部目录
- **GIVEN** 二进制文件部署在目录 `/opt/inu/`
- **AND** 该目录中不存在 `pkg/web/static/` 文件夹
- **WHEN** 启动服务器 `./inu web --admin-token test`
- **THEN** 服务器应该正常启动
- **AND** Web UI 应该完全可访问
- **AND** 所有静态资源应该正常加载

#### Scenario: 静态资源缓存控制
- **WHEN** 客户端请求静态资源
- **THEN** 响应应该包含适当的缓存头
- **AND** 支持 ETag 或 Last-Modified 进行缓存验证
- **AND** 客户端可以通过 If-None-Match 或 If-Modified-Since 实现条件请求

### Requirement: 嵌入资源版本绑定
系统 SHALL 确保静态资源版本与二进制文件版本完全一致。

#### Scenario: 静态资源自动更新
- **GIVEN** 开发者修改了 `pkg/web/static/app.js`
- **WHEN** 重新编译二进制文件
- **THEN** 新二进制文件应该包含更新后的 `app.js`
- **AND** 旧二进制文件仍使用旧版本的 `app.js`
- **AND** 不存在版本不一致的风险

#### Scenario: 二进制文件独立性
- **GIVEN** 用户从 GitHub Release 下载二进制文件
- **WHEN** 用户在任意目录运行二进制文件
- **THEN** Web UI 应该完全可用
- **AND** 无需下载或配置额外的静态文件
- **AND** 无需关心静态文件的路径或位置

#### Scenario: 多版本共存
- **GIVEN** 用户在不同目录运行不同版本的二进制文件
  - v1.0.0 在 `/opt/inu/v1/`
  - v1.1.0 在 `/opt/inu/v2/`
- **WHEN** 分别启动两个服务器（不同端口）
- **THEN** 每个服务器应该使用其编译时嵌入的静态资源版本
- **AND** 两个服务器的 UI 不会互相干扰

### Requirement: 还原视图
UI SHALL 提供还原视图,显示实体映射,包含只读的脱敏文本、可编辑的输入区域和还原按钮。

#### 场景: 切换到还原模式
```
UI状态:
- 用户在脱敏视图完成脱敏
- "切换到还原模式" 按钮可见

用户操作:
- 点击 "切换到还原模式" 按钮

预期:
- UI 切换到还原视图
- 顶部显示实体映射: "[PERSON_1] → 张三  [ORG_1] → 阿里巴巴  [EMAIL_1] → zhangsan@example.com"
- 左侧只读区域显示脱敏文本: "[PERSON_1]在[ORG_1]工作,邮箱是[EMAIL_1]"
- 右侧可编辑区域为空
- "还原" 和 "返回脱敏" 按钮可见
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
- 实体映射保持不变 (从首次脱敏获取)
- 可以无限次还原不同文本
```

#### 场景: 返回脱敏视图
```
UI状态:
- 用户在还原视图

用户操作:
- 点击 "返回脱敏" 按钮

预期:
- UI 切换回脱敏视图
- 左侧输入框保留之前的原始文本 (提供更好的用户体验)
- 右侧输出区域保留脱敏结果
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
前端 SHALL 使用 sessionStorage 存储实体映射、脱敏文本和原始文本,实现视图间状态共享。

#### 场景: 存储脱敏结果
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
- 用户完成脱敏后刷新页面

预期:
- JavaScript 从 sessionStorage.getItem("inuState") 读取数据
- 如果数据存在:
  - 左侧输入框显示原始文本
  - 右侧显示脱敏结果
  - "切换到还原模式" 按钮可见
  - 实体类型选择器恢复之前选择
- 如果数据不存在:
  - UI 恢复初始状态 (空输入框, 等待脱敏)
```

#### 场景: 清除状态
```
时机: 用户开始新的脱敏操作

用户操作:
- 用户在脱敏视图点击 "脱敏" 按钮 (新文本)

预期:
- sessionStorage.removeItem("inuState") 清除旧数据
- 新的脱敏结果覆盖 sessionStorage
- 之前的实体映射被丢弃
```

### Requirement: 响应式设计
UI SHALL 适配不同屏幕尺寸,提供桌面和移动端优化布局。

#### 场景: 桌面端布局 (>768px)
```
设备: 1920x1080 桌面浏览器

预期:
- 脱敏视图: 左右分栏 (输入在左, 输出在右)
- 还原视图: 左右分栏 (脱敏文本在左, 输入在右)
- 实体类型选择器: 横向排列
- 实体映射显示: 横向滚动条 (超出时)
```

#### 场景: 移动端布局 (<768px)
```
设备: 375x667 iPhone SE

预期:
- 脱敏视图: 上下堆叠 (输入在上, 输出在下)
- 还原视图: 上下堆叠 (脱敏文本在上, 输入在下)
- 实体类型选择器: 下拉菜单 (节省空间)
- 实体映射显示: 纵向滚动
- 按钮: 全宽显示 (易于点击)
```

### Requirement: 错误处理和用户反馈
UI SHALL 提供清晰的详细错误提示和加载状态,改善用户体验。保留详细错误消息以便内部调试。

#### 场景: 网络连接失败
```
UI状态: 用户点击 "脱敏" 按钮

后端状态: 服务器不可达 (网络断开或服务器停止)

预期:
- JavaScript fetch() 抛出 TypeError 或 Network Error
- 右侧输出区域显示详细错误: "无法连接到服务器,请检查网络或服务器状态"
- "脱敏" 按钮恢复可点击状态
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

用户操作: 点击 "脱敏" 或 "还原" 按钮

预期:
- 不调用 API (前端验证)
- 输入框边框高亮红色
- 显示提示信息: "请输入要处理的文本"
```

