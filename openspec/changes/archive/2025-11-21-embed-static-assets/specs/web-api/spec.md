# web-api 规格变更

## ADDED Requirements

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

## MODIFIED Requirements

### Requirement: Web 服务器命令（修改）
系统 SHALL 提供 `web` 子命令来启动包含嵌入式 Web UI 的 HTTP API 服务器。

#### Scenario: 启动信息包含嵌入资源说明（新增）
- **WHEN** Web 服务器成功启动
- **THEN** 启动日志应该说明 Web UI 已内置
- **AND** 无需提示用户配置静态文件目录
- **AND** 显示 Web UI 访问地址（如 `http://127.0.0.1:8080/`）

## REMOVED Requirements

无移除的需求。
