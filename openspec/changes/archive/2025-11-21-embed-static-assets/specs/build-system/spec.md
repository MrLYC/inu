# build-system 规格变更

## ADDED Requirements

### Requirement: 静态资源编译时嵌入
系统 SHALL 在编译时将 Web UI 静态资源自动嵌入到二进制文件中，无需额外的构建步骤或工具。

#### Scenario: 编译自动包含静态资源
- **GIVEN** 存在 `pkg/web/static/` 目录包含 HTML/CSS/JS 文件
- **WHEN** 执行 `make build` 或 `go build ./cmd/inu`
- **THEN** 生成的二进制文件应该包含所有静态资源
- **AND** 二进制文件可以独立运行，无需外部静态文件目录

#### Scenario: 静态文件修改后重新编译
- **GIVEN** 修改了 `pkg/web/static/` 中的任何文件
- **WHEN** 重新执行 `make build`
- **THEN** 新生成的二进制文件应该包含更新后的静态资源
- **AND** 运行新二进制时应该加载新版本的静态文件

#### Scenario: 多平台编译包含静态资源
- **WHEN** 执行 `make build-all` 进行交叉编译
- **THEN** 所有平台的二进制文件应该都包含静态资源
- **AND** 每个平台的二进制文件都可以独立运行 Web UI

#### Scenario: 编译失败当静态文件缺失
- **GIVEN** `pkg/web/static/` 目录不存在或缺少关键文件
- **WHEN** 执行 `make build`
- **THEN** 编译应该失败并提示静态文件缺失错误
- **AND** 错误信息应该明确指出缺失的文件或目录

### Requirement: 二进制文件大小管理
系统 SHALL 合理管理嵌入静态资源后的二进制文件大小。

#### Scenario: 嵌入静态资源后的文件大小
- **GIVEN** 静态资源总大小约为 10-20KB
- **WHEN** 编译包含嵌入资源的二进制文件
- **THEN** 二进制文件大小增加应该接近静态资源大小（±10%）
- **AND** 文件大小增加对分发和部署无显著影响

#### Scenario: 构建输出显示文件大小
- **WHEN** 执行 `make build` 或 `make build-all`
- **THEN** 构建完成后应该显示二进制文件的大小信息
- **AND** 开发者可以验证嵌入资源的影响

## MODIFIED Requirements

无修改的需求。

## REMOVED Requirements

无移除的需求。
