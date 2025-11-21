# web-api Specification

## Purpose
定义 Inu 的 Web API 功能规范，提供基于 HTTP 的文本匿名化和还原服务。

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

#### Scenario: 环境变量未配置时报错
- **WHEN** 用户执行 `inu web --admin-token secret123` 但未设置 OPENAI_API_KEY
- **THEN** 系统应该显示友好的错误信息，说明需要配置的环境变量

### Requirement: 匿名化 API 端点
系统 SHALL 提供 `POST /api/v1/anonymize` 端点来匿名化文本。

#### Scenario: 成功匿名化单个实体
- **GIVEN** 客户端已通过 HTTP Basic Auth 认证
- **WHEN** 客户端发送 POST 请求到 `/api/v1/anonymize`：
  ```json
  {
    "text": "张三的电话是 13800138000"
  }
  ```
- **THEN** 系统应该返回 200 OK 和匿名化结果：
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
- **THEN** 系统应该只识别和匿名化 "个人信息" 类型的实体
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
- **WHEN** 匿名化过程中 LLM API 调用失败
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
系统 SHALL 提供 `POST /api/v1/restore` 端点来还原匿名化文本。

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

#### Scenario: 空匿名化文本
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
